package server

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sepuka/vkbotserver/config"
	"github.com/sepuka/vkbotserver/domain"
	"github.com/sepuka/vkbotserver/message"
	"github.com/sepuka/vkbotserver/middleware"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/fcgi"
	"net/http/httputil"
	"os"
	"os/signal"
	"syscall"
)

type SocketServer struct {
	cfg      config.Config
	logger   *zap.SugaredLogger
	messages message.HandlerMap
	handler  middleware.HandlerFunc
}

func NewSocketServer(
	cfg config.Config,
	logger *zap.SugaredLogger,
	messages message.HandlerMap,
	handler middleware.HandlerFunc,
) *SocketServer {
	return &SocketServer{
		cfg:      cfg,
		logger:   logger,
		messages: messages,
		handler:  handler,
	}
}

func (s *SocketServer) Listen() error {
	var (
		socket   = s.cfg.Socket
		signals  = make(chan os.Signal, 1)
		stop     = make(chan error, 1)
		listener net.Listener
		err      error
	)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	listener, err = net.Listen(`unix`, socket)
	if err != nil {
		s.logger.Errorf(`cannot listen to unix socket: %s`, err)
		return err
	}

	if err := os.Chmod(socket, 0775); err != nil {
		return err
	}

	defer func() error {
		return os.Remove(socket)
	}()

	go func() {
		<-signals
		if err = listener.Close(); err != nil {
			stop <- errors.Wrap(err, `unable to close HTTP connection`)
		}
	}()

	go s.server(listener, stop)

	err = <-stop

	return err
}

func (s *SocketServer) server(listener net.Listener, c chan<- error) {
	if err := fcgi.Serve(listener, s); err != nil {
		s.logger.Errorf(`cannot serve accept connections: %s`, err)
		c <- err
	}
}

func (s *SocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		callback = &domain.Request{}
		output   = []byte(`ok`)
		clone    []byte
		err      error
	)

	defer r.Body.Close()
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(`500 SocketServer error`))
		}
	}()

	if clone, err = httputil.DumpRequest(r, true); err != nil {
		s.
			logger.
			Errorf(`unable to dump request: %s`, err)
		panic(`invalid request`)
	}

	s.
		logger.
		With(
			zap.String(`host`, r.Host),
			zap.ByteString(`body`, clone),
		).
		Infof(`incoming %s-request to %s`, r.Method, r.URL.Path)

	if err = json.NewDecoder(r.Body).Decode(callback); err != nil {
		if _, err = w.Write([]byte(`invalid json`)); err != nil {
			s.logger.Errorf(`cannot write error message about invalid incoming json %s`, err)
		}
		w.WriteHeader(400)

		return
	}

	if finalHandler, ok := s.messages[callback.Type]; ok {
		if err = s.handler(finalHandler, callback, w); err != nil {
			s.logger.Errorf(`error while handling request: %s`, err)
		}
	} else {
		if _, err = w.Write(output); err != nil {
			s.logger.Errorf(`cannot write error message about unknown type field %s`, err)
		}
		w.WriteHeader(400)
	}
}

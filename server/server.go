package server

import (
	"github.com/mailru/easyjson"
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

const (
	invalidJSON = `invalid json`
)

type SocketServer struct {
	cfg      config.Config
	logger   *zap.SugaredLogger
	messages message.HandlerMap
	handler  middleware.HandlerFunc
}

// NewSocketServer constructor
func NewSocketServer(
	cfg config.Config,
	messages message.HandlerMap,
	handler middleware.HandlerFunc,
	logger *zap.SugaredLogger,
) *SocketServer {
	return &SocketServer{
		cfg:      cfg,
		logger:   logger,
		messages: messages,
		handler:  handler,
	}
}

// Listen listens unix socket which created by webserver
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
		_ = s.logger.Sync()
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
		clone    []byte
		err      error
	)

	defer r.Body.Close()
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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

	if err = easyjson.UnmarshalFromReader(r.Body, callback); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err = w.Write([]byte(invalidJSON)); err != nil {
			s.logger.Errorf(`cannot write error message about invalid incoming json %s`, err)
		}

		return
	}

	if finalHandler, ok := s.messages[callback.Type]; ok {
		if err = s.handler(finalHandler, callback, w); err != nil {
			s.logger.Errorf(`error while handling request: %s`, err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

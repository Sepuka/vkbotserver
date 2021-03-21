package api

import (
	"encoding/json"
	"fmt"
	"github.com/sepuka/vkbotserver/api/button"
	"github.com/sepuka/vkbotserver/config"
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/google/go-querystring/query"
)

const (
	defaultOutput         = `ok`
	Endpoint              = `https://api.vk.com/method`
	apiMethodMessagesSend = `messages.send`
	Version               = `5.120`
)

type (
	OutcomeMessage struct {
		Keyboard    string `url:"keyboard"`
		Message     string `url:"message"`
		AccessToken string `url:"access_token"`
		ApiVersion  string `url:"v"`
		PeerId      int    `url:"peer_id"`
		RandomId    int64  `url:"random_id"`
		Attachment  string `url:"attachment"`
	}

	HTTPClient interface {
		Do(req *http.Request) (*http.Response, error)
	}

	Api struct {
		logger *zap.SugaredLogger
		cfg    config.Config
		client HTTPClient
		rnd    Rnder
	}
)

// Just returns "ok"
func DefaultResponseBody() []byte {
	return []byte(defaultOutput)
}

// Creates API gate
func NewApi(logger *zap.SugaredLogger, cfg config.Config, client HTTPClient, rnd Rnder) *Api {
	return &Api{
		logger: logger,
		cfg:    cfg,
		client: client,
		rnd:    rnd,
	}
}

// Sends custom message with VK attachment
func (a *Api) SendMessageWithAttachmentAndButton(peerId int, msg string, attachment string, keyboard button.Keyboard) error {
	var (
		payload = OutcomeMessage{
			Message:     msg,
			AccessToken: a.cfg.Api.Token,
			ApiVersion:  Version,
			PeerId:      peerId,
			RandomId:    a.rnd.Rnd(),
			Attachment:  attachment,
		}
		err error
		js  []byte
	)

	if js, err = json.Marshal(keyboard); err != nil {
		a.
			logger.
			With(
				zap.Any(`request`, keyboard),
				zap.Error(err),
			).
			Errorf(`build keyboard query string error`)

		return err
	}

	payload.Keyboard = string(js)

	return a.send(payload)
}

// Sends message with keyboard
func (a *Api) SendMessageWithButton(peerId int, msg string, keyboard button.Keyboard) error {
	var (
		payload = OutcomeMessage{
			Message:     msg,
			AccessToken: a.cfg.Api.Token,
			ApiVersion:  Version,
			PeerId:      peerId,
			RandomId:    a.rnd.Rnd(),
		}
		err error
		js  []byte
	)

	if js, err = json.Marshal(keyboard); err != nil {
		a.
			logger.
			With(
				zap.Any(`request`, keyboard),
				zap.Error(err),
			).
			Errorf(`build keyboard query string error`)

		return err
	}

	payload.Keyboard = string(js)

	return a.send(payload)
}

func (a *Api) send(msgStruct OutcomeMessage) error {
	var (
		request      *http.Request
		response     *http.Response
		answer       = &Response{}
		dumpResponse []byte
		err          error
		params       url.Values
		maskedParams string
		endpoint     string
	)

	if params, err = query.Values(msgStruct); err != nil {
		a.
			logger.
			With(zap.Error(err)).
			Errorf(`build request query string error`)

		return err
	}

	endpoint = fmt.Sprintf(`%s/%s?%s`, Endpoint, apiMethodMessagesSend, params.Encode())
	maskedParams = a.cfg.Api.MaskedToken(endpoint)

	if request, err = http.NewRequest(`POST`, endpoint, nil); err != nil {
		a.
			logger.
			With(
				zap.String(`request`, maskedParams),
				zap.Error(err),
			).
			Errorf(`build Api request error`)

		return err
	}

	if response, err = a.client.Do(request); err != nil {
		a.
			logger.
			With(
				zap.String(`request`, maskedParams),
				zap.Error(err),
			).
			Errorf(`send Api request error`)

		return err
	}

	if dumpResponse, err = httputil.DumpResponse(response, true); err != nil {
		a.
			logger.
			With(
				zap.String(`request`, maskedParams),
				zap.Error(err),
			).
			Errorf(`dump Api response error`)

		return err
	}

	a.
		logger.
		With(
			zap.String(`request`, maskedParams),
			zap.ByteString(`response`, dumpResponse),
		).
		Info(`Api message sent`)

	if err = json.NewDecoder(response.Body).Decode(answer); err != nil {
		a.
			logger.
			With(
				zap.Error(err),
				zap.ByteString(`response`, dumpResponse),
			).
			Error(`error while decoding Api response`)

		return err
	}

	if len(answer.Error.Message) > 0 {
		a.
			logger.
			With(
				zap.Int32(`code`, answer.Error.Code),
				zap.String(`message`, answer.Error.Message),
			).
			Error(`failed Api answer`)
	}

	return nil
}

package server

import (
	"fmt"
	"github.com/sepuka/vkbotserver/config"
	"github.com/sepuka/vkbotserver/domain"
	"github.com/sepuka/vkbotserver/message"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestSocketServer_ServeHTTP(t *testing.T) {
	const (
		validConfirmationOutput = `this_is_a_valid_confirmation_output`

		validConfirmationMsg = `{"type": "confirmation", "group_id": 123}`
		unknownTypeMsg       = `{"type": "???", "group_id": 123}`
		invalidJsonMsg       = `{`
		emptyJsonMsg         = ``
	)

	var (
		errMsg string
		req    *http.Request
		resp   *httptest.ResponseRecorder
		cfg    = config.Config{
			Confirmation: validConfirmationOutput,
			Logger:       config.Logger{},
		}

		handler = func(handler message.Executor, req *domain.Request, resp http.ResponseWriter) error {
			return handler.Exec(req, resp)
		}
		handlerMap = message.HandlerMap{
			`confirmation`: message.NewConfirmation(cfg),
		}
		server = NewSocketServer(cfg, handlerMap, handler)

		tests = map[string]struct {
			server       *SocketServer
			incomingMsg  string
			expectedBody []byte
			expectedCode int
		}{
			`empty JSON`: {
				server:       server,
				incomingMsg:  emptyJsonMsg,
				expectedBody: []byte(invalidJSON),
				expectedCode: http.StatusBadRequest,
			},
			`invalid JSON`: {
				server:       server,
				incomingMsg:  invalidJsonMsg,
				expectedBody: []byte(invalidJSON),
				expectedCode: http.StatusBadRequest,
			},
			`unknown message type`: {
				server:       server,
				incomingMsg:  unknownTypeMsg,
				expectedBody: []byte(defaultOutput),
				expectedCode: http.StatusOK,
			},
			`valid confirmation msg`: {
				server:       server,
				incomingMsg:  validConfirmationMsg,
				expectedBody: []byte(validConfirmationOutput),
				expectedCode: http.StatusOK,
			},
		}
	)

	for testName, testCase := range tests {
		errMsg = fmt.Sprintf(`there is an unexpected error "%s"`, testName)

		resp = httptest.NewRecorder()
		req = &http.Request{
			Method: "GET",
			Host:   "vk.com",
			URL:    &url.URL{Path: "/"},
			Header: http.Header{},
			Body:   ioutil.NopCloser(strings.NewReader(testCase.incomingMsg)),
		}

		testCase.server.ServeHTTP(resp, req)

		responseBody, _ := ioutil.ReadAll(resp.Body)
		assert.Equal(t, testCase.expectedBody, responseBody, errMsg)
		assert.Equal(t, testCase.expectedCode, resp.Code, errMsg)
	}
}

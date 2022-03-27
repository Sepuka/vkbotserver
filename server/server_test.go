package server

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sepuka/vkbotserver/api"
	"github.com/sepuka/vkbotserver/api/mocks"
	"github.com/sepuka/vkbotserver/config"
	"github.com/sepuka/vkbotserver/domain"
	mocks2 "github.com/sepuka/vkbotserver/domain/mocks"
	"github.com/sepuka/vkbotserver/message"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type mistakenHandler struct{}

func (m mistakenHandler) Exec(req *domain.Request, resp http.ResponseWriter) error {
	return errors.New(`there is a persistent error`)
}

func TestSocketServer_ServeHTTP(t *testing.T) {
	const (
		validConfirmationOutput = `this_is_a_valid_confirmation_output`

		validConfirmationMsg = `{"type": "confirmation", "group_id": 123}`
		handlerWithErrorMsg  = `{"type": "mistakenHandler", "group_id": 123}`
		unknownTypeMsg       = `{"type": "???", "group_id": 123}`
		invalidJsonMsg       = `{`
		emptyJsonMsg         = ``
	)

	var (
		errMsg string
		req    *http.Request
		resp   *httptest.ResponseRecorder
		logger = zap.NewNop().Sugar()
		cfg    = config.Config{
			Confirmation: validConfirmationOutput,
			Logger:       config.Logger{},
		}
		emptyAnswer = []byte(``)

		handler = func(handler message.Executor, req *domain.Request, resp http.ResponseWriter) error {
			return handler.Exec(req, resp)
		}
		handlerMap = message.HandlerMap{
			`confirmation`:    message.NewConfirmation(cfg),
			`mistakenHandler`: mistakenHandler{},
		}
		server = NewSocketServer(cfg, handlerMap, handler, logger)

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
				expectedBody: emptyAnswer,
				expectedCode: http.StatusBadRequest,
			},
			`valid confirmation msg`: {
				server:       server,
				incomingMsg:  validConfirmationMsg,
				expectedBody: []byte(validConfirmationOutput),
				expectedCode: http.StatusOK,
			},
			`handler with error`: {
				server:       server,
				incomingMsg:  handlerWithErrorMsg,
				expectedBody: emptyAnswer,
				expectedCode: http.StatusInternalServerError,
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

func TestSocketServer_ServeHTTP_VkOauth(t *testing.T) {
	const (
		tokenResponse = `{
  "access_token": "533bacf01e11f55b536a565b57531ac114461ae8736d6506a3",
  "expires_in": 43200,
  "user_id": 66748,
  "email": "email@host.com"
}`
		cookie = `533bacf01e11f55b536a565b57531ac114461ae8736d6506a3`
	)
	var (
		errMsg          string
		incomeRequest   *http.Request
		vkTokenRequest  *http.Request
		resp            *httptest.ResponseRecorder
		vkTokenResponse *http.Response
		logger          = zap.NewNop().Sugar()
		client          = mocks.HTTPClient{}
		userRepo        = mocks2.UserRepository{}
		cfg             = config.Config{
			Logger: config.Logger{},
			VkOauth: config.VkOauth{
				ClientId:     `client_id`,
				ClientSecret: `client_secret`,
				RedirectUri:  `https://host.domain/path?args`,
				VkPath:       `vk_auth`,
			},
			PathPrefix: `/myza/`,
		}
		user *domain.User

		handler = func(handler message.Executor, req *domain.Request, resp http.ResponseWriter) error {
			return handler.Exec(req, resp)
		}
		handlerMap = message.HandlerMap{
			`vk_auth`: message.NewAuthVk(cfg.VkOauth, &client, logger, &userRepo, []domain.Callback{}),
		}
		server = NewSocketServer(cfg, handlerMap, handler, logger)
	)

	vkTokenResponse = &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(tokenResponse))),
	}
	vkTokenRequest, _ = http.NewRequest(`GET`, `https://oauth.vk.com/access_token?client_id=client_id&client_secret=client_secret&redirect_uri=https://host.domain/path?args&code=777`, nil)
	client.On(`Do`, vkTokenRequest).Return(vkTokenResponse, nil)

	user = &domain.User{Token: cookie, LastName: `some last name`, FirstName: `some first name`}
	userRepo.On(`GetByExternalId`, domain.OAuthVk, `66748`).Return(user, nil)
	// updating users`s token
	userRepo.On(`Update`, user).Return(nil)

	resp = httptest.NewRecorder()
	incomeRequest = &http.Request{
		Method: "GET",
		Host:   "vk.com",
		URL: &url.URL{
			Path:     `/myza/vk_auth`,
			RawQuery: `code=777&state=https://sepuka.github.io/somepath/`,
		},
		Header: http.Header{},
		Body:   ioutil.NopCloser(bytes.NewReader(api.DefaultResponseBody())),
	}

	server.ServeHTTP(resp, incomeRequest)

	assert.Equal(t, http.StatusFound, resp.Code, errMsg)
}

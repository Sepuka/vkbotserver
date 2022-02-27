package message

import (
	"github.com/sepuka/vkbotserver/api/mocks"
	"github.com/sepuka/vkbotserver/config"
	"github.com/sepuka/vkbotserver/domain"
	mocks2 "github.com/sepuka/vkbotserver/domain/mocks"
	"github.com/sepuka/vkbotserver/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http/httptest"
	"testing"
)

func TestYaAuth_Exec(t *testing.T) {
	var (
		logger   = zap.NewNop().Sugar()
		userRepo = mocks2.UserRepository{}
		client   = mocks.HTTPClient{}
		resp     = httptest.NewRecorder()

		incomeReq = &domain.Request{
			Type:    domain.OauthYaHandlerName,
			Context: `ya_auth#error=invalid_request&error_description=%D0%92%D1%8B%D0%B1%D1%80%D0%B0%D0%BD%D0%BD%D1%8B%D0%B9%20redirect_uri%20%D0%BD%D0%B5%D0%B1%D0%B5%D0%B7%D0%BE%D0%BF%D0%B0%D1%81%D0%B5%D0%BD`,
		}

		cfg      = config.Config{}
		executor Executor
	)

	executor = NewYaAuth(cfg.YaOauth, &client, logger, &userRepo)

	assert.ErrorIs(t, executor.Exec(incomeReq, resp), errors.OauthError)
}

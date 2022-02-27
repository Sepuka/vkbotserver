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
			Context: `ya_auth#access_token=&expires_in=86400&token_type=bearer&error=access_denied&error_description=description`,
		}

		cfg      = config.Config{}
		executor Executor
	)

	executor = NewYaAuth(cfg.YaOauth, &client, logger, &userRepo)

	assert.ErrorIs(t, executor.Exec(incomeReq, resp), errors.OauthError)
}

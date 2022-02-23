package message

import (
	"bytes"
	"fmt"
	"github.com/sepuka/vkbotserver/api/mocks"
	"github.com/sepuka/vkbotserver/config"
	"github.com/sepuka/vkbotserver/domain"
	mocks2 "github.com/sepuka/vkbotserver/domain/mocks"
	"github.com/sepuka/vkbotserver/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVkAuth_Exec_ClientIdTrouble(t *testing.T) {
	const (
		invalidClientResponse = `{"error":"invalid_client","error_description":"client_id is undefined"}`
	)

	var (
		incomeReq = &domain.Request{
			Type:    domain.OauthVkHandlerName,
			Context: `code=054d68fed17c35c307&state=https://sepuka.github.io/myzaapp/`,
		}
		expectedOutcomeReq = &http.Request{}
		expectedIncomeResp = &http.Response{}
		logger             = zap.NewNop().Sugar()
		userRepo           = mocks2.UserRepository{}
		client             = mocks.HTTPClient{}
		resp               = httptest.NewRecorder()
		oauthCfg           = config.VkOauth{
			VkPath:       `vk_auth`,
			ClientId:     `123`,
			ClientSecret: `secret`,
			RedirectUri:  `https://acc.github.io/yourapp/`,
		}
		cfg = config.Config{
			VkOauth: oauthCfg,
		}
		executor Executor
		tokenUrl string
	)

	expectedIncomeResp = &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(invalidClientResponse))),
	}
	tokenUrl = fmt.Sprintf(tokenUrlTemplate, oauthCfg.ClientId, oauthCfg.ClientSecret, oauthCfg.RedirectUri, `054d68fed17c35c307`)
	expectedOutcomeReq, _ = http.NewRequest(`GET`, tokenUrl, nil)
	client.On(`Do`, expectedOutcomeReq).Once().Return(expectedIncomeResp, nil)

	executor = NewVkAuth(cfg.VkOauth, &client, logger, &userRepo)

	assert.ErrorIs(t, executor.Exec(incomeReq, resp), errors.OauthVkError)
}

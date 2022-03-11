package users

import (
	"bytes"
	"fmt"
	"github.com/sepuka/vkbotserver/api"
	"github.com/sepuka/vkbotserver/api/mocks"
	"github.com/sepuka/vkbotserver/domain"
	mocks2 "github.com/sepuka/vkbotserver/domain/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestVkAuth_Exec_FillUser_FailedWithoutToken(t *testing.T) {
	const (
		responseUsersGet = `{"error":{"error_code":5,"error_msg":"User authorization failed: no access_token passed.","request_params":[{"key":"api_token","value":"c991fd1144d1de516fab2cb6d671ebe16a785f44ca281304414dbf0157fe0b5c0b2973d58bbc40886e1a5"},{"key":"v","value":"5.170"},{"key":"method","value":"users.get"},{"key":"oauth","value":"1"}]}}`
	)

	var (
		expectedOutcomeReq          = &http.Request{}
		expectedIncomeResp          = &http.Response{}
		logger                      = zap.NewNop().Sugar()
		userRepo                    = mocks2.UserRepository{}
		client                      = mocks.HTTPClient{}
		tokenUrl                    string
		someExistsUserWithEmptyName = &domain.User{}
	)

	expectedIncomeResp = &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(responseUsersGet))),
	}

	tokenUrl = fmt.Sprintf(apiPathTmpl, api.Endpoint, someExistsUserWithEmptyName.Token, api.Version)
	expectedOutcomeReq, _ = http.NewRequest(`GET`, tokenUrl, nil)
	client.On(`Do`, expectedOutcomeReq).Once().Return(expectedIncomeResp, nil)
	// Do not update user`s props because an error was occurred
	userRepo.On(`Update`, someExistsUserWithEmptyName).Times(0)

	userGetter := NewGet(&client, logger, &userRepo)

	userGetter.FillUser(someExistsUserWithEmptyName)
}

func TestVkAuth_Exec_FillUser_Update(t *testing.T) {
	const (
		responseUsersGet = `{"response":[{"id":557404793,"first_name":"Максим","last_name":"Шломин","can_access_closed":true,"is_closed":false}]}`
	)

	var (
		expectedOutcomeReq = &http.Request{}
		expectedIncomeResp = &http.Response{}
		logger             = zap.NewNop().Sugar()
		userRepo           = mocks2.UserRepository{}
		client             = mocks.HTTPClient{}
		tokenUrl           string
		someUser           = &domain.User{}
	)

	expectedIncomeResp = &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(responseUsersGet))),
	}

	tokenUrl = fmt.Sprintf(apiPathTmpl, api.Endpoint, someUser.Token, api.Version)
	expectedOutcomeReq, _ = http.NewRequest(`GET`, tokenUrl, nil)
	client.On(`Do`, expectedOutcomeReq).Once().Return(expectedIncomeResp, nil)
	// Do not update user`s props because an error was occurred
	userRepo.On(`Update`, someUser).Once().Return(nil)

	userGetter := NewGet(&client, logger, &userRepo)

	userGetter.FillUser(someUser)
	assert.Equal(t, `Шломин`, someUser.LastName)
	assert.Equal(t, `Максим`, someUser.FirstName)
}

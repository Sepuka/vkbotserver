package users

import (
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/sepuka/vkbotserver/api"
	"github.com/sepuka/vkbotserver/domain"
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"
)

const (
	apiPathTmpl = `%s/users.get?access_token=%s&v=%s`
)

type Get struct {
	client   api.HTTPClient
	logger   *zap.SugaredLogger
	userRepo domain.UserRepository
}

func NewGet(
	client api.HTTPClient,
	logger *zap.SugaredLogger,
	userRepo domain.UserRepository,
) *Get {
	return &Get{
		client:   client,
		logger:   logger,
		userRepo: userRepo,
	}
}

func (o *Get) FillUser(user *domain.User) {
	var (
		err         error
		path        = fmt.Sprintf(apiPathTmpl, api.Endpoint, user.Token, api.Version)
		response    *http.Response
		request     *http.Request
		dump        []byte
		apiResponse = &domain.ApiResponse{}
		apiUser     *domain.VkUser
	)

	if request, err = http.NewRequest(`GET`, path, nil); err != nil {
		o.
			logger.
			With(
				zap.Error(err),
				zap.String(`url`, path),
				zap.String(`api`, `users.get`),
			).
			Error(`Build API request error`)

		return
	}

	if response, err = o.client.Do(request); err != nil {
		o.
			logger.
			With(
				zap.Error(err),
				zap.String(`url`, path),
				zap.String(`api`, `users.get`),
			).
			Error(`Send API request error`)

		return
	}

	if dump, err = httputil.DumpResponse(response, true); err != nil {
		o.
			logger.
			With(
				zap.Error(err),
				zap.Int64(`size`, response.ContentLength),
				zap.Int(`code`, response.StatusCode),
				zap.String(`api`, `users.get`),
			).
			Error(`Dump response error`)

		return
	}

	o.
		logger.
		With(
			zap.String(`api`, `users.get`),
			zap.ByteString(`response`, dump),
		).
		Info(`VK API response`)

	if err = easyjson.UnmarshalFromReader(response.Body, apiResponse); err != nil {
		o.
			logger.
			With(
				zap.Error(err),
				zap.String(`api`, `users.get`),
			).
			Error(`Unmarshalling oauth response error`)

		return
	}

	if apiResponse.Error.ErrorCode > 0 {
		o.
			logger.
			With(
				zap.String(`message`, apiResponse.Error.ErrorMessage),
				zap.Int(`code`, apiResponse.Error.ErrorCode),
				zap.String(`api`, `users.get`),
			).
			Error(`Response has an error`)

		return
	}

	apiUser = &apiResponse.Response[0]
	user.FirstName = apiUser.FirstName
	user.LastName = apiUser.LastName

	if err = o.userRepo.Update(user); err != nil {
		o.
			logger.
			With(
				zap.Error(err),
				zap.String(`api`, `users.get`),
			).
			Error(`Update user info error`)
	}
}

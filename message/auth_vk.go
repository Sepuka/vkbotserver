package message

import (
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/sepuka/vkbotserver/api"
	"github.com/sepuka/vkbotserver/config"
	"github.com/sepuka/vkbotserver/domain"
	"github.com/sepuka/vkbotserver/errors"
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

const (
	OauthVkLogKey    = `VK`
	ver              = `5.170`
	tokenUrlTemplate = `https://oauth.vk.com/access_token?client_id=%s&client_secret=%s&redirect_uri=%s&code=%s`
)

type authVk struct {
	cfg      config.VkOauth
	client   api.HTTPClient
	logger   *zap.SugaredLogger
	userRepo domain.UserRepository
}

// NewAuthVk creates an instance VK VkOauth handler
func NewAuthVk(
	cfg config.VkOauth,
	client api.HTTPClient,
	logger *zap.SugaredLogger,
	userRepo domain.UserRepository,
) *authVk {
	return &authVk{
		cfg:      cfg,
		client:   client,
		logger:   logger,
		userRepo: userRepo,
	}
}

func (o *authVk) Exec(req *domain.Request, resp http.ResponseWriter) error {
	const (
		urlPartCode  = `code`
		urlPartState = `state`
	)

	var (
		tokenUrl          string
		err               error
		args              url.Values
		tokenHttpResponse *http.Response
		tokenHttpRequest  *http.Request
		dumpResponse      []byte
		tokenResponse     = &domain.OauthVkTokenResponse{}
		user              *domain.User
		redirectUrl       string
	)

	if args, err = url.ParseQuery(req.Context.(string)); err != nil {
		return err
	}

	tokenUrl = fmt.Sprintf(tokenUrlTemplate, o.cfg.ClientId, o.cfg.ClientSecret, o.cfg.RedirectUri, args[urlPartCode][0])

	if tokenHttpRequest, err = http.NewRequest(`GET`, tokenUrl, nil); err != nil {
		o.
			logger.
			With(
				zap.Error(err),
				zap.String(`oauth`, OauthVkLogKey),
				zap.String(`url`, tokenUrl),
			).
			Error(`Build oauth token API request error`)

		return err
	}

	if tokenHttpResponse, err = o.client.Do(tokenHttpRequest); err != nil {
		o.
			logger.
			With(
				zap.Error(err),
				zap.String(`oauth`, OauthVkLogKey),
				zap.String(`url`, tokenUrl),
			).
			Error(`Send oauth token API request error`)

		return err
	}

	if dumpResponse, err = httputil.DumpResponse(tokenHttpResponse, true); err != nil {
		o.
			logger.
			With(
				zap.Error(err),
				zap.String(`oauth`, OauthVkLogKey),
				zap.Int64(`size`, tokenHttpResponse.ContentLength),
				zap.Int(`code`, tokenHttpResponse.StatusCode),
			).
			Error(`Dump API oauth response error`)

		return err
	}

	o.
		logger.
		With(
			zap.String(`oauth`, OauthVkLogKey),
			zap.ByteString(`response`, dumpResponse),
		).
		Info(`Oauth API response`)

	if err = easyjson.UnmarshalFromReader(tokenHttpResponse.Body, tokenResponse); err != nil {
		o.
			logger.
			With(
				zap.Error(err),
				zap.String(`oauth`, OauthVkLogKey),
			).
			Error(`Unmarshalling oauth response error`)

		return err
	}

	if len(tokenResponse.Error) > 0 {
		o.
			logger.
			With(
				zap.String(`oauth`, OauthVkLogKey),
				zap.String(`description`, tokenResponse.ErrorDescription),
			).
			Error(`could not authorize`)

		return errors.NewOauthError(tokenResponse.Error)
	}

	if user, err = o.userRepo.GetByExternalId(domain.OAuthVk, strconv.Itoa(tokenResponse.UserId)); err != nil {
		if err == errors.NoUserFound {
			user = &domain.User{
				CreatedAt:  time.Now(),
				OAuth:      domain.OAuthVk,
				ExternalId: strconv.Itoa(tokenResponse.UserId),
				Email:      tokenResponse.Email,
				Token:      tokenResponse.Token,
			}
			if err = o.userRepo.Create(user); err != nil {
				o.
					logger.
					With(
						zap.String(`oauth`, OauthVkLogKey),
						zap.Error(err),
					).
					Error(`could not create oauth user`)

				return err
			}
		} else {
			o.
				logger.
				With(
					zap.String(`oauth`, OauthVkLogKey),
					zap.Error(err),
				).
				Error(`could not find oauth user`)

			return err
		}
	}

	go o.fillUser(user)

	redirectUrl = fmt.Sprintf(`%s?token=%s`, args[urlPartState][0], user.Token)
	http.Redirect(resp, &http.Request{}, redirectUrl, http.StatusFound)

	return nil
}

func (o *authVk) fillUser(user *domain.User) {
	const (
		apiPathTmpl = `https://api.vk.com/method/users.get?api_token=%s&v=%s`
	)

	if user.IsFilledPersonalData() {
		return
	}

	var (
		err      error
		path     = fmt.Sprintf(apiPathTmpl, user.Token, ver)
		response *http.Response
		request  *http.Request
		dump     []byte
		data     = &domain.UsersGetResponse{}
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
		Info(`Oauth API response`)

	if err = easyjson.UnmarshalFromReader(response.Body, data); err != nil {
		o.
			logger.
			With(
				zap.Error(err),
				zap.String(`api`, `users.get`),
			).
			Error(`Unmarshalling oauth response error`)

		return
	}

	user.FistName = data.FirstName
	user.LastName = data.LastName

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

func (o *authVk) String() string {
	return domain.OauthVkHandlerName
}

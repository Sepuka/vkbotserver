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
	urlPartCode      = `code`
	tokenUrlTemplate = `https://oauth.vk.com/access_token?client_id=%s&client_secret=%s&redirect_uri=%s&code=%s`
)

type vkAuth struct {
	cfg      config.VkOauth
	client   api.HTTPClient
	logger   *zap.SugaredLogger
	userRepo domain.UserRepository
}

// NewVkAuth creates an instance VK VkOauth handler
func NewVkAuth(
	cfg config.VkOauth,
	client api.HTTPClient,
	logger *zap.SugaredLogger,
	userRepo domain.UserRepository,
) *vkAuth {
	return &vkAuth{
		cfg:      cfg,
		client:   client,
		logger:   logger,
		userRepo: userRepo,
	}
}

func (o *vkAuth) Exec(req *domain.Request, resp http.ResponseWriter) error {
	var (
		tokenUrl          string
		err               error
		args              url.Values
		tokenHttpResponse *http.Response
		tokenHttpRequest  *http.Request
		dumpResponse      []byte
		tokenResponse     = &domain.OauthVkTokenResponse{}
		user              *domain.User
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
				zap.String(`oauth`, `VK`),
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
				zap.String(`oauth`, `VK`),
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
				zap.String(`oauth`, `VK`),
				zap.Int64(`size`, tokenHttpResponse.ContentLength),
				zap.Int(`code`, tokenHttpResponse.StatusCode),
			).
			Error(`Dump API oauth response error`)

		return err
	}

	o.
		logger.
		With(
			zap.String(`oauth`, `VK`),
			zap.ByteString(`response`, dumpResponse),
		).
		Info(`Oauth API response`)

	if err = easyjson.UnmarshalFromReader(tokenHttpResponse.Body, tokenResponse); err != nil {
		o.
			logger.
			With(
				zap.Error(err),
				zap.String(`oauth`, `VK`),
			).
			Error(`Unmarshalling oauth response error`)

		return err
	}

	if len(tokenResponse.Error) > 0 {
		o.
			logger.
			With(
				zap.String(`oauth`, `VK`),
				zap.String(`description`, tokenResponse.ErrorDescription),
			).
			Error(`could not authorize`)

		return errors.NewOauthVkError(tokenResponse.Error)
	}

	if user, err = o.userRepo.GetByExternalId(domain.OAuthVk, strconv.Itoa(tokenResponse.UserId)); err != nil {
		if err == errors.NoUserFound {
			user = &domain.User{
				CreatedAt:  time.Time{},
				UpdatedAt:  time.Time{},
				OAuth:      domain.OAuthVk,
				ExternalId: strconv.Itoa(tokenResponse.UserId),
				Email:      tokenResponse.Email,
				Token:      tokenResponse.Token,
			}
			if err = o.userRepo.Create(user); err != nil {
				o.
					logger.
					With(
						zap.String(`oauth`, `VK`),
						zap.Error(err),
					).
					Error(`could not create oauth user`)

				return err
			}
		} else {
			o.
				logger.
				With(
					zap.String(`oauth`, `VK`),
					zap.Error(err),
				).
				Error(`could not find oauth user`)

			return err
		}
	}

	http.SetCookie(resp, &http.Cookie{Name: domain.CookieName, Value: user.Token})

	if _, err = resp.Write(api.DefaultResponseBody()); err != nil {
		o.
			logger.
			With(
				zap.Error(err),
				zap.String(`oauth`, `VK`),
			).
			Error(`Could not save response body`)
	}

	return nil
}

func (o *vkAuth) String() string {
	return domain.OauthVkHandlerName
}

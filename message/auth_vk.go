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
	tokenUrlTemplate = `https://oauth.vk.com/access_token?client_id=%s&client_secret=%s&redirect_uri=%s&code=%s`
)

type authVk struct {
	cfg       config.VkOauth
	client    api.HTTPClient
	logger    *zap.SugaredLogger
	userRepo  domain.UserRepository
	callbacks []domain.Callback
}

// NewAuthVk creates an instance VK VkOauth handler
func NewAuthVk(
	cfg config.VkOauth,
	client api.HTTPClient,
	logger *zap.SugaredLogger,
	userRepo domain.UserRepository,
	callbacks []domain.Callback,
) *authVk {
	return &authVk{
		cfg:       cfg,
		client:    client,
		logger:    logger,
		userRepo:  userRepo,
		callbacks: callbacks,
	}
}

func (o *authVk) Exec(req *domain.Request, resp http.ResponseWriter) error {
	const (
		urlPartCode = `code`
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
		siteUrl           *url.URL
		callback          domain.Callback
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
	} else {
		user.Token = tokenResponse.Token
		if err = o.userRepo.Update(user); err != nil {
			o.
				logger.
				With(
					zap.Error(err),
				).
				Error(`Update user's token error`)
		}
	}

	for _, callback = range o.callbacks {
		go callback()
	}

	if siteUrl, err = url.Parse(o.cfg.RedirectUri); err != nil {
		o.
			logger.
			With(
				zap.Error(err),
				zap.String(`cfg url`, o.cfg.RedirectUri),
			).
			Error(`Could not build redirect url`)

		return nil
	}
	siteUrl.Path = `/`

	http.SetCookie(resp, &http.Cookie{Name: domain.CookieName, Value: user.Token})
	http.Redirect(resp, &http.Request{}, siteUrl.String(), http.StatusFound)

	return nil
}

func (o *authVk) String() string {
	return domain.OauthVkHandlerName
}

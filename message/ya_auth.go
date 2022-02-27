package message

import (
	"github.com/sepuka/vkbotserver/api"
	"github.com/sepuka/vkbotserver/config"
	"github.com/sepuka/vkbotserver/domain"
	"github.com/sepuka/vkbotserver/errors"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

const (
	OauthYaLogKey = `YA`
)

type YaAuth struct {
	cfg      config.YaOauth
	client   api.HTTPClient
	logger   *zap.SugaredLogger
	userRepo domain.UserRepository
}

func NewYaAuth(
	cfg config.YaOauth,
	client api.HTTPClient,
	logger *zap.SugaredLogger,
	userRepo domain.UserRepository,
) *YaAuth {
	return &YaAuth{
		cfg:      cfg,
		client:   client,
		logger:   logger,
		userRepo: userRepo,
	}
}

func (o *YaAuth) Exec(req *domain.Request, resp http.ResponseWriter) error {
	const (
		urlPartToken       = `access_token`
		errPartCode        = `error`
		errPartDescription = `error_description`
	)

	var (
		args      url.Values
		err       error
		errorCode string
		//user *domain.User
	)

	if args, err = url.ParseQuery(req.Context.(string)); err != nil {
		return err
	}

	errorCode = args.Get(errPartCode)

	if len(errorCode) > 0 {
		return errors.NewOauthError(errorCode)
	}

	o.
		logger.
		With(
			zap.String(`oauth`, OauthYaLogKey),
		).
		Info(`Oauth API response`)

	//if user, err = o.userRepo.GetByExternalId(domain.OAuthVk, strconv.Itoa(tokenResponse.UserId)); err != nil {
	//	if err == errors.NoUserFound {
	//		user = &domain.User{
	//			CreatedAt:  time.Now(),
	//			OAuth:      domain.OAuthVk,
	//			ExternalId: strconv.Itoa(tokenResponse.UserId),
	//			Email:      tokenResponse.Email,
	//			Token:      tokenResponse.Token,
	//		}
	//		if err = o.userRepo.Create(user); err != nil {
	//			o.
	//				logger.
	//				With(
	//					zap.String(`oauth`, OauthYaLogKey),
	//					zap.Error(err),
	//				).
	//				Error(`could not create oauth user`)
	//
	//			return err
	//		}
	//	} else {
	//		o.
	//			logger.
	//			With(
	//				zap.String(`oauth`, OauthYaLogKey),
	//				zap.Error(err),
	//			).
	//			Error(`could not find oauth user`)
	//
	//		return err
	//	}
	//}
	//
	//http.SetCookie(resp, &http.Cookie{Name: domain.CookieName, Value: user.Token})
	//http.Redirect(resp, &http.Request{}, args[urlPartState][0], http.StatusFound)

	return nil
}

func (o *YaAuth) String() string {
	return domain.OauthYaHandlerName
}

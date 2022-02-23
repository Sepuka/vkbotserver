package domain

const (
	CookieName         = `token`
	OauthVkHandlerName = `vk_auth`

	OAuthVk Oauth = 1 + iota
)

type (
	Oauth      uint8
	OauthStore interface {
		GetToken(authCookie string) (authToken interface{}, err error)
		SetToken(authToken interface{}) (cookie string, err error)
	}
)

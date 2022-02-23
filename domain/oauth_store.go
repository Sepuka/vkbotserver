package domain

const (
	CookieName         = `token`
	OauthVkHandlerName = `vk_auth`

	AuthVk = iota
)

type (
	OauthStore interface {
		GetToken(authCookie string) (authToken interface{}, err error)
		SetToken(authToken interface{}) (cookie string, err error)
	}
)

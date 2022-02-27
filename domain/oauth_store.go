package domain

import "time"

const (
	CookieName         = `token`
	OauthVkHandlerName = `vk_auth`
)

const (
	OAuthVk Oauth = iota
	OAuthYandex
)

type (
	Oauth uint8

	User struct {
		UserId     int       `sql:",pk"`
		CreatedAt  time.Time `pg:"notnull"`
		UpdatedAt  time.Time `pg:"default:now(),notnull"`
		OAuth      Oauth     `pg:"notnull"`
		ExternalId string    `pg:"notnull"`
		Email      string    `pg:"notnull"`
		Token      string
	}

	UserRepository interface {
		GetByExternalId(auth Oauth, id string) (*User, error)
		Create(user *User) error
	}
)

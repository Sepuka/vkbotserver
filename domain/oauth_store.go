package domain

import "time"

const (
	CookieName               = `token`
	OauthVkHandlerName       = `vk_auth`
	OauthYaHandlerName       = `ya_auth`
	OAuthYa            Oauth = 1
	OAuthVk            Oauth = 2
)

type (
	Oauth uint8

	// User refers to internal client who's bound with some Oauth network
	User struct {
		UserId     int       `sql:",pk"`
		CreatedAt  time.Time `pg:"notnull"`
		UpdatedAt  time.Time `pg:"default:now(),notnull"`
		OAuth      Oauth     `pg:"notnull"`
		ExternalId string    `pg:"notnull"`
		Email      string    `pg:"notnull"`
		Token      string
		FirstName  string
		LastName   string
	}

	// UserRepository offers an interface for create and fetch clients
	UserRepository interface {
		GetByExternalId(auth Oauth, id string) (*User, error)
		Create(user *User) error
		Update(user *User) error
	}
)

func (u *User) IsFilledPersonalData() bool {
	return u.LastName != `` || u.FistName != ``
}

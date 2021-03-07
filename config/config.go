package config

import (
	"fmt"
	"strings"
)

type (
	// Logger writes to stdout
	Logger struct {
		Prod bool
	}
	// API config
	Api struct {
		Token string `default:"???_there_is_the_access_api_token"`
	}
)

// Config is struct which is filling by config from App path like /etc/app.yml
type Config struct {
	Confirmation string
	Socket       string `default:"/var/run/vkbotserver.sock"`
	Logger       Logger
	Api          Api
}

func (api *Api) MaskedToken(params string) string {
	var maskedToken = fmt.Sprintf(`%s...`, api.Token[0:3])

	return strings.Replace(params, api.Token, maskedToken, 1)
}

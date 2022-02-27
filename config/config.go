package config

import (
	"fmt"
	"strings"
	"time"
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

	// Cache requests
	// Ttl in ns, default is 1 sec
	Cache struct {
		Enabled bool
		Ttl     time.Duration `default:"1000000000"`
	}

	VkOauth struct {
		VkPath       string
		ClientId     string
		ClientSecret string
		RedirectUri  string
	}

	YaOauth struct {
		Path string
	}
)

// Config is the struct which is filling by config from App path like /etc/app.yml
type Config struct {
	Confirmation string
	Socket       string `default:"/var/run/vkbotserver.sock"`
	Logger       Logger
	Api          Api
	Cache        Cache
	VkOauth      VkOauth
	YaOauth      YaOauth
	// if your web-server configured to handle VKbot-requests with some prefix
	// like /mybot/ rewrite this opt
	PathPrefix string `default:"/"`
}

func (api *Api) MaskedToken(params string) string {
	var maskedToken = fmt.Sprintf(`%s...`, api.Token[0:3])

	return strings.Replace(params, api.Token, maskedToken, 1)
}

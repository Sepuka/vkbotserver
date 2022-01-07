package message

import (
	"github.com/sepuka/vkbotserver/api/button"
	"github.com/sepuka/vkbotserver/domain"
)

type (
	Handler interface {
		Handle(*domain.Request, *button.Payload) error
	}
)

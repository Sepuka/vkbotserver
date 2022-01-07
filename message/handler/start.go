package handler

import (
	"github.com/sepuka/vkbotserver/api"
	"github.com/sepuka/vkbotserver/api/button"
	"github.com/sepuka/vkbotserver/domain"
)

const (
	msg = `Hello world!`
)

type (
	startHandler struct {
		api *api.Api
	}
)

func NewStartHandler(api *api.Api) *startHandler {
	return &startHandler{api: api}
}

func (h *startHandler) Handle(req *domain.Request, payload *button.Payload) error {
	var (
		peerId = int(req.Object.Message.FromId)
	)

	return h.api.SendMessage(peerId, msg)
}

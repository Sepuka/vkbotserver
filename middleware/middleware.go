package middleware

import (
	"github.com/sepuka/vkbotserver/domain"
	"github.com/sepuka/vkbotserver/message"
	"net/http"
)

type HandlerFunc func(message.Executor, *domain.Request, http.ResponseWriter) error

func final(handler message.Executor, req *domain.Request, resp http.ResponseWriter) error {
	return handler.Exec(req, resp)
}

func BuildHandlerChain(handlers []func(HandlerFunc) HandlerFunc) HandlerFunc {
	if len(handlers) == 0 {
		return final
	}

	return handlers[0](BuildHandlerChain(handlers[1:]))
}

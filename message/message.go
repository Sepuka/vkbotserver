package message

import (
	"github.com/sepuka/vkbotserver/domain"
	"net/http"
)

type HandlerMap map[string]Executor

type Executor interface {
	Exec(*domain.Request, http.ResponseWriter) error
}

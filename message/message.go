package message

import (
	"github.com/sepuka/vkbotserver/domain"
	"net/http"
)

// each msg handler must implement this interface
type HandlerMap map[string]Executor

type Executor interface {
	Exec(*domain.Request, http.ResponseWriter) error
}

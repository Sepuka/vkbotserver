package middleware

import (
	"errors"
	"fmt"
	"github.com/sepuka/vkbotserver/domain"
	"github.com/sepuka/vkbotserver/message"
	"net/http"
	"runtime/debug"
)

func Panic(next HandlerFunc) HandlerFunc {
	return func(exec message.Executor, req *domain.Request, writer http.ResponseWriter) error {
		var err error

		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("panic: %s\n"+
					"command `%s`\n"+
					"stacktrace from panic: %s\n",
					r, req.Type, string(debug.Stack()))
				err = errors.New(`internal error`)
			}
		}()

		err = next(exec, req, writer)

		return err
	}
}

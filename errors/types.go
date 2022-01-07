package errors

import "errors"

var (
	InvalidJson = errors.New(`invalid JSON`)
)

// NewInvalidJsonError instance an InvalidJson error
func NewInvalidJsonError(msg string, originalErr error) BotError {
	return BotError{
		err:           InvalidJson,
		message:       msg,
		originalError: originalErr,
	}
}

package errors

import "errors"

var (
	InvalidJson       = errors.New(`invalid JSON`)
	NotIsOAuthRequest = errors.New(`not is an OAuth request`)
)

// NewInvalidJsonError instance an InvalidJson error
func NewInvalidJsonError(msg string, originalErr error) BotError {
	return BotError{
		err:           InvalidJson,
		message:       msg,
		originalError: originalErr,
	}
}

func NewNotIsOAuthReqError() BotError {
	return BotError{
		err: NotIsOAuthRequest,
	}
}

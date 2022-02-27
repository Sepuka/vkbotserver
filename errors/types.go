package errors

import "errors"

var (
	InvalidJson       = errors.New(`invalid JSON`)
	NotIsOAuthRequest = errors.New(`not is an OAuth request`)
	OauthError        = errors.New(`oauth error`)
	NoUserFound       = errors.New(`there are any user was found`)
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

func NewOauthError(msg string) BotError {
	return BotError{
		err:     OauthError,
		message: msg,
	}
}

func NewNoUserFound() BotError {
	return BotError{
		err: NoUserFound,
	}
}

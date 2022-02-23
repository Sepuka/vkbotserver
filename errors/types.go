package errors

import "errors"

var (
	InvalidJson       = errors.New(`invalid JSON`)
	NotIsOAuthRequest = errors.New(`not is an OAuth request`)
	OauthVkError      = errors.New(`oauth VK error`)
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

func NewOauthVkError(msg string) BotError {
	return BotError{
		err:     OauthVkError,
		message: msg,
	}
}

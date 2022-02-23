package errors

import "fmt"

type BotError struct {
	err           error
	message       string
	originalError error
	context       map[string]string
}

func (e BotError) Error() string {
	if e.originalError != nil {
		return fmt.Sprintf(`%s (%s)`, e.message, e.originalError)
	}

	return fmt.Sprintf(`%s`, e.message)
}

func (e BotError) Is(target error) bool {
	return e.err == target
}

func (e BotError) Unwrap() error {
	return e.err
}

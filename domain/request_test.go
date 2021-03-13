package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequest_IsKeyboardButton(t *testing.T) {
	var (
		test = map[string]struct {
			req      Request
			isButton bool
		}{
			`text only`: {
				req: Request{
					Object: Object{
						Message: Message{},
					},
				},
				isButton: false,
			},
			`start button`: {
				req: Request{
					Object: Object{
						Message: Message{
							Payload: `{"command":"start"}`,
						},
					},
				},
				isButton: true,
			},
		}
	)

	for _, testCase := range test {
		assert.Equal(t, testCase.isButton, testCase.req.IsKeyboardButton())
	}
}

package message

import (
	"github.com/sepuka/vkbotserver/config"
	"github.com/sepuka/vkbotserver/domain"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestConfirmation_Exec(t *testing.T) {
	const confirmationString = `this_is_confirmation_string`

	var (
		req  = &domain.Request{}
		resp = httptest.NewRecorder()
		cfg  = config.Config{
			Confirmation: confirmationString,
		}
		executor = NewConfirmation(cfg)
	)

	assert.Nil(t, executor.Exec(req, resp))
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, confirmationString, string(body))
}

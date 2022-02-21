package message

import (
	"github.com/sepuka/vkbotserver/config"
	"github.com/sepuka/vkbotserver/domain"
	"net/http"
)

// handler for confirmation-requests
type confirmation struct {
	cfg config.Config
}

// NewConfirmation creates a confirmation handler
func NewConfirmation(cfg config.Config) *confirmation {
	return &confirmation{
		cfg: cfg,
	}
}

func (o *confirmation) Exec(req *domain.Request, resp http.ResponseWriter) error {
	_, err := resp.Write([]byte(o.cfg.Confirmation))

	return err
}

func (o *confirmation) String() string {
	return `confirmation`
}

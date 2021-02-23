package message

import (
	"github.com/sepuka/vkbotserver/config"
	"github.com/sepuka/vkbotserver/domain"
	"net/http"
)

type Confirmation struct {
	cfg config.Config
}

func NewConfirmation(cfg config.Config) *Confirmation {
	return &Confirmation{
		cfg: cfg,
	}
}

func (o *Confirmation) Exec(req *domain.Request, resp http.ResponseWriter) error {
	_, err := resp.Write([]byte(o.cfg.Confirmation))

	return err
}

func (o *Confirmation) String() string {
	return `confirmation`
}

package provider

import (
	"fmt"
	"net/smtp"

	"github.com/kubex-ecosystem/kbx/types"
)

type MicrosoftProvider struct{}

func (m MicrosoftProvider) Send(cfg *types.SMTPConfig, msg *types.Email) error {
	auth := smtp.PlainAuth("", cfg.User, cfg.Pass, cfg.Host)

	c, err := smtp.Dial(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		return err
	}
	defer c.Close()

	if err := c.Auth(auth); err != nil {
		return err
	}

	return sendSMTPMessage(c, msg)
}

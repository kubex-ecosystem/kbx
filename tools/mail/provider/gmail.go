package provider

import (
	"crypto/tls"
	"net/smtp"

	"github.com/kubex-ecosystem/kbx/types"
)

type GmailProvider struct{}

func (g GmailProvider) Send(cfg *types.SMTPConfig, msg *types.Email) error {
	auth := smtp.PlainAuth("", cfg.User, cfg.Pass, cfg.Host)

	tlsCfg := &tls.Config{
		ServerName: cfg.Host,
	}

	conn, err := tls.Dial("tcp", formatAddr(cfg.Host, cfg.Port), tlsCfg)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, cfg.Host)
	if err != nil {
		return err
	}

	if err := client.Auth(auth); err != nil {
		return err
	}

	defer client.Close()

	return sendSMTPMessage(client, msg)
}

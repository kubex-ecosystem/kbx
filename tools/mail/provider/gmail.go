package provider

import (
	"crypto/tls"
	"net/smtp"

	"github.com/kubex-ecosystem/kbx/types"
)

type GmailProvider struct{}

func (g GmailProvider) Send(cfg *types.MailConnection, msg *types.Email) error {
	auth := smtp.PlainAuth("", cfg.User, cfg.Pass, cfg.Host)
	addr := formatAddr(cfg.Host, cfg.Port)

	// Porta 465 ou use_ssl => TLS direto
	if cfg.SSL || cfg.Port == 465 {
		tlsCfg := &tls.Config{ServerName: cfg.Host}
		conn, err := tls.Dial("tcp", addr, tlsCfg)
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

	// Porta 587 ou use_tls => STARTTLS
	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Hello(cfg.Host); err != nil {
		return err
	}
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsCfg := &tls.Config{ServerName: cfg.Host}
		if err := client.StartTLS(tlsCfg); err != nil {
			return err
		}
	}

	if err := client.Auth(auth); err != nil {
		return err
	}

	return sendSMTPMessage(client, msg)
}

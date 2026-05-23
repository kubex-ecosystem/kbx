// Package mail provides email sending functionality.
package mail

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/kubex-ecosystem/kbx/types"
)

// Send sends an email using the given SMTP connection configuration.
func Send(cfg *types.MailConnection, msg *types.Email) error {
	if cfg == nil {
		return fmt.Errorf("mail: connection config is nil")
	}
	if msg == nil {
		return fmt.Errorf("mail: message is nil")
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	auth := smtp.PlainAuth("", cfg.User, cfg.Pass, cfg.Host)

	to := msg.To
	if len(to) == 0 {
		return fmt.Errorf("mail: no recipients specified")
	}

	body := buildMessage(cfg, msg)
	return smtp.SendMail(addr, auth, msg.From, to, []byte(body))
}

func buildMessage(cfg *types.MailConnection, msg *types.Email) string {
	var sb strings.Builder
	from := msg.From
	if from == "" {
		from = cfg.User
	}
	sb.WriteString("From: " + from + "\r\n")
	sb.WriteString("To: " + strings.Join(msg.To, ", ") + "\r\n")
	sb.WriteString("Subject: " + msg.Subject + "\r\n")
	if msg.HTML != "" {
		sb.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		sb.WriteString("\r\n")
		sb.WriteString(msg.HTML)
	} else {
		sb.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		sb.WriteString("\r\n")
		sb.WriteString(msg.Text)
	}
	return sb.String()
}

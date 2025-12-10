package mailing

import (
	"context"
	"errors"
	"time"

	"github.com/kubex-ecosystem/kbx/mailing/templates"
	"github.com/kubex-ecosystem/kbx/tools"
	"github.com/kubex-ecosystem/kbx/tools/mail"
	"github.com/kubex-ecosystem/kbx/types"
)

var errNilRequest = errors.New("mailing: mail request is nil")

// Config parametriza o envio com retry/timeout via tools.Retry.
type Config struct {
	SMTP  types.SMTPConfig
	Retry tools.RetryConfig
}

// Mailer expõe a API única usada pelo backend.
type Mailer struct {
	cfg Config
}

// NewMailer cria um Mailer com defaults para retry/timeout se não informados.
func NewMailer(cfg Config) *Mailer {
	if cfg.Retry.Retries <= 0 {
		cfg.Retry.Retries = 3
	}
	if cfg.Retry.Delay == 0 {
		cfg.Retry.Delay = 2 * time.Second
	}
	if cfg.Retry.Timeout == 0 {
		cfg.Retry.Timeout = 5 * time.Second
	}
	return &Mailer{cfg: cfg}
}

// Send dispara um e-mail convertendo MailRequest -> types.Email e delegando para tools/mail.
// Respeita retry/timeout configurados via tools.Retry.
func (m *Mailer) Send(ctx context.Context, req *MailRequest) error {
	if req == nil {
		return errNilRequest
	}
	email := req.ToEmail()

	_, err := tools.Retry(func() (struct{}, error) {
		select {
		case <-ctx.Done():
			return struct{}{}, ctx.Err()
		default:
		}
		return struct{}{}, mail.Send(m.cfg.SMTP, email)
	},
		tools.WithRetries(m.cfg.Retry.Retries),
		tools.WithDelay(m.cfg.Retry.Delay),
		tools.WithTimeout(m.cfg.Retry.Timeout),
	)

	return err
}

// SendTemplate aplica o template loader + render e envia.
func (m *Mailer) SendTemplate(ctx context.Context, loader templates.TemplateLoader, name string, data any, to string, subject string, from string) error {
	if loader == nil {
		return errNilRequest
	}
	htmlTmpl, err := loader.LoadHTML(name)
	if err != nil {
		return err
	}
	html, err := RenderHTML(htmlTmpl, data)
	if err != nil {
		return err
	}
	req := &MailRequest{
		From:    from,
		To:      []string{to},
		Subject: subject,
		HTML:    html,
	}
	return m.Send(ctx, req)
}

package mailing

import (
	"context"
	"errors"
	"time"

	"github.com/kubex-ecosystem/kbx/load"
	"github.com/kubex-ecosystem/kbx/mailing/templates"
	"github.com/kubex-ecosystem/kbx/tools"
	"github.com/kubex-ecosystem/kbx/tools/mail"
	"github.com/kubex-ecosystem/kbx/types"

	gl "github.com/kubex-ecosystem/logz"
)

var errNilRequest = errors.New("mailing: mail request is nil")

// MailConfig parametriza o envio com retry/timeout via tools.Retry.
type MailConfig = load.MailConfig

// Mailer expõe a API única usada pelo backend.
type Mailer struct {
	*MailConfig `json:",inline" yaml:",inline" xml:"-" toml:",inline" mapstructure:",squash"`
	Sender      types.MailProvider `json:"-" yaml:"-" xml:"-" toml:"-" mapstructure:"-"`
}

// NewMailer cria um Mailer com defaults para retry/timeout se não informados.
func NewMailer(cfg *MailConfig) *Mailer {
	for _, conn := range cfg.Connections {
		if conn.RetryCount <= 0 {
			conn.RetryCount = 3
		}
		if conn.RetryInterval == 0 {
			conn.RetryInterval = 2 * time.Second
		}
		if conn.Timeout == 0 {
			conn.Timeout = 5 * time.Second
		}
	}
	return &Mailer{MailConfig: cfg}
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
		return struct{}{}, mail.Send(m.GetSMTPConnection(), email)
	},
		tools.WithRetries(m.GetSMTPConnection().RetryCount),
		tools.WithDelay(m.GetSMTPConnection().RetryInterval),
		tools.WithTimeout(m.GetSMTPConnection().Timeout),
	)

	return err
}

func (m *Mailer) GetSMTPConnection() *types.MailConnection {
	if m == nil || m.Connections == nil {
		conn := types.NewMailConnection()
		conn.Protocol = "smtp"
		conn.RetryCount = 3
		conn.RetryInterval = 2 * time.Second
		conn.Timeout = 5 * time.Second
		gl.Warn("mailing: no SMTP connection found in config, using default parameters")
		return conn
	}
	for _, conn := range m.Connections {
		if conn.Protocol == "smtp" || conn.Protocol == "" {
			return conn
		}
	}
	return nil
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

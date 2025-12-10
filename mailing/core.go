package mailing

import "github.com/kubex-ecosystem/kbx/types"

// MailRequest é o DTO de envio usado pelo BE.
type MailRequest struct {
	From        string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	HTML        string
	Text        string
	Attachments []types.Attachment
}

// ToEmail converte para o tipo nativo de tools/mail.
func (r MailRequest) ToEmail() *types.Email {
	return &types.Email{
		From:        r.From,
		To:          r.To,
		Cc:          r.Cc,
		Bcc:         r.Bcc,
		Subject:     r.Subject,
		Text:        r.Text,
		HTML:        r.HTML,
		Attachments: r.Attachments,
	}
}

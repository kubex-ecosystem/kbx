// Package mailing fornece funcionalidades relacionadas ao envio de emails.
package mailing

import "github.com/kubex-ecosystem/kbx/types"

// MailRequest é o DTO de envio usado pelo BE.
type MailRequest struct {
	Name        string             `json:"from_name,omitempty" yaml:"from_name,omitempty" xml:"from_name" toml:"from_name" mapstructure:"from_name"`
	From        string             `json:"from_email,omitempty" yaml:"from_email,omitempty" xml:"from_email" toml:"from_email" mapstructure:"from_email"`
	To          []string           `json:"to,omitempty" yaml:"to,omitempty" xml:"to" toml:"to" mapstructure:"to"`
	Cc          []string           `json:"cc,omitempty" yaml:"cc,omitempty" xml:"cc" toml:"cc" mapstructure:"cc"`
	Bcc         []string           `json:"bcc,omitempty" yaml:"bcc,omitempty" xml:"bcc" toml:"bcc" mapstructure:"bcc"`
	Subject     string             `json:"subject,omitempty" yaml:"subject,omitempty" xml:"subject" toml:"subject" mapstructure:"subject"`
	HTML        string             `json:"html,omitempty" yaml:"html,omitempty" xml:"html" toml:"html" mapstructure:"html"`
	Text        string             `json:"text,omitempty" yaml:"text,omitempty" xml:"text" toml:"text" mapstructure:"text"`
	Attachments []types.Attachment `json:"attachments,omitempty" yaml:"attachments,omitempty" xml:"attachments" toml:"attachments" mapstructure:"attachments"`
}

// ToEmail converte para o tipo nativo de tools/mail.
func (r MailRequest) ToEmail() *types.Email {
	return &types.Email{
		Name:        r.Name,
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

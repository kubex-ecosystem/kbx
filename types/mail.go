// Package types defines types and interfaces for sending emails using various SMTP providers.
package types

import "github.com/google/uuid"

type Attachment struct {
	UID      uint32 `json:"attach_id,omitempty" yaml:"attach_id" xml:"attach_id" toml:"attach_id" mapstructure:"attach_id"`
	Filename string `json:"filename,omitempty" yaml:"filename" xml:"filename" toml:"filename" mapstructure:"filename"`
	Data     []byte `json:"data,omitempty" yaml:"data" xml:"data" toml:"data" mapstructure:"data"`
	Size     int    `json:"size,omitempty" yaml:"size" xml:"size" toml:"size" mapstructure:"size"`
	Mime     string `json:"mime,omitempty" yaml:"mime" xml:"mime" toml:"mime" mapstructure:"mime"`
}
type Email struct {
	UID         uint32       `json:"id,omitempty" yaml:"id" xml:"id" toml:"id" mapstructure:"id"`
	Name        string       `json:"from_name,omitempty" yaml:"from_name" xml:"from_name" toml:"from_name" mapstructure:"from_name"`
	From        string       `json:"from_email,omitempty" yaml:"from_email" xml:"from_email" toml:"from_email" mapstructure:"from_email"`
	To          []string     `json:"to,omitempty" yaml:"to" xml:"to" toml:"to" mapstructure:"to"`
	Cc          []string     `json:"cc,omitempty" yaml:"cc" xml:"cc" toml:"cc" mapstructure:"cc"`
	Bcc         []string     `json:"bcc,omitempty" yaml:"bcc" xml:"bcc" toml:"bcc" mapstructure:"bcc"`
	Subject     string       `json:"subject,omitempty" yaml:"subject" xml:"subject" toml:"subject" mapstructure:"subject"`
	Text        string       `json:"text,omitempty" yaml:"text" xml:"text" toml:"text" mapstructure:"text"`
	HTML        string       `json:"html,omitempty" yaml:"html" xml:"html" toml:"html" mapstructure:"html"`
	Attachments []Attachment `json:"attachments,omitempty" yaml:"attachments" xml:"attachments" toml:"attachments" mapstructure:"attachments"`
}

func NewEmail() *Email {
	return &Email{
		UID:         uuid.New().ID(),
		Attachments: []Attachment{},
	}
}

func NewAttachment() *Attachment {
	return &Attachment{
		UID: uuid.New().ID(),
	}
}

type MailProvider interface {
	Send(cfg *MailConnection, msg *Email) error
}

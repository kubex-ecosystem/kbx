// Package types defines types and interfaces for sending emails using various SMTP providers.
package types

import (
	"time"
)

type Email struct {
	Name        string
	From        string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Text        string
	HTML        string
	Attachments []Attachment
}

type Attachment struct {
	Filename string
	Data     []byte
	Mime     string
}

type SMTPConfig struct {
	Provider string // "gmail", "outlook", "microsoft", "sendmail"
	Host     string
	Port     int
	User     string
	Pass     string
	SSL      bool
	TLS      bool
	Timeout  time.Duration
}

type MailProvider interface {
	Send(cfg SMTPConfig, msg *Email) error
}

type MailParams struct {
	*SMTPConfig
	*Email
	Provider MailProvider
}

func NewMailParams() *MailParams {
	return &MailParams{
		&SMTPConfig{},
		&Email{},
		nil,
	}
}

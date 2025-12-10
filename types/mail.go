// Package types defines types and interfaces for sending emails using various SMTP providers.
package types

import (
	"net/smtp"
	"time"
)

type Email struct {
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

type Attachment struct {
	Filename string `json:"filename,omitempty" yaml:"filename" xml:"filename" toml:"filename" mapstructure:"filename"`
	Data     []byte `json:"data,omitempty" yaml:"data" xml:"data" toml:"data" mapstructure:"data"`
	Mime     string `json:"mime,omitempty" yaml:"mime" xml:"mime" toml:"mime" mapstructure:"mime"`
}

type SMTPSender struct {
	Host string    `json:"host" yaml:"host" xml:"host" toml:"host" mapstructure:"host"`
	Port int       `json:"port" yaml:"port" xml:"port" toml:"port" mapstructure:"port"`
	User string    `json:"username" yaml:"username" xml:"username" toml:"username" mapstructure:"username"`
	Pass string    `json:"password" yaml:"password" xml:"password" toml:"password" mapstructure:"password"`
	Auth smtp.Auth `json:"-" yaml:"-" xml:"-" toml:"-" mapstructure:"-"`
}

type SMTPConfig struct {
	*SMTPSender `json:",inline" yaml:",inline" xml:"inline" toml:",inline" mapstructure:"squash"`

	Provider string        `json:"provider,omitempty" yaml:"provider" xml:"provider" toml:"provider" mapstructure:"provider"`
	SSL      bool          `json:"use_ssl,omitempty" yaml:"ssl" xml:"ssl" toml:"ssl" mapstructure:"ssl"`
	TLS      bool          `json:"use_tls,omitempty" yaml:"tls" xml:"tls" toml:"tls" mapstructure:"tls"`
	Timeout  time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty" xml:"timeout,omitempty" toml:"timeout,omitempty" mapstructure:"timeout,omitempty"`
}

type MailProvider interface {
	Send(cfg *SMTPConfig, msg *Email) error
}

type MailParams struct {
	*SMTPConfig `json:"smtp_config,omitempty" yaml:"smtp_config" xml:"smtp_config" toml:"smtp_config" mapstructure:"smtp_config"`
	*Email      `json:"email,omitempty" yaml:"email" xml:"email" toml:"email" mapstructure:"email"`
	Provider    MailProvider `json:"-" yaml:"-" xml:"-" toml:"-" mapstructure:"-"`
}

func NewMailParams() *MailParams {
	return &MailParams{
		&SMTPConfig{
			SMTPSender: &SMTPSender{},
		},
		&Email{},
		nil,
	}
}

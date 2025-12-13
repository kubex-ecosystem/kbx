package types

import (
	"net/smtp"
	"time"
)

type MailGeneralConfig struct {
	MailBox          string        `json:"mailbox,omitempty" yaml:"mailbox" xml:"mailbox" toml:"mailbox" mapstructure:"mailbox"`
	DefaultFilter    string        `json:"filter,omitempty" yaml:"filter" xml:"filter" toml:"filter" mapstructure:"filter"`
	ImportanceCutoff int           `json:"importance_cutoff,omitempty" yaml:"importance_cutoff" xml:"importance_cutoff" toml:"importance_cutoff" mapstructure:"importance_cutoff"`
	MaxEmailsPerRun  int           `json:"max_emails_per_run,omitempty" yaml:"max_emails_per_run" xml:"max_emails_per_run" toml:"max_emails_per_run" mapstructure:"max_emails_per_run"`
	Timezone         string        `json:"timezone,omitempty" yaml:"timezone" xml:"timezone" toml:"timezone" mapstructure:"timezone"`
	RetryCount       int           `json:"retry_count,omitempty" yaml:"retry_count" xml:"retry_count" toml:"retry_count" mapstructure:"retry_count"`
	RetryInterval    time.Duration `json:"retry_interval,omitempty" yaml:"retry_interval" xml:"retry_interval" toml:"retry_interval" mapstructure:"retry_interval"`
	KeepSentCopies   bool          `json:"keep_sent_copies,omitempty" yaml:"keep_sent_copies" xml:"keep_sent_copies" toml:"keep_sent_copies" mapstructure:"keep_sent_copies"`
	Signature        string        `json:"signature,omitempty" yaml:"signature" xml:"signature" toml:"signature" mapstructure:"signature"`
}

type MailAuthConfig struct {
	Provider      string    `json:"provider,omitempty" yaml:"provider" xml:"provider" toml:"provider" mapstructure:"provider"`
	Host          string    `json:"host,omitempty" yaml:"host" xml:"host" toml:"host" mapstructure:"host"`
	Port          int       `json:"port,omitempty" yaml:"port" xml:"port" toml:"port" mapstructure:"port"`
	User          string    `json:"username,omitempty" yaml:"username" xml:"username" toml:"username" mapstructure:"username"`
	Pass          string    `json:"password,omitempty" yaml:"password" xml:"password" toml:"password" mapstructure:"password"`
	Auth          smtp.Auth `json:"-" yaml:"-" xml:"-" toml:"-" mapstructure:"-"`
	From          string    `json:"from_email,omitempty" yaml:"from_email" xml:"from_email" toml:"from_email" mapstructure:"from_email"`
	Name          string    `json:"from_name,omitempty" yaml:"from_name" xml:"from_name" toml:"from_name" mapstructure:"from_name"`
	TestRecipient string    `json:"test_recipient,omitempty" yaml:"test_recipient" xml:"test_recipient" toml:"test_recipient" mapstructure:"test_recipient"`
}
type MailProtocolConfig struct {
	Protocol string        `json:"protocol,omitempty" yaml:"protocol,omitempty" xml:"protocol,omitempty" toml:"protocol,omitempty" mapstructure:"protocol,omitempty"` // "smtp" (default) ou "imap"
	SSL      bool          `json:"use_ssl,omitempty" yaml:"ssl" xml:"ssl" toml:"ssl" mapstructure:"ssl"`
	TLS      bool          `json:"use_tls,omitempty" yaml:"tls" xml:"tls" toml:"tls" mapstructure:"tls"`
	Timeout  time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty" xml:"timeout,omitempty" toml:"timeout,omitempty" mapstructure:"timeout,omitempty"`
}
type MailConnection struct {
	*MailGeneralConfig  `json:",inline,omitempty" yaml:",inline,omitempty" xml:"general,omitempty" toml:",inline,omitempty" mapstructure:"squash,omitempty"`
	*MailAuthConfig     `json:",inline" yaml:",inline" xml:"auth" toml:",inline" mapstructure:"squash"`
	*MailProtocolConfig `json:",inline" yaml:",inline" xml:"protocol" toml:",inline" mapstructure:"squash"`
}
type MailConfig struct {
	ConfigPath  string            `json:"config_path,omitempty" yaml:"config_path,omitempty" xml:"config_path,omitempty" toml:"config_path,omitempty" mapstructure:"config_path,omitempty"`
	Provider    string            `json:"provider,omitempty" yaml:"provider" xml:"provider" toml:"provider" mapstructure:"provider"`
	Connections []*MailConnection `json:"connections,omitempty" yaml:"connections" xml:"connections" toml:"connections" mapstructure:"connections"`
}

func NewMailConfig(provider string) *MailConfig {
	return &MailConfig{
		Provider:    provider,
		Connections: make([]*MailConnection, 0),
	}
}

func NewMailConnection() *MailConnection {
	return &MailConnection{
		MailGeneralConfig:  &MailGeneralConfig{},
		MailAuthConfig:     &MailAuthConfig{},
		MailProtocolConfig: &MailProtocolConfig{},
	}
}

type MailSrvParams struct {
	ConfigPath  string `json:"config_path,omitempty" yaml:"config_path" xml:"config_path" toml:"config_path" mapstructure:"config_path"`
	*MailConfig `json:",inline" yaml:",inline" xml:"-" toml:",inline" mapstructure:",squash"`
	*Email      `json:",inline" yaml:",inline" xml:"-" toml:",inline" mapstructure:",squash"`
}

func NewMailSrvParams(configPath string) *MailSrvParams {
	return &MailSrvParams{
		ConfigPath: configPath,
		MailConfig: NewMailConfig(""),
		Email:      &Email{},
	}
}

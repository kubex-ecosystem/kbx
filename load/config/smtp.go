// Package config defines configuration structures for the application.
package config

import (
	"path/filepath"

	"github.com/kubex-ecosystem/kbx/tools"
)

type SMTP struct {
	Protocol string `json:"protocol,omitempty" yaml:"protocol,omitempty" db:"protocol,omitempty" bson:"protocol,omitempty" form:"protocol,omitempty" mapstructure:"protocol,omitempty"`
	Mail     string `json:"mail,omitempty" yaml:"mail,omitempty" db:"mail,omitempty" bson:"mail,omitempty" form:"mail,omitempty" mapstructure:"mail,omitempty"`
	From     string `json:"from_email,omitempty" yaml:"from_email,omitempty" db:"from_email,omitempty" bson:"from_email,omitempty" form:"from_email,omitempty" mapstructure:"from_email,omitempty"`
	Name     string `json:"from_name,omitempty" yaml:"from_name,omitempty" db:"from_name,omitempty" bson:"from_name,omitempty" form:"from_name,omitempty" mapstructure:"from_name,omitempty"`

	Host string `json:"host,omitempty" yaml:"host,omitempty" db:"host,omitempty" bson:"host,omitempty" form:"host,omitempty" mapstructure:"host,omitempty"`
	Port string `json:"port,omitempty" yaml:"port,omitempty" db:"port,omitempty" bson:"port,omitempty" form:"port,omitempty" mapstructure:"port,omitempty"`
	User string `json:"username,omitempty" yaml:"username,omitempty" db:"username,omitempty" bson:"username,omitempty" form:"username,omitempty" mapstructure:"username,omitempty"`
	Pass string `json:"password,omitempty" yaml:"password,omitempty" db:"password,omitempty" bson:"password,omitempty" form:"password,omitempty" mapstructure:"password,omitempty"`

	SSL  bool   `json:"use_ssl,omitempty" yaml:"use_ssl,omitempty" db:"use_ssl,omitempty" bson:"use_ssl,omitempty" form:"use_ssl,omitempty" mapstructure:"use_ssl,omitempty"`
	TLS  bool   `json:"use_tls,omitempty" yaml:"use_tls,omitempty" db:"use_tls,omitempty" bson:"use_tls,omitempty" form:"use_tls,omitempty" mapstructure:"use_tls,omitempty"`
	Mode string `json:"mode,omitempty" yaml:"mode,omitempty" db:"mode,omitempty" bson:"mode,omitempty" form:"mode,omitempty" mapstructure:"mode,omitempty"`

	Meta map[string]any `json:"meta,omitempty" yaml:"meta,omitempty" db:"meta,omitempty" bson:"meta,omitempty" form:"meta,omitempty" mapstructure:"meta,omitempty"`
}

func BasicSMTP(from string) (*SMTP, error) {
	loader := tools.NewEmptyMapperType[SMTP](from)
	return loader.DeserializeFromFile(filepath.Ext(from)[1:])
}

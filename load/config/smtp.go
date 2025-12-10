// Package config defines configuration structures for the application.
package config

import (
	"path/filepath"

	"github.com/kubex-ecosystem/kbx/tools"
)

type SMTP struct {
	Mail string `json:"mail,omitempty" yaml:"mail,omitempty" db:"mail,omitempty" bson:"mail,omitempty" form:"mail,omitempty" mapstructure:"mail,omitempty"`
	From string `json:"from,omitempty" yaml:"from,omitempty" db:"from,omitempty" bson:"from,omitempty" form:"from,omitempty" mapstructure:"from,omitempty"`
	Name string `json:"name,omitempty" yaml:"name,omitempty" db:"name,omitempty" bson:"name,omitempty" form:"name,omitempty" mapstructure:"name,omitempty"`

	Host string `json:"host,omitempty" yaml:"host,omitempty" db:"host,omitempty" bson:"host,omitempty" form:"host,omitempty" mapstructure:"host,omitempty"`
	Port string `json:"port,omitempty" yaml:"port,omitempty" db:"port,omitempty" bson:"port,omitempty" form:"port,omitempty" mapstructure:"port,omitempty"`
	User string `json:"user,omitempty" yaml:"user,omitempty" db:"user,omitempty" bson:"user,omitempty" form:"user,omitempty" mapstructure:"user,omitempty"`
	Pass string `json:"pass,omitempty" yaml:"pass,omitempty" db:"pass,omitempty" bson:"pass,omitempty" form:"pass,omitempty" mapstructure:"pass,omitempty"`

	SSL  bool   `json:"ssl,omitempty" yaml:"ssl,omitempty" db:"ssl,omitempty" bson:"ssl,omitempty" form:"ssl,omitempty" mapstructure:"ssl,omitempty"`
	TLS  bool   `json:"tls,omitempty" yaml:"tls,omitempty" db:"tls,omitempty" bson:"tls,omitempty" form:"tls,omitempty" mapstructure:"tls,omitempty"`
	Mode string `json:"mode,omitempty" yaml:"mode,omitempty" db:"mode,omitempty" bson:"mode,omitempty" form:"mode,omitempty" mapstructure:"mode,omitempty"`

	Meta map[string]any `json:"meta,omitempty" yaml:"meta,omitempty" db:"meta,omitempty" bson:"meta,omitempty" form:"meta,omitempty" mapstructure:"meta,omitempty"`
}

func BasicSMTP(from string) (*SMTP, error) {
	loader := tools.NewEmptyMapperType[SMTP](from)
	return loader.DeserializeFromFile(filepath.Ext(from)[1:])
}

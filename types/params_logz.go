// Package types provides utilities for working with initialization arguments.
package types

import (
	"github.com/google/uuid"

	gl "github.com/kubex-ecosystem/logz"
)

type DBType string

const (
	DBTypePostgres DBType = "postgres"
	DBTypeRabbitMQ DBType = "rabbitmq"
	DBTypeRedis    DBType = "redis"
	DBTypeMongoDB  DBType = "mongodb"
	DBTypeMySQL    DBType = "mysql"
	DBTypeMSSQL    DBType = "mssql"
	DBTypeSQLite   DBType = "sqlite"
	DBTypeOracle   DBType = "oracle"
)

type LogzConfig struct {
	ID uuid.UUID

	*gl.LogzGeneralOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*gl.LogzFormatOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*gl.LogzOutputOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*gl.LogzRotatingOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*gl.LogzBufferingOptions `json:",inline" yaml:",inline" mapstructure:",squash"`
}

// RootParams representa o arquivo de configuração do DS.
type RootParams struct {
	Name     string `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`
	FilePath string `json:"file_path,omitempty" yaml:"file_path,omitempty" mapstructure:"file_path,omitempty"`
	Enabled  *bool  `json:"enabled,omitempty" yaml:"enabled,omitempty" mapstructure:"enabled,omitempty" default:"true"`
}

func NewLogzConfig() *LogzConfig {
	return &LogzConfig{
		ID:                   uuid.New(),
		LogzGeneralOptions:   &gl.LogzGeneralOptions{},
		LogzFormatOptions:    &gl.LogzFormatOptions{},
		LogzOutputOptions:    &gl.LogzOutputOptions{},
		LogzRotatingOptions:  &gl.LogzRotatingOptions{},
		LogzBufferingOptions: &gl.LogzBufferingOptions{},
	}
}

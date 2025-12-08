// Package kbx provides utilities for working with initialization arguments.
package kbx

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

type Params struct {
	ID uuid.UUID

	*gl.LogzGeneralOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*gl.LogzFormatOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*gl.LogzOutputOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*gl.LogzRotatingOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*gl.LogzBufferingOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	// *LogzAdvancedOptions `json:",inline" yaml:",inline" mapstructure:",squash"`
}

// RootConfig representa o arquivo de configuração do DS.
type RootConfig struct {
	Name     string `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`
	FilePath string `json:"file_path,omitempty" yaml:"file_path,omitempty" mapstructure:"file_path,omitempty"`
	Enabled  *bool  `json:"enabled,omitempty" yaml:"enabled,omitempty" mapstructure:"enabled,omitempty" default:"true"`
}

var LoggerArgs *Params = &Params{
	ID:                   uuid.New(),
	LogzGeneralOptions:   &gl.LogzGeneralOptions{},
	LogzFormatOptions:    &gl.LogzFormatOptions{},
	LogzOutputOptions:    &gl.LogzOutputOptions{},
	LogzRotatingOptions:  &gl.LogzRotatingOptions{},
	LogzBufferingOptions: &gl.LogzBufferingOptions{},
}

func ParseLoggerArgs(level string, minLevel string, maxLevel string, output string) *Params {
	LoggerArgs.Level = gl.Level(GetValueOrDefaultSimple(level, "info"))
	LoggerArgs.MinLevel = gl.Level(GetValueOrDefaultSimple(minLevel, "debug"))
	LoggerArgs.MaxLevel = gl.Level(GetValueOrDefaultSimple(maxLevel, "fatal"))
	// LoggerArgs.Output = GetValueOrDefaultSimple(gl.NewLogzWriter(output, os.Stdout), gl.NewLogzIOWriter(os.Stdout))
	return LoggerArgs
}

func init() {
	if LoggerArgs == nil {
		LoggerArgs = &Params{
			ID:                   uuid.New(),
			LogzGeneralOptions:   &gl.LogzGeneralOptions{},
			LogzFormatOptions:    &gl.LogzFormatOptions{},
			LogzOutputOptions:    &gl.LogzOutputOptions{},
			LogzRotatingOptions:  &gl.LogzRotatingOptions{},
			LogzBufferingOptions: &gl.LogzBufferingOptions{},
		}
	}
}

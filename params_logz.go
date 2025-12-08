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

type LogzParams struct {
	ID uuid.UUID

	*gl.LogzGeneralOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*gl.LogzFormatOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*gl.LogzOutputOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*gl.LogzRotatingOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*gl.LogzBufferingOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	// *LogzAdvancedOptions `json:",inline" yaml:",inline" mapstructure:",squash"`
}

// RootParams representa o arquivo de configuração do DS.
type RootParams struct {
	Name     string `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`
	FilePath string `json:"file_path,omitempty" yaml:"file_path,omitempty" mapstructure:"file_path,omitempty"`
	Enabled  *bool  `json:"enabled,omitempty" yaml:"enabled,omitempty" mapstructure:"enabled,omitempty" default:"true"`
}

var LogzArgs *LogzParams = &LogzParams{
	ID:                   uuid.New(),
	LogzGeneralOptions:   &gl.LogzGeneralOptions{},
	LogzFormatOptions:    &gl.LogzFormatOptions{},
	LogzOutputOptions:    &gl.LogzOutputOptions{},
	LogzRotatingOptions:  &gl.LogzRotatingOptions{},
	LogzBufferingOptions: &gl.LogzBufferingOptions{},
}

func ParseLogzArgs(level string, minLevel string, maxLevel string, output string) *LogzParams {
	LogzArgs.Level = gl.Level(GetValueOrDefaultSimple(level, "info"))
	LogzArgs.MinLevel = gl.Level(GetValueOrDefaultSimple(minLevel, "debug"))
	LogzArgs.MaxLevel = gl.Level(GetValueOrDefaultSimple(maxLevel, "fatal"))
	// LogzArgs.Output = GetValueOrDefaultSimple(gl.NewLogzWriter(output, os.Stdout), gl.NewLogzIOWriter(os.Stdout))
	return LogzArgs
}

func init() {
	if LogzArgs == nil {
		LogzArgs = &LogzParams{
			ID:                   uuid.New(),
			LogzGeneralOptions:   &gl.LogzGeneralOptions{},
			LogzFormatOptions:    &gl.LogzFormatOptions{},
			LogzOutputOptions:    &gl.LogzOutputOptions{},
			LogzRotatingOptions:  &gl.LogzRotatingOptions{},
			LogzBufferingOptions: &gl.LogzBufferingOptions{},
		}
	}
}

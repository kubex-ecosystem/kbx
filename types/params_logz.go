// Package types provides utilities for working with initialization arguments.
package types

import "github.com/google/uuid"

// LogLevel is a string alias for log level values.
type LogLevel = string

// LogzConfig represents the configuration for the logger.
type LogzConfig struct {
	ID       uuid.UUID `json:"id,omitempty" yaml:"id,omitempty" mapstructure:"id,omitempty"`
	Level    LogLevel  `json:"level,omitempty" yaml:"level,omitempty" mapstructure:"level,omitempty"`
	MinLevel LogLevel  `json:"min_level,omitempty" yaml:"min_level,omitempty" mapstructure:"min_level,omitempty"`
	MaxLevel LogLevel  `json:"max_level,omitempty" yaml:"max_level,omitempty" mapstructure:"max_level,omitempty"`
	Debug    bool      `json:"debug,omitempty" yaml:"debug,omitempty" mapstructure:"debug,omitempty"`
	Format   string    `json:"format,omitempty" yaml:"format,omitempty" mapstructure:"format,omitempty"`
	Output   string    `json:"output,omitempty" yaml:"output,omitempty" mapstructure:"output,omitempty"`
}

// NewLogzConfig creates a new LogzConfig.
func NewLogzConfig() *LogzConfig {
	return &LogzConfig{
		ID: uuid.New(),
	}
}

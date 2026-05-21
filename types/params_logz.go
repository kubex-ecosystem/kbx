// Package types provides utilities for working with initialization arguments.
package types

import (
	"github.com/google/uuid"

	"github.com/kubex-ecosystem/logz"
	gl "github.com/kubex-ecosystem/logz"
)

// LogzConfig represents the configuration for the logger.
type LogzConfig struct {
	ID uuid.UUID

	*gl.LogzGeneralOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*gl.LogzFormatOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*gl.LogzOutputOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*gl.LogzRotatingOptions `json:",inline" yaml:",inline" mapstructure:",squash"`

	*gl.LogzBufferingOptions `json:",inline" yaml:",inline" mapstructure:",squash"`
}

// NewLogzConfig creates a new LogzConfig.
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

// No Go, quando um struct contém campos que são ponteiros para outros structs,
// os campos dos structs apontados são promovidos para o struct pai.
// Isso significa que os campos podem ser acessados diretamente como se fossem
// campos do struct pai, mas eles ainda são ponteiros.
//
// Para mais informações, veja:
//   - https://pkg.go.dev/fmt#Sprintf
//   - https://pkg.go.dev/fmt# %#v
func ExemploPonteiroPromovido(l *LogzConfig) {
	if !l.Debug || !l.LogzGeneralOptions.Debug {
		return
	}
	logz.Successf("O objeto aparenta ser um 'mutante' em função dos ponteiros.. (%v) rsrs", l)
}

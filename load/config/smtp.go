// Package config defines configuration structures for the application.
package config

import (
	"path/filepath"

	"github.com/kubex-ecosystem/kbx/load"
	"github.com/kubex-ecosystem/kbx/types"

	gl "github.com/kubex-ecosystem/logz"
)

type SMTP = types.MailConnection

func BasicSMTP(from string) (*SMTP, error) {
	mailerConfig, err := load.LoadConfigOrDefault[types.MailConfig](from, true)
	if err != nil {
		return nil, err
	}
	for _, smtp := range mailerConfig.Connections {
		if smtp.Protocol == "smtp" || smtp.Protocol == "" {
			return smtp, nil
		}
	}
	return nil, gl.Errorf("no SMTP configuration found in %s", filepath.Base(from))
}

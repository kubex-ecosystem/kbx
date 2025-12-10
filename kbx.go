// Package kbx provides utilities for working with initialization arguments.
package kbx

import (
	"os"
	"path/filepath"
	"sync"
)

const (
	// defaultEnvVarPrefix = "KBX_"
	// defaultManifestPath   = "internal/module/info/manifest.json"

	defaultSMTPConfigPath = ".kubex/mainling/config/smtp.json"
	defaultTemplatePath   = ".kubex/mainling/email_templates"
	defaultEnvFilePath    = ".kubex/.env"
)

var (
	kubexEcosystemCwd string
	once              sync.Once
)

func getKubexEcosystemCwd() string {
	once.Do(func() {
		kubexEcosystemCwd = os.ExpandEnv(
			GetValErrOrDefault(
				os.UserHomeDir,
				GetValErrOrDefault(
					os.Getwd,
					GetEnvOrDefault(
						"KUBEX_ECOSYSTEM_CWD",
						".",
					),
				),
			),
		)
	})
	return kubexEcosystemCwd
}

func getFullExpandedPath(path string) string {
	return os.ExpandEnv(filepath.Join(getKubexEcosystemCwd(), path))
}

func DefaultSMTPConfigPath() string { return getFullExpandedPath(defaultSMTPConfigPath) }
func DefaultTemplatePath() string   { return getFullExpandedPath(defaultTemplatePath) }
func DefaultEnvFilePath() string    { return getFullExpandedPath(defaultEnvFilePath) }

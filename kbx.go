// Package kbx provides utilities for working with initialization arguments.
package kbx

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/kubex-ecosystem/kbx/get"
	"github.com/kubex-ecosystem/kbx/load"
	"github.com/kubex-ecosystem/kbx/types"
)

const (
	// defaultEnvVarPrefix = "KBX_"
	// defaultManifestPath   = "internal/module/info/manifest.json"

	// defaultRootPath is the default directory name under which
	// Kubex Ecosystem will store its configuration and data files
	// within the user's home directory or specified directory.
	defaultRootPath = ".kubex"

	// ------------------------------- Default Paths -----------------------------------//

	defaultSMTPConfigPath = "mailing/config/smtp.json"
	defaultTemplatePath   = "mailing/email_templates"
	defaultEnvFilePath    = ".env"
)

var (
	kubexEcosystemCwd string
	once              sync.Once
)

// getKubexEcosystemCwd returns where Kubex Ecosystem will store
// its configuration and data files. It first checks the KUBEX_ECOSYSTEM_CWD
// environment variable. If not set, it falls back to the user's home directory.
// If the home directory cannot be determined, it uses the current working directory.
func getKubexEcosystemCwd() string {
	once.Do(func() {
		kubexEcosystemCwd = os.ExpandEnv(
			get.ValErrOr(
				os.UserHomeDir,
				get.ValErrOr(
					os.Getwd,
					get.EnvOr("KUBEX_ECOSYSTEM_CWD", "."),
				),
			),
		)
	})
	return kubexEcosystemCwd
}

// getFullExpandedPath constructs a full path by joining the root path,
// current working directory, and the provided path, then expands any environment
// variables in the resulting path. Root path will be determined at runtime.
func getFullExpandedPath(path string) string {
	return os.ExpandEnv(
		filepath.Join(
			get.EnvOr("KUBEX_ECOSYSTEM_ROOT", defaultRootPath),
			getKubexEcosystemCwd(),
			path,
		),
	)
}

// ------------------------------- Default Paths Functions -----------------------------//

func DefaultSMTPConfigPath() string { return getFullExpandedPath(defaultSMTPConfigPath) }
func DefaultTemplatePath() string   { return getFullExpandedPath(defaultTemplatePath) }
func DefaultEnvFilePath() string    { return getFullExpandedPath(defaultEnvFilePath) }

// ------------------------------- New Mail Params Functions -----------------------------//

type MailSrvParams = load.MailSrvParams
type MailConfig = load.MailConfig
type MailConnection = types.MailConnection
type MailAttachment = types.Attachment
type Email = types.Email
type MManifest = types.MManifest
type Manifest = load.Manifest

type LogzConfig = types.LogzConfig
type SrvConfig = types.SrvConfig
type VendorAuthConfig = load.VendorAuthConfig
type AuthOAuthClientConfig = load.AuthOAuthClientConfig
type AuthClientConfig = load.AuthClientConfig
type AuthProvidersConfig = load.AuthProvidersConfig
type GlobalRef = load.GlobalRef

func NewMailSrvParams(cfgPath string) *MailSrvParams { return load.NewMailSrvParams(cfgPath) }
func NewMailConfig(cfgPath string) *MailConfig       { return load.NewMailConfig(cfgPath) }
func NewMailConnection() *MailConnection             { return types.NewMailConnection() }
func NewMailAttachment() *MailAttachment             { return &MailAttachment{} }
func NewEmail() *Email                               { return &Email{} }
func NewManifestType() *MManifest                    { return load.NewManifestType() }
func NewManifest() Manifest                          { return load.NewManifest() }

// func NewMailSender(params *MailSrvParams) MailSender { return nil }

func NewLogzParams() *types.LogzConfig   { return load.NewLogzParams() }
func NewSrvArgs() types.SrvConfig        { return load.NewSrvArgs() }
func NewGlobalRef(name string) GlobalRef { return load.NewGlobalRef(name) }

func ParseLogzArgs(level string, minLevel string, maxLevel string, output string) *types.LogzConfig {
	return load.ParseLogzArgs(level, minLevel, maxLevel, output)
}
func ParseSrvArgs(bind, port, pubCertKeyPath, pubKeyPath, privKeyPath string, accessTokenTTL int, refreshTokenTTL int, issuer string) types.SrvConfig {
	return load.ParseSrvArgs(bind, port, pubCertKeyPath, pubKeyPath, privKeyPath, accessTokenTTL, refreshTokenTTL, issuer)
}

func LoadConfig[T any](path string) (T, error) { return load.LoadConfig[T](path) }
func LoadConfigOrDefault[T MailConfig | MailConnection | LogzConfig | SrvConfig | MailSrvParams | Email | MManifest | VendorAuthConfig | AuthOAuthClientConfig](cfgPath string, genFile bool) (*T, error) {
	return load.LoadConfigOrDefault[T](cfgPath, genFile)
}

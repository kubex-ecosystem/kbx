// Package kbx provides utilities for working with initialization arguments.
package kbx

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kubex-ecosystem/kbx/get"
	"github.com/kubex-ecosystem/kbx/tools"
	"github.com/kubex-ecosystem/kbx/tools/mail/provider"
	"github.com/kubex-ecosystem/kbx/types"

	gl "github.com/kubex-ecosystem/logz"
)

const (
	// defaultEnvVarPrefix = "KBX_"
	// defaultManifestPath   = "internal/module/info/manifest.json"

	// defaultRootPath is the default directory name under which
	// Kubex Ecosystem will store its configuration and data files
	// within the user's home directory or specified directory.
	defaultRootPath = ".kubex"

	// ------------------------------- Default Paths -----------------------------------//

	defaultSMTPConfigPath = "mainling/config/smtp.json"
	defaultTemplatePath   = "mainling/email_templates"
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

type MailParams struct {
	ConfigPath        string `json:"config_path,omitempty"`
	*types.MailParams `json:",inline" mapstructure:",squash"`
}

func NewMailParams(configPath string) *MailParams {
	return &MailParams{ConfigPath: configPath, MailParams: types.NewMailParams()}
}

func (mp *MailParams) LoadSMTPConfig(cfgFilePath string) (*types.SMTPConfig, error) {
	cfgLoader := get.Loader[types.SMTPConfig](get.ValOrType(cfgFilePath, mp.ConfigPath))
	return cfgLoader.DeserializeFromFile(get.FileExt(get.ValOrType(cfgFilePath, mp.ConfigPath)))
}

// ------------------------------- New Logz Params Functions -----------------------------//

func NewLogzParams() *types.LogzParams { return types.NewLogzParams() }

func ParseLogzArgs(level string, minLevel string, maxLevel string, output string) *types.LogzParams {
	LogzArgs := NewLogzParams()
	LogzArgs.Level = gl.Level(get.ValOrType(level, "info"))
	LogzArgs.MinLevel = gl.Level(get.ValOrType(minLevel, "debug"))
	LogzArgs.MaxLevel = gl.Level(get.ValOrType(maxLevel, "fatal"))
	return LogzArgs
}
func LoadLogzConfig(cfgPath string) (*types.LogzParams, error) {
	cfgMapper := tools.NewEmptyMapperType[types.LogzParams](cfgPath)
	return cfgMapper.DeserializeFromFile(filepath.Ext(cfgPath)[1:])
}

// ------------------------------- New Srv Params Functions -----------------------------//

func NewSrvArgs() *types.SrvParams { return types.NewSrvParams() }

func ParseSrvArgs(bind string, pubCertKeyPath string, pubKeyPath string, privKeyPath string, accessTokenTTL int, refreshTokenTTL int, issuer string) *types.SrvParams {
	SrvArgs := NewSrvArgs()
	SrvArgs.Bind = get.ValOrType(bind, ":8080")
	SrvArgs.PubCertKeyPath = get.ValOrType(pubCertKeyPath, "")
	SrvArgs.PubKeyPath = get.ValOrType(pubKeyPath, "")
	SrvArgs.PrivKeyPath = get.ValOrType(privKeyPath, "")
	SrvArgs.AccessTokenTTL = time.Duration(get.ValOrType(accessTokenTTL, 15)) * time.Minute
	SrvArgs.RefreshTokenTTL = time.Duration(get.ValOrType(refreshTokenTTL, 60)) * time.Minute
	SrvArgs.Issuer = get.ValOrType(issuer, "kubex-ecosystem")
	return SrvArgs
}

type MailSvc = types.MailProvider

func NewMailSvc(cfgPath string) MailSvc {
	msvc, err := provider.NewProvider[MailSvc](cfgPath)
	if err != nil {
		return nil
	}
	return msvc
}

type GlobalRef struct {
	ID   uuid.UUID `json:"id,omitempty"`
	Name string    `json:"name,omitempty"`
}

func NewGlobalRef(name string) *GlobalRef {
	return &GlobalRef{
		ID:   uuid.New(),
		Name: name,
	}
}

func (gr *GlobalRef) GetGlobalRef() *GlobalRef { return gr }
func (gr *GlobalRef) GetName() string          { return gr.Name }
func (gr *GlobalRef) GetID() uuid.UUID         { return gr.ID }
func (gr *GlobalRef) SetName(name string)      { gr.Name = name }
func (gr *GlobalRef) SetID(id uuid.UUID)       { gr.ID = id }
func (gr *GlobalRef) String() string {
	return gr.Name + "-" + gr.ID.String()
}

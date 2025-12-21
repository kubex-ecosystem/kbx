// Package load provides functions to load configuration and environment settings.
package load

import (
	"os"
	"reflect"
	"time"

	"github.com/kubex-ecosystem/kbx/get"
	"github.com/kubex-ecosystem/kbx/is"
	"github.com/kubex-ecosystem/kbx/tools"
	"github.com/kubex-ecosystem/kbx/types"

	gl "github.com/kubex-ecosystem/logz"
)

type MailConfig = types.MailConfig
type MailConnection = types.MailConnection
type Email = types.Email

// ------------------------------- New Mail Srv Params Functions -----------------------------//

type MailSrvParams struct {
	ConfigPath         string `json:"config_path,omitempty" yaml:"config_path,omitempty" xml:"config_path,omitempty" toml:"config_path,omitempty" mapstructure:"config_path,omitempty"`
	types.Attachment   `json:",inline" yaml:",inline" xml:"-" toml:",inline" mapstructure:",squash"`
	types.Email        `json:",inline" yaml:",inline" xml:"-" toml:",inline" mapstructure:",squash"`
	types.MailConfig   `json:",inline" yaml:",inline" xml:"-" toml:",inline" mapstructure:",squash"`
	types.MailProvider `json:"-" yaml:"-" xml:"-" toml:"-" mapstructure:"-"`
}

func NewMailSrvParams(configPath string) *MailSrvParams {
	mailCfg := types.NewMailConfig(configPath)
	return &MailSrvParams{ConfigPath: configPath, MailConfig: mailCfg, Attachment: types.Attachment{}, Email: types.Email{}}
}

// ------------------------------- New Mail Params Functions -----------------------------//

func NewMailConfig(configPath string) *MailConfig {
	return &MailConfig{
		ConfigPath:  configPath,
		Provider:    "",
		Connections: make([]MailConnection, 0),
	}
}

// ------------------------------- New Logz Params Functions -----------------------------//

type LogzConfig = types.LogzConfig

func NewLogzParams() *LogzConfig { return &LogzConfig{} }

func ParseLogzArgs(level string, minLevel string, maxLevel string, output string) *LogzConfig {
	LogzArgs := NewLogzParams()
	LogzArgs.Level = gl.Level(get.ValOrType(level, "info"))
	LogzArgs.MinLevel = gl.Level(get.ValOrType(minLevel, "info"))
	LogzArgs.MaxLevel = gl.Level(get.ValOrType(maxLevel, "fatal"))
	return LogzArgs
}

// ------------------------------- New Srv Params Functions -----------------------------//

type SrvConfig = types.SrvConfig

func NewSrvArgs() SrvConfig { return types.NewSrvConfig() }

func ParseSrvArgs(bind string, port string, pubCertKeyPath string, pubKeyPath string, privKeyPath string, accessTokenTTL int, refreshTokenTTL int, issuer string) SrvConfig {
	SrvArgs := NewSrvArgs()
	SrvArgs.Runtime.Bind = get.ValOrType(bind, "0.0.0.0")
	SrvArgs.Runtime.Port = get.ValOrType(port, "4000")
	SrvArgs.Runtime.PubCertKeyPath = get.ValOrType(pubCertKeyPath, "")
	SrvArgs.Runtime.PubKeyPath = get.ValOrType(pubKeyPath, "")
	SrvArgs.Runtime.PrivKeyPath = get.ValOrType(privKeyPath, "")
	SrvArgs.Runtime.AccessTokenTTL = time.Duration(get.ValOrType(accessTokenTTL, 15)) * time.Minute
	SrvArgs.Runtime.RefreshTokenTTL = time.Duration(get.ValOrType(refreshTokenTTL, 60)) * time.Minute
	SrvArgs.Runtime.Issuer = get.ValOrType(issuer, "kubex-ecosystem")
	return SrvArgs
}

func NewSrvDefaultConfig(defaults map[string]any) SrvConfig {
	baseURL := get.ValueOrIf((get.EnvOr("CANALIZE_ENV", "development") == "production"),
		"https://api.canalize.app",
		"http://localhost:4000",
	)
	defaultTTL := get.EnvOrType("INVITE_EXPIRATION", 7*24*time.Hour)
	configPath := os.ExpandEnv(get.EnvOr("CANALIZE_BE_CONFIG_PATH", "$HOME/mvp/canalize_be_latest/.be.config.json"))
	pubKeyPath := os.ExpandEnv(get.EnvOrType("CANALIZE_BE_PUBLIC_KEY_PATH", defaults["default_canalyze_be_cert_path"].(string)))
	privKeyPath := os.ExpandEnv(get.EnvOrType("CANALIZE_BE_PRIVATE_KEY_PATH", defaults["default_canalyze_be_key_path"].(string)))

	Cfg := types.NewSrvConfig()
	Cfg.Files.ConfigFile = os.ExpandEnv(configPath)
	Cfg.Files.DBConfigFile = os.ExpandEnv(get.EnvOr("CANALIZE_DS_CONFIG_PATH", "$HOME/.canalize/canalize_ds/config/config.json"))
	Cfg.Files.EnvFile = os.ExpandEnv(get.EnvOr("CANALIZE_BE_ENV_PATH", "$HOME/mvp/canalize_be_latest/.be.env"))
	Cfg.Files.LogFile = os.ExpandEnv(get.EnvOr("CANALIZE_BE_LOG_FILE_PATH", "$HOME/mvp/canalize_be_latest/canalize_be.log"))
	Cfg.GlobalRef = types.NewGlobalRef(get.EnvOr("CANALIZE_BE_PROCESS_NAME", "canalize_be"))
	Cfg.Basic.Debug = get.EnvOrType("CANALIZE_BE_DEBUG_MODE", false)
	Cfg.Basic.ReleaseMode = get.EnvOrType("CANALIZE_BE_RELEASE_MODE", false)
	Cfg.Basic.IsConfidential = get.EnvOrType("CANALIZE_BE_CONFIDENCIAL_MODE", false)
	Cfg.Runtime.Port = get.EnvOrType("CANALIZE_BE_PORT", "4000")
	Cfg.Runtime.Host = baseURL
	Cfg.Runtime.PrivKeyPath = privKeyPath
	Cfg.Runtime.PubKeyPath = pubKeyPath
	Cfg.Runtime.PubCertKeyPath = pubKeyPath
	Cfg.Basic.CORSEnabled = get.EnvOrType("CANALIZE_BE_ENABLE_CORS", true)
	Cfg.Basic.Debug = get.EnvOrType("CANALIZE_BE_DEBUG_MODE", false)
	Cfg.Files.ProvidersConfig = os.ExpandEnv(get.EnvOr("CANALIZE_BE_PROVIDERS_CONFIG_PATH", ""))
	Cfg.Runtime.RefreshTokenTTL = defaultTTL

	return Cfg
}

func NewSrvConfigFromParams(params *SrvConfig) SrvConfig {
	Cfg := types.NewSrvConfig()
	Cfg.Files.ConfigFile = get.ValOrType(params.Files.ConfigFile, Cfg.Files.ConfigFile)
	Cfg.Files.DBConfigFile = get.ValOrType(params.Files.DBConfigFile, Cfg.Files.DBConfigFile)
	Cfg.Files.EnvFile = get.ValOrType(params.Files.EnvFile, Cfg.Files.EnvFile)
	Cfg.Files.LogFile = get.ValOrType(params.Files.LogFile, Cfg.Files.LogFile)
	Cfg.GlobalRef = get.ValOrType(params.GlobalRef, Cfg.GlobalRef)
	Cfg.Basic.Debug = get.ValOrType(params.Basic.Debug, Cfg.Basic.Debug)
	Cfg.Basic.ReleaseMode = get.ValOrType(params.Basic.ReleaseMode, Cfg.Basic.ReleaseMode)
	Cfg.Basic.IsConfidential = get.ValOrType(params.Basic.IsConfidential, Cfg.Basic.IsConfidential)
	Cfg.Runtime.Port = get.ValOrType(params.Runtime.Port, Cfg.Runtime.Port)
	Cfg.Runtime.Host = get.ValOrType(params.Runtime.Host, Cfg.Runtime.Host)
	Cfg.Runtime.PrivKeyPath = get.ValOrType(params.Runtime.PrivKeyPath, Cfg.Runtime.PrivKeyPath)
	Cfg.Runtime.PubKeyPath = get.ValOrType(params.Runtime.PubKeyPath, Cfg.Runtime.PubKeyPath)
	Cfg.Runtime.PubCertKeyPath = get.ValOrType(params.Runtime.PubCertKeyPath, Cfg.Runtime.PubCertKeyPath)
	Cfg.Basic.CORSEnabled = params.Basic.CORSEnabled
	Cfg.Files.ProvidersConfig = get.ValOrType(params.Files.ProvidersConfig, Cfg.Files.ProvidersConfig)
	Cfg.Runtime.RefreshTokenTTL = get.ValOrType(params.Runtime.RefreshTokenTTL, Cfg.Runtime.RefreshTokenTTL)
	return Cfg
}

type GlobalRef = types.GlobalRef

func NewGlobalRef(name string) GlobalRef { return types.NewGlobalRef(name) }

// ------------------------------- Google Auth Config Functions -----------------------------//

type VendorAuthConfig struct {
	AuthProvider string `json:"auth_provider,omitempty" yaml:"auth_provider,omitempty" xml:"auth_provider,omitempty" toml:"auth_provider,omitempty" mapstructure:"auth_provider,omitempty"`
	types.AuthClientConfig
	ConfigPath string `json:"config_path,omitempty" yaml:"config_path,omitempty" xml:"config_path,omitempty" toml:"config_path,omitempty" mapstructure:"config_path,omitempty"`
}

func NewVendorAuthConfig(cfgPath string) VendorAuthConfig {
	return VendorAuthConfig{
		AuthClientConfig: types.AuthClientConfig{
			AuthProvider: "google",
			// Web default config
			Web: types.AuthOAuthClientConfig{
				ClientID:                "",
				ClientSecret:            "",
				RedirectURL:             "",
				AuthURI:                 "",
				TokenURI:                "",
				AuthProviderX509CertURL: "",
				Scopes:                  []string{"openid", "email", "profile"},
				RedirectURIs:            make([]string, 0),
				JavaScriptOrigins:       make([]string, 0),
				MapUserInfo:             false,
				MetadataOnly:            false,
				Metadata:                make(map[string]any),
			},
			Options: make(map[string]any),
		},
		ConfigPath: cfgPath,
	}
}

// ------------------------------- New Manifest Functions -----------------------------//

type Manifest = types.Manifest
type MManifest = types.MManifest

func NewManifestType() *MManifest {
	bin, _ := os.Executable()

	return &MManifest{
		Version:      "1.0.0",
		Name:         "kubex-manifest",
		Description:  "Kubex Ecosystem Manifest File",
		GoVersion:    "1.25.5",
		Private:      true,
		Author:       "Rafael Mori",
		License:      "MIT",
		Published:    false,
		Aliases:      []string{"kbx-manifest"},
		Homepage:     "https://kubex.world",
		Repository:   "github.com/kubex-ecosystem/kbx",
		Keywords:     []string{"kubex", "kbx", "manifest", "configuration", "ecosystem"},
		Bin:          bin,
		Organization: "Kubex Ecosystem",
		Application:  "kbx",
		Main:         "cmd",
		Platforms: []string{
			"linux/amd64",
			"linux/arm64",
			"darwin/amd64",
			"darwin/arm64",
			"windows/amd64",
		},
		Dependencies: []string{
			"tar",
			"gzip",
			"curl",
			"git",
			"zip",
			"unzip",
			"jq",
			"upx",
		},
	}
}

func NewManifest() Manifest {
	return NewManifestType()
}

func EnsureGlobalManifest(n, c *MManifest) {
	if n == nil && c == nil {
		gl.Fatal("No manifest available")
	}
	if c == nil {
		c = n
	} else if n != nil && n.GetVersion() != c.GetVersion() {
		// Merge new manifest into existing one
		*c = *n
	}
	types.KubexManifest = c
}

// ------------------------------- KBX Config Registry -----------------------------//

var configRegistry = map[reflect.Type]bool{
	reflect.TypeFor[MailSrvParams]():    true,
	reflect.TypeFor[MailConfig]():       true,
	reflect.TypeFor[LogzConfig]():       true,
	reflect.TypeFor[SrvConfig]():        true,
	reflect.TypeFor[MManifest]():        true,
	reflect.TypeFor[VendorAuthConfig](): true,
	reflect.TypeFor[Email]():            true,
	reflect.TypeFor[MailConnection]():   true,
}

var defaultFactories = map[reflect.Type]func() any{
	reflect.TypeFor[MailSrvParams]():    func() any { return NewMailSrvParams("") },
	reflect.TypeFor[MailConfig]():       func() any { return NewMailConfig("") },
	reflect.TypeFor[LogzConfig]():       func() any { return NewLogzParams() },
	reflect.TypeFor[SrvConfig]():        func() any { return NewSrvArgs() },
	reflect.TypeFor[MManifest]():        func() any { return NewManifestType() },
	reflect.TypeFor[VendorAuthConfig](): func() any { return NewVendorAuthConfig("") },
	reflect.TypeFor[Email]():            func() any { return types.NewEmail() },
	reflect.TypeFor[MailConnection]():   func() any { return types.NewMailConnection() },
}

// LoadConfig loads a configuration of type T from the specified file path.

func LoadConfig[T any](cfgPath string) (T, error) {
	var zero T
	var okob bool
	if configRegistry[reflect.TypeFor[T]()] {
		cfgLoader := get.Loader[T](cfgPath)
		obj, err := cfgLoader.DeserializeFromFile(get.FileExt(cfgPath))
		if err != nil && !os.IsNotExist(err) {
			return zero, err
		} else if os.IsNotExist(err) {
			gl.Warnf("configuration file '%s' does not exist", cfgPath)
			return zero, nil
		}
		if reflect.TypeFor[T]() == reflect.TypeFor[MManifest]() {
			var b *MManifest
			o := *obj
			b, okob = any(o).(*MManifest)
			if !okob {
				return zero, gl.Errorf("loaded object is not of type MManifest")
			}
			EnsureGlobalManifest(b, types.KubexManifest)
		}
		return *obj, nil
	}
	return zero, gl.Errorf("configuration type not registered")
}

func LoadConfigOrDefault[T MailConfig | MailConnection | LogzConfig | SrvConfig | MailSrvParams | Email | MManifest | VendorAuthConfig](cfgPath string, genFile bool) (*T, error) {
	if cfgPath == "" {
		gl.Fatalf("config path is empty")
	}

	// Só entra aqui se o tipo for algum já registrado, então não me preocupo em checar o erro, só logo retorno o default
	cfgMapper := tools.NewEmptyMapperType[T](cfgPath)
	cfg, err := cfgMapper.DeserializeFromFile(get.FileExt(cfgPath))
	if err == nil {
		return cfg, nil
	}
	gl.Warnf("failed to load config from '%s', using default: %v", cfgPath, err)
	defaultCfg := defaultFactories[reflect.TypeFor[T]()]().(T)
	if !is.PtrOf[T](defaultCfg) {
		if genFile {
			cfgMapper.SetValue(&defaultCfg)
			cfgMapper.SerializeToFile(get.FileExt(cfgPath))
		}
		return &defaultCfg, nil
	} else {
		d := any(defaultCfg).(*T)
		if genFile {
			cfgMapper.SetValue(d)
			cfgMapper.SerializeToFile(get.FileExt(cfgPath))
		}
		return d, nil
	}
}

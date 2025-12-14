// Package load provides functions to load configuration and environment settings.
package load

import (
	"os"
	"reflect"
	"time"

	"github.com/kubex-ecosystem/kbx/get"
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
	*types.Attachment  `json:",inline" yaml:",inline" xml:"-" toml:",inline" mapstructure:",squash"`
	*types.Email       `json:",inline" yaml:",inline" xml:"-" toml:",inline" mapstructure:",squash"`
	*types.MailConfig  `json:",inline" yaml:",inline" xml:"-" toml:",inline" mapstructure:",squash"`
	types.MailProvider `json:"-" yaml:"-" xml:"-" toml:"-" mapstructure:"-"`
}

func NewMailSrvParams(configPath string) *MailSrvParams {
	return &MailSrvParams{ConfigPath: configPath, MailConfig: types.NewMailConfig(configPath), Attachment: &types.Attachment{}, Email: &types.Email{}}
}

// ------------------------------- New Mail Params Functions -----------------------------//

func NewMailConfig(configPath string) *MailConfig {
	return &MailConfig{
		ConfigPath:  configPath,
		Provider:    "",
		Connections: make([]*MailConnection, 0),
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

func NewSrvArgs() *SrvConfig { return &SrvConfig{} }

func ParseSrvArgs(bind string, pubCertKeyPath string, pubKeyPath string, privKeyPath string, accessTokenTTL int, refreshTokenTTL int, issuer string) *SrvConfig {
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

func NewSrvDefaultConfig(defaults map[string]any) *SrvConfig {
	baseURL := get.ValueOrIf((get.EnvOr("CANALIZE_ENV", "development") == "production"),
		"https://api.canalize.app",
		"http://localhost:4000",
	)
	defaultTTL := get.EnvOrType("INVITE_EXPIRATION", 7*24*time.Hour)
	configPath := os.ExpandEnv(get.EnvOr("CANALIZE_BE_CONFIG_PATH", "/ALL/CANALIZE/projects/BACKEND/canalize_be/configs/config.json"))
	pubKeyPath := os.ExpandEnv(get.EnvOrType[string]("CANALIZE_BE_PUBLIC_KEY_PATH", defaults["default_canalyze_be_cert_path"].(string)))
	privKeyPath := os.ExpandEnv(get.EnvOrType[string]("CANALIZE_BE_PRIVATE_KEY_PATH", defaults["default_canalyze_be_key_path"].(string)))

	Cfg := types.NewSrvConfig()
	Cfg.ConfigFile = os.ExpandEnv(configPath)
	Cfg.DBConfigFile = os.ExpandEnv(get.EnvOr("CANALIZE_DS_CONFIG_PATH", "/ALL/CANALIZE/projects/DATABASE/canalize_ds/configs/config.json"))
	Cfg.EnvFile = os.ExpandEnv(get.EnvOr("CANALIZE_BE_ENV_PATH", "/ALL/CANALIZE/projects/BACKEND/canalize_be/.env"))
	Cfg.LogFile = os.ExpandEnv(get.EnvOr("CANALIZE_BE_LOG_FILE_PATH", "/ALL/CANALIZE/logs/canalize_be.log"))
	Cfg.GlobalRef = types.NewGlobalRef(get.EnvOr("CANALIZE_BE_PROCESS_NAME", "canalize_be")).GetGlobalRef()
	Cfg.Debug = get.EnvOrType("CANALIZE_BE_DEBUG_MODE", false)
	Cfg.ReleaseMode = get.EnvOrType("CANALIZE_BE_RELEASE_MODE", false)
	Cfg.IsConfidential = get.EnvOrType("CANALIZE_BE_CONFIDENCIAL_MODE", false)
	Cfg.Port = get.EnvOrType("CANALIZE_BE_PORT", "4000")
	Cfg.Host = baseURL
	Cfg.PrivKeyPath = privKeyPath
	Cfg.PubKeyPath = pubKeyPath
	Cfg.PubCertKeyPath = pubKeyPath
	Cfg.CORSEnabled = get.EnvOrType("CANALIZE_BE_ENABLE_CORS", true)
	Cfg.Debug = get.EnvOrType("CANALIZE_BE_DEBUG_MODE", false)
	Cfg.ProvidersConfig = os.ExpandEnv(get.EnvOr("CANALIZE_BE_PROVIDERS_CONFIG_PATH",
		"/ALL/CANALIZE/projects/BACKEND/canalize_be/configs/providers.yaml"))
	Cfg.RefreshTokenTTL = defaultTTL

	return Cfg
}

func NewSrvConfigFromParams(params *SrvConfig) *SrvConfig {
	Cfg := types.NewSrvConfig()

	Cfg.ConfigFile = get.ValOrType(params.ConfigFile, Cfg.ConfigFile)
	Cfg.DBConfigFile = get.ValOrType(params.DBConfigFile, Cfg.DBConfigFile)
	Cfg.EnvFile = get.ValOrType(params.EnvFile, Cfg.EnvFile)
	Cfg.LogFile = get.ValOrType(params.LogFile, Cfg.LogFile)
	Cfg.GlobalRef = get.ValOrType(params.GlobalRef, Cfg.GlobalRef)
	Cfg.Debug = get.ValOrType(params.Debug, Cfg.Debug)
	Cfg.ReleaseMode = get.ValOrType(params.ReleaseMode, Cfg.ReleaseMode)
	Cfg.IsConfidential = get.ValOrType(params.IsConfidential, Cfg.IsConfidential)
	Cfg.Port = get.ValOrType(params.Port, Cfg.Port)
	Cfg.Host = get.ValOrType(params.Host, Cfg.Host)
	Cfg.PrivKeyPath = get.ValOrType(params.PrivKeyPath, Cfg.PrivKeyPath)
	Cfg.PubKeyPath = get.ValOrType(params.PubKeyPath, Cfg.PubKeyPath)
	Cfg.PubCertKeyPath = get.ValOrType(params.PubCertKeyPath, Cfg.PubCertKeyPath)
	Cfg.CORSEnabled = params.CORSEnabled
	Cfg.ProvidersConfig = get.ValOrType(params.ProvidersConfig, Cfg.ProvidersConfig)
	Cfg.RefreshTokenTTL = get.ValOrType(params.RefreshTokenTTL, Cfg.RefreshTokenTTL)

	return Cfg
}

type GlobalRef = types.GlobalRef

func NewGlobalRef(name string) *GlobalRef { return types.NewGlobalRef(name) }

// ------------------------------- New Manifest Functions -----------------------------//

type Manifest = types.Manifest
type ManifestImpl = types.ManifestImpl

func NewManifestType() *ManifestImpl {
	return &ManifestImpl{}
}

func NewManifest() Manifest {
	return NewManifestType()
}

func EnsureGlobalManifest(n, c *ManifestImpl) {
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
	reflect.TypeFor[MailSrvParams](): true,
	reflect.TypeFor[MailConfig]():    true,
	reflect.TypeFor[LogzConfig]():    true,
	reflect.TypeFor[SrvConfig]():     true,
	reflect.TypeFor[GlobalRef]():     true,
	reflect.TypeFor[ManifestImpl]():  true,
}

var defaultFactories = map[reflect.Type]func() any{
	reflect.TypeFor[MailSrvParams](): func() any { return NewMailSrvParams("") },
	reflect.TypeFor[MailConfig]():    func() any { return NewMailConfig("") },
	reflect.TypeFor[LogzConfig]():    func() any { return NewLogzParams() },
	reflect.TypeFor[SrvConfig]():     func() any { return NewSrvArgs() },
	reflect.TypeFor[GlobalRef]():     func() any { return NewGlobalRef("default") },
	reflect.TypeFor[ManifestImpl]():  func() any { return NewManifestType() },
}

// LoadConfig loads a configuration of type T from the specified file path.

func LoadConfig[T any](cfgPath string) (*T, error) {
	if configRegistry[reflect.TypeFor[T]()] {
		cfgLoader := get.Loader[T](cfgPath)
		obj, err := cfgLoader.DeserializeFromFile(get.FileExt(cfgPath))
		if err != nil {
			return nil, err
		}
		if reflect.TypeFor[T]() == reflect.TypeFor[ManifestImpl]() {
			EnsureGlobalManifest(any(obj).(*ManifestImpl), types.KubexManifest)
		}
		return obj, nil
	}
	return nil, gl.Errorf("configuration type not registered")
}

func LoadConfigOrDefault[T MailConfig | MailConnection | LogzConfig | SrvConfig | MailSrvParams | Email | ManifestImpl](cfgPath string, genFile bool) (*T, error) {
	// Só entra aqui se o tipo for algum já registrado, então não me preocupo em checar o erro, só logo retorno o default
	cfgMapper := tools.NewEmptyMapperType[T](cfgPath)
	cfg, err := cfgMapper.DeserializeFromFile(get.FileExt(cfgPath))
	if err == nil {
		return cfg, nil
	}
	gl.Warnf("failed to load config from '%s', using default: %v", cfgPath, err)
	defaultCfg := defaultFactories[reflect.TypeFor[T]()]().(*T)
	if genFile {
		cfgMapper.SetValue(defaultCfg)
		cfgMapper.SerializeToFile(get.FileExt(cfgPath))
	}
	return defaultCfg, nil
}

// Package load provides functions to load configuration and environment settings.
package load

import (
	"net"
	"os"
	"reflect"
	"strings"
	"time"

	// "net"
	"net/url"

	"github.com/kubex-ecosystem/kbx/get"
	"github.com/kubex-ecosystem/kbx/is"
	"github.com/kubex-ecosystem/kbx/tools"
	"github.com/kubex-ecosystem/kbx/types"
	"golang.org/x/oauth2"

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

func ParseSrvArgs(bind string, port string, pubCertKeyPath string, pubKeyPath string, privKeyPath string, accessTokenTTL int, refreshTokenTTL int, issuer string, defaults map[string]any) SrvConfig {
	SrvArgs := NewSrvArgs()
	SrvArgs.Runtime.Bind = os.ExpandEnv(get.ValOrType(bind, get.EnvOr(defaults["DefaultServerHost"].(string), "0.0.0.0")))
	SrvArgs.Runtime.Port = get.ValOrType(port, get.EnvOr(defaults["DefaultServerPort"].(string), "5000"))
	SrvArgs.Runtime.PubCertKeyPath = os.ExpandEnv(get.ValOrType(pubCertKeyPath, get.EnvOr(defaults["DefaultGNyxPubCertKeyPath"].(string), "")))
	SrvArgs.Runtime.PubKeyPath = os.ExpandEnv(get.ValOrType(pubKeyPath, get.EnvOr(defaults["DefaultGNyxPubKeyPath"].(string), "")))
	SrvArgs.Runtime.PrivKeyPath = os.ExpandEnv(get.ValOrType(privKeyPath, get.EnvOr(defaults["DefaultCanalizeBEPrivKeyPath"].(string), "")))
	SrvArgs.Runtime.AccessTokenTTL = time.Duration(get.ValOrType(accessTokenTTL, 15)) * time.Minute
	SrvArgs.Runtime.RefreshTokenTTL = time.Duration(get.ValOrType(refreshTokenTTL, 60)) * time.Minute
	SrvArgs.Runtime.Issuer = get.ValOrType(issuer, "kubex-ecosystem")
	return SrvArgs
}

func NewSrvDefaultConfig(defaults map[string]any) SrvConfig {
	scheme := os.ExpandEnv(get.EnvOr("KUBEX_GNYX_SCHEME", "http"))
	host := os.ExpandEnv(get.EnvOr("KUBEX_GNYX_HOST", defaults["DefaultServerHost"].(string)))
	addr := net.JoinHostPort(host, get.EnvOr("KUBEX_GNYX_PORT", defaults["DefaultServerPort"].(string)))
	url := url.URL{Scheme: scheme, Host: addr}
	baseURL := get.ValueOrIf(get.EnvOr("KUBEX_ENV", "development") == "production",
		"https://api.kubex.world",
		url.String(),
	)
	defaultTTL := get.EnvOrType("INVITE_EXPIRATION", 7*24*time.Hour)
	configPath := os.ExpandEnv(get.EnvOr("KUBEX_GNYX_CONFIG_PATH", get.ValOrType(defaults["default_kubex_gnyx_config_path"].(string), "")))
	pubKeyPath := os.ExpandEnv(get.EnvOrType("KUBEX_GNYX_PUBLIC_KEY_PATH", get.ValOrType(defaults["default_kubex_gnyx_cert_path"].(string), "")))
	privKeyPath := os.ExpandEnv(get.EnvOrType("KUBEX_GNYX_PRIVATE_KEY_PATH", get.ValOrType(defaults["default_kubex_gnyx_key_path"].(string), "")))

	Cfg := types.NewSrvConfig()
	Cfg.Files.ConfigFile = os.ExpandEnv(configPath)
	Cfg.Files.DBConfigFile = os.ExpandEnv(get.EnvOr("KUBEX_DOMUS_CONFIG_PATH", get.ValOrType(defaults["default_kubex_domus_config_path"].(string), "$HOME/.kubex/domus/config/config.json")))
	Cfg.Files.EnvFile = os.ExpandEnv(get.EnvOr("KUBEX_GNYX_ENV_PATH", get.ValOrType(defaults["default_kubex_gnyx_env_path"].(string), "$HOME/mvp/kubex_gnyx_latest/.be.env")))
	Cfg.Files.LogFile = os.ExpandEnv(get.EnvOr("KUBEX_GNYX_LOG_FILE_PATH", get.ValOrType(defaults["default_kubex_gnyx_log_file_path"].(string), "$HOME/mvp/kubex_gnyx_latest/kubex_gnyx.log")))
	Cfg.GlobalRef = types.NewGlobalRef(get.EnvOr("KUBEX_GNYX_PROCESS_NAME", get.ValOrType(defaults["default_kubex_gnyx_process_name"].(string), "kubex_gnyx")))
	Cfg.Basic.Debug = get.EnvOrType("KUBEX_GNYX_DEBUG_MODE", false)
	Cfg.Basic.ReleaseMode = get.EnvOrType("KUBEX_GNYX_RELEASE_MODE", false)
	Cfg.Basic.IsConfidential = get.EnvOrType("KUBEX_GNYX_CONFIDENCIAL_MODE", false)
	Cfg.Runtime.Port = get.EnvOrType("KUBEX_GNYX_PORT", get.ValOrType(defaults["default_kubex_gnyx_port"].(string), "5000"))
	Cfg.Runtime.Host = baseURL
	Cfg.Runtime.PrivKeyPath = privKeyPath   // pragma: allowlist secret
	Cfg.Runtime.PubKeyPath = pubKeyPath     // pragma: allowlist secret
	Cfg.Runtime.PubCertKeyPath = pubKeyPath // pragma: allowlist secret
	Cfg.Basic.CORSEnabled = get.EnvOrType("KUBEX_GNYX_ENABLE_CORS", true)
	Cfg.Basic.Debug = get.EnvOrType("KUBEX_GNYX_DEBUG_MODE", false)
	Cfg.Files.ProvidersConfig = os.ExpandEnv(get.EnvOr("KUBEX_GNYX_PROVIDERS_CONFIG_PATH", get.ValOrType(defaults["default_kubex_gnyx_providers_config_path"].(string), "")))
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

// ------------------------------- New LLM Config Functions -----------------------------//

type LLMConfig = types.LLMConfig
type LLMProviderConfig = types.LLMProviderConfig
type LLMDevelopmentConfig = types.LLMDevelopmentConfig

func NewLLMConfig() LLMConfig                       { return LLMConfig{} }
func NewLLMProviderConfig() LLMProviderConfig       { return LLMProviderConfig{} }
func NewLLMDevelopmentConfig() LLMDevelopmentConfig { return LLMDevelopmentConfig{} }

func ParseLLMConfig(providers map[string]types.LLMProviderConfig, development types.LLMDevelopmentConfig) LLMConfig {
	LLMArgs := NewLLMConfig()

	LLMArgs.GlobalRef = types.NewGlobalRef(get.EnvOr("KUBEX_GNYX_PROCESS_NAME", "kubex_gnyx"))

	// LLMArgs.Providers =
	// LLMArgs.Development =
	// etc... (principalmente ptr's e/ou interfaces)

	return LLMArgs
}
func NewLLMConfigDefaultValues() LLMConfig {
	cfg := NewLLMConfig()

	cfg.GlobalRef = types.NewGlobalRef(get.EnvOr("KUBEX_GNYX_PROCESS_NAME", "kubex_gnyx"))

	// As configs serão carregadas de arquivos/envvars/args em flags,
	// aqui seriam os defaults totalmente básicos e genéricos. Eles podem ser
	// usados como um "full-load" de defaults para geração de arquivos de config,
	// ou como fallback caso alguma config específica não seja encontrada nos arquivos/envvars/args.
	cfg.Providers = map[string]*types.LLMProviderConfig{
		"groq": types.NewLLMProviderConfigType(
			"groq",
			"https://api.groq.com/v1",
			"GROQ_API_KEY",
			"",
		),
		"gemini": types.NewLLMProviderConfigType(
			"gemini",
			"https://ai.google.dev",
			"GEMINI_API_KEY",
			"",
		),
		"openai": types.NewLLMProviderConfigType(
			"openai",
			"https://api.openai.com",
			"OPENAI_API_KEY",
			"",
		),
		"anthropic": types.NewLLMProviderConfigType(
			"anthropic",
			"https://api.anthropic.com",
			"ANTHROPIC_API_KEY",
			"",
		),
		"azure": types.NewLLMProviderConfigType(
			"azure",
			"https://azure.microsoft.com/en-us/services/cognitive-services/openai-service/",
			"AZURE_API_KEY",
			"",
		),
		"deepseek": types.NewLLMProviderConfigType(
			"deepseek",
			"https://www.deepseek.com/api",
			"DEEPSEEK_API_KEY",
			"",
		),
		"custom": types.NewLLMProviderConfigType(
			// Para providers custom, a ideia é que o usuário forneça a implementação da interface `LLMProvider` e registre ela no registry, e aí a config do provider custom seria mais pra guardar informações como baseURL, chave de API, modelo default, etc... que seriam usadas pela implementação custom na hora de fazer as requisições pros endpoints do provider。
			"custom",
			"",
			"",
			"",
		),
	}

	perProviderCfg := types.LLMTokenBucket{
		Capacity:   1000,
		RefillRate: (1000 / 60), // 1000 tokens per minute
	}
	cfg.Development = types.LLMDevelopmentConfig{
		LoggingLevel: "info",
		HealthCheck: types.LLMHealthCheckConfig{
			Enabled:     true,
			IntervalSec: 60,
			TimeoutSec:  5,
		},
		CircuitBreaker: types.LLMCircuitBreakerConfig{
			Enabled: true,
			PerProvider: map[string]types.LLMCircuitBreakerRule{
				"groq": {
					MaxFailures:      5,
					ResetTimeoutSec:  300,
					SuccessThreshold: 3,
				},
				"gemini": {
					MaxFailures:      5,
					ResetTimeoutSec:  300,
					SuccessThreshold: 3,
				},
				"openai": {
					MaxFailures:      5,
					ResetTimeoutSec:  300,
					SuccessThreshold: 3,
				},
			},
			Default: types.LLMCircuitBreakerRule{
				MaxFailures:      5,
				ResetTimeoutSec:  300,
				SuccessThreshold: 3,
			},
		},
		RateLimit: types.LLMRateLimitConfig{
			Enabled: true,
			PerProvider: map[string]types.LLMTokenBucket{
				"groq":      perProviderCfg,
				"gemini":    perProviderCfg,
				"openai":    perProviderCfg,
				"anthropic": perProviderCfg,
				"azure":     perProviderCfg,
				"deepseek":  perProviderCfg,
				"custom": {
					Capacity:   10000,
					RefillRate: (10000 / 15), // 10000 tokens per 15 minutes
				},
			},
			Default: perProviderCfg,
		},
		Retry: types.LLMRetryConfig{
			Enabled:     true,
			MaxRetries:  3,
			BaseDelayMS: 500,  // 500 milliseconds
			MaxDelayMS:  5000, // 5 seconds
			Multiplier:  2,
			// TODO: implementar backoff exponencial com jitter na config (talvez já haja coisa pronta no `tools/retry.go`, lembro de ter feito algo nesse sentido)
			// Backoff: types.LLMBackoffConfig{
			// 	InitialDelay: 1000000000, // 1 second in nanoseconds
			// 	MaxDelay:     1000000000, // 1 second in nanoseconds
			// 	Multiplier:   2,
			// },
		},
		Defaults: types.LLMRequestDefaults{
			MaxTokens:        1000,
			Temperature:      0.7,
			TopP:             1,
			FrequencyPenalty: 0,
			PresencePenalty:  0,
			Stream:           false,
			TimeoutSec:       30,
			TenantID:         "",
			UserID:           "",
		},
	}

	defProvProdCfg := types.LLMProviderProductionConfig{
		TimeoutSec:  30,
		Priority:    "standard",
		MaxRetries:  3,
		BaseDelayMS: 500,
		MaxDelayMS:  5000,
		Multiplier:  2,
	}
	cfg.ProviderProduction = map[string]types.LLMProviderProductionConfig{
		"groq":      defProvProdCfg,
		"gemini":    defProvProdCfg,
		"openai":    defProvProdCfg,
		"anthropic": defProvProdCfg,
		"azure":     defProvProdCfg,
		"deepseek":  defProvProdCfg,
		"custom": {
			TimeoutSec:  120,
			Priority:    "low",
			MaxRetries:  3,
			BaseDelayMS: 500,
			MaxDelayMS:  5000,
			Multiplier:  1,
			// Para providers custom, a ideia é que o usuário forneça a implementação da interface `LLMProvider` e registre ela no registry, e aí a config do provider custom seria mais pra guardar informações como baseURL, chave de API, modelo default, etc... que seriam usadas pela implementação custom na hora de fazer as requisições pros endpoints do provider.
		},
	}

	cfg.Security = types.LLMSecurityConfig{
		EnableHTTPS:    false,
		AllowedOrigins: []string{"*"},
		JWTSecret:      "",
		APIKeys:        []string{},
	}
	for _, provCfg := range cfg.Providers {
		if get.EnvOr(provCfg.KeyRef(), "") != "" {
			cfg.Security.APIKeys = append(cfg.Security.APIKeys, provCfg.KeyRef())
			// Note: aqui a gente tá assumindo que o valor da chave de API é o próprio
			// nome da variável de ambiente onde a chave tá armazenada (ex: "OPENAI_API_KEY", "GROQ_API_KEY", etc...).
			// Todas as apikeys carregadas na aplicação (seja por env, arquivos de config, args, etc...) serão
			// definidas como variáveis de ambiente, assim o valor sempre será lido em tempo de execução.
			// Se o user inserir uma chave de API na config diretamente, a aplicação irá definir uma variável de ambiente
			// com o padrão de nomeclatura "KUBEX_GNYX_LLM_PROVIDER_{PROVIDER_NAME}_API_KEY" (ex: "KUBEX_GNYX_LLM_PROVIDER_OPENAI_API_KEY", "KUBEX_GNYX_LLM_PROVIDER_GROQ_API_KEY", etc...)
			// com o valor da chave de API, e aí o valor da config do provider custom seria
			// a variável de ambiente onde a chave tá armazenada, seguindo o mesmo padrão dos providers pré-definidos.
		}
	}

	cfg.License = "MIT"
	cfg.Monitoring = types.LLMMonitoringConfig{EnableMetrics: true}
	cfg.Repository = "github.com/kubex-ecosystem/gnyx"
	cfg.Authors = []string{"Rafael Mori"}
	cfg.Version = "0.0.1"

	return cfg
}
func NewLLMConfigFromParams(params *LLMConfig) LLMConfig {
	LLMArgs := NewLLMConfig()
	LLMArgs.GlobalRef = get.ValOrType(params.GlobalRef, LLMArgs.GlobalRef)
	LLMArgs.Providers = get.ValOrType(params.Providers, LLMArgs.Providers)
	LLMArgs.Development = get.ValOrType(params.Development, LLMArgs.Development)

	// etc... (principalmente ptr's e/ou interfaces)

	return LLMArgs
}

// ------------------------------- New Global Ref Functions -----------------------------//

type GlobalRef = types.GlobalRef

func NewGlobalRef(name string) GlobalRef { return types.NewGlobalRef(name) }

// ------------------------------- Google Auth Config Functions -----------------------------//

type AuthOAuthClientConfig = types.AuthOAuthClientConfig
type AuthClientConfig = types.AuthClientConfig
type AuthProvidersConfig = types.AuthProvidersConfig
type VendorAuthConfig = types.VendorAuthConfig

func NewVendorAuthConfig(cfgPath string) VendorAuthConfig {
	return VendorAuthConfig{
		AuthClientConfig: AuthClientConfig{
			AuthProvider: "google",
			// Web default config
			Web: AuthOAuthClientConfig{
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
				ProjectID:               "",
				Metadata:                make(map[string]any),
				Config:                  &oauth2.Config{},
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
	reflect.TypeFor[MailSrvParams]():         true,
	reflect.TypeFor[MailConfig]():            true,
	reflect.TypeFor[LogzConfig]():            true,
	reflect.TypeFor[SrvConfig]():             true,
	reflect.TypeFor[LLMConfig]():             true,
	reflect.TypeFor[LLMProviderConfig]():     true,
	reflect.TypeFor[LLMDevelopmentConfig]():  true,
	reflect.TypeFor[MManifest]():             true,
	reflect.TypeFor[VendorAuthConfig]():      true,
	reflect.TypeFor[AuthOAuthClientConfig](): true,
	reflect.TypeFor[Email]():                 true,
	reflect.TypeFor[MailConnection]():        true,
}

var defaultFactories = map[reflect.Type]func() any{
	reflect.TypeFor[MailSrvParams]():         func() any { return NewMailSrvParams("") },
	reflect.TypeFor[MailConfig]():            func() any { return NewMailConfig("") },
	reflect.TypeFor[LogzConfig]():            func() any { return NewLogzParams() },
	reflect.TypeFor[SrvConfig]():             func() any { return NewSrvArgs() },
	reflect.TypeFor[MManifest]():             func() any { return NewManifestType() },
	reflect.TypeFor[LLMConfig]():             func() any { return NewLLMConfigDefaultValues() },
	reflect.TypeFor[LLMProviderConfig]():     func() any { return NewLLMProviderConfig() },
	reflect.TypeFor[LLMDevelopmentConfig]():  func() any { return NewLLMDevelopmentConfig() },
	reflect.TypeFor[VendorAuthConfig]():      func() any { return NewVendorAuthConfig("") },
	reflect.TypeFor[AuthOAuthClientConfig](): func() any { return NewVendorAuthConfig("").Web },
	reflect.TypeFor[Email]():                 func() any { return types.NewEmail() },
	reflect.TypeFor[MailConnection]():        func() any { return types.NewMailConnection() },
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

func LoadConfigOrDefault[
	T MailConfig |
		MailConnection |
		LogzConfig |
		SrvConfig |
		LLMConfig |
		LLMProviderConfig |
		LLMDevelopmentConfig |
		MailSrvParams |
		Email |
		MManifest |
		VendorAuthConfig |
		AuthOAuthClientConfig](cfgPath string, genFile bool) (*T, error) {
	cfgPath = os.ExpandEnv(strings.TrimSpace(strings.ToValidUTF8(cfgPath, "")))
	if cfgPath == "" {
		return nil, gl.Errorf("configuration path cannot be empty")
	}

	// Só entra aqui se o tipo for algum já registrado, então não me preocupo em checar o erro, só logo retorno o default
	cfgMapper := tools.NewEmptyMapperType[T](cfgPath)
	cfg, err := cfgMapper.DeserializeFromFile(get.FileExt(cfgPath))
	if err == nil {
		return cfg, nil
	}
	gl.Warnf("failed to load config from '%s', using default: %v", cfgPath, err)

	if !is.Compatible[T](cfg) {
		gl.Warnf("loaded config is not compatible with expected type '%s'", reflect.TypeFor[T]())
		return nil, gl.Errorf("loaded config is not compatible with expected type '%s'", reflect.TypeFor[T]())
	}

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

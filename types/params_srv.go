package types

import (
	"time"

	"github.com/google/uuid"
)

type SrvBasicParams struct {
	// Basic options
	CompanyName  string `yaml:"company_name,omitempty" json:"company_name,omitempty" mapstructure:"company_name,omitempty"`
	FriendlyName string `yaml:"friendly_name,omitempty" json:"friendly_name,omitempty" mapstructure:"friendly_name,omitempty"`
	AppName      string `yaml:"app_name,omitempty" json:"app_name,omitempty" mapstructure:"app_name,omitempty"`
	AppVersion   string `yaml:"app_version,omitempty" json:"app_version,omitempty" mapstructure:"app_version,omitempty"`
	Environment  string `yaml:"environment,omitempty" json:"environment,omitempty" mapstructure:"environment,omitempty"`
	ContactEmail string `yaml:"contact_email,omitempty" json:"contact_email,omitempty" mapstructure:"contact_email,omitempty"`
	SupportEmail string `yaml:"support_email,omitempty" json:"support_email,omitempty" mapstructure:"support_email,omitempty"`

	Debug          bool     `yaml:"debug" json:"debug" mapstructure:"debug"`
	ReleaseMode    bool     `yaml:"release_mode" json:"release_mode" mapstructure:"release_mode"`
	IsConfidential bool     `yaml:"is_confidential" json:"is_confidential" mapstructure:"is_confidential"`
	CORSEnabled    bool     `yaml:"enable_cors" json:"enable_cors" mapstructure:"enable_cors"`
	TrustedProxies []string `yaml:"trusted_proxies" json:"trusted_proxies" mapstructure:"trusted_proxies"`
}

func NewSrvBasic() SrvBasicParams { return SrvBasicParams{} }

func NewSrvBasicDefault() SrvBasicParams {
	return SrvBasicParams{
		CompanyName:    "Canalize",
		FriendlyName:   "CanalizeBE",
		AppName:        "canalize",
		AppVersion:     "v0.1.0",
		Environment:    "development",
		ContactEmail:   "contato@canalize.app",
		SupportEmail:   "ti@canalize.app",
		Debug:          false,
		ReleaseMode:    false,
		IsConfidential: true,
		CORSEnabled:    true,
		TrustedProxies: []string{},
	}
}

type SrvFilesParams struct {
	// Paths and files
	Cwd              string `yaml:"cwd,omitempty" json:"cwd,omitempty" mapstructure:"cwd,omitempty"`
	LogFile          string `yaml:"log_file,omitempty" json:"log_file,omitempty" mapstructure:"log_file,omitempty"`
	EnvFile          string `yaml:"env_file,omitempty" json:"env_file,omitempty" mapstructure:"env_file,omitempty"`
	ConfigFile       string `yaml:"config_file,omitempty" json:"config_file,omitempty" mapstructure:"config_file,omitempty"`
	MainDBName       string `yaml:"main_db_name,omitempty" json:"main_db_name,omitempty" mapstructure:"main_db_name,omitempty"`
	DBConfigFile     string `yaml:"db_config_file,omitempty" json:"db_config_file,omitempty" mapstructure:"db_config_file,omitempty"`
	TemplatesDir     string `yaml:"templates_dir,omitempty" json:"templates_dir,omitempty" mapstructure:"templates_dir,omitempty"`
	MailerConfigFile string `yaml:"mailer_config_file,omitempty" json:"mailer_config_file,omitempty" mapstructure:"mailer_config_file,omitempty"`
	ProvidersConfig  string `yaml:"providers_config,omitempty" json:"providers_config,omitempty" mapstructure:"providers_config,omitempty"`
}

func NewSrvFiles() SrvFilesParams { return SrvFilesParams{} }

func NewSrvFilesDefault() SrvFilesParams {
	return SrvFilesParams{
		Cwd:              "",
		LogFile:          "",
		EnvFile:          "",
		ConfigFile:       "",
		MainDBName:       "canalize_db",
		DBConfigFile:     "",
		TemplatesDir:     "",
		MailerConfigFile: "",
		ProvidersConfig:  "",
	}
}

type SrvRuntimeParams struct {
	// Runtime options

	Host            string        `yaml:"host,omitempty" json:"host,omitempty" mapstructure:"host,omitempty"`
	Port            string        `yaml:"port,omitempty" json:"port,omitempty" mapstructure:"port,omitempty"`
	Bind            string        `yaml:"bind,omitempty" json:"bind,omitempty" mapstructure:"bind,omitempty"`
	PubCertKeyPath  string        `yaml:"pub_cert_key_path,omitempty" json:"pub_cert_key_path,omitempty" mapstructure:"pub_cert_key_path,omitempty"`
	PubKeyPath      string        `yaml:"pub_key_path,omitempty" json:"pub_key_path,omitempty" mapstructure:"pub_key_path,omitempty"`
	PrivKeyPath     string        `yaml:"priv_key_path,omitempty" json:"priv_key_path,omitempty" mapstructure:"priv_key_path,omitempty"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl,omitempty" json:"access_token_ttl,omitempty" mapstructure:"access_token_ttl,omitempty"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl,omitempty" json:"refresh_token_ttl,omitempty" mapstructure:"refresh_token_ttl,omitempty"`
	Issuer          string        `yaml:"issuer,omitempty" json:"issuer,omitempty" mapstructure:"issuer,omitempty"`
}

func NewSrvRuntime() SrvRuntimeParams { return SrvRuntimeParams{} }

func NewSrvRuntimeDefault() SrvRuntimeParams {
	return SrvRuntimeParams{
		Host:            "localhost",
		Port:            "4000",
		Bind:            ":8080",
		PubCertKeyPath:  "",
		PubKeyPath:      "",
		PrivKeyPath:     "",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 60 * time.Minute,
		Issuer:          "kubex-ecosystem",
	}
}

type SrvAdvancedParams struct {
	// Advanced options

	Context    string            `yaml:"context,omitempty" json:"context,omitempty" mapstructure:"context,omitempty"`
	Command    string            `yaml:"command,omitempty" json:"command,omitempty" mapstructure:"command,omitempty"`
	Subcommand string            `yaml:"subcommand,omitempty" json:"subcommand,omitempty" mapstructure:"subcommand,omitempty"`
	Args       string            `yaml:"args,omitempty" json:"args,omitempty" mapstructure:"args,omitempty"`
	EnvVars    map[string]string `yaml:"env_vars,omitempty" json:"env_vars,omitempty" mapstructure:"env_vars,omitempty"`
}

func NewSrvAdvanced() SrvAdvancedParams { return SrvAdvancedParams{} }

func NewSrvAdvancedDefault() SrvAdvancedParams {
	return SrvAdvancedParams{
		Context:    "",
		Command:    "",
		Subcommand: "",
		Args:       "",
		EnvVars:    make(map[string]string),
	}
}

type SrvFlagsParams struct {
	// Flags

	FailFast  bool `yaml:"fail_fast,omitempty" json:"fail_fast,omitempty" mapstructure:"fail_fast,omitempty"`
	Verbose   bool `yaml:"verbose,omitempty" json:"verbose,omitempty" mapstructure:"verbose,omitempty"`
	BatchMode bool `yaml:"batch_mode,omitempty" json:"batch_mode,omitempty" mapstructure:"batch_mode,omitempty"`
	NoColor   bool `yaml:"no_color,omitempty" json:"no_color,omitempty" mapstructure:"no_color,omitempty"`
	TraceMode bool `yaml:"trace_mode,omitempty" json:"trace_mode,omitempty" mapstructure:"trace_mode,omitempty"`
	RootMode  bool `yaml:"root_mode,omitempty" json:"root_mode,omitempty" mapstructure:"root_mode,omitempty"`
}

func NewSrvFlags() SrvFlagsParams { return SrvFlagsParams{} }

func NewSrvFlagsDefault() SrvFlagsParams {
	return SrvFlagsParams{
		FailFast:  false,
		Verbose:   false,
		BatchMode: false,
		NoColor:   false,
		TraceMode: false,
		RootMode:  false,
	}
}

type SrvPerformanceParams struct {
	// Performance options

	MaxProcs  int    `yaml:"max_procs,omitempty" json:"max_procs,omitempty" mapstructure:"max_procs,omitempty"`
	TimeoutMS int    `yaml:"timeout_ms,omitempty" json:"timeout_ms,omitempty" mapstructure:"timeout_ms,omitempty"`
	Hash      string `yaml:"hash,omitempty" json:"hash,omitempty" mapstructure:"hash,omitempty"`
}

func NewSrvPerformance() SrvPerformanceParams { return SrvPerformanceParams{} }

func NewSrvPerformanceDefault() SrvPerformanceParams {
	return SrvPerformanceParams{
		MaxProcs:  0,
		TimeoutMS: 5000,
		Hash:      "",
	}
}

// InviteConfig controla opções de envio e branding.
type InviteConfig struct {
	BaseURL     string        `json:"base_url,omitempty" yaml:"base_url,omitempty" toml:"base_url,omitempty" mapstructure:"base_url,omitempty"`
	SenderName  string        `json:"sender_name,omitempty" yaml:"sender_name,omitempty" toml:"sender_name,omitempty" mapstructure:"sender_name,omitempty"`
	SenderEmail string        `json:"sender_email,omitempty" yaml:"sender_email,omitempty" toml:"sender_email,omitempty" mapstructure:"sender_email,omitempty"`
	CompanyName string        `json:"company_name,omitempty" yaml:"company_name,omitempty" toml:"company_name,omitempty" mapstructure:"company_name,omitempty"`
	DefaultTTL  time.Duration `json:"default_ttl,omitempty" yaml:"default_ttl,omitempty" toml:"default_ttl,omitempty" mapstructure:"default_ttl,omitempty"`
}

type AuthConfig struct {
	JWTSecret           string        `json:"jwt_secret,omitempty" yaml:"jwt_secret,omitempty" toml:"jwt_secret,omitempty" mapstructure:"jwt_secret,omitempty"`
	AccessTokenTTL      time.Duration `json:"access_token_ttl,omitempty" yaml:"access_token_ttl,omitempty" toml:"access_token_ttl,omitempty" mapstructure:"access_token_ttl,omitempty"`
	RefreshTokenTTL     time.Duration `json:"refresh_token_ttl,omitempty" yaml:"refresh_token_ttl,omitempty" toml:"refresh_token_ttl,omitempty" mapstructure:"refresh_token_ttl,omitempty"`
	PasswordSaltRounds  int           `json:"password_salt_rounds,omitempty" yaml:"password_salt_rounds,omitempty" toml:"password_salt_rounds,omitempty" mapstructure:"password_salt_rounds,omitempty"`
	EnableEmailVerified bool          `json:"enable_email_verified,omitempty" yaml:"enable_email_verified,omitempty" toml:"enable_email_verified,omitempty" mapstructure:"enable_email_verified,omitempty"`
	Invite              InviteConfig  `json:"invite,omitempty" yaml:"invite,omitempty" toml:"invite,omitempty" mapstructure:"invite,omitempty"`
}

type SrvConfig struct {
	GlobalRef   `json:",inline" yaml:",inline" mapstructure:",squash"`
	Basic       SrvBasicParams       `json:",inline" yaml:",inline" mapstructure:",squash"`
	Files       SrvFilesParams       `json:",inline" yaml:",inline" mapstructure:",squash"`
	Runtime     SrvRuntimeParams     `json:",inline" yaml:",inline" mapstructure:",squash"`
	Advanced    SrvAdvancedParams    `json:",inline" yaml:",inline" mapstructure:",squash"`
	Flags       SrvFlagsParams       `json:",inline" yaml:",inline" mapstructure:",squash"`
	Performance SrvPerformanceParams `json:",inline" yaml:",inline" mapstructure:",squash"`
	Auth        AuthConfig           `json:"auth,omitempty" yaml:"auth,omitempty" mapstructure:"auth,omitempty"`
}

func NewSrvConfig() SrvConfig {
	return SrvConfig{
		GlobalRef:   GlobalRef{ID: uuid.New()},
		Basic:       NewSrvBasic(),
		Files:       NewSrvFiles(),
		Runtime:     NewSrvRuntime(),
		Advanced:    NewSrvAdvanced(),
		Flags:       NewSrvFlags(),
		Performance: NewSrvPerformance(),
		Auth:        AuthConfig{},
	}
}

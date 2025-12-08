package kbx

import "time"

type ServerParams struct {
	// Basic options

	Debug          bool     `yaml:"debug" json:"debug" mapstructure:"debug"`
	ReleaseMode    bool     `yaml:"release_mode" json:"release_mode" mapstructure:"release_mode"`
	IsConfidential bool     `yaml:"is_confidential" json:"is_confidential" mapstructure:"is_confidential"`
	CORSEnabled    bool     `yaml:"enable_cors" json:"enable_cors" mapstructure:"enable_cors"`
	TrustedProxies []string `yaml:"trusted_proxies" json:"trusted_proxies" mapstructure:"trusted_proxies"`

	// Paths and files

	Cwd             string `yaml:"cwd,omitempty" json:"cwd,omitempty" mapstructure:"cwd,omitempty"`
	LogFile         string `yaml:"log_file,omitempty" json:"log_file,omitempty" mapstructure:"log_file,omitempty"`
	EnvFile         string `yaml:"env_file,omitempty" json:"env_file,omitempty" mapstructure:"env_file,omitempty"`
	ConfigFile      string `yaml:"config_file,omitempty" json:"config_file,omitempty" mapstructure:"config_file,omitempty"`
	DBConfigFile    string `yaml:"db_config_file,omitempty" json:"db_config_file,omitempty" mapstructure:"db_config_file,omitempty"`
	ProvidersConfig string `yaml:"providers_config,omitempty" json:"providers_config,omitempty" mapstructure:"providers_config,omitempty"`

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

	// Advanced options

	Context    string            `yaml:"context,omitempty" json:"context,omitempty" mapstructure:"context,omitempty"`
	Command    string            `yaml:"command,omitempty" json:"command,omitempty" mapstructure:"command,omitempty"`
	Subcommand string            `yaml:"subcommand,omitempty" json:"subcommand,omitempty" mapstructure:"subcommand,omitempty"`
	Args       string            `yaml:"args,omitempty" json:"args,omitempty" mapstructure:"args,omitempty"`
	EnvVars    map[string]string `yaml:"env_vars,omitempty" json:"env_vars,omitempty" mapstructure:"env_vars,omitempty"`

	// Flags

	FailFast  bool `yaml:"fail_fast,omitempty" json:"fail_fast,omitempty" mapstructure:"fail_fast,omitempty"`
	Verbose   bool `yaml:"verbose,omitempty" json:"verbose,omitempty" mapstructure:"verbose,omitempty"`
	BatchMode bool `yaml:"batch_mode,omitempty" json:"batch_mode,omitempty" mapstructure:"batch_mode,omitempty"`
	NoColor   bool `yaml:"no_color,omitempty" json:"no_color,omitempty" mapstructure:"no_color,omitempty"`
	TraceMode bool `yaml:"trace_mode,omitempty" json:"trace_mode,omitempty" mapstructure:"trace_mode,omitempty"`
	RootMode  bool `yaml:"root_mode,omitempty" json:"root_mode,omitempty" mapstructure:"root_mode,omitempty"`

	// Performance options

	MaxProcs  int    `yaml:"max_procs,omitempty" json:"max_procs,omitempty" mapstructure:"max_procs,omitempty"`
	TimeoutMS int    `yaml:"timeout_ms,omitempty" json:"timeout_ms,omitempty" mapstructure:"timeout_ms,omitempty"`
	Hash      string `yaml:"hash,omitempty" json:"hash,omitempty" mapstructure:"hash,omitempty"`
}

func NewServerParamsRaw() *ServerParams {
	return &ServerParams{
		Debug:          false,
		ReleaseMode:    false,
		IsConfidential: true,
		CORSEnabled:    false,
		TrustedProxies: []string{},

		Host:            "localhost",
		Port:            "8080",
		Bind:            "0.0.0.0:8080",
		PubCertKeyPath:  "",
		PubKeyPath:      "",
		PrivKeyPath:     "",
		AccessTokenTTL:  0,
		RefreshTokenTTL: 0,
		Issuer:          "",

		Context:    "",
		Command:    "",
		Subcommand: "",
		Args:       "",
		EnvVars:    map[string]string{},

		FailFast:  false,
		Verbose:   false,
		BatchMode: false,
		NoColor:   false,
		TraceMode: false,
		RootMode:  false,

		MaxProcs:  0,
		TimeoutMS: 0,
		Hash:      "",
	}
}

func NewServerParamsEmpty() *ServerParams {
	return &ServerParams{}
}

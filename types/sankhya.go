package types

import (
	"time"

	kbxMod "github.com/kubex-ecosystem/kbx/internal/module/kbx"
)

// SankhyaServerBase is the base configuration for a Sankhya server.
type SankhyaServerBase struct {
	Host               string        `yaml:"host" json:"host" mapstructure:"host"`
	ServiceManagerPath string        `yaml:"service_manager_path" json:"service_manager_path" mapstructure:"service_manager_path"`
	DefaultTimeout     time.Duration `yaml:"default_timeout" json:"default_timeout" mapstructure:"default_timeout"`
}

// SankhyaJSessionAuth is the jsession_id authentication configuration for a Sankhya server.
type SankhyaJSessionAuth struct {
	Username   string `yaml:"username" json:"username" mapstructure:"username"`
	Password   string `yaml:"password" json:"password" mapstructure:"password"`
	JSessionID string `yaml:"jsession_id,omitempty" json:"jsession_id,omitempty" mapstructure:"jsession_id,omitempty"`
}

// SankhyaSimpleAuth is the simple authentication configuration for a Sankhya server.
type SankhyaSimpleAuth struct {
	Username  string        `yaml:"username" json:"username" mapstructure:"username"`
	Password  string        `yaml:"password" json:"password" mapstructure:"password"`
	AppKey    string        `yaml:"app_key" json:"app_key" mapstructure:"app_key"`
	Token     string        `yaml:"token,omitempty" json:"token,omitempty" mapstructure:"token,omitempty"`
	ExpiresIn time.Duration `yaml:"expires_in,omitempty" json:"expires_in,omitempty" mapstructure:"expires_in,omitempty"`
}

// SankhyaOauth2Auth is the oauth2 authentication configuration for a Sankhya server.
type SankhyaOauth2Auth struct {
	ClientID     string        `yaml:"client_id" json:"client_id" mapstructure:"client_id"`
	ClientSecret string        `yaml:"client_secret" json:"client_secret" mapstructure:"client_secret"`
	RedirectURI  string        `yaml:"redirect_uri" json:"redirect_uri" mapstructure:"redirect_uri"`
	Scope        string        `yaml:"scope" json:"scope" mapstructure:"scope"`
	TokenURL     string        `yaml:"token_url" json:"token_url" mapstructure:"token_url"`
	AuthCode     string        `yaml:"-" json:"-" mapstructure:"-"`
	TokenType    string        `yaml:"-" json:"-" mapstructure:"-"`
	AccessToken  string        `yaml:"-" json:"-" mapstructure:"-"`
	RefreshToken string        `yaml:"-" json:"-" mapstructure:"-"`
	IDToken      string        `yaml:"-" json:"-" mapstructure:"-"`
	ExpiresIn    time.Duration `yaml:"expires_in,omitempty" json:"expires_in,omitempty" mapstructure:"expires_in,omitempty"`
}

// SankhyaAuthMethod is the authentication method for a Sankhya server.
type SankhyaAuthMethod struct {
	Method     string               `yaml:"method" json:"method" mapstructure:"method"`
	JSessionID *SankhyaJSessionAuth `yaml:"jsession_id,omitempty" json:"jsession_id,omitempty" mapstructure:"jsession_id,omitempty"`
	Simple     *SankhyaSimpleAuth   `yaml:"simple,omitempty" json:"simple,omitempty" mapstructure:"simple,omitempty"`
	OAuth2     *SankhyaOauth2Auth   `yaml:"oauth2,omitempty" json:"oauth2,omitempty" mapstructure:"oauth2,omitempty"`
}

// SankhyaConfig is the main configuration for a Sankhya server.
type SankhyaConfig struct {
	// Embed server base configuration
	SankhyaServerBase `json:",inline" yaml:",inline" toml:",inline" mapstructure:",squash"`

	// Authentication method
	AuthenticationMethod SankhyaAuthMethod `yaml:"authentication_method" json:"authentication_method" mapstructure:"authentication_method"`
}

// NewSankhyaConfig creates a new Sankhya configuration.
// To get default values, pass nil for the parameters s and a.
func NewSankhyaConfig(s *SankhyaServerBase, a *SankhyaAuthMethod) *SankhyaConfig {
	return &SankhyaConfig{
		SankhyaServerBase: kbxMod.GetValueOrDefaultSimple(
			*s, //If s is nil, the default value will be used.
			SankhyaServerBase{
				ServiceManagerPath: "/mge/service.sbr",
				DefaultTimeout:     30 * time.Second,
			},
		),
		AuthenticationMethod: kbxMod.GetValueOrDefaultSimple(
			*a, //If a is nil, the default value will be used.
			SankhyaAuthMethod{
				Method: "oauth2",             //Default authentication method if none is provided. (Safest for production)
				Simple: &SankhyaSimpleAuth{}, //Use this only for quick tests. Not recommended for production.

				OAuth2: &SankhyaOauth2Auth{},
			},
		),
	}
}

// SankhyaContext is the context for a Sankhya server.
type SankhyaContext struct {
	SankhyaConfig
}

// SankhyaServers is the list of Sankhya servers configuration.
type SankhyaServers struct {
	Servers []SankhyaContext `yaml:"servers" json:"servers" mapstructure:"servers"`
}

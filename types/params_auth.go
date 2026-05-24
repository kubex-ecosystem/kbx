package types

import (
	"time"

	"golang.org/x/oauth2"
)

// InviteConfig controla opções de envio e branding.
type InviteConfig struct {
	BaseURL     string        `json:"base_url,omitempty" yaml:"base_url,omitempty" toml:"base_url,omitempty" mapstructure:"base_url,omitempty"`
	SenderName  string        `json:"sender_name,omitempty" yaml:"sender_name,omitempty" toml:"sender_name,omitempty" mapstructure:"sender_name,omitempty"`
	SenderEmail string        `json:"sender_email,omitempty" yaml:"sender_email,omitempty" toml:"sender_email,omitempty" mapstructure:"sender_email,omitempty"`
	CompanyName string        `json:"company_name,omitempty" yaml:"company_name,omitempty" toml:"company_name,omitempty" mapstructure:"company_name,omitempty"`
	DefaultTTL  time.Duration `json:"default_ttl,omitempty" yaml:"default_ttl,omitempty" toml:"default_ttl,omitempty" mapstructure:"default_ttl,omitempty"`
}

// NewInviteConfig cria uma nova instância de InviteConfig.
func NewInviteConfig() InviteConfig { return InviteConfig{} }

// NewInviteConfigDefault cria uma nova instância de InviteConfig com valores padrão.
func NewInviteConfigDefault() InviteConfig {
	return InviteConfig{
		BaseURL:     "https://gnyx.kubex.world",
		SenderName:  "Kubex Team",
		SenderEmail: "contact@kubex.world",
		CompanyName: "Kubex",
		DefaultTTL:  7 * 24 * time.Hour,
	}
}

// OptionsMap represents a map of client options.
type OptionsMap map[string]any

// BasicAuth represents a simple client configuration.
type BasicAuth struct {
	Username  string        `yaml:"username" json:"username" mapstructure:"username"`
	Password  string        `yaml:"password" json:"password" mapstructure:"password"`
	AppKey    string        `yaml:"app_key" json:"app_key" mapstructure:"app_key"`
	Token     string        `yaml:"token,omitempty" json:"token,omitempty" mapstructure:"token,omitempty"`
	ExpiresIn time.Duration `yaml:"expires_in,omitempty" json:"expires_in,omitempty" mapstructure:"expires_in,omitempty"`
}

// AuthProviderOptions represents a configuration for authentication options.
type AuthProviderOptions struct {
	JWTSecret           string                `json:"jwt_secret,omitempty" yaml:"jwt_secret,omitempty" toml:"jwt_secret,omitempty" mapstructure:"jwt_secret,omitempty"`
	AccessTokenTTL      time.Duration         `json:"access_token_ttl,omitempty" yaml:"access_token_ttl,omitempty" toml:"access_token_ttl,omitempty" mapstructure:"access_token_ttl,omitempty"`
	RefreshTokenTTL     time.Duration         `json:"refresh_token_ttl,omitempty" yaml:"refresh_token_ttl,omitempty" toml:"refresh_token_ttl,omitempty" mapstructure:"refresh_token_ttl,omitempty"`
	PasswordSaltRounds  int                   `json:"password_salt_rounds,omitempty" yaml:"password_salt_rounds,omitempty" toml:"password_salt_rounds,omitempty" mapstructure:"password_salt_rounds,omitempty"`
	EnableEmailVerified bool                  `json:"enable_email_verified,omitempty" yaml:"enable_email_verified,omitempty" toml:"enable_email_verified,omitempty" mapstructure:"enable_email_verified,omitempty"`
	Invite              InviteConfig          `json:"invite" yaml:"invite,omitempty" toml:"invite,omitempty" mapstructure:"invite,omitempty"`
	AuthProviders       AuthProvidersRegistry `json:"auth_providers" yaml:"auth_providers,omitempty" toml:"auth_providers,omitempty" mapstructure:"auth_providers,omitempty"`
	Options             OptionsMap            `json:"options,omitempty" yaml:"options,omitempty" toml:"options,omitempty" mapstructure:"options,omitempty"`
}

// AuthProvider represents an OAuth client configuration.
// (*oauth2.Config embeddado no struct para herança)
type AuthProvider struct {
	*oauth2.Config `json:",inline,omitempty" yaml:",inline,omitempty" toml:",inline,omitempty" mapstructure:",squash,omitempty"`
	*BasicAuth     `json:",inline,omitempty" yaml:",inline,omitempty" toml:",inline,omitempty" mapstructure:",squash,omitempty"`

	ConfigPath              string   `json:"config_path" yaml:"config_path" toml:"config_path" mapstructure:"config_path"`
	ProjectID               string   `json:"project_id,omitempty" env:"GOOGLE_PROJECT_ID"`
	APIKey                  string   `json:"-" yaml:"-" toml:"-" mapstructure:"-"`
	ClientID                string   `json:"client_id,omitempty" env:"GOOGLE_CLIENT_ID"`
	ClientSecret            string   `json:"client_secret,omitempty" env:"GOOGLE_CLIENT_SECRET"`
	RedirectURL             string   `json:"redirect_url,omitempty" env:"GOOGLE_REDIRECT_URL"`
	AuthURI                 string   `json:"auth_uri"`
	TokenURI                string   `json:"token_uri"`
	AuthProviderX509CertURL string   `json:"auth_provider_x509_cert_url"`
	MapUserInfo             bool     `json:"map_user_info,omitempty" env:"GOOGLE_MAP_USER_INFO"`
	MetadataOnly            bool     `json:"metadata_only,omitempty" env:"GOOGLE_METADATA_ONLY"`
	Scopes                  []string `json:"scopes,omitempty" env:"GOOGLE_SCOPES"`
	RedirectURIs            []string `json:"redirect_uris,omitempty" env:"GOOGLE_REDIRECT_URIS"`
	JavaScriptOrigins       []string `json:"javascript_origins,omitempty" env:"GOOGLE_JAVASCRIPT_ORIGINS"`

	Metadata map[string]any      `json:"metadata,omitempty" env:"GOOGLE_METADATA"`
	Options  AuthProviderOptions `json:"options,omitempty" yaml:"options,omitempty" toml:"options,omitempty" mapstructure:"options,omitempty"`
}

// NewAuthClient returns a new instance of AuthClient.
func NewAuthClient(c *AuthProvider) *AuthProvider {
	if c != nil {
		return c
	}
	return &AuthProvider{}
}

// AuthProviderWrapper represents a authentication configuration.
type AuthProviderWrapper struct {
	Web      *AuthProvider          `json:"web" yaml:"web,omitempty" toml:"web,omitempty" mapstructure:"web,omitempty"`
	Mobile   *AuthProvider          `json:"mobile" yaml:"mobile,omitempty" toml:"mobile,omitempty" mapstructure:"mobile,omitempty"`
	API      *AuthProvider          `json:"api" yaml:"api,omitempty" toml:"api,omitempty" mapstructure:"api,omitempty"`
	Internal *AuthProvider          `json:"internal" yaml:"internal,omitempty" toml:"internal,omitempty" mapstructure:"internal,omitempty"`
	Invite   *AuthProvider          `json:"invite" yaml:"invite,omitempty" toml:"invite,omitempty" mapstructure:"invite,omitempty"`
	Config   *AuthProvidersRegistry `json:"auth_providers_config" yaml:"auth_providers_config,omitempty" toml:"auth_providers_config,omitempty" mapstructure:"auth_providers_config,omitempty"`
}

// NewAuthClientWrapper creates a new instance of AuthClientWrapper.
func NewAuthClientWrapper(w *AuthProviderWrapper) *AuthProviderWrapper {
	if w != nil {
		return w
	}
	return &AuthProviderWrapper{}
}

// AuthProvidersRegistry represents a configuration for multiple authentication providers.
type AuthProvidersRegistry struct {
	Sankhya   *AuthProviderWrapper   `json:"sankhya,omitempty" env:"SANKHYA_AUTH_CONFIG"`
	Google    *AuthProviderWrapper   `json:"google,omitempty" env:"GOOGLE_AUTH_CONFIG"`
	Firebase  *AuthProviderWrapper   `json:"firebase,omitempty" env:"FIREBASE_AUTH_CONFIG"`
	Microsoft *AuthProviderWrapper   `json:"microsoft,omitempty" env:"MICROSOFT_AUTH_CONFIG"`
	Facebook  *AuthProviderWrapper   `json:"facebook,omitempty" env:"FACEBOOK_AUTH_CONFIG"`
	LinkedIn  *AuthProviderWrapper   `json:"linkedin,omitempty" env:"LINKEDIN_AUTH_CONFIG"`
	Twitter   *AuthProviderWrapper   `json:"twitter,omitempty" env:"TWITTER_AUTH_CONFIG"`
	Github    *AuthProviderWrapper   `json:"github,omitempty" env:"GITHUB_AUTH_CONFIG"`
	Custom    []*AuthProviderWrapper `json:"custom,omitempty" env:"CUSTOM_AUTH_CONFIG"`
}

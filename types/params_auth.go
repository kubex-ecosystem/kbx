package types

import "time"

type AuthClientOptions map[string]any

type AuthOAuthClientConfig struct {
	ProjectID               string         `json:"project_id,omitempty" env:"GOOGLE_PROJECT_ID"`
	ClientID                string         `json:"client_id,omitempty" env:"GOOGLE_CLIENT_ID"`
	ClientSecret            string         `json:"client_secret,omitempty" env:"GOOGLE_CLIENT_SECRET"` // Cuidado com esse log!
	RedirectURL             string         `json:"redirect_url,omitempty" env:"GOOGLE_REDIRECT_URL"`
	AuthURI                 string         `json:"auth_uri"`
	TokenURI                string         `json:"token_uri"`
	AuthProviderX509CertURL string         `json:"auth_provider_x509_cert_url"`
	MapUserInfo             bool           `json:"map_user_info,omitempty" env:"GOOGLE_MAP_USER_INFO"`
	MetadataOnly            bool           `json:"metadata_only,omitempty" env:"GOOGLE_METADATA_ONLY"`
	Scopes                  []string       `json:"scopes,omitempty" env:"GOOGLE_SCOPES"`
	RedirectURIs            []string       `json:"redirect_uris,omitempty" env:"GOOGLE_REDIRECT_URIS"`
	JavaScriptOrigins       []string       `json:"javascript_origins,omitempty" env:"GOOGLE_JAVASCRIPT_ORIGINS"`
	Metadata                map[string]any `json:"metadata,omitempty" env:"GOOGLE_METADATA"`
}

type AuthClientConfig struct {
	Web          AuthOAuthClientConfig `json:"web" yaml:"web,omitempty" toml:"web,omitempty" mapstructure:"web,omitempty"`
	AuthProvider string                `json:"auth_provider,omitempty" env:"AUTH_PROVIDER"`
	Options      AuthClientOptions     `json:"options,omitempty" env:"AUTH_OPTIONS"`
}

type AuthProvidersConfig struct {
	Google   AuthClientConfig `json:"google,omitempty" env:"GOOGLE_AUTH_CONFIG"`
	Facebook AuthClientConfig `json:"facebook,omitempty" env:"FACEBOOK_AUTH_CONFIG"`
	Github   AuthClientConfig `json:"github,omitempty" env:"GITHUB_AUTH_CONFIG"`
}

type AuthConfig struct {
	JWTSecret           string              `json:"jwt_secret,omitempty" yaml:"jwt_secret,omitempty" toml:"jwt_secret,omitempty" mapstructure:"jwt_secret,omitempty"`
	AccessTokenTTL      time.Duration       `json:"access_token_ttl,omitempty" yaml:"access_token_ttl,omitempty" toml:"access_token_ttl,omitempty" mapstructure:"access_token_ttl,omitempty"`
	RefreshTokenTTL     time.Duration       `json:"refresh_token_ttl,omitempty" yaml:"refresh_token_ttl,omitempty" toml:"refresh_token_ttl,omitempty" mapstructure:"refresh_token_ttl,omitempty"`
	PasswordSaltRounds  int                 `json:"password_salt_rounds,omitempty" yaml:"password_salt_rounds,omitempty" toml:"password_salt_rounds,omitempty" mapstructure:"password_salt_rounds,omitempty"`
	EnableEmailVerified bool                `json:"enable_email_verified,omitempty" yaml:"enable_email_verified,omitempty" toml:"enable_email_verified,omitempty" mapstructure:"enable_email_verified,omitempty"`
	Invite              InviteConfig        `json:"invite" yaml:"invite,omitempty" toml:"invite,omitempty" mapstructure:"invite,omitempty"`
	AuthProvidersConfig AuthProvidersConfig `json:"auth_providers_config" yaml:"auth_providers_config,omitempty" toml:"auth_providers_config,omitempty" mapstructure:"auth_providers_config,omitempty"`
}

type VendorAuthConfig struct {
	AuthClientConfig
	AuthProvider string `json:"auth_provider,omitempty" yaml:"auth_provider,omitempty" xml:"auth_provider,omitempty" toml:"auth_provider,omitempty" mapstructure:"auth_provider,omitempty"`
	ConfigPath   string `json:"config_path,omitempty" yaml:"config_path,omitempty" xml:"config_path,omitempty" toml:"config_path,omitempty" mapstructure:"config_path,omitempty"`
}

// type AuthClientPlatform struct {
// 	ID     string `json:"id,omitempty" env:"AUTH_PLATFORM_ID"`
// 	Name   string `json:"name,omitempty" env:"AUTH_PLATFORM_NAME"`
// 	Icon   string `json:"icon,omitempty" env:"AUTH_PLATFORM_ICON"`
// 	APIURL string `json:"api_url,omitempty" env:"AUTH_PLATFORM_API_URL"`
// }

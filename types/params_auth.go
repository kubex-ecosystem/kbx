package types

type AuthClientOptions map[string]any

// type AuthSession struct {
// 	GlobalRef `json:",inline" yaml:",inline" mapstructure:",squash"`

// 	// Identification
// 	Provider    string `json:"-" yaml:"-" mapstructure:"-"`
// 	Issuer      string `json:"-" yaml:"-" mapstructure:"-"`
// 	AccessType  string `json:"-" yaml:"-" mapstructure:"-"`
// 	AuthURL     string `json:"-" yaml:"-" mapstructure:"-"`
// 	TokenURL    string `json:"-" yaml:"-" mapstructure:"-"`
// 	RedirectURL string `json:"-" yaml:"-" mapstructure:"-"`
// 	Locale      string `json:"-" yaml:"-" mapstructure:"-"`
// 	Subject     string `json:"-" yaml:"-" mapstructure:"-"`

// 	// Tokens
// 	AccessToken   string         `json:"-" yaml:"-" mapstructure:"-"`
// 	RefreshToken  string         `json:"-" yaml:"-" mapstructure:"-"`
// 	IDToken       string         `json:"-" yaml:"-" mapstructure:"-"`
// 	RawIDToken    string         `json:"-" yaml:"-" mapstructure:"-"`
// 	Nonce         string         `json:"-" yaml:"-" mapstructure:"-"`
// 	AuthCode      string         `json:"-" yaml:"-" mapstructure:"-"`
// 	IDTokenClaims map[string]any `json:"-" yaml:"-" mapstructure:"-"`

// 	// Metadata
// 	ExpiresIn int64     `json:"-" yaml:"-" mapstructure:"-"`
// 	Expiry    time.Time `json:"-" yaml:"-" mapstructure:"-"`
// 	CreatedAt time.Time `json:"-" yaml:"-" mapstructure:"-"`
// 	UpdatedAt time.Time `json:"-" yaml:"-" mapstructure:"-"`

// 	Extra map[string]any `json:"-" yaml:"-" mapstructure:"-"`
// }

type AuthOAuthClientConfig struct {
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

type AuthClientPlatform struct {
	ID     string `json:"id,omitempty" env:"AUTH_PLATFORM_ID"`
	Name   string `json:"name,omitempty" env:"AUTH_PLATFORM_NAME"`
	Icon   string `json:"icon,omitempty" env:"AUTH_PLATFORM_ICON"`
	APIURL string `json:"api_url,omitempty" env:"AUTH_PLATFORM_API_URL"`
}

type AuthProvidersConfig struct {
	Google   AuthClientConfig `json:"google,omitempty" env:"GOOGLE_AUTH_CONFIG"`
	Facebook AuthClientConfig `json:"facebook,omitempty" env:"FACEBOOK_AUTH_CONFIG"`
	Github   AuthClientConfig `json:"github,omitempty" env:"GITHUB_AUTH_CONFIG"`
}

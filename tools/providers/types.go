package registry

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type GlobalRef struct {
	ID   uuid.UUID `json:"id,omitempty"`
	Name string    `json:"name,omitempty"`
}

func NewGlobalRef(name string) GlobalRef {
	return GlobalRef{
		ID:   uuid.New(),
		Name: name,
	}
}

func (gr *GlobalRef) GetGlobalRef() GlobalRef { return *gr }
func (gr *GlobalRef) GetName() string         { return gr.Name }
func (gr *GlobalRef) GetID() uuid.UUID        { return gr.ID }
func (gr *GlobalRef) SetName(name string)     { gr.Name = name }
func (gr *GlobalRef) SetID(id uuid.UUID)      { gr.ID = id }
func (gr *GlobalRef) String() string {
	return gr.Name + "-" + gr.ID.String()
}

type ToolCall struct {
	Name string `json:"name"`
	Args any    `json:"args"` // geralmente map[string]any
}

type ChatRequest struct {
	Headers  map[string]string `json:"-"`
	Provider string            `json:"provider"`
	Model    string            `json:"model"`
	Messages []Message         `json:"messages"`
	Temp     float32           `json:"temperature"`
	Stream   bool              `json:"stream"`
	Meta     map[string]any    `json:"meta"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	Completion int     `json:"completion_tokens"`
	Prompt     int     `json:"prompt_tokens"`
	Tokens     int     `json:"tokens"`
	Ms         int64   `json:"latency_ms"`
	CostUSD    float64 `json:"cost_usd"`
	Provider   string  `json:"provider"`
	Model      string  `json:"model"`
}

type ChatChunk struct {
	Content  string    `json:"content,omitempty"`
	Done     bool      `json:"done"`
	Usage    *Usage    `json:"usage,omitempty"`
	Error    string    `json:"error,omitempty"`
	ToolCall *ToolCall `json:"toolCall,omitempty"`
}

type LLMRequestDefaults struct {
	MaxTokens        int     `yaml:"max_tokens,omitempty" json:"max_tokens,omitempty" mapstructure:"max_tokens,omitempty"`
	Temperature      float64 `yaml:"temperature,omitempty" json:"temperature,omitempty" mapstructure:"temperature,omitempty"`
	TopP             float64 `yaml:"top_p,omitempty" json:"top_p,omitempty" mapstructure:"top_p,omitempty"`
	FrequencyPenalty float64 `yaml:"frequency_penalty,omitempty" json:"frequency_penalty,omitempty" mapstructure:"frequency_penalty,omitempty"`
	PresencePenalty  float64 `yaml:"presence_penalty,omitempty" json:"presence_penalty,omitempty" mapstructure:"presence_penalty,omitempty"`
	Stream           bool    `yaml:"stream,omitempty" json:"stream,omitempty" mapstructure:"stream,omitempty"`
	TimeoutSec       int     `yaml:"timeout_sec,omitempty" json:"timeout_sec,omitempty" mapstructure:"timeout_sec,omitempty"`
	TenantID         string  `yaml:"tenant_id,omitempty" json:"tenant_id,omitempty" mapstructure:"tenant_id,omitempty"`
	UserID           string  `yaml:"user_id,omitempty" json:"user_id,omitempty" mapstructure:"user_id,omitempty"`
}

type LLMTokenBucket struct {
	Capacity   int `yaml:"capacity,omitempty" json:"capacity,omitempty" mapstructure:"capacity,omitempty"`
	RefillRate int `yaml:"refill_rate,omitempty" json:"refill_rate,omitempty" mapstructure:"refill_rate,omitempty"`
}

type LLMRateLimitConfig struct {
	Enabled     bool                      `yaml:"enabled,omitempty" json:"enabled,omitempty" mapstructure:"enabled,omitempty"`
	Default     LLMTokenBucket            `yaml:"default,omitempty" json:"default,omitempty" mapstructure:"default,omitempty"`
	PerProvider map[string]LLMTokenBucket `yaml:"per_provider,omitempty" json:"per_provider,omitempty" mapstructure:"per_provider,omitempty"`
}

type LLMCircuitBreakerRule struct {
	MaxFailures      int `yaml:"max_failures,omitempty" json:"max_failures,omitempty" mapstructure:"max_failures,omitempty"`
	ResetTimeoutSec  int `yaml:"reset_timeout_sec,omitempty" json:"reset_timeout_sec,omitempty" mapstructure:"reset_timeout_sec,omitempty"`
	SuccessThreshold int `yaml:"success_threshold,omitempty" json:"success_threshold,omitempty" mapstructure:"success_threshold,omitempty"`
}

type LLMCircuitBreakerConfig struct {
	Enabled     bool                             `yaml:"enabled,omitempty" json:"enabled,omitempty" mapstructure:"enabled,omitempty"`
	Default     LLMCircuitBreakerRule            `yaml:"default,omitempty" json:"default,omitempty" mapstructure:"default,omitempty"`
	PerProvider map[string]LLMCircuitBreakerRule `yaml:"per_provider,omitempty" json:"per_provider,omitempty" mapstructure:"per_provider,omitempty"`
}

type LLMHealthCheckConfig struct {
	Enabled     bool `yaml:"enabled,omitempty" json:"enabled,omitempty" mapstructure:"enabled,omitempty"`
	IntervalSec int  `yaml:"interval_sec,omitempty" json:"interval_sec,omitempty" mapstructure:"interval_sec,omitempty"`
	TimeoutSec  int  `yaml:"timeout_sec,omitempty" json:"timeout_sec,omitempty" mapstructure:"timeout_sec,omitempty"`
}

type LLMRetryConfig struct {
	Enabled     bool    `yaml:"enabled,omitempty" json:"enabled,omitempty" mapstructure:"enabled,omitempty"`
	MaxRetries  int     `yaml:"max_retries,omitempty" json:"max_retries,omitempty" mapstructure:"max_retries,omitempty"`
	BaseDelayMS int     `yaml:"base_delay_ms,omitempty" json:"base_delay_ms,omitempty" mapstructure:"base_delay_ms,omitempty"`
	MaxDelayMS  int     `yaml:"max_delay_ms,omitempty" json:"max_delay_ms,omitempty" mapstructure:"max_delay_ms,omitempty"`
	Multiplier  float64 `yaml:"multiplier,omitempty" json:"multiplier,omitempty" mapstructure:"multiplier,omitempty"`
}

type LLMDevelopmentConfig struct {
	LoggingLevel   string                  `yaml:"logging_level,omitempty" json:"logging_level,omitempty" mapstructure:"logging_level,omitempty"`
	Defaults       LLMRequestDefaults      `yaml:"defaults,omitempty" json:"defaults,omitempty" mapstructure:"defaults,omitempty"`
	RateLimit      LLMRateLimitConfig      `yaml:"rate_limit,omitempty" json:"rate_limit,omitempty" mapstructure:"rate_limit,omitempty"`
	CircuitBreaker LLMCircuitBreakerConfig `yaml:"circuit_breaker,omitempty" json:"circuit_breaker,omitempty" mapstructure:"circuit_breaker,omitempty"`
	HealthCheck    LLMHealthCheckConfig    `yaml:"health_check,omitempty" json:"health_check,omitempty" mapstructure:"health_check,omitempty"`
	Retry          LLMRetryConfig          `yaml:"retry,omitempty" json:"retry,omitempty" mapstructure:"retry,omitempty"`
}

type LLMProviderProductionConfig struct {
	TimeoutSec  int     `yaml:"timeout_sec,omitempty" json:"timeout_sec,omitempty" mapstructure:"timeout_sec,omitempty"`
	Priority    string  `yaml:"priority,omitempty" json:"priority,omitempty" mapstructure:"priority,omitempty"`
	MaxRetries  int     `yaml:"max_retries,omitempty" json:"max_retries,omitempty" mapstructure:"max_retries,omitempty"`
	BaseDelayMS int     `yaml:"base_delay_ms,omitempty" json:"base_delay_ms,omitempty" mapstructure:"base_delay_ms,omitempty"`
	MaxDelayMS  int     `yaml:"max_delay_ms,omitempty" json:"max_delay_ms,omitempty" mapstructure:"max_delay_ms,omitempty"`
	Multiplier  float64 `yaml:"multiplier,omitempty" json:"multiplier,omitempty" mapstructure:"multiplier,omitempty"`
}

type LLMSecurityConfig struct {
	EnableHTTPS    bool     `yaml:"enable_https,omitempty" json:"enable_https,omitempty" mapstructure:"enable_https,omitempty"`
	AllowedOrigins []string `yaml:"allowed_origins,omitempty" json:"allowed_origins,omitempty" mapstructure:"allowed_origins,omitempty"`
	JWTSecret      string   `yaml:"jwt_secret,omitempty" json:"jwt_secret,omitempty" mapstructure:"jwt_secret,omitempty"`
	APIKeys        []string `yaml:"api_keys,omitempty" json:"api_keys,omitempty" mapstructure:"api_keys,omitempty"`
}

type LLMMonitoringConfig struct {
	EnableMetrics bool `yaml:"enable_metrics,omitempty" json:"enable_metrics,omitempty" mapstructure:"enable_metrics,omitempty"`
}

type LLMProvidersExtMap map[string]ProviderExt

type LLMProvidersMap map[string]*LLMProviderConfig

type LLMConfig struct {
	GlobalRef          `json:",inline" yaml:",inline" mapstructure:",squash"`
	FilePath           string                                 `json:"file_path,omitempty" yaml:"file_path,omitempty" mapstructure:"file_path,omitempty"`
	Development        LLMDevelopmentConfig                   `yaml:"development,omitempty" json:"development,omitempty" mapstructure:"development,omitempty"`
	Providers          LLMProvidersMap                        `yaml:"providers,omitempty" json:"providers,omitempty" mapstructure:"providers,omitempty"`
	ProviderProduction map[string]LLMProviderProductionConfig `yaml:"provider_production,omitempty" json:"provider_production,omitempty" mapstructure:"provider_production,omitempty"`
	Security           LLMSecurityConfig                      `yaml:"security,omitempty" json:"security,omitempty" mapstructure:"security,omitempty"`
	Monitoring         LLMMonitoringConfig                    `yaml:"monitoring,omitempty" json:"monitoring,omitempty" mapstructure:"monitoring,omitempty"`
	Repository         string                                 `yaml:"repository,omitempty" json:"repository,omitempty" mapstructure:"repository,omitempty"`
	Version            string                                 `yaml:"version,omitempty" json:"version,omitempty" mapstructure:"version,omitempty"`
	Authors            []string                               `yaml:"authors,omitempty" json:"authors,omitempty" mapstructure:"authors,omitempty"`
	License            string                                 `yaml:"license,omitempty" json:"license,omitempty" mapstructure:"license,omitempty"`
}

type LLMProviderConfig struct {
	name         string `json:"-"`
	typ          string `json:"-"`
	BaseURL      string `yaml:"base_url,omitempty" json:"base_url,omitempty" mapstructure:"base_url,omitempty"`
	KeyEnv       string `yaml:"key_env,omitempty" json:"key_env,omitempty" mapstructure:"key_env,omitempty"`
	DefaultModel string `yaml:"default_model,omitempty" json:"default_model,omitempty" mapstructure:"default_model,omitempty"`
}

type Provider interface {
	Name() string
	Type() string
	Chat(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error)
	Available() error

	Notify(ctx context.Context, event NotificationEvent) error
}

type ProviderExt interface {
	Provider

	KeyRef() string
	URLBase() string

	ListModels(ctx context.Context) ([]string, error)

	ModelInfo(ctx context.Context) (map[string]any, error)
	SetModel(ctx context.Context, model string) error

	HealthCheck(ctx context.Context) error
}

type NotificationEvent struct {
	Type      string         `json:"type"` // "discord", "whatsapp", "email"
	Recipient string         `json:"recipient"`
	Subject   string         `json:"subject"`
	Content   string         `json:"content"`
	Priority  string         `json:"priority"` // "low", "medium", "high", "critical"
	Metadata  map[string]any `json:"metadata"`
	CreatedAt time.Time      `json:"created_at"`
}

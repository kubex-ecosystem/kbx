package types

import (
	"context"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/kubex-ecosystem/kbx/internal/module/kbx"
	load "github.com/kubex-ecosystem/kbx/tools"

	gl "github.com/kubex-ecosystem/logz"
)

type ToolCall struct {
	Name string `json:"name"`
	Args any    `json:"args"` // geralmente map[string]any
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Headers  map[string]string `json:"-"`
	Provider string            `json:"provider"`
	Model    string            `json:"model"`
	Messages []Message         `json:"messages"`
	Temp     float32           `json:"temperature"`
	Stream   bool              `json:"stream"`
	Meta     map[string]any    `json:"meta"`
}

func (r ChatRequest) Validate() error {
	if strings.TrimSpace(r.Provider) == "" {
		return gl.Error("Provider is required")
	}
	return nil
}

func (r ChatRequest) GetModel() string { return r.Model }

func (r ChatRequest) Read(ctx context.Context) (ChatChunk, error) {
	if r.Stream {
		return ChatChunk{}, gl.Error("streaming not implemented in this method")
	}

	var cnk ChatChunk
	if err := r.Validate(); err != nil {
		return cnk, err
	}

	var p ProviderExt
	p, err := getLLMProviderByName(nil, r.Provider)
	if err != nil {
		return cnk, gl.Errorf("failed to get provider '%s': %v", r.Provider, err)
	}
	resp, err := p.Chat(
		ctx,
		ChatRequest{
			Provider: r.Provider,
			Model:    r.Model,
			Messages: r.Messages,
			Temp:     r.Temp,
			Stream:   r.Stream,
			Meta:     r.Meta,
		},
	)
	if err != nil {
		return cnk, err
	}

	cnk = <-resp

	return cnk, nil
}

// Message represents a single chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Usage represents token usage and cost information
type Usage struct {
	Completion int     `json:"completion_tokens"`
	Prompt     int     `json:"prompt_tokens"`
	Tokens     int     `json:"tokens"`
	Ms         int64   `json:"latency_ms"`
	CostUSD    float64 `json:"cost_usd"`
	Provider   string  `json:"provider"`
	Model      string  `json:"model"`
}

// ChatChunk represents a streaming response chunk
type ChatChunk struct {
	Content  string    `json:"content,omitempty"`
	Done     bool      `json:"done"`
	Usage    *Usage    `json:"usage,omitempty"`
	Error    string    `json:"error,omitempty"`
	ToolCall *ToolCall `json:"toolCall,omitempty"`
}

func (c ChatChunk) IsSuccess() bool   { return c.Error == "" }
func (c ChatChunk) IsError() bool     { return c.Error != "" }
func (c ChatChunk) IsDone() bool      { return c.Done }
func (c ChatChunk) HasContent() bool  { return len(c.Content) > 0 }
func (c ChatChunk) HasToolCall() bool { return c.ToolCall != nil }

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

func (m LLMProvidersMap) GetProvider(name string) (ProviderExt, bool) {
	p, ok := m[name]
	return p, ok
}

func (m LLMProvidersMap) GetAllProviders() LLMProvidersExtMap {
	pMap := make(LLMProvidersExtMap, len(m))
	for k, v := range m {
		pMap[k] = v
	}
	return pMap
}

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

func NewLLMConfig(path string, name string, version string, p map[string]*LLMProviderConfig) LLMConfig {
	if len(name) == 0 {
		name = "kubex-llm-provider-config"
	}
	if len(version) == 0 {
		version = "1.0.0"
	}
	if len(path) == 0 {
		path = getDefaultLLMConfigPath()
	}
	return LLMConfig{
		GlobalRef:          NewGlobalRef(name),
		FilePath:           path,
		Version:            version,
		Providers:          p,
		Development:        LLMDevelopmentConfig{},
		ProviderProduction: map[string]LLMProviderProductionConfig{},
		Security:           LLMSecurityConfig{},
		Monitoring:         LLMMonitoringConfig{},
		Authors:            []string{},
	}
}

func NewLLMConfigDefault() LLMConfig {
	cfg := NewLLMConfig(getDefaultLLMConfigPath(), "kubex-llm-provider-config", "1.0.0", nil)
	cfg.Name = "kubex-llm-provider-config"
	cfg.Providers = make(map[string]*LLMProviderConfig, 0)
	pvds := map[string]*LLMProviderConfig{
		"groq": &LLMProviderConfig{
			name:         "groq",
			typ:          "groq",
			BaseURL:      "https://api.groq.com",
			KeyEnv:       "GROQ_API_KEY",
			DefaultModel: "",
		},
		"gemini": &LLMProviderConfig{
			name:         "gemini",
			typ:          "gemini",
			BaseURL:      "https://generativelanguage.googleapis.com",
			KeyEnv:       "GEMINI_API_KEY",
			DefaultModel: "",
		},
		"openai": &LLMProviderConfig{
			name:         "openai",
			typ:          "openai",
			BaseURL:      "https://api.openai.com",
			KeyEnv:       "OPENAI_API_KEY",
			DefaultModel: "",
		},
	}
	for k, v := range pvds {
		cfg.Providers[k] = v
	}
	cfg.Development = LLMDevelopmentConfig{
		LoggingLevel: "DEBUG",
		Defaults: LLMRequestDefaults{
			MaxTokens:        2048,
			Temperature:      0.7,
			TopP:             0.9,
			FrequencyPenalty: 0.0,
			PresencePenalty:  0.0,
			Stream:           false,
			TimeoutSec:       30,
			TenantID:         "default",
			UserID:           "anonymous",
		},
		RateLimit: LLMRateLimitConfig{
			Enabled: true,
			Default: LLMTokenBucket{
				Capacity:   100,
				RefillRate: 10,
			},
			PerProvider: map[string]LLMTokenBucket{
				"groq":   {Capacity: 200, RefillRate: 20},
				"gemini": {Capacity: 50, RefillRate: 8},
				"openai": {Capacity: 30, RefillRate: 5},
			},
		},
		CircuitBreaker: LLMCircuitBreakerConfig{
			Enabled: true,
			Default: LLMCircuitBreakerRule{
				MaxFailures:      5,
				ResetTimeoutSec:  60,
				SuccessThreshold: 3,
			},
			PerProvider: map[string]LLMCircuitBreakerRule{
				"groq":   {MaxFailures: 3, ResetTimeoutSec: 30, SuccessThreshold: 2},
				"gemini": {MaxFailures: 5, ResetTimeoutSec: 60, SuccessThreshold: 3},
				"openai": {MaxFailures: 4, ResetTimeoutSec: 90, SuccessThreshold: 3},
			},
		},
		HealthCheck: LLMHealthCheckConfig{
			Enabled:     true,
			IntervalSec: 30,
			TimeoutSec:  10,
		},
		Retry: LLMRetryConfig{
			Enabled:     true,
			MaxRetries:  3,
			BaseDelayMS: 100,
			MaxDelayMS:  5000,
			Multiplier:  2.0,
		},
	}
	cfg.ProviderProduction = map[string]LLMProviderProductionConfig{
		"groq": {
			TimeoutSec:  30,
			Priority:    "high",
			MaxRetries:  3,
			BaseDelayMS: 50,
			MaxDelayMS:  3000,
			Multiplier:  1.5,
		},
		"gemini": {
			TimeoutSec:  60,
			Priority:    "medium",
			MaxRetries:  3,
			BaseDelayMS: 100,
			MaxDelayMS:  5000,
			Multiplier:  2.0,
		},
		"openai": {
			TimeoutSec:  120,
			Priority:    "medium",
			MaxRetries:  4,
			BaseDelayMS: 200,
			MaxDelayMS:  6000,
			Multiplier:  2.5,
		},
	}
	cfg.Security = LLMSecurityConfig{
		EnableHTTPS:    true,
		AllowedOrigins: []string{"https://kubex.world"},
		JWTSecret:      "",
		APIKeys:        []string{},
	}
	cfg.Monitoring = LLMMonitoringConfig{
		EnableMetrics: true,
	}
	cfg.Repository = "https://github.com/kubex-ecosystem/kubex-gemx-gnyx"
	cfg.Version = "1.0.0"
	cfg.Authors = []string{"Kubex Dev Team <dev@kubex.world>"}
	cfg.License = "MIT"

	return cfg
}

func (cfg *LLMConfig) GetProvider(name string) (ProviderExt, bool) {
	p, ok := cfg.Providers[name]
	return p, ok
}

func (cfg *LLMConfig) GetProviders() LLMProvidersExtMap {
	pMap := make(LLMProvidersExtMap, len(cfg.Providers))
	for k, v := range cfg.Providers {
		pMap[k] = v
	}
	return pMap
}

func (cfg *LLMConfig) Validate() error {
	for name, provider := range cfg.Providers {
		if err := provider.Available(); err != nil {
			return gl.Errorf("provider '%s' is not available: %v", name, err)
		}
	}
	return nil
}

func (cfg *LLMConfig) GetCurrentProvider() (ProviderExt, error) {
	if !kbx.IsObjValid(cfg.Providers) || len(cfg.Providers) == 0 {
		cfg.Providers = make(map[string]*LLMProviderConfig, 0)
		gl.Warn("no providers configured, using defaults")
		pp := NewLLMConfig(
			os.ExpandEnv(
				kbx.GetEnvOrDefaultWithType(
					"LLM_CONFIG_PATH",
					kbx.GetValueOrDefaultSimple(
						cfg.FilePath,
						getDefaultLLMConfigPath(),
					),
				),
			), "kubex-llm-provider-config", "1.0.0", nil)
		cfg.Providers = pp.Providers
	}
	if len(cfg.Providers) == 0 {
		return nil, gl.Errorf("no providers configured")
	} else {
		for name, provider := range cfg.Providers {
			if err := provider.Available(); err == nil {
				return provider, nil
			} else {
				gl.Warnf("provider '%s' is not available: %v", name, err)
			}
		}
	}
	return nil, gl.Errorf("no available providers found in configuration")
}

func (cfg *LLMConfig) SetProvider(name string, provider ProviderExt) error {
	if cfg.Providers == nil {
		cfg.Providers = make(map[string]*LLMProviderConfig)
	}
	if provider == nil {
		return gl.Errorf("provider cannot be nil")
	}
	if len(name) == 0 {
		return gl.Errorf("provider name cannot be empty")
	}
	if _, exists := cfg.Providers[name]; !exists {
		return gl.Errorf("provider with name '%s' does not exist", name)
	}
	if p, ok := cfg.Providers[name]; ok || p != nil {
		p = &LLMProviderConfig{
			name:         kbx.GetValueOrDefaultSimple(provider.Name(), p.name),
			typ:          kbx.GetValueOrDefaultSimple(provider.Type(), p.typ),
			BaseURL:      kbx.GetValueOrDefaultSimple(provider.URLBase(), p.BaseURL),
			KeyEnv:       kbx.GetValueOrDefaultSimple(provider.KeyRef(), p.KeyEnv),
			DefaultModel: kbx.GetValueOrDefaultSimple("", p.DefaultModel), // DefaultModel is not part of ProviderExt, so we keep existing value
		}
		cfg.Providers[name] = p
	}
	return nil
}

func (cfg *LLMConfig) AddProvider(name string, provider ProviderExt) error {
	if cfg.Providers == nil {
		cfg.Providers = make(map[string]*LLMProviderConfig)
	}
	if _, exists := cfg.Providers[name]; exists {
		return gl.Errorf("provider with name '%s' already exists", name)
	}
	p := &LLMProviderConfig{
		name:         provider.Name(),
		typ:          provider.Type(),
		BaseURL:      provider.URLBase(),
		KeyEnv:       provider.KeyRef(),
		DefaultModel: "",
	}
	cfg.Providers[name] = p
	return nil
}

func (cfg *LLMConfig) RemoveProvider(name string) error {
	if cfg.Providers == nil {
		return gl.Errorf("no providers configured")
	}
	if _, exists := cfg.Providers[name]; !exists {
		return gl.Errorf("provider with name '%s' does not exist", name)
	}
	delete(cfg.Providers, name)
	return nil
}

// Provider interface defines the contract for AI providers
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

type LLMProviderConfig struct {
	name         string `json:"-"`
	typ          string `json:"-"`
	BaseURL      string `yaml:"base_url,omitempty" json:"base_url,omitempty" mapstructure:"base_url,omitempty"`
	KeyEnv       string `yaml:"key_env,omitempty" json:"key_env,omitempty" mapstructure:"key_env,omitempty"`
	DefaultModel string `yaml:"default_model,omitempty" json:"default_model,omitempty" mapstructure:"default_model,omitempty"`
}

// NewLLMProviderConfigType exports concrete implementation of Provider interface for LLMProviderConfig to be used with caution
// This is a simple implementation to allow LLMProviderConfig to satisfy the Provider interface.
// Each provider type (e.g., GroqProvider, GeminiProvider) should ideally have its own struct implementing Provider
// with specific logic for Chat, Available, and Notify methods. This is just a placeholder to allow the config struct
// to be used as fallback/default provider implementation, increasing coupling but resilience also.
func NewLLMProviderConfigType(name, baseurl, keyenv, defaultmodel string) *LLMProviderConfig {
	return &LLMProviderConfig{
		name:         name,
		typ:          strings.ToLower(name),
		BaseURL:      baseurl,
		KeyEnv:       keyenv,
		DefaultModel: defaultmodel,
	}
}

// NewLLMProviderConfig is a helper function to create a new LLMProviderConfig with the given parameters.
// This can be used in tests or when programmatically creating provider configs.
func NewLLMProviderConfig(name, baseurl, keyenv, defaultmodel string) Provider {
	return NewLLMProviderConfigType(name, baseurl, keyenv, defaultmodel)
}

// NewLLMProviderConfigExt is a helper function to create a new LLMProviderConfig and return it as a ProviderExt interface.
func NewLLMProviderConfigExt(name, baseurl, keyenv, defaultmodel string) ProviderExt {
	return NewLLMProviderConfigType(name, baseurl, keyenv, defaultmodel)
}

func (pc *LLMProviderConfig) Name() string    { return pc.name }
func (pc *LLMProviderConfig) Type() string    { return pc.typ }
func (pc *LLMProviderConfig) URLBase() string { return pc.BaseURL }
func (pc *LLMProviderConfig) Available() error {
	if pc.BaseURL == "" || pc.KeyEnv == "" {
		return gl.Errorf("provider '%s' is not properly configured", pc.typ)
	}
	if _, ok := os.LookupEnv(pc.KeyEnv); !ok {
		return gl.Errorf("environment variable '%s' for provider '%s' is not set", pc.KeyEnv, pc.typ)
	}
	if pc.DefaultModel == "" {
		return gl.Errorf("provider '%s' does not have a default model configured", pc.typ)
	}
	if _, err := url.Parse(pc.BaseURL); err != nil {
		return gl.Errorf("provider '%s' has an invalid base URL '%s': %v", pc.typ, pc.BaseURL, err)
	}
	return nil
}
func (pc *LLMProviderConfig) Chat(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error) {
	ch := make(chan ChatChunk)
	go func() {
		defer close(ch)
		// Simulate a response for demonstration purposes
		ch <- ChatChunk{Content: "Hello from " + pc.typ + "!", Done: true}
	}()
	return ch, nil
}
func (pc *LLMProviderConfig) Notify(ctx context.Context, event NotificationEvent) error {
	// Basic generic implementation
	switch event.Type {
	case "rate_limit_exceeded":
		gl.Warnf("Provider '%s' has exceeded its rate limit", pc.typ)
	case "provider_error":
		gl.Errorf("Provider '%s' encountered an error: %v", pc.typ, event.Content)
	default:
		switch event.Type {
		case "error":
			gl.Errorf("Provider '%s' error: %v", pc.typ, event.Content)
		case "warning":
			gl.Warnf("Provider '%s' warning: %v", pc.typ, event.Content)
		case "info":
			gl.Infof("Provider '%s' info: %v", pc.typ, event.Content)
		default:
			switch event.Subject {
			case "chat":
				gl.Infof("Provider '%s' chat event: %v", pc.typ, event.Content)
			case "tool_call":
				gl.Infof("Provider '%s' tool call event: %v", pc.typ, event.Content)
			default:
				gl.Infof("Provider '%s' event: %v", pc.typ, event.Content)
			}
		}
	}

	return nil
}
func (pc *LLMProviderConfig) KeyRef() string {
	return os.ExpandEnv(pc.KeyEnv)
}
func (pc *LLMProviderConfig) ListModels(ctx context.Context) ([]string, error) {
	// This is a placeholder implementation. In a real implementation, this would make an API call to the provider to retrieve available models.
	return []string{pc.DefaultModel}, nil
}
func (pc *LLMProviderConfig) ModelInfo(ctx context.Context) (map[string]any, error) {
	// Placeholder implementation. In a real implementation, this would retrieve detailed information about the current model from the provider.
	return map[string]any{
		"name":        pc.DefaultModel,
		"description": "Default model for " + pc.typ,
		"max_tokens":  2048,
	}, nil
}
func (pc *LLMProviderConfig) SetModel(ctx context.Context, model string) error {
	// Placeholder implementation. In a real implementation, this might validate the model against the provider's available models and set it for future requests.
	pc.DefaultModel = model
	return nil
}
func (pc *LLMProviderConfig) HealthCheck(ctx context.Context) error {
	// Placeholder implementation. In a real implementation, this would make a lightweight API call to the provider to check its health/status.
	if err := pc.Available(); err != nil {
		return gl.Errorf("health check failed for provider '%s': %v", pc.typ, err)
	}
	return nil
}

func getDefaultLLMConfigPath() string {
	p1 := os.Getenv("KUBEX_LLM_PROVIDER_CONFIG_PATH")
	if len(p1) == 0 {
		p1 = os.ExpandEnv(kbx.DefaultConfigFile)
	}
	p1 = strings.TrimSpace(filepath.Clean(strings.ToValidUTF8(p1, "")))
	if len(p1) == 0 {
		p1 = os.ExpandEnv(kbx.DefaultConfigFile)
	}
	return p1
}

func getLLMProviderByName(cfg *LLMConfig, name string) (ProviderExt, error) {
	if cfg.Providers == nil {
		cfg.Providers = make(map[string]*LLMProviderConfig)
	}
	if p, ok := cfg.Providers[name]; ok && p != nil {
		return p, nil
	} else {
		return nil, gl.Errorf("provider with name '%s' not found in configuration", name)
	}
}

func loadLLMConfigFromFile(path string) (LLMConfig, error) {
	cfg := NewLLMConfigDefault()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		gl.Warnf("LLM config file '%s' does not exist, using defaults", path)
		return cfg, nil
	}
	pc := load.NewEmptyMapperType[LLMProviderConfig](path)
	pcd, err := pc.DeserializeFromFile(filepath.Ext(path)[1:])
	if err != nil {
		return cfg, gl.Errorf("failed to load LLM config from file '%s': %v", path, err)
	}
	cfg.Providers[pcd.Name()] = pcd
	return cfg, nil
}

func getLLMConfig(cfg *LLMConfig) (*LLMConfig, error) {
	if cfg == nil {
		cfg = &LLMConfig{}
	}
	if len(cfg.FilePath) == 0 {
		cfg.FilePath = getDefaultLLMConfigPath()
	}
	loadedCfg, err := loadLLMConfigFromFile(cfg.FilePath)
	if err != nil {
		return &loadedCfg, err
	}
	// Merge loaded config with defaults, giving precedence to loaded values
	if loadedCfg.Providers != nil && len(loadedCfg.Providers) > 0 {
		cfg.Providers = loadedCfg.Providers
	} else if cfg.Providers == nil || len(cfg.Providers) == 0 {
		cfg.Providers = make(map[string]*LLMProviderConfig, 0)
		gl.Warn("no providers found in loaded config, using defaults")
		pp := NewLLMConfig(
			os.ExpandEnv(
				kbx.GetEnvOrDefaultWithType(
					"LLM_CONFIG_PATH",
					kbx.GetValueOrDefaultSimple(
						cfg.FilePath,
						getDefaultLLMConfigPath(),
					),
				),
			), "kubex-llm-provider-config", "1.0.0", nil)
		cfg.Providers = pp.Providers
	}
	return cfg, nil
}

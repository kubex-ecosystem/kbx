// Package registry provides provider registration and resolution functionality.
package registry

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	kbx "github.com/kubex-ecosystem/kbx"
	kbxMod "github.com/kubex-ecosystem/kbx/internal/module/kbx"
	kbxTypes "github.com/kubex-ecosystem/kbx/types"

	gl "github.com/kubex-ecosystem/logz"
)

// Registry manages provider registration, configuration, and runtime resolution.
type Registry struct {
	cfg       *kbxTypes.LLMConfig
	providers map[string]kbxTypes.ProviderExt
}

// -------------------------------- REGISTRY CONSTRUCTORS --------------------------------

func NewRegistry(c *kbxTypes.LLMConfig) *Registry {
	cfg := c
	if cfg == nil {
		defaultCfg := kbxTypes.NewLLMConfigDefault()
		cfg = &defaultCfg
	}
	if cfg.Providers == nil {
		cfg.Providers = make(kbxTypes.LLMProvidersMap)
	}
	return &Registry{
		cfg:       cfg,
		providers: make(map[string]kbxTypes.ProviderExt, len(cfg.Providers)),
	}
}

func Load(path string) (*Registry, error) {
	path = strings.TrimSpace(filepath.Clean(strings.ToValidUTF8(path, "")))
	if path == "" {
		return nil, gl.Errorf("provider config path cannot be empty")
	}

	gl.Debugf("Loading provider configuration from %s", path)

	loadedCfg, err := kbx.LoadConfigOrDefault[kbxTypes.LLMConfig](path, true)
	if err != nil {
		gl.Errorf("Failed to load config: %v", err)
		return nil, gl.Errorf("failed to load provider config: %w", err)
	}

	cfg := buildRuntimeConfig(path, loadedCfg)
	rg := NewRegistry(&cfg)
	rg.instantiateProviders()

	return rg, nil
}

// -------------------------------- REGISTRY GENERAL METHODS --------------------------------

func (r *Registry) Config() kbxTypes.LLMConfig {
	if r == nil || r.cfg == nil {
		gl.Warn("Provider registry config is nil. Returning default config.")
		return kbxTypes.NewLLMConfigDefault()
	}
	return *r.cfg
}

func (r *Registry) GetProviderConfig(name string) *kbxTypes.LLMProviderConfig {
	if r == nil || r.cfg == nil || r.cfg.Providers == nil {
		return nil
	}
	return r.cfg.Providers[normalizeProviderName(name)]
}

func (r *Registry) Providers() kbxTypes.LLMProvidersExtMap {
	providers := make(map[string]kbxTypes.ProviderExt, len(r.providers))
	for name, provider := range r.providers {
		providers[name] = provider
	}
	return providers
}

func (r *Registry) ListProviders() []string {
	if r == nil || len(r.providers) == 0 {
		return []string{}
	}
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// -------------------------------- PROVIDER INTERFACE IMPLEMENTATION --------------------------------

func (r *Registry) Resolve(name string) kbxTypes.Provider {
	return r.ResolveProvider(name)
}

func (r *Registry) ResolveProvider(name string) kbxTypes.ProviderExt {
	if r == nil || len(r.providers) == 0 {
		gl.Warn("Provider registry is empty. No providers available for resolution.")
		return nil
	}
	provider, ok := r.providers[normalizeProviderName(name)]
	if !ok {
		gl.Warnf("Provider '%s' not found in registry.", name)
		return nil
	}
	return provider
}

func (r *Registry) Chat(ctx context.Context, req kbxTypes.ChatRequest) (<-chan kbxTypes.ChatChunk, error) {
	p := r.ResolveProvider(req.Provider)
	if p == nil {
		return nil, gl.Errorf("provider '%s' not found", req.Provider)
	}
	return p.Chat(ctx, req)
}

func (r *Registry) Notify(ctx context.Context, event kbxTypes.NotificationEvent) error {
	p := r.ResolveProvider(event.Type)
	if p == nil {
		return gl.Errorf("provider '%s' not found", event.Type)
	}
	return p.Notify(ctx, event)
}

// -------------------------------- PRIVATE INTERNAL METHODS --------------------------------

func (r *Registry) instantiateProviders() {
	if r == nil || r.cfg == nil {
		return
	}
	if r.providers == nil {
		r.providers = make(map[string]kbxTypes.ProviderExt)
	}

	for name, pc := range r.cfg.Providers {
		providerType := normalizeProviderType(name, pc)
		constructor, ok := providerConstructors[providerType]
		if !ok {
			gl.Warnf("Skipping provider '%s' - unsupported type '%s'", name, providerType)
			continue
		}

		key := resolveAPIKey(name, pc)
		if key == "" {
			gl.Warnf("Skipping provider '%s' - no API key found in %s", name, pc.KeyEnv)
			continue
		}

		provider, err := constructor(name, strings.TrimSpace(pc.BaseURL), key, strings.TrimSpace(pc.DefaultModel))
		if err != nil {
			gl.Warnf("Failed to initialize provider '%s': %v. This provider will be unavailable for use.", name, err)
			continue
		}

		r.providers[name] = provider

		if info, err := provider.ModelInfo(context.Background()); err == nil {
			if modelName, _ := info["name"].(string); modelName != "" {
				gl.Debugf("Provider '%s' model info: %v", name, info)
			}
		}
	}
}

var providerConstructors = map[string]func(name, baseURL, key, model string) (kbxTypes.ProviderExt, error){
	"openai":    NewOpenAIProvider,
	"gemini":    NewGeminiProvider,
	"anthropic": NewAnthropicProvider,
	"groq":      NewGroqProvider,
}

func buildRuntimeConfig(path string, loaded *kbxTypes.LLMConfig) kbxTypes.LLMConfig {
	cfg := kbxTypes.NewLLMConfigDefault()
	cfg.FilePath = path

	if loaded == nil {
		return cfg
	}

	if loaded.Name != "" {
		cfg.GlobalRef = loaded.GlobalRef
	}
	cfg.Development = loaded.Development
	if loaded.ProviderProduction != nil {
		cfg.ProviderProduction = loaded.ProviderProduction
	}
	cfg.Security = loaded.Security
	cfg.Monitoring = loaded.Monitoring
	if loaded.Repository != "" {
		cfg.Repository = loaded.Repository
	}
	if loaded.Version != "" {
		cfg.Version = loaded.Version
	}
	if len(loaded.Authors) > 0 {
		cfg.Authors = loaded.Authors
	}
	if loaded.License != "" {
		cfg.License = loaded.License
	}

	if len(loaded.Providers) > 0 {
		cfg.Providers = make(kbxTypes.LLMProvidersMap, len(loaded.Providers))
		for rawName, providerCfg := range loaded.Providers {
			name := normalizeProviderName(rawName)
			cfg.Providers[name] = normalizeProviderConfig(name, providerCfg, cfg.Providers[name])
		}
	}

	return cfg
}

func normalizeProviderConfig(name string, providerCfg *kbxTypes.LLMProviderConfig, fallback *kbxTypes.LLMProviderConfig) *kbxTypes.LLMProviderConfig {
	baseURL := ""
	keyEnv := ""
	defaultModel := ""

	if fallback != nil {
		baseURL = strings.TrimSpace(fallback.BaseURL)
		keyEnv = strings.TrimSpace(fallback.KeyEnv)
		defaultModel = strings.TrimSpace(fallback.DefaultModel)
	}
	if providerCfg != nil {
		if strings.TrimSpace(providerCfg.BaseURL) != "" {
			baseURL = strings.TrimSpace(providerCfg.BaseURL)
		}
		if strings.TrimSpace(providerCfg.KeyEnv) != "" {
			keyEnv = strings.TrimSpace(providerCfg.KeyEnv)
		}
		if strings.TrimSpace(providerCfg.DefaultModel) != "" {
			defaultModel = strings.TrimSpace(providerCfg.DefaultModel)
		}
	}

	return kbxTypes.NewLLMProviderConfigType(name, baseURL, keyEnv, defaultModel)
}

func normalizeProviderName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func normalizeProviderType(name string, providerCfg *kbxTypes.LLMProviderConfig) string {
	if providerCfg == nil {
		return normalizeProviderName(name)
	}
	if providerType := strings.TrimSpace(providerCfg.Type()); providerType != "" {
		return strings.ToLower(providerType)
	}
	return normalizeProviderName(name)
}

func resolveAPIKey(name string, providerCfg *kbxTypes.LLMProviderConfig) string {
	for _, candidate := range apiKeyCandidates(name, providerCfg) {
		if value := resolveCandidateValue(candidate); value != "" {
			return value
		}
	}
	return ""
}

func apiKeyCandidates(name string, providerCfg *kbxTypes.LLMProviderConfig) []string {
	candidates := []string{}
	appendCandidate := func(value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		for _, existing := range candidates {
			if existing == value {
				return
			}
		}
		candidates = append(candidates, value)
	}

	if providerCfg != nil {
		appendCandidate(providerCfg.KeyEnv)
	}
	appendCandidate(strings.ToUpper(normalizeProviderName(name)) + "_API_KEY")
	appendCandidate(defaultKeyEnv(normalizeProviderType(name, providerCfg)))

	return candidates
}

func resolveCandidateValue(candidate string) string {
	candidate = strings.TrimSpace(candidate)
	if candidate == "" {
		return ""
	}

	if expanded := strings.TrimSpace(os.ExpandEnv(candidate)); expanded != "" && expanded != candidate {
		return expanded
	}
	if value := strings.TrimSpace(os.Getenv(candidate)); value != "" {
		return value
	}
	if !looksLikeEnvName(candidate) {
		return candidate
	}

	return ""
}

func looksLikeEnvName(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		switch {
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '_':
		default:
			return false
		}
	}
	return true
}

func defaultKeyEnv(providerType string) string {
	switch strings.ToLower(strings.TrimSpace(providerType)) {
	case "openai":
		return kbxMod.DefaultLLMOpenAIKeyEnv
	case "gemini":
		return kbxMod.DefaultLLMGeminiKeyEnv
	case "anthropic":
		return kbxMod.DefaultLLMAnthropicKeyEnv
	case "groq":
		return kbxMod.DefaultLLMGroqKeyEnv
	default:
		return ""
	}
}

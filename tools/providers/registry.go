package providers

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/kubex-ecosystem/kbx/load"
	kbxTypes "github.com/kubex-ecosystem/kbx/types"
)

// Registry manages runtime LLM provider instances.
type Registry struct {
	cfg       *kbxTypes.LLMConfig
	providers map[string]kbxTypes.ProviderExt
}

// NewRegistry creates a Registry from the given LLMConfig. Passing nil returns an empty registry.
func NewRegistry(cfg *kbxTypes.LLMConfig) *Registry {
	r := &Registry{
		cfg:       cfg,
		providers: make(map[string]kbxTypes.ProviderExt),
	}
	if cfg != nil {
		for rawName, pc := range cfg.Providers {
			name := strings.ToLower(strings.TrimSpace(rawName))
			if pc != nil {
				r.providers[name] = pc
			}
		}
	}
	return r
}

// Load loads an LLMConfig from a file path and returns a populated Registry.
func Load(path string) (*Registry, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, fmt.Errorf("providers: config path cannot be empty")
	}
	cfg, err := load.Config[kbxTypes.LLMConfig](path)
	if err != nil {
		return nil, fmt.Errorf("providers: failed to load config from %s: %w", path, err)
	}
	return NewRegistry(&cfg), nil
}

// Register adds a named provider to the registry.
func (r *Registry) Register(name string, p kbxTypes.ProviderExt) {
	if r != nil {
		r.providers[strings.ToLower(strings.TrimSpace(name))] = p
	}
}

// Config returns the underlying LLMConfig (empty default if none was loaded).
func (r *Registry) Config() kbxTypes.LLMConfig {
	if r == nil || r.cfg == nil {
		return kbxTypes.NewLLMConfigDefault()
	}
	return *r.cfg
}

// GetProviderConfig returns the config for the named provider, or nil if not found.
func (r *Registry) GetProviderConfig(name string) *kbxTypes.LLMProviderConfig {
	if r == nil || r.cfg == nil || r.cfg.Providers == nil {
		return nil
	}
	return r.cfg.Providers[strings.ToLower(strings.TrimSpace(name))]
}

// ResolveProvider returns the named provider, or nil if not found.
func (r *Registry) ResolveProvider(name string) kbxTypes.ProviderExt {
	if r == nil {
		return nil
	}
	return r.providers[strings.ToLower(strings.TrimSpace(name))]
}

// ListProviders returns a sorted list of registered provider names.
func (r *Registry) ListProviders() []string {
	if r == nil {
		return nil
	}
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Chat dispatches a ChatRequest to the provider named in req.Provider.
func (r *Registry) Chat(ctx context.Context, req kbxTypes.ChatRequest) (<-chan kbxTypes.ChatChunk, error) {
	p := r.ResolveProvider(req.Provider)
	if p == nil {
		return nil, fmt.Errorf("providers: provider %q not found", req.Provider)
	}
	return p.Chat(ctx, req)
}

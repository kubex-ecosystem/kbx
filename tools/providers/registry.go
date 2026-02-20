// Package registry provides provider registration and resolution functionality.
package registry

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	kbx "github.com/kubex-ecosystem/kbx"
	kbxGet "github.com/kubex-ecosystem/kbx/get"
	kbxMod "github.com/kubex-ecosystem/kbx/internal/module/kbx"
	kbxIs "github.com/kubex-ecosystem/kbx/is"
	kbxTypes "github.com/kubex-ecosystem/kbx/types"
	gl "github.com/kubex-ecosystem/logz"
)

// Registry manages provider registration and resolution
type Registry struct {
	cfg *kbxTypes.LLMConfig
}

func NewRegistry(c *kbxTypes.LLMConfig) *Registry {
	return &Registry{
		cfg: c,
	}
}

// Load creates a new registry from a YAML configuration file
func Load(path string) (*Registry, error) {
	path = strings.TrimSpace(filepath.Clean(strings.ToValidUTF8(path, "")))
	if len(path) == 0 {
		gl.Warn("No provider config path specified. AI services will be unavailable.")
		return nil, nil
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		gl.Warnf("Provider config file not found at %s. AI services will be unavailable.", path)
		return nil, nil
	} else if err != nil {
		return nil, gl.Errorf("error checking provider config file at %s: %v", path, err)
	}

	gl.Debugf("Loading provider configuration from %s", path)

	// Load or create config file with kbx method
	rgCfg, err := kbx.LoadConfigOrDefault[kbx.LLMConfig](path, true)
	if err != nil {
		gl.Errorf("Failed to load config: %v", err)
		return nil, fmt.Errorf("failed to load provider config: %w", err)
	} else if rgCfg == nil {
		gl.Noticef("No config file found, proceeding with default auto-generated config at %s", path)
		rgCfg = kbx.NewLLMConfig()
	}
	rg := &Registry{
		cfg: rgCfg,
	}

	// Initialize providers based on configuration
	for name, pc := range rg.cfg.Providers {
		cmiMap, err := pc.ModelInfo(context.Background())
		if err != nil {
			gl.Warnf("Failed to get model info for provider '%s': %v. Using default model from config.", name, err)
		}
		cmiStrName := cmiMap["name"].(string)
		if cmiStrName != "" {
			gl.Debugf("Provider '%s' model info: %v", name, cmiMap)
		} else {
			gl.Debugf("Provider '%s' does not have model info available. Using default model from config.", name)
		}

		var providerConstructor func(name, baseURL, key, model string) (kbxTypes.ProviderExt, error)
		tp := strings.ToLower(pc.Type())
		key := os.ExpandEnv(kbxGet.EnvOr(pc.KeyRef(), kbxGet.ValueOrIf(kbxIs.Map[string](rg.cfg.Providers[name].KeyRef()), kbxGet.EnvOr(rg.cfg.Providers[name].KeyRef(), ""), "")))

		switch tp {
		case "openai":
			lKey := kbxGet.EnvOr(
				kbxGet.EnvOr(strings.ToUpper(name)+"_API_KEY", ""), // Se não achar nada de jeito nenhum até aqui, tenta buscar pela variável de ambiente genérica GEMINI_API_KEY (para compatibilidade com versões anteriores ou configuração sem especificar KeyEnv)
				kbxMod.DefaultLLMOpenAIKeyEnv,                      // Por fim, tenta buscar pela variável de ambiente padrão para o kubex, definida no módulo kbx e documentada, padronizada em todo ecossistema para fins de fallback e resiliência (ex: KUBEX_GNYX_GEMINI_KEY)
			)
			key := rg.getEnvKeyValueWithFallback(*rg.cfg, pc, lKey)
			if key == "" {
				gl.Log("warning", fmt.Sprintf("Skipping OpenAI provider '%s' - no API key found in %s", name, pc.KeyRef()))
				continue
			}
			providerConstructor = NewOpenAIProvider
		case "gemini":
			lKey := kbxGet.EnvOr(
				kbxGet.EnvOr(strings.ToUpper(name)+"_API_KEY", ""), // Se não achar nada de jeito nenhum até aqui, tenta buscar pela variável de ambiente genérica GEMINI_API_KEY (para compatibilidade com versões anteriores ou configuração sem especificar KeyEnv)
				kbxMod.DefaultLLMGeminiKeyEnv,                      // Por fim, tenta buscar pela variável de ambiente padrão para o kubex, definida no módulo kbx e documentada, padronizada em todo ecossistema para fins de fallback e resiliência (ex: KUBEX_GNYX_GEMINI_KEY)
			)
			key := rg.getEnvKeyValueWithFallback(*rg.cfg, pc, lKey)
			if key == "" {
				gl.Log("warning", fmt.Sprintf("Skipping Gemini provider '%s' - no API key found in %s", name, pc.KeyRef()))
				continue
			}
			providerConstructor = NewGeminiProvider
		case "anthropic":
			lKey := kbxGet.EnvOr(
				kbxGet.EnvOr(strings.ToUpper(name)+"_API_KEY", ""), // Se não achar nada de jeito nenhum até aqui, tenta buscar pela variável de ambiente genérica ANTHROPIC_API_KEY (para compatibilidade com versões anteriores ou configuração sem especificar KeyEnv)
				kbxMod.DefaultLLMAnthropicKeyEnv,                   // Por fim, tenta buscar pela variável de ambiente padrão para o kubex, definida no módulo kbx e documentada, padronizada em todo ecossistema para fins de fallback e resiliência (ex: KUBEX_GNYX_ANTHROPIC_KEY)
			)
			key := rg.getEnvKeyValueWithFallback(*rg.cfg, pc, lKey)
			if key == "" {
				gl.Log("warning", fmt.Sprintf("Skipping Anthropic provider '%s' - no API key found in %s", name, pc.KeyRef()))
				continue
			}
			providerConstructor = NewAnthropicProvider
		case "groq":
			lKey := kbxGet.EnvOr(
				kbxGet.EnvOr(strings.ToUpper(name)+"_API_KEY", ""), // Se não achar nada de jeito nenhum até aqui, tenta buscar pela variável de ambiente genérica GROQ_API_KEY (para compatibilidade com versões anteriores ou configuração sem especificar KeyEnv)
				kbxMod.DefaultLLMGroqKeyEnv,                        // Por fim, tenta buscar pela variável de ambiente padrão para o kubex, definida no módulo kbx e documentada, padronizada em todo ecossistema para fins de fallback e resiliência (ex: KUBEX_GNYX_GROQ_KEY)
			)
			key := rg.getEnvKeyValueWithFallback(*rg.cfg, pc, lKey)
			if key == "" {
				gl.Log("warning", fmt.Sprintf("Skipping Groq provider '%s' - no API key found in %s", name, pc.KeyRef()))
				continue
			}
			providerConstructor = NewGroqProvider
		case "openrouter":
			// TODO: Implement OpenRouter provider
			return nil, gl.Errorf("openrouter provider not yet implemented")
		case "ollama":
			// TODO: Implement Ollama provider
			return nil, gl.Errorf("ollama provider not yet implemented")
		default:
			gl.Errorf("unknown provider type: %s", tp)
			continue
		}

		p, err := providerConstructor(name, pc.URLBase(), key, cmiStrName)
		if err != nil {
			gl.Warnf("Failed to initialize provider '%s': %v. This provider will be unavailable for use.", name, err)
			continue
		}
		nPrv, ok := p.(*kbxTypes.LLMProviderConfig)
		if !ok {
			gl.Warnf("Failed to assert provider '%s' to *kbxTypes.LLMProviderConfig", name)
			continue
		}

		rg.cfg.Providers[name] = nPrv
	}

	return rg, nil
}

func (r *Registry) Providers() kbxTypes.LLMProvidersExtMap {
	providers := make(map[string]kbxTypes.ProviderExt)
	for name, pc := range r.cfg.Providers {
		p := r.ResolveProvider(name)
		if p != nil {
			providers[name] = p
		} else if pc != nil {
			providers[name] = pc
		} else {
			gl.Warnf("Provider '%s' not found in registry and failed to load from config. This provider will be unavailable for use.", name)
		}
	}
	return providers
}

func (r *Registry) GetProvider(name string) kbxTypes.ProviderExt {
	p := r.ResolveProvider(name)
	if p == nil {
		gl.Warnf("Provider '%s' not found in registry. Trying to recover from process history and/or config...", name)
		if r.cfg != nil {
			if pc, ok := r.cfg.Providers[name]; ok {
				gl.Infof("Provider '%s' found in config. Adding to registry.", name)
				r.cfg.Providers[name] = pc
				return pc
			}
		}
		gl.Warnf("Provider '%s' not found in config. Unable to recover provider.", name)
	}
	return nil
}

func (r *Registry) GetProviderConfig(name string) *kbxTypes.LLMProviderConfig {
	if r.cfg != nil {
		if pc, ok := r.cfg.Providers[name]; ok {
			return pc
		}
	}
	return nil
}

func (r *Registry) getEnvKeyValueWithFallback(c kbxTypes.LLMConfig, p kbxTypes.ProviderExt, fb string) string {
	pt := r.ResolveProvider(p.Name())
	pp, ok := pt.(*kbxTypes.LLMProviderConfig)
	if !ok {
		gl.Warnf("Provider '%s' does not have a valid LLMProviderConfig, cannot resolve API key with fallback logic", p.Name())
		return fb
	}

	lKey := kbxGet.EnvOr(
		pp.KeyRef(), // Busca pela variável de ambiente específica do provider (Ex: GEMINI_API_KEY)
		kbxGet.ValueOrIf(
			kbxIs.Map[string](c.Providers[p.Name()].KeyRef()), // Verifica se a chave do provider é um mapa (para múltiplas chaves por provider) e se o valor tem o tipo correto caso ele exista
			kbxGet.EnvOr(
				c.Providers[p.Name()].KeyRef(), // Se passar no teste do mapa anterior, tenta buscar a variável de ambiente usando o valor do campo KeyEnv do provider (ex: GEMINI_API_KEY_1)
				kbxGet.ValOrType(fb, ""),       // Se não for um mapa ou tiver tipo errado, usa o valor lido anteriormente (que já considera o fallback para a variável de ambiente genérica)
			),
			kbxGet.ValOrType(fb, ""), // Se não for um mapa ou tiver tipo errado, usa o valor lido anteriormente (que já considera o fallback para a variável de ambiente genérica)
		),
	)

	key := os.ExpandEnv(
		kbxGet.EnvOr(
			p.KeyRef(), // Busca pela variável de ambiente específica do provider (Ex: GEMINI_API_KEY)
			kbxGet.ValueOrIf(
				kbxIs.Map[string](p), // Verifica se a chave do provider é um mapa (para múltiplas chaves por provider) e se o valor tem o tipo correto caso ele exista
				kbxGet.EnvOr(
					p.KeyRef(),                 // Se passar no teste do mapa anterior, tenta buscar a variável de ambiente usando o valor do campo KeyEnv do provider (ex: GEMINI_API_KEY_1)
					kbxGet.ValOrType(lKey, ""), // Se não for um mapa ou tiver tipo errado, usa o valor lido anteriormente (que já considera o fallback para a variável de ambiente genérica)
				),
				kbxGet.ValOrType(lKey, ""), // Se não for um mapa ou tiver tipo errado, usa o valor lido anteriormente (que já considera o fallback para a variável de ambiente genérica)
			),
		),
	)

	return key
}

// Resolve returns a provider by name
func (r *Registry) Resolve(name string) kbxTypes.Provider {
	if len(r.cfg.Providers) == 0 {
		r.cfg.Providers = make(map[string]*kbxTypes.LLMProviderConfig)
	}
	p, ok := r.cfg.Providers[name]
	if !ok {
		gl.Warnf("Provider '%s' found in registry but does not implement Provider interface. This provider will be unavailable for use.", name)
		return nil
	}
	if p == nil {
		gl.Warnf("Provider '%s' not found in registry.", name)
		return nil
	}
	return p
}

// ListProviders returns all available provider names
func (r *Registry) ListProviders() []string {
	if !kbxIs.Valid(r.cfg.Providers) {
		gl.Warn("Provider registry is empty. No providers available for listing.")
		r.cfg.Providers = make(map[string]*kbxTypes.LLMProviderConfig)
		cfg := r.GetConfig()
		if cfg.Providers != nil {
			for name := range cfg.Providers {
				gl.Debugf("Provider '%s' found in config during ListProviders. Adding to registry.", name)
				r.cfg.Providers[name] = cfg.Providers[name]
			}
		} else {
			gl.Warn("No providers found in config during ListProviders.")
		}
		return []string{}
	}
	names := make([]string, 0, len(r.cfg.Providers))
	for name := range r.cfg.Providers {
		names = append(names, name)
	}
	return names
}

// GetConfig returns the provider configuration
func (r *Registry) GetConfig() kbxTypes.LLMConfig {
	if r.cfg == nil {
		gl.Warn("Provider registry config is nil. Returning empty config.")
		return kbxTypes.LLMConfig{}
	}
	return *r.cfg
}

func (r *Registry) ResolveProvider(name string) kbxTypes.ProviderExt {
	if len(r.cfg.Providers) == 0 {
		gl.Warnf("Provider registry is empty. No providers available for resolution.")
		return nil
	}
	if p, ok := r.cfg.Providers[name]; ok {
		return p
	}
	gl.Warnf("Provider '%s' not found in registry.", name)
	return nil
}

func (r *Registry) Config() kbxTypes.LLMConfig {
	if r.cfg == nil {
		gl.Warn("Provider registry config is nil. Returning empty config.")
		return kbxTypes.LLMConfig{}
	}
	return *r.cfg
} // <- usado por /v1/providers

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

// /v1/chat/completions — SSE endpoints

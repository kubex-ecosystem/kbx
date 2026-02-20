package types

import (
	"context"
	"reflect"
	"sync"
	"time"

	gl "github.com/kubex-ecosystem/logz"
)

// NotifierEventBase represents a notification to be sent
type NotifierEventBase[P NotifierConfig[P] | LLMProviderConfig | ChatRequest | LLMConfig | any] interface {
	// Common methods for all notifier events
	Ref() GlobalRef
	NType() (reflect.Type, string)

	Type(context.Context) string
	Recipient(context.Context) string
	Subject(context.Context) string
	Content(context.Context) string
	Priority(context.Context) string
	Metadata(context.Context) map[string]any
	CreatedAt(context.Context) time.Time
}

type NotifierEvent[P NotifierConfig[P] | LLMProviderConfig | ChatRequest | LLMConfig | any] interface {
	NotifierEventBase[P]
	Dispatch(ctx context.Context) <-chan error
	Done(context.Context) <-chan struct{}
	Cancel(context.Context) <-chan struct{}
	Error() error
}

type NotifierEventExt[P NotifierConfig[P] | LLMProviderConfig | ChatRequest | LLMConfig | any] interface {
	NotifierEventBase[P]
	NotifierEvent[P]

	DoneWithError(context.Context) <-chan error

	Reset(context.Context) error
	Retry(context.Context) <-chan NotifierEvent[P]

	Wait(ctx context.Context) error
	WaitWithTimeout(ctx context.Context, timeout time.Duration) error
	WaitWithCondition(ctx context.Context, cond sync.Cond) error

	Timeout(context.Context) time.Duration
}

// Notifier interface defines the contract for sending notifications
type Notifier interface {
	Send(ctx context.Context, event NotifierEvent[any]) error
}

// NotifierProvider interface defines the contract for notification providers
type NotifierProvider interface {
	Name() string
	Notify(ctx context.Context, event NotifierEvent[any]) error
}

// NotifierConfig holds configuration for a specific notifier provider
type NotifierConfig[P NotifierConfig[P] | LLMProviderConfig | ChatRequest | LLMConfig | any] struct {
	Type       string         `yaml:"type"` // "discord", "whatsapp", "email"
	Recipient  string         `yaml:"recipient"`
	Subject    string         `yaml:"subject"`
	Content    string         `yaml:"content"`
	Priority   string         `yaml:"priority"` // "low", "medium", "high", "critical"
	Metadata   map[string]any `yaml:"metadata"`
	CreatedAt  time.Time      `yaml:"created_at"`
	Provider   string         `yaml:"provider"`   // e.g., "sendgrid", "twilio", "discord"
	Parameters map[string]any `yaml:"parameters"` // provider-specific parameters
}

// NotifierRegistry holds registered notifier providers
type NotifierRegistry struct {
	providers map[string]NotifierProvider
}

// Register adds a new notifier provider to the registry
func (r *NotifierRegistry) Register(provider NotifierProvider) {
	if r.providers == nil {
		r.providers = make(map[string]NotifierProvider)
	}
	r.providers[provider.Name()] = provider
}

// Notify sends a notification using the appropriate provider based on the event type
func (r *NotifierRegistry) Notify(ctx context.Context, event NotifierEvent[any]) error {
	provider, exists := r.providers[event.Type(ctx)]
	if !exists {
		return gl.Errorf("no notifier provider registered for type '%s'", event.Type(ctx))
	}
	return provider.Notify(ctx, event)
}

// ListProviders returns the names of all registered notifier providers
func (r *NotifierRegistry) ListProviders() []string {
	providers := make([]string, 0, len(r.providers))
	for name := range r.providers {
		providers = append(providers, name)
	}
	return providers
}

func (r *NotifierRegistry) GetProvider(name string) (NotifierProvider, bool) {
	provider, exists := r.providers[name]
	return provider, exists
}

// NotificationEvent represents a notification to be sent
type NotificationEvent struct {
	Type      string         `json:"type"` // "discord", "whatsapp", "email"
	Recipient string         `json:"recipient"`
	Subject   string         `json:"subject"`
	Content   string         `json:"content"`
	Priority  string         `json:"priority"` // "low", "medium", "high", "critical"
	Metadata  map[string]any `json:"metadata"`
	CreatedAt time.Time      `json:"created_at"`
}

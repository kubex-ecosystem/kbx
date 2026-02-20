// Package kbx has default configuration values
package kbx

const (
	KeyringService        = "kubex"
	DefaultKubexConfigDir = "$HOME/.kubex"

	DefaultGNyxKeyPath    = "$HOME/.kubex/gnyx/kubex_gnyx-key.pem"
	DefaultGNyxCertPath   = "$HOME/.kubex/gnyx/kubex_gnyx-cert.pem"
	DefaultGNyxCAPath     = "$HOME/.kubex/gnyx/ca-cert.pem"
	DefaultGNyxConfigPath = "$HOME/.kubex/gnyx/config/config.json"

	DefaultConfigDir       = "$HOME/.kubex/domus/config"
	DefaultConfigFile      = "$HOME/.kubex/domus/config.json"
	DefaultDomusConfigPath = "$HOME/.kubex/domus/config/config.json"
)

const (
	DefaultVolumesDir     = "$HOME/.kubex/volumes"
	DefaultRedisVolume    = "$HOME/.kubex/volumes/redis"
	DefaultPostgresVolume = "$HOME/.kubex/volumes/postgresql"
	DefaultMongoDBVolume  = "$HOME/.kubex/volumes/mongodb"
	DefaultMongoVolume    = "$HOME/.kubex/volumes/mongo"
	DefaultRabbitMQVolume = "$HOME/.kubex/volumes/rabbitmq"
)

const (
	DefaultRateLimitLimit  = 100
	DefaultRateLimitBurst  = 100
	DefaultRequestWindow   = 1 * 60 * 1000 // 1 minute
	DefaultRateLimitJitter = 0.1
)

const (
	DefaultMaxRetries = 3
	DefaultRetryDelay = 1 * 1000 // 1 second
)

const (
	DefaultMaxIdleConns          = 100
	DefaultMaxIdleConnsPerHost   = 100
	DefaultIdleConnTimeout       = 90 * 1000 // 90 seconds
	DefaultTLSHandshakeTimeout   = 10 * 1000 // 10 seconds
	DefaultExpectContinueTimeout = 1 * 1000  // 1 second
	DefaultResponseHeaderTimeout = 5 * 1000  // 5 seconds
	DefaultTimeout               = 30 * 1000 // 30 seconds
	DefaultKeepAlive             = 30 * 1000 // 30 seconds
	DefaultMaxConnsPerHost       = 100
)

const (
	DefaultLLMOpenAIKeyEnv       = "OPENAI_API_KEY"
	DefaultLLMGoogleKeyEnv       = "GOOGLE_API_KEY"
	DefaultLLMAzureKeyEnv        = "AZURE_API_KEY"
	DefaultLLMAnthropicKeyEnv    = "ANTHROPIC_API_KEY"
	DefaultLLMGeminiKeyEnv       = "GEMINI_API_KEY"
	DefaultLLMOllamaKeyEnv       = "OLLAMA_API_KEY"
	DefaultLLMChatGPTKeyEnv      = "CHATGPT_API_KEY"
	DefaultLLMDeepseekKeyEnv     = "DEEPSEEK_API_KEY"
	DefaultLLMCohereKeyEnv       = "COHERE_API_KEY"
	DefaultLLMGroqKeyEnv         = "GROQ_API_KEY"
	DefaultLLMGrokKeyEnv         = "GROK_API_KEY"
	DefaultLLMMistralKeyEnv      = "MISTRAL_API_KEY"
	DefaultLLMCustomKeyEnv       = "CUSTOM_API_KEY"
	DefaultLLMMetaKeyEnv         = "META_API_KEY"
	DefaultLLMClaudeKeyEnv       = "CLAUDE_API_KEY"
	DefaultLLMErnieKeyEnv        = "ERNIE_API_KEY"
	DefaultLLMCustomKeyEnvPrefix = "CUSTOM_"
	DefaultLLMCustomKeyEnvSuffix = "_KEY_ENV"
)

const (
	DefaultLLMProvider    = "gemini"
	DefaultLLMModel       = "gemini-2.0-flash"
	DefaultLLMMaxTokens   = 1024
	DefaultLLMTemperature = 0.3
)

const (
	DefaultApprovalRequireForResponses = false
	DefaultApprovalTimeoutMinutes      = 15
)

const (
	DefaultServerPort = "5000"
	DefaultServerHost = "0.0.0.0"
)

type ValidationError struct {
	Field   string
	Message string
}

func (v *ValidationError) Error() string {
	return v.Message
}
func (v *ValidationError) FieldError() map[string]string {
	return map[string]string{v.Field: v.Message}
}
func (v *ValidationError) FieldsError() map[string]string {
	return map[string]string{v.Field: v.Message}
}
func (v *ValidationError) ErrorOrNil() error {
	return v
}

var (
	ErrUsernameRequired = &ValidationError{Field: "username", Message: "Username is required"}
	ErrPasswordRequired = &ValidationError{Field: "password", Message: "Password is required"}
	ErrEmailRequired    = &ValidationError{Field: "email", Message: "Email is required"}
	ErrDBNotProvided    = &ValidationError{Field: "db", Message: "Database not provided"}
	ErrModelNotFound    = &ValidationError{Field: "model", Message: "Model not found"}
)

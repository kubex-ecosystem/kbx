package tools

import (
	"context"
	"errors"
	"time"
)

// ErrRetryExhausted is returned when all retry attempts have been exhausted.
var ErrRetryExhausted = errors.New("retry limits exhausted")

// Retryer struct
type Retryer struct {
	cfg RetryConfig
}

// NewRetryer creates a new retryer.
func NewRetryer(cfg RetryConfig) *Retryer {
	return &Retryer{cfg: cfg}
}

// DoVoid executes a function with retry logic, ignoring the result.
func (r *Retryer) DoVoid(fn func() error) error {
	_, err := Retry(func() (struct{}, error) {
		return struct{}{}, fn()
	},
		WithRetries(r.cfg.Retries),
		WithDelay(r.cfg.Delay),
		WithTimeout(r.cfg.Timeout),
	)
	return err
}

// Do executes a function with retry logic.
func (r *Retryer) Do(fn func(a any) (any, error)) (any, error) {
	return Retry[any](
		func() (any, error) {
			return fn(nil)
		},
		WithRetries(r.cfg.Retries),
		WithDelay(r.cfg.Delay),
		WithTimeout(r.cfg.Timeout),
	)
}

// RetryResult holds the result of a retry operation.
type RetryResult[T any] struct {
	Result T
	Error  error
}

// RetryStats holds the statistics of a retry operation.
type RetryStats struct {
	Attempts    int
	Successful  bool
	TotalTime   time.Duration
	MaxAttempts int
}

// RetryOption holds the options for the retry mechanism.
type RetryOption struct {
	Retries       int
	Delay         time.Duration
	Timeout       time.Duration
	MaxDelay      time.Duration
	InitialDelay  time.Duration
	MaxRetries    int
	BackoffFactor float64
}

// RetryConfig holds the configuration for the retry mechanism.
type RetryConfig struct {
	Retries       int
	Delay         time.Duration
	Timeout       time.Duration
	MaxDelay      time.Duration
	InitialDelay  time.Duration
	BackoffFactor float64
}

// WithRetries returns a RetryOption retries value.
func WithRetries(n int) *RetryOption {
	return &RetryOption{Retries: n}
}

// WithDelay returns a RetryOption delay value.
func WithDelay(d time.Duration) *RetryOption {
	return &RetryOption{Delay: d}
}

// WithMaxAttempts returns a RetryOption max attempts value.
func WithMaxAttempts(n int) *RetryOption {
	return &RetryOption{MaxRetries: n}
}

// WithInitialDelay returns a RetryOption initial delay value.
func WithInitialDelay(d time.Duration) *RetryOption {
	return &RetryOption{InitialDelay: d}
}

// WithMaxDelay returns a RetryOption max delay value.
func WithMaxDelay(d time.Duration) *RetryOption {
	return &RetryOption{MaxDelay: d}
}

// WithBackoffFactor returns a RetryOption backoff factor value.
func WithBackoffFactor(f float64) *RetryOption {
	return &RetryOption{BackoffFactor: f}
}

// WithTimeout returns a RetryOption timeout value.
func WithTimeout(d time.Duration) *RetryOption {
	return &RetryOption{Timeout: d}
}

// Retry unifica TUDO. Se for Void, basta usar Retry[struct{}] e ignorar o retorno.
func Retry[T any](fn func() (T, error), opts ...*RetryOption) (T, error) {
	// Configuração Default
	config := RetryConfig{
		Retries: 3,
		Delay:   0,
		Timeout: 0,
	}
	for _, opt := range opts {
		if opt.Retries > 0 {
			config.Retries = opt.Retries
		}
		if opt.Delay > 0 {
			config.Delay = opt.Delay
		}
		if opt.Timeout > 0 {
			config.Timeout = opt.Timeout
		}
		// if opt.MaxDelay > 0 {
		// 	config.MaxDelay = opt.MaxDelay
		// }
		// if opt.InitialDelay > 0 {
		// 	config.InitialDelay = opt.InitialDelay
		// }
		// if opt.BackoffFactor > 0 {
		// 	config.BackoffFactor = opt.BackoffFactor
		// }
	}

	var lastErr error
	var result T

	// Controle de Timeout Global via Contexto (Jeito Go moderno)
	var ctx context.Context
	var cancel context.CancelFunc

	if config.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), config.Timeout)
		defer cancel()
	} else {
		ctx = context.Background()
	}

	for i := 0; i < config.Retries; i++ {
		// Checa se o tempo global estourou antes de tentar
		select {
		case <-ctx.Done():
			// Retorna erro de timeout, não nil!
			return result, ctx.Err()
		default:
		}

		result, lastErr = fn()
		if lastErr == nil {
			return result, nil
		}

		// Se não for a última tentativa, espera
		if i < config.Retries-1 && config.Delay > 0 {
			// Sleep respeitando o contexto (se cancelar no meio do sleep, acorda)
			select {
			case <-ctx.Done():
				return result, ctx.Err()
			case <-time.After(config.Delay):
				// continua
			}
		}
	}

	// Se saiu do loop, falhou todas. Retorna o último erro.
	if lastErr != nil {
		return result, lastErr
	}
	return result, ErrRetryExhausted
}

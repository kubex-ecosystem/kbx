package tools

import (
	"context"
	"errors"
	"time"
)

var ErrRetryExhausted = errors.New("retry limits exhausted")

type Retryer struct {
	cfg RetryConfig
}

func NewRetryer(cfg RetryConfig) *Retryer {
	return &Retryer{cfg: cfg}
}

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

type RetryResult[T any] struct {
	Result T
	Error  error
}

type RetryStats struct {
	Attempts    int
	Successful  bool
	TotalTime   time.Duration
	MaxAttempts int
}
type RetryOption struct {
	Retries       int
	Delay         time.Duration
	Timeout       time.Duration
	MaxDelay      time.Duration
	InitialDelay  time.Duration
	MaxRetries    int
	BackoffFactor float64
}

type RetryConfig struct {
	Retries       int
	Delay         time.Duration
	Timeout       time.Duration
	MaxDelay      time.Duration
	InitialDelay  time.Duration
	BackoffFactor float64
}

func WithRetries(n int) *RetryOption {
	return &RetryOption{Retries: n}
}

func WithDelay(d time.Duration) *RetryOption {
	return &RetryOption{Delay: d}
}

func WithMaxAttempts(n int) *RetryOption {
	return &RetryOption{MaxRetries: n}
}

func WithInitialDelay(d time.Duration) *RetryOption {
	return &RetryOption{InitialDelay: d}
}

func WithMaxDelay(d time.Duration) *RetryOption {
	return &RetryOption{MaxDelay: d}
}

func WithBackoffFactor(f float64) *RetryOption {
	return &RetryOption{BackoffFactor: f}
}

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

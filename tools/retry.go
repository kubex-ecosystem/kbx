package tools

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"time"

	"github.com/kubex-ecosystem/logz"
)

// ErrRetryExhausted is returned when all retry attempts have been exhausted.
var ErrRetryExhausted = errors.New("retry limits exhausted")

// Retryer struct
type Retryer struct {
	cfg *RetryOptionWithBackoff
}

// NewRetryer creates a new retryer.
func NewRetryer(opts ...*RetryOption) *Retryer {
	config := DefaultRetryConfig()
	for _, opt := range opts {
		if opt.Retries > 0 {
			config.Retries = opt.Retries
		}
		if opt.InitialDelay > 0 {
			config.InitialDelay = opt.InitialDelay
		}
		if opt.MaxDelay > 0 {
			config.MaxDelay = opt.MaxDelay
		}
	}
	return &Retryer{cfg: config}
}

func NewRetryerWithBackoff(opts ...*RetryOptionWithBackoff) *Retryer {
	config := DefaultRetryConfig()
	for _, opt := range opts {
		if opt.Retries > 0 {
			config.Retries = opt.Retries
		}
		if opt.InitialDelay > 0 {
			config.InitialDelay = opt.InitialDelay
		}
		if opt.MaxDelay > 0 {
			config.MaxDelay = opt.MaxDelay
		}
		if opt.BackoffFactor > 0 {
			config.BackoffFactor = opt.BackoffFactor
		}
	}
	return &Retryer{cfg: config}
}

// DoVoid executes a function with retry logic, ignoring the result.
func (r *Retryer) DoVoid(fn func() error) error {
	_, err := Retry(func() (struct{}, error) {
		return struct{}{}, fn()
	},
		WithRetries(r.cfg.Retries),
		WithDelay(r.cfg.InitialDelay),
		WithTimeout(r.cfg.MaxDelay),
	)
	return err
}

// Do executes a function with retry logic.
func (r *Retryer) Do(fn func(a any) (any, error)) (any, error) {
	return Retry(
		func() (any, error) {
			return fn(nil)
		},
		WithRetries(r.cfg.Retries),
		WithDelay(r.cfg.InitialDelay),
		WithTimeout(r.cfg.MaxDelay),
	)
}

// RetryOptionWithBackoff holds the options for the retry mechanism.
type RetryOption struct {
	Retries      int
	MaxDelay     time.Duration
	InitialDelay time.Duration
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() *RetryOptionWithBackoff {
	return &RetryOptionWithBackoff{
		RetryOption: &RetryOption{
			Retries:      5,
			MaxDelay:     5 * time.Second,
			InitialDelay: 100 * time.Millisecond,
		},
		Timeout:       5 * time.Second,
		BackoffFactor: 1.5,
	}
}

// RetryOption holds the options for the retry mechanism.
type RetryOptionWithBackoff struct {
	*RetryOption
	Timeout       time.Duration
	BackoffFactor float64
}

// WithRetries returns a RetryOption retries value.
func WithRetries(n int) *RetryOption {
	return &RetryOption{Retries: n}
}

// WithDelay returns a RetryOption delay value.
func WithDelay(d time.Duration) *RetryOption {
	return &RetryOption{InitialDelay: d}
}

// WithTimeout returns a RetryOption timeout value.
func WithTimeout(d time.Duration) *RetryOption {
	return &RetryOption{MaxDelay: d}
}

// WithMaxAttempts returns a RetryOption max attempts value.
func WithMaxAttempts(n int) *RetryOption {
	return &RetryOption{
		Retries: n,
	}
}

// Retry executes fn up to the configured number of times, returning on first success.
func Retry[T any](fn func() (T, error), opts ...*RetryOption) (T, error) {
	config := &RetryOption{
		Retries:      3,
		InitialDelay: 0,
		MaxDelay:     0,
	}
	for _, opt := range opts {
		if opt.Retries > 0 {
			config.Retries = opt.Retries
		}
		if opt.InitialDelay > 0 {
			config.InitialDelay = opt.InitialDelay
		}
		if opt.MaxDelay > 0 {
			config.MaxDelay = opt.MaxDelay
		}
	}

	var lastErr error
	var result T

	var ctx context.Context
	var cancel context.CancelFunc

	if config.MaxDelay > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), config.MaxDelay)
		defer cancel()
	} else {
		ctx = context.Background()
	}

	for i := 0; i < config.Retries; i++ {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}

		result, lastErr = fn()
		if lastErr == nil {
			return result, nil
		}

		if i < config.Retries-1 && config.InitialDelay > 0 {
			select {
			case <-ctx.Done():
				return result, ctx.Err()
			case <-time.After(config.InitialDelay):
			}
		}
	}

	if lastErr != nil {
		return result, lastErr
	}
	return result, ErrRetryExhausted
}

// RetryWithBackoff executes a function with exponential backoff retry logic
func RetryWithBackoff(ctx context.Context, config RetryOptionWithBackoff, operation func() error) error {
	var lastErr error

	for attempt := 0; attempt <= config.Retries; attempt++ {
		// Try the operation
		err := operation()
		if err == nil {
			return nil // Success!
		}

		lastErr = err

		// Don't wait after the last attempt
		if attempt == config.Retries {
			break
		}

		// Calculate delay with exponential backoff
		delay := time.Duration(float64(config.InitialDelay) * math.Pow(config.BackoffFactor, float64(attempt)))
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}

		// Add a jitter (variation of 0 to 20%)
		jitter := time.Duration(rand.Int63n(int64(delay) / 5))
		delay += jitter

		// Log the retry attempt
		logz.Warnf("Attempt %d failed: %v. Retrying in %v...", attempt+1, err, delay)

		// Wait for the delay or context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	return logz.Errorf("operation failed after %d attempts: %v", config.Retries+1, lastErr)
}

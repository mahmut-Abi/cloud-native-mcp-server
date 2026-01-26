package retry

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// RetryConfig holds configuration for retry logic
type RetryConfig struct {
	MaxAttempts   int                          // Maximum number of retry attempts
	InitialDelay  time.Duration                // Initial delay before first retry
	MaxDelay      time.Duration                // Maximum delay between retries
	Multiplier    float64                      // Multiplier for exponential backoff
	Jitter        bool                         // Add random jitter to delay
	RetryableFunc func(error) bool             // Function to determine if error is retryable
	OnRetry       func(attempt int, err error) // Callback called on each retry
}

// DefaultRetryConfig returns default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:   3,
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      30 * time.Second,
		Multiplier:    2.0,
		Jitter:        true,
		RetryableFunc: DefaultRetryableFunc,
	}
}

// DefaultRetryableFunc determines if an error is retryable
func DefaultRetryableFunc(err error) bool {
	if err == nil {
		return false
	}

	// Check if error is wrapped in a ServiceError
	type serviceError interface {
		Code() string
		HTTPStatus() int
	}

	if se, ok := err.(serviceError); ok {
		// Retry on 5xx errors, 429 (rate limit), and timeout errors
		if se.HTTPStatus() >= 500 || se.HTTPStatus() == 429 {
			return true
		}
		// Retry on timeout errors
		if se.Code() == "TIMEOUT" || se.Code() == "PROMETHEUS_QUERY_TIMEOUT" {
			return true
		}
		return false
	}

	return false
}

// Retry executes a function with retry logic
func Retry(ctx context.Context, config RetryConfig, fn func() error) error {
	if config.MaxAttempts <= 0 {
		config.MaxAttempts = 1
	}

	var lastErr error
	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		if attempt > 0 {
			// Calculate delay with exponential backoff
			delay := calculateDelay(config, attempt)

			// Add jitter if enabled
			if config.Jitter {
				delay = addJitter(delay)
			}

			// Call on retry callback
			if config.OnRetry != nil {
				config.OnRetry(attempt, lastErr)
			}

			// Wait before retry
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return fmt.Errorf("retry cancelled: %w", ctx.Err())
			}
		}

		// Execute the function
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if error is retryable
		if !config.RetryableFunc(err) {
			return err
		}

		// Check if we should stop retrying
		if attempt >= config.MaxAttempts-1 {
			break
		}
	}

	return fmt.Errorf("max retry attempts (%d) reached, last error: %w", config.MaxAttempts, lastErr)
}

// RetryWithContext executes a function with retry logic and context
func RetryWithContext(ctx context.Context, config RetryConfig, fn func(context.Context) error) error {
	if config.MaxAttempts <= 0 {
		config.MaxAttempts = 1
	}

	var lastErr error
	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		if attempt > 0 {
			// Calculate delay with exponential backoff
			delay := calculateDelay(config, attempt)

			// Add jitter if enabled
			if config.Jitter {
				delay = addJitter(delay)
			}

			// Call on retry callback
			if config.OnRetry != nil {
				config.OnRetry(attempt, lastErr)
			}

			// Wait before retry
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return fmt.Errorf("retry cancelled: %w", ctx.Err())
			}
		}

		// Execute the function
		err := fn(ctx)
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if error is retryable
		if !config.RetryableFunc(err) {
			return err
		}

		// Check if we should stop retrying
		if attempt >= config.MaxAttempts-1 {
			break
		}
	}

	return fmt.Errorf("max retry attempts (%d) reached, last error: %w", config.MaxAttempts, lastErr)
}

// calculateDelay calculates the delay for a given attempt using exponential backoff
func calculateDelay(config RetryConfig, attempt int) time.Duration {
	delay := float64(config.InitialDelay) * math.Pow(config.Multiplier, float64(attempt-1))

	// Cap at max delay
	if delay > float64(config.MaxDelay) {
		delay = float64(config.MaxDelay)
	}

	return time.Duration(delay)
}

// addJitter adds random jitter to the delay (up to 25%)
func addJitter(delay time.Duration) time.Duration {
	// Add up to 25% random jitter
	jitter := rand.Float64() * 0.25
	jitteredDelay := float64(delay) * (1 + jitter)
	return time.Duration(jitteredDelay)
}

// RetryableFunc creates a retryable function for specific HTTP status codes
func RetryableFunc(statusCodes ...int) func(error) bool {
	return func(err error) bool {
		if err == nil {
			return false
		}

		type serviceError interface {
			HTTPStatus() int
		}

		if se, ok := err.(serviceError); ok {
			for _, code := range statusCodes {
				if se.HTTPStatus() == code {
					return true
				}
			}
		}

		return false
	}
}

// RetryableFuncForCode creates a retryable function for specific error codes
func RetryableFuncForCode(codes ...string) func(error) bool {
	return func(err error) bool {
		if err == nil {
			return false
		}

		type serviceError interface {
			Code() string
		}

		if se, ok := err.(serviceError); ok {
			for _, code := range codes {
				if se.Code() == code {
					return true
				}
			}
		}

		return false
	}
}

// Do executes a function with default retry configuration
func Do(ctx context.Context, fn func() error) error {
	return Retry(ctx, DefaultRetryConfig(), fn)
}

// DoWithContext executes a function with default retry configuration and context
func DoWithContext(ctx context.Context, fn func(context.Context) error) error {
	return RetryWithContext(ctx, DefaultRetryConfig(), fn)
}

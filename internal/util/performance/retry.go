package optimize

import (
	"context"
	stderrs "errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultRetryMaxRetries = 2
	DefaultRetryBaseDelay  = 200 * time.Millisecond
	DefaultRetryMaxDelay   = 2 * time.Second
)

// HTTPRetryEvent describes one retry event emitted from DoWithHTTPRetry.
type HTTPRetryEvent struct {
	Attempt    int
	Delay      time.Duration
	StatusCode int
	Err        error
}

// NormalizeRetryConfig clamps retry settings to safe defaults.
func NormalizeRetryConfig(maxRetries int, baseDelay, maxDelay time.Duration) (int, time.Duration, time.Duration) {
	if maxRetries < 0 {
		maxRetries = 0
	}
	if maxRetries == 0 {
		maxRetries = DefaultRetryMaxRetries
	}

	if baseDelay <= 0 {
		baseDelay = DefaultRetryBaseDelay
	}
	if maxDelay <= 0 {
		maxDelay = DefaultRetryMaxDelay
	}
	if maxDelay < baseDelay {
		maxDelay = baseDelay
	}

	return maxRetries, baseDelay, maxDelay
}

// IsRetryableMethod returns true for methods that are safe to retry by default.
func IsRetryableMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead:
		return true
	default:
		return false
	}
}

// ShouldRetryStatusCode reports whether an HTTP status should be retried.
func ShouldRetryStatusCode(statusCode int) bool {
	switch statusCode {
	case http.StatusRequestTimeout, http.StatusTooManyRequests, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

// ShouldRetryTransportError reports whether an error is likely transient.
func ShouldRetryTransportError(err error) bool {
	if err == nil {
		return false
	}

	var netErr net.Error
	if stderrs.As(err, &netErr) {
		if netErr.Timeout() {
			return true
		}
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "timeout") ||
		strings.Contains(msg, "temporarily unavailable") ||
		strings.Contains(msg, "try again")
}

// NextRetryDelay computes exponential backoff with bounded jitter.
func NextRetryDelay(baseDelay, maxDelay time.Duration, attempt int) time.Duration {
	delay := baseDelay
	for i := 1; i < attempt; i++ {
		delay *= 2
		if delay >= maxDelay {
			delay = maxDelay
			break
		}
	}
	if delay > maxDelay {
		delay = maxDelay
	}

	// Add bounded jitter to avoid synchronized retries.
	jitter := time.Duration(rand.Int63n(int64(delay/4 + 1)))
	return delay + jitter
}

// WaitForRetry blocks for a retry delay or exits early when the context is done.
func WaitForRetry(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// DoWithHTTPRetry executes an HTTP request with idempotent retries and backoff.
// Callers should pass normalized retry values (see NormalizeRetryConfig).
func DoWithHTTPRetry(
	ctx context.Context,
	method string,
	maxRetries int,
	baseDelay, maxDelay time.Duration,
	do func(attempt int) (*http.Response, error),
	onRetry func(event HTTPRetryEvent),
) (*http.Response, error) {
	allowRetry := IsRetryableMethod(method)
	totalAttempts := 1
	if allowRetry {
		totalAttempts = maxRetries + 1
	}

	for attempt := 1; attempt <= totalAttempts; attempt++ {
		resp, err := do(attempt)
		if err == nil {
			if allowRetry && ShouldRetryStatusCode(resp.StatusCode) && attempt < totalAttempts {
				_ = resp.Body.Close()
				delay := NextRetryDelay(baseDelay, maxDelay, attempt)
				if onRetry != nil {
					onRetry(HTTPRetryEvent{
						Attempt:    attempt,
						Delay:      delay,
						StatusCode: resp.StatusCode,
					})
				}
				if waitErr := WaitForRetry(ctx, delay); waitErr != nil {
					return nil, waitErr
				}
				continue
			}
			return resp, nil
		}

		if stderrs.Is(err, context.Canceled) || stderrs.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		if allowRetry && ShouldRetryTransportError(err) && attempt < totalAttempts {
			delay := NextRetryDelay(baseDelay, maxDelay, attempt)
			if onRetry != nil {
				onRetry(HTTPRetryEvent{
					Attempt: attempt,
					Delay:   delay,
					Err:     err,
				})
			}
			if waitErr := WaitForRetry(ctx, delay); waitErr != nil {
				return nil, waitErr
			}
			continue
		}
		return nil, err
	}

	return nil, fmt.Errorf("retry attempts exhausted for request %s", method)
}

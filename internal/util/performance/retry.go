package optimize

import (
	stderrs "errors"
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
		return netErr.Timeout() || netErr.Temporary()
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "timeout")
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

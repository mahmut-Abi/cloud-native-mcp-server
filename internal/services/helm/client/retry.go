// Package client provides Helm client operations for the MCP server.
package client

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// RetryStrategy defines the retry behavior
type RetryStrategy struct {
	MaxRetries        int
	InitialBackoff    time.Duration
	MaxBackoff        time.Duration
	BackoffMultiplier float64
}

// DefaultRetryStrategy returns a default retry strategy
func DefaultRetryStrategy() *RetryStrategy {
	return &RetryStrategy{
		MaxRetries:        3,
		InitialBackoff:    1 * time.Second,
		MaxBackoff:        10 * time.Second,
		BackoffMultiplier: 2.0,
	}
}

// HelmError represents a Helm operation error
type HelmError struct {
	Code       string
	Message    string
	Suggestion string
	Retryable  bool
	Details    map[string]string
}

// Error implements error interface
func (he *HelmError) Error() string {
	return fmt.Sprintf("[%s] %s", he.Code, he.Message)
}

// IsRetryable checks if an error is retryable
func IsRetryable(err error) bool {
	if helmErr, ok := err.(*HelmError); ok {
		return helmErr.Retryable
	}
	return false
}

// Execute retries the function with exponential backoff
func (rs *RetryStrategy) Execute(fn func() error) error {
	var lastErr error
	backoff := rs.InitialBackoff

	for attempt := 0; attempt <= rs.MaxRetries; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// Don't retry if not retryable
		if !IsRetryable(err) {
			return err
		}

		// Don't retry if we've exhausted retries
		if attempt == rs.MaxRetries {
			break
		}

		logrus.WithField("attempt", attempt+1).WithField("backoff", backoff).WithError(err).Debug("Retrying operation")
		time.Sleep(backoff)

		// Exponential backoff with max limit
		backoff = time.Duration(float64(backoff) * rs.BackoffMultiplier)
		if backoff > rs.MaxBackoff {
			backoff = rs.MaxBackoff
		}
	}

	return lastErr
}

// NewHelmError creates a new Helm error
func NewHelmError(code, message, suggestion string, retryable bool) *HelmError {
	return &HelmError{
		Code:       code,
		Message:    message,
		Suggestion: suggestion,
		Retryable:  retryable,
		Details:    make(map[string]string),
	}
}

// Common Helm errors
var (
	ErrRepoUnreachable  = NewHelmError("REPO_UNREACHABLE", "Repository is unreachable", "Check network connectivity and repository URL", true)
	ErrRepoIndexCorrupt = NewHelmError("REPO_INDEX_CORRUPT", "Repository index is corrupted", "Update repository with 'helm repo update'", true)
	ErrTimeout          = NewHelmError("TIMEOUT", "Operation timed out", "Try again or increase timeout value", true)
	ErrReleaseNotFound  = NewHelmError("RELEASE_NOT_FOUND", "Release not found", "Check release name and namespace", false)
	ErrChartNotFound    = NewHelmError("CHART_NOT_FOUND", "Chart not found", "Search charts with 'helm search'", false)
)

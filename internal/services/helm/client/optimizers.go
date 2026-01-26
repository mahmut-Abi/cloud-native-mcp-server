// Package client provides Helm client operations for the MCP server.
package client

import (
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// RepositoryOptimizer handles optimization for Helm repository access in China.
type RepositoryOptimizer struct {
	mirrors   map[string]string
	timeout   time.Duration
	maxRetry  int
	useMirror bool
}

// ChineseMirrors defines common Chinese mirrors for overseas repositories.
var ChineseMirrors = map[string]string{}

// NewRepositoryOptimizer creates a new repository optimizer.
func NewRepositoryOptimizer(mirrors map[string]string, timeoutSec int, maxRetry int, useMirror bool) *RepositoryOptimizer {
	if timeoutSec <= 0 {
		timeoutSec = 300
	}
	if maxRetry <= 0 {
		maxRetry = 3
	}

	optimizer := &RepositoryOptimizer{
		mirrors:   make(map[string]string),
		timeout:   time.Duration(timeoutSec) * time.Second,
		maxRetry:  maxRetry,
		useMirror: useMirror,
	}

	if useMirror {
		for k, v := range ChineseMirrors {
			optimizer.mirrors[k] = v
		}
	}

	for k, v := range mirrors {
		optimizer.mirrors[k] = v
	}

	return optimizer
}

// ResolveRepositoryURL resolves a repository URL using mirrors if available.
func (o *RepositoryOptimizer) ResolveRepositoryURL(originalURL string) string {
	if !o.useMirror {
		return originalURL
	}

	if mirrorURL, ok := o.mirrors[originalURL]; ok {
		logrus.Debugf("Using mirror %q for repository %q", mirrorURL, originalURL)
		return mirrorURL
	}

	return originalURL
}

// CreateOptimizedHTTPClient creates an HTTP client with optimized settings.
func (o *RepositoryOptimizer) CreateOptimizedHTTPClient() *http.Client {
	dialer := &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	return &http.Client{
		Timeout: o.timeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     32,
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
			DialContext:         dialer.DialContext,
		},
	}
}

// GetTimeout returns the configured timeout duration.
func (o *RepositoryOptimizer) GetTimeout() time.Duration {
	return o.timeout
}

// GetMaxRetry returns the configured maximum retry count.
func (o *RepositoryOptimizer) GetMaxRetry() int {
	return o.maxRetry
}

// HasMirror returns whether a mirror is configured for the given URL.
func (o *RepositoryOptimizer) HasMirror(url string) bool {
	_, ok := o.mirrors[url]
	return ok
}

// ListMirrors returns all configured mirrors.
func (o *RepositoryOptimizer) ListMirrors() map[string]string {
	return o.mirrors
}

// IsMirrorEnabled returns whether mirrors are enabled.
func (o *RepositoryOptimizer) IsMirrorEnabled() bool {
	return o.useMirror
}

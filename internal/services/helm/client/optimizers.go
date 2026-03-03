// Package client provides Helm client operations for the MCP server.
package client

import (
	"net/http"
	"net/url"
	"time"

	"helm.sh/helm/v3/pkg/getter"
)

// RepositoryOptimizer holds network configuration for Helm repository operations.
type RepositoryOptimizer struct {
	timeout   time.Duration
	maxRetry  int
	httpProxy string
}

// NewRepositoryOptimizer creates a new repository optimizer.
func NewRepositoryOptimizer(timeoutSec int, maxRetry int, httpProxy string) *RepositoryOptimizer {
	if timeoutSec <= 0 {
		timeoutSec = 300
	}
	if maxRetry <= 0 {
		maxRetry = 3
	}

	return &RepositoryOptimizer{
		timeout:   time.Duration(timeoutSec) * time.Second,
		maxRetry:  maxRetry,
		httpProxy: httpProxy,
	}
}

// GetterOptions returns getter options that apply timeout and optional HTTP proxy.
func (o *RepositoryOptimizer) GetterOptions() ([]getter.Option, error) {
	options := []getter.Option{getter.WithTimeout(o.timeout)}

	if o.httpProxy == "" {
		return options, nil
	}

	parsedProxy, err := url.Parse(o.httpProxy)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyURL(parsedProxy),
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		MaxConnsPerHost:       32,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	options = append(options, getter.WithTransport(transport))
	return options, nil
}

// GetTimeout returns the configured timeout duration.
func (o *RepositoryOptimizer) GetTimeout() time.Duration {
	return o.timeout
}

// GetMaxRetry returns the configured maximum retry count.
func (o *RepositoryOptimizer) GetMaxRetry() int {
	return o.maxRetry
}

// GetHTTPProxy returns the configured dedicated Helm HTTP proxy.
func (o *RepositoryOptimizer) GetHTTPProxy() string {
	return o.httpProxy
}

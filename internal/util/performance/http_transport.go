package optimize

import (
	"net"
	"net/http"
	"time"
)

// HTTPClientConfig holds configuration for HTTP client
type HTTPClientConfig struct {
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	MaxConnsPerHost     int
	IdleConnTimeout     time.Duration
	DialTimeout         time.Duration
	KeepAliveTimeout    time.Duration
	TLSHandshakeTimeout time.Duration
	ResponseTimeout     time.Duration
	ClientTimeout       time.Duration
	EnableHTTP2         bool
}

// DefaultHTTPClientConfig returns sensible defaults for HTTP client
func DefaultHTTPClientConfig() HTTPClientConfig {
	return HTTPClientConfig{
		MaxIdleConns:        256,
		MaxIdleConnsPerHost: 256,
		MaxConnsPerHost:     256,
		IdleConnTimeout:     120 * time.Second,
		DialTimeout:         15 * time.Second,
		KeepAliveTimeout:    60 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		ResponseTimeout:     30 * time.Second,
		ClientTimeout:       60 * time.Second,
		EnableHTTP2:         true,
	}
}

// NewOptimizedHTTPClient creates an HTTP client with optimized transport settings
// for better performance with connection pooling and keepalive
func NewOptimizedHTTPClient() *http.Client {
	return NewConfigurableHTTPClient(DefaultHTTPClientConfig())
}

// NewConfigurableHTTPClient creates an HTTP client with custom configuration
func NewConfigurableHTTPClient(config HTTPClientConfig) *http.Client {
	transport := &http.Transport{
		// Connection pooling
		MaxIdleConns:        config.MaxIdleConns,
		MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
		MaxConnsPerHost:     config.MaxConnsPerHost,

		// Keep-alive settings
		IdleConnTimeout:   config.IdleConnTimeout,
		DisableKeepAlives: false,

		// Timeouts
		DialContext: (&net.Dialer{
			Timeout:   config.DialTimeout,
			KeepAlive: config.KeepAliveTimeout,
		}).DialContext,
		TLSHandshakeTimeout:   config.TLSHandshakeTimeout,
		ExpectContinueTimeout: 1 * time.Second,

		// Response body reading
		ResponseHeaderTimeout: config.ResponseTimeout,

		// Compression
		DisableCompression: false,

		// HTTP/2 support
		ForceAttemptHTTP2: config.EnableHTTP2,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   config.ClientTimeout,
	}
}

// NewOptimizedHTTPClientWithTimeout creates an HTTP client with custom timeout
func NewOptimizedHTTPClientWithTimeout(timeout time.Duration) *http.Client {
	config := DefaultHTTPClientConfig()
	config.ClientTimeout = timeout
	return NewConfigurableHTTPClient(config)
}

// NewHTTPClientForWorkload creates an HTTP client optimized for specific workloads
func NewHTTPClientForWorkload(workloadType string) *http.Client {
	config := DefaultHTTPClientConfig()

	switch workloadType {
	case "highConcurrency":
		// Optimize for many simultaneous connections
		config.MaxIdleConns = 500
		config.MaxIdleConnsPerHost = 100
		config.MaxConnsPerHost = 500
		config.IdleConnTimeout = 180 * time.Second
	case "longRunning":
		// Optimize for long-running connections
		config.IdleConnTimeout = 300 * time.Second
		config.KeepAliveTimeout = 120 * time.Second
		config.ClientTimeout = 300 * time.Second
	case "lowLatency":
		// Optimize for low latency responses
		config.MaxIdleConnsPerHost = 50
		config.DialTimeout = 5 * time.Second
		config.ResponseTimeout = 10 * time.Second
		config.ClientTimeout = 20 * time.Second
	case "lowMemory":
		// Optimize for low memory usage
		config.MaxIdleConns = 50
		config.MaxIdleConnsPerHost = 10
		config.MaxConnsPerHost = 50
	}

	return NewConfigurableHTTPClient(config)
}

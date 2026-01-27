// Package constants defines common constants used across the k8s-mcp-server
package constants

import "time"

// Cache constants
const (
	// DefaultCacheSize is the default maximum number of entries in the cache
	DefaultCacheSize = 10000

	// DefaultCacheTTL is the default time-to-live for cache entries
	DefaultCacheTTL = 15 * time.Minute

	// CacheCleanupInterval is the interval between cache cleanup operations
	CacheCleanupInterval = 5 * time.Minute
)

// Rate limiting constants
const (
	// DefaultRateLimitRPS is the default requests per second limit
	DefaultRateLimitRPS = 10

	// DefaultRateLimitBurst is the default burst size for rate limiting
	DefaultRateLimitBurst = 20

	// RateLimitCleanupThreshold is the number of entries that triggers cleanup
	RateLimitCleanupThreshold = 10000

	// RateLimitCleanupBatch is the number of entries to clean in one batch
	RateLimitCleanupBatch = 1000

	// RateLimitStaleDuration is the duration after which an entry is considered stale
	RateLimitStaleDuration = time.Hour

	// RateLimitCleanupInterval is the interval between rate limit cleanup operations
	RateLimitCleanupInterval = 5 * time.Minute
)

// Kubernetes constants
const (
	// DefaultTailLines is the default number of log lines to retrieve
	DefaultTailLines = 50

	// DefaultLimit is the default number of resources to list
	DefaultLimit = 30

	// MaxLimit is the maximum number of resources to return in a single request
	MaxLimit = 80

	// WarningLimit is the limit at which a warning is logged for large requests
	WarningLimit = 40

	// MaxLogLines is the maximum number of log lines that can be requested
	MaxLogLines = 200

	// LogSizeThreshold is the size threshold (in bytes) for log truncation
	LogSizeThreshold = 10000 // 10KB

	// TruncatedLogLines is the maximum number of lines to keep after truncation
	TruncatedLogLines = 200

	// MaxLogCharacters is the maximum number of characters in log output
	MaxLogCharacters = 50000 // 50KB
)

// HTTP constants
const (
	// DefaultHTTPTimeout is the default timeout for HTTP requests
	DefaultHTTPTimeout = 30 * time.Second

	// DefaultDialTimeout is the default timeout for establishing TCP connections
	DefaultDialTimeout = 15 * time.Second

	// DefaultKeepAliveTimeout is the default keep-alive timeout for connections
	DefaultKeepAliveTimeout = 60 * time.Second

	// DefaultIdleConnTimeout is the default timeout for idle connections
	DefaultIdleConnTimeout = 120 * time.Second

	// DefaultTLSHandshakeTimeout is the default timeout for TLS handshakes
	DefaultTLSHandshakeTimeout = 10 * time.Second

	// MaxIdleConns is the maximum number of idle connections across all hosts
	MaxIdleConns = 256

	// MaxIdleConnsPerHost is the maximum number of idle connections per host
	MaxIdleConnsPerHost = 256

	// MaxConnsPerHost is the maximum number of connections per host
	MaxConnsPerHost = 256
)

// Pagination constants
const (
	// DefaultPageSize is the default page size for paginated results
	DefaultPageSize = 20

	// MaxPageSize is the maximum page size allowed
	MaxPageSize = 100

	// WarningPageSize is the limit at which a warning is logged for large requests
	WarningPageSize = 50

	// DefaultPageSizeNodes is the default page size for node-related requests
	DefaultPageSizeNodes = 50
)

// Circuit breaker constants
const (
	// DefaultMaxFailures is the default number of failures before opening the circuit
	DefaultMaxFailures = 5

	// DefaultCircuitTimeout is the default timeout for the circuit breaker
	DefaultCircuitTimeout = 30 * time.Second

	// DefaultHalfOpenAttempts is the default number of attempts in half-open state
	DefaultHalfOpenAttempts = 3
)

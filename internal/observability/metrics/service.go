package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// ToolCallsTotal counts total tool calls
	ToolCallsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tool_calls_total",
			Help: "Total number of tool calls",
		},
		[]string{"service_name", "tool_name", "status"},
	)

	// ToolCallDuration tracks tool call duration
	ToolCallDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tool_call_duration_seconds",
			Help:    "Tool call duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service_name", "tool_name", "status"},
	)

	// ExternalAPICallsTotal counts total external API calls
	ExternalAPICallsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "external_api_calls_total",
			Help: "Total number of external API calls",
		},
		[]string{"service_name", "api_name", "status"},
	)

	// ExternalAPICallDuration tracks external API call duration
	ExternalAPICallDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "external_api_call_duration_seconds",
			Help:    "External API call duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service_name", "api_name", "status"},
	)

	// CacheHitsTotal counts cache hits
	CacheHitsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"service_name", "cache_type"},
	)

	// CacheMissesTotal counts cache misses
	CacheMissesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"service_name", "cache_type"},
	)

	// CircuitBreakerState tracks circuit breaker state
	CircuitBreakerState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "circuit_breaker_state",
			Help: "Circuit breaker state (0=closed, 1=open, 2=half_open)",
		},
		[]string{"service_name", "circuit_breaker_name"},
	)

	// CircuitBreakerFailures counts circuit breaker failures
	CircuitBreakerFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "circuit_breaker_failures_total",
			Help: "Total number of circuit breaker failures",
		},
		[]string{"service_name", "circuit_breaker_name"},
	)
)

// RecordToolCall records a tool call metric
func RecordToolCall(serviceName, toolName, status string, duration float64) {
	ToolCallsTotal.WithLabelValues(serviceName, toolName, status).Inc()
	ToolCallDuration.WithLabelValues(serviceName, toolName, status).Observe(duration)
}

// RecordExternalAPICall records an external API call metric
func RecordExternalAPICall(serviceName, apiName, status string, duration float64) {
	ExternalAPICallsTotal.WithLabelValues(serviceName, apiName, status).Inc()
	ExternalAPICallDuration.WithLabelValues(serviceName, apiName, status).Observe(duration)
}

// RecordCacheHit records a cache hit
func RecordCacheHit(serviceName, cacheType string) {
	CacheHitsTotal.WithLabelValues(serviceName, cacheType).Inc()
}

// RecordCacheMiss records a cache miss
func RecordCacheMiss(serviceName, cacheType string) {
	CacheMissesTotal.WithLabelValues(serviceName, cacheType).Inc()
}

// SetCircuitBreakerState sets the circuit breaker state
func SetCircuitBreakerState(serviceName, circuitBreakerName string, state float64) {
	CircuitBreakerState.WithLabelValues(serviceName, circuitBreakerName).Set(state)
}

// RecordCircuitBreakerFailure records a circuit breaker failure
func RecordCircuitBreakerFailure(serviceName, circuitBreakerName string) {
	CircuitBreakerFailures.WithLabelValues(serviceName, circuitBreakerName).Inc()
}

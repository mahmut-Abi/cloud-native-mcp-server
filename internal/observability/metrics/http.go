package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTPRequestsTotal counts total HTTP requests
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status_code", "service"},
	)

	// HTTPRequestDuration tracks HTTP request duration
	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status_code", "service"},
	)

	// HTTPRequestSize tracks HTTP request size
	HTTPRequestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 7), // 100B, 1KB, 10KB, 100KB, 1MB, 10MB, 100MB
		},
		[]string{"method", "path", "service"},
	)

	// HTTPResponseSize tracks HTTP response size
	HTTPResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 7), // 100B, 1KB, 10KB, 100KB, 1MB, 10MB, 100MB
		},
		[]string{"method", "path", "status_code", "service"},
	)

	// HTTPConnectionsActive tracks active HTTP connections
	HTTPConnectionsActive = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_connections_active",
			Help: "Number of active HTTP connections",
		},
	)
)

// Handler returns the Prometheus metrics HTTP handler
func Handler() http.Handler {
	return promhttp.HandlerFor(Registry, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	})
}

// RecordHTTPRequest records an HTTP request metric
func RecordHTTPRequest(method, path, statusCode, service string, duration float64, requestSize, responseSize int64) {
	HTTPRequestsTotal.WithLabelValues(method, path, statusCode, service).Inc()
	HTTPRequestDuration.WithLabelValues(method, path, statusCode, service).Observe(duration)

	if requestSize > 0 {
		HTTPRequestSize.WithLabelValues(method, path, service).Observe(float64(requestSize))
	}

	if responseSize > 0 {
		HTTPResponseSize.WithLabelValues(method, path, statusCode, service).Observe(float64(responseSize))
	}
}

// IncActiveConnections increments the active connections counter
func IncActiveConnections() {
	HTTPConnectionsActive.Inc()
}

// DecActiveConnections decrements the active connections counter
func DecActiveConnections() {
	HTTPConnectionsActive.Dec()
}

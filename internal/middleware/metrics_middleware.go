package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/mahmut-Abi/k8s-mcp-server/internal/observability/metrics"
	"github.com/sirupsen/logrus"
)

// MetricsMiddleware is an HTTP middleware that records Prometheus metrics
func MetricsMiddleware(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Track active connections
			metrics.IncActiveConnections()
			defer metrics.DecActiveConnections()

			// Wrap response writer to capture status code and response size
			rw := &responseWriterWrapper{
				ResponseWriter: w,
				headerWritten:  false,
				statusCode:     http.StatusOK, // Default to 200
				size:           0,             // Initialize size to 0
			}

			// Get request size
			requestSize := r.ContentLength
			if requestSize < 0 {
				requestSize = 0
			}

			// Call next handler
			next.ServeHTTP(rw, r)

			// Calculate duration
			duration := time.Since(start).Seconds()

			// Get response size
			responseSize := int64(rw.size)

			// Extract service name from path if not provided
			service := serviceName
			if service == "" {
				service = extractServiceName(r.URL.Path)
			}

			// Record metrics
			statusCode := strconv.Itoa(rw.statusCode)
			metrics.RecordHTTPRequest(
				r.Method,
				r.URL.Path,
				statusCode,
				service,
				duration,
				requestSize,
				responseSize,
			)

			logrus.WithFields(logrus.Fields{
				"component": "metrics",
				"method":    r.Method,
				"path":      r.URL.Path,
				"status":    rw.statusCode,
				"duration":  duration,
				"service":   service,
			}).Debug("HTTP request recorded")
		})
	}
}

// extractServiceName extracts service name from URL path
func extractServiceName(path string) string {
	// Extract service name from path like /api/kubernetes/sse -> kubernetes
	// /api/grafana/sse -> grafana
	// /api/aggregate/sse -> aggregate
	// /health -> health
	// /metrics -> metrics

	if path == "/health" || path == "/metrics" || path == "/ready" {
		return "system"
	}

	// Check for API paths
	if len(path) > 5 && path[:5] == "/api/" {
		remaining := path[5:]
		// Find the next slash
		for i, c := range remaining {
			if c == '/' {
				return remaining[:i]
			}
		}
		return remaining
	}

	return "unknown"
}

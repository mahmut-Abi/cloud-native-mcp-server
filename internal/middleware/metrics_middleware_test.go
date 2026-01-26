package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mahmut-Abi/k8s-mcp-server/internal/observability/metrics"
	dto "github.com/prometheus/client_model/go"
)

func init() {
	// Initialize metrics for all tests
	metrics.Init("test", "test", "go1.24", "sse", "0.0.0.0:8080")
}

func TestMetricsMiddleware(t *testing.T) {
	// Create a simple handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test response"))
	})

	// Wrap with metrics middleware (no service name, will extract from path)
	middleware := MetricsMiddleware("")(handler)

	tests := []struct {
		name           string
		method         string
		path           string
		wantStatus     int
		expectedLabels map[string]string
	}{
		{
			name:       "GET request",
			method:     "GET",
			path:       "/api/test",
			wantStatus: http.StatusOK,
			expectedLabels: map[string]string{
				"service": "test",
			},
		},
		{
			name:       "POST request",
			method:     "POST",
			path:       "/api/data",
			wantStatus: http.StatusOK,
			expectedLabels: map[string]string{
				"service": "data",
			},
		},
		{
			name: "health endpoint", method: "GET",
			path:       "/health",
			wantStatus: http.StatusOK,
			expectedLabels: map[string]string{
				"service": "system",
			},
		},
		{
			name:       "metrics endpoint",
			method:     "GET",
			path:       "/metrics",
			wantStatus: http.StatusOK,
			expectedLabels: map[string]string{
				"service": "system",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			middleware.ServeHTTP(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("ServeHTTP() status = %v, want %v", resp.StatusCode, tt.wantStatus)
			}

			// Verify metrics were recorded
			metricsList, err := metrics.Registry.Gather()
			if err != nil {
				t.Fatalf("Failed to gather metrics: %v", err)
			}

			// Check http_requests_total
			found := false
			for _, m := range metricsList {
				if m.GetName() == "http_requests_total" {
					for _, metric := range m.GetMetric() {
						labels := metric.GetLabel()
						if hasLabels(labels, "method", tt.method) &&
							hasLabels(labels, "path", tt.path) &&
							hasLabels(labels, "status_code", "200") &&
							hasLabelsFromMap(labels, tt.expectedLabels) {
							found = true
							if metric.Counter.GetValue() >= 1 {
								// Counter may be > 1 due to previous tests
								break
							}
						}
					}
					break
				}
			}

			if !found {
				t.Error("http_requests_total metric not found")
			}
		})
	}
}

func TestMetricsMiddleware_WithRequestBody(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("created"))
	})

	middleware := MetricsMiddleware("")(handler)

	body := strings.NewReader(`{"test": "data"}`)
	req := httptest.NewRequest("POST", "/api/create", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("ServeHTTP() status = %v, want %v", resp.StatusCode, http.StatusCreated)
	}

	// Verify metrics
	metricsList, err := metrics.Registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// Check http_request_size_bytes was recorded
	found := false
	for _, m := range metricsList {
		if m.GetName() == "http_request_size_bytes" {
			found = true
			for _, metric := range m.GetMetric() {
				labels := metric.GetLabel()
				if hasLabels(labels, "method", "POST") &&
					hasLabels(labels, "path", "/api/create") &&
					hasLabels(labels, "service", "create") {
					if metric.Histogram.GetSampleCount() != 1 {
						t.Errorf("http_request_size_bytes count = %v, want 1", metric.Histogram.GetSampleCount())
					}
					break
				}
			}
			break
		}
	}
	if !found {
		t.Error("http_request_size_bytes metric not found")
	}
}

func TestMetricsMiddleware_WithErrorResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error"))
	})

	middleware := MetricsMiddleware("")(handler)

	req := httptest.NewRequest("GET", "/api/error", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("ServeHTTP() status = %v, want %v", resp.StatusCode, http.StatusInternalServerError)
	}

	// Verify metrics
	metricsList, err := metrics.Registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// Check http_requests_total with 500 status
	found := false
	for _, m := range metricsList {
		if m.GetName() == "http_requests_total" {
			for _, metric := range m.GetMetric() {
				labels := metric.GetLabel()
				if hasLabels(labels, "method", "GET") &&
					hasLabels(labels, "path", "/api/error") &&
					hasLabels(labels, "status_code", "500") &&
					hasLabels(labels, "service", "error") {
					found = true
					if metric.Counter.GetValue() >= 1 {
						break
					}
				}
			}
			break
		}
	}

	if !found {
		t.Error("http_requests_total metric with 500 status not found")
	}
}

func TestMetricsMiddleware_ActiveConnections(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := MetricsMiddleware("")(handler)

	// Make multiple concurrent requests
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/api/test", nil)
		w := httptest.NewRecorder()
		middleware.ServeHTTP(w, req)
	}

	// Check that active connections metric exists and was incremented/decremented
	metricsList, err := metrics.Registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	found := false
	for _, m := range metricsList {
		if m.GetName() == "http_connections_active" {
			found = true
			// After all requests complete, should be back to 0
			if len(m.GetMetric()) > 0 {
				value := m.GetMetric()[0].GetGauge().GetValue()
				if value != 0 {
					t.Errorf("http_connections_active = %v, want 0 after all requests", value)
				}
			}
			break
		}
	}
	if !found {
		t.Error("http_connections_active metric not found")
	}
}

func TestExtractServiceName(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "kubernetes API",
			path: "/api/kubernetes/pods",
			want: "kubernetes",
		},
		{
			name: "grafana API",
			path: "/api/grafana/dashboards",
			want: "grafana",
		},
		{
			name: "aggregate API",
			path: "/api/aggregate/sse",
			want: "aggregate",
		},
		{
			name: "health endpoint",
			path: "/health",
			want: "system",
		},
		{
			name: "metrics endpoint",
			path: "/metrics",
			want: "system",
		},
		{
			name: "root path",
			path: "/",
			want: "unknown",
		},
		{
			name: "unknown path",
			path: "/unknown/path",
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractServiceName(tt.path)
			if got != tt.want {
				t.Errorf("extractServiceName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricsMiddleware_ResponseSize(t *testing.T) {
	responseBody := "test response body"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responseBody))
	})

	middleware := MetricsMiddleware("")(handler)

	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	// Verify http_response_size_bytes was recorded
	metricsList, err := metrics.Registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	found := false
	for _, m := range metricsList {
		if m.GetName() == "http_response_size_bytes" {
			found = true
			for _, metric := range m.GetMetric() {
				labels := metric.GetLabel()
				if hasLabels(labels, "method", "GET") &&
					hasLabels(labels, "path", "/api/test") &&
					hasLabels(labels, "status_code", "200") &&
					hasLabels(labels, "service", "test") {
					// Check that we have at least one sample
					if metric.Histogram.GetSampleCount() < 1 {
						t.Errorf("http_response_size_bytes count = %v, want >= 1", metric.Histogram.GetSampleCount())
					}
					break
				}
			}
			break
		}
	}
	if !found {
		t.Error("http_response_size_bytes metric not found")
	}
}

func hasLabels(labels []*dto.LabelPair, name, value string) bool {
	for _, label := range labels {
		if label.GetName() == name && label.GetValue() == value {
			return true
		}
	}
	return false
}

func hasLabelsFromMap(labels []*dto.LabelPair, expected map[string]string) bool {
	for k, v := range expected {
		if !hasLabels(labels, k, v) {
			return false
		}
	}
	return true
}

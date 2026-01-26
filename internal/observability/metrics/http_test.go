package metrics

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func TestHandler(t *testing.T) {
	// Reset registry for clean test
	Registry = prometheus.NewRegistry()
	Init("test", "test", "go1.25", "sse", "0.0.0.0:8080")

	handler := Handler()
	if handler == nil {
		t.Fatal("Handler() returned nil")
	}

	// Create a test request
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Check response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Handler() status = %v, want %v", resp.StatusCode, http.StatusOK)
	}

	body, _ := io.ReadAll(resp.Body)
	if len(body) == 0 {
		t.Error("Handler() returned empty body")
	}

	// Check that response contains expected metrics
	bodyStr := string(body)
	expectedMetrics := []string{
		"build_info",
		"server_info",
		"go_goroutines",
		"go_threads",
	}

	for _, metric := range expectedMetrics {
		if !strings.Contains(bodyStr, metric) {
			t.Errorf("Handler() response missing metric: %s", metric)
		}
	}
}

func TestRecordHTTPRequest(t *testing.T) {
	// Reset registry for clean test
	Registry = prometheus.NewRegistry()
	Init("test", "test", "go1.25", "sse", "0.0.0.0:8080")

	tests := []struct {
		name         string
		method       string
		path         string
		statusCode   string
		service      string
		duration     float64
		requestSize  int64
		responseSize int64
	}{
		{
			name:         "successful GET request",
			method:       "GET",
			path:         "/api/kubernetes/pods",
			statusCode:   "200",
			service:      "kubernetes",
			duration:     0.123,
			requestSize:  100,
			responseSize: 1000,
		},
		{
			name:         "failed POST request",
			method:       "POST",
			path:         "/api/grafana/dashboards",
			statusCode:   "500",
			service:      "grafana",
			duration:     0.456,
			requestSize:  500,
			responseSize: 0,
		},
		{
			name:         "zero sizes",
			method:       "GET",
			path:         "/health",
			statusCode:   "200",
			service:      "system",
			duration:     0.001,
			requestSize:  0,
			responseSize: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RecordHTTPRequest(tt.method, tt.path, tt.statusCode, tt.service, tt.duration, tt.requestSize, tt.responseSize)

			// Gather metrics and verify
			metrics, err := Registry.Gather()
			if err != nil {
				t.Fatalf("Failed to gather metrics: %v", err)
			}

			// Check http_requests_total
			foundRequests := false
			for _, m := range metrics {
				if m.GetName() == "http_requests_total" {
					foundRequests = true
					found := false
					for _, metric := range m.GetMetric() {
						labels := metric.GetLabel()
						if hasLabels(labels, "method", tt.method) &&
							hasLabels(labels, "path", tt.path) &&
							hasLabels(labels, "status_code", tt.statusCode) &&
							hasLabels(labels, "service", tt.service) {
							found = true
							if metric.Counter.GetValue() != 1 {
								t.Errorf("http_requests_total = %v, want 1", metric.Counter.GetValue())
							}
							break
						}
					}
					if !found {
						t.Error("http_requests_total metric with expected labels not found")
					}
					break
				}
			}
			if !foundRequests {
				t.Error("http_requests_total metric not found")
			}

			// Check http_request_duration_seconds
			foundDuration := false
			for _, m := range metrics {
				if m.GetName() == "http_request_duration_seconds" {
					foundDuration = true
					found := false
					for _, metric := range m.GetMetric() {
						labels := metric.GetLabel()
						if hasLabels(labels, "method", tt.method) &&
							hasLabels(labels, "path", tt.path) &&
							hasLabels(labels, "status_code", tt.statusCode) &&
							hasLabels(labels, "service", tt.service) {
							found = true
							if metric.Histogram.GetSampleSum() != tt.duration {
								t.Errorf("http_request_duration_seconds sum = %v, want %v", metric.Histogram.GetSampleSum(), tt.duration)
							}
							break
						}
					}
					if !found {
						t.Error("http_request_duration_seconds metric with expected labels not found")
					}
					break
				}
			}
			if !foundDuration {
				t.Error("http_request_duration_seconds metric not found")
			}
		})
	}
}

func TestIncDecActiveConnections(t *testing.T) {
	// Reset registry for clean test
	Registry = prometheus.NewRegistry()
	Init("test", "test", "go1.25", "sse", "0.0.0.0:8080")

	// Initial value should be 0
	initialValue := getGaugeValue("http_connections_active")
	if initialValue != 0 {
		t.Errorf("Initial connections = %v, want 0", initialValue)
	}

	// Increment
	IncActiveConnections()
	value1 := getGaugeValue("http_connections_active")
	if value1 != 1 {
		t.Errorf("After increment, connections = %v, want 1", value1)
	}

	// Increment again
	IncActiveConnections()
	value2 := getGaugeValue("http_connections_active")
	if value2 != 2 {
		t.Errorf("After second increment, connections = %v, want 2", value2)
	}

	// Decrement
	DecActiveConnections()
	value3 := getGaugeValue("http_connections_active")
	if value3 != 1 {
		t.Errorf("After decrement, connections = %v, want 1", value3)
	}

	// Decrement to zero
	DecActiveConnections()
	value4 := getGaugeValue("http_connections_active")
	if value4 != 0 {
		t.Errorf("After second decrement, connections = %v, want 0", value4)
	}
}

func TestRecordHTTPRequest_MultipleRequests(t *testing.T) {
	// Reset registry for clean test
	Registry = prometheus.NewRegistry()
	Init("test", "test", "go1.25", "sse", "0.0.0.0:8080")

	// Record multiple requests
	for i := 0; i < 5; i++ {
		RecordHTTPRequest("GET", "/api/test", "200", "test", 0.1, 100, 1000)
	}

	// Check that counter is 5
	metrics, err := Registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	for _, m := range metrics {
		if m.GetName() == "http_requests_total" {
			for _, metric := range m.GetMetric() {
				labels := metric.GetLabel()
				if hasLabels(labels, "method", "GET") &&
					hasLabels(labels, "path", "/api/test") &&
					hasLabels(labels, "status_code", "200") &&
					hasLabels(labels, "service", "test") {
					if metric.Counter.GetValue() != 5 {
						t.Errorf("http_requests_total = %v, want 5", metric.Counter.GetValue())
					}
					return
				}
			}
		}
	}
	t.Error("Expected metric not found")
}

func hasLabels(labels []*dto.LabelPair, name, value string) bool {
	for _, label := range labels {
		if label.GetName() == name && label.GetValue() == value {
			return true
		}
	}
	return false
}

func getGaugeValue(metricName string) float64 {
	metrics, _ := Registry.Gather()
	for _, m := range metrics {
		if m.GetName() == metricName {
			if len(m.GetMetric()) > 0 {
				return m.GetMetric()[0].GetGauge().GetValue()
			}
		}
	}
	return 0
}

package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestRecordToolCall(t *testing.T) {
	// Reset registry for clean test
	Registry = prometheus.NewRegistry()
	Init("test", "test", "go1.24", "sse", "0.0.0.0:8080")

	tests := []struct {
		name        string
		serviceName string
		toolName    string
		status      string
		duration    float64
	}{
		{
			name:        "successful tool call",
			serviceName: "kubernetes",
			toolName:    "kubernetes_list_pods",
			status:      "success",
			duration:    0.123,
		},
		{
			name:        "failed tool call",
			serviceName: "helm",
			toolName:    "helm_list_releases",
			status:      "error",
			duration:    0.456,
		},
		{
			name:        "unknown service",
			serviceName: "unknown",
			toolName:    "unknown_tool",
			status:      "success",
			duration:    0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RecordToolCall(tt.serviceName, tt.toolName, tt.status, tt.duration)

			// Gather metrics and verify
			metrics, err := Registry.Gather()
			if err != nil {
				t.Fatalf("Failed to gather metrics: %v", err)
			}

			// Check tool_calls_total
			foundTotal := false
			for _, m := range metrics {
				if m.GetName() == "tool_calls_total" {
					foundTotal = true
					found := false
					for _, metric := range m.GetMetric() {
						labels := metric.GetLabel()
						if hasLabels(labels, "service_name", tt.serviceName) &&
							hasLabels(labels, "tool_name", tt.toolName) &&
							hasLabels(labels, "status", tt.status) {
							found = true
							if metric.Counter.GetValue() != 1 {
								t.Errorf("tool_calls_total = %v, want 1", metric.Counter.GetValue())
							}
							break
						}
					}
					if !found {
						t.Error("tool_calls_total metric with expected labels not found")
					}
					break
				}
			}
			if !foundTotal {
				t.Error("tool_calls_total metric not found")
			}

			// Check tool_call_duration_seconds
			foundDuration := false
			for _, m := range metrics {
				if m.GetName() == "tool_call_duration_seconds" {
					foundDuration = true
					found := false
					for _, metric := range m.GetMetric() {
						labels := metric.GetLabel()
						if hasLabels(labels, "service_name", tt.serviceName) &&
							hasLabels(labels, "tool_name", tt.toolName) &&
							hasLabels(labels, "status", tt.status) {
							found = true
							if metric.Histogram.GetSampleSum() != tt.duration {
								t.Errorf("tool_call_duration_seconds sum = %v, want %v", metric.Histogram.GetSampleSum(), tt.duration)
							}
							break
						}
					}
					if !found {
						t.Error("tool_call_duration_seconds metric with expected labels not found")
					}
					break
				}
			}
			if !foundDuration {
				t.Error("tool_call_duration_seconds metric not found")
			}
		})
	}
}

func TestRecordToolCall_MultipleCalls(t *testing.T) {
	// Reset registry for clean test
	Registry = prometheus.NewRegistry()
	Init("test", "test", "go1.24", "sse", "0.0.0.0:8080")

	// Record multiple calls
	for i := 0; i < 3; i++ {
		RecordToolCall("kubernetes", "kubernetes_list_pods", "success", 0.1)
	}

	// Check that counter is 3
	metrics, err := Registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	for _, m := range metrics {
		if m.GetName() == "tool_calls_total" {
			for _, metric := range m.GetMetric() {
				labels := metric.GetLabel()
				if hasLabels(labels, "service_name", "kubernetes") &&
					hasLabels(labels, "tool_name", "kubernetes_list_pods") &&
					hasLabels(labels, "status", "success") {
					if metric.Counter.GetValue() < 3 {
						t.Errorf("tool_calls_total = %v, want >= 3", metric.Counter.GetValue())
					}
					return
				}
			}
		}
	}
	t.Error("Expected metric not found")
}

func TestRecordExternalAPICall(t *testing.T) {
	// Reset registry for clean test
	Registry = prometheus.NewRegistry()
	Init("test", "test", "go1.24", "sse", "0.0.0.0:8080")

	tests := []struct {
		name        string
		serviceName string
		apiName     string
		status      string
		duration    float64
	}{
		{
			name:        "successful API call",
			serviceName: "grafana",
			apiName:     "dashboards",
			status:      "success",
			duration:    0.234,
		},
		{
			name:        "failed API call",
			serviceName: "prometheus",
			apiName:     "query",
			status:      "error",
			duration:    0.567,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RecordExternalAPICall(tt.serviceName, tt.apiName, tt.status, tt.duration)

			// Gather metrics and verify
			metrics, err := Registry.Gather()
			if err != nil {
				t.Fatalf("Failed to gather metrics: %v", err)
			}

			// Check external_api_calls_total
			found := false
			for _, m := range metrics {
				if m.GetName() == "external_api_calls_total" {
					found = true
					for _, metric := range m.GetMetric() {
						labels := metric.GetLabel()
						if hasLabels(labels, "service_name", tt.serviceName) &&
							hasLabels(labels, "api_name", tt.apiName) &&
							hasLabels(labels, "status", tt.status) {
							if metric.Counter.GetValue() != 1 {
								t.Errorf("external_api_calls_total = %v, want 1", metric.Counter.GetValue())
							}
							return
						}
					}
					t.Error("external_api_calls_total metric with expected labels not found")
					break
				}
			}
			if !found {
				t.Error("external_api_calls_total metric not found")
			}
		})
	}
}

func TestRecordCacheHit(t *testing.T) {
	// Reset registry for clean test
	Registry = prometheus.NewRegistry()
	Init("test", "test", "go1.24", "sse", "0.0.0.0:8080")

	RecordCacheHit("kubernetes", "tools")

	metrics, err := Registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	found := false
	for _, m := range metrics {
		if m.GetName() == "cache_hits_total" {
			found = true
			for _, metric := range m.GetMetric() {
				labels := metric.GetLabel()
				if hasLabels(labels, "service_name", "kubernetes") &&
					hasLabels(labels, "cache_type", "tools") {
					if metric.Counter.GetValue() != 1 {
						t.Errorf("cache_hits_total = %v, want 1", metric.Counter.GetValue())
					}
					return
				}
			}
			t.Error("cache_hits_total metric with expected labels not found")
			break
		}
	}
	if !found {
		t.Error("cache_hits_total metric not found")
	}
}

func TestRecordCacheMiss(t *testing.T) {
	// Reset registry for clean test
	Registry = prometheus.NewRegistry()
	Init("test", "test", "go1.24", "sse", "0.0.0.0:8080")

	RecordCacheMiss("grafana", "dashboards")

	metrics, err := Registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	found := false
	for _, m := range metrics {
		if m.GetName() == "cache_misses_total" {
			found = true
			for _, metric := range m.GetMetric() {
				labels := metric.GetLabel()
				if hasLabels(labels, "service_name", "grafana") &&
					hasLabels(labels, "cache_type", "dashboards") {
					if metric.Counter.GetValue() != 1 {
						t.Errorf("cache_misses_total = %v, want 1", metric.Counter.GetValue())
					}
					return
				}
			}
			t.Error("cache_misses_total metric with expected labels not found")
			break
		}
	}
	if !found {
		t.Error("cache_misses_total metric not found")
	}
}

func TestSetCircuitBreakerState(t *testing.T) {
	// Reset registry for clean test
	Registry = prometheus.NewRegistry()
	Init("test", "test", "go1.24", "sse", "0.0.0.0:8080")

	tests := []struct {
		name               string
		serviceName        string
		circuitBreakerName string
		state              float64
	}{
		{
			name:               "closed state",
			serviceName:        "kubernetes",
			circuitBreakerName: "api",
			state:              0.0,
		},
		{
			name:               "open state",
			serviceName:        "grafana",
			circuitBreakerName: "dashboards",
			state:              1.0,
		},
		{
			name:               "half_open state",
			serviceName:        "prometheus",
			circuitBreakerName: "query",
			state:              2.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetCircuitBreakerState(tt.serviceName, tt.circuitBreakerName, tt.state)

			metrics, err := Registry.Gather()
			if err != nil {
				t.Fatalf("Failed to gather metrics: %v", err)
			}

			found := false
			for _, m := range metrics {
				if m.GetName() == "circuit_breaker_state" {
					found = true
					for _, metric := range m.GetMetric() {
						labels := metric.GetLabel()
						if hasLabels(labels, "service_name", tt.serviceName) &&
							hasLabels(labels, "circuit_breaker_name", tt.circuitBreakerName) {
							if metric.Gauge.GetValue() != tt.state {
								t.Errorf("circuit_breaker_state = %v, want %v", metric.Gauge.GetValue(), tt.state)
							}
							return
						}
					}
					t.Error("circuit_breaker_state metric with expected labels not found")
					break
				}
			}
			if !found {
				t.Error("circuit_breaker_state metric not found")
			}
		})
	}
}

func TestRecordCircuitBreakerFailure(t *testing.T) {
	// Reset registry for clean test
	Registry = prometheus.NewRegistry()
	Init("test", "test", "go1.24", "sse", "0.0.0.0:8080")

	RecordCircuitBreakerFailure("kubernetes", "api")
	RecordCircuitBreakerFailure("kubernetes", "api")
	RecordCircuitBreakerFailure("kubernetes", "api")

	metrics, err := Registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	found := false
	for _, m := range metrics {
		if m.GetName() == "circuit_breaker_failures_total" {
			found = true
			for _, metric := range m.GetMetric() {
				labels := metric.GetLabel()
				if hasLabels(labels, "service_name", "kubernetes") &&
					hasLabels(labels, "circuit_breaker_name", "api") {
					if metric.Counter.GetValue() != 3 {
						t.Errorf("circuit_breaker_failures_total = %v, want 3", metric.Counter.GetValue())
					}
					return
				}
			}
			t.Error("circuit_breaker_failures_total metric with expected labels not found")
			break
		}
	}
	if !found {
		t.Error("circuit_breaker_failures_total metric not found")
	}
}

func TestCacheMetrics_MultipleOperations(t *testing.T) {
	// Reset registry for clean test
	Registry = prometheus.NewRegistry()
	Init("test", "test", "go1.24", "sse", "0.0.0.0:8080")

	// Record multiple cache operations
	RecordCacheHit("kubernetes", "tools")
	RecordCacheHit("kubernetes", "tools")
	RecordCacheMiss("kubernetes", "tools")
	RecordCacheMiss("kubernetes", "tools")
	RecordCacheMiss("kubernetes", "tools")

	metrics, err := Registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// Check cache hits
	hitCount := 0
	for _, m := range metrics {
		if m.GetName() == "cache_hits_total" {
			for _, metric := range m.GetMetric() {
				labels := metric.GetLabel()
				if hasLabels(labels, "service_name", "kubernetes") &&
					hasLabels(labels, "cache_type", "tools") {
					hitCount = int(metric.Counter.GetValue())
				}
			}
		}
	}

	// Check cache misses
	missCount := 0
	for _, m := range metrics {
		if m.GetName() == "cache_misses_total" {
			for _, metric := range m.GetMetric() {
				labels := metric.GetLabel()
				if hasLabels(labels, "service_name", "kubernetes") &&
					hasLabels(labels, "cache_type", "tools") {
					missCount = int(metric.Counter.GetValue())
				}
			}
		}
	}

	if hitCount < 2 {
		t.Errorf("cache hits = %d, want at least 2", hitCount)
	}
	if missCount < 3 {
		t.Errorf("cache misses = %d, want at least 3", missCount)
	}
}

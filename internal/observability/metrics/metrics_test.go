package metrics

import (
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name       string
		version    string
		commit     string
		goVersion  string
		mode       string
		addr       string
		wantLabels map[string]string
	}{
		{
			name:      "basic initialization",
			version:   "1.0.0",
			commit:    "abc123",
			goVersion: "go1.24",
			mode:      "sse",
			addr:      "0.0.0.0:8080",
			wantLabels: map[string]string{
				"version":    "1.0.0",
				"commit":     "abc123",
				"go_version": "go1.24",
			},
		},
		{
			name:      "streamable-http mode",
			version:   "2.0.0",
			commit:    "def456",
			goVersion: "go1.23",
			mode:      "streamable-http",
			addr:      "0.0.0.0:9090",
			wantLabels: map[string]string{
				"version":    "2.0.0",
				"commit":     "def456",
				"go_version": "go1.23",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.version, tt.commit, tt.goVersion, tt.mode, tt.addr)

			// Check build_info metric
			metric, err := Registry.Gather()
			if err != nil {
				t.Fatalf("Failed to gather metrics: %v", err)
			}

			found := false
			for _, m := range metric {
				if m.GetName() == "build_info" {
					found = true
					if len(m.GetMetric()) < 1 {
						t.Errorf("Expected at least 1 metric, got %d", len(m.GetMetric()))
						break
					}

					// Find the metric with matching version label
					for _, metric := range m.GetMetric() {
						labels := metric.GetLabel()
						labelsMap := make(map[string]string)
						for _, label := range labels {
							labelsMap[label.GetName()] = label.GetValue()
						}

						// Check if this metric has the expected version
						if labelsMap["version"] == tt.wantLabels["version"] {
							// Verify all expected labels
							for k, v := range tt.wantLabels {
								if labelsMap[k] != v {
									t.Errorf("Label %s: got %s, want %s", k, labelsMap[k], v)
								}
							}
							break
						}
					}
					break
				}
			}

			if !found {
				t.Error("build_info metric not found")
			}
		})
	}
}

func TestSetServiceStatus(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		enabled     bool
		wantValue   float64
	}{
		{
			name:        "enable service",
			serviceName: "kubernetes",
			enabled:     true,
			wantValue:   1.0,
		},
		{
			name:        "disable service",
			serviceName: "grafana",
			enabled:     false,
			wantValue:   0.0,
		},
		{
			name:        "enable then disable",
			serviceName: "prometheus",
			enabled:     true,
			wantValue:   1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetServiceStatus(tt.serviceName, tt.enabled)

			// Gather metrics
			metric, err := Registry.Gather()
			if err != nil {
				t.Fatalf("Failed to gather metrics: %v", err)
			}

			found := false
			for _, m := range metric {
				if m.GetName() == "service_status" {
					for _, metricFamily := range m.GetMetric() {
						labels := metricFamily.GetLabel()
						for _, label := range labels {
							if label.GetName() == "service_name" && label.GetValue() == tt.serviceName {
								found = true
								gauge := metricFamily.GetGauge()
								if gauge.GetValue() != tt.wantValue {
									t.Errorf("SetServiceStatus() = %v, want %v", gauge.GetValue(), tt.wantValue)
								}
								break
							}
						}
					}
					break
				}
			}

			if !found {
				t.Errorf("service_status metric for service %s not found", tt.serviceName)
			}
		})
	}
}

func TestSetServiceStatus_Update(t *testing.T) {
	serviceName := "test_service"

	// Enable service
	SetServiceStatus(serviceName, true)

	// Check value is 1.0
	var gaugeValue float64
	metric, err := Registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	for _, m := range metric {
		if m.GetName() == "service_status" {
			for _, metricFamily := range m.GetMetric() {
				labels := metricFamily.GetLabel()
				for _, label := range labels {
					if label.GetName() == "service_name" && label.GetValue() == serviceName {
						gaugeValue = metricFamily.GetGauge().GetValue()
					}
				}
			}
		}
	}

	if gaugeValue != 1.0 {
		t.Errorf("After enabling, expected 1.0, got %v", gaugeValue)
	}

	// Disable service
	SetServiceStatus(serviceName, false)

	// Check value is 0.0
	metric, err = Registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	for _, m := range metric {
		if m.GetName() == "service_status" {
			for _, metricFamily := range m.GetMetric() {
				labels := metricFamily.GetLabel()
				for _, label := range labels {
					if label.GetName() == "service_name" && label.GetValue() == serviceName {
						gaugeValue = metricFamily.GetGauge().GetValue()
					}
				}
			}
		}
	}

	if gaugeValue != 0.0 {
		t.Errorf("After disabling, expected 0.0, got %v", gaugeValue)
	}
}

func TestRegistryNotNil(t *testing.T) {
	if Registry == nil {
		t.Error("Registry should not be nil")
	}
}

func TestBuildInfoMetric(t *testing.T) {
	metric, err := Registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	found := false
	for _, m := range metric {
		if strings.Contains(m.GetName(), "build_info") {
			found = true
			break
		}
	}

	if !found {
		t.Error("build_info metric not found in registry")
	}
}

func TestServerInfoMetric(t *testing.T) {
	metric, err := Registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	found := false
	for _, m := range metric {
		if strings.Contains(m.GetName(), "server_info") {
			found = true
			break
		}
	}

	if !found {
		t.Error("server_info metric not found in registry")
	}
}

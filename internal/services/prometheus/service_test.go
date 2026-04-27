package prometheus

import (
	"testing"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/prometheus/client"
)

func TestPrometheusServiceNew(t *testing.T) {
	svc := NewService()
	if svc == nil {
		t.Error("NewService should return non-nil service")
	}
}

func TestPrometheusServiceIsEnabled(t *testing.T) {
	svc := NewService()
	enabled := svc.IsEnabled()
	_ = enabled
}

func TestPrometheusServiceInitialize(t *testing.T) {
	svc := NewService()
	err := svc.Initialize(nil)
	_ = err
}

func TestPrometheusServiceGetTools(t *testing.T) {
	svc := NewService()
	tools := svc.GetTools()
	if len(tools) > 0 {
		for _, tool := range tools {
			if tool.Name == "" {
				t.Error("Tool name should not be empty")
			}
		}
	}
}

func TestPrometheusServiceGetHandlers(t *testing.T) {
	svc := NewService()
	handlers := svc.GetHandlers()
	if len(handlers) > 0 {
		_ = handlers
	}
}

func TestPrometheusServiceName(t *testing.T) {
	svc := NewService()
	name := svc.Name()
	if name != "prometheus" {
		t.Errorf("Expected service name 'prometheus', got '%s'", name)
	}
}

func TestPrometheusServiceRegistersExtendedToolsAndHandlers(t *testing.T) {
	svc := NewService()
	svc.enabled = true
	svc.client = &client.Client{}

	tools := svc.GetTools()
	handlers := svc.GetHandlers()

	expectedTools := []string{
		"prometheus_targets_summary",
		"prometheus_alerts_summary",
		"prometheus_rules_summary",
		"prometheus_get_metrics_metadata",
		"prometheus_get_target_metadata",
		"prometheus_get_tsdb_stats",
		"prometheus_get_tsdb_status",
		"prometheus_create_snapshot",
		"prometheus_get_wal_replay_status",
		"prometheus_get_server_info",
		"prometheus_get_runtime_info",
	}

	toolNames := make(map[string]bool, len(tools))
	for _, tool := range tools {
		toolNames[tool.Name] = true
	}

	for _, name := range expectedTools {
		if !toolNames[name] {
			t.Fatalf("expected tool %q to be registered", name)
		}
		if _, ok := handlers[name]; !ok {
			t.Fatalf("expected handler %q to be registered", name)
		}
	}
}

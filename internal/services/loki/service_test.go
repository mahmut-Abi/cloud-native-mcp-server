package loki

import (
	"testing"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/loki/client"
)

func TestLokiServiceNew(t *testing.T) {
	svc := NewService()
	if svc == nil {
		t.Fatal("NewService() returned nil")
	}
}

func TestLokiServiceName(t *testing.T) {
	svc := NewService()
	if svc.Name() != "loki" {
		t.Fatalf("expected service name loki, got %s", svc.Name())
	}
}

func TestLokiServiceInitialize(t *testing.T) {
	svc := NewService()
	cfg := &config.AppConfig{}
	cfg.Loki.Enabled = false
	if err := svc.Initialize(cfg); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
}

func TestLokiServiceRegistersToolsAndHandlers(t *testing.T) {
	svc := NewService()
	svc.enabled = true
	svc.client = &client.Client{}

	tools := svc.GetTools()
	handlers := svc.GetHandlers()

	expected := []string{
		"loki_query_logs_summary",
		"loki_query",
		"loki_query_range",
		"loki_get_label_names",
		"loki_get_label_values",
		"loki_get_series",
		"loki_test_connection",
	}

	toolNames := make(map[string]bool, len(tools))
	for _, tool := range tools {
		toolNames[tool.Name] = true
	}

	for _, name := range expected {
		if !toolNames[name] {
			t.Fatalf("expected tool %q to be registered", name)
		}
		if _, ok := handlers[name]; !ok {
			t.Fatalf("expected handler %q to be registered", name)
		}
	}
}

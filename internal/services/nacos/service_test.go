package nacos

import (
	"testing"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	clientpkg "github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/nacos/client"
)

func TestNacosServiceNew(t *testing.T) {
	svc := NewService()
	if svc == nil {
		t.Fatal("NewService() returned nil")
	}
}

func TestNacosServiceName(t *testing.T) {
	svc := NewService()
	if svc.Name() != "nacos" {
		t.Fatalf("expected service name nacos, got %q", svc.Name())
	}
}

func TestNacosServiceDisabledByDefault(t *testing.T) {
	svc := NewService()
	if svc.IsEnabled() {
		t.Fatal("service should be disabled by default")
	}
	if tools := svc.GetTools(); len(tools) != 0 {
		t.Fatalf("expected no tools when disabled, got %d", len(tools))
	}
	if handlers := svc.GetHandlers(); len(handlers) != 0 {
		t.Fatalf("expected no handlers when disabled, got %d", len(handlers))
	}
}

func TestNacosServiceInitializeNilConfig(t *testing.T) {
	svc := NewService()
	if err := svc.Initialize(nil); err != nil {
		t.Fatalf("Initialize(nil) returned error: %v", err)
	}
	if svc.IsEnabled() {
		t.Fatal("service should remain disabled without config")
	}
}

func TestNacosServiceRegistersToolsAndHandlers(t *testing.T) {
	svc := NewService()
	svc.enabled = true
	svc.client = &clientpkg.Client{}
	svc.defaultGroup = "DEFAULT_GROUP"
	svc.defaultNamespaceID = "public"

	cfg := &config.AppConfig{}
	cfg.Nacos.Group = "DEFAULT_GROUP"
	cfg.Nacos.NamespaceID = "public"

	tools := svc.GetTools()
	handlers := svc.GetHandlers()

	expected := []string{
		"nacos_test_connection",
		"nacos_list_namespaces",
		"nacos_list_configs_summary",
		"nacos_get_config",
		"nacos_list_services_summary",
		"nacos_get_service",
		"nacos_list_instances",
		"nacos_list_cluster_nodes",
		"nacos_get_system_metrics",
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

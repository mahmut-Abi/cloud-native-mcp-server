package argocd

import (
	"testing"
)

func TestArgoCDServiceNew(t *testing.T) {
	svc := NewService()
	if svc == nil {
		t.Fatal("NewService() returned nil")
	}
}

func TestArgoCDServiceName(t *testing.T) {
	svc := NewService()
	if svc.Name() != "argocd" {
		t.Fatalf("expected service name argocd, got %q", svc.Name())
	}
}

func TestArgoCDServiceDisabledByDefault(t *testing.T) {
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

func TestArgoCDServiceRegistersToolsAndHandlers(t *testing.T) {
	svc := NewService()
	svc.enabled = true

	expected := []string{
		"argocd_test_connection",
		"argocd_list_applications_summary",
		"argocd_get_application",
		"argocd_get_application_manifests",
		"argocd_list_projects",
		"argocd_get_project",
		"argocd_list_clusters",
	}

	toolNames := make(map[string]bool)
	for _, tool := range svc.GetTools() {
		toolNames[tool.Name] = true
	}
	for _, name := range expected {
		if !toolNames[name] {
			t.Fatalf("expected tool %q to be registered", name)
		}
		if _, ok := svc.GetHandlers()[name]; !ok {
			t.Fatalf("expected handler %q to be registered", name)
		}
	}
}

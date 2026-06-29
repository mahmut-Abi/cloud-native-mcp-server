package langfuse

import (
	"testing"
)

func TestLangfuseServiceNew(t *testing.T) {
	svc := NewService()
	if svc == nil {
		t.Fatal("NewService() returned nil")
	}
}

func TestLangfuseServiceName(t *testing.T) {
	svc := NewService()
	if svc.Name() != "langfuse" {
		t.Fatalf("expected service name langfuse, got %q", svc.Name())
	}
}

func TestLangfuseServiceDisabledByDefault(t *testing.T) {
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

func TestLangfuseServiceInitializeNilConfig(t *testing.T) {
	svc := NewService()
	if err := svc.Initialize(nil); err != nil {
		t.Fatalf("Initialize(nil) returned error: %v", err)
	}
	if svc.IsEnabled() {
		t.Fatal("service should remain disabled without config")
	}
}

func TestLangfuseServiceExposesProjectManagementTools(t *testing.T) {
	svc := NewService()
	svc.enabled = true

	tools := svc.GetTools()
	if len(tools) != 37 {
		t.Fatalf("expected 37 tools, got %d", len(tools))
	}

	handlers := svc.GetHandlers()
	for _, name := range []string{
		"langfuse_get_project",
		"langfuse_list_organization_projects",
		"langfuse_create_project",
		"langfuse_update_project",
		"langfuse_delete_project",
		"langfuse_list_project_memberships",
		"langfuse_upsert_project_membership",
		"langfuse_delete_project_membership",
		"langfuse_list_organization_api_keys",
		"langfuse_list_project_api_keys",
		"langfuse_create_project_api_key",
		"langfuse_delete_project_api_key",
	} {
		if _, ok := handlers[name]; !ok {
			t.Fatalf("expected handler %q to be registered", name)
		}
	}
}

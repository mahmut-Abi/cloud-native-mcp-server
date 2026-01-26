package helm

import (
	"testing"
)

func TestHelmServiceNew(t *testing.T) {
	svc := NewService()
	if svc == nil {
		t.Error("NewService should return non-nil service")
	}
}

func TestHelmServiceIsEnabled(t *testing.T) {
	svc := NewService()
	enabled := svc.IsEnabled()
	_ = enabled
}

func TestHelmServiceInitialize(t *testing.T) {
	svc := NewService()
	err := svc.Initialize(nil)
	_ = err
}

func TestHelmServiceGetTools(t *testing.T) {
	svc := NewService()
	err := svc.Initialize(nil)
	if err != nil {
		t.Fatalf("Failed to initialize service: %v", err)
	}
	tools := svc.GetTools()
	if len(tools) > 0 {
		for _, tool := range tools {
			if tool.Name == "" {
				t.Error("Tool name should not be empty")
			}
		}
	}
}

func TestHelmServiceGetHandlers(t *testing.T) {
	svc := NewService()
	err := svc.Initialize(nil)
	if err != nil {
		t.Fatalf("Failed to initialize service: %v", err)
	}
	handlers := svc.GetHandlers()
	if len(handlers) > 0 {
		_ = handlers
	}
}

func TestHelmServiceGetMirrorConfigurationTool(t *testing.T) {
	svc := NewService()
	err := svc.Initialize(nil)
	if err != nil {
		t.Fatalf("Failed to initialize service: %v", err)
	}
	tools := svc.GetTools()

	found := false
	for _, tool := range tools {
		if tool.Name == "helm_get_mirror_configuration" {
			found = true
			break
		}
	}

	if !found {
		t.Error("helm_get_mirror_configuration tool not found in service tools")
	}
}

func TestHelmServiceGetMirrorConfigurationHandler(t *testing.T) {
	svc := NewService()
	err := svc.Initialize(nil)
	if err != nil {
		t.Fatalf("Failed to initialize service: %v", err)
	}
	handlers := svc.GetHandlers()

	if _, ok := handlers["helm_get_mirror_configuration"]; !ok {
		t.Error("helm_get_mirror_configuration handler not found in service handlers")
	}
}

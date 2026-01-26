package grafana

import (
	"testing"
)

func TestGrafanaServiceNew(t *testing.T) {
	svc := NewService()
	if svc == nil {
		t.Error("NewService should return non-nil service")
	}
}

func TestGrafanaServiceIsEnabled(t *testing.T) {
	svc := NewService()
	enabled := svc.IsEnabled()
	_ = enabled
}

func TestGrafanaServiceInitialize(t *testing.T) {
	svc := NewService()
	err := svc.Initialize(nil)
	_ = err
}

func TestGrafanaServiceGetTools(t *testing.T) {
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

func TestGrafanaServiceGetHandlers(t *testing.T) {
	svc := NewService()
	handlers := svc.GetHandlers()
	if len(handlers) > 0 {
		_ = handlers
	}
}

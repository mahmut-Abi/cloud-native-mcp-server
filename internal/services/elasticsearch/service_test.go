package elasticsearch

import (
	"testing"
)

func TestElasticsearchServiceNew(t *testing.T) {
	svc := NewService()
	if svc == nil {
		t.Error("NewService should return non-nil service")
	}
}

func TestElasticsearchServiceIsEnabled(t *testing.T) {
	svc := NewService()
	enabled := svc.IsEnabled()
	_ = enabled
}

func TestElasticsearchServiceInitialize(t *testing.T) {
	svc := NewService()
	err := svc.Initialize(nil)
	_ = err
}

func TestElasticsearchServiceGetTools(t *testing.T) {
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

func TestElasticsearchServiceGetHandlers(t *testing.T) {
	svc := NewService()
	handlers := svc.GetHandlers()
	if len(handlers) > 0 {
		_ = handlers
	}
}

func TestElasticsearchServiceName(t *testing.T) {
	svc := NewService()
	name := svc.Name()
	if name != "elasticsearch" {
		t.Errorf("Expected service name 'elasticsearch', got '%s'", name)
	}
}
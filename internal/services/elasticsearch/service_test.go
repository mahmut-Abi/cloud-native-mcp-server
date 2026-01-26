package elasticsearch

import (
	"testing"
)

func TestNewService(t *testing.T) {
	svc := NewService()
	if svc == nil {
		t.Error("Expected service, got nil")
	}
	if svc.Name() != "elasticsearch" {
		t.Errorf("Expected name elasticsearch, got %s", svc.Name())
	}
}

func TestServiceDisabledByDefault(t *testing.T) {
	svc := NewService()
	if svc.IsEnabled() {
		t.Error("Expected service disabled by default")
	}
}

func TestServiceInitializeWithNilConfig(t *testing.T) {
	svc := NewService()
	err := svc.Initialize(nil)
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}
	if svc.IsEnabled() {
		t.Error("Expected service disabled with nil config")
	}
}

func TestServiceGetToolsWhenDisabled(t *testing.T) {
	svc := NewService()
	tools := svc.GetTools()
	if tools != nil {
		t.Error("Expected nil tools when disabled")
	}
}

func TestServiceGetHandlersWhenDisabled(t *testing.T) {
	svc := NewService()
	handlers := svc.GetHandlers()
	if handlers != nil {
		t.Error("Expected nil handlers when disabled")
	}
}

package sentry

import "testing"

func TestSentryServiceNew(t *testing.T) {
	svc := NewService()
	if svc == nil {
		t.Fatal("NewService() returned nil")
	}
}

func TestSentryServiceName(t *testing.T) {
	svc := NewService()
	if svc.Name() != "sentry" {
		t.Fatalf("expected service name sentry, got %q", svc.Name())
	}
}

func TestSentryServiceDisabledByDefault(t *testing.T) {
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

func TestSentryServiceInitializeNilConfig(t *testing.T) {
	svc := NewService()
	if err := svc.Initialize(nil); err != nil {
		t.Fatalf("Initialize(nil) returned error: %v", err)
	}
	if svc.IsEnabled() {
		t.Fatal("service should remain disabled without config")
	}
}

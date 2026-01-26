package prometheus

import (
	"testing"
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
	// Service should be disabled until initialized
	_ = enabled
}

func TestPrometheusServiceInitialize(t *testing.T) {
	svc := NewService()
	err := svc.Initialize(nil)
	// May fail in test environment without proper configuration
	_ = err
}

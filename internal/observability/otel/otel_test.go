package otel

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/trace"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}

	if cfg.Enabled {
		t.Error("DefaultConfig should have disabled by default")
	}

	if cfg.ServiceName != "cloud-native-mcp-server" {
		t.Errorf("Expected service name 'cloud-native-mcp-server', got '%s'", cfg.ServiceName)
	}

	if cfg.TracingConfig.SampleRate <= 0 || cfg.TracingConfig.SampleRate > 1.0 {
		t.Errorf("Invalid sample rate: %f", cfg.TracingConfig.SampleRate)
	}

	if cfg.MetricsConfig.ExportInterval == 0 {
		t.Error("Export interval should not be zero")
	}
}

func TestInit(t *testing.T) {
	cfg := &Config{
		Enabled:        false,
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		Environment:    "test",
		Endpoint:       "http://localhost:4317",
		Insecure:       true,
		TracingConfig: TracingConfig{
			Enabled:            true,
			SampleRate:         1.0,
			ExportTimeout:      5 * time.Second,
			BatchTimeout:       1 * time.Second,
			MaxExportBatchSize: 10,
		},
		MetricsConfig: MetricsConfig{
			Enabled:        true,
			ExportInterval: 5 * time.Second,
			ExportTimeout:  5 * time.Second,
			Temporality:    "cumulative",
		},
	}

	err := Init(cfg)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Shutdown to clean up
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = Shutdown(ctx)
	if err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}
}

func TestInitDisabled(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Enabled = false

	err := Init(cfg)
	if err != nil {
		t.Fatalf("Init with disabled config failed: %v", err)
	}

	if GetTracer() != nil {
		t.Error("Tracer should be nil when OTEL is disabled")
	}

	if GetMeter() != nil {
		t.Error("Meter should be nil when OTEL is disabled")
	}
}

func TestSpanHelper(t *testing.T) {
	helper := NewSpanHelper()

	if helper == nil {
		t.Fatal("NewSpanHelper returned nil")
	}

	ctx := context.Background()

	// Test with disabled OTEL
	_, _ = helper.StartSpan(ctx, "test-span")
}

func TestWithSpan(t *testing.T) {
	ctx := context.Background()

	err := WithSpan(ctx, "test-span", func(ctx context.Context, span trace.Span) error {
		// This should work even when OTEL is disabled
		return nil
	})

	if err != nil {
		t.Errorf("WithSpan failed: %v", err)
	}
}

func TestWithSpanWithError(t *testing.T) {
	ctx := context.Background()
	expectedErr := "test error"

	err := WithSpan(ctx, "test-span", func(ctx context.Context, span trace.Span) error {
		return &TestError{msg: expectedErr}
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

// TestError is a simple error for testing
type TestError struct {
	msg string
}

func (e *TestError) Error() string {
	return e.msg
}
package otel

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"go.opentelemetry.io/otel/attribute"
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

func TestFromAppConfig(t *testing.T) {
	appConfig := &config.AppConfig{}
	appConfig.OTEL.Enabled = true
	appConfig.OTEL.ServiceName = "test-service"
	appConfig.OTEL.ServiceVersion = "2.0.0"
	appConfig.OTEL.Environment = "test"
	appConfig.OTEL.Endpoint = "http://localhost:4318"
	appConfig.OTEL.Insecure = true
	appConfig.OTEL.Tracing.Enabled = true
	appConfig.OTEL.Tracing.SampleRate = 0.5
	appConfig.OTEL.Tracing.ExportTimeoutSec = 10
	appConfig.OTEL.Tracing.BatchTimeoutSec = 2
	appConfig.OTEL.Tracing.MaxExportBatchSize = 100
	appConfig.OTEL.Metrics.Enabled = true
	appConfig.OTEL.Metrics.ExportIntervalSec = 10
	appConfig.OTEL.Metrics.ExportTimeoutSec = 10
	appConfig.OTEL.Metrics.Temporality = "delta"

	cfg := FromAppConfig(appConfig, "1.0.0")
	if cfg == nil {
		t.Fatal("FromAppConfig returned nil")
	}

	if cfg.ServiceName != "test-service" {
		t.Errorf("Expected service name 'test-service', got '%s'", cfg.ServiceName)
	}

	if cfg.ServiceVersion != "2.0.0" {
		t.Errorf("Expected service version '2.0.0', got '%s'", cfg.ServiceVersion)
	}

	if cfg.Endpoint != "http://localhost:4318" {
		t.Errorf("Expected endpoint 'http://localhost:4318', got '%s'", cfg.Endpoint)
	}

	if cfg.TracingConfig.SampleRate != 0.5 {
		t.Errorf("Expected sample rate 0.5, got %f", cfg.TracingConfig.SampleRate)
	}

	if cfg.MetricsConfig.Temporality != "delta" {
		t.Errorf("Expected temporality 'delta', got '%s'", cfg.MetricsConfig.Temporality)
	}
}

func TestFromAppConfigNil(t *testing.T) {
	cfg := FromAppConfig(nil, "1.0.0")
	if cfg != nil {
		t.Error("FromAppConfig should return nil for nil input")
	}
}

func TestFromAppConfigDisabled(t *testing.T) {
	appConfig := &config.AppConfig{}
	appConfig.OTEL.Enabled = false

	cfg := FromAppConfig(appConfig, "1.0.0")
	if cfg != nil {
		t.Error("FromAppConfig should return nil when disabled")
	}
}

func TestMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mw := Middleware("test-service")
	if mw == nil {
		t.Fatal("Middleware returned nil")
	}

	wrappedHandler := mw(handler)
	if wrappedHandler == nil {
		t.Fatal("Wrapped handler is nil")
	}

	// Test that the middleware works
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestMiddlewareWithDisabledOTEL(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := Middleware("test-service")
	wrappedHandler := mw(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestInitWithTracing(t *testing.T) {
	cfg := &Config{
		Enabled:        true,
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
			Enabled:        false, // Only test tracing
			ExportInterval: 5 * time.Second,
			ExportTimeout:  5 * time.Second,
		},
	}

	err := Init(cfg)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if GetTracer() == nil {
		t.Error("Tracer should not be nil when tracing is enabled")
	}

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = Shutdown(ctx)
}

func TestInitWithMetrics(t *testing.T) {
	cfg := &Config{
		Enabled:        true,
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		Environment:    "test",
		Endpoint:       "http://localhost:4317",
		Insecure:       true,
		TracingConfig: TracingConfig{
			Enabled: false, // Only test metrics
		},
		MetricsConfig: MetricsConfig{
			Enabled:        true,
			ExportInterval: 5 * time.Second,
			ExportTimeout:  5 * time.Second,
		},
	}

	err := Init(cfg)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if GetMeter() == nil {
		t.Error("Meter should not be nil when metrics is enabled")
	}

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = Shutdown(ctx)
}

func TestSpanHelperRecordError(t *testing.T) {
	helper := NewSpanHelper()

	if helper == nil {
		t.Fatal("NewSpanHelper returned nil")
	}

	// Test with disabled OTEL
	testErr := &TestError{msg: "test error"}
	helper.RecordError(nil, testErr, attribute.String("key", "value"))
	// Should not panic
}

func TestSpanHelperSetAttributes(t *testing.T) {
	helper := NewSpanHelper()

	if helper == nil {
		t.Fatal("NewSpanHelper returned nil")
	}

	// Test with disabled OTEL
	attrs := []attribute.KeyValue{
		attribute.String("service", "test"),
		attribute.Int("count", 42),
	}
	helper.SetAttributes(nil, attrs...)
	// Should not panic
}

func TestSpanHelperAddEvent(t *testing.T) {
	helper := NewSpanHelper()

	if helper == nil {
		t.Fatal("NewSpanHelper returned nil")
	}

	// Test with disabled OTEL
	attrs := []attribute.KeyValue{
		attribute.String("event.type", "test"),
	}
	helper.AddEvent(nil, "test-event", attrs...)
	// Should not panic
}

func TestWithSpanAsync(t *testing.T) {
	ctx := context.Background()

	// Test with disabled OTEL
	WithSpanAsync(ctx, "test-async-span", func(ctx context.Context, span trace.Span) {
		// This should work even when OTEL is disabled
	})

	// Give it a moment to complete
	time.Sleep(10 * time.Millisecond)
}

func TestShutdownWithNilProviders(t *testing.T) {
	ctx := context.Background()

	// Test shutdown when providers are nil
	err := Shutdown(ctx)
	// Shutdown may have errors but should not panic
	if err != nil {
		t.Logf("Shutdown with nil providers returned error (expected): %v", err)
	}
}

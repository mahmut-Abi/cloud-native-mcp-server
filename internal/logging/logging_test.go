package logging

import (
	"bytes"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if config == nil {
		t.Error("Expected config, got nil")
		return
	}
	if config.Output != os.Stdout {
		t.Error("Expected Stdout as default output")
	}
	if config.Level != logrus.InfoLevel {
		t.Error("Expected InfoLevel as default level")
	}
	if config.UseJSONFormat {
		t.Error("Expected JSON format disabled by default")
	}
}

func TestInitStdoutLogger(t *testing.T) {
	// Should not panic
	InitStdoutLogger()
}

func TestInitLoggerWithConfig(t *testing.T) {
	buf := &bytes.Buffer{}
	config := &LogConfig{
		Output:          buf,
		Level:           logrus.InfoLevel,
		UseJSONFormat:   false,
		TimestampFormat: time.RFC3339,
		EnableColors:    false,
	}
	InitLoggerWithConfig(config)
	logrus.SetOutput(buf)
}

func TestSetLogLevel(t *testing.T) {
	SetLogLevel(logrus.DebugLevel)
	if logrus.GetLevel() != logrus.DebugLevel {
		t.Error("Expected DebugLevel")
	}
	SetLogLevel(logrus.InfoLevel)
}

func TestEnableJSONFormat(t *testing.T) {
	EnableJSONFormat()
}

func TestInitLoggerWithJSONFormat(t *testing.T) {
	buf := &bytes.Buffer{}
	config := &LogConfig{
		Output:          buf,
		Level:           logrus.InfoLevel,
		UseJSONFormat:   true,
		TimestampFormat: time.RFC3339,
		EnableColors:    false,
	}
	InitLoggerWithConfig(config)
}

func TestLogResourceError(t *testing.T) {
	ctx := ResourceContext{
		Component: "kubernetes",
		Kind:      "Pod",
		Name:      "test-pod",
		Namespace: "default",
		Action:    "get",
	}

	err := errors.New("connection timeout")
	LogResourceError(ctx, err) // Should not panic
}

func TestWrapResourceError(t *testing.T) {
	tests := []struct {
		name   string
		ctx    ResourceContext
		input  error
		expect string
	}{
		{
			name: "with namespace and name",
			ctx: ResourceContext{
				Kind:      "Pod",
				Name:      "my-pod",
				Namespace: "default",
				Action:    "get",
			},
			input:  errors.New("not found"),
			expect: "failed to get Pod default/my-pod",
		},
		{
			name: "without namespace",
			ctx: ResourceContext{
				Kind:   "Node",
				Name:   "node-1",
				Action: "describe",
			},
			input:  errors.New("unavailable"),
			expect: "failed to describe Node node-1",
		},
		{
			name: "nil error",
			ctx: ResourceContext{
				Kind:   "Service",
				Action: "list",
			},
			input:  nil,
			expect: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WrapResourceError(tt.ctx, tt.input)
			if tt.expect == "" && err != nil {
				t.Errorf("expected nil, got %v", err)
			}
			if tt.expect != "" && err == nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}

func TestLogAndWrapResourceError(t *testing.T) {
	ctx := ResourceContext{
		Component: "grafana",
		Kind:      "Dashboard",
		Name:      "test-dash",
		Action:    "create",
	}

	err := LogAndWrapResourceError(ctx, errors.New("auth failed"))
	if err == nil {
		t.Error("expected error, got nil")
	}
}

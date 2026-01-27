package opentelemetry

import (
	"testing"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
)

func TestOpenTelemetryServiceNew(t *testing.T) {
	svc := NewService()
	if svc == nil {
		t.Error("NewService should return non-nil service")
	}
}

func TestOpenTelemetryServiceIsEnabled(t *testing.T) {
	svc := NewService()
	enabled := svc.IsEnabled()
	_ = enabled
}

func TestOpenTelemetryServiceInitialize(t *testing.T) {
	svc := NewService()
	
	appConfig := &config.AppConfig{}
	appConfig.OpenTelemetry.Enabled = false
	
	err := svc.Initialize(appConfig)
	_ = err
}

func TestOpenTelemetryServiceGetTools(t *testing.T) {
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

func TestOpenTelemetryServiceGetHandlers(t *testing.T) {
	svc := NewService()
	handlers := svc.GetHandlers()
	if len(handlers) > 0 {
		_ = handlers
	}
}

func TestOpenTelemetryServiceGetResources(t *testing.T) {
	svc := NewService()
	resources := svc.GetResources()
	if len(resources) > 0 {
		for _, resource := range resources {
			if resource.URI == "" {
				t.Error("Resource URI should not be empty")
			}
		}
	}
}

func TestOpenTelemetryServiceGetResourceHandlers(t *testing.T) {
	svc := NewService()
	handlers := svc.GetResourceHandlers()
	if len(handlers) > 0 {
		_ = handlers
	}
}

func TestOpenTelemetryServiceGetClient(t *testing.T) {
	svc := NewService()
	client := svc.GetClient()
	_ = client
}

func TestOpenTelemetryServiceName(t *testing.T) {
	svc := NewService()
	name := svc.Name()
	if name != "opentelemetry" {
		t.Errorf("Expected service name 'opentelemetry', got '%s'", name)
	}
}
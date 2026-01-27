package jaeger

import (
	"testing"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
)

func TestJaegerServiceNew(t *testing.T) {
	svc := NewService()
	if svc == nil {
		t.Error("NewService should return non-nil service")
	}
}

func TestJaegerServiceIsEnabled(t *testing.T) {
	svc := NewService()
	enabled := svc.IsEnabled()
	_ = enabled
}

func TestJaegerServiceInitialize(t *testing.T) {
	svc := NewService()
	
	appConfig := &config.AppConfig{}
	appConfig.Jaeger.Enabled = false
	
	err := svc.Initialize(appConfig)
	_ = err
}

func TestJaegerServiceGetTools(t *testing.T) {
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

func TestJaegerServiceGetHandlers(t *testing.T) {
	svc := NewService()
	handlers := svc.GetHandlers()
	if len(handlers) > 0 {
		_ = handlers
	}
}

func TestJaegerServiceName(t *testing.T) {
	svc := NewService()
	name := svc.Name()
	if name != "jaeger" {
		t.Errorf("Expected service name 'jaeger', got '%s'", name)
	}
}

func TestJaegerServiceGetToolsCache(t *testing.T) {
	svc := NewService()
	cache := svc.GetToolsCache()
	if cache == nil {
		t.Error("GetToolsCache should return non-nil cache")
	}
}

func TestJaegerServiceGetClient(t *testing.T) {
	svc := NewService()
	client := svc.GetClient()
	_ = client
}
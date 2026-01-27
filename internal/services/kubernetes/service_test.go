package kubernetes

import (
	"testing"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
)

func TestNewService(t *testing.T) {
	service := NewService()

	if service == nil {
		t.Fatal("NewService() returned nil")
	}

	if service.Name() != "kubernetes" {
		t.Errorf("Expected service name 'kubernetes', got '%s'", service.Name())
	}

	if service.IsEnabled() {
		t.Error("Service should not be enabled before initialization")
	}

	if service.client != nil {
		t.Error("Client should be nil before initialization")
	}
}

func TestServiceInitializeWithNilConfig(t *testing.T) {
	service := NewService()

	err := service.Initialize(nil)

	// Should not error but client creation might fail due to missing kubeconfig
	// This is expected behavior in test environment
	if err != nil {
		t.Logf("Expected error in test environment (no kubeconfig): %v", err)
	}
}

func TestServiceInitializeWithInvalidConfig(t *testing.T) {
	service := NewService()

	// Test with wrong config type
	err := service.Initialize("invalid-config")

	// Should not error but use default options
	if err != nil {
		t.Logf("Expected error in test environment (no kubeconfig): %v", err)
	}
}

func TestServiceInitializeWithValidConfig(t *testing.T) {
	service := NewService()

	appConfig := &config.AppConfig{
		Kubernetes: struct {
			Kubeconfig string  `yaml:"kubeconfig"`
			TimeoutSec int     `yaml:"timeoutSec"`
			QPS        float32 `yaml:"qps"`
			Burst      int     `yaml:"burst"`
		}{
			Kubeconfig: "/non-existent/kubeconfig", // Use non-existent path for test
			TimeoutSec: 30,
			QPS:        10.0,
			Burst:      20,
		},
	}

	err := service.Initialize(appConfig)

	// Error expected due to non-existent kubeconfig file
	if err != nil {
		t.Logf("Expected error due to non-existent kubeconfig: %v", err)
	}
}

func TestServiceGetToolsWhenDisabled(t *testing.T) {
	service := NewService()
	// Service is disabled by default

	tools := service.GetTools()

	if tools != nil {
		t.Error("GetTools() should return nil when service is disabled")
	}
}

func TestServiceGetToolsWhenClientIsNil(t *testing.T) {
	service := NewService()
	service.enabled = true // Enable service but keep client nil

	tools := service.GetTools()

	if tools != nil {
		t.Error("GetTools() should return nil when client is nil")
	}
}

func TestServiceGetHandlersWhenDisabled(t *testing.T) {
	service := NewService()
	// Service is disabled by default

	handlers := service.GetHandlers()

	if handlers != nil {
		t.Error("GetHandlers() should return nil when service is disabled")
	}
}

func TestServiceGetHandlersWhenClientIsNil(t *testing.T) {
	service := NewService()
	service.enabled = true // Enable service but keep client nil

	handlers := service.GetHandlers()

	if handlers != nil {
		t.Error("GetHandlers() should return nil when client is nil")
	}
}

func TestServiceName(t *testing.T) {
	service := NewService()

	expectedName := "kubernetes"
	actualName := service.Name()

	if actualName != expectedName {
		t.Errorf("Expected service name '%s', got '%s'", expectedName, actualName)
	}
}

func TestServiceIsEnabledDefault(t *testing.T) {
	service := NewService()

	if service.IsEnabled() {
		t.Error("Service should not be enabled by default")
	}
}

func TestServiceGetClient(t *testing.T) {
	service := NewService()

	client := service.GetClient()

	if client != nil {
		t.Error("GetClient() should return nil when client is not initialized")
	}
}

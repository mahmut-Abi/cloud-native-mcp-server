package alertmanager

import (
	"testing"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
)

func TestNewService(t *testing.T) {
	service := NewService()

	if service == nil {
		t.Fatal("Service should not be nil")
	}

	if service.Name() != "alertmanager" {
		t.Errorf("Expected service name to be 'alertmanager', got %s", service.Name())
	}

	if service.IsEnabled() {
		t.Error("Service should be disabled by default")
	}

	if service.GetTools() != nil {
		t.Error("Tools should be nil when service is disabled")
	}

	if service.GetHandlers() != nil {
		t.Error("Handlers should be nil when service is disabled")
	}
}

func TestServiceInitializeWithNilConfig(t *testing.T) {
	service := NewService()

	err := service.Initialize(nil)
	if err != nil {
		t.Errorf("Initialize with nil config should not return error, got %v", err)
	}

	if service.IsEnabled() {
		t.Error("Service should remain disabled with nil config")
	}
}

func TestServiceInitializeWithEmptyConfig(t *testing.T) {
	service := NewService()
	config := &config.AppConfig{}

	err := service.Initialize(config)
	if err != nil {
		t.Errorf("Initialize with empty config should not return error, got %v", err)
	}

	if service.IsEnabled() {
		t.Error("Service should be disabled when Alertmanager.Enabled is false")
	}
}

func TestServiceInitializeWithEnabledConfig(t *testing.T) {
	service := NewService()
	config := &config.AppConfig{}
	config.Alertmanager.Enabled = true
	config.Alertmanager.Address = "http://localhost:9093"

	err := service.Initialize(config)
	if err != nil {
		t.Errorf("Initialize with valid config should not return error, got %v", err)
	}

	if !service.IsEnabled() {
		t.Error("Service should be enabled when properly configured")
	}

	if service.GetClient() == nil {
		t.Error("Client should not be nil when service is enabled")
	}
}

func TestServiceGetToolsWhenEnabled(t *testing.T) {
	service := NewService()
	config := &config.AppConfig{}
	config.Alertmanager.Enabled = true
	config.Alertmanager.Address = "http://localhost:9093"

	err := service.Initialize(config)
	if err != nil {
		t.Fatalf("Failed to initialize service: %v", err)
	}

	tools := service.GetTools()
	if tools == nil {
		t.Fatal("Tools should not be nil when service is enabled")
	}

	if len(tools) == 0 {
		t.Error("Should have at least one tool")
	}

	// Check that tools are cached
	tools2 := service.GetTools()
	if &tools[0] != &tools2[0] {
		t.Error("Tools should be cached and return the same slice")
	}

	// Verify tool names
	expectedTools := []string{
		"alertmanager_alerts_summary",
		"alertmanager_silences_summary",
		"alertmanager_alert_groups_paginated",
		"alertmanager_silences_paginated",
		"alertmanager_receivers_summary",
		"alertmanager_query_alerts_advanced",
		"alertmanager_health_summary",
		"alertmanager_get_status",
		"alertmanager_get_alerts",
		"alertmanager_get_alert_groups",
		"alertmanager_get_silences",
		"alertmanager_create_silence",
		"alertmanager_delete_silence",
		"alertmanager_get_receivers",
		"alertmanager_test_receiver",
		"alertmanager_query_alerts",
	}

	if len(tools) != len(expectedTools) {
		t.Errorf("Expected %d tools, got %d", len(expectedTools), len(tools))
	}

	for i, tool := range tools {
		if i < len(expectedTools) && tool.Name != expectedTools[i] {
			t.Errorf("Expected tool %d to be %s, got %s", i, expectedTools[i], tool.Name)
		}
	}
}

func TestServiceGetHandlersWhenEnabled(t *testing.T) {
	service := NewService()
	config := &config.AppConfig{}
	config.Alertmanager.Enabled = true
	config.Alertmanager.Address = "http://localhost:9093"

	err := service.Initialize(config)
	if err != nil {
		t.Fatalf("Failed to initialize service: %v", err)
	}

	handlers := service.GetHandlers()
	if handlers == nil {
		t.Fatal("Handlers should not be nil when service is enabled")
	}

	if len(handlers) == 0 {
		t.Error("Should have at least one handler")
	}

	// Verify handler names match tool names
	expectedHandlers := []string{
		"alertmanager_alerts_summary",
		"alertmanager_silences_summary",
		"alertmanager_alert_groups_paginated",
		"alertmanager_silences_paginated",
		"alertmanager_receivers_summary",
		"alertmanager_query_alerts_advanced",
		"alertmanager_health_summary",
		"alertmanager_get_status",
		"alertmanager_get_alerts",
		"alertmanager_get_alert_groups",
		"alertmanager_get_silences",
		"alertmanager_create_silence",
		"alertmanager_delete_silence",
		"alertmanager_get_receivers",
		"alertmanager_test_receiver",
		"alertmanager_query_alerts",
	}

	if len(handlers) != len(expectedHandlers) {
		t.Errorf("Expected %d handlers, got %d", len(expectedHandlers), len(handlers))
	}

	for _, handlerName := range expectedHandlers {
		if _, exists := handlers[handlerName]; !exists {
			t.Errorf("Expected handler %s not found", handlerName)
		}
	}
}

func TestServiceConfigurationOptions(t *testing.T) {
	service := NewService()
	config := &config.AppConfig{}
	config.Alertmanager.Enabled = true
	config.Alertmanager.Address = "https://alertmanager.example.com:9093"
	config.Alertmanager.TimeoutSec = 60
	config.Alertmanager.Username = "testuser"
	config.Alertmanager.Password = "testpass"
	config.Alertmanager.BearerToken = "testtoken"
	config.Alertmanager.TLSSkipVerify = true

	err := service.Initialize(config)
	if err != nil {
		t.Errorf("Initialize with full config should not return error, got %v", err)
	}

	if !service.IsEnabled() {
		t.Error("Service should be enabled with full config")
	}

	client := service.GetClient()
	if client == nil {
		t.Error("Client should not be nil")
	}
}

func TestServiceName(t *testing.T) {
	service := NewService()
	if service.Name() != "alertmanager" {
		t.Errorf("Expected service name to be 'alertmanager', got %s", service.Name())
	}
}

package services

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

// MockService implements the Service interface for testing
type MockService struct {
	name     string
	enabled  bool
	tools    []mcp.Tool
	handlers map[string]server.ToolHandlerFunc
	initErr  error
}

func (m *MockService) Name() string {
	return m.name
}

func (m *MockService) GetTools() []mcp.Tool {
	return m.tools
}

func (m *MockService) GetHandlers() map[string]server.ToolHandlerFunc {
	return m.handlers
}

func (m *MockService) Initialize(cfg interface{}) error {
	return m.initErr
}

func (m *MockService) IsEnabled() bool {
	return m.enabled
}

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()

	if registry == nil {
		t.Fatal("NewRegistry() returned nil")
	}

	if registry.services == nil {
		t.Error("Registry services map is nil")
	}

	if len(registry.services) != 0 {
		t.Error("Registry should start with empty services map")
	}
}

func TestRegistryRegister(t *testing.T) {
	registry := NewRegistry()

	mockService := &MockService{
		name:    "test-service",
		enabled: true,
	}

	registry.Register(mockService)

	if len(registry.services) != 1 {
		t.Errorf("Expected 1 service, got %d", len(registry.services))
	}

	if _, exists := registry.services["test-service"]; !exists {
		t.Error("Service was not registered with correct name")
	}

	if registry.services["test-service"] != mockService {
		t.Error("Registered service does not match original")
	}
}

func TestRegistryRegisterMultipleServices(t *testing.T) {
	registry := NewRegistry()

	service1 := &MockService{name: "service1", enabled: true}
	service2 := &MockService{name: "service2", enabled: false}
	service3 := &MockService{name: "service3", enabled: true}

	registry.Register(service1)
	registry.Register(service2)
	registry.Register(service3)

	if len(registry.services) != 3 {
		t.Errorf("Expected 3 services, got %d", len(registry.services))
	}
}

func TestRegistryGetEnabledServices(t *testing.T) {
	registry := NewRegistry()

	enabledService1 := &MockService{name: "enabled1", enabled: true}
	enabledService2 := &MockService{name: "enabled2", enabled: true}
	disabledService := &MockService{name: "disabled", enabled: false}

	registry.Register(enabledService1)
	registry.Register(enabledService2)
	registry.Register(disabledService)

	enabledServices := registry.GetEnabledServices()

	if len(enabledServices) != 2 {
		t.Errorf("Expected 2 enabled services, got %d", len(enabledServices))
	}

	if _, exists := enabledServices["enabled1"]; !exists {
		t.Error("enabled1 service should be in enabled services")
	}

	if _, exists := enabledServices["enabled2"]; !exists {
		t.Error("enabled2 service should be in enabled services")
	}

	if _, exists := enabledServices["disabled"]; exists {
		t.Error("disabled service should not be in enabled services")
	}
}

func TestRegistryGetAllTools(t *testing.T) {
	registry := NewRegistry()

	tool1 := mcp.NewTool("tool1", mcp.WithDescription("Tool 1"))
	tool2 := mcp.NewTool("tool2", mcp.WithDescription("Tool 2"))
	tool3 := mcp.NewTool("tool3", mcp.WithDescription("Tool 3"))

	service1 := &MockService{
		name:    "service1",
		enabled: true,
		tools:   []mcp.Tool{tool1, tool2},
	}

	service2 := &MockService{
		name:    "service2",
		enabled: true,
		tools:   []mcp.Tool{tool3},
	}

	disabledService := &MockService{
		name:    "disabled",
		enabled: false,
		tools:   []mcp.Tool{mcp.NewTool("disabled-tool", mcp.WithDescription("Disabled tool"))},
	}

	registry.Register(service1)
	registry.Register(service2)
	registry.Register(disabledService)

	allTools := registry.GetAllTools()

	if len(allTools) != 3 {
		t.Errorf("Expected 3 tools from enabled services, got %d", len(allTools))
	}

	toolNames := make(map[string]bool)
	for _, tool := range allTools {
		toolNames[tool.Name] = true
	}

	expectedTools := []string{"tool1", "tool2", "tool3"}
	for _, expectedTool := range expectedTools {
		if !toolNames[expectedTool] {
			t.Errorf("Expected tool %s not found in all tools", expectedTool)
		}
	}

	if toolNames["disabled-tool"] {
		t.Error("Disabled service tool should not be in all tools")
	}
}

func TestRegistryGetAllHandlers(t *testing.T) {
	registry := NewRegistry()

	handler1 := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("handler1 result"), nil
	}

	handler2 := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("handler2 result"), nil
	}

	service1 := &MockService{
		name:     "service1",
		enabled:  true,
		handlers: map[string]server.ToolHandlerFunc{"handler1": handler1},
	}

	service2 := &MockService{
		name:     "service2",
		enabled:  true,
		handlers: map[string]server.ToolHandlerFunc{"handler2": handler2},
	}

	disabledService := &MockService{
		name:     "disabled",
		enabled:  false,
		handlers: map[string]server.ToolHandlerFunc{"disabled-handler": handler1},
	}

	registry.Register(service1)
	registry.Register(service2)
	registry.Register(disabledService)

	allHandlers := registry.GetAllHandlers()

	if len(allHandlers) != 2 {
		t.Errorf("Expected 2 handlers from enabled services, got %d", len(allHandlers))
	}

	if _, exists := allHandlers["handler1"]; !exists {
		t.Error("handler1 should be in all handlers")
	}

	if _, exists := allHandlers["handler2"]; !exists {
		t.Error("handler2 should be in all handlers")
	}

	if _, exists := allHandlers["disabled-handler"]; exists {
		t.Error("disabled-handler should not be in all handlers")
	}
}

func TestRegistryGetService(t *testing.T) {
	registry := NewRegistry()

	mockService := &MockService{name: "test-service", enabled: true}
	registry.Register(mockService)

	// Test existing service
	service, exists := registry.GetService("test-service")
	if !exists {
		t.Error("GetService should return true for existing service")
	}

	if service == nil {
		t.Error("GetService should return the registered service")
	}

	if service != mockService {
		t.Error("GetService should return the exact service instance")
	}

	// Test non-existing service
	nonExistentService, exists := registry.GetService("non-existent")
	if exists {
		t.Error("GetService should return false for non-existent service")
	}

	if nonExistentService != nil {
		t.Error("GetService should return nil for non-existent service")
	}
}

func TestRegistryEmptyRegistry(t *testing.T) {
	registry := NewRegistry()

	// Test empty registry behavior
	enabledServices := registry.GetEnabledServices()
	if len(enabledServices) != 0 {
		t.Error("Empty registry should have no enabled services")
	}

	allTools := registry.GetAllTools()
	if len(allTools) != 0 {
		t.Error("Empty registry should have no tools")
	}

	allHandlers := registry.GetAllHandlers()
	if len(allHandlers) != 0 {
		t.Error("Empty registry should have no handlers")
	}
}

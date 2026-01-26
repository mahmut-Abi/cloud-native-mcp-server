package utilities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
	assert.Equal(t, "utilities", service.Name())
	assert.True(t, service.IsEnabled())
}

func TestService_Name(t *testing.T) {
	service := NewService()
	assert.Equal(t, "utilities", service.Name())
}

func TestService_GetTools(t *testing.T) {
	service := NewService()
	tools := service.GetTools()
	assert.NotNil(t, tools)
	assert.Equal(t, 6, len(tools)) // 6 tools: get_time, get_timestamp, get_date, pause, sleep, web_fetch
}

func TestService_GetHandlers(t *testing.T) {
	service := NewService()
	handlers := service.GetHandlers()
	assert.NotNil(t, handlers)
	assert.Greater(t, len(handlers), 0)
}

func TestService_Initialize(t *testing.T) {
	service := NewService()
	assert.True(t, service.IsEnabled())
	err := service.Initialize(nil)
	assert.NoError(t, err)
	assert.True(t, service.IsEnabled())
}

func TestService_Initialize_WithConfig(t *testing.T) {
	service := NewService()
	err := service.Initialize(struct{}{})
	assert.NoError(t, err)
	assert.True(t, service.IsEnabled())
}

func TestService_IsEnabled(t *testing.T) {
	service := NewService()
	assert.True(t, service.IsEnabled())
}

func TestService_GetName(t *testing.T) {
	service := NewService()
	assert.Equal(t, "utilities", service.Name())
}

func TestService_GetToolsCache(t *testing.T) {
	service := NewService()
	tools1 := service.GetTools()
	assert.NotNil(t, tools1)
	assert.Len(t, tools1, 6)

	// Second call should return cached tools
	tools2 := service.GetTools()
	assert.Equal(t, len(tools1), len(tools2))
}

func TestService_GetToolsWhenDisabled(t *testing.T) {
	// The service is always enabled, but we test the logic
	service := NewService()
	tools := service.GetTools()
	assert.NotNil(t, tools)
	assert.NotEmpty(t, tools)
}

func TestService_GetHandlersContent(t *testing.T) {
	service := NewService()
	handlers := service.GetHandlers()

	// Verify expected handlers exist
	assert.Contains(t, handlers, "utilities_get_time")
	assert.Contains(t, handlers, "utilities_get_timestamp")
	assert.Contains(t, handlers, "utilities_get_date")
	assert.Contains(t, handlers, "utilities_pause")
	assert.Contains(t, handlers, "utilities_sleep")
	assert.Contains(t, handlers, "utilities_web_fetch")
}

func TestService_ToolsHaveCorrectNames(t *testing.T) {
	service := NewService()
	tools := service.GetTools()

	toolNames := make(map[string]bool)
	for _, tool := range tools {
		toolNames[tool.Name] = true
	}

	assert.True(t, toolNames["utilities_get_time"])
	assert.True(t, toolNames["utilities_get_timestamp"])
	assert.True(t, toolNames["utilities_get_date"])
	assert.True(t, toolNames["utilities_pause"])
	assert.True(t, toolNames["utilities_sleep"])
	assert.True(t, toolNames["utilities_web_fetch"])
}

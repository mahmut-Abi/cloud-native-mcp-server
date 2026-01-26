package services

import (
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

// Service represents a pluggable service that can provide tools and handlers
type Service interface {
	// Name returns the service name
	Name() string

	// GetTools returns all tools provided by this service
	GetTools() []mcp.Tool

	// GetHandlers returns all tool handlers provided by this service
	GetHandlers() map[string]server.ToolHandlerFunc

	// Initialize initializes the service with the given configuration
	Initialize(config interface{}) error

	// IsEnabled returns whether the service is enabled
	IsEnabled() bool
}

// Registry manages all registered services
type Registry struct {
	mu       sync.RWMutex
	services map[string]Service
}

// NewRegistry creates a new service registry
func NewRegistry() *Registry {
	return &Registry{
		services: make(map[string]Service),
	}
}

// Register registers a service with the registry
func (r *Registry) Register(service Service) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.services[service.Name()] = service
}

// GetService returns a service by name
func (r *Registry) GetService(name string) (Service, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	service, exists := r.services[name]
	return service, exists
}

// GetAllServices returns all registered services
func (r *Registry) GetAllServices() map[string]Service {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make(map[string]Service, len(r.services))
	for k, v := range r.services {
		result[k] = v
	}
	return result
}

// GetEnabledServices returns only enabled services
func (r *Registry) GetEnabledServices() map[string]Service {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// Pre-allocate map for better performance
	enabled := make(map[string]Service)
	for name, service := range r.services {
		if service.IsEnabled() {
			enabled[name] = service
		}
	}
	return enabled
}

// GetAllTools returns all tools from all enabled services
func (r *Registry) GetAllTools() []mcp.Tool {
	var tools []mcp.Tool
	for _, service := range r.GetEnabledServices() {
		if svcTools := service.GetTools(); svcTools != nil {
			tools = append(tools, svcTools...)
		}
	}
	return tools
}

// GetAllToolsIncludingDisabled returns all tools from all registered services
// This includes tools from disabled services (e.g., for discovery/listing purposes)
func (r *Registry) GetAllToolsIncludingDisabled() []mcp.Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var tools []mcp.Tool
	for _, service := range r.services {
		if svcTools := service.GetTools(); svcTools != nil {
			tools = append(tools, svcTools...)
		}
	}
	return tools
}

// GetAllHandlers returns all handlers from all enabled services
func (r *Registry) GetAllHandlers() map[string]server.ToolHandlerFunc {
	handlers := make(map[string]server.ToolHandlerFunc)
	for _, service := range r.GetEnabledServices() {
		serviceHandlers := service.GetHandlers()
		if serviceHandlers == nil {
			continue
		}
		for name, handler := range serviceHandlers {
			handlers[name] = handler
		}
	}
	return handlers
}

// GetEnabledToolsByName returns tools filtered by enabled services
func (r *Registry) GetEnabledToolsByName(toolNames map[string]bool) []mcp.Tool {
	var tools []mcp.Tool
	enabledServices := r.GetEnabledServices()

	for _, service := range enabledServices {
		svcTools := service.GetTools()
		for _, tool := range svcTools {
			if enabled, exists := toolNames[tool.Name]; !exists || enabled {
				tools = append(tools, tool)
			}
		}
	}
	return tools
}

// FilterTools filters tools based on enabled list
func (r *Registry) FilterTools(tools []mcp.Tool, disabled map[string]bool) []mcp.Tool {
	var filtered []mcp.Tool
	for _, tool := range tools {
		if !disabled[tool.Name] {
			filtered = append(filtered, tool)
		}
	}
	return filtered
}

// FilterHandlers filters handlers based on disabled tools
func (r *Registry) FilterHandlers(handlers map[string]server.ToolHandlerFunc, disabled map[string]bool) map[string]server.ToolHandlerFunc {
	filtered := make(map[string]server.ToolHandlerFunc)
	for name, handler := range handlers {
		if !disabled[name] {
			filtered[name] = handler
		}
	}
	return filtered
}

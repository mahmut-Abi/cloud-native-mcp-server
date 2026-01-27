// Package utilities provides general-purpose utility tools for the MCP server.
// It implements tools for time operations, pausing, web fetching, and other common tasks.
package utilities

import (
	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/cache"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/utilities/handlers"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/utilities/tools"
)

// Service implements the utilities service for MCP server integration.
// It provides general-purpose tools like time operations, pausing, and web fetching.
type Service struct {
	enabled    bool              // Whether the service is enabled
	toolsCache *cache.ToolsCache // Cached tools to avoid recreation
}

// NewService creates a new utilities service instance.
// The service is enabled by default and requires no external configuration.
func NewService() *Service {
	return &Service{
		enabled:    true, // Always enabled
		toolsCache: cache.NewToolsCache(),
	}
}

// Name returns the service identifier used for registration and logging.
func (s *Service) Name() string {
	return "utilities"
}

// Initialize configures the utilities service.
// No external configuration is required for this service.
func (s *Service) Initialize(cfg interface{}) error {
	// Utilities service doesn't require any configuration
	// It's always enabled and ready to use
	s.enabled = true
	return nil
}

// IsEnabled returns whether the service is enabled.
func (s *Service) IsEnabled() bool {
	return s.enabled
}

// GetTools returns all available utilities MCP tools.
func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled {
		return nil
	}

	// Use unified cache
	return s.toolsCache.Get(func() []mcp.Tool {
		return []mcp.Tool{
			// Time operations
			tools.GetTimeTool(),
			tools.GetTimestampTool(),
			tools.GetDateTool(),

			// Pause/Sleep operations
			tools.PauseTool(),
			tools.SleepTool(),

			// Web operations
			tools.WebFetchTool(),
		}
	})
}

// GetHandlers returns all tool handlers mapped to their respective tool names.
func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
	return map[string]server.ToolHandlerFunc{
		// Time tools
		"utilities_get_time":      handlers.HandleGetTime,
		"utilities_get_timestamp": handlers.HandleGetTimestamp,
		"utilities_get_date":      handlers.HandleGetDate,

		// Pause/Sleep tools
		"utilities_pause": handlers.HandlePause,
		"utilities_sleep": handlers.HandleSleep,

		// Web tools
		"utilities_web_fetch": handlers.HandleWebFetch,
	}
}

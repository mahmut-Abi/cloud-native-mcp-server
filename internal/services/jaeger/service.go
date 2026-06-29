package jaeger

import (
	"context"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/cache"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/framework"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/jaeger/handlers"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/jaeger/tools"
	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
)

// Service provides Jaeger distributed tracing functionality.
// The backend client is not stored — it is created per-request from HTTP headers.
type Service struct {
	enabled       bool
	toolsCache    *cache.ToolsCache
	initFramework *framework.CommonServiceInit
}

// NewService creates a new Jaeger service.
func NewService() *Service {
	checker := framework.NewServiceEnabled(
		func(cfg *config.AppConfig) bool { return true },
		func(cfg *config.AppConfig) string { return "header-based-auth" },
	)

	initConfig := &framework.InitConfig{
		Required:     false,
		URLValidator: framework.SimpleURLValidator,
		ClientBuilder: nil,
	}

	return &Service{
		enabled:       false,
		toolsCache:    cache.NewToolsCache(),
		initFramework: framework.NewCommonServiceInit("Jaeger", initConfig, checker),
	}
}

// Name returns the service name.
func (s *Service) Name() string {
	return "jaeger"
}

// GetTools returns all available tools for the Jaeger service.
func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled {
		return []mcp.Tool{}
	}

	return s.toolsCache.Get(func() []mcp.Tool {
		return []mcp.Tool{
			tools.GetTracesSummaryTool(),
			tools.GetServicesSummaryTool(),
			tools.GetTracesTool(),
			tools.GetTraceTool(),
			tools.GetServicesTool(),
			tools.GetServiceOperationsTool(),
			tools.SearchTracesTool(),
			tools.GetDependenciesTool(),
		}
	})
}

// GetHandlers returns all tool handlers for the Jaeger service.
func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
	handlersMap := map[string]server.ToolHandlerFunc{
		"jaeger_get_traces":           handlers.GetTracesHandler(),
		"jaeger_get_trace":            handlers.GetTraceHandler(),
		"jaeger_get_services":         handlers.GetServicesHandler(),
		"jaeger_get_service_ops":      handlers.GetServiceOperationsHandler(),
		"jaeger_search_traces":        handlers.SearchTracesHandler(),
		"jaeger_get_dependencies":     handlers.GetDependenciesHandler(),
		"jaeger_get_traces_summary":   handlers.GetTracesSummaryHandler(),
		"jaeger_get_services_summary": handlers.GetServicesSummaryHandler(),
	}

	for name, handler := range handlersMap {
		handlersMap[name] = s.wrapToolErrors(name, handler)
	}

	return handlersMap
}

// Initialize configures the Jaeger service with the provided application configuration.
// The backend client is created per-request from HTTP headers (see client/config.go).
func (s *Service) Initialize(cfg interface{}) error {
	return s.initFramework.Initialize(cfg,
		func(enabled bool) { s.enabled = enabled },
		func(_ interface{}) {
			// Backend client is created per-request from HTTP headers.
			// The backend auth handler was registered in client/config.go init().
		},
	)
}

// IsEnabled returns whether the Jaeger service is enabled.
func (s *Service) IsEnabled() bool {
	return s.enabled
}

// GetToolsCache returns the tools cache.
func (s *Service) GetToolsCache() *cache.ToolsCache {
	return s.toolsCache
}

func (s *Service) wrapToolErrors(toolName string, handler server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handler(ctx, request)
		if err != nil {
			logrus.WithError(err).WithField("tool", toolName).Warn("Tool execution failed")
			return mcp.NewToolResultError(err.Error()), nil
		}
		return result, nil
	}
}

package jaeger

import (
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/cache"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/framework"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/jaeger/client"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/jaeger/handlers"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/jaeger/tools"
	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

// Service provides Jaeger distributed tracing functionality.
type Service struct {
	client        *client.Client
	enabled       bool
	toolsCache    *cache.ToolsCache
	initFramework *framework.CommonServiceInit
}

// NewService creates a new Jaeger service.
func NewService() *Service {
	// Create service enable checker
	checker := framework.NewServiceEnabled(
		func(cfg *config.AppConfig) bool { return cfg.Jaeger.Enabled },
		func(cfg *config.AppConfig) string { return cfg.Jaeger.Address },
	)

	// Create init configuration
	initConfig := &framework.InitConfig{
		Required:     false,
		URLValidator: framework.SimpleURLValidator,
		ClientBuilder: func(cfg *config.AppConfig) (interface{}, error) {
			return client.NewClient(&client.ClientOptions{
				BaseURL: cfg.Jaeger.Address,
				Timeout: time.Duration(cfg.Jaeger.TimeoutSec) * time.Second,
			})
		},
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
	// Return empty tools if service is not enabled
	if !s.IsEnabled() {
		return []mcp.Tool{}
	}

	return s.toolsCache.Get(func() []mcp.Tool {
		return []mcp.Tool{
			// Summary tools (recommended for LLM efficiency)
			tools.GetTracesSummaryTool(),
			tools.GetServicesSummaryTool(),

			// Standard tools
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
	return map[string]server.ToolHandlerFunc{
		"jaeger_get_traces":           handlers.GetTracesHandler(s),
		"jaeger_get_trace":            handlers.GetTraceHandler(s),
		"jaeger_get_services":         handlers.GetServicesHandler(s),
		"jaeger_get_service_ops":      handlers.GetServiceOperationsHandler(s),
		"jaeger_search_traces":        handlers.SearchTracesHandler(s),
		"jaeger_get_dependencies":     handlers.GetDependenciesHandler(s),
		"jaeger_get_traces_summary":   handlers.GetTracesSummaryHandler(s),
		"jaeger_get_services_summary": handlers.GetServicesSummaryHandler(s),
	}
}

// Initialize initializes the Jaeger service with the given configuration.
func (s *Service) Initialize(cfg interface{}) error {
	return s.initFramework.Initialize(cfg,
		func(enabled bool) { s.enabled = enabled },
		func(clientIface interface{}) {
			if jaegerClient, ok := clientIface.(*client.Client); ok {
				s.client = jaegerClient
			}
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

// GetClient returns the Jaeger client.
func (s *Service) GetClient() *client.Client {
	return s.client
}

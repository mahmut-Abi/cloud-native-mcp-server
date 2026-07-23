// Package victoriametrics provides VictoriaMetrics monitoring integration for the MCP server.
// It implements tools for querying VictoriaMetrics metrics using MetricsQL, bulk export/import,
// and managing VictoriaMetrics instances.
package victoriametrics

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/cache"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/framework"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/victoriametrics/handlers"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/victoriametrics/tools"
)

// Service implements the VictoriaMetrics service for MCP server integration.
// It provides tools and handlers for interacting with VictoriaMetrics instances.
// The backend client is not stored — it is created per-request from HTTP headers.
type Service struct {
	enabled       bool                         // Whether the service is enabled
	toolsCache    *cache.ToolsCache            // Cached tools to avoid recreation
	initFramework *framework.CommonServiceInit // Common initialization framework
}

// NewService creates a new VictoriaMetrics service instance.
// The service is disabled by default and requires initialization before use.
func NewService() *Service {
	// Create service enable checker
	checker := framework.NewServiceEnabled(
		func(cfg *config.AppConfig) bool { return true },
		func(cfg *config.AppConfig) string { return "header-based-auth" },
	)

	// Create init configuration
	initConfig := &framework.InitConfig{
		Required:      false,
		URLValidator:  func(url string) bool { return url != "" },
		ClientBuilder: nil,
	}

	return &Service{
		enabled:       false, // Default disabled until configured
		toolsCache:    cache.NewToolsCache(),
		initFramework: framework.NewCommonServiceInit("VictoriaMetrics", initConfig, checker),
	}
}

// Name returns the service identifier used for registration and logging.
func (s *Service) Name() string {
	return "victoriametrics"
}

// Initialize configures the VictoriaMetrics service with the provided application configuration.
// The backend client is created per-request from HTTP headers (see client/config.go).
func (s *Service) Initialize(cfg interface{}) error {
	return s.initFramework.Initialize(cfg,
		func(enabled bool) { s.enabled = enabled },
		func(_ interface{}) {
			// Backend client is created per-request from HTTP headers.
			// The backend auth middleware (client/config.go) injects the client into the request context.
		},
	)
}

// IsEnabled returns whether the VictoriaMetrics service is currently enabled.
func (s *Service) IsEnabled() bool {
	return s.enabled
}

// GetTools returns all MCP tools provided by the VictoriaMetrics service.
// Tools are cached for performance and include query, export, import, and management operations.
func (s *Service) GetTools() []mcp.Tool {
	return s.toolsCache.Get(tools.AllTools)
}

// GetHandlers returns all MCP tool handlers for the VictoriaMetrics service.
func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
	return map[string]server.ToolHandlerFunc{
		"victoriametrics_query":              handlers.HandleQuery(),
		"victoriametrics_query_range":        handlers.HandleQueryRange(),
		"victoriametrics_export":             handlers.HandleExport(),
		"victoriametrics_import":             handlers.HandleImport(),
		"victoriametrics_series_count":       handlers.HandleSeriesCount(),
		"victoriametrics_get_label_names":    handlers.HandleLabelNames(),
		"victoriametrics_get_label_values":   handlers.HandleLabelValues(),
		"victoriametrics_get_tsdb_stats":     handlers.HandleTSDBStats(),
		"victoriametrics_get_active_queries": handlers.HandleActiveQueries(),
		"victoriametrics_delete_series":      handlers.HandleDeleteSeries(),
		"victoriametrics_health_check":       handlers.HandleHealthCheck(),
		"victoriametrics_get_metadata":       handlers.HandleGetMetadata(),
		"victoriametrics_list_tools":         handlers.HandleListTools(),
	}
}

// Ensure json import is used for proper MCP result handling.
var _ = json.Marshal

// Ensure context and logrus are available for service operations.
var _ = context.Background
var _ = logrus.Info

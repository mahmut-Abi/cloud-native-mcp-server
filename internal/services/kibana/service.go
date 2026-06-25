// Package kibana provides Kibana integration for the MCP server.
// It implements tools for managing Kibana spaces, dashboards, visualizations, and saved objects.
package kibana

import (
	"context"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/cache"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/framework"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/kibana/client"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/kibana/handlers"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/kibana/tools"
)

// Service implements the Kibana service for MCP server integration.
// It provides tools and handlers for interacting with Kibana instances.
// The backend client is not stored — it is created per-request from HTTP headers.
type Service struct {
	enabled       bool                         // Whether the service is enabled
	toolsCache    *cache.ToolsCache            // Cached tools to avoid recreation
	initFramework *framework.CommonServiceInit // Common initialization framework
}

// NewService creates a new Kibana service instance.
// The service is disabled by default and requires initialization before use.
func NewService() *Service {
	// Create service enable checker
	checker := framework.NewServiceEnabled(
		func(cfg *config.AppConfig) bool { return cfg.Kibana.Enabled },
		func(cfg *config.AppConfig) string { return cfg.Kibana.URL },
	)

	// Create init configuration
	initConfig := &framework.InitConfig{
		Required:     false,
		URLValidator: framework.SimpleURLValidator,
		ClientBuilder: func(cfg *config.AppConfig) (interface{}, error) {
			return client.NewClient(&client.ClientOptions{
				URL:        cfg.Kibana.URL,
				APIKey:     cfg.Kibana.APIKey,
				Username:   cfg.Kibana.Username,
				Password:   cfg.Kibana.Password,
				Timeout:    time.Duration(cfg.Kibana.TimeoutSec) * time.Second,
				SkipVerify: cfg.Kibana.SkipVerify,
				Space:      cfg.Kibana.Space,
			})
		},
	}

	return &Service{
		enabled:       false, // Default disabled until configured
		toolsCache:    cache.NewToolsCache(),
		initFramework: framework.NewCommonServiceInit("Kibana", initConfig, checker),
	}
}

// Name returns the service identifier used for registration and logging.
func (s *Service) Name() string {
	return "kibana"
}

// Initialize configures the Kibana service with the provided application configuration.
// It uses the common service framework for standardized initialization.
func (s *Service) Initialize(cfg interface{}) error {
	return s.initFramework.Initialize(cfg,
		func(enabled bool) { s.enabled = enabled },
		func(_ interface{}) {
			// Backend client is created per-request from HTTP headers.
			// The backend auth handler was registered in client/config.go init().
		},
	)
}

// GetTools returns all available Kibana MCP tools.
// Tools are only returned if the service is enabled and properly initialized.
// The tools include space management, dashboard operations, and saved object management capabilities.
func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled {
		return nil
	}

	// Use unified cache
	return s.toolsCache.Get(func() []mcp.Tool {
		// Legacy tools (maintained for compatibility)
		legacyTools := []mcp.Tool{
			tools.GetSpacesTool(),
			tools.GetSpaceTool(),
			tools.GetIndexPatternsTool(),
			tools.GetDashboardsTool(),
			tools.GetDashboardTool(),
			tools.GetVisualizationsTool(),
			tools.GetVisualizationTool(),
			tools.GetIndexPatternTool(),
			tools.GetSavedSearchesTool(),
			tools.GetSavedSearchTool(),
			tools.SearchSavedObjectsTool(),
			tools.TestConnectionTool(),
			tools.GetKibanaStatusTool(),
		}

		// ⚠️ PRIORITY: New optimized tools for LLM efficiency
		optimizedTools := []mcp.Tool{
			tools.GetSpacesSummaryTool(),
			tools.GetDashboardsSummaryTool(),
			tools.GetVisualizationsSummaryTool(),
			tools.GetIndexPatternsSummaryTool(),
			tools.GetDashboardsPaginatedTool(),
			tools.GetVisualizationsPaginatedTool(),
			tools.GetSavedObjectsAdvancedTool(),
			tools.GetDashboardDetailAdvancedTool(),
			tools.GetKibanaHealthSummaryTool(),

			// Analysis & Discovery tools
			tools.QueryLogsTool(),
			tools.GetCanvasWorkpadsTool(),
			tools.GetLensObjectsTool(),
			tools.GetMapsTool(),
			tools.GetKibanaAlertsTool(),
			tools.GetIndexPatternFieldsTool(),

			// ============ Write Operations: Spaces ============
			tools.CreateSpaceTool(),
			tools.UpdateSpaceTool(),
			tools.DeleteSpaceTool(),

			// ============ Write Operations: Index Patterns ============
			tools.CreateIndexPatternTool(),
			tools.UpdateIndexPatternTool(),
			tools.DeleteIndexPatternTool(),
			tools.SetDefaultIndexPatternTool(),
			tools.RefreshIndexPatternFieldsTool(),

			// ============ Write Operations: Dashboards ============
			tools.CreateDashboardTool(),
			tools.UpdateDashboardTool(),
			tools.DeleteDashboardTool(),
			tools.CloneDashboardTool(),

			// ============ Write Operations: Visualizations ============
			tools.CreateVisualizationTool(),
			tools.UpdateVisualizationTool(),
			tools.DeleteVisualizationTool(),
			tools.CloneVisualizationTool(),

			// ============ Write Operations: Saved Objects (Generic) ============
			tools.CreateSavedObjectTool(),
			tools.UpdateSavedObjectTool(),
			tools.DeleteSavedObjectTool(),
			tools.BulkDeleteSavedObjectsTool(),
			tools.ExportSavedObjectsTool(),
			tools.ImportSavedObjectsTool(),

			// ============ Alert Rules ============
			tools.GetAlertRulesTool(),
			tools.GetAlertRuleTool(),
			tools.CreateAlertRuleTool(),
			tools.UpdateAlertRuleTool(),
			tools.DeleteAlertRuleTool(),
			tools.EnableAlertRuleTool(),
			tools.DisableAlertRuleTool(),
			tools.MuteAlertRuleTool(),
			tools.UnmuteAlertRuleTool(),
			tools.GetAlertRuleTypesTool(),
			tools.GetAlertRuleHistoryTool(),

			// ============ Connectors ============
			tools.GetConnectorsTool(),
			tools.GetConnectorTool(),
			tools.CreateConnectorTool(),
			tools.UpdateConnectorTool(),
			tools.DeleteConnectorTool(),
			tools.TestConnectorTool(),
			tools.GetConnectorTypesTool(),

			// ============ Data Views ============
			tools.GetDataViewsTool(),
			tools.GetDataViewTool(),
			tools.CreateDataViewTool(),
			tools.UpdateDataViewTool(),
			tools.DeleteDataViewTool(),
		}

		// Combine all tools - optimized tools first for better visibility
		return append(optimizedTools, legacyTools...)
	})
}

// GetHandlers returns all tool handlers mapped to their respective tool names.
// Handlers are only returned if the service is enabled and properly initialized.
func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
	if !s.enabled {
		return nil
	}

	// Legacy handlers (maintained for compatibility)
	legacyHandlers := map[string]server.ToolHandlerFunc{
		"kibana_get_spaces":           handlers.HandleGetSpaces(),
		"kibana_get_space":            handlers.HandleGetSpace(),
		"kibana_get_index_patterns":   handlers.HandleGetIndexPatterns(),
		"kibana_get_index_pattern":    handlers.HandleGetIndexPattern(),
		"kibana_get_dashboards":       handlers.HandleGetDashboards(),
		"kibana_get_dashboard":        handlers.HandleGetDashboard(),
		"kibana_get_visualizations":   handlers.HandleGetVisualizations(),
		"kibana_get_visualization":    handlers.HandleGetVisualization(),
		"kibana_get_saved_searches":   handlers.HandleGetSavedSearches(),
		"kibana_get_saved_search":     handlers.HandleGetSavedSearch(),
		"kibana_search_saved_objects": handlers.HandleSearchSavedObjects(),
		"kibana_test_connection":      handlers.HandleTestConnection(),
		"kibana_get_status":           handlers.HandleGetKibanaStatus(),
	}

	// ⚠️ PRIORITY: New optimized handlers for LLM efficiency
	optimizedHandlers := map[string]server.ToolHandlerFunc{
		// Summary tools
		"kibana_spaces_summary":         handlers.HandleSpacesSummary(),
		"kibana_dashboards_summary":     handlers.HandleDashboardsPaginated(),
		"kibana_visualizations_summary": handlers.HandleVisualizationsPaginated(),
		"kibana_index_patterns_summary": handlers.HandleGetIndexPatterns(),

		// Advanced tools
		"kibana_dashboards_paginated":          handlers.HandleDashboardsPaginated(),
		"kibana_visualizations_paginated":      handlers.HandleVisualizationsPaginated(),
		"kibana_search_saved_objects_advanced": handlers.HandleSearchSavedObjectsAdvanced(),
		"kibana_get_dashboard_detail_advanced": handlers.HandleGetDashboardDetailAdvanced(),
		"kibana_health_summary":                handlers.HandleGetHealthSummary(),

		// Analysis & Discovery handlers
		"kibana_query_logs":               handlers.HandleQueryLogs(),
		"kibana_get_canvas_workpads":      handlers.HandleGetCanvasWorkpads(),
		"kibana_get_lens_objects":         handlers.HandleGetLensObjects(),
		"kibana_get_maps":                 handlers.HandleGetMaps(),
		"kibana_get_alerts":               handlers.HandleGetKibanaAlerts(),
		"kibana_get_index_pattern_fields": handlers.HandleGetIndexPatternFields(),

		// ============ Write Operations: Spaces ============
		"kibana_create_space": handlers.HandleCreateSpace(),
		"kibana_update_space": handlers.HandleUpdateSpace(),
		"kibana_delete_space": handlers.HandleDeleteSpace(),

		// ============ Write Operations: Index Patterns ============
		"kibana_create_index_pattern":         handlers.HandleCreateIndexPattern(),
		"kibana_update_index_pattern":         handlers.HandleUpdateIndexPattern(),
		"kibana_delete_index_pattern":         handlers.HandleDeleteIndexPattern(),
		"kibana_set_default_index_pattern":    handlers.HandleSetDefaultIndexPattern(),
		"kibana_refresh_index_pattern_fields": handlers.HandleRefreshIndexPatternFields(),

		// ============ Write Operations: Dashboards ============
		"kibana_create_dashboard": handlers.HandleCreateDashboard(),
		"kibana_update_dashboard": handlers.HandleUpdateDashboard(),
		"kibana_delete_dashboard": handlers.HandleDeleteDashboard(),
		"kibana_clone_dashboard":  handlers.HandleCloneDashboard(),

		// ============ Write Operations: Visualizations ============
		"kibana_create_visualization": handlers.HandleCreateVisualization(),
		"kibana_update_visualization": handlers.HandleUpdateVisualization(),
		"kibana_delete_visualization": handlers.HandleDeleteVisualization(),
		"kibana_clone_visualization":  handlers.HandleCloneVisualization(),

		// ============ Write Operations: Saved Objects (Generic) ============
		"kibana_create_saved_object":       handlers.HandleCreateSavedObject(),
		"kibana_update_saved_object":       handlers.HandleUpdateSavedObject(),
		"kibana_delete_saved_object":       handlers.HandleDeleteSavedObject(),
		"kibana_bulk_delete_saved_objects": handlers.HandleBulkDeleteSavedObjects(),
		"kibana_export_saved_objects":      handlers.HandleExportSavedObjects(),
		"kibana_import_saved_objects":      handlers.HandleImportSavedObjects(),

		// ============ Alert Rules ============
		"kibana_get_alert_rules":        handlers.HandleGetAlertRules(),
		"kibana_get_alert_rule":         handlers.HandleGetAlertRule(),
		"kibana_create_alert_rule":      handlers.HandleCreateAlertRule(),
		"kibana_update_alert_rule":      handlers.HandleUpdateAlertRule(),
		"kibana_delete_alert_rule":      handlers.HandleDeleteAlertRule(),
		"kibana_enable_alert_rule":      handlers.HandleEnableAlertRule(),
		"kibana_disable_alert_rule":     handlers.HandleDisableAlertRule(),
		"kibana_mute_alert_rule":        handlers.HandleMuteAlertRule(),
		"kibana_unmute_alert_rule":      handlers.HandleUnmuteAlertRule(),
		"kibana_get_alert_rule_types":   handlers.HandleGetAlertRuleTypes(),
		"kibana_get_alert_rule_history": handlers.HandleGetAlertRuleHistory(),

		// ============ Connectors ============
		"kibana_get_connectors":      handlers.HandleGetConnectors(),
		"kibana_get_connector":       handlers.HandleGetConnector(),
		"kibana_create_connector":    handlers.HandleCreateConnector(),
		"kibana_update_connector":    handlers.HandleUpdateConnector(),
		"kibana_delete_connector":    handlers.HandleDeleteConnector(),
		"kibana_test_connector":      handlers.HandleTestConnector(),
		"kibana_get_connector_types": handlers.HandleGetConnectorTypes(),

		// ============ Data Views ============
		"kibana_get_data_views":   handlers.HandleGetDataViews(),
		"kibana_get_data_view":    handlers.HandleGetDataView(),
		"kibana_create_data_view": handlers.HandleCreateDataView(),
		"kibana_update_data_view": handlers.HandleUpdateDataView(),
		"kibana_delete_data_view": handlers.HandleDeleteDataView(),
	}

	// Combine all handlers
	allHandlers := make(map[string]server.ToolHandlerFunc)
	for k, v := range optimizedHandlers {
		allHandlers[k] = s.wrapToolErrors(k, v)
	}
	for k, v := range legacyHandlers {
		allHandlers[k] = s.wrapToolErrors(k, v)
	}

	return allHandlers
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

// IsEnabled returns whether the service is enabled and ready for use.
// A service is considered enabled if it's marked as enabled and has a valid client.
func (s *Service) IsEnabled() bool {
	return s.enabled
}

// GetClient returns the underlying Kibana client for advanced operations.
// This method is primarily used for testing and internal service communication.
func (s *Service) GetClient() *client.Client {
	return nil // Backend client is created per-request from HTTP headers
}

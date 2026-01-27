// Package kibana provides Kibana integration for the MCP server.
// It implements tools for managing Kibana spaces, dashboards, visualizations, and saved objects.
package kibana

import (
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/cache"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/framework"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/kibana/client"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/kibana/handlers"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/kibana/tools"
)

// Service implements the Kibana service for MCP server integration.
// It provides tools and handlers for interacting with Kibana instances.
type Service struct {
	client        *client.Client               // Kibana HTTP client for API operations
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
		func(clientIface interface{}) {
			if kibanaClient, ok := clientIface.(*client.Client); ok {
				s.client = kibanaClient
			}
		},
	)
}

// GetTools returns all available Kibana MCP tools.
// Tools are only returned if the service is enabled and properly initialized.
// The tools include space management, dashboard operations, and saved object management capabilities.
func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled || s.client == nil {
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
	if !s.enabled || s.client == nil {
		return nil
	}

	// Legacy handlers (maintained for compatibility)
	legacyHandlers := map[string]server.ToolHandlerFunc{
		"kibana_get_spaces":           handlers.HandleGetSpaces(s.client),
		"kibana_get_space":            handlers.HandleGetSpace(s.client),
		"kibana_get_index_patterns":   handlers.HandleGetIndexPatterns(s.client),
		"kibana_get_index_pattern":    handlers.HandleGetIndexPattern(s.client),
		"kibana_get_dashboards":       handlers.HandleGetDashboards(s.client),
		"kibana_get_dashboard":        handlers.HandleGetDashboard(s.client),
		"kibana_get_visualizations":   handlers.HandleGetVisualizations(s.client),
		"kibana_get_visualization":    handlers.HandleGetVisualization(s.client),
		"kibana_get_saved_searches":   handlers.HandleGetSavedSearches(s.client),
		"kibana_get_saved_search":     handlers.HandleGetSavedSearch(s.client),
		"kibana_search_saved_objects": handlers.HandleSearchSavedObjects(s.client),
		"kibana_test_connection":      handlers.HandleTestConnection(s.client),
		"kibana_get_status":           handlers.HandleGetKibanaStatus(s.client),
	}

	// ⚠️ PRIORITY: New optimized handlers for LLM efficiency
	optimizedHandlers := map[string]server.ToolHandlerFunc{
		// Summary tools
		"kibana_spaces_summary":         handlers.HandleSpacesSummary(s.client),
		"kibana_dashboards_summary":     handlers.HandleDashboardsPaginated(s.client),
		"kibana_visualizations_summary": handlers.HandleVisualizationsPaginated(s.client),
		"kibana_index_patterns_summary": handlers.HandleGetIndexPatterns(s.client),

		// Advanced tools
		"kibana_dashboards_paginated":          handlers.HandleDashboardsPaginated(s.client),
		"kibana_visualizations_paginated":      handlers.HandleVisualizationsPaginated(s.client),
		"kibana_search_saved_objects_advanced": handlers.HandleSearchSavedObjectsAdvanced(s.client),
		"kibana_get_dashboard_detail_advanced": handlers.HandleGetDashboardDetailAdvanced(s.client),
		"kibana_health_summary":                handlers.HandleGetHealthSummary(s.client),

		// Analysis & Discovery handlers
		"kibana_query_logs":               handlers.HandleQueryLogs(s.client),
		"kibana_get_canvas_workpads":      handlers.HandleGetCanvasWorkpads(s.client),
		"kibana_get_lens_objects":         handlers.HandleGetLensObjects(s.client),
		"kibana_get_maps":                 handlers.HandleGetMaps(s.client),
		"kibana_get_alerts":               handlers.HandleGetKibanaAlerts(s.client),
		"kibana_get_index_pattern_fields": handlers.HandleGetIndexPatternFields(s.client),

		// ============ Write Operations: Spaces ============
		"kibana_create_space": handlers.HandleCreateSpace(s.client),
		"kibana_update_space": handlers.HandleUpdateSpace(s.client),
		"kibana_delete_space": handlers.HandleDeleteSpace(s.client),

		// ============ Write Operations: Index Patterns ============
		"kibana_create_index_pattern":         handlers.HandleCreateIndexPattern(s.client),
		"kibana_update_index_pattern":         handlers.HandleUpdateIndexPattern(s.client),
		"kibana_delete_index_pattern":         handlers.HandleDeleteIndexPattern(s.client),
		"kibana_set_default_index_pattern":    handlers.HandleSetDefaultIndexPattern(s.client),
		"kibana_refresh_index_pattern_fields": handlers.HandleRefreshIndexPatternFields(s.client),

		// ============ Write Operations: Dashboards ============
		"kibana_create_dashboard": handlers.HandleCreateDashboard(s.client),
		"kibana_update_dashboard": handlers.HandleUpdateDashboard(s.client),
		"kibana_delete_dashboard": handlers.HandleDeleteDashboard(s.client),
		"kibana_clone_dashboard":  handlers.HandleCloneDashboard(s.client),

		// ============ Write Operations: Visualizations ============
		"kibana_create_visualization": handlers.HandleCreateVisualization(s.client),
		"kibana_update_visualization": handlers.HandleUpdateVisualization(s.client),
		"kibana_delete_visualization": handlers.HandleDeleteVisualization(s.client),
		"kibana_clone_visualization":  handlers.HandleCloneVisualization(s.client),

		// ============ Write Operations: Saved Objects (Generic) ============
		"kibana_create_saved_object":       handlers.HandleCreateSavedObject(s.client),
		"kibana_update_saved_object":       handlers.HandleUpdateSavedObject(s.client),
		"kibana_delete_saved_object":       handlers.HandleDeleteSavedObject(s.client),
		"kibana_bulk_delete_saved_objects": handlers.HandleBulkDeleteSavedObjects(s.client),
		"kibana_export_saved_objects":      handlers.HandleExportSavedObjects(s.client),
		"kibana_import_saved_objects":      handlers.HandleImportSavedObjects(s.client),

		// ============ Alert Rules ============
		"kibana_get_alert_rules":        handlers.HandleGetAlertRules(s.client),
		"kibana_get_alert_rule":         handlers.HandleGetAlertRule(s.client),
		"kibana_create_alert_rule":      handlers.HandleCreateAlertRule(s.client),
		"kibana_update_alert_rule":      handlers.HandleUpdateAlertRule(s.client),
		"kibana_delete_alert_rule":      handlers.HandleDeleteAlertRule(s.client),
		"kibana_enable_alert_rule":      handlers.HandleEnableAlertRule(s.client),
		"kibana_disable_alert_rule":     handlers.HandleDisableAlertRule(s.client),
		"kibana_mute_alert_rule":        handlers.HandleMuteAlertRule(s.client),
		"kibana_unmute_alert_rule":      handlers.HandleUnmuteAlertRule(s.client),
		"kibana_get_alert_rule_types":   handlers.HandleGetAlertRuleTypes(s.client),
		"kibana_get_alert_rule_history": handlers.HandleGetAlertRuleHistory(s.client),

		// ============ Connectors ============
		"kibana_get_connectors":      handlers.HandleGetConnectors(s.client),
		"kibana_get_connector":       handlers.HandleGetConnector(s.client),
		"kibana_create_connector":    handlers.HandleCreateConnector(s.client),
		"kibana_update_connector":    handlers.HandleUpdateConnector(s.client),
		"kibana_delete_connector":    handlers.HandleDeleteConnector(s.client),
		"kibana_test_connector":      handlers.HandleTestConnector(s.client),
		"kibana_get_connector_types": handlers.HandleGetConnectorTypes(s.client),

		// ============ Data Views ============
		"kibana_get_data_views":   handlers.HandleGetDataViews(s.client),
		"kibana_get_data_view":    handlers.HandleGetDataView(s.client),
		"kibana_create_data_view": handlers.HandleCreateDataView(s.client),
		"kibana_update_data_view": handlers.HandleUpdateDataView(s.client),
		"kibana_delete_data_view": handlers.HandleDeleteDataView(s.client),
	}

	// Combine all handlers
	allHandlers := make(map[string]server.ToolHandlerFunc)
	for k, v := range optimizedHandlers {
		allHandlers[k] = v
	}
	for k, v := range legacyHandlers {
		allHandlers[k] = v
	}

	return allHandlers
}

// IsEnabled returns whether the service is enabled and ready for use.
// A service is considered enabled if it's marked as enabled and has a valid client.
func (s *Service) IsEnabled() bool {
	return s.enabled && s.client != nil
}

// GetClient returns the underlying Kibana client for advanced operations.
// This method is primarily used for testing and internal service communication.
func (s *Service) GetClient() *client.Client {
	return s.client
}

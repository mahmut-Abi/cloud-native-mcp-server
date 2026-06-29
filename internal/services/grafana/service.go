// Package grafana provides Grafana dashboard and monitoring integration for the MCP server.
// It implements tools for managing Grafana dashboards, data sources, folders, and alert rules.
package grafana

import (

	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/cache"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/framework"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/grafana/client"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/grafana/handlers"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/grafana/tools"
)

// Service implements the Grafana service for MCP server integration.
// It provides tools and handlers for interacting with Grafana instances.
// The backend client is not stored — it is created per-request from HTTP headers.
type Service struct {
	enabled       bool                         // Whether the service is enabled
	toolsCache    *cache.ToolsCache            // Cached tools to avoid recreation
	initFramework *framework.CommonServiceInit // Common initialization framework
}

// NewService creates a new Grafana service instance.
// The service is disabled by default and requires initialization before use.
func NewService() *Service {
	// Create service enable checker
	checker := framework.NewServiceEnabled(
		func(cfg *config.AppConfig) bool { return true },
		func(cfg *config.AppConfig) string { return "header-based-auth" },
	)

	// Create init configuration
	initConfig := &framework.InitConfig{
		Required:     false,
		URLValidator: framework.SimpleURLValidator,
		ClientBuilder: nil,
	}

	return &Service{
		enabled:       false, // Default disabled until configured
		toolsCache:    cache.NewToolsCache(),
		initFramework: framework.NewCommonServiceInit("Grafana", initConfig, checker),
	}
}

// Name returns the service identifier used for registration and logging.
func (s *Service) Name() string {
	return "grafana"
}

// Initialize configures the Grafana service with the provided application configuration.
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

// GetTools returns all available Grafana MCP tools with optimization focus.
// Tools are only returned if the service is enabled and properly initialized.
// Organized with summary tools first for better user experience.
func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled {
		return nil
	}

	// Use unified cache
	return s.toolsCache.Get(func() []mcp.Tool {
		return []mcp.Tool{
			// 🎯 SUMMARY TOOLS - Recommended first (95% of use cases)
			tools.GetDashboardsSummaryTool(),  // RECOMMENDED for dashboard discovery
			tools.GetDataSourcesSummaryTool(), // RECOMMENDED for datasource discovery
			tools.GetPluginsSummaryTool(),     // RECOMMENDED for plugin discovery

			// 📋 STANDARD TOOLS - For detailed information when needed
			tools.GetDashboardsTool(),
			tools.GetDataSourcesTool(),
			tools.GetPluginsTool(),
			tools.GetFoldersTool(),
			tools.CreateFolderTool(),
			tools.UpdateFolderTool(),
			tools.DeleteFolderTool(),
			tools.SearchDashboardsTool(),

			// 🔍 SPECIFIC RESOURCE TOOLS - For detailed inspection
			tools.GetDashboardTool(),
			tools.GetFolderTool(),
			tools.GetDataSourceTool(),
			tools.GetDataSourceByNameTool(),
			tools.GetPluginTool(),

			// 🚨 MONITORING TOOLS
			tools.GetAlertRulesTool(),
			tools.GetAlertRuleByUIDTool(),
			tools.ListContactPointsTool(),

			// 📝 ANNOTATION TOOLS
			tools.GetAnnotationsTool(),
			tools.CreateAnnotationTool(),
			tools.UpdateAnnotationTool(),
			tools.PatchAnnotationTool(),
			tools.DeleteAnnotationTool(),
			tools.GetAnnotationTagsTool(),

			// 🔧 UTILITY TOOLS
			tools.TestConnectionTool(),
			tools.GetCurrentUserTool(),
			tools.GetUsersTool(),
			tools.GetOrganizationTool(),
			tools.CheckDatasourceHealthTool(),

			// 🔐 ADMIN TOOLS
			tools.ListTeamsTool(),
			tools.ListAllRolesTool(),
			tools.GetRoleDetailsTool(),
			tools.GetRoleAssignmentsTool(),
			tools.ListUserRolesTool(),
			tools.ListTeamRolesTool(),
			tools.GetResourcePermissionsTool(),
			tools.GetResourceDescriptionTool(),

			// 📊 DASHBOARD MANAGEMENT TOOLS
			tools.UpdateDashboardTool(),
			tools.GetDashboardVersionsTool(),
			tools.GetDashboardVersionTool(),
			tools.RestoreDashboardVersionTool(),
			tools.DeleteDashboardTool(),
			tools.GetDashboardPanelQueriesTool(),
			tools.GetDashboardPropertyTool(),

			// 🚨 ALERTING TOOLS
			tools.CreateAlertRuleTool(),
			tools.UpdateAlertRuleTool(),
			tools.DeleteAlertRuleTool(),

			// 🔗 NAVIGATION TOOLS
			tools.GenerateDeeplinkTool(),
			tools.GenerateLogsDrilldownLinkTool(),

			// 🎨 RENDERING TOOLS
			tools.RenderPanelImageTool(),

			// 📝 GRAPHITE ANNOTATION TOOLS
			tools.CreateGraphiteAnnotationTool(),

			// 🔧 DATASOURCE MANAGEMENT TOOLS
			tools.CreateDatasourceTool(),
			tools.UpdateDatasourceTool(),
			tools.DeleteDatasourceTool(),
		}
	})
}

// GetHandlers returns all tool handlers mapped to their respective tool names.
// Handlers are only returned if the service is enabled and properly initialized.
// Organized with same order as GetTools for consistency.
func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
	if !s.enabled {
		return nil
	}

	return map[string]server.ToolHandlerFunc{
		// 🎯 SUMMARY TOOLS - Recommended first (95% of use cases)
		"grafana_dashboards_summary":  handlers.HandleGetDashboardsSummary(),
		"grafana_datasources_summary": handlers.HandleGetDataSourcesSummary(),
		"grafana_plugins_summary":     handlers.HandleGetPluginsSummary(),

		// 📋 STANDARD TOOLS - For detailed information when needed
		"grafana_dashboards":        handlers.HandleGetDashboards(),
		"grafana_datasources":       handlers.HandleGetDataSources(),
		"grafana_plugins":           handlers.HandleGetPlugins(),
		"grafana_folders":           handlers.HandleGetFolders(),
		"grafana_create_folder":     handlers.HandleCreateFolder(),
		"grafana_update_folder":     handlers.HandleUpdateFolder(),
		"grafana_delete_folder":     handlers.HandleDeleteFolder(),
		"grafana_search_dashboards": handlers.HandleSearchDashboards(),

		// 🔍 SPECIFIC RESOURCE TOOLS - For detailed inspection
		"grafana_dashboard":              handlers.HandleGetDashboard(),
		"grafana_folder_detail":          handlers.HandleGetFolder(),
		"grafana_datasource_detail":      handlers.HandleGetDataSource(),
		"grafana_get_datasource_by_name": handlers.HandleGetDataSourceByName(),
		"grafana_plugin_detail":          handlers.HandleGetPlugin(),

		// 🚨 MONITORING TOOLS
		"grafana_alerts":                handlers.HandleGetAlertRules(),
		"grafana_get_alert_rule_by_uid": handlers.HandleGetAlertRuleByUID(),
		"grafana_list_contact_points":   handlers.HandleListContactPoints(),

		// 📝 ANNOTATION TOOLS
		"grafana_get_annotations":     handlers.HandleGetAnnotations(),
		"grafana_create_annotation":   handlers.HandleCreateAnnotation(),
		"grafana_update_annotation":   handlers.HandleUpdateAnnotation(),
		"grafana_patch_annotation":    handlers.HandlePatchAnnotation(),
		"grafana_delete_annotation":   handlers.HandleDeleteAnnotation(),
		"grafana_get_annotation_tags": handlers.HandleGetAnnotationTags(),

		// 🔧 UTILITY TOOLS
		"grafana_test_connection":         handlers.HandleTestConnection(),
		"grafana_current_user":            handlers.HandleGetCurrentUser(),
		"grafana_users":                   handlers.HandleGetUsers(),
		"grafana_organization":            handlers.HandleGetOrganization(),
		"grafana_check_datasource_health": handlers.HandleCheckDatasourceHealth(),

		// 🔐 ADMIN TOOLS
		"grafana_list_teams":               handlers.HandleListTeams(),
		"grafana_list_all_roles":           handlers.HandleListAllRoles(),
		"grafana_get_role_details":         handlers.HandleGetRoleDetails(),
		"grafana_get_role_assignments":     handlers.HandleGetRoleAssignments(),
		"grafana_list_user_roles":          handlers.HandleListUserRoles(),
		"grafana_list_team_roles":          handlers.HandleListTeamRoles(),
		"grafana_get_resource_permissions": handlers.HandleGetResourcePermissions(),
		"grafana_get_resource_description": handlers.HandleGetResourceDescription(),

		// 📊 DASHBOARD MANAGEMENT TOOLS
		"grafana_update_dashboard":            handlers.HandleUpdateDashboard(),
		"grafana_get_dashboard_versions":      handlers.HandleGetDashboardVersions(),
		"grafana_get_dashboard_version":       handlers.HandleGetDashboardVersion(),
		"grafana_restore_dashboard_version":   handlers.HandleRestoreDashboardVersion(),
		"grafana_delete_dashboard":            handlers.HandleDeleteDashboard(),
		"grafana_get_dashboard_panel_queries": handlers.HandleGetDashboardPanelQueries(),
		"grafana_get_dashboard_property":      handlers.HandleGetDashboardProperty(),

		// 🚨 ALERTING TOOLS
		"grafana_create_alert_rule": handlers.HandleCreateAlertRule(),
		"grafana_update_alert_rule": handlers.HandleUpdateAlertRule(),
		"grafana_delete_alert_rule": handlers.HandleDeleteAlertRule(),

		// 🔗 NAVIGATION TOOLS
		"grafana_generate_deeplink":            handlers.HandleGenerateDeeplink(),
		"grafana_generate_logs_drilldown_link": handlers.HandleGenerateLogsDrilldownLink(),

		// 🎨 RENDERING TOOLS
		"grafana_render_panel_image": handlers.HandleRenderPanelImage(),

		// 📝 GRAPHITE ANNOTATION TOOLS
		"grafana_create_graphite_annotation": handlers.HandleCreateGraphiteAnnotation(),

		// 🔧 DATASOURCE MANAGEMENT TOOLS
		"grafana_create_datasource": handlers.HandleCreateDatasource(),
		"grafana_update_datasource": handlers.HandleUpdateDatasource(),
		"grafana_delete_datasource": handlers.HandleDeleteDatasource(),
	}
}

// IsEnabled returns whether the service is enabled and ready for use.
// A service is considered enabled if it's marked as enabled and has a valid client.
func (s *Service) IsEnabled() bool {
	return s.enabled
}

// GetClient returns the underlying Grafana client for advanced operations.
// This method is primarily used for testing and internal service communication.
func (s *Service) GetClient() *client.Client {
	return nil // Backend client is created per-request from HTTP headers
}

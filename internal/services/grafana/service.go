// Package grafana provides Grafana dashboard and monitoring integration for the MCP server.
// It implements tools for managing Grafana dashboards, data sources, folders, and alert rules.
package grafana

import (
	"time"

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
type Service struct {
	client        *client.Client               // Grafana HTTP client for API operations
	enabled       bool                         // Whether the service is enabled
	toolsCache    *cache.ToolsCache            // Cached tools to avoid recreation
	initFramework *framework.CommonServiceInit // Common initialization framework
}

// NewService creates a new Grafana service instance.
// The service is disabled by default and requires initialization before use.
func NewService() *Service {
	// Create service enable checker
	checker := framework.NewServiceEnabled(
		func(cfg *config.AppConfig) bool { return cfg.Grafana.Enabled },
		func(cfg *config.AppConfig) string { return cfg.Grafana.URL },
	)

	// Create init configuration
	initConfig := &framework.InitConfig{
		Required:     false,
		URLValidator: framework.SimpleURLValidator,
		ClientBuilder: func(cfg *config.AppConfig) (interface{}, error) {
			return client.NewClient(&client.ClientOptions{
				URL:      cfg.Grafana.URL,
				APIKey:   cfg.Grafana.APIKey,
				Username: cfg.Grafana.Username,
				Password: cfg.Grafana.Password,
				Timeout:  time.Duration(cfg.Grafana.TimeoutSec) * time.Second,
			})
		},
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
		func(clientIface interface{}) {
			if grafanaClient, ok := clientIface.(*client.Client); ok {
				s.client = grafanaClient
			}
		},
	)
}

// GetTools returns all available Grafana MCP tools with optimization focus.
// Tools are only returned if the service is enabled and properly initialized.
// Organized with summary tools first for better user experience.
func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled || s.client == nil {
		return nil
	}

	// Use unified cache
	return s.toolsCache.Get(func() []mcp.Tool {
		return []mcp.Tool{
			// 🎯 SUMMARY TOOLS - Recommended first (95% of use cases)
			tools.GetDashboardsSummaryTool(),  // RECOMMENDED for dashboard discovery
			tools.GetDataSourcesSummaryTool(), // RECOMMENDED for datasource discovery

			// 📋 STANDARD TOOLS - For detailed information when needed
			tools.GetDashboardsTool(),
			tools.GetDataSourcesTool(),
			tools.GetFoldersTool(),
			tools.SearchDashboardsTool(),

			// 🔍 SPECIFIC RESOURCE TOOLS - For detailed inspection
			tools.GetDashboardTool(),
			tools.GetFolderTool(),
			tools.GetDataSourceTool(),
			tools.GetDataSourceByNameTool(),

			// 🚨 MONITORING TOOLS
			tools.GetAlertRulesTool(),
			tools.GetAlertRuleByUIDTool(),
			tools.ListContactPointsTool(),

			// 📝 ANNOTATION TOOLS
			tools.GetAnnotationsTool(),
			tools.CreateAnnotationTool(),
			tools.UpdateAnnotationTool(),
			tools.PatchAnnotationTool(),
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
			tools.GetDashboardPanelQueriesTool(),
			tools.GetDashboardPropertyTool(),

			// 🚨 ALERTING TOOLS
			tools.CreateAlertRuleTool(),
			tools.UpdateAlertRuleTool(),
			tools.DeleteAlertRuleTool(),

			// 🔗 NAVIGATION TOOLS
			tools.GenerateDeeplinkTool(),

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
	if !s.enabled || s.client == nil {
		return nil
	}

	return map[string]server.ToolHandlerFunc{
		// 🎯 SUMMARY TOOLS - Recommended first (95% of use cases)
		"grafana_dashboards_summary":  handlers.HandleGetDashboardsSummary(s.client),
		"grafana_datasources_summary": handlers.HandleGetDataSourcesSummary(s.client),

		// 📋 STANDARD TOOLS - For detailed information when needed
		"grafana_dashboards":        handlers.HandleGetDashboards(s.client),
		"grafana_datasources":       handlers.HandleGetDataSources(s.client),
		"grafana_folders":           handlers.HandleGetFolders(s.client),
		"grafana_search_dashboards": handlers.HandleSearchDashboards(s.client),

		// 🔍 SPECIFIC RESOURCE TOOLS - For detailed inspection
		"grafana_dashboard":              handlers.HandleGetDashboard(s.client),
		"grafana_folder_detail":          handlers.HandleGetFolder(s.client),
		"grafana_datasource_detail":      handlers.HandleGetDataSource(s.client),
		"grafana_get_datasource_by_name": handlers.HandleGetDataSourceByName(s.client),

		// 🚨 MONITORING TOOLS
		"grafana_alerts":                handlers.HandleGetAlertRules(s.client),
		"grafana_get_alert_rule_by_uid": handlers.HandleGetAlertRuleByUID(s.client),
		"grafana_list_contact_points":   handlers.HandleListContactPoints(s.client),

		// 📝 ANNOTATION TOOLS
		"grafana_get_annotations":     handlers.HandleGetAnnotations(s.client),
		"grafana_create_annotation":   handlers.HandleCreateAnnotation(s.client),
		"grafana_update_annotation":   handlers.HandleUpdateAnnotation(s.client),
		"grafana_patch_annotation":    handlers.HandlePatchAnnotation(s.client),
		"grafana_get_annotation_tags": handlers.HandleGetAnnotationTags(s.client),

		// 🔧 UTILITY TOOLS
		"grafana_test_connection":         handlers.HandleTestConnection(s.client),
		"grafana_current_user":            handlers.HandleGetCurrentUser(s.client),
		"grafana_users":                   handlers.HandleGetUsers(s.client),
		"grafana_organization":            handlers.HandleGetOrganization(s.client),
		"grafana_check_datasource_health": handlers.HandleCheckDatasourceHealth(s.client),

		// 🔐 ADMIN TOOLS
		"grafana_list_teams":               handlers.HandleListTeams(s.client),
		"grafana_list_all_roles":           handlers.HandleListAllRoles(s.client),
		"grafana_get_role_details":         handlers.HandleGetRoleDetails(s.client),
		"grafana_get_role_assignments":     handlers.HandleGetRoleAssignments(s.client),
		"grafana_list_user_roles":          handlers.HandleListUserRoles(s.client),
		"grafana_list_team_roles":          handlers.HandleListTeamRoles(s.client),
		"grafana_get_resource_permissions": handlers.HandleGetResourcePermissions(s.client),
		"grafana_get_resource_description": handlers.HandleGetResourceDescription(s.client),

		// 📊 DASHBOARD MANAGEMENT TOOLS
		"grafana_update_dashboard":            handlers.HandleUpdateDashboard(s.client),
		"grafana_get_dashboard_panel_queries": handlers.HandleGetDashboardPanelQueries(s.client),
		"grafana_get_dashboard_property":      handlers.HandleGetDashboardProperty(s.client),

		// 🚨 ALERTING TOOLS
		"grafana_create_alert_rule": handlers.HandleCreateAlertRule(s.client),
		"grafana_update_alert_rule": handlers.HandleUpdateAlertRule(s.client),
		"grafana_delete_alert_rule": handlers.HandleDeleteAlertRule(s.client),

		// 🔗 NAVIGATION TOOLS
		"grafana_generate_deeplink": handlers.HandleGenerateDeeplink(s.client),

		// 🎨 RENDERING TOOLS
		"grafana_render_panel_image": handlers.HandleRenderPanelImage(s.client),

		// 📝 GRAPHITE ANNOTATION TOOLS
		"grafana_create_graphite_annotation": handlers.HandleCreateGraphiteAnnotation(s.client),

		// 🔧 DATASOURCE MANAGEMENT TOOLS
		"grafana_create_datasource": handlers.HandleCreateDatasource(s.client),
		"grafana_update_datasource": handlers.HandleUpdateDatasource(s.client),
		"grafana_delete_datasource": handlers.HandleDeleteDatasource(s.client),
	}
}

// IsEnabled returns whether the service is enabled and ready for use.
// A service is considered enabled if it's marked as enabled and has a valid client.
func (s *Service) IsEnabled() bool {
	return s.enabled && s.client != nil
}

// GetClient returns the underlying Grafana client for advanced operations.
// This method is primarily used for testing and internal service communication.
func (s *Service) GetClient() *client.Client {
	return s.client
}

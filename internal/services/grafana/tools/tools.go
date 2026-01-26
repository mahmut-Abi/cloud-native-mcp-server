// Package tools provides Grafana MCP tools for dashboard and data source management.
// These tools allow interaction with Grafana instances through the MCP protocol.
// NOTE: For large clusters, consider using a smaller limit parameter to avoid overwhelming outputs
package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
)

// GetDashboardsTool retrieves all dashboards from Grafana with intelligent limits.
func GetDashboardsTool() mcp.Tool {
	logrus.Debug("Creating GetDashboardsTool")
	return mcp.NewTool("grafana_dashboards",
		mcp.WithDescription("📋 List all dashboards with intelligent size limits and metadata. This tool returns dashboard summaries (id, uid, title, tags, folder information) but automatically removes heavy dashboard configurations to prevent context overflow. Use this tool when you need to: discover available dashboards, understand dashboard organization, find dashboards by browsing, or get dashboard metadata. The response includes pagination info and warnings if data was truncated. For complete dashboard data, use 'grafana_dashboard' with specific UID."),
		mcp.WithString("limit",
			mcp.Description("Maximum dashboards to return (default: 20, max: 100). Use lower values for faster responses and to prevent context overflow.")),
		mcp.WithString("debug",
			mcp.Description("Enable verbose debug output for troubleshooting the API call. Set to 'true' to see detailed information about the Grafana API request, response processing, size optimization, and any errors encountered.")),
	)
}

// GetDashboardsSummaryTool retrieves all dashboards with minimal output (RECOMMENDED)
func GetDashboardsSummaryTool() mcp.Tool {
	logrus.Debug("Creating GetDashboardsSummaryTool")
	return mcp.NewTool("grafana_dashboards_summary",
		mcp.WithDescription("🎯 RECOMMENDED: List all dashboards with minimal output for efficient discovery. Returns only essential fields (id, uid, title, folderUID, tags) with 70-85% smaller response size than detailed version. Perfect for: quick dashboard browsing, finding dashboards by name, understanding dashboard landscape, or as first step before getting detailed info. Includes pagination for large collections."),
		mcp.WithString("limit",
			mcp.Description("Maximum dashboards to return (default: 20, max: 100). Lower values recommended for better performance.")),
		mcp.WithString("offset",
			mcp.Description("Pagination offset for browsing through large dashboard collections (default: 0)")),
	)
}

// GetDataSourcesSummaryTool retrieves all data sources with minimal output (RECOMMENDED)
func GetDataSourcesSummaryTool() mcp.Tool {
	logrus.Debug("Creating GetDataSourcesSummaryTool")
	return mcp.NewTool("grafana_datasources_summary",
		mcp.WithDescription("🎯 RECOMMENDED: List all data sources with minimal output for quick discovery. Returns only essential fields (id, uid, name, type, isDefault) with 70-85% smaller response size. Perfect for: understanding available data sources, checking data source types, identifying default data source, or as first step before getting detailed configuration. Automatically limits to 15 items to prevent context overflow."),
		mcp.WithString("limit",
			mcp.Description("Maximum data sources to return (default: 15, max: 50). Data sources contain configurations, so conservative limits are applied.")),
	)
}

// GetDashboardTool retrieves a specific dashboard from Grafana.
func GetDashboardTool() mcp.Tool {
	logrus.Debug("Creating GetDashboardTool")
	return mcp.NewTool("grafana_dashboard",
		mcp.WithDescription("Retrieve a specific dashboard from Grafana by its unique identifier (UID). This tool returns the complete dashboard configuration including all panels, queries, variables, and settings. Use this tool when you need to: examine the detailed configuration of a specific dashboard, export a dashboard for backup or migration, analyze dashboard structure and panel configurations, debug dashboard issues or performance problems, or prepare for dashboard modifications. The response includes the full dashboard JSON along with metadata such as version, permissions, and folder information. This is essential for dashboard management, troubleshooting, and configuration analysis."),
		mcp.WithString("uid", mcp.Required(),
			mcp.Description("Unique identifier (UID) of the dashboard to retrieve. The UID is a string that uniquely identifies the dashboard in Grafana and remains constant even if the dashboard is moved between folders or renamed. You can find the UID in the dashboard URL (e.g., /d/dashboard-uid/dashboard-name) or by using the get_grafana_dashboards tool to list all dashboards with their UIDs. The UID is case-sensitive and must be provided exactly as it appears in Grafana. Common UID formats include short alphanumeric strings like 'abc123def' or longer identifiers with hyphens.")),
		mcp.WithString("debug",
			mcp.Description("Enable comprehensive debug output for troubleshooting dashboard retrieval. Set to 'true' to see detailed information about the API request, response parsing, and any errors during the dashboard fetch process. Set to 'false' or omit for normal output showing only the dashboard configuration. Debug mode is helpful when dashboards cannot be found, when there are authentication issues, or when you need to understand the API interaction for automation purposes.")),
	)
}

// GetDataSourcesTool retrieves all data sources from Grafana with smart filtering.
func GetDataSourcesTool() mcp.Tool {
	logrus.Debug("CreatingGetDataSourcesTool")
	return mcp.NewTool("grafana_datasources",
		mcp.WithDescription("📋 List all data sources with intelligent filtering and security. Returns data source summaries (id, uid, name, type, URL, access) but automatically removes sensitive data (passwords) and large configurations (jsonData) to prevent context overflow and protect security. Use this tool when you need to: inventory data sources, understand data source types, verify connectivity, or get basic configuration. For complete config including sensitive data, use 'grafana_datasource_detail' with specific UID."),
		mcp.WithString("limit",
			mcp.Description("Maximum data sources to return (default: 15, max: 100). Conservative limits applied as datasources contain configuration data.")),
		mcp.WithString("debug",
			mcp.Description("Enable detailed debug output for troubleshooting. Sensitive data will still be filtered in output.")),
	)
}

// GetFoldersTool retrieves all folders from Grafana with limits.
func GetFoldersTool() mcp.Tool {
	logrus.Debug("Creating GetFoldersTool")
	return mcp.NewTool("grafana_folders",
		mcp.WithDescription("📁 List all folders from Grafana with intelligent limits. Returns folder metadata (id, uid, title, permissions) to understand organizational structure. Folders organize dashboards and control access permissions. Use this tool to: understand Grafana organization, find appropriate folders for dashboards, audit permissions, or navigate hierarchy. Response includes pagination info for large folder collections."),
		mcp.WithString("limit",
			mcp.Description("Maximum folders to return (default: 20, max: 100). Folders are typically small but limits still apply.")),
		mcp.WithString("debug",
			mcp.Description("Enable verbose debug output for troubleshooting folder retrieval and API interaction.")),
	)
}

// GetAlertRulesTool retrieves alert rules from Grafana with limits.
func GetAlertRulesTool() mcp.Tool {
	logrus.Debug("Creating GetAlertRulesTool")
	return mcp.NewTool("grafana_alerts",
		mcp.WithDescription("🚨 List alert rules and alerting configuration with intelligent limits. Returns alert rules with conditions, thresholds, and notification settings. Use this tool to: monitor alerting health, audit configurations, troubleshoot notification issues, review thresholds, understand active alerts, or analyze alerting patterns. Response includes pagination metadata and size warnings to prevent context overflow with complex alert configurations."),
		mcp.WithString("limit",
			mcp.Description("Maximum alert rules to return (default: 20, max: 100). Alert rules can be complex, so conservative limits are applied.")),
		mcp.WithString("debug",
			mcp.Description("Enable comprehensive debug output for troubleshooting alert rule retrieval and alerting system analysis.")),
	)
}

// TestConnectionTool tests the connection to Grafana.
func TestConnectionTool() mcp.Tool {
	logrus.Debug("Creating TestConnectionTool")
	return mcp.NewTool("grafana_test_connection",
		mcp.WithDescription("Test the connection and authentication to a Grafana instance. This tool verifies that the Grafana server is accessible, authentication credentials are valid, and the API is responding correctly. Use this tool when you need to: verify Grafana server connectivity and availability, validate API credentials and authentication setup, troubleshoot connection issues before performing other operations, check if the Grafana service is healthy and responsive, or diagnose network or configuration problems. This is typically the first tool to use when setting up Grafana integration or when experiencing connectivity issues. The tool performs a simple health check without retrieving sensitive data, making it safe for initial connection validation."),
		mcp.WithString("debug",
			mcp.Description("Enable detailed debug output for connection testing and troubleshooting. Set to 'true' to see comprehensive information about the connection attempt, HTTP request details, authentication process, response analysis, and any network or API errors encountered. Set to 'false' or omit for normal output showing only the connection test result. Debug mode is invaluable when diagnosing connectivity issues, authentication problems, network configuration issues, or SSL/TLS certificate problems.")),
	)
}

// SearchDashboardsTool searches for dashboards in Grafana with intelligent limits.
func SearchDashboardsTool() mcp.Tool {
	logrus.Debug("Creating SearchDashboardsTool")
	return mcp.NewTool("grafana_search_dashboards",
		mcp.WithDescription("🔍 Search dashboards using criteria (title, tags, folder) with intelligent size limits. Allows finding specific dashboards without retrieving entire list. Returns dashboard summaries (id, uid, title, tags) with search criteria metadata. Perfect for: finding dashboards by name pattern, locating dashboards with specific tags, searching within folders, filtering by starred status, or efficiently discovering Dashboards in large instances."),
		mcp.WithString("query",
			mcp.Description("Search query string to match dashboard titles (case-insensitive partial match). E.g., 'cpu' finds 'CPU Usage', 'Server CPU'. Leave empty for all.")),
		mcp.WithString("tag",
			mcp.Description("Filter dashboards by specific tag (exact match). E.g., 'production', 'kubernetes', 'backend-team'. Leave empty to ignore.")),
		mcp.WithString("folderUID",
			mcp.Description("Filter dashboards within specific folder by UID. Use grafana_folders to get UIDs. Leave empty to search all folders.")),
		mcp.WithBoolean("starred",
			mcp.Description("Filter starred dashboards (true = only starred, false = only non-starred). Leave unset for both.")),
		mcp.WithString("limit",
			mcp.Description("Maximum results to return (default: 20, max: 100). Search efficiency improves with specific criteria.")),
		mcp.WithString("debug",
			mcp.Description("Enable detailed debug output for troubleshooting search operations and filtering logic.")),
	)
}

// GetFolderTool retrieves a specific folder by UID.
func GetFolderTool() mcp.Tool {
	logrus.Debug("Creating GetFolderTool")
	return mcp.NewTool("grafana_folder_detail",
		mcp.WithDescription("Retrieve detailed information about a specific Grafana folder by UID, including permissions and ACL settings."),
		mcp.WithString("uid", mcp.Required(),
			mcp.Description("Unique identifier (UID) of the folder to retrieve.")),
		mcp.WithString("debug",
			mcp.Description("Enable debug output for troubleshooting.")),
	)
}

// GetDataSourceTool retrieves a specific datasource by UID.
func GetDataSourceTool() mcp.Tool {
	logrus.Debug("Creating GetDataSourceTool")
	return mcp.NewTool("grafana_datasource_detail",
		mcp.WithDescription("Retrieve detailed configuration of a specific Grafana data source by UID, including connection parameters and settings."),
		mcp.WithString("uid", mcp.Required(),
			mcp.Description("Unique identifier (UID) of the data source to retrieve.")),
		mcp.WithString("debug",
			mcp.Description("Enable debug output for troubleshooting.")),
	)
}

// CheckDatasourceHealthTool checks the health of a datasource.
func CheckDatasourceHealthTool() mcp.Tool {
	logrus.Debug("Creating CheckDatasourceHealthTool")
	return mcp.NewTool("grafana_check_datasource_health",
		mcp.WithDescription("Test the connectivity and health status of a specific Grafana data source. Returns status and diagnostic messages."),
		mcp.WithString("uid", mcp.Required(),
			mcp.Description("Unique identifier (UID) of the data source to check.")),
		mcp.WithString("debug",
			mcp.Description("Enable debug output for troubleshooting.")),
	)
}

// GetCurrentUserTool retrieves current authenticated user.
func GetCurrentUserTool() mcp.Tool {
	logrus.Debug("Creating GetCurrentUserTool")
	return mcp.NewTool("grafana_current_user",
		mcp.WithDescription("Retrieve information about the currently authenticated Grafana user, including username, email, and permissions."),
		mcp.WithString("debug",
			mcp.Description("Enable debug output for troubleshooting.")),
	)
}

// GetUsersTool retrieves all users.
func GetUsersTool() mcp.Tool {
	logrus.Debug("Creating GetUsersTool")
	return mcp.NewTool("grafana_users",
		mcp.WithDescription("Retrieve all users in the Grafana organization. Returns user list with information including username, email, admin status, and last activity."),
		mcp.WithString("debug",
			mcp.Description("Enable debug output for troubleshooting.")),
	)
}

// GetOrganizationTool retrieves organization information.
func GetOrganizationTool() mcp.Tool {
	logrus.Debug("Creating GetOrganizationTool")
	return mcp.NewTool("grafana_organization",
		mcp.WithDescription("Retrieve current Grafana organization information including name and ID."),
		mcp.WithString("debug",
			mcp.Description("Enable debug output for troubleshooting.")),
	)
}

// ============ Admin Tools ============

// ListTeamsTool lists all teams in the organization.
func ListTeamsTool() mcp.Tool {
	logrus.Debug("Creating ListTeamsTool")
	return mcp.NewTool("grafana_list_teams",
		mcp.WithDescription("List all teams in the Grafana organization. Returns team information including name, email, and member count."),
	)
}

// ListAllRolesTool lists all available roles.
func ListAllRolesTool() mcp.Tool {
	logrus.Debug("Creating ListAllRolesTool")
	return mcp.NewTool("grafana_list_all_roles",
		mcp.WithDescription("List all available Grafana roles with their permissions. Requires admin permissions."),
	)
}

// GetRoleDetailsTool retrieves details of a specific role.
func GetRoleDetailsTool() mcp.Tool {
	logrus.Debug("Creating GetRoleDetailsTool")
	return mcp.NewTool("grafana_get_role_details",
		mcp.WithDescription("Retrieve detailed information about a specific Grafana role, including its permissions."),
		mcp.WithString("roleUID", mcp.Required(),
			mcp.Description("Unique identifier of the role to retrieve.")),
	)
}

// GetRoleAssignmentsTool retrieves assignments for a specific role.
func GetRoleAssignmentsTool() mcp.Tool {
	logrus.Debug("Creating GetRoleAssignmentsTool")
	return mcp.NewTool("grafana_get_role_assignments",
		mcp.WithDescription("List all users and teams assigned to a specific role."),
		mcp.WithString("roleUID", mcp.Required(),
			mcp.Description("Unique identifier of the role.")),
	)
}

// ListUserRolesTool lists roles for a specific user.
func ListUserRolesTool() mcp.Tool {
	logrus.Debug("Creating ListUserRolesTool")
	return mcp.NewTool("grafana_list_user_roles",
		mcp.WithDescription("List all roles assigned to a specific user."),
		mcp.WithString("userID", mcp.Required(),
			mcp.Description("User ID to query roles for.")),
	)
}

// ListTeamRolesTool lists roles for a specific team.
func ListTeamRolesTool() mcp.Tool {
	logrus.Debug("Creating ListTeamRolesTool")
	return mcp.NewTool("grafana_list_team_roles",
		mcp.WithDescription("List all roles assigned to a specific team."),
		mcp.WithString("teamID", mcp.Required(),
			mcp.Description("Team ID to query roles for.")),
	)
}

// GetResourcePermissionsTool retrieves permissions for a resource.
func GetResourcePermissionsTool() mcp.Tool {
	logrus.Debug("Creating GetResourcePermissionsTool")
	return mcp.NewTool("grafana_get_resource_permissions",
		mcp.WithDescription("List all permissions defined for a specific resource."),
		mcp.WithString("resourceType", mcp.Required(),
			mcp.Description("Resource type (e.g., dashboards, datasources, folders).")),
		mcp.WithString("resourceUID",
			mcp.Description("Unique identifier of the specific resource. Leave empty for all resources of this type.")),
	)
}

// GetResourceDescriptionTool describes a Grafana resource type.
func GetResourceDescriptionTool() mcp.Tool {
	logrus.Debug("Creating GetResourceDescriptionTool")
	return mcp.NewTool("grafana_get_resource_description",
		mcp.WithDescription("Describe a Grafana resource type, including available permissions and assignment capabilities."),
		mcp.WithString("resourceType", mcp.Required(),
			mcp.Description("Resource type to describe (e.g., dashboards, datasources, folders).")),
	)
}

// ============ Dashboard Update Tools ============

// UpdateDashboardTool creates or updates a dashboard.
func UpdateDashboardTool() mcp.Tool {
	logrus.Debug("Creating UpdateDashboardTool")
	return mcp.NewTool("grafana_update_dashboard",
		mcp.WithDescription("Create a new dashboard or update an existing one. Requires dashboard write permissions."),
		mcp.WithAny("dashboard", mcp.Required(),
			mcp.Description("Dashboard JSON model to create or update.")),
		mcp.WithString("folderUID",
			mcp.Description("UID of the folder to save the dashboard in. Defaults to General folder.")),
		mcp.WithBoolean("overwrite",
			mcp.Description("Overwrite existing dashboard with the same name/UID. Default: false.")),
		mcp.WithString("message",
			mcp.Description("Commit message for the dashboard change.")),
	)
}

// GetDashboardPanelQueriesTool retrieves panel queries and datasource info.
func GetDashboardPanelQueriesTool() mcp.Tool {
	logrus.Debug("Creating GetDashboardPanelQueriesTool")
	return mcp.NewTool("grafana_get_dashboard_panel_queries",
		mcp.WithDescription("Get the title, query string, and datasource information from every panel in a dashboard."),
		mcp.WithString("dashboardUID", mcp.Required(),
			mcp.Description("Unique identifier of the dashboard.")),
	)
}

// GetDashboardPropertyTool extracts specific parts of a dashboard using JSONPath.
func GetDashboardPropertyTool() mcp.Tool {
	logrus.Debug("Creating GetDashboardPropertyTool")
	return mcp.NewTool("grafana_get_dashboard_property",
		mcp.WithDescription("Extract specific parts of a dashboard using JSONPath expressions. Reduces context window usage."),
		mcp.WithString("dashboardUID", mcp.Required(),
			mcp.Description("Unique identifier of the dashboard.")),
		mcp.WithString("propertyPath", mcp.Required(),
			mcp.Description("JSONPath-like expression (e.g., 'title', 'panels[0].title', 'tags').")),
	)
}

// ============ Alerting Tools ============

// GetAlertRuleByUIDTool retrieves a specific alert rule.
func GetAlertRuleByUIDTool() mcp.Tool {
	logrus.Debug("Creating GetAlertRuleByUIDTool")
	return mcp.NewTool("grafana_get_alert_rule_by_uid",
		mcp.WithDescription("Retrieve details of a specific alert rule by its UID."),
		mcp.WithString("ruleUID", mcp.Required(),
			mcp.Description("Unique identifier of the alert rule.")),
	)
}

// CreateAlertRuleTool creates a new alert rule.
func CreateAlertRuleTool() mcp.Tool {
	logrus.Debug("Creating CreateAlertRuleTool")
	return mcp.NewTool("grafana_create_alert_rule",
		mcp.WithDescription("Create a new alert rule in Grafana. Requires alert rule write permissions."),
		mcp.WithString("title", mcp.Required(),
			mcp.Description("Name of the alert rule.")),
		mcp.WithString("condition", mcp.Required(),
			mcp.Description("Alert condition (e.g., 'F' for firing).")),
		mcp.WithString("folderUID", mcp.Required(),
			mcp.Description("UID of the folder to store the alert rule.")),
		mcp.WithString("ruleGroup",
			mcp.Description("Name of the rule group. Default: 'default'.")),
		mcp.WithNumber("intervalSeconds",
			mcp.Description("Evaluation interval in seconds. Default: 60.")),
		mcp.WithAny("data",
			mcp.Description("Query data for the alert rule.")),
		mcp.WithObject("annotations",
			mcp.Description("Annotations to add to alerts.")),
		mcp.WithObject("labels",
			mcp.Description("Labels to add to alerts.")),
	)
}

// UpdateAlertRuleTool updates an existing alert rule.
func UpdateAlertRuleTool() mcp.Tool {
	logrus.Debug("Creating UpdateAlertRuleTool")
	return mcp.NewTool("grafana_update_alert_rule",
		mcp.WithDescription("Update an existing alert rule. Requires alert rule write permissions."),
		mcp.WithString("ruleUID", mcp.Required(),
			mcp.Description("Unique identifier of the alert rule to update.")),
		mcp.WithString("title",
			mcp.Description("New name for the alert rule.")),
		mcp.WithString("condition",
			mcp.Description("New alert condition.")),
		mcp.WithString("folderUID",
			mcp.Description("New folder UID for the alert rule.")),
	)
}

// DeleteAlertRuleTool deletes an alert rule.
func DeleteAlertRuleTool() mcp.Tool {
	logrus.Debug("Creating DeleteAlertRuleTool")
	return mcp.NewTool("grafana_delete_alert_rule",
		mcp.WithDescription("Delete an alert rule by its UID. Requires alert rule write permissions."),
		mcp.WithString("ruleUID", mcp.Required(),
			mcp.Description("Unique identifier of the alert rule to delete.")),
	)
}

// ListContactPointsTool lists notification contact points.
func ListContactPointsTool() mcp.Tool {
	logrus.Debug("Creating ListContactPointsTool")
	return mcp.NewTool("grafana_list_contact_points",
		mcp.WithDescription("List all notification contact points (Grafana-managed and Alertmanager)."),
	)
}

// ============ Annotation Tools ============

// GetAnnotationsTool retrieves annotations with filters.
func GetAnnotationsTool() mcp.Tool {
	logrus.Debug("Creating GetAnnotationsTool")
	return mcp.NewTool("grafana_get_annotations",
		mcp.WithDescription("Query annotations with optional filters. Supports time range, dashboard UID, tags, and match mode."),
		mcp.WithString("dashboardUID",
			mcp.Description("Filter annotations by dashboard UID.")),
		mcp.WithString("from",
			mcp.Description("Start time in RFC3339 format or Unix timestamp.")),
		mcp.WithString("to",
			mcp.Description("End time in RFC3339 format or Unix timestamp.")),
		mcp.WithString("tags",
			mcp.Description("Filter by tags (comma-separated).")),
		mcp.WithString("limit",
			mcp.Description("Maximum number of annotations to return.")),
	)
}

// CreateAnnotationTool creates a new annotation.
func CreateAnnotationTool() mcp.Tool {
	logrus.Debug("Creating CreateAnnotationTool")
	return mcp.NewTool("grafana_create_annotation",
		mcp.WithDescription("Create a new annotation on a dashboard or panel."),
		mcp.WithString("text", mcp.Required(),
			mcp.Description("Annotation text.")),
		mcp.WithNumber("time",
			mcp.Description("Timestamp in milliseconds. Default: now.")),
		mcp.WithNumber("timeEnd",
			mcp.Description("End timestamp in milliseconds for range annotations.")),
		mcp.WithString("dashboardUID",
			mcp.Description("Dashboard UID to annotate.")),
		mcp.WithNumber("panelID",
			mcp.Description("Panel ID to annotate.")),
		mcp.WithArray("tags",
			mcp.Description("Tags for the annotation.")),
	)
}

// UpdateAnnotationTool updates an annotation.
func UpdateAnnotationTool() mcp.Tool {
	logrus.Debug("Creating UpdateAnnotationTool")
	return mcp.NewTool("grafana_update_annotation",
		mcp.WithDescription("Replace all fields of an existing annotation (full update)."),
		mcp.WithString("annotationID", mcp.Required(),
			mcp.Description("ID of the annotation to update.")),
		mcp.WithString("text",
			mcp.Description("New annotation text.")),
		mcp.WithNumber("time",
			mcp.Description("New timestamp in milliseconds.")),
		mcp.WithNumber("timeEnd",
			mcp.Description("New end timestamp in milliseconds.")),
	)
}

// PatchAnnotationTool partially updates an annotation.
func PatchAnnotationTool() mcp.Tool {
	logrus.Debug("Creating PatchAnnotationTool")
	return mcp.NewTool("grafana_patch_annotation",
		mcp.WithDescription("Update only specific fields of an annotation (partial update)."),
		mcp.WithString("annotationID", mcp.Required(),
			mcp.Description("ID of the annotation to patch.")),
		mcp.WithString("text",
			mcp.Description("Text to update.")),
		mcp.WithNumber("time",
			mcp.Description("Timestamp to update.")),
	)
}

// GetAnnotationTagsTool lists available annotation tags.
func GetAnnotationTagsTool() mcp.Tool {
	logrus.Debug("Creating GetAnnotationTagsTool")
	return mcp.NewTool("grafana_get_annotation_tags",
		mcp.WithDescription("List available annotation tags with optional filtering."),
		mcp.WithString("tag",
			mcp.Description("Filter by specific tag.")),
	)
}

// ============ Navigation Tools ============

// GenerateDeeplinkTool generates accurate deeplink URLs for Grafana resources.
func GenerateDeeplinkTool() mcp.Tool {
	logrus.Debug("Creating GenerateDeeplinkTool")
	return mcp.NewTool("grafana_generate_deeplink",
		mcp.WithDescription("Generate accurate deeplink URLs for Grafana resources instead of relying on URL guessing."),
		mcp.WithString("resourceType", mcp.Required(),
			mcp.Description("Resource type: dashboard, panel, or explore.")),
		mcp.WithString("resourceUID", mcp.Required(),
			mcp.Description("Unique identifier of the resource.")),
		mcp.WithString("panelID",
			mcp.Description("Panel ID for panel links.")),
		mcp.WithString("datasource",
			mcp.Description("Datasource UID for explore links.")),
		mcp.WithString("from",
			mcp.Description("Start time (e.g., now-1h).")),
		mcp.WithString("to",
			mcp.Description("End time (e.g., now).")),
	)
}

// ============ Datasource Tools ============

// GetDataSourceByNameTool retrieves a datasource by name.
func GetDataSourceByNameTool() mcp.Tool {
	logrus.Debug("Creating GetDataSourceByNameTool")
	return mcp.NewTool("grafana_get_datasource_by_name",
		mcp.WithDescription("Retrieve a Grafana data source by its name."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("Name of the data source to retrieve.")),
	)
}

// ============ Rendering Tools ============

// RenderPanelImageTool renders a dashboard panel to an image.
func RenderPanelImageTool() mcp.Tool {
	logrus.Debug("Creating RenderPanelImageTool")
	return mcp.NewTool("grafana_render_panel_image",
		mcp.WithDescription("Render a Grafana dashboard panel to an image (PNG)."),
		mcp.WithString("dashboardUID", mcp.Required(),
			mcp.Description("Dashboard UID containing the panel.")),
		mcp.WithNumber("panelID", mcp.Required(),
			mcp.Description("ID of the panel to render.")),
		mcp.WithNumber("width",
			mcp.Description("Image width in pixels. Default: 800.")),
		mcp.WithNumber("height",
			mcp.Description("Image height in pixels. Default: 400.")),
		mcp.WithString("from",
			mcp.Description("Start time for the panel data (e.g., now-1h).")),
		mcp.WithString("to",
			mcp.Description("End time for the panel data (e.g., now).")),
		mcp.WithString("timeout",
			mcp.Description("Render timeout in seconds. Default: 30.")),
	)
}

// ============ Graphite Annotation Tool ============

// CreateGraphiteAnnotationTool creates a Graphite annotation.
func CreateGraphiteAnnotationTool() mcp.Tool {
	logrus.Debug("Creating CreateGraphiteAnnotationTool")
	return mcp.NewTool("grafana_create_graphite_annotation",
		mcp.WithDescription("Create an annotation via Graphite API."),
		mcp.WithString("what", mcp.Required(),
			mcp.Description("What happened (event title).")),
		mcp.WithString("data",
			mcp.Description("Additional data or description.")),
		mcp.WithNumber("timestamp",
			mcp.Description("Unix timestamp for the event. Default: now.")),
		mcp.WithString("tags",
			mcp.Description("Comma-separated tags for the annotation.")),
	)
}

// ============ Datasource Management Tools ============

// CreateDatasourceTool creates a new datasource in Grafana.
func CreateDatasourceTool() mcp.Tool {
	logrus.Debug("Creating CreateDatasourceTool")
	return mcp.NewTool("grafana_create_datasource",
		mcp.WithDescription("Create a new datasource in Grafana. Requires datasource write permissions."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("Name of the datasource to create.")),
		mcp.WithString("type", mcp.Required(),
			mcp.Description("Type of the datasource (e.g., prometheus, influxdb, elasticsearch, postgresql, mysql, etc.).")),
		mcp.WithString("url", mcp.Required(),
			mcp.Description("Connection URL for the datasource.")),
		mcp.WithString("access",
			mcp.Description("Access mode: proxy (default) or direct.")),
		mcp.WithString("database",
			mcp.Description("Database name (for SQL-based datasources).")),
		mcp.WithString("user",
			mcp.Description("Username for authentication.")),
		mcp.WithString("password",
			mcp.Description("Password for authentication.")),
		mcp.WithObject("jsonData",
			mcp.Description("Additional JSON configuration data specific to the datasource type.")),
		mcp.WithObject("secureJsonData",
			mcp.Description("Secure JSON data (passwords, tokens, etc.) that will be encrypted.")),
		mcp.WithBoolean("isDefault",
			mcp.Description("Set as default datasource. Default: false.")),
	)
}

// UpdateDatasourceTool updates an existing datasource in Grafana.
func UpdateDatasourceTool() mcp.Tool {
	logrus.Debug("Creating UpdateDatasourceTool")
	return mcp.NewTool("grafana_update_datasource",
		mcp.WithDescription("Update an existing datasource in Grafana. Requires datasource write permissions."),
		mcp.WithString("uid", mcp.Required(),
			mcp.Description("Unique identifier (UID) of the datasource to update.")),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("New name for the datasource.")),
		mcp.WithString("type", mcp.Required(),
			mcp.Description("Type of the datasource (e.g., prometheus, influxdb, elasticsearch, postgresql, mysql, etc.).")),
		mcp.WithString("url", mcp.Required(),
			mcp.Description("New connection URL for the datasource.")),
		mcp.WithString("access",
			mcp.Description("Access mode: proxy or direct.")),
		mcp.WithString("database",
			mcp.Description("Database name (for SQL-based datasources).")),
		mcp.WithString("user",
			mcp.Description("Username for authentication.")),
		mcp.WithString("password",
			mcp.Description("Password for authentication.")),
		mcp.WithObject("jsonData",
			mcp.Description("Additional JSON configuration data specific to the datasource type.")),
		mcp.WithObject("secureJsonData",
			mcp.Description("Secure JSON data (passwords, tokens, etc.) that will be encrypted.")),
		mcp.WithBoolean("isDefault",
			mcp.Description("Set as default datasource.")),
	)
}

// DeleteDatasourceTool deletes a datasource from Grafana.
func DeleteDatasourceTool() mcp.Tool {
	logrus.Debug("Creating DeleteDatasourceTool")
	return mcp.NewTool("grafana_delete_datasource",
		mcp.WithDescription("Delete a datasource from Grafana by its UID. Requires datasource write permissions. This action cannot be undone."),
		mcp.WithString("uid", mcp.Required(),
			mcp.Description("Unique identifier (UID) of the datasource to delete.")),
	)
}

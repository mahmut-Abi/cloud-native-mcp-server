// Package tools provides MCP tool definitions for Kibana operations.
// It defines the structure and parameters for all Kibana-related tools.
package tools

import "github.com/mark3labs/mcp-go/mcp"

// GetSpacesTool returns the tool definition for retrieving Kibana spaces.
func GetSpacesTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_spaces",
		Description: "Retrieve all Kibana spaces with their configurations and metadata",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// GetSpaceTool returns the tool definition for retrieving a specific Kibana space.
func GetSpaceTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_space",
		Description: "Retrieve details of a specific Kibana space by ID",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"space_id": map[string]interface{}{
					"type":        "string",
					"description": "The unique identifier of the Kibana space to retrieve",
				},
			},
			Required: []string{"space_id"},
		},
	}
}

// GetIndexPatternsTool returns the tool definition for retrieving Kibana index patterns.
func GetIndexPatternsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_index_patterns",
		Description: "Retrieve all index patterns from Kibana with their field mappings and configurations",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// GetDashboardsTool returns the tool definition for retrieving Kibana dashboards.
func GetDashboardsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_dashboards",
		Description: "Retrieve all dashboards from Kibana with their metadata and basic configuration",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// GetDashboardTool returns the tool definition for retrieving a specific Kibana dashboard.
func GetDashboardTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_dashboard",
		Description: "Retrieve detailed information about a specific Kibana dashboard by ID",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"dashboard_id": map[string]interface{}{
					"type":        "string",
					"description": "The unique identifier of the dashboard to retrieve",
				},
			},
			Required: []string{"dashboard_id"},
		},
	}
}

// GetVisualizationsTool returns the tool definition for retrieving Kibana visualizations.
func GetVisualizationsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_visualizations",
		Description: "Retrieve all visualizations from Kibana with their configurations and metadata",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// SearchSavedObjectsTool returns the tool definition for searching Kibana saved objects.
func SearchSavedObjectsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_search_saved_objects",
		Description: "Search and filter Kibana saved objects by type, name, or other criteria with pagination support",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"type": map[string]interface{}{
					"type":        "string",
					"description": "The type of saved object to search for (e.g., 'dashboard', 'visualization', 'index-pattern', 'search')",
				},
				"search": map[string]interface{}{
					"type":        "string",
					"description": "Search term to filter saved objects by title or other text fields",
				},
				"page": map[string]interface{}{
					"type":        "number",
					"description": "Page number for pagination (starts from 1)",
					"default":     1,
				},
				"per_page": map[string]interface{}{
					"type":        "number",
					"description": "Number of results per page (max 100)",
					"default":     20,
				},
			},
		},
	}
}

// TestConnectionTool returns the tool definition for testing Kibana connection.
func TestConnectionTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_test_connection",
		Description: "Test the connection to the Kibana server and verify API access",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// GetVisualizationTool returns the tool definition for retrieving a specific visualization.
func GetVisualizationTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_visualization",
		Description: "Retrieve detailed information about a specific Kibana visualization by ID",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"visualization_id": map[string]interface{}{
					"type":        "string",
					"description": "The unique identifier of the visualization to retrieve",
				},
			},
			Required: []string{"visualization_id"},
		},
	}
}

// GetIndexPatternTool returns the tool definition for retrieving a specific index pattern.
func GetIndexPatternTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_index_pattern",
		Description: "Retrieve detailed information about a specific Kibana index pattern by ID including fields and configuration",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"index_pattern_id": map[string]interface{}{
					"type":        "string",
					"description": "The unique identifier of the index pattern to retrieve",
				},
			},
			Required: []string{"index_pattern_id"},
		},
	}
}

// GetSavedSearchesTool returns the tool definition for retrieving all saved searches.
func GetSavedSearchesTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_saved_searches",
		Description: "Retrieve all saved searches from Kibana with their configurations and metadata",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// GetSavedSearchTool returns the tool definition for retrieving a specific saved search.

// GetDashboardsSummaryTool returns tool definition for getting Kibana dashboards summary
func GetDashboardsSummaryTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_dashboards_summary",
		Description: "List Kibana dashboards summary (id, title). 70-85% smaller output.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"page": map[string]interface{}{
					"type":    "number",
					"default": 1,
				},
				"per_page": map[string]interface{}{
					"type":    "number",
					"default": 20,
				},
			},
		},
	}
}

// GetVisualizationsSummaryTool returns tool definition for getting Kibana visualizations summary
func GetVisualizationsSummaryTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_visualizations_summary",
		Description: "List Kibana visualizations summary (id, title, type). 70-85% smaller output.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"page": map[string]interface{}{
					"type":    "number",
					"default": 1,
				},
			},
		},
	}
}

// GetIndexPatternsSummaryTool returns tool definition for getting Kibana index patterns summary
func GetIndexPatternsSummaryTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_index_patterns_summary",
		Description: "List Kibana index patterns summary (id, title, pattern). 70-85% smaller output.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

func GetSavedSearchTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_saved_search",
		Description: "Retrieve detailed information about a specific Kibana saved search by ID",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"search_id": map[string]interface{}{
					"type":        "string",
					"description": "The unique identifier of the saved search to retrieve",
				},
			},
			Required: []string{"search_id"},
		},
	}
}

// GetKibanaStatusTool returns the tool definition for retrieving Kibana status.
func GetKibanaStatusTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_status",
		Description: "Retrieve Kibana health status and version information including server state and metrics",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// ⚠️ PRIORITY: Optimized tools for LLM efficiency

// GetSpacesSummaryTool returns tool definition for getting Kibana spaces summary
func GetSpacesSummaryTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_spaces_summary",
		Description: "⚠️ PRIORITY: Get Kibana spaces summary (id, name, description). 80-90% smaller output. Optimized for LLM efficiency.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Maximum number of spaces to return (default: 50, max: 100)",
				},
			},
		},
	}
}

// GetDashboardsPaginatedTool returns tool definition for paginated dashboards listing
func GetDashboardsPaginatedTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_dashboards_paginated",
		Description: "⚠️ PRIORITY: Optimized for LLM efficiency: List Kibana dashboards with pagination and summary output. 80-90% smaller than full listing.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"page": map[string]interface{}{
					"type":        "number",
					"description": "Page number for pagination (starts from 1)",
					"default":     1,
				},
				"per_page": map[string]interface{}{
					"type":        "number",
					"description": "Number of results per page (default: 20, max: 100)",
					"default":     20,
				},
				"search": map[string]interface{}{
					"type":        "string",
					"description": "Search term to filter dashboards by title",
				},
				"include_description": map[string]interface{}{
					"type":        "boolean",
					"description": "Include dashboard description (adds minimal data). Default: false",
					"default":     false,
				},
			},
		},
	}
}

// GetVisualizationsPaginatedTool returns tool definition for paginated visualizations listing
func GetVisualizationsPaginatedTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_visualizations_paginated",
		Description: "⚠️ PRIORITY: Optimized for LLM efficiency: List Kibana visualizations with pagination and summary output. 80-90% smaller than full listing.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"page": map[string]interface{}{
					"type":        "number",
					"description": "Page number for pagination (starts from 1)",
					"default":     1,
				},
				"per_page": map[string]interface{}{
					"type":        "number",
					"description": "Number of results per page (default: 20, max: 100)",
					"default":     20,
				},
				"search": map[string]interface{}{
					"type":        "string",
					"description": "Search term to filter visualizations by title",
				},
				"type": map[string]interface{}{
					"type":        "string",
					"description": "Filter by visualization type (e.g., 'histogram', 'pie', 'table')",
				},
			},
		},
	}
}

// GetSavedObjectsAdvancedTool returns tool definition for advanced saved objects search
func GetSavedObjectsAdvancedTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_search_saved_objects_advanced",
		Description: "🔍 Advanced search Kibana saved objects with enhanced filters and pagination. Optimized for finding specific objects.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"type": map[string]interface{}{
					"type":        "string",
					"description": "The type of saved object to search for (e.g., 'dashboard', 'visualization', 'index-pattern', 'search')",
				},
				"search": map[string]interface{}{
					"type":        "string",
					"description": "Search term to filter saved objects by title or other text fields",
				},
				"page": map[string]interface{}{
					"type":        "number",
					"description": "Page number for pagination (starts from 1)",
					"default":     1,
				},
				"per_page": map[string]interface{}{
					"type":        "number",
					"description": "Number of results per page (default: 30, max: 200)",
					"default":     30,
				},
				"sort_field": map[string]interface{}{
					"type":        "string",
					"description": "Field to sort by: title, updated_at, created_at (default: title)",
					"default":     "title",
				},
				"sort_order": map[string]interface{}{
					"type":        "string",
					"description": "Sort order: asc or desc (default: asc)",
					"default":     "asc",
				},
				"has_reference": map[string]interface{}{
					"type":        "string",
					"description": "Filter objects that reference a specific object ID",
				},
				"fields": map[string]interface{}{
					"type":        "array",
					"description": "Specific fields to return in results (reduces output size)",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
			},
		},
	}
}

// GetDashboardDetailAdvancedTool returns tool definition for advanced dashboard details
func GetDashboardDetailAdvancedTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_dashboard_detail_advanced",
		Description: "🔍 Advanced dashboard detail retrieval with enhanced formatting and optional components. Use when comprehensive analysis needed.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"dashboard_id": map[string]interface{}{
					"type":        "string",
					"description": "The unique identifier of the dashboard to retrieve",
				},
				"include_panels": map[string]interface{}{
					"type":        "boolean",
					"description": "Include detailed panel information. Default: true",
					"default":     true,
				},
				"include_ui_state": map[string]interface{}{
					"type":        "boolean",
					"description": "Include UI state configuration. Default: false",
					"default":     false,
				},
				"include_time_options": map[string]interface{}{
					"type":        "boolean",
					"description": "Include time range options. Default: true",
					"default":     true,
				},
				"output_format": map[string]interface{}{
					"type":        "string",
					"description": "Output format: structured (default), compact, or verbose",
					"default":     "structured",
				},
			},
			Required: []string{"dashboard_id"},
		},
	}
}

// GetKibanaHealthSummaryTool returns tool definition for Kibana health summary
func GetKibanaHealthSummaryTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_health_summary",
		Description: "⚠️ PRIORITY: Get Kibana health and status summary (status, version, metrics). Lightweight health overview. Optimized for monitoring.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"level": map[string]interface{}{
					"type":        "string",
					"description": "Health detail level: basic (default), detailed, or metrics",
					"default":     "basic",
				},
				"include_saved_objects": map[string]interface{}{
					"type":        "boolean",
					"description": "Include saved objects statistics (adds data). Default: false",
					"default":     false,
				},
			},
		},
	}
}

// ============ Analysis & Discovery Tools ============

// QueryLogsTool returns tool definition for log search.
func QueryLogsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_query_logs",
		Description: "🔍 Direct log search through Kibana. Quickest way to view and analyze logs. Supports query syntax, sorting, and result limiting.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"indexPattern": map[string]interface{}{
					"type":        "string",
					"description": "Optional index pattern to search (e.g., 'logs-*'). If not specified, searches all indices",
				},
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Lucene query string. Default: '*' (match all)",
					"default":     "*",
				},
				"size": map[string]interface{}{
					"type":        "number",
					"description": "Number of log entries to return. Default: 20, max: 1000",
					"default":     20,
				},
				"sortBy": map[string]interface{}{
					"type":        "string",
					"description": "Field to sort by. Default: @timestamp",
					"default":     "@timestamp",
				},
				"sortOrder": map[string]interface{}{
					"type":        "string",
					"description": "Sort order: asc or desc. Default: desc",
					"default":     "desc",
				},
			},
		},
	}
}

// GetCanvasWorkpadsTool returns tool definition for Canvas workpads.
func GetCanvasWorkpadsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_canvas_workpads",
		Description: "🎨 Retrieve all Canvas workpads. Canvas provides pixel-perfect, self-contained canvas reports for custom visualizations.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// GetLensObjectsTool returns tool definition for Lens visualizations.
func GetLensObjectsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_lens_objects",
		Description: "📊 Retrieve all Lens visualizations. Lens is Kibana's drag-and-drop visualization tool for creating complex queries easily.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// GetMapsTool returns tool definition for Kibana Maps.
func GetMapsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_maps",
		Description: "🗺️ Retrieve all Kibana Maps. Maps enable visualization of geospatial data and provide rich mapping capabilities.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// GetKibanaAlertsTool returns tool definition for Kibana alerting rules.
func GetKibanaAlertsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_alerts",
		Description: "🚨 Retrieve all Kibana alerting rules. Monitor your data and get notified when conditions are met.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// GetIndexPatternFieldsTool returns tool definition for index pattern fields.
func GetIndexPatternFieldsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_index_pattern_fields",
		Description: "📋 Retrieve all fields for an index pattern. Essential for building queries and understanding available data.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"patternID": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the index pattern to retrieve fields for",
				},
			},
			Required: []string{"patternID"},
		},
	}
}

// ============ Write Operations: Spaces ============

// CreateSpaceTool returns tool definition for creating a new Kibana space
func CreateSpaceTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_create_space",
		Description: "🏗️ Create a new Kibana space. Spaces help organize your dashboards, visualizations, and other saved objects into logical groupings.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": "Unique identifier for the space (recommended: lowercase, hyphens only)",
				},
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Display name for the space",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Optional description of the space purpose",
				},
				"color": map[string]interface{}{
					"type":        "string",
					"description": "Optional hex color code for space identification (e.g., '#FF5733')",
				},
				"initials": map[string]interface{}{
					"type":        "string",
					"description": "Optional initials for the space avatar (2 characters)",
				},
				"disabledFeatures": map[string]interface{}{
					"type":        "array",
					"description": "List of features to disable in this space (e.g., ['visualize', 'dashboard'])",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
			},
			Required: []string{"id", "name"},
		},
	}
}

// UpdateSpaceTool returns tool definition for updating an existing Kibana space
func UpdateSpaceTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_update_space",
		Description: "✏️ Update an existing Kibana space. Modify the name, description, color, or other properties of a space.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"space_id": map[string]interface{}{
					"type":        "string",
					"description": "The unique identifier of the space to update",
				},
				"name": map[string]interface{}{
					"type":        "string",
					"description": "New display name for the space",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "New description for the space",
				},
				"color": map[string]interface{}{
					"type":        "string",
					"description": "New hex color code for the space",
				},
				"initials": map[string]interface{}{
					"type":        "string",
					"description": "New initials for the space avatar",
				},
				"disabledFeatures": map[string]interface{}{
					"type":        "array",
					"description": "Updated list of disabled features",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
			},
			Required: []string{"space_id"},
		},
	}
}

// DeleteSpaceTool returns tool definition for deleting a Kibana space
func DeleteSpaceTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_delete_space",
		Description: "🗑️ Delete a Kibana space and all its saved objects. ⚠️ This action cannot be undone!",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"space_id": map[string]interface{}{
					"type":        "string",
					"description": "The unique identifier of the space to delete",
				},
				"force": map[string]interface{}{
					"type":        "boolean",
					"description": "Force deletion even if space has saved objects (default: false)",
					"default":     false,
				},
			},
			Required: []string{"space_id"},
		},
	}
}

// ============ Write Operations: Index Patterns ============

// CreateIndexPatternTool returns tool definition for creating a new index pattern
func CreateIndexPatternTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_create_index_pattern",
		Description: "🔍 Create a new index pattern in Kibana. Index patterns define how to access your Elasticsearch indices for Discover, visualizations, and dashboards.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"title": map[string]interface{}{
					"type":        "string",
					"description": "Index pattern title (supports wildcards, e.g., 'logs-*', 'metricbeat-*')",
				},
				"timeField": map[string]interface{}{
					"type":        "string",
					"description": "Optional name of the time field for time-based queries (e.g., '@timestamp', 'timestamp')",
				},
				"customSource": map[string]interface{}{
					"type":        "map",
					"description": "Optional custom source configuration",
				},
				"fieldFormatMap": map[string]interface{}{
					"type":        "map",
					"description": "Optional field format mappings for custom display formats",
				},
			},
			Required: []string{"title"},
		},
	}
}

// UpdateIndexPatternTool returns tool definition for updating an index pattern
func UpdateIndexPatternTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_update_index_pattern",
		Description: "✏️ Update an existing index pattern. Modify the title, time field, or format configurations.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"index_pattern_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the index pattern to update",
				},
				"title": map[string]interface{}{
					"type":        "string",
					"description": "New index pattern title",
				},
				"timeField": map[string]interface{}{
					"type":        "string",
					"description": "New time field name",
				},
				"fieldFormatMap": map[string]interface{}{
					"type":        "map",
					"description": "Updated field format mappings",
				},
			},
			Required: []string{"index_pattern_id"},
		},
	}
}

// DeleteIndexPatternTool returns tool definition for deleting an index pattern
func DeleteIndexPatternTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_delete_index_pattern",
		Description: "🗑️ Delete an index pattern from Kibana. This removes the pattern but does not affect the underlying Elasticsearch indices.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"index_pattern_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the index pattern to delete",
				},
			},
			Required: []string{"index_pattern_id"},
		},
	}
}

// SetDefaultIndexPatternTool returns tool definition for setting the default index pattern
func SetDefaultIndexPatternTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_set_default_index_pattern",
		Description: "⭐ Set an index pattern as the default for new visualizations and dashboards.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"index_pattern_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the index pattern to set as default",
				},
			},
			Required: []string{"index_pattern_id"},
		},
	}
}

// RefreshIndexPatternFieldsTool returns tool definition for refreshing index pattern fields
func RefreshIndexPatternFieldsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_refresh_index_pattern_fields",
		Description: "🔄 Refresh and sync the fields for an index pattern from Elasticsearch. Use this when new fields are added to your indices.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"index_pattern_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the index pattern to refresh fields for",
				},
			},
			Required: []string{"index_pattern_id"},
		},
	}
}

// ============ Write Operations: Dashboards ============

// CreateDashboardTool returns tool definition for creating a new dashboard
func CreateDashboardTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_create_dashboard",
		Description: "📊 Create a new dashboard in Kibana. Dashboards display collections of visualizations and saved searches.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"title": map[string]interface{}{
					"type":        "string",
					"description": "Dashboard title",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Optional description of the dashboard",
				},
				"timeRestore": map[string]interface{}{
					"type":        "boolean",
					"description": "Enable time picker and time-based filtering (default: true)",
					"default":     true,
				},
				"timeFrom": map[string]interface{}{
					"type":        "string",
					"description": "Default time range start (e.g., 'now-24h')",
				},
				"timeTo": map[string]interface{}{
					"type":        "string",
					"description": "Default time range end (e.g., 'now')",
				},
				"refreshInterval": map[string]interface{}{
					"type":        "map",
					"description": "Auto-refresh interval configuration",
				},
			},
			Required: []string{"title"},
		},
	}
}

// UpdateDashboardTool returns tool definition for updating a dashboard
func UpdateDashboardTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_update_dashboard",
		Description: "✏️ Update an existing dashboard. Modify the title, description, panels, or time settings.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"dashboard_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the dashboard to update",
				},
				"title": map[string]interface{}{
					"type":        "string",
					"description": "New dashboard title",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "New dashboard description",
				},
				"panelsJSON": map[string]interface{}{
					"type":        "string",
					"description": "JSON string of dashboard panels layout",
				},
				"timeFrom": map[string]interface{}{
					"type":        "string",
					"description": "Updated default time range start",
				},
				"timeTo": map[string]interface{}{
					"type":        "string",
					"description": "Updated default time range end",
				},
				"version": map[string]interface{}{
					"type":        "number",
					"description": "Current version number for optimistic locking",
				},
			},
			Required: []string{"dashboard_id"},
		},
	}
}

// DeleteDashboardTool returns tool definition for deleting a dashboard
func DeleteDashboardTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_delete_dashboard",
		Description: "🗑️ Delete a dashboard from Kibana. This removes the dashboard but not its referenced visualizations.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"dashboard_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the dashboard to delete",
				},
			},
			Required: []string{"dashboard_id"},
		},
	}
}

// CloneDashboardTool returns tool definition for cloning a dashboard
func CloneDashboardTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_clone_dashboard",
		Description: "📋 Create a copy of an existing dashboard with a new title.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"dashboard_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the dashboard to clone",
				},
				"new_title": map[string]interface{}{
					"type":        "string",
					"description": "Title for the new cloned dashboard",
				},
			},
			Required: []string{"dashboard_id", "new_title"},
		},
	}
}

// ============ Write Operations: Visualizations ============

// CreateVisualizationTool returns tool definition for creating a new visualization
func CreateVisualizationTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_create_visualization",
		Description: "📈 Create a new visualization in Kibana. Visualizations display data from Elasticsearch indices in various chart types.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"title": map[string]interface{}{
					"type":        "string",
					"description": "Visualization title",
				},
				"visState": map[string]interface{}{
					"type":        "map",
					"description": "Visualization state configuration (type, aggs, params)",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Optional description of the visualization",
				},
				"savedSearchRefName": map[string]interface{}{
					"type":        "string",
					"description": "Reference to an existing saved search (optional)",
				},
			},
			Required: []string{"title"},
		},
	}
}

// UpdateVisualizationTool returns tool definition for updating a visualization
func UpdateVisualizationTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_update_visualization",
		Description: "✏️ Update an existing visualization. Modify the title, visState, or description.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"visualization_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the visualization to update",
				},
				"title": map[string]interface{}{
					"type":        "string",
					"description": "New visualization title",
				},
				"visState": map[string]interface{}{
					"type":        "map",
					"description": "Updated visualization state configuration",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "New description",
				},
				"version": map[string]interface{}{
					"type":        "number",
					"description": "Current version number for optimistic locking",
				},
			},
			Required: []string{"visualization_id"},
		},
	}
}

// DeleteVisualizationTool returns tool definition for deleting a visualization
func DeleteVisualizationTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_delete_visualization",
		Description: "🗑️ Delete a visualization from Kibana.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"visualization_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the visualization to delete",
				},
			},
			Required: []string{"visualization_id"},
		},
	}
}

// CloneVisualizationTool returns tool definition for cloning a visualization
func CloneVisualizationTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_clone_visualization",
		Description: "📋 Create a copy of an existing visualization with a new title.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"visualization_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the visualization to clone",
				},
				"new_title": map[string]interface{}{
					"type":        "string",
					"description": "Title for the new cloned visualization",
				},
			},
			Required: []string{"visualization_id", "new_title"},
		},
	}
}

// ============ Write Operations: Saved Objects (Generic) ============

// CreateSavedObjectTool returns tool definition for creating a generic saved object
func CreateSavedObjectTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_create_saved_object",
		Description: "💾 Create a generic saved object in Kibana. Supports any saved object type (dashboard, visualization, search, etc.).",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"type": map[string]interface{}{
					"type":        "string",
					"description": "Saved object type (e.g., 'dashboard', 'visualization', 'search', 'index-pattern', 'lens', 'map', 'canvas-workpad')",
				},
				"attributes": map[string]interface{}{
					"type":        "map",
					"description": "Object attributes (title, description, visState, etc.)",
				},
				"references": map[string]interface{}{
					"type":        "array",
					"description": "Array of object references for linking related objects",
					"items": map[string]interface{}{
						"type": "object",
					},
				},
				"initialObjectType": map[string]interface{}{
					"type":        "string",
					"description": "Optional initial object type for migration",
				},
			},
			Required: []string{"type", "attributes"},
		},
	}
}

// UpdateSavedObjectTool returns tool definition for updating a saved object
func UpdateSavedObjectTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_update_saved_object",
		Description: "✏️ Update attributes of a saved object. Use this for partial updates.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"type": map[string]interface{}{
					"type":        "string",
					"description": "Saved object type",
				},
				"object_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the object to update",
				},
				"attributes": map[string]interface{}{
					"type":        "map",
					"description": "Updated attributes to merge with existing",
				},
				"references": map[string]interface{}{
					"type":        "array",
					"description": "Updated object references",
					"items": map[string]interface{}{
						"type": "object",
					},
				},
				"version": map[string]interface{}{
					"type":        "string",
					"description": "Current version for optimistic locking",
				},
			},
			Required: []string{"type", "object_id"},
		},
	}
}

// DeleteSavedObjectTool returns tool definition for deleting a saved object
func DeleteSavedObjectTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_delete_saved_object",
		Description: "🗑️ Delete a saved object from Kibana by type and ID.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"type": map[string]interface{}{
					"type":        "string",
					"description": "Saved object type",
				},
				"object_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the object to delete",
				},
				"force": map[string]interface{}{
					"type":        "boolean",
					"description": "Force deletion even if object is referenced (default: false)",
					"default":     false,
				},
			},
			Required: []string{"type", "object_id"},
		},
	}
}

// BulkDeleteSavedObjectsTool returns tool definition for deleting multiple saved objects
func BulkDeleteSavedObjectsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_bulk_delete_saved_objects",
		Description: "🗑️ Delete multiple saved objects in a single operation for efficiency.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"objects": map[string]interface{}{
					"type":        "array",
					"description": "Array of objects to delete with type and id",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"type": map[string]interface{}{
								"type": "string",
							},
							"id": map[string]interface{}{
								"type": "string",
							},
						},
					},
				},
			},
			Required: []string{"objects"},
		},
	}
}

// ExportSavedObjectsTool returns tool definition for exporting saved objects
func ExportSavedObjectsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_export_saved_objects",
		Description: "📦 Export saved objects to a JSON file for backup or migration to another Kibana instance.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"objects": map[string]interface{}{
					"type":        "array",
					"description": "Array of objects to export with type and id",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"type": map[string]interface{}{
								"type": "string",
							},
							"id": map[string]interface{}{
								"type": "string",
							},
						},
					},
				},
				"includeReferences": map[string]interface{}{
					"type":        "boolean",
					"description": "Include all referenced objects in the export (default: true)",
					"default":     true,
				},
			},
			Required: []string{"objects"},
		},
	}
}

// ImportSavedObjectsTool returns tool definition for importing saved objects
func ImportSavedObjectsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_import_saved_objects",
		Description: "📥 Import saved objects from a JSON file. Supports objects exported from kibana_export_saved_objects.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"file": map[string]interface{}{
					"type":        "string",
					"description": "Base64-encoded JSON content of the exported objects",
				},
				"createNewCopies": map[string]interface{}{
					"type":        "boolean",
					"description": "Create new copies with new IDs instead of overwriting (default: false)",
					"default":     false,
				},
			},
			Required: []string{"file"},
		},
	}
}

// ============ Alerting Rules ============

// GetAlertRulesTool returns tool definition for listing alert rules
func GetAlertRulesTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_alert_rules",
		Description: "🚨 List all alert rules with filtering and pagination. Returns rule summaries including name, type, schedule, and enabled status.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"page": map[string]interface{}{
					"type":        "number",
					"description": "Page number for pagination (default: 1)",
					"default":     1,
				},
				"per_page": map[string]interface{}{
					"type":        "number",
					"description": "Number of results per page (default: 20, max: 100)",
					"default":     20,
				},
				"filter": map[string]interface{}{
					"type":        "string",
					"description": "Filter rules by name or type",
				},
				"enabled": map[string]interface{}{
					"type":        "boolean",
					"description": "Filter by enabled status (true/false, omit for all)",
				},
			},
		},
	}
}

// GetAlertRuleTool returns tool definition for getting a specific alert rule
func GetAlertRuleTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_alert_rule",
		Description: "📋 Get detailed information about a specific alert rule including its configuration, schedule, actions, and recent execution status.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"rule_id": map[string]interface{}{
					"type":        "string",
					"description": "The unique identifier of the alert rule",
				},
			},
			Required: []string{"rule_id"},
		},
	}
}

// CreateAlertRuleTool returns tool definition for creating an alert rule
func CreateAlertRuleTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_create_alert_rule",
		Description: "🔔 Create a new alert rule in Kibana. Define when and how to be notified when conditions are met.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Display name for the alert rule",
				},
				"alertTypeId": map[string]interface{}{
					"type":        "string",
					"description": "Alert type identifier (e.g., 'threshold', 'inventory_changed', 'metrics_threshold')",
				},
				"schedule": map[string]interface{}{
					"type":        "map",
					"description": "Alert schedule (e.g., {\"interval\": \"5m\"})",
				},
				"params": map[string]interface{}{
					"type":        "map",
					"description": "Alert-specific parameters",
				},
				"actions": map[string]interface{}{
					"type":        "array",
					"description": "Actions to execute when alert triggers",
					"items": map[string]interface{}{
						"type": "object",
					},
				},
				"tags": map[string]interface{}{
					"type":        "array",
					"description": "Tags for organizing alert rules",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
			},
			Required: []string{"name", "alertTypeId", "schedule"},
		},
	}
}

// UpdateAlertRuleTool returns tool definition for updating an alert rule
func UpdateAlertRuleTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_update_alert_rule",
		Description: "✏️ Update an existing alert rule. Modify the name, schedule, actions, or parameters.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"rule_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the alert rule to update",
				},
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Updated display name",
				},
				"schedule": map[string]interface{}{
					"type":        "map",
					"description": "Updated schedule",
				},
				"params": map[string]interface{}{
					"type":        "map",
					"description": "Updated alert parameters",
				},
				"actions": map[string]interface{}{
					"type":        "array",
					"description": "Updated actions",
					"items": map[string]interface{}{
						"type": "object",
					},
				},
				"tags": map[string]interface{}{
					"type":        "array",
					"description": "Updated tags",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
			},
			Required: []string{"rule_id"},
		},
	}
}

// DeleteAlertRuleTool returns tool definition for deleting an alert rule
func DeleteAlertRuleTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_delete_alert_rule",
		Description: "🗑️ Delete an alert rule from Kibana. This also removes all associated history and execution records.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"rule_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the alert rule to delete",
				},
			},
			Required: []string{"rule_id"},
		},
	}
}

// EnableAlertRuleTool returns tool definition for enabling an alert rule
func EnableAlertRuleTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_enable_alert_rule",
		Description: "✅ Enable a disabled alert rule. The rule will start executing according to its schedule.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"rule_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the alert rule to enable",
				},
			},
			Required: []string{"rule_id"},
		},
	}
}

// DisableAlertRuleTool returns tool definition for disabling an alert rule
func DisableAlertRuleTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_disable_alert_rule",
		Description: "⏸️ Disable an alert rule. The rule will stop executing but is retained for future use.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"rule_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the alert rule to disable",
				},
			},
			Required: []string{"rule_id"},
		},
	}
}

// MuteAlertRuleTool returns tool definition for muting an alert rule
func MuteAlertRuleTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_mute_alert_rule",
		Description: "🔇 Mute all alerts for a specific rule for a defined time period.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"rule_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the alert rule to mute",
				},
				"duration": map[string]interface{}{
					"type":        "string",
					"description": "Duration to mute (e.g., '1h', '30m', '7d')",
				},
			},
			Required: []string{"rule_id", "duration"},
		},
	}
}

// UnmuteAlertRuleTool returns tool definition for unmuting an alert rule
func UnmuteAlertRuleTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_unmute_alert_rule",
		Description: "🔊 Unmute a previously muted alert rule.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"rule_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the alert rule to unmute",
				},
			},
			Required: []string{"rule_id"},
		},
	}
}

// GetAlertRuleTypesTool returns tool definition for listing available alert rule types
func GetAlertRuleTypesTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_alert_rule_types",
		Description: "📋 List all available alert rule types with their parameters and configuration options.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// GetAlertRuleHistoryTool returns tool definition for getting alert rule execution history
func GetAlertRuleHistoryTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_alert_rule_history",
		Description: "📊 Get execution history for a specific alert rule including trigger times, status, and error messages.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"rule_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the alert rule",
				},
				"page": map[string]interface{}{
					"type":        "number",
					"description": "Page number (default: 1)",
					"default":     1,
				},
				"per_page": map[string]interface{}{
					"type":        "number",
					"description": "Results per page (default: 20, max: 100)",
					"default":     20,
				},
			},
			Required: []string{"rule_id"},
		},
	}
}

// ============ Connectors ============

// GetConnectorsTool returns tool definition for listing connectors
func GetConnectorsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_connectors",
		Description: "🔗 List all configured action connectors (Slack, email, webhook, etc.) for sending alert notifications.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"page": map[string]interface{}{
					"type":        "number",
					"description": "Page number (default: 1)",
					"default":     1,
				},
				"per_page": map[string]interface{}{
					"type":        "number",
					"description": "Results per page (default: 20, max: 100)",
					"default":     20,
				},
			},
		},
	}
}

// GetConnectorTool returns tool definition for getting a specific connector
func GetConnectorTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_connector",
		Description: "🔗 Get detailed information about a specific connector including its configuration.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"connector_id": map[string]interface{}{
					"type":        "string",
					"description": "The unique identifier of the connector",
				},
			},
			Required: []string{"connector_id"},
		},
	}
}

// CreateConnectorTool returns tool definition for creating a connector
func CreateConnectorTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_create_connector",
		Description: "🔧 Create a new action connector for sending alert notifications via Slack, email, webhook, etc.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Display name for the connector",
				},
				"connectorTypeId": map[string]interface{}{
					"type":        "string",
					"description": "Connector type (e.g., '.slack', '.email', '.webhook', '.pagerduty', '.teams')",
				},
				"config": map[string]interface{}{
					"type":        "map",
					"description": "Connector-specific configuration",
				},
				"secrets": map[string]interface{}{
					"type":        "map",
					"description": "Sensitive configuration (API keys, tokens, etc.)",
				},
			},
			Required: []string{"name", "connectorTypeId"},
		},
	}
}

// UpdateConnectorTool returns tool definition for updating a connector
func UpdateConnectorTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_update_connector",
		Description: "✏️ Update an existing connector's configuration or secrets.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"connector_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the connector to update",
				},
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Updated display name",
				},
				"config": map[string]interface{}{
					"type":        "map",
					"description": "Updated configuration",
				},
				"secrets": map[string]interface{}{
					"type":        "map",
					"description": "Updated secrets",
				},
			},
			Required: []string{"connector_id"},
		},
	}
}

// DeleteConnectorTool returns tool definition for deleting a connector
func DeleteConnectorTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_delete_connector",
		Description: "🗑️ Delete a connector from Kibana. This will also remove it from any alert rules using it.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"connector_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the connector to delete",
				},
			},
			Required: []string{"connector_id"},
		},
	}
}

// TestConnectorTool returns tool definition for testing a connector
func TestConnectorTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_test_connector",
		Description: "🧪 Test a connector by sending a test notification. Verifies the connector is properly configured.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"connector_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the connector to test",
				},
				"body": map[string]interface{}{
					"type":        "map",
					"description": "Optional custom message body to send",
				},
			},
			Required: []string{"connector_id"},
		},
	}
}

// GetConnectorTypesTool returns tool definition for listing available connector types
func GetConnectorTypesTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_connector_types",
		Description: "📋 List all available connector types with their configuration requirements.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// ============ Data Views (Index Patterns v2) ============

// GetDataViewsTool returns tool definition for listing data views
func GetDataViewsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_data_views",
		Description: "📊 List all data views (formerly index patterns) in Kibana.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"page": map[string]interface{}{
					"type":        "number",
					"description": "Page number (default: 1)",
					"default":     1,
				},
				"per_page": map[string]interface{}{
					"type":        "number",
					"description": "Results per page (default: 20, max: 100)",
					"default":     20,
				},
			},
		},
	}
}

// GetDataViewTool returns tool definition for getting a specific data view
func GetDataViewTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_get_data_view",
		Description: "📋 Get detailed information about a specific data view including field mappings and runtime fields.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"data_view_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the data view",
				},
			},
			Required: []string{"data_view_id"},
		},
	}
}

// CreateDataViewTool returns tool definition for creating a data view
func CreateDataViewTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_create_data_view",
		Description: "🔍 Create a new data view in Kibana. Defines how to access one or more Elasticsearch indices.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"title": map[string]interface{}{
					"type":        "string",
					"description": "Data view title/pattern (e.g., 'logs-*', 'metricbeat-*')",
				},
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Optional custom name for the data view",
				},
				"timeField": map[string]interface{}{
					"type":        "string",
					"description": "Optional time field name for time-based queries",
				},
				"sourceFilters": map[string]interface{}{
					"type":        "array",
					"description": "Optional source filters to exclude fields",
					"items": map[string]interface{}{
						"type": "object",
					},
				},
				"fieldFormats": map[string]interface{}{
					"type":        "map",
					"description": "Optional custom field format configurations",
				},
			},
			Required: []string{"title"},
		},
	}
}

// UpdateDataViewTool returns tool definition for updating a data view
func UpdateDataViewTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_update_data_view",
		Description: "✏️ Update an existing data view's configuration.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"data_view_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the data view to update",
				},
				"title": map[string]interface{}{
					"type":        "string",
					"description": "Updated title/pattern",
				},
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Updated name",
				},
				"timeField": map[string]interface{}{
					"type":        "string",
					"description": "Updated time field",
				},
			},
			Required: []string{"data_view_id"},
		},
	}
}

// DeleteDataViewTool returns tool definition for deleting a data view
func DeleteDataViewTool() mcp.Tool {
	return mcp.Tool{
		Name:        "kibana_delete_data_view",
		Description: "🗑️ Delete a data view from Kibana. This does not affect the underlying Elasticsearch indices.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"data_view_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the data view to delete",
				},
			},
			Required: []string{"data_view_id"},
		},
	}
}

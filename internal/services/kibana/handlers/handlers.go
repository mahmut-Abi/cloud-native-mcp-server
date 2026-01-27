// Package handlers provides HTTP handlers for Kibana MCP operations.
// It implements request handling for Kibana spaces, dashboards, visualizations, and saved objects.
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/kibana/client"
	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
)

// marshalIndentJSON performs indented JSON encoding using object pool
func marshalIndentJSON(data interface{}) ([]byte, error) {
	// First encode to compact format using object pool
	compactBytes, err := optimize.GlobalJSONPool.MarshalToBytes(data)
	if err != nil {
		return nil, err
	}

	// For scenarios requiring indented display, still use standard library but reduce allocations
	// This is a trade-off between performance and readability
	var result bytes.Buffer
	err = json.Indent(&result, compactBytes, "", "  ")
	return result.Bytes(), err
}

// HandleGetSpaces handles Kibana spaces retrieval requests.
func HandleGetSpaces(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get spaces handler")

		// Get spaces
		spaces, err := c.GetSpaces(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get spaces: %v", err)),
				},
			}, nil
		}

		// Format result using optimized JSON pool
		resultJSON, err := marshalIndentJSON(spaces)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format spaces: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleGetSpace handles specific Kibana space retrieval requests.
func HandleGetSpace(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get space handler")

		// Extract parameters
		args := req.Params.Arguments
		if args == nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("No arguments provided"),
				},
			}, nil
		}

		// Get space ID parameter
		spaceID, ok := req.GetArguments()["space_id"].(string)
		if !ok || spaceID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Space ID parameter is required"),
				},
			}, nil
		}

		// Get space
		space, err := c.GetSpace(ctx, spaceID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get space: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(space)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format space: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleGetIndexPatterns handles Kibana index patterns retrieval requests.
func HandleGetIndexPatterns(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get index patterns handler")

		// Get index patterns
		indexPatterns, err := c.GetIndexPatterns(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get index patterns: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(indexPatterns)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format index patterns: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleGetDashboards handles Kibana dashboards retrieval requests.
func HandleGetDashboards(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get dashboards handler")

		// Get dashboards
		dashboards, err := c.GetDashboards(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get dashboards: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(dashboards)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format dashboards: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleGetDashboard handles specific Kibana dashboard retrieval requests.
func HandleGetDashboard(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get dashboard handler")

		// Extract parameters
		args := req.Params.Arguments
		if args == nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("No arguments provided"),
				},
			}, nil
		}

		// Get dashboard ID parameter
		dashboardID, ok := req.GetArguments()["dashboard_id"].(string)
		if !ok || dashboardID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Dashboard ID parameter is required"),
				},
			}, nil
		}

		// Get dashboard
		dashboard, err := c.GetDashboard(ctx, dashboardID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get dashboard: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(dashboard)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format dashboard: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleGetVisualizations handles Kibana visualizations retrieval requests.
func HandleGetVisualizations(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get visualizations handler")

		// Get visualizations
		visualizations, err := c.GetVisualizations(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get visualizations: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(visualizations)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format visualizations: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleSearchSavedObjects handles Kibana saved objects search requests.
func HandleSearchSavedObjects(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana search saved objects handler")

		// Get optional parameters
		objectType, _ := req.GetArguments()["type"].(string)
		search, _ := req.GetArguments()["search"].(string)

		page := 1
		if p, exists := req.GetArguments()["page"]; exists {
			if pageStr, ok := p.(string); ok {
				if parsed, err := strconv.Atoi(pageStr); err == nil && parsed > 0 {
					page = parsed
				}
			} else if pageFloat, ok := p.(float64); ok && pageFloat > 0 {
				page = int(pageFloat)
			}
		}

		perPage := 20
		if pp, exists := req.GetArguments()["per_page"]; exists {
			if perPageStr, ok := pp.(string); ok {
				if parsed, err := strconv.Atoi(perPageStr); err == nil && parsed > 0 {
					perPage = parsed
				}
			} else if perPageFloat, ok := pp.(float64); ok && perPageFloat > 0 {
				perPage = int(perPageFloat)
			}
		}

		// Search saved objects
		result, err := c.SearchSavedObjects(ctx, objectType, search, page, perPage)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to search saved objects: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(result)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format search results: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleTestConnection handles Kibana connection test requests.
func HandleTestConnection(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Testing Kibana connection")

		// Test connection
		err := c.TestConnection(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Connection test failed: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Connection test successful"),
			},
		}, nil
	}
}

// Helper function to validate and parse limit parameter with warnings
func parseLimitWithWarnings(request mcp.CallToolRequest, toolName string) int {
	limit := 20
	if v, ok := request.GetArguments()["limit"]; ok {
		if f, ok := v.(float64); ok {
			limit = int(f)
			if limit <= 0 {
				limit = 20
			} else if limit > 100 {
				logrus.WithField("requested", limit).WithField("max", 100).Warn("Limit too high, resetting to safe maximum")
				limit = 100
			}
		}
	}

	if limit > 50 {
		logrus.WithFields(logrus.Fields{
			"tool":  toolName,
			"limit": limit,
		}).Warn("Large limit may cause context overflow, consider using pagination")
	}

	return limit
}

// Helper function to get optional numeric parameter
func getOptionalIntParam(request mcp.CallToolRequest, param string, defaultValue int) int {
	if v, ok := request.GetArguments()[param]; ok {
		if f, ok := v.(float64); ok {
			val := int(f)
			if val > 0 {
				return val
			}
		} else if s, ok := v.(string); ok {
			if val, err := strconv.Atoi(s); err == nil && val > 0 {
				return val
			}
		}
	}
	return defaultValue
}

// Helper function to get optional boolean parameter
func getOptionalBoolParam(request mcp.CallToolRequest, param string) *bool {
	if value, ok := request.GetArguments()[param].(bool); ok {
		return &value
	}
	return nil
}

// Helper function to get optional string parameter
func getOptionalStringParam(request mcp.CallToolRequest, param string) string {
	value, _ := request.GetArguments()[param].(string)
	return value
}

// Helper function to marshal optimized response with size warning
func marshalOptimizedResponse(data any, toolName string) (*mcp.CallToolResult, error) {
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize response: %w", err)
	}

	// Add size warning for large responses
	if len(jsonResponse) > 100*1024 { // 100KB
		logrus.WithFields(logrus.Fields{
			"tool":      toolName,
			"sizeBytes": len(jsonResponse),
			"sizeKB":    len(jsonResponse) / 1024,
		}).Warn("Large response generated")
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(jsonResponse)),
		},
	}, nil
}

// ⚠️ PRIORITY: Optimized handlers for LLM efficiency

// HandleSpacesSummary handles getting spaces summary with LLM optimization
func HandleSpacesSummary(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		limit := parseLimitWithWarnings(request, "kibana_spaces_summary")

		logrus.WithFields(logrus.Fields{
			"tool":  "kibana_spaces_summary",
			"limit": limit,
		}).Debug("Handler invoked")

		spaces, err := c.SpacesSummary(ctx, limit)
		if err != nil {
			return nil, fmt.Errorf("failed to get spaces summary: %w", err)
		}

		response := map[string]interface{}{
			"spaces": spaces,
			"count":  len(spaces),
			"metadata": map[string]interface{}{
				"tool":         "kibana_spaces_summary",
				"optimizedFor": "LLM efficiency",
			},
		}

		return marshalOptimizedResponse(response, "kibana_spaces_summary")
	}
}

// HandleDashboardsPaginated handles paginated dashboards listing with LLM optimization
func HandleDashboardsPaginated(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		page := getOptionalIntParam(request, "page", 1)
		perPage := getOptionalIntParam(request, "per_page", 20)
		search := getOptionalStringParam(request, "search")
		includeDescription := getOptionalBoolParam(request, "include_description")
		if includeDescription == nil {
			defaultDesc := false
			includeDescription = &defaultDesc
		}

		logrus.WithFields(logrus.Fields{
			"tool":               "kibana_dashboards_paginated",
			"page":               page,
			"perPage":            perPage,
			"search":             search,
			"includeDescription": *includeDescription,
		}).Debug("Handler invoked")

		dashboards, pagination, err := c.DashboardsPaginated(ctx, page, perPage, search, *includeDescription)
		if err != nil {
			return nil, fmt.Errorf("failed to list dashboards paginated: %w", err)
		}

		response := map[string]interface{}{
			"dashboards": dashboards,
			"count":      len(dashboards),
			"pagination": pagination,
			"searchCriteria": map[string]interface{}{
				"search":             search,
				"includeDescription": *includeDescription,
			},
			"metadata": map[string]interface{}{
				"tool":         "kibana_dashboards_paginated",
				"optimizedFor": "LLM efficiency",
			},
		}

		return marshalOptimizedResponse(response, "kibana_dashboards_paginated")
	}
}

// HandleVisualizationsPaginated handles paginated visualizations listing with LLM optimization
func HandleVisualizationsPaginated(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		page := getOptionalIntParam(request, "page", 1)
		perPage := getOptionalIntParam(request, "per_page", 20)
		search := getOptionalStringParam(request, "search")
		visType := getOptionalStringParam(request, "type")

		logrus.WithFields(logrus.Fields{
			"tool":    "kibana_visualizations_paginated",
			"page":    page,
			"perPage": perPage,
			"search":  search,
			"visType": visType,
		}).Debug("Handler invoked")

		visualizations, pagination, err := c.VisualizationsPaginated(ctx, page, perPage, search, visType)
		if err != nil {
			return nil, fmt.Errorf("failed to list visualizations paginated: %w", err)
		}

		response := map[string]interface{}{
			"visualizations": visualizations,
			"count":          len(visualizations),
			"pagination":     pagination,
			"searchCriteria": map[string]interface{}{
				"search":  search,
				"visType": visType,
			},
			"metadata": map[string]interface{}{
				"tool":         "kibana_visualizations_paginated",
				"optimizedFor": "LLM efficiency",
			},
		}

		return marshalOptimizedResponse(response, "kibana_visualizations_paginated")
	}
}

// HandleSearchSavedObjectsAdvanced handles advanced saved objects search with enhanced filters
func HandleSearchSavedObjectsAdvanced(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		objectType := getOptionalStringParam(request, "type")
		search := getOptionalStringParam(request, "search")
		page := getOptionalIntParam(request, "page", 1)
		perPage := getOptionalIntParam(request, "per_page", 30)
		sortField := getOptionalStringParam(request, "sort_field")
		sortOrder := getOptionalStringParam(request, "sort_order")
		hasReference := getOptionalStringParam(request, "has_reference")

		var fields []string
		if f, ok := request.GetArguments()["fields"].([]interface{}); ok {
			for _, field := range f {
				if fieldStr, ok := field.(string); ok {
					fields = append(fields, fieldStr)
				}
			}
		}

		logrus.WithFields(logrus.Fields{
			"tool":         "kibana_search_saved_objects_advanced",
			"objectType":   objectType,
			"search":       search,
			"page":         page,
			"perPage":      perPage,
			"sortField":    sortField,
			"sortOrder":    sortOrder,
			"hasReference": hasReference,
			"fields":       fields,
		}).Debug("Handler invoked")

		result, err := c.SearchSavedObjectsAdvanced(ctx, objectType, search, page, perPage, sortField, sortOrder, hasReference, fields)
		if err != nil {
			return nil, fmt.Errorf("failed to search saved objects advanced: %w", err)
		}

		response := map[string]interface{}{
			"savedObjects": result.SavedObjects,
			"count":        len(result.SavedObjects),
			"searchCriteria": map[string]interface{}{
				"objectType":   objectType,
				"search":       search,
				"sortField":    sortField,
				"sortOrder":    sortOrder,
				"hasReference": hasReference,
				"fields":       fields,
			},
			"pagination": map[string]interface{}{
				"currentPage": result.Page,
				"perPage":     result.PerPage,
				"total":       result.Total,
			},
			"metadata": map[string]interface{}{
				"tool":         "kibana_search_saved_objects_advanced",
				"optimizedFor": "finding specific objects",
			},
		}

		return marshalOptimizedResponse(response, "kibana_search_saved_objects_advanced")
	}
}

// HandleGetDashboardDetailAdvanced handles getting advanced dashboard details
func HandleGetDashboardDetailAdvanced(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dashboardID, err := requireStringParam(request, "dashboard_id")
		if err != nil {
			return nil, err
		}

		includePanels := getOptionalBoolParam(request, "include_panels")
		if includePanels == nil {
			defaultPanels := true
			includePanels = &defaultPanels
		}

		includeUIState := getOptionalBoolParam(request, "include_ui_state")
		if includeUIState == nil {
			defaultUIState := false
			includeUIState = &defaultUIState
		}

		includeTimeOptions := getOptionalBoolParam(request, "include_time_options")
		if includeTimeOptions == nil {
			defaultTimeOptions := true
			includeTimeOptions = &defaultTimeOptions
		}

		outputFormat := getOptionalStringParam(request, "output_format")
		if outputFormat == "" {
			outputFormat = "structured"
		}

		logrus.WithFields(logrus.Fields{
			"tool":               "kibana_get_dashboard_detail_advanced",
			"dashboardID":        dashboardID,
			"includePanels":      *includePanels,
			"includeUIState":     *includeUIState,
			"includeTimeOptions": *includeTimeOptions,
			"outputFormat":       outputFormat,
		}).Debug("Handler invoked")

		detail, err := c.GetDashboardDetailAdvanced(ctx, dashboardID, *includePanels, *includeUIState, *includeTimeOptions, outputFormat)
		if err != nil {
			return nil, fmt.Errorf("failed to get dashboard detail advanced: %w", err)
		}

		response := map[string]interface{}{
			"dashboardDetail": detail,
			"metadata": map[string]interface{}{
				"tool":         "kibana_get_dashboard_detail_advanced",
				"dashboardID":  dashboardID,
				"outputFormat": outputFormat,
				"optimizedFor": "comprehensive analysis",
			},
		}

		return marshalOptimizedResponse(response, "kibana_get_dashboard_detail_advanced")
	}
}

// HandleGetHealthSummary handles getting Kibana health summary
func HandleGetHealthSummary(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		level := getOptionalStringParam(request, "level")
		if level == "" {
			level = "basic"
		}

		includeSavedObjects := getOptionalBoolParam(request, "include_saved_objects")
		if includeSavedObjects == nil {
			defaultSavedObjects := false
			includeSavedObjects = &defaultSavedObjects
		}

		logrus.WithFields(logrus.Fields{
			"tool":                "kibana_health_summary",
			"level":               level,
			"includeSavedObjects": *includeSavedObjects,
		}).Debug("Handler invoked")

		health, err := c.GetHealthSummary(ctx, level, *includeSavedObjects)
		if err != nil {
			return nil, fmt.Errorf("failed to get health summary: %w", err)
		}

		response := map[string]interface{}{
			"health": health,
			"metadata": map[string]interface{}{
				"tool":                "kibana_health_summary",
				"level":               level,
				"includeSavedObjects": *includeSavedObjects,
				"optimizedFor":        "monitoring and LLM efficiency",
			},
		}

		return marshalOptimizedResponse(response, "kibana_health_summary")
	}
}

// requireStringParam helper validates required string parameter
func requireStringParam(request mcp.CallToolRequest, param string) (string, error) {
	value, ok := request.GetArguments()[param].(string)
	if !ok || value == "" {
		return "", fmt.Errorf("missing required parameter: %s", param)
	}
	return value, nil
}

// HandleGetVisualization handles specific Kibana visualization retrieval requests.
func HandleGetVisualization(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get visualization handler")

		// Extract parameters
		args := req.Params.Arguments
		if args == nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("No arguments provided"),
				},
			}, nil
		}

		// Get visualization ID parameter
		visualizationID, ok := req.GetArguments()["visualization_id"].(string)
		if !ok || visualizationID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Visualization ID parameter is required"),
				},
			}, nil
		}

		// Get visualization
		visualization, err := c.GetVisualization(ctx, visualizationID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get visualization: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(visualization)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format visualization: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleGetIndexPattern handles specific Kibana index pattern retrieval requests.
func HandleGetIndexPattern(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get index pattern handler")

		// Extract parameters
		args := req.Params.Arguments
		if args == nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("No arguments provided"),
				},
			}, nil
		}

		// Get index pattern ID parameter
		indexPatternID, ok := req.GetArguments()["index_pattern_id"].(string)
		if !ok || indexPatternID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Index pattern ID parameter is required"),
				},
			}, nil
		}

		// Get index pattern
		indexPattern, err := c.GetIndexPattern(ctx, indexPatternID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get index pattern: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(indexPattern)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format index pattern: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleGetSavedSearches handles Kibana saved searches retrieval requests.
func HandleGetSavedSearches(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get saved searches handler")

		// Get saved searches
		savedSearches, err := c.GetSavedSearches(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get saved searches: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(savedSearches)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format saved searches: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleGetSavedSearch handles specific Kibana saved search retrieval requests.
func HandleGetSavedSearch(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get saved search handler")

		// Extract parameters
		args := req.Params.Arguments
		if args == nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("No arguments provided"),
				},
			}, nil
		}

		// Get search ID parameter
		searchID, ok := req.GetArguments()["search_id"].(string)
		if !ok || searchID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Search ID parameter is required"),
				},
			}, nil
		}

		// Get saved search
		savedSearch, err := c.GetSavedSearch(ctx, searchID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get saved search: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(savedSearch)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format saved search: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleGetKibanaStatus handles Kibana status retrieval requests.
func HandleGetKibanaStatus(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get status handler")

		// Get Kibana status
		status, err := c.GetKibanaStatus(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get Kibana status: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(status)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format Kibana status: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// ============ Write Operations: Spaces ============

// HandleCreateSpace handles creating a new Kibana space
func HandleCreateSpace(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana create space handler")

		spaceID, _ := req.GetArguments()["id"].(string)
		name, _ := req.GetArguments()["name"].(string)
		description, _ := req.GetArguments()["description"].(string)
		color, _ := req.GetArguments()["color"].(string)
		initials, _ := req.GetArguments()["initials"].(string)

		var disabledFeatures []string
		if feats, ok := req.GetArguments()["disabledFeatures"].([]interface{}); ok {
			for _, f := range feats {
				if s, ok := f.(string); ok {
					disabledFeatures = append(disabledFeatures, s)
				}
			}
		}

		space := client.Space{
			ID:               spaceID,
			Name:             name,
			Description:      description,
			Color:            color,
			Initials:         initials,
			DisabledFeatures: disabledFeatures,
		}

		created, err := c.CreateSpace(ctx, space)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to create space: %v", err)),
				},
			}, nil
		}

		resultJSON, err := json.Marshal(created)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format response: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleUpdateSpace handles updating an existing Kibana space
func HandleUpdateSpace(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana update space handler")

		spaceID, _ := req.GetArguments()["space_id"].(string)
		name, _ := req.GetArguments()["name"].(string)
		description, _ := req.GetArguments()["description"].(string)
		color, _ := req.GetArguments()["color"].(string)
		initials, _ := req.GetArguments()["initials"].(string)

		var disabledFeatures []string
		if feats, ok := req.GetArguments()["disabledFeatures"].([]interface{}); ok {
			for _, f := range feats {
				if s, ok := f.(string); ok {
					disabledFeatures = append(disabledFeatures, s)
				}
			}
		}

		space := client.Space{
			Name:             name,
			Description:      description,
			Color:            color,
			Initials:         initials,
			DisabledFeatures: disabledFeatures,
		}

		updated, err := c.UpdateSpace(ctx, spaceID, space)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to update space: %v", err)),
				},
			}, nil
		}

		resultJSON, err := json.Marshal(updated)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format response: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleDeleteSpace handles deleting a Kibana space
func HandleDeleteSpace(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana delete space handler")

		spaceID, _ := req.GetArguments()["space_id"].(string)
		force := false
		if f, ok := req.GetArguments()["force"].(bool); ok {
			force = f
		}

		if spaceID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("space_id is required"),
				},
			}, nil
		}

		err := c.DeleteSpace(ctx, spaceID, force)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to delete space: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully deleted space: %s", spaceID)),
			},
		}, nil
	}
}

// ============ Write Operations: Index Patterns ============

// HandleCreateIndexPattern handles creating a new index pattern
func HandleCreateIndexPattern(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana create index pattern handler")

		title, _ := req.GetArguments()["title"].(string)
		timeField, _ := req.GetArguments()["timeField"].(string)

		if title == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("title is required"),
				},
			}, nil
		}

		pattern, err := c.CreateIndexPattern(ctx, title, timeField)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to create index pattern: %v", err)),
				},
			}, nil
		}

		resultJSON, err := json.Marshal(pattern)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format response: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleUpdateIndexPattern handles updating an index pattern
func HandleUpdateIndexPattern(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana update index pattern handler")

		patternID, _ := req.GetArguments()["index_pattern_id"].(string)
		title, _ := req.GetArguments()["title"].(string)
		timeField, _ := req.GetArguments()["timeField"].(string)

		if patternID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("index_pattern_id is required"),
				},
			}, nil
		}

		pattern, err := c.UpdateIndexPattern(ctx, patternID, title, timeField)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to update index pattern: %v", err)),
				},
			}, nil
		}

		resultJSON, err := json.Marshal(pattern)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format response: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleDeleteIndexPattern handles deleting an index pattern
func HandleDeleteIndexPattern(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana delete index pattern handler")

		patternID, _ := req.GetArguments()["index_pattern_id"].(string)

		if patternID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("index_pattern_id is required"),
				},
			}, nil
		}

		err := c.DeleteIndexPattern(ctx, patternID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to delete index pattern: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully deleted index pattern: %s", patternID)),
			},
		}, nil
	}
}

// HandleSetDefaultIndexPattern handles setting the default index pattern
func HandleSetDefaultIndexPattern(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana set default index pattern handler")

		patternID, _ := req.GetArguments()["index_pattern_id"].(string)

		if patternID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("index_pattern_id is required"),
				},
			}, nil
		}

		err := c.SetDefaultIndexPattern(ctx, patternID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to set default index pattern: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully set default index pattern: %s", patternID)),
			},
		}, nil
	}
}

// HandleRefreshIndexPatternFields handles refreshing index pattern fields
func HandleRefreshIndexPatternFields(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana refresh index pattern fields handler")

		patternID, _ := req.GetArguments()["index_pattern_id"].(string)

		if patternID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("index_pattern_id is required"),
				},
			}, nil
		}

		err := c.RefreshIndexPatternFields(ctx, patternID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to refresh index pattern fields: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully refreshed fields for index pattern: %s", patternID)),
			},
		}, nil
	}
}

// ============ Write Operations: Dashboards ============

// HandleCreateDashboard handles creating a new dashboard
func HandleCreateDashboard(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana create dashboard handler")

		title, _ := req.GetArguments()["title"].(string)
		description, _ := req.GetArguments()["description"].(string)

		timeRestore := true
		if tr, ok := req.GetArguments()["timeRestore"].(bool); ok {
			timeRestore = tr
		}

		timeFrom, _ := req.GetArguments()["timeFrom"].(string)
		timeTo, _ := req.GetArguments()["timeTo"].(string)

		var refreshInterval map[string]interface{}
		if ri, ok := req.GetArguments()["refreshInterval"].(map[string]interface{}); ok {
			refreshInterval = ri
		}

		if title == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("title is required"),
				},
			}, nil
		}

		dashboard, err := c.CreateDashboard(ctx, title, description, timeRestore, timeFrom, timeTo, refreshInterval)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to create dashboard: %v", err)),
				},
			}, nil
		}

		resultJSON, err := json.Marshal(dashboard)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format response: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleUpdateDashboard handles updating a dashboard
func HandleUpdateDashboard(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana update dashboard handler")

		dashboardID, _ := req.GetArguments()["dashboard_id"].(string)
		title, _ := req.GetArguments()["title"].(string)
		description, _ := req.GetArguments()["description"].(string)
		panelsJSON, _ := req.GetArguments()["panelsJSON"].(string)
		timeFrom, _ := req.GetArguments()["timeFrom"].(string)
		timeTo, _ := req.GetArguments()["timeTo"].(string)

		version := 0
		if v, ok := req.GetArguments()["version"].(float64); ok {
			version = int(v)
		}

		if dashboardID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("dashboard_id is required"),
				},
			}, nil
		}

		dashboard, err := c.UpdateDashboard(ctx, dashboardID, title, description, panelsJSON, timeFrom, timeTo, version)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to update dashboard: %v", err)),
				},
			}, nil
		}

		resultJSON, err := json.Marshal(dashboard)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format response: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleDeleteDashboard handles deleting a dashboard
func HandleDeleteDashboard(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana delete dashboard handler")

		dashboardID, _ := req.GetArguments()["dashboard_id"].(string)

		if dashboardID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("dashboard_id is required"),
				},
			}, nil
		}

		err := c.DeleteDashboard(ctx, dashboardID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to delete dashboard: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully deleted dashboard: %s", dashboardID)),
			},
		}, nil
	}
}

// HandleCloneDashboard handles cloning a dashboard
func HandleCloneDashboard(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana clone dashboard handler")

		dashboardID, _ := req.GetArguments()["dashboard_id"].(string)
		newTitle, _ := req.GetArguments()["new_title"].(string)

		if dashboardID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("dashboard_id is required"),
				},
			}, nil
		}

		if newTitle == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("new_title is required"),
				},
			}, nil
		}

		dashboard, err := c.CloneDashboard(ctx, dashboardID, newTitle)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to clone dashboard: %v", err)),
				},
			}, nil
		}

		resultJSON, err := json.Marshal(dashboard)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format response: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// ============ Write Operations: Visualizations ============

// HandleCreateVisualization handles creating a new visualization
func HandleCreateVisualization(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana create visualization handler")

		title, _ := req.GetArguments()["title"].(string)
		description, _ := req.GetArguments()["description"].(string)
		savedSearchRefName, _ := req.GetArguments()["savedSearchRefName"].(string)

		var visState map[string]interface{}
		if vs, ok := req.GetArguments()["visState"].(map[string]interface{}); ok {
			visState = vs
		}

		if title == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("title is required"),
				},
			}, nil
		}

		visualization, err := c.CreateVisualization(ctx, title, visState, description, savedSearchRefName)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to create visualization: %v", err)),
				},
			}, nil
		}

		resultJSON, err := json.Marshal(visualization)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format response: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleUpdateVisualization handles updating a visualization
func HandleUpdateVisualization(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana update visualization handler")

		visualizationID, _ := req.GetArguments()["visualization_id"].(string)
		title, _ := req.GetArguments()["title"].(string)
		description, _ := req.GetArguments()["description"].(string)

		var visState map[string]interface{}
		if vs, ok := req.GetArguments()["visState"].(map[string]interface{}); ok {
			visState = vs
		}

		version := 0
		if v, ok := req.GetArguments()["version"].(float64); ok {
			version = int(v)
		}

		if visualizationID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("visualization_id is required"),
				},
			}, nil
		}

		visualization, err := c.UpdateVisualization(ctx, visualizationID, title, visState, description, version)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to update visualization: %v", err)),
				},
			}, nil
		}

		resultJSON, err := json.Marshal(visualization)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format response: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleDeleteVisualization handles deleting a visualization
func HandleDeleteVisualization(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana delete visualization handler")

		visualizationID, _ := req.GetArguments()["visualization_id"].(string)

		if visualizationID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("visualization_id is required"),
				},
			}, nil
		}

		err := c.DeleteVisualization(ctx, visualizationID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to delete visualization: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully deleted visualization: %s", visualizationID)),
			},
		}, nil
	}
}

// HandleCloneVisualization handles cloning a visualization
func HandleCloneVisualization(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana clone visualization handler")

		visualizationID, _ := req.GetArguments()["visualization_id"].(string)
		newTitle, _ := req.GetArguments()["new_title"].(string)

		if visualizationID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("visualization_id is required"),
				},
			}, nil
		}

		if newTitle == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("new_title is required"),
				},
			}, nil
		}

		visualization, err := c.CloneVisualization(ctx, visualizationID, newTitle)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to clone visualization: %v", err)),
				},
			}, nil
		}

		resultJSON, err := json.Marshal(visualization)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format response: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// ============ Write Operations: Saved Objects (Generic) ============

// HandleCreateSavedObject handles creating a generic saved object
func HandleCreateSavedObject(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana create saved object handler")

		objectType, _ := req.GetArguments()["type"].(string)
		initialObjectType, _ := req.GetArguments()["initialObjectType"].(string)

		var attributes map[string]interface{}
		if attrs, ok := req.GetArguments()["attributes"].(map[string]interface{}); ok {
			attributes = attrs
		}

		var references []client.Reference
		if refs, ok := req.GetArguments()["references"].([]interface{}); ok {
			for _, r := range refs {
				if refMap, ok := r.(map[string]interface{}); ok {
					references = append(references, client.Reference{
						Name: getStringFieldFromMap(refMap, "name"),
						Type: getStringFieldFromMap(refMap, "type"),
						ID:   getStringFieldFromMap(refMap, "id"),
					})
				}
			}
		}

		if objectType == "" || len(attributes) == 0 {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("type and attributes are required"),
				},
			}, nil
		}

		obj, err := c.CreateSavedObject(ctx, objectType, attributes, references, initialObjectType)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to create saved object: %v", err)),
				},
			}, nil
		}

		resultJSON, err := json.Marshal(obj)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format response: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleUpdateSavedObject handles updating a saved object
func HandleUpdateSavedObject(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana update saved object handler")

		objectType, _ := req.GetArguments()["type"].(string)
		objectID, _ := req.GetArguments()["object_id"].(string)
		version, _ := req.GetArguments()["version"].(string)

		var attributes map[string]interface{}
		if attrs, ok := req.GetArguments()["attributes"].(map[string]interface{}); ok {
			attributes = attrs
		}

		var references []client.Reference
		if refs, ok := req.GetArguments()["references"].([]interface{}); ok {
			for _, r := range refs {
				if refMap, ok := r.(map[string]interface{}); ok {
					references = append(references, client.Reference{
						Name: getStringFieldFromMap(refMap, "name"),
						Type: getStringFieldFromMap(refMap, "type"),
						ID:   getStringFieldFromMap(refMap, "id"),
					})
				}
			}
		}

		if objectType == "" || objectID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("type and object_id are required"),
				},
			}, nil
		}

		obj, err := c.UpdateSavedObject(ctx, objectType, objectID, attributes, references, version)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to update saved object: %v", err)),
				},
			}, nil
		}

		resultJSON, err := json.Marshal(obj)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format response: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleDeleteSavedObject handles deleting a saved object
func HandleDeleteSavedObject(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana delete saved object handler")

		objectType, _ := req.GetArguments()["type"].(string)
		objectID, _ := req.GetArguments()["object_id"].(string)

		force := false
		if f, ok := req.GetArguments()["force"].(bool); ok {
			force = f
		}

		if objectType == "" || objectID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("type and object_id are required"),
				},
			}, nil
		}

		err := c.DeleteSavedObject(ctx, objectType, objectID, force)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to delete saved object: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully deleted saved object: %s/%s", objectType, objectID)),
			},
		}, nil
	}
}

// HandleBulkDeleteSavedObjects handles bulk deleting saved objects
func HandleBulkDeleteSavedObjects(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana bulk delete saved objects handler")

		var objects []client.SavedObject
		if objs, ok := req.GetArguments()["objects"].([]interface{}); ok {
			for _, o := range objs {
				if objMap, ok := o.(map[string]interface{}); ok {
					objects = append(objects, client.SavedObject{
						Type: getStringFieldFromMap(objMap, "type"),
						ID:   getStringFieldFromMap(objMap, "id"),
					})
				}
			}
		}

		if len(objects) == 0 {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("objects array is required"),
				},
			}, nil
		}

		err := c.BulkDeleteSavedObjects(ctx, objects)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to bulk delete saved objects: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully deleted %d saved objects", len(objects))),
			},
		}, nil
	}
}

// HandleExportSavedObjects handles exporting saved objects
func HandleExportSavedObjects(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana export saved objects handler")

		var objects []client.SavedObject
		if objs, ok := req.GetArguments()["objects"].([]interface{}); ok {
			for _, o := range objs {
				if objMap, ok := o.(map[string]interface{}); ok {
					objects = append(objects, client.SavedObject{
						Type: getStringFieldFromMap(objMap, "type"),
						ID:   getStringFieldFromMap(objMap, "id"),
					})
				}
			}
		}

		includeReferences := true
		if ir, ok := req.GetArguments()["includeReferences"].(bool); ok {
			includeReferences = ir
		}

		if len(objects) == 0 {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("objects array is required"),
				},
			}, nil
		}

		data, err := c.ExportSavedObjects(ctx, objects, includeReferences)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to export saved objects: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(data)),
			},
		}, nil
	}
}

// HandleImportSavedObjects handles importing saved objects
func HandleImportSavedObjects(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana import saved objects handler")

		fileContent, _ := req.GetArguments()["file"].(string)

		createNewCopies := false
		if cnc, ok := req.GetArguments()["createNewCopies"].(bool); ok {
			createNewCopies = cnc
		}

		if fileContent == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("file is required"),
				},
			}, nil
		}

		err := c.ImportSavedObjects(ctx, fileContent, createNewCopies)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to import saved objects: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Successfully imported saved objects"),
			},
		}, nil
	}
}

// Helper function to get string field from map
func getStringFieldFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

// ============ Analysis & Discovery Handlers ============

// HandleQueryLogs handles log search requests.
func HandleQueryLogs(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		indexPattern := ""
		if v, ok := req.GetArguments()["indexPattern"]; ok {
			if s, ok := v.(string); ok {
				indexPattern = s
			}
		}

		query := "*"
		if v, ok := req.GetArguments()["query"]; ok {
			if s, ok := v.(string); ok {
				query = s
			}
		}

		size := 20
		if v, ok := req.GetArguments()["size"]; ok {
			if f, ok := v.(float64); ok {
				size = int(f)
			}
		}

		sortBy := "@timestamp"
		if v, ok := req.GetArguments()["sortBy"]; ok {
			if s, ok := v.(string); ok {
				sortBy = s
			}
		}

		sortOrder := "desc"
		if v, ok := req.GetArguments()["sortOrder"]; ok {
			if s, ok := v.(string); ok {
				sortOrder = s
			}
		}

		logrus.WithFields(logrus.Fields{
			"indexPattern": indexPattern,
			"query":        query,
			"size":         size,
		}).Debug("Executing Kibana log query handler")

		result, err := c.QueryLogs(ctx, indexPattern, query, size, sortBy, sortOrder)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to query logs: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(result)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format log results: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleGetCanvasWorkpads handles Canvas workpad retrieval requests.
func HandleGetCanvasWorkpads(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get Canvas workpads handler")

		workpads, err := c.GetCanvasWorkpads(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get Canvas workpads: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(workpads)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format workpads: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleGetLensObjects handles Lens visualization retrieval requests.
func HandleGetLensObjects(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get Lens objects handler")

		lenses, err := c.GetLensObjects(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get Lens objects: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(lenses)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format Lens objects: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleGetMaps handles Kibana Maps retrieval requests.
func HandleGetMaps(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get Maps handler")

		maps, err := c.GetMaps(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get Maps: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(maps)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format Maps: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleGetKibanaAlerts handles Kibana alerting rules retrieval requests.
func HandleGetKibanaAlerts(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get alerts handler")

		alerts, err := c.GetAlerts(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get Kibana alerts: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(alerts)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format alerts: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleGetIndexPatternFields handles index pattern fields retrieval requests.
func HandleGetIndexPatternFields(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		patternID := ""
		if v, ok := req.GetArguments()["patternID"]; ok {
			if s, ok := v.(string); ok {
				patternID = s
			}
		}

		if patternID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("patternID is required"),
				},
			}, nil
		}

		logrus.WithField("patternID", patternID).Debug("Executing Kibana get index pattern fields handler")

		fields, err := c.GetIndexPatternFields(ctx, patternID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get index pattern fields: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(fields)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format fields: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// ============ Alert Rules Handlers ============

// HandleGetAlertRules handles listing alert rules.
func HandleGetAlertRules(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		page := getOptionalIntParam(req, "page", 1)
		perPage := getOptionalIntParam(req, "per_page", 20)
		filter := getOptionalStringParam(req, "filter")
		var enabled *bool
		if e, exists := req.GetArguments()["enabled"]; exists {
			if eBool, ok := e.(bool); ok {
				enabled = &eBool
			}
		}

		logrus.WithFields(logrus.Fields{
			"page":    page,
			"perPage": perPage,
			"filter":  filter,
		}).Debug("Executing Kibana get alert rules handler")

		rules, err := c.GetAlertRules(ctx, page, perPage, filter, enabled)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get alert rules: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(rules)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format alert rules: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleGetAlertRule handles getting a specific alert rule.
func HandleGetAlertRule(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ruleID, err := requireStringParam(req, "rule_id")
		if err != nil {
			return nil, err
		}

		logrus.WithField("rule_id", ruleID).Debug("Executing Kibana get alert rule handler")

		rule, err := c.GetAlertRule(ctx, ruleID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get alert rule: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(rule)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format alert rule: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleCreateAlertRule handles creating a new alert rule.
func HandleCreateAlertRule(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := getOptionalStringParam(req, "name")
		alertTypeID := getOptionalStringParam(req, "alertTypeId")

		var schedule, params map[string]interface{}
		if s, ok := req.GetArguments()["schedule"].(map[string]interface{}); ok {
			schedule = s
		}
		if p, ok := req.GetArguments()["params"].(map[string]interface{}); ok {
			params = p
		}

		var actions []map[string]interface{}
		if a, ok := req.GetArguments()["actions"].([]interface{}); ok {
			for _, item := range a {
				if actionMap, ok := item.(map[string]interface{}); ok {
					actions = append(actions, actionMap)
				}
			}
		}

		var tags []string
		if t, ok := req.GetArguments()["tags"].([]interface{}); ok {
			for _, tag := range t {
				if tagStr, ok := tag.(string); ok {
					tags = append(tags, tagStr)
				}
			}
		}

		if name == "" || alertTypeID == "" || schedule == nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("name, alertTypeId, and schedule are required"),
				},
			}, nil
		}

		logrus.WithField("name", name).Debug("Executing Kibana create alert rule handler")

		rule, err := c.CreateAlertRule(ctx, name, alertTypeID, schedule, params, actions, tags)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to create alert rule: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(rule)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format response: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleUpdateAlertRule handles updating an existing alert rule.
func HandleUpdateAlertRule(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ruleID, err := requireStringParam(req, "rule_id")
		if err != nil {
			return nil, err
		}

		name := getOptionalStringParam(req, "name")
		schedule := getOptionalStringParam(req, "schedule")

		var params, actions map[string]interface{}
		if p, ok := req.GetArguments()["params"].(map[string]interface{}); ok {
			params = p
		}
		if a, ok := req.GetArguments()["actions"].(map[string]interface{}); ok {
			actions = a
		}

		var tags []string
		if t, ok := req.GetArguments()["tags"].([]interface{}); ok {
			for _, tag := range t {
				if tagStr, ok := tag.(string); ok {
					tags = append(tags, tagStr)
				}
			}
		}

		logrus.WithField("rule_id", ruleID).Debug("Executing Kibana update alert rule handler")

		rule, err := c.UpdateAlertRule(ctx, ruleID, name, schedule, params, actions, tags)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to update alert rule: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(rule)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format response: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleDeleteAlertRule handles deleting an alert rule.
func HandleDeleteAlertRule(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ruleID, err := requireStringParam(req, "rule_id")
		if err != nil {
			return nil, err
		}

		logrus.WithField("rule_id", ruleID).Debug("Executing Kibana delete alert rule handler")

		err = c.DeleteAlertRule(ctx, ruleID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to delete alert rule: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully deleted alert rule: %s", ruleID)),
			},
		}, nil
	}
}

// HandleEnableAlertRule handles enabling an alert rule.
func HandleEnableAlertRule(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ruleID, err := requireStringParam(req, "rule_id")
		if err != nil {
			return nil, err
		}

		logrus.WithField("rule_id", ruleID).Debug("Executing Kibana enable alert rule handler")

		err = c.EnableAlertRule(ctx, ruleID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to enable alert rule: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully enabled alert rule: %s", ruleID)),
			},
		}, nil
	}
}

// HandleDisableAlertRule handles disabling an alert rule.
func HandleDisableAlertRule(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ruleID, err := requireStringParam(req, "rule_id")
		if err != nil {
			return nil, err
		}

		logrus.WithField("rule_id", ruleID).Debug("Executing Kibana disable alert rule handler")

		err = c.DisableAlertRule(ctx, ruleID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to disable alert rule: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully disabled alert rule: %s", ruleID)),
			},
		}, nil
	}
}

// HandleMuteAlertRule handles muting an alert rule.
func HandleMuteAlertRule(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ruleID, err := requireStringParam(req, "rule_id")
		if err != nil {
			return nil, err
		}

		duration := getOptionalStringParam(req, "duration")
		if duration == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("duration is required"),
				},
			}, nil
		}

		logrus.WithFields(logrus.Fields{
			"rule_id":  ruleID,
			"duration": duration,
		}).Debug("Executing Kibana mute alert rule handler")

		err = c.MuteAlertRule(ctx, ruleID, duration)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to mute alert rule: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully muted alert rule: %s for %s", ruleID, duration)),
			},
		}, nil
	}
}

// HandleUnmuteAlertRule handles unmuting an alert rule.
func HandleUnmuteAlertRule(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ruleID, err := requireStringParam(req, "rule_id")
		if err != nil {
			return nil, err
		}

		logrus.WithField("rule_id", ruleID).Debug("Executing Kibana unmute alert rule handler")

		err = c.UnmuteAlertRule(ctx, ruleID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to unmute alert rule: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully unmuted alert rule: %s", ruleID)),
			},
		}, nil
	}
}

// HandleGetAlertRuleTypes handles listing available alert rule types.
func HandleGetAlertRuleTypes(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get alert rule types handler")

		ruleTypes, err := c.GetAlertRuleTypes(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get alert rule types: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(ruleTypes)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format alert rule types: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleGetAlertRuleHistory handles getting alert rule execution history.
func HandleGetAlertRuleHistory(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ruleID, err := requireStringParam(req, "rule_id")
		if err != nil {
			return nil, err
		}

		page := getOptionalIntParam(req, "page", 1)
		perPage := getOptionalIntParam(req, "per_page", 20)

		logrus.WithFields(logrus.Fields{
			"rule_id": ruleID,
			"page":    page,
			"perPage": perPage,
		}).Debug("Executing Kibana get alert rule history handler")

		history, err := c.GetAlertRuleHistory(ctx, ruleID, page, perPage)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get alert rule history: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(history)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format alert rule history: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// ============ Connectors Handlers ============

// HandleGetConnectors handles listing connectors.
func HandleGetConnectors(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		page := getOptionalIntParam(req, "page", 1)
		perPage := getOptionalIntParam(req, "per_page", 20)

		logrus.WithFields(logrus.Fields{
			"page":    page,
			"perPage": perPage,
		}).Debug("Executing Kibana get connectors handler")

		connectors, err := c.GetConnectors(ctx, page, perPage)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get connectors: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(connectors)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format connectors: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleGetConnector handles getting a specific connector.
func HandleGetConnector(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		connectorID, err := requireStringParam(req, "connector_id")
		if err != nil {
			return nil, err
		}

		logrus.WithField("connector_id", connectorID).Debug("Executing Kibana get connector handler")

		connector, err := c.GetConnector(ctx, connectorID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get connector: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(connector)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format connector: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleCreateConnector handles creating a new connector.
func HandleCreateConnector(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := getOptionalStringParam(req, "name")
		connectorTypeID := getOptionalStringParam(req, "connectorTypeId")

		var config, secrets map[string]interface{}
		if c, ok := req.GetArguments()["config"].(map[string]interface{}); ok {
			config = c
		}
		if s, ok := req.GetArguments()["secrets"].(map[string]interface{}); ok {
			secrets = s
		}

		if name == "" || connectorTypeID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("name and connectorTypeId are required"),
				},
			}, nil
		}

		logrus.WithField("name", name).Debug("Executing Kibana create connector handler")

		connector, err := c.CreateConnector(ctx, name, connectorTypeID, config, secrets)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to create connector: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(connector)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format response: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleUpdateConnector handles updating an existing connector.
func HandleUpdateConnector(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		connectorID, err := requireStringParam(req, "connector_id")
		if err != nil {
			return nil, err
		}

		name := getOptionalStringParam(req, "name")

		var config, secrets map[string]interface{}
		if c, ok := req.GetArguments()["config"].(map[string]interface{}); ok {
			config = c
		}
		if s, ok := req.GetArguments()["secrets"].(map[string]interface{}); ok {
			secrets = s
		}

		logrus.WithField("connector_id", connectorID).Debug("Executing Kibana update connector handler")

		connector, err := c.UpdateConnector(ctx, connectorID, name, config, secrets)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to update connector: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(connector)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format response: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleDeleteConnector handles deleting a connector.
func HandleDeleteConnector(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		connectorID, err := requireStringParam(req, "connector_id")
		if err != nil {
			return nil, err
		}

		logrus.WithField("connector_id", connectorID).Debug("Executing Kibana delete connector handler")

		err = c.DeleteConnector(ctx, connectorID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to delete connector: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully deleted connector: %s", connectorID)),
			},
		}, nil
	}
}

// HandleTestConnector handles testing a connector.
func HandleTestConnector(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		connectorID, err := requireStringParam(req, "connector_id")
		if err != nil {
			return nil, err
		}

		var body map[string]interface{}
		if b, ok := req.GetArguments()["body"].(map[string]interface{}); ok {
			body = b
		}

		logrus.WithField("connector_id", connectorID).Debug("Executing Kibana test connector handler")

		err = c.TestConnector(ctx, connectorID, body)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to test connector: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully tested connector: %s", connectorID)),
			},
		}, nil
	}
}

// HandleGetConnectorTypes handles listing available connector types.
func HandleGetConnectorTypes(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get connector types handler")

		connectorTypes, err := c.GetConnectorTypes(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get connector types: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(connectorTypes)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format connector types: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// ============ Data Views Handlers ============

// HandleGetDataViews handles listing data views.
func HandleGetDataViews(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		page := getOptionalIntParam(req, "page", 1)
		perPage := getOptionalIntParam(req, "per_page", 20)

		logrus.WithFields(logrus.Fields{
			"page":    page,
			"perPage": perPage,
		}).Debug("Executing Kibana get data views handler")

		dataViews, err := c.GetDataViews(ctx, page, perPage)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get data views: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(dataViews)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format data views: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleGetDataView handles getting a specific data view.
func HandleGetDataView(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dataViewID, err := requireStringParam(req, "data_view_id")
		if err != nil {
			return nil, err
		}

		logrus.WithField("data_view_id", dataViewID).Debug("Executing Kibana get data view handler")

		dataView, err := c.GetDataView(ctx, dataViewID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get data view: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(dataView)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format data view: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleCreateDataView handles creating a new data view.
func HandleCreateDataView(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title := getOptionalStringParam(req, "title")
		name := getOptionalStringParam(req, "name")
		timeField := getOptionalStringParam(req, "timeField")
		allowNoIndex := false
		if ani, ok := req.GetArguments()["allowNoIndex"].(bool); ok {
			allowNoIndex = ani
		}

		if title == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("title is required"),
				},
			}, nil
		}

		logrus.WithField("title", title).Debug("Executing Kibana create data view handler")

		dataView, err := c.CreateDataView(ctx, title, name, timeField, nil, nil, allowNoIndex)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to create data view: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(dataView)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format response: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleUpdateDataView handles updating an existing data view.
func HandleUpdateDataView(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dataViewID, err := requireStringParam(req, "data_view_id")
		if err != nil {
			return nil, err
		}

		title := getOptionalStringParam(req, "title")
		name := getOptionalStringParam(req, "name")
		timeField := getOptionalStringParam(req, "timeField")

		logrus.WithField("data_view_id", dataViewID).Debug("Executing Kibana update data view handler")

		dataView, err := c.UpdateDataView(ctx, dataViewID, title, name, timeField)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to update data view: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(dataView)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format response: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleDeleteDataView handles deleting a data view.
func HandleDeleteDataView(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dataViewID, err := requireStringParam(req, "data_view_id")
		if err != nil {
			return nil, err
		}

		logrus.WithField("data_view_id", dataViewID).Debug("Executing Kibana delete data view handler")

		err = c.DeleteDataView(ctx, dataViewID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to delete data view: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully deleted data view: %s", dataViewID)),
			},
		}, nil
	}
}

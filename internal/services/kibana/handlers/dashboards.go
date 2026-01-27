// Package handlers provides HTTP handlers for Kibana MCP operations.
// This file contains dashboard-related handlers.
package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/kibana/client"
)

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

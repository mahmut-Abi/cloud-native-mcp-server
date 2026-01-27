// Package handlers provides HTTP handlers for Kibana MCP operations.
// This file contains visualization-related handlers.
package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/kibana/client"
)

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

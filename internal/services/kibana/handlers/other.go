// Package handlers provides HTTP handlers for Kibana MCP operations.
// This file contains other handlers (Status, Logs, Canvas, Lens, Maps, Health).
package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/kibana/client"
)

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

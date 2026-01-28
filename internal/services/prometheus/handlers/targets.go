// Package handlers provides HTTP handlers for Prometheus MCP operations.
// It implements request handling for Prometheus queries, metrics, targets, and alerts.
package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/prometheus/client"
)

// HandleGetTargets handles Prometheus targets retrieval requests.
func HandleGetTargets(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Prometheus get targets handler")

		// Extract optional state filter
		var state string
		if req.GetArguments() != nil {
			if s, exists := req.GetArguments()["state"]; exists {
				if stateStr, ok := s.(string); ok {
					state = stateStr
				}
			}
		}

		// Get targets
		targets, err := c.GetTargets(ctx, state)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get targets: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(targets)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format targets: %v", err)),
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

// HandleGetTargetsSummary handles Prometheus targets summary requests (optimized version).
func HandleGetTargetsSummary(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		state := "any"
		if s, ok := req.GetArguments()["state"].(string); ok {
			state = s
		}

		logrus.WithField("state", state).Debug("Executing Prometheus targets summary handler")

		targets, err := c.GetTargets(ctx, state)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get targets: %v", err)),
				},
			}, nil
		}

		// Create summary (much smaller than full result)
		summary := make([]map[string]interface{}, 0, len(targets))
		for _, t := range targets {
			summary = append(summary, map[string]interface{}{
				"job":       t.Labels["job"],
				"instance":  t.Labels["instance"],
				"health":    t.Health,
				"scrapeUrl": t.ScrapeURL,
			})
		}

		result := map[string]interface{}{
			"count":   len(summary),
			"targets": summary,
		}

		resultJSON, err := marshalIndentJSON(result)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format result: %v", err)),
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
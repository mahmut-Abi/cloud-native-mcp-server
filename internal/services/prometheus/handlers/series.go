// Package handlers provides HTTP handlers for Prometheus MCP operations.
// It implements request handling for Prometheus queries, metrics, targets, and alerts.
package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/prometheus/client"
)

// HandleGetSeries handles Prometheus series retrieval requests.
func HandleGetSeries(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Prometheus get series handler")

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

		// Get matches parameter
		var matches []string
		if m, exists := req.GetArguments()["match"]; exists {
			switch v := m.(type) {
			case string:
				matches = []string{v}
			case []interface{}:
				for _, item := range v {
					if str, ok := item.(string); ok {
						matches = append(matches, str)
					}
				}
			}
		}

		if len(matches) == 0 {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("At least one match selector is required"),
				},
			}, nil
		}

		// Parse optional time range
		var start, end *time.Time
		if s, exists := req.GetArguments()["start"]; exists {
			if startStr, ok := s.(string); ok && startStr != "" {
				if parsed, err := time.Parse(time.RFC3339, startStr); err == nil {
					start = &parsed
				}
			}
		}
		if e, exists := req.GetArguments()["end"]; exists {
			if endStr, ok := e.(string); ok && endStr != "" {
				if parsed, err := time.Parse(time.RFC3339, endStr); err == nil {
					end = &parsed
				}
			}
		}

		// Get series
		series, err := c.GetSeries(ctx, matches, start, end)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get series: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(series)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format series: %v", err)),
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

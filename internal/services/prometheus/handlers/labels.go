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

// HandleGetLabelNames handles Prometheus label names retrieval requests.
func HandleGetLabelNames(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Prometheus get label names handler")

		// Parse optional time range
		var start, end *time.Time
		if req.GetArguments() != nil {
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
		}

		// Get label names
		labelNames, err := c.GetLabelNames(ctx, start, end)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get label names: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(labelNames)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format label names: %v", err)),
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

// HandleGetLabelValues handles Prometheus label values retrieval requests.
func HandleGetLabelValues(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Prometheus get label values handler")

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

		// Get label name parameter
		labelName, ok := req.GetArguments()["label"].(string)
		if !ok || labelName == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Label name parameter is required"),
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

		// Get label values
		labelValues, err := c.GetLabelValues(ctx, labelName, start, end)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get label values: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(labelValues)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format label values: %v", err)),
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
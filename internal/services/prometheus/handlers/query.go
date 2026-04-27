// Package handlers provides HTTP handlers for Prometheus MCP operations.
// It implements request handling for Prometheus queries, metrics, targets, and alerts.
package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	svccommon "github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/common"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/prometheus/client"
)

// HandleQuery handles Prometheus instant query requests.
func HandleQuery(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Prometheus query handler")

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

		// Get query parameter
		query, ok := svccommon.GetStringArg(req.GetArguments(), "query")
		if !ok || query == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Query parameter is required"),
				},
			}, nil
		}

		// Parse optional timestamp
		timestamp, err := svccommon.GetRFC3339TimeArg(req.GetArguments(), "time")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Execute query
		result, err := c.Query(ctx, query, timestamp)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to execute query: %v", err)),
				},
			}, nil
		}

		// Format result
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

// HandleQueryRange handles Prometheus range query requests.
func HandleQueryRange(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Prometheus range query handler")

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

		// Get required parameters
		query, ok := svccommon.GetStringArg(req.GetArguments(), "query")
		if !ok || query == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Query parameter is required"),
				},
			}, nil
		}

		startStr, ok := svccommon.GetStringArg(req.GetArguments(), "start")
		if !ok || startStr == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Start time parameter is required"),
				},
			}, nil
		}

		endStr, ok := svccommon.GetStringArg(req.GetArguments(), "end")
		if !ok || endStr == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("End time parameter is required"),
				},
			}, nil
		}

		step, ok := svccommon.GetStringArg(req.GetArguments(), "step")
		if !ok || step == "" {
			step = "15s" // Default step
		}

		// Parse timestamps
		start, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Invalid start time format: %v", err)),
				},
			}, nil
		}

		end, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Invalid end time format: %v", err)),
				},
			}, nil
		}

		// Execute range query
		result, err := c.QueryRange(ctx, query, start, end, step)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to execute range query: %v", err)),
				},
			}, nil
		}

		// Format result
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

// Package handlers provides HTTP handlers for Prometheus MCP operations.
// It implements request handling for Prometheus queries, metrics, targets, and alerts.
package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	svccommon "github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/common"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/prometheus/client"
)

// HandleGetLabelNames handles Prometheus label names retrieval requests.
func HandleGetLabelNames(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Prometheus get label names handler")

		// Parse optional time range
		start, err := svccommon.GetRFC3339TimeArg(req.GetArguments(), "start")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		end, err := svccommon.GetRFC3339TimeArg(req.GetArguments(), "end")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
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
		labelName, ok := svccommon.GetStringArg(req.GetArguments(), "label")
		if !ok || labelName == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Label name parameter is required"),
				},
			}, nil
		}

		// Parse optional time range
		start, err := svccommon.GetRFC3339TimeArg(req.GetArguments(), "start")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		end, err := svccommon.GetRFC3339TimeArg(req.GetArguments(), "end")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
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

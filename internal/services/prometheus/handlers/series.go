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
		matches, _, err := svccommon.GetStringSliceArg(req.GetArguments(), "match")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
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
		start, err := svccommon.GetRFC3339TimeArg(req.GetArguments(), "start")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		end, err := svccommon.GetRFC3339TimeArg(req.GetArguments(), "end")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
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

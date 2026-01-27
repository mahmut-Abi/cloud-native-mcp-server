// Package handlers provides tool handlers for OpenTelemetry operations.
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/opentelemetry/client"
)

// HandleGetMetrics handles the opentelemetry_get_metrics tool.
func HandleGetMetrics(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		var metricName *string
		if name, ok := args["metric_name"].(string); ok && name != "" {
			metricName = &name
		}

		var startTime, endTime *time.Time
		if startStr, ok := args["start_time"].(string); ok && startStr != "" {
			if t, err := time.Parse(time.RFC3339, startStr); err == nil {
				startTime = &t
			}
		}
		if endStr, ok := args["end_time"].(string); ok && endStr != "" {
			if t, err := time.Parse(time.RFC3339, endStr); err == nil {
				endTime = &t
			}
		}

		metrics, err := c.GetMetrics(ctx, metricName, startTime, endTime)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to retrieve metrics: %v", err)), nil
		}

		resultJSON, err := json.MarshalIndent(metrics, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize metrics: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// HandleQueryMetrics handles the opentelemetry_query_metrics tool.
func HandleQueryMetrics(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		query, ok := args["query"].(string)
		if !ok || query == "" {
			return mcp.NewToolResultError("Query parameter is required"), nil
		}

		var queryTime *time.Time
		if timeStr, ok := args["time"].(string); ok && timeStr != "" {
			if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
				queryTime = &t
			}
		}

		result, err := c.QueryMetrics(ctx, query, queryTime)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to query metrics: %v", err)), nil
		}

		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize query result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// HandleGetTraces handles the opentelemetry_get_traces tool.
func HandleGetTraces(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		var traceID, service *string
		if id, ok := args["trace_id"].(string); ok && id != "" {
			traceID = &id
		}
		if svc, ok := args["service"].(string); ok && svc != "" {
			service = &svc
		}

		var startTime, endTime *time.Time
		if startStr, ok := args["start_time"].(string); ok && startStr != "" {
			if t, err := time.Parse(time.RFC3339, startStr); err == nil {
				startTime = &t
			}
		}
		if endStr, ok := args["end_time"].(string); ok && endStr != "" {
			if t, err := time.Parse(time.RFC3339, endStr); err == nil {
				endTime = &t
			}
		}

		var limit *int
		if l, ok := args["limit"].(float64); ok && l > 0 {
			limitVal := int(l)
			limit = &limitVal
		}

		traces, err := c.GetTraces(ctx, traceID, service, startTime, endTime, limit)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to retrieve traces: %v", err)), nil
		}

		resultJSON, err := json.MarshalIndent(traces, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize traces: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// HandleQueryTraces handles the opentelemetry_query_traces tool.
func HandleQueryTraces(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		query, ok := args["query"].(string)
		if !ok || query == "" {
			return mcp.NewToolResultError("Query parameter is required"), nil
		}

		var service *string
		if svc, ok := args["service"].(string); ok && svc != "" {
			service = &svc
		}

		var startTime, endTime *time.Time
		if startStr, ok := args["start_time"].(string); ok && startStr != "" {
			if t, err := time.Parse(time.RFC3339, startStr); err == nil {
				startTime = &t
			}
		}
		if endStr, ok := args["end_time"].(string); ok && endStr != "" {
			if t, err := time.Parse(time.RFC3339, endStr); err == nil {
				endTime = &t
			}
		}

		var limit *int
		if l, ok := args["limit"].(float64); ok && l > 0 {
			limitVal := int(l)
			limit = &limitVal
		}

		result, err := c.QueryTraces(ctx, query, service, startTime, endTime, limit)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to query traces: %v", err)), nil
		}

		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize query result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// HandleGetLogs handles the opentelemetry_get_logs tool.
func HandleGetLogs(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		var service, level *string
		if svc, ok := args["service"].(string); ok && svc != "" {
			service = &svc
		}
		if lvl, ok := args["level"].(string); ok && lvl != "" {
			level = &lvl
		}

		var startTime, endTime *time.Time
		if startStr, ok := args["start_time"].(string); ok && startStr != "" {
			if t, err := time.Parse(time.RFC3339, startStr); err == nil {
				startTime = &t
			}
		}
		if endStr, ok := args["end_time"].(string); ok && endStr != "" {
			if t, err := time.Parse(time.RFC3339, endStr); err == nil {
				endTime = &t
			}
		}

		var limit *int
		if l, ok := args["limit"].(float64); ok && l > 0 {
			limitVal := int(l)
			limit = &limitVal
		}

		logs, err := c.GetLogs(ctx, service, level, startTime, endTime, limit)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to retrieve logs: %v", err)), nil
		}

		resultJSON, err := json.MarshalIndent(logs, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize logs: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// HandleQueryLogs handles the opentelemetry_query_logs tool.
func HandleQueryLogs(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		query, ok := args["query"].(string)
		if !ok || query == "" {
			return mcp.NewToolResultError("Query parameter is required"), nil
		}

		var service, level *string
		if svc, ok := args["service"].(string); ok && svc != "" {
			service = &svc
		}
		if lvl, ok := args["level"].(string); ok && lvl != "" {
			level = &lvl
		}

		var startTime, endTime *time.Time
		if startStr, ok := args["start_time"].(string); ok && startStr != "" {
			if t, err := time.Parse(time.RFC3339, startStr); err == nil {
				startTime = &t
			}
		}
		if endStr, ok := args["end_time"].(string); ok && endStr != "" {
			if t, err := time.Parse(time.RFC3339, endStr); err == nil {
				endTime = &t
			}
		}

		var limit *int
		if l, ok := args["limit"].(float64); ok && l > 0 {
			limitVal := int(l)
			limit = &limitVal
		}

		result, err := c.QueryLogs(ctx, query, service, level, startTime, endTime, limit)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to query logs: %v", err)), nil
		}

		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize query result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// HandleGetHealth handles the opentelemetry_get_health tool.
func HandleGetHealth(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		health, err := c.GetHealth(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to retrieve health status: %v", err)), nil
		}

		resultJSON, err := json.MarshalIndent(health, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize health status: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// HandleGetStatus handles the opentelemetry_get_status tool.
func HandleGetStatus(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		status, err := c.GetStatus(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to retrieve status: %v", err)), nil
		}

		resultJSON, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize status: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// HandleGetConfig handles the opentelemetry_get_config tool.
func HandleGetConfig(c *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		config, err := c.GetConfig(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to retrieve configuration: %v", err)), nil
		}

		resultJSON, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize configuration: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

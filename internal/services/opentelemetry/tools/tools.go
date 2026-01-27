// Package tools provides MCP tool definitions for OpenTelemetry operations.
package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
)

// GetMetricsTool returns a tool definition for retrieving metrics from OpenTelemetry Collector.
func GetMetricsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "opentelemetry_get_metrics",
		Description: "Retrieve metrics from OpenTelemetry Collector. Can filter by metric name and time range.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"metric_name": map[string]interface{}{
					"type":        "string",
					"description": "Optional: Filter by specific metric name",
				},
				"start_time": map[string]interface{}{
					"type":        "string",
					"description": "Optional: Start time in RFC3339 format (e.g., 2024-01-01T00:00:00Z)",
				},
				"end_time": map[string]interface{}{
					"type":        "string",
					"description": "Optional: End time in RFC3339 format (e.g., 2024-01-02T00:00:00Z)",
				},
			},
		},
	}
}

// QueryMetricsTool returns a tool definition for querying metrics with PromQL-style syntax.
func QueryMetricsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "opentelemetry_query_metrics",
		Description: "Execute a PromQL-style query against OpenTelemetry Collector metrics. Useful for aggregations, filtering, and complex metric analysis.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "PromQL-style query string (e.g., 'rate(http_requests_total[5m])')",
				},
				"time": map[string]interface{}{
					"type":        "string",
					"description": "Optional: Query timestamp in RFC3339 format. Defaults to current time.",
				},
			},
			Required: []string{"query"},
		},
	}
}

// GetTracesTool returns a tool definition for retrieving traces from OpenTelemetry Collector.
func GetTracesTool() mcp.Tool {
	return mcp.Tool{
		Name:        "opentelemetry_get_traces",
		Description: "Retrieve traces from OpenTelemetry Collector. Can filter by trace ID, service name, and time range.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"trace_id": map[string]interface{}{
					"type":        "string",
					"description": "Optional: Filter by specific trace ID",
				},
				"service": map[string]interface{}{
					"type":        "string",
					"description": "Optional: Filter by service name",
				},
				"start_time": map[string]interface{}{
					"type":        "string",
					"description": "Optional: Start time in RFC3339 format",
				},
				"end_time": map[string]interface{}{
					"type":        "string",
					"description": "Optional: End time in RFC3339 format",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Optional: Maximum number of traces to return",
				},
			},
		},
	}
}

// QueryTracesTool returns a tool definition for searching traces with custom criteria.
func QueryTracesTool() mcp.Tool {
	return mcp.Tool{
		Name:        "opentelemetry_query_traces",
		Description: "Search for traces matching custom criteria in OpenTelemetry Collector. Supports filtering by service, tags, and time range.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search query string (e.g., 'service=my-service AND error=true')",
				},
				"service": map[string]interface{}{
					"type":        "string",
					"description": "Optional: Filter by service name",
				},
				"start_time": map[string]interface{}{
					"type":        "string",
					"description": "Optional: Start time in RFC3339 format",
				},
				"end_time": map[string]interface{}{
					"type":        "string",
					"description": "Optional: End time in RFC3339 format",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Optional: Maximum number of traces to return",
				},
			},
			Required: []string{"query"},
		},
	}
}

// GetLogsTool returns a tool definition for retrieving logs from OpenTelemetry Collector.
func GetLogsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "opentelemetry_get_logs",
		Description: "Retrieve logs from OpenTelemetry Collector. Can filter by service name, log level, and time range.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"service": map[string]interface{}{
					"type":        "string",
					"description": "Optional: Filter by service name",
				},
				"level": map[string]interface{}{
					"type":        "string",
					"description": "Optional: Filter by log level (e.g., debug, info, warn, error)",
				},
				"start_time": map[string]interface{}{
					"type":        "string",
					"description": "Optional: Start time in RFC3339 format",
				},
				"end_time": map[string]interface{}{
					"type":        "string",
					"description": "Optional: End time in RFC3339 format",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Optional: Maximum number of log entries to return",
				},
			},
		},
	}
}

// QueryLogsTool returns a tool definition for searching logs with custom criteria.
func QueryLogsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "opentelemetry_query_logs",
		Description: "Search for logs matching custom criteria in OpenTelemetry Collector. Supports filtering by service, level, message content, and time range.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search query string (e.g., 'error' or 'service=my-service AND error=true')",
				},
				"service": map[string]interface{}{
					"type":        "string",
					"description": "Optional: Filter by service name",
				},
				"level": map[string]interface{}{
					"type":        "string",
					"description": "Optional: Filter by log level",
				},
				"start_time": map[string]interface{}{
					"type":        "string",
					"description": "Optional: Start time in RFC3339 format",
				},
				"end_time": map[string]interface{}{
					"type":        "string",
					"description": "Optional: End time in RFC3339 format",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Optional: Maximum number of log entries to return",
				},
			},
			Required: []string{"query"},
		},
	}
}

// GetHealthTool returns a tool definition for checking OpenTelemetry Collector health.
func GetHealthTool() mcp.Tool {
	return mcp.Tool{
		Name:        "opentelemetry_get_health",
		Description: "Check the health status of OpenTelemetry Collector. Returns overall health and component status.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// GetStatusTool returns a tool definition for retrieving OpenTelemetry Collector status.
func GetStatusTool() mcp.Tool {
	return mcp.Tool{
		Name:        "opentelemetry_get_status",
		Description: "Retrieve detailed status information about OpenTelemetry Collector, including components, pipelines, and configuration.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// GetConfigTool returns a tool definition for retrieving OpenTelemetry Collector configuration.
func GetConfigTool() mcp.Tool {
	return mcp.Tool{
		Name:        "opentelemetry_get_config",
		Description: "Retrieve the current configuration of OpenTelemetry Collector. Shows pipelines, receivers, processors, exporters, and extensions.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

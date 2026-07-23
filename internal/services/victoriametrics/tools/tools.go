// Package tools provides MCP tool definitions for VictoriaMetrics operations.
// It defines the structure and parameters for all VictoriaMetrics-related tools.
package tools

import "github.com/mark3labs/mcp-go/mcp"

// QueryTool returns the tool definition for VictoriaMetrics instant queries using MetricsQL.
func QueryTool() mcp.Tool {
	return mcp.Tool{
		Name:        "victoriametrics_query",
		Description: "Execute a VictoriaMetrics instant query using MetricsQL (backward-compatible with PromQL, with extensions like WITH templates, keep_last_value, drop_last_samples). Use this for the current value of a metric or expression; prefer `victoriametrics_query_range` only when you need a time series.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "MetricsQL or PromQL expression, for example `up`, `rate(http_requests_total[5m])`, or `sum by (pod) (container_memory_usage_bytes)`.",
				},
				"time": map[string]interface{}{
					"type":        "string",
					"description": "Optional RFC3339 timestamp. If omitted, the current server time is used.",
					"format":      "date-time",
				},
			},
			Required: []string{"query"},
		},
	}
}

// QueryRangeTool returns the tool definition for VictoriaMetrics range queries.
func QueryRangeTool() mcp.Tool {
	return mcp.Tool{
		Name:        "victoriametrics_query_range",
		Description: "Execute a VictoriaMetrics range query using MetricsQL and return a time series. Keep the time window narrow and the query specific to avoid very large responses.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "MetricsQL or PromQL expression to evaluate across a time range.",
				},
				"start": map[string]interface{}{
					"type":        "string",
					"description": "Required RFC3339 start timestamp.",
					"format":      "date-time",
				},
				"end": map[string]interface{}{
					"type":        "string",
					"description": "Required RFC3339 end timestamp.",
					"format":      "date-time",
				},
				"step": map[string]interface{}{
					"type":        "string",
					"description": "Query resolution step width in duration format (e.g., '15s', '1m', '1h').",
				},
			},
			Required: []string{"query", "start", "end", "step"},
		},
	}
}

// ExportTool returns the tool definition for VictoriaMetrics bulk export.
func ExportTool() mcp.Tool {
	return mcp.Tool{
		Name:        "victoriametrics_export",
		Description: "Export raw metrics data in bulk from VictoriaMetrics. Much more efficient than query_range for bulk data retrieval. Returns JSON lines format.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"match": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "string",
					},
					"description": "Time series selectors to match, e.g. ['{__name__=~\"http_requests_total\"}', '{job=\"api-server\"}'].",
				},
				"start": map[string]interface{}{
					"type":        "string",
					"description": "Required RFC3339 start timestamp for export window.",
					"format":      "date-time",
				},
				"end": map[string]interface{}{
					"type":        "string",
					"description": "Required RFC3339 end timestamp for export window.",
					"format":      "date-time",
				},
			},
			Required: []string{"match", "start", "end"},
		},
	}
}

// ImportTool returns the tool definition for VictoriaMetrics bulk import.
func ImportTool() mcp.Tool {
	return mcp.Tool{
		Name:        "victoriametrics_import",
		Description: "Import metrics data in bulk into VictoriaMetrics using the native import API. Accepts JSON lines format.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"data": map[string]interface{}{
					"type":        "string",
					"description": "JSON lines data to import. Each line should be a JSON object with metric, timestamps, and values.",
				},
			},
			Required: []string{"data"},
		},
	}
}

// SeriesCountTool returns the tool definition for VictoriaMetrics series count.
func SeriesCountTool() mcp.Tool {
	return mcp.Tool{
		Name:        "victoriametrics_series_count",
		Description: "Get the total number of time series currently stored in VictoriaMetrics. Much faster than counting via query.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// LabelNamesTool returns the tool definition for getting label names.
func LabelNamesTool() mcp.Tool {
	return mcp.Tool{
		Name:        "victoriametrics_get_label_names",
		Description: "Get all label names from VictoriaMetrics.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"start": map[string]interface{}{
					"type":        "string",
					"description": "Optional RFC3339 start timestamp for time range filter.",
					"format":      "date-time",
				},
				"end": map[string]interface{}{
					"type":        "string",
					"description": "Optional RFC3339 end timestamp for time range filter.",
					"format":      "date-time",
				},
			},
		},
	}
}

// LabelValuesTool returns the tool definition for getting label values.
func LabelValuesTool() mcp.Tool {
	return mcp.Tool{
		Name:        "victoriametrics_get_label_values",
		Description: "Get all values for a specific label name in VictoriaMetrics.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"label": map[string]interface{}{
					"type":        "string",
					"description": "The label name to query values for (e.g., 'job', 'instance').",
				},
				"start": map[string]interface{}{
					"type":        "string",
					"description": "Optional RFC3339 start timestamp.",
					"format":      "date-time",
				},
				"end": map[string]interface{}{
					"type":        "string",
					"description": "Optional RFC3339 end timestamp.",
					"format":      "date-time",
				},
			},
			Required: []string{"label"},
		},
	}
}

// TSDBStatsTool returns the tool definition for TSDB statistics.
func TSDBStatsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "victoriametrics_get_tsdb_stats",
		Description: "Get TSDB (time series database) storage statistics from VictoriaMetrics, including series count, storage usage, and index metrics.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// ActiveQueriesTool returns the tool definition for active queries.
func ActiveQueriesTool() mcp.Tool {
	return mcp.Tool{
		Name:        "victoriametrics_get_active_queries",
		Description: "Get a list of currently running/active queries in VictoriaMetrics. Useful for diagnosing slow queries.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// DeleteSeriesTool returns the tool definition for deleting series.
func DeleteSeriesTool() mcp.Tool {
	return mcp.Tool{
		Name:        "victoriametrics_delete_series",
		Description: "Delete time series matching the given label matchers from VictoriaMetrics. This is a destructive operation — use with caution.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"match": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "string",
					},
					"description": "Time series selectors to match for deletion, e.g. ['{__name__=~\"temp_metric.*\"}', '{job=\"test\"}'].",
				},
			},
			Required: []string{"match"},
		},
	}
}

// HealthCheckTool returns the tool definition for health check.
func HealthCheckTool() mcp.Tool {
	return mcp.Tool{
		Name:        "victoriametrics_health_check",
		Description: "Check VictoriaMetrics health status. Returns OK if the instance is healthy.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// GetMetadataTool returns the tool definition for metric metadata.
func GetMetadataTool() mcp.Tool {
	return mcp.Tool{
		Name:        "victoriametrics_get_metadata",
		Description: "Get metadata about metrics currently scraped by VictoriaMetrics, including type, help text, and unit information.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"metric": map[string]interface{}{
					"type":        "string",
					"description": "Optional metric name to filter metadata for.",
				},
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Optional maximum number of metadata entries to return.",
				},
			},
		},
	}
}

// ListToolsTool returns the tool definition for listing available VictoriaMetrics tools.
func ListToolsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "victoriametrics_list_tools",
		Description: "List all available VictoriaMetrics tools with their descriptions.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// AllTools returns all VictoriaMetrics tool definitions.
func AllTools() []mcp.Tool {
	return []mcp.Tool{
		QueryTool(),
		QueryRangeTool(),
		ExportTool(),
		ImportTool(),
		SeriesCountTool(),
		LabelNamesTool(),
		LabelValuesTool(),
		TSDBStatsTool(),
		ActiveQueriesTool(),
		DeleteSeriesTool(),
		HealthCheckTool(),
		GetMetadataTool(),
		ListToolsTool(),
	}
}

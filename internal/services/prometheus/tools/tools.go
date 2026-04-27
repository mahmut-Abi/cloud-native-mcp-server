// Package tools provides MCP tool definitions for Prometheus operations.
// It defines the structure and parameters for all Prometheus-related tools.
package tools

import "github.com/mark3labs/mcp-go/mcp"

// QueryTool returns the tool definition for Prometheus instant queries.
func QueryTool() mcp.Tool {
	return mcp.Tool{
		Name:        "prometheus_query",
		Description: "Execute a Prometheus instant query. Use this for the current value of a metric or expression; prefer `prometheus_query_range` only when you need a time series.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "PromQL expression, for example `up`, `rate(http_requests_total[5m])`, or `sum by (pod) (container_memory_usage_bytes)`.",
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

// QueryRangeTool returns the tool definition for Prometheus range queries.
func QueryRangeTool() mcp.Tool {
	return mcp.Tool{
		Name:        "prometheus_query_range",
		Description: "Execute a Prometheus range query and return a time series. Keep the time window narrow and the query specific to avoid very large responses.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "PromQL expression to evaluate across a time range.",
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
					"description": "Step duration such as `30s`, `1m`, or `5m`. Defaults to `15s`.",
					"default":     "15s",
				},
			},
			Required: []string{"query", "start", "end"},
		},
	}
}

// GetTargetsTool returns the tool definition for retrieving Prometheus targets.
func GetTargetsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "prometheus_get_targets",
		Description: "List Prometheus scrape targets. Use `state` to narrow results when you only need `active` or `dropped` targets.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"state": map[string]interface{}{
					"type":        "string",
					"description": "Optional state filter: `active`, `dropped`, or `any`.",
					"enum":        []string{"active", "dropped", "any"},
					"default":     "any",
				},
			},
		},
	}
}

// GetAlertsTool returns the tool definition for retrieving Prometheus alerts.
func GetAlertsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "prometheus_get_alerts",
		Description: "Retrieve the current active alerts from Prometheus",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// GetRulesTool returns the tool definition for retrieving Prometheus rules.
func GetRulesTool() mcp.Tool {
	return mcp.Tool{
		Name:        "prometheus_get_rules",
		Description: "List Prometheus recording and alerting rules. Use `type` to narrow results when you only want one rule class.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"type": map[string]interface{}{
					"type":        "string",
					"description": "Optional rule type filter: `alert` or `record`.",
					"enum":        []string{"alert", "record"},
				},
			},
		},
	}
}

// GetLabelNamesTool returns the tool definition for retrieving Prometheus label names.
func GetLabelNamesTool() mcp.Tool {
	return mcp.Tool{
		Name:        "prometheus_get_label_names",
		Description: "Retrieve all available label names from Prometheus",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"start": map[string]interface{}{
					"type":        "string",
					"description": "Optional start timestamp in RFC3339 format to limit label names to a time range",
					"format":      "date-time",
				},
				"end": map[string]interface{}{
					"type":        "string",
					"description": "Optional end timestamp in RFC3339 format to limit label names to a time range",
					"format":      "date-time",
				},
			},
		},
	}
}

// GetLabelValuesTool returns the tool definition for retrieving Prometheus label values.

// GetTargetsSummaryTool returns tool definition for getting Prometheus targets summary
func GetTargetsSummaryTool() mcp.Tool {
	return mcp.Tool{
		Name:        "prometheus_targets_summary",
		Description: "Get Prometheus targets summary (job, instance, health). 70-80% smaller than detailed version.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"state": map[string]interface{}{
					"type":        "string",
					"description": "Filter targets by state: active, dropped, any",
					"default":     "any",
				},
			},
		},
	}
}

// GetAlertsSummaryTool returns tool definition for getting Prometheus alerts summary
func GetAlertsSummaryTool() mcp.Tool {
	return mcp.Tool{
		Name:        "prometheus_alerts_summary",
		Description: "Get Prometheus alerts summary (alertname, state). 70-80% smaller than detailed.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// GetRulesSummaryTool returns tool definition for getting Prometheus rules summary
func GetRulesSummaryTool() mcp.Tool {
	return mcp.Tool{
		Name:        "prometheus_rules_summary",
		Description: "Get Prometheus rules summary (name, type). 70-80% smaller than detailed.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"type": map[string]interface{}{
					"type":        "string",
					"description": "Filter rules by type: alert, record",
				},
			},
		},
	}
}

func GetLabelValuesTool() mcp.Tool {
	return mcp.Tool{
		Name:        "prometheus_get_label_values",
		Description: "List all values for a specific Prometheus label. Use this after `prometheus_get_label_names` when you need valid values for a selector.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"label": map[string]interface{}{
					"type":        "string",
					"description": "Required label name such as `job`, `instance`, or `namespace`.",
				},
				"start": map[string]interface{}{
					"type":        "string",
					"description": "Optional start timestamp in RFC3339 format to limit label values to a time range",
					"format":      "date-time",
				},
				"end": map[string]interface{}{
					"type":        "string",
					"description": "Optional end timestamp in RFC3339 format to limit label values to a time range",
					"format":      "date-time",
				},
			},
			Required: []string{"label"},
		},
	}
}

// GetSeriesTool returns the tool definition for retrieving Prometheus series.
func GetSeriesTool() mcp.Tool {
	return mcp.Tool{
		Name:        "prometheus_get_series",
		Description: "List Prometheus series for one or more match selectors. This can be large, so prefer specific selectors such as `up{job=\"api\"}` instead of broad queries.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"match": map[string]interface{}{
					"type":        "array",
					"description": "One selector string or an array of selector strings. The server also accepts a JSON array string for compatibility.",
					"items": map[string]interface{}{
						"type": "string",
					},
					"minItems": 1,
				},
				"start": map[string]interface{}{
					"type":        "string",
					"description": "Optional start timestamp in RFC3339 format to limit series to a time range",
					"format":      "date-time",
				},
				"end": map[string]interface{}{
					"type":        "string",
					"description": "Optional end timestamp in RFC3339 format to limit series to a time range",
					"format":      "date-time",
				},
			},
			Required: []string{"match"},
		},
	}
}

// TestConnectionTool returns the tool definition for testing Prometheus connection.
func TestConnectionTool() mcp.Tool {
	return mcp.Tool{
		Name:        "prometheus_test_connection",
		Description: "Test the connection to the Prometheus server and verify API access",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// GetServerInfoTool returns the tool definition for retrieving Prometheus server info.
func GetServerInfoTool() mcp.Tool {
	return mcp.Tool{
		Name:        "prometheus_get_server_info",
		Description: "Retrieve Prometheus server information including build details and configuration",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// GetMetricsMetadataTool returns the tool definition for retrieving metrics metadata.
func GetMetricsMetadataTool() mcp.Tool {
	return mcp.Tool{
		Name:        "prometheus_get_metrics_metadata",
		Description: "Retrieve metadata for Prometheus metrics including type and help text",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"metric": map[string]interface{}{
					"type":        "string",
					"description": "Optional metric name to filter metadata. If not specified, returns metadata for all metrics",
				},
			},
		},
	}
}

// GetTargetMetadataTool returns the tool definition for retrieving target metadata.
func GetTargetMetadataTool() mcp.Tool {
	return mcp.Tool{
		Name:        "prometheus_get_target_metadata",
		Description: "Retrieve metadata about scrape targets for specific metrics",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"metric": map[string]interface{}{
					"type":        "string",
					"description": "Optional metric name to filter by. If not specified, returns metadata for all targets",
				},
			},
		},
	}
}

// GetTSDBStatsTool returns the tool definition for retrieving TSDB statistics.
func GetTSDBStatsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "prometheus_get_tsdb_stats",
		Description: "Retrieve Prometheus TSDB (time-series database) statistics including series count, chunk count, and storage info",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// GetTSDBStatusTool returns the tool definition for retrieving TSDB status.
func GetTSDBStatusTool() mcp.Tool {
	return mcp.Tool{
		Name:        "prometheus_get_tsdb_status",
		Description: "Retrieve Prometheus TSDB status information including head stats and block information",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// GetRuntimeInfoTool returns the tool definition for retrieving runtime information.
func GetRuntimeInfoTool() mcp.Tool {
	return mcp.Tool{
		Name:        "prometheus_get_runtime_info",
		Description: "Retrieve Prometheus runtime and build information including version, Go version, and build details",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

// CreateSnapshotTool returns the tool definition for creating TSDB snapshots.
func CreateSnapshotTool() mcp.Tool {
	return mcp.Tool{
		Name:        "prometheus_create_snapshot",
		Description: "Create a snapshot of all Prometheus data in TSDB. Returns the snapshot directory name and path",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"skipHead": map[string]interface{}{
					"type":        "boolean",
					"description": "Optional. If true, skip data in the WAL head and only snapshot data that has been compacted to disk",
				},
			},
		},
	}
}

// GetWALReplayStatusTool returns the tool definition for retrieving WAL replay status.
func GetWALReplayStatusTool() mcp.Tool {
	return mcp.Tool{
		Name:        "prometheus_get_wal_replay_status",
		Description: "Retrieve the status of WAL (Write-Ahead Log) replay on startup. Useful for debugging startup issues",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

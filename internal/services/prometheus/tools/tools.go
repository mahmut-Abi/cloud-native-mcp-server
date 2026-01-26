// Package tools provides MCP tool definitions for Prometheus operations.
// It defines the structure and parameters for all Prometheus-related tools.
package tools

import "github.com/mark3labs/mcp-go/mcp"

// QueryTool returns the tool definition for Prometheus instant queries.
func QueryTool() mcp.Tool {
	return mcp.Tool{
		Name:        "prometheus_query",
		Description: "Execute a Prometheus instant query to retrieve metric values at a specific point in time",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "PromQL query string to execute",
				},
				"time": map[string]interface{}{
					"type":        "string",
					"description": "Optional timestamp in RFC3339 format. If not specified, uses current time",
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
		Description: "⚠️ Returns extensive time series data which may be very large. Use smaller time ranges and specific queries",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "PromQL query string to execute",
				},
				"start": map[string]interface{}{
					"type":        "string",
					"description": "Start timestamp in RFC3339 format",
					"format":      "date-time",
				},
				"end": map[string]interface{}{
					"type":        "string",
					"description": "End timestamp in RFC3339 format",
					"format":      "date-time",
				},
				"step": map[string]interface{}{
					"type":        "string",
					"description": "Query resolution step width in duration format (e.g., '30s', '1m', '5m'). Defaults to '15s'",
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
		Description: "Retrieve the current state of Prometheus service discovery targets",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"state": map[string]interface{}{
					"type":        "string",
					"description": "Filter targets by state. Options: 'active', 'dropped', 'any'. Defaults to 'any'",
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
		Description: "Retrieve recording and alerting rules from Prometheus",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"type": map[string]interface{}{
					"type":        "string",
					"description": "Filter rules by type. Options: 'alert', 'record'. If not specified, returns all rules",
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
		Description: "Retrieve all available values for a specific label name from Prometheus",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"label": map[string]interface{}{
					"type":        "string",
					"description": "The label name to get values for",
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
		Description: "⚠️ May return large number of series. Consider specific label selectors to limit results",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"match": map[string]interface{}{
					"type":        "array",
					"description": "Label selector(s) to match series. Can be a single string or array of strings",
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

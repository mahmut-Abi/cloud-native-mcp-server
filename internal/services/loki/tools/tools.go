package tools

import "github.com/mark3labs/mcp-go/mcp"

// QueryLogsSummaryTool returns the tool definition for Loki log summaries.
func QueryLogsSummaryTool() mcp.Tool {
	return mcp.Tool{
		Name:        "loki_query_logs_summary",
		Description: "Recommended first step for LogQL exploration. Run a Loki range query and return a compact per-stream summary with line counts and a few sample lines.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Required LogQL query, usually a stream selector or pipeline such as `{namespace=\"prod\"} |= \"error\"`.",
				},
				"start": map[string]interface{}{
					"type":        "string",
					"description": "Optional RFC3339 start timestamp. Defaults to 1 hour ago.",
					"format":      "date-time",
				},
				"end": map[string]interface{}{
					"type":        "string",
					"description": "Optional RFC3339 end timestamp. Defaults to now.",
					"format":      "date-time",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Maximum log lines to inspect across returned streams. Defaults to 50, max 500.",
					"default":     50,
				},
				"direction": map[string]interface{}{
					"type":        "string",
					"description": "Result order: `backward` for newest-first or `forward` for oldest-first. Defaults to `backward`.",
					"enum":        []string{"backward", "forward"},
					"default":     "backward",
				},
			},
			Required: []string{"query"},
		},
	}
}

// QueryTool returns the tool definition for Loki instant queries.
func QueryTool() mcp.Tool {
	return mcp.Tool{
		Name:        "loki_query",
		Description: "Execute a Loki instant query. Use this for point-in-time LogQL evaluation or metric-style LogQL expressions.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Required LogQL expression.",
				},
				"time": map[string]interface{}{
					"type":        "string",
					"description": "Optional RFC3339 timestamp. If omitted, Loki uses the current time.",
					"format":      "date-time",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Maximum log lines to return for stream results. Defaults to 100, max 500.",
					"default":     100,
				},
				"direction": map[string]interface{}{
					"type":        "string",
					"description": "Result order for stream results.",
					"enum":        []string{"backward", "forward"},
					"default":     "backward",
				},
			},
			Required: []string{"query"},
		},
	}
}

// QueryRangeTool returns the tool definition for Loki range queries.
func QueryRangeTool() mcp.Tool {
	return mcp.Tool{
		Name:        "loki_query_range",
		Description: "Execute a Loki range query across a time window. Keep the selector specific and the time range narrow to avoid very large responses.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Required LogQL expression.",
				},
				"start": map[string]interface{}{
					"type":        "string",
					"description": "Optional RFC3339 start timestamp. Defaults to 1 hour ago.",
					"format":      "date-time",
				},
				"end": map[string]interface{}{
					"type":        "string",
					"description": "Optional RFC3339 end timestamp. Defaults to now.",
					"format":      "date-time",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Maximum log lines to return. Defaults to 100, max 500.",
					"default":     100,
				},
				"direction": map[string]interface{}{
					"type":        "string",
					"description": "Result order.",
					"enum":        []string{"backward", "forward"},
					"default":     "backward",
				},
				"step": map[string]interface{}{
					"type":        "string",
					"description": "Optional step duration such as `30s` or `1m`. Mainly useful for metric-style LogQL range queries.",
				},
			},
			Required: []string{"query"},
		},
	}
}

// GetLabelNamesTool returns the tool definition for Loki label names.
func GetLabelNamesTool() mcp.Tool {
	return mcp.Tool{
		Name:        "loki_get_label_names",
		Description: "List available Loki label names. Use this before building stream selectors when you do not yet know the indexed label keys.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Optional stream selector to narrow label discovery, for example `{namespace=\"prod\"}`.",
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
		},
	}
}

// GetLabelValuesTool returns the tool definition for Loki label values.
func GetLabelValuesTool() mcp.Tool {
	return mcp.Tool{
		Name:        "loki_get_label_values",
		Description: "List values for a specific Loki label. Use this after `loki_get_label_names` to discover valid selector values.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"label": map[string]interface{}{
					"type":        "string",
					"description": "Required Loki label name such as `namespace`, `pod`, or `container`.",
				},
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Optional stream selector to narrow the value set.",
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

// GetSeriesTool returns the tool definition for Loki series lookup.
func GetSeriesTool() mcp.Tool {
	return mcp.Tool{
		Name:        "loki_get_series",
		Description: "List indexed Loki series that match one or more stream selectors. Use this to confirm label combinations before running larger log queries.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"matchers": map[string]interface{}{
					"type":        "array",
					"description": "Required list of one or more stream selectors, for example [`{namespace=\"prod\"}`, `{app=\"api\"}`].",
					"items": map[string]interface{}{
						"type": "string",
					},
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
			Required: []string{"matchers"},
		},
	}
}

// TestConnectionTool returns the tool definition for checking Loki connectivity.
func TestConnectionTool() mcp.Tool {
	return mcp.Tool{
		Name:        "loki_test_connection",
		Description: "Check whether the configured Loki endpoint is reachable and authentication works.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
}

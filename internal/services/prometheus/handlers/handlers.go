// Package handlers provides HTTP handlers for Prometheus MCP operations.
// It implements request handling for Prometheus queries, metrics, targets, and alerts.
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/prometheus/client"
	optimize "github.com/mahmut-Abi/k8s-mcp-server/internal/util/performance"
)

// marshalIndentJSON performs indented JSON encoding using object pool
func marshalIndentJSON(data interface{}) ([]byte, error) {
	// First encode to compact format using object pool
	compactBytes, err := optimize.GlobalJSONPool.MarshalToBytes(data)
	if err != nil {
		return nil, err
	}

	// For scenarios requiring indented display, still use standard library but reduce allocations
	// This is a trade-off between performance and readability
	var result bytes.Buffer
	err = json.Indent(&result, compactBytes, "", "  ")
	return result.Bytes(), err
}

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
		query, ok := req.GetArguments()["query"].(string)
		if !ok || query == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Query parameter is required"),
				},
			}, nil
		}

		// Parse optional timestamp
		var timestamp *time.Time
		if ts, exists := req.GetArguments()["time"]; exists {
			if tsStr, ok := ts.(string); ok && tsStr != "" {
				if parsed, err := time.Parse(time.RFC3339, tsStr); err == nil {
					timestamp = &parsed
				}
			}
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
		query, ok := req.GetArguments()["query"].(string)
		if !ok || query == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Query parameter is required"),
				},
			}, nil
		}

		startStr, ok := req.GetArguments()["start"].(string)
		if !ok || startStr == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Start time parameter is required"),
				},
			}, nil
		}

		endStr, ok := req.GetArguments()["end"].(string)
		if !ok || endStr == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("End time parameter is required"),
				},
			}, nil
		}

		step, ok := req.GetArguments()["step"].(string)
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

// HandleGetTargets handles Prometheus targets retrieval requests.
func HandleGetTargets(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Prometheus get targets handler")

		// Extract optional state filter
		var state string
		if req.GetArguments() != nil {
			if s, exists := req.GetArguments()["state"]; exists {
				if stateStr, ok := s.(string); ok {
					state = stateStr
				}
			}
		}

		// Get targets
		targets, err := c.GetTargets(ctx, state)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get targets: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(targets)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format targets: %v", err)),
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

// HandleGetAlerts handles Prometheus alerts retrieval requests.
func HandleGetAlerts(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Prometheus get alerts handler")

		// Get alerts
		alerts, err := c.GetAlerts(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get alerts: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(alerts)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format alerts: %v", err)),
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

// HandleGetRules handles Prometheus rules retrieval requests.
func HandleGetRules(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Prometheus get rules handler")

		// Extract optional rule type filter
		var ruleType string
		if req.GetArguments() != nil {
			if t, exists := req.GetArguments()["type"]; exists {
				if typeStr, ok := t.(string); ok {
					ruleType = typeStr
				}
			}
		}

		// Get rules
		rules, err := c.GetRules(ctx, ruleType)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get rules: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(rules)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format rules: %v", err)),
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
		var matches []string
		if m, exists := req.GetArguments()["match"]; exists {
			switch v := m.(type) {
			case string:
				matches = []string{v}
			case []interface{}:
				for _, item := range v {
					if str, ok := item.(string); ok {
						matches = append(matches, str)
					}
				}
			}
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

// HandleTestConnection handles Prometheus connection test requests.
func HandleTestConnection(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Testing Prometheus connection")

		// Test connection
		err := c.TestConnection(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Connection test failed: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Connection test successful"),
			},
		}, nil
	}
}

// HandleGetServerInfo handles Prometheus server info requests.
func HandleGetServerInfo(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Prometheus server info handler")

		info, err := c.GetServerInfo(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get server info: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(info)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format server info: %v", err)),
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

// HandleGetMetricsMetadata handles metrics metadata requests.
func HandleGetMetricsMetadata(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Prometheus metrics metadata handler")

		metric := ""
		if m, exists := req.GetArguments()["metric"]; exists {
			if metricStr, ok := m.(string); ok {
				metric = metricStr
			}
		}

		metadata, err := c.GetMetricsMetadata(ctx, metric)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get metrics metadata: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(metadata)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format metrics metadata: %v", err)),
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

// HandleGetTargetMetadata handles target metadata requests.
func HandleGetTargetMetadata(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Prometheus target metadata handler")

		metric := ""
		if m, exists := req.GetArguments()["metric"]; exists {
			if metricStr, ok := m.(string); ok {
				metric = metricStr
			}
		}

		metadata, err := c.GetTargetMetadata(ctx, metric)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get target metadata: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(metadata)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format target metadata: %v", err)),
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

// ============ TSDB Handlers ============

// HandleGetTSDBStats handles TSDB stats requests.
func HandleGetTSDBStats(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Prometheus TSDB stats handler")

		stats, err := c.GetTSDBStats(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get TSDB stats: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(stats)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format TSDB stats: %v", err)),
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

// HandleGetTSDBStatus handles TSDB status requests.
func HandleGetTSDBStatus(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Prometheus TSDB status handler")

		status, err := c.GetTSDBStatus(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get TSDB status: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(status)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format TSDB status: %v", err)),
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

// HandleGetRuntimeInfo handles runtime info requests.
func HandleGetRuntimeInfo(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Prometheus runtime info handler")

		info, err := c.GetRuntimeInfo(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get runtime info: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(info)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format runtime info: %v", err)),
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

// HandleCreateSnapshot handles TSDB snapshot creation requests.
func HandleCreateSnapshot(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Prometheus TSDB snapshot handler")

		skipHead := false
		if skipHeadArg, exists := req.GetArguments()["skipHead"]; exists {
			if skipHeadBool, ok := skipHeadArg.(bool); ok {
				skipHead = skipHeadBool
			}
		}

		result, err := c.CreateSnapshot(ctx, skipHead)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to create snapshot: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(result)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format snapshot result: %v", err)),
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

// HandleGetWALReplayStatus handles WAL replay status requests.
func HandleGetWALReplayStatus(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Prometheus WAL replay status handler")

		status, err := c.GetWALReplayStatus(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get WAL replay status: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(status)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format WAL replay status: %v", err)),
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

// HandleGetTargetsSummary handles Prometheus targets summary requests (optimized version).
func HandleGetTargetsSummary(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		state := "any"
		if s, ok := req.GetArguments()["state"].(string); ok {
			state = s
		}

		logrus.WithField("state", state).Debug("Executing Prometheus targets summary handler")

		targets, err := c.GetTargets(ctx, state)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get targets: %v", err)),
				},
			}, nil
		}

		// Create summary (much smaller than full result)
		summary := make([]map[string]interface{}, 0, len(targets))
		for _, t := range targets {
			summary = append(summary, map[string]interface{}{
				"job":       t.Labels["job"],
				"instance":  t.Labels["instance"],
				"health":    t.Health,
				"scrapeUrl": t.ScrapeURL,
			})
		}

		result := map[string]interface{}{
			"count":   len(summary),
			"targets": summary,
		}

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

// HandleGetAlertsSummary handles Prometheus alerts summary requests (optimized version).
func HandleGetAlertsSummary(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Prometheus alerts summary handler")

		alerts, err := c.GetAlerts(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get alerts: %v", err)),
				},
			}, nil
		}

		// Create summary
		summary := make([]map[string]interface{}, 0, len(alerts))
		for _, a := range alerts {
			summary = append(summary, map[string]interface{}{
				"alertname": a.Labels["alertname"],
				"state":     a.State,
				"labels":    a.Labels,
			})
		}

		result := map[string]interface{}{
			"count":  len(summary),
			"alerts": summary,
		}

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

// HandleGetRulesSummary handles Prometheus rules summary requests (optimized version).
func HandleGetRulesSummary(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ruleType := ""
		if t, ok := req.GetArguments()["type"].(string); ok {
			ruleType = t
		}

		logrus.WithField("type", ruleType).Debug("Executing Prometheus rules summary handler")

		rules, err := c.GetRules(ctx, ruleType)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get rules: %v", err)),
				},
			}, nil
		}

		// Create summary
		summary := make([]map[string]interface{}, 0)
		for _, rg := range rules {
			for _, r := range rg.Rules {
				summary = append(summary, map[string]interface{}{
					"name":   r.Name,
					"type":   r.Type,
					"health": r.Health,
				})
			}
		}

		result := map[string]interface{}{
			"count": len(summary),
			"rules": summary,
		}

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

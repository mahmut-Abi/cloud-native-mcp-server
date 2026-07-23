// Package handlers provides MCP handler functions for VictoriaMetrics operations.
// It implements request handling for VictoriaMetrics queries, exports, imports,
// and instance management.
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	svccommon "github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/common"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/victoriametrics/client"
	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
)

// marshalIndentJSON performs indented JSON encoding using object pool.
func marshalIndentJSON(data interface{}) ([]byte, error) {
	compactBytes, err := optimize.GlobalJSONPool.MarshalToBytes(data)
	if err != nil {
		return nil, err
	}

	var result bytes.Buffer
	err = json.Indent(&result, compactBytes, "", "  ")
	return result.Bytes(), err
}

// HandleQuery handles VictoriaMetrics instant query requests using MetricsQL.
func HandleQuery() func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		logrus.Debug("Executing VictoriaMetrics query handler")

		args := req.Params.Arguments
		if args == nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("No arguments provided"),
				},
			}, nil
		}

		query, ok := svccommon.GetStringArg(req.GetArguments(), "query")
		if !ok || query == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Query parameter is required"),
				},
			}, nil
		}

		timestamp, err := svccommon.GetRFC3339TimeArg(req.GetArguments(), "time")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		result, err := c.Query(ctx, query, timestamp)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to execute query: %v", err)),
				},
			}, nil
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

// HandleQueryRange handles VictoriaMetrics range query requests.
func HandleQueryRange() func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		logrus.Debug("Executing VictoriaMetrics range query handler")

		query, ok := svccommon.GetStringArg(req.GetArguments(), "query")
		if !ok || query == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Query parameter is required"),
				},
			}, nil
		}

		start, ok := svccommon.GetStringArg(req.GetArguments(), "start")
		if !ok || start == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Start parameter is required"),
				},
			}, nil
		}

		end, ok := svccommon.GetStringArg(req.GetArguments(), "end")
		if !ok || end == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("End parameter is required"),
				},
			}, nil
		}

		step, ok := svccommon.GetStringArg(req.GetArguments(), "step")
		if !ok || step == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Step parameter is required"),
				},
			}, nil
		}

		result, err := c.QueryRange(ctx, query, start, end, step)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to execute range query: %v", err)),
				},
			}, nil
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

// HandleExport handles VictoriaMetrics bulk export requests.
func HandleExport() func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		logrus.Debug("Executing VictoriaMetrics export handler")

		match, ok, err := svccommon.GetStringSliceArg(req.GetArguments(), "match")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid match parameter: %v", err)), nil
		}
		if !ok || len(match) == 0 {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Match parameter is required (array of time series selectors)"),
				},
			}, nil
		}

		startStr, ok := svccommon.GetStringArg(req.GetArguments(), "start")
		if !ok {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Start parameter is required"),
				},
			}, nil
		}
		start, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid start timestamp: %v", err)), nil
		}

		endStr, ok := svccommon.GetStringArg(req.GetArguments(), "end")
		if !ok {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("End parameter is required"),
				},
			}, nil
		}
		end, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid end timestamp: %v", err)), nil
		}

		result, err := c.Export(ctx, match, start, end)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to export metrics: %v", err)),
				},
			}, nil
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

// HandleImport handles VictoriaMetrics bulk import requests.
func HandleImport() func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		logrus.Debug("Executing VictoriaMetrics import handler")

		dataStr, ok := svccommon.GetStringArg(req.GetArguments(), "data")
		if !ok || dataStr == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Data parameter is required"),
				},
			}, nil
		}

		// Parse JSON lines
		var importData []client.ImportRequest
		if err := json.Unmarshal([]byte(dataStr), &importData); err != nil {
			// Try as single object
			var single client.ImportRequest
			if err2 := json.Unmarshal([]byte(dataStr), &single); err2 != nil {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{
						mcp.NewTextContent(fmt.Sprintf("Failed to parse import data: %v", err)),
					},
				}, nil
			}
			importData = []client.ImportRequest{single}
		}

		if err := c.Import(ctx, importData); err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to import metrics: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully imported %d metrics", len(importData))),
			},
		}, nil
	}
}

// HandleSeriesCount handles VictoriaMetrics series count requests.
func HandleSeriesCount() func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		logrus.Debug("Executing VictoriaMetrics series count handler")

		count, err := c.GetSeriesCount(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get series count: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Total series count: %d", count)),
			},
		}, nil
	}
}

// HandleLabelNames handles VictoriaMetrics label name queries.
func HandleLabelNames() func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		logrus.Debug("Executing VictoriaMetrics label names handler")

		var startPtr, endPtr *time.Time
		if startStr, ok := svccommon.GetStringArg(req.GetArguments(), "start"); ok && startStr != "" {
			t, err := time.Parse(time.RFC3339, startStr)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Invalid start time: %v", err)), nil
			}
			startPtr = &t
		}
		if endStr, ok := svccommon.GetStringArg(req.GetArguments(), "end"); ok && endStr != "" {
			t, err := time.Parse(time.RFC3339, endStr)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Invalid end time: %v", err)), nil
			}
			endPtr = &t
		}

		labels, err := c.GetLabelNames(ctx, startPtr, endPtr)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get label names: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(labels)
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

// HandleLabelValues handles VictoriaMetrics label value queries.
func HandleLabelValues() func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		logrus.Debug("Executing VictoriaMetrics label values handler")

		label, ok := svccommon.GetStringArg(req.GetArguments(), "label")
		if !ok || label == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Label parameter is required"),
				},
			}, nil
		}

		var startPtr, endPtr *time.Time
		if startStr, ok := svccommon.GetStringArg(req.GetArguments(), "start"); ok && startStr != "" {
			t, err := time.Parse(time.RFC3339, startStr)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Invalid start time: %v", err)), nil
			}
			startPtr = &t
		}
		if endStr, ok := svccommon.GetStringArg(req.GetArguments(), "end"); ok && endStr != "" {
			t, err := time.Parse(time.RFC3339, endStr)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Invalid end time: %v", err)), nil
			}
			endPtr = &t
		}

		values, err := c.GetLabelValues(ctx, label, startPtr, endPtr)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get label values: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(values)
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

// HandleTSDBStats handles VictoriaMetrics TSDB statistics requests.
func HandleTSDBStats() func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		logrus.Debug("Executing VictoriaMetrics TSDB stats handler")

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

// HandleActiveQueries handles VictoriaMetrics active queries requests.
func HandleActiveQueries() func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		logrus.Debug("Executing VictoriaMetrics active queries handler")

		queries, err := c.GetActiveQueries(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get active queries: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(queries)
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

// HandleDeleteSeries handles VictoriaMetrics series deletion requests.
func HandleDeleteSeries() func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		logrus.Debug("Executing VictoriaMetrics delete series handler")

		match, ok, err := svccommon.GetStringSliceArg(req.GetArguments(), "match")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid match parameter: %v", err)), nil
		}
		if !ok || len(match) == 0 {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Match parameter is required (array of time series selectors)"),
				},
			}, nil
		}

		deleted, err := c.DeleteSeries(ctx, match)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to delete series: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully deleted %d time series", deleted)),
			},
		}, nil
	}
}

// HandleHealthCheck handles VictoriaMetrics health check requests.
func HandleHealthCheck() func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		logrus.Debug("Checking VictoriaMetrics health")

		health, err := c.HealthCheck(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Health check failed: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("VictoriaMetrics health: %s", health.Status)),
			},
		}, nil
	}
}

// HandleGetMetadata handles VictoriaMetrics metric metadata requests.
func HandleGetMetadata() func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		logrus.Debug("Executing VictoriaMetrics metadata handler")

		metric, _ := svccommon.GetStringArg(req.GetArguments(), "metric")
		limit := svccommon.GetIntArg(req.GetArguments(), 0, "limit")

		metadata, err := c.GetMetadata(ctx, metric, limit)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get metadata: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(metadata)
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

// HandleListTools handles listing available VictoriaMetrics tools.
func HandleListTools() func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Listing VictoriaMetrics tools")

		toolList := []map[string]string{
			{"name": "victoriametrics_query", "description": "Execute a VictoriaMetrics instant query using MetricsQL (backward-compatible with PromQL)."},
			{"name": "victoriametrics_query_range", "description": "Execute a VictoriaMetrics range query using MetricsQL over a time range."},
			{"name": "victoriametrics_export", "description": "Export raw metrics data in bulk (more efficient than query_range for bulk data)."},
			{"name": "victoriametrics_import", "description": "Import metrics data in bulk into VictoriaMetrics."},
			{"name": "victoriametrics_series_count", "description": "Get the total number of time series currently stored."},
			{"name": "victoriametrics_get_label_names", "description": "Get all label names from VictoriaMetrics."},
			{"name": "victoriametrics_get_label_values", "description": "Get all values for a specific label name."},
			{"name": "victoriametrics_get_tsdb_stats", "description": "Get TSDB storage statistics including series count and disk usage."},
			{"name": "victoriametrics_get_active_queries", "description": "Get currently running/active queries. Useful for diagnosing slow queries."},
			{"name": "victoriametrics_delete_series", "description": "Delete time series matching given label matchers. Destructive operation - use with caution."},
			{"name": "victoriametrics_health_check", "description": "Check VictoriaMetrics health status."},
			{"name": "victoriametrics_get_metadata", "description": "Get metadata about metrics including type, help text, and unit."},
			{"name": "victoriametrics_list_tools", "description": "List all available VictoriaMetrics tools."},
		}

		resultJSON, err := marshalIndentJSON(toolList)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format tool list: %v", err)),
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

package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/common"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/loki/client"
	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

const (
	defaultQueryLimit       = 100
	defaultSummaryLimit     = 50
	maxQueryLimit           = 500
	defaultQueryRangeWindow = time.Hour
)

// ServiceInterface defines the interface for the Loki service.
type ServiceInterface interface {
	GetClient() *client.Client
}

// QueryHandler handles loki_query.
func QueryHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getRequestArguments(request)

		query, err := common.RequireStringArg(args, "query", "expr", "expression", "logql")
		if err != nil {
			return nil, err
		}

		queryTime, err := common.GetRFC3339TimeArg(args, "time", "query_time", "queryTime")
		if err != nil {
			return nil, err
		}

		lokiClient, err := getLokiClient(service)
		if err != nil {
			return nil, err
		}

		result, err := lokiClient.Query(
			ctx,
			query,
			queryTime,
			boundedLimit(common.GetIntArg(args, defaultQueryLimit, "limit", "max_lines", "maxLines"), defaultQueryLimit),
			normalizeDirection(optionalString(args, "direction")),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to execute Loki instant query: %w", err)
		}

		return marshalResult(result)
	}
}

// QueryRangeHandler handles loki_query_range.
func QueryRangeHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getRequestArguments(request)

		query, err := common.RequireStringArg(args, "query", "expr", "expression", "logql")
		if err != nil {
			return nil, err
		}

		start, end, err := getQueryWindow(args)
		if err != nil {
			return nil, err
		}

		lokiClient, err := getLokiClient(service)
		if err != nil {
			return nil, err
		}

		result, err := lokiClient.QueryRange(
			ctx,
			query,
			start,
			end,
			boundedLimit(common.GetIntArg(args, defaultQueryLimit, "limit", "max_lines", "maxLines"), defaultQueryLimit),
			normalizeDirection(optionalString(args, "direction")),
			optionalString(args, "step"),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to execute Loki range query: %w", err)
		}

		return marshalResult(result)
	}
}

// QueryLogsSummaryHandler handles loki_query_logs_summary.
func QueryLogsSummaryHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getRequestArguments(request)

		query, err := common.RequireStringArg(args, "query", "expr", "expression", "logql")
		if err != nil {
			return nil, err
		}

		start, end, err := getQueryWindow(args)
		if err != nil {
			return nil, err
		}

		lokiClient, err := getLokiClient(service)
		if err != nil {
			return nil, err
		}

		result, err := lokiClient.QueryRange(
			ctx,
			query,
			start,
			end,
			boundedLimit(common.GetIntArg(args, defaultSummaryLimit, "limit", "max_lines", "maxLines"), defaultSummaryLimit),
			normalizeDirection(optionalString(args, "direction")),
			"",
		)
		if err != nil {
			return nil, fmt.Errorf("failed to execute Loki range query: %w", err)
		}

		summary, err := summarizeQueryResult(query, result)
		if err != nil {
			return nil, err
		}
		return marshalResult(summary)
	}
}

// GetLabelNamesHandler handles loki_get_label_names.
func GetLabelNamesHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getRequestArguments(request)

		start, err := common.GetRFC3339TimeArg(args, "start", "start_time", "startTime")
		if err != nil {
			return nil, err
		}
		end, err := common.GetRFC3339TimeArg(args, "end", "end_time", "endTime")
		if err != nil {
			return nil, err
		}

		lokiClient, err := getLokiClient(service)
		if err != nil {
			return nil, err
		}

		labels, err := lokiClient.GetLabelNames(ctx, optionalString(args, "query", "selector"), start, end)
		if err != nil {
			return nil, fmt.Errorf("failed to get Loki label names: %w", err)
		}

		return marshalResult(map[string]interface{}{
			"count":  len(labels),
			"labels": labels,
		})
	}
}

// GetLabelValuesHandler handles loki_get_label_values.
func GetLabelValuesHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getRequestArguments(request)

		label, err := common.RequireStringArg(args, "label", "label_name", "labelName", "name")
		if err != nil {
			return nil, err
		}
		start, err := common.GetRFC3339TimeArg(args, "start", "start_time", "startTime")
		if err != nil {
			return nil, err
		}
		end, err := common.GetRFC3339TimeArg(args, "end", "end_time", "endTime")
		if err != nil {
			return nil, err
		}

		lokiClient, err := getLokiClient(service)
		if err != nil {
			return nil, err
		}

		values, err := lokiClient.GetLabelValues(ctx, label, optionalString(args, "query", "selector"), start, end)
		if err != nil {
			return nil, fmt.Errorf("failed to get Loki label values: %w", err)
		}

		return marshalResult(map[string]interface{}{
			"label":  label,
			"count":  len(values),
			"values": values,
		})
	}
}

// GetSeriesHandler handles loki_get_series.
func GetSeriesHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := getRequestArguments(request)

		matchers, ok, err := common.GetStringSliceArg(args, "matchers", "match", "matches", "selectors", "selector")
		if err != nil {
			return nil, err
		}
		if !ok || len(matchers) == 0 {
			return nil, fmt.Errorf("missing required parameter: matchers")
		}

		start, err := common.GetRFC3339TimeArg(args, "start", "start_time", "startTime")
		if err != nil {
			return nil, err
		}
		end, err := common.GetRFC3339TimeArg(args, "end", "end_time", "endTime")
		if err != nil {
			return nil, err
		}

		lokiClient, err := getLokiClient(service)
		if err != nil {
			return nil, err
		}

		series, err := lokiClient.GetSeries(ctx, matchers, start, end)
		if err != nil {
			return nil, fmt.Errorf("failed to get Loki series: %w", err)
		}

		return marshalResult(map[string]interface{}{
			"count":   len(series),
			"series":  series,
			"filters": matchers,
		})
	}
}

// TestConnectionHandler handles loki_test_connection.
func TestConnectionHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		lokiClient, err := getLokiClient(service)
		if err != nil {
			return nil, err
		}

		if err := lokiClient.TestConnection(ctx); err != nil {
			return nil, fmt.Errorf("failed to connect to Loki: %w", err)
		}

		return marshalResult(map[string]interface{}{
			"status":  "ok",
			"message": "Loki connection successful",
		})
	}
}

func getLokiClient(service ServiceInterface) (*client.Client, error) {
	lokiClient := service.GetClient()
	if lokiClient == nil {
		return nil, fmt.Errorf("loki client is not initialized")
	}
	return lokiClient, nil
}

func getRequestArguments(request mcp.CallToolRequest) map[string]interface{} {
	args := request.GetArguments()
	if args == nil {
		return nil
	}

	nested, ok := args["params"].(map[string]interface{})
	if !ok || len(nested) == 0 {
		return args
	}

	merged := make(map[string]interface{}, len(args)+len(nested))
	for key, value := range nested {
		merged[key] = value
	}
	for key, value := range args {
		if key == "params" {
			continue
		}
		merged[key] = value
	}
	return merged
}

func optionalString(args map[string]interface{}, keys ...string) string {
	value, _ := common.GetStringArg(args, keys...)
	return value
}

func getQueryWindow(args map[string]interface{}) (time.Time, time.Time, error) {
	end, err := common.GetRFC3339TimeArg(args, "end", "end_time", "endTime")
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	start, err := common.GetRFC3339TimeArg(args, "start", "start_time", "startTime")
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endValue := time.Now().UTC()
	if end != nil {
		endValue = end.UTC()
	}

	startValue := endValue.Add(-defaultQueryRangeWindow)
	if start != nil {
		startValue = start.UTC()
	}
	if !startValue.Before(endValue) {
		return time.Time{}, time.Time{}, fmt.Errorf("start must be before end")
	}

	return startValue, endValue, nil
}

func boundedLimit(limit int, defaultValue int) int {
	if limit <= 0 {
		return defaultValue
	}
	if limit > maxQueryLimit {
		return maxQueryLimit
	}
	return limit
}

func normalizeDirection(direction string) string {
	switch strings.ToLower(strings.TrimSpace(direction)) {
	case "forward":
		return "forward"
	case "backward", "":
		return "backward"
	default:
		return "backward"
	}
}

func summarizeQueryResult(query string, result map[string]interface{}) (map[string]interface{}, error) {
	data, _ := result["data"].(map[string]interface{})
	resultType, _ := data["resultType"].(string)
	items, _ := data["result"].([]interface{})

	streams := make([]map[string]interface{}, 0, len(items))
	totalLines := 0
	for _, item := range items {
		streamItem, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		labels := toStringMap(streamItem["stream"])
		summary := map[string]interface{}{
			"labels": labels,
		}

		var rawLines []interface{}
		if values, ok := streamItem["values"].([]interface{}); ok {
			rawLines = values
		} else if value, ok := streamItem["value"].([]interface{}); ok {
			rawLines = []interface{}{value}
		}

		lineCount := len(rawLines)
		totalLines += lineCount
		summary["line_count"] = lineCount

		sampleLines := make([]string, 0, 3)
		for index, line := range rawLines {
			pair, ok := line.([]interface{})
			if !ok || len(pair) < 2 {
				continue
			}

			timestamp := strings.TrimSpace(fmt.Sprintf("%v", pair[0]))
			message := strings.TrimSpace(fmt.Sprintf("%v", pair[1]))
			if index == 0 {
				summary["first_timestamp"] = formatLokiTimestamp(timestamp)
			}
			summary["last_timestamp"] = formatLokiTimestamp(timestamp)
			if len(sampleLines) < 3 && message != "" {
				sampleLines = append(sampleLines, message)
			}
		}
		summary["sample_lines"] = sampleLines
		streams = append(streams, summary)
	}

	return map[string]interface{}{
		"query":        query,
		"result_type":  resultType,
		"stream_count": len(streams),
		"total_lines":  totalLines,
		"streams":      streams,
	}, nil
}

func toStringMap(raw interface{}) map[string]string {
	typed, ok := raw.(map[string]interface{})
	if !ok {
		return map[string]string{}
	}

	result := make(map[string]string, len(typed))
	for key, value := range typed {
		result[key] = strings.TrimSpace(fmt.Sprintf("%v", value))
	}
	return result
}

func formatLokiTimestamp(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return raw
	}
	return time.Unix(0, value).UTC().Format(time.RFC3339Nano)
}

func marshalResult(result interface{}) (*mcp.CallToolResult, error) {
	payload, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize response: %w", err)
	}
	return mcp.NewToolResultText(string(payload)), nil
}

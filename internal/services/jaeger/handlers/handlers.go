package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/jaeger/client"
	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

// ServiceInterface defines the interface for Jaeger service.
type ServiceInterface interface {
	GetClient() *client.Client
}

// GetTracesHandler handles the jaeger_get_traces tool.
func GetTracesHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		serviceName := getStringArg(args, "service")
		operation := getStringArg(args, "operation")
		startTime := getStringArg(args, "start_time")
		endTime := getStringArg(args, "end_time")
		limit := getBoundedIntArg(args, "limit", 20, 100)
		minDuration := getStringArg(args, "min_duration")
		maxDuration := getStringArg(args, "max_duration")

		if startTime == "" {
			startTime = time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
		}
		if endTime == "" {
			endTime = time.Now().Format(time.RFC3339)
		}

		jaegerClient, err := getJaegerClient(service)
		if err != nil {
			return nil, err
		}

		params := client.TraceQueryParameters{
			Service:     serviceName,
			Operation:   operation,
			StartTime:   startTime,
			EndTime:     endTime,
			Limit:       limit,
			MinDuration: minDuration,
			MaxDuration: maxDuration,
		}

		traces, err := jaegerClient.SearchTraces(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to search traces: %w", err)
		}

		return marshalResult(map[string]interface{}{
			"count":  len(traces),
			"traces": traces,
		})
	}
}

// GetTraceHandler handles the jaeger_get_trace tool.
func GetTraceHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		traceID, err := getRequiredStringArg(args, "trace_id")
		if err != nil {
			return nil, err
		}

		jaegerClient, err := getJaegerClient(service)
		if err != nil {
			return nil, err
		}

		trace, err := jaegerClient.GetTrace(ctx, traceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get trace: %w", err)
		}

		return marshalResult(map[string]interface{}{
			"trace": trace,
		})
	}
}

// GetServicesHandler handles the jaeger_get_services tool.
func GetServicesHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		jaegerClient, err := getJaegerClient(service)
		if err != nil {
			return nil, err
		}

		services, err := jaegerClient.GetServices(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get services: %w", err)
		}

		return marshalResult(map[string]interface{}{
			"count":    len(services),
			"services": services,
		})
	}
}

// GetServiceOperationsHandler handles the jaeger_get_service_ops tool.
func GetServiceOperationsHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		serviceName, err := getRequiredStringArg(args, "service")
		if err != nil {
			return nil, err
		}

		jaegerClient, err := getJaegerClient(service)
		if err != nil {
			return nil, err
		}

		operations, err := jaegerClient.GetOperations(ctx, serviceName)
		if err != nil {
			return nil, fmt.Errorf("failed to get service operations: %w", err)
		}

		return marshalResult(map[string]interface{}{
			"service":    serviceName,
			"count":      len(operations),
			"operations": operations,
		})
	}
}

// SearchTracesHandler handles the jaeger_search_traces tool.
func SearchTracesHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		serviceName := getStringArg(args, "service")
		operation := getStringArg(args, "operation")
		startTime := getStringArg(args, "start_time")
		endTime := getStringArg(args, "end_time")
		limit := getBoundedIntArg(args, "limit", 20, 100)
		minDuration := getStringArg(args, "min_duration")
		maxDuration := getStringArg(args, "max_duration")
		tags := getStringMapArg(args, "tags")

		if startTime == "" {
			startTime = time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
		}
		if endTime == "" {
			endTime = time.Now().Format(time.RFC3339)
		}

		jaegerClient, err := getJaegerClient(service)
		if err != nil {
			return nil, err
		}

		params := client.TraceQueryParameters{
			Service:     serviceName,
			Operation:   operation,
			Tags:        tags,
			StartTime:   startTime,
			EndTime:     endTime,
			Limit:       limit,
			MinDuration: minDuration,
			MaxDuration: maxDuration,
		}

		traces, err := jaegerClient.SearchTraces(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to search traces: %w", err)
		}

		return marshalResult(map[string]interface{}{
			"count":  len(traces),
			"traces": traces,
		})
	}
}

// GetDependenciesHandler handles the jaeger_get_dependencies tool.
func GetDependenciesHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		startTime := getStringArg(args, "start_time")
		endTime := getStringArg(args, "end_time")

		if startTime == "" {
			startTime = time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
		}
		if endTime == "" {
			endTime = time.Now().Format(time.RFC3339)
		}

		jaegerClient, err := getJaegerClient(service)
		if err != nil {
			return nil, err
		}

		dependencies, err := jaegerClient.GetDependencies(ctx, startTime, endTime)
		if err != nil {
			return nil, fmt.Errorf("failed to get dependencies: %w", err)
		}

		return marshalResult(map[string]interface{}{
			"count":        len(dependencies),
			"dependencies": dependencies,
		})
	}
}

// GetTracesSummaryHandler handles the jaeger_get_traces_summary tool.
func GetTracesSummaryHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		serviceName := getStringArg(args, "service")
		operation := getStringArg(args, "operation")
		startTime := getStringArg(args, "start_time")
		endTime := getStringArg(args, "end_time")
		limit := getBoundedIntArg(args, "limit", 20, 100)
		minDuration := getStringArg(args, "min_duration")

		if startTime == "" {
			startTime = time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
		}
		if endTime == "" {
			endTime = time.Now().Format(time.RFC3339)
		}

		jaegerClient, err := getJaegerClient(service)
		if err != nil {
			return nil, err
		}

		params := client.TraceQueryParameters{
			Service:     serviceName,
			Operation:   operation,
			StartTime:   startTime,
			EndTime:     endTime,
			Limit:       limit,
			MinDuration: minDuration,
		}

		traces, err := jaegerClient.SearchTraces(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to search traces: %w", err)
		}

		summaries := make([]map[string]interface{}, 0, len(traces))
		for _, trace := range traces {
			summary := map[string]interface{}{
				"trace_id":    trace.TraceID,
				"span_count":  len(trace.Spans),
				"duration_us": 0,
			}

			if len(trace.Spans) > 0 {
				summary["operation"] = trace.Spans[0].OperationName
				summary["duration_us"] = trace.Spans[0].Duration
			}
			if len(trace.Processes) > 0 {
				summary["service"] = trace.Processes[0].ServiceName
			}

			summaries = append(summaries, summary)
		}

		return marshalResult(map[string]interface{}{
			"count":  len(summaries),
			"traces": summaries,
		})
	}
}

// GetServicesSummaryHandler handles the jaeger_get_services_summary tool.
func GetServicesSummaryHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		jaegerClient, err := getJaegerClient(service)
		if err != nil {
			return nil, err
		}

		services, err := jaegerClient.GetServices(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get services: %w", err)
		}

		summaries := make([]map[string]interface{}, 0, len(services))
		for _, svc := range services {
			summaries = append(summaries, map[string]interface{}{
				"name":            svc.Name,
				"operation_count": len(svc.Operations),
			})
		}

		return marshalResult(map[string]interface{}{
			"count":    len(summaries),
			"services": summaries,
		})
	}
}

func getJaegerClient(service ServiceInterface) (*client.Client, error) {
	jaegerClient := service.GetClient()
	if jaegerClient == nil {
		return nil, fmt.Errorf("jaeger client is not initialized")
	}
	return jaegerClient, nil
}

func getRequiredStringArg(args map[string]interface{}, key string) (string, error) {
	value := getStringArg(args, key)
	if value == "" {
		return "", fmt.Errorf("missing required parameter: %s", key)
	}
	return value, nil
}

func getStringArg(args map[string]interface{}, key string) string {
	value, _ := args[key].(string)
	return strings.TrimSpace(value)
}

func getBoundedIntArg(args map[string]interface{}, key string, def, max int) int {
	value := def
	if raw, ok := args[key]; ok {
		switch typed := raw.(type) {
		case float64:
			value = int(typed)
		case float32:
			value = int(typed)
		case int:
			value = typed
		case int64:
			value = int(typed)
		case string:
			if parsed, err := strconv.Atoi(strings.TrimSpace(typed)); err == nil {
				value = parsed
			}
		}
	}

	if value <= 0 {
		value = def
	}
	if value > max {
		value = max
	}
	return value
}

func getStringMapArg(args map[string]interface{}, key string) map[string]string {
	typedMap, ok := args[key].(map[string]interface{})
	if !ok {
		return map[string]string{}
	}

	result := make(map[string]string, len(typedMap))
	for k, v := range typedMap {
		value := strings.TrimSpace(fmt.Sprintf("%v", v))
		if value != "" {
			result[k] = value
		}
	}
	return result
}

func marshalResult(result map[string]interface{}) (*mcp.CallToolResult, error) {
	jsonResponse, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize response: %w", err)
	}
	return mcp.NewToolResultText(string(jsonResponse)), nil
}

package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/jaeger/client"
	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

// ServiceInterface defines the interface for Jaeger service
type ServiceInterface interface {
	GetClient() *client.Client
}

// GetTracesHandler handles the jaeger_get_traces tool.
func GetTracesHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		// Parse parameters
		serviceName := args["service"].(string)
		operation := args["operation"].(string)
		startTime := args["start_time"].(string)
		endTime := args["end_time"].(string)
		limit := int(args["limit"].(float64))
		minDuration := args["min_duration"].(string)
		maxDuration := args["max_duration"].(string)

		// Set defaults
		if limit == 0 {
			limit = 20
		}
		if limit > 100 {
			limit = 100
		}
		if startTime == "" {
			startTime = time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
		}
		if endTime == "" {
			endTime = time.Now().Format(time.RFC3339)
		}

		// Build query parameters
		params := client.TraceQueryParameters{
			Service:     serviceName,
			Operation:   operation,
			StartTime:   startTime,
			EndTime:     endTime,
			Limit:       limit,
			MinDuration: minDuration,
			MaxDuration: maxDuration,
		}

		// Search traces
		traces, err := service.GetClient().SearchTraces(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to search traces: %w", err)
		}

		// Serialize response
		result := map[string]interface{}{
			"count":  len(traces),
			"traces": traces,
		}

		jsonResponse, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize response: %w", err)
		}

		return mcp.NewToolResultText(string(jsonResponse)), nil
	}
}

// GetTraceHandler handles the jaeger_get_trace tool.
func GetTraceHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		traceID := args["trace_id"].(string)

		// Get trace
		trace, err := service.GetClient().GetTrace(ctx, traceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get trace: %w", err)
		}

		// Serialize response
		result := map[string]interface{}{
			"trace": trace,
		}

		jsonResponse, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize response: %w", err)
		}

		return mcp.NewToolResultText(string(jsonResponse)), nil
	}
}

// GetServicesHandler handles the jaeger_get_services tool.
func GetServicesHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get services
		services, err := service.GetClient().GetServices(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get services: %w", err)
		}

		// Serialize response
		result := map[string]interface{}{
			"count":    len(services),
			"services": services,
		}

		jsonResponse, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize response: %w", err)
		}

		return mcp.NewToolResultText(string(jsonResponse)), nil
	}
}

// GetServiceOperationsHandler handles the jaeger_get_service_ops tool.
func GetServiceOperationsHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		serviceName := args["service"].(string)

		// Get operations
		operations, err := service.GetClient().GetOperations(ctx, serviceName)
		if err != nil {
			return nil, fmt.Errorf("failed to get service operations: %w", err)
		}

		// Serialize response
		result := map[string]interface{}{
			"service":    serviceName,
			"count":      len(operations),
			"operations": operations,
		}

		jsonResponse, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize response: %w", err)
		}

		return mcp.NewToolResultText(string(jsonResponse)), nil
	}
}

// SearchTracesHandler handles the jaeger_search_traces tool.
func SearchTracesHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		// Parse parameters
		serviceName := args["service"].(string)
		operation := args["operation"].(string)
		startTime := args["start_time"].(string)
		endTime := args["end_time"].(string)
		limit := int(args["limit"].(float64))
		minDuration := args["min_duration"].(string)
		maxDuration := args["max_duration"].(string)

		// Parse tags
		tags := make(map[string]string)
		if tagsArg, ok := args["tags"].(map[string]interface{}); ok {
			for key, value := range tagsArg {
				tags[key] = fmt.Sprintf("%v", value)
			}
		}

		// Set defaults
		if limit == 0 {
			limit = 20
		}
		if limit > 100 {
			limit = 100
		}
		if startTime == "" {
			startTime = time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
		}
		if endTime == "" {
			endTime = time.Now().Format(time.RFC3339)
		}

		// Build query parameters
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

		// Search traces
		traces, err := service.GetClient().SearchTraces(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to search traces: %w", err)
		}

		// Serialize response
		result := map[string]interface{}{
			"count":  len(traces),
			"traces": traces,
		}

		jsonResponse, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize response: %w", err)
		}

		return mcp.NewToolResultText(string(jsonResponse)), nil
	}
}

// GetDependenciesHandler handles the jaeger_get_dependencies tool.
func GetDependenciesHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		startTime := args["start_time"].(string)
		endTime := args["end_time"].(string)

		// Set defaults
		if startTime == "" {
			startTime = time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
		}
		if endTime == "" {
			endTime = time.Now().Format(time.RFC3339)
		}

		// Get dependencies
		dependencies, err := service.GetClient().GetDependencies(ctx, startTime, endTime)
		if err != nil {
			return nil, fmt.Errorf("failed to get dependencies: %w", err)
		}

		// Serialize response
		result := map[string]interface{}{
			"count":        len(dependencies),
			"dependencies": dependencies,
		}

		jsonResponse, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize response: %w", err)
		}

		return mcp.NewToolResultText(string(jsonResponse)), nil
	}
}

// GetTracesSummaryHandler handles the jaeger_get_traces_summary tool.
func GetTracesSummaryHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		// Parse parameters
		serviceName := args["service"].(string)
		operation := args["operation"].(string)
		startTime := args["start_time"].(string)
		endTime := args["end_time"].(string)
		limit := int(args["limit"].(float64))
		minDuration := args["min_duration"].(string)

		// Set defaults
		if limit == 0 {
			limit = 20
		}
		if limit > 100 {
			limit = 100
		}
		if startTime == "" {
			startTime = time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
		}
		if endTime == "" {
			endTime = time.Now().Format(time.RFC3339)
		}

		// Build query parameters
		params := client.TraceQueryParameters{
			Service:     serviceName,
			Operation:   operation,
			StartTime:   startTime,
			EndTime:     endTime,
			Limit:       limit,
			MinDuration: minDuration,
		}

		// Search traces
		traces, err := service.GetClient().SearchTraces(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to search traces: %w", err)
		}

		// Build summary
		summaries := make([]map[string]interface{}, 0, len(traces))
		for _, trace := range traces {
			summary := map[string]interface{}{
				"trace_id":    trace.TraceID,
				"span_count":  len(trace.Spans),
				"duration_us": 0,
			}

			// Get service and operation from first span
			if len(trace.Spans) > 0 {
				summary["operation"] = trace.Spans[0].OperationName
				summary["duration_us"] = trace.Spans[0].Duration
			}

			// Get service from processes
			if len(trace.Processes) > 0 {
				summary["service"] = trace.Processes[0].ServiceName
			}

			summaries = append(summaries, summary)
		}

		// Serialize response
		result := map[string]interface{}{
			"count":  len(summaries),
			"traces": summaries,
		}

		jsonResponse, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize response: %w", err)
		}

		return mcp.NewToolResultText(string(jsonResponse)), nil
	}
}

// GetServicesSummaryHandler handles the jaeger_get_services_summary tool.
func GetServicesSummaryHandler(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get services
		services, err := service.GetClient().GetServices(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get services: %w", err)
		}

		// Build summary
		summaries := make([]map[string]interface{}, 0, len(services))
		for _, svc := range services {
			summary := map[string]interface{}{
				"name":            svc.Name,
				"operation_count": len(svc.Operations),
			}
			summaries = append(summaries, summary)
		}

		// Serialize response
		result := map[string]interface{}{
			"count":    len(summaries),
			"services": summaries,
		}

		jsonResponse, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize response: %w", err)
		}

		return mcp.NewToolResultText(string(jsonResponse)), nil
	}
}

package tools

import "github.com/mark3labs/mcp-go/mcp"

// GetTracesTool retrieves traces from Jaeger.
func GetTracesTool() mcp.Tool {
	return mcp.NewTool("jaeger_get_traces",
		mcp.WithDescription("Retrieve traces from Jaeger with optional filtering. Use this tool when you need to: analyze distributed tracing data, debug performance issues, understand request flow across microservices, or investigate latency problems. The tool returns trace data including spans, operations, and timing information."),
		mcp.WithString("service",
			mcp.Description("Filter traces by service name. If not specified, returns traces from all services.")),
		mcp.WithString("operation",
			mcp.Description("Filter traces by operation name. Requires service parameter.")),
		mcp.WithString("start_time",
			mcp.Description("Start time for trace search in RFC3339 format (e.g., '2024-01-01T00:00:00Z'). Defaults to 24 hours ago.")),
		mcp.WithString("end_time",
			mcp.Description("End time for trace search in RFC3339 format (e.g., '2024-01-02T00:00:00Z'). Defaults to now.")),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of traces to return (default: 20, max: 100).")),
		mcp.WithString("min_duration",
			mcp.Description("Minimum duration in microseconds (e.g., '1000000' for 1 second).")),
		mcp.WithString("max_duration",
			mcp.Description("Maximum duration in microseconds (e.g., '5000000' for 5 seconds).")),
	)
}

// GetTraceTool retrieves a specific trace by ID.
func GetTraceTool() mcp.Tool {
	return mcp.NewTool("jaeger_get_trace",
		mcp.WithDescription("Retrieve a specific trace by its trace ID. Use this tool when you need to: analyze a specific trace in detail, debug a particular request, or examine the full span tree for a trace. The trace ID is typically 32 hexadecimal characters."+
			"⚠️ PRIORITY: Use this to get detailed trace information after searching with jaeger_get_traces_summary."),
		mcp.WithString("trace_id",
			mcp.Required(),
			mcp.Description("The trace ID to retrieve. This is a 32-character hexadecimal string (e.g., '1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7'). You can find trace IDs from search results or logs.")),
	)
}

// GetServicesTool retrieves all services from Jaeger.
func GetServicesTool() mcp.Tool {
	return mcp.NewTool("jaeger_get_services",
		mcp.WithDescription("Retrieve all services registered in Jaeger. Use this tool when you need to: discover available services, understand service topology, or get a list of services for filtering traces. Returns service names and their operations."),
	)
}

// GetServiceOperationsTool retrieves operations for a specific service.
func GetServiceOperationsTool() mcp.Tool {
	return mcp.NewTool("jaeger_get_service_ops",
		mcp.WithDescription("Retrieve all operations (endpoints) for a specific service. Use this tool when you need to: understand what operations a service provides, filter traces by operation, or analyze service endpoints."),
		mcp.WithString("service",
			mcp.Required(),
			mcp.Description("The service name to get operations for. Use jaeger_get_services to get a list of available services.")),
	)
}

// SearchTracesTool searches for traces based on query parameters.
func SearchTracesTool() mcp.Tool {
	return mcp.NewTool("jaeger_search_traces",
		mcp.WithDescription("Search for traces in Jaeger with advanced filtering. Use this tool when you need to: find traces matching specific criteria, investigate performance issues, or analyze traces with specific tags or duration ranges."+
			"⚠️ PRIORITY: Use jaeger_get_traces_summary for faster, LLM-optimized results first."),
		mcp.WithString("service",
			mcp.Description("Filter traces by service name.")),
		mcp.WithString("operation",
			mcp.Description("Filter traces by operation name. Requires service parameter.")),
		mcp.WithString("start_time",
			mcp.Description("Start time for trace search in RFC3339 format.")),
		mcp.WithString("end_time",
			mcp.Description("End time for trace search in RFC3339 format.")),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of traces to return (default: 20, max: 100).")),
		mcp.WithString("min_duration",
			mcp.Description("Minimum duration in microseconds.")),
		mcp.WithString("max_duration",
			mcp.Description("Maximum duration in microseconds.")),
		mcp.WithObject("tags",
			mcp.Description("Filter traces by tags as key-value pairs (e.g., {'http.method': 'GET', 'http.status_code': '500'}).")),
	)
}

// GetDependenciesTool retrieves service dependencies.
func GetDependenciesTool() mcp.Tool {
	return mcp.NewTool("jaeger_get_dependencies",
		mcp.WithDescription("Retrieve service dependency graph from Jaeger. Use this tool when you need to: understand service relationships, visualize service topology, or analyze service dependencies and call patterns. Returns parent-child relationships with call counts."),
		mcp.WithString("start_time",
			mcp.Description("Start time for dependency analysis in RFC3339 format (e.g., '2024-01-01T00:00:00Z'). Defaults to 24 hours ago.")),
		mcp.WithString("end_time",
			mcp.Description("End time for dependency analysis in RFC3339 format (e.g., '2024-01-02T00:00:00Z'). Defaults to now.")),
	)
}

// GetTracesSummaryTool retrieves traces summary with minimal output (RECOMMENDED).
func GetTracesSummaryTool() mcp.Tool {
	return mcp.NewTool("jaeger_get_traces_summary",
		mcp.WithDescription("🎯 RECOMMENDED: List traces with minimal output for efficient discovery. Returns only essential fields (trace_id, service, operation, duration, span_count) with 70-85% smaller response size. Perfect for: quick trace browsing, finding slow traces, understanding trace landscape, or as first step before getting detailed info. Includes pagination for large collections."),
		mcp.WithString("service",
			mcp.Description("Filter traces by service name.")),
		mcp.WithString("operation",
			mcp.Description("Filter traces by operation name.")),
		mcp.WithString("start_time",
			mcp.Description("Start time for trace search in RFC3339 format.")),
		mcp.WithString("end_time",
			mcp.Description("End time for trace search in RFC3339 format.")),
		mcp.WithNumber("limit",
			mcp.Description("Maximum traces to return (default: 20, max: 100).")),
		mcp.WithString("min_duration",
			mcp.Description("Minimum duration in microseconds (e.g., '1000000' for 1 second).")),
	)
}

// GetServicesSummaryTool retrieves services summary with minimal output (RECOMMENDED).
func GetServicesSummaryTool() mcp.Tool {
	return mcp.NewTool("jaeger_get_services_summary",
		mcp.WithDescription("🎯 RECOMMENDED: List all services with minimal output for quick discovery. Returns only essential fields (name, operation_count) with 70-85% smaller response size. Perfect for: understanding available services, quick service inventory, or as first step before getting detailed operations."),
	)
}

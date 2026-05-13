package tools

import "github.com/mark3labs/mcp-go/mcp"

func TestConnectionTool() mcp.Tool {
	return mcp.NewTool("nacos_test_connection",
		mcp.WithDescription("Check whether the configured Nacos server is reachable and authentication works."),
	)
}

func ListNamespacesTool() mcp.Tool {
	return mcp.NewTool("nacos_list_namespaces",
		mcp.WithDescription("List namespaces from Nacos. Use this first when you need available namespace IDs before querying configs or services."),
	)
}

func ListConfigsSummaryTool() mcp.Tool {
	return mcp.NewTool("nacos_list_configs_summary",
		mcp.WithDescription("RECOMMENDED: List Nacos config entries in a compact summary view before fetching full config content."),
		mcp.WithString("namespace_id",
			mcp.Description("Optional namespace ID (tenant). Falls back to configured default namespace when available.")),
		mcp.WithString("group",
			mcp.Description("Optional config group. Falls back to configured default group when available.")),
		mcp.WithString("query",
			mcp.Description("Optional fuzzy search term for dataId or appName.")),
		mcp.WithNumber("page",
			mcp.Description("Page number starting at 1. Defaults to 1.")),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of config entries to return per page. Defaults to 20.")),
	)
}

func GetConfigTool() mcp.Tool {
	return mcp.NewTool("nacos_get_config",
		mcp.WithDescription("Get the full content of one Nacos configuration entry by dataId, group, and optional namespace."),
		mcp.WithString("data_id", mcp.Required(),
			mcp.Description("Config dataId to fetch.")),
		mcp.WithString("group",
			mcp.Description("Optional config group. Falls back to configured default group when available.")),
		mcp.WithString("namespace_id",
			mcp.Description("Optional namespace ID (tenant). Falls back to configured default namespace when available.")),
	)
}

func ListServicesSummaryTool() mcp.Tool {
	return mcp.NewTool("nacos_list_services_summary",
		mcp.WithDescription("RECOMMENDED: List Nacos services in a compact summary view before fetching full service details."),
		mcp.WithString("namespace_id",
			mcp.Description("Optional namespace ID. Falls back to configured default namespace when available.")),
		mcp.WithString("group_name",
			mcp.Description("Optional group name. Falls back to configured default group when available.")),
		mcp.WithString("service_name",
			mcp.Description("Optional fuzzy filter for service name.")),
		mcp.WithNumber("page",
			mcp.Description("Page number starting at 1. Defaults to 1.")),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of services to return per page. Defaults to 20.")),
	)
}

func GetServiceTool() mcp.Tool {
	return mcp.NewTool("nacos_get_service",
		mcp.WithDescription("Get one Nacos service including metadata, thresholds, selectors, and instance summary."),
		mcp.WithString("service_name", mcp.Required(),
			mcp.Description("Service name to fetch.")),
		mcp.WithString("group_name",
			mcp.Description("Optional group name. Falls back to configured default group when available.")),
		mcp.WithString("namespace_id",
			mcp.Description("Optional namespace ID. Falls back to configured default namespace when available.")),
	)
}

func ListInstancesTool() mcp.Tool {
	return mcp.NewTool("nacos_list_instances",
		mcp.WithDescription("List instances for one Nacos service. Useful for confirming registration, health, and cluster placement."),
		mcp.WithString("service_name", mcp.Required(),
			mcp.Description("Service name to inspect.")),
		mcp.WithString("group_name",
			mcp.Description("Optional group name. Falls back to configured default group when available.")),
		mcp.WithString("namespace_id",
			mcp.Description("Optional namespace ID. Falls back to configured default namespace when available.")),
		mcp.WithString("cluster_name",
			mcp.Description("Optional cluster name filter.")),
		mcp.WithBoolean("healthy_only",
			mcp.Description("Return only healthy instances when true. Defaults to false.")),
	)
}

func ListClusterNodesTool() mcp.Tool {
	return mcp.NewTool("nacos_list_cluster_nodes",
		mcp.WithDescription("List Nacos cluster nodes and their reported state."),
	)
}

func GetSystemMetricsTool() mcp.Tool {
	return mcp.NewTool("nacos_get_system_metrics",
		mcp.WithDescription("Get Nacos server metrics and runtime counters exposed by the operator metrics endpoint."),
	)
}

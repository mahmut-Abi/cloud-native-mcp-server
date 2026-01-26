// Package tools provides MCP tool definitions for the Helm service.
package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
)

// GetCacheClearTool returns the helm_clear_cache tool
func GetCacheClearTool() mcp.Tool {
	return mcp.NewTool("helm_clear_cache",
		mcp.WithDescription("Clear Helm cache to force fresh queries"),
	)
}

// GetCacheStatsTool returns the helm_cache_stats tool
func GetCacheStatsTool() mcp.Tool {
	return mcp.NewTool("helm_cache_stats",
		mcp.WithDescription("Get cache statistics including hit rate, misses, and evictions"),
	)
}

// GetQuickInfoTool returns the helm_quick_info tool
func GetQuickInfoTool() mcp.Tool {
	return mcp.NewTool("helm_quick_info",
		mcp.WithDescription("Get quick overview of all Helm releases in all namespaces (summary view)"),
	)
}

// GetReleaseSummaryTool returns the helm_get_release_summary tool
func GetReleaseSummaryTool() mcp.Tool {
	return mcp.NewTool("helm_get_release_summary",
		mcp.WithDescription("Get a brief summary of a Helm release (faster than detailed version)"),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("Release name")),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("Release namespace")),
	)
}

// GetListReleasesSummaryTool returns the helm_list_releases_summary tool
func GetListReleasesSummaryTool() mcp.Tool {
	return mcp.NewTool("helm_list_releases_summary",
		mcp.WithDescription("List all Helm releases with summary information (compact output)"),
		mcp.WithString("namespace",
			mcp.Description("Filter by namespace (optional)")),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of results (default 50)")),
		mcp.WithNumber("offset",
			mcp.Description("Pagination offset (default 0)")),
	)
}

// GetFindReleasesByChartTool returns the helm_find_releases_by_chart tool
func GetFindReleasesByChartTool() mcp.Tool {
	return mcp.NewTool("helm_find_releases_by_chart",
		mcp.WithDescription("Find all releases using a specific chart"),
		mcp.WithString("chart_name", mcp.Required(),
			mcp.Description("Chart name to search for")),
		mcp.WithString("chart_version",
			mcp.Description("Chart version (optional, filters by exact version)")),
	)
}

// GetListReleasesByNamespaceTool returns the helm_list_releases_in_namespace tool
func GetListReleasesByNamespaceTool() mcp.Tool {
	return mcp.NewTool("helm_list_releases_in_namespace",
		mcp.WithDescription("List all releases in a specific namespace with summary data"),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("Kubernetes namespace")),
	)
}

// GetFindBrokenReleases returns the helm_find_broken_releases tool
func GetFindBrokenReleases() mcp.Tool {
	return mcp.NewTool("helm_find_broken_releases",
		mcp.WithDescription("Find releases with failed or pending status"),
		mcp.WithString("namespace",
			mcp.Description("Filter by namespace (optional)")),
	)
}

// GetValidateReleaseTool returns the helm_validate_release tool
func GetValidateReleaseTool() mcp.Tool {
	return mcp.NewTool("helm_validate_release",
		mcp.WithDescription("Validate a release configuration"),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("Release name")),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("Release namespace")),
	)
}

// GetListReleasesPaginatedTool returns the helm_list_releases_paginated tool
func GetListReleasesPaginatedTool() mcp.Tool {
	return mcp.NewTool("helm_list_releases_paginated",
		mcp.WithDescription("⚠️ PRIORITY: Optimized for LLM efficiency: List Helm releases with pagination and summary output. 80-90% smaller than full listing. 🚀 Best for: quick browsing, resource discovery, health checks. 📋 Use cases: release_inventory, health_check, cluster_monitoring. 🔄 Workflow: use this tool to discover releases → use helm_get_release_summary for details → use helm_get_release for deep analysis when needed."),
		mcp.WithString("namespace",
			mcp.Description("Optional namespace filter. Omit for cluster-wide listing across all namespaces.")),
		mcp.WithString("status",
			mcp.Description("Filter by release status (deployed, failed, pending, superseded, uninstalled, uninstalling)")),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of releases to return (default: 50, max: 100). This enables server-side pagination to prevent context overflow.")),
		mcp.WithString("continueToken",
			mcp.Description("Pagination token from previous response to fetch the next page. When response indicates 'hasMore': true, use the provided 'continueToken' to get the next batch.")),
		mcp.WithString("includeLabels",
			mcp.Description("Optional comma-separated label keys to include in the summary output")),
	)
}

// GetReleaseHistoryPaginatedTool returns the helm_get_release_history_paginated tool
func GetReleaseHistoryPaginatedTool() mcp.Tool {
	return mcp.NewTool("helm_get_release_history_paginated",
		mcp.WithDescription("Get release history with pagination support. Returns revision history with summary information to prevent context overflow."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("The name of the release to get history for.")),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("The namespace where the release is deployed.")),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of revisions to return (default: 20, max: 50)")),
		mcp.WithString("continueToken",
			mcp.Description("Pagination token from previous response")),
		mcp.WithBoolean("includeStatus",
			mcp.Description("Include detailed status information for each revision (default: false)")),
	)
}

// GetRecentFailuresTool returns the helm_get_recent_failures tool
func GetRecentFailuresTool() mcp.Tool {
	return mcp.NewTool("helm_get_recent_failures",
		mcp.WithDescription("⚠️ PRIORITY: Get recent failed or problematic releases only. Optimized for troubleshooting - returns only releases with issues. 80-90% smaller than full listing. 🎯 Best for: problem diagnosis, troubleshooting, cluster health monitoring. 📋 Use cases: troubleshooting, health_check, failure_analysis."),
		mcp.WithString("namespace",
			mcp.Description("Optional namespace filter. Omit to check all namespaces.")),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of failed releases to return (default: 20, max: 50)")),
		mcp.WithBoolean("includePending",
			mcp.Description("Include pending releases in addition to failed ones (default: false)")),
	)
}

// GetClusterOverviewTool returns the helm_cluster_overview tool
func GetClusterOverviewTool() mcp.Tool {
	return mcp.NewTool("helm_cluster_overview",
		mcp.WithDescription("⚠️ PRIORITY: Get cluster-wide Helm overview with minimal output. Returns only essential statistics and summary across all namespaces. 95% smaller than detailed listings. 🎯 Best for: cluster overview, resource statistics, quick monitoring. 📋 Use cases: cluster_overview, resource_monitoring, capacity_planning."),
		mcp.WithBoolean("includeNodes",
			mcp.Description("Include node information in the overview (default: false)")),
		mcp.WithBoolean("includeStorage",
			mcp.Description("Include storage information in the overview (default: false)")),
	)
}

// GetFindReleasesByLabelsTool returns the helm_find_releases_by_labels tool
func GetFindReleasesByLabelsTool() mcp.Tool {
	return mcp.NewTool("helm_find_releases_by_labels",
		mcp.WithDescription("Find releases by label selectors with summary output"),
		mcp.WithString("labelSelector", mcp.Required(),
			mcp.Description("Label selector in format 'key=value' or 'key!=value'. Multiple selectors can be combined with commas.")),
		mcp.WithString("namespace",
			mcp.Description("Optional namespace filter")),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of results (default: 30, max: 100)")),
	)
}

// GetResourcesOfReleaseTool returns the helm_get_resources_of_release tool
func GetResourcesOfReleaseTool() mcp.Tool {
	return mcp.NewTool("helm_get_resources_of_release",
		mcp.WithDescription("Get summarized list of Kubernetes resources managed by a Helm release"),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("Release name")),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("Release namespace")),
		mcp.WithBoolean("includeStatus",
			mcp.Description("Include resource status information (default: false)")),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of resources to return (default: 50, max: 200)")),
	)
}

// HelmHealthCheckTool returns a tool for diagnosing Helm service issues
func HelmHealthCheckTool() mcp.Tool {
	return mcp.NewTool("helm_health_check",
		mcp.WithDescription("🩺 Diagnose Helm service health and configuration. Use this tool when Helm tools are not responding or returning errors. Checks: client initialization, Kubernetes connection, repository configuration, and cache status."),
		mcp.WithBoolean("checkClient",
			mcp.Description("Check Helm client initialization status (default: true)")),
		mcp.WithBoolean("checkKubernetes",
			mcp.Description("Check Kubernetes cluster connectivity (default: true)")),
		mcp.WithBoolean("checkRepositories",
			mcp.Description("Check Helm repository configuration (default: true)")),
		mcp.WithBoolean("checkCache",
			mcp.Description("Check cache status and statistics (default: true)")),
	)
}

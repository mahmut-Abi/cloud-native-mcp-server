package tools

import "github.com/mark3labs/mcp-go/mcp"

func HealthCheckTool() mcp.Tool {
	return mcp.NewTool("elasticsearch_health",
		mcp.WithDescription("Check Elasticsearch cluster health status"),
		mcp.WithString("level", mcp.Description("Health level: indices, cluster, or shards")))
}

func ListIndicesTool() mcp.Tool {
	return mcp.NewTool("elasticsearch_list_indices",
		mcp.WithDescription("List all indices in the cluster"))
}

func GetIndexStatsTool() mcp.Tool {
	return mcp.NewTool("elasticsearch_index_stats",
		mcp.WithDescription("Get index statistics"),
		mcp.WithString("index", mcp.Required(), mcp.Description("Index name")))
}

func GetNodesTool() mcp.Tool {
	return mcp.NewTool("elasticsearch_nodes",
		mcp.WithDescription("Get cluster nodes information"))
}

// GetIndicesSummaryTool returns tool definition for getting Elasticsearch indices summary
func GetIndicesSummaryTool() mcp.Tool {
	return mcp.NewTool("elasticsearch_indices_summary",
		mcp.WithDescription("List Elasticsearch indices summary (name, health, status, docs_count). 75-90% smaller output."),
	)
}

func GetInfoTool() mcp.Tool {
	return mcp.NewTool("elasticsearch_info",
		mcp.WithDescription("Get cluster information"))
}

// ⚠️ PRIORITY: Pagination and summary tools - LLM optimized version

// ListIndicesPaginatedTool returns tool definition for paginated indices listing
func ListIndicesPaginatedTool() mcp.Tool {
	return mcp.NewTool("elasticsearch_list_indices_paginated",
		mcp.WithDescription("⚠️ PRIORITY: Optimized for LLM efficiency: List Elasticsearch indices with pagination and summary output. 80-90% smaller than full listing."),
		mcp.WithString("continueToken", mcp.Description("Pagination token from previous response to fetch the next page. Empty for first page.")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of indices to return (default: 20, max: 100)")),
		mcp.WithString("indexPattern", mcp.Description("Index pattern filter (e.g., 'logs-*', 'metrics-*'). Use '*' for wildcard.")),
		mcp.WithBoolean("includeHealth", mcp.Description("Include health status in summary (adds minimal data). Default: false")),
	)
}

// GetNodesSummaryTool returns tool definition for getting Elasticsearch nodes summary
func GetNodesSummaryTool() mcp.Tool {
	return mcp.NewTool("elasticsearch_nodes_summary",
		mcp.WithDescription("⚠️ PRIORITY: Get Elasticsearch nodes summary (node_id, name, role, metrics). 75-90% smaller output. Optimized for LLM efficiency."),
		mcp.WithString("role", mcp.Description("Filter by node role: master, data, ingest, coordinating, or empty for all")),
		mcp.WithBoolean("includeMetrics", mcp.Description("Include basic performance metrics (CPU, memory). Default: false")),
		mcp.WithNumber("limit", mcp.Description("Maximum nodes to return (default: 50, max: 100)")),
	)
}

// GetClusterHealthSummaryTool returns tool definition for getting cluster health summary
func GetClusterHealthSummaryTool() mcp.Tool {
	return mcp.NewTool("elasticsearch_cluster_health_summary",
		mcp.WithDescription("⚠️ PRIORITY: Get Elasticsearch cluster health summary (status, nodes, active_shards, indices). Lightweight health overview. Optimized for monitoring."),
		mcp.WithString("level", mcp.Description("Health detail level: basic (default), detailed, or indices")),
		mcp.WithBoolean("includeIndices", mcp.Description("Include per-index health summary (adds data). Default: false")),
	)
}

// GetIndexDetailAdvancedTool returns tool definition for advanced index details
func GetIndexDetailAdvancedTool() mcp.Tool {
	return mcp.NewTool("elasticsearch_get_index_detail_advanced",
		mcp.WithDescription("🔍 Advanced index detail retrieval with enhanced formatting and optional components. Use when comprehensive analysis needed."),
		mcp.WithString("index", mcp.Required(), mcp.Description("Index name")),
		mcp.WithBoolean("includeMappings", mcp.Description("Include field mappings and schema. Default: false")),
		mcp.WithBoolean("includeSettings", mcp.Description("Include index settings and configuration. Default: false")),
		mcp.WithBoolean("includeStats", mcp.Description("Include detailed statistics and metrics. Default: true")),
		mcp.WithBoolean("includeSegments", mcp.Description("Include segment information. Default: false")),
		mcp.WithString("outputFormat", mcp.Description("Output format: structured (default), compact, or verbose")),
	)
}

// GetClusterDetailAdvancedTool returns tool definition for advanced cluster details
func GetClusterDetailAdvancedTool() mcp.Tool {
	return mcp.NewTool("elasticsearch_get_cluster_detail_advanced",
		mcp.WithDescription("🔍 Advanced cluster detail retrieval with comprehensive information. Use for deep cluster analysis."),
		mcp.WithBoolean("includeNodes", mcp.Description("Include detailed node information. Default: true")),
		mcp.WithBoolean("includeIndices", mcp.Description("Include comprehensive index overview. Default: false")),
		mcp.WithBoolean("includeSettings", mcp.Description("Include cluster-wide settings. Default: false")),
		mcp.WithBoolean("includeStats", mcp.Description("Include detailed cluster statistics. Default: true")),
		mcp.WithBoolean("includeShards", mcp.Description("Include shard allocation details. Default: false")),
		mcp.WithString("outputFormat", mcp.Description("Output format: structured (default), compact, or verbose")),
	)
}

// SearchIndicesTool returns tool definition for searching indices with filters
func SearchIndicesTool() mcp.Tool {
	return mcp.NewTool("elasticsearch_search_indices",
		mcp.WithDescription("🔍 Search Elasticsearch indices with advanced filters and pagination. Optimized for finding specific indices."),
		mcp.WithString("query", mcp.Description("Search query string (matches index names). Supports wildcards.")),
		mcp.WithString("healthStatus", mcp.Description("Filter by health status: green, yellow, red, or empty for all")),
		mcp.WithString("indexStatus", mcp.Description("Filter by index status: open, close, or empty for all")),
		mcp.WithNumber("minDocCount", mcp.Description("Minimum document count filter")),
		mcp.WithNumber("maxDocCount", mcp.Description("Maximum document count filter")),
		mcp.WithString("sortBy", mcp.Description("Sort by: name, docs, size, health (default: name)")),
		mcp.WithString("sortOrder", mcp.Description("Sort order: asc or desc (default: asc)")),
		mcp.WithNumber("limit", mcp.Description("Maximum results to return (default: 30, max: 200)")),
		mcp.WithString("continueToken", mcp.Description("Pagination token for next page")),
	)
}

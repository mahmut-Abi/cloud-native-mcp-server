package handlers

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/k8s-mcp-server/internal/constants"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/elasticsearch/client"
	optimize "github.com/mahmut-Abi/k8s-mcp-server/internal/util/performance"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/util/sanitize"
)

var (
	ErrMissingRequiredParam = errors.New("missing required parameter")
	ErrInvalidParameter     = errors.New("invalid parameter")
)

// Response size limits for Elasticsearch - similar to other optimizations
// Use constants from internal/constants package for consistency
const (
	defaultLimit      = constants.DefaultPageSize
	maxLimit          = constants.MaxPageSize
	warningLimit      = constants.WarningPageSize
	defaultLimitNodes = constants.DefaultPageSizeNodes
)

// Helper function to validate and parse limit parameter with warnings
func parseLimitWithWarnings(request mcp.CallToolRequest, toolName string) int64 {
	limit := int64(defaultLimit)
	if v, ok := request.GetArguments()["limit"]; ok {
		if f, ok := v.(float64); ok {
			limit = int64(f)
			if limit <= 0 || limit > maxLimit {
				if limit > maxLimit {
					logrus.WithField("requested", limit).WithField("max", maxLimit).Warn("Limit too high, resetting to safe maximum")
					limit = maxLimit
				} else {
					limit = defaultLimit
				}
			}
			if limit > warningLimit {
				logrus.WithFields(logrus.Fields{
					"tool":  toolName,
					"limit": limit,
				}).Warn("Large limit may cause context overflow, consider using summary tools or pagination")
			}
		}
	}
	return limit
}

// Helper function to validate required string parameter
func requireStringParam(request mcp.CallToolRequest, param string) (string, error) {
	value, ok := request.GetArguments()[param].(string)
	if !ok || value == "" {
		return "", fmt.Errorf("%w: %s", ErrMissingRequiredParam, param)
	}
	return value, nil
}

// Helper function to get optional string parameter
func getOptionalStringParam(request mcp.CallToolRequest, param string) string {
	value, _ := request.GetArguments()[param].(string)
	return value
}

// Helper function to get optional boolean parameter
func getOptionalBoolParam(request mcp.CallToolRequest, param string) *bool {
	if value, ok := request.GetArguments()[param].(bool); ok {
		return &value
	}
	return nil
}

// Helper function to marshal JSON response with size optimization
func marshalOptimizedResponse(data any, toolName string) (*mcp.CallToolResult, error) {
	jsonResponse, err := optimize.GlobalJSONPool.MarshalToBytes(data)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize response: %w", err)
	}

	// Add size warning for large responses
	if len(jsonResponse) > 100*1024 { // 100KB
		logrus.WithFields(logrus.Fields{
			"tool":      toolName,
			"sizeBytes": len(jsonResponse),
			"sizeKB":    len(jsonResponse) / 1024,
		}).Warn("Large response generated")
	}

	return mcp.NewToolResultText(string(jsonResponse)), nil
}

// Existing handlers (unchanged for compatibility)
func HandleHealthCheck(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Handling es_health tool")
		health, err := c.Health(ctx)
		if err != nil {
			logrus.WithError(err).Error("Failed to get cluster health")
			return mcp.NewToolResultError("Failed to get cluster health: " + err.Error()), nil
		}
		data, _ := optimize.GlobalJSONPool.MarshalToBytes(health)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func HandleListIndices(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Handling es_list_indices tool")
		indices, err := c.Indices(ctx)
		if err != nil {
			logrus.WithError(err).Error("Failed to list indices")
			return mcp.NewToolResultError("Failed to list indices: " + err.Error()), nil
		}
		data, _ := optimize.GlobalJSONPool.MarshalToBytes(indices)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func HandleIndexStats(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		indexName, ok := request.GetArguments()["index"].(string)
		if !ok || indexName == "" {
			return mcp.NewToolResultError("index parameter is required"), nil
		}
		// Sanitize index name to prevent injection
		indexName = sanitize.SanitizeFilterValue(indexName)
		logrus.WithField("index", indexName).Debug("Handling es_index_stats tool")
		stats, err := c.IndexStats(ctx, indexName)
		if err != nil {
			logrus.WithError(err).Error("Failed to get index stats")
			return mcp.NewToolResultError("Failed to get index stats: " + err.Error()), nil
		}
		data, _ := optimize.GlobalJSONPool.MarshalToBytes(stats)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func HandleNodes(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Handling es_nodes tool")
		nodes, err := c.Nodes(ctx)
		if err != nil {
			logrus.WithError(err).Error("Failed to get nodes")
			return mcp.NewToolResultError("Failed to get nodes: " + err.Error()), nil
		}
		data, _ := optimize.GlobalJSONPool.MarshalToBytes(nodes)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func HandleInfo(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Handling es_info tool")
		info, err := c.Info(ctx)
		if err != nil {
			logrus.WithError(err).Error("Failed to get cluster info")
			return mcp.NewToolResultError("Failed to get cluster info: " + err.Error()), nil
		}
		data, _ := optimize.GlobalJSONPool.MarshalToBytes(info)
		return mcp.NewToolResultText(string(data)), nil
	}
}

// ⚠️ PRIORITY: Optimized handlers for LLM efficiency

// HandleListIndicesPaginated handles paginated indices listing with LLM optimization
func HandleListIndicesPaginated(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		continueToken := getOptionalStringParam(request, "continueToken")
		limit := parseLimitWithWarnings(request, "elasticsearch_list_indices_paginated")
		indexPattern := getOptionalStringParam(request, "indexPattern")
		includeHealth := getOptionalBoolParam(request, "includeHealth")
		if includeHealth == nil {
			defaultHealth := false
			includeHealth = &defaultHealth
		}

		logrus.WithFields(logrus.Fields{
			"tool":          "elasticsearch_list_indices_paginated",
			"continueToken": continueToken,
			"limit":         limit,
			"indexPattern":  indexPattern,
			"includeHealth": *includeHealth,
		}).Debug("Handler invoked")

		indices, pagination, err := c.IndicesPaginated(ctx, continueToken, int(limit), indexPattern, *includeHealth)
		if err != nil {
			return nil, fmt.Errorf("failed to list indices paginated: %w", err)
		}

		response := map[string]interface{}{
			"indices":    indices,
			"count":      len(indices),
			"pagination": pagination,
			"metadata": map[string]interface{}{
				"tool":          "elasticsearch_list_indices_paginated",
				"indexPattern":  indexPattern,
				"includeHealth": *includeHealth,
				"optimizedFor":  "LLM efficiency",
			},
		}

		return marshalOptimizedResponse(response, "elasticsearch_list_indices_paginated")
	}
}

// HandleGetNodesSummary handles getting nodes summary with role filtering
func HandleGetNodesSummary(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		role := getOptionalStringParam(request, "role")
		includeMetrics := getOptionalBoolParam(request, "includeMetrics")
		if includeMetrics == nil {
			defaultMetrics := false
			includeMetrics = &defaultMetrics
		}
		limit := parseLimitWithWarnings(request, "elasticsearch_nodes_summary")
		if limit == 0 {
			limit = defaultLimitNodes
		}

		logrus.WithFields(logrus.Fields{
			"tool":           "elasticsearch_nodes_summary",
			"role":           role,
			"includeMetrics": *includeMetrics,
			"limit":          limit,
		}).Debug("Handler invoked")

		nodes, err := c.NodesSummary(ctx, role, *includeMetrics, int(limit))
		if err != nil {
			return nil, fmt.Errorf("failed to get nodes summary: %w", err)
		}

		response := map[string]interface{}{
			"nodes": nodes,
			"count": len(nodes),
			"filter": map[string]interface{}{
				"role":           role,
				"includeMetrics": *includeMetrics,
				"limit":          limit,
			},
			"metadata": map[string]interface{}{
				"tool":         "elasticsearch_nodes_summary",
				"optimizedFor": "LLM efficiency",
			},
		}

		return marshalOptimizedResponse(response, "elasticsearch_nodes_summary")
	}
}

// HandleGetClusterHealthSummary handles getting cluster health summary
func HandleGetClusterHealthSummary(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		level := getOptionalStringParam(request, "level")
		if level == "" {
			level = "basic"
		}
		includeIndices := getOptionalBoolParam(request, "includeIndices")
		if includeIndices == nil {
			defaultIndices := false
			includeIndices = &defaultIndices
		}

		logrus.WithFields(logrus.Fields{
			"tool":           "elasticsearch_cluster_health_summary",
			"level":          level,
			"includeIndices": *includeIndices,
		}).Debug("Handler invoked")

		health, err := c.ClusterHealthSummary(ctx, level, *includeIndices)
		if err != nil {
			return nil, fmt.Errorf("failed to get cluster health summary: %w", err)
		}

		response := map[string]interface{}{
			"health": health,
			"metadata": map[string]interface{}{
				"tool":           "elasticsearch_cluster_health_summary",
				"level":          level,
				"includeIndices": *includeIndices,
				"optimizedFor":   "monitoring and LLM efficiency",
			},
		}

		return marshalOptimizedResponse(response, "elasticsearch_cluster_health_summary")
	}
}

// HandleGetIndexDetailAdvanced handles getting advanced index details
func HandleGetIndexDetailAdvanced(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		indexName, err := requireStringParam(request, "index")
		if err != nil {
			return nil, err
		}

		includeMappings := getOptionalBoolParam(request, "includeMappings")
		if includeMappings == nil {
			defaultMappings := false
			includeMappings = &defaultMappings
		}
		includeSettings := getOptionalBoolParam(request, "includeSettings")
		if includeSettings == nil {
			defaultSettings := false
			includeSettings = &defaultSettings
		}
		includeStats := getOptionalBoolParam(request, "includeStats")
		if includeStats == nil {
			defaultStats := true
			includeStats = &defaultStats
		}
		includeSegments := getOptionalBoolParam(request, "includeSegments")
		if includeSegments == nil {
			defaultSegments := false
			includeSegments = &defaultSegments
		}
		outputFormat := getOptionalStringParam(request, "outputFormat")
		if outputFormat == "" {
			outputFormat = "structured"
		}

		logrus.WithFields(logrus.Fields{
			"tool":            "elasticsearch_get_index_detail_advanced",
			"index":           indexName,
			"includeMappings": *includeMappings,
			"includeSettings": *includeSettings,
			"includeStats":    *includeStats,
			"includeSegments": *includeSegments,
			"outputFormat":    outputFormat,
		}).Debug("Handler invoked")

		detail, err := c.GetIndexDetailAdvanced(ctx, indexName, *includeMappings, *includeSettings, *includeStats, *includeSegments, outputFormat)
		if err != nil {
			return nil, fmt.Errorf("failed to get index detail advanced: %w", err)
		}

		response := map[string]interface{}{
			"indexDetail": detail,
			"metadata": map[string]interface{}{
				"tool":         "elasticsearch_get_index_detail_advanced",
				"index":        indexName,
				"outputFormat": outputFormat,
				"optimizedFor": "comprehensive analysis",
			},
		}

		return marshalOptimizedResponse(response, "elasticsearch_get_index_detail_advanced")
	}
}

// HandleGetClusterDetailAdvanced handles getting advanced cluster details
func HandleGetClusterDetailAdvanced(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		includeNodes := getOptionalBoolParam(request, "includeNodes")
		if includeNodes == nil {
			defaultNodes := true
			includeNodes = &defaultNodes
		}
		includeIndices := getOptionalBoolParam(request, "includeIndices")
		if includeIndices == nil {
			defaultIndices := false
			includeIndices = &defaultIndices
		}
		includeSettings := getOptionalBoolParam(request, "includeSettings")
		if includeSettings == nil {
			defaultSettings := false
			includeSettings = &defaultSettings
		}
		includeStats := getOptionalBoolParam(request, "includeStats")
		if includeStats == nil {
			defaultStats := true
			includeStats = &defaultStats
		}
		includeShards := getOptionalBoolParam(request, "includeShards")
		if includeShards == nil {
			defaultShards := false
			includeShards = &defaultShards
		}
		outputFormat := getOptionalStringParam(request, "outputFormat")
		if outputFormat == "" {
			outputFormat = "structured"
		}

		logrus.WithFields(logrus.Fields{
			"tool":            "elasticsearch_get_cluster_detail_advanced",
			"includeNodes":    *includeNodes,
			"includeIndices":  *includeIndices,
			"includeSettings": *includeSettings,
			"includeStats":    *includeStats,
			"includeShards":   *includeShards,
			"outputFormat":    outputFormat,
		}).Debug("Handler invoked")

		detail, err := c.GetClusterDetailAdvanced(ctx, *includeNodes, *includeIndices, *includeSettings, *includeStats, *includeShards, outputFormat)
		if err != nil {
			return nil, fmt.Errorf("failed to get cluster detail advanced: %w", err)
		}

		response := map[string]interface{}{
			"clusterDetail": detail,
			"metadata": map[string]interface{}{
				"tool":         "elasticsearch_get_cluster_detail_advanced",
				"outputFormat": outputFormat,
				"optimizedFor": "deep cluster analysis",
			},
		}

		return marshalOptimizedResponse(response, "elasticsearch_get_cluster_detail_advanced")
	}
}

// HandleSearchIndices handles searching indices with filters and pagination
func HandleSearchIndices(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := getOptionalStringParam(request, "query")
		healthStatus := getOptionalStringParam(request, "healthStatus")
		indexStatus := getOptionalStringParam(request, "indexStatus")
		minDocCountStr := getOptionalStringParam(request, "minDocCount")
		maxDocCountStr := getOptionalStringParam(request, "maxDocCount")
		sortBy := getOptionalStringParam(request, "sortBy")
		if sortBy == "" {
			sortBy = "name"
		}
		sortOrder := getOptionalStringParam(request, "sortOrder")
		if sortOrder == "" {
			sortOrder = "asc"
		}
		limit := parseLimitWithWarnings(request, "elasticsearch_search_indices")
		continueToken := getOptionalStringParam(request, "continueToken")

		// Parse numeric parameters
		var minDocCount, maxDocCount int
		if minDocCountStr != "" {
			if val, err := strconv.Atoi(minDocCountStr); err == nil {
				minDocCount = val
			}
		}
		if maxDocCountStr != "" {
			if val, err := strconv.Atoi(maxDocCountStr); err == nil {
				maxDocCount = val
			}
		}

		logrus.WithFields(logrus.Fields{
			"tool":          "elasticsearch_search_indices",
			"query":         query,
			"healthStatus":  healthStatus,
			"indexStatus":   indexStatus,
			"minDocCount":   minDocCount,
			"maxDocCount":   maxDocCount,
			"sortBy":        sortBy,
			"sortOrder":     sortOrder,
			"limit":         limit,
			"continueToken": continueToken,
		}).Debug("Handler invoked")

		indices, pagination, err := c.SearchIndices(ctx, query, healthStatus, indexStatus, minDocCount, maxDocCount, sortBy, sortOrder, int(limit), continueToken)
		if err != nil {
			return nil, fmt.Errorf("failed to search indices: %w", err)
		}

		response := map[string]interface{}{
			"indices": indices,
			"count":   len(indices),
			"searchCriteria": map[string]interface{}{
				"query":        query,
				"healthStatus": healthStatus,
				"indexStatus":  indexStatus,
				"minDocCount":  minDocCount,
				"maxDocCount":  maxDocCount,
				"sortBy":       sortBy,
				"sortOrder":    sortOrder,
			},
			"pagination": pagination,
			"metadata": map[string]interface{}{
				"tool":         "elasticsearch_search_indices",
				"optimizedFor": "finding specific indices",
			},
		}

		return marshalOptimizedResponse(response, "elasticsearch_search_indices")
	}
}

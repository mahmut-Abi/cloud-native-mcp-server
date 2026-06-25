// Package kubernetes provides Kubernetes API integration for the MCP server.
// It implements tools for managing Kubernetes resources, including CRUD operations,
// resource discovery, and cluster interaction capabilities.
package kubernetes

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/cache"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/kubernetes/client"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/kubernetes/handlers"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/kubernetes/tools"
)

// Service implements the Kubernetes service for MCP server integration.
// It provides tools and handlers for interacting with Kubernetes clusters.
// The backend client is not stored — it is created per-request from HTTP headers.
type Service struct {
	enabled    bool              // Whether the service is enabled
	toolsCache *cache.ToolsCache // Cached tools to avoid recreation
}

// NewService creates a new Kubernetes service instance.
// The service is enabled by default and requires initialization before use.
func NewService() *Service {
	return &Service{
		enabled:    true, // Default enabled
		toolsCache: cache.NewToolsCache(),
	}
}

// Name returns the service identifier used for registration and logging.
func (s *Service) Name() string {
	return "kubernetes"
}

// Initialize configures the Kubernetes service with the provided application configuration.
func (s *Service) Initialize(cfg interface{}) error {
	logrus.Debug("Initializing Kubernetes service")
	// Kubernetes is always enabled by default; client is created per-request from headers.
	_ = cfg
	return nil
}

// GetTools returns all available Kubernetes MCP tools.
// Tools are only returned if the service is enabled.
// The tools include resource management, cluster interaction, and diagnostic capabilities.
func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled {
		return nil
	}

	// Use unified cache
	return s.toolsCache.Get(func() []mcp.Tool {
		return []mcp.Tool{
			// Core resource operations (optimized for LLM efficiency)
			tools.GetResourceSummaryTool(),
			tools.GetResourceTool(),
			tools.ListResourcesSummaryTool(), // Summary-first approach
			tools.ListResourcesTool(),
			tools.GetResourcesDetailTool(),

			// Full detail tools (use sparingly)
			tools.ListResourcesFullTool(),

			// Resource creation and management
			tools.CreateResourceTool(),
			tools.PatchResourceTool(),
			tools.DeleteResourceTool(),

			// Resource discovery and inspection
			tools.DescribeResourceTool(),
			tools.GetResourceDetailsTool(),
			tools.GetResourceDetailAdvancedTool(), // Advanced detail tool
			tools.GetAPIVersionsTool(),
			tools.GetAPIResourcesTool(),

			// Cluster operations
			tools.ScaleResourceTool(),
			tools.GetRolloutStatusTool(),
			tools.CordonNodeTool(),
			tools.UncordonNodeTool(),
			tools.DrainNodeTool(),
			tools.WaitForResourceTool(),
			tools.RestartWorkloadTool(),
			tools.PortForwardTool(),

			// Container and pod operations
			tools.ContainerLogsTool(),
			tools.ContainerExecTool(),
			tools.CheckPermissionsTool(),

			// Event monitoring (optimized vs detailed)
			tools.GetRecentEventsTool(), // Optimized for critical events
			tools.GetEventsTool(),       // Standard events
			tools.GetEventsDetailTool(), // Full detailed events

			// Resource monitoring
			tools.GetResourceUsageTool(),

			// Troubleshooting and diagnostics
			tools.GetUnhealthyResourcesTool(),
			tools.GetNodeConditionsTool(),
			tools.AnalyzeIssueTool(),

			// Search and discovery
			tools.SearchResourcesTool(),

			// Testing and validation
			tools.TestTool(),
		}
	})
}

// GetHandlers returns all tool handlers mapped to their respective tool names.
// Handlers are only returned if the service is enabled.
func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
	if !s.enabled {
		return nil
	}

	handlersMap := map[string]server.ToolHandlerFunc{
		// Core resource operations (optimized for LLM efficiency)
		"kubernetes_get_resource_summary":   s.wrapWithCache("kubernetes_get_resource_summary", handlers.HandleGetResourceSummary()),
		"kubernetes_get_resource":           handlers.HandleGetResource(),
		"kubernetes_list_resources_summary": s.wrapWithCache("kubernetes_list_resources_summary", handlers.HandleListResourcesSummary()), // Summary-first with cache
		"kubernetes_list_resources":         handlers.HandleListResources(),
		"kubernetes_get_resources_detail":   handlers.HandleGetResourcesDetail(),

		// Full detail tools (use sparingly)
		"kubernetes_list_resources_full": handlers.HandleListResourcesFull(),

		// Resource creation and management
		"kubernetes_create_resource": handlers.HandleCreateResource(),
		"kubernetes_patch_resource":  handlers.HandlePatchResource(),
		"kubernetes_delete_resource": handlers.HandleDeleteResource(),

		// Resource discovery and inspection
		"kubernetes_describe_resource":            handlers.HandleDescribeResource(),
		"kubernetes_get_resource_details":         handlers.HandleGetResourceDetails(),
		"kubernetes_get_resource_detail_advanced": handlers.HandleGetResourceDetailAdvanced(), // Advanced detail handler
		"kubernetes_get_api_versions":             s.wrapWithCache("kubernetes_get_api_versions", handlers.HandleGetAPIVersions()),
		"kubernetes_get_api_resources":            s.wrapWithCache("kubernetes_get_api_resources", handlers.HandleGetAPIResources()),

		// Cluster operations
		"kubernetes_scale_resource":     handlers.HandleScaleResource(),
		"kubernetes_get_rollout_status": handlers.HandleGetRolloutStatus(),
		"kubernetes_cordon_node":        handlers.HandleCordonNode(),
		"kubernetes_uncordon_node":      handlers.HandleUncordonNode(),
		"kubernetes_drain_node":         handlers.HandleDrainNode(),
		"kubernetes_wait_for_resource":  handlers.HandleWaitForResource(),
		"kubernetes_restart_workload":   handlers.HandleRestartWorkload(),
		"kubernetes_port_forward":       handlers.HandlePortForward(),

		// Container and pod operations
		"kubernetes_get_pod_logs":      handlers.HandleContainerLogs(),
		"kubernetes_pod_exec":          handlers.HandleContainerExec(),
		"kubernetes_check_permissions": s.wrapWithCache("kubernetes_check_permissions", handlers.HandleCheckPermissions()),

		// Event monitoring (optimized vs detailed)
		"kubernetes_get_recent_events": s.wrapWithCache("kubernetes_get_recent_events", handlers.HandleGetRecentEvents()), // Optimized for critical events with cache
		"kubernetes_get_events":        handlers.HandleGetEvents(),                                                        // Standard events
		"kubernetes_get_events_detail": handlers.HandleGetEventsDetail(),                                                  // Full detailed events

		// Resource monitoring
		"kubernetes_get_resource_usage": handlers.HandleGetResourceUsage(),

		// Troubleshooting and diagnostics
		"kubernetes_get_unhealthy_resources": handlers.HandleGetUnhealthyResources(),
		"kubernetes_get_node_conditions":     handlers.HandleGetNodeConditions(),
		"kubernetes_analyze_issue":           handlers.HandleAnalyzeIssue(),

		// Search and discovery
		"kubernetes_search_resources": handlers.HandleSearchResources(),

		// Testing and validation
		"kubernetes_test_tool": handlers.HandleTest(),
	}

	for name, handler := range handlersMap {
		handlersMap[name] = s.wrapWithToolErrors(name, handler)
	}

	return handlersMap
}

// wrapWithCache wraps a handler with caching if the tool is cacheable
func (s *Service) wrapWithCache(toolName string, handler server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Check if this tool should be cached
		if !handlers.IsToolCacheable(toolName) {
			return handler(ctx, request)
		}

		// Get TTL for this tool
		ttl := handlers.GetTTLForTool(toolName)

		// Filter parameters for cache key generation
		params := handlers.CacheParamsFilter(toolName, request.GetArguments())

		// Wrap the handler execution with cache
		result, _, err := handlers.CacheToolResponse(
			handlers.DefaultSmartCache,
			toolName,
			params,
			func() (interface{}, error) {
				return handler(ctx, request)
			},
			ttl,
		)

		if err != nil {
			return nil, err
		}

		// Convert result to CallToolResult
		if toolResult, ok := result.(*mcp.CallToolResult); ok {
			return toolResult, nil
		}

		// Fallback for other result types
		return mcp.NewToolResultText(fmt.Sprintf("%v", result)), nil
	}
}

func (s *Service) wrapWithToolErrors(toolName string, handler server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handler(ctx, request)
		if err != nil {
			logrus.WithError(err).WithField("tool", toolName).Warn("Tool execution failed")
			return mcp.NewToolResultError(err.Error()), nil
		}
		return result, nil
	}
}

// IsEnabled returns whether the service is enabled and ready for use.
func (s *Service) IsEnabled() bool {
	return s.enabled
}

// GetClient returns the underlying Kubernetes client.
// The client is no longer stored in the service — use client.FromContext instead.
func (s *Service) GetClient() *client.Client {
	return nil
}

// Package kubernetes provides Kubernetes API integration for the MCP server.
// It implements tools for managing Kubernetes resources, including CRUD operations,
// resource discovery, and cluster interaction capabilities.
package kubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	"github.com/mahmut-Abi/k8s-mcp-server/internal/config"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/cache"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/framework"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/kubernetes/client"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/kubernetes/handlers"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/kubernetes/tools"
)

// Service implements the Kubernetes service for MCP server integration.
// It provides tools and handlers for interacting with Kubernetes clusters.
type Service struct {
	client        *client.Client               // Kubernetes client for API operations
	enabled       bool                         // Whether the service is enabled
	toolsCache    *cache.ToolsCache            // Cached tools to avoid recreation
	initFramework *framework.CommonServiceInit // Common initialization framework
}

// NewService creates a new Kubernetes service instance.
// The service is enabled by default and requires initialization before use.
func NewService() *Service {
	// Create service enable checker - Kubernetes is always enabled by default
	checker := framework.NewServiceEnabled(
		func(cfg *config.AppConfig) bool { return true }, // Always enabled
		func(cfg *config.AppConfig) string { return cfg.Kubernetes.Kubeconfig },
	)

	// Create init configuration
	initConfig := &framework.InitConfig{
		Required:     false,
		URLValidator: func(url string) bool { return true }, // Kubeconfig can be empty (use in-cluster config)
		ClientBuilder: func(cfg *config.AppConfig) (interface{}, error) {
			opts := client.DefaultClientOptions()
			if cfg.Kubernetes.Kubeconfig != "" {
				opts.KubeconfigPath = cfg.Kubernetes.Kubeconfig
			}
			if cfg.Kubernetes.TimeoutSec > 0 {
				opts.Timeout = time.Duration(cfg.Kubernetes.TimeoutSec) * time.Second
			}
			if cfg.Kubernetes.QPS > 0 {
				opts.QPS = cfg.Kubernetes.QPS
			}
			if cfg.Kubernetes.Burst > 0 {
				opts.Burst = cfg.Kubernetes.Burst
			}
			return client.NewClientWithOptions(opts)
		},
	}

	return &Service{
		enabled:       true, // Default enabled
		toolsCache:    cache.NewToolsCache(),
		initFramework: framework.NewCommonServiceInit("Kubernetes", initConfig, checker),
	}
}

// Name returns the service identifier used for registration and logging.
func (s *Service) Name() string {
	return "kubernetes"
}

// Initialize configures the Kubernetes service with the provided application configuration.
// It creates and configures the underlying Kubernetes client with appropriate timeouts,
// rate limiting, and authentication settings.
func (s *Service) Initialize(cfg interface{}) error {
	return s.initFramework.Initialize(cfg,
		func(enabled bool) { s.enabled = enabled },
		func(clientIface interface{}) {
			if k8sClient, ok := clientIface.(*client.Client); ok {
				s.client = k8sClient
			}
		},
	)
}

// GetTools returns all available Kubernetes MCP tools.
// Tools are only returned if the service is enabled and properly initialized.
// The tools include resource management, cluster interaction, and diagnostic capabilities.
func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled || s.client == nil {
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
			tools.UpdateResourceTool(),
			tools.DeleteResourceTool(),

			// Resource discovery and inspection
			tools.DescribeResourceTool(),
			tools.GetResourceDetailsTool(),
			tools.GetResourceDetailAdvancedTool(), // Advanced detail tool
			tools.GetAPIVersionsTool(),
			tools.GetAPIResourcesTool(),

			// Cluster operations
			tools.ScaleResourceTool(),
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
// Handlers are only returned if the service is enabled and properly initialized.
func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
	if !s.enabled || s.client == nil {
		return nil
	}

	return map[string]server.ToolHandlerFunc{
		// Core resource operations (optimized for LLM efficiency)
		"kubernetes_get_resource_summary":   s.wrapWithCache("kubernetes_get_resource_summary", handlers.HandleGetResourceSummary(s.client)),
		"kubernetes_get_resource":           handlers.HandleGetResource(s.client),
		"kubernetes_list_resources_summary": s.wrapWithCache("kubernetes_list_resources_summary", handlers.HandleListResourcesSummary(s.client)), // Summary-first with cache
		"kubernetes_list_resources":         handlers.HandleListResources(s.client),
		"kubernetes_get_resources_detail":   handlers.HandleGetResourcesDetail(s.client),

		// Full detail tools (use sparingly)
		"kubernetes_list_resources_full": handlers.HandleListResourcesFull(s.client),

		// Resource creation and management
		"kubernetes_create_resource": handlers.HandleCreateResource(s.client),
		"kubernetes_update_resource": handlers.HandleUpdateResource(s.client),
		"kubernetes_delete_resource": handlers.HandleDeleteResource(s.client),

		// Resource discovery and inspection
		"kubernetes_describe_resource":            handlers.HandleDescribeResource(s.client),
		"kubernetes_get_resource_details":         handlers.HandleGetResourceDetails(s.client),
		"kubernetes_get_resource_detail_advanced": handlers.HandleGetResourceDetailAdvanced(s.client), // Advanced detail handler
		"kubernetes_get_api_versions":             s.wrapWithCache("kubernetes_get_api_versions", handlers.HandleGetAPIVersions(s.client)),
		"kubernetes_get_api_resources":            s.wrapWithCache("kubernetes_get_api_resources", handlers.HandleGetAPIResources(s.client)),

		// Cluster operations
		"kubernetes_scale_resource": handlers.HandleScaleResource(s.client),
		"kubernetes_port_forward":   handlers.HandlePortForward(s.client),

		// Container and pod operations
		"kubernetes_get_pod_logs":      handlers.HandleContainerLogs(s.client),
		"kubernetes_pod_exec":          handlers.HandleContainerExec(s.client),
		"kubernetes_check_permissions": s.wrapWithCache("kubernetes_check_permissions", handlers.HandleCheckPermissions(s.client)),

		// Event monitoring (optimized vs detailed)
		"kubernetes_get_recent_events": s.wrapWithCache("kubernetes_get_recent_events", handlers.HandleGetRecentEvents(s.client)), // Optimized for critical events with cache
		"kubernetes_get_events":        handlers.HandleGetEvents(s.client),                                                        // Standard events
		"kubernetes_get_events_detail": handlers.HandleGetEventsDetail(s.client),                                                  // Full detailed events

		// Resource monitoring
		"kubernetes_get_resource_usage": handlers.HandleGetResourceUsage(s.client),

		// Troubleshooting and diagnostics
		"kubernetes_get_unhealthy_resources": handlers.HandleGetUnhealthyResources(s.client),
		"kubernetes_get_node_conditions":     handlers.HandleGetNodeConditions(s.client),
		"kubernetes_analyze_issue":           handlers.HandleAnalyzeIssue(s.client),

		// Search and discovery
		"kubernetes_search_resources": handlers.HandleSearchResources(s.client),

		// Testing and validation
		"kubernetes_test_tool": handlers.HandleTest(s.client),
	}
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

// IsEnabled returns whether the service is enabled and ready for use.
// A service is considered enabled if it's marked as enabled and has a valid client.
func (s *Service) IsEnabled() bool {
	return s.enabled && s.client != nil
}

// GetClient returns the underlying Kubernetes client for advanced operations.
// This method is primarily used for testing and internal service communication.
func (s *Service) GetClient() *client.Client {
	return s.client
}

// Package helm provides Helm chart and release management integration for the MCP server.
// It implements tools for managing Helm releases, charts, repositories, and their integration
// with Kubernetes and Grafana.
package helm

import (
	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/cache"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/helm/client"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/helm/handlers"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/helm/tools"
)

// Service implements the Helm service for MCP server integration.
// The backend client is not stored — it is created per-request from HTTP headers.
type Service struct {
	enabled              bool              // Whether the service is enabled
	toolsCache           *cache.ToolsCache // Cached basic tools to avoid recreation
	additionalToolsCache *cache.ToolsCache // Cached additional tools to avoid recreation
}

// NewService creates a new Helm service instance.
func NewService() *Service {
	return &Service{
		enabled:              true, // Default enabled
		toolsCache:           cache.NewToolsCache(),
		additionalToolsCache: cache.NewToolsCache(),
	}
}

// Name returns the service identifier.
func (s *Service) Name() string {
	return "helm"
}

// Initialize configures the Helm service with the provided application configuration.
func (s *Service) Initialize(cfg interface{}) error {
	logrus.Debug("Initializing Helm service")
	// Helm is enabled by default; client is created per-request from headers.
	_ = cfg
	return nil
}

// GetTools returns basic Helm MCP tools.
func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled {
		return nil
	}

	// Use unified cache
	return s.toolsCache.Get(func() []mcp.Tool {
		return []mcp.Tool{
			// Optimized tools for LLM efficiency
			tools.GetListReleasesPaginatedTool(), // ⚠️ Recommended first choice
			tools.GetReleaseStatusTool(),         // ⚠️ Recommended for status checks
			tools.GetRecentFailuresTool(),        // ⚠️ Recommended for troubleshooting
			tools.GetClusterOverviewTool(),       // ⚠️ Recommended for overview

			// Standard tools (for specific use cases)
			tools.ListReleasesTool(),
			tools.GetReleaseTool(),
			tools.ListRepositoriesTool(),
			tools.GetReleaseValuesTool(),
			tools.GetReleaseManifestTool(),
			tools.GetReleaseHistoryTool(),
			tools.SearchChartsTool(),
			tools.GetChartInfoTool(),
			tools.TemplateChartTool(),
			tools.CompareReleaseVersionsTool(),
			tools.AddRepositoryTool(),
			tools.RemoveRepositoryTool(),
			tools.UpdateRepositoriesTool(),

			// Additional specialized tools
			tools.GetReleaseHistoryPaginatedTool(),
			tools.GetFindReleasesByLabelsTool(),
			tools.GetResourcesOfReleaseTool(),

			// Additional tools from tools_extension.go
			tools.GetCacheClearTool(),
			tools.GetCacheStatsTool(),
			tools.GetQuickInfoTool(),
			tools.GetReleaseSummaryTool(),
			tools.GetListReleasesSummaryTool(),
			tools.GetFindReleasesByChartTool(),
			tools.GetListReleasesByNamespaceTool(),
			tools.GetFindBrokenReleases(),
			tools.GetValidateReleaseTool(),
			tools.HelmHealthCheckTool(),
		}
	})
}

// GetHandlers returns tool handlers for basic operations.
func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
	if !s.enabled {
		return nil
	}

	return map[string]server.ToolHandlerFunc{
		// Optimized tools for LLM efficiency
		"helm_list_releases_paginated": handlers.HandleListReleasesPaginated(), // ⚠️ Recommended first choice
		"helm_get_release_status":      handlers.HandleGetReleaseStatus(),      // ⚠️ Recommended for status checks
		"helm_get_recent_failures":     handlers.HandleGetRecentFailures(),     // ⚠️ Recommended for troubleshooting
		"helm_cluster_overview":        handlers.HandleGetClusterOverview(),    // ⚠️ Recommended for overview

		// Standard tools (for specific use cases)
		"helm_list_releases":        handlers.HandleListReleases(),
		"helm_get_release":          handlers.HandleGetRelease(),
		"helm_list_repos":           handlers.HandleListRepositories(),
		"helm_get_release_values":   handlers.HandleGetReleaseValues(),
		"helm_get_release_manifest": handlers.HandleGetReleaseManifest(),
		"helm_get_release_history":  handlers.HandleGetReleaseHistory(),
		"helm_search_charts":        handlers.HandleSearchCharts(),
		"helm_get_chart_info":       handlers.HandleGetChartInfo(),
		"helm_template_chart":       handlers.HandleTemplateChart(),
		"helm_compare_revisions":    handlers.HandleCompareRevisions(),
		"helm_add_repository":       handlers.HandleAddRepository(),
		"helm_remove_repository":    handlers.HandleRemoveRepository(),
		"helm_update_repositories":  handlers.HandleUpdateRepositories(),

		// Additional specialized tools
		"helm_get_release_history_paginated": handlers.HandleGetReleaseHistoryPaginated(),
		"helm_find_releases_by_labels":       handlers.HandleFindReleasesByLabels(),
		"helm_get_resources_of_release":      handlers.HandleGetResourcesOfRelease(),

		// Additional handlers from tools_extension.go
		"helm_clear_cache":                handlers.HandleClearCache(),
		"helm_cache_stats":                handlers.HandleGetCacheStats(),
		"helm_quick_info":                 handlers.HandleGetQuickInfo(),
		"helm_get_release_summary":        handlers.HandleGetReleaseSummary(),
		"helm_list_releases_summary":      handlers.HandleGetListReleasesSummary(),
		"helm_find_releases_by_chart":     handlers.HandleFindReleasesByChart(),
		"helm_list_releases_in_namespace": handlers.HandleListReleasesByNamespace(),
		"helm_find_broken_releases":       handlers.HandleFindBrokenReleases(),
		"helm_validate_release":           handlers.HandleValidateRelease(),
		"helm_health_check":               handlers.HandleHelmHealthCheck(),
	}
}

// IsEnabled returns whether the service is enabled and ready for use.
func (s *Service) IsEnabled() bool {
	return s.enabled
}

// GetClient returns the underlying Helm client.
// The client is no longer stored in the service — use client.FromContext instead.
func (s *Service) GetClient() *client.Client {
	return nil
}

// GetAdditionalTools returns additional Helm tools for advanced operations.
func (s *Service) GetAdditionalTools() []mcp.Tool {
	if !s.enabled {
		return nil
	}

	// Use unified cache
	return s.additionalToolsCache.Get(func() []mcp.Tool {
		return []mcp.Tool{
			tools.InstallReleaseTool(),
			tools.UninstallReleaseTool(),
			tools.UpgradeReleaseTool(),
			tools.RollbackReleaseTool(),
		}
	})
}

// GetAdditionalHandlers returns handlers for additional tools.
func (s *Service) GetAdditionalHandlers() map[string]server.ToolHandlerFunc {
	if !s.enabled {
		return nil
	}

	return map[string]server.ToolHandlerFunc{
		"helm_install_release":   handlers.HandleInstallRelease(),
		"helm_uninstall_release": handlers.HandleUninstallRelease(),
		"helm_upgrade_release":   handlers.HandleUpgradeRelease(),
		"helm_rollback_release":  handlers.HandleRollbackRelease(),
	}
}

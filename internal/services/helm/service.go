// Package helm provides Helm chart and release management integration for the MCP server.
// It implements tools for managing Helm releases, charts, repositories, and their integration
// with Kubernetes and Grafana.
package helm

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/cache"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/helm/client"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/helm/handlers"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/helm/tools"
)

// Service implements the Helm service for MCP server integration.
type Service struct {
	client               *client.Client    // Helm client for managing releases and charts
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

	var opts *client.ClientOptions

	if cfg == nil {
		// Use default options if no config provided
		opts = &client.ClientOptions{}
	} else {
		appConfig, ok := cfg.(*config.AppConfig)
		if !ok {
			s.enabled = false
			return fmt.Errorf("invalid configuration type provided: %T", cfg)
		}

		// Build client options from config
		opts = &client.ClientOptions{
			KubeconfigPath: appConfig.Helm.KubeconfigPath,
			Namespace:      appConfig.Helm.Namespace,
		}

		// If no Helm kubeconfig specified, use the Kubernetes kubeconfig

		// Create optimizer from config
		optimizer := client.NewRepositoryOptimizer(
			appConfig.Helm.Mirrors,
			appConfig.Helm.TimeoutSec,
			appConfig.Helm.MaxRetries,
			appConfig.Helm.UseMirrors,
		)
		opts.Optimizer = optimizer
		if opts.KubeconfigPath == "" && appConfig.Kubernetes.Kubeconfig != "" {
			opts.KubeconfigPath = appConfig.Kubernetes.Kubeconfig
		}
	}

	// Create Helm client
	var err error
	s.client, err = client.NewClient(opts)
	if err != nil {
		logrus.Errorf("Failed to create Helm client: %v", err)
		s.enabled = false
		return err
	}

	return nil
}

// GetTools returns basic Helm MCP tools.

func (s *Service) GetTools() []mcp.Tool {

	if !s.enabled || s.client == nil {

		return nil

	}

	// Use unified cache

	return s.toolsCache.Get(func() []mcp.Tool {

		return []mcp.Tool{

			// Optimized tools for LLM efficiency

			tools.GetListReleasesPaginatedTool(), // ⚠️ Recommended first choice

			tools.GetReleaseStatusTool(), // ⚠️ Recommended for status checks

			tools.GetRecentFailuresTool(), // ⚠️ Recommended for troubleshooting

			tools.GetClusterOverviewTool(), // ⚠️ Recommended for overview

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

			tools.GetMirrorConfigurationTool(),

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
	if !s.enabled || s.client == nil {
		return nil
	}

	return map[string]server.ToolHandlerFunc{
		// Optimized tools for LLM efficiency
		"helm_list_releases_paginated": handlers.HandleListReleasesPaginated(s.client), // ⚠️ Recommended first choice
		"helm_get_release_status":      handlers.HandleGetReleaseStatus(s.client),      // ⚠️ Recommended for status checks
		"helm_get_recent_failures":     handlers.HandleGetRecentFailures(s.client),     // ⚠️ Recommended for troubleshooting
		"helm_cluster_overview":        handlers.HandleGetClusterOverview(s.client),    // ⚠️ Recommended for overview

		// Standard tools (for specific use cases)
		"helm_list_releases":            handlers.HandleListReleases(s.client),
		"helm_get_release":              handlers.HandleGetRelease(s.client),
		"helm_list_repos":               handlers.HandleListRepositories(s.client),
		"helm_get_release_values":       handlers.HandleGetReleaseValues(s.client),
		"helm_get_release_manifest":     handlers.HandleGetReleaseManifest(s.client),
		"helm_get_release_history":      handlers.HandleGetReleaseHistory(s.client),
		"helm_search_charts":            handlers.HandleSearchCharts(s.client),
		"helm_get_chart_info":           handlers.HandleGetChartInfo(s.client),
		"helm_template_chart":           handlers.HandleTemplateChart(s.client),
		"helm_compare_revisions":        handlers.HandleCompareRevisions(s.client),
		"helm_add_repository":           handlers.HandleAddRepository(s.client),
		"helm_remove_repository":        handlers.HandleRemoveRepository(s.client),
		"helm_update_repositories":      handlers.HandleUpdateRepositories(s.client),
		"helm_get_mirror_configuration": handlers.HandleGetMirrorConfiguration(s.client),

		// Additional specialized tools
		"helm_get_release_history_paginated": handlers.HandleGetReleaseHistoryPaginated(s.client),
		"helm_find_releases_by_labels":       handlers.HandleFindReleasesByLabels(s.client),
		"helm_get_resources_of_release":      handlers.HandleGetResourcesOfRelease(s.client),

		// Additional handlers from tools_extension.go
		"helm_clear_cache":                handlers.HandleClearCache(s.client),
		"helm_cache_stats":                handlers.HandleGetCacheStats(s.client),
		"helm_quick_info":                 handlers.HandleGetQuickInfo(s.client),
		"helm_get_release_summary":        handlers.HandleGetReleaseSummary(s.client),
		"helm_list_releases_summary":      handlers.HandleGetListReleasesSummary(s.client),
		"helm_find_releases_by_chart":     handlers.HandleFindReleasesByChart(s.client),
		"helm_list_releases_in_namespace": handlers.HandleListReleasesByNamespace(s.client),
		"helm_find_broken_releases":       handlers.HandleFindBrokenReleases(s.client),
		"helm_validate_release":           handlers.HandleValidateRelease(s.client),
		"helm_health_check":               handlers.HandleHelmHealthCheck(s.client),
	}
}

// IsEnabled returns whether the service is enabled and ready for use.
func (s *Service) IsEnabled() bool {
	return s.enabled && s.client != nil
}

// GetClient returns the underlying Helm client.
func (s *Service) GetClient() *client.Client {
	return s.client
}

// GetAdditionalTools returns additional Helm tools for advanced operations.
func (s *Service) GetAdditionalTools() []mcp.Tool {
	if !s.enabled || s.client == nil {
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
	if !s.enabled || s.client == nil {
		return nil
	}

	return map[string]server.ToolHandlerFunc{
		"helm_install_release":   handlers.HandleInstallRelease(s.client),
		"helm_uninstall_release": handlers.HandleUninstallRelease(s.client),
		"helm_upgrade_release":   handlers.HandleUpgradeRelease(s.client),
		"helm_rollback_release":  handlers.HandleRollbackRelease(s.client),
	}
}

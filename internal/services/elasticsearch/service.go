package elasticsearch

import (
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/cache"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/elasticsearch/client"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/elasticsearch/handlers"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/elasticsearch/tools"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/framework"
)

// Service implements the Elasticsearch service for MCP server integration.
// It provides tools and handlers for interacting with Elasticsearch instances.
// The backend client is not stored — it is created per-request from HTTP headers.
type Service struct {
	enabled       bool                         // Whether the service is enabled
	toolsCache    *cache.ToolsCache            // Cached tools to avoid recreation
	initFramework *framework.CommonServiceInit // Common initialization framework
}

func NewService() *Service {
	// Create service enable checker
	checker := framework.NewServiceEnabled(
		func(cfg *config.AppConfig) bool { return cfg.Elasticsearch.Enabled },
		func(cfg *config.AppConfig) string {
			if len(cfg.Elasticsearch.Addresses) > 0 {
				return cfg.Elasticsearch.Addresses[0]
			}
			return cfg.Elasticsearch.Address
		},
	)

	// Create init configuration
	initConfig := &framework.InitConfig{
		Required:     false,
		URLValidator: framework.SimpleURLValidator,
		ClientBuilder: func(cfg *config.AppConfig) (interface{}, error) {
			addresses := cfg.Elasticsearch.Addresses
			if len(addresses) == 0 && cfg.Elasticsearch.Address != "" {
				addresses = []string{cfg.Elasticsearch.Address}
			}

			timeout := 30 * time.Second
			if cfg.Elasticsearch.TimeoutSec > 0 {
				timeout = time.Duration(cfg.Elasticsearch.TimeoutSec) * time.Second
			}

			return client.NewClient(&client.ClientOptions{
				Addresses:     addresses,
				Username:      cfg.Elasticsearch.Username,
				Password:      cfg.Elasticsearch.Password,
				BearerToken:   cfg.Elasticsearch.BearerToken,
				APIKey:        cfg.Elasticsearch.APIKey,
				TLSSkipVerify: cfg.Elasticsearch.TLSSkipVerify,
				Timeout:       timeout,
			})
		},
	}

	return &Service{
		enabled:       false, // Default disabled until configured
		toolsCache:    cache.NewToolsCache(),
		initFramework: framework.NewCommonServiceInit("Elasticsearch", initConfig, checker),
	}
}

func (s *Service) Name() string {
	return "elasticsearch"
}

// Initialize configures the Elasticsearch service with the provided application configuration.
// The backend client is created per-request from HTTP headers (see client/config.go).
func (s *Service) Initialize(cfg interface{}) error {
	return s.initFramework.Initialize(cfg,
		func(enabled bool) { s.enabled = enabled },
		func(_ interface{}) {
			// Backend client is created per-request from HTTP headers.
			// The backend auth handler was registered in client/config.go init().
		},
	)
}

func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled {
		return nil
	}

	// Use unified cache
	return s.toolsCache.Get(func() []mcp.Tool {
		// Legacy tools (maintained for compatibility)
		legacyTools := []mcp.Tool{
			tools.HealthCheckTool(),
			tools.ListIndicesTool(),
			tools.GetIndexStatsTool(),
			tools.GetNodesTool(),
			tools.GetInfoTool(),
		}

		// ⚠️ PRIORITY: New optimized tools for LLM efficiency
		optimizedTools := []mcp.Tool{
			tools.GetIndicesSummaryTool(),
			tools.GetNodesSummaryTool(),
			tools.GetClusterHealthSummaryTool(),
			tools.ListIndicesPaginatedTool(),
			tools.GetIndexDetailAdvancedTool(),
			tools.GetClusterDetailAdvancedTool(),
			tools.SearchIndicesTool(),
		}

		// Combine all tools - optimized tools first
		return append(optimizedTools, legacyTools...)
	})
}

func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
	if !s.enabled {
		return nil
	}

	// Legacy handlers (maintained for compatibility)
	legacyHandlers := map[string]server.ToolHandlerFunc{
		"elasticsearch_health":       handlers.HandleHealthCheck(),
		"elasticsearch_list_indices": handlers.HandleListIndices(),
		"elasticsearch_index_stats":  handlers.HandleIndexStats(),
		"elasticsearch_nodes":        handlers.HandleNodes(),
		"elasticsearch_info":         handlers.HandleInfo(),
	}

	// New optimized handlers for LLM efficiency
	optimizedHandlers := map[string]server.ToolHandlerFunc{
		// Summary tools
		"elasticsearch_indices_summary":        handlers.HandleListIndicesPaginated(),
		"elasticsearch_nodes_summary":          handlers.HandleGetNodesSummary(),
		"elasticsearch_cluster_health_summary": handlers.HandleGetClusterHealthSummary(),

		// Advanced tools
		"elasticsearch_list_indices_paginated":      handlers.HandleListIndicesPaginated(),
		"elasticsearch_get_index_detail_advanced":   handlers.HandleGetIndexDetailAdvanced(),
		"elasticsearch_get_cluster_detail_advanced": handlers.HandleGetClusterDetailAdvanced(),
		"elasticsearch_search_indices":              handlers.HandleSearchIndices(),
	}

	// Combine all handlers
	allHandlers := make(map[string]server.ToolHandlerFunc)
	for k, v := range optimizedHandlers {
		allHandlers[k] = v
	}
	for k, v := range legacyHandlers {
		allHandlers[k] = v
	}

	return allHandlers
}

func (s *Service) IsEnabled() bool {
	return s.enabled
}

func (s *Service) GetClient() *client.Client {
	return nil
}

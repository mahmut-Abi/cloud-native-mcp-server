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

type Service struct {
	client        *client.Client               // Elasticsearch HTTP client for API operations
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

func (s *Service) Initialize(cfg interface{}) error {
	return s.initFramework.Initialize(cfg,
		func(enabled bool) { s.enabled = enabled },
		func(clientIface interface{}) {
			if esClient, ok := clientIface.(*client.Client); ok {
				s.client = esClient
			}
		},
	)
}

func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled || s.client == nil {
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
	if !s.enabled || s.client == nil {
		return nil
	}

	// Legacy handlers (maintained for compatibility)
	legacyHandlers := map[string]server.ToolHandlerFunc{
		"elasticsearch_health":       handlers.HandleHealthCheck(s.client),
		"elasticsearch_list_indices": handlers.HandleListIndices(s.client),
		"elasticsearch_index_stats":  handlers.HandleIndexStats(s.client),
		"elasticsearch_nodes":        handlers.HandleNodes(s.client),
		"elasticsearch_info":         handlers.HandleInfo(s.client),
	}

	// ⚠️ PRIORITY: New optimized handlers for LLM efficiency
	optimizedHandlers := map[string]server.ToolHandlerFunc{
		// Summary tools
		"elasticsearch_indices_summary":        handlers.HandleListIndicesPaginated(s.client), // Same handler for summary tool
		"elasticsearch_nodes_summary":          handlers.HandleGetNodesSummary(s.client),
		"elasticsearch_cluster_health_summary": handlers.HandleGetClusterHealthSummary(s.client),

		// Advanced tools
		"elasticsearch_list_indices_paginated":      handlers.HandleListIndicesPaginated(s.client),
		"elasticsearch_get_index_detail_advanced":   handlers.HandleGetIndexDetailAdvanced(s.client),
		"elasticsearch_get_cluster_detail_advanced": handlers.HandleGetClusterDetailAdvanced(s.client),
		"elasticsearch_search_indices":              handlers.HandleSearchIndices(s.client),
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
	return s.enabled && s.client != nil
}

func (s *Service) GetClient() *client.Client {
	return s.client
}

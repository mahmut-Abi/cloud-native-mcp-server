// Package alertmanager provides Alertmanager integration for the MCP server.
// It implements tools for managing alerts, silences, and notification receivers.
package alertmanager

import (
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/alertmanager/client"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/alertmanager/handlers"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/alertmanager/tools"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/cache"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/framework"
)

// Service implements the Alertmanager service for MCP server integration.
// It provides tools and handlers for interacting with Alertmanager instances.
type Service struct {
	client        *client.Client               // Alertmanager HTTP client for API operations
	enabled       bool                         // Whether the service is enabled
	toolsCache    *cache.ToolsCache            // Cached tools to avoid recreation
	initFramework *framework.CommonServiceInit // Common initialization framework
}

// NewService creates a new Alertmanager service instance.
// The service is disabled by default and requires initialization before use.
func NewService() *Service {
	// Create service enable checker
	checker := framework.NewServiceEnabled(
		func(cfg *config.AppConfig) bool { return cfg.Alertmanager.Enabled },
		func(cfg *config.AppConfig) string {
			if cfg.Alertmanager.Address != "" {
				return cfg.Alertmanager.Address
			}
			return "http://localhost:9093" // Default address
		},
	)

	// Create init configuration
	initConfig := &framework.InitConfig{
		Required:     false,
		URLValidator: framework.SimpleURLValidator,
		ClientBuilder: func(cfg *config.AppConfig) (interface{}, error) {
			opts := client.DefaultClientOptions()
			if cfg.Alertmanager.Address != "" {
				opts.Address = cfg.Alertmanager.Address
			}
			if cfg.Alertmanager.TimeoutSec > 0 {
				opts.Timeout = time.Duration(cfg.Alertmanager.TimeoutSec) * time.Second
			}
			if cfg.Alertmanager.Username != "" {
				opts.Username = cfg.Alertmanager.Username
			}
			if cfg.Alertmanager.Password != "" {
				opts.Password = cfg.Alertmanager.Password
			}
			if cfg.Alertmanager.BearerToken != "" {
				opts.BearerToken = cfg.Alertmanager.BearerToken
			}
			opts.TLSSkipVerify = cfg.Alertmanager.TLSSkipVerify
			if cfg.Alertmanager.TLSCertFile != "" {
				opts.TLSCertFile = cfg.Alertmanager.TLSCertFile
			}
			if cfg.Alertmanager.TLSKeyFile != "" {
				opts.TLSKeyFile = cfg.Alertmanager.TLSKeyFile
			}

			return client.NewClientWithOptions(opts)
		},
	}

	return &Service{
		enabled:       false, // Default disabled until configured
		toolsCache:    cache.NewToolsCache(),
		initFramework: framework.NewCommonServiceInit("AlertManager", initConfig, checker),
	}
}

// Name returns the service identifier used for registration and logging.
func (s *Service) Name() string {
	return "alertmanager"
}

// Initialize configures the Alertmanager service with the provided application configuration.
// It uses the common service framework for standardized initialization.
func (s *Service) Initialize(cfg interface{}) error {
	return s.initFramework.Initialize(cfg,
		func(enabled bool) { s.enabled = enabled },
		func(clientIface interface{}) {
			if alertmanagerClient, ok := clientIface.(*client.Client); ok {
				s.client = alertmanagerClient
			}
		},
	)
}

// GetTools returns all available Alertmanager MCP tools.
// Tools are only returned if the service is enabled and properly initialized.
// The tools include alert management, silence operations, and receiver testing.
func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled || s.client == nil {
		return nil
	}

	// Use unified cache
	return s.toolsCache.Get(func() []mcp.Tool {
		// Legacy tools (maintained for compatibility)
		legacyTools := []mcp.Tool{
			tools.GetStatusTool(),
			tools.GetAlertsTool(),
			tools.GetAlertGroupsTool(),
			tools.GetSilencesTool(),
			tools.CreateSilenceTool(),
			tools.DeleteSilenceTool(),
			tools.GetReceiversTool(),
			tools.TestReceiverTool(),
			tools.QueryAlertsTool(),
		}

		// ⚠️ PRIORITY: New optimized tools for LLM efficiency
		optimizedTools := []mcp.Tool{
			tools.GetAlertsSummaryTool(),
			tools.GetSilencesSummaryTool(),
			tools.GetAlertGroupsPaginatedTool(),
			tools.GetSilencesPaginatedTool(),
			tools.GetReceiversSummaryTool(),
			tools.QueryAlertsAdvancedTool(),
			tools.GetHealthStatusSummaryTool(),
		}

		// Combine all tools - optimized tools first for better visibility
		return append(optimizedTools, legacyTools...)
	})
}

// GetHandlers returns all tool handlers mapped to their respective tool names.
// Handlers are only returned if the service is enabled and properly initialized.
func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
	if !s.enabled || s.client == nil {
		return nil
	}

	// Legacy handlers (maintained for compatibility)
	legacyHandlers := map[string]server.ToolHandlerFunc{
		"alertmanager_get_status":       handlers.HandleGetStatus(s.client),
		"alertmanager_get_alerts":       handlers.HandleGetAlerts(s.client),
		"alertmanager_get_alert_groups": handlers.HandleGetAlertGroups(s.client),
		"alertmanager_get_silences":     handlers.HandleGetSilences(s.client),
		"alertmanager_create_silence":   handlers.HandleCreateSilence(s.client),
		"alertmanager_delete_silence":   handlers.HandleDeleteSilence(s.client),
		"alertmanager_get_receivers":    handlers.HandleGetReceivers(s.client),
		"alertmanager_test_receiver":    handlers.HandleTestReceiver(s.client),
		"alertmanager_query_alerts":     handlers.HandleQueryAlerts(s.client),
	}

	// ⚠️ PRIORITY: New optimized handlers for LLM efficiency
	optimizedHandlers := map[string]server.ToolHandlerFunc{
		"alertmanager_alerts_summary":         handlers.HandleAlertsSummary(s.client),
		"alertmanager_silences_summary":       handlers.HandleSilencesSummary(s.client),
		"alertmanager_alert_groups_paginated": handlers.HandleAlertGroupsPaginated(s.client),
		"alertmanager_silences_paginated":     handlers.HandleSilencesPaginated(s.client),
		"alertmanager_receivers_summary":      handlers.HandleReceiversSummary(s.client),
		"alertmanager_query_alerts_advanced":  handlers.HandleQueryAlertsAdvanced(s.client),
		"alertmanager_health_summary":         handlers.HandleHealthSummary(s.client),
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

// IsEnabled returns whether the service is enabled and ready for use.
// A service is considered enabled if it's marked as enabled and has a valid client.
func (s *Service) IsEnabled() bool {
	return s.enabled && s.client != nil
}

// GetClient returns the underlying Alertmanager client for advanced operations.
// This method is primarily used for testing and internal service communication.
func (s *Service) GetClient() *client.Client {
	return s.client
}

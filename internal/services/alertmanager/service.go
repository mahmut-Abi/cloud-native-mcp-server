// Package alertmanager provides Alertmanager integration for the MCP server.
// It implements tools for managing alerts, silences, and notification receivers.
package alertmanager

import (

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
// The backend client is not stored — it is created per-request from HTTP headers.
type Service struct {
	enabled       bool                         // Whether the service is enabled
	toolsCache    *cache.ToolsCache            // Cached tools to avoid recreation
	initFramework *framework.CommonServiceInit // Common initialization framework
}

// NewService creates a new Alertmanager service instance.
// The service is disabled by default and requires initialization before use.
func NewService() *Service {
	// Create service enable checker
	checker := framework.NewServiceEnabled(
		func(cfg *config.AppConfig) bool { return true },
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
		ClientBuilder: nil,
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

// GetTools returns all available Alertmanager MCP tools.
// Tools are only returned if the service is enabled and properly initialized.
// The tools include alert management, silence operations, and receiver testing.
func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled {
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
	if !s.enabled {
		return nil
	}

	// Legacy handlers (maintained for compatibility)
	legacyHandlers := map[string]server.ToolHandlerFunc{
		"alertmanager_get_status":       handlers.HandleGetStatus(),
		"alertmanager_get_alerts":       handlers.HandleGetAlerts(),
		"alertmanager_get_alert_groups": handlers.HandleGetAlertGroups(),
		"alertmanager_get_silences":     handlers.HandleGetSilences(),
		"alertmanager_create_silence":   handlers.HandleCreateSilence(),
		"alertmanager_delete_silence":   handlers.HandleDeleteSilence(),
		"alertmanager_get_receivers":    handlers.HandleGetReceivers(),
		"alertmanager_test_receiver":    handlers.HandleTestReceiver(),
		"alertmanager_query_alerts":     handlers.HandleQueryAlerts(),
	}

	// ⚠️ PRIORITY: New optimized handlers for LLM efficiency
	optimizedHandlers := map[string]server.ToolHandlerFunc{
		"alertmanager_alerts_summary":         handlers.HandleAlertsSummary(),
		"alertmanager_silences_summary":       handlers.HandleSilencesSummary(),
		"alertmanager_alert_groups_paginated": handlers.HandleAlertGroupsPaginated(),
		"alertmanager_silences_paginated":     handlers.HandleSilencesPaginated(),
		"alertmanager_receivers_summary":      handlers.HandleReceiversSummary(),
		"alertmanager_query_alerts_advanced":  handlers.HandleQueryAlertsAdvanced(),
		"alertmanager_health_summary":         handlers.HandleHealthSummary(),
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
func (s *Service) IsEnabled() bool {
	return s.enabled
}

// GetClient returns the underlying Alertmanager client for advanced operations.
func (s *Service) GetClient() *client.Client {
	return nil
}

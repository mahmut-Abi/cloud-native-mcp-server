// Package prometheus provides Prometheus monitoring integration for the MCP server.
// It implements tools for querying Prometheus metrics, managing targets, and monitoring alerts.
package prometheus

import (
	"context"
	"encoding/json"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	"github.com/mahmut-Abi/k8s-mcp-server/internal/config"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/cache"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/framework"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/prometheus/client"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/prometheus/handlers"
	"github.com/mahmut-Abi/k8s-mcp-server/internal/services/prometheus/tools"
)

// Service implements the Prometheus service for MCP server integration.
// It provides tools and handlers for interacting with Prometheus instances.
type Service struct {
	client        *client.Client               // Prometheus HTTP client for API operations
	enabled       bool                         // Whether the service is enabled
	toolsCache    *cache.ToolsCache            // Cached tools to avoid recreation
	initFramework *framework.CommonServiceInit // Common initialization framework
}

// NewService creates a new Prometheus service instance.
// The service is disabled by default and requires initialization before use.
func NewService() *Service {
	// Create service enable checker
	checker := framework.NewServiceEnabled(
		func(cfg *config.AppConfig) bool { return cfg.Prometheus.Enabled },
		func(cfg *config.AppConfig) string { return cfg.Prometheus.Address },
	)

	// Create init configuration
	initConfig := &framework.InitConfig{
		Required:     false,
		URLValidator: func(url string) bool { return url != "" },
		ClientBuilder: func(cfg *config.AppConfig) (interface{}, error) {
			opts := &client.ClientOptions{
				Address:       cfg.Prometheus.Address,
				Username:      cfg.Prometheus.Username,
				Password:      cfg.Prometheus.Password,
				BearerToken:   cfg.Prometheus.BearerToken,
				Timeout:       30 * time.Second,
				TLSSkipVerify: cfg.Prometheus.TLSSkipVerify,
				TLSCertFile:   cfg.Prometheus.TLSCertFile,
				TLSKeyFile:    cfg.Prometheus.TLSKeyFile,
				TLSCAFile:     cfg.Prometheus.TLSCAFile,
			}

			if cfg.Prometheus.TimeoutSec > 0 {
				opts.Timeout = time.Duration(cfg.Prometheus.TimeoutSec) * time.Second
			}

			return client.NewClient(opts)
		},
	}

	return &Service{
		enabled:       false, // Default disabled until configured
		toolsCache:    cache.NewToolsCache(),
		initFramework: framework.NewCommonServiceInit("Prometheus", initConfig, checker),
	}
}

// Name returns the service identifier used for registration and logging.
func (s *Service) Name() string {
	return "prometheus"
}

// Initialize configures the Prometheus service with the provided application configuration.
// It uses the common service framework for standardized initialization.
func (s *Service) Initialize(cfg interface{}) error {
	return s.initFramework.Initialize(cfg,
		func(enabled bool) { s.enabled = enabled },
		func(clientIface interface{}) {
			if prometheusClient, ok := clientIface.(*client.Client); ok {
				s.client = prometheusClient
			}
		},
	)
}

// GetTools returns all available Prometheus MCP tools.
// Tools are only returned if the service is enabled and properly initialized.
// The tools include query operations, target monitoring, and alert management capabilities.
func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled || s.client == nil {
		return nil
	}

	// Use unified cache
	return s.toolsCache.Get(func() []mcp.Tool {
		return []mcp.Tool{
			// Query operations
			tools.QueryTool(),
			tools.QueryRangeTool(),

			// Target and service discovery
			tools.GetTargetsTool(),

			// Alert and rule management
			tools.GetAlertsTool(),
			tools.GetRulesTool(),

			// Label and series operations
			tools.GetLabelNamesTool(),
			tools.GetLabelValuesTool(),
			tools.GetSeriesTool(),

			// Connection and health
			tools.TestConnectionTool(),
		}
	})
}

// GetToolsEnhanced returns all available Prometheus MCP tools with enhancements.
func (s *Service) GetToolsEnhanced() []mcp.Tool {
	if !s.enabled || s.client == nil {
		return nil
	}

	result := s.GetTools()

	// Add enhanced tools
	result = append(result, tools.GetServerInfoTool())
	result = append(result, tools.GetMetricsMetadataTool())
	result = append(result, tools.GetTargetMetadataTool())

	// Add TSDB tools
	result = append(result, tools.GetTSDBStatsTool())
	result = append(result, tools.GetTSDBStatusTool())
	result = append(result, tools.GetRuntimeInfoTool())
	result = append(result, tools.CreateSnapshotTool())
	result = append(result, tools.GetWALReplayStatusTool())

	// Add summary tools (optimized versions)
	result = append(result, tools.GetTargetsSummaryTool())
	result = append(result, tools.GetAlertsSummaryTool())
	result = append(result, tools.GetRulesSummaryTool())

	return result
}

// GetHandlers returns all tool handlers mapped to their respective tool names.
// Handlers are only returned if the service is enabled and properly initialized.
func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
	if !s.enabled || s.client == nil {
		return nil
	}

	return map[string]server.ToolHandlerFunc{
		// Query operations
		"prometheus_query":       handlers.HandleQuery(s.client),
		"prometheus_query_range": handlers.HandleQueryRange(s.client),

		// Target and service discovery
		"prometheus_get_targets": handlers.HandleGetTargets(s.client),

		// Alert and rule management
		"prometheus_get_alerts": handlers.HandleGetAlerts(s.client),
		"prometheus_get_rules":  handlers.HandleGetRules(s.client),

		// Label and series operations
		"prometheus_get_label_names":  handlers.HandleGetLabelNames(s.client),
		"prometheus_get_label_values": handlers.HandleGetLabelValues(s.client),
		"prometheus_get_series":       handlers.HandleGetSeries(s.client),

		// Connection and health
		"prometheus_test_connection": handlers.HandleTestConnection(s.client),
	}
}

// GetHandlersEnhanced returns enhanced handlers.
func (s *Service) GetHandlersEnhanced() map[string]server.ToolHandlerFunc {
	if !s.enabled || s.client == nil {
		return nil
	}

	h := s.GetHandlers()
	h["prometheus_get_server_info"] = handlers.HandleGetServerInfo(s.client)
	h["prometheus_get_metrics_metadata"] = handlers.HandleGetMetricsMetadata(s.client)
	h["prometheus_get_target_metadata"] = handlers.HandleGetTargetMetadata(s.client)

	// Add TSDB handlers
	h["prometheus_get_tsdb_stats"] = handlers.HandleGetTSDBStats(s.client)
	h["prometheus_get_tsdb_status"] = handlers.HandleGetTSDBStatus(s.client)
	h["prometheus_get_runtime_info"] = handlers.HandleGetRuntimeInfo(s.client)
	h["prometheus_create_snapshot"] = handlers.HandleCreateSnapshot(s.client)
	h["prometheus_get_wal_replay_status"] = handlers.HandleGetWALReplayStatus(s.client)

	// Add summary tool handlers (optimized versions)
	h["prometheus_targets_summary"] = handlers.HandleGetTargetsSummary(s.client)
	h["prometheus_alerts_summary"] = handlers.HandleGetAlertsSummary(s.client)
	h["prometheus_rules_summary"] = handlers.HandleGetRulesSummary(s.client)

	return h
}

// IsEnabled returns whether the service is enabled and ready for use.
// A service is considered enabled if it's marked as enabled and has a valid client.
func (s *Service) IsEnabled() bool {
	return s.enabled && s.client != nil
}

// GetClient returns the underlying Prometheus client for advanced operations.
// This method is primarily used for testing and internal service communication.
func (s *Service) GetClient() *client.Client {
	return s.client
}

// GetResources returns all available Prometheus MCP resources.
// Resources provide static/semi-static data that gives context to AI models.
func (s *Service) GetResources() []mcp.Resource {
	if !s.enabled || s.client == nil {
		return nil
	}

	return []mcp.Resource{
		{
			URI:         "prometheus://targets",
			Name:        "Prometheus Targets",
			Description: "Current Prometheus scrape targets and their health status",
			MIMEType:    "application/json",
		},
		{
			URI:         "prometheus://rules",
			Name:        "Prometheus Rules",
			Description: "Current Prometheus recording and alerting rules",
			MIMEType:    "application/json",
		},
		{
			URI:         "prometheus://alerts",
			Name:        "Prometheus Alerts",
			Description: "Currently active Prometheus alerts",
			MIMEType:    "application/json",
		},
		{
			URI:         "prometheus://label-names",
			Name:        "Prometheus Label Names",
			Description: "Available label names in the Prometheus instance",
			MIMEType:    "application/json",
		},
	}
}

// GetResourceHandlers returns all resource handlers mapped to their respective resource URIs.
// Handlers are only returned if the service is enabled and properly initialized.
func (s *Service) GetResourceHandlers() map[string]server.ResourceHandlerFunc {
	if !s.enabled || s.client == nil {
		return nil
	}

	return map[string]server.ResourceHandlerFunc{
		"prometheus://targets":     s.handleTargetsResource,
		"prometheus://rules":       s.handleRulesResource,
		"prometheus://alerts":      s.handleAlertsResource,
		"prometheus://label-names": s.handleLabelNamesResource,
	}
}

// handleTargetsResource provides current Prometheus targets as a resource.
func (s *Service) handleTargetsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	targets, err := s.client.GetTargets(ctx, "")
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "Failed to fetch Prometheus targets: " + err.Error(),
			},
		}, nil
	}

	targetsJSON, err := json.MarshalIndent(targets, "", "  ")
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "Failed to serialize targets data: " + err.Error(),
			},
		}, nil
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:  request.Params.URI,
			Text: string(targetsJSON),
		},
	}, nil
}

// handleRulesResource provides current Prometheus rules as a resource.
func (s *Service) handleRulesResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	rules, err := s.client.GetRules(ctx, "")
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "Failed to fetch Prometheus rules: " + err.Error(),
			},
		}, nil
	}

	rulesJSON, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "Failed to serialize rules data: " + err.Error(),
			},
		}, nil
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:  request.Params.URI,
			Text: string(rulesJSON),
		},
	}, nil
}

// handleAlertsResource provides current Prometheus alerts as a resource.
func (s *Service) handleAlertsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	alerts, err := s.client.GetAlerts(ctx)
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "Failed to fetch Prometheus alerts: " + err.Error(),
			},
		}, nil
	}

	alertsJSON, err := json.MarshalIndent(alerts, "", "  ")
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "Failed to serialize alerts data: " + err.Error(),
			},
		}, nil
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:  request.Params.URI,
			Text: string(alertsJSON),
		},
	}, nil
}

// handleLabelNamesResource provides available label names as a resource.
func (s *Service) handleLabelNamesResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	labelNames, err := s.client.GetLabelNames(ctx, nil, nil)
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "Failed to fetch Prometheus label names: " + err.Error(),
			},
		}, nil
	}

	labelNamesJSON, err := json.MarshalIndent(labelNames, "", "  ")
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "Failed to serialize label names data: " + err.Error(),
			},
		}, nil
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:  request.Params.URI,
			Text: string(labelNamesJSON),
		},
	}, nil
}

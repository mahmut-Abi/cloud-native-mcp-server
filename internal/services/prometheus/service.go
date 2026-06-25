// Package prometheus provides Prometheus monitoring integration for the MCP server.
// It implements tools for querying Prometheus metrics, managing targets, and monitoring alerts.
package prometheus

import (
	"context"
	"encoding/json"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/cache"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/framework"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/prometheus/client"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/prometheus/handlers"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/prometheus/tools"
)

// Service implements the Prometheus service for MCP server integration.
// It provides tools and handlers for interacting with Prometheus instances.
// The backend client is not stored — it is created per-request from HTTP headers.
type Service struct {
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

// GetTools returns all available Prometheus MCP tools.
// Tools are only returned if the service is enabled and properly initialized.
// The tools include query operations, target monitoring, and alert management capabilities.
func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled {
		return nil
	}

	// Use unified cache
	return s.toolsCache.Get(func() []mcp.Tool {
		return []mcp.Tool{
			// Query operations
			tools.QueryTool(),
			tools.QueryRangeTool(),

			// Target and service discovery
			tools.GetTargetsSummaryTool(),
			tools.GetTargetsTool(),

			// Alert and rule management
			tools.GetAlertsSummaryTool(),
			tools.GetAlertsTool(),
			tools.GetRulesSummaryTool(),
			tools.GetRulesTool(),

			// Label and series operations
			tools.GetLabelNamesTool(),
			tools.GetLabelValuesTool(),
			tools.GetSeriesTool(),

			// Metadata and diagnostics
			tools.GetMetricsMetadataTool(),
			tools.GetTargetMetadataTool(),
			tools.GetTSDBStatsTool(),
			tools.GetTSDBStatusTool(),
			tools.CreateSnapshotTool(),
			tools.GetWALReplayStatusTool(),

			// Connection and health
			tools.TestConnectionTool(),
			tools.GetServerInfoTool(),
			tools.GetRuntimeInfoTool(),
		}
	})
}

// GetHandlers returns all tool handlers mapped to their respective tool names.
// Handlers are only returned if the service is enabled and properly initialized.
func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
	if !s.enabled {
		return nil
	}

	handlersMap := map[string]server.ToolHandlerFunc{
		// Query operations
		"prometheus_query":       handlers.HandleQuery(),
		"prometheus_query_range": handlers.HandleQueryRange(),

		// Target and service discovery
		"prometheus_targets_summary": handlers.HandleGetTargetsSummary(),
		"prometheus_get_targets":     handlers.HandleGetTargets(),

		// Alert and rule management
		"prometheus_alerts_summary": handlers.HandleGetAlertsSummary(),
		"prometheus_get_alerts":     handlers.HandleGetAlerts(),
		"prometheus_rules_summary":  handlers.HandleGetRulesSummary(),
		"prometheus_get_rules":      handlers.HandleGetRules(),

		// Label and series operations
		"prometheus_get_label_names":  handlers.HandleGetLabelNames(),
		"prometheus_get_label_values": handlers.HandleGetLabelValues(),
		"prometheus_get_series":       handlers.HandleGetSeries(),

		// Metadata and diagnostics
		"prometheus_get_metrics_metadata":  handlers.HandleGetMetricsMetadata(),
		"prometheus_get_target_metadata":   handlers.HandleGetTargetMetadata(),
		"prometheus_get_tsdb_stats":        handlers.HandleGetTSDBStats(),
		"prometheus_get_tsdb_status":       handlers.HandleGetTSDBStatus(),
		"prometheus_create_snapshot":       handlers.HandleCreateSnapshot(),
		"prometheus_get_wal_replay_status": handlers.HandleGetWALReplayStatus(),

		// Connection and health
		"prometheus_test_connection":  handlers.HandleTestConnection(),
		"prometheus_get_server_info":  handlers.HandleGetServerInfo(),
		"prometheus_get_runtime_info": handlers.HandleGetRuntimeInfo(),
	}

	for name, handler := range handlersMap {
		handlersMap[name] = s.wrapToolErrors(name, handler)
	}

	return handlersMap
}

func (s *Service) IsEnabled() bool {
	return s.enabled
}

func (s *Service) GetClient() *client.Client {
	return nil // Backend client is created per-request from HTTP headers
}

func (s *Service) wrapToolErrors(toolName string, handler server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handler(ctx, request)
		if err != nil {
			logrus.WithError(err).WithField("tool", toolName).Warn("Tool execution failed")
			return mcp.NewToolResultError(err.Error()), nil
		}
		return result, nil
	}
}

// GetResources returns all available Prometheus MCP resources.
// Resources provide static/semi-static data that gives context to AI models.
func (s *Service) GetResources() []mcp.Resource {
	if !s.enabled {
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
	if !s.enabled {
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
	c, err := client.FromContext(ctx)
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "Failed to create Prometheus client: " + err.Error(),
			},
		}, nil
	}
	targets, err := c.GetTargets(ctx, "")
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
	c, err := client.FromContext(ctx)
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "Failed to create Prometheus client: " + err.Error(),
			},
		}, nil
	}
	rules, err := c.GetRules(ctx, "")
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
	c, err := client.FromContext(ctx)
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "Failed to create Prometheus client: " + err.Error(),
			},
		}, nil
	}
	alerts, err := c.GetAlerts(ctx)
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
	c, err := client.FromContext(ctx)
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "Failed to create Prometheus client: " + err.Error(),
			},
		}, nil
	}
	labelNames, err := c.GetLabelNames(ctx, nil, nil)
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

// Package opentelemetry provides OpenTelemetry integration for the MCP server.
// It implements tools for querying OpenTelemetry Collector metrics, traces, and logs.
package opentelemetry

import (
	"context"
	"encoding/json"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/cache"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/framework"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/opentelemetry/client"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/opentelemetry/handlers"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/opentelemetry/tools"
)

// Service implements the OpenTelemetry service for MCP server integration.
// It provides tools and handlers for interacting with OpenTelemetry Collector instances.
// The backend client is not stored — it is created per-request from HTTP headers.
type Service struct {
	enabled       bool                         // Whether the service is enabled
	toolsCache    *cache.ToolsCache            // Cached tools to avoid recreation
	initFramework *framework.CommonServiceInit // Common initialization framework
}

// NewService creates a new OpenTelemetry service instance.
// The service is disabled by default and requires initialization before use.
func NewService() *Service {
	// Create service enable checker
	checker := framework.NewServiceEnabled(
		func(cfg *config.AppConfig) bool { return cfg.OpenTelemetry.Enabled },
		func(cfg *config.AppConfig) string { return cfg.OpenTelemetry.Address },
	)

	// Create init configuration
	initConfig := &framework.InitConfig{
		Required:     false,
		URLValidator: func(url string) bool { return url != "" },
		ClientBuilder: func(cfg *config.AppConfig) (interface{}, error) {
			opts := &client.ClientOptions{
				Address:       cfg.OpenTelemetry.Address,
				Username:      cfg.OpenTelemetry.Username,
				Password:      cfg.OpenTelemetry.Password,
				BearerToken:   cfg.OpenTelemetry.BearerToken,
				Timeout:       30 * time.Second,
				TLSSkipVerify: cfg.OpenTelemetry.TLSSkipVerify,
				TLSCertFile:   cfg.OpenTelemetry.TLSCertFile,
				TLSKeyFile:    cfg.OpenTelemetry.TLSKeyFile,
				TLSCAFile:     cfg.OpenTelemetry.TLSCAFile,
			}

			if cfg.OpenTelemetry.TimeoutSec > 0 {
				opts.Timeout = time.Duration(cfg.OpenTelemetry.TimeoutSec) * time.Second
			}

			return client.NewClient(opts)
		},
	}

	return &Service{
		enabled:       false, // Default disabled until configured
		toolsCache:    cache.NewToolsCache(),
		initFramework: framework.NewCommonServiceInit("OpenTelemetry", initConfig, checker),
	}
}

// Name returns the service identifier used for registration and logging.
func (s *Service) Name() string {
	return "opentelemetry"
}

// Initialize configures the OpenTelemetry service with the provided application configuration.
// It uses the common service framework for standardized initialization.
func (s *Service) Initialize(cfg interface{}) error {
	return s.initFramework.Initialize(cfg,
		func(enabled bool) { s.enabled = enabled },
		func(_ interface{}) {
			// Backend client is created per-request from HTTP headers.
			// The backend auth handler was registered in client/config.go init().
		},
	)
}

// GetTools returns all available OpenTelemetry MCP tools.
// Tools are only returned if the service is enabled and properly initialized.
func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled {
		return nil
	}

	// Use unified cache
	return s.toolsCache.Get(func() []mcp.Tool {
		return []mcp.Tool{
			// Metrics operations
			tools.GetMetricsTool(),
			tools.QueryMetricsTool(),

			// Traces operations
			tools.GetTracesTool(),
			tools.QueryTracesTool(),

			// Logs operations
			tools.GetLogsTool(),
			tools.QueryLogsTool(),

			// Health and status
			tools.GetHealthTool(),
			tools.GetStatusTool(),

			// Configuration
			tools.GetConfigTool(),
			tools.GetConfigSummaryTool(),
			tools.GetCollectorSummaryTool(),
			tools.AnalyzePipelineStatusTool(),
		}
	})
}

// GetHandlers returns all tool handlers mapped to their respective tool names.
// Handlers are only returned if the service is enabled and properly initialized.
func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
	if !s.enabled {
		return nil
	}

	return map[string]server.ToolHandlerFunc{
		// Metrics operations
		"opentelemetry_get_metrics":   handlers.HandleGetMetrics(),
		"opentelemetry_query_metrics": handlers.HandleQueryMetrics(),

		// Traces operations
		"opentelemetry_get_traces":   handlers.HandleGetTraces(),
		"opentelemetry_query_traces": handlers.HandleQueryTraces(),

		// Logs operations
		"opentelemetry_get_logs":   handlers.HandleGetLogs(),
		"opentelemetry_query_logs": handlers.HandleQueryLogs(),

		// Health and status
		"opentelemetry_get_health": handlers.HandleGetHealth(),
		"opentelemetry_get_status": handlers.HandleGetStatus(),

		// Configuration
		"opentelemetry_get_config":              handlers.HandleGetConfig(),
		"opentelemetry_get_config_summary":      handlers.HandleGetConfigSummary(),
		"opentelemetry_get_collector_summary":   handlers.HandleGetCollectorSummary(),
		"opentelemetry_analyze_pipeline_status": handlers.HandleAnalyzePipelineStatus(),
	}
}

// IsEnabled returns whether the service is enabled and ready for use.
func (s *Service) IsEnabled() bool {
	return s.enabled
}

// GetClient returns the underlying OpenTelemetry client for advanced operations.
func (s *Service) GetClient() *client.Client {
	return nil
}

// GetResources returns all available OpenTelemetry MCP resources.
func (s *Service) GetResources() []mcp.Resource {
	if !s.enabled {
		return nil
	}

	return []mcp.Resource{
		{
			URI:         "opentelemetry://metrics",
			Name:        "OpenTelemetry Metrics",
			Description: "Current metrics from OpenTelemetry Collector",
			MIMEType:    "application/json",
		},
		{
			URI:         "opentelemetry://traces",
			Name:        "OpenTelemetry Traces",
			Description: "Recent traces from OpenTelemetry Collector",
			MIMEType:    "application/json",
		},
		{
			URI:         "opentelemetry://logs",
			Name:        "OpenTelemetry Logs",
			Description: "Recent logs from OpenTelemetry Collector",
			MIMEType:    "application/json",
		},
		{
			URI:         "opentelemetry://health",
			Name:        "OpenTelemetry Health",
			Description: "OpenTelemetry Collector health status",
			MIMEType:    "application/json",
		},
	}
}

// GetResourceHandlers returns all resource handlers mapped to their respective resource URIs.
func (s *Service) GetResourceHandlers() map[string]server.ResourceHandlerFunc {
	if !s.enabled {
		return nil
	}

	return map[string]server.ResourceHandlerFunc{
		"opentelemetry://metrics": s.handleMetricsResource,
		"opentelemetry://traces":  s.handleTracesResource,
		"opentelemetry://logs":    s.handleLogsResource,
		"opentelemetry://health":  s.handleHealthResource,
	}
}

// handleMetricsResource provides current metrics as a resource.
func (s *Service) handleMetricsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	c, err := client.FromContext(ctx)
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "OpenTelemetry client not available: " + err.Error(),
			},
		}, nil
	}
	metrics, err := c.GetMetrics(ctx, nil, nil, nil)
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "Failed to fetch OpenTelemetry metrics: " + err.Error(),
			},
		}, nil
	}

	metricsJSON, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "Failed to serialize metrics data: " + err.Error(),
			},
		}, nil
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:  request.Params.URI,
			Text: string(metricsJSON),
		},
	}, nil
}

// handleTracesResource provides recent traces as a resource.
func (s *Service) handleTracesResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	c, err := client.FromContext(ctx)
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "OpenTelemetry client not available: " + err.Error(),
			},
		}, nil
	}
	traces, err := c.GetTraces(ctx, nil, nil, nil, nil, nil)
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "Failed to fetch OpenTelemetry traces: " + err.Error(),
			},
		}, nil
	}

	tracesJSON, err := json.MarshalIndent(traces, "", "  ")
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "Failed to serialize traces data: " + err.Error(),
			},
		}, nil
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:  request.Params.URI,
			Text: string(tracesJSON),
		},
	}, nil
}

// handleLogsResource provides recent logs as a resource.
func (s *Service) handleLogsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	c, err := client.FromContext(ctx)
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "OpenTelemetry client not available: " + err.Error(),
			},
		}, nil
	}
	logs, err := c.GetLogs(ctx, nil, nil, nil, nil, nil)
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "Failed to fetch OpenTelemetry logs: " + err.Error(),
			},
		}, nil
	}

	logsJSON, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "Failed to serialize logs data: " + err.Error(),
			},
		}, nil
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:  request.Params.URI,
			Text: string(logsJSON),
		},
	}, nil
}

// handleHealthResource provides health status as a resource.
func (s *Service) handleHealthResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	c, err := client.FromContext(ctx)
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "OpenTelemetry client not available: " + err.Error(),
			},
		}, nil
	}
	health, err := c.GetHealth(ctx)
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "Failed to fetch OpenTelemetry health status: " + err.Error(),
			},
		}, nil
	}

	healthJSON, err := json.MarshalIndent(health, "", "  ")
	if err != nil {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:  request.Params.URI,
				Text: "Failed to serialize health data: " + err.Error(),
			},
		}, nil
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:  request.Params.URI,
			Text: string(healthJSON),
		},
	}, nil
}

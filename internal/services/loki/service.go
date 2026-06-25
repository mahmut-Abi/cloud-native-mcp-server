package loki

import (
	"context"
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/cache"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/framework"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/loki/client"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/loki/handlers"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/loki/tools"
	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
)

// Service implements the Loki service for MCP server integration.
// The backend client is not stored — it is created per-request from HTTP headers.
type Service struct {
	enabled       bool
	toolsCache    *cache.ToolsCache
	initFramework *framework.CommonServiceInit
}

// NewService creates a new Loki service instance.
func NewService() *Service {
	checker := framework.NewServiceEnabled(
		func(cfg *config.AppConfig) bool { return cfg.Loki.Enabled },
		func(cfg *config.AppConfig) string { return cfg.Loki.Address },
	)

	initConfig := &framework.InitConfig{
		Required:     false,
		URLValidator: framework.SimpleURLValidator,
		ClientBuilder: func(cfg *config.AppConfig) (interface{}, error) {
			opts := &client.ClientOptions{
				Address:       cfg.Loki.Address,
				Username:      cfg.Loki.Username,
				Password:      cfg.Loki.Password,
				BearerToken:   cfg.Loki.BearerToken,
				Timeout:       30 * time.Second,
				TLSSkipVerify: cfg.Loki.TLSSkipVerify,
				TLSCertFile:   cfg.Loki.TLSCertFile,
				TLSKeyFile:    cfg.Loki.TLSKeyFile,
				TLSCAFile:     cfg.Loki.TLSCAFile,
			}
			if cfg.Loki.TimeoutSec > 0 {
				opts.Timeout = time.Duration(cfg.Loki.TimeoutSec) * time.Second
			}
			return client.NewClient(opts)
		},
	}

	return &Service{
		enabled:       false,
		toolsCache:    cache.NewToolsCache(),
		initFramework: framework.NewCommonServiceInit("Loki", initConfig, checker),
	}
}

// Name returns the service identifier.
func (s *Service) Name() string {
	return "loki"
}

// Initialize configures the Loki service.
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

// GetTools returns all available Loki tools.
func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled {
		return nil
	}

	return s.toolsCache.Get(func() []mcp.Tool {
		return []mcp.Tool{
			tools.QueryLogsSummaryTool(),
			tools.QueryTool(),
			tools.QueryRangeTool(),
			tools.GetLabelNamesTool(),
			tools.GetLabelValuesTool(),
			tools.GetSeriesTool(),
			tools.TestConnectionTool(),
		}
	})
}

// GetHandlers returns all Loki tool handlers.
func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
	if !s.enabled {
		return nil
	}

	handlersMap := map[string]server.ToolHandlerFunc{
		"loki_query_logs_summary": handlers.QueryLogsSummaryHandler(),
		"loki_query":              handlers.QueryHandler(),
		"loki_query_range":        handlers.QueryRangeHandler(),
		"loki_get_label_names":    handlers.GetLabelNamesHandler(),
		"loki_get_label_values":   handlers.GetLabelValuesHandler(),
		"loki_get_series":         handlers.GetSeriesHandler(),
		"loki_test_connection":    handlers.TestConnectionHandler(),
	}

	for name, handler := range handlersMap {
		handlersMap[name] = s.wrapToolErrors(name, handler)
	}

	return handlersMap
}

// IsEnabled returns whether the Loki service is enabled and ready.
func (s *Service) IsEnabled() bool {
	return s.enabled
}

// GetClient returns the underlying Loki client.
func (s *Service) GetClient() *client.Client {
	return nil
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

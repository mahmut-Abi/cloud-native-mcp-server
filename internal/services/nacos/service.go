package nacos

import (
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/cache"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/framework"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/nacos/client"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/nacos/handlers"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/nacos/tools"
)

// Service implements the Nacos MCP service.
// The backend client is not stored — it is created per-request from HTTP headers.
type Service struct {
	enabled            bool
	defaultNamespaceID string
	defaultGroup       string
	toolsCache         *cache.ToolsCache
	initFramework      *framework.CommonServiceInit
}

// NewService creates a new Nacos service instance.
func NewService() *Service {
	checker := framework.NewServiceEnabled(
		func(cfg *config.AppConfig) bool { return cfg.Nacos.Enabled },
		func(cfg *config.AppConfig) string { return cfg.Nacos.URL },
	)

	initConfig := &framework.InitConfig{
		Required:     false,
		URLValidator: framework.SimpleURLValidator,
		ClientBuilder: func(cfg *config.AppConfig) (interface{}, error) {
			return client.NewClient(&client.ClientOptions{
				URL:         cfg.Nacos.URL,
				Username:    cfg.Nacos.Username,
				Password:    cfg.Nacos.Password,
				AccessToken: cfg.Nacos.AccessToken,
				Timeout:     time.Duration(cfg.Nacos.TimeoutSec) * time.Second,
			})
		},
	}

	return &Service{
		enabled:       false,
		toolsCache:    cache.NewToolsCache(),
		initFramework: framework.NewCommonServiceInit("Nacos", initConfig, checker),
	}
}

// Name returns the service identifier.
func (s *Service) Name() string {
	return "nacos"
}

// Initialize configures the Nacos service.
// The backend client is created per-request from HTTP headers (see client/config.go).
func (s *Service) Initialize(cfg interface{}) error {
	appConfig, _ := cfg.(*config.AppConfig)
	if appConfig != nil {
		s.defaultNamespaceID = appConfig.Nacos.NamespaceID
		s.defaultGroup = appConfig.Nacos.Group
	}

	return s.initFramework.Initialize(cfg,
		func(enabled bool) { s.enabled = enabled },
		func(_ interface{}) {
			// Backend client is created per-request from HTTP headers.
			// The backend auth handler was registered in client/config.go init().
		},
	)
}

// GetTools returns all Nacos tools.
func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled {
		return nil
	}

	return s.toolsCache.Get(func() []mcp.Tool {
		return []mcp.Tool{
			tools.TestConnectionTool(),
			tools.ListNamespacesTool(),
			tools.ListConfigsSummaryTool(),
			tools.GetConfigTool(),
			tools.ListServicesSummaryTool(),
			tools.GetServiceTool(),
			tools.ListInstancesTool(),
			tools.ListClusterNodesTool(),
			tools.GetSystemMetricsTool(),
		}
	})
}

// GetHandlers returns all Nacos handlers.
func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
	if !s.enabled {
		return nil
	}

	return map[string]server.ToolHandlerFunc{
		"nacos_test_connection":       handlers.HandleTestConnection(s),
		"nacos_list_namespaces":       handlers.HandleListNamespaces(s),
		"nacos_list_configs_summary":  handlers.HandleListConfigsSummary(s),
		"nacos_get_config":            handlers.HandleGetConfig(s),
		"nacos_list_services_summary": handlers.HandleListServicesSummary(s),
		"nacos_get_service":           handlers.HandleGetService(s),
		"nacos_list_instances":        handlers.HandleListInstances(s),
		"nacos_list_cluster_nodes":    handlers.HandleListClusterNodes(s),
		"nacos_get_system_metrics":    handlers.HandleGetSystemMetrics(s),
	}
}

// IsEnabled returns whether the service is enabled.
func (s *Service) IsEnabled() bool {
	return s.enabled
}

// GetDefaultNamespaceID returns the configured default namespace ID.
func (s *Service) GetDefaultNamespaceID() string {
	return s.defaultNamespaceID
}

// GetDefaultGroup returns the configured default group.
func (s *Service) GetDefaultGroup() string {
	return s.defaultGroup
}

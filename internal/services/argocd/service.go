package argocd

import (
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/argocd/client"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/argocd/handlers"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/argocd/tools"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/cache"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/framework"
)

// Service implements the Argo CD MCP service.
type Service struct {
	client        *client.Client
	enabled       bool
	toolsCache    *cache.ToolsCache
	initFramework *framework.CommonServiceInit
}

// NewService creates a new Argo CD service instance.
func NewService() *Service {
	checker := framework.NewServiceEnabled(
		func(cfg *config.AppConfig) bool { return cfg.ArgoCD.Enabled },
		func(cfg *config.AppConfig) string { return cfg.ArgoCD.URL },
	)

	initConfig := &framework.InitConfig{
		Required:     false,
		URLValidator: framework.SimpleURLValidator,
		ClientBuilder: func(cfg *config.AppConfig) (interface{}, error) {
			return client.NewClient(&client.ClientOptions{
				URL:       cfg.ArgoCD.URL,
				Username:  cfg.ArgoCD.Username,
				Password:  cfg.ArgoCD.Password,
				AuthToken: cfg.ArgoCD.AuthToken,
				Timeout:   time.Duration(cfg.ArgoCD.TimeoutSec) * time.Second,
			})
		},
	}

	return &Service{
		enabled:       false,
		toolsCache:    cache.NewToolsCache(),
		initFramework: framework.NewCommonServiceInit("ArgoCD", initConfig, checker),
	}
}

// Name returns the service identifier.
func (s *Service) Name() string {
	return "argocd"
}

// Initialize configures the Argo CD service.
func (s *Service) Initialize(cfg interface{}) error {
	return s.initFramework.Initialize(cfg,
		func(enabled bool) { s.enabled = enabled },
		func(clientIface interface{}) {
			if argocdClient, ok := clientIface.(*client.Client); ok {
				s.client = argocdClient
			}
		},
	)
}

// GetTools returns all Argo CD tools.
func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled || s.client == nil {
		return nil
	}

	return s.toolsCache.Get(func() []mcp.Tool {
		return []mcp.Tool{
			tools.TestConnectionTool(),
			tools.ListApplicationsSummaryTool(),
			tools.GetApplicationTool(),
			tools.GetApplicationManifestsTool(),
			tools.ListProjectsTool(),
			tools.GetProjectTool(),
			tools.ListClustersTool(),
		}
	})
}

// GetHandlers returns all Argo CD handlers.
func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
	if !s.enabled || s.client == nil {
		return nil
	}

	return map[string]server.ToolHandlerFunc{
		"argocd_test_connection":           handlers.HandleTestConnection(s),
		"argocd_list_applications_summary": handlers.HandleListApplicationsSummary(s),
		"argocd_get_application":           handlers.HandleGetApplication(s),
		"argocd_get_application_manifests": handlers.HandleGetApplicationManifests(s),
		"argocd_list_projects":             handlers.HandleListProjects(s),
		"argocd_get_project":               handlers.HandleGetProject(s),
		"argocd_list_clusters":             handlers.HandleListClusters(s),
	}
}

// IsEnabled returns whether the service is enabled and ready.
func (s *Service) IsEnabled() bool {
	return s.enabled && s.client != nil
}

// GetClient exposes the underlying Argo CD client.
func (s *Service) GetClient() *client.Client {
	return s.client
}

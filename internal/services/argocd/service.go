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
// The backend client is not stored — it is created per-request from HTTP headers.
type Service struct {
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

// GetTools returns all Argo CD tools.
func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled {
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
	if !s.enabled {
		return nil
	}

	return map[string]server.ToolHandlerFunc{
		"argocd_test_connection":           handlers.HandleTestConnection(),
		"argocd_list_applications_summary": handlers.HandleListApplicationsSummary(),
		"argocd_get_application":           handlers.HandleGetApplication(),
		"argocd_get_application_manifests": handlers.HandleGetApplicationManifests(),
		"argocd_list_projects":             handlers.HandleListProjects(),
		"argocd_get_project":               handlers.HandleGetProject(),
		"argocd_list_clusters":             handlers.HandleListClusters(),
	}
}

// IsEnabled returns whether the service is enabled.
func (s *Service) IsEnabled() bool {
	return s.enabled
}

package sentry

import (
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/cache"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/framework"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/sentry/client"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/sentry/handlers"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/sentry/tools"
)

// Service implements the Sentry MCP service.
// The backend client is not stored — it is created per-request from HTTP headers.
type Service struct {
	enabled             bool
	defaultOrganization string
	defaultProject      string
	toolsCache          *cache.ToolsCache
	initFramework       *framework.CommonServiceInit
}

// NewService creates a new Sentry service instance.
func NewService() *Service {
	checker := framework.NewServiceEnabled(
		func(cfg *config.AppConfig) bool { return cfg.Sentry.Enabled },
		func(cfg *config.AppConfig) string { return cfg.Sentry.URL },
	)

	initConfig := &framework.InitConfig{
		Required:     false,
		URLValidator: framework.SimpleURLValidator,
		ClientBuilder: func(cfg *config.AppConfig) (interface{}, error) {
			return client.NewClient(&client.ClientOptions{
				URL:       cfg.Sentry.URL,
				AuthToken: cfg.Sentry.AuthToken,
				Timeout:   time.Duration(cfg.Sentry.TimeoutSec) * time.Second,
			})
		},
	}

	return &Service{
		enabled:       false,
		toolsCache:    cache.NewToolsCache(),
		initFramework: framework.NewCommonServiceInit("Sentry", initConfig, checker),
	}
}

// Name returns the service identifier.
func (s *Service) Name() string {
	return "sentry"
}

// Initialize configures the Sentry service.
// The backend client is created per-request from HTTP headers (see client/config.go).
func (s *Service) Initialize(cfg interface{}) error {
	appConfig, _ := cfg.(*config.AppConfig)
	if appConfig != nil {
		s.defaultOrganization = appConfig.Sentry.Organization
		s.defaultProject = appConfig.Sentry.Project
	}

	return s.initFramework.Initialize(cfg,
		func(enabled bool) { s.enabled = enabled },
		func(_ interface{}) {
			// Backend client is created per-request from HTTP headers.
			// The backend auth handler was registered in client/config.go init().
		},
	)
}

// GetTools returns all Sentry tools.
func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled {
		return nil
	}

	return s.toolsCache.Get(func() []mcp.Tool {
		return []mcp.Tool{
			tools.TestConnectionTool(),
			tools.ListOrganizationsTool(),
			tools.ListProjectsTool(),
			tools.GetProjectTool(),
			tools.ListIssuesSummaryTool(),
			tools.ListIssuesTool(),
			tools.GetIssueTool(),
			tools.ListIssueEventsTool(),
			tools.GetIssueEventTool(),
		}
	})
}

// GetHandlers returns all Sentry handlers.
func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
	if !s.enabled {
		return nil
	}

	return map[string]server.ToolHandlerFunc{
		"sentry_test_connection":     handlers.HandleTestConnection(s),
		"sentry_list_organizations":  handlers.HandleListOrganizations(s),
		"sentry_list_projects":       handlers.HandleListProjects(s),
		"sentry_get_project":         handlers.HandleGetProject(s),
		"sentry_list_issues_summary": handlers.HandleListIssuesSummary(s),
		"sentry_list_issues":         handlers.HandleListIssues(s),
		"sentry_get_issue":           handlers.HandleGetIssue(s),
		"sentry_list_issue_events":   handlers.HandleListIssueEvents(s),
		"sentry_get_issue_event":     handlers.HandleGetIssueEvent(s),
	}
}

// IsEnabled returns whether the service is enabled.
func (s *Service) IsEnabled() bool {
	return s.enabled
}

// GetDefaultOrganization returns the configured default organization slug.
func (s *Service) GetDefaultOrganization() string {
	return s.defaultOrganization
}

// GetDefaultProject returns the configured default project slug.
func (s *Service) GetDefaultProject() string {
	return s.defaultProject
}

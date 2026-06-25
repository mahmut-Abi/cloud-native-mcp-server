package langfuse

import (
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/cache"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/framework"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/langfuse/client"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/langfuse/handlers"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/langfuse/tools"
)

// Service implements the Langfuse MCP service.
// The backend client is not stored — it is created per-request from HTTP headers.
type Service struct {
	enabled       bool
	toolsCache    *cache.ToolsCache
	initFramework *framework.CommonServiceInit
}

// NewService creates a new Langfuse service instance.
func NewService() *Service {
	checker := framework.NewServiceEnabled(
		func(cfg *config.AppConfig) bool { return cfg.Langfuse.Enabled },
		func(cfg *config.AppConfig) string { return cfg.Langfuse.URL },
	)

	initConfig := &framework.InitConfig{
		Required:     false,
		URLValidator: framework.SimpleURLValidator,
		ClientBuilder: func(cfg *config.AppConfig) (interface{}, error) {
			return client.NewClient(&client.ClientOptions{
				URL:      cfg.Langfuse.URL,
				Username: cfg.Langfuse.Username,
				Password: cfg.Langfuse.Password,
				Timeout:  time.Duration(cfg.Langfuse.TimeoutSec) * time.Second,
			})
		},
	}

	return &Service{
		enabled:       false,
		toolsCache:    cache.NewToolsCache(),
		initFramework: framework.NewCommonServiceInit("Langfuse", initConfig, checker),
	}
}

// Name returns the service identifier.
func (s *Service) Name() string {
	return "langfuse"
}

// Initialize configures the Langfuse service.
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

// GetTools returns all Langfuse tools.
func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled {
		return nil
	}

	return s.toolsCache.Get(func() []mcp.Tool {
		return []mcp.Tool{
			tools.CheckHealthTool(),
			tools.ListTracesSummaryTool(),
			tools.ListTracesTool(),
			tools.GetTraceTool(),
			tools.ListAnnotationQueuesTool(),
			tools.GetAnnotationQueueTool(),
			tools.ListAnnotationQueueItemsTool(),
			tools.ListDatasetsTool(),
			tools.GetDatasetTool(),
			tools.ListDatasetRunsTool(),
			tools.GetDatasetRunTool(),
			tools.ListLLMConnectionsTool(),
			tools.ListModelsTool(),
			tools.GetModelTool(),
			tools.ListSessionsTool(),
			tools.GetSessionTool(),
			tools.ListObservationsTool(),
			tools.GetObservationTool(),
			tools.ListPromptsTool(),
			tools.GetPromptTool(),
			tools.ListScoreConfigsTool(),
			tools.GetScoreConfigTool(),
			tools.ListScoresTool(),
			tools.GetScoreTool(),
			tools.GetMetricsTool(),
			tools.GetProjectTool(),
			tools.ListOrganizationProjectsTool(),
			tools.CreateProjectTool(),
			tools.UpdateProjectTool(),
			tools.DeleteProjectTool(),
			tools.ListProjectMembershipsTool(),
			tools.UpsertProjectMembershipTool(),
			tools.DeleteProjectMembershipTool(),
			tools.ListOrganizationAPIKeysTool(),
			tools.ListProjectAPIKeysTool(),
			tools.CreateProjectAPIKeyTool(),
			tools.DeleteProjectAPIKeyTool(),
		}
	})
}

// GetHandlers returns all Langfuse handlers.
func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
	if !s.enabled {
		return nil
	}

	return map[string]server.ToolHandlerFunc{
		"langfuse_check_health":                handlers.HandleCheckHealth(),
		"langfuse_list_traces_summary":         handlers.HandleListTracesSummary(),
		"langfuse_list_traces":                 handlers.HandleListTraces(),
		"langfuse_get_trace":                   handlers.HandleGetTrace(),
		"langfuse_list_annotation_queues":      handlers.HandleListAnnotationQueues(),
		"langfuse_get_annotation_queue":        handlers.HandleGetAnnotationQueue(),
		"langfuse_list_annotation_queue_items": handlers.HandleListAnnotationQueueItems(),
		"langfuse_list_datasets":               handlers.HandleListDatasets(),
		"langfuse_get_dataset":                 handlers.HandleGetDataset(),
		"langfuse_list_dataset_runs":           handlers.HandleListDatasetRuns(),
		"langfuse_get_dataset_run":             handlers.HandleGetDatasetRun(),
		"langfuse_list_llm_connections":        handlers.HandleListLLMConnections(),
		"langfuse_list_models":                 handlers.HandleListModels(),
		"langfuse_get_model":                   handlers.HandleGetModel(),
		"langfuse_list_sessions":               handlers.HandleListSessions(),
		"langfuse_get_session":                 handlers.HandleGetSession(),
		"langfuse_list_observations":           handlers.HandleListObservations(),
		"langfuse_get_observation":             handlers.HandleGetObservation(),
		"langfuse_list_prompts":                handlers.HandleListPrompts(),
		"langfuse_get_prompt":                  handlers.HandleGetPrompt(),
		"langfuse_list_score_configs":          handlers.HandleListScoreConfigs(),
		"langfuse_get_score_config":            handlers.HandleGetScoreConfig(),
		"langfuse_list_scores":                 handlers.HandleListScores(),
		"langfuse_get_score":                   handlers.HandleGetScore(),
		"langfuse_get_metrics":                 handlers.HandleGetMetrics(),
		"langfuse_get_project":                 handlers.HandleGetProject(),
		"langfuse_list_organization_projects":  handlers.HandleListOrganizationProjects(),
		"langfuse_create_project":              handlers.HandleCreateProject(),
		"langfuse_update_project":              handlers.HandleUpdateProject(),
		"langfuse_delete_project":              handlers.HandleDeleteProject(),
		"langfuse_list_project_memberships":    handlers.HandleListProjectMemberships(),
		"langfuse_upsert_project_membership":   handlers.HandleUpsertProjectMembership(),
		"langfuse_delete_project_membership":   handlers.HandleDeleteProjectMembership(),
		"langfuse_list_organization_api_keys":  handlers.HandleListOrganizationAPIKeys(),
		"langfuse_list_project_api_keys":       handlers.HandleListProjectAPIKeys(),
		"langfuse_create_project_api_key":      handlers.HandleCreateProjectAPIKey(),
		"langfuse_delete_project_api_key":      handlers.HandleDeleteProjectAPIKey(),
	}
}

// IsEnabled returns whether the service is enabled.
func (s *Service) IsEnabled() bool {
	return s.enabled
}

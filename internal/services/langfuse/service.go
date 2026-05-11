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
type Service struct {
	client        *client.Client
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
				URL:       cfg.Langfuse.URL,
				PublicKey: cfg.Langfuse.PublicKey,
				SecretKey: cfg.Langfuse.SecretKey,
				Timeout:   time.Duration(cfg.Langfuse.TimeoutSec) * time.Second,
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
func (s *Service) Initialize(cfg interface{}) error {
	return s.initFramework.Initialize(cfg,
		func(enabled bool) { s.enabled = enabled },
		func(clientIface interface{}) {
			if langfuseClient, ok := clientIface.(*client.Client); ok {
				s.client = langfuseClient
			}
		},
	)
}

// GetTools returns all Langfuse tools.
func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled || s.client == nil {
		return nil
	}

	return s.toolsCache.Get(func() []mcp.Tool {
		return []mcp.Tool{
			tools.CheckHealthTool(),
			tools.ListTracesSummaryTool(),
			tools.ListTracesTool(),
			tools.GetTraceTool(),
			tools.ListSessionsTool(),
			tools.GetSessionTool(),
			tools.ListObservationsTool(),
			tools.GetObservationTool(),
			tools.ListPromptsTool(),
			tools.GetPromptTool(),
			tools.ListScoresTool(),
			tools.GetScoreTool(),
			tools.GetMetricsTool(),
		}
	})
}

// GetHandlers returns all Langfuse handlers.
func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
	if !s.enabled || s.client == nil {
		return nil
	}

	return map[string]server.ToolHandlerFunc{
		"langfuse_check_health":        handlers.HandleCheckHealth(s),
		"langfuse_list_traces_summary": handlers.HandleListTracesSummary(s),
		"langfuse_list_traces":         handlers.HandleListTraces(s),
		"langfuse_get_trace":           handlers.HandleGetTrace(s),
		"langfuse_list_sessions":       handlers.HandleListSessions(s),
		"langfuse_get_session":         handlers.HandleGetSession(s),
		"langfuse_list_observations":   handlers.HandleListObservations(s),
		"langfuse_get_observation":     handlers.HandleGetObservation(s),
		"langfuse_list_prompts":        handlers.HandleListPrompts(s),
		"langfuse_get_prompt":          handlers.HandleGetPrompt(s),
		"langfuse_list_scores":         handlers.HandleListScores(s),
		"langfuse_get_score":           handlers.HandleGetScore(s),
		"langfuse_get_metrics":         handlers.HandleGetMetrics(s),
	}
}

// IsEnabled returns whether the service is enabled and ready.
func (s *Service) IsEnabled() bool {
	return s.enabled && s.client != nil
}

// GetClient exposes the underlying Langfuse client.
func (s *Service) GetClient() *client.Client {
	return s.client
}

package dify

import (
	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/cache"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/dify/handlers"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/dify/tools"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/framework"
)

type Service struct {
	enabled       bool
	toolsCache    *cache.ToolsCache
	initFramework *framework.CommonServiceInit
}

func NewService() *Service {
	checker := framework.NewServiceEnabled(
		func(cfg *config.AppConfig) bool { return true },
		func(cfg *config.AppConfig) string { return "header-based-auth" },
	)

	initConfig := &framework.InitConfig{
		Required:      false,
		URLValidator:  framework.SimpleURLValidator,
		ClientBuilder: nil,
	}

	return &Service{
		enabled:       false,
		toolsCache:    cache.NewToolsCache(),
		initFramework: framework.NewCommonServiceInit("Dify", initConfig, checker),
	}
}

func (s *Service) Name() string {
	return "dify"
}

func (s *Service) Initialize(cfg interface{}) error {
	return s.initFramework.Initialize(cfg,
		func(enabled bool) { s.enabled = enabled },
		func(_ interface{}) {},
	)
}

func (s *Service) GetTools() []mcp.Tool {
	if !s.enabled {
		return nil
	}

	return s.toolsCache.Get(func() []mcp.Tool {
		return []mcp.Tool{
			tools.DifyAppInfoTool(),
			tools.DifyAppMetaTool(),
			tools.DifyAppParametersTool(),
			tools.DifyAppSiteTool(),
			tools.DifyAPIRequestTool(),
			tools.DifyChatMessageTool(),
			tools.DifyCompletionMessageTool(),
			tools.DifyListConversationsTool(),
			tools.DifyListMessagesTool(),
			tools.DifyRetrieveDatasetTool(),
			tools.DifyRunWorkflowTool(),
			tools.ConsoleAPIRequestTool(),
			tools.ListAppsTool(),
			tools.GetAppTool(),
			tools.CreateAppTool(),
			tools.SetAppAPIStatusTool(),
			tools.ListAppAPIKeysTool(),
			tools.CreateAppAPIKeyTool(),
			tools.GetAppTraceStatusTool(),
			tools.SetAppTraceStatusTool(),
			tools.GetAppTraceConfigTool(),
			tools.CreateAppTraceConfigTool(),
			tools.UpdateAppTraceConfigTool(),
			tools.DeleteAppTraceConfigTool(),
			tools.GetDraftWorkflowTool(),
			tools.SyncDraftWorkflowTool(),
			tools.GetDraftWorkflowEnvVarsTool(),
			tools.UpdateDraftWorkflowEnvVarsTool(),
			tools.GetDraftWorkflowConvVarsTool(),
			tools.UpdateDraftWorkflowConvVarsTool(),
			tools.GetPublishedWorkflowTool(),
			tools.PublishWorkflowTool(),
			tools.ListPublishedWorkflowsTool(),
			tools.RestoreWorkflowToDraftTool(),
			tools.RunDraftWorkflowTool(),
			tools.RunAdvancedChatDraftWorkflowTool(),
			tools.ListDatasetsTool(),
			tools.GetDatasetTool(),
			tools.CreateDatasetTool(),
			tools.DeleteDatasetTool(),
			tools.UpdateDatasetTool(),
			tools.ListDatasetDocumentsTool(),
			tools.GetDatasetDocumentTool(),
			tools.SetDatasetAPIStatusTool(),
			tools.ListDatasetAPIKeysTool(),
			tools.CreateDatasetAPIKeyTool(),
		}
	})
}

func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
	if !s.enabled {
		return nil
	}

	return map[string]server.ToolHandlerFunc{
		"dify_dify_app_info":               handlers.HandleDifyAppInfo(s),
		"dify_dify_app_meta":               handlers.HandleDifyAppMeta(s),
		"dify_dify_app_parameters":         handlers.HandleDifyAppParameters(s),
		"dify_dify_app_site":               handlers.HandleDifyAppSite(s),
		"dify_dify_api_request":            handlers.HandleDifyApiRequest(s),
		"dify_dify_chat_message":           handlers.HandleDifyChatMessage(s),
		"dify_dify_completion_message":     handlers.HandleDifyCompletionMessage(s),
		"dify_dify_list_conversations":     handlers.HandleDifyListConversations(s),
		"dify_dify_list_messages":          handlers.HandleDifyListMessages(s),
		"dify_dify_retrieve_dataset":       handlers.HandleDifyRetrieveDataset(s),
		"dify_dify_run_workflow":           handlers.HandleDifyRunWorkflow(s),
		"dify_console_api_request":         handlers.HandleDifyConsoleApiRequest(s),
		"dify_list_apps":                   handlers.HandleDifyListApps(s),
		"dify_get_app":                     handlers.HandleDifyGetApp(s),
		"dify_create_app":                  handlers.HandleDifyCreateApp(s),
		"dify_set_app_api_status":          handlers.HandleDifySetAppApiStatus(s),
		"dify_list_app_api_keys":           handlers.HandleDifyListAppApiKeys(s),
		"dify_create_app_api_key":          handlers.HandleDifyCreateAppApiKey(s),
		"dify_get_app_trace_status":        handlers.HandleDifyGetAppTraceStatus(s),
		"dify_set_app_trace_status":        handlers.HandleDifySetAppTraceStatus(s),
		"dify_get_app_trace_config":        handlers.HandleDifyGetAppTraceConfig(s),
		"dify_create_app_trace_config":     handlers.HandleDifyCreateAppTraceConfig(s),
		"dify_update_app_trace_config":     handlers.HandleDifyUpdateAppTraceConfig(s),
		"dify_delete_app_trace_config":     handlers.HandleDifyDeleteAppTraceConfig(s),
		"dify_get_draft_workflow":          handlers.HandleDifyGetDraftWorkflow(s),
		"dify_sync_draft_workflow":         handlers.HandleDifySyncDraftWorkflow(s),
		"dify_get_draft_workflow_environment_variables": handlers.HandleDifyGetDraftWorkflowEnvVars(s),
		"dify_update_draft_workflow_environment_variables": handlers.HandleDifyUpdateDraftWorkflowEnvVars(s),
		"dify_get_draft_workflow_conversation_variables":   handlers.HandleDifyGetDraftWorkflowConvVars(s),
		"dify_update_draft_workflow_conversation_variables": handlers.HandleDifyUpdateDraftWorkflowConvVars(s),
		"dify_get_published_workflow":          handlers.HandleDifyGetPublishedWorkflow(s),
		"dify_publish_workflow":                handlers.HandleDifyPublishWorkflow(s),
		"dify_list_published_workflows":        handlers.HandleDifyListPublishedWorkflows(s),
		"dify_restore_workflow_to_draft":       handlers.HandleDifyRestoreWorkflowToDraft(s),
		"dify_run_draft_workflow":              handlers.HandleDifyRunDraftWorkflow(s),
		"dify_run_advanced_chat_draft_workflow": handlers.HandleDifyRunAdvancedChatDraftWorkflow(s),
		"dify_list_datasets":                    handlers.HandleDifyListDatasets(s),
		"dify_get_dataset":                      handlers.HandleDifyGetDataset(s),
		"dify_create_dataset":                   handlers.HandleDifyCreateDataset(s),
		"dify_delete_dataset":                   handlers.HandleDifyDeleteDataset(s),
		"dify_update_dataset":                   handlers.HandleDifyUpdateDataset(s),
		"dify_list_dataset_documents":           handlers.HandleDifyListDatasetDocuments(s),
		"dify_get_dataset_document":             handlers.HandleDifyGetDatasetDocument(s),
		"dify_set_dataset_api_status":           handlers.HandleDifySetDatasetApiStatus(s),
		"dify_list_dataset_api_keys":            handlers.HandleDifyListDatasetApiKeys(s),
		"dify_create_dataset_api_key":           handlers.HandleDifyCreateDatasetApiKey(s),
	}
}

func (s *Service) IsEnabled() bool {
	return s.enabled
}

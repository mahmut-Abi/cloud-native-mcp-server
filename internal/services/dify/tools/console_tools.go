package tools

import "github.com/mark3labs/mcp-go/mcp"

// ConsoleAPIRequestTool returns a tool for calling any relative Dify Console API path.
func ConsoleAPIRequestTool() mcp.Tool {
	return mcp.NewTool("dify_console_api_request",
		mcp.WithDescription("Call any relative Dify Console API path using DIFY_CONSOLE_EMAIL/DIFY_CONSOLE_PASSWORD."),
		mcp.WithString("method",
			mcp.Description("HTTP method."),
			mcp.Enum("GET", "POST", "PUT", "PATCH", "DELETE"),
			mcp.DefaultString("GET")),
		mcp.WithString("path", mcp.Required(),
			mcp.Description("Relative /console/api path, for example /apps or /datasets.")),
		mcp.WithObject("query",
			mcp.Description("Query string parameters.")),
		mcp.WithObject("body",
			mcp.Description("JSON body for non-GET requests.")),
	)
}

// ListAppsTool returns the app listing tool.
func ListAppsTool() mcp.Tool {
	return mcp.NewTool("dify_list_apps",
		mcp.WithDescription("List Dify apps globally from the Console API."),
		mcp.WithNumber("page",
			mcp.Description("Page number (starts at 1)."),
			mcp.DefaultNumber(1),
			mcp.Min(1)),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of records to return."),
			mcp.DefaultNumber(20),
			mcp.Min(1),
			mcp.Max(100)),
		mcp.WithString("mode",
			mcp.Description("App mode filter."),
			mcp.Enum("completion", "chat", "advanced-chat", "workflow", "agent-chat", "channel", "all"),
			mcp.DefaultString("all")),
		mcp.WithString("name",
			mcp.Description("App name filter.")),
		mcp.WithBoolean("is_created_by_me",
			mcp.Description("Filter to only apps created by the current user.")),
	)
}

// GetAppTool returns the app detail tool.
func GetAppTool() mcp.Tool {
	return mcp.NewTool("dify_get_app",
		mcp.WithDescription("Get details for any Dify app by app ID."),
		mcp.WithString("app_id", mcp.Required(),
			mcp.Description("Dify app UUID.")),
	)
}

// CreateAppTool returns the app creation tool.
func CreateAppTool() mcp.Tool {
	return mcp.NewTool("dify_create_app",
		mcp.WithDescription("Create a Dify app globally from the Console API."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("App name."),
			mcp.MinLength(1)),
		mcp.WithString("mode", mcp.Required(),
			mcp.Description("App mode."),
			mcp.Enum("chat", "agent-chat", "advanced-chat", "workflow", "completion")),
		mcp.WithString("description",
			mcp.Description("App description.")),
		mcp.WithString("icon_type",
			mcp.Description("Icon type.")),
		mcp.WithString("icon",
			mcp.Description("Icon identifier.")),
		mcp.WithString("icon_background",
			mcp.Description("Icon background color.")),
	)
}

// SetAppAPIStatusTool returns the app API status toggle tool.
func SetAppAPIStatusTool() mcp.Tool {
	return mcp.NewTool("dify_set_app_api_status",
		mcp.WithDescription("Enable or disable Service API access for any Dify app."),
		mcp.WithString("app_id", mcp.Required(),
			mcp.Description("Dify app UUID.")),
		mcp.WithBoolean("enable_api", mcp.Required(),
			mcp.Description("Whether to enable or disable the Service API.")),
	)
}

// ListAppAPIKeysTool returns the app API key listing tool.
func ListAppAPIKeysTool() mcp.Tool {
	return mcp.NewTool("dify_list_app_api_keys",
		mcp.WithDescription("List Service API keys for any Dify app."),
		mcp.WithString("app_id", mcp.Required(),
			mcp.Description("Dify app UUID.")),
	)
}

// CreateAppAPIKeyTool returns the app API key creation tool.
func CreateAppAPIKeyTool() mcp.Tool {
	return mcp.NewTool("dify_create_app_api_key",
		mcp.WithDescription("Create a Service API key for any Dify app."),
		mcp.WithString("app_id", mcp.Required(),
			mcp.Description("Dify app UUID.")),
	)
}

// GetAppTraceStatusTool returns the app trace status tool.
func GetAppTraceStatusTool() mcp.Tool {
	return mcp.NewTool("dify_get_app_trace_status",
		mcp.WithDescription("Get whether app tracing is enabled and which tracing provider is active."),
		mcp.WithString("app_id", mcp.Required(),
			mcp.Description("Dify app UUID.")),
	)
}

// SetAppTraceStatusTool returns the app trace status toggle tool.
func SetAppTraceStatusTool() mcp.Tool {
	return mcp.NewTool("dify_set_app_trace_status",
		mcp.WithDescription("Enable or disable app tracing. When enabling, tracing_provider is required."),
		mcp.WithString("app_id", mcp.Required(),
			mcp.Description("Dify app UUID.")),
		mcp.WithBoolean("enabled", mcp.Required(),
			mcp.Description("Whether tracing is enabled.")),
		mcp.WithString("tracing_provider",
			mcp.Description("Tracing provider name."),
			mcp.Enum("aliyun", "arize", "databricks", "langfuse", "langsmith", "mlflow", "opik", "phoenix", "tencent", "weave")),
	)
}

// GetAppTraceConfigTool returns the app trace config tool.
func GetAppTraceConfigTool() mcp.Tool {
	return mcp.NewTool("dify_get_app_trace_config",
		mcp.WithDescription("Get a tracing provider config for an app."),
		mcp.WithString("app_id", mcp.Required(),
			mcp.Description("Dify app UUID.")),
		mcp.WithString("tracing_provider", mcp.Required(),
			mcp.Description("Tracing provider name."),
			mcp.Enum("aliyun", "arize", "databricks", "langfuse", "langsmith", "mlflow", "opik", "phoenix", "tencent", "weave")),
	)
}

// CreateAppTraceConfigTool returns the app trace config creation tool.
func CreateAppTraceConfigTool() mcp.Tool {
	return mcp.NewTool("dify_create_app_trace_config",
		mcp.WithDescription("Create a tracing provider config for an app. Body sent to Dify is { tracing_provider, tracing_config }."),
		mcp.WithString("app_id", mcp.Required(),
			mcp.Description("Dify app UUID.")),
		mcp.WithString("tracing_provider", mcp.Required(),
			mcp.Description("Tracing provider name."),
			mcp.Enum("aliyun", "arize", "databricks", "langfuse", "langsmith", "mlflow", "opik", "phoenix", "tencent", "weave")),
		mcp.WithObject("tracing_config", mcp.Required(),
			mcp.Description("Provider-specific config, for example { public_key, secret_key, host } for langfuse.")),
	)
}

// UpdateAppTraceConfigTool returns the app trace config update tool.
func UpdateAppTraceConfigTool() mcp.Tool {
	return mcp.NewTool("dify_update_app_trace_config",
		mcp.WithDescription("Update an existing tracing provider config for an app. Body sent to Dify is { tracing_provider, tracing_config }."),
		mcp.WithString("app_id", mcp.Required(),
			mcp.Description("Dify app UUID.")),
		mcp.WithString("tracing_provider", mcp.Required(),
			mcp.Description("Tracing provider name."),
			mcp.Enum("aliyun", "arize", "databricks", "langfuse", "langsmith", "mlflow", "opik", "phoenix", "tencent", "weave")),
		mcp.WithObject("tracing_config", mcp.Required(),
			mcp.Description(`Provider-specific config. Use "*" for unchanged secret values when Dify supports secret preservation.`)),
	)
}

// DeleteAppTraceConfigTool returns the app trace config deletion tool.
func DeleteAppTraceConfigTool() mcp.Tool {
	return mcp.NewTool("dify_delete_app_trace_config",
		mcp.WithDescription("Delete a tracing provider config for an app."),
		mcp.WithString("app_id", mcp.Required(),
			mcp.Description("Dify app UUID.")),
		mcp.WithString("tracing_provider", mcp.Required(),
			mcp.Description("Tracing provider name."),
			mcp.Enum("aliyun", "arize", "databricks", "langfuse", "langsmith", "mlflow", "opik", "phoenix", "tencent", "weave")),
	)
}

// GetDraftWorkflowTool returns the draft workflow tool.
func GetDraftWorkflowTool() mcp.Tool {
	return mcp.NewTool("dify_get_draft_workflow",
		mcp.WithDescription("Get the draft workflow graph/features for any workflow or advanced-chat app."),
		mcp.WithString("app_id", mcp.Required(),
			mcp.Description("Dify app UUID.")),
	)
}

// SyncDraftWorkflowTool returns the draft workflow sync tool.
func SyncDraftWorkflowTool() mcp.Tool {
	return mcp.NewTool("dify_sync_draft_workflow",
		mcp.WithDescription("Synchronize the full draft workflow graph/features for any workflow or advanced-chat app. Prefer dify_update_draft_workflow_environment_variables for environment-variable-only changes."),
		mcp.WithString("app_id", mcp.Required(),
			mcp.Description("Dify app UUID.")),
		mcp.WithObject("graph",
			mcp.Description("Workflow graph object. Required for full draft synchronization.")),
		mcp.WithObject("features",
			mcp.Description("Workflow features object. Required for full draft synchronization.")),
		mcp.WithString("hash",
			mcp.Description("Current draft workflow hash for optimistic locking.")),
		mcp.WithObject("conversation_variables",
			mcp.Description("Full conversation variable list. If graph/features are omitted, this tool updates only conversation variables.")),
		mcp.WithObject("environment_variables",
			mcp.Description("Full environment variable list. If graph/features are omitted, this tool updates only environment variables.")),
	)
}

// GetDraftWorkflowEnvVarsTool returns the draft workflow env vars tool.
func GetDraftWorkflowEnvVarsTool() mcp.Tool {
	return mcp.NewTool("dify_get_draft_workflow_environment_variables",
		mcp.WithDescription("Get environment variables for a draft workflow."),
		mcp.WithString("app_id", mcp.Required(),
			mcp.Description("Dify app UUID.")),
	)
}

// UpdateDraftWorkflowEnvVarsTool returns the draft workflow env vars update tool.
func UpdateDraftWorkflowEnvVarsTool() mcp.Tool {
	return mcp.NewTool("dify_update_draft_workflow_environment_variables",
		mcp.WithDescription("Preferred tool to update only draft workflow environment variables without changing the graph or features."),
		mcp.WithString("app_id", mcp.Required(),
			mcp.Description("Dify app UUID.")),
		mcp.WithObject("environment_variables", mcp.Required(),
			mcp.Description("Full environment variable list. Each item needs name, value, and value_type.")),
	)
}

// GetDraftWorkflowConvVarsTool returns the draft workflow conversation vars tool.
func GetDraftWorkflowConvVarsTool() mcp.Tool {
	return mcp.NewTool("dify_get_draft_workflow_conversation_variables",
		mcp.WithDescription("Get conversation variables for a draft workflow."),
		mcp.WithString("app_id", mcp.Required(),
			mcp.Description("Dify app UUID.")),
	)
}

// UpdateDraftWorkflowConvVarsTool returns the draft workflow conversation vars update tool.
func UpdateDraftWorkflowConvVarsTool() mcp.Tool {
	return mcp.NewTool("dify_update_draft_workflow_conversation_variables",
		mcp.WithDescription("Preferred tool to update only draft workflow conversation variables without changing the graph or features."),
		mcp.WithString("app_id", mcp.Required(),
			mcp.Description("Dify app UUID.")),
		mcp.WithObject("conversation_variables", mcp.Required(),
			mcp.Description("Full conversation variable list. Each item needs name, value, and value_type.")),
	)
}

// GetPublishedWorkflowTool returns the published workflow tool.
func GetPublishedWorkflowTool() mcp.Tool {
	return mcp.NewTool("dify_get_published_workflow",
		mcp.WithDescription("Get the currently published workflow for any workflow or advanced-chat app."),
		mcp.WithString("app_id", mcp.Required(),
			mcp.Description("Dify app UUID.")),
	)
}

// PublishWorkflowTool returns the workflow publish tool.
func PublishWorkflowTool() mcp.Tool {
	return mcp.NewTool("dify_publish_workflow",
		mcp.WithDescription("Publish the draft workflow for any workflow or advanced-chat app."),
		mcp.WithString("app_id", mcp.Required(),
			mcp.Description("Dify app UUID.")),
		mcp.WithString("marked_name",
			mcp.Description("Version name (max 20 characters)."),
			mcp.MaxLength(20)),
		mcp.WithString("marked_comment",
			mcp.Description("Version comment (max 100 characters)."),
			mcp.MaxLength(100)),
	)
}

// ListPublishedWorkflowsTool returns the published workflow listing tool.
func ListPublishedWorkflowsTool() mcp.Tool {
	return mcp.NewTool("dify_list_published_workflows",
		mcp.WithDescription("List published workflow versions for any workflow or advanced-chat app."),
		mcp.WithString("app_id", mcp.Required(),
			mcp.Description("Dify app UUID.")),
		mcp.WithNumber("page",
			mcp.Description("Page number (starts at 1)."),
			mcp.DefaultNumber(1),
			mcp.Min(1)),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of records to return."),
			mcp.Min(1),
			mcp.Max(100),
			mcp.DefaultNumber(20)),
		mcp.WithBoolean("named_only",
			mcp.Description("Only return named versions.")),
	)
}

// RestoreWorkflowToDraftTool returns the workflow restore tool.
func RestoreWorkflowToDraftTool() mcp.Tool {
	return mcp.NewTool("dify_restore_workflow_to_draft",
		mcp.WithDescription("Restore a published workflow version into the draft workflow."),
		mcp.WithString("app_id", mcp.Required(),
			mcp.Description("Dify app UUID.")),
		mcp.WithString("workflow_id", mcp.Required(),
			mcp.Description("Published workflow ID to restore.")),
	)
}

// RunDraftWorkflowTool returns the draft workflow run tool.
func RunDraftWorkflowTool() mcp.Tool {
	return mcp.NewTool("dify_run_draft_workflow",
		mcp.WithDescription("Run a workflow app draft from the Console debugger API."),
		mcp.WithString("app_id", mcp.Required(),
			mcp.Description("Dify app UUID.")),
		mcp.WithObject("inputs",
			mcp.Description("Workflow input variables.")),
		mcp.WithObject("files",
			mcp.Description("Dify file descriptors.")),
	)
}

// RunAdvancedChatDraftWorkflowTool returns the advanced chat draft workflow run tool.
func RunAdvancedChatDraftWorkflowTool() mcp.Tool {
	return mcp.NewTool("dify_run_advanced_chat_draft_workflow",
		mcp.WithDescription("Run an advanced-chat app draft workflow from the Console debugger API."),
		mcp.WithString("app_id", mcp.Required(),
			mcp.Description("Dify app UUID.")),
		mcp.WithString("query",
			mcp.Description("User message text."),
			mcp.DefaultString("")),
		mcp.WithObject("inputs",
			mcp.Description("App input variables.")),
		mcp.WithString("conversation_id",
			mcp.Description("Existing conversation UUID.")),
		mcp.WithString("parent_message_id",
			mcp.Description("Parent message ID for threading.")),
		mcp.WithObject("files",
			mcp.Description("Dify file descriptors.")),
	)
}

// ListDatasetsTool returns the dataset listing tool.
func ListDatasetsTool() mcp.Tool {
	return mcp.NewTool("dify_list_datasets",
		mcp.WithDescription("List Dify knowledge bases globally from the Console API."),
		mcp.WithNumber("page",
			mcp.Description("Page number (starts at 1)."),
			mcp.DefaultNumber(1),
			mcp.Min(1)),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of records to return."),
			mcp.Min(1),
			mcp.Max(100),
			mcp.DefaultNumber(20)),
		mcp.WithString("keyword",
			mcp.Description("Search keyword.")),
		mcp.WithBoolean("include_all",
			mcp.Description("Include all datasets regardless of permissions."),
			mcp.DefaultBool(false)),
		mcp.WithObject("ids",
			mcp.Description("Filter by dataset IDs (array of strings).")),
		mcp.WithObject("tag_ids",
			mcp.Description("Filter by tag IDs (array of strings).")),
	)
}

// GetDatasetTool returns the dataset detail tool.
func GetDatasetTool() mcp.Tool {
	return mcp.NewTool("dify_get_dataset",
		mcp.WithDescription("Get knowledge base details by dataset ID."),
		mcp.WithString("dataset_id", mcp.Required(),
			mcp.Description("Dataset UUID.")),
	)
}

// CreateDatasetTool returns the dataset creation tool.
func CreateDatasetTool() mcp.Tool {
	return mcp.NewTool("dify_create_dataset",
		mcp.WithDescription("Create an empty Dify knowledge base."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("Knowledge base name."),
			mcp.MinLength(1),
			mcp.MaxLength(40)),
		mcp.WithString("description",
			mcp.Description("Knowledge base description."),
			mcp.MaxLength(400),
			mcp.DefaultString("")),
		mcp.WithString("indexing_technique",
			mcp.Description("Indexing technique.")),
		mcp.WithString("permission",
			mcp.Description("Access permission."),
			mcp.Enum("only_me", "all_team_members", "partial_members"),
			mcp.DefaultString("only_me")),
		mcp.WithString("provider",
			mcp.Description("Embedding provider."),
			mcp.DefaultString("vendor")),
	)
}

// DeleteDatasetTool returns the dataset deletion tool.
func DeleteDatasetTool() mcp.Tool {
	return mcp.NewTool("dify_delete_dataset",
		mcp.WithDescription("Delete a knowledge base by dataset ID."),
		mcp.WithString("dataset_id", mcp.Required(),
			mcp.Description("Dataset UUID.")),
	)
}

// UpdateDatasetTool returns the dataset update tool.
func UpdateDatasetTool() mcp.Tool {
	return mcp.NewTool("dify_update_dataset",
		mcp.WithDescription("Update knowledge base settings."),
		mcp.WithString("dataset_id", mcp.Required(),
			mcp.Description("Dataset UUID.")),
		mcp.WithObject("update", mcp.Required(),
			mcp.Description("DatasetUpdatePayload fields.")),
	)
}

// ListDatasetDocumentsTool returns the dataset document listing tool.
func ListDatasetDocumentsTool() mcp.Tool {
	return mcp.NewTool("dify_list_dataset_documents",
		mcp.WithDescription("List documents in a knowledge base."),
		mcp.WithString("dataset_id", mcp.Required(),
			mcp.Description("Dataset UUID.")),
		mcp.WithNumber("page",
			mcp.Description("Page number (starts at 1)."),
			mcp.DefaultNumber(1),
			mcp.Min(1)),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of records to return."),
			mcp.Min(1),
			mcp.Max(100),
			mcp.DefaultNumber(20)),
		mcp.WithString("keyword",
			mcp.Description("Search keyword.")),
		mcp.WithBoolean("fetch",
			mcp.Description("Whether to fetch document content."),
			mcp.DefaultBool(false)),
		mcp.WithString("status",
			mcp.Description("Document status filter.")),
		mcp.WithString("sort",
			mcp.Description("Sort field."),
			mcp.DefaultString("-created_at")),
	)
}

// GetDatasetDocumentTool returns the dataset document detail tool.
func GetDatasetDocumentTool() mcp.Tool {
	return mcp.NewTool("dify_get_dataset_document",
		mcp.WithDescription("Get a knowledge base document by document ID."),
		mcp.WithString("dataset_id", mcp.Required(),
			mcp.Description("Dataset UUID.")),
		mcp.WithString("document_id", mcp.Required(),
			mcp.Description("Document UUID.")),
		mcp.WithString("metadata",
			mcp.Description("Metadata inclusion mode."),
			mcp.Enum("all", "only", "without"),
			mcp.DefaultString("all")),
	)
}

// SetDatasetAPIStatusTool returns the dataset API status toggle tool.
func SetDatasetAPIStatusTool() mcp.Tool {
	return mcp.NewTool("dify_set_dataset_api_status",
		mcp.WithDescription("Enable or disable dataset API access for a knowledge base."),
		mcp.WithString("dataset_id", mcp.Required(),
			mcp.Description("Dataset UUID.")),
		mcp.WithBoolean("enabled", mcp.Required(),
			mcp.Description("Whether to enable or disable dataset API access.")),
	)
}

// ListDatasetAPIKeysTool returns the dataset API key listing tool.
func ListDatasetAPIKeysTool() mcp.Tool {
	return mcp.NewTool("dify_list_dataset_api_keys",
		mcp.WithDescription("List workspace-level dataset API keys."),
	)
}

// CreateDatasetAPIKeyTool returns the dataset API key creation tool.
func CreateDatasetAPIKeyTool() mcp.Tool {
	return mcp.NewTool("dify_create_dataset_api_key",
		mcp.WithDescription("Create a workspace-level dataset API key."),
	)
}

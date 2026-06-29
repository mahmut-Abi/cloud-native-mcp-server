package tools

import "github.com/mark3labs/mcp-go/mcp"

// DifyAPIRequestTool returns a tool for calling any relative Dify Service API /v1 path.
func DifyAPIRequestTool() mcp.Tool {
	return mcp.NewTool("dify_dify_api_request",
		mcp.WithDescription("Call any relative Dify Service API /v1 path. Console API paths such as /apps or /datasets are automatically routed through the Console username/password session."),
		mcp.WithString("method",
			mcp.Description("HTTP method."),
			mcp.Enum("GET", "POST", "PUT", "PATCH", "DELETE"),
			mcp.DefaultString("GET")),
		mcp.WithString("path", mcp.Required(),
			mcp.Description("Relative path, for example /info for Service API or /apps for Console API.")),
		mcp.WithString("token_kind",
			mcp.Description("Use app token or dataset token."),
			mcp.Enum("app", "dataset"),
			mcp.DefaultString("app")),
		mcp.WithObject("query",
			mcp.Description("Query string parameters.")),
		mcp.WithObject("body",
			mcp.Description("JSON body for non-GET requests.")),
	)
}

// DifyAppInfoTool returns the basic app info tool.
func DifyAppInfoTool() mcp.Tool {
	return mcp.NewTool("dify_dify_app_info",
		mcp.WithDescription("Get basic Dify app information via GET /v1/info."),
	)
}

// DifyAppMetaTool returns the app metadata tool.
func DifyAppMetaTool() mcp.Tool {
	return mcp.NewTool("dify_dify_app_meta",
		mcp.WithDescription("Get Dify app metadata via GET /v1/meta."),
	)
}

// DifyAppParametersTool returns the app parameters tool.
func DifyAppParametersTool() mcp.Tool {
	return mcp.NewTool("dify_dify_app_parameters",
		mcp.WithDescription("Get Dify app input parameters and configuration via GET /v1/parameters."),
	)
}

// DifyAppSiteTool returns the app site configuration tool.
func DifyAppSiteTool() mcp.Tool {
	return mcp.NewTool("dify_dify_app_site",
		mcp.WithDescription("Get Dify app site configuration via GET /v1/site."),
	)
}

// DifyChatMessageTool returns the chat message tool.
func DifyChatMessageTool() mcp.Tool {
	return mcp.NewTool("dify_dify_chat_message",
		mcp.WithDescription("Send a chat message to a Dify chat, agent chat, or advanced chat app."),
		mcp.WithString("query", mcp.Required(),
			mcp.Description("User message text.")),
		mcp.WithString("conversation_id",
			mcp.Description("Existing conversation UUID.")),
		mcp.WithString("user",
			mcp.Description("End-user identifier. Defaults to DIFY_USER or mcp-user.")),
		mcp.WithObject("inputs",
			mcp.Description("Dify app input variables.")),
		mcp.WithObject("files",
			mcp.Description("Dify file descriptors.")),
		mcp.WithString("response_mode",
			mcp.Description("Use blocking for normal MCP tool calls."),
			mcp.Enum("blocking", "streaming"),
			mcp.DefaultString("blocking")),
		mcp.WithBoolean("auto_generate_name",
			mcp.Description("Whether to auto-generate conversation name."),
			mcp.DefaultBool(true)),
		mcp.WithString("workflow_id",
			mcp.Description("Workflow ID for advanced chat apps.")),
		mcp.WithString("retriever_from",
			mcp.Description("Retriever source."),
			mcp.DefaultString("dev")),
	)
}

// DifyCompletionMessageTool returns the completion message tool.
func DifyCompletionMessageTool() mcp.Tool {
	return mcp.NewTool("dify_dify_completion_message",
		mcp.WithDescription("Create a completion in a Dify completion app."),
		mcp.WithString("query",
			mcp.Description("Optional query text."),
			mcp.DefaultString("")),
		mcp.WithString("user",
			mcp.Description("End-user identifier. Defaults to DIFY_USER or mcp-user.")),
		mcp.WithObject("inputs",
			mcp.Description("Dify app input variables.")),
		mcp.WithObject("files",
			mcp.Description("Dify file descriptors.")),
		mcp.WithString("response_mode",
			mcp.Description("Use blocking for normal MCP tool calls."),
			mcp.Enum("blocking", "streaming"),
			mcp.DefaultString("blocking")),
		mcp.WithString("retriever_from",
			mcp.Description("Retriever source."),
			mcp.DefaultString("dev")),
	)
}

// DifyListConversationsTool returns the conversation listing tool.
func DifyListConversationsTool() mcp.Tool {
	return mcp.NewTool("dify_dify_list_conversations",
		mcp.WithDescription("List conversations for the configured Dify chat app and user."),
		mcp.WithString("user",
			mcp.Description("End-user identifier. Defaults to DIFY_USER or mcp-user.")),
		mcp.WithString("last_id",
			mcp.Description("Last conversation ID for pagination.")),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of records to return."),
			mcp.Min(1),
			mcp.Max(100),
			mcp.DefaultNumber(20)),
		mcp.WithString("sort_by",
			mcp.Description("Sort field."),
			mcp.Enum("created_at", "-created_at", "updated_at", "-updated_at"),
			mcp.DefaultString("-updated_at")),
	)
}

// DifyListMessagesTool returns the message listing tool.
func DifyListMessagesTool() mcp.Tool {
	return mcp.NewTool("dify_dify_list_messages",
		mcp.WithDescription("List messages in a Dify conversation."),
		mcp.WithString("conversation_id", mcp.Required(),
			mcp.Description("Conversation UUID.")),
		mcp.WithString("user",
			mcp.Description("End-user identifier. Defaults to DIFY_USER or mcp-user.")),
		mcp.WithString("first_id",
			mcp.Description("First message ID for pagination.")),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of records to return."),
			mcp.Min(1),
			mcp.Max(100),
			mcp.DefaultNumber(20)),
	)
}

// DifyRetrieveDatasetTool returns the dataset retrieval tool.
func DifyRetrieveDatasetTool() mcp.Tool {
	return mcp.NewTool("dify_dify_retrieve_dataset",
		mcp.WithDescription("Retrieve matching chunks from a Dify dataset via POST /v1/datasets/{dataset_id}/retrieve."),
		mcp.WithString("dataset_id", mcp.Required(),
			mcp.Description("Dataset UUID.")),
		mcp.WithString("query", mcp.Required(),
			mcp.Description("Retrieval query.")),
		mcp.WithObject("retrieval_model",
			mcp.Description("Optional Dify retrieval_model object.")),
		mcp.WithObject("external_retrieval_model",
			mcp.Description("Optional external retrieval model config.")),
		mcp.WithObject("attachment_ids",
			mcp.Description("Optional attachment IDs (array of strings).")),
	)
}

// DifyRunWorkflowTool returns the workflow run tool.
func DifyRunWorkflowTool() mcp.Tool {
	return mcp.NewTool("dify_dify_run_workflow",
		mcp.WithDescription("Run a Dify workflow app. Optionally pass workflow_id to run a specific published workflow version, or app_id to use that app API key from Console instead of DIFY_API_KEY."),
		mcp.WithObject("inputs",
			mcp.Description("Workflow input variables.")),
		mcp.WithString("user",
			mcp.Description("End-user identifier. Defaults to DIFY_USER or mcp-user.")),
		mcp.WithObject("files",
			mcp.Description("Dify file descriptors.")),
		mcp.WithString("response_mode",
			mcp.Description("Use blocking for normal MCP tool calls."),
			mcp.Enum("blocking", "streaming"),
			mcp.DefaultString("blocking")),
		mcp.WithString("app_id",
			mcp.Description("Optional Dify app UUID. When set, the tool uses Console credentials to find this app API key instead of DIFY_API_KEY.")),
		mcp.WithString("workflow_id",
			mcp.Description("Optional published workflow ID. When set, calls POST /v1/workflows/{workflow_id}/run.")),
		mcp.WithBoolean("create_api_key_if_missing",
			mcp.Description("Only used with app_id. Create an app API key if none exists."),
			mcp.DefaultBool(false)),
		mcp.WithBoolean("enable_api_if_disabled",
			mcp.Description("Only used with app_id. Enable Service API for the app if disabled."),
			mcp.DefaultBool(false)),
	)
}

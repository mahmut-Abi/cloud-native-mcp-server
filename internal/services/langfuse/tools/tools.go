package tools

import "github.com/mark3labs/mcp-go/mcp"

// CheckHealthTool returns the Langfuse health check tool.
func CheckHealthTool() mcp.Tool {
	return mcp.NewTool("langfuse_check_health",
		mcp.WithDescription("Check Langfuse API and database health. Use this first when Langfuse requests may be failing or timing out."),
	)
}

// ListTracesSummaryTool returns a compact trace discovery tool.
func ListTracesSummaryTool() mcp.Tool {
	options := append([]mcp.ToolOption{
		mcp.WithDescription("RECOMMENDED: List Langfuse traces with a compact summary view for fast discovery before fetching full trace details."),
	}, commonTraceListOptions()...)
	return mcp.NewTool("langfuse_list_traces_summary", options...)
}

// ListTracesTool returns the full trace listing tool.
func ListTracesTool() mcp.Tool {
	options := append([]mcp.ToolOption{
		mcp.WithDescription("List Langfuse traces with optional filters such as user, session, tags, environment, time range, and field groups."),
	}, commonTraceListOptions()...)
	return mcp.NewTool("langfuse_list_traces", options...)
}

// GetTraceTool returns the trace detail tool.
func GetTraceTool() mcp.Tool {
	return mcp.NewTool("langfuse_get_trace",
		mcp.WithDescription("Get a specific Langfuse trace by ID. Optionally limit the returned field groups for smaller responses."),
		mcp.WithString("trace_id", mcp.Required(),
			mcp.Description("Langfuse trace ID.")),
		mcp.WithString("fields",
			mcp.Description("Optional comma-separated field groups such as `core`, `io`, `scores`, `observations`, or `metrics`.")),
	)
}

// ListSessionsTool returns the session listing tool.
func ListSessionsTool() mcp.Tool {
	options := append([]mcp.ToolOption{
		mcp.WithDescription("List Langfuse sessions with optional time range and environment filters."),
		mcp.WithString("from_timestamp",
			mcp.Description("Only include sessions created at or after this RFC3339 timestamp.")),
		mcp.WithString("to_timestamp",
			mcp.Description("Only include sessions created before this RFC3339 timestamp.")),
		stringArrayOption("environment", "Optional list of environments to match."),
	}, paginationOptions()...)
	return mcp.NewTool("langfuse_list_sessions", options...)
}

// GetSessionTool returns the session detail tool.
func GetSessionTool() mcp.Tool {
	return mcp.NewTool("langfuse_get_session",
		mcp.WithDescription("Get a specific Langfuse session by ID. Useful for correlating a user session with its traces."),
		mcp.WithString("session_id", mcp.Required(),
			mcp.Description("Langfuse session ID.")),
	)
}

// ListObservationsTool returns the observation listing tool.
func ListObservationsTool() mcp.Tool {
	options := append([]mcp.ToolOption{
		mcp.WithDescription("List Langfuse observations such as spans, generations, and events with filtering by trace, level, type, environment, and time."),
		mcp.WithString("name", mcp.Description("Observation name filter.")),
		mcp.WithString("user_id", mcp.Description("Trace user ID filter.")),
		mcp.WithString("type", mcp.Description("Observation type such as `SPAN`, `GENERATION`, or `EVENT`.")),
		mcp.WithString("trace_id", mcp.Description("Associated trace ID filter.")),
		mcp.WithString("level", mcp.Description("Observation level filter such as `DEBUG`, `DEFAULT`, `WARNING`, or `ERROR`.")),
		mcp.WithString("parent_observation_id", mcp.Description("Parent observation ID filter.")),
		stringArrayOption("environment", "Optional list of environments to match."),
		mcp.WithString("from_start_time", mcp.Description("Only include observations starting at or after this RFC3339 timestamp.")),
		mcp.WithString("to_start_time", mcp.Description("Only include observations starting before this RFC3339 timestamp.")),
		mcp.WithString("version", mcp.Description("Version filter.")),
		filterOption("Advanced filter array. Structured arrays are preferred; JSON strings are also accepted."),
	}, paginationOptions()...)
	return mcp.NewTool("langfuse_list_observations", options...)
}

// GetObservationTool returns the observation detail tool.
func GetObservationTool() mcp.Tool {
	return mcp.NewTool("langfuse_get_observation",
		mcp.WithDescription("Get a specific Langfuse observation by ID."),
		mcp.WithString("observation_id", mcp.Required(),
			mcp.Description("Langfuse observation ID.")),
	)
}

// ListPromptsTool returns the prompt listing tool.
func ListPromptsTool() mcp.Tool {
	options := append([]mcp.ToolOption{
		mcp.WithDescription("List Langfuse prompt versions with optional name, label, tag, and update time filters."),
		mcp.WithString("name", mcp.Description("Prompt name filter.")),
		mcp.WithString("label", mcp.Description("Prompt label filter.")),
		mcp.WithString("tag", mcp.Description("Prompt tag filter.")),
		mcp.WithString("from_updated_at", mcp.Description("Only include prompt versions created or updated at or after this RFC3339 timestamp.")),
		mcp.WithString("to_updated_at", mcp.Description("Only include prompt versions created or updated before this RFC3339 timestamp.")),
	}, paginationOptions()...)
	return mcp.NewTool("langfuse_list_prompts", options...)
}

// ListAnnotationQueuesTool returns the annotation queue listing tool.
func ListAnnotationQueuesTool() mcp.Tool {
	options := append([]mcp.ToolOption{
		mcp.WithDescription("List Langfuse annotation queues for manual review workflows."),
	}, paginationOptions()...)
	return mcp.NewTool("langfuse_list_annotation_queues", options...)
}

// GetAnnotationQueueTool returns the annotation queue detail tool.
func GetAnnotationQueueTool() mcp.Tool {
	return mcp.NewTool("langfuse_get_annotation_queue",
		mcp.WithDescription("Get a specific Langfuse annotation queue by ID."),
		mcp.WithString("queue_id", mcp.Required(),
			mcp.Description("Langfuse annotation queue ID.")),
	)
}

// ListAnnotationQueueItemsTool returns the annotation queue item listing tool.
func ListAnnotationQueueItemsTool() mcp.Tool {
	options := append([]mcp.ToolOption{
		mcp.WithDescription("List items inside a specific Langfuse annotation queue."),
		mcp.WithString("queue_id", mcp.Required(),
			mcp.Description("Langfuse annotation queue ID.")),
		mcp.WithString("status",
			mcp.Description("Optional queue item status filter.")),
	}, paginationOptions()...)
	return mcp.NewTool("langfuse_list_annotation_queue_items", options...)
}

// ListDatasetsTool returns the dataset listing tool.
func ListDatasetsTool() mcp.Tool {
	options := append([]mcp.ToolOption{
		mcp.WithDescription("List Langfuse datasets for evaluation and experimentation workflows."),
	}, paginationOptions()...)
	return mcp.NewTool("langfuse_list_datasets", options...)
}

// GetDatasetTool returns the dataset detail tool.
func GetDatasetTool() mcp.Tool {
	return mcp.NewTool("langfuse_get_dataset",
		mcp.WithDescription("Get a specific Langfuse dataset by name."),
		mcp.WithString("dataset_name", mcp.Required(),
			mcp.Description("Langfuse dataset name.")),
	)
}

// ListDatasetRunsTool returns the dataset run listing tool.
func ListDatasetRunsTool() mcp.Tool {
	options := append([]mcp.ToolOption{
		mcp.WithDescription("List runs for a specific Langfuse dataset."),
		mcp.WithString("dataset_name", mcp.Required(),
			mcp.Description("Langfuse dataset name.")),
	}, paginationOptions()...)
	return mcp.NewTool("langfuse_list_dataset_runs", options...)
}

// GetDatasetRunTool returns the dataset run detail tool.
func GetDatasetRunTool() mcp.Tool {
	return mcp.NewTool("langfuse_get_dataset_run",
		mcp.WithDescription("Get a specific Langfuse dataset run by dataset name and run name."),
		mcp.WithString("dataset_name", mcp.Required(),
			mcp.Description("Langfuse dataset name.")),
		mcp.WithString("run_name", mcp.Required(),
			mcp.Description("Dataset run name.")),
	)
}

// ListLLMConnectionsTool returns the LLM connection listing tool.
func ListLLMConnectionsTool() mcp.Tool {
	options := append([]mcp.ToolOption{
		mcp.WithDescription("List Langfuse LLM connections configured for model gateways or providers."),
	}, paginationOptions()...)
	return mcp.NewTool("langfuse_list_llm_connections", options...)
}

// ListModelsTool returns the model listing tool.
func ListModelsTool() mcp.Tool {
	options := append([]mcp.ToolOption{
		mcp.WithDescription("List Langfuse model definitions."),
	}, paginationOptions()...)
	return mcp.NewTool("langfuse_list_models", options...)
}

// GetModelTool returns the model detail tool.
func GetModelTool() mcp.Tool {
	return mcp.NewTool("langfuse_get_model",
		mcp.WithDescription("Get a specific Langfuse model definition by ID."),
		mcp.WithString("model_id", mcp.Required(),
			mcp.Description("Langfuse model ID.")),
	)
}

// ListScoreConfigsTool returns the score config listing tool.
func ListScoreConfigsTool() mcp.Tool {
	options := append([]mcp.ToolOption{
		mcp.WithDescription("List Langfuse score configurations."),
	}, paginationOptions()...)
	return mcp.NewTool("langfuse_list_score_configs", options...)
}

// GetScoreConfigTool returns the score config detail tool.
func GetScoreConfigTool() mcp.Tool {
	return mcp.NewTool("langfuse_get_score_config",
		mcp.WithDescription("Get a specific Langfuse score configuration by ID."),
		mcp.WithString("config_id", mcp.Required(),
			mcp.Description("Langfuse score configuration ID.")),
	)
}

// GetPromptTool returns the prompt detail tool.
func GetPromptTool() mcp.Tool {
	return mcp.NewTool("langfuse_get_prompt",
		mcp.WithDescription("Get a Langfuse prompt by name and optional label or version."),
		mcp.WithString("prompt_name", mcp.Required(),
			mcp.Description("Prompt name. If the prompt lives in a folder, pass the full path such as `folder/subfolder/prompt-name`.")),
		mcp.WithNumber("version",
			mcp.Description("Optional prompt version number.")),
		mcp.WithString("label",
			mcp.Description("Optional label. If neither label nor version is set, Langfuse defaults to `production`.")),
		mcp.WithBoolean("resolve",
			mcp.Description("Whether Langfuse should resolve prompt dependencies before returning the prompt. Defaults to true on the server side.")),
	)
}

// ListScoresTool returns the score listing tool.
func ListScoresTool() mcp.Tool {
	options := append([]mcp.ToolOption{
		mcp.WithDescription("List Langfuse scores with optional filters for trace, session, observation, config, environment, source, and time range."),
		mcp.WithString("user_id", mcp.Description("Trace user ID filter.")),
		mcp.WithString("name", mcp.Description("Score name filter.")),
		mcp.WithString("from_timestamp", mcp.Description("Only include scores created at or after this RFC3339 timestamp.")),
		mcp.WithString("to_timestamp", mcp.Description("Only include scores created before this RFC3339 timestamp.")),
		stringArrayOption("environment", "Optional list of environments to match."),
		mcp.WithString("source", mcp.Description("Score source filter.")),
		mcp.WithString("trace_id", mcp.Description("Trace ID filter.")),
		mcp.WithString("session_id", mcp.Description("Session ID filter.")),
		mcp.WithString("observation_id", mcp.Description("Observation ID filter.")),
		mcp.WithString("config_id", mcp.Description("Score configuration ID filter.")),
		stringArrayOption("trace_tags", "Only scores linked to traces containing all of these tags are returned."),
		mcp.WithString("fields", mcp.Description("Optional comma-separated field groups, for example `score` or `score,trace`.")),
		filterOption("Advanced filter array. Structured arrays are preferred; JSON strings are also accepted."),
	}, paginationOptions()...)
	return mcp.NewTool("langfuse_list_scores", options...)
}

// GetScoreTool returns the score detail tool.
func GetScoreTool() mcp.Tool {
	return mcp.NewTool("langfuse_get_score",
		mcp.WithDescription("Get a specific Langfuse score by ID."),
		mcp.WithString("score_id", mcp.Required(),
			mcp.Description("Langfuse score ID.")),
	)
}

// GetMetricsTool returns the metrics query tool.
func GetMetricsTool() mcp.Tool {
	return mcp.NewTool("langfuse_get_metrics",
		mcp.WithDescription("Execute a Langfuse metrics query. Pass the `query` object using the same shape as the Langfuse metrics API."),
		mcp.WithObject("query", mcp.Required(),
			mcp.Description("Metrics query object. Example: {\"view\":\"traces\",\"dimensions\":[{\"field\":\"name\"}],\"metrics\":[{\"measure\":\"count\",\"aggregation\":\"count\"}]}.")),
	)
}

// GetProjectTool returns the project associated with the configured credentials.
func GetProjectTool() mcp.Tool {
	return mcp.NewTool("langfuse_get_project",
		mcp.WithDescription("Get the Langfuse project associated with the configured project-scoped credentials."),
	)
}

// ListOrganizationProjectsTool returns organization project summaries.
func ListOrganizationProjectsTool() mcp.Tool {
	return mcp.NewTool("langfuse_list_organization_projects",
		mcp.WithDescription("List Langfuse projects in the organization. Requires organization-scoped Langfuse credentials configured as username/password."),
	)
}

// CreateProjectTool returns a tool for creating Langfuse projects.
func CreateProjectTool() mcp.Tool {
	return mcp.NewTool("langfuse_create_project",
		mcp.WithDescription("Create a Langfuse project. Requires organization-scoped credentials. retention_days defaults to 0, which means no retention limit."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("Project name.")),
		mcp.WithObject("metadata",
			mcp.Description("Optional project metadata object.")),
		mcp.WithNumber("retention_days",
			mcp.Description("Number of days to retain data. Use 0 for no retention limit; otherwise Langfuse requires at least 3 days.")),
	)
}

// UpdateProjectTool returns a tool for updating Langfuse projects.
func UpdateProjectTool() mcp.Tool {
	return mcp.NewTool("langfuse_update_project",
		mcp.WithDescription("Update a Langfuse project by ID. Requires organization-scoped credentials."),
		mcp.WithString("project_id", mcp.Required(),
			mcp.Description("Langfuse project ID.")),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("Updated project name.")),
		mcp.WithObject("metadata",
			mcp.Description("Optional project metadata object.")),
		mcp.WithNumber("retention_days",
			mcp.Description("Optional number of days to retain data. Use 0 for no retention limit; otherwise Langfuse requires at least 3 days.")),
	)
}

// DeleteProjectTool returns a tool for deleting Langfuse projects.
func DeleteProjectTool() mcp.Tool {
	return mcp.NewTool("langfuse_delete_project",
		mcp.WithDescription("Delete a Langfuse project by ID. Requires organization-scoped credentials. Project deletion is asynchronous and should be explicitly confirmed before use."),
		mcp.WithString("project_id", mcp.Required(),
			mcp.Description("Langfuse project ID.")),
	)
}

// ListProjectMembershipsTool returns project membership listing.
func ListProjectMembershipsTool() mcp.Tool {
	return mcp.NewTool("langfuse_list_project_memberships",
		mcp.WithDescription("List memberships for a Langfuse project. Requires organization-scoped Langfuse credentials configured as username/password."),
		mcp.WithString("project_id", mcp.Required(),
			mcp.Description("Langfuse project ID.")),
	)
}

// UpsertProjectMembershipTool returns a tool for creating or updating project memberships.
func UpsertProjectMembershipTool() mcp.Tool {
	return mcp.NewTool("langfuse_upsert_project_membership",
		mcp.WithDescription("Create or update a Langfuse project membership. Requires organization-scoped credentials."),
		mcp.WithString("project_id", mcp.Required(),
			mcp.Description("Langfuse project ID.")),
		mcp.WithString("user_id", mcp.Required(),
			mcp.Description("Langfuse user ID to add or update.")),
		mcp.WithString("role", mcp.Required(),
			mcp.Description("Membership role. Use OWNER, ADMIN, MEMBER, or VIEWER.")),
	)
}

// DeleteProjectMembershipTool returns a tool for deleting project memberships.
func DeleteProjectMembershipTool() mcp.Tool {
	return mcp.NewTool("langfuse_delete_project_membership",
		mcp.WithDescription("Delete a Langfuse project membership by user ID. Requires organization-scoped credentials. Confirm the project_id and user_id before calling."),
		mcp.WithString("project_id", mcp.Required(),
			mcp.Description("Langfuse project ID.")),
		mcp.WithString("user_id", mcp.Required(),
			mcp.Description("Langfuse user ID to remove.")),
	)
}

// ListOrganizationAPIKeysTool returns organization-scoped API key summaries.
func ListOrganizationAPIKeysTool() mcp.Tool {
	return mcp.NewTool("langfuse_list_organization_api_keys",
		mcp.WithDescription("List Langfuse organization API keys. Requires organization-scoped Langfuse credentials configured as username/password."),
	)
}

// ListProjectAPIKeysTool returns project API key summaries.
func ListProjectAPIKeysTool() mcp.Tool {
	return mcp.NewTool("langfuse_list_project_api_keys",
		mcp.WithDescription("List Langfuse project API keys for one project. Requires organization-scoped Langfuse credentials configured as username/password."),
		mcp.WithString("project_id", mcp.Required(),
			mcp.Description("Langfuse project ID.")),
	)
}

// CreateProjectAPIKeyTool returns a tool for creating project API keys.
func CreateProjectAPIKeyTool() mcp.Tool {
	return mcp.NewTool("langfuse_create_project_api_key",
		mcp.WithDescription("Create a Langfuse project API key. Requires organization-scoped credentials. The returned secretKey is only available in the creation response."),
		mcp.WithString("project_id", mcp.Required(),
			mcp.Description("Langfuse project ID.")),
		mcp.WithString("note",
			mcp.Description("Optional note for the API key.")),
		mcp.WithString("public_key",
			mcp.Description("Optional predefined public key. Must start with `pk-lf-`; if set, secret_key is required.")),
		mcp.WithString("secret_key",
			mcp.Description("Optional predefined secret key. Must start with `sk-lf-`; if set, public_key is required.")),
	)
}

// DeleteProjectAPIKeyTool returns a tool for deleting project API keys.
func DeleteProjectAPIKeyTool() mcp.Tool {
	return mcp.NewTool("langfuse_delete_project_api_key",
		mcp.WithDescription("Delete a Langfuse project API key by ID. Requires organization-scoped credentials. Confirm the project_id and api_key_id before calling."),
		mcp.WithString("project_id", mcp.Required(),
			mcp.Description("Langfuse project ID.")),
		mcp.WithString("api_key_id", mcp.Required(),
			mcp.Description("Langfuse API key ID to delete.")),
	)
}

func paginationOptions() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithNumber("page",
			mcp.Description("Page number starting at 1.")),
		mcp.WithNumber("limit",
			mcp.Description("Items per page. Prefer smaller values when exploring large projects.")),
	}
}

func commonTraceListOptions() []mcp.ToolOption {
	options := paginationOptions()
	options = append(options,
		mcp.WithString("user_id", mcp.Description("Trace user ID filter.")),
		mcp.WithString("name", mcp.Description("Trace name filter.")),
		mcp.WithString("session_id", mcp.Description("Session ID filter.")),
		mcp.WithString("from_timestamp", mcp.Description("Only include traces whose timestamp is at or after this RFC3339 timestamp.")),
		mcp.WithString("to_timestamp", mcp.Description("Only include traces whose timestamp is before this RFC3339 timestamp.")),
		mcp.WithString("order_by", mcp.Description("Sort expression such as `timestamp.desc` or `name.asc`.")),
		stringArrayOption("tags", "Only traces containing all of these tags are returned."),
		mcp.WithString("version", mcp.Description("Trace version filter.")),
		mcp.WithString("release", mcp.Description("Trace release filter.")),
		stringArrayOption("environment", "Optional list of environments to match."),
		mcp.WithString("fields", mcp.Description("Optional comma-separated field groups such as `core`, `io`, `scores`, `observations`, or `metrics`.")),
		filterOption("Advanced filter array. Structured arrays are preferred; JSON strings are also accepted."),
	)
	return options
}

func stringArrayOption(name, description string) mcp.ToolOption {
	return mcp.WithArray(name,
		mcp.Description(description),
		mcp.Items(map[string]any{
			"type": "string",
		}),
	)
}

func filterOption(description string) mcp.ToolOption {
	return mcp.WithArray("filter",
		mcp.Description(description),
		mcp.Items(map[string]any{
			"type": "object",
		}),
	)
}

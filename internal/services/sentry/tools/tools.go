package tools

import "github.com/mark3labs/mcp-go/mcp"

// TestConnectionTool returns the Sentry connection check tool.
func TestConnectionTool() mcp.Tool {
	return mcp.NewTool("sentry_test_connection",
		mcp.WithDescription("Check whether the configured Sentry token and base URL work. If an organization is provided, it verifies organization/project access as well."),
		mcp.WithString("organization",
			mcp.Description("Optional organization slug. Falls back to the configured default organization.")),
		mcp.WithString("project",
			mcp.Description("Optional project slug. Falls back to the configured default project.")),
	)
}

// ListOrganizationsTool returns the organization listing tool.
func ListOrganizationsTool() mcp.Tool {
	options := append([]mcp.ToolOption{
		mcp.WithDescription("List organizations visible to the current Sentry token. Useful when you need to discover available organization slugs."),
		mcp.WithString("cursor",
			mcp.Description("Opaque cursor value from a previous response to continue pagination.")),
	}, paginationOptions()...)
	return mcp.NewTool("sentry_list_organizations", options...)
}

// ListProjectsTool returns the project listing tool.
func ListProjectsTool() mcp.Tool {
	options := append([]mcp.ToolOption{
		mcp.WithDescription("List Sentry projects for an organization."),
		mcp.WithString("organization",
			mcp.Description("Organization slug. If omitted in some clients, the configured default organization may be used by the handler.")),
		mcp.WithString("query",
			mcp.Description("Optional text query to filter projects.")),
		mcp.WithString("stats_period",
			mcp.Description("Optional Sentry stats period such as `24h`, `7d`, or `14d`.")),
		mcp.WithString("sort_by",
			mcp.Description("Optional sort key such as `date`, `name`, or `issueCount`.")),
		mcp.WithString("cursor",
			mcp.Description("Opaque cursor value from a previous response to continue pagination.")),
	}, paginationOptions()...)
	return mcp.NewTool("sentry_list_projects", options...)
}

// GetProjectTool returns the project detail tool.
func GetProjectTool() mcp.Tool {
	return mcp.NewTool("sentry_get_project",
		mcp.WithDescription("Get a specific Sentry project by organization slug and project slug."),
		mcp.WithString("organization",
			mcp.Description("Organization slug. Falls back to configured default organization when available.")),
		mcp.WithString("project",
			mcp.Description("Project slug. Falls back to configured default project when available.")),
	)
}

// ListIssuesSummaryTool returns the compact issue listing tool.
func ListIssuesSummaryTool() mcp.Tool {
	options := append([]mcp.ToolOption{
		mcp.WithDescription("RECOMMENDED: List Sentry issues in a compact summary view before fetching full issue details."),
		mcp.WithString("organization",
			mcp.Description("Organization slug. Falls back to configured default organization when available.")),
	}, commonIssueListOptions()...)
	return mcp.NewTool("sentry_list_issues_summary", options...)
}

// ListIssuesTool returns the full issue listing tool.
func ListIssuesTool() mcp.Tool {
	options := append([]mcp.ToolOption{
		mcp.WithDescription("List Sentry issues for an organization with optional search, project, environment, and sorting filters."),
		mcp.WithString("organization",
			mcp.Description("Organization slug. Falls back to configured default organization when available.")),
	}, commonIssueListOptions()...)
	return mcp.NewTool("sentry_list_issues", options...)
}

// GetIssueTool returns the issue detail tool.
func GetIssueTool() mcp.Tool {
	return mcp.NewTool("sentry_get_issue",
		mcp.WithDescription("Get a specific Sentry issue by issue ID."),
		mcp.WithString("issue_id", mcp.Required(),
			mcp.Description("Numeric or string Sentry issue ID.")),
	)
}

// ListIssueEventsTool returns the issue events listing tool.
func ListIssueEventsTool() mcp.Tool {
	options := append([]mcp.ToolOption{
		mcp.WithDescription("List events for a Sentry issue."),
		mcp.WithString("issue_id", mcp.Required(),
			mcp.Description("Sentry issue ID.")),
		mcp.WithString("query",
			mcp.Description("Optional search query for issue events.")),
		mcp.WithString("sort",
			mcp.Description("Optional sort field such as `date`, `new`, or `freq`.")),
		stringArrayOption("environment", "Optional list of Sentry environments to filter by."),
		mcp.WithBoolean("full",
			mcp.Description("Whether to return full event payloads instead of compact event records.")),
		mcp.WithBoolean("sample",
			mcp.Description("Whether Sentry may return a sample of events.")),
		mcp.WithString("cursor",
			mcp.Description("Opaque cursor value from a previous response to continue pagination.")),
	}, paginationOptions()...)
	return mcp.NewTool("sentry_list_issue_events", options...)
}

// GetIssueEventTool returns the issue event detail tool.
func GetIssueEventTool() mcp.Tool {
	return mcp.NewTool("sentry_get_issue_event",
		mcp.WithDescription("Get a specific event for a Sentry issue."),
		mcp.WithString("issue_id", mcp.Required(),
			mcp.Description("Sentry issue ID.")),
		mcp.WithString("event_id", mcp.Required(),
			mcp.Description("Event ID from the issue event list.")),
	)
}

func paginationOptions() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of records to request when the endpoint supports it.")),
	}
}

func commonIssueListOptions() []mcp.ToolOption {
	options := paginationOptions()
	options = append(options,
		mcp.WithString("query",
			mcp.Description("Sentry search query, for example `is:unresolved level:error project:backend`.")),
		mcp.WithString("sort",
			mcp.Description("Sort field such as `date`, `new`, `freq`, `user`, or `inbox`.")),
		mcp.WithString("stats_period",
			mcp.Description("Optional Sentry stats period such as `24h`, `7d`, or `14d`.")),
		mcp.WithString("cursor",
			mcp.Description("Opaque cursor value from a previous response to continue pagination.")),
		stringArrayOption("environment", "Optional list of Sentry environments to filter by."),
		stringArrayOption("project_ids", "Optional list of Sentry project IDs to filter the issue list."),
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

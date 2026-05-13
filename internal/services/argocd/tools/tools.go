package tools

import "github.com/mark3labs/mcp-go/mcp"

func TestConnectionTool() mcp.Tool {
	return mcp.NewTool("argocd_test_connection",
		mcp.WithDescription("Check whether the configured Argo CD server is reachable and authentication works."),
	)
}

func ListApplicationsSummaryTool() mcp.Tool {
	return mcp.NewTool("argocd_list_applications_summary",
		mcp.WithDescription("RECOMMENDED: List Argo CD applications in a compact summary view before fetching full application details."),
		mcp.WithString("project",
			mcp.Description("Optional Argo CD project name filter.")),
		mcp.WithString("selector",
			mcp.Description("Optional label selector understood by Argo CD.")),
		mcp.WithString("repo",
			mcp.Description("Optional repository URL filter.")),
	)
}

func GetApplicationTool() mcp.Tool {
	return mcp.NewTool("argocd_get_application",
		mcp.WithDescription("Get one Argo CD application including sync status, health, destination, source, and operation state."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("Application name.")),
		mcp.WithString("app_namespace",
			mcp.Description("Optional application namespace when Argo CD requires disambiguation.")),
	)
}

func GetApplicationManifestsTool() mcp.Tool {
	return mcp.NewTool("argocd_get_application_manifests",
		mcp.WithDescription("Get rendered manifests for one Argo CD application. Useful for diff inspection and troubleshooting generated resources."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("Application name.")),
		mcp.WithString("app_namespace",
			mcp.Description("Optional application namespace when Argo CD requires disambiguation.")),
		mcp.WithString("revision",
			mcp.Description("Optional Git revision or source revision to render.")),
		mcp.WithString("namespace",
			mcp.Description("Optional destination namespace override for manifest generation.")),
	)
}

func ListProjectsTool() mcp.Tool {
	return mcp.NewTool("argocd_list_projects",
		mcp.WithDescription("List Argo CD projects."),
	)
}

func GetProjectTool() mcp.Tool {
	return mcp.NewTool("argocd_get_project",
		mcp.WithDescription("Get one Argo CD project by name."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("Project name.")),
	)
}

func ListClustersTool() mcp.Tool {
	return mcp.NewTool("argocd_list_clusters",
		mcp.WithDescription("List clusters known to Argo CD, including server addresses and connection state."),
	)
}

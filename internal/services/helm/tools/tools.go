// Package tools provides MCP tool definitions for the Helm service.
// It implements tools for managing Helm releases, charts, repositories, and their integration with other services.
package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
)

// ListReleasesTool returns a tool definition for listing Helm releases.
func ListReleasesTool() mcp.Tool {
	logrus.Debug("Creating ListReleasesTool")
	return mcp.NewTool("helm_list_releases",
		mcp.WithDescription("List all Helm releases in the configured namespace or across all namespaces. This tool provides information about installed Helm charts, including release name, namespace, status, chart version, and application version."),
		mcp.WithString("namespace",
			mcp.Description("The namespace to list releases from. If not specified, the configured default namespace will be used.")),
		mcp.WithBoolean("all_namespaces",
			mcp.Description("List releases across all namespaces. Set to true to get a cluster-wide view of all Helm releases.")),
	)
}

// GetReleaseTool returns a tool definition for getting a Helm release.
func GetReleaseTool() mcp.Tool {
	logrus.Debug("Creating GetReleaseTool")
	return mcp.NewTool("helm_get_release",
		mcp.WithDescription("Get detailed information about a specific Helm release. This tool retrieves comprehensive details about a deployed Helm release including its status, chart version, application version, values used, and deployment history."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("The exact name of the Helm release to retrieve.")),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("The namespace where the release is deployed.")),
	)
}

// ListRepositoriesTool returns a tool definition for listing Helm repositories.
func ListRepositoriesTool() mcp.Tool {
	logrus.Debug("Creating ListRepositoriesTool")
	return mcp.NewTool("helm_list_repos",
		mcp.WithDescription("List all configured Helm repositories. This tool retrieves all Helm chart repositories configured in the system, including their names and URLs."),
	)
}

// InstallReleaseTool returns a tool definition for installing Helm releases.
func InstallReleaseTool() mcp.Tool {
	logrus.Debug("Creating InstallReleaseTool")
	return mcp.NewTool("helm_install_release",
		mcp.WithDescription("Install a Helm chart and create a new release. This tool deploys a chart to the cluster with specified configuration."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("The name of the release to install.")),
		mcp.WithString("chart", mcp.Required(),
			mcp.Description("The chart to install (e.g., 'bitnami/nginx', '/path/to/local/chart').")),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("The namespace to install the release into.")),
		mcp.WithString("values_file",
			mcp.Description("Path to a values file to use during installation.")),
	)
}

// UninstallReleaseTool returns a tool definition for uninstalling Helm releases.
func UninstallReleaseTool() mcp.Tool {
	logrus.Debug("Creating UninstallReleaseTool")
	return mcp.NewTool("helm_uninstall_release",
		mcp.WithDescription("Uninstall a Helm release from the cluster. This tool removes a release and all its associated resources."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("The name of the release to uninstall.")),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("The namespace where the release is deployed.")),
	)
}

// UpgradeReleaseTool returns a tool definition for upgrading Helm releases.
func UpgradeReleaseTool() mcp.Tool {
	logrus.Debug("Creating UpgradeReleaseTool")
	return mcp.NewTool("helm_upgrade_release",
		mcp.WithDescription("Upgrade an existing Helm release to a new version or with new values. This tool updates a release with new configuration."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("The name of the release to upgrade.")),
		mcp.WithString("chart", mcp.Required(),
			mcp.Description("The chart to upgrade to.")),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("The namespace where the release is deployed.")),
		mcp.WithString("values_file",
			mcp.Description("Path to a values file to use during upgrade.")),
	)
}

// RollbackReleaseTool returns a tool definition for rolling back Helm releases.
func RollbackReleaseTool() mcp.Tool {
	logrus.Debug("Creating RollbackReleaseTool")
	return mcp.NewTool("helm_rollback_release",
		mcp.WithDescription("Rollback a Helm release to a previous revision. This tool reverts a release to an earlier state."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("The name of the release to rollback.")),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("The namespace where the release is deployed.")),
		mcp.WithNumber("revision",
			mcp.Description("The revision to rollback to. If not specified, rolls back to the previous revision.")),
	)
}

// GetReleaseValuesTool returns a tool definition for getting Helm release values.
func GetReleaseValuesTool() mcp.Tool {
	logrus.Debug("Creating GetReleaseValuesTool")
	return mcp.NewTool("helm_get_release_values",
		mcp.WithDescription("Get the values of a specific Helm release. Returns the actual values used during installation, including both user-provided and default values."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("The name of the release to get values from.")),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("The namespace where the release is deployed.")),
		mcp.WithBoolean("all",
			mcp.Description("If true, return both user-provided and computed values. If false, return only user-provided values.")),
	)
}

// GetReleaseManifestTool returns a tool definition for getting Helm release manifest.
func GetReleaseManifestTool() mcp.Tool {
	logrus.Debug("Creating GetReleaseManifestTool")
	return mcp.NewTool("helm_get_release_manifest",
		mcp.WithDescription("Get the rendered Kubernetes manifest of a Helm release. Shows the actual YAML that was deployed to the cluster."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("The name of the release to get the manifest for.")),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("The namespace where the release is deployed.")),
	)
}

// GetReleaseHistoryTool returns a tool definition for getting Helm release history.
func GetReleaseHistoryTool() mcp.Tool {
	logrus.Debug("Creating GetReleaseHistoryTool")
	return mcp.NewTool("helm_get_release_history",
		mcp.WithDescription("Get the release history (all revisions) of a Helm release. Shows all previous versions and their status."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("The name of the release to get history for.")),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("The namespace where the release is deployed.")),
		mcp.WithNumber("max",
			mcp.Description("The maximum number of revisions to return. If not specified, returns all revisions.")),
	)
}

// SearchChartsTool returns a tool definition for searching Helm charts.
func SearchChartsTool() mcp.Tool {
	logrus.Debug("Creating SearchChartsTool")
	return mcp.NewTool("helm_search_charts",
		mcp.WithDescription("Search for Helm charts in configured repositories. Allows discovering available charts and their versions."),
		mcp.WithString("keyword", mcp.Required(),
			mcp.Description("The keyword or chart name to search for.")),
		mcp.WithBoolean("devel",
			mcp.Description("Include development versions (e.g., beta, alpha) in the search results.")),
	)
}

// GetChartInfoTool returns a tool definition for getting Helm chart information.
func GetChartInfoTool() mcp.Tool {
	logrus.Debug("Creating GetChartInfoTool")
	return mcp.NewTool("helm_get_chart_info",
		mcp.WithDescription("Get detailed information about a Helm chart, including description, version, maintainers, and dependencies."),
		mcp.WithString("chart", mcp.Required(),
			mcp.Description("The chart reference (e.g., 'bitnami/nginx' or path to local chart).")),
	)
}

// GetReleaseStatusTool returns a tool definition for getting Helm release status.
func GetReleaseStatusTool() mcp.Tool {
	logrus.Debug("Creating GetReleaseStatusTool")
	return mcp.NewTool("helm_get_release_status",
		mcp.WithDescription("Get the status of a Helm release. Shows whether the release and its resources are deployed successfully."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("The name of the release to check the status for.")),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("The namespace where the release is deployed.")),
	)
}

// TemplateChartTool returns a tool definition for rendering a Helm chart template.
func TemplateChartTool() mcp.Tool {
	logrus.Debug("Creating TemplateChartTool")
	return mcp.NewTool("helm_template_chart",
		mcp.WithDescription("Render a Helm chart template locally without installing it to the cluster. Useful for previewing what will be deployed."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("The release name to use for templating.")),
		mcp.WithString("chart", mcp.Required(),
			mcp.Description("The chart to template (e.g., 'bitnami/nginx' or path to local chart).")),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("The namespace to use for templating.")),
		mcp.WithString("values_file",
			mcp.Description("Path to a values file to use during templating.")),
	)
}

// CompareReleaseVersionsTool returns a tool definition for comparing release versions.
func CompareReleaseVersionsTool() mcp.Tool {
	logrus.Debug("Creating CompareReleaseVersionsTool")
	return mcp.NewTool("helm_compare_revisions",
		mcp.WithDescription("Compare two revisions of a Helm release to see what changed between them."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("The name of the release to compare.")),
		mcp.WithString("namespace", mcp.Required(),
			mcp.Description("The namespace where the release is deployed.")),
		mcp.WithNumber("revision1", mcp.Required(),
			mcp.Description("The first revision number to compare.")),
		mcp.WithNumber("revision2", mcp.Required(),
			mcp.Description("The second revision number to compare.")),
	)
}

// AddRepositoryTool returns a tool definition for adding a Helm repository.
func AddRepositoryTool() mcp.Tool {
	logrus.Debug("Creating AddRepositoryTool")
	return mcp.NewTool("helm_add_repository",
		mcp.WithDescription("Add a new Helm chart repository to the system. Allows accessing charts from the new repository."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("The name to give the repository.")),
		mcp.WithString("url", mcp.Required(),
			mcp.Description("The URL of the Helm repository.")),
	)
}

// RemoveRepositoryTool returns a tool definition for removing a Helm repository.
func RemoveRepositoryTool() mcp.Tool {
	logrus.Debug("Creating RemoveRepositoryTool")
	return mcp.NewTool("helm_remove_repository",
		mcp.WithDescription("Remove a Helm chart repository from the system."),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("The name of the repository to remove.")),
	)
}

// UpdateRepositoriesTool returns a tool definition for updating Helm repositories.
func UpdateRepositoriesTool() mcp.Tool {
	logrus.Debug("Creating UpdateRepositoriesTool")
	return mcp.NewTool("helm_update_repositories",
		mcp.WithDescription("Update all configured Helm repositories to get the latest chart information. Should be run periodically to see new chart versions."),
	)
}

// GetMirrorConfigurationTool returns a tool definition for checking mirror configuration.
func GetMirrorConfigurationTool() mcp.Tool {
	logrus.Debug("Creating GetMirrorConfigurationTool")
	return mcp.NewTool("helm_get_mirror_configuration",
		mcp.WithDescription("Get information about the configured Helm repository mirrors. Shows whether mirrors are enabled, which repositories are mirrored, and optimization settings."),
	)
}

package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/release"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/helm/client"
	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
)

// HandleListReleases returns a handler function for listing Helm releases.
func HandleListReleases(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "list_helm_releases").Debug("Handler invoked")

		// Get optional parameters
		allNamespaces := getOptionalBoolParam(request, "all_namespaces")
		namespace := getOptionalStringParam(request, "namespace")

		// Call Helm client with parameters
		releases, err := c.ListReleasesAsMap(allNamespaces, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to list releases: %w", err)
		}

		logrus.WithField("count", len(releases)).Debug("list_helm_releases succeeded")

		// Serialize to JSON for better readability
		jsonData, err := marshalIndentJSON(releases)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleGetRelease returns a handler function for getting a Helm release.
func HandleGetRelease(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "get_helm_release").Debug("Handler invoked")

		// Validate required parameters
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}

		// Get namespace parameter (required for production use)
		namespace, err := requireStringParam(request, "namespace")
		if err != nil {
			return nil, err
		}

		// Call Helm client
		release, err := c.GetReleaseAsMap(name, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to get release %s in namespace %s: %w", name, namespace, err)
		}

		logrus.WithField("release", name).Debug("get_helm_release succeeded")

		// Serialize to JSON for better readability
		jsonData, err := marshalIndentJSON(release)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleInstallRelease installs a Helm release.
func HandleInstallRelease(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "install_helm_release").Debug("Handler invoked")

		// Validate required parameters
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}
		chart, err := requireStringParam(request, "chart")
		if err != nil {
			return nil, err
		}
		namespace, err := requireStringParam(request, "namespace")
		if err != nil {
			return nil, err
		}

		// Get values parameter if provided
		values := make(map[string]interface{})
		if valuesArg, ok := request.GetArguments()["values"].(map[string]interface{}); ok {
			values = valuesArg
		}

		// Create a context with 5 minute timeout for installation
		installCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()

		// Create a channel to receive the result
		resultChan := make(chan struct {
			rel *release.Release
			err error
		}, 1)

		// Run the install operation in a goroutine
		go func() {
			rel, err := c.InstallRelease(name, chart, namespace, values)
			resultChan <- struct {
				rel *release.Release
				err error
			}{rel, err}
		}()

		// Wait for either the operation to complete or the context to timeout
		var rel *release.Release
		select {
		case result := <-resultChan:
			rel, err = result.rel, result.err
			if err != nil {
				return nil, fmt.Errorf("failed to install release %s from chart %s in namespace %s: %w", name, chart, namespace, err)
			}
		case <-installCtx.Done():
			return nil, fmt.Errorf("release installation timed out after 5 minutes")
		}

		logrus.WithField("release", name).Debug("install_helm_release succeeded")

		// Convert release to map for better readability
		releaseMap := client.ReleaseToMap(rel)

		// Serialize to JSON
		jsonData, err := marshalIndentJSON(releaseMap)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleUninstallRelease uninstalls a Helm release.
func HandleUninstallRelease(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "uninstall_helm_release").Debug("Handler invoked")

		// Validate required parameters
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}

		// Get namespace parameter (required for production use)
		namespace, err := requireStringParam(request, "namespace")
		if err != nil {
			return nil, err
		}

		// Call Helm client
		result, err := c.UninstallRelease(name, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to uninstall release %s in namespace %s: %w", name, namespace, err)
		}

		logrus.WithField("release", name).Debug("uninstall_helm_release succeeded")

		// Return success message
		message := fmt.Sprintf("Successfully uninstalled release %s in namespace %s", name, namespace)
		if result != nil && result.Release != nil {
			message = fmt.Sprintf("Successfully uninstalled release %s (v%d) in namespace %s", name, result.Release.Version, namespace)
		}
		return mcp.NewToolResultText(message), nil
	}
}

// HandleUpgradeRelease upgrades a Helm release.
func HandleUpgradeRelease(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "upgrade_helm_release").Debug("Handler invoked")

		// Validate required parameters
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}
		chart, err := requireStringParam(request, "chart")
		if err != nil {
			return nil, err
		}

		// Get namespace parameter (required for production use)
		namespace, err := requireStringParam(request, "namespace")
		if err != nil {
			return nil, err
		}

		// Get values parameter if provided
		values := make(map[string]interface{})
		if valuesArg, ok := request.GetArguments()["values"].(map[string]interface{}); ok {
			values = valuesArg
		}

		// Call Helm client
		release, err := c.UpgradeRelease(name, chart, namespace, values)
		if err != nil {
			return nil, fmt.Errorf("failed to upgrade release %s to chart %s in namespace %s: %w", name, chart, namespace, err)
		}

		logrus.WithField("release", name).Debug("upgrade_helm_release succeeded")

		// Convert release to map for better readability
		releaseMap := client.ReleaseToMap(release)

		// Serialize to JSON
		jsonData, err := marshalIndentJSON(releaseMap)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleRollbackRelease rolls back a Helm release.
func HandleRollbackRelease(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "rollback_helm_release").Debug("Handler invoked")

		// Validate required parameters
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}

		// Get namespace parameter (required for production use)
		namespace, err := requireStringParam(request, "namespace")
		if err != nil {
			return nil, err
		}

		// Get optional parameters
		revision := getOptionalIntParam(request, "revision")

		// Call Helm client
		if revision <= 0 {
			revision = 0 // 0 means rollback to previous release
		}
		err = c.RollbackRelease(name, revision, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to rollback release %s in namespace %s: %w", name, namespace, err)
		}

		logrus.WithField("release", name).Debug("rollback_helm_release succeeded")

		// Return success message
		message := fmt.Sprintf("Successfully rolled back release %s in namespace %s", name, namespace)
		if revision > 0 {
			message = fmt.Sprintf("Successfully rolled back release %s to revision %d in namespace %s", name, revision, namespace)
		}
		return mcp.NewToolResultText(message), nil
	}
}

// HandleGetReleaseHistory returns a handler function for getting Helm release history.
func HandleGetReleaseHistory(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_get_release_history").Debug("Handler invoked")

		// Validate required parameters
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}

		// Get namespace parameter (required for production use)
		namespace, err := requireStringParam(request, "namespace")
		if err != nil {
			return nil, err
		}

		max := getOptionalIntParam(request, "max")
		history, err := c.GetReleaseHistoryAsMap(name, namespace, max)
		if err != nil {
			return nil, fmt.Errorf("failed to get release history for %s in namespace %s: %w", name, namespace, err)
		}
		logrus.WithField("release", name).Debug("helm_get_release_history succeeded")
		jsonData, err := marshalIndentJSON(history)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleGetReleaseStatus returns a handler function for getting Helm release status.
func HandleGetReleaseStatus(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_get_release_status").Debug("Handler invoked")

		// Validate required parameters
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}

		// Get namespace parameter (required for production use)
		namespace, err := requireStringParam(request, "namespace")
		if err != nil {
			return nil, err
		}

		status, err := c.GetReleaseStatusAsMap(name, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to get release status for %s in namespace %s: %w", name, namespace, err)
		}
		logrus.WithField("release", name).Debug("helm_get_release_status succeeded")
		jsonData, err := marshalIndentJSON(status)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleGetReleaseManifest returns a handler function for getting Helm release manifest.
func HandleGetReleaseManifest(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_get_release_manifest").Debug("Handler invoked")

		// Validate required parameters
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}

		// Get namespace parameter (required for production use)
		namespace, err := requireStringParam(request, "namespace")
		if err != nil {
			return nil, err
		}

		manifest, err := c.GetReleaseManifest(name, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to get release manifest for %s in namespace %s: %w", name, namespace, err)
		}
		logrus.WithField("release", name).Debug("helm_get_release_manifest succeeded")
		return mcp.NewToolResultText(manifest), nil
	}
}

// HandleCompareRevisions returns a handler function for comparing Helm release revisions.
func HandleCompareRevisions(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_compare_revisions").Debug("Handler invoked")

		// Validate required parameters
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}

		// Get namespace parameter (required for production use)
		namespace, err := requireStringParam(request, "namespace")
		if err != nil {
			return nil, err
		}

		revision1 := getOptionalIntParam(request, "revision1")
		revision2 := getOptionalIntParam(request, "revision2")
		if revision1 <= 0 || revision2 <= 0 {
			return nil, fmt.Errorf("revision1 and revision2 must be positive integers")
		}
		history, err := c.GetReleaseHistoryAsMap(name, namespace, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to get release history for %s in namespace %s: %w", name, namespace, err)
		}
		var rev1, rev2 map[string]interface{}
		for _, relMap := range history {
			if v, ok := relMap["version"].(float64); ok {
				if int(v) == revision1 {
					rev1 = relMap
				}
				if int(v) == revision2 {
					rev2 = relMap
				}
			}
		}
		comparison := map[string]interface{}{
			"release":    name,
			"namespace":  namespace,
			"revision_1": rev1,
			"revision_2": rev2,
		}
		logrus.WithField("release", name).Debug("helm_compare_revisions succeeded")
		jsonData, err := marshalIndentJSON(comparison)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleListReleasesPaginated returns a handler function for listing Helm releases with pagination and optimization.
func HandleListReleasesPaginated(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_list_releases_paginated").Debug("Handler invoked")

		// Get parameters with conservative defaults
		namespace := getOptionalStringParam(request, "namespace")
		status := getOptionalStringParam(request, "status")
		continueToken := getOptionalStringParam(request, "continueToken")
		includeLabels := getOptionalStringParam(request, "includeLabels")

		limit := getOptionalIntParam(request, "limit")
		if limit <= 0 || limit > 100 {
			limit = 50 // Default to 50 to prevent context overflow
		}

		// Call Helm client with pagination
		releases, hasMore, newContinueToken, err := c.ListReleasesPaginated(limit, continueToken, namespace, status)
		if err != nil {
			return nil, fmt.Errorf("failed to list releases: %w", err)
		}

		// Apply label filtering if requested
		if includeLabels != "" {
			releases = c.FilterReleasesByLabels(releases, includeLabels)
		}

		// Create summary response
		response := map[string]interface{}{
			"releases": releases,
			"pagination": map[string]interface{}{
				"continueToken":   newContinueToken,
				"hasMore":         hasMore,
				"currentPageSize": len(releases),
			},
			"count": len(releases),
		}

		logrus.WithFields(logrus.Fields{
			"count":   len(releases),
			"hasMore": hasMore,
		}).Debug("helm_list_releases_paginated succeeded")

		jsonData, err := optimize.GlobalJSONPool.MarshalToBytes(response)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleGetReleaseHistoryPaginated returns a handler function for getting release history with pagination.
func HandleGetReleaseHistoryPaginated(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_get_release_history_paginated").Debug("Handler invoked")

		// Validate required parameters
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}
		namespace, err := requireStringParam(request, "namespace")
		if err != nil {
			return nil, err
		}

		// Get optional parameters
		continueToken := getOptionalStringParam(request, "continueToken")
		includeStatus := getOptionalBoolParam(request, "includeStatus")

		limit := getOptionalIntParam(request, "limit")
		if limit <= 0 || limit > 50 {
			limit = 20 // Default to 20 for history
		}

		// Get paginated history
		history, hasMore, newContinueToken, err := c.GetReleaseHistoryPaginated(name, namespace, limit, continueToken, includeStatus)
		if err != nil {
			return nil, fmt.Errorf("failed to get release history for %s in namespace %s: %w", name, namespace, err)
		}

		response := map[string]interface{}{
			"name":    name,
			"history": history,
			"pagination": map[string]interface{}{
				"continueToken":   newContinueToken,
				"hasMore":         hasMore,
				"currentPageSize": len(history),
			},
			"count": len(history),
		}

		logrus.WithFields(logrus.Fields{
			"release": name,
			"count":   len(history),
		}).Debug("helm_get_release_history_paginated succeeded")

		jsonData, err := optimize.GlobalJSONPool.MarshalToBytes(response)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleGetReleaseSummary returns a brief summary of a Helm release
func HandleGetReleaseSummary(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_get_release_summary").Debug("Handler invoked")

		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}
		namespace, err := requireStringParam(request, "namespace")
		if err != nil {
			return nil, err
		}

		summary, err := c.GetReleaseSummary(name, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to get release summary: %w", err)
		}

		logrus.WithField("release", name).Debug("helm_get_release_summary succeeded")
		jsonData, err := optimize.GlobalJSONPool.MarshalToBytes(summary)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleGetListReleasesSummary returns a list of release summaries
func HandleGetListReleasesSummary(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_list_releases_summary").Debug("Handler invoked")

		namespace := getOptionalStringParam(request, "namespace")
		limit := getOptionalIntParam(request, "limit")
		offset := getOptionalIntParam(request, "offset")

		summaries, err := c.GetListReleasesSummary(namespace, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to get releases summary: %w", err)
		}

		response := map[string]interface{}{
			"releases": summaries,
			"count":    len(summaries),
		}

		logrus.WithField("count", len(summaries)).Debug("helm_list_releases_summary succeeded")
		jsonData, err := optimize.GlobalJSONPool.MarshalToBytes(response)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleListReleasesByNamespace lists all releases in a specific namespace
func HandleListReleasesByNamespace(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_list_releases_in_namespace").Debug("Handler invoked")

		namespace, err := requireStringParam(request, "namespace")
		if err != nil {
			return nil, err
		}

		summaries, err := c.ListReleasesByNamespace(namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to list releases in namespace: %w", err)
		}

		response := map[string]interface{}{
			"namespace": namespace,
			"releases":  summaries,
			"count":     len(summaries),
		}

		logrus.WithField("count", len(summaries)).Debug("helm_list_releases_in_namespace succeeded")
		jsonData, err := optimize.GlobalJSONPool.MarshalToBytes(response)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleFindReleasesByLabels returns a handler function for finding releases by labels.
func HandleFindReleasesByLabels(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_find_releases_by_labels").Debug("Handler invoked")

		// Validate required parameters
		labelSelector, err := requireStringParam(request, "labelSelector")
		if err != nil {
			return nil, err
		}

		// Get optional parameters
		namespace := getOptionalStringParam(request, "namespace")

		limit := getOptionalIntParam(request, "limit")
		if limit <= 0 || limit > 100 {
			limit = 30 // Default to 30 for label searches
		}

		// Find releases by labels
		releases, err := c.FindReleasesByLabels(labelSelector, limit, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to find releases by labels: %w", err)
		}

		response := map[string]interface{}{
			"labelSelector": labelSelector,
			"releases":      releases,
			"count":         len(releases),
			"namespace":     namespace,
		}

		logrus.WithFields(logrus.Fields{
			"selector": labelSelector,
			"count":    len(releases),
		}).Debug("helm_find_releases_by_labels succeeded")

		jsonData, err := optimize.GlobalJSONPool.MarshalToBytes(response)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleFindReleasesByChart finds releases using a specific chart
func HandleFindReleasesByChart(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_find_releases_by_chart").Debug("Handler invoked")

		chartName, err := requireStringParam(request, "chart_name")
		if err != nil {
			return nil, err
		}
		chartVersion := getOptionalStringParam(request, "chart_version")
		limit := getOptionalIntParam(request, "limit")

		releases, err := c.FindReleasesByChart(chartName, chartVersion, limit)
		if err != nil {
			return nil, fmt.Errorf("failed to find releases by chart: %w", err)
		}

		response := map[string]interface{}{
			"chart_name":    chartName,
			"chart_version": chartVersion,
			"releases":      releases,
			"count":         len(releases),
		}

		logrus.WithFields(logrus.Fields{
			"chart": chartName,
			"count": len(releases),
		}).Debug("helm_find_releases_by_chart succeeded")

		jsonData, err := optimize.GlobalJSONPool.MarshalToBytes(response)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleFindBrokenReleases finds releases with failed or pending status
func HandleFindBrokenReleases(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_find_broken_releases").Debug("Handler invoked")

		namespace := getOptionalStringParam(request, "namespace")
		limit := getOptionalIntParam(request, "limit")

		releases, err := c.FindBrokenReleases(namespace, limit)
		if err != nil {
			return nil, fmt.Errorf("failed to find broken releases: %w", err)
		}

		response := map[string]interface{}{
			"brokenReleases": releases,
			"count":          len(releases),
			"namespace":      namespace,
		}

		logrus.WithField("count", len(releases)).Debug("helm_find_broken_releases succeeded")
		jsonData, err := optimize.GlobalJSONPool.MarshalToBytes(response)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleValidateRelease validates a release configuration
func HandleValidateRelease(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_validate_release").Debug("Handler invoked")

		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}
		namespace, err := requireStringParam(request, "namespace")
		if err != nil {
			return nil, err
		}

		result, err := c.ValidateRelease(name, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to validate release: %w", err)
		}

		logrus.WithField("release", name).Debug("helm_validate_release succeeded")
		jsonData, err := optimize.GlobalJSONPool.MarshalToBytes(result)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleGetRecentFailures returns a handler function for getting recent failed releases.
func HandleGetRecentFailures(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_get_recent_failures").Debug("Handler invoked")

		// Get parameters
		namespace := getOptionalStringParam(request, "namespace")
		includePending := getOptionalBoolParam(request, "includePending")

		limit := getOptionalIntParam(request, "limit")
		if limit <= 0 || limit > 50 {
			limit = 20 // Default to 20 for failures
		}

		// Get failed releases
		failures, err := c.GetFailedReleases(limit, namespace, includePending)
		if err != nil {
			return nil, fmt.Errorf("failed to get failed releases: %w", err)
		}

		response := map[string]interface{}{
			"failedReleases": failures,
			"count":          len(failures),
			"filters": map[string]interface{}{
				"namespace":      namespace,
				"includePending": includePending,
			},
		}

		logrus.WithField("count", len(failures)).Debug("helm_get_recent_failures succeeded")

		jsonData, err := optimize.GlobalJSONPool.MarshalToBytes(response)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleGetClusterOverview returns a handler function for getting cluster-wide Helm overview.
func HandleGetClusterOverview(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_cluster_overview").Debug("Handler invoked")

		// Get parameters
		includeNodes := getOptionalBoolParam(request, "includeNodes")
		includeStorage := getOptionalBoolParam(request, "includeStorage")

		// Get cluster overview
		overview, err := c.GetClusterOverview(includeNodes, includeStorage)
		if err != nil {
			return nil, fmt.Errorf("failed to get cluster overview: %w", err)
		}

		logrus.Debug("helm_cluster_overview succeeded")

		jsonData, err := optimize.GlobalJSONPool.MarshalToBytes(overview)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleGetResourcesOfRelease returns a handler function for getting resources managed by a Helm release.
func HandleGetResourcesOfRelease(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_get_resources_of_release").Debug("Handler invoked")

		// Validate required parameters
		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}
		namespace, err := requireStringParam(request, "namespace")
		if err != nil {
			return nil, err
		}

		// Get optional parameters
		includeStatus := getOptionalBoolParam(request, "includeStatus")

		limit := getOptionalIntParam(request, "limit")
		if limit <= 0 || limit > 200 {
			limit = 50 // Default to 50 for resource lists
		}

		// Get resources of release
		resources, err := c.GetResourcesOfRelease(name, namespace, includeStatus, limit)
		if err != nil {
			return nil, fmt.Errorf("failed to get resources of release: %w", err)
		}

		response := map[string]interface{}{
			"releaseName": name,
			"namespace":   namespace,
			"resources":   resources,
			"count":       len(resources),
		}

		logrus.WithFields(logrus.Fields{
			"release": name,
			"count":   len(resources),
		}).Debug("helm_get_resources_of_release succeeded")

		jsonData, err := optimize.GlobalJSONPool.MarshalToBytes(response)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// HandleGetQuickInfo returns a quick overview of all Helm releases
func HandleGetQuickInfo(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.WithField("tool", "helm_quick_info").Debug("Handler invoked")

		info, err := c.GetQuickInfo()
		if err != nil {
			return nil, fmt.Errorf("failed to get quick info: %w", err)
		}

		logrus.Debug("helm_quick_info succeeded")
		jsonData, err := optimize.GlobalJSONPool.MarshalToBytes(info)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize results: %w", err)
		}
		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

// Package handlers provides HTTP request handlers for Grafana MCP tools.
// These handlers process MCP requests and interact with the Grafana client.
package handlers

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/constants"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/grafana/client"
	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
)

var (
	ErrMissingRequiredParam = errors.New("missing required parameter")
	ErrInvalidParameter     = errors.New("invalid parameter")
)

// Response size limits for Grafana - similar to Kubernetes optimizations
// Use constants from internal/constants package for consistency
const (
	defaultLimit    = constants.DefaultPageSize
	maxLimit        = constants.MaxPageSize
	warningLimit    = constants.WarningPageSize
	datasourceLimit = 15 // Lower limit for datasources (they contain configs)
)

// Helper function to validate and parse limit parameter with warnings
func parseLimitWithWarnings(request mcp.CallToolRequest, toolName string) int64 {
	limit := int64(defaultLimit)
	if v, ok := request.GetArguments()["limit"]; ok {
		if f, ok := v.(float64); ok {
			limit = int64(f)
			if limit <= 0 || limit > maxLimit {
				if limit > maxLimit {
					logrus.WithField("requested", limit).WithField("max", maxLimit).Warn("Limit too high, resetting to safe maximum")
					limit = maxLimit
				} else {
					limit = defaultLimit
				}
			}
			if limit > warningLimit {
				logrus.WithFields(logrus.Fields{
					"tool":  toolName,
					"limit": limit,
				}).Warn("Large limit may cause context overflow, consider using summary tools or pagination")
			}
		}
	}
	return limit
}

// Helper function to validate required string parameter
func requireStringParam(request mcp.CallToolRequest, param string) (string, error) {
	value, ok := request.GetArguments()[param].(string)
	if !ok || value == "" {
		return "", fmt.Errorf("%w: %s", ErrMissingRequiredParam, param)
	}
	return value, nil
}

// Helper function to get optional string parameter
func getOptionalStringParam(request mcp.CallToolRequest, param string) string {
	value, _ := request.GetArguments()[param].(string)
	return value
}

// Helper function to get optional boolean parameter
func getOptionalBoolParam(request mcp.CallToolRequest, param string) *bool {
	if value, ok := request.GetArguments()[param].(bool); ok {
		return &value
	}
	return nil
}

// Helper function to marshal JSON response using pooled encoder
func marshalJSONResponse(data interface{}) (*mcp.CallToolResult, error) {
	jsonResponse, err := optimize.GlobalJSONPool.MarshalToBytes(data)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize response: %w", err)
	}
	return mcp.NewToolResultText(string(jsonResponse)), nil
}

// Helper function to marshal optimized JSON response for LLM
func marshalOptimizedResponse(data any, toolName string) (*mcp.CallToolResult, error) {
	// For now, same as JSON response - utils.go will add optimizations
	return marshalJSONResponse(data)
}

// HandleGetDashboards handles dashboard listing requests with intelligent limits.
func HandleGetDashboards() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		debug := getOptionalStringParam(request, "debug")
		limit := parseLimitWithWarnings(request, "grafana_dashboards")

		logrus.WithFields(logrus.Fields{
			"tool":  "grafana_dashboards",
			"debug": debug,
			"limit": limit,
		}).Debug("Handler invoked")

		// Get all dashboards (Grafana API doesn't support pagination for dashboards)
		dashboards, err := c.GetDashboards(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get dashboards: %w", err)
		}

		// Apply server-side limit and remove heavy dashboard data
		result := make([]client.Dashboard, 0)
		count := 0

		for _, dashboard := range dashboards {
			if count >= int(limit) {
				break
			}

			// Create summary version without heavy dashboard data
			summaryDashboard := client.Dashboard{
				ID:        dashboard.ID,
				UID:       dashboard.UID,
				Title:     dashboard.Title,
				Tags:      dashboard.Tags,
				FolderID:  dashboard.FolderID,
				FolderUID: dashboard.FolderUID,
				URL:       dashboard.URL,
				Version:   dashboard.Version,
				// Omit dashboard.Dashboard (the JSON data) to save space
				// Omit dashboard.Meta to save space
			}
			result = append(result, summaryDashboard)
			count++
		}

		response := map[string]interface{}{
			"dashboards":     result,
			"count":          len(result),
			"totalAvailable": len(dashboards),
			"hasMore":        len(dashboards) > int(limit),
			"metadata": map[string]interface{}{
				"limit":   limit,
				"warning": "Dashboard configurations removed to save space. Use grafana_dashboard tool for full data.",
			},
		}

		logrus.WithFields(logrus.Fields{
			"returned": len(result),
			"total":    len(dashboards),
			"hasMore":  len(dashboards) > int(limit),
		}).Debug("grafana_dashboards succeeded")

		return marshalOptimizedResponse(response, "grafana_dashboards")
	}
}

// HandleGetDashboardsSummary handles getting dashboards with minimal output
func HandleGetDashboardsSummary() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		limit := parseLimitWithWarnings(request, "grafana_dashboards_summary")
		offsetStr := getOptionalStringParam(request, "offset")

		offset := 0
		if offsetStr != "" {
			if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
				offset = o
			}
		}

		// Get all dashboards
		dashboards, err := c.GetDashboards(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Type: "text", Text: fmt.Sprintf("error: %s", err.Error())},
				},
			}, nil
		}

		// Apply offset and limit
		var summaries []map[string]interface{}
		totalCount := len(dashboards)

		start := offset
		if start > totalCount {
			start = totalCount
		}

		end := start + int(limit)
		if end > totalCount {
			end = totalCount
		}

		for i := start; i < end; i++ {
			db := dashboards[i]
			summaries = append(summaries, map[string]interface{}{
				"id":        db.ID,
				"uid":       db.UID,
				"title":     db.Title,
				"folderUID": db.FolderUID,
				"tags":      db.Tags,
			})
		}

		response := map[string]interface{}{
			"dashboards":     summaries,
			"count":          len(summaries),
			"offset":         offset,
			"limit":          limit,
			"totalAvailable": totalCount,
			"hasMore":        end < totalCount,
			"pagination": map[string]interface{}{
				"currentPage": (offset / int(limit)) + 1,
				"totalPages":  (totalCount + int(limit) - 1) / int(limit),
				"nextOffset":  end,
				"hasNext":     end < totalCount,
			},
		}

		return marshalOptimizedResponse(response, "grafana_dashboards_summary")
	}
}

// HandleGetDashboard handles specific dashboard retrieval requests.
func HandleGetDashboard() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		uid, err := requireStringParam(request, "uid")
		if err != nil {
			return nil, err
		}
		debug := getOptionalStringParam(request, "debug")
		logrus.WithFields(logrus.Fields{
			"tool":  "grafana_dashboard",
			"uid":   uid,
			"debug": debug,
		}).Debug("Handler invoked")

		dashboard, err := c.GetDashboard(ctx, uid)
		if err != nil {
			return nil, fmt.Errorf("failed to get dashboard: %w", err)
		}

		// Check if dashboard is too large and truncate if necessary
		dashboardSize := 0
		if dashboardBytes, err := optimize.GlobalJSONPool.MarshalToBytes(dashboard); err == nil {
			dashboardSize = len(dashboardBytes)
		}

		response := map[string]interface{}{
			"dashboard": dashboard,
			"metadata": map[string]interface{}{
				"uid":       uid,
				"sizeBytes": dashboardSize,
				"warning":   "Complete dashboard data returned - may be large",
			},
		}

		// Add size warning
		if dashboardSize > 512*1024 { // 512KB
			logrus.WithFields(logrus.Fields{
				"uid":       uid,
				"sizeBytes": dashboardSize,
				"sizeKB":    dashboardSize / 1024,
			}).Warn("Large dashboard returned")
			response["metadata"].(map[string]interface{})["largeWarning"] = "This dashboard is very large and may cause context issues"
		}

		logrus.WithFields(logrus.Fields{
			"uid":       uid,
			"sizeBytes": dashboardSize,
		}).Debug("grafana_dashboard succeeded")

		return marshalOptimizedResponse(response, "grafana_dashboard")
	}
}

// HandleGetDataSources handles data source listing requests with smart limits.
func HandleGetDataSources() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		debug := getOptionalStringParam(request, "debug")
		limit := parseLimitWithWarnings(request, "grafana_datasources")

		// Use more conservative limit for datasources since they contain configurations
		if limit > datasourceLimit {
			limit = datasourceLimit
			logrus.WithField("limit", limit).Info("Using conservative limit for datasources to prevent overflow")
		}

		logrus.WithFields(logrus.Fields{
			"tool":  "grafana_datasources",
			"debug": debug,
			"limit": limit,
		}).Debug("Handler invoked")

		dataSources, err := c.GetDataSources(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get data sources: %w", err)
		}

		// Apply limit and remove sensitive/heavy data
		result := make([]client.DataSource, 0)
		count := 0

		for _, dataSource := range dataSources {
			if count >= int(limit) {
				break
			}

			// Create summary version without heavy/sensitive config data
			summaryDataSource := client.DataSource{
				ID:        dataSource.ID,
				UID:       dataSource.UID,
				Name:      dataSource.Name,
				Type:      dataSource.Type,
				URL:       dataSource.URL,
				Access:    dataSource.Access,
				Database:  dataSource.Database,
				IsDefault: dataSource.IsDefault,
				// Omit Password for security
				// Omit JSONData to save space (can be large)
			}
			result = append(result, summaryDataSource)
			count++
		}

		response := map[string]interface{}{
			"datasources":    result,
			"count":          len(result),
			"totalAvailable": len(dataSources),
			"hasMore":        len(dataSources) > int(limit),
			"metadata": map[string]interface{}{
				"limit":   limit,
				"warning": "Sensitive configuration data removed. Use grafana_datasource_detail for full config.",
			},
		}

		logrus.WithFields(logrus.Fields{
			"returned": len(result),
			"total":    len(dataSources),
			"hasMore":  len(dataSources) > int(limit),
		}).Debug("grafana_datasources succeeded")

		return marshalOptimizedResponse(response, "grafana_datasources")
	}
}

// HandleGetDataSourcesSummary handles getting data sources with minimal output
func HandleGetDataSourcesSummary() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		limit := parseLimitWithWarnings(request, "grafana_datasources_summary")

		// Use even more conservative limit for summary
		if limit > datasourceLimit {
			limit = datasourceLimit
		}

		datasources, err := c.GetDataSources(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Type: "text", Text: fmt.Sprintf("error: %s", err.Error())},
				},
			}, nil
		}

		var summaries []map[string]interface{}
		count := 0

		for _, ds := range datasources {
			if count >= int(limit) {
				break
			}

			summaries = append(summaries, map[string]interface{}{
				"id":         ds.ID,
				"uid":        ds.UID,
				"name":       ds.Name,
				"type":       ds.Type,
				"is_default": ds.IsDefault,
			})
			count++
		}

		response := map[string]interface{}{
			"datasources":    summaries,
			"count":          len(summaries),
			"limit":          limit,
			"totalAvailable": len(datasources),
			"hasMore":        len(datasources) > int(limit),
		}

		return marshalOptimizedResponse(response, "grafana_datasources_summary")
	}
}

// HandleGetPluginsSummary handles plugin summary requests.
func HandleGetPluginsSummary() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		limit := parseLimitWithWarnings(request, "grafana_plugins_summary")

		plugins, err := c.GetPlugins(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get plugins: %w", err)
		}

		summaries := make([]map[string]interface{}, 0, len(plugins))
		for _, plugin := range plugins {
			summaries = append(summaries, map[string]interface{}{
				"id":            plugin.ID,
				"name":          plugin.Name,
				"type":          plugin.Type,
				"enabled":       plugin.Enabled,
				"pinned":        plugin.Pinned,
				"hasUpdate":     plugin.HasUpdate,
				"defaultNavUrl": plugin.DefaultNavURL,
				"signature":     plugin.Signature,
			})
		}

		if int(limit) < len(summaries) {
			summaries = summaries[:limit]
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"plugins": summaries,
			"count":   len(summaries),
			"total":   len(plugins),
			"hasMore": len(plugins) > len(summaries),
		}, "grafana_plugins_summary")
	}
}

// HandleGetFolders handles folder listing requests with limits.
func HandleGetFolders() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		debug := getOptionalStringParam(request, "debug")
		limit := parseLimitWithWarnings(request, "grafana_folders")

		logrus.WithFields(logrus.Fields{
			"tool":  "grafana_folders",
			"debug": debug,
			"limit": limit,
		}).Debug("Handler invoked")

		folders, err := c.GetFolders(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get folders: %w", err)
		}

		// Apply limit (folders are usually small, but we still enforce limits)
		result := make([]client.Folder, 0)
		count := 0

		for _, folder := range folders {
			if count >= int(limit) {
				break
			}
			result = append(result, folder)
			count++
		}

		response := map[string]interface{}{
			"folders":        result,
			"count":          len(result),
			"totalAvailable": len(folders),
			"hasMore":        len(folders) > int(limit),
			"metadata": map[string]interface{}{
				"limit": limit,
			},
		}

		logrus.WithFields(logrus.Fields{
			"returned": len(result),
			"total":    len(folders),
			"hasMore":  len(folders) > int(limit),
		}).Debug("grafana_folders succeeded")

		return marshalOptimizedResponse(response, "grafana_folders")
	}
}

// HandleGetFolder handles single folder retrieval requests.
func HandleGetFolder() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		uid, err := requireStringParam(request, "uid")
		if err != nil {
			return nil, err
		}

		logrus.WithField("uid", uid).Debug("Handler invoked: grafana_folder_detail")

		folder, err := c.GetFolder(ctx, uid)
		if err != nil {
			return nil, fmt.Errorf("failed to get folder: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"folder": folder,
		}, "grafana_folder_detail")
	}
}

// HandleCreateFolder handles creating a folder.
func HandleCreateFolder() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		title, err := requireStringParam(request, "title")
		if err != nil {
			return nil, err
		}

		req := client.FolderCreateRequest{
			Title: title,
			UID:   getOptionalStringParam(request, "uid"),
		}

		logrus.WithFields(logrus.Fields{
			"title": req.Title,
			"uid":   req.UID,
		}).Debug("Handler invoked: grafana_create_folder")

		folder, err := c.CreateFolder(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to create folder: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"folder":  folder,
			"message": "Folder created successfully",
		}, "grafana_create_folder")
	}
}

// HandleUpdateFolder handles updating a folder.
func HandleUpdateFolder() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		uid, err := requireStringParam(request, "uid")
		if err != nil {
			return nil, err
		}
		title, err := requireStringParam(request, "title")
		if err != nil {
			return nil, err
		}

		req := client.FolderUpdateRequest{
			UID:       uid,
			Title:     title,
			Overwrite: false,
		}
		if v, ok := request.GetArguments()["version"].(float64); ok {
			req.Version = int(v)
		}
		if v, ok := request.GetArguments()["overwrite"].(bool); ok {
			req.Overwrite = v
		}

		logrus.WithFields(logrus.Fields{
			"uid":       req.UID,
			"title":     req.Title,
			"version":   req.Version,
			"overwrite": req.Overwrite,
		}).Debug("Handler invoked: grafana_update_folder")

		folder, err := c.UpdateFolder(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to update folder: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"folder":  folder,
			"message": "Folder updated successfully",
		}, "grafana_update_folder")
	}
}

// HandleDeleteFolder handles deleting a folder.
func HandleDeleteFolder() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		uid, err := requireStringParam(request, "uid")
		if err != nil {
			return nil, err
		}

		forceDeleteRules := false
		if v, ok := request.GetArguments()["forceDeleteRules"].(bool); ok {
			forceDeleteRules = v
		}

		logrus.WithFields(logrus.Fields{
			"uid":              uid,
			"forceDeleteRules": forceDeleteRules,
		}).Debug("Handler invoked: grafana_delete_folder")

		result, err := c.DeleteFolder(ctx, uid, forceDeleteRules)
		if err != nil {
			return nil, fmt.Errorf("failed to delete folder: %w", err)
		}

		return marshalOptimizedResponse(result, "grafana_delete_folder")
	}
}

// HandleGetDataSource handles single datasource retrieval requests.
func HandleGetDataSource() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		uid, err := requireStringParam(request, "uid")
		if err != nil {
			return nil, err
		}

		logrus.WithField("uid", uid).Debug("Handler invoked: grafana_datasource_detail")

		dataSource, err := c.GetDataSource(ctx, uid)
		if err != nil {
			return nil, fmt.Errorf("failed to get datasource: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"datasource": dataSource,
		}, "grafana_datasource_detail")
	}
}

// HandleGetPlugins handles plugin enumeration requests.
func HandleGetPlugins() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		limit := parseLimitWithWarnings(request, "grafana_plugins")

		plugins, err := c.GetPlugins(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get plugins: %w", err)
		}

		if int(limit) < len(plugins) {
			plugins = plugins[:limit]
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"plugins": plugins,
			"count":   len(plugins),
		}, "grafana_plugins")
	}
}

// HandleGetPlugin handles plugin detail requests.
func HandleGetPlugin() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		pluginID, err := requireStringParam(request, "pluginID")
		if err != nil {
			return nil, err
		}

		plugin, err := c.GetPlugin(ctx, pluginID)
		if err != nil {
			return nil, fmt.Errorf("failed to get plugin: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"plugin": plugin,
		}, "grafana_plugin_detail")
	}
}

// HandleGetCurrentUser handles current user retrieval requests.
func HandleGetCurrentUser() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		logrus.Debug("Handler invoked: grafana_current_user")

		user, err := c.GetCurrentUser(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get current user: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"user": user,
		}, "grafana_current_user")
	}
}

// HandleGetUsers handles users listing requests.
func HandleGetUsers() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		logrus.Debug("Handler invoked: grafana_users")

		users, err := c.GetUsers(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get users: %w", err)
		}

		response := map[string]interface{}{
			"users": users,
			"count": len(users),
		}

		return marshalOptimizedResponse(response, "grafana_users")
	}
}

// HandleGetOrganization handles organization retrieval requests.
func HandleGetOrganization() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		logrus.Debug("Handler invoked: grafana_organization")

		org, err := c.GetOrganization(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get organization: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"organization": org,
		}, "grafana_organization")
	}
}

// HandleCheckDatasourceHealth handles datasource health check requests.
func HandleCheckDatasourceHealth() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		uid, err := requireStringParam(request, "uid")
		if err != nil {
			return nil, err
		}

		logrus.WithField("uid", uid).Debug("Handler invoked: grafana_check_datasource_health")

		health, err := c.CheckDatasourceHealth(ctx, uid)
		if err != nil {
			return nil, fmt.Errorf("failed to check datasource health: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"uid":    uid,
			"health": health,
		}, "grafana_check_datasource_health")
	}
}

// HandleGetAlertRules handles alert rules retrieval requests with limits.
func HandleGetAlertRules() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		debug := getOptionalStringParam(request, "debug")
		limit := parseLimitWithWarnings(request, "grafana_alerts")

		logrus.WithFields(logrus.Fields{
			"tool":  "grafana_alerts",
			"debug": debug,
			"limit": limit,
		}).Debug("Handler invoked")

		alertRules, err := c.GetAlertRules(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get alert rules: %w", err)
		}

		// Apply limit (alert rules can be complex)
		result := make([]client.AlertRule, 0)
		count := 0

		for _, rule := range alertRules {
			if count >= int(limit) {
				break
			}
			result = append(result, rule)
			count++
		}

		response := map[string]interface{}{
			"alertRules":     result,
			"count":          len(result),
			"totalAvailable": len(alertRules),
			"hasMore":        len(alertRules) > int(limit),
			"metadata": map[string]interface{}{
				"limit": limit,
			},
		}

		logrus.WithFields(logrus.Fields{
			"returned": len(result),
			"total":    len(alertRules),
			"hasMore":  len(alertRules) > int(limit),
		}).Debug("grafana_alerts succeeded")

		return marshalOptimizedResponse(response, "grafana_alerts")
	}
}

// HandleTestConnection handles connection testing requests.
func HandleTestConnection() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		debug := getOptionalStringParam(request, "debug")
		logrus.WithFields(logrus.Fields{
			"tool":  "grafana_test_connection",
			"debug": debug,
		}).Debug("Handler invoked")

		err := c.TestConnection(ctx)
		if err != nil {
			return nil, fmt.Errorf("connection test failed: %w", err)
		}

		result := map[string]interface{}{
			"status":  "success",
			"message": "Grafana connection test successful",
		}

		logrus.Debug("test_grafana_connection succeeded")
		return marshalJSONResponse(result)
	}
}

// HandleSearchDashboards handles dashboard search requests with results limiting.
func HandleSearchDashboards() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		query := getOptionalStringParam(request, "query")
		tag := getOptionalStringParam(request, "tag")
		folderUID := getOptionalStringParam(request, "folderUID")
		starred := getOptionalBoolParam(request, "starred")
		debug := getOptionalStringParam(request, "debug")

		limit := parseLimitWithWarnings(request, "grafana_search_dashboards")

		logrus.WithFields(logrus.Fields{
			"tool":      "grafana_search_dashboards",
			"query":     query,
			"tag":       tag,
			"folderUID": folderUID,
			"starred":   starred,
			"debug":     debug,
			"limit":     limit,
		}).Debug("Handler invoked")

		// Get all dashboards and filter
		dashboards, err := c.GetDashboards(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to search dashboards: %w", err)
		}

		// Apply client-side filtering with limit
		filteredDashboards := make([]client.Dashboard, 0)
		count := 0

		for _, dashboard := range dashboards {
			if count >= int(limit) {
				break
			}

			// Apply query filter
			if query != "" && !strings.Contains(strings.ToLower(dashboard.Title), strings.ToLower(query)) {
				continue
			}

			// Apply tag filter
			if tag != "" {
				hasTag := false
				for _, t := range dashboard.Tags {
					if t == tag {
						hasTag = true
						break
					}
				}
				if !hasTag {
					continue
				}
			}

			// Apply folder filter
			if folderUID != "" && dashboard.FolderUID != folderUID {
				continue
			}

			// Apply starred filter
			if starred != nil && dashboard.IsStarred != *starred {
				continue
			}

			// Add summary version (without heavy dashboard data)
			summaryDashboard := client.Dashboard{
				ID:        dashboard.ID,
				UID:       dashboard.UID,
				Title:     dashboard.Title,
				Tags:      dashboard.Tags,
				FolderID:  dashboard.FolderID,
				FolderUID: dashboard.FolderUID,
				IsStarred: dashboard.IsStarred,
				URL:       dashboard.URL,
				Version:   dashboard.Version,
			}
			filteredDashboards = append(filteredDashboards, summaryDashboard)
			count++
		}

		response := map[string]interface{}{
			"dashboards": filteredDashboards,
			"count":      len(filteredDashboards),
			"searchCriteria": map[string]interface{}{
				"query":     query,
				"tag":       tag,
				"folderUID": folderUID,
				"starred":   starred,
			},
			"limit":   limit,
			"hasMore": len(dashboards) > len(filteredDashboards),
		}

		logrus.WithField("count", len(filteredDashboards)).Debug("search_grafana_dashboards succeeded")
		return marshalOptimizedResponse(response, "grafana_search_dashboards")
	}
}

// ============ Admin Handlers ============

// HandleListTeams handles listing teams.
func HandleListTeams() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		logrus.Debug("Handler invoked: grafana_list_teams")

		teams, err := c.GetTeams(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list teams: %w", err)
		}

		response := map[string]interface{}{
			"teams": teams,
			"count": len(teams),
		}

		return marshalOptimizedResponse(response, "grafana_list_teams")
	}
}

// HandleListAllRoles handles listing all roles.
func HandleListAllRoles() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		logrus.Debug("Handler invoked: grafana_list_all_roles")

		roles, err := c.GetRoles(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list roles: %w", err)
		}

		response := map[string]interface{}{
			"roles": roles,
			"count": len(roles),
		}

		return marshalOptimizedResponse(response, "grafana_list_all_roles")
	}
}

// HandleGetRoleDetails handles getting role details.
func HandleGetRoleDetails() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		roleUID, err := requireStringParam(request, "roleUID")
		if err != nil {
			return nil, err
		}

		logrus.WithField("roleUID", roleUID).Debug("Handler invoked: grafana_get_role_details")

		role, err := c.GetRoleDetails(ctx, roleUID)
		if err != nil {
			return nil, fmt.Errorf("failed to get role details: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"role": role,
		}, "grafana_get_role_details")
	}
}

// HandleGetRoleAssignments handles getting role assignments.
func HandleGetRoleAssignments() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		roleUID, err := requireStringParam(request, "roleUID")
		if err != nil {
			return nil, err
		}

		logrus.WithField("roleUID", roleUID).Debug("Handler invoked: grafana_get_role_assignments")

		assignments, err := c.GetRoleAssignments(ctx, roleUID)
		if err != nil {
			return nil, fmt.Errorf("failed to get role assignments: %w", err)
		}

		response := map[string]interface{}{
			"assignments": assignments,
			"count":       len(assignments),
		}

		return marshalOptimizedResponse(response, "grafana_get_role_assignments")
	}
}

// HandleListUserRoles handles listing roles for a user.
func HandleListUserRoles() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		userIDStr, err := requireStringParam(request, "userID")
		if err != nil {
			return nil, err
		}

		var userID int
		if _, parseErr := fmt.Sscanf(userIDStr, "%d", &userID); parseErr != nil {
			return nil, fmt.Errorf("%w: userID must be a valid integer", ErrInvalidParameter)
		}

		logrus.WithField("userID", userID).Debug("Handler invoked: grafana_list_user_roles")

		roles, err := c.GetUserRoles(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to list user roles: %w", err)
		}

		response := map[string]interface{}{
			"userID": userID,
			"roles":  roles,
			"count":  len(roles),
		}

		return marshalOptimizedResponse(response, "grafana_list_user_roles")
	}
}

// HandleListTeamRoles handles listing roles for a team.
func HandleListTeamRoles() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		teamIDStr, err := requireStringParam(request, "teamID")
		if err != nil {
			return nil, err
		}

		var teamID int
		if _, parseErr := fmt.Sscanf(teamIDStr, "%d", &teamID); parseErr != nil {
			return nil, fmt.Errorf("%w: teamID must be a valid integer", ErrInvalidParameter)
		}

		logrus.WithField("teamID", teamID).Debug("Handler invoked: grafana_list_team_roles")

		roles, err := c.GetTeamRoles(ctx, teamID)
		if err != nil {
			return nil, fmt.Errorf("failed to list team roles: %w", err)
		}

		response := map[string]interface{}{
			"teamID": teamID,
			"roles":  roles,
			"count":  len(roles),
		}

		return marshalOptimizedResponse(response, "grafana_list_team_roles")
	}
}

// HandleGetResourcePermissions handles getting permissions for a resource.
func HandleGetResourcePermissions() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		resourceType, err := requireStringParam(request, "resourceType")
		if err != nil {
			return nil, err
		}
		resourceUID := getOptionalStringParam(request, "resourceUID")

		logrus.WithFields(logrus.Fields{
			"resourceType": resourceType,
			"resourceUID":  resourceUID,
		}).Debug("Handler invoked: grafana_get_resource_permissions")

		permissions, err := c.GetResourcePermissions(ctx, resourceType, resourceUID)
		if err != nil {
			return nil, fmt.Errorf("failed to get resource permissions: %w", err)
		}

		response := map[string]interface{}{
			"resourceType": resourceType,
			"resourceUID":  resourceUID,
			"permissions":  permissions,
			"count":        len(permissions),
		}

		return marshalOptimizedResponse(response, "grafana_get_resource_permissions")
	}
}

// HandleGetResourceDescription handles getting resource description.
func HandleGetResourceDescription() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		resourceType, err := requireStringParam(request, "resourceType")
		if err != nil {
			return nil, err
		}

		logrus.WithField("resourceType", resourceType).Debug("Handler invoked: grafana_get_resource_description")

		desc, err := c.GetResourceDescription(ctx, resourceType)
		if err != nil {
			return nil, fmt.Errorf("failed to get resource description: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"description": desc,
		}, "grafana_get_resource_description")
	}
}

// ============ Dashboard Update Handlers ============

// HandleUpdateDashboard handles creating or updating a dashboard.
func HandleUpdateDashboard() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		dashboardJSON, ok := request.GetArguments()["dashboard"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%w: dashboard is required", ErrMissingRequiredParam)
		}

		folderUID := getOptionalStringParam(request, "folderUID")
		overwrite := false
		if v, ok := request.GetArguments()["overwrite"].(bool); ok {
			overwrite = v
		}
		message := getOptionalStringParam(request, "message")

		logrus.WithFields(logrus.Fields{
			"folderUID": folderUID,
			"overwrite": overwrite,
		}).Debug("Handler invoked: grafana_update_dashboard")

		req := client.DashboardUpdateRequest{
			Dashboard: dashboardJSON,
			FolderUID: folderUID,
			Overwrite: overwrite,
			Message:   message,
		}

		dashboard, err := c.UpdateDashboard(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to update dashboard: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"dashboard": dashboard,
			"message":   "Dashboard created/updated successfully",
		}, "grafana_update_dashboard")
	}
}

// HandleGetDashboardVersions handles retrieving dashboard version history.
func HandleGetDashboardVersions() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		dashboardUID, err := requireStringParam(request, "dashboardUID")
		if err != nil {
			return nil, err
		}

		limit := 0
		if v, ok := request.GetArguments()["limit"].(float64); ok {
			limit = int(v)
		}
		start := 0
		if v, ok := request.GetArguments()["start"].(float64); ok {
			start = int(v)
		}

		logrus.WithFields(logrus.Fields{
			"dashboardUID": dashboardUID,
			"limit":        limit,
			"start":        start,
		}).Debug("Handler invoked: grafana_get_dashboard_versions")

		versions, err := c.GetDashboardVersions(ctx, dashboardUID, limit, start)
		if err != nil {
			return nil, fmt.Errorf("failed to get dashboard versions: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"dashboardUID": dashboardUID,
			"versions":     versions,
			"count":        len(versions),
		}, "grafana_get_dashboard_versions")
	}
}

// HandleGetDashboardVersion handles retrieving a specific dashboard version.
func HandleGetDashboardVersion() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		dashboardUID, err := requireStringParam(request, "dashboardUID")
		if err != nil {
			return nil, err
		}
		versionValue, ok := request.GetArguments()["version"].(float64)
		if !ok {
			return nil, fmt.Errorf("%w: version", ErrMissingRequiredParam)
		}
		version := int(versionValue)

		logrus.WithFields(logrus.Fields{
			"dashboardUID": dashboardUID,
			"version":      version,
		}).Debug("Handler invoked: grafana_get_dashboard_version")

		dashboardVersion, err := c.GetDashboardVersion(ctx, dashboardUID, version)
		if err != nil {
			return nil, fmt.Errorf("failed to get dashboard version: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"dashboardUID": dashboardUID,
			"version":      dashboardVersion,
		}, "grafana_get_dashboard_version")
	}
}

// HandleRestoreDashboardVersion handles restoring a dashboard version.
func HandleRestoreDashboardVersion() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		dashboardUID, err := requireStringParam(request, "dashboardUID")
		if err != nil {
			return nil, err
		}
		versionValue, ok := request.GetArguments()["version"].(float64)
		if !ok {
			return nil, fmt.Errorf("%w: version", ErrMissingRequiredParam)
		}
		version := int(versionValue)

		logrus.WithFields(logrus.Fields{
			"dashboardUID": dashboardUID,
			"version":      version,
		}).Debug("Handler invoked: grafana_restore_dashboard_version")

		result, err := c.RestoreDashboardVersion(ctx, dashboardUID, version)
		if err != nil {
			return nil, fmt.Errorf("failed to restore dashboard version: %w", err)
		}

		return marshalOptimizedResponse(result, "grafana_restore_dashboard_version")
	}
}

// HandleDeleteDashboard handles deleting a dashboard by UID.
func HandleDeleteDashboard() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		uid, err := requireStringParam(request, "uid")
		if err != nil {
			return nil, err
		}

		logrus.WithField("uid", uid).Debug("Handler invoked: grafana_delete_dashboard")

		result, err := c.DeleteDashboard(ctx, uid)
		if err != nil {
			return nil, fmt.Errorf("failed to delete dashboard: %w", err)
		}

		return marshalOptimizedResponse(result, "grafana_delete_dashboard")
	}
}

// HandleGetDashboardPanelQueries handles getting panel queries from a dashboard.
func HandleGetDashboardPanelQueries() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		dashboardUID, err := requireStringParam(request, "dashboardUID")
		if err != nil {
			return nil, err
		}

		logrus.WithField("dashboardUID", dashboardUID).Debug("Handler invoked: grafana_get_dashboard_panel_queries")

		panels, err := c.GetDashboardPanelQueries(ctx, dashboardUID)
		if err != nil {
			return nil, fmt.Errorf("failed to get dashboard panel queries: %w", err)
		}

		response := map[string]interface{}{
			"dashboardUID": dashboardUID,
			"panels":       panels,
			"count":        len(panels),
		}

		return marshalOptimizedResponse(response, "grafana_get_dashboard_panel_queries")
	}
}

// HandleGetDashboardProperty handles extracting specific dashboard properties.
func HandleGetDashboardProperty() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		dashboardUID, err := requireStringParam(request, "dashboardUID")
		if err != nil {
			return nil, err
		}

		propertyPath, err := requireStringParam(request, "propertyPath")
		if err != nil {
			return nil, err
		}

		logrus.WithFields(logrus.Fields{
			"dashboardUID": dashboardUID,
			"propertyPath": propertyPath,
		}).Debug("Handler invoked: grafana_get_dashboard_property")

		value, err := c.GetDashboardProperty(ctx, dashboardUID, propertyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get dashboard property: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"dashboardUID": dashboardUID,
			"propertyPath": propertyPath,
			"value":        value,
		}, "grafana_get_dashboard_property")
	}
}

// ============ Alerting Handlers ============

// HandleGetAlertRuleByUID handles getting a specific alert rule.
func HandleGetAlertRuleByUID() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		ruleUID, err := requireStringParam(request, "ruleUID")
		if err != nil {
			return nil, err
		}

		logrus.WithField("ruleUID", ruleUID).Debug("Handler invoked: grafana_get_alert_rule_by_uid")

		rule, err := c.GetAlertRuleByUID(ctx, ruleUID)
		if err != nil {
			return nil, fmt.Errorf("failed to get alert rule: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"alertRule": rule,
		}, "grafana_get_alert_rule_by_uid")
	}
}

// HandleCreateAlertRule handles creating a new alert rule.
func HandleCreateAlertRule() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		title, err := requireStringParam(request, "title")
		if err != nil {
			return nil, err
		}

		condition, err := requireStringParam(request, "condition")
		if err != nil {
			return nil, err
		}

		folderUID, err := requireStringParam(request, "folderUID")
		if err != nil {
			return nil, err
		}

		ruleGroup := getOptionalStringParam(request, "ruleGroup")
		if ruleGroup == "" {
			ruleGroup = "default"
		}

		intervalSeconds := 60
		if v, ok := request.GetArguments()["intervalSeconds"].(float64); ok {
			intervalSeconds = int(v)
		}

		// Parse data (queries)
		var data []map[string]interface{}
		if dataArg, ok := request.GetArguments()["data"].([]interface{}); ok {
			for _, item := range dataArg {
				if m, ok := item.(map[string]interface{}); ok {
					data = append(data, m)
				}
			}
		}

		// Parse annotations
		var annotations map[string]string
		if annArg, ok := request.GetArguments()["annotations"].(map[string]interface{}); ok {
			annotations = make(map[string]string)
			for k, v := range annArg {
				if s, ok := v.(string); ok {
					annotations[k] = s
				}
			}
		}

		// Parse labels
		var labels map[string]string
		if labelArg, ok := request.GetArguments()["labels"].(map[string]interface{}); ok {
			labels = make(map[string]string)
			for k, v := range labelArg {
				if s, ok := v.(string); ok {
					labels[k] = s
				}
			}
		}

		logrus.WithField("title", title).Debug("Handler invoked: grafana_create_alert_rule")

		req := client.CreateAlertRuleRequest{
			Title:           title,
			Condition:       condition,
			Data:            data,
			IntervalSeconds: intervalSeconds,
			FolderUID:       folderUID,
			RuleGroup:       ruleGroup,
			Annotations:     annotations,
			Labels:          labels,
		}

		rule, err := c.CreateAlertRule(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to create alert rule: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"alertRule": rule,
			"message":   "Alert rule created successfully",
		}, "grafana_create_alert_rule")
	}
}

// HandleUpdateAlertRule handles updating an existing alert rule.
func HandleUpdateAlertRule() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		ruleUID, err := requireStringParam(request, "ruleUID")
		if err != nil {
			return nil, err
		}

		title := getOptionalStringParam(request, "title")
		condition := getOptionalStringParam(request, "condition")
		folderUID := getOptionalStringParam(request, "folderUID")

		logrus.WithField("ruleUID", ruleUID).Debug("Handler invoked: grafana_update_alert_rule")

		req := client.UpdateAlertRuleRequest{
			Title:     title,
			Condition: condition,
			FolderUID: folderUID,
		}

		rule, err := c.UpdateAlertRule(ctx, ruleUID, req)
		if err != nil {
			return nil, fmt.Errorf("failed to update alert rule: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"alertRule": rule,
			"message":   "Alert rule updated successfully",
		}, "grafana_update_alert_rule")
	}
}

// HandleDeleteAlertRule handles deleting an alert rule.
func HandleDeleteAlertRule() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		ruleUID, err := requireStringParam(request, "ruleUID")
		if err != nil {
			return nil, err
		}

		logrus.WithField("ruleUID", ruleUID).Debug("Handler invoked: grafana_delete_alert_rule")

		if err := c.DeleteAlertRule(ctx, ruleUID); err != nil {
			return nil, fmt.Errorf("failed to delete alert rule: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"ruleUID": ruleUID,
			"message": "Alert rule deleted successfully",
		}, "grafana_delete_alert_rule")
	}
}

// HandleListContactPoints handles listing contact points.
func HandleListContactPoints() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		logrus.Debug("Handler invoked: grafana_list_contact_points")

		contactPoints, err := c.GetContactPoints(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list contact points: %w", err)
		}

		response := map[string]interface{}{
			"contactPoints": contactPoints,
			"count":         len(contactPoints),
		}

		return marshalOptimizedResponse(response, "grafana_list_contact_points")
	}
}

// ============ Annotation Handlers ============

// HandleGetAnnotations handles getting annotations.
func HandleGetAnnotations() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		params := make(map[string]string)

		if v := getOptionalStringParam(request, "dashboardUID"); v != "" {
			params["dashboardUid"] = v
		}
		if v := getOptionalStringParam(request, "from"); v != "" {
			params["from"] = v
		}
		if v := getOptionalStringParam(request, "to"); v != "" {
			params["to"] = v
		}
		if v := getOptionalStringParam(request, "tags"); v != "" {
			params["tags"] = v
		}
		if v := getOptionalStringParam(request, "limit"); v != "" {
			params["limit"] = v
		}

		logrus.WithField("params", params).Debug("Handler invoked: grafana_get_annotations")

		annotations, err := c.GetAnnotations(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to get annotations: %w", err)
		}

		response := map[string]interface{}{
			"annotations": annotations,
			"count":       len(annotations),
		}

		return marshalOptimizedResponse(response, "grafana_get_annotations")
	}
}

// HandleCreateAnnotation handles creating an annotation.
func HandleCreateAnnotation() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		text, err := requireStringParam(request, "text")
		if err != nil {
			return nil, err
		}

		var annotationTime int64
		if t, ok := request.GetArguments()["time"].(float64); ok {
			annotationTime = int64(t)
		} else {
			annotationTime = time.Now().Unix() * 1000
		}

		dashboardUID := getOptionalStringParam(request, "dashboardUID")
		timeEnd := int64(0)
		if t, ok := request.GetArguments()["timeEnd"].(float64); ok {
			timeEnd = int64(t)
		}

		var tags []string
		if tagsArg, ok := request.GetArguments()["tags"].([]interface{}); ok {
			for _, tag := range tagsArg {
				if s, ok := tag.(string); ok {
					tags = append(tags, s)
				}
			}
		}

		var panelID int
		if p, ok := request.GetArguments()["panelID"].(float64); ok {
			panelID = int(p)
		}

		logrus.WithField("text", text).Debug("Handler invoked: grafana_create_annotation")

		req := client.CreateAnnotationRequest{
			DashboardUID: dashboardUID,
			PanelID:      panelID,
			Time:         annotationTime,
			TimeEnd:      timeEnd,
			Text:         text,
			Tags:         tags,
		}

		annotation, err := c.CreateAnnotation(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to create annotation: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"annotation": annotation,
			"message":    "Annotation created successfully",
		}, "grafana_create_annotation")
	}
}

// HandleUpdateAnnotation handles updating an annotation.
func HandleUpdateAnnotation() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		annotationIDStr, err := requireStringParam(request, "annotationID")
		if err != nil {
			return nil, err
		}

		var annotationID int64
		if _, parseErr := fmt.Sscanf(annotationIDStr, "%d", &annotationID); parseErr != nil {
			return nil, fmt.Errorf("%w: annotationID must be a valid integer", ErrInvalidParameter)
		}

		req := client.UpdateAnnotationRequest{
			Text: getOptionalStringParam(request, "text"),
		}

		if t, ok := request.GetArguments()["time"].(float64); ok {
			req.Time = int64(t)
		}
		if t, ok := request.GetArguments()["timeEnd"].(float64); ok {
			req.TimeEnd = int64(t)
		}

		logrus.WithField("annotationID", annotationID).Debug("Handler invoked: grafana_update_annotation")

		annotation, err := c.UpdateAnnotation(ctx, annotationID, req)
		if err != nil {
			return nil, fmt.Errorf("failed to update annotation: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"annotation": annotation,
			"message":    "Annotation updated successfully",
		}, "grafana_update_annotation")
	}
}

// HandlePatchAnnotation handles patching an annotation.
func HandlePatchAnnotation() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		annotationIDStr, err := requireStringParam(request, "annotationID")
		if err != nil {
			return nil, err
		}

		var annotationID int64
		if _, parseErr := fmt.Sscanf(annotationIDStr, "%d", &annotationID); parseErr != nil {
			return nil, fmt.Errorf("%w: annotationID must be a valid integer", ErrInvalidParameter)
		}

		req := client.PatchAnnotationRequest{}

		logrus.WithField("annotationID", annotationID).Debug("Handler invoked: grafana_patch_annotation")

		annotation, err := c.PatchAnnotation(ctx, annotationID, req)
		if err != nil {
			return nil, fmt.Errorf("failed to patch annotation: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"annotation": annotation,
			"message":    "Annotation patched successfully",
		}, "grafana_patch_annotation")
	}
}

// HandleDeleteAnnotation handles deleting an annotation.
func HandleDeleteAnnotation() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		annotationIDStr, err := requireStringParam(request, "annotationID")
		if err != nil {
			return nil, err
		}

		annotationID, err := strconv.ParseInt(annotationIDStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: annotationID must be a valid integer", ErrInvalidParameter)
		}

		logrus.WithField("annotationID", annotationID).Debug("Handler invoked: grafana_delete_annotation")

		if err := c.DeleteAnnotation(ctx, annotationID); err != nil {
			return nil, fmt.Errorf("failed to delete annotation: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"annotationID": annotationID,
			"message":      "Annotation deleted successfully",
		}, "grafana_delete_annotation")
	}
}

// HandleGetAnnotationTags handles getting annotation tags.
func HandleGetAnnotationTags() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		tag := getOptionalStringParam(request, "tag")

		logrus.WithField("tag", tag).Debug("Handler invoked: grafana_get_annotation_tags")

		tags, err := c.GetAnnotationTags(ctx, tag)
		if err != nil {
			return nil, fmt.Errorf("failed to get annotation tags: %w", err)
		}

		response := map[string]interface{}{
			"tags":  tags,
			"count": len(tags),
		}

		return marshalOptimizedResponse(response, "grafana_get_annotation_tags")
	}
}

// HandleGenerateDeeplink handles generating deeplink URLs.
func HandleGenerateDeeplink() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		resourceType, err := requireStringParam(request, "resourceType")
		if err != nil {
			return nil, err
		}

		resourceUID, err := requireStringParam(request, "resourceUID")
		if err != nil {
			return nil, err
		}

		params := make(map[string]string)
		if v := getOptionalStringParam(request, "panelID"); v != "" {
			params["panelId"] = v
		}
		if v := getOptionalStringParam(request, "datasource"); v != "" {
			params["datasource"] = v
		}
		if v := getOptionalStringParam(request, "from"); v != "" {
			params["from"] = v
		}
		if v := getOptionalStringParam(request, "to"); v != "" {
			params["to"] = v
		}

		logrus.WithFields(logrus.Fields{
			"resourceType": resourceType,
			"resourceUID":  resourceUID,
		}).Debug("Handler invoked: grafana_generate_deeplink")

		deeplink, err := c.GenerateDeeplink(ctx, resourceType, resourceUID, params)
		if err != nil {
			return nil, fmt.Errorf("failed to generate deeplink: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"deeplink": deeplink,
		}, "grafana_generate_deeplink")
	}
}

// HandleGenerateLogsDrilldownLink handles generating Logs Drilldown deeplink URLs.
func HandleGenerateLogsDrilldownLink() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		params := map[string]string{}
		if from := getOptionalStringParam(request, "from"); from != "" {
			params["from"] = from
		}
		if to := getOptionalStringParam(request, "to"); to != "" {
			params["to"] = to
		}
		if datasourceUID := getOptionalStringParam(request, "datasourceUid"); datasourceUID != "" {
			params["datasourceUid"] = datasourceUID
		}
		if datasourceName := getOptionalStringParam(request, "datasourceName"); datasourceName != "" {
			params["datasourceName"] = datasourceName
		}
		if orgID := getOptionalStringParam(request, "orgId"); orgID != "" {
			params["orgId"] = orgID
		}

		deeplink, err := c.GenerateLogsDrilldownLink(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to generate logs drilldown link: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"deeplink": deeplink,
		}, "grafana_generate_logs_drilldown_link")
	}
}

// HandleGetDataSourceByName handles getting a datasource by name.
func HandleGetDataSourceByName() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}

		logrus.WithField("name", name).Debug("Handler invoked: grafana_get_datasource_by_name")

		dataSource, err := c.GetDataSourceByName(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("failed to get datasource by name: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"datasource": dataSource,
		}, "grafana_get_datasource_by_name")
	}
}

// ============ Panel Image Rendering Handler ============

// HandleRenderPanelImage handles rendering a dashboard panel to an image.
func HandleRenderPanelImage() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		dashboardUID, err := requireStringParam(request, "dashboardUID")
		if err != nil {
			return nil, err
		}

		panelIDFloat, ok := request.GetArguments()["panelID"].(float64)
		if !ok {
			return nil, fmt.Errorf("%w: panelID is required", ErrMissingRequiredParam)
		}
		panelID := int(panelIDFloat)

		// Optional parameters
		params := make(map[string]string)

		if width, ok := request.GetArguments()["width"].(float64); ok {
			params["width"] = fmt.Sprintf("%d", int(width))
		}
		if height, ok := request.GetArguments()["height"].(float64); ok {
			params["height"] = fmt.Sprintf("%d", int(height))
		}
		if from := getOptionalStringParam(request, "from"); from != "" {
			params["from"] = from
		}
		if to := getOptionalStringParam(request, "to"); to != "" {
			params["to"] = to
		}
		if timeout := getOptionalStringParam(request, "timeout"); timeout != "" {
			params["timeout"] = timeout
		}

		logrus.WithFields(logrus.Fields{
			"dashboardUID": dashboardUID,
			"panelID":      panelID,
		}).Debug("Handler invoked: grafana_render_panel_image")

		image, err := c.RenderDashboardPanel(ctx, dashboardUID, panelID, params)
		if err != nil {
			return nil, fmt.Errorf("failed to render panel: %w", err)
		}

		encoded := base64.StdEncoding.EncodeToString(image.ImageData)
		return marshalOptimizedResponse(map[string]interface{}{
			"contentType": image.ContentType,
			"size":        image.Size,
			"width":       image.Width,
			"height":      image.Height,
			"imageData":   encoded,
			"dataUrl": fmt.Sprintf("data:%s;base64,%s",
				image.ContentType,
				encoded,
			),
		}, "grafana_render_panel_image")
	}
}

// ============ Graphite Annotation Handler ============

// HandleCreateGraphiteAnnotation handles creating a Graphite annotation.
func HandleCreateGraphiteAnnotation() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		what, err := requireStringParam(request, "what")
		if err != nil {
			return nil, err
		}

		data := getOptionalStringParam(request, "data")
		tags := getOptionalStringParam(request, "tags")

		var timestamp int64
		if t, ok := request.GetArguments()["timestamp"].(float64); ok {
			timestamp = int64(t)
		} else {
			timestamp = time.Now().Unix()
		}

		logrus.WithField("what", what).Debug("Handler invoked: grafana_create_graphite_annotation")

		annotation, err := c.CreateGraphiteAnnotation(ctx, what, data, timestamp, tags)
		if err != nil {
			return nil, fmt.Errorf("failed to create Graphite annotation: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"annotation": annotation,
			"message":    "Graphite annotation created successfully",
		}, "grafana_create_graphite_annotation")
	}
}

// ============ Datasource Management Handlers ============

// HandleCreateDatasource handles creating a new datasource.
func HandleCreateDatasource() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}

		datasourceType, err := requireStringParam(request, "type")
		if err != nil {
			return nil, err
		}

		url, err := requireStringParam(request, "url")
		if err != nil {
			return nil, err
		}

		req := client.CreateDatasourceRequest{
			Name: name,
			Type: datasourceType,
			URL:  url,
		}

		// Optional parameters
		if access := getOptionalStringParam(request, "access"); access != "" {
			req.Access = access
		} else {
			req.Access = "proxy" // Default access mode
		}

		if database := getOptionalStringParam(request, "database"); database != "" {
			req.Database = database
		}

		if user := getOptionalStringParam(request, "user"); user != "" {
			req.User = user
		}

		if password := getOptionalStringParam(request, "password"); password != "" {
			req.Password = password
		}

		if jsonData, ok := request.GetArguments()["jsonData"].(map[string]interface{}); ok {
			req.JSONData = jsonData
		}

		if secureJsonData, ok := request.GetArguments()["secureJsonData"].(map[string]interface{}); ok {
			req.SecureJSONData = secureJsonData
		}

		if isDefault, ok := request.GetArguments()["isDefault"].(bool); ok {
			req.IsDefault = isDefault
		}

		logrus.WithField("name", name).Debug("Handler invoked: grafana_create_datasource")

		dataSource, err := c.CreateDatasource(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to create datasource: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"datasource": dataSource,
			"message":    "Datasource created successfully",
		}, "grafana_create_datasource")
	}
}

// HandleUpdateDatasource handles updating an existing datasource.
func HandleUpdateDatasource() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		uid, err := requireStringParam(request, "uid")
		if err != nil {
			return nil, err
		}

		name, err := requireStringParam(request, "name")
		if err != nil {
			return nil, err
		}

		datasourceType, err := requireStringParam(request, "type")
		if err != nil {
			return nil, err
		}

		url, err := requireStringParam(request, "url")
		if err != nil {
			return nil, err
		}

		req := client.UpdateDatasourceRequest{
			UID:  uid,
			Name: name,
			Type: datasourceType,
			URL:  url,
		}

		// Optional parameters
		if access := getOptionalStringParam(request, "access"); access != "" {
			req.Access = access
		}

		if database := getOptionalStringParam(request, "database"); database != "" {
			req.Database = database
		}

		if user := getOptionalStringParam(request, "user"); user != "" {
			req.User = user
		}

		if password := getOptionalStringParam(request, "password"); password != "" {
			req.Password = password
		}

		if jsonData, ok := request.GetArguments()["jsonData"].(map[string]interface{}); ok {
			req.JSONData = jsonData
		}

		if secureJsonData, ok := request.GetArguments()["secureJsonData"].(map[string]interface{}); ok {
			req.SecureJSONData = secureJsonData
		}

		if isDefault, ok := request.GetArguments()["isDefault"].(bool); ok {
			req.IsDefault = isDefault
		}

		logrus.WithField("uid", uid).Debug("Handler invoked: grafana_update_datasource")

		dataSource, err := c.UpdateDatasource(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to update datasource: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"datasource": dataSource,
			"message":    "Datasource updated successfully",
		}, "grafana_update_datasource")
	}
}

// HandleDeleteDatasource handles deleting a datasource.
func HandleDeleteDatasource() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, cerr := client.FromContext(ctx)
		if cerr != nil {
			return mcp.NewToolResultError(cerr.Error()), nil
		}

		uid, err := requireStringParam(request, "uid")
		if err != nil {
			return nil, err
		}

		logrus.WithField("uid", uid).Debug("Handler invoked: grafana_delete_datasource")

		if err := c.DeleteDatasource(ctx, uid); err != nil {
			return nil, fmt.Errorf("failed to delete datasource: %w", err)
		}

		return marshalOptimizedResponse(map[string]interface{}{
			"uid":     uid,
			"message": "Datasource deleted successfully",
		}, "grafana_delete_datasource")
	}
}

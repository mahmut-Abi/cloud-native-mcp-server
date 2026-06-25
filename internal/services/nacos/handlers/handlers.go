package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	svccommon "github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/common"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/nacos/client"
	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

// ServiceInterface is the subset of methods required by handlers.
type ServiceInterface interface {
	GetDefaultNamespaceID() string
	GetDefaultGroup() string
}

func HandleTestConnection(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nacosClient, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		result, err := nacosClient.TestConnection(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to test nacos connection: %w", err)
		}
		return marshalResult(result)
	}
}

func HandleListNamespaces(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nacosClient, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		result, err := nacosClient.ListNamespaces(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list nacos namespaces: %w", err)
		}
		return marshalResult(map[string]interface{}{
			"count":      len(result),
			"namespaces": result,
		})
	}
}

func HandleListConfigsSummary(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		params := url.Values{}
		params.Set("search", "blur")
		params.Set("pageNo", strconv.Itoa(getInt(args, 1, "page")))
		params.Set("pageSize", strconv.Itoa(getInt(args, 20, "limit")))

		if namespaceID := resolveNamespaceID(args, service); namespaceID != "" {
			params.Set("tenant", namespaceID)
		}
		if group := resolveConfigGroup(args, service); group != "" {
			params.Set("group", group)
		}
		if query, ok := svccommon.GetStringArg(args, "query", "data_id"); ok {
			params.Set("dataId", query)
		}

		nacosClient, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		result, err := nacosClient.ListConfigs(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list nacos configs: %w", err)
		}

		items, _ := mapSlice(result["pageItems"])
		summaries := make([]map[string]interface{}, 0, len(items))
		for _, item := range items {
			summaries = append(summaries, compactConfig(item))
		}

		return marshalResult(map[string]interface{}{
			"count":      len(summaries),
			"page":       result["pageNumber"],
			"totalCount": result["totalCount"],
			"data":       summaries,
		})
	}
}

func HandleGetConfig(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		dataID, err := svccommon.RequireStringArg(args, "data_id")
		if err != nil {
			return nil, err
		}

		params := url.Values{}
		params.Set("dataId", dataID)
		if group := resolveConfigGroup(args, service); group != "" {
			params.Set("group", group)
		}
		if namespaceID := resolveNamespaceID(args, service); namespaceID != "" {
			params.Set("tenant", namespaceID)
		}

		nacosClient, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		result, err := nacosClient.GetConfig(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to get nacos config: %w", err)
		}
		return marshalResult(result)
	}
}

func HandleListServicesSummary(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		params := url.Values{}
		params.Set("pageNo", strconv.Itoa(getInt(args, 1, "page")))
		params.Set("pageSize", strconv.Itoa(getInt(args, 20, "limit")))
		if namespaceID := resolveNamespaceID(args, service); namespaceID != "" {
			params.Set("namespaceId", namespaceID)
		}
		if groupName := resolveNamingGroup(args, service); groupName != "" {
			params.Set("groupName", groupName)
		}

		nacosClient, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		result, err := nacosClient.ListServices(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list nacos services: %w", err)
		}

		doms, _ := stringSlice(result["doms"])
		filter := strings.TrimSpace(firstString(args, "service_name", "query"))
		summaries := make([]map[string]interface{}, 0, len(doms))
		for _, item := range doms {
			if filter != "" && !strings.Contains(item, filter) {
				continue
			}
			summaries = append(summaries, map[string]interface{}{
				"serviceName": item,
				"groupName":   resolveNamingGroup(args, service),
				"namespaceId": resolveNamespaceID(args, service),
			})
		}

		return marshalResult(map[string]interface{}{
			"count":      len(summaries),
			"totalCount": result["count"],
			"data":       summaries,
		})
	}
}

func HandleGetService(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		serviceName, err := svccommon.RequireStringArg(args, "service_name")
		if err != nil {
			return nil, err
		}

		params := url.Values{}
		params.Set("serviceName", serviceName)
		if namespaceID := resolveNamespaceID(args, service); namespaceID != "" {
			params.Set("namespaceId", namespaceID)
		}
		if groupName := resolveNamingGroup(args, service); groupName != "" {
			params.Set("groupName", groupName)
		}

		nacosClient, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		result, err := nacosClient.GetService(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to get nacos service: %w", err)
		}
		return marshalResult(result)
	}
}

func HandleListInstances(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		serviceName, err := svccommon.RequireStringArg(args, "service_name")
		if err != nil {
			return nil, err
		}

		params := url.Values{}
		params.Set("serviceName", serviceName)
		if namespaceID := resolveNamespaceID(args, service); namespaceID != "" {
			params.Set("namespaceId", namespaceID)
		}
		if groupName := resolveNamingGroup(args, service); groupName != "" {
			params.Set("groupName", groupName)
		}
		if clusterName, ok := svccommon.GetStringArg(args, "cluster_name"); ok {
			params.Set("clusterName", clusterName)
		}
		if healthyOnly, ok := getBool(args, "healthy_only"); ok {
			params.Set("healthyOnly", strconv.FormatBool(healthyOnly))
		}

		nacosClient, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		result, err := nacosClient.ListInstances(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list nacos instances: %w", err)
		}
		return marshalResult(result)
	}
}

func HandleListClusterNodes(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nacosClient, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		result, err := nacosClient.ListClusterNodes(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list nacos cluster nodes: %w", err)
		}
		return marshalResult(result)
	}
}

func HandleGetSystemMetrics(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nacosClient, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		result, err := nacosClient.GetSystemMetrics(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get nacos system metrics: %w", err)
		}
		return marshalResult(result)
	}
}

func marshalResult(data interface{}) (*mcp.CallToolResult, error) {
	body, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}
	return mcp.NewToolResultText(string(body)), nil
}

func resolveNamespaceID(args map[string]interface{}, service ServiceInterface) string {
	if value, ok := svccommon.GetStringArg(args, "namespace_id", "tenant"); ok {
		return value
	}
	return service.GetDefaultNamespaceID()
}

func resolveConfigGroup(args map[string]interface{}, service ServiceInterface) string {
	if value, ok := svccommon.GetStringArg(args, "group"); ok {
		return value
	}
	if group := service.GetDefaultGroup(); group != "" {
		return group
	}
	return "DEFAULT_GROUP"
}

func resolveNamingGroup(args map[string]interface{}, service ServiceInterface) string {
	if value, ok := svccommon.GetStringArg(args, "group_name"); ok {
		return value
	}
	return service.GetDefaultGroup()
}

func compactConfig(item map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for _, key := range []string{"id", "dataId", "group", "tenant", "appName", "type", "lastModifiedTime"} {
		if value, ok := item[key]; ok {
			result[key] = value
		}
	}
	return result
}

func getInt(args map[string]interface{}, defaultValue int, keys ...string) int {
	return svccommon.GetIntArg(args, defaultValue, keys...)
}

func getBool(args map[string]interface{}, keys ...string) (bool, bool) {
	value, ok := svccommon.LookupArg(args, keys...)
	if !ok {
		return false, false
	}
	switch typed := value.(type) {
	case bool:
		return typed, true
	case string:
		lower := strings.TrimSpace(strings.ToLower(typed))
		switch lower {
		case "true", "1", "yes", "on":
			return true, true
		case "false", "0", "no", "off":
			return false, true
		}
	}
	return false, false
}

func firstString(args map[string]interface{}, keys ...string) string {
	if value, ok := svccommon.GetStringArg(args, keys...); ok {
		return value
	}
	return ""
}

func mapSlice(value interface{}) ([]map[string]interface{}, bool) {
	items, ok := value.([]interface{})
	if !ok {
		return nil, false
	}
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		if typed, ok := item.(map[string]interface{}); ok {
			result = append(result, typed)
		}
	}
	return result, true
}

func stringSlice(value interface{}) ([]string, bool) {
	items, ok := value.([]interface{})
	if !ok {
		return nil, false
	}
	result := make([]string, 0, len(items))
	for _, item := range items {
		if typed, ok := item.(string); ok {
			result = append(result, typed)
		}
	}
	return result, true
}

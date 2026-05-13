package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/argocd/client"
	svccommon "github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/common"
	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

// ServiceInterface is the subset of methods required by handlers.
type ServiceInterface interface {
	GetClient() *client.Client
}

func HandleTestConnection(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		argocdClient, err := getClient(service)
		if err != nil {
			return nil, err
		}
		result, err := argocdClient.TestConnection(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to test argocd connection: %w", err)
		}
		return marshalResult(result)
	}
}

func HandleListApplicationsSummary(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		params := url.Values{}
		if project, ok := svccommon.GetStringArg(args, "project"); ok {
			params.Set("project", project)
		}
		if selector, ok := svccommon.GetStringArg(args, "selector"); ok {
			params.Set("selector", selector)
		}
		if repo, ok := svccommon.GetStringArg(args, "repo"); ok {
			params.Set("repo", repo)
		}

		argocdClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := argocdClient.ListApplications(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list argocd applications: %w", err)
		}
		items, _ := client.ItemsSlice(result)
		summaries := make([]map[string]interface{}, 0, len(items))
		for _, item := range items {
			summaries = append(summaries, compactApplication(item))
		}
		return marshalResult(map[string]interface{}{
			"count": len(summaries),
			"data":  summaries,
		})
	}
}

func HandleGetApplication(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		name, err := svccommon.RequireStringArg(args, "name")
		if err != nil {
			return nil, err
		}

		params := url.Values{}
		if appNamespace, ok := svccommon.GetStringArg(args, "app_namespace"); ok {
			params.Set("appNamespace", appNamespace)
		}

		argocdClient, err := getClient(service)
		if err != nil {
			return nil, err
		}
		result, err := argocdClient.GetApplication(ctx, name, params)
		if err != nil {
			return nil, fmt.Errorf("failed to get argocd application: %w", err)
		}
		return marshalResult(result)
	}
}

func HandleGetApplicationManifests(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		name, err := svccommon.RequireStringArg(args, "name")
		if err != nil {
			return nil, err
		}

		params := url.Values{}
		if appNamespace, ok := svccommon.GetStringArg(args, "app_namespace"); ok {
			params.Set("appNamespace", appNamespace)
		}
		if revision, ok := svccommon.GetStringArg(args, "revision"); ok {
			params.Set("revision", revision)
		}
		if namespace, ok := svccommon.GetStringArg(args, "namespace"); ok {
			params.Set("namespace", namespace)
		}

		argocdClient, err := getClient(service)
		if err != nil {
			return nil, err
		}
		result, err := argocdClient.GetApplicationManifests(ctx, name, params)
		if err != nil {
			return nil, fmt.Errorf("failed to get argocd application manifests: %w", err)
		}
		return marshalResult(result)
	}
}

func HandleListProjects(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		argocdClient, err := getClient(service)
		if err != nil {
			return nil, err
		}
		result, err := argocdClient.ListProjects(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list argocd projects: %w", err)
		}
		items, _ := client.ItemsSlice(result)
		return marshalResult(map[string]interface{}{
			"count": len(items),
			"data":  items,
		})
	}
}

func HandleGetProject(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := svccommon.RequireStringArg(request.GetArguments(), "name")
		if err != nil {
			return nil, err
		}
		argocdClient, err := getClient(service)
		if err != nil {
			return nil, err
		}
		result, err := argocdClient.GetProject(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("failed to get argocd project: %w", err)
		}
		return marshalResult(result)
	}
}

func HandleListClusters(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		argocdClient, err := getClient(service)
		if err != nil {
			return nil, err
		}
		result, err := argocdClient.ListClusters(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list argocd clusters: %w", err)
		}
		items, _ := client.ItemsSlice(result)
		return marshalResult(map[string]interface{}{
			"count": len(items),
			"data":  items,
		})
	}
}

func marshalResult(data interface{}) (*mcp.CallToolResult, error) {
	body, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}
	return mcp.NewToolResultText(string(body)), nil
}

func getClient(service ServiceInterface) (*client.Client, error) {
	argocdClient := service.GetClient()
	if argocdClient == nil {
		return nil, fmt.Errorf("argocd client is not initialized")
	}
	return argocdClient, nil
}

func compactApplication(item map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	if metadata, ok := item["metadata"].(map[string]interface{}); ok {
		for _, key := range []string{"name", "namespace"} {
			if value, ok := metadata[key]; ok {
				result[key] = value
			}
		}
	}
	if spec, ok := item["spec"].(map[string]interface{}); ok {
		for _, key := range []string{"project"} {
			if value, ok := spec[key]; ok {
				result[key] = value
			}
		}
		if destination, ok := spec["destination"].(map[string]interface{}); ok {
			if value, ok := destination["namespace"]; ok {
				result["destinationNamespace"] = value
			}
			if value, ok := destination["server"]; ok {
				result["destinationServer"] = value
			}
		}
		if source, ok := spec["source"].(map[string]interface{}); ok {
			if value, ok := source["repoURL"]; ok {
				result["repoURL"] = value
			}
			if value, ok := source["targetRevision"]; ok {
				result["targetRevision"] = value
			}
			if value, ok := source["path"]; ok {
				result["path"] = value
			}
		}
	}
	if status, ok := item["status"].(map[string]interface{}); ok {
		if sync, ok := status["sync"].(map[string]interface{}); ok {
			if value, ok := sync["status"]; ok {
				result["syncStatus"] = value
			}
		}
		if health, ok := status["health"].(map[string]interface{}); ok {
			if value, ok := health["status"]; ok {
				result["healthStatus"] = value
			}
		}
	}
	return result
}

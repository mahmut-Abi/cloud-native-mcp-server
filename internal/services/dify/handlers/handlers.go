package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/dify/client"
	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

type ServiceInterface interface{}

// --- Service API handlers ---

func HandleDifyAppInfo(service ServiceInterface) server.ToolHandlerFunc {
	return serviceGetHandler("GET", "/v1/info", "dify app info")
}

func HandleDifyAppMeta(service ServiceInterface) server.ToolHandlerFunc {
	return serviceGetHandler("GET", "/v1/meta", "dify app meta")
}

func HandleDifyAppParameters(service ServiceInterface) server.ToolHandlerFunc {
	return serviceGetHandler("GET", "/v1/parameters", "dify app parameters")
}

func HandleDifyAppSite(service ServiceInterface) server.ToolHandlerFunc {
	return serviceGetHandler("GET", "/v1/site", "dify app site")
}

func HandleDifyApiRequest(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		args := request.GetArguments()
		path := getStringArg(args, "path")
		method := getStringArg(args, "method")
		if method == "" {
			method = "GET"
		}
		q := parseQueryArgs(args["query"])
		var body interface{}
		if b, ok := args["body"]; ok {
			body = b
		}
		if isConsolePath(path) {
			resp, err := c.ConsoleRequest(ctx, method, path, q, body)
			if err != nil {
				return nil, fmt.Errorf("dify api request: %w", err)
			}
			return mcp.NewToolResultText(string(resp)), nil
		}
		resp, err := c.ServiceRequest(ctx, method, path, q, body)
		if err != nil {
			return nil, fmt.Errorf("dify api request: %w", err)
		}
		return mcp.NewToolResultText(string(resp)), nil
	}
}

func HandleDifyChatMessage(service ServiceInterface) server.ToolHandlerFunc {
	return servicePostHandler("/v1/chat-messages", "dify chat message")
}

func HandleDifyCompletionMessage(service ServiceInterface) server.ToolHandlerFunc {
	return servicePostHandler("/v1/completion-messages", "dify completion message")
}

func HandleDifyListConversations(service ServiceInterface) server.ToolHandlerFunc {
	return serviceGetHandler("GET", "/v1/conversations", "dify list conversations")
}

func HandleDifyListMessages(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		args := request.GetArguments()
		convID := getStringArg(args, "conversation_id")
		query := url.Values{}
		if user := getStringArg(args, "user"); user != "" {
			query.Set("user", user)
		}
		if firstID := getStringArg(args, "first_id"); firstID != "" {
			query.Set("first_id", firstID)
		}
		if limit, ok := args["limit"]; ok {
			query.Set("limit", fmt.Sprint(limit))
		}
		resp, err := c.ServiceRequest(ctx, "GET", "/v1/messages?conversation_id="+convID, query, nil)
		if err != nil {
			return nil, fmt.Errorf("dify list messages: %w", err)
		}
		return mcp.NewToolResultText(string(resp)), nil
	}
}

func HandleDifyRetrieveDataset(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		args := request.GetArguments()
		datasetID := getStringArg(args, "dataset_id")
		query := getStringArg(args, "query")
		body := map[string]interface{}{"query": query}
		if rm, ok := args["retrieval_model"]; ok {
			body["retrieval_model"] = rm
		}
		bodyJSON, _ := json.Marshal(body)
		var bodyMap interface{}
		json.Unmarshal(bodyJSON, &bodyMap)
		resp, err := c.ServiceRequest(ctx, "POST", "/v1/datasets/"+datasetID+"/retrieve", nil, bodyMap)
		if err != nil {
			return nil, fmt.Errorf("dify retrieve dataset: %w", err)
		}
		return mcp.NewToolResultText(string(resp)), nil
	}
}

func HandleDifyRunWorkflow(service ServiceInterface) server.ToolHandlerFunc {
	return servicePostHandler("/v1/workflows/run", "dify run workflow")
}

// --- Console API handlers ---

func HandleDifyConsoleApiRequest(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleHandler("")
}

func HandleDifyListApps(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		args := request.GetArguments()
		query := url.Values{}
		if page, ok := args["page"]; ok {
			query.Set("page", fmt.Sprint(page))
		}
		if limit, ok := args["limit"]; ok {
			query.Set("limit", fmt.Sprint(limit))
		}
		if mode := getStringArg(args, "mode"); mode != "" {
			query.Set("mode", mode)
		}
		if name := getStringArg(args, "name"); name != "" {
			query.Set("name", name)
		}
		if isMe, ok := args["is_created_by_me"]; ok {
			query.Set("is_created_by_me", fmt.Sprint(isMe))
		}
		resp, err := c.ConsoleRequest(ctx, "GET", "/api/apps", query, nil)
		if err != nil {
			return nil, fmt.Errorf("dify list apps: %w", err)
		}
		return mcp.NewToolResultText(string(resp)), nil
	}
}

func HandleDifyGetApp(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDHandler("/api/apps", "app_id", "dify get app")
}

func HandleDifyCreateApp(service ServiceInterface) server.ToolHandlerFunc {
	return createConsolePostHandler("/api/apps", "dify create app")
}

func HandleDifySetAppApiStatus(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDPostHandler("/api/apps", "app_id", "/api/enable", "dify set app api status")
}

func HandleDifyListAppApiKeys(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDHandlerFlat("/api/apps", "app_id", "/api-keys", "dify list app api keys")
}

func HandleDifyCreateAppApiKey(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDPostFlat("/api/apps", "app_id", "/api-keys", "dify create app api key")
}

func HandleDifyGetAppTraceStatus(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDHandlerFlat("/api/apps", "app_id", "/trace", "dify get app trace status")
}

func HandleDifySetAppTraceStatus(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDPostHandler("/api/apps", "app_id", "/trace", "dify set app trace status")
}

func HandleDifyGetAppTraceConfig(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDQueryHandler("/api/apps", "app_id", "/trace-config", "dify get app trace config")
}

func HandleDifyCreateAppTraceConfig(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDPostHandler("/api/apps", "app_id", "/trace-config", "dify create app trace config")
}

func HandleDifyUpdateAppTraceConfig(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDPostHandler("/api/apps", "app_id", "/trace-config", "dify update app trace config")
}

func HandleDifyDeleteAppTraceConfig(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDDeleteHandler("/api/apps", "app_id", "/trace-config", "dify delete app trace config")
}

func HandleDifyGetDraftWorkflow(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDHandlerFlat("/api/apps", "app_id", "/workflows/draft", "dify get draft workflow")
}

func HandleDifySyncDraftWorkflow(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDPostHandler("/api/apps", "app_id", "/workflows/draft", "dify sync draft workflow")
}

func HandleDifyGetDraftWorkflowEnvVars(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDHandlerFlat("/api/apps", "app_id", "/workflows/draft/environment-variables", "dify get draft workflow env vars")
}

func HandleDifyUpdateDraftWorkflowEnvVars(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDPostHandler("/api/apps", "app_id", "/workflows/draft/environment-variables", "dify update draft workflow env vars")
}

func HandleDifyGetDraftWorkflowConvVars(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDHandlerFlat("/api/apps", "app_id", "/workflows/draft/conversation-variables", "dify get draft workflow conv vars")
}

func HandleDifyUpdateDraftWorkflowConvVars(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDPostHandler("/api/apps", "app_id", "/workflows/draft/conversation-variables", "dify update draft workflow conv vars")
}

func HandleDifyGetPublishedWorkflow(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDHandlerFlat("/api/apps", "app_id", "/workflows/publish", "dify get published workflow")
}

func HandleDifyPublishWorkflow(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDPostHandler("/api/apps", "app_id", "/workflows/publish", "dify publish workflow")
}

func HandleDifyListPublishedWorkflows(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDListHandler("/api/apps", "app_id", "/workflows/publish", "dify list published workflows")
}

func HandleDifyRestoreWorkflowToDraft(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDPostHandler("/api/apps", "app_id", "/workflows/restore", "dify restore workflow to draft")
}

func HandleDifyRunDraftWorkflow(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDPostHandler("/api/apps", "app_id", "/workflows/run", "dify run draft workflow")
}

func HandleDifyRunAdvancedChatDraftWorkflow(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDPostHandler("/api/apps", "app_id", "/advanced-chat/workflows/draft/run", "dify run advanced chat draft workflow")
}

func HandleDifyListDatasets(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleListHandler("/api/datasets", "dify list datasets")
}

func HandleDifyGetDataset(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDHandler("/api/datasets", "dataset_id", "dify get dataset")
}

func HandleDifyCreateDataset(service ServiceInterface) server.ToolHandlerFunc {
	return createConsolePostHandler("/api/datasets", "dify create dataset")
}

func HandleDifyDeleteDataset(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDDeleteHandler("/api/datasets", "dataset_id", "", "dify delete dataset")
}

func HandleDifyUpdateDataset(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		args := request.GetArguments()
		datasetID := getStringArg(args, "dataset_id")
		var body interface{}
		if update, ok := args["update"]; ok {
			body = update
		}
		resp, err := c.ConsoleRequest(ctx, "PATCH", "/api/datasets/"+datasetID, nil, body)
		if err != nil {
			return nil, fmt.Errorf("dify update dataset: %w", err)
		}
		return mcp.NewToolResultText(string(resp)), nil
	}
}

func HandleDifyListDatasetDocuments(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDListHandler("/api/datasets", "dataset_id", "/documents", "dify list dataset documents")
}

func HandleDifyGetDatasetDocument(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		args := request.GetArguments()
		datasetID := getStringArg(args, "dataset_id")
		docID := getStringArg(args, "document_id")
		query := url.Values{}
		if meta := getStringArg(args, "metadata"); meta != "" {
			query.Set("metadata", meta)
		}
		resp, err := c.ConsoleRequest(ctx, "GET", "/api/datasets/"+datasetID+"/documents/"+docID, query, nil)
		if err != nil {
			return nil, fmt.Errorf("dify get dataset document: %w", err)
		}
		return mcp.NewToolResultText(string(resp)), nil
	}
}

func HandleDifySetDatasetApiStatus(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleIDPostHandler("/api/datasets", "dataset_id", "/api-status", "dify set dataset api status")
}

func HandleDifyListDatasetApiKeys(service ServiceInterface) server.ToolHandlerFunc {
	return createConsoleFlatHandler("/api/datasets/api-keys", "dify list dataset api keys")
}

func HandleDifyCreateDatasetApiKey(service ServiceInterface) server.ToolHandlerFunc {
	return createConsolePostHandler("/api/datasets/api-keys", "dify create dataset api key")
}

// --- Helper functions ---

func serviceGetHandler(method, path, label string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		resp, err := c.ServiceRequest(ctx, method, path, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", label, err)
		}
		return mcp.NewToolResultText(string(resp)), nil
	}
}

func servicePostHandler(path, label string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		args := request.GetArguments()
		body := buildServiceBody(args)
		bodyJSON, _ := json.Marshal(body)
		var bodyMap interface{}
		json.Unmarshal(bodyJSON, &bodyMap)
		resp, err := c.ServiceRequest(ctx, "POST", path, nil, bodyMap)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", label, err)
		}
		return mcp.NewToolResultText(string(resp)), nil
	}
}

func buildServiceBody(args map[string]interface{}) map[string]interface{} {
	body := map[string]interface{}{"response_mode": "blocking"}
	if query := getStringArg(args, "query"); query != "" {
		body["query"] = query
	}
	if inputs, ok := args["inputs"]; ok {
		body["inputs"] = inputs
	}
	if user := getStringArg(args, "user"); user != "" {
		body["user"] = user
	}
	if convID := getStringArg(args, "conversation_id"); convID != "" {
		body["conversation_id"] = convID
	}
	return body
}

func createConsoleHandler(idField string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		args := request.GetArguments()
		path := getStringArg(args, "path")
		method := getStringArg(args, "method")
		if method == "" {
			method = "GET"
		}
		query := parseQueryArgs(args["query"])
		var body interface{}
		if b, ok := args["body"]; ok {
			body = b
		}
		resp, err := c.ConsoleRequest(ctx, method, path, query, body)
		if err != nil {
			return nil, fmt.Errorf("dify console request: %w", err)
		}
		return mcp.NewToolResultText(string(resp)), nil
	}
}

func createConsoleIDHandler(basePath, idField, label string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		args := request.GetArguments()
		id := getStringArg(args, idField)
		resp, err := c.ConsoleRequest(ctx, "GET", basePath+"/"+id, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", label, err)
		}
		return mcp.NewToolResultText(string(resp)), nil
	}
}

func createConsoleIDHandlerFlat(basePath, idField, suffix, label string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		args := request.GetArguments()
		id := getStringArg(args, idField)
		resp, err := c.ConsoleRequest(ctx, "GET", basePath+"/"+id+suffix, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", label, err)
		}
		return mcp.NewToolResultText(string(resp)), nil
	}
}

func createConsoleIDQueryHandler(basePath, idField, suffix, label string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		args := request.GetArguments()
		id := getStringArg(args, idField)
		query := parseQueryArgs(args["query"])
		resp, err := c.ConsoleRequest(ctx, "GET", basePath+"/"+id+suffix, query, nil)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", label, err)
		}
		return mcp.NewToolResultText(string(resp)), nil
	}
}

func createConsolePostHandler(basePath, label string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		args := request.GetArguments()
		body := extractBody(args, "body")
		resp, err := c.ConsoleRequest(ctx, "POST", basePath, nil, body)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", label, err)
		}
		return mcp.NewToolResultText(string(resp)), nil
	}
}

func createConsoleIDPostHandler(basePath, idField, suffix, label string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		args := request.GetArguments()
		id := getStringArg(args, idField)
		path := basePath + "/" + id + suffix
		body := extractBody(args, "body")
		resp, err := c.ConsoleRequest(ctx, "POST", path, nil, body)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", label, err)
		}
		return mcp.NewToolResultText(string(resp)), nil
	}
}

func createConsoleIDPostFlat(basePath, idField, suffix, label string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		args := request.GetArguments()
		id := getStringArg(args, idField)
		resp, err := c.ConsoleRequest(ctx, "POST", basePath+"/"+id+suffix, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", label, err)
		}
		return mcp.NewToolResultText(string(resp)), nil
	}
}

func createConsoleIDDeleteHandler(basePath, idField, suffix, label string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		args := request.GetArguments()
		id := getStringArg(args, idField)
		path := basePath + "/" + id
		if suffix != "" {
			path += suffix
		}
		resp, err := c.ConsoleRequest(ctx, "DELETE", path, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", label, err)
		}
		return mcp.NewToolResultText(string(resp)), nil
	}
}

func createConsoleListHandler(basePath, label string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		args := request.GetArguments()
		query := parseQueryArgs(args["query"])
		resp, err := c.ConsoleRequest(ctx, "GET", basePath, query, nil)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", label, err)
		}
		return mcp.NewToolResultText(string(resp)), nil
	}
}

func createConsoleIDListHandler(basePath, idField, suffix, label string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		args := request.GetArguments()
		id := getStringArg(args, idField)
		query := parseQueryArgs(args["query"])
		path := basePath + "/" + id
		if suffix != "" {
			path += suffix
		}
		resp, err := c.ConsoleRequest(ctx, "GET", path, query, nil)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", label, err)
		}
		return mcp.NewToolResultText(string(resp)), nil
	}
}

func createConsoleFlatHandler(basePath, label string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		c, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		resp, err := c.ConsoleRequest(ctx, "GET", basePath, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", label, err)
		}
		return mcp.NewToolResultText(string(resp)), nil
	}
}

func parseQueryArgs(raw interface{}) url.Values {
	qv := url.Values{}
	if raw == nil {
		return qv
	}
	if m, ok := raw.(map[string]interface{}); ok {
		for k, v := range m {
			qv.Set(k, fmt.Sprint(v))
		}
	}
	return qv
}

func extractBody(args map[string]interface{}, key string) interface{} {
	if b, ok := args[key]; ok {
		return b
	}
	return nil
}

func getStringArg(args map[string]interface{}, key string) string {
	if v, ok := args[key]; ok {
		return fmt.Sprint(v)
	}
	return ""
}

func isConsolePath(path string) bool {
	p := path
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	p = strings.Replace(p, "/console/api", "", 1)
	switch {
	case p == "/apps" || strings.HasPrefix(p, "/apps/"):
		return true
	case p == "/datasets" || strings.HasPrefix(p, "/datasets/"):
		return !strings.HasSuffix(p, "/retrieve")
	case strings.HasPrefix(p, "/rag/"):
		return true
	default:
		return false
	}
}

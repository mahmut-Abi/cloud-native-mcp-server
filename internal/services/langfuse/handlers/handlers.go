package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	svccommon "github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/common"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/langfuse/client"
	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

// ServiceInterface is the subset of service methods required by handlers.
type ServiceInterface interface {
	GetClient() *client.Client
}

// HandleCheckHealth handles Langfuse health checks.
func HandleCheckHealth(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.CheckHealth(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to check langfuse health: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleListTracesSummary returns compact trace summaries.
func HandleListTracesSummary(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params, err := buildTraceListParams(request.GetArguments())
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.ListTraces(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list langfuse traces: %w", err)
		}

		summaries := make([]map[string]interface{}, 0)
		if rawItems, ok := result["data"].([]interface{}); ok {
			for _, rawItem := range rawItems {
				item, ok := rawItem.(map[string]interface{})
				if !ok {
					continue
				}
				summaries = append(summaries, compactTrace(item))
			}
		}

		return marshalResult(map[string]interface{}{
			"data": summaries,
			"meta": result["meta"],
		})
	}
}

// HandleListTraces handles Langfuse trace listing.
func HandleListTraces(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params, err := buildTraceListParams(request.GetArguments())
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.ListTraces(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list langfuse traces: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleGetTrace handles Langfuse trace detail lookups.
func HandleGetTrace(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		traceID, err := svccommon.RequireStringArg(args, "trace_id")
		if err != nil {
			return nil, err
		}

		fields, _ := svccommon.GetStringArg(args, "fields")

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.GetTrace(ctx, traceID, fields)
		if err != nil {
			return nil, fmt.Errorf("failed to get langfuse trace: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleListSessions handles Langfuse session listing.
func HandleListSessions(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params, err := buildSessionListParams(request.GetArguments())
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.ListSessions(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list langfuse sessions: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleGetSession handles Langfuse session detail lookups.
func HandleGetSession(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sessionID, err := svccommon.RequireStringArg(request.GetArguments(), "session_id")
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.GetSession(ctx, sessionID)
		if err != nil {
			return nil, fmt.Errorf("failed to get langfuse session: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleListObservations handles Langfuse observation listing.
func HandleListObservations(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params, err := buildObservationListParams(request.GetArguments())
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.ListObservations(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list langfuse observations: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleGetObservation handles Langfuse observation detail lookups.
func HandleGetObservation(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		observationID, err := svccommon.RequireStringArg(request.GetArguments(), "observation_id")
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.GetObservation(ctx, observationID)
		if err != nil {
			return nil, fmt.Errorf("failed to get langfuse observation: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleListPrompts handles Langfuse prompt metadata listing.
func HandleListPrompts(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params, err := buildPromptListParams(request.GetArguments())
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.ListPrompts(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list langfuse prompts: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleListAnnotationQueues handles annotation queue listing.
func HandleListAnnotationQueues(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params, err := buildPaginationParams(request.GetArguments())
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.ListAnnotationQueues(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list langfuse annotation queues: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleGetAnnotationQueue handles annotation queue detail lookups.
func HandleGetAnnotationQueue(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		queueID, err := svccommon.RequireStringArg(request.GetArguments(), "queue_id")
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.GetAnnotationQueue(ctx, queueID)
		if err != nil {
			return nil, fmt.Errorf("failed to get langfuse annotation queue: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleListAnnotationQueueItems handles queue item listing.
func HandleListAnnotationQueueItems(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		queueID, err := svccommon.RequireStringArg(args, "queue_id")
		if err != nil {
			return nil, err
		}
		params, err := buildPaginationParams(args)
		if err != nil {
			return nil, err
		}
		addStringParam(params, args, "status", "status")

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.ListAnnotationQueueItems(ctx, queueID, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list langfuse annotation queue items: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleListDatasets handles dataset listing.
func HandleListDatasets(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params, err := buildPaginationParams(request.GetArguments())
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.ListDatasets(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list langfuse datasets: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleGetDataset handles dataset detail lookups.
func HandleGetDataset(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		datasetName, err := svccommon.RequireStringArg(request.GetArguments(), "dataset_name")
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.GetDataset(ctx, datasetName)
		if err != nil {
			return nil, fmt.Errorf("failed to get langfuse dataset: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleListDatasetRuns handles dataset run listing.
func HandleListDatasetRuns(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		datasetName, err := svccommon.RequireStringArg(args, "dataset_name")
		if err != nil {
			return nil, err
		}
		params, err := buildPaginationParams(args)
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.ListDatasetRuns(ctx, datasetName, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list langfuse dataset runs: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleGetDatasetRun handles dataset run detail lookups.
func HandleGetDatasetRun(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		datasetName, err := svccommon.RequireStringArg(args, "dataset_name")
		if err != nil {
			return nil, err
		}
		runName, err := svccommon.RequireStringArg(args, "run_name")
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.GetDatasetRun(ctx, datasetName, runName)
		if err != nil {
			return nil, fmt.Errorf("failed to get langfuse dataset run: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleListLLMConnections handles LLM connection listing.
func HandleListLLMConnections(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params, err := buildPaginationParams(request.GetArguments())
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.ListLLMConnections(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list langfuse llm connections: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleListModels handles model listing.
func HandleListModels(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params, err := buildPaginationParams(request.GetArguments())
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.ListModels(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list langfuse models: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleGetModel handles model detail lookups.
func HandleGetModel(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		modelID, err := svccommon.RequireStringArg(request.GetArguments(), "model_id")
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.GetModel(ctx, modelID)
		if err != nil {
			return nil, fmt.Errorf("failed to get langfuse model: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleListScoreConfigs handles score configuration listing.
func HandleListScoreConfigs(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params, err := buildPaginationParams(request.GetArguments())
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.ListScoreConfigs(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list langfuse score configs: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleGetScoreConfig handles score configuration detail lookups.
func HandleGetScoreConfig(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		configID, err := svccommon.RequireStringArg(request.GetArguments(), "config_id")
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.GetScoreConfig(ctx, configID)
		if err != nil {
			return nil, fmt.Errorf("failed to get langfuse score config: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleGetPrompt handles Langfuse prompt detail lookups.
func HandleGetPrompt(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		promptName, err := svccommon.RequireStringArg(args, "prompt_name")
		if err != nil {
			return nil, err
		}

		params := url.Values{}
		if version := normalizePositiveInt(svccommon.GetIntArg(args, 0, "version")); version > 0 {
			params.Set("version", fmt.Sprintf("%d", version))
		}
		if label, ok := svccommon.GetStringArg(args, "label"); ok {
			params.Set("label", label)
		}
		resolve, err := svccommon.GetBoolArg(args, "resolve")
		if err != nil {
			return nil, err
		}
		if resolve != nil {
			params.Set("resolve", fmt.Sprintf("%t", *resolve))
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.GetPrompt(ctx, promptName, params)
		if err != nil {
			return nil, fmt.Errorf("failed to get langfuse prompt: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleListScores handles Langfuse score listing.
func HandleListScores(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params, err := buildScoreListParams(request.GetArguments())
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.ListScores(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list langfuse scores: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleGetScore handles Langfuse score detail lookups.
func HandleGetScore(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		scoreID, err := svccommon.RequireStringArg(request.GetArguments(), "score_id")
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.GetScore(ctx, scoreID)
		if err != nil {
			return nil, fmt.Errorf("failed to get langfuse score: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleGetMetrics handles Langfuse metrics queries.
func HandleGetMetrics(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, ok, err := svccommon.GetJSONStringArg(request.GetArguments(), "query")
		if err != nil {
			return nil, err
		}
		if !ok || query == "" {
			return nil, fmt.Errorf("missing required parameter: query")
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.GetMetrics(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("failed to query langfuse metrics: %w", err)
		}
		return marshalResult(result)
	}
}

func buildTraceListParams(args map[string]interface{}) (url.Values, error) {
	params, err := buildPaginationParams(args)
	if err != nil {
		return nil, err
	}

	addStringParam(params, args, "userId", "user_id", "userId")
	addStringParam(params, args, "name", "name")
	addStringParam(params, args, "sessionId", "session_id", "sessionId")
	addStringParam(params, args, "fromTimestamp", "from_timestamp", "fromTimestamp")
	addStringParam(params, args, "toTimestamp", "to_timestamp", "toTimestamp")
	addStringParam(params, args, "orderBy", "order_by", "orderBy")
	addStringParam(params, args, "version", "version")
	addStringParam(params, args, "release", "release")
	addStringParam(params, args, "fields", "fields")
	if err := addStringSliceParam(params, args, "tags", "tags"); err != nil {
		return nil, err
	}
	if err := addStringSliceParam(params, args, "environment", "environment"); err != nil {
		return nil, err
	}
	if err := addJSONParam(params, args, "filter", "filter"); err != nil {
		return nil, err
	}

	return params, nil
}

func buildSessionListParams(args map[string]interface{}) (url.Values, error) {
	params, err := buildPaginationParams(args)
	if err != nil {
		return nil, err
	}

	addStringParam(params, args, "fromTimestamp", "from_timestamp", "fromTimestamp")
	addStringParam(params, args, "toTimestamp", "to_timestamp", "toTimestamp")
	if err := addStringSliceParam(params, args, "environment", "environment"); err != nil {
		return nil, err
	}

	return params, nil
}

func buildObservationListParams(args map[string]interface{}) (url.Values, error) {
	params, err := buildPaginationParams(args)
	if err != nil {
		return nil, err
	}

	addStringParam(params, args, "name", "name")
	addStringParam(params, args, "userId", "user_id", "userId")
	addStringParam(params, args, "type", "type")
	addStringParam(params, args, "traceId", "trace_id", "traceId")
	addStringParam(params, args, "level", "level")
	addStringParam(params, args, "parentObservationId", "parent_observation_id", "parentObservationId")
	addStringParam(params, args, "fromStartTime", "from_start_time", "fromStartTime")
	addStringParam(params, args, "toStartTime", "to_start_time", "toStartTime")
	addStringParam(params, args, "version", "version")
	if err := addStringSliceParam(params, args, "environment", "environment"); err != nil {
		return nil, err
	}
	if err := addJSONParam(params, args, "filter", "filter"); err != nil {
		return nil, err
	}

	return params, nil
}

func buildPromptListParams(args map[string]interface{}) (url.Values, error) {
	params, err := buildPaginationParams(args)
	if err != nil {
		return nil, err
	}

	addStringParam(params, args, "name", "name")
	addStringParam(params, args, "label", "label")
	addStringParam(params, args, "tag", "tag")
	addStringParam(params, args, "fromUpdatedAt", "from_updated_at", "fromUpdatedAt")
	addStringParam(params, args, "toUpdatedAt", "to_updated_at", "toUpdatedAt")

	return params, nil
}

func buildScoreListParams(args map[string]interface{}) (url.Values, error) {
	params, err := buildPaginationParams(args)
	if err != nil {
		return nil, err
	}

	addStringParam(params, args, "userId", "user_id", "userId")
	addStringParam(params, args, "name", "name")
	addStringParam(params, args, "fromTimestamp", "from_timestamp", "fromTimestamp")
	addStringParam(params, args, "toTimestamp", "to_timestamp", "toTimestamp")
	addStringParam(params, args, "source", "source")
	addStringParam(params, args, "traceId", "trace_id", "traceId")
	addStringParam(params, args, "sessionId", "session_id", "sessionId")
	addStringParam(params, args, "observationId", "observation_id", "observationId")
	addStringParam(params, args, "configId", "config_id", "configId")
	addStringParam(params, args, "fields", "fields")
	if err := addStringSliceParam(params, args, "environment", "environment"); err != nil {
		return nil, err
	}
	if err := addStringSliceParam(params, args, "traceTags", "trace_tags", "traceTags"); err != nil {
		return nil, err
	}
	if err := addJSONParam(params, args, "filter", "filter"); err != nil {
		return nil, err
	}

	return params, nil
}

func buildPaginationParams(args map[string]interface{}) (url.Values, error) {
	params := url.Values{}

	if page := normalizePositiveInt(svccommon.GetIntArg(args, 0, "page")); page > 0 {
		params.Set("page", fmt.Sprintf("%d", page))
	}
	if limit := normalizePositiveInt(svccommon.GetIntArg(args, 0, "limit")); limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}

	return params, nil
}

func addStringParam(params url.Values, args map[string]interface{}, queryKey string, keys ...string) {
	if value, ok := svccommon.GetStringArg(args, keys...); ok {
		params.Set(queryKey, value)
	}
}

func addStringSliceParam(params url.Values, args map[string]interface{}, queryKey string, keys ...string) error {
	values, ok, err := svccommon.GetStringSliceArg(args, keys...)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	for _, value := range values {
		params.Add(queryKey, value)
	}
	return nil
}

func addJSONParam(params url.Values, args map[string]interface{}, queryKey string, keys ...string) error {
	value, ok, err := svccommon.GetJSONStringArg(args, keys...)
	if err != nil {
		return err
	}
	if !ok || value == "" {
		return nil
	}
	params.Set(queryKey, value)
	return nil
}

func compactTrace(item map[string]interface{}) map[string]interface{} {
	summary := map[string]interface{}{}
	for _, key := range []string{"id", "name", "userId", "sessionId", "timestamp", "environment", "latency", "totalCost", "public"} {
		if value, ok := item[key]; ok && value != nil {
			summary[key] = value
		}
	}
	if value, ok := item["tags"]; ok && value != nil {
		summary["tags"] = value
	}
	return summary
}

func normalizePositiveInt(value int) int {
	if value < 1 {
		return 0
	}
	return value
}

func getClient(service ServiceInterface) (*client.Client, error) {
	langfuseClient := service.GetClient()
	if langfuseClient == nil {
		return nil, fmt.Errorf("langfuse client is not initialized")
	}
	return langfuseClient, nil
}

func marshalResult(result interface{}) (*mcp.CallToolResult, error) {
	jsonResponse, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize response: %w", err)
	}
	return mcp.NewToolResultText(string(jsonResponse)), nil
}

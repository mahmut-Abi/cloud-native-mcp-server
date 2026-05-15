package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

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
		rawQuery, ok, err := svccommon.GetObjectArg(request.GetArguments(), "query")
		if err != nil {
			return nil, err
		}
		if !ok || len(rawQuery) == 0 {
			return nil, fmt.Errorf("missing required parameter: query")
		}

		normalizedQuery, err := normalizeMetricsQuery(rawQuery)
		if err != nil {
			return nil, err
		}
		queryBytes, err := json.Marshal(normalizedQuery)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize normalized metrics query: %w", err)
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.GetMetrics(ctx, string(queryBytes))
		if err != nil {
			return nil, fmt.Errorf("failed to query langfuse metrics: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleGetProject handles current project lookup.
func HandleGetProject(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.GetProject(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get langfuse project: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleListOrganizationProjects handles organization project listing.
func HandleListOrganizationProjects(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.ListOrganizationProjects(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list langfuse organization projects: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleCreateProject handles project creation.
func HandleCreateProject(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		name, err := svccommon.RequireStringArg(args, "name")
		if err != nil {
			return nil, err
		}
		metadata, hasMetadata, err := svccommon.GetObjectArg(args, "metadata")
		if err != nil {
			return nil, fmt.Errorf("invalid metadata: %w", err)
		}
		if !hasMetadata {
			metadata = nil
		}
		retentionDays := normalizeNonNegativeInt(svccommon.GetIntArg(args, 0, "retention_days", "retentionDays", "retention"))

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.CreateProject(ctx, name, metadata, retentionDays)
		if err != nil {
			return nil, fmt.Errorf("failed to create langfuse project: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleUpdateProject handles project updates.
func HandleUpdateProject(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		projectID, err := svccommon.RequireStringArg(args, "project_id", "projectId")
		if err != nil {
			return nil, err
		}
		name, err := svccommon.RequireStringArg(args, "name")
		if err != nil {
			return nil, err
		}
		metadata, hasMetadata, err := svccommon.GetObjectArg(args, "metadata")
		if err != nil {
			return nil, fmt.Errorf("invalid metadata: %w", err)
		}
		if !hasMetadata {
			metadata = nil
		}

		var retentionDays *int
		if _, ok := svccommon.LookupArg(args, "retention_days", "retentionDays", "retention"); ok {
			value := normalizeNonNegativeInt(svccommon.GetIntArg(args, 0, "retention_days", "retentionDays", "retention"))
			retentionDays = &value
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.UpdateProject(ctx, projectID, name, metadata, retentionDays)
		if err != nil {
			return nil, fmt.Errorf("failed to update langfuse project: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleDeleteProject handles project deletion.
func HandleDeleteProject(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectID, err := svccommon.RequireStringArg(request.GetArguments(), "project_id", "projectId")
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.DeleteProject(ctx, projectID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete langfuse project: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleListProjectMemberships handles project membership listing.
func HandleListProjectMemberships(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectID, err := svccommon.RequireStringArg(request.GetArguments(), "project_id", "projectId")
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.ListProjectMemberships(ctx, projectID)
		if err != nil {
			return nil, fmt.Errorf("failed to list langfuse project memberships: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleUpsertProjectMembership handles project membership creation or update.
func HandleUpsertProjectMembership(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		projectID, err := svccommon.RequireStringArg(args, "project_id", "projectId")
		if err != nil {
			return nil, err
		}
		userID, err := svccommon.RequireStringArg(args, "user_id", "userId")
		if err != nil {
			return nil, err
		}
		role, err := svccommon.RequireStringArg(args, "role")
		if err != nil {
			return nil, err
		}
		role, err = normalizeMembershipRole(role)
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.UpsertProjectMembership(ctx, projectID, userID, role)
		if err != nil {
			return nil, fmt.Errorf("failed to upsert langfuse project membership: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleDeleteProjectMembership handles project membership deletion.
func HandleDeleteProjectMembership(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		projectID, err := svccommon.RequireStringArg(args, "project_id", "projectId")
		if err != nil {
			return nil, err
		}
		userID, err := svccommon.RequireStringArg(args, "user_id", "userId")
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.DeleteProjectMembership(ctx, projectID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete langfuse project membership: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleListOrganizationAPIKeys handles organization API key listing.
func HandleListOrganizationAPIKeys(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.ListOrganizationAPIKeys(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list langfuse organization api keys: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleListProjectAPIKeys handles project API key listing.
func HandleListProjectAPIKeys(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectID, err := svccommon.RequireStringArg(request.GetArguments(), "project_id", "projectId")
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.ListProjectAPIKeys(ctx, projectID)
		if err != nil {
			return nil, fmt.Errorf("failed to list langfuse project api keys: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleCreateProjectAPIKey handles project API key creation.
func HandleCreateProjectAPIKey(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		projectID, err := svccommon.RequireStringArg(args, "project_id", "projectId")
		if err != nil {
			return nil, err
		}
		note, _ := svccommon.GetStringArg(args, "note")
		publicKey, _ := svccommon.GetStringArg(args, "public_key", "publicKey")
		secretKey, _ := svccommon.GetStringArg(args, "secret_key", "secretKey")

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.CreateProjectAPIKey(ctx, projectID, note, publicKey, secretKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create langfuse project api key: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleDeleteProjectAPIKey handles project API key deletion.
func HandleDeleteProjectAPIKey(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		projectID, err := svccommon.RequireStringArg(args, "project_id", "projectId")
		if err != nil {
			return nil, err
		}
		apiKeyID, err := svccommon.RequireStringArg(args, "api_key_id", "apiKeyId")
		if err != nil {
			return nil, err
		}

		langfuseClient, err := getClient(service)
		if err != nil {
			return nil, err
		}

		result, err := langfuseClient.DeleteProjectAPIKey(ctx, projectID, apiKeyID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete langfuse project api key: %w", err)
		}
		return marshalResult(result)
	}
}

func normalizeMetricsQuery(query map[string]interface{}) (map[string]interface{}, error) {
	normalized := cloneObject(query)
	filters, _, err := svccommon.GetObjectSliceArg(query, "filters")
	if err != nil {
		return nil, fmt.Errorf("invalid metrics query filters: %w", err)
	}
	normalizedFilters := make([]map[string]interface{}, 0, len(filters))

	for _, filter := range filters {
		current := cloneObject(filter)
		column := strings.TrimSpace(firstStringValue(current, "column", "field"))
		if column == "" {
			normalizedFilters = append(normalizedFilters, current)
			continue
		}

		if _, exists := current["column"]; !exists {
			current["column"] = column
		}
		delete(current, "field")

		if strings.EqualFold(column, "timestamp") {
			operator := strings.TrimSpace(firstStringValue(current, "operator"))
			value := strings.TrimSpace(firstStringValue(current, "value"))
			switch operator {
			case ">", ">=", "after", "gte":
				if value != "" && firstStringValue(normalized, "fromTimestamp", "from_timestamp") == "" {
					normalized["fromTimestamp"] = value
				}
				continue
			case "<", "<=", "before", "lte":
				if value != "" && firstStringValue(normalized, "toTimestamp", "to_timestamp") == "" {
					normalized["toTimestamp"] = value
				}
				continue
			}
		}

		if _, exists := current["type"]; !exists {
			current["type"] = inferMetricsFilterType(current)
		}
		normalizedFilters = append(normalizedFilters, current)
	}

	if len(normalizedFilters) > 0 {
		normalized["filters"] = normalizedFilters
	} else {
		delete(normalized, "filters")
	}

	if from := firstStringValue(normalized, "from_timestamp"); from != "" && firstStringValue(normalized, "fromTimestamp") == "" {
		normalized["fromTimestamp"] = from
	}
	if to := firstStringValue(normalized, "to_timestamp"); to != "" && firstStringValue(normalized, "toTimestamp") == "" {
		normalized["toTimestamp"] = to
	}
	delete(normalized, "from_timestamp")
	delete(normalized, "to_timestamp")

	if firstStringValue(normalized, "fromTimestamp") == "" {
		return nil, fmt.Errorf("metrics query requires fromTimestamp (or equivalent Timestamp >= filter)")
	}
	if firstStringValue(normalized, "toTimestamp") == "" {
		return nil, fmt.Errorf("metrics query requires toTimestamp (or equivalent Timestamp <= filter)")
	}

	return normalized, nil
}

func inferMetricsFilterType(filter map[string]interface{}) string {
	if _, ok := filter["key"]; ok {
		return "stringObject"
	}

	if raw, ok := filter["value"]; ok {
		switch typed := raw.(type) {
		case bool:
			return "boolean"
		case float64, float32, int, int32, int64, uint, uint32, uint64:
			return "number"
		case string:
			if _, err := strconv.ParseFloat(strings.TrimSpace(typed), 64); err == nil {
				return "number"
			}
			return "string"
		}
	}

	return "string"
}

func normalizeNonNegativeInt(value int) int {
	if value < 0 {
		return 0
	}
	return value
}

func normalizeMembershipRole(role string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(role))
	switch normalized {
	case "OWNER", "ADMIN", "MEMBER", "VIEWER":
		return normalized, nil
	default:
		return "", fmt.Errorf("invalid role %q: expected OWNER, ADMIN, MEMBER, or VIEWER", role)
	}
}

func cloneObject(input map[string]interface{}) map[string]interface{} {
	if input == nil {
		return map[string]interface{}{}
	}
	result := make(map[string]interface{}, len(input))
	for key, value := range input {
		result[key] = value
	}
	return result
}

func firstStringValue(values map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		raw, ok := values[key]
		if !ok || raw == nil {
			continue
		}
		switch typed := raw.(type) {
		case string:
			if strings.TrimSpace(typed) != "" {
				return strings.TrimSpace(typed)
			}
		default:
			value := strings.TrimSpace(fmt.Sprintf("%v", typed))
			if value != "" {
				return value
			}
		}
	}
	return ""
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

package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	svccommon "github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/common"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/sentry/client"
	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

// ServiceInterface is the subset of service methods required by handlers.
type ServiceInterface interface {
	GetDefaultOrganization() string
	GetDefaultProject() string
}

// HandleTestConnection verifies that Sentry connectivity works.
func HandleTestConnection(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		sentryClient, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		organization, err := resolveOrganization(args, service, false)
		if err != nil {
			return nil, err
		}
		project, _ := resolveProject(args, service, false)

		if organization == "" {
			orgs, pagination, err := sentryClient.ListOrganizations(ctx, url.Values{})
			if err != nil {
				return nil, fmt.Errorf("failed to list sentry organizations: %w", err)
			}
			return marshalResult(map[string]interface{}{
				"status":        "ok",
				"mode":          "organizations",
				"count":         len(orgs),
				"organizations": orgs,
				"pagination":    pagination,
			})
		}

		if project != "" {
			result, err := sentryClient.GetProject(ctx, organization, project)
			if err != nil {
				return nil, fmt.Errorf("failed to get sentry project: %w", err)
			}
			return marshalResult(map[string]interface{}{
				"status":       "ok",
				"mode":         "project",
				"organization": organization,
				"project":      result,
			})
		}

		projects, pagination, err := sentryClient.ListProjects(ctx, organization, url.Values{})
		if err != nil {
			return nil, fmt.Errorf("failed to list sentry projects: %w", err)
		}
		return marshalResult(map[string]interface{}{
			"status":       "ok",
			"mode":         "projects",
			"organization": organization,
			"count":        len(projects),
			"projects":     projects,
			"pagination":   pagination,
		})
	}
}

// HandleListOrganizations handles organization listing.
func HandleListOrganizations(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params, err := buildOrganizationListParams(request.GetArguments())
		if err != nil {
			return nil, err
		}

		sentryClient, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		result, pagination, err := sentryClient.ListOrganizations(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list sentry organizations: %w", err)
		}
		return marshalResult(map[string]interface{}{
			"data":       result,
			"count":      len(result),
			"pagination": pagination,
		})
	}
}

// HandleListProjects handles project listing.
func HandleListProjects(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		organization, err := resolveOrganization(args, service, true)
		if err != nil {
			return nil, err
		}
		params, err := buildProjectListParams(args)
		if err != nil {
			return nil, err
		}

		sentryClient, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		result, pagination, err := sentryClient.ListProjects(ctx, organization, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list sentry projects: %w", err)
		}
		return marshalResult(map[string]interface{}{
			"organization": organization,
			"data":         result,
			"count":        len(result),
			"pagination":   pagination,
		})
	}
}

// HandleGetProject handles project retrieval.
func HandleGetProject(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		organization, err := resolveOrganization(args, service, true)
		if err != nil {
			return nil, err
		}
		project, err := resolveProject(args, service, true)
		if err != nil {
			return nil, err
		}

		sentryClient, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		result, err := sentryClient.GetProject(ctx, organization, project)
		if err != nil {
			return nil, fmt.Errorf("failed to get sentry project: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleListIssuesSummary handles compact issue discovery.
func HandleListIssuesSummary(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		organization, err := resolveOrganization(args, service, true)
		if err != nil {
			return nil, err
		}
		params, err := buildIssueListParams(args)
		if err != nil {
			return nil, err
		}

		sentryClient, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		result, pagination, err := sentryClient.ListIssues(ctx, organization, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list sentry issues: %w", err)
		}

		summaries := make([]map[string]interface{}, 0, len(result))
		for _, issue := range result {
			summaries = append(summaries, compactIssue(issue))
		}

		return marshalResult(map[string]interface{}{
			"organization": organization,
			"data":         summaries,
			"count":        len(summaries),
			"pagination":   pagination,
		})
	}
}

// HandleListIssues handles full issue listing.
func HandleListIssues(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		organization, err := resolveOrganization(args, service, true)
		if err != nil {
			return nil, err
		}
		params, err := buildIssueListParams(args)
		if err != nil {
			return nil, err
		}

		sentryClient, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		result, pagination, err := sentryClient.ListIssues(ctx, organization, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list sentry issues: %w", err)
		}
		return marshalResult(map[string]interface{}{
			"organization": organization,
			"data":         result,
			"count":        len(result),
			"pagination":   pagination,
		})
	}
}

// HandleGetIssue handles issue retrieval.
func HandleGetIssue(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		issueID, err := svccommon.RequireStringArg(request.GetArguments(), "issue_id")
		if err != nil {
			return nil, err
		}

		sentryClient, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		result, err := sentryClient.GetIssue(ctx, issueID)
		if err != nil {
			return nil, fmt.Errorf("failed to get sentry issue: %w", err)
		}
		return marshalResult(result)
	}
}

// HandleListIssueEvents handles issue event listing.
func HandleListIssueEvents(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		issueID, err := svccommon.RequireStringArg(args, "issue_id")
		if err != nil {
			return nil, err
		}
		params, err := buildIssueEventListParams(args)
		if err != nil {
			return nil, err
		}

		sentryClient, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		result, pagination, err := sentryClient.ListIssueEvents(ctx, issueID, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list sentry issue events: %w", err)
		}
		return marshalResult(map[string]interface{}{
			"issue_id":   issueID,
			"data":       result,
			"count":      len(result),
			"pagination": pagination,
		})
	}
}

// HandleGetIssueEvent handles issue event retrieval.
func HandleGetIssueEvent(service ServiceInterface) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		issueID, err := svccommon.RequireStringArg(args, "issue_id")
		if err != nil {
			return nil, err
		}
		eventID, err := svccommon.RequireStringArg(args, "event_id")
		if err != nil {
			return nil, err
		}

		sentryClient, err := client.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		result, err := sentryClient.GetIssueEvent(ctx, issueID, eventID)
		if err != nil {
			return nil, fmt.Errorf("failed to get sentry issue event: %w", err)
		}
		return marshalResult(result)
	}
}

func buildOrganizationListParams(args map[string]interface{}) (url.Values, error) {
	params := url.Values{}
	addStringParam(params, args, "cursor", "cursor")
	addLimitParam(params, args)
	return params, nil
}

func buildProjectListParams(args map[string]interface{}) (url.Values, error) {
	params := url.Values{}
	addStringParam(params, args, "query", "query")
	addStringParam(params, args, "statsPeriod", "stats_period", "statsPeriod")
	addStringParam(params, args, "sortBy", "sort_by", "sortBy")
	addStringParam(params, args, "cursor", "cursor")
	addLimitParam(params, args)
	return params, nil
}

func buildIssueListParams(args map[string]interface{}) (url.Values, error) {
	params := url.Values{}
	addStringParam(params, args, "query", "query")
	addStringParam(params, args, "sort", "sort")
	addStringParam(params, args, "statsPeriod", "stats_period", "statsPeriod")
	addStringParam(params, args, "cursor", "cursor")
	addLimitParam(params, args)
	if err := addStringSliceParam(params, args, "environment", "environment"); err != nil {
		return nil, err
	}
	if err := addStringSliceParam(params, args, "project", "project_ids", "projectIds"); err != nil {
		return nil, err
	}
	return params, nil
}

func buildIssueEventListParams(args map[string]interface{}) (url.Values, error) {
	params := url.Values{}
	addStringParam(params, args, "query", "query")
	addStringParam(params, args, "sort", "sort")
	addStringParam(params, args, "cursor", "cursor")
	addLimitParam(params, args)
	if err := addStringSliceParam(params, args, "environment", "environment"); err != nil {
		return nil, err
	}
	full, err := svccommon.GetBoolArg(args, "full")
	if err != nil {
		return nil, err
	}
	if full != nil {
		params.Set("full", strconv.FormatBool(*full))
	}
	sample, err := svccommon.GetBoolArg(args, "sample")
	if err != nil {
		return nil, err
	}
	if sample != nil {
		params.Set("sample", strconv.FormatBool(*sample))
	}
	return params, nil
}

func addLimitParam(params url.Values, args map[string]interface{}) {
	limit := svccommon.GetIntArg(args, 0, "limit")
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
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

func resolveOrganization(args map[string]interface{}, service ServiceInterface, required bool) (string, error) {
	if organization, ok := svccommon.GetStringArg(args, "organization"); ok {
		return organization, nil
	}
	if organization := service.GetDefaultOrganization(); organization != "" {
		return organization, nil
	}
	if required {
		return "", fmt.Errorf("missing required parameter: organization")
	}
	return "", nil
}

func resolveProject(args map[string]interface{}, service ServiceInterface, required bool) (string, error) {
	if project, ok := svccommon.GetStringArg(args, "project"); ok {
		return project, nil
	}
	if project := service.GetDefaultProject(); project != "" {
		return project, nil
	}
	if required {
		return "", fmt.Errorf("missing required parameter: project")
	}
	return "", nil
}

func compactIssue(issue map[string]interface{}) map[string]interface{} {
	summary := map[string]interface{}{}
	for _, key := range []string{"id", "shortId", "title", "culprit", "level", "status", "count", "userCount", "permalink", "firstSeen", "lastSeen"} {
		if value, ok := issue[key]; ok && value != nil {
			summary[key] = value
		}
	}
	if value, ok := issue["project"]; ok && value != nil {
		summary["project"] = value
	}
	return summary
}

func marshalResult(result interface{}) (*mcp.CallToolResult, error) {
	jsonResponse, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize response: %w", err)
	}
	return mcp.NewToolResultText(string(jsonResponse)), nil
}

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/alertmanager/client"
	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/sanitize"
)

var logger = logrus.WithField("component", "alertmanager-handlers")

// marshalIndentJSON performs indented JSON encoding using object pool
func marshalIndentJSON(data interface{}) ([]byte, error) {
	// First encode to compact format using object pool
	compactBytes, err := optimize.GlobalJSONPool.MarshalToBytes(data)
	if err != nil {
		return nil, err
	}

	// For scenarios requiring indented display, still use standard library but reduce allocations
	var result bytes.Buffer
	err = json.Indent(&result, compactBytes, "", "  ")
	return result.Bytes(), err
}

// HandleGetStatus handles the alertmanager_get_status tool
func HandleGetStatus(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logger.WithField("tool", "alertmanager_get_status").Debug("Handling get status request")

		status, err := client.GetStatus(ctx)
		if err != nil {
			logger.WithError(err).Error("Failed to get Alertmanager status")
			return nil, fmt.Errorf("failed to get Alertmanager status: %w", err)
		}

		content, err := marshalIndentJSON(status)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal status: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(content)),
			},
		}, nil
	}
}

// HandleGetAlerts handles the alertmanager_get_alerts tool
func HandleGetAlerts(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logger.WithField("tool", "alertmanager_get_alerts").Debug("Handling get alerts request")

		// Parse filters from arguments
		var filters map[string]string
		if filtersArg, ok := request.GetArguments()["filters"]; ok {
			if filtersMap, ok := filtersArg.(map[string]interface{}); ok {
				filters = make(map[string]string)
				for k, v := range filtersMap {
					// Sanitize filter values to prevent injection attacks
					if s, ok := v.(string); ok {
						filters[k] = sanitize.SanitizeFilterValue(s)
					}
				}
			}
		}

		alerts, err := client.GetAlerts(ctx, filters)
		if err != nil {
			logger.WithError(err).Error("Failed to get alerts")
			return nil, fmt.Errorf("failed to get alerts: %w", err)
		}

		content, err := marshalIndentJSON(alerts)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal alerts: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Found %d alerts:\n%s", len(alerts), string(content))),
			},
		}, nil
	}
}

// HandleGetAlertGroups handles the alertmanager_get_alert_groups tool
func HandleGetAlertGroups(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logger.WithField("tool", "alertmanager_get_alert_groups").Debug("Handling get alert groups request")

		groups, err := client.GetAlertGroups(ctx)
		if err != nil {
			logger.WithError(err).Error("Failed to get alert groups")
			return nil, fmt.Errorf("failed to get alert groups: %w", err)
		}

		content, err := marshalIndentJSON(groups)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal alert groups: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Found %d alert groups:\n%s", len(groups), string(content))),
			},
		}, nil
	}
}

// HandleGetSilences handles the alertmanager_get_silences tool
func HandleGetSilences(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logger.WithField("tool", "alertmanager_get_silences").Debug("Handling get silences request")

		silences, err := client.GetSilences(ctx)
		if err != nil {
			logger.WithError(err).Error("Failed to get silences")
			return nil, fmt.Errorf("failed to get silences: %w", err)
		}

		content, err := marshalIndentJSON(silences)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal silences: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Found %d silences:\n%s", len(silences), string(content))),
			},
		}, nil
	}
}

// HandleCreateSilence handles the alertmanager_create_silence tool
func HandleCreateSilence(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logger.WithField("tool", "alertmanager_create_silence").Debug("Handling create silence request")

		// Parse silence configuration
		silence := make(map[string]interface{})

		// Parse matchers
		if matchersArg, ok := request.GetArguments()["matchers"]; ok {
			silence["matchers"] = matchersArg
		} else {
			return nil, fmt.Errorf("matchers are required")
		}

		// Parse times
		if startsAtArg, ok := request.GetArguments()["startsAt"]; ok {
			silence["startsAt"] = startsAtArg
		} else {
			// Default to current time
			silence["startsAt"] = time.Now().Format(time.RFC3339)
		}

		if endsAtArg, ok := request.GetArguments()["endsAt"]; ok {
			silence["endsAt"] = endsAtArg
		} else {
			return nil, fmt.Errorf("endsAt is required")
		}

		// Parse comment and creator
		if commentArg, ok := request.GetArguments()["comment"]; ok {
			silence["comment"] = commentArg
		} else {
			return nil, fmt.Errorf("comment is required")
		}

		if createdByArg, ok := request.GetArguments()["createdBy"]; ok {
			silence["createdBy"] = createdByArg
		} else {
			return nil, fmt.Errorf("createdBy is required")
		}

		result, err := client.CreateSilence(ctx, silence)
		if err != nil {
			logger.WithError(err).Error("Failed to create silence")
			return nil, fmt.Errorf("failed to create silence: %w", err)
		}

		content, err := marshalIndentJSON(result)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal result: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Silence created successfully:\n%s", string(content))),
			},
		}, nil
	}
}

// HandleDeleteSilence handles the alertmanager_delete_silence tool
func HandleDeleteSilence(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logger.WithField("tool", "alertmanager_delete_silence").Debug("Handling delete silence request")

		// Parse silence ID
		silenceID, ok := request.GetArguments()["silenceId"].(string)
		if !ok {
			return nil, fmt.Errorf("silenceId is required")
		}

		err := client.DeleteSilence(ctx, silenceID)
		if err != nil {
			logger.WithError(err).WithField("silenceId", silenceID).Error("Failed to delete silence")
			return nil, fmt.Errorf("failed to delete silence: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Silence %s deleted successfully", silenceID)),
			},
		}, nil
	}
}

// HandleGetReceivers handles the alertmanager_get_receivers tool
func HandleGetReceivers(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logger.WithField("tool", "alertmanager_get_receivers").Debug("Handling get receivers request")

		receivers, err := client.GetReceivers(ctx)
		if err != nil {
			logger.WithError(err).Error("Failed to get receivers")
			return nil, fmt.Errorf("failed to get receivers: %w", err)
		}

		content, err := marshalIndentJSON(receivers)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal receivers: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Found %d receivers:\n%s", len(receivers), string(content))),
			},
		}, nil
	}
}

// HandleTestReceiver handles the alertmanager_test_receiver tool
func HandleTestReceiver(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logger.WithField("tool", "alertmanager_test_receiver").Debug("Handling test receiver request")

		// Parse receiver configuration
		receiverArg, ok := request.GetArguments()["receiver"]
		if !ok {
			return nil, fmt.Errorf("receiver configuration is required")
		}

		// Parse alerts if provided
		testData := map[string]interface{}{
			"receiver": receiverArg,
		}

		if alertsArg, ok := request.GetArguments()["alerts"]; ok {
			testData["alerts"] = alertsArg
		}

		result, err := client.TestReceiver(ctx, testData)
		if err != nil {
			logger.WithError(err).Error("Failed to test receiver")
			return nil, fmt.Errorf("failed to test receiver: %w", err)
		}

		content, err := marshalIndentJSON(result)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal result: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Receiver test result:\n%s", string(content))),
			},
		}, nil
	}
}

// HandleQueryAlerts handles the alertmanager_query_alerts tool
func HandleQueryAlerts(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logger.WithField("tool", "alertmanager_query_alerts").Debug("Handling query alerts request")

		// Build filters from arguments
		filters := make(map[string]string)

		// Standard filters
		if receiver, ok := request.GetArguments()["receiver"].(string); ok {
			filters["receiver"] = receiver
		}

		if silenced, ok := request.GetArguments()["silenced"].(bool); ok {
			filters["silenced"] = fmt.Sprintf("%t", silenced)
		}

		if active, ok := request.GetArguments()["active"].(bool); ok {
			filters["active"] = fmt.Sprintf("%t", active)
		}

		if unprocessed, ok := request.GetArguments()["unprocessed"].(bool); ok {
			filters["unprocessed"] = fmt.Sprintf("%t", unprocessed)
		}

		if inhibited, ok := request.GetArguments()["inhibited"].(bool); ok {
			filters["inhibited"] = fmt.Sprintf("%t", inhibited)
		}

		if filter, ok := request.GetArguments()["filter"].(string); ok {
			filters["filter"] = filter
		}

		alerts, err := client.GetAlerts(ctx, filters)
		if err != nil {
			logger.WithError(err).Error("Failed to query alerts")
			return nil, fmt.Errorf("failed to query alerts: %w", err)
		}

		// Implement sorting if sortBy is specified
		if sortBy, ok := request.GetArguments()["sortBy"].(string); ok {
			sortAlerts(alerts, sortBy)
		}

		content, err := marshalIndentJSON(alerts)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal alerts: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Query returned %d alerts:\n%s", len(alerts), string(content))),
			},
		}, nil
	}
}

// sortAlerts sorts alerts based on the specified field
func sortAlerts(alerts []map[string]interface{}, sortBy string) {
	if len(alerts) == 0 {
		return
	}

	// Supported sort fields
	switch sortBy {
	case "severity", "severity_desc":
		sortBySeverity(alerts, sortBy == "severity_desc")
	case "startsAt", "startsAt_desc":
		sortByStartsAt(alerts, sortBy == "startsAt_desc")
	case "endsAt", "endsAt_desc":
		sortByEndsAt(alerts, sortBy == "endsAt_desc")
	case "fingerprint", "fingerprint_desc":
		sortByFingerprint(alerts, sortBy == "fingerprint_desc")
	default:
		logger.WithField("sortBy", sortBy).Warn("Unsupported sort field, skipping sort")
	}
}

// sortBySeverity sorts alerts by severity (critical > warning > info)
func sortBySeverity(alerts []map[string]interface{}, desc bool) {
	sort.Slice(alerts, func(i, j int) bool {
		severityI := getSeverity(alerts[i])
		severityJ := getSeverity(alerts[j])

		if desc {
			return severityI > severityJ
		}
		return severityI < severityJ
	})
}

// getSeverity extracts severity from alert labels
func getSeverity(alert map[string]interface{}) int {
	labels, ok := alert["labels"].(map[string]interface{})
	if !ok {
		return 0
	}

	severity, ok := labels["severity"].(string)
	if !ok {
		return 0
	}

	// Severity priority: critical=3, warning=2, info=1, other=0
	switch severity {
	case "critical":
		return 3
	case "warning":
		return 2
	case "info":
		return 1
	default:
		return 0
	}
}

// sortByStartsAt sorts alerts by start time
func sortByStartsAt(alerts []map[string]interface{}, desc bool) {
	sort.Slice(alerts, func(i, j int) bool {
		startsAtI, _ := time.Parse(time.RFC3339, alerts[i]["startsAt"].(string))
		startsAtJ, _ := time.Parse(time.RFC3339, alerts[j]["startsAt"].(string))

		if desc {
			return startsAtI.After(startsAtJ)
		}
		return startsAtI.Before(startsAtJ)
	})
}

// sortByEndsAt sorts alerts by end time
func sortByEndsAt(alerts []map[string]interface{}, desc bool) {
	sort.Slice(alerts, func(i, j int) bool {
		endsAtI, _ := time.Parse(time.RFC3339, alerts[i]["endsAt"].(string))
		endsAtJ, _ := time.Parse(time.RFC3339, alerts[j]["endsAt"].(string))

		if desc {
			return endsAtI.After(endsAtJ)
		}
		return endsAtI.Before(endsAtJ)
	})
}

// sortByFingerprint sorts alerts by fingerprint
func sortByFingerprint(alerts []map[string]interface{}, desc bool) {
	sort.Slice(alerts, func(i, j int) bool {
		fingerprintI := alerts[i]["fingerprint"].(string)
		fingerprintJ := alerts[j]["fingerprint"].(string)

		if desc {
			return fingerprintI > fingerprintJ
		}
		return fingerprintI < fingerprintJ
	})
}

// Helper function to validate and parse limit parameter with warnings
func parseLimitWithWarnings(request mcp.CallToolRequest, toolName string) int {
	limit := 20
	if v, ok := request.GetArguments()["limit"]; ok {
		if f, ok := v.(float64); ok {
			limit = int(f)
			if limit <= 0 {
				limit = 20
			} else if limit > 100 {
				logrus.WithField("requested", limit).WithField("max", 100).Warn("Limit too high, resetting to safe maximum")
				limit = 100
			}
		}
	}

	if limit > 50 {
		logrus.WithFields(logrus.Fields{
			"tool":  toolName,
			"limit": limit,
		}).Warn("Large limit may cause context overflow, consider using pagination")
	}

	return limit
}

// Helper function to get optional numeric parameter
func getOptionalIntParam(request mcp.CallToolRequest, param string, defaultValue int) int {
	if v, ok := request.GetArguments()[param]; ok {
		if f, ok := v.(float64); ok {
			val := int(f)
			if val > 0 {
				return val
			}
		} else if s, ok := v.(string); ok {
			if val, err := strconv.Atoi(s); err == nil && val > 0 {
				return val
			}
		}
	}
	return defaultValue
}

// Helper function to get optional boolean parameter
func getOptionalBoolParam(request mcp.CallToolRequest, param string) *bool {
	if value, ok := request.GetArguments()[param].(bool); ok {
		return &value
	}
	return nil
}

// Helper function to get optional string parameter
func getOptionalStringParam(request mcp.CallToolRequest, param string) string {
	value, _ := request.GetArguments()[param].(string)
	return value
}

// Helper function to marshal optimized response with size warning
func marshalOptimizedResponse(data any, toolName string) (*mcp.CallToolResult, error) {
	jsonResponse, err := optimize.GlobalJSONPool.MarshalToBytes(data)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize response: %w", err)
	}

	// Add size warning for large responses
	if len(jsonResponse) > 100*1024 { // 100KB
		logrus.WithFields(logrus.Fields{
			"tool":      toolName,
			"sizeBytes": len(jsonResponse),
			"sizeKB":    len(jsonResponse) / 1024,
		}).Warn("Large response generated")
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(jsonResponse)),
		},
	}, nil
}

// ⚠️ PRIORITY: Optimized handlers for LLM efficiency

// HandleAlertsSummary handles getting alerts summary with LLM optimization
func HandleAlertsSummary(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filter := getOptionalStringParam(request, "filter")
		receiver := getOptionalStringParam(request, "receiver")
		silenced := getOptionalBoolParam(request, "silenced")
		activeOnly := getOptionalBoolParam(request, "active_only")
		limit := parseLimitWithWarnings(request, "alertmanager_alerts_summary")

		logger.WithFields(logrus.Fields{
			"tool":       "alertmanager_alerts_summary",
			"filter":     filter,
			"receiver":   receiver,
			"silenced":   silenced,
			"activeOnly": activeOnly,
			"limit":      limit,
		}).Debug("Handler invoked")

		alerts, err := client.AlertsSummary(ctx, filter, receiver, silenced, activeOnly, limit)
		if err != nil {
			return nil, fmt.Errorf("failed to get alerts summary: %w", err)
		}

		response := map[string]interface{}{
			"alerts": alerts,
			"count":  len(alerts),
			"metadata": map[string]interface{}{
				"tool":         "alertmanager_alerts_summary",
				"optimizedFor": "LLM efficiency",
			},
		}

		return marshalOptimizedResponse(response, "alertmanager_alerts_summary")
	}
}

// HandleSilencesSummary handles getting silences summary with LLM optimization
func HandleSilencesSummary(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		status := getOptionalStringParam(request, "status")
		limit := parseLimitWithWarnings(request, "alertmanager_silences_summary")

		logger.WithFields(logrus.Fields{
			"tool":   "alertmanager_silences_summary",
			"status": status,
			"limit":  limit,
		}).Debug("Handler invoked")

		silences, err := client.SilencesSummary(ctx, status, limit)
		if err != nil {
			return nil, fmt.Errorf("failed to get silences summary: %w", err)
		}

		response := map[string]interface{}{
			"silences": silences,
			"count":    len(silences),
			"metadata": map[string]interface{}{
				"tool":         "alertmanager_silences_summary",
				"optimizedFor": "LLM efficiency",
			},
		}

		return marshalOptimizedResponse(response, "alertmanager_silences_summary")
	}
}

// HandleAlertGroupsPaginated handles paginated alert groups listing with LLM optimization
func HandleAlertGroupsPaginated(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		page := getOptionalIntParam(request, "page", 1)
		perPage := getOptionalIntParam(request, "per_page", 20)
		receiver := getOptionalStringParam(request, "receiver")
		activeOnly := getOptionalBoolParam(request, "active_only")
		sortBy := getOptionalStringParam(request, "sort_by")

		logger.WithFields(logrus.Fields{
			"tool":       "alertmanager_alert_groups_paginated",
			"page":       page,
			"perPage":    perPage,
			"receiver":   receiver,
			"activeOnly": activeOnly,
			"sortBy":     sortBy,
		}).Debug("Handler invoked")

		groups, pagination, err := client.AlertGroupsPaginated(ctx, page, perPage, receiver, activeOnly, sortBy)
		if err != nil {
			return nil, fmt.Errorf("failed to list alert groups paginated: %w", err)
		}

		response := map[string]interface{}{
			"alertGroups": groups,
			"count":       len(groups),
			"pagination":  pagination,
			"searchCriteria": map[string]interface{}{
				"receiver":   receiver,
				"activeOnly": activeOnly,
				"sortBy":     sortBy,
			},
			"metadata": map[string]interface{}{
				"tool":         "alertmanager_alert_groups_paginated",
				"optimizedFor": "LLM efficiency",
			},
		}

		return marshalOptimizedResponse(response, "alertmanager_alert_groups_paginated")
	}
}

// HandleSilencesPaginated handles paginated silences listing with LLM optimization
func HandleSilencesPaginated(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		page := getOptionalIntParam(request, "page", 1)
		perPage := getOptionalIntParam(request, "per_page", 20)
		status := getOptionalStringParam(request, "status")
		createdBy := getOptionalStringParam(request, "created_by")
		commentFilter := getOptionalStringParam(request, "comment_filter")

		logger.WithFields(logrus.Fields{
			"tool":          "alertmanager_silences_paginated",
			"page":          page,
			"perPage":       perPage,
			"status":        status,
			"createdBy":     createdBy,
			"commentFilter": commentFilter,
		}).Debug("Handler invoked")

		silences, pagination, err := client.SilencesPaginated(ctx, page, perPage, status, createdBy, commentFilter)
		if err != nil {
			return nil, fmt.Errorf("failed to list silences paginated: %w", err)
		}

		response := map[string]interface{}{
			"silences":   silences,
			"count":      len(silences),
			"pagination": pagination,
			"searchCriteria": map[string]interface{}{
				"status":        status,
				"createdBy":     createdBy,
				"commentFilter": commentFilter,
			},
			"metadata": map[string]interface{}{
				"tool":         "alertmanager_silences_paginated",
				"optimizedFor": "LLM efficiency",
			},
		}

		return marshalOptimizedResponse(response, "alertmanager_silences_paginated")
	}
}

// HandleReceiversSummary handles getting receivers summary with LLM optimization
func HandleReceiversSummary(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		testInfo := getOptionalBoolParam(request, "test_info")

		logger.WithFields(logrus.Fields{
			"tool":     "alertmanager_receivers_summary",
			"testInfo": testInfo,
		}).Debug("Handler invoked")

		receivers, err := client.ReceiversSummary(ctx, testInfo)
		if err != nil {
			return nil, fmt.Errorf("failed to get receivers summary: %w", err)
		}

		response := map[string]interface{}{
			"receivers": receivers,
			"count":     len(receivers),
			"metadata": map[string]interface{}{
				"tool":         "alertmanager_receivers_summary",
				"optimizedFor": "LLM efficiency",
			},
		}

		return marshalOptimizedResponse(response, "alertmanager_receivers_summary")
	}
}

// HandleQueryAlertsAdvanced handles advanced alert querying with enhanced filters
func HandleQueryAlertsAdvanced(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filter := getOptionalStringParam(request, "filter")
		receiver := getOptionalStringParam(request, "receiver")
		silenced := getOptionalBoolParam(request, "silenced")
		active := getOptionalBoolParam(request, "active")
		inhibited := getOptionalBoolParam(request, "inhibited")
		timeRange := getOptionalStringParam(request, "time_range")
		page := getOptionalIntParam(request, "page", 1)
		perPage := getOptionalIntParam(request, "per_page", 30)
		sortBy := getOptionalStringParam(request, "sort_by")
		sortOrder := getOptionalStringParam(request, "sort_order")
		includeLabels := getOptionalBoolParam(request, "include_labels")

		logger.WithFields(logrus.Fields{
			"tool":          "alertmanager_query_alerts_advanced",
			"filter":        filter,
			"receiver":      receiver,
			"silenced":      silenced,
			"active":        active,
			"inhibited":     inhibited,
			"timeRange":     timeRange,
			"page":          page,
			"perPage":       perPage,
			"sortBy":        sortBy,
			"sortOrder":     sortOrder,
			"includeLabels": includeLabels,
		}).Debug("Handler invoked")

		alerts, pagination, err := client.QueryAlertsAdvanced(ctx, filter, receiver, silenced, active, inhibited, timeRange, page, perPage, sortBy, sortOrder, includeLabels)
		if err != nil {
			return nil, fmt.Errorf("failed to query alerts advanced: %w", err)
		}

		response := map[string]interface{}{
			"alerts": alerts,
			"count":  len(alerts),
			"searchCriteria": map[string]interface{}{
				"filter":        filter,
				"receiver":      receiver,
				"silenced":      silenced,
				"active":        active,
				"inhibited":     inhibited,
				"timeRange":     timeRange,
				"sortBy":        sortBy,
				"sortOrder":     sortOrder,
				"includeLabels": includeLabels,
			},
			"pagination": pagination,
			"metadata": map[string]interface{}{
				"tool":         "alertmanager_query_alerts_advanced",
				"optimizedFor": "finding specific alerts",
			},
		}

		return marshalOptimizedResponse(response, "alertmanager_query_alerts_advanced")
	}
}

// HandleHealthSummary handles getting Alertmanager health summary
func HandleHealthSummary(client *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		level := getOptionalStringParam(request, "level")
		if level == "" {
			level = "basic"
		}

		includeCluster := getOptionalBoolParam(request, "include_cluster")

		logger.WithFields(logrus.Fields{
			"tool":           "alertmanager_health_summary",
			"level":          level,
			"includeCluster": includeCluster,
		}).Debug("Handler invoked")

		health, err := client.GetHealthSummary(ctx, level, includeCluster)
		if err != nil {
			return nil, fmt.Errorf("failed to get health summary: %w", err)
		}

		response := map[string]interface{}{
			"health": health,
			"metadata": map[string]interface{}{
				"tool":           "alertmanager_health_summary",
				"level":          level,
				"includeCluster": includeCluster,
				"optimizedFor":   "monitoring and LLM efficiency",
			},
		}

		return marshalOptimizedResponse(response, "alertmanager_health_summary")
	}
}

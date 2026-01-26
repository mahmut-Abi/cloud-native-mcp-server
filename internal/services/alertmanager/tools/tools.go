// Package tools defines MCP tools for Alertmanager operations.
// It provides tool definitions for alert management, silence operations, and monitoring.
package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
)

// GetStatusTool returns a tool for retrieving Alertmanager status
func GetStatusTool() mcp.Tool {
	return mcp.NewTool("alertmanager_get_status",
		mcp.WithDescription("Get Alertmanager status and configuration information"),
	)
}

// GetAlertsTool returns a tool for retrieving current alerts
func GetAlertsTool() mcp.Tool {
	return mcp.NewTool("alertmanager_get_alerts",
		mcp.WithDescription("Get current alerts from Alertmanager with optional filtering"),
		mcp.WithObject("filters",
			mcp.Description("Optional filters for alerts (e.g., receiver, silenced, active, etc.)"),
		),
	)
}

// GetAlertGroupsTool returns a tool for retrieving alert groups
func GetAlertGroupsTool() mcp.Tool {
	return mcp.NewTool("alertmanager_get_alert_groups",
		mcp.WithDescription("Get alert groups from Alertmanager"),
	)
}

// GetSilencesTool returns a tool for retrieving current silences
func GetSilencesTool() mcp.Tool {
	return mcp.NewTool("alertmanager_get_silences",
		mcp.WithDescription("Get current silences from Alertmanager"),
	)
}

// CreateSilenceTool returns a tool for creating alert silences
func CreateSilenceTool() mcp.Tool {
	return mcp.NewTool("alertmanager_create_silence",
		mcp.WithDescription("Create a new silence in Alertmanager to suppress matching alerts"),
		mcp.WithArray("matchers", mcp.Required(),
			mcp.Description("List of matchers to identify which alerts to silence")),
		mcp.WithString("startsAt",
			mcp.Description("Start time of the silence (RFC3339 format)")),
		mcp.WithString("endsAt", mcp.Required(),
			mcp.Description("End time of the silence (RFC3339 format)")),
		mcp.WithString("comment", mcp.Required(),
			mcp.Description("Comment explaining the reason for the silence")),
		mcp.WithString("createdBy", mcp.Required(),
			mcp.Description("User who created the silence")),
	)
}

// DeleteSilenceTool returns a tool for deleting silences
func DeleteSilenceTool() mcp.Tool {
	return mcp.NewTool("alertmanager_delete_silence",
		mcp.WithDescription("Delete an existing silence from Alertmanager"),
		mcp.WithString("silenceId", mcp.Required(),
			mcp.Description("ID of the silence to delete")),
	)
}

// GetReceiversTool returns a tool for retrieving configured receivers
func GetReceiversTool() mcp.Tool {
	return mcp.NewTool("alertmanager_get_receivers",
		mcp.WithDescription("Get configured receivers from Alertmanager"),
	)
}

// TestReceiverTool returns a tool for testing receiver configurations
func TestReceiverTool() mcp.Tool {
	return mcp.NewTool("alertmanager_test_receiver",
		mcp.WithDescription("Test a receiver configuration by sending test notifications"),
		mcp.WithObject("receiver", mcp.Required(),
			mcp.Description("Receiver configuration to test")),
		mcp.WithArray("alerts",
			mcp.Description("Test alerts to send to the receiver")),
	)
}

// QueryAlertsTool returns a tool for querying alerts with advanced filters
func QueryAlertsTool() mcp.Tool {
	return mcp.NewTool("alertmanager_query_alerts",
		mcp.WithDescription("Query alerts with advanced filtering and sorting options"),
		mcp.WithString("receiver",
			mcp.Description("Filter by receiver name")),
		mcp.WithBoolean("silenced",
			mcp.Description("Filter by silenced status")),
		mcp.WithBoolean("active",
			mcp.Description("Filter by active status")),
		mcp.WithBoolean("unprocessed",
			mcp.Description("Filter by unprocessed status")),
		mcp.WithBoolean("inhibited",
			mcp.Description("Filter by inhibited status")),
		mcp.WithString("filter",
			mcp.Description("Label filter expression (e.g., 'alertname=\"HighCPU\"')")),
		mcp.WithString("sortBy",
			mcp.Description("Sort field: startsAt, endsAt, or updatedAt")),
		mcp.WithString("sortOrder",
			mcp.Description("Sort order: asc or desc")),
	)
}

// ⚠️ PRIORITY: Optimized tools for LLM efficiency

// GetAlertsSummaryTool returns tool definition for getting alerts summary
func GetAlertsSummaryTool() mcp.Tool {
	return mcp.NewTool("alertmanager_alerts_summary",
		mcp.WithDescription("⚠️ PRIORITY: Get Alertmanager alerts summary (state, count, receiver). 85-95% smaller output. Optimized for LLM efficiency."),
		mcp.WithString("filter"),
		mcp.WithString("receiver"),
		mcp.WithBoolean("silenced"),
		mcp.WithBoolean("active_only"),
		mcp.WithNumber("limit"),
	)
}

// GetSilencesSummaryTool returns tool definition for getting silences summary
func GetSilencesSummaryTool() mcp.Tool {
	return mcp.NewTool("alertmanager_silences_summary",
		mcp.WithDescription("⚠️ PRIORITY: Get Alertmanager silences summary (id, status, comment, duration). 85-95% smaller output. Optimized for LLM efficiency."),
		mcp.WithString("status"),
		mcp.WithNumber("limit"),
	)
}

// GetAlertGroupsPaginatedTool returns tool definition for paginated alert groups listing
func GetAlertGroupsPaginatedTool() mcp.Tool {
	return mcp.NewTool("alertmanager_alert_groups_paginated",
		mcp.WithDescription("⚠️ PRIORITY: Optimized for LLM efficiency: List Alertmanager alert groups with pagination and summary output. 80-90% smaller than full listing."),
		mcp.WithNumber("page"),
		mcp.WithNumber("per_page"),
		mcp.WithString("receiver"),
		mcp.WithBoolean("active_only"),
		mcp.WithString("sort_by"),
	)
}

// GetSilencesPaginatedTool returns tool definition for paginated silences listing
func GetSilencesPaginatedTool() mcp.Tool {
	return mcp.NewTool("alertmanager_silences_paginated",
		mcp.WithDescription("⚠️ PRIORITY: Optimized for LLM efficiency: List Alertmanager silences with pagination and summary output. 80-90% smaller than full listing."),
		mcp.WithNumber("page"),
		mcp.WithNumber("per_page"),
		mcp.WithString("status"),
		mcp.WithString("created_by"),
		mcp.WithString("comment_filter"),
	)
}

// GetReceiversSummaryTool returns tool definition for getting receivers summary
func GetReceiversSummaryTool() mcp.Tool {
	return mcp.NewTool("alertmanager_receivers_summary",
		mcp.WithDescription("⚠️ PRIORITY: Get Alertmanager receivers summary (name, type, status). 85-95% smaller output. Optimized for LLM efficiency."),
		mcp.WithBoolean("test_info"),
	)
}

// QueryAlertsAdvancedTool returns tool definition for advanced alert querying
func QueryAlertsAdvancedTool() mcp.Tool {
	return mcp.NewTool("alertmanager_query_alerts_advanced",
		mcp.WithDescription("🔍 Advanced Alertmanager query with enhanced filters, pagination, and sorting. Optimized for finding specific alerts."),
		mcp.WithString("filter"),
		mcp.WithString("receiver"),
		mcp.WithBoolean("silenced"),
		mcp.WithBoolean("active"),
		mcp.WithBoolean("inhibited"),
		mcp.WithString("time_range"),
		mcp.WithNumber("page"),
		mcp.WithNumber("per_page"),
		mcp.WithString("sort_by"),
		mcp.WithString("sort_order"),
		mcp.WithBoolean("include_labels"),
	)
}

// GetHealthStatusSummaryTool returns tool definition for health status summary
func GetHealthStatusSummaryTool() mcp.Tool {
	return mcp.NewTool("alertmanager_health_summary",
		mcp.WithDescription("⚠️ PRIORITY: Get Alertmanager health and status summary (uptime, cluster state, performance metrics). Lightweight status overview. Optimized for monitoring."),
		mcp.WithString("level"),
		mcp.WithBoolean("include_cluster"),
	)
}

package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryTool(t *testing.T) {
	tool := QueryTool()
	assert.Equal(t, "prometheus_query", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestQueryRangeTool(t *testing.T) {
	tool := QueryRangeTool()
	assert.Equal(t, "prometheus_query_range", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetTargetsTool(t *testing.T) {
	tool := GetTargetsTool()
	assert.Equal(t, "prometheus_get_targets", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetAlertsTool(t *testing.T) {
	tool := GetAlertsTool()
	assert.Equal(t, "prometheus_get_alerts", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetRulesTool(t *testing.T) {
	tool := GetRulesTool()
	assert.Equal(t, "prometheus_get_rules", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetLabelNamesTool(t *testing.T) {
	tool := GetLabelNamesTool()
	assert.Equal(t, "prometheus_get_label_names", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetLabelValuesTool(t *testing.T) {
	tool := GetLabelValuesTool()
	assert.Equal(t, "prometheus_get_label_values", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetSeriesTool(t *testing.T) {
	tool := GetSeriesTool()
	assert.Equal(t, "prometheus_get_series", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestTestConnectionTool(t *testing.T) {
	tool := TestConnectionTool()
	assert.Equal(t, "prometheus_test_connection", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetServerInfoTool(t *testing.T) {
	tool := GetServerInfoTool()
	assert.Equal(t, "prometheus_get_server_info", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetMetricsMetadataTool(t *testing.T) {
	tool := GetMetricsMetadataTool()
	assert.Equal(t, "prometheus_get_metrics_metadata", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetTargetMetadataTool(t *testing.T) {
	tool := GetTargetMetadataTool()
	assert.Equal(t, "prometheus_get_target_metadata", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetTSDBStatsTool(t *testing.T) {
	tool := GetTSDBStatsTool()
	assert.Equal(t, "prometheus_get_tsdb_stats", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetTSDBStatusTool(t *testing.T) {
	tool := GetTSDBStatusTool()
	assert.Equal(t, "prometheus_get_tsdb_status", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetRuntimeInfoTool(t *testing.T) {
	tool := GetRuntimeInfoTool()
	assert.Equal(t, "prometheus_get_runtime_info", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestCreateSnapshotTool(t *testing.T) {
	tool := CreateSnapshotTool()
	assert.Equal(t, "prometheus_create_snapshot", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetWALReplayStatusTool(t *testing.T) {
	tool := GetWALReplayStatusTool()
	assert.Equal(t, "prometheus_get_wal_replay_status", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetTargetsSummaryTool(t *testing.T) {
	tool := GetTargetsSummaryTool()
	assert.Equal(t, "prometheus_targets_summary", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetAlertsSummaryTool(t *testing.T) {
	tool := GetAlertsSummaryTool()
	assert.Equal(t, "prometheus_alerts_summary", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetRulesSummaryTool(t *testing.T) {
	tool := GetRulesSummaryTool()
	assert.Equal(t, "prometheus_rules_summary", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}
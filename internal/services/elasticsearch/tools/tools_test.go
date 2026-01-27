package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckTool(t *testing.T) {
	tool := HealthCheckTool()
	assert.Equal(t, "elasticsearch_health", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestListIndicesTool(t *testing.T) {
	tool := ListIndicesTool()
	assert.Equal(t, "elasticsearch_list_indices", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetIndexStatsTool(t *testing.T) {
	tool := GetIndexStatsTool()
	assert.Equal(t, "elasticsearch_index_stats", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetNodesTool(t *testing.T) {
	tool := GetNodesTool()
	assert.Equal(t, "elasticsearch_nodes", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetIndicesSummaryTool(t *testing.T) {
	tool := GetIndicesSummaryTool()
	assert.Equal(t, "elasticsearch_indices_summary", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetInfoTool(t *testing.T) {
	tool := GetInfoTool()
	assert.Equal(t, "elasticsearch_info", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestListIndicesPaginatedTool(t *testing.T) {
	tool := ListIndicesPaginatedTool()
	assert.Equal(t, "elasticsearch_list_indices_paginated", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetNodesSummaryTool(t *testing.T) {
	tool := GetNodesSummaryTool()
	assert.Equal(t, "elasticsearch_nodes_summary", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetClusterHealthSummaryTool(t *testing.T) {
	tool := GetClusterHealthSummaryTool()
	assert.Equal(t, "elasticsearch_cluster_health_summary", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetIndexDetailAdvancedTool(t *testing.T) {
	tool := GetIndexDetailAdvancedTool()
	assert.Equal(t, "elasticsearch_get_index_detail_advanced", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestGetClusterDetailAdvancedTool(t *testing.T) {
	tool := GetClusterDetailAdvancedTool()
	assert.Equal(t, "elasticsearch_get_cluster_detail_advanced", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}

func TestSearchIndicesTool(t *testing.T) {
	tool := SearchIndicesTool()
	assert.Equal(t, "elasticsearch_search_indices", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.InputSchema)
}
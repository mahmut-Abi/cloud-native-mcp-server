// Package tools provides comprehensive tests for the tools optimization framework.
package tools

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// TestSchemaBuilder_StringParam tests string parameter addition.
func TestSchemaBuilder_StringParam(t *testing.T) {
	schema := NewSchemaBuilder().
		StringParam("test", "test description", true).
		Build()

	if len(schema.Required) != 1 || schema.Required[0] != "test" {
		t.Error("Expected 'test' in required parameters")
	}

	if schema.Properties["test"] == nil {
		t.Error("Expected 'test' property in schema")
	}
}

// TestSchemaBuilder_MultipleParams tests adding multiple parameters.
func TestSchemaBuilder_MultipleParams(t *testing.T) {
	schema := NewSchemaBuilder().
		StringParam("kind", "Kind description", true).
		StringParam("name", "Name description", true).
		StringParam("namespace", "Namespace description", false).
		StringParam("debug", "Debug description", false).
		Build()

	if len(schema.Properties) != 4 {
		t.Errorf("Expected 4 properties, got %d", len(schema.Properties))
	}

	if len(schema.Required) != 2 {
		t.Errorf("Expected 2 required parameters, got %d", len(schema.Required))
	}
}

// TestSchemaBuilder_EnumParam tests enum parameter addition.
func TestSchemaBuilder_EnumParam(t *testing.T) {
	enumValues := []string{"active", "dropped", "any"}
	schema := NewSchemaBuilder().
		EnumParam("state", "State description", enumValues, true).
		Build()

	if len(schema.Required) != 1 {
		t.Error("Expected 'state' in required parameters")
	}

	prop := schema.Properties["state"].(map[string]interface{})
	if prop["type"] != "string" {
		t.Error("Expected string type for enum")
	}
}

// TestSchemaBuilder_Fluent tests fluent interface chaining.
func TestSchemaBuilder_Fluent(t *testing.T) {
	builder := NewSchemaBuilder()
	result := builder.
		StringParam("a", "A", true).
		StringParam("b", "B", false).
		BooleanParam("c", "C").
		Build()

	if result.Type != "object" {
		t.Error("Expected object type")
	}
}

// TestCommonDescriptions tests common description access.
func TestCommonDescriptions_Access(t *testing.T) {
	tests := []string{
		"k8s_kind",
		"k8s_name",
		"k8s_namespace",
		"k8s_debug",
		"prom_query",
		"grafana_uid",
		"kibana_space_id",
	}

	for _, test := range tests {
		desc := GetCommonDescription(test)
		if desc == "" {
			t.Errorf("Expected description for %s", test)
		}
	}
}

// TestToolCache_GetOrCreate tests tool caching.
func TestToolCache_GetOrCreate(t *testing.T) {
	cache := &ToolCache{
		tools: make(map[string]mcp.Tool),
	}

	callCount := 0
	builder := func() mcp.Tool {
		callCount++
		return mcp.Tool{Name: "test_tool"}
	}

	tool1 := cache.GetOrCreate("test", builder)
	tool2 := cache.GetOrCreate("test", builder)

	if callCount != 1 {
		t.Errorf("Expected builder called once, was called %d times", callCount)
	}

	if tool1.Name != tool2.Name {
		t.Error("Expected same tool from cache")
	}
}

// TestNewStringProperty tests string property creation.
func TestNewStringProperty(t *testing.T) {
	prop := NewStringProperty("test description")

	if prop["type"] != "string" {
		t.Error("Expected string type")
	}

	if prop["description"] != "test description" {
		t.Error("Expected description")
	}
}

// TestNewNumberProperty tests number property creation.
func TestNewNumberProperty(t *testing.T) {
	enum := []float64{1, 2, 3}
	prop := NewNumberProperty("test", enum)

	if prop["type"] != "number" {
		t.Error("Expected number type")
	}
}

// TestNewBooleanProperty tests boolean property creation.
func TestNewBooleanProperty(t *testing.T) {
	prop := NewBooleanProperty("test")

	if prop["type"] != "boolean" {
		t.Error("Expected boolean type")
	}
}

// TestNewEnumProperty tests enum property creation.
func TestNewEnumProperty(t *testing.T) {
	enumValues := []string{"a", "b", "c"}
	prop := NewEnumProperty("test", enumValues)

	if prop["type"] != "string" {
		t.Error("Expected string type")
	}

	if len(prop["enum"].([]string)) != 3 {
		t.Error("Expected 3 enum values")
	}
}

// TestCreateObjectSchema tests object schema creation.
func TestCreateObjectSchema(t *testing.T) {
	params := map[string]map[string]interface{}{
		"param1": {"type": "string"},
	}
	required := []string{"param1"}

	schema := CreateObjectSchema(params, required)

	if schema.Type != "object" {
		t.Error("Expected object type")
	}

	if len(schema.Properties) != 1 {
		t.Error("Expected 1 property")
	}
}

// TestKubernetesResourceParams tests Kubernetes resource parameters helper.
func TestKubernetesResourceParams(t *testing.T) {
	schema := KubernetesResourceParams().Build()

	if schema.Properties["kind"] == nil {
		t.Error("Expected 'kind' parameter")
	}

	if schema.Properties["name"] == nil {
		t.Error("Expected 'name' parameter")
	}

	if schema.Properties["namespace"] == nil {
		t.Error("Expected 'namespace' parameter")
	}
}

// TestPrometheusQueryParams tests Prometheus query parameters helper.
func TestPrometheusQueryParams(t *testing.T) {
	schema := PrometheusQueryParams(true).Build()

	if schema.Properties["query"] == nil {
		t.Error("Expected 'query' parameter")
	}

	if schema.Properties["time"] == nil {
		t.Error("Expected 'time' parameter")
	}
}

// TestPrometheusRangeParams tests Prometheus range parameters helper.
func TestPrometheusRangeParams(t *testing.T) {
	schema := PrometheusRangeParams().Build()

	if schema.Properties["query"] == nil {
		t.Error("Expected 'query' parameter")
	}

	if schema.Properties["start"] == nil {
		t.Error("Expected 'start' parameter")
	}

	if schema.Properties["end"] == nil {
		t.Error("Expected 'end' parameter")
	}

	if schema.Properties["step"] == nil {
		t.Error("Expected 'step' parameter")
	}
}

// Package tools provides schema helper utilities to reduce boilerplate in tool definitions.
package tools

import "github.com/mark3labs/mcp-go/mcp"

// SchemaBuilder provides a fluent interface for building tool input schemas.
type SchemaBuilder struct {
	properties map[string]map[string]interface{}
	required   []string
}

// NewSchemaBuilder creates a new schema builder instance.
func NewSchemaBuilder() *SchemaBuilder {
	return &SchemaBuilder{
		properties: make(map[string]map[string]interface{}),
		required:   []string{},
	}
}

// StringParam adds a string parameter to the schema.
func (sb *SchemaBuilder) StringParam(name, description string, required bool) *SchemaBuilder {
	sb.properties[name] = NewStringProperty(description)
	if required {
		sb.required = append(sb.required, name)
	}
	return sb
}

// NumberParam adds a number parameter to the schema.
func (sb *SchemaBuilder) NumberParam(name, description string, enum []float64, required bool) *SchemaBuilder {
	sb.properties[name] = NewNumberProperty(description, enum)
	if required {
		sb.required = append(sb.required, name)
	}
	return sb
}

// BooleanParam adds a boolean parameter to the schema.
func (sb *SchemaBuilder) BooleanParam(name, description string) *SchemaBuilder {
	sb.properties[name] = NewBooleanProperty(description)
	return sb
}

// EnumParam adds an enum parameter to the schema.
func (sb *SchemaBuilder) EnumParam(name, description string, values []string, required bool) *SchemaBuilder {
	sb.properties[name] = NewEnumProperty(description, values)
	if required {
		sb.required = append(sb.required, name)
	}
	return sb
}

// CustomParam adds a custom parameter property to the schema.
func (sb *SchemaBuilder) CustomParam(name string, property map[string]interface{}, required bool) *SchemaBuilder {
	sb.properties[name] = property
	if required {
		sb.required = append(sb.required, name)
	}
	return sb
}

// Build returns the completed ToolInputSchema.
// Build returns the completed ToolInputSchema.
func (sb *SchemaBuilder) Build() mcp.ToolInputSchema {
	properties := make(map[string]any)
	for k, v := range sb.properties {
		properties[k] = v
	}
	return mcp.ToolInputSchema{
		Type:       "object",
		Properties: properties,
		Required:   sb.required,
	}
}

// Kubernetes Schema Helpers (reduce duplication in k8s tools)

// KubernetesResourceParams returns a pre-built schema for standard Kubernetes resource parameters.
func KubernetesResourceParams() *SchemaBuilder {
	return NewSchemaBuilder().
		StringParam("kind", GetCommonDescription("k8s_kind"), true).
		StringParam("name", GetCommonDescription("k8s_name"), true).
		StringParam("namespace", GetCommonDescription("k8s_namespace"), false)
}

// KubernetesNamespacedParams returns params for namespaced resource operations.
func KubernetesNamespacedParams() *SchemaBuilder {
	return KubernetesResourceParams().
		StringParam("debug", GetCommonDescription("k8s_debug"), false)
}

// PrometheusQueryParams returns a pre-built schema for Prometheus query parameters.
func PrometheusQueryParams(includeTime bool) *SchemaBuilder {
	sb := NewSchemaBuilder().
		StringParam("query", GetCommonDescription("prom_query"), true)

	if includeTime {
		sb.StringParam("time", GetCommonDescription("prom_time"), false)
	}

	return sb
}

// PrometheusRangeParams returns a pre-built schema for Prometheus range query parameters.
func PrometheusRangeParams() *SchemaBuilder {
	return NewSchemaBuilder().
		StringParam("query", GetCommonDescription("prom_query"), true).
		StringParam("start", GetCommonDescription("start_time"), true).
		StringParam("end", GetCommonDescription("end_time"), true).
		StringParam("step", GetCommonDescription("prom_step"), false)
}

// GrafanaParams returns a pre-built schema for common Grafana parameters.
func GrafanaParams() *SchemaBuilder {
	return NewSchemaBuilder().
		StringParam("debug", GetCommonDescription("debug_verbose"), false)
}

// KibanaParams returns a pre-built schema for common Kibana parameters.
func KibanaParams() *SchemaBuilder {
	return NewSchemaBuilder().
		StringParam("debug", GetCommonDescription("debug_verbose"), false)
}

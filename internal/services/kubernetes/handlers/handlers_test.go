package handlers

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestRequireRawStringParamPreservesJSON(t *testing.T) {
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"metadata": "{\"name\":\"test-otel-debug\"}",
			},
		},
	}

	got, err := requireRawStringParam(req, "metadata")
	if err != nil {
		t.Fatalf("requireRawStringParam returned error: %v", err)
	}

	want := "{\"name\":\"test-otel-debug\"}"
	if got != want {
		t.Fatalf("requireRawStringParam = %q, want %q", got, want)
	}
}

func TestRequireStringParamSanitizesJSON(t *testing.T) {
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"metadata": "{\"name\":\"test-otel-debug\"}",
			},
		},
	}

	got, err := requireStringParam(req, "metadata")
	if err != nil {
		t.Fatalf("requireStringParam returned error: %v", err)
	}

	want := "{name:test-otel-debug}"
	if got != want {
		t.Fatalf("requireStringParam = %q, want %q", got, want)
	}
}

func TestGetOptionalRawStringParamPreservesJSON(t *testing.T) {
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"spec": "{\"replicas\":1}",
			},
		},
	}

	got := getOptionalRawStringParam(req, "spec")
	want := "{\"replicas\":1}"
	if got != want {
		t.Fatalf("getOptionalRawStringParam = %q, want %q", got, want)
	}
}

func TestRequireJSONObjectParamSupportsObjectAndJSONString(t *testing.T) {
	tests := []struct {
		name string
		arg  interface{}
	}{
		{
			name: "object input",
			arg: map[string]interface{}{
				"name": "test-otel-debug",
			},
		},
		{
			name: "json string input",
			arg:  "{\"name\":\"test-otel-debug\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{
						"metadata": tt.arg,
					},
				},
			}

			got, err := requireJSONObjectParam(req, "metadata")
			if err != nil {
				t.Fatalf("requireJSONObjectParam returned error: %v", err)
			}
			if got["name"] != "test-otel-debug" {
				t.Fatalf("requireJSONObjectParam name = %v, want test-otel-debug", got["name"])
			}
		})
	}
}

func TestRequireRawJSONParamSupportsObjectArrayAndString(t *testing.T) {
	tests := []struct {
		name string
		arg  interface{}
		want string
	}{
		{
			name: "object input",
			arg: map[string]interface{}{
				"spec": map[string]interface{}{"replicas": float64(3)},
			},
			want: `{"spec":{"replicas":3}}`,
		},
		{
			name: "array input",
			arg: []interface{}{
				map[string]interface{}{"op": "replace", "path": "/spec/replicas", "value": float64(3)},
			},
			want: `[{"op":"replace","path":"/spec/replicas","value":3}]`,
		},
		{
			name: "string input",
			arg:  `{"metadata":{"labels":{"app":"demo"}}}`,
			want: `{"metadata":{"labels":{"app":"demo"}}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{
						"patch": tt.arg,
					},
				},
			}

			got, err := requireRawJSONParam(req, "patch")
			if err != nil {
				t.Fatalf("requireRawJSONParam returned error: %v", err)
			}
			if string(got) != tt.want {
				t.Fatalf("requireRawJSONParam = %s, want %s", string(got), tt.want)
			}
		})
	}
}

func TestGetOptionalStringArrayParamSupportsStructuredAndLegacyInputs(t *testing.T) {
	tests := []struct {
		name string
		arg  interface{}
		want []string
	}{
		{
			name: "array input",
			arg:  []interface{}{"metadata.name", "status.phase"},
			want: []string{"metadata.name", "status.phase"},
		},
		{
			name: "json string array input",
			arg:  `["metadata.name","status.phase"]`,
			want: []string{"metadata.name", "status.phase"},
		},
		{
			name: "csv input",
			arg:  "metadata.name,status.phase",
			want: []string{"metadata.name", "status.phase"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{
						"jsonpaths": tt.arg,
					},
				},
			}

			got, err := getOptionalStringArrayParam(req, "jsonpaths")
			if err != nil {
				t.Fatalf("getOptionalStringArrayParam returned error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("getOptionalStringArrayParam length = %d, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("getOptionalStringArrayParam[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestNormalizeJSONPathExpression(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "{.metadata.name}", want: "{.metadata.name}"},
		{input: ".metadata.name", want: "{.metadata.name}"},
		{input: "metadata.name", want: "{.metadata.name}"},
	}

	for _, tt := range tests {
		if got := normalizeJSONPathExpression(tt.input); got != tt.want {
			t.Fatalf("normalizeJSONPathExpression(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestGetRequestArgumentsSupportsNestedParams(t *testing.T) {
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"params": map[string]interface{}{
					"kind":      "Pod",
					"namespace": "open-telemetry",
				},
				"debug": "true",
			},
		},
	}

	args := getRequestArguments(req)
	if args["kind"] != "Pod" {
		t.Fatalf("expected nested params kind to be merged, got %#v", args["kind"])
	}
	if args["namespace"] != "open-telemetry" {
		t.Fatalf("expected nested params namespace to be merged, got %#v", args["namespace"])
	}
	if args["debug"] != "true" {
		t.Fatalf("expected top-level debug to be preserved, got %#v", args["debug"])
	}
}

func TestGetOptionalSearchKinds(t *testing.T) {
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"resourceTypes": []interface{}{"pods", "deployments"},
			},
		},
	}

	kinds, err := getOptionalSearchKinds(req)
	if err != nil {
		t.Fatalf("getOptionalSearchKinds returned error: %v", err)
	}
	if len(kinds) != 2 || kinds[0] != "pods" || kinds[1] != "deployments" {
		t.Fatalf("unexpected kinds: %#v", kinds)
	}
}

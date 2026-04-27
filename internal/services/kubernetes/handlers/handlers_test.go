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

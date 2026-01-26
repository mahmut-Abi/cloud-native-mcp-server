package handlers

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestRequireStringParam_Success(t *testing.T) {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{"name": "pod-1"}
	v, err := requireStringParam(req, "name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "pod-1" {
		t.Fatalf("expected pod-1, got %s", v)
	}
}

func TestRequireStringParam_Missing(t *testing.T) {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{"x": "y"}
	_, err := requireStringParam(req, "name")
	if err == nil {
		t.Fatalf("expected error for missing param")
	}
}

func TestGetOptionalStringParam(t *testing.T) {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{"ns": "default"}
	if got := getOptionalStringParam(req, "ns"); got != "default" {
		t.Fatalf("expected default, got %s", got)
	}
	if got := getOptionalStringParam(req, "missing"); got != "" {
		t.Fatalf("expected empty for missing, got %s", got)
	}
}

func TestMarshalJSONResponse(t *testing.T) {
	ctx := context.Background()
	_ = ctx
	data := map[string]any{"ok": true}
	res, err := marshalJSONResponse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatalf("expected non-nil result")
	}
}

func TestCreateErrorResponse(t *testing.T) {
	res := createErrorResponse("test error")
	if res == nil {
		t.Fatalf("expected non-nil result")
	}
}

func TestGetNestedString(t *testing.T) {
	obj := map[string]any{
		"metadata": map[string]any{
			"name": "test-pod",
		},
	}

	// Test valid path
	if got := getNestedString(obj, "metadata.name"); got != "test-pod" {
		t.Fatalf("expected test-pod, got %s", got)
	}

	// Test missing path
	if got := getNestedString(obj, "missing.path"); got != "" {
		t.Fatalf("expected empty, got %s", got)
	}
}

func TestMustMarshalJSON(t *testing.T) {
	data := map[string]any{"key": "value"}
	result := mustMarshalJSON(data)
	if result == nil {
		t.Fatalf("expected non-nil result")
	}
}

func TestMarshalOptimizedResponse(t *testing.T) {
	data := map[string]any{"items": []any{"item1", "item2"}}
	res, err := marshalOptimizedResponse(data, "test_tool")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatalf("expected non-nil result")
	}
}

func TestApplyJSONPath(t *testing.T) {
	input := map[string]any{
		"items": []any{
			map[string]any{"name": "pod1"},
			map[string]any{"name": "pod2"},
		},
	}

	// Test valid JSONPath
	result, err := applyJSONPath(input, "$.items[*].name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatalf("expected non-nil result")
	}
}

func TestHandleDescribeResource(t *testing.T) {
	handler := HandleDescribeResource(nil)

	// Test missing required parameters
	reqMissing := mcp.CallToolRequest{}
	reqMissing.Params.Arguments = map[string]any{"name": "test-pod"}

	ctx := context.Background()
	_, err := handler(ctx, reqMissing)
	if err == nil {
		t.Fatalf("expected error for missing kind parameter")
	}
}

func TestHandleGetResourceUsage(t *testing.T) {
	handler := HandleGetResourceUsage(nil)

	// Test missing required parameter
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{}

	ctx := context.Background()
	_, err := handler(ctx, req)
	if err == nil {
		t.Fatalf("expected error for missing resourceType parameter")
	}
}

func TestHandleGetResourceDetails(t *testing.T) {
	handler := HandleGetResourceDetails(nil)

	// Test missing required parameters
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{"name": "test-pod"}

	ctx := context.Background()
	_, err := handler(ctx, req)
	if err == nil {
		t.Fatalf("expected error for missing kind parameter")
	}
}

func TestHandleContainerLogs(t *testing.T) {
	handler := HandleContainerLogs(nil)

	// Test missing required parameters
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{"namespace": "default"}

	ctx := context.Background()
	_, err := handler(ctx, req)
	if err == nil {
		t.Fatalf("expected error for missing pod parameter")
	}
}

func TestHandlePortForward(t *testing.T) {
	handler := HandlePortForward(nil)

	// Test missing required parameters
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{"namespace": "default", "pod": "test-pod"}

	ctx := context.Background()
	_, err := handler(ctx, req)
	if err == nil {
		t.Fatalf("expected error for missing localPort parameter")
	}
}

func TestHandleCreateResource(t *testing.T) {
	handler := HandleCreateResource(nil)

	// Test missing required parameters
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{"kind": "Pod", "namespace": "default"}

	ctx := context.Background()
	_, err := handler(ctx, req)
	if err == nil {
		t.Fatalf("expected error for missing manifest parameter")
	}
}

func TestHandleUpdateResource(t *testing.T) {
	handler := HandleUpdateResource(nil)

	// Test missing required parameters
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{"name": "test-pod", "namespace": "default"}

	ctx := context.Background()
	_, err := handler(ctx, req)
	if err == nil {
		t.Fatalf("expected error for missing kind parameter")
	}
}

func TestHandleContainerExec(t *testing.T) {
	handler := HandleContainerExec(nil)

	// Test missing required parameters
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{"namespace": "default"}

	ctx := context.Background()
	_, err := handler(ctx, req)
	if err == nil {
		t.Fatalf("expected error for missing pod parameter")
	}
}

func TestHandleGetResourceSummary(t *testing.T) {
	handler := HandleGetResourceSummary(nil)

	// Test missing required parameters
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{}

	ctx := context.Background()
	_, err := handler(ctx, req)
	if err == nil {
		t.Fatalf("expected error for missing kind parameter")
	}
}

func TestHandleGetResource(t *testing.T) {
	handler := HandleGetResource(nil)

	// Test missing required parameters
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{"kind": "Pod"}

	ctx := context.Background()
	_, err := handler(ctx, req)
	if err == nil {
		t.Fatalf("expected error for missing name parameter")
	}
}

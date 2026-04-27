package tools

import (
	"testing"
)

func TestGetResourceTool_Definition(t *testing.T) {
	tool := GetResourceTool()
	if tool.Name != "kubernetes_get_resource" {
		t.Fatalf("unexpected name: %s", tool.Name)
	}
}

func TestListResourcesTool_Definition(t *testing.T) {
	tool := ListResourcesTool()
	if tool.Name != "kubernetes_list_resources" {
		t.Fatalf("unexpected name: %s", tool.Name)
	}
}

func TestScaleResourceTool_Definition(t *testing.T) {
	tool := ScaleResourceTool()
	if tool.Name != "kubernetes_scale_resource" {
		t.Fatalf("unexpected name: %s", tool.Name)
	}
}

func TestRolloutAndNodeOperationTools_Definition(t *testing.T) {
	tests := []struct {
		name string
		tool string
	}{
		{"rollout", GetRolloutStatusTool().Name},
		{"cordon", CordonNodeTool().Name},
		{"uncordon", UncordonNodeTool().Name},
		{"drain", DrainNodeTool().Name},
	}

	expected := map[string]string{
		"rollout":  "kubernetes_get_rollout_status",
		"cordon":   "kubernetes_cordon_node",
		"uncordon": "kubernetes_uncordon_node",
		"drain":    "kubernetes_drain_node",
	}

	for _, tt := range tests {
		if tt.tool != expected[tt.name] {
			t.Fatalf("%s tool unexpected name: %s", tt.name, tt.tool)
		}
	}
}

func TestPatchResourceTool_Definition(t *testing.T) {
	tool := PatchResourceTool()
	if tool.Name != "kubernetes_patch_resource" {
		t.Fatalf("unexpected name: %s", tool.Name)
	}
	if patch, ok := tool.InputSchema.Properties["patch"].(map[string]any); !ok || len(patch) == 0 {
		t.Fatalf("patch schema should be present")
	}
}

func TestCreateResourceTool_Definition(t *testing.T) {
	tool := CreateResourceTool()
	if tool.Name != "kubernetes_create_resource" {
		t.Fatalf("unexpected name: %s", tool.Name)
	}
	metadata, ok := tool.InputSchema.Properties["metadata"].(map[string]any)
	if !ok || metadata["type"] != "object" {
		t.Fatalf("metadata schema should be object, got %#v", tool.InputSchema.Properties["metadata"])
	}
	spec, ok := tool.InputSchema.Properties["spec"].(map[string]any)
	if !ok || spec["type"] != "object" {
		t.Fatalf("spec schema should be object, got %#v", tool.InputSchema.Properties["spec"])
	}
}

func TestListResourcesTool_JSONPathsSchema(t *testing.T) {
	tool := ListResourcesTool()
	jsonpaths, ok := tool.InputSchema.Properties["jsonpaths"].(map[string]any)
	if !ok || jsonpaths["type"] != "array" {
		t.Fatalf("jsonpaths schema should be array, got %#v", tool.InputSchema.Properties["jsonpaths"])
	}
}

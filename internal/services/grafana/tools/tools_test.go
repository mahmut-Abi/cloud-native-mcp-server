package tools

import "testing"

func TestCreateAnnotationTool_TagsSchema(t *testing.T) {
	tool := CreateAnnotationTool()

	tags, ok := tool.InputSchema.Properties["tags"].(map[string]any)
	if !ok || tags["type"] != "array" {
		t.Fatalf("tags schema should be array, got %#v", tool.InputSchema.Properties["tags"])
	}

	items, ok := tags["items"].(map[string]any)
	if !ok || items["type"] != "string" {
		t.Fatalf("tags items schema should be string, got %#v", tags["items"])
	}
}

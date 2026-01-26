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

func TestUpdateResourceTool_Definition(t *testing.T) {
	tool := UpdateResourceTool()
	if tool.Name != "kubernetes_update_resource" {
		t.Fatalf("unexpected name: %s", tool.Name)
	}
}

func TestCreateResourceTool_Definition(t *testing.T) {
	tool := CreateResourceTool()
	if tool.Name != "kubernetes_create_resource" {
		t.Fatalf("unexpected name: %s", tool.Name)
	}
}

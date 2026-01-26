package tools

import (
	"testing"
)

func TestListReleasesTool(t *testing.T) {
	tool := ListReleasesTool()
	if tool.Name != "helm_list_releases" {
		t.Errorf("Expected tool name 'list_helm_releases', got '%s'", tool.Name)
	}
	if tool.Description == "" {
		t.Error("Expected tool description to be non-empty")
	}
}

func TestGetReleaseTool(t *testing.T) {
	tool := GetReleaseTool()
	if tool.Name != "helm_get_release" {
		t.Errorf("Expected tool name 'get_helm_release', got '%s'", tool.Name)
	}
	if tool.Description == "" {
		t.Error("Expected tool description to be non-empty")
	}
}

func TestListRepositoriesTool(t *testing.T) {
	tool := ListRepositoriesTool()
	if tool.Name != "helm_list_repos" {
		t.Errorf("Expected tool name 'list_helm_repos', got '%s'", tool.Name)
	}
	if tool.Description == "" {
		t.Error("Expected tool description to be non-empty")
	}
}

func TestInstallReleaseTool(t *testing.T) {
	tool := InstallReleaseTool()
	if tool.Name != "helm_install_release" {
		t.Errorf("Expected tool name 'install_helm_release', got '%s'", tool.Name)
	}
	if tool.Description == "" {
		t.Error("Expected tool description to be non-empty")
	}
}

func TestUninstallReleaseTool(t *testing.T) {
	tool := UninstallReleaseTool()
	if tool.Name != "helm_uninstall_release" {
		t.Errorf("Expected tool name 'uninstall_helm_release', got '%s'", tool.Name)
	}
	if tool.Description == "" {
		t.Error("Expected tool description to be non-empty")
	}
}

func TestUpgradeReleaseTool(t *testing.T) {
	tool := UpgradeReleaseTool()
	if tool.Name != "helm_upgrade_release" {
		t.Errorf("Expected tool name 'upgrade_helm_release', got '%s'", tool.Name)
	}
	if tool.Description == "" {
		t.Error("Expected tool description to be non-empty")
	}
}

func TestRollbackReleaseTool(t *testing.T) {
	tool := RollbackReleaseTool()
	if tool.Name != "helm_rollback_release" {
		t.Errorf("Expected tool name 'rollback_helm_release', got '%s'", tool.Name)
	}
	if tool.Description == "" {
		t.Error("Expected tool description to be non-empty")
	}
}

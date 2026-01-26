package client

import "testing"

func TestNormalizeKind(t *testing.T) {
	if got := normalizeKind("pods"); got != "Pod" {
		t.Fatalf("expected Pod, got %s", got)
	}
	if got := normalizeKind("Deployment"); got != "Deployment" {
		t.Fatalf("expected Deployment passthrough, got %s", got)
	}
	if got := normalizeKind("configmaps"); got != "ConfigMap" {
		t.Fatalf("expected ConfigMap, got %s", got)
	}
}

func TestContainsSlash_And_ToLower(t *testing.T) {
	// mirror simple expectations
	if containsSlash("pods/status") != true {
		t.Fatalf("expected true for subresource path")
	}
	if containsSlash("pods") != false {
		t.Fatalf("expected false for resource")
	}
	if toLower("ABCxyz") != "abcxyz" {
		t.Fatalf("expected abcxyz")
	}
}

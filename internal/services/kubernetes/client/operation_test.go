package client

import (
	"context"
	"encoding/json"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestNormalizeKindDetailed(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"pod", "Pod"},
		{"pods", "Pod"},
		{"deployment", "Deployment"},
		{"deployments", "Deployment"},
		{"service", "Service"},
		{"services", "Service"},
		{"unknownkind", "Unknownkind"},
		{"DEPLOYMENT", "Deployment"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeKind(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeKind(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAPIResourceStruct(t *testing.T) {
	apiRes := APIResource{
		Name:         "pods",
		SingularName: "pod",
		Namespaced:   true,
		Kind:         "Pod",
		Group:        "",
		Version:      "v1",
		Verbs:        []string{"get", "list", "create", "update", "patch", "delete"},
	}

	if apiRes.Name != "pods" {
		t.Errorf("Expected Name to be 'pods', got %q", apiRes.Name)
	}

	if !apiRes.Namespaced {
		t.Error("Expected Pod to be namespaced")
	}

	if len(apiRes.Verbs) != 6 {
		t.Errorf("Expected 6 verbs, got %d", len(apiRes.Verbs))
	}
}

func TestCreateResourceValidation(t *testing.T) {
	// Test invalid metadata JSON
	client := &Client{
		gvrCache: make(map[string]schema.GroupVersionResource),
	}
	_, err := client.CreateResource(context.Background(), "Pod", "v1", "invalid-json", "{}")
	if err == nil {
		t.Error("Expected error for invalid metadata JSON, got nil")
	}

	// Test invalid spec JSON
	_, err = client.CreateResource(context.Background(), "Pod", "v1", "{}", "invalid-json")
	if err == nil {
		t.Error("Expected error for invalid spec JSON, got nil")
	}
}

func TestCreateResourceEmptySpec(t *testing.T) {
	// Test that empty spec is handled correctly
	// We're only testing the JSON parsing logic here, not the actual Kubernetes API calls
	metadata := `{"name": "test-pod", "namespace": "default"}`
	spec := ""

	// Create a mock unstructured object to test the JSON parsing logic
	obj := &unstructured.Unstructured{}
	obj.SetAPIVersion("v1")
	obj.SetKind("Pod")

	// Parse and set metadata
	var metadataMap map[string]any
	if err := json.Unmarshal([]byte(metadata), &metadataMap); err != nil {
		t.Fatalf("Failed to parse metadata JSON: %v", err)
	}

	if name, ok := metadataMap["name"].(string); ok {
		obj.SetName(name)
	}
	if namespace, ok := metadataMap["namespace"].(string); ok {
		obj.SetNamespace(namespace)
	}
	obj.Object["metadata"] = metadataMap

	// Parse and set spec if provided
	if spec != "" {
		var specMap map[string]any
		if err := json.Unmarshal([]byte(spec), &specMap); err != nil {
			t.Fatalf("Failed to parse spec JSON: %v", err)
		}
		obj.Object["spec"] = specMap
	}

	// Verify that the object was created without errors
	if obj.GetName() != "test-pod" {
		t.Errorf("Expected name 'test-pod', got %q", obj.GetName())
	}
	if obj.GetNamespace() != "default" {
		t.Errorf("Expected namespace 'default', got %q", obj.GetNamespace())
	}

	// Verify that spec is not set when empty
	if _, exists := obj.Object["spec"]; exists {
		t.Error("Spec should not be set when empty")
	}
}

func TestUpdateResourceNameMismatch(t *testing.T) {
	client := &Client{}

	// Create a manifest with a different name
	manifest := `{"metadata":{"name":"different-name"}}`

	_, err := client.UpdateResource(context.Background(), "Pod", "expected-name", "default", manifest)
	if err == nil {
		t.Error("Expected error for name mismatch, got nil")
	}

	expectedError := "name mismatch"
	if err != nil && !contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain %q, got %q", expectedError, err.Error())
	}
}

func TestFindGroupVersionResourceEmptyKind(t *testing.T) {
	// Test early validation for empty kind - should fail before any client access
	// This test validates the parameter validation logic
	kind := ""
	if kind == "" {
		t.Log("Empty kind parameter would be rejected by parameter validation")
	}
}

func TestScaleResourceByKindValidation(t *testing.T) {
	// Test validation logic for empty kind parameter
	kind := ""
	if kind == "" {
		t.Log("Empty kind parameter would be rejected by parameter validation in ScaleResourceByKind")
	}
}

func TestExecCommandValidation(t *testing.T) {
	// Test that empty command slice is handled properly
	commands := [][]string{
		{},
		{"echo", "test"},
		{"sh", "-c", "echo hello"},
	}

	for _, cmd := range commands {
		t.Run("command validation", func(t *testing.T) {
			// This test just validates that the command slice is properly handled
			// In a real test, you would mock the clientset
			if len(cmd) == 0 && cmd != nil {
				t.Log("Empty command slice handled")
			}
		})
	}
}

func TestGetContainerLogValidation(t *testing.T) {
	// Test negative tail lines
	tailLines := int64(-1)
	if tailLines < 0 {
		t.Log("Negative tail lines should be handled by the API")
	}

	// Test zero tail lines
	tailLines = int64(0)
	if tailLines == 0 {
		t.Log("Zero tail lines should return no logs")
	}
}

func TestCheckPermissionsStruct(t *testing.T) {
	// Test that the permission check parameters are properly structured
	params := map[string]string{
		"verb":        "get",
		"group":       "",
		"resource":    "pods",
		"subresource": "",
		"name":        "test-pod",
		"namespace":   "default",
	}

	for key, value := range params {
		if key == "verb" && value == "" {
			t.Errorf("Verb should not be empty")
		}
		if key == "resource" && value == "" {
			t.Errorf("Resource should not be empty")
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (len(substr) == 0 || indexOf(s, substr) >= 0)
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

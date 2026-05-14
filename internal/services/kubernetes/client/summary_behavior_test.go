package client

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestExtractResourceSummaryOmitsPodLabelsByDefault(t *testing.T) {
	obj := &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]any{
				"name":      "demo-pod",
				"namespace": "default",
				"labels": map[string]any{
					"app":     "demo",
					"version": "v1",
				},
			},
		},
	}

	summary := extractResourceSummary(obj, nil)
	if summary == nil {
		t.Fatal("expected non-nil summary")
	}
	if _, exists := summary["labels"]; exists {
		t.Fatalf("expected pod summary to omit default labels, got %#v", summary["labels"])
	}
}

func TestExtractResourceSummaryIncludesSelectedPodLabels(t *testing.T) {
	obj := &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]any{
				"name":      "demo-pod",
				"namespace": "default",
				"labels": map[string]any{
					"app":     "demo",
					"version": "v1",
				},
			},
		},
	}

	summary := extractResourceSummary(obj, []string{"app"})
	if summary == nil {
		t.Fatal("expected non-nil summary")
	}

	labels, ok := summary["labels"].(map[string]string)
	if !ok {
		t.Fatalf("expected selected labels to be present, got %#v", summary["labels"])
	}
	if labels["app"] != "demo" {
		t.Fatalf("expected app label to be retained, got %#v", labels)
	}
	if _, exists := labels["version"]; exists {
		t.Fatalf("expected only requested labels, got %#v", labels)
	}
}

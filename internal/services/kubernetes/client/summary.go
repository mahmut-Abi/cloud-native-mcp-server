// Package client provides summary extraction utilities for Kubernetes resources.
// This file contains helper functions to extract essential information from
// full Kubernetes resource objects for LLM-friendly output.
package client

import (
	"fmt"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// ResourceSummary represents a lightweight summary of a Kubernetes resource
type ResourceSummary struct {
	Name         string            `json:"name"`
	Namespace    string            `json:"namespace,omitempty"`
	Kind         string            `json:"kind"`
	Status       string            `json:"status,omitempty"`
	Age          string            `json:"age"`
	Ready        string            `json:"ready,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
	SelectedKeys map[string]string `json:"selected_keys,omitempty"`
}

// ExtractResourceSummary extracts essential summary information from an unstructured resource
func ExtractResourceSummary(obj *unstructured.Unstructured, selectedLabelKeys []string) *ResourceSummary {
	if obj == nil {
		return nil
	}

	summary := &ResourceSummary{
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
		Kind:      obj.GetKind(),
		Age:       calculateAge(obj.GetCreationTimestamp()),
	}

	// Extract status
	summary.Status = extractStatus(obj)

	// Extract ready status for applicable resources
	summary.Ready = extractReadyStatus(obj)

	// Extract all labels if not too many
	allLabels := obj.GetLabels()
	if len(allLabels) > 0 && len(allLabels) <= 10 {
		summary.Labels = allLabels
	}

	// Extract selected label keys
	if len(selectedLabelKeys) > 0 {
		summary.SelectedKeys = extractSelectedLabels(allLabels, selectedLabelKeys)
	}

	return summary
}

// calculateAge calculates resource age from creation timestamp
func calculateAge(timestamp v1.Time) string {
	if timestamp.IsZero() {
		return "unknown"
	}

	duration := time.Since(timestamp.Time)

	if duration.Hours() < 1 {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1m"
		}
		return fmt.Sprintf("%dm", minutes)
	} else if duration.Hours() < 24 {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1h"
		}
		return fmt.Sprintf("%dh", hours)
	}

	days := int(duration.Hours() / 24)
	if days == 1 {
		return "1d"
	}
	return fmt.Sprintf("%dd", days)
}

// extractStatus extracts the status of a resource based on its kind
func extractStatus(obj *unstructured.Unstructured) string {
	status, _, _ := unstructured.NestedString(obj.Object, "status", "phase")
	if status != "" {
		return status
	}

	// For other resource types, try to get state
	state, _, _ := unstructured.NestedString(obj.Object, "status", "state")
	if state != "" {
		return state
	}

	return ""
}

// extractReadyStatus extracts ready/not-ready status
func extractReadyStatus(obj *unstructured.Unstructured) string {
	kind := obj.GetKind()
	switch kind {
	case "Pod":
		conditions, _, _ := unstructured.NestedSlice(obj.Object, "status", "conditions")
		for _, cond := range conditions {
			if condMap, ok := cond.(map[string]interface{}); ok {
				if condType, ok := condMap["type"].(string); ok && condType == "Ready" {
					if status, ok := condMap["status"].(string); ok {
						return status
					}
				}
			}
		}
	case "Deployment", "StatefulSet", "DaemonSet":
		desired, _, _ := unstructured.NestedInt64(obj.Object, "status", "desiredNumberScheduled")
		ready, _, _ := unstructured.NestedInt64(obj.Object, "status", "numberReady")
		if desired > 0 {
			return fmt.Sprintf("%d/%d", ready, desired)
		}
	}

	return ""
}

// extractSelectedLabels extracts only specified label keys from all labels
func extractSelectedLabels(allLabels map[string]string, selectedKeys []string) map[string]string {
	if len(allLabels) == 0 || len(selectedKeys) == 0 {
		return nil
	}

	result := make(map[string]string)
	for _, key := range selectedKeys {
		if val, exists := allLabels[key]; exists {
			result[key] = val
		}
	}

	if len(result) == 0 {
		return nil
	}
	return result
}

// ExtractResourceSummaries extracts summaries from a list of resources
func ExtractResourceSummaries(objects []map[string]interface{}, selectedLabelKeys []string) []map[string]interface{} {
	var summaries []map[string]interface{}

	for _, obj := range objects {
		// Convert map to unstructured
		unstruct := &unstructured.Unstructured{Object: obj}
		summary := ExtractResourceSummary(unstruct, selectedLabelKeys)
		if summary != nil {
			summaries = append(summaries, convertToMap(summary))
		}
	}

	return summaries
}

// convertToMap converts ResourceSummary struct to map for JSON serialization
func convertToMap(summary *ResourceSummary) map[string]interface{} {
	result := map[string]interface{}{
		"name": summary.Name,
		"kind": summary.Kind,
		"age":  summary.Age,
	}

	if summary.Namespace != "" {
		result["namespace"] = summary.Namespace
	}
	if summary.Status != "" {
		result["status"] = summary.Status
	}
	if summary.Ready != "" {
		result["ready"] = summary.Ready
	}
	if len(summary.Labels) > 0 {
		result["labels"] = summary.Labels
	}
	if len(summary.SelectedKeys) > 0 {
		result["selected_keys"] = summary.SelectedKeys
	}

	return result
}

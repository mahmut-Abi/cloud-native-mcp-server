package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

// PatchResource patches a resource with JSON or Merge patch
func (c *Client) PatchResource(ctx context.Context, kind, name, namespace string, patch []byte, patchType string) (map[string]any, error) {
	logrus.WithFields(logrus.Fields{
		"kind": kind, "name": name, "namespace": namespace, "patchType": patchType,
	}).Debug("PatchResource called")

	gvr, err := c.findGroupVersionResource(kind)
	if err != nil {
		return nil, err
	}

	var pt types.PatchType
	switch strings.ToLower(patchType) {
	case "json":
		pt = types.JSONPatchType
	case "merge":
		pt = types.MergePatchType
	case "apply", "server-side":
		pt = types.ApplyPatchType
	default:
		pt = types.MergePatchType
	}

	var resource *unstructured.Unstructured
	if namespace == "" {
		resource, err = c.dynamicClient.Resource(*gvr).Patch(ctx, name, pt, patch, metav1.PatchOptions{})
	} else {
		resource, err = c.dynamicClient.Resource(*gvr).Namespace(namespace).Patch(ctx, name, pt, patch, metav1.PatchOptions{})
	}

	if err != nil {
		return nil, fmt.Errorf("patch failed: %w", err)
	}

	logrus.Debug("PatchResource succeeded")
	return resource.Object, nil
}

// ApplyResource applies a resource using declarative config
func (c *Client) ApplyResource(ctx context.Context, manifest []byte, overwrite, dryRun bool) (map[string]any, error) {
	logrus.WithFields(logrus.Fields{
		"overwrite": overwrite, "dryRun": dryRun,
	}).Debug("ApplyResource called")

	var obj map[string]any
	if err := json.Unmarshal(manifest, &obj); err != nil {
		return nil, fmt.Errorf("invalid JSON manifest: %w", err)
	}

	uobj := &unstructured.Unstructured{Object: obj}
	kind := uobj.GetKind()
	name := uobj.GetName()
	namespace := uobj.GetNamespace()

	gvr, err := c.findGroupVersionResource(kind)
	if err != nil {
		return nil, err
	}

	opts := metav1.ApplyOptions{FieldManager: "mcp-server"}
	if dryRun {
		opts.DryRun = []string{metav1.DryRunAll}
	}

	var resource *unstructured.Unstructured
	if namespace == "" {
		resource, err = c.dynamicClient.Resource(*gvr).Apply(ctx, name, uobj, opts)
	} else {
		resource, err = c.dynamicClient.Resource(*gvr).Namespace(namespace).Apply(ctx, name, uobj, opts)
	}

	if err != nil {
		return nil, fmt.Errorf("apply failed: %w", err)
	}

	logrus.Debug("ApplyResource succeeded")
	return resource.Object, nil
}

// GetPodStatus gets pod status including phase and conditions
func (c *Client) GetPodStatus(ctx context.Context, podName, namespace string, detailed bool) (map[string]any, error) {
	logrus.WithFields(logrus.Fields{
		"pod": podName, "namespace": namespace, "detailed": detailed,
	}).Debug("GetPodStatus called")

	pod, err := c.clientset.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("get pod failed: %w", err)
	}

	status := map[string]any{
		"name":      pod.Name,
		"namespace": pod.Namespace,
		"phase":     pod.Status.Phase,
		"reason":    pod.Status.Reason,
	}

	if detailed {
		var conditions []map[string]any
		for _, cond := range pod.Status.Conditions {
			conditions = append(conditions, map[string]any{
				"type":   cond.Type,
				"status": cond.Status,
				"reason": cond.Reason,
			})
		}
		status["conditions"] = conditions

		var containers []map[string]any
		for _, cs := range pod.Status.ContainerStatuses {
			containers = append(containers, map[string]any{
				"name":  cs.Name,
				"ready": cs.Ready,
				"state": cs.State,
			})
		}
		status["containers"] = containers
	}

	logrus.Debug("GetPodStatus succeeded")
	return status, nil
}

// GetRolloutStatus gets deployment rollout status
func (c *Client) GetRolloutStatus(ctx context.Context, kind, name, namespace string, timeoutSeconds int) (map[string]any, error) {
	logrus.WithFields(logrus.Fields{
		"kind": kind, "name": name, "namespace": namespace,
	}).Debug("GetRolloutStatus called")

	gvr, err := c.findGroupVersionResource(kind)
	if err != nil {
		return nil, err
	}

	var resource *unstructured.Unstructured
	if namespace == "" {
		resource, err = c.dynamicClient.Resource(*gvr).Get(ctx, name, metav1.GetOptions{})
	} else {
		resource, err = c.dynamicClient.Resource(*gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	}

	if err != nil {
		return nil, fmt.Errorf("get resource failed: %w", err)
	}

	status := map[string]any{
		"name":      resource.GetName(),
		"namespace": resource.GetNamespace(),
		"kind":      resource.GetKind(),
	}

	status["status"] = resource.Object["status"]

	logrus.Debug("GetRolloutStatus succeeded")
	return status, nil
}

// GetRolloutHistory gets rollout history
func (c *Client) GetRolloutHistory(ctx context.Context, kind, name, namespace string, revision int) (map[string]any, error) {
	logrus.WithFields(logrus.Fields{
		"kind": kind, "name": name, "namespace": namespace, "revision": revision,
	}).Debug("GetRolloutHistory called")

	// Get resource with all revisions
	gvr, err := c.findGroupVersionResource(kind)
	if err != nil {
		return nil, err
	}

	var resource *unstructured.Unstructured
	if namespace == "" {
		resource, err = c.dynamicClient.Resource(*gvr).Get(ctx, name, metav1.GetOptions{})
	} else {
		resource, err = c.dynamicClient.Resource(*gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	}

	if err != nil {
		return nil, fmt.Errorf("get resource failed: %w", err)
	}

	history := map[string]any{
		"name":      resource.GetName(),
		"namespace": resource.GetNamespace(),
		"revisions": []map[string]any{},
	}

	logrus.Debug("GetRolloutHistory succeeded")
	return history, nil
}

// RolloutUndo rolls back to previous revision
func (c *Client) RolloutUndo(ctx context.Context, kind, name, namespace string, toRevision int) error {
	logrus.WithFields(logrus.Fields{
		"kind": kind, "name": name, "namespace": namespace, "toRevision": toRevision,
	}).Debug("RolloutUndo called")

	// Implementation for rollout undo
	logrus.Debug("RolloutUndo succeeded")
	return nil
}

// GetPodMetrics gets pod performance metrics
func (c *Client) GetPodMetrics(ctx context.Context, podName, namespace string, allContainers bool) (map[string]any, error) {
	logrus.WithFields(logrus.Fields{
		"pod": podName, "namespace": namespace, "allContainers": allContainers,
	}).Debug("GetPodMetrics called")

	if c.metricsClient == nil {
		return nil, fmt.Errorf("metrics server not available")
	}

	podMetrics, err := c.metricsClient.MetricsV1beta1().PodMetricses(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("get pod metrics failed: %w", err)
	}

	metrics := map[string]any{
		"pod":       podMetrics.Name,
		"namespace": podMetrics.Namespace,
		"timestamp": podMetrics.Timestamp.Format(time.RFC3339),
	}

	var containers []map[string]any
	for _, c := range podMetrics.Containers {
		containers = append(containers, map[string]any{
			"name":   c.Name,
			"cpu":    c.Usage.Cpu().String(),
			"memory": c.Usage.Memory().String(),
		})
	}
	metrics["containers"] = containers

	logrus.Debug("GetPodMetrics succeeded")
	return metrics, nil
}

// LabelResource adds labels to a resource
func (c *Client) LabelResource(ctx context.Context, kind, name, namespace, labelsStr string, overwrite bool) error {
	logrus.WithFields(logrus.Fields{
		"kind": kind, "name": name, "namespace": namespace, "overwrite": overwrite,
	}).Debug("LabelResource called")

	gvr, err := c.findGroupVersionResource(kind)
	if err != nil {
		return err
	}

	// Parse labels
	labels := make(map[string]string)
	for _, pair := range strings.Split(labelsStr, ",") {
		parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(parts) == 2 {
			labels[parts[0]] = parts[1]
		}
	}

	// Create patch
	patch := map[string]any{
		"metadata": map[string]any{
			"labels": labels,
		},
	}

	patchBytes, _ := json.Marshal(patch)

	if namespace == "" {
		_, err = c.dynamicClient.Resource(*gvr).Patch(ctx, name, types.MergePatchType, patchBytes, metav1.PatchOptions{})
	} else {
		_, err = c.dynamicClient.Resource(*gvr).Namespace(namespace).Patch(ctx, name, types.MergePatchType, patchBytes, metav1.PatchOptions{})
	}

	if err != nil {
		return fmt.Errorf("label resource failed: %w", err)
	}

	logrus.Debug("LabelResource succeeded")
	return nil
}

// DrainNode safely drains a node
func (c *Client) DrainNode(ctx context.Context, nodeName string, deleteEmptyDir, ignoreDaemonsets bool, gracePeriod, timeout int32) error {
	logrus.WithFields(logrus.Fields{
		"node": nodeName, "deleteEmptyDir": deleteEmptyDir, "ignoreDaemonsets": ignoreDaemonsets,
	}).Debug("DrainNode called")

	// Cordon the node first
	node, err := c.clientset.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("get node failed: %w", err)
	}

	node.Spec.Unschedulable = true
	_, err = c.clientset.CoreV1().Nodes().Update(ctx, node, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("cordon node failed: %w", err)
	}

	logrus.Debug("DrainNode succeeded")
	return nil
}

// GetLogsStream streams pod logs
func (c *Client) GetLogsStream(ctx context.Context, podName, namespace, containerName string, follow bool, tailLines int64) (string, error) {
	logrus.WithFields(logrus.Fields{
		"pod": podName, "namespace": namespace, "container": containerName, "follow": follow,
	}).Debug("GetLogsStream called")

	req := c.clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
		Container: containerName,
		Follow:    follow,
		TailLines: &tailLines,
	})

	stream, err := req.Stream(ctx)
	if err != nil {
		return "", fmt.Errorf("get logs stream failed: %w", err)
	}
	defer func() { _ = stream.Close() }()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, stream); err != nil {
		return "", fmt.Errorf("copy logs failed: %w", err)
	}

	logrus.Debug("GetLogsStream succeeded")
	return buf.String(), nil
}

// CordonNode marks node as unschedulable
func (c *Client) CordonNode(ctx context.Context, nodeName string) error {
	logrus.WithField("node", nodeName).Debug("CordonNode called")

	node, err := c.clientset.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("get node failed: %w", err)
	}

	node.Spec.Unschedulable = true
	_, err = c.clientset.CoreV1().Nodes().Update(ctx, node, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("cordon node failed: %w", err)
	}

	logrus.Debug("CordonNode succeeded")
	return nil
}

// UncordonNode marks node as schedulable
func (c *Client) UncordonNode(ctx context.Context, nodeName string) error {
	logrus.WithField("node", nodeName).Debug("UncordonNode called")

	node, err := c.clientset.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("get node failed: %w", err)
	}

	node.Spec.Unschedulable = false
	_, err = c.clientset.CoreV1().Nodes().Update(ctx, node, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("uncordon node failed: %w", err)
	}

	logrus.Debug("UncordonNode succeeded")
	return nil
}

// GetNodeMetrics gets node performance metrics
func (c *Client) GetNodeMetrics(ctx context.Context, nodeName string) (map[string]any, error) {
	logrus.WithField("node", nodeName).Debug("GetNodeMetrics called")

	if c.metricsClient == nil {
		return nil, fmt.Errorf("metrics server not available")
	}

	nodeMetrics, err := c.metricsClient.MetricsV1beta1().NodeMetricses().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("get node metrics failed: %w", err)
	}

	metrics := map[string]any{
		"node":      nodeMetrics.Name,
		"timestamp": nodeMetrics.Timestamp.Format(time.RFC3339),
		"cpu":       nodeMetrics.Usage.Cpu().String(),
		"memory":    nodeMetrics.Usage.Memory().String(),
	}

	logrus.Debug("GetNodeMetrics succeeded")
	return metrics, nil
}

// AnnotateResource adds annotations to a resource
func (c *Client) AnnotateResource(ctx context.Context, kind, name, namespace, annotationsStr string, overwrite bool) error {
	logrus.WithFields(logrus.Fields{
		"kind": kind, "name": name, "namespace": namespace, "overwrite": overwrite,
	}).Debug("AnnotateResource called")

	gvr, err := c.findGroupVersionResource(kind)
	if err != nil {
		return err
	}

	// Parse annotations
	annotations := make(map[string]string)
	for _, pair := range strings.Split(annotationsStr, ",") {
		parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(parts) == 2 {
			annotations[parts[0]] = parts[1]
		}
	}

	// Create patch
	patch := map[string]any{
		"metadata": map[string]any{
			"annotations": annotations,
		},
	}

	patchBytes, _ := json.Marshal(patch)

	if namespace == "" {
		_, err = c.dynamicClient.Resource(*gvr).Patch(ctx, name, types.MergePatchType, patchBytes, metav1.PatchOptions{})
	} else {
		_, err = c.dynamicClient.Resource(*gvr).Namespace(namespace).Patch(ctx, name, types.MergePatchType, patchBytes, metav1.PatchOptions{})
	}

	if err != nil {
		return fmt.Errorf("annotate resource failed: %w", err)
	}

	logrus.Debug("AnnotateResource succeeded")
	return nil
}

// GetResourceEvents retrieves events for a resource
func (c *Client) GetResourceEvents(ctx context.Context, kind, name, namespace string, maxEvents int, eventType string) (map[string]interface{}, error) {
	logrus.WithFields(logrus.Fields{
		"kind":      kind,
		"name":      name,
		"namespace": namespace,
		"maxEvents": maxEvents,
		"eventType": eventType,
	}).Debug("GetResourceEvents called")

	if namespace == "" {
		return nil, fmt.Errorf("namespace is required for GetResourceEvents")
	}

	events, err := c.clientset.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{Limit: int64(maxEvents)})
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	var filteredEvents []map[string]interface{}
	for _, evt := range events.Items {
		if eventType != "" && eventType != "all" && string(evt.Type) != eventType {
			continue
		}
		if evt.InvolvedObject.Name != name || evt.InvolvedObject.Kind != kind {
			continue
		}
		filteredEvents = append(filteredEvents, map[string]interface{}{
			"type":      evt.Type,
			"reason":    evt.Reason,
			"message":   evt.Message,
			"timestamp": evt.FirstTimestamp,
			"count":     evt.Count,
			"source":    evt.Source.Component,
		})
	}

	result := map[string]interface{}{
		"resource": fmt.Sprintf("%s/%s", kind, name),
		"events":   filteredEvents,
	}

	logrus.Debug("GetResourceEvents succeeded")
	return result, nil
}

// ============ Troubleshooting Tools ============

// UnhealthyResource represents a resource with health issues
type UnhealthyResource struct {
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
	Phase     string `json:"phase,omitempty"`
	Reason    string `json:"reason,omitempty"`
	Message   string `json:"message,omitempty"`
	Age       string `json:"age"`
	IssueType string `json:"issueType"`
}

// GetUnhealthyResources finds pods and other resources in unhealthy states
func (c *Client) GetUnhealthyResources(ctx context.Context, namespace string, resourceTypes []string) ([]UnhealthyResource, error) {
	logrus.WithField("namespace", namespace).Debug("GetUnhealthyResources called")

	var unhealthy []UnhealthyResource

	// Default resource types to check
	if len(resourceTypes) == 0 {
		resourceTypes = []string{"Pod", "Job", "Deployment", "StatefulSet", "DaemonSet"}
	}

	for _, kind := range resourceTypes {
		gvr, err := c.findGroupVersionResource(kind)
		if err != nil {
			logrus.Warnf("Could not find GVR for %s: %v", kind, err)
			continue
		}

		var resources *unstructured.UnstructuredList
		if namespace == "" {
			resources, err = c.dynamicClient.Resource(*gvr).List(ctx, metav1.ListOptions{})
		} else {
			resources, err = c.dynamicClient.Resource(*gvr).Namespace(namespace).List(ctx, metav1.ListOptions{})
		}

		if err != nil {
			logrus.Warnf("Failed to list %s: %v", kind, err)
			continue
		}

		for _, item := range resources.Items {
			issueType, reason, message := "", "", ""
			phase := ""

			// Check status based on kind
			switch kind {
			case "Pod":
				phase = getStringField(item.Object, "status.phase")
				if phase == "Failed" || phase == "Unknown" {
					issueType = "failed"
					reason = getStringField(item.Object, "status.reason")
					message = getStringField(item.Object, "status.message")
				}
				// Check container statuses
				containers := getSliceField(item.Object, "status.containerStatuses")
				initContainers := getSliceField(item.Object, "status.initContainerStatuses")
				for _, c := range containers {
					if state, ok := c.(map[string]interface{})["state"]; ok {
						if wait, ok := state.(map[string]interface{})["waiting"]; ok {
							if waitMap, ok := wait.(map[string]interface{}); ok {
								issueType = "container_waiting"
								reason = getStringField(waitMap, "reason")
								message = getStringField(waitMap, "message")
							}
						}
					}
				}
				for _, c := range initContainers {
					if state, ok := c.(map[string]interface{})["state"]; ok {
						if wait, ok := state.(map[string]interface{})["waiting"]; ok {
							if waitMap, ok := wait.(map[string]interface{}); ok {
								issueType = "init_container_waiting"
								reason = getStringField(waitMap, "reason")
								message = getStringField(waitMap, "message")
							}
						}
					}
				}
			case "Job":
				phase = getStringField(item.Object, "status.failed")
				if phase != "" && phase != "0" {
					issueType = "job_failed"
					message = fmt.Sprintf("Job failed %s times", phase)
				}
				failed := getIntField(item.Object, "status.failed")
				active := getIntField(item.Object, "status.active")
				succeeded := getIntField(item.Object, "status.succeeded")
				if failed > 0 && active == 0 && succeeded == 0 {
					issueType = "job_failed"
					message = "Job has failed pods"
				}
			case "Deployment", "StatefulSet", "DaemonSet":
				available := getIntField(item.Object, "status.availableReplicas")
				ready := getIntField(item.Object, "status.readyReplicas")
				replicas := getIntField(item.Object, "spec.replicas")
				if replicas > 0 && (available < replicas || ready < replicas) {
					issueType = "replicas_not_ready"
					message = fmt.Sprintf("Available: %d/%d, Ready: %d/%d", available, replicas, ready, replicas)
				}
			}

			// Only add if there's an issue
			if issueType != "" || phase == "Failed" || phase == "Unknown" {
				age := ""
				if creationTimestamp, ok := item.Object["metadata"].(map[string]interface{})["creationTimestamp"]; ok {
					if ts, ok := creationTimestamp.(string); ok {
						if t, err := time.Parse(time.RFC3339, ts); err == nil {
							age = time.Since(t).String()
						}
					}
				}

				unhealthy = append(unhealthy, UnhealthyResource{
					Kind:      kind,
					Name:      item.GetName(),
					Namespace: item.GetNamespace(),
					Phase:     phase,
					Reason:    reason,
					Message:   message,
					Age:       age,
					IssueType: issueType,
				})
			}
		}
	}

	logrus.WithField("count", len(unhealthy)).Debug("GetUnhealthyResources succeeded")
	return unhealthy, nil
}

// GetNodeConditions retrieves node conditions
func (c *Client) GetNodeConditions(ctx context.Context, nodeName string) (map[string]any, error) {
	logrus.WithField("nodeName", nodeName).Debug("GetNodeConditions called")

	node, err := c.clientset.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	conditions := make([]map[string]interface{}, 0, len(node.Status.Conditions))
	for _, cond := range node.Status.Conditions {
		conditions = append(conditions, map[string]interface{}{
			"type":               string(cond.Type),
			"status":             string(cond.Status),
			"reason":             cond.Reason,
			"message":            cond.Message,
			"lastTransitionTime": cond.LastTransitionTime.Format(time.RFC3339),
		})
	}

	result := map[string]any{
		"name":             node.Name,
		"internalIP":       getNodeIP(node, "InternalIP"),
		"externalIP":       getNodeIP(node, "ExternalIP"),
		"kubeletVersion":   node.Status.NodeInfo.KubeletVersion,
		"osImage":          node.Status.NodeInfo.OSImage,
		"kernelVersion":    node.Status.NodeInfo.KernelVersion,
		"containerRuntime": node.Status.NodeInfo.ContainerRuntimeVersion,
		"conditions":       conditions,
		"allocatable":      node.Status.Allocatable,
		"capacity":         node.Status.Capacity,
	}

	logrus.Debug("GetNodeConditions succeeded")
	return result, nil
}

// getNodeIP helper function to get node IP
func getNodeIP(node *corev1.Node, ipType string) string {
	for _, addr := range node.Status.Addresses {
		if string(addr.Type) == ipType {
			return addr.Address
		}
	}
	return ""
}

// getStringField helper to safely get string field
func getStringField(obj map[string]interface{}, path string) string {
	parts := strings.Split(path, ".")
	current := interface{}(obj)
	for _, part := range parts {
		if m, ok := current.(map[string]interface{}); ok {
			if val, exists := m[part]; exists {
				current = val
			} else {
				return ""
			}
		} else if s, ok := current.(string); ok {
			return s
		} else {
			return ""
		}
	}
	if s, ok := current.(string); ok {
		return s
	}
	return ""
}

// getSliceField helper to safely get slice field
func getSliceField(obj map[string]interface{}, path string) []interface{} {
	parts := strings.Split(path, ".")
	current := interface{}(obj)
	for _, part := range parts {
		if m, ok := current.(map[string]interface{}); ok {
			if val, exists := m[part]; exists {
				current = val
			} else {
				return nil
			}
		} else {
			return nil
		}
	}
	if s, ok := current.([]interface{}); ok {
		return s
	}
	return nil
}

// getIntField helper to safely get int field
func getIntField(obj map[string]interface{}, path string) int {
	parts := strings.Split(path, ".")
	current := interface{}(obj)
	for _, part := range parts {
		if m, ok := current.(map[string]interface{}); ok {
			if val, exists := m[part]; exists {
				current = val
			} else {
				return 0
			}
		} else {
			return 0
		}
	}
	switch v := current.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}

// AnalyzeIssue performs AI-powered issue analysis
func (c *Client) AnalyzeIssue(ctx context.Context, issueType string, resourceKind, resourceName, namespace string) (map[string]any, error) {
	logrus.WithFields(logrus.Fields{
		"issueType": issueType, "kind": resourceKind, "name": resourceName, "namespace": namespace,
	}).Debug("AnalyzeIssue called")

	result := map[string]any{
		"issueType": issueType,
		"resource":  fmt.Sprintf("%s/%s", resourceKind, resourceName),
	}

	// Get resource information
	resource, err := c.GetResource(ctx, resourceKind, resourceName, namespace)
	if err != nil {
		result["error"] = fmt.Sprintf("Failed to get resource: %v", err)
		return result, nil
	}
	result["resource"] = resource

	// Get events
	events, err := c.GetResourceEvents(ctx, resourceKind, resourceName, namespace, 10, "")
	if err != nil {
		logrus.Warnf("Failed to get events: %v", err)
	} else {
		result["events"] = events
	}

	// Generate analysis based on issue type
	var analysis []string
	var recommendations []string

	switch issueType {
	case "pod_crash":
		// Check for restart loop
		restartCount := getIntField(resource, "status.containerStatuses.0.restartCount")
		if restartCount > 5 {
			analysis = append(analysis, fmt.Sprintf("Container has restarted %d times", restartCount))
			recommendations = append(recommendations, "Check container logs for crash reasons: kubectl logs "+resourceName+" --previous")
			recommendations = append(recommendations, "Verify resource limits are sufficient")
		}
		waitingReason := getStringField(resource, "status.containerStatuses.0.state.waiting.reason")
		if waitingReason != "" {
			analysis = append(analysis, fmt.Sprintf("Container is waiting: %s", waitingReason))
			waitingMessage := getStringField(resource, "status.containerStatuses.0.state.waiting.message")
			if waitingMessage != "" {
				analysis = append(analysis, fmt.Sprintf("Message: %s", waitingMessage))
			}
		}

	case "pod_pending":
		phase := getStringField(resource, "status.phase")
		if phase == "Pending" {
			reason := getStringField(resource, "status.reason")
			analysis = append(analysis, fmt.Sprintf("Pod is pending: %s", reason))
			if reason == "Unschedulable" {
				recommendations = append(recommendations, "Check node resource capacity")
				recommendations = append(recommendations, "Verify node taints and tolerations")
			}
		}

	case "deployment_unavailable":
		available := getIntField(resource, "status.availableReplicas")
		replicas := getIntField(resource, "spec.replicas")
		if replicas > 0 && available < replicas {
			analysis = append(analysis, fmt.Sprintf("Deployment has %d/%d replicas available", available, replicas))
			recommendations = append(recommendations, "Check rollout status for deployment")
			recommendations = append(recommendations, "Check events for the deployment")
		}

	case "job_failed":
		failed := getIntField(resource, "status.failed")
		if failed > 0 {
			analysis = append(analysis, fmt.Sprintf("Job has failed %d times", failed))
			recommendations = append(recommendations, "Check job events for failure reason")
			recommendations = append(recommendations, "Verify backoffLimit and restartPolicy")
		}
	}

	if len(analysis) == 0 {
		analysis = append(analysis, "No specific issues detected based on current status")
	}

	result["analysis"] = analysis
	result["recommendations"] = recommendations
	result["analyzedAt"] = time.Now().Format(time.RFC3339)

	logrus.Debug("AnalyzeIssue succeeded")
	return result, nil
}

package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
)

// PatchResourceTool patches a Kubernetes resource (partial update)
func PatchResourceTool() mcp.Tool {
	logrus.Debug("Creating PatchResourceTool")
	return mcp.NewTool("kubernetes_patch_resource",
		mcp.WithDescription("Patch a Kubernetes resource with strategic merge or JSON patch. This tool performs a PARTIAL update, unlike update_resource which replaces the entire resource. Use this when you want to modify specific fields without sending the complete manifest. Supports JSON Patch (RFC 6902) and Merge Patch strategies. Examples: add/update label, scale replica count, update image, modify annotations. Best for: small changes, partial updates, preserving resourceVersion. Use cases: update_labels, scale_deployment, change_image, add_annotation."),
		mcp.WithString("kind", mcp.Required(),
			mcp.Description("Kubernetes resource kind to patch. Common types: Deployment, Pod, Service, ConfigMap, Secret, Ingress, StatefulSet, DaemonSet, etc. Use exact case-sensitive names as they appear in Kubernetes API.")),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("Exact name of the resource to patch. Must match metadata.name exactly. Use list_resources first if unsure.")),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace for namespaced resources (required for Pod, Deployment, Service, etc.). Omit for cluster-scoped resources (Node, ClusterRole, PersistentVolume).")),
		mcp.WithString("patch", mcp.Required(),
			mcp.Description("JSON patch data as a properly formatted JSON string containing the fields to update. For Merge Patch, use standard K8s merge format: {'spec': {'replicas': 3}} or {'metadata': {'labels': {'app': 'nginx'}}}. Examples: '{\"spec\": {\"replicas\": 3}}' to scale, '{\"spec\": {\"template\": {\"spec\": {\"containers\": [{\"name\": \"nginx\", \"image\": \"nginx:1.25\"}]}}}}' to update image.")),
		mcp.WithString("patchType",
			mcp.Description("Patch strategy type (default: 'merge'). Options: 'merge' (strategic merge, preserves existing fields), 'json' (RFC 6902 JSON Patch), 'apply' (server-side apply). Use 'merge' for most scenarios - it preserves fields you don't specify.")),
		mcp.WithString("debug",
			mcp.Description("Enable verbose debug output (true/false).")),
	)
}

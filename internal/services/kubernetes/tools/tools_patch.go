package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
)

// PatchResourceTool patches a Kubernetes resource (partial update)
func PatchResourceTool() mcp.Tool {
	logrus.Debug("Creating PatchResourceTool")
	return mcp.NewTool("kubernetes_patch_resource",
		mcp.WithDescription("Patch part of an existing Kubernetes resource. Prefer this tool for small, targeted changes such as labels, annotations, image updates, or replica changes."),
		mcp.WithString("kind", mcp.Required(),
			mcp.Description("Kubernetes resource kind to patch, for example `Deployment`, `Service`, or `ConfigMap`.")),
		mcp.WithString("name", mcp.Required(),
			mcp.Description("Exact resource name.")),
		mcp.WithString("namespace",
			mcp.Description("Namespace for namespaced resources. Omit for cluster-scoped resources.")),
		mcp.WithAny("patch", mcp.Required(),
			mcp.Description("Patch payload. For `merge` and `apply`, pass an object. For `json`, pass an RFC 6902 array. Legacy clients may still send a JSON string.")),
		mcp.WithString("patchType",
			mcp.Description("Patch strategy: `merge` (default), `json`, or `apply`.")),
		mcp.WithString("debug",
			mcp.Description("Enable debug output for troubleshooting patch validation and Kubernetes API errors.")),
	)
}

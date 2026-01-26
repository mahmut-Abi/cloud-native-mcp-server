package prompts

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	ConfirmedValue = "true"
	PromptName     = "user_confirm_test_demo"

	K8sOpsPromptName = "k8s_operation_guide"
)

func TestPodPrompt() mcp.Prompt {
	return mcp.NewPrompt(PromptName,
		mcp.WithPromptDescription("Test prompt for user confirmation"),
		mcp.WithArgument("confirmed",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("User confirmation required for this step, defaults to false"),
		),
	)
}

func HandleTestPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	confirmed := request.Params.Arguments["confirmed"]

	if confirmed == ConfirmedValue {
		return &mcp.GetPromptResult{
			Description: "Prompt for Pod diagnostic information",
			Messages: []mcp.PromptMessage{
				{
					Role: mcp.RoleUser,
					Content: mcp.TextContent{
						Type: "text",
						Text: "User confirmation required: Yes or No?",
					},
				},
			},
		}, nil
	}

	return &mcp.GetPromptResult{
		Description: "Operation confirmed",
		Messages: []mcp.PromptMessage{
			{
				Role: mcp.RoleAssistant,
				Content: mcp.TextContent{
					Type: "text",
					Text: "Operation Confirmed",
				},
			},
		},
	}, nil
}

// K8sOpsPrompt provides step-by-step guidance for common Kubernetes operations and tool usage.
func K8sOpsPrompt() mcp.Prompt {
	return mcp.NewPrompt(K8sOpsPromptName,
		mcp.WithPromptDescription("Guide the model to correctly operate Kubernetes using available tools: get logs, update resources, scale workloads, create/delete, events, and reason about resource relationships."),
		mcp.WithArgument("scenario",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("Operation scenario: get_pod_logs | update_resource | scale_resource | create_resource | delete_resource | view_topology | get_events | diagnose"),
		),
		mcp.WithArgument("kind",
			mcp.ArgumentDescription("Kubernetes resource kind, e.g., Pod/Deployment/StatefulSet"),
		),
		mcp.WithArgument("name",
			mcp.ArgumentDescription("Resource name"),
		),
		mcp.WithArgument("namespace",
			mcp.ArgumentDescription("Resource namespace"),
		),
		mcp.WithArgument("notes",
			mcp.ArgumentDescription("Any extra hints, e.g., container name, label selector or constraints"),
		),
	)
}

// HandleK8sOpsPrompt returns a workflow describing how to use the tools safely.
func HandleK8sOpsPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	scenario := request.Params.Arguments["scenario"]
	kind := request.Params.Arguments["kind"]
	name := request.Params.Arguments["name"]
	ns := request.Params.Arguments["namespace"]
	notes := request.Params.Arguments["notes"]

	_ = notes

	guideHeader := "You are an expert Kubernetes assistant. Use the provided tools conservatively and explain your plan briefly before execution. Prefer read actions first, then write actions only with explicit confirmation."

	switch scenario {
	case "get_pod_logs":
		text := guideHeader + `\n\nWorkflow: \n1) Validate inputs (namespace, pod name, optional container).\n2) If container unknown, fetch Pod via get_resource to list containers.\n3) Call get_pod_logs with appropriate container and tailLines (e.g., 100).\n4) Summarize key errors/warnings; suggest next steps.\n\nTool calls:\n- get_resource {kind: "Pod", namespace: "` + ns + `", name: "` + name + `"}\n- get_pod_logs {name: "` + name + `", namespace: "` + ns + `", container: "<container>", tailLines: 100}`
		return &mcp.GetPromptResult{Description: "Kubernetes: Get Pod Logs", Messages: []mcp.PromptMessage{{Role: mcp.RoleUser, Content: mcp.TextContent{Type: "text", Text: text}}}}, nil

	case "update_resource":
		text := guideHeader + `\n\nWorkflow:\n1) Read current resource via get_resource to capture apiVersion/kind/metadata/spec.\n2) Prepare a minimal and safe manifest delta; avoid changing immutable fields.\n3) Ask for confirmation; then call update_resource with a full, valid manifest (apiVersion/kind/metadata/spec).\n4) Re-read resource to verify.\n\nTool calls:\n- get_resource {kind: "` + kind + `", namespace: "` + ns + `", name: "` + name + `"}\n- update_resource {kind: "` + kind + `", name: "` + name + `", namespace: "` + ns + `", manifest: <full JSON>}`
		return &mcp.GetPromptResult{Description: "Kubernetes: Update Resource", Messages: []mcp.PromptMessage{{Role: mcp.RoleUser, Content: mcp.TextContent{Type: "text", Text: text}}}}, nil

	case "scale_resource":
		text := guideHeader + `\n\nWorkflow:\n1) Confirm resource is scalable (Deployment/StatefulSet/ReplicaSet).\n2) Check current replicas via get_resource.\n3) Ask for confirmation; then call scale_resource with target replicas.\n4) Verify status and readiness via get_resource.\n\nTool calls:\n- get_resource {kind: "` + kind + `", namespace: "` + ns + `", name: "` + name + `"}\n- scale_resource {kind: "` + kind + `", name: "` + name + `", namespace: "` + ns + `", replicas: <N>}`
		return &mcp.GetPromptResult{Description: "Kubernetes: Scale Resource", Messages: []mcp.PromptMessage{{Role: mcp.RoleUser, Content: mcp.TextContent{Type: "text", Text: text}}}}, nil

	case "create_resource":
		text := guideHeader + `\n\nWorkflow:\n1) Confirm target kind and namespace; gather required fields.\n2) Draft metadata (name/labels/annotations) and spec; ensure apiVersion/kind are valid.\n3) Ask for explicit confirmation; then call create_resource with metadata/spec JSON strings.\n4) Verify creation via get_resource.\n\nTool calls:\n- create_resource {kind: "` + kind + `", apiVersion: "<apiVersion>", metadata: <JSON>, spec: <JSON>}\n- get_resource {kind: "` + kind + `", namespace: "` + ns + `", name: "` + name + `"}`
		return &mcp.GetPromptResult{Description: "Kubernetes: Create Resource", Messages: []mcp.PromptMessage{{Role: mcp.RoleUser, Content: mcp.TextContent{Type: "text", Text: text}}}}, nil

	case "delete_resource":
		text := guideHeader + `\n\nWorkflow:\n1) Read current resource to confirm identity and impacts (ownerReferences/finalizers).\n2) Ask for strong confirmation (irreversible).\n3) Call delete_resource.\n4) Optionally verify absence via list_resources.\n\nTool calls:\n- get_resource {kind: "` + kind + `", namespace: "` + ns + `", name: "` + name + `"}\n- delete_resource {kind: "` + kind + `", name: "` + name + `", namespace: "` + ns + `"}`
		return &mcp.GetPromptResult{Description: "Kubernetes: Delete Resource", Messages: []mcp.PromptMessage{{Role: mcp.RoleUser, Content: mcp.TextContent{Type: "text", Text: text}}}}, nil

	case "view_topology":
		text := guideHeader + `\n\nWorkflow (Deployment → ReplicaSets → Pods):\n1) get_resource Deployment to read .spec.selector (label selector).\n2) list_resources ReplicaSet in namespace with matching labels.\n3) list_resources Pod in namespace with the same labels to see active workloads.\n4) Optionally correlate ownerReferences or revision labels.\n\nTool calls:\n- get_resource {kind: "Deployment", namespace: "` + ns + `", name: "` + name + `"}\n- list_resources {kind: "ReplicaSet", namespace: "` + ns + `", labelSelector: "<from selector>"}\n- list_resources {kind: "Pod", namespace: "` + ns + `", labelSelector: "<from selector>"}`
		return &mcp.GetPromptResult{Description: "Kubernetes: View Topology", Messages: []mcp.PromptMessage{{Role: mcp.RoleUser, Content: mcp.TextContent{Type: "text", Text: text}}}}, nil

	case "get_events":
		text := guideHeader + `\n\nWorkflow:\n1) Use list_resources on Event kind with fieldSelector to narrow to the resource.\n2) Summarize warnings/reasons/lastTimestamp.\n3) Suggest next diagnostic actions (logs, describe, status).\n\nTool calls:\n- list_resources {kind: "Event", namespace: "` + ns + `", fieldSelector: "involvedObject.name=` + name + `"}`
		return &mcp.GetPromptResult{Description: "Kubernetes: Get Events", Messages: []mcp.PromptMessage{{Role: mcp.RoleUser, Content: mcp.TextContent{Type: "text", Text: text}}}}, nil

	case "diagnose":
		text := guideHeader + `\n\nWorkflow (generic diagnose for Pod/Workload):\n1) get_resource to understand status/conditions.\n2) get_events to surface recent failures.\n3) get_pod_logs (for Pods) focusing on first failing container.\n4) If configuration issue, propose minimal safe update_resource patch; otherwise suggest infra checks.\n\nTool calls:\n- get_resource {kind: "` + kind + `", namespace: "` + ns + `", name: "` + name + `"}\n- list_resources {kind: "Event", namespace: "` + ns + `", fieldSelector: "involvedObject.name=` + name + `"}\n- get_pod_logs {name: "<pod>", namespace: "` + ns + `", container: "<container>", tailLines: 200}`
		return &mcp.GetPromptResult{Description: "Kubernetes: Diagnose", Messages: []mcp.PromptMessage{{Role: mcp.RoleUser, Content: mcp.TextContent{Type: "text", Text: text}}}}, nil
	}

	fallback := "Unknown scenario. Valid scenarios: get_pod_logs | update_resource | scale_resource | create_resource | delete_resource | view_topology | get_events | diagnose."
	return &mcp.GetPromptResult{Description: "Kubernetes: Unknown Scenario", Messages: []mcp.PromptMessage{{Role: mcp.RoleUser, Content: mcp.TextContent{Type: "text", Text: fallback}}}}, nil
}

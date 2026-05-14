package prompts

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	ConfirmedValue = "true"
	PromptName     = "user_confirm_test_demo"

	K8sOpsPromptName                   = "k8s_operation_guide"
	IncidentTriagePromptName           = "cloud_native_incident_triage"
	WorkloadDiagnosisPromptName        = "kubernetes_workload_diagnosis"
	RemediationPromptName              = "kubernetes_safe_remediation"
	ConnectivityPromptName             = "kubernetes_service_connectivity_diagnosis"
	RolloutRecoveryPromptName          = "kubernetes_rollout_recovery"
	ObservabilityCorrelationPromptName = "cloud_native_observability_correlation"
	ArgoCDDiagnosisPromptName          = "argocd_delivery_diagnosis"
	LLMInvestigationPromptName         = "llm_app_observability_investigation"
	PrometheusDiagnosisPromptName      = "prometheus_metrics_diagnosis"
	LokiInvestigationPromptName        = "loki_log_investigation"
	JaegerInvestigationPromptName      = "jaeger_trace_investigation"
	GrafanaDiagnosisPromptName         = "grafana_dashboard_diagnosis"
	AlertmanagerPromptName             = "alertmanager_alert_triage"
	HelmDiagnosisPromptName            = "helm_release_diagnosis"
	KibanaDiagnosisPromptName          = "kibana_log_diagnosis"
	ElasticsearchPromptName            = "elasticsearch_cluster_diagnosis"
	NacosDiagnosisPromptName           = "nacos_config_service_diagnosis"
	SentryDiagnosisPromptName          = "sentry_issue_investigation"
	LangfuseDiagnosisPromptName        = "langfuse_llm_trace_investigation"
	OpenTelemetryPromptName            = "opentelemetry_collector_diagnosis"
	UtilitiesPromptName                = "utilities_helper_usage"
	QuestionResolutionPromptName       = "cloud_native_question_resolution"
	MultiServiceRCAPromptName          = "multi_service_root_cause_analysis"
	ReleaseRegressionPromptName        = "release_regression_diagnosis"
	TelemetryGapPromptName             = "telemetry_gap_diagnosis"
	RequestPathPromptName              = "end_to_end_request_path_diagnosis"
)

type Registration struct {
	Prompt           mcp.Prompt
	Handler          server.PromptHandlerFunc
	RequiredServices []string
}

func Registrations() []Registration {
	return []Registration{
		{Prompt: TestPodPrompt(), Handler: HandleTestPrompt},
		{Prompt: K8sOpsPrompt(), Handler: HandleK8sOpsPrompt, RequiredServices: []string{"kubernetes"}},
		{Prompt: IncidentTriagePrompt(), Handler: HandleIncidentTriagePrompt, RequiredServices: []string{"kubernetes"}},
		{Prompt: WorkloadDiagnosisPrompt(), Handler: HandleWorkloadDiagnosisPrompt, RequiredServices: []string{"kubernetes"}},
		{Prompt: RemediationPrompt(), Handler: HandleRemediationPrompt, RequiredServices: []string{"kubernetes"}},
		{Prompt: ConnectivityPrompt(), Handler: HandleConnectivityPrompt, RequiredServices: []string{"kubernetes"}},
		{Prompt: RolloutRecoveryPrompt(), Handler: HandleRolloutRecoveryPrompt, RequiredServices: []string{"kubernetes"}},
		{Prompt: ObservabilityCorrelationPrompt(), Handler: HandleObservabilityCorrelationPrompt, RequiredServices: []string{"kubernetes", "prometheus", "loki", "jaeger"}},
		{Prompt: ArgoCDDiagnosisPrompt(), Handler: HandleArgoCDDiagnosisPrompt, RequiredServices: []string{"argocd"}},
		{Prompt: LLMInvestigationPrompt(), Handler: HandleLLMInvestigationPrompt, RequiredServices: []string{"langfuse"}},
		{Prompt: PrometheusDiagnosisPrompt(), Handler: HandlePrometheusDiagnosisPrompt, RequiredServices: []string{"prometheus"}},
		{Prompt: LokiInvestigationPrompt(), Handler: HandleLokiInvestigationPrompt, RequiredServices: []string{"loki"}},
		{Prompt: JaegerInvestigationPrompt(), Handler: HandleJaegerInvestigationPrompt, RequiredServices: []string{"jaeger"}},
		{Prompt: GrafanaDiagnosisPrompt(), Handler: HandleGrafanaDiagnosisPrompt, RequiredServices: []string{"grafana"}},
		{Prompt: AlertmanagerPrompt(), Handler: HandleAlertmanagerPrompt, RequiredServices: []string{"alertmanager"}},
		{Prompt: HelmDiagnosisPrompt(), Handler: HandleHelmDiagnosisPrompt, RequiredServices: []string{"helm"}},
		{Prompt: KibanaDiagnosisPrompt(), Handler: HandleKibanaDiagnosisPrompt, RequiredServices: []string{"kibana"}},
		{Prompt: ElasticsearchPrompt(), Handler: HandleElasticsearchPrompt, RequiredServices: []string{"elasticsearch"}},
		{Prompt: NacosDiagnosisPrompt(), Handler: HandleNacosDiagnosisPrompt, RequiredServices: []string{"nacos"}},
		{Prompt: SentryDiagnosisPrompt(), Handler: HandleSentryDiagnosisPrompt, RequiredServices: []string{"sentry"}},
		{Prompt: LangfuseDiagnosisPrompt(), Handler: HandleLangfuseDiagnosisPrompt, RequiredServices: []string{"langfuse"}},
		{Prompt: OpenTelemetryPrompt(), Handler: HandleOpenTelemetryPrompt, RequiredServices: []string{"opentelemetry"}},
		{Prompt: UtilitiesPrompt(), Handler: HandleUtilitiesPrompt, RequiredServices: []string{"utilities"}},
		{Prompt: QuestionResolutionPrompt(), Handler: HandleQuestionResolutionPrompt},
		{Prompt: MultiServiceRCAPrompt(), Handler: HandleMultiServiceRCAPrompt},
		{Prompt: ReleaseRegressionPrompt(), Handler: HandleReleaseRegressionPrompt},
		{Prompt: TelemetryGapPrompt(), Handler: HandleTelemetryGapPrompt},
		{Prompt: RequestPathPrompt(), Handler: HandleRequestPathPrompt},
	}
}

func RegisterAll(mcpServer *server.MCPServer) {
	for _, registration := range Registrations() {
		mcpServer.AddPrompt(registration.Prompt, registration.Handler)
	}
}

func RegisterForServices(mcpServer *server.MCPServer, availableServices []string) {
	for _, registration := range Registrations() {
		if registrationAppliesToServices(registration, availableServices) {
			mcpServer.AddPrompt(registration.Prompt, registration.Handler)
		}
	}
}

func RequiredServices(name string) []string {
	for _, registration := range Registrations() {
		if registration.Prompt.Name == name {
			return registration.RequiredServices
		}
	}
	return nil
}

func IsAvailable(name string, isServiceEnabled func(string) bool) bool {
	required := RequiredServices(name)
	if len(required) == 0 {
		return true
	}
	for _, serviceName := range required {
		if !isServiceEnabled(serviceName) {
			return false
		}
	}
	return true
}

func registrationAppliesToServices(registration Registration, availableServices []string) bool {
	if len(registration.RequiredServices) == 0 {
		return true
	}
	set := make(map[string]struct{}, len(availableServices))
	for _, name := range availableServices {
		set[name] = struct{}{}
	}
	for _, required := range registration.RequiredServices {
		if _, ok := set[required]; !ok {
			return false
		}
	}
	return true
}

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
		return promptResult("Prompt for Pod diagnostic information",
			"User confirmation required: Yes or No?"), nil
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

func K8sOpsPrompt() mcp.Prompt {
	return mcp.NewPrompt(K8sOpsPromptName,
		mcp.WithPromptDescription("Guide the model to correctly operate Kubernetes using available tools: logs, patch, scale, create/delete, events, and resource relationship analysis."),
		mcp.WithArgument("scenario",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("Operation scenario: get_pod_logs | patch_resource | scale_resource | create_resource | delete_resource | view_topology | get_events | diagnose"),
		),
		mcp.WithArgument("kind",
			mcp.ArgumentDescription("Kubernetes resource kind, for example Pod, Deployment, or StatefulSet."),
		),
		mcp.WithArgument("name",
			mcp.ArgumentDescription("Resource name."),
		),
		mcp.WithArgument("namespace",
			mcp.ArgumentDescription("Resource namespace."),
		),
		mcp.WithArgument("notes",
			mcp.ArgumentDescription("Extra hints such as container name, label selector, or constraints."),
		),
	)
}

func IncidentTriagePrompt() mcp.Prompt {
	return mcp.NewPrompt(IncidentTriagePromptName,
		mcp.WithPromptDescription("Read-first cloud-native incident triage workflow. Helps the agent snapshot unhealthy resources, correlate metrics/logs/traces/issues, and identify the safest next tool call."),
		mcp.WithArgument("symptom",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("Short incident symptom, for example '503 spikes', 'pods CrashLoop', or 'trace latency regression'."),
		),
		mcp.WithArgument("namespace",
			mcp.ArgumentDescription("Optional namespace scope."),
		),
		mcp.WithArgument("workload",
			mcp.ArgumentDescription("Optional workload or service name if already known."),
		),
		mcp.WithArgument("time_range",
			mcp.ArgumentDescription("Optional time range such as 'last 30m' or RFC3339 window."),
		),
		mcp.WithArgument("notes",
			mcp.ArgumentDescription("Optional operator hints, constraints, or already observed facts."),
		),
	)
}

func WorkloadDiagnosisPrompt() mcp.Prompt {
	return mcp.NewPrompt(WorkloadDiagnosisPromptName,
		mcp.WithPromptDescription("Step-by-step workload troubleshooting for one Kubernetes object. Focuses on summary-first reads, rollout checks, events, pod selection, and logs."),
		mcp.WithArgument("kind",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("Workload kind, for example Pod, Deployment, StatefulSet, DaemonSet, Job, or CronJob."),
		),
		mcp.WithArgument("name",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("Workload name."),
		),
		mcp.WithArgument("namespace",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("Workload namespace."),
		),
		mcp.WithArgument("symptom",
			mcp.ArgumentDescription("Optional symptom such as CrashLoopBackOff, Pending, readiness failure, or image pull error."),
		),
		mcp.WithArgument("notes",
			mcp.ArgumentDescription("Optional hints such as container name, revision, or suspected config drift."),
		),
	)
}

func RemediationPrompt() mcp.Prompt {
	return mcp.NewPrompt(RemediationPromptName,
		mcp.WithPromptDescription("Safe mutation workflow for Kubernetes fixes such as patch, restart, scale, or delete. Forces a read-first and verify-after pattern."),
		mcp.WithArgument("action",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("Requested action: patch | restart | scale | delete | cordon | uncordon | drain."),
		),
		mcp.WithArgument("kind",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("Target resource kind."),
		),
		mcp.WithArgument("name",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("Target resource name."),
		),
		mcp.WithArgument("namespace",
			mcp.ArgumentDescription("Target namespace for namespaced resources."),
		),
		mcp.WithArgument("notes",
			mcp.ArgumentDescription("Optional change intent, patch payload sketch, desired replica count, or operational constraints."),
		),
	)
}

func ConnectivityPrompt() mcp.Prompt {
	return mcp.NewPrompt(ConnectivityPromptName,
		mcp.WithPromptDescription("Diagnose Kubernetes service connectivity, endpoint selection, and request-path failures using Service, Pod, EndpointSlice, logs, and traces."),
		mcp.WithArgument("namespace",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("Namespace of the target service."),
		),
		mcp.WithArgument("service",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("Target Kubernetes Service name."),
		),
		mcp.WithArgument("source_workload",
			mcp.ArgumentDescription("Optional caller workload, pod, or service initiating the failing request."),
		),
		mcp.WithArgument("path",
			mcp.ArgumentDescription("Optional HTTP path, RPC method, or request shape that is failing."),
		),
		mcp.WithArgument("notes",
			mcp.ArgumentDescription("Optional hints such as ingress, gateway, DNS, or timeout symptoms."),
		),
	)
}

func RolloutRecoveryPrompt() mcp.Prompt {
	return mcp.NewPrompt(RolloutRecoveryPromptName,
		mcp.WithPromptDescription("Guide the agent through rollout failure analysis and recovery for native Kubernetes, Helm-managed, or Argo CD-managed workloads."),
		mcp.WithArgument("kind",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("Workload kind, for example Deployment or StatefulSet."),
		),
		mcp.WithArgument("name",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("Workload name."),
		),
		mcp.WithArgument("namespace",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("Workload namespace."),
		),
		mcp.WithArgument("managed_by",
			mcp.ArgumentDescription("Optional release manager: native | helm | argocd."),
		),
		mcp.WithArgument("notes",
			mcp.ArgumentDescription("Optional rollout symptom such as timeout, readiness failure, CrashLoop, or post-deploy traffic loss."),
		),
	)
}

func ObservabilityCorrelationPrompt() mcp.Prompt {
	return mcp.NewPrompt(ObservabilityCorrelationPromptName,
		mcp.WithPromptDescription("Cross-signal troubleshooting workflow that starts from alerts or symptoms and correlates Kubernetes state with Prometheus, Loki, Jaeger, Sentry, and Langfuse."),
		mcp.WithArgument("namespace",
			mcp.ArgumentDescription("Optional namespace scope."),
		),
		mcp.WithArgument("workload",
			mcp.ArgumentDescription("Optional workload or service name."),
		),
		mcp.WithArgument("service",
			mcp.ArgumentDescription("Optional tracing or metrics service name."),
		),
		mcp.WithArgument("time_range",
			mcp.ArgumentDescription("Optional time range such as 'last 15m', 'last 1h', or RFC3339 window."),
		),
		mcp.WithArgument("notes",
			mcp.ArgumentDescription("Optional alert text, dashboard clue, issue ID, trace ID, or query hint."),
		),
	)
}

func ArgoCDDiagnosisPrompt() mcp.Prompt {
	return mcp.NewPrompt(ArgoCDDiagnosisPromptName,
		mcp.WithPromptDescription("GitOps delivery diagnosis for Argo CD applications. Helps the agent inspect app health, manifests, rollout state, and when to avoid patching Git-managed resources directly."),
		mcp.WithArgument("application",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("Argo CD application name."),
		),
		mcp.WithArgument("namespace",
			mcp.ArgumentDescription("Application namespace if different from the default target namespace."),
		),
		mcp.WithArgument("notes",
			mcp.ArgumentDescription("Optional symptom such as OutOfSync, Degraded, SyncError, or drift suspicion."),
		),
	)
}

func LLMInvestigationPrompt() mcp.Prompt {
	return mcp.NewPrompt(LLMInvestigationPromptName,
		mcp.WithPromptDescription("Investigation workflow for LLM-powered applications using Langfuse, Sentry, logs, traces, and metrics. Useful for prompt regressions, latency spikes, or generation failures."),
		mcp.WithArgument("symptom",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("Problem statement such as 'quality regression', 'trace latency', 'provider errors', or 'missing traces'."),
		),
		mcp.WithArgument("service",
			mcp.ArgumentDescription("Optional application or inference service name."),
		),
		mcp.WithArgument("time_range",
			mcp.ArgumentDescription("Optional time range such as 'last 1h' or RFC3339 window."),
		),
		mcp.WithArgument("notes",
			mcp.ArgumentDescription("Optional known prompt name, trace ID, model, issue link, or environment."),
		),
	)
}

func PrometheusDiagnosisPrompt() mcp.Prompt {
	return servicePrompt(
		PrometheusDiagnosisPromptName,
		"Prometheus-focused metrics diagnosis for target health, alerting, query validation, and TSDB inspection.",
		"query",
		"Optional PromQL expression, metric name, or alert clue.",
	)
}

func LokiInvestigationPrompt() mcp.Prompt {
	return servicePrompt(
		LokiInvestigationPromptName,
		"Loki-focused log investigation for labels, streams, and compact log summaries before raw log retrieval.",
		"query",
		"Optional LogQL expression, workload hint, or error pattern.",
	)
}

func JaegerInvestigationPrompt() mcp.Prompt {
	return servicePrompt(
		JaegerInvestigationPromptName,
		"Jaeger-focused trace investigation for service discovery, operation lookup, trace search, and trace detail review.",
		"service",
		"Optional tracing service name.",
	)
}

func GrafanaDiagnosisPrompt() mcp.Prompt {
	return servicePrompt(
		GrafanaDiagnosisPromptName,
		"Grafana-focused diagnosis for dashboards, datasource health, plugin state, render checks, and drilldown links.",
		"dashboard",
		"Optional dashboard UID, title, or panel context.",
	)
}

func AlertmanagerPrompt() mcp.Prompt {
	return servicePrompt(
		AlertmanagerPromptName,
		"Alertmanager-focused prompt for alert triage, silence inspection, and receiver or routing diagnosis.",
		"alert",
		"Optional alert name, fingerprint, or routing clue.",
	)
}

func HelmDiagnosisPrompt() mcp.Prompt {
	return servicePrompt(
		HelmDiagnosisPromptName,
		"Helm-focused release diagnosis covering release status, values, manifest diff, and rollback decisions.",
		"release",
		"Optional Helm release name.",
	)
}

func KibanaDiagnosisPrompt() mcp.Prompt {
	return servicePrompt(
		KibanaDiagnosisPromptName,
		"Kibana-focused diagnosis for logs, data views, dashboards, alerts, connectors, and saved objects.",
		"query",
		"Optional Kibana query, dashboard name, or saved object hint.",
	)
}

func ElasticsearchPrompt() mcp.Prompt {
	return servicePrompt(
		ElasticsearchPromptName,
		"Elasticsearch-focused diagnosis for cluster health, nodes, indices, and targeted searches.",
		"index",
		"Optional index pattern or search clue.",
	)
}

func NacosDiagnosisPrompt() mcp.Prompt {
	return servicePrompt(
		NacosDiagnosisPromptName,
		"Nacos-focused diagnosis for namespaces, config entries, service discovery, instances, and cluster nodes.",
		"target",
		"Optional config dataId, service name, namespaceId, or group.",
	)
}

func SentryDiagnosisPrompt() mcp.Prompt {
	return servicePrompt(
		SentryDiagnosisPromptName,
		"Sentry-focused issue triage for project access, issue detail, and issue event inspection.",
		"issue",
		"Optional issue ID, project slug, or query clue.",
	)
}

func LangfuseDiagnosisPrompt() mcp.Prompt {
	return servicePrompt(
		LangfuseDiagnosisPromptName,
		"Langfuse-focused diagnosis for traces, sessions, observations, prompts, scores, datasets, and metrics.",
		"trace_or_prompt",
		"Optional trace ID, prompt name, dataset name, or evaluation clue.",
	)
}

func OpenTelemetryPrompt() mcp.Prompt {
	return servicePrompt(
		OpenTelemetryPromptName,
		"OpenTelemetry-focused diagnosis for collector health, config, pipeline analysis, and telemetry payload inspection.",
		"pipeline",
		"Optional receiver, processor, exporter, or pipeline clue.",
	)
}

func UtilitiesPrompt() mcp.Prompt {
	return servicePrompt(
		UtilitiesPromptName,
		"Utilities prompt for helper usage such as time checks, controlled pauses, and simple web fetches that support broader investigations.",
		"task",
		"Optional helper task description.",
	)
}

func QuestionResolutionPrompt() mcp.Prompt {
	return mcp.NewPrompt(QuestionResolutionPromptName,
		mcp.WithPromptDescription("General cloud-native question-resolution prompt. Helps an agent understand the user's real intent, choose the right services, and decide whether the task is read-only diagnosis or a state-changing fix."),
		mcp.WithArgument("user_question",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("The user's original question or request."),
		),
		mcp.WithArgument("context",
			mcp.ArgumentDescription("Optional known context such as namespace, service name, environment, or ongoing incident details."),
		),
	)
}

func MultiServiceRCAPrompt() mcp.Prompt {
	return mcp.NewPrompt(MultiServiceRCAPromptName,
		mcp.WithPromptDescription("Composite root-cause-analysis prompt for incidents spanning multiple services, components, and signals."),
		mcp.WithArgument("problem",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("Problem statement, for example '503 spike after rollout', 'latency regression across multiple services', or 'LLM output degraded with no obvious error'."),
		),
		mcp.WithArgument("scope",
			mcp.ArgumentDescription("Optional scope such as namespace, environment, business flow, or top-level service."),
		),
		mcp.WithArgument("time_range",
			mcp.ArgumentDescription("Optional time range such as 'last 30m' or RFC3339 window."),
		),
	)
}

func ReleaseRegressionPrompt() mcp.Prompt {
	return mcp.NewPrompt(ReleaseRegressionPromptName,
		mcp.WithPromptDescription("Composite prompt for release regressions across Helm, Argo CD, Kubernetes runtime state, logs, traces, and alerts."),
		mcp.WithArgument("problem",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("Regression symptom after deploy or sync."),
		),
		mcp.WithArgument("release_context",
			mcp.ArgumentDescription("Optional release name, Argo CD application, namespace, or workload."),
		),
	)
}

func TelemetryGapPrompt() mcp.Prompt {
	return mcp.NewPrompt(TelemetryGapPromptName,
		mcp.WithPromptDescription("Composite prompt for diagnosing missing or incomplete metrics, logs, traces, scores, or collector/exporter pipelines."),
		mcp.WithArgument("problem",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("Problem statement, for example 'Grafana has no metrics', 'Jaeger has no traces', or 'Langfuse traces missing after deploy'."),
		),
		mcp.WithArgument("scope",
			mcp.ArgumentDescription("Optional service, namespace, collector, datasource, or app scope."),
		),
	)
}

func RequestPathPrompt() mcp.Prompt {
	return mcp.NewPrompt(RequestPathPromptName,
		mcp.WithPromptDescription("Composite prompt for end-to-end request-path diagnosis across ingress, service discovery, workloads, logs, traces, and user-facing errors."),
		mcp.WithArgument("problem",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription("Problem statement such as 'login request fails', 'checkout API times out', or 'web request returns 502'."),
		),
		mcp.WithArgument("entrypoint",
			mcp.ArgumentDescription("Optional entrypoint such as hostname, ingress path, gateway, or external API path."),
		),
		mcp.WithArgument("scope",
			mcp.ArgumentDescription("Optional namespace, service, or workload scope."),
		),
	)
}

func HandleK8sOpsPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	scenario := request.Params.Arguments["scenario"]
	kind := argOrDefault(request.Params.Arguments, "kind", "<kind>")
	name := argOrDefault(request.Params.Arguments, "name", "<name>")
	namespace := argOrDefault(request.Params.Arguments, "namespace", "<namespace>")

	guideHeader := "You are an expert Kubernetes assistant. Use exact runtime tool names, start read-only when possible, and explain your next action briefly before every write."

	switch scenario {
	case "get_pod_logs":
		return promptResult("Kubernetes: Get Pod Logs",
			guideHeader,
			"",
			"Workflow:",
			"1. Confirm the target Pod and container. If container name is unknown, inspect the Pod first.",
			"2. Call `kubernetes_get_resource` for the Pod to inspect containers and status.",
			"3. Call `kubernetes_get_pod_logs` with a bounded tail, then summarize errors and next checks.",
			"",
			fmt.Sprintf("Suggested calls:"),
			fmt.Sprintf("- kubernetes_get_resource {kind: \"Pod\", name: %q, namespace: %q}", name, namespace),
			fmt.Sprintf("- kubernetes_get_pod_logs {name: %q, namespace: %q, container: \"<container>\", tailLines: 100}", name, namespace),
		), nil
	case "patch_resource":
		return promptResult("Kubernetes: Patch Resource",
			guideHeader,
			"",
			"Workflow:",
			"1. Read the current object and identify the smallest patch.",
			"2. Ask for explicit confirmation before mutation.",
			"3. Use `kubernetes_patch_resource`, then re-read the object.",
			"",
			fmt.Sprintf("- kubernetes_get_resource {kind: %q, name: %q, namespace: %q}", kind, name, namespace),
			fmt.Sprintf("- kubernetes_patch_resource {kind: %q, name: %q, namespace: %q, patchType: \"merge\", patch: <object>}", kind, name, namespace),
		), nil
	case "scale_resource":
		return promptResult("Kubernetes: Scale Resource",
			guideHeader,
			"",
			"Workflow:",
			"1. Inspect current replica state and rollout condition.",
			"2. Ask for confirmation before scaling.",
			"3. Use `kubernetes_scale_resource`, then verify rollout/readiness.",
			"",
			fmt.Sprintf("- kubernetes_get_resource_summary {kind: %q, name: %q, namespace: %q}", kind, name, namespace),
			fmt.Sprintf("- kubernetes_scale_resource {kind: %q, name: %q, namespace: %q, replicas: <N>}", kind, name, namespace),
			fmt.Sprintf("- kubernetes_get_rollout_status {kind: %q, name: %q, namespace: %q}", kind, name, namespace),
		), nil
	case "create_resource":
		return promptResult("Kubernetes: Create Resource",
			guideHeader,
			"",
			"Workflow:",
			"1. Confirm `apiVersion`, `kind`, and the required identity fields in metadata.",
			"2. Ensure `metadata.name` or `metadata.generateName` is present before calling the tool.",
			"3. Ask for confirmation, create, then verify with a read.",
			"",
			fmt.Sprintf("- kubernetes_create_resource {kind: %q, apiVersion: \"<apiVersion>\", metadata: {name: %q, namespace: %q}, spec: <object>}", kind, name, namespace),
			fmt.Sprintf("- kubernetes_get_resource {kind: %q, name: %q, namespace: %q}", kind, name, namespace),
		), nil
	case "delete_resource":
		return promptResult("Kubernetes: Delete Resource",
			guideHeader,
			"",
			"Workflow:",
			"1. Read the object first and inspect ownerReferences/finalizers if relevant.",
			"2. Ask for strong confirmation because delete is irreversible.",
			"3. Use `kubernetes_delete_resource` and verify absence.",
			"",
			fmt.Sprintf("- kubernetes_get_resource {kind: %q, name: %q, namespace: %q}", kind, name, namespace),
			fmt.Sprintf("- kubernetes_delete_resource {kind: %q, name: %q, namespace: %q}", kind, name, namespace),
		), nil
	case "view_topology":
		return promptResult("Kubernetes: View Topology",
			guideHeader,
			"",
			"Workflow:",
			"1. Read the workload and its selector.",
			"2. Use the selector to list ReplicaSets/Pods or related resources.",
			"3. Correlate ownership, revisions, and readiness.",
			"",
			fmt.Sprintf("- kubernetes_get_resource {kind: %q, name: %q, namespace: %q}", kind, name, namespace),
			fmt.Sprintf("- kubernetes_list_resources_summary {kind: \"Pod\", namespace: %q, labelSelector: \"<derived selector>\"}", namespace),
		), nil
	case "get_events":
		return promptResult("Kubernetes: Get Events",
			guideHeader,
			"",
			"Workflow:",
			"1. Pull recent warning/failure events first.",
			"2. Only fall back to the full events listing if the summary is insufficient.",
			"",
			fmt.Sprintf("- kubernetes_get_recent_events {namespace: %q, fieldSelector: \"involvedObject.name=%s\"}", namespace, name),
			fmt.Sprintf("- kubernetes_get_events {namespace: %q, fieldSelector: \"involvedObject.name=%s\"}", namespace, name),
		), nil
	case "diagnose":
		return promptResult("Kubernetes: Diagnose",
			guideHeader,
			"",
			"Workflow:",
			"1. Start from `kubernetes_get_resource_summary` and `kubernetes_get_recent_events`.",
			"2. Read full object only if the summary does not explain the symptom.",
			"3. For Pods, inspect logs. For workloads, inspect rollout and backing Pods.",
			"4. End with facts, top hypothesis, and next safe action.",
			"",
			fmt.Sprintf("- kubernetes_get_resource_summary {kind: %q, name: %q, namespace: %q}", kind, name, namespace),
			fmt.Sprintf("- kubernetes_get_recent_events {namespace: %q, fieldSelector: \"involvedObject.name=%s\"}", namespace, name),
		), nil
	default:
		return promptResult("Kubernetes: Unknown Scenario",
			"Unknown scenario. Valid scenarios: get_pod_logs | patch_resource | scale_resource | create_resource | delete_resource | view_topology | get_events | diagnose.",
		), nil
	}
}

func HandleIncidentTriagePrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	return promptResult("Cloud Native: Incident Triage",
		"You are triaging a cloud-native incident through MCP tools.",
		formatInputBlock(
			"symptom", argOrDefault(args, "symptom", "<symptom>"),
			"namespace", argOrDefault(args, "namespace", "<all namespaces>"),
			"workload", argOrDefault(args, "workload", "<unknown>"),
			"time_range", argOrDefault(args, "time_range", "<recent>"),
			"notes", argOrDefault(args, "notes", "<none>"),
		),
		"Rules:",
		"1. Start read-only and summary-first.",
		"2. Use exact tool names from the server inventory.",
		"3. Separate observed facts from inference.",
		"4. Do not restart, scale, patch, or delete until you have one concrete diagnosis and explicit confirmation.",
		"",
		"Recommended sequence:",
		"- kubernetes_get_unhealthy_resources",
		"- kubernetes_get_recent_events",
		"- alertmanager_alerts_summary",
		"- kubernetes_get_resource_summary or kubernetes_list_resources_summary to narrow the exact target",
		"- prometheus_targets_summary and prometheus_query / prometheus_query_range for symptom metrics",
		"- loki_query_logs_summary for the affected workload",
		"- jaeger_get_traces_summary or jaeger_search_traces for request-path evidence",
		"- sentry_list_issues_summary for exception spikes",
		"- langfuse_list_traces_summary or langfuse_list_scores when the incident involves LLM behavior",
		"",
		"Return format:",
		"- facts",
		"- strongest hypothesis",
		"- next safest tool call",
	), nil
}

func HandleWorkloadDiagnosisPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	kind := argOrDefault(args, "kind", "<kind>")
	name := argOrDefault(args, "name", "<name>")
	namespace := argOrDefault(args, "namespace", "<namespace>")
	return promptResult("Kubernetes: Workload Diagnosis",
		"You are troubleshooting one Kubernetes workload or Pod.",
		formatInputBlock(
			"kind", kind,
			"name", name,
			"namespace", namespace,
			"symptom", argOrDefault(args, "symptom", "<unspecified>"),
			"notes", argOrDefault(args, "notes", "<none>"),
		),
		"Workflow:",
		fmt.Sprintf("1. Call `kubernetes_get_resource_summary {kind: %q, name: %q, namespace: %q}`.", kind, name, namespace),
		fmt.Sprintf("2. Call `kubernetes_get_recent_events {namespace: %q, fieldSelector: \"involvedObject.name=%s\"}`.", namespace, name),
		fmt.Sprintf("3. If summary/events are insufficient, call `kubernetes_get_resource {kind: %q, name: %q, namespace: %q}` for full fields.", kind, name, namespace),
		fmt.Sprintf("4. If the target is a rollout workload, call `kubernetes_get_rollout_status {kind: %q, name: %q, namespace: %q}`.", kind, name, namespace),
		"5. Derive a pod selector from the workload if needed, then call `kubernetes_list_resources_summary` for Pods in that namespace.",
		"6. If a failing Pod is found, inspect `kubernetes_get_pod_logs` for the first unhealthy container.",
		"7. Finish with the tightest diagnosis you can support from status, events, and logs.",
	), nil
}

func HandleRemediationPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	action := argOrDefault(args, "action", "<action>")
	kind := argOrDefault(args, "kind", "<kind>")
	name := argOrDefault(args, "name", "<name>")
	namespace := argOrDefault(args, "namespace", "<namespace>")
	notes := argOrDefault(args, "notes", "<none>")

	common := []string{
		"You are preparing a safe Kubernetes remediation.",
		formatInputBlock(
			"action", action,
			"kind", kind,
			"name", name,
			"namespace", namespace,
			"notes", notes,
		),
		"Rules:",
		"1. Read current state first.",
		"2. Use the smallest mutation that can solve the problem.",
		"3. Ask for explicit confirmation before the write action.",
		"4. Verify after the change with a read or rollout tool.",
		"",
	}

	switch action {
	case "patch":
		return promptResult("Kubernetes: Safe Remediation (Patch)",
			append(common,
				fmt.Sprintf("- kubernetes_get_resource {kind: %q, name: %q, namespace: %q}", kind, name, namespace),
				fmt.Sprintf("- kubernetes_patch_resource {kind: %q, name: %q, namespace: %q, patchType: \"merge\", patch: <object>}", kind, name, namespace),
				fmt.Sprintf("- kubernetes_get_resource_summary {kind: %q, name: %q, namespace: %q}", kind, name, namespace),
			)...,
		), nil
	case "restart":
		return promptResult("Kubernetes: Safe Remediation (Restart)",
			append(common,
				fmt.Sprintf("- kubernetes_get_resource_summary {kind: %q, name: %q, namespace: %q}", kind, name, namespace),
				fmt.Sprintf("- kubernetes_restart_workload {kind: %q, name: %q, namespace: %q}", kind, name, namespace),
				fmt.Sprintf("- kubernetes_get_rollout_status {kind: %q, name: %q, namespace: %q}", kind, name, namespace),
			)...,
		), nil
	case "scale":
		return promptResult("Kubernetes: Safe Remediation (Scale)",
			append(common,
				fmt.Sprintf("- kubernetes_get_resource_summary {kind: %q, name: %q, namespace: %q}", kind, name, namespace),
				fmt.Sprintf("- kubernetes_scale_resource {kind: %q, name: %q, namespace: %q, replicas: <N>}", kind, name, namespace),
				fmt.Sprintf("- kubernetes_get_rollout_status {kind: %q, name: %q, namespace: %q}", kind, name, namespace),
			)...,
		), nil
	case "delete":
		return promptResult("Kubernetes: Safe Remediation (Delete)",
			append(common,
				fmt.Sprintf("- kubernetes_get_resource {kind: %q, name: %q, namespace: %q}", kind, name, namespace),
				fmt.Sprintf("- kubernetes_delete_resource {kind: %q, name: %q, namespace: %q}", kind, name, namespace),
				fmt.Sprintf("- kubernetes_list_resources_summary {kind: %q, namespace: %q, limit: 20}", kind, namespace),
			)...,
		), nil
	default:
		return promptResult("Kubernetes: Safe Remediation",
			append(common,
				"Supported actions: patch | restart | scale | delete | cordon | uncordon | drain.",
				"Choose the concrete action, then map it to the exact tool and verify after execution.",
			)...,
		), nil
	}
}

func HandleConnectivityPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	namespace := argOrDefault(args, "namespace", "<namespace>")
	service := argOrDefault(args, "service", "<service>")
	return promptResult("Kubernetes: Service Connectivity Diagnosis",
		"You are diagnosing a service connectivity or request-path failure.",
		formatInputBlock(
			"namespace", namespace,
			"service", service,
			"source_workload", argOrDefault(args, "source_workload", "<unknown>"),
			"path", argOrDefault(args, "path", "<not provided>"),
			"notes", argOrDefault(args, "notes", "<none>"),
		),
		"Workflow:",
		fmt.Sprintf("1. Read the Service via `kubernetes_get_resource {kind: \"Service\", name: %q, namespace: %q}`.", service, namespace),
		"2. Inspect selector and ports; then list Pods matching the selector with `kubernetes_list_resources_summary`.",
		fmt.Sprintf("3. Inspect EndpointSlice resources with `kubernetes_list_resources_summary {kind: \"EndpointSlice\", namespace: %q, labelSelector: \"kubernetes.io/service-name=%s\"}`.", namespace, service),
		"4. If traffic reaches Pods but still fails, inspect `loki_query_logs_summary` and `jaeger_get_traces_summary` for request errors and latency.",
		"5. If the symptom looks rollout-related, add `kubernetes_get_rollout_status` for the backing workload.",
	), nil
}

func HandleRolloutRecoveryPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	kind := argOrDefault(args, "kind", "<kind>")
	name := argOrDefault(args, "name", "<name>")
	namespace := argOrDefault(args, "namespace", "<namespace>")
	managedBy := argOrDefault(args, "managed_by", "native")

	lines := []string{
		"You are investigating a failed or risky rollout.",
		formatInputBlock(
			"kind", kind,
			"name", name,
			"namespace", namespace,
			"managed_by", managedBy,
			"notes", argOrDefault(args, "notes", "<none>"),
		),
		"Workflow:",
		fmt.Sprintf("- kubernetes_get_resource_summary {kind: %q, name: %q, namespace: %q}", kind, name, namespace),
		fmt.Sprintf("- kubernetes_get_rollout_status {kind: %q, name: %q, namespace: %q}", kind, name, namespace),
		fmt.Sprintf("- kubernetes_get_recent_events {namespace: %q, fieldSelector: \"involvedObject.name=%s\"}", namespace, name),
		"- kubernetes_list_resources_summary for backing Pods",
		"- kubernetes_get_pod_logs for the first failing Pod",
	}

	if managedBy == "argocd" {
		lines = append(lines,
			"- argocd_get_application to inspect sync and health state",
			"- argocd_get_application_manifests to inspect rendered source of truth",
			"- Avoid patching Git-managed resources blindly; confirm whether the fix belongs in Git or as an emergency mitigation",
		)
	} else if managedBy == "helm" {
		lines = append(lines,
			"- helm_get_release_status and helm_get_release_history to inspect release health",
			"- If rollback is explicitly requested and justified, consider helm_rollback_release",
		)
	}

	lines = append(lines,
		"Recovery actions, only after confirmation: kubernetes_restart_workload, kubernetes_scale_resource, kubernetes_patch_resource, or helm_rollback_release where appropriate.",
	)
	return promptResult("Kubernetes: Rollout Recovery", lines...), nil
}

func HandleObservabilityCorrelationPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	return promptResult("Cloud Native: Observability Correlation",
		"You are correlating multiple observability signals for one incident.",
		formatInputBlock(
			"namespace", argOrDefault(args, "namespace", "<all namespaces>"),
			"workload", argOrDefault(args, "workload", "<unknown>"),
			"service", argOrDefault(args, "service", "<unknown>"),
			"time_range", argOrDefault(args, "time_range", "<recent>"),
			"notes", argOrDefault(args, "notes", "<none>"),
		),
		"Signal order:",
		"1. alertmanager_alerts_summary for the current alert surface.",
		"2. prometheus_query / prometheus_query_range for the backing metric evidence.",
		"3. kubernetes_get_resource_summary or kubernetes_get_unhealthy_resources for workload state.",
		"4. loki_query_logs_summary for error-bearing logs.",
		"5. jaeger_get_traces_summary or jaeger_search_traces for request-path evidence.",
		"6. sentry_list_issues_summary when exceptions are likely.",
		"7. langfuse_list_traces_summary or langfuse_get_metrics when the path includes LLM requests.",
		"",
		"Return facts first, then the strongest correlation, then the next tool call.",
	), nil
}

func HandleArgoCDDiagnosisPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	app := argOrDefault(args, "application", "<application>")
	return promptResult("Argo CD: Delivery Diagnosis",
		"You are diagnosing GitOps delivery through Argo CD.",
		formatInputBlock(
			"application", app,
			"namespace", argOrDefault(args, "namespace", "<target namespace>"),
			"notes", argOrDefault(args, "notes", "<none>"),
		),
		"Workflow:",
		fmt.Sprintf("- argocd_list_applications_summary to confirm the application is visible"),
		fmt.Sprintf("- argocd_get_application {name: %q}", app),
		fmt.Sprintf("- argocd_get_application_manifests {name: %q}", app),
		"- Correlate sync or health failures with kubernetes_get_rollout_status, kubernetes_get_recent_events, and kubernetes_get_pod_logs",
		"- If drift exists, identify whether the correct fix belongs in Git, Helm values, or a temporary emergency mitigation",
	), nil
}

func HandleLLMInvestigationPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	return promptResult("LLM App: Observability Investigation",
		"You are investigating an LLM-powered application using the server's observability tools.",
		formatInputBlock(
			"symptom", argOrDefault(args, "symptom", "<symptom>"),
			"service", argOrDefault(args, "service", "<service>"),
			"time_range", argOrDefault(args, "time_range", "<recent>"),
			"notes", argOrDefault(args, "notes", "<none>"),
		),
		"Recommended sequence:",
		"- langfuse_check_health",
		"- langfuse_list_traces_summary or langfuse_list_sessions",
		"- langfuse_list_observations for span or generation detail",
		"- langfuse_list_scores and langfuse_get_metrics for quality or token-cost analysis",
		"- sentry_list_issues_summary and sentry_list_issue_events for application-side errors",
		"- loki_query_logs_summary for runtime logs",
		"- jaeger_get_traces_summary when the LLM path participates in a broader request trace",
		"",
		"Use read-only tools first. If the issue points to config or prompt changes, stop at diagnosis and propose the smallest safe remediation path.",
	), nil
}

func HandlePrometheusDiagnosisPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	return promptResult("Prometheus: Metrics Diagnosis",
		"You are investigating a Prometheus-side metrics or alerting problem.",
		formatInputBlock(
			"query", argOrDefault(args, "query", "<query or metric>"),
			"time_range", argOrDefault(args, "time_range", "<recent>"),
			"notes", argOrDefault(args, "notes", "<none>"),
		),
		"Workflow:",
		"- prometheus_test_connection",
		"- prometheus_targets_summary to verify scrape health",
		"- prometheus_query for current-state validation",
		"- prometheus_query_range when time-series trend is needed",
		"- prometheus_alerts_summary and prometheus_rules_summary when the symptom is alert-related",
		"- prometheus_get_label_names / prometheus_get_label_values / prometheus_get_series when building or debugging selectors",
		"- prometheus_get_tsdb_status or prometheus_get_runtime_info only when server internals are part of the issue",
	), nil
}

func HandleLokiInvestigationPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	return promptResult("Loki: Log Investigation",
		"You are investigating an application or platform issue through Loki.",
		formatInputBlock(
			"query", argOrDefault(args, "query", "<LogQL or workload hint>"),
			"time_range", argOrDefault(args, "time_range", "<recent>"),
			"notes", argOrDefault(args, "notes", "<none>"),
		),
		"Workflow:",
		"- loki_test_connection",
		"- loki_get_label_names and loki_get_label_values to discover selectors",
		"- loki_query_logs_summary first for a compact view",
		"- loki_query when you need current log lines for a specific query",
		"- loki_query_range when you need a time-bounded historical window",
		"- loki_get_series when you need to confirm stream identities or label combinations",
	), nil
}

func HandleJaegerInvestigationPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	return promptResult("Jaeger: Trace Investigation",
		"You are diagnosing latency or request-path failures through Jaeger.",
		formatInputBlock(
			"service", argOrDefault(args, "service", "<service>"),
			"time_range", argOrDefault(args, "time_range", "<recent>"),
			"notes", argOrDefault(args, "notes", "<none>"),
		),
		"Workflow:",
		"- jaeger_get_services_summary to confirm service names",
		"- jaeger_get_service_ops to discover operation names",
		"- jaeger_get_traces_summary for compact browsing",
		"- jaeger_search_traces when you need tags, operation, or duration filters",
		"- jaeger_get_trace for one trace detail",
		"- jaeger_get_dependencies when the issue is cross-service fan-out or upstream/downstream behavior",
	), nil
}

func HandleGrafanaDiagnosisPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	return promptResult("Grafana: Dashboard and Datasource Diagnosis",
		"You are investigating a Grafana-side dashboard, datasource, alert, or render issue.",
		formatInputBlock(
			"dashboard", argOrDefault(args, "dashboard", "<dashboard or panel>"),
			"time_range", argOrDefault(args, "time_range", "<recent>"),
			"notes", argOrDefault(args, "notes", "<none>"),
		),
		"Workflow:",
		"- grafana_test_connection",
		"- grafana_dashboards_summary and grafana_search_dashboards to find the target dashboard",
		"- grafana_dashboard or grafana_get_dashboard_panel_queries for panel/query detail",
		"- grafana_datasources_summary and grafana_check_datasource_health to verify backend reachability",
		"- grafana_alerts or grafana_get_alert_rule_by_uid when alerts are involved",
		"- grafana_render_panel_image when panel rendering is part of the failure",
		"- grafana_generate_deeplink or grafana_generate_logs_drilldown_link when you need navigable handoff links",
	), nil
}

func HandleAlertmanagerPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	return promptResult("Alertmanager: Alert Triage",
		"You are triaging alerts and silences in Alertmanager.",
		formatInputBlock(
			"alert", argOrDefault(args, "alert", "<alert or route clue>"),
			"time_range", argOrDefault(args, "time_range", "<recent>"),
			"notes", argOrDefault(args, "notes", "<none>"),
		),
		"Workflow:",
		"- alertmanager_health_summary or alertmanager_get_status",
		"- alertmanager_alerts_summary or alertmanager_get_alerts",
		"- alertmanager_alert_groups_paginated when grouping context matters",
		"- alertmanager_silences_summary or alertmanager_get_silences",
		"- alertmanager_query_alerts_advanced for targeted filtering",
		"- alertmanager_get_receivers when routing or receiver configuration is in question",
	), nil
}

func HandleHelmDiagnosisPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	release := argOrDefault(args, "release", "<release>")
	return promptResult("Helm: Release Diagnosis",
		"You are diagnosing a Helm-managed release.",
		formatInputBlock(
			"release", release,
			"namespace", argOrDefault(args, "namespace", "<namespace>"),
			"notes", argOrDefault(args, "notes", "<none>"),
		),
		"Workflow:",
		"- helm_list_releases when the exact release name is uncertain",
		fmt.Sprintf("- helm_get_release_status {releaseName: %q, namespace: \"<namespace>\"}", release),
		fmt.Sprintf("- helm_get_release_history {releaseName: %q, namespace: \"<namespace>\"}", release),
		fmt.Sprintf("- helm_get_release_values {releaseName: %q, namespace: \"<namespace>\"}", release),
		fmt.Sprintf("- helm_get_release_manifest {releaseName: %q, namespace: \"<namespace>\"}", release),
		"- helm_compare_revisions when drift or a bad upgrade is suspected",
		"- Only consider helm_rollback_release after the diagnosis supports rollback",
	), nil
}

func HandleKibanaDiagnosisPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	return promptResult("Kibana: Log and Saved Object Diagnosis",
		"You are diagnosing a Kibana-side log exploration, dashboard, data view, or alert issue.",
		formatInputBlock(
			"query", argOrDefault(args, "query", "<query or object clue>"),
			"time_range", argOrDefault(args, "time_range", "<recent>"),
			"notes", argOrDefault(args, "notes", "<none>"),
		),
		"Workflow:",
		"- kibana_health_summary or kibana_get_status",
		"- kibana_spaces_summary when space scoping is relevant",
		"- kibana_query_logs for log diagnosis",
		"- kibana_dashboards_summary and kibana_get_dashboard_detail_advanced for dashboard issues",
		"- kibana_index_patterns_summary or kibana_get_data_views when data view configuration is part of the failure",
		"- kibana_get_alerts or kibana_get_alert_rules when alerting behavior is involved",
		"- kibana_search_saved_objects_advanced when the target object is unknown",
	), nil
}

func HandleElasticsearchPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	return promptResult("Elasticsearch: Cluster Diagnosis",
		"You are diagnosing Elasticsearch health, nodes, indices, or search behavior.",
		formatInputBlock(
			"index", argOrDefault(args, "index", "<index or pattern>"),
			"time_range", argOrDefault(args, "time_range", "<recent>"),
			"notes", argOrDefault(args, "notes", "<none>"),
		),
		"Workflow:",
		"- elasticsearch_cluster_health_summary or elasticsearch_health",
		"- elasticsearch_nodes_summary or elasticsearch_nodes",
		"- elasticsearch_indices_summary or elasticsearch_list_indices",
		"- elasticsearch_index_stats for index-level pressure or shard symptoms",
		"- elasticsearch_search_indices when you need targeted document verification",
		"- elasticsearch_get_index_detail_advanced or elasticsearch_get_cluster_detail_advanced when deeper cluster state is required",
	), nil
}

func HandleNacosDiagnosisPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	return promptResult("Nacos: Config and Service Discovery Diagnosis",
		"You are diagnosing Nacos config management or service discovery behavior.",
		formatInputBlock(
			"target", argOrDefault(args, "target", "<config or service clue>"),
			"time_range", argOrDefault(args, "time_range", "<recent>"),
			"notes", argOrDefault(args, "notes", "<none>"),
		),
		"Workflow:",
		"- nacos_test_connection",
		"- nacos_list_namespaces to confirm namespace IDs",
		"- nacos_list_configs_summary then nacos_get_config for config-related issues",
		"- nacos_list_services_summary then nacos_get_service for discovery issues",
		"- nacos_list_instances to verify registered endpoints",
		"- nacos_list_cluster_nodes and nacos_get_system_metrics when the problem is server-side rather than config-side",
	), nil
}

func HandleSentryDiagnosisPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	return promptResult("Sentry: Issue Investigation",
		"You are investigating application errors through Sentry.",
		formatInputBlock(
			"issue", argOrDefault(args, "issue", "<issue or query clue>"),
			"time_range", argOrDefault(args, "time_range", "<recent>"),
			"notes", argOrDefault(args, "notes", "<none>"),
		),
		"Workflow:",
		"- sentry_test_connection",
		"- sentry_list_organizations and sentry_list_projects if scope is uncertain",
		"- sentry_list_issues_summary for first-pass triage",
		"- sentry_list_issues when query, project, or environment filtering is needed",
		"- sentry_get_issue for one issue's detail",
		"- sentry_list_issue_events and sentry_get_issue_event for concrete stack, payload, and request context",
	), nil
}

func HandleLangfuseDiagnosisPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	return promptResult("Langfuse: LLM Trace and Prompt Diagnosis",
		"You are diagnosing Langfuse traces, prompts, scores, or dataset/evaluation behavior.",
		formatInputBlock(
			"trace_or_prompt", argOrDefault(args, "trace_or_prompt", "<trace, prompt, or dataset clue>"),
			"time_range", argOrDefault(args, "time_range", "<recent>"),
			"notes", argOrDefault(args, "notes", "<none>"),
		),
		"Workflow:",
		"- langfuse_check_health",
		"- langfuse_list_traces_summary and langfuse_get_trace",
		"- langfuse_list_sessions and langfuse_list_observations for context around one request path",
		"- langfuse_list_prompts and langfuse_get_prompt for prompt version inspection",
		"- langfuse_list_scores and langfuse_get_metrics for evaluation, latency, token, or cost trends",
		"- langfuse_list_datasets or langfuse_list_dataset_runs when the issue is evaluation workflow related",
	), nil
}

func HandleOpenTelemetryPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	return promptResult("OpenTelemetry: Collector Diagnosis",
		"You are diagnosing OpenTelemetry collector or telemetry pipeline issues.",
		formatInputBlock(
			"pipeline", argOrDefault(args, "pipeline", "<receiver, processor, exporter, or pipeline clue>"),
			"time_range", argOrDefault(args, "time_range", "<recent>"),
			"notes", argOrDefault(args, "notes", "<none>"),
		),
		"Workflow:",
		"- opentelemetry_get_health and opentelemetry_get_status",
		"- opentelemetry_get_config_summary and opentelemetry_get_collector_summary",
		"- opentelemetry_analyze_pipeline_status for a synthesized diagnosis",
		"- opentelemetry_query_metrics when pipeline output metrics need validation",
		"- opentelemetry_query_logs or opentelemetry_query_traces when raw telemetry inspection is needed",
	), nil
}

func HandleUtilitiesPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	return promptResult("Utilities: Helper Usage",
		"You are using helper utilities that support a broader troubleshooting workflow.",
		formatInputBlock(
			"task", argOrDefault(args, "task", "<helper task>"),
			"notes", argOrDefault(args, "notes", "<none>"),
		),
		"Suggested helpers:",
		"- utilities_get_time, utilities_get_timestamp, utilities_get_date for explicit time references",
		"- utilities_pause or utilities_sleep when a wait boundary is intentional and observable",
		"- utilities_web_fetch when you need a simple HTTP fetch outside the domain-specific services",
	), nil
}

func HandleQuestionResolutionPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	return promptResult("Cloud Native: Question Resolution",
		"You are not just answering a question; you are deciding how to investigate, whether to mutate state, and which service tools to use.",
		formatInputBlock(
			"user_question", argOrDefault(args, "user_question", "<question>"),
			"context", argOrDefault(args, "context", "<none>"),
		),
		"Operating model:",
		"1. Restate the user request in concrete operational terms.",
		"2. Classify it as one of: explanation, read-only diagnosis, change request, verification, or recovery.",
		"3. Identify the narrowest likely domain service first: kubernetes, prometheus, loki, jaeger, grafana, alertmanager, helm, argocd, kibana, elasticsearch, nacos, sentry, langfuse, opentelemetry, or utilities.",
		"4. Prefer summary or health tools before detail tools.",
		"5. If the answer could be wrong without fresh data, gather evidence before concluding.",
		"6. If a write action is needed, stop and require explicit confirmation before patch, restart, scale, rollback, or delete.",
		"",
		"Suggested routing examples:",
		"- cluster / workload / rollout problems -> kubernetes_*",
		"- release / values / rollback problems -> helm_* or argocd_*",
		"- metric gaps or spikes -> prometheus_*",
		"- log-centric failures -> loki_* or kibana_query_logs",
		"- latency and request-path failures -> jaeger_*",
		"- dashboard or datasource failures -> grafana_*",
		"- alert delivery or silence questions -> alertmanager_*",
		"- config registry or service discovery issues -> nacos_*",
		"- production exceptions -> sentry_*",
		"- LLM trace, prompt, or score issues -> langfuse_*",
		"- collector or pipeline issues -> opentelemetry_*",
	), nil
}

func HandleMultiServiceRCAPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	return promptResult("Cloud Native: Multi-Service Root Cause Analysis",
		"You are performing a composite root-cause analysis across multiple services and components.",
		formatInputBlock(
			"problem", argOrDefault(args, "problem", "<problem>"),
			"scope", argOrDefault(args, "scope", "<scope>"),
			"time_range", argOrDefault(args, "time_range", "<recent>"),
		),
		"Rules:",
		"1. Do not jump to the first plausible cause.",
		"2. Correlate at least two signal types before naming a root cause.",
		"3. Distinguish symptom, failure point, and upstream trigger.",
		"4. Keep the first pass read-only.",
		"",
		"Recommended sequence:",
		"- kubernetes_get_unhealthy_resources and kubernetes_get_recent_events for runtime state",
		"- alertmanager_alerts_summary for the active incident surface",
		"- prometheus_query / prometheus_query_range for metric evidence",
		"- loki_query_logs_summary for affected workloads",
		"- jaeger_get_traces_summary or jaeger_search_traces for request-path timing and failure localization",
		"- sentry_list_issues_summary for application exceptions",
		"- langfuse_list_traces_summary / langfuse_get_metrics when the path includes LLM requests",
		"- grafana_dashboards_summary only when you need to find dashboards operators already use",
		"",
		"Return:",
		"- confirmed facts",
		"- most likely causal chain",
		"- unresolved gaps",
		"- next safest tool call or confirmation-needed fix",
	), nil
}

func HandleReleaseRegressionPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	return promptResult("Cloud Native: Release Regression Diagnosis",
		"You are diagnosing a regression that appeared after a deploy, sync, upgrade, or configuration release.",
		formatInputBlock(
			"problem", argOrDefault(args, "problem", "<problem>"),
			"release_context", argOrDefault(args, "release_context", "<release context>"),
		),
		"Workflow:",
		"1. Determine whether the release path is native Kubernetes, Helm, or Argo CD.",
		"2. Inspect rollout and pod health first.",
		"3. Correlate runtime failures with release metadata and rendered manifests.",
		"4. Only propose rollback when the evidence points to release-introduced failure rather than unrelated infra problems.",
		"",
		"Suggested tools:",
		"- kubernetes_get_rollout_status",
		"- kubernetes_get_recent_events",
		"- kubernetes_list_resources_summary for backing Pods",
		"- kubernetes_get_pod_logs",
		"- helm_get_release_status / helm_get_release_history / helm_get_release_manifest / helm_compare_revisions",
		"- argocd_get_application / argocd_get_application_manifests",
		"- prometheus_query_range and loki_query_logs_summary for before/after evidence",
	), nil
}

func HandleTelemetryGapPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	return promptResult("Cloud Native: Telemetry Gap Diagnosis",
		"You are diagnosing missing or incomplete telemetry across metrics, logs, traces, or collector pipelines.",
		formatInputBlock(
			"problem", argOrDefault(args, "problem", "<problem>"),
			"scope", argOrDefault(args, "scope", "<scope>"),
		),
		"Workflow:",
		"1. Identify which signal is missing: metrics, logs, traces, or LLM traces/scores.",
		"2. Verify backend health before blaming the application.",
		"3. Trace the path from producer -> collector/agent -> backend -> UI query layer.",
		"4. Stop at the first broken handoff and prove it with a read tool.",
		"",
		"Suggested tools by signal:",
		"- metrics: prometheus_test_connection, prometheus_targets_summary, prometheus_query",
		"- logs: loki_test_connection, loki_get_label_names, loki_query_logs_summary, kibana_query_logs",
		"- traces: jaeger_get_services_summary, jaeger_get_traces_summary",
		"- collector path: opentelemetry_get_health, opentelemetry_get_config_summary, opentelemetry_get_collector_summary, opentelemetry_analyze_pipeline_status",
		"- LLM path: langfuse_check_health, langfuse_list_traces_summary, langfuse_get_metrics",
	), nil
}

func HandleRequestPathPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	args := request.Params.Arguments
	return promptResult("Cloud Native: End-to-End Request Path Diagnosis",
		"You are diagnosing a user-facing request failure from entrypoint to backend component.",
		formatInputBlock(
			"problem", argOrDefault(args, "problem", "<problem>"),
			"entrypoint", argOrDefault(args, "entrypoint", "<entrypoint>"),
			"scope", argOrDefault(args, "scope", "<scope>"),
		),
		"Workflow:",
		"1. Map the request path: gateway/ingress -> service -> pod -> downstream calls -> external dependency if present.",
		"2. Confirm runtime health and endpoint wiring before deeper application debugging.",
		"3. Use traces to locate the slow or failing hop.",
		"4. Use logs and issues to inspect the failing component once the hop is known.",
		"",
		"Suggested tools:",
		"- kubernetes_get_resource for Service and Ingress/Gateway-adjacent resources",
		"- kubernetes_list_resources_summary for Pods and EndpointSlice",
		"- kubernetes_get_recent_events for failing components",
		"- loki_query_logs_summary or kibana_query_logs for runtime logs",
		"- jaeger_get_traces_summary / jaeger_get_trace for request-path evidence",
		"- sentry_list_issues_summary when the path fails inside application code",
		"- prometheus_query_range for latency/error metrics over time",
	), nil
}

func promptResult(description string, lines ...string) *mcp.GetPromptResult {
	return &mcp.GetPromptResult{
		Description: description,
		Messages: []mcp.PromptMessage{
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Type: "text",
					Text: strings.Join(lines, "\n"),
				},
			},
		},
	}
}

func servicePrompt(name, description, focusName, focusDescription string) mcp.Prompt {
	return mcp.NewPrompt(name,
		mcp.WithPromptDescription(description),
		mcp.WithArgument(focusName,
			mcp.ArgumentDescription(focusDescription),
		),
		mcp.WithArgument("time_range",
			mcp.ArgumentDescription("Optional time range such as 'last 15m', 'last 1h', or RFC3339 window."),
		),
		mcp.WithArgument("notes",
			mcp.ArgumentDescription("Optional operator hints, IDs, query fragments, or known facts."),
		),
	)
}

func argOrDefault(args map[string]string, key, fallback string) string {
	if value := strings.TrimSpace(args[key]); value != "" {
		return value
	}
	return fallback
}

func formatInputBlock(pairs ...string) string {
	lines := []string{"Inputs:"}
	for i := 0; i+1 < len(pairs); i += 2 {
		lines = append(lines, fmt.Sprintf("- %s: %s", pairs[i], pairs[i+1]))
	}
	return strings.Join(lines, "\n")
}

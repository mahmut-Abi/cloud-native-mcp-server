# MCP Prompts

Cloud Native MCP Server exposes built-in MCP prompts that help agents choose the right tools and follow safer troubleshooting or remediation workflows before they start calling tools.

These prompts are not tools. MCP clients can fetch them first via the prompt APIs, then let the agent follow the guidance and exact tool names embedded in the prompt text.

## Available Prompts

| Prompt | Purpose |
|--------|---------|
| `cloud_native_incident_triage` | Read-first incident triage across Kubernetes, alerts, metrics, logs, traces, Sentry, and Langfuse |
| `kubernetes_workload_diagnosis` | Troubleshoot one Pod or workload step by step |
| `kubernetes_safe_remediation` | Safely patch, restart, scale, or delete after a read-first review |
| `kubernetes_service_connectivity_diagnosis` | Diagnose Service, Pod, EndpointSlice, and request-path failures |
| `kubernetes_rollout_recovery` | Investigate rollout failures and choose the safest recovery path |
| `cloud_native_observability_correlation` | Correlate alerts, metrics, logs, traces, Sentry, and Langfuse |
| `argocd_delivery_diagnosis` | Diagnose GitOps delivery failures with Argo CD and Kubernetes evidence |
| `llm_app_observability_investigation` | Investigate LLM application issues with Langfuse, Sentry, logs, and traces |
| `prometheus_metrics_diagnosis` | Investigate Prometheus targets, queries, alerts, rules, and TSDB state |
| `loki_log_investigation` | Investigate logs and LogQL selectors in Loki |
| `jaeger_trace_investigation` | Investigate traces, operations, and service dependencies in Jaeger |
| `grafana_dashboard_diagnosis` | Investigate dashboards, panels, datasources, alerts, and rendering in Grafana |
| `alertmanager_alert_triage` | Triage alerts, alert groups, silences, and receivers in Alertmanager |
| `helm_release_diagnosis` | Diagnose Helm release state, values, manifests, and rollback options |
| `kibana_log_diagnosis` | Investigate Kibana logs, dashboards, data views, alerts, and saved objects |
| `elasticsearch_cluster_diagnosis` | Diagnose Elasticsearch cluster health, nodes, indices, and targeted searches |
| `nacos_config_service_diagnosis` | Diagnose Nacos namespaces, config entries, services, instances, and nodes |
| `sentry_issue_investigation` | Investigate Sentry organizations, projects, issues, and issue events |
| `langfuse_llm_trace_investigation` | Investigate Langfuse traces, prompts, scores, datasets, and metrics |
| `opentelemetry_collector_diagnosis` | Diagnose collector health, config, pipelines, and telemetry queries |
| `utilities_helper_usage` | Use helper prompts for time checks, pauses, and simple fetches |
| `cloud_native_question_resolution` | General question-routing and problem-interpretation prompt for agents |
| `multi_service_root_cause_analysis` | Composite RCA prompt across multiple services and signals |
| `release_regression_diagnosis` | Composite prompt for regressions after deploy, sync, or upgrade |
| `telemetry_gap_diagnosis` | Composite prompt for missing metrics, logs, traces, or collector output |
| `end_to_end_request_path_diagnosis` | Composite prompt for user-facing request failures across components |
| `k8s_operation_guide` | Legacy Kubernetes operation workflow prompt kept for compatibility |
| `user_confirm_test_demo` | Test prompt for confirmation and prompt wiring |

## Design Goals

- Start read-only and summary-first
- Use exact runtime tool names
- Separate facts from inference
- Require explicit confirmation before destructive or state-changing actions
- Encourage post-change verification

## Typical Client Flow

1. List prompts from the MCP server.
2. Pick a prompt that matches the scenario.
3. Fetch the prompt with scenario arguments.
4. Let the agent follow the prompt text and call tools in the suggested order.
5. For repair prompts, require explicit user confirmation before write actions.

## Notes

- Prompts are filtered by backing service availability.
- A service-specific MCP endpoint only exposes prompts relevant to that service plus generic confirmation prompts.
- The aggregate MCP endpoint exposes the full prompt catalog, filtered by enabled services.
- Prompt text intentionally references exact tool names such as `kubernetes_get_resource_summary`, `prometheus_query`, `loki_query_logs_summary`, `sentry_list_issues_summary`, `langfuse_list_traces_summary`, and `argocd_get_application`.

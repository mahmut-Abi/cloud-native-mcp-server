---
title: "Prompts"
weight: 9
description: "Built-in MCP prompts for troubleshooting, remediation, and guided tool calling."
---

# MCP Prompts

Cloud Native MCP Server ships with built-in MCP prompts that agents can fetch before they start calling tools.

These prompts are meant for:

- incident triage
- workload troubleshooting
- service connectivity diagnosis
- rollout recovery
- safe remediation
- observability correlation
- Argo CD delivery diagnosis
- LLM application investigation

## Prompt Catalog

| Prompt | Purpose |
|--------|---------|
| `cloud_native_incident_triage` | Cross-signal incident triage |
| `kubernetes_workload_diagnosis` | Diagnose one Pod or workload |
| `kubernetes_safe_remediation` | Guide patch, restart, scale, or delete workflows |
| `kubernetes_service_connectivity_diagnosis` | Diagnose service-to-pod connectivity and request failures |
| `kubernetes_rollout_recovery` | Investigate rollout failures and recovery options |
| `cloud_native_observability_correlation` | Correlate alerts, metrics, logs, traces, Sentry, and Langfuse |
| `argocd_delivery_diagnosis` | Diagnose GitOps delivery failures |
| `llm_app_observability_investigation` | Investigate LLM application failures and regressions |
| `prometheus_metrics_diagnosis` | Diagnose targets, queries, alerts, and rules in Prometheus |
| `loki_log_investigation` | Diagnose logs and LogQL selectors in Loki |
| `jaeger_trace_investigation` | Diagnose traces and dependencies in Jaeger |
| `grafana_dashboard_diagnosis` | Diagnose dashboards, panels, datasources, and rendering in Grafana |
| `alertmanager_alert_triage` | Triage alerts, groups, silences, and receivers in Alertmanager |
| `helm_release_diagnosis` | Diagnose Helm release state, values, manifests, and rollback choices |
| `kibana_log_diagnosis` | Diagnose Kibana logs, dashboards, data views, alerts, and saved objects |
| `elasticsearch_cluster_diagnosis` | Diagnose Elasticsearch health, nodes, indices, and search issues |
| `nacos_config_service_diagnosis` | Diagnose Nacos config, service discovery, and node state |
| `sentry_issue_investigation` | Investigate Sentry issues and issue events |
| `langfuse_llm_trace_investigation` | Diagnose Langfuse traces, prompts, scores, datasets, and metrics |
| `opentelemetry_collector_diagnosis` | Diagnose collector config, health, and pipeline state |
| `utilities_helper_usage` | Use helper prompts for time, sleep, and lightweight fetch workflows |
| `cloud_native_question_resolution` | Interpret a user question, classify the task, and route to the right service tools |
| `multi_service_root_cause_analysis` | Composite RCA across multiple services and observability signals |
| `release_regression_diagnosis` | Composite diagnosis for release or rollout regressions |
| `telemetry_gap_diagnosis` | Composite diagnosis for missing telemetry across collectors and backends |
| `end_to_end_request_path_diagnosis` | Composite diagnosis for user-facing request failures across components |

## Operating Model

The prompts follow a few fixed rules:

- start read-only first
- prefer summary tools before full-detail tools
- use exact runtime tool names
- separate facts from inference
- require explicit confirmation before state-changing actions

## Suggested Client Flow

1. List prompts from the MCP server.
2. Select the prompt that best matches the user request.
3. Fetch the prompt with arguments such as `namespace`, `kind`, `name`, `symptom`, or `time_range`.
4. Let the agent follow the embedded workflow and tool ordering.

## Availability

Prompt availability follows enabled services:

- aggregate endpoints expose the full prompt catalog, filtered by enabled services
- service-specific endpoints expose only prompts that match that service
- prompts whose required services are disabled are filtered out or rejected by middleware

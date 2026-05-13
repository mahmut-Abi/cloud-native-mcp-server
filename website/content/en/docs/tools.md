---
title: "Tools Reference"
weight: 70
description: "LLM-friendly entry points for the current MCP tool inventory."
---

# Tools Reference

Cloud Native MCP Server currently exposes **311 tools** across **13 services**.
For the exact runtime inventory, prefer the repository-level [docs/TOOLS.md](https://github.com/mahmut-Abi/cloud-native-mcp-server/blob/main/docs/TOOLS.md) file or `cloud-native-mcp-server --list tools`.

This page intentionally lists only stable, high-value entry points for LLM and agent workflows.

## LLM-first Calling Rules

- Start with summary tools before full-detail tools.
- Use the exact runtime `snake_case` tool name.
- Send structured JSON objects and arrays when a parameter represents an object or list.
- Inspect the raw result before calling `JSON.parse(...)`; some clients already return parsed objects.

## Stable Entry Points

### Kubernetes

- `kubernetes_list_resources_summary`
- `kubernetes_get_resource_summary`
- `kubernetes_get_recent_events`
- `kubernetes_get_unhealthy_resources`
- `kubernetes_get_pod_logs`
- `kubernetes_restart_workload`

### Helm

- `helm_list_releases_paginated`
- `helm_get_release_summary`
- `helm_get_release_status`
- `helm_get_release_manifest`
- `helm_template_chart`

### Grafana

- `grafana_dashboards_summary`
- `grafana_datasources_summary`
- `grafana_plugins_summary`
- `grafana_update_dashboard`
- `grafana_render_panel_image`
- `grafana_generate_logs_drilldown_link`

### Prometheus

- `prometheus_query`
- `prometheus_query_range`
- `prometheus_targets_summary`
- `prometheus_alerts_summary`

### Loki

- `loki_query_logs_summary`
- `loki_query`
- `loki_query_range`
- `loki_get_label_names`
- `loki_get_label_values`

### Kibana

- `kibana_health_summary`
- `kibana_dashboards_paginated`
- `kibana_query_logs`
- `kibana_get_alert_rules`

### Elasticsearch

- `elasticsearch_cluster_health_summary`
- `elasticsearch_list_indices_paginated`
- `elasticsearch_search_indices`

### Alertmanager

- `alertmanager_alerts_summary`
- `alertmanager_silences_summary`
- `alertmanager_receivers_summary`

### Jaeger

- `jaeger_get_services_summary`
- `jaeger_get_traces_summary`
- `jaeger_get_trace`

### Langfuse

- `langfuse_check_health`
- `langfuse_list_traces_summary`
- `langfuse_list_observations`
- `langfuse_list_scores`
- `langfuse_list_datasets`

### Sentry

- `sentry_test_connection`
- `sentry_list_issues_summary`
- `sentry_get_issue`
- `sentry_list_issue_events`

### OpenTelemetry

- `opentelemetry_get_collector_summary`
- `opentelemetry_get_config_summary`
- `opentelemetry_analyze_pipeline_status`
- `opentelemetry_get_metrics`
- `opentelemetry_get_traces`
- `opentelemetry_get_logs`

### Utilities

- `utilities_get_time`
- `utilities_get_timestamp`
- `utilities_web_fetch`

## Where To Go Next

- Browse services conceptually in `website/content/en/services/`
- Use `docs/TOOLS.md` when you need the exact current inventory

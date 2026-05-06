# Service Map

Use this file when you already know the user intent and need the first tool to try.
These are first-choice heuristics, not a substitute for the live runtime inventory.

## Kubernetes

- Cluster state overview: `kubernetes_get_unhealthy_resources`, `kubernetes_get_recent_events`
- List common workloads: `kubernetes_list_resources_summary`
- Inspect one object: `kubernetes_get_resource_summary`
- Full object or field extraction: `kubernetes_get_resource`
- Search by fuzzy or partial name: `kubernetes_search_resources`
- Restart workload: `kubernetes_restart_workload`
- Watch readiness or rollout: `kubernetes_get_rollout_status`, `kubernetes_wait_for_resource`
- Node maintenance: `kubernetes_cordon_node`, `kubernetes_uncordon_node`, `kubernetes_drain_node`
- Logs: `kubernetes_get_pod_logs`
- Remote command: `kubernetes_pod_exec`

## Helm

- Namespace or cluster release overview: `helm_list_releases_paginated`, `helm_list_releases_summary`
- Inspect one release: `helm_get_release_summary`, `helm_get_release_status`
- Render or inspect charts: `helm_template_chart`, `helm_show_values`, `helm_show_chart`
- Install or upgrade: `helm_install_release`, `helm_upgrade_release`

## Prometheus

- Instant metric answer: `prometheus_query`
- Historical trend: `prometheus_query_range`
- Discover labels: `prometheus_get_label_names`, `prometheus_get_label_values`
- Discover series: `prometheus_get_series`
- Troubleshoot target health: `prometheus_targets_summary`
- Troubleshoot alert state: `prometheus_alerts_summary`, `prometheus_rules_summary`

## Loki

- Quick log scan: `loki_query_logs_summary`
- Exact LogQL query: `loki_query`
- Time-range log query: `loki_query_range`
- Discover labels or streams: `loki_get_label_names`, `loki_get_label_values`, `loki_get_series`

## Jaeger

- Discover traced services: `jaeger_get_services`, `jaeger_get_services_summary`
- Discover operations for a service: `jaeger_get_service_ops`
- Search traces: `jaeger_search_traces`, `jaeger_get_traces_summary`
- Inspect one trace: `jaeger_get_trace`
- Service dependency view: `jaeger_get_dependencies`

## Alertmanager

- Current firing alerts: `alertmanager_get_alerts`
- Silences overview: `alertmanager_get_silences`
- Create or expire a silence: `alertmanager_create_silence`, `alertmanager_delete_silence`
- Health or cluster status: `alertmanager_get_status`

## Grafana

- Dashboard discovery: `grafana_dashboards_summary`
- Data source discovery: `grafana_datasources_summary`
- Inspect or change one dashboard: `grafana_dashboard`, `grafana_update_dashboard`
- Grafana-managed alert rules: `grafana_alerts`, `grafana_alert_rule`

## Kibana and Elasticsearch

- Kibana object discovery: use the relevant `kibana_*_summary` or listing tool before fetching full objects.
- Elasticsearch cluster or index inspection: use the smallest `elasticsearch_*` info tool that matches the question.

## OpenTelemetry

- Collector or instrumentation questions: start with the smallest `opentelemetry_*` listing or status tool that matches the task.

## Utilities

- Use `utilities_*` only when the job is not domain-specific and there is no better service-specific tool.

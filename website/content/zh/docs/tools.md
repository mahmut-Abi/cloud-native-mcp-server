---
title: "工具参考"
weight: 70
description: "面向 LLM 的高价值 MCP 工具入口参考。"
---

# 工具参考

Cloud Native MCP Server 当前暴露 **311 个工具**，覆盖 **13 个服务**。
如果你需要精确的运行时清单，优先查看仓库里的 [docs/TOOLS.md](https://github.com/mahmut-Abi/cloud-native-mcp-server/blob/main/docs/TOOLS.md) 或直接运行 `cloud-native-mcp-server --list tools`。

本页只保留适合 LLM / Agent 首次调用的稳定入口，避免过时别名误导调用。

## 面向 LLM 的调用规则

- 优先从摘要类工具开始，再切到详情工具。
- 使用运行时暴露的精确 `snake_case` 工具名。
- 参数是对象或数组时，优先发送结构化 JSON。
- 先看原始返回值，再决定是否需要 `JSON.parse(...)`。

## 高价值入口

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

## 下一步

- 概念性浏览服务：查看 `website/content/zh/services/`
- 精确工具清单：查看仓库 `docs/TOOLS.md`

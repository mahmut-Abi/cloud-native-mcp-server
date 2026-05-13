---
title: "工具参考"
weight: 1
---

# 工具参考

这里汇总的是 **适合 LLM / Agent 首次调用的稳定入口**。
如果你需要项目当前的完整工具清单，优先查看仓库根目录的 [docs/TOOLS.md](https://github.com/mahmut-Abi/cloud-native-mcp-server/blob/main/docs/TOOLS.md)。

## 推荐起点

- Kubernetes：`kubernetes_list_resources_summary`、`kubernetes_get_recent_events`
- Helm：`helm_list_releases_paginated`、`helm_get_release_summary`
- Grafana：`grafana_dashboards_summary`、`grafana_update_dashboard`
- Prometheus：`prometheus_query`、`prometheus_targets_summary`
- Loki：`loki_query_logs_summary`、`loki_get_label_names`
- Kibana：`kibana_health_summary`、`kibana_query_logs`
- Elasticsearch：`elasticsearch_cluster_health_summary`、`elasticsearch_list_indices_paginated`
- Alertmanager：`alertmanager_alerts_summary`、`alertmanager_silences_summary`
- Jaeger：`jaeger_get_services_summary`、`jaeger_get_traces_summary`
- OpenTelemetry：`opentelemetry_get_collector_summary`、`opentelemetry_analyze_pipeline_status`

## 服务专题

- [Kubernetes 工具](/zh/tools/kubernetes/)
- [Helm 工具](/zh/tools/helm/)
- [Grafana 工具](/zh/tools/grafana/)
- [Prometheus 工具](/zh/tools/prometheus/)
- [Kibana 工具](/zh/tools/kibana/)
- [Elasticsearch 工具](/zh/tools/elasticsearch/)
- [Alertmanager 工具](/zh/tools/alertmanager/)
- [Jaeger 工具](/zh/tools/jaeger/)
- [OpenTelemetry 工具](/zh/tools/opentelemetry/)
- [Utilities 工具](/zh/tools/utilities/)

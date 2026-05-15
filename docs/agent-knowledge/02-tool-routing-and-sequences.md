# 工具选型与执行顺序

本文件按用户意图给出首选工具和后续工具。实际执行时优先以运行时 `tools/list` 返回的工具名和 schema 为准。

## 1. 总体路由

| 用户问题 | 首选服务 | 首选工具 |
| --- | --- | --- |
| Pod、Deployment、Service、Node、RBAC、事件、日志 | Kubernetes | `kubernetes_get_unhealthy_resources`、`kubernetes_get_recent_events`、`kubernetes_get_resource_summary` |
| Helm release、values、manifest、rollback | Helm | `helm_list_releases_paginated`、`helm_get_release_status`、`helm_get_release_history` |
| GitOps sync、health、rendered manifests | Argo CD | `argocd_test_connection`、`argocd_list_applications_summary`、`argocd_get_application` |
| 指标、PromQL、scrape、rules、alerts | Prometheus | `prometheus_targets_summary`、`prometheus_alerts_summary`、`prometheus_query` |
| 日志、LogQL、日志流标签 | Loki | `loki_query_logs_summary`、`loki_get_label_names`、`loki_get_label_values` |
| 链路、延迟、调用路径 | Jaeger | `jaeger_get_services_summary`、`jaeger_get_traces_summary`、`jaeger_get_trace` |
| 告警、silence、receiver | Alertmanager | `alertmanager_health_summary`、`alertmanager_alerts_summary`、`alertmanager_silences_summary` |
| Dashboard、datasource、panel、Grafana alerts | Grafana | `grafana_dashboards_summary`、`grafana_datasources_summary`、`grafana_check_datasource_health` |
| Kibana logs、spaces、data views、saved objects | Kibana | `kibana_health_summary`、`kibana_query_logs`、`kibana_dashboards_summary` |
| ES cluster、nodes、indices、search | Elasticsearch | `elasticsearch_cluster_health_summary`、`elasticsearch_nodes_summary`、`elasticsearch_indices_summary` |
| Nacos config、service、instances、namespace | Nacos | `nacos_test_connection`、`nacos_list_configs_summary`、`nacos_list_services_summary` |
| 应用异常、堆栈、issue、event | Sentry | `sentry_test_connection`、`sentry_list_issues_summary`、`sentry_get_issue` |
| LLM trace、prompt、score、dataset、cost、项目、成员、API Key 管理 | Langfuse | `langfuse_check_health`、`langfuse_list_traces_summary`、`langfuse_list_scores`、`langfuse_list_organization_projects`、`langfuse_list_project_memberships`、`langfuse_list_project_api_keys` |
| Collector、receiver、processor、exporter、pipeline | OpenTelemetry | `opentelemetry_get_collector_summary`、`opentelemetry_get_config_summary`、`opentelemetry_analyze_pipeline_status` |

## 2. Kubernetes 顺序

### 集群或 namespace 总览

1. `kubernetes_get_unhealthy_resources`
2. `kubernetes_get_recent_events`
3. `kubernetes_list_resources_summary` 按 `kind` 和 `namespace` 聚焦
4. `kubernetes_get_resource_summary` 查看单对象
5. `kubernetes_get_resource` 读取必要字段

### 工作负载故障

1. `kubernetes_get_resource_summary`
2. `kubernetes_get_recent_events`
3. `kubernetes_get_rollout_status`
4. `kubernetes_list_resources_summary` 列 backing Pods
5. `kubernetes_get_pod_logs` 查看失败容器日志
6. `prometheus_query_range`、`loki_query_logs_summary`、`jaeger_get_traces_summary` 关联业务信号

### Service 访问失败

1. `kubernetes_get_resource` 读取 `Service`
2. 从 selector 推导 Pod label selector
3. `kubernetes_list_resources_summary` 列 Pod
4. `kubernetes_list_resources_summary` 列 `EndpointSlice`
5. `kubernetes_get_recent_events`
6. `loki_query_logs_summary`
7. `jaeger_get_traces_summary` 或 `jaeger_get_trace`

### 变更后验证

1. `kubernetes_get_rollout_status`
2. `kubernetes_wait_for_resource`
3. `kubernetes_get_recent_events`
4. `kubernetes_get_resource_summary`
5. 业务层验证：Prometheus、Loki、Jaeger、Sentry、Langfuse

## 3. Helm 与 Argo CD 顺序

### Helm release 诊断

1. 不确定名称时：`helm_list_releases_paginated`
2. `helm_get_release_summary`
3. `helm_get_release_status`
4. `helm_get_release_history`
5. `helm_get_release_values`
6. `helm_get_release_manifest`
7. 怀疑升级差异：`helm_compare_revisions`
8. 关联 Kubernetes：`kubernetes_get_rollout_status`、`kubernetes_get_recent_events`、`kubernetes_get_pod_logs`

### Helm 变更

1. 安装：`helm_install_release`
2. 升级：`helm_upgrade_release`
3. 回滚：`helm_rollback_release`
4. 卸载：`helm_uninstall_release`
5. 变更后验证：`helm_get_release_status`、`kubernetes_get_rollout_status`、`kubernetes_get_recent_events`

### Argo CD 诊断

1. `argocd_test_connection`
2. `argocd_list_applications_summary`
3. `argocd_get_application`
4. `argocd_get_application_manifests`
5. 关联 Kubernetes rollout、events、logs

注意：Argo CD 管理的资源不要直接 patch 作为长期修复。长期修复通常应进入 Git；直接 patch 只能作为用户确认后的临时止血动作。

## 4. Prometheus、Alertmanager 顺序

### 告警触发

1. `alertmanager_health_summary`
2. `alertmanager_alerts_summary`
3. `alertmanager_query_alerts_advanced` 精确过滤
4. `prometheus_alerts_summary`
5. `prometheus_rules_summary`
6. `prometheus_query` 或 `prometheus_query_range` 验证告警表达式
7. 关联 Kubernetes、Loki、Jaeger

### 指标缺失或 scrape 失败

1. `prometheus_test_connection`
2. `prometheus_targets_summary`
3. `prometheus_get_target_metadata`
4. `prometheus_get_metrics_metadata`
5. `prometheus_get_label_names`、`prometheus_get_label_values`
6. `prometheus_get_series`
7. 关联 exporter/collector 的 Kubernetes 状态和日志

## 5. Loki、Kibana 日志顺序

### Loki

1. `loki_test_connection`
2. `loki_get_label_names`
3. `loki_get_label_values`
4. `loki_query_logs_summary`
5. `loki_query_range`
6. `loki_get_series`

### Kibana

1. `kibana_health_summary`
2. 有 space 时：`kibana_spaces_summary`
3. 日志查询：`kibana_query_logs`
4. dashboard：`kibana_dashboards_summary`、`kibana_get_dashboard_detail_advanced`
5. data view：`kibana_index_patterns_summary`、`kibana_get_data_views`
6. alert：`kibana_get_alerts`、`kibana_get_alert_rules`
7. saved object：`kibana_search_saved_objects_advanced`

## 6. Jaeger、OpenTelemetry 顺序

### Jaeger trace 诊断

1. `jaeger_get_services_summary`
2. `jaeger_get_service_ops`
3. `jaeger_get_traces_summary`
4. `jaeger_search_traces`
5. `jaeger_get_trace`
6. `jaeger_get_dependencies`

### OpenTelemetry collector 诊断

1. `opentelemetry_get_health`
2. `opentelemetry_get_collector_summary`
3. `opentelemetry_get_config_summary`
4. `opentelemetry_analyze_pipeline_status`
5. `opentelemetry_query_metrics`
6. `opentelemetry_query_logs`
7. `opentelemetry_query_traces`

## 7. Grafana 顺序

1. `grafana_test_connection`
2. `grafana_dashboards_summary` 或 `grafana_search_dashboards`
3. `grafana_dashboard`
4. `grafana_get_dashboard_panel_queries`
5. `grafana_datasources_summary`
6. `grafana_check_datasource_health`
7. 告警相关：`grafana_alerts`、`grafana_get_alert_rule_by_uid`
8. 渲染问题：`grafana_render_panel_image`
9. 交接链接：`grafana_generate_deeplink`、`grafana_generate_logs_drilldown_link`

## 8. Elasticsearch 顺序

1. `elasticsearch_cluster_health_summary`
2. `elasticsearch_nodes_summary`
3. `elasticsearch_indices_summary` 或 `elasticsearch_list_indices_paginated`
4. `elasticsearch_index_stats`
5. `elasticsearch_search_indices`
6. 深度分析：`elasticsearch_get_index_detail_advanced`、`elasticsearch_get_cluster_detail_advanced`

## 9. Nacos 顺序

1. `nacos_test_connection`
2. `nacos_list_namespaces`
3. 配置问题：`nacos_list_configs_summary`、`nacos_get_config`
4. 服务发现问题：`nacos_list_services_summary`、`nacos_get_service`、`nacos_list_instances`
5. 服务端问题：`nacos_list_cluster_nodes`、`nacos_get_system_metrics`

## 10. Sentry、Langfuse 顺序

### Sentry

1. `sentry_test_connection`
2. scope 不明确时：`sentry_list_organizations`、`sentry_list_projects`
3. `sentry_list_issues_summary`
4. `sentry_list_issues`
5. `sentry_get_issue`
6. `sentry_list_issue_events`
7. `sentry_get_issue_event`

### Langfuse

1. `langfuse_check_health`
2. `langfuse_list_traces_summary`
3. `langfuse_get_trace`
4. `langfuse_list_sessions`
5. `langfuse_list_observations`
6. prompt 问题：`langfuse_list_prompts`、`langfuse_get_prompt`
7. 质量或成本问题：`langfuse_list_scores`、`langfuse_get_metrics`
8. 评测问题：`langfuse_list_datasets`、`langfuse_list_dataset_runs`
9. 项目管理：`langfuse_get_project`、`langfuse_list_organization_projects`、`langfuse_create_project`、`langfuse_update_project`、`langfuse_delete_project`
10. 项目成员管理：`langfuse_list_project_memberships`、`langfuse_upsert_project_membership`、`langfuse_delete_project_membership`
11. API Key 管理：`langfuse_list_organization_api_keys`、`langfuse_list_project_api_keys`、`langfuse_create_project_api_key`、`langfuse_delete_project_api_key`

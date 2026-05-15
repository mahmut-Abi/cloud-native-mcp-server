# 用户提问样例与解决路径

本文件适合作为知识库检索入口。每条包含用户问题、意图分类、推荐工具链和输出重点。

| 用户提问 | 意图 | 推荐工具链 | 输出重点 |
| --- | --- | --- | --- |
| “prod 的 api Pod 一直重启，帮我看原因。” | 只读诊断 | `kubernetes_get_resource_summary` -> `kubernetes_get_recent_events` -> `kubernetes_get_pod_logs` -> `kubernetes_get_rollout_status` | 状态、事件、日志错误、是否发布相关 |
| “checkout Deployment rollout 卡住了。” | 只读诊断 | `kubernetes_get_rollout_status` -> `kubernetes_get_recent_events` -> `kubernetes_list_resources_summary` -> `kubernetes_get_pod_logs` | 卡在哪个 revision/Pod，失败原因 |
| “把 api 重启一下。” | 变更请求 | `kubernetes_get_resource_summary` -> 要求确认 -> `kubernetes_restart_workload` -> `kubernetes_get_rollout_status` | 影响范围、确认、验证 |
| “把 worker 副本数调到 5。” | 变更请求 | `kubernetes_get_resource_summary` -> 要求确认 -> `kubernetes_scale_resource` -> `kubernetes_get_rollout_status` | 当前副本、目标副本、HPA 风险 |
| “Service 返回 503。” | 只读诊断 | `kubernetes_get_resource` -> `kubernetes_list_resources_summary` for Pods -> `kubernetes_list_resources_summary` for EndpointSlice -> `loki_query_logs_summary` -> `jaeger_get_traces_summary` | selector、endpoint、Pod readiness、请求路径 |
| “Pod 一直 Pending。” | 只读诊断 | `kubernetes_get_resource_summary` -> `kubernetes_get_recent_events` -> `kubernetes_get_node_conditions` -> `kubernetes_get_resource_usage` | 调度失败原因、资源/亲和/污点/PVC |
| “为什么镜像拉不下来？” | 只读诊断 | `kubernetes_get_resource_summary` -> `kubernetes_get_recent_events` -> `kubernetes_get_resource` | tag、registry、auth、网络、imagePullSecrets |
| “节点 NotReady，帮我判断影响。” | 只读诊断 | `kubernetes_get_node_conditions` -> `kubernetes_get_recent_events` -> `kubernetes_list_resources_summary` -> `kubernetes_get_resource_usage` | 节点条件、受影响 Pod、是否需要 cordon/drain |
| “这个 service account 为什么 forbidden？” | 只读诊断 | `kubernetes_get_recent_events` -> `kubernetes_check_permissions` -> `kubernetes_get_resource` for RBAC | 缺少的 verb/resource、绑定错误 |
| “HighErrorRate 告警为什么触发？” | 只读诊断 | `alertmanager_alerts_summary` -> `prometheus_alerts_summary` -> `prometheus_rules_summary` -> `prometheus_query_range` -> `loki_query_logs_summary` | 告警表达式、指标证据、关联日志 |
| “现在有哪些告警被 silence 了？” | 只读诊断 | `alertmanager_silences_summary` -> `alertmanager_alerts_summary` | silence matcher、影响告警 |
| “给维护窗口加 silence。” | 变更请求 | `alertmanager_alerts_summary` -> `alertmanager_silences_summary` -> 要求确认 -> `alertmanager_create_silence` -> `alertmanager_silences_summary` | matcher、时间、原因、影响范围 |
| “Prometheus 没有这个服务的指标。” | 只读诊断 | `prometheus_targets_summary` -> `prometheus_get_label_values` -> `prometheus_get_series` -> `kubernetes_get_resource_summary` -> `loki_query_logs_summary` | target 是否存在、scrape 是否成功、selector 是否正确 |
| “Grafana dashboard 空白。” | 只读诊断 | `grafana_dashboards_summary` -> `grafana_dashboard` -> `grafana_get_dashboard_panel_queries` -> `grafana_datasources_summary` -> `grafana_check_datasource_health` | dashboard、panel query、datasource health |
| “Loki 查不到日志。” | 只读诊断 | `loki_test_connection` -> `loki_get_label_names` -> `loki_get_label_values` -> `loki_query_logs_summary` -> `loki_get_series` | selector、stream、采集链路 |
| “接口最近变慢。” | 只读诊断 | `prometheus_query_range` -> `jaeger_get_services_summary` -> `jaeger_get_traces_summary` -> `jaeger_get_trace` -> `loki_query_logs_summary` | 慢在哪个 span，指标趋势，日志错误 |
| “某个 trace 失败了。” | 只读诊断 | `jaeger_get_trace` -> `loki_query_logs_summary` -> `sentry_list_issues_summary` -> `kubernetes_get_resource_summary` | 失败 span、相关日志、应用异常 |
| “Sentry 登录接口错误变多。” | 只读诊断 | `sentry_test_connection` -> `sentry_list_issues_summary` -> `sentry_get_issue` -> `sentry_list_issue_events` -> `loki_query_logs_summary` | issue、release、event、stack、影响范围 |
| “LLM 回答质量下降。” | 只读诊断 | `langfuse_check_health` -> `langfuse_list_traces_summary` -> `langfuse_get_trace` -> `langfuse_list_prompts` -> `langfuse_list_scores` | prompt 版本、trace、score、模型/上下文变化 |
| “token 成本突然升高。” | 只读诊断 | `langfuse_get_metrics` -> `langfuse_list_traces_summary` -> `langfuse_list_observations` -> `langfuse_list_prompts` | 成本来源、token、模型、重试、prompt |
| “OpenTelemetry Collector 没有发 trace。” | 只读诊断 | `opentelemetry_get_health` -> `opentelemetry_get_collector_summary` -> `opentelemetry_get_config_summary` -> `opentelemetry_analyze_pipeline_status` -> `jaeger_get_services_summary` | receiver/pipeline/exporter、后端是否收到 |
| “Helm 升级失败了。” | 只读诊断 | `helm_get_release_status` -> `helm_get_release_history` -> `helm_get_release_values` -> `helm_get_release_manifest` -> `kubernetes_get_recent_events` | release 状态、revision、manifest、K8s 事件 |
| “帮我回滚 Helm release。” | 恢复请求 | `helm_get_release_status` -> `helm_get_release_history` -> `helm_compare_revisions` -> 要求确认 -> `helm_rollback_release` -> `helm_get_release_status` -> `kubernetes_get_rollout_status` | 回滚 revision、证据、风险、验证 |
| “Argo CD app OutOfSync。” | 只读诊断 | `argocd_test_connection` -> `argocd_list_applications_summary` -> `argocd_get_application` -> `argocd_get_application_manifests` -> `kubernetes_get_rollout_status` | sync/health、diff 来源、Git 与集群关系 |
| “Nacos 配置改了没生效。” | 只读诊断 | `nacos_test_connection` -> `nacos_list_namespaces` -> `nacos_list_configs_summary` -> `nacos_get_config` -> `kubernetes_get_pod_logs` | namespace/group/dataId、配置内容、应用是否动态加载 |
| “Nacos 服务实例不健康。” | 只读诊断 | `nacos_list_services_summary` -> `nacos_get_service` -> `nacos_list_instances` -> `nacos_list_cluster_nodes` -> `kubernetes_get_resource_summary` | 服务、实例、心跳、Nacos 节点 |
| “Elasticsearch 集群 red。” | 只读诊断 | `elasticsearch_cluster_health_summary` -> `elasticsearch_nodes_summary` -> `elasticsearch_indices_summary` -> `elasticsearch_get_cluster_detail_advanced` | red index、unassigned primary、节点和磁盘 |
| “Kibana dashboard 没数据。” | 只读诊断 | `kibana_health_summary` -> `kibana_spaces_summary` -> `kibana_dashboards_summary` -> `kibana_get_dashboard_detail_advanced` -> `kibana_get_data_views` -> `elasticsearch_indices_summary` | space、data view、时间字段、index 数据 |
| “Kibana alert 没触发。” | 只读诊断 | `kibana_health_summary` -> `kibana_get_alert_rules` -> `kibana_get_alerts` -> `kibana_query_logs` | rule 状态、mute/disabled、query、connector |
| “Pod 里解析 service 域名失败。” | 只读诊断 | `kubernetes_get_resource` -> `kubernetes_list_resources_summary` for Service/EndpointSlice -> `kubernetes_list_resources_summary` for CoreDNS Pods -> `kubernetes_get_pod_logs` | DNS 配置、Service 是否存在、CoreDNS 健康 |
| “两个服务之间突然访问超时。” | 只读诊断 | `kubernetes_get_resource` for Service/Pod -> `kubernetes_list_resources_summary` for NetworkPolicy -> `loki_query_logs_summary` -> `jaeger_get_traces_summary` | NetworkPolicy、endpoint、是否到达目标 |
| “Ingress 返回 502。” | 只读诊断 | `kubernetes_get_resource` for Ingress/Service -> `kubernetes_list_resources_summary` for EndpointSlice -> `kubernetes_get_recent_events` -> `kubernetes_get_pod_logs` | host/path、backend、端口、controller 日志 |
| “HTTPS 证书过期了。” | 只读诊断 | `kubernetes_get_resource` for Ingress/Secret -> `kubernetes_list_resources_summary` for Certificate/Issuer -> `kubernetes_get_recent_events` -> `kubernetes_get_pod_logs` | Secret、Certificate Ready、签发失败原因 |
| “PVC 一直 Pending。” | 只读诊断 | `kubernetes_get_resource_summary` -> `kubernetes_get_recent_events` -> `kubernetes_get_resource` for PVC/PV/StorageClass -> `kubernetes_get_pod_logs` for CSI | StorageClass、容量、accessMode、CSI |
| “StatefulSet 升级卡住。” | 只读诊断 | `kubernetes_get_resource_summary` -> `kubernetes_get_rollout_status` -> `kubernetes_list_resources_summary` for Pods/PVC -> `kubernetes_get_pod_logs` | ordinal、PVC、readiness、有状态依赖 |
| “HPA 没有扩容。” | 只读诊断 | `kubernetes_get_resource` for HPA -> `kubernetes_get_recent_events` -> `prometheus_query_range` -> `prometheus_targets_summary` | 指标可用性、maxReplicas、目标指标 |
| “HPA 副本数一直抖动。” | 只读诊断 | `kubernetes_get_resource` for HPA -> `prometheus_query_range` -> `kubernetes_get_resource_summary` | 指标波动、behavior、阈值和窗口 |
| “节点 drain 卡住了。” | 只读诊断 | `kubernetes_get_node_conditions` -> `kubernetes_list_resources_summary` for Pods/PDB -> `kubernetes_get_resource` for PDB -> `kubernetes_get_recent_events` | PDB、local storage、关键 Pod |
| “CronJob 没有按时执行。” | 只读诊断 | `kubernetes_get_resource_summary` -> `kubernetes_get_resource` for CronJob -> `kubernetes_list_resources_summary` for Jobs/Pods -> `kubernetes_get_recent_events` | schedule、suspend、deadline、concurrencyPolicy |
| “Job 一直失败重试。” | 只读诊断 | `kubernetes_get_resource_summary` -> `kubernetes_list_resources_summary` for Pods -> `kubernetes_get_recent_events` -> `kubernetes_get_pod_logs` | backoffLimit、命令、镜像、权限、依赖 |
| “DaemonSet 没跑满所有节点。” | 只读诊断 | `kubernetes_get_resource_summary` -> `kubernetes_get_resource` -> `kubernetes_list_resources_summary` -> `kubernetes_get_node_conditions` | taint/toleration、nodeSelector、affinity |
| “创建资源被 webhook 拒绝。” | 只读诊断 | `kubernetes_get_recent_events` -> `kubernetes_list_resources_summary` for webhook configs -> `kubernetes_get_resource_summary` for webhook backend -> `kubernetes_get_pod_logs` | failurePolicy、后端 Service、TLS、超时 |
| “Namespace quota exceeded。” | 只读诊断 | `kubernetes_get_recent_events` -> `kubernetes_list_resources_summary` for ResourceQuota/LimitRange -> `kubernetes_get_resource` -> `kubernetes_get_resource_usage` | used/hard、LimitRange、资源清理或扩额 |
| “Secret 更新后应用还是旧密码。” | 只读诊断/恢复请求 | `kubernetes_get_resource` for Pod/Secret -> `kubernetes_get_pod_logs` -> 要求确认 -> `kubernetes_restart_workload` | env 注入、是否需要重启、密钥不泄露 |
| “告警一直抖动。” | 只读诊断 | `alertmanager_alerts_summary` -> `prometheus_rules_summary` -> `prometheus_query_range` -> `loki_query_logs_summary` | 阈值、for 时长、scrape 抖动、真实尖峰 |
| “Prometheus 查询很慢。” | 只读诊断 | `prometheus_get_runtime_info` -> `prometheus_get_tsdb_status` -> `prometheus_get_tsdb_stats` -> `grafana_get_dashboard_panel_queries` | 高基数、查询范围、step、recording rule |
| “Prometheus rule 没生效。” | 只读诊断 | `prometheus_rules_summary` -> `prometheus_get_rules` -> `prometheus_query` -> `kubernetes_get_pod_logs` | rule 是否加载、evaluation error、reload |
| “日志延迟好几分钟才出现。” | 只读诊断 | `loki_query_logs_summary` -> `loki_get_series` -> `kubernetes_get_resource_summary` for collector -> `kubernetes_get_pod_logs` -> `prometheus_query_range` | 采集器队列、后端限流、pipeline |
| “LogQL 查询太慢。” | 只读诊断 | `loki_get_label_names` -> `loki_get_label_values` -> `loki_get_series` -> `loki_query_logs_summary` | selector 过宽、高基数、日志量 |
| “Jaeger 里没有 trace。” | 只读诊断 | `jaeger_get_services_summary` -> `opentelemetry_get_collector_summary` -> `opentelemetry_get_config_summary` -> `opentelemetry_analyze_pipeline_status` | SDK、service.name、collector exporter |
| “Trace 里缺少下游 span。” | 只读诊断 | `jaeger_get_trace` -> `jaeger_get_dependencies` -> `loki_query_logs_summary` -> `opentelemetry_query_traces` | propagation、sampling、下游 instrumentation |
| “Dashboard 变量选了以后没数据。” | 只读诊断 | `grafana_dashboard` -> `grafana_get_dashboard_panel_queries` -> `grafana_check_datasource_health` -> `prometheus_query` | 变量 query、regex、datasource UID |
| “Elasticsearch index 变只读了。” | 只读诊断/变更请求 | `elasticsearch_cluster_health_summary` -> `elasticsearch_nodes_summary` -> `elasticsearch_indices_summary` -> `elasticsearch_index_stats` | flood watermark、磁盘、先修根因再解除只读 |
| “Kibana connector 测试失败。” | 只读诊断 | `kibana_health_summary` -> `kibana_get_connectors` -> `kibana_get_connector` -> `kibana_test_connector` | URL、认证、外部通知平台、task manager |
| “Sentry issue 没有关联 release。” | 只读诊断 | `sentry_list_issues_summary` -> `sentry_get_issue` -> `sentry_list_issue_events` -> `kubernetes_get_resource` | SDK release、environment、镜像 tag |
| “Langfuse trace 不完整。” | 只读诊断 | `langfuse_list_traces_summary` -> `langfuse_get_trace` -> `langfuse_list_sessions` -> `langfuse_list_observations` -> `loki_query_logs_summary` | session id、observation、SDK flush |
| “列出 Langfuse 组织内的项目。” | 只读诊断 | `langfuse_list_organization_projects` | project id、name、organization、retentionDays |
| “帮我创建一个 Langfuse 项目。” | 变更请求 | `langfuse_list_organization_projects` -> 要求确认 -> `langfuse_create_project` -> `langfuse_list_organization_projects` | name、metadata、retention_days、组织级凭据要求 |
| “更新 Langfuse 项目名称或 metadata。” | 变更请求 | `langfuse_list_organization_projects` -> 要求确认 -> `langfuse_update_project` -> `langfuse_list_organization_projects` | project_id、name、metadata、retention_days |
| “删除 Langfuse 项目。” | 变更请求 | `langfuse_list_organization_projects` -> 要求强确认 -> `langfuse_delete_project` -> `langfuse_list_organization_projects` | 删除异步且高风险，确认 project_id |
| “列出 Langfuse 项目成员。” | 只读诊断 | `langfuse_list_project_memberships` | project_id、userId、role、email |
| “给某个用户添加或修改 Langfuse 项目角色。” | 变更请求 | `langfuse_list_project_memberships` -> 要求确认 -> `langfuse_upsert_project_membership` -> `langfuse_list_project_memberships` | project_id、user_id、role 只能是 OWNER/ADMIN/MEMBER/VIEWER |
| “移除 Langfuse 项目成员。” | 变更请求 | `langfuse_list_project_memberships` -> 要求强确认 -> `langfuse_delete_project_membership` -> `langfuse_list_project_memberships` | project_id、user_id、删除不可逆、确认不会移除最后管理员 |
| “帮我给 Langfuse 项目创建一组 AK/SK。” | 变更请求 | `langfuse_list_project_api_keys` -> 要求确认 -> `langfuse_create_project_api_key` -> `langfuse_list_project_api_keys` | project_id、note、secretKey 只在创建响应返回、组织级凭据要求 |
| “列出 Langfuse 项目的 API Keys。” | 只读诊断 | `langfuse_list_project_api_keys` | api key id、publicKey、displaySecretKey、lastUsedAt |
| “删除 Langfuse 项目的某个 AK/SK。” | 变更请求 | `langfuse_list_project_api_keys` -> 要求确认 -> `langfuse_delete_project_api_key` -> `langfuse_list_project_api_keys` | project_id、api_key_id、删除不可逆、验证已删除 |

## 输出模板

对于上面的任意问题，agent 都应尽量输出：

```text
问题类型：<解释|只读诊断|变更请求|验证请求|恢复请求>
作用范围：<cluster|namespace|workload|service|release|dashboard|alert|index>
已用工具：<tool list>
观察事实：<facts>
推断：<inference>
下一步：<next tool or remediation>
是否需要确认：<yes/no>
```

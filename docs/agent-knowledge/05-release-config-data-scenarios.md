# 发布、配置与数据平台场景

本文件覆盖 Helm、Argo CD、Nacos、Elasticsearch、Kibana 等场景。发布和配置类问题要特别注意“源头在哪里”：Git、Helm values、集群临时状态、配置中心或外部系统。

## 1. Helm 升级失败或 release 异常

用户提问示例：

- “Helm upgrade 后服务不可用，帮我看是不是这个版本的问题。”
- “release 一直 pending-upgrade，怎么处理？”

推荐工具顺序：

1. `helm_list_releases_paginated`
2. `helm_get_release_summary`
3. `helm_get_release_status`
4. `helm_get_release_history`
5. `helm_get_release_values`
6. `helm_get_release_manifest`
7. `helm_compare_revisions`
8. 关联 `kubernetes_get_rollout_status`、`kubernetes_get_recent_events`、`kubernetes_get_pod_logs`

判断方法：

- release 状态 failed/pending：查看 hook、manifest、Kubernetes events。
- values 变化导致配置错误：比较 revision。
- manifest 正确但 Pod 异常：继续走 Kubernetes workload 场景。
- 历史 revision 中存在稳定版本：可评估 rollback。

修复建议：

- chart 或 values 错误：修正 values 后重新 upgrade。
- 明确坏版本且用户确认：使用 `helm_rollback_release`。
- pending 状态：谨慎确认是否存在并发操作或锁。

验证：

- `helm_get_release_status`
- `kubernetes_get_rollout_status`
- `prometheus_query_range`
- `loki_query_logs_summary`

## 2. Argo CD OutOfSync 或 Degraded

用户提问示例：

- “Argo CD 应用 OutOfSync，帮我看差异在哪。”
- “Application health degraded，但是 Kubernetes 看起来有 Pod，怎么查？”

推荐工具顺序：

1. `argocd_test_connection`
2. `argocd_list_applications_summary`
3. `argocd_get_application`
4. `argocd_get_application_manifests`
5. `kubernetes_get_rollout_status`
6. `kubernetes_get_recent_events`
7. `kubernetes_get_pod_logs`

判断方法：

- Sync OutOfSync：Git 目标状态与集群实际状态不一致。
- Health Degraded：资源创建成功但 runtime 不健康。
- Manifest 渲染错误：源配置、Helm values、Kustomize、插件问题。
- 手工 patch 被 Argo 覆盖：长期修复必须进入 Git。

修复建议：

- 配置源错误：修改 Git。
- 紧急止血：用户确认后可做临时 Kubernetes patch，但必须说明会被 GitOps 覆盖。
- 应用 runtime 故障：按 Kubernetes 场景修复。

验证：

- `argocd_get_application`
- `kubernetes_get_rollout_status`
- `kubernetes_get_recent_events`

## 3. 发布后回归

用户提问示例：

- “昨晚发布后错误率升高，帮我判断要不要回滚。”
- “新版本上线后延迟升高，给出证据链。”

推荐工具顺序：

1. 判断发布路径：Kubernetes 原生、Helm、Argo CD。
2. `kubernetes_get_rollout_status`
3. `kubernetes_get_recent_events`
4. `kubernetes_list_resources_summary` 查新旧 Pod。
5. `prometheus_query_range` 对比发布前后错误率、延迟、流量。
6. `loki_query_logs_summary` 查新版本日志。
7. `sentry_list_issues_summary` 查异常。
8. Helm：`helm_get_release_history`、`helm_compare_revisions`
9. Argo CD：`argocd_get_application`、`argocd_get_application_manifests`

判断方法：

- 错误从新 revision 开始并集中在新 Pod：发布回归证据强。
- 指标异常早于发布：可能是外部依赖或基础设施问题。
- trace 显示下游慢：不一定应回滚当前服务。

修复建议：

- 证据指向坏版本：用户确认后回滚。
- 配置问题：修正 values、ConfigMap、Nacos 或 Secret。
- 资源瓶颈：扩容或调整资源。

验证：

- 回滚后查看 `kubernetes_get_rollout_status`
- `prometheus_query_range` 对比错误率/延迟恢复
- `sentry_list_issues_summary` 查看 issue 是否停止增长

## 4. Nacos 配置未生效

用户提问示例：

- “Nacos 配置改了但是应用没生效。”
- “某个服务读取到的配置和预期不一致。”

推荐工具顺序：

1. `nacos_test_connection`
2. `nacos_list_namespaces`
3. `nacos_list_configs_summary`
4. `nacos_get_config`
5. `kubernetes_get_resource` 查 Pod env、启动参数、配置引用。
6. `kubernetes_get_pod_logs` 或 `loki_query_logs_summary`

判断方法：

- namespace、group、dataId 不匹配。
- 应用未监听动态配置，需要重启。
- 配置内容正确但应用报错：格式、字段名或兼容性问题。
- 多环境配置混用。

修复建议：

- 修正 namespace/group/dataId。
- 修正配置格式和字段。
- 对非动态配置，用户确认后重启工作负载。

验证：

- `nacos_get_config`
- `kubernetes_get_pod_logs`
- `kubernetes_get_rollout_status`

## 5. Nacos 服务发现异常

用户提问示例：

- “Nacos 上服务实例缺失，调用方找不到服务。”
- “服务注册了但是实例不健康。”

推荐工具顺序：

1. `nacos_test_connection`
2. `nacos_list_namespaces`
3. `nacos_list_services_summary`
4. `nacos_get_service`
5. `nacos_list_instances`
6. `nacos_list_cluster_nodes`
7. `nacos_get_system_metrics`
8. 关联 Kubernetes Pod 状态和日志

判断方法：

- 服务不存在：注册失败或 namespace/group 错误。
- 实例不健康：应用心跳、网络、端口或健康检查失败。
- Nacos 节点异常：服务端集群问题。

修复建议：

- 修正注册配置。
- 修复应用健康检查或网络。
- 修复 Nacos 节点状态。

验证：

- `nacos_list_instances`
- `nacos_get_service`
- 调用方日志或 trace

## 6. Elasticsearch 集群 yellow/red

用户提问示例：

- “ES 集群 red，帮我判断影响。”
- “为什么有 index yellow？”

推荐工具顺序：

1. `elasticsearch_cluster_health_summary`
2. `elasticsearch_nodes_summary`
3. `elasticsearch_indices_summary`
4. `elasticsearch_list_indices_paginated`
5. `elasticsearch_index_stats`
6. `elasticsearch_get_cluster_detail_advanced`
7. `elasticsearch_get_index_detail_advanced`

判断方法：

- red：primary shard 未分配，数据读写可能受影响。
- yellow：replica shard 未分配，通常可读写但冗余不足。
- 节点少于 replica 要求：单节点集群常见 yellow。
- 磁盘水位高：分片分配受限。

修复建议：

- 扩容节点或调整 replica。
- 释放磁盘、调整 ILM、删除过期 index。
- 修复异常节点。

验证：

- `elasticsearch_cluster_health_summary`
- `elasticsearch_indices_summary`
- `elasticsearch_nodes_summary`

## 7. Elasticsearch 查询结果异常或 index 缺失

用户提问示例：

- “为什么这个 index 没有新数据？”
- “搜索结果不符合预期，帮我查 index 状态。”

推荐工具顺序：

1. `elasticsearch_indices_summary`
2. `elasticsearch_search_indices`
3. `elasticsearch_index_stats`
4. `elasticsearch_get_index_detail_advanced`
5. 如果来自 Kibana：`kibana_get_data_views`、`kibana_query_logs`

判断方法：

- index 不存在：ingest pipeline、采集器或命名规则问题。
- docs count 不增长：写入链路中断。
- mapping 不符合查询：字段类型或 analyzer 问题。
- Kibana 查不到但 ES 有数据：data view 或时间字段问题。

修复建议：

- 修复 ingest、采集器或 index template。
- 修正 data view。
- 修正查询 DSL 或字段。

验证：

- `elasticsearch_search_indices`
- `kibana_query_logs`

## 8. Kibana Dashboard、Data View 或 Alert 异常

用户提问示例：

- “Kibana dashboard 没数据，但是 ES index 有数据。”
- “Kibana alert 没触发，帮我看配置。”

推荐工具顺序：

1. `kibana_health_summary`
2. `kibana_spaces_summary`
3. Dashboard：`kibana_dashboards_summary`、`kibana_get_dashboard_detail_advanced`
4. Data View：`kibana_index_patterns_summary`、`kibana_get_data_views`
5. Logs：`kibana_query_logs`
6. Alerts：`kibana_get_alerts`、`kibana_get_alert_rules`、`kibana_get_alert_rule`
7. Saved objects：`kibana_search_saved_objects_advanced`
8. Elasticsearch：`elasticsearch_indices_summary`

判断方法：

- space 错误：对象在另一个 space。
- data view 时间字段或 index pattern 错误：Discover/dashboard 无数据。
- alert rule disabled/muted：不会触发通知。
- connector 失败：通知链路问题。

修复建议：

- 修正 data view 或 index pattern。
- 修复 dashboard query。
- enable/unmute/update alert rule 前要求用户确认。
- 修复 connector。

验证：

- `kibana_query_logs`
- `kibana_get_alert_rules`
- `elasticsearch_indices_summary`


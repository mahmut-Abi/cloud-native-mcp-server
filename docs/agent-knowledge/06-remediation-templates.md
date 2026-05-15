# 修复动作模板

本文件只描述修复动作的安全执行模板。任何写操作前都必须先读当前状态并要求用户明确确认。

## 1. Kubernetes Patch

适用：

- 修改 probe、resources、labels、annotations、env、replicas 以外的 spec 局部字段。
- 临时修复 GitOps/Helm 管理资源时必须说明会被控制器覆盖。

执行前工具：

1. `kubernetes_get_resource`
2. `kubernetes_get_recent_events`
3. 如果是 workload：`kubernetes_get_rollout_status`

确认话术：

```text
我将对 <namespace>/<kind>/<name> 执行 merge patch。
Patch 内容为 <patch>。
风险：可能触发 rollout 或被 GitOps/Helm 覆盖。
验证方式：patch 后检查 resource summary、rollout status 和 recent events。
请确认是否执行。
```

执行工具：

- `kubernetes_patch_resource`

验证工具：

1. `kubernetes_get_resource_summary`
2. `kubernetes_get_rollout_status`
3. `kubernetes_get_recent_events`
4. 业务验证：Prometheus、Loki、Jaeger、Sentry

## 2. Kubernetes Restart

适用：

- 配置已修复但应用只在启动时加载。
- Pod 卡死、连接池异常等需要临时恢复，但必须说明不是根因修复。

执行前工具：

1. `kubernetes_get_resource_summary`
2. `kubernetes_get_recent_events`
3. `kubernetes_get_rollout_status`

确认话术：

```text
我将重启 <namespace>/<kind>/<name>。
影响：会触发滚动重启，期间可能出现短暂容量下降。
验证方式：检查 rollout status、Pod readiness、错误率和日志。
请确认是否执行。
```

执行工具：

- `kubernetes_restart_workload`

验证工具：

1. `kubernetes_get_rollout_status`
2. `kubernetes_wait_for_resource`
3. `kubernetes_get_recent_events`
4. `prometheus_query_range`
5. `loki_query_logs_summary`

## 3. Kubernetes Scale

适用：

- 临时流量升高、资源压力、容量不足。
- 回滚前需要临时降低坏版本影响时，要先判断控制器策略。

执行前工具：

1. `kubernetes_get_resource_summary`
2. `kubernetes_get_resource_usage`
3. `prometheus_query_range`

确认话术：

```text
我将把 <namespace>/<kind>/<name> 的副本数从 <current> 调整到 <desired>。
影响：会改变容量和资源消耗；如果存在 HPA，可能被 HPA 后续覆盖。
请确认是否执行。
```

执行工具：

- `kubernetes_scale_resource`

验证工具：

1. `kubernetes_get_rollout_status`
2. `kubernetes_list_resources_summary` 查 Pods
3. `prometheus_query_range` 查容量和错误率

## 4. Kubernetes Delete

适用：

- 删除明确错误创建的资源。
- 删除卡住的 Pod 让控制器重建。

禁止：

- 未确认 ownerReferences、finalizers、namespace 和资源名就删除。
- 删除 PV、PVC、Secret、Namespace 等高风险资源而没有强确认。

执行前工具：

1. `kubernetes_get_resource`
2. `kubernetes_get_recent_events`
3. 必要时 `kubernetes_list_resources_summary` 查相关资源

确认话术：

```text
我将删除 <namespace>/<kind>/<name>。
这是不可逆动作。当前 ownerReferences/finalizers 为 <summary>。
预期结果：<expected>。
请明确确认删除。
```

执行工具：

- `kubernetes_delete_resource`

验证工具：

1. `kubernetes_list_resources_summary`
2. `kubernetes_get_recent_events`
3. 如果由控制器重建：`kubernetes_get_rollout_status`

## 5. Node Cordon、Drain、Uncordon

适用：

- 节点维护、节点 NotReady、资源压力、磁盘压力。

执行前工具：

1. `kubernetes_get_node_conditions`
2. `kubernetes_get_resource_usage`
3. `kubernetes_list_resources_summary` 查该节点上的 Pod
4. `kubernetes_get_recent_events`

确认话术：

```text
我将对节点 <node> 执行 <cordon|drain|uncordon>。
影响：drain 会驱逐 Pod，可能受 PDB、DaemonSet、local storage 影响。
请确认是否执行。
```

执行工具：

- `kubernetes_cordon_node`
- `kubernetes_drain_node`
- `kubernetes_uncordon_node`

验证工具：

1. `kubernetes_get_node_conditions`
2. `kubernetes_list_resources_summary`
3. `kubernetes_get_recent_events`

## 6. Helm Rollback

适用：

- 证据明确表明当前 Helm revision 导致错误。
- 存在已知稳定 revision。

执行前工具：

1. `helm_get_release_status`
2. `helm_get_release_history`
3. `helm_compare_revisions`
4. `helm_get_release_manifest`
5. `kubernetes_get_rollout_status`
6. `prometheus_query_range` 和 `loki_query_logs_summary`

确认话术：

```text
我将把 Helm release <namespace>/<release> 回滚到 revision <revision>。
证据：<evidence>。
风险：会恢复旧版本配置和镜像，可能覆盖当前 revision 的修复。
验证方式：release status、rollout、错误率、日志。
请确认是否执行。
```

执行工具：

- `helm_rollback_release`

验证工具：

1. `helm_get_release_status`
2. `kubernetes_get_rollout_status`
3. `kubernetes_get_recent_events`
4. `prometheus_query_range`
5. `sentry_list_issues_summary`

## 7. Helm Upgrade 或 Install

执行前工具：

1. `helm_search_charts` 或 `helm_get_chart_info`
2. `helm_template_chart`
3. 已存在 release 时：`helm_get_release_values`、`helm_get_release_status`
4. `helm_validate_release`

执行工具：

- `helm_install_release`
- `helm_upgrade_release`

验证工具：

1. `helm_get_release_status`
2. `kubernetes_get_rollout_status`
3. `kubernetes_get_recent_events`
4. `kubernetes_get_pod_logs`

## 8. Alertmanager Silence

适用：

- 维护窗口、已知故障、重复告警临时降噪。

执行前工具：

1. `alertmanager_alerts_summary`
2. `alertmanager_query_alerts_advanced`
3. `alertmanager_silences_summary`

确认话术：

```text
我将创建 silence：
- matcher: <matchers>
- startsAt: <start>
- endsAt: <end>
- createdBy: <user>
- comment: <reason>
影响：匹配这些标签的告警不会通知。
请确认是否执行。
```

执行工具：

- `alertmanager_create_silence`

删除 silence：

1. 先 `alertmanager_silences_summary`
2. 确认 silence id
3. 用户确认后 `alertmanager_delete_silence`

验证工具：

- `alertmanager_silences_summary`
- `alertmanager_alerts_summary`

## 9. Grafana Dashboard 或 Datasource 更新

执行前工具：

1. `grafana_test_connection`
2. `grafana_dashboard`
3. `grafana_get_dashboard_versions`
4. `grafana_datasources_summary`
5. `grafana_check_datasource_health`

执行工具：

- `grafana_update_dashboard`
- `grafana_update_datasource`
- 必要时 `grafana_restore_dashboard_version`

验证工具：

1. `grafana_dashboard`
2. `grafana_get_dashboard_panel_queries`
3. `grafana_check_datasource_health`
4. 对应后端查询工具

## 10. Kibana 对象或 Alert 更新

执行前工具：

1. `kibana_health_summary`
2. `kibana_search_saved_objects_advanced`
3. `kibana_get_dashboard_detail_advanced`
4. `kibana_get_data_views`
5. `kibana_get_alert_rules`

执行工具示例：

- `kibana_update_dashboard`
- `kibana_update_data_view`
- `kibana_update_alert_rule`
- `kibana_enable_alert_rule`
- `kibana_disable_alert_rule`

验证工具：

1. `kibana_query_logs`
2. `kibana_get_alert_rules`
3. `elasticsearch_indices_summary`


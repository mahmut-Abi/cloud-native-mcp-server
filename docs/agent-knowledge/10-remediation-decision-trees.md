# 修复决策树

本文件用于帮助 agent 在已有证据基础上选择最小修复动作。原则：没有证据不修复；能配置修复就不重启；能单对象 patch 就不大范围替换；能验证就不凭感觉结束。

## 1. 是否应该重启

适合重启：

- 配置源已修复，但应用只在启动时加载配置。
- Pod 卡死或连接池状态异常，有日志/指标证明重启可恢复。
- Secret 或 ConfigMap 通过 env 注入，更新后需要新进程读取。
- 临时止血，且用户接受不是根因修复。

不适合重启：

- 镜像错误、配置错误、PVC 挂载失败、资源不足、网络策略阻断。
- 所有 Pod 都因同一错误 CrashLoop，重启只会重复失败。
- HPA/Deployment 正在 rollout，重启会扩大扰动。

执行顺序：

1. `kubernetes_get_resource_summary`
2. `kubernetes_get_recent_events`
3. `kubernetes_get_pod_logs`
4. 用户确认
5. `kubernetes_restart_workload`
6. `kubernetes_get_rollout_status`
7. `prometheus_query_range` 或 `loki_query_logs_summary`

## 2. 是否应该扩容

适合扩容：

- 指标显示 CPU、内存、队列、连接池或请求量达到容量瓶颈。
- 错误率与流量上升强相关。
- 当前副本数低于安全容量，且资源池可承载。

不适合扩容：

- 错误来自配置、镜像、权限、下游故障、数据库锁或外部依赖。
- HPA 已经达到 maxReplicas，需要先确认上限和资源。
- 新副本也会因同样配置失败。

执行顺序：

1. `kubernetes_get_resource_summary`
2. `kubernetes_get_resource_usage`
3. `prometheus_query_range`
4. 如果有 HPA：`kubernetes_get_resource` 查 `HorizontalPodAutoscaler`
5. 用户确认
6. `kubernetes_scale_resource`
7. `kubernetes_get_rollout_status`

## 3. 是否应该 patch

适合 patch：

- 需要修改单个字段：probe、resource request/limit、annotation、label、env、selector、端口等。
- 修复范围明确，且 patch 后能立即验证。
- 临时止血需要小范围调整。

不适合 patch：

- GitOps/Helm 管理资源的长期修复，应进入 Git 或 values。
- 需要大范围结构性变更。
- 不清楚当前对象实际状态。

执行顺序：

1. `kubernetes_get_resource`
2. `kubernetes_get_recent_events`
3. 说明 patch 内容、风险和回滚方式。
4. 用户确认
5. `kubernetes_patch_resource`
6. `kubernetes_get_resource_summary`
7. `kubernetes_get_rollout_status` 或 `kubernetes_wait_for_resource`

## 4. 是否应该回滚

适合回滚：

- 错误率、延迟或 CrashLoop 从某个 release/revision 后开始。
- 旧 revision 已知稳定。
- 日志、trace、Sentry 或指标指向新版本代码或配置。
- 修复当前版本需要较长时间，业务需要先恢复。

不适合回滚：

- 问题早于发布发生。
- 根因是外部依赖、节点、网络、存储或告警误报。
- 旧版本也受同一外部故障影响。

Helm 回滚顺序：

1. `helm_get_release_status`
2. `helm_get_release_history`
3. `helm_compare_revisions`
4. `prometheus_query_range`
5. `loki_query_logs_summary`
6. 用户确认
7. `helm_rollback_release`
8. `helm_get_release_status`
9. `kubernetes_get_rollout_status`

Argo CD 场景：

- 使用 `argocd_get_application` 和 `argocd_get_application_manifests` 查源状态。
- 长期回滚应通过 Git revert 或回退目标 revision。
- 直接 patch 只能作为用户确认后的临时止血。

## 5. 是否应该删除 Pod

适合删除 Pod：

- Pod 被控制器管理，删除后会自动重建。
- Pod 卡在异常本地状态，配置和镜像已确认正确。
- 节点或容器运行时问题导致单个 Pod 异常。

不适合删除 Pod：

- 裸 Pod 或不确定 ownerReferences。
- 有状态 Pod 挂载重要卷，删除可能导致业务中断。
- 删除 PVC、PV、Secret、Namespace 等数据或关键资源。

执行顺序：

1. `kubernetes_get_resource`
2. 检查 ownerReferences、finalizers、volumes。
3. 用户强确认
4. `kubernetes_delete_resource`
5. `kubernetes_list_resources_summary`
6. `kubernetes_get_recent_events`

## 6. 是否应该创建 silence

适合 silence：

- 维护窗口内的已知告警。
- 已确认故障正在处理，需要临时减少重复通知。
- 告警噪声影响值班判断，且 silence 范围可精确限制。

不适合 silence：

- 未确认是否真实故障。
- matcher 过宽，可能屏蔽无关严重告警。
- 没有结束时间或原因。

执行顺序：

1. `alertmanager_alerts_summary`
2. `alertmanager_query_alerts_advanced`
3. `alertmanager_silences_summary`
4. 用户确认 matcher、时间、原因
5. `alertmanager_create_silence`
6. `alertmanager_silences_summary`

## 7. 是否应该 drain 节点

适合 drain：

- 节点维护、下线、硬件故障。
- 节点 NotReady 或 DiskPressure/MemoryPressure 持续影响业务。
- 业务副本和 PDB 允许迁移。

不适合 drain：

- PDB 不允许中断。
- 节点上有关键单副本或本地存储 Pod。
- 问题是全局资源不足，drain 会让其他节点更拥挤。

执行顺序：

1. `kubernetes_get_node_conditions`
2. `kubernetes_get_resource_usage`
3. `kubernetes_list_resources_summary` 查节点 Pod。
4. `kubernetes_list_resources_summary` 查 `PodDisruptionBudget`。
5. 用户确认
6. `kubernetes_cordon_node`
7. `kubernetes_drain_node`
8. `kubernetes_get_node_conditions`

## 8. 是否应该修改告警规则

适合修改：

- 告警表达式与业务目标不匹配。
- 指标噪声导致抖动，且有数据证明阈值或窗口不合理。
- route/receiver 需要调整以符合通知策略。

不适合修改：

- 告警真实反映故障，应该修业务。
- 没有足够历史数据证明误报。
- 临时为了消除告警而降低敏感度。

诊断顺序：

1. `prometheus_rules_summary`
2. `prometheus_get_rules`
3. `prometheus_query_range`
4. `alertmanager_alerts_summary`
5. `alertmanager_receivers_summary`

修复建议：

- 规则变更走配置发布流程。
- 临时静默用 silence，不直接删除规则。
- 修改后验证 firing 状态和通知状态。

## 9. 是否应该修改 Dashboard

适合修改：

- datasource UID 错误。
- panel query 与实际指标标签不匹配。
- 变量 query、regex、默认值错误。
- 面板时间范围导致误判。

不适合修改：

- 后端确实没有数据。
- 指标、日志、trace 采集链路故障。
- Dashboard 只是反映真实故障。

执行顺序：

1. `grafana_dashboard`
2. `grafana_get_dashboard_versions`
3. `grafana_get_dashboard_panel_queries`
4. `grafana_check_datasource_health`
5. 后端 query 验证
6. 用户确认
7. `grafana_update_dashboard`
8. `grafana_dashboard`

## 10. 是否应该修复采集链路

适合修采集链路：

- 后端健康，但某个 signal 缺失。
- 应用有日志/metrics/traces 产生，但后端查不到。
- Collector/exporter 日志有发送失败、丢弃、重试。

诊断顺序：

1. 后端健康：
   - metrics：`prometheus_test_connection`、`prometheus_targets_summary`
   - logs：`loki_test_connection`、`kibana_health_summary`
   - traces：`jaeger_get_services_summary`
2. collector：
   - `opentelemetry_get_collector_summary`
   - `opentelemetry_get_config_summary`
   - `opentelemetry_analyze_pipeline_status`
3. workload：
   - `kubernetes_get_resource_summary`
   - `kubernetes_get_pod_logs`
4. 验证：
   - 对应后端 query 工具

修复建议：

- 修正 receiver、processor、exporter 或 pipeline 引用。
- 修复鉴权、TLS、网络和 endpoint。
- 调整 batch、memory_limiter、sampling。

## 11. 快速症状到修复映射

| 症状 | 优先诊断 | 常见修复 | 不建议 |
| --- | --- | --- | --- |
| CrashLoopBackOff | 事件、lastState、日志 | 修配置、修 probe、调资源、必要时重启 | 盲目扩容 |
| Pending | 调度事件、节点、PVC、quota | 调资源、修 toleration/affinity、修 PVC | 重启 Deployment |
| ImagePullBackOff | 事件、Pod spec、values | 修 tag、registry、imagePullSecrets | 删除 Pod 反复重试 |
| Service 503 | Service、EndpointSlice、Pod readiness | 修 selector、端口、readiness、后端 | 只重启入口 |
| DNS no such host | Pod DNS、CoreDNS、Service | 修域名、CoreDNS、上游 DNS | 改业务重试掩盖 |
| HPA 不扩容 | HPA、metrics、Prometheus | 修指标链路、调 HPA、临时 scale | 只改 maxReplicas 不看指标 |
| PVC Pending | PVC、PV、StorageClass、事件 | 修 StorageClass、容量、CSI | 删除 PVC |
| Alert 未通知 | Prometheus alert、Alertmanager route/silence | 修 route/receiver、删除错误 silence | 关闭规则 |
| Dashboard 空白 | panel query、datasource、后端数据 | 修 query、变量、datasource | 直接改后端 |
| Trace 缺失 | Jaeger service、OTel pipeline、SDK logs | 修 instrumentation、collector、sampling | 只查应用日志 |


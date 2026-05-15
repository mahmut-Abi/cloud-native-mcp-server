# Kubernetes 进阶场景

本文件补充网络、DNS、Ingress/Gateway、存储、自动伸缩、安全准入、任务类工作负载等更细的排查和修复方法。

## 1. DNS 解析失败

用户提问示例：

- “Pod 里解析 service 域名失败。”
- “应用日志里大量 `no such host`，帮我排查 DNS。”

推荐工具顺序：

1. `kubernetes_get_resource_summary` 查业务 Pod。
2. `kubernetes_get_pod_logs` 查应用错误。
3. `kubernetes_get_resource` 查 Pod 的 `dnsPolicy`、`dnsConfig`、namespace。
4. `kubernetes_list_resources_summary` 查 `Service` 和 `EndpointSlice`。
5. `kubernetes_list_resources_summary` 查 kube-system 中 CoreDNS Pod。
6. `kubernetes_get_pod_logs` 查 CoreDNS 日志。
7. `prometheus_query_range` 查 DNS 错误率、CoreDNS 延迟和重启。

判断方法：

- 只有某个 namespace 失败：可能是 NetworkPolicy、DNS search path、服务名写错。
- 所有 Pod 都失败：CoreDNS、node local DNS、集群网络问题。
- 解析成功但连接失败：问题不在 DNS，转向 Service endpoint 或 NetworkPolicy。
- CoreDNS 日志出现 upstream timeout：上游 DNS 或网络出口问题。

修复建议：

- 修正服务名、namespace、FQDN。
- 修正 Pod `dnsPolicy` 或 `dnsConfig`。
- 修复 CoreDNS 配置或上游 DNS。
- 如果 CoreDNS Pod 不健康，先诊断后在用户确认下重启相关 Deployment。

验证：

- `kubernetes_get_resource_summary` 查 CoreDNS Ready。
- `kubernetes_get_pod_logs` 确认 DNS 错误停止。
- `prometheus_query_range` 确认 DNS 错误率下降。

## 2. NetworkPolicy 阻断流量

用户提问示例：

- “两个服务之前能通，今天开始互相访问超时。”
- “是不是 NetworkPolicy 把流量挡住了？”

推荐工具顺序：

1. `kubernetes_get_resource` 查源 Pod 和目标 Service。
2. `kubernetes_list_resources_summary` 查源/目标 namespace 的 `NetworkPolicy`。
3. `kubernetes_get_resource` 读取相关 NetworkPolicy 规则。
4. `kubernetes_list_resources_summary` 查目标 Service 的 EndpointSlice。
5. `loki_query_logs_summary` 查源服务连接超时和目标服务访问日志。
6. `jaeger_get_traces_summary` 查请求是否到达目标服务。

判断方法：

- 源日志显示 timeout，目标无日志：可能在网络层被拒。
- 目标 Service 有 endpoint，但源无法访问：NetworkPolicy、CNI 或安全组。
- 只有某些 Pod 受影响：label selector 不匹配或新版本 Pod label 变化。
- trace 没有目标 span：请求未到达目标。

修复建议：

- 修正 NetworkPolicy 的 podSelector、namespaceSelector、ports。
- 恢复被发布改掉的 labels。
- 临时放通规则属于高风险变更，必须用户确认并说明作用范围。

验证：

- `loki_query_logs_summary` 查源超时消失。
- `jaeger_get_traces_summary` 查目标 span 出现。
- `prometheus_query_range` 查请求成功率恢复。

## 3. Ingress 或 Gateway 路由异常

用户提问示例：

- “域名访问 404/502/503，帮我查入口路由。”
- “Ingress 配了但是没有转到后端服务。”

推荐工具顺序：

1. `kubernetes_get_resource` 查 `Ingress` 或 `Gateway`。
2. `kubernetes_get_resource` 查后端 `Service`。
3. `kubernetes_list_resources_summary` 查 `EndpointSlice`。
4. `kubernetes_get_recent_events` 查 Ingress/Gateway controller 事件。
5. `kubernetes_get_pod_logs` 查 ingress controller 或 gateway controller 日志。
6. `prometheus_query_range` 查入口 4xx/5xx、后端 5xx。
7. `jaeger_get_traces_summary` 查请求是否进入应用。

判断方法：

- 404：host/path 没匹配，或 controller 未加载规则。
- 502：后端连接失败、端口错误、TLS upstream 错误。
- 503：Service 无可用 endpoint 或后端全 NotReady。
- TLS 证书错误：Secret、Certificate、Issuer、域名不匹配。

修复建议：

- 修正 host/path、backend service name、service port。
- 修复后端 Pod readiness。
- 修正 TLS Secret 或 cert-manager 资源。
- controller 配置变更前先查当前对象和事件，确认影响范围。

验证：

- `kubernetes_get_recent_events`
- `kubernetes_list_resources_summary` for EndpointSlice
- `prometheus_query_range` 查入口错误率下降
- `loki_query_logs_summary` 查入口 controller 错误消失

## 4. TLS 证书过期或证书签发失败

用户提问示例：

- “Ingress HTTPS 证书过期了。”
- “cert-manager 一直签不出证书。”

推荐工具顺序：

1. `kubernetes_get_resource` 查 `Ingress` 的 TLS Secret 引用。
2. `kubernetes_get_resource` 查 `Secret` 元数据。
3. `kubernetes_list_resources_summary` 查 `Certificate`、`Issuer`、`ClusterIssuer`。
4. `kubernetes_get_resource` 查具体 `Certificate` 状态。
5. `kubernetes_get_recent_events` 查 cert-manager 事件。
6. `kubernetes_get_pod_logs` 查 cert-manager controller 日志。

判断方法：

- Secret 不存在或证书过期：证书未签发或未更新。
- Certificate `Ready=False`：看 reason 和 message。
- ACME HTTP-01 失败：Ingress path、DNS、外部访问问题。
- DNS-01 失败：DNS provider 凭证或权限。

修复建议：

- 修正 Issuer/ClusterIssuer。
- 修正 DNS 或 HTTP-01 solver。
- 修正 Secret 引用。
- 临时换证或 patch Ingress TLS 必须用户确认。

验证：

- `kubernetes_get_resource` 查 Certificate Ready。
- `kubernetes_get_recent_events` 无新失败事件。
- 入口错误率或 TLS 错误日志下降。

## 5. PVC Pending 或挂载失败

用户提问示例：

- “Pod 因为 PVC 一直起不来。”
- “volume mount failed，帮我看是存储还是权限问题。”

推荐工具顺序：

1. `kubernetes_get_resource_summary` 查 Pod。
2. `kubernetes_get_recent_events` 查 FailedMount、FailedAttachVolume、ProvisioningFailed。
3. `kubernetes_get_resource` 查 `PersistentVolumeClaim`。
4. `kubernetes_get_resource` 查 `PersistentVolume`。
5. `kubernetes_get_resource` 查 `StorageClass`。
6. `kubernetes_get_pod_logs` 查 CSI controller/node 日志。

判断方法：

- PVC Pending：StorageClass、provisioner、容量、accessMode。
- FailedMount：节点挂载失败、权限、文件系统、Secret。
- Multi-Attach：卷只能挂一个节点，但 Pod 被调度到多个节点。
- PV reclaim policy 风险：删除 PVC 前必须确认数据影响。

修复建议：

- 修正 StorageClass、accessMode、容量。
- 修复 CSI driver 或底层存储。
- 对 StatefulSet 存储问题，不要直接删除 PVC。
- 数据类资源删除必须强确认，并说明不可逆风险。

验证：

- `kubernetes_get_resource_summary` 查 Pod Running。
- `kubernetes_get_resource` 查 PVC Bound。
- `kubernetes_get_recent_events` 无新的 mount/attach 失败。

## 6. StatefulSet 启动或滚动更新卡住

用户提问示例：

- “StatefulSet 第 2 个副本一直起不来。”
- “有状态服务升级卡在一个 Pod 上。”

推荐工具顺序：

1. `kubernetes_get_resource_summary`
2. `kubernetes_get_rollout_status`
3. `kubernetes_list_resources_summary` 查 Pods 和 PVC。
4. `kubernetes_get_recent_events`
5. `kubernetes_get_pod_logs`
6. `kubernetes_get_resource` 查 StatefulSet updateStrategy、volumeClaimTemplates。

判断方法：

- OrderedReady 阻塞：前一个 Pod 未 Ready，后续不会继续。
- PVC 或 PV 问题：Pod 无法挂载。
- readiness 失败：应用未达到有状态服务健康条件。
- 分片或主从切换问题：需要应用层日志和指标支持。

修复建议：

- 先恢复阻塞的 ordinal Pod。
- 不要盲目删除 PVC 或强制并行更新。
- 明确坏版本时，按 Helm/Argo CD 发布回退流程。

验证：

- `kubernetes_get_rollout_status`
- `kubernetes_list_resources_summary`
- `prometheus_query_range` 查有状态服务健康指标

## 7. HPA 不扩容或异常扩缩容

用户提问示例：

- “流量上来了 HPA 没扩容。”
- “HPA 一直抖动，副本数忽高忽低。”

推荐工具顺序：

1. `kubernetes_get_resource` 查 `HorizontalPodAutoscaler`。
2. `kubernetes_get_recent_events` 查 HPA 事件。
3. `kubernetes_get_resource_summary` 查目标 Deployment/StatefulSet。
4. `prometheus_query_range` 查 CPU、内存、业务指标。
5. `prometheus_targets_summary` 查 metrics pipeline。
6. `opentelemetry_get_collector_summary` 或相关 collector 状态。

判断方法：

- metrics unavailable：metrics-server、Prometheus Adapter 或自定义指标缺失。
- current utilization 低但业务拥塞：HPA 指标不代表瓶颈。
- maxReplicas 太低：达到上限。
- scale behavior 设置过激：抖动。

修复建议：

- 修复指标链路。
- 调整 minReplicas、maxReplicas、target utilization 或 behavior。
- 对紧急容量问题，可在用户确认下临时 `kubernetes_scale_resource`，并说明 HPA 可能覆盖。

验证：

- `kubernetes_get_resource` 查 HPA current/desired。
- `kubernetes_get_rollout_status`
- `prometheus_query_range` 查容量和错误率。

## 8. PDB 阻止驱逐或节点 drain 卡住

用户提问示例：

- “drain 节点卡住了，是 PDB 吗？”
- “为什么 Pod 一直不能被驱逐？”

推荐工具顺序：

1. `kubernetes_get_node_conditions`
2. `kubernetes_list_resources_summary` 查节点上的 Pod。
3. `kubernetes_list_resources_summary` 查 namespace 内 `PodDisruptionBudget`。
4. `kubernetes_get_resource` 查相关 PDB。
5. `kubernetes_get_recent_events`

判断方法：

- PDB `disruptionsAllowed=0`：驱逐被保护。
- 副本数太低或可用副本不足：PDB 无法允许中断。
- Pod 非控制器管理或有 local storage：drain 风险更高。

修复建议：

- 先扩容工作负载，再 drain。
- 临时调整 PDB 需要用户确认，并说明可用性风险。
- 不要强制删除关键 Pod，除非用户明确接受风险。

验证：

- `kubernetes_get_resource` 查 PDB 状态。
- `kubernetes_drain_node` 后用 `kubernetes_get_node_conditions` 和 Pod summary 验证。

## 9. Job 或 CronJob 失败

用户提问示例：

- “定时任务没跑起来。”
- “Job 一直失败重试，帮我看原因。”

推荐工具顺序：

1. `kubernetes_get_resource_summary` 查 `CronJob` 或 `Job`。
2. `kubernetes_get_resource` 查 schedule、suspend、backoffLimit、activeDeadlineSeconds。
3. `kubernetes_list_resources_summary` 查 Job 生成的 Pod。
4. `kubernetes_get_recent_events`
5. `kubernetes_get_pod_logs`

判断方法：

- CronJob suspend=true：不会调度。
- schedule/timezone 错误：触发时间不符合预期。
- startingDeadlineSeconds 太短：错过调度。
- Job Pod 失败：镜像、命令、配置、权限或依赖。
- backoffLimit 达到：Job 终止。

修复建议：

- 修正 schedule、suspend、deadline、concurrencyPolicy。
- 修正 Job template。
- 手动创建补偿 Job 属于变更操作，需要用户确认。

验证：

- `kubernetes_get_resource_summary`
- `kubernetes_list_resources_summary` 查 Job/Pod。
- `kubernetes_get_pod_logs`

## 10. DaemonSet 节点覆盖不完整

用户提问示例：

- “日志采集 DaemonSet 没跑满所有节点。”
- “为什么有些节点没有 node exporter Pod？”

推荐工具顺序：

1. `kubernetes_get_resource_summary` 查 `DaemonSet`。
2. `kubernetes_get_resource` 查 selector、template、tolerations、nodeSelector、affinity。
3. `kubernetes_list_resources_summary` 查 DaemonSet Pods。
4. `kubernetes_get_node_conditions`
5. `kubernetes_get_recent_events`

判断方法：

- 节点 taint 未 toleration：Pod 不会调度到该节点。
- nodeSelector/affinity 限制：只覆盖部分节点。
- Pod 在部分节点 CrashLoop：节点环境或权限差异。
- DaemonSet 更新卡住：maxUnavailable 或 Pod readiness 问题。

修复建议：

- 增加 toleration 或调整 nodeSelector。
- 修复节点条件或镜像/权限。
- 对采集组件修复后验证 telemetry 是否恢复。

验证：

- `kubernetes_get_resource_summary`
- `kubernetes_list_resources_summary`
- `prometheus_targets_summary` 或 `loki_query_logs_summary`

## 11. Admission Webhook 阻塞资源创建

用户提问示例：

- “创建 Deployment 被 webhook 拒绝。”
- “集群突然无法创建资源，提示 validating webhook 超时。”

推荐工具顺序：

1. `kubernetes_get_recent_events`
2. `kubernetes_list_resources_summary` 查 `ValidatingWebhookConfiguration` 和 `MutatingWebhookConfiguration`。
3. `kubernetes_get_resource` 查具体 webhook 配置。
4. `kubernetes_get_resource_summary` 查 webhook 后端 Service/Deployment。
5. `kubernetes_get_pod_logs` 查 webhook Pod 日志。
6. `prometheus_query_range` 查 admission latency/error。

判断方法：

- webhook 后端 Service 无 endpoint：所有匹配资源可能阻塞。
- failurePolicy=Fail：后端不可用会拒绝创建。
- timeoutSeconds 过短或后端慢：请求超时。
- caBundle 或 TLS 错误：apiserver 无法调用 webhook。

修复建议：

- 修复 webhook 后端 Pod/Service。
- 临时调整 failurePolicy 或 scope 风险很高，必须用户确认。
- 修复证书或 caBundle。

验证：

- `kubernetes_get_resource_summary` 查 webhook 后端 Ready。
- `kubernetes_get_recent_events` 无新拒绝事件。
- `prometheus_query_range` 查 admission 错误下降。

## 12. Namespace ResourceQuota 或 LimitRange 导致创建失败

用户提问示例：

- “为什么新 Pod 创建失败，提示 quota exceeded？”
- “开发环境不能创建 PVC 了。”

推荐工具顺序：

1. `kubernetes_get_recent_events`
2. `kubernetes_list_resources_summary` 查 `ResourceQuota`。
3. `kubernetes_get_resource` 查具体 ResourceQuota used/hard。
4. `kubernetes_list_resources_summary` 查 `LimitRange`。
5. `kubernetes_get_resource_usage` 查资源使用。

判断方法：

- quota exceeded：namespace 配额达到上限。
- LimitRange 要求 request/limit：Pod spec 不符合。
- PVC 数量或存储容量超限：存储配额问题。

修复建议：

- 调整 workload request/limit。
- 清理无用资源。
- 提高 ResourceQuota 需要用户确认和容量评估。

验证：

- `kubernetes_get_resource` 查 ResourceQuota used/hard。
- `kubernetes_get_resource_summary` 查新资源创建成功。

## 13. Secret 缺失或密钥轮换失败

用户提问示例：

- “应用启动报 secret not found。”
- “数据库密码轮换后服务连不上。”

推荐工具顺序：

1. `kubernetes_get_resource_summary` 查 Pod。
2. `kubernetes_get_recent_events`
3. `kubernetes_get_resource` 查 Pod env/envFrom/volume secret 引用。
4. `kubernetes_get_resource` 查 `Secret` 元数据。
5. `kubernetes_get_pod_logs` 或 `loki_query_logs_summary`
6. `sentry_list_issues_summary` 查应用连接异常。

判断方法：

- Secret 不存在或 key 缺失：配置错误。
- Secret 更新后 Pod 未重启：环境变量仍是旧值。
- 外部密钥同步器失败：同步 controller 日志和事件。
- 应用报认证失败：密钥与外部服务不匹配。

修复建议：

- 修复 Secret 或外部密钥源。
- 对 env 引用的 Secret，用户确认后重启工作负载。
- 不在回答中泄露 Secret 值，只报告 key 是否存在、版本或元数据。

验证：

- `kubernetes_get_resource_summary`
- `kubernetes_get_pod_logs`
- `sentry_list_issues_summary`


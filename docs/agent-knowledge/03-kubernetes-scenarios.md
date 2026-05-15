# Kubernetes 常见场景

每个场景都包含用户提问示例、工具顺序、判断方法和修复建议。除非用户明确要求执行，修复建议只停留在方案；真正写操作前必须要求确认。

## 1. Pod CrashLoopBackOff

用户提问示例：

- “prod namespace 里的 api pod 一直 CrashLoop，帮我看原因。”
- “Deployment `checkout` 发布后 Pod 重启，怎么处理？”

推荐工具顺序：

1. `kubernetes_get_resource_summary`：确认目标状态、重启次数、ready 状态。
2. `kubernetes_get_recent_events`：看 probe failed、image pull、OOM、mount、scheduler 事件。
3. `kubernetes_get_resource`：必要时看 container statuses、lastState、env、volume。
4. `kubernetes_get_pod_logs`：查当前容器日志；如果工具支持 previous 日志，应优先看上一次崩溃日志。
5. `kubernetes_get_rollout_status`：如果来自 Deployment/StatefulSet，检查 rollout。
6. `loki_query_logs_summary` 或 `sentry_list_issues_summary`：关联应用错误。

判断方法：

- `OOMKilled`：优先查资源限制、内存指标、近期发布。
- `Error` 且日志有配置错误：查 ConfigMap、Secret、Nacos 配置或环境变量。
- `CrashLoopBackOff` 且 probe failed：检查应用启动时间、readiness/liveness path、端口。
- 发布后才出现：关联 Helm/Argo CD release 和 rollout。

修复建议：

- 配置错误：建议修正 ConfigMap、Secret、Nacos 或 Helm values。
- 探针过严：建议 patch probe 的 timeout、period、failureThreshold 或 startupProbe。
- 内存不足：建议调整 request/limit 或排查内存泄漏。
- 需要临时恢复：用户确认后可 `kubernetes_restart_workload` 或 `kubernetes_scale_resource`，但不能替代根因修复。

验证：

1. `kubernetes_get_rollout_status`
2. `kubernetes_wait_for_resource`
3. `kubernetes_get_recent_events`
4. `kubernetes_get_pod_logs`
5. `prometheus_query_range` 查看重启率、错误率是否下降

## 2. Pod Pending 或调度失败

用户提问示例：

- “新 Pod 一直 Pending，帮我查一下。”
- “为什么 Job 创建了但是没有跑起来？”

推荐工具顺序：

1. `kubernetes_get_resource_summary`
2. `kubernetes_get_recent_events`
3. `kubernetes_get_resource`
4. `kubernetes_get_node_conditions`
5. `kubernetes_get_resource_usage`
6. 必要时 `kubernetes_check_permissions`

判断方法：

- `Insufficient cpu/memory`：资源不足或 request 过大。
- `node(s) had taint`：缺少 toleration。
- `node affinity/selector mismatch`：节点标签不满足。
- `pod has unbound immediate PersistentVolumeClaims`：PVC/PV 绑定问题。
- `quota exceeded`：namespace ResourceQuota 或 LimitRange 限制。

修复建议：

- 调整 request/limit 或副本数。
- 增加 toleration、修改 nodeSelector/affinity。
- 修复 PVC、StorageClass 或容量。
- 调整 namespace quota。

验证：

- `kubernetes_get_resource_summary` 查看 Pod 进入 Running。
- `kubernetes_get_recent_events` 确认无新的 FailedScheduling。

## 3. ImagePullBackOff 或 ErrImagePull

用户提问示例：

- “Pod 拉镜像失败，帮我定位是镜像名还是密钥问题。”
- “发布后所有 Pod 都 ImagePullBackOff。”

推荐工具顺序：

1. `kubernetes_get_resource_summary`
2. `kubernetes_get_recent_events`
3. `kubernetes_get_resource`
4. 如果由 Helm/Argo CD 管理，继续查 `helm_get_release_values` 或 `argocd_get_application_manifests`

判断方法：

- `not found`：镜像 tag、repository 或 registry 路径错误。
- `unauthorized`：imagePullSecrets 或 registry 凭证错误。
- `TLS handshake`、`connection refused`：registry 网络或证书问题。
- 只有部分节点失败：节点到 registry 网络、DNS 或证书差异。

修复建议：

- 修正镜像 tag 或 values。
- 修复 imagePullSecrets。
- 检查节点网络和 registry 可达性。

验证：

- `kubernetes_get_recent_events`
- `kubernetes_get_resource_summary`
- `kubernetes_get_rollout_status`

## 4. Deployment Rollout 卡住

用户提问示例：

- “`payment` Deployment 一直 rollout 不完成。”
- “发布后新版本没接流量，帮我看哪里卡住。”

推荐工具顺序：

1. `kubernetes_get_resource_summary`
2. `kubernetes_get_rollout_status`
3. `kubernetes_get_recent_events`
4. `kubernetes_list_resources_summary` 列 ReplicaSet 和 Pod
5. `kubernetes_get_pod_logs` 查第一个失败 Pod
6. Helm 管理时：`helm_get_release_status`、`helm_get_release_history`
7. Argo CD 管理时：`argocd_get_application`、`argocd_get_application_manifests`

判断方法：

- 新 ReplicaSet 副本无法 Ready：看探针、日志、资源、配置。
- old ReplicaSet 不缩容：看 PDB、readiness、maxUnavailable、finalizer。
- Argo CD OutOfSync：判断 Git 与集群状态差异。
- Helm failed/pending：查看 release history 和 hook 状态。

修复建议：

- 配置或镜像错误：修正发布源。
- 探针错误：调整 probe。
- 资源不足：调资源或扩容节点。
- 明确是坏版本：用户确认后考虑 `helm_rollback_release`。

验证：

- `kubernetes_get_rollout_status`
- `kubernetes_wait_for_resource`
- `prometheus_query_range` 查看 5xx、latency、traffic
- `loki_query_logs_summary` 查看错误是否停止

## 5. Service 503、连接失败或无 Endpoint

用户提问示例：

- “Service `api` 返回 503，帮我查链路。”
- “Ingress 到后端 Service 不通，怎么排？”

推荐工具顺序：

1. `kubernetes_get_resource` 读取 `Service`
2. `kubernetes_list_resources_summary` 用 Service selector 查 Pod
3. `kubernetes_list_resources_summary` 查 `EndpointSlice`
4. `kubernetes_get_recent_events`
5. `kubernetes_get_resource_summary` 查 backing workload
6. `loki_query_logs_summary`
7. `jaeger_get_traces_summary`、`jaeger_get_trace`
8. `prometheus_query_range` 查 5xx、latency、request rate

判断方法：

- Service selector 不匹配 Pod labels：无 endpoint。
- Pod not Ready：EndpointSlice 不会包含可用 endpoint。
- targetPort 错误：Service 指向不存在或错误端口。
- Ingress/Gateway 配置错误：入口规则没路由到正确 Service。
- 应用内部 503：Kubernetes endpoint 正常，但日志/trace 显示下游失败。

修复建议：

- 修正 Service selector 或 Pod labels。
- 修正 targetPort 或 containerPort。
- 修复 Pod readiness。
- 修正 Ingress/Gateway 路由。
- 如果是下游依赖故障，继续排查对应下游服务。

验证：

- `kubernetes_list_resources_summary` 查 EndpointSlice 有 ready endpoints。
- `prometheus_query_range` 查 5xx 下降。
- `jaeger_get_traces_summary` 查失败 hop 恢复。

## 6. 节点 NotReady 或资源压力

用户提问示例：

- “集群里有节点 NotReady，影响了哪些工作负载？”
- “节点内存压力导致 Pod 驱逐，帮我排查。”

推荐工具顺序：

1. `kubernetes_get_node_conditions`
2. `kubernetes_get_recent_events`
3. `kubernetes_get_resource_usage`
4. `kubernetes_list_resources_summary` 列受影响节点上的 Pod
5. `prometheus_query_range` 查 node CPU、memory、disk、network

判断方法：

- `MemoryPressure`：内存不足或 Pod limit 配置问题。
- `DiskPressure`：镜像、容器日志、emptyDir、磁盘容量。
- `NetworkUnavailable`：CNI 或节点网络问题。
- `Ready=False/Unknown`：kubelet、节点心跳、底层机器问题。

修复建议：

- 临时隔离：用户确认后 `kubernetes_cordon_node`。
- 维护排空：用户确认后 `kubernetes_drain_node`。
- 恢复调度：修复后 `kubernetes_uncordon_node`。
- 业务修复：调整 Pod request/limit、日志策略、存储。

验证：

- `kubernetes_get_node_conditions`
- `kubernetes_get_recent_events`
- `kubernetes_get_resource_usage`

## 7. OOMKilled 或 CPU Throttling

用户提问示例：

- “服务频繁 OOMKilled，帮我判断要不要加内存。”
- “接口变慢是不是 CPU throttling？”

推荐工具顺序：

1. `kubernetes_get_resource_summary`
2. `kubernetes_get_resource`
3. `kubernetes_get_resource_usage`
4. `prometheus_query_range` 查询容器内存、CPU、重启、throttling 指标
5. `loki_query_logs_summary`

判断方法：

- 内存曲线接近 limit 后重启：大概率 OOMKilled。
- CPU usage 接近 limit 且 throttling 高：CPU limit 过低或应用需要优化。
- request 太低：调度时资源保障不足。

修复建议：

- 调整 request/limit。
- 优化应用内存或 CPU 热点。
- 对突发型服务考虑 HPA 或副本扩展。

验证：

- `kubernetes_get_resource_summary`
- `prometheus_query_range` 查看 OOM、restart、throttling 是否下降。

## 8. RBAC 权限不足

用户提问示例：

- “这个 service account 为什么不能 list pods？”
- “应用报 forbidden，帮我查 RBAC。”

推荐工具顺序：

1. `kubernetes_get_recent_events`
2. `kubernetes_get_resource` 查 ServiceAccount、Role、RoleBinding、ClusterRole、ClusterRoleBinding
3. `kubernetes_check_permissions`
4. `loki_query_logs_summary` 或 `kubernetes_get_pod_logs` 查看应用报错

判断方法：

- `forbidden` 明确指向 verb/resource/apiGroup。
- RoleBinding 绑定错 namespace 或 subject。
- ClusterRole 缺少对应资源或 verb。

修复建议：

- 最小权限增加对应 verb/resource。
- 修改 RoleBinding subject 或 namespace。
- 不要直接给 cluster-admin，除非用户明确要求并接受风险。

验证：

- `kubernetes_check_permissions`
- 应用日志不再出现 forbidden。

## 9. ConfigMap、Secret 或环境变量错误

用户提问示例：

- “发布后应用启动失败，怀疑配置有问题。”
- “Secret 改了为什么服务没生效？”

推荐工具顺序：

1. `kubernetes_get_resource_summary` 查工作负载和 Pod。
2. `kubernetes_get_resource` 查 Pod env、envFrom、volumeMounts。
3. `kubernetes_get_resource` 查 ConfigMap 或 Secret 元数据。
4. `kubernetes_get_pod_logs`
5. 若使用 Nacos：`nacos_list_configs_summary`、`nacos_get_config`
6. 若由 Helm 管理：`helm_get_release_values`

判断方法：

- 环境变量缺失或 key 不存在。
- Secret 更新后 Pod 未重启，应用仍使用旧环境变量。
- 配置内容与发布 values 不一致。
- Nacos namespace/group/dataId 错误。

修复建议：

- 修正配置源。
- 对只在启动时读取的配置，用户确认后重启工作负载。
- 对 GitOps/Helm 管理资源，长期修复应改 Git 或 values。

验证：

- `kubernetes_get_rollout_status`
- `kubernetes_get_pod_logs`
- Nacos 场景用 `nacos_get_config` 再核对。


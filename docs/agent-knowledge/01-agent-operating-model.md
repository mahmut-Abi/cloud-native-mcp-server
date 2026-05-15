# Agent 操作模型

## 1. 请求分类

收到用户问题后，先分类再调用工具。

| 类型 | 用户常见说法 | Agent 行为 |
| --- | --- | --- |
| 解释 | “这个告警是什么意思？”“为什么 Pod 会 Pending？” | 可以先解释概念；如果涉及当前环境，继续读 MCP 数据验证 |
| 只读诊断 | “帮我看看 prod api 为什么 503” | 只使用 summary、health、list、query、logs、trace 等读工具 |
| 变更请求 | “重启这个 Deployment”“把副本数改成 5” | 先读当前状态，确认目标和风险，要求明确确认，再执行 |
| 验证请求 | “确认刚才恢复了吗？” | 使用与故障同源的验证工具，例如 rollout、metrics、logs、alerts |
| 恢复请求 | “帮我恢复服务”“回滚这个版本” | 先诊断，再提出最小恢复动作；回滚、删除、patch 前必须确认 |

## 2. 默认诊断流程

1. 明确症状：CrashLoop、Pending、503、延迟、无指标、无日志、告警、发布失败等。
2. 明确范围：cluster、namespace、workload、pod、service、release、application、dashboard、alert、trace、index。
3. 先用 summary 或 health 工具建立快照。
4. 找到最小失败单元：一个 Pod、一组 Pods、一个 Service、一个 release、一个 alert、一个 query、一个 trace。
5. 至少关联两个信号后再下结论。
6. 输出事实和推断，说明下一步最小工具调用。
7. 只有在证据充分且用户确认后执行修复。

## 3. 读工具优先级

优先级从高到低：

1. health、summary、paginated
2. list、search、query
3. single object detail
4. full object、manifest、raw logs、large range query
5. mutation

示例：

- 查 Pod：先 `kubernetes_get_resource_summary`，不够再 `kubernetes_get_resource`
- 查日志：先 `loki_query_logs_summary`，不够再 `loki_query_range`
- 查 release：先 `helm_get_release_summary` 或 `helm_get_release_status`，不够再 `helm_get_release_manifest`
- 查 dashboard：先 `grafana_dashboards_summary` 或 `grafana_search_dashboards`，不够再 `grafana_dashboard`

## 4. 参数规则

优先使用扁平 JSON：

```json
{
  "kind": "Deployment",
  "name": "api",
  "namespace": "prod"
}
```

Kubernetes 注意事项：

- `kind` 用单数 Kind，例如 `Pod`、`Deployment`、`StatefulSet`、`Service`、`EndpointSlice`、`Ingress`、`ConfigMap`、`Secret`、`Node`。
- 用户说“deployments”“pods”“svc”时，转换成 `Deployment`、`Pod`、`Service`。
- 名称不确定时先用 `kubernetes_search_resources` 或 `kubernetes_list_resources_summary`。
- 查工作负载对应 Pod 时，从 workload selector 推导 label selector，再列 Pod。

时间范围注意事项：

- 用户说“最近半小时”，agent 应尽量转成明确窗口。
- Prometheus range query、Jaeger trace、OpenTelemetry query 优先使用 RFC3339 起止时间。
- 如果需要当前时间，使用 `utilities_get_time`、`utilities_get_timestamp` 或 `utilities_get_date`。

## 5. 结果解析规则

不要假设工具返回值一定是 JSON 字符串。

可能形态：

- 已解析对象或数组
- MCP envelope，payload 在 `content[0].text`
- JSON 字符串
- 文本摘要

处理顺序：

1. 先检查原始返回类型。
2. 如果已经是 object 或 array，直接使用。
3. 如果是 envelope，取 `content[0].text` 后再判断是否需要解析。
4. 如果出现 `[object Object]` 或 `Unexpected token o`，通常是重复 parse。

## 6. 变更安全规则

写操作包括但不限于：

- `kubernetes_patch_resource`
- `kubernetes_restart_workload`
- `kubernetes_scale_resource`
- `kubernetes_delete_resource`
- `kubernetes_cordon_node`
- `kubernetes_uncordon_node`
- `kubernetes_drain_node`
- `helm_install_release`
- `helm_upgrade_release`
- `helm_rollback_release`
- `helm_uninstall_release`
- `alertmanager_create_silence`
- `alertmanager_delete_silence`
- `grafana_update_dashboard`
- `grafana_update_datasource`
- `kibana_update_dashboard`
- `kibana_delete_saved_object`

写操作前：

1. 读取当前状态。
2. 确认目标名称、namespace、集群或外部系统 scope。
3. 说明将执行的最小变更。
4. 说明风险和回滚方式。
5. 要求用户明确确认。

写操作后：

1. 用读工具验证对象状态。
2. 用 rollout、wait、metrics、logs、trace 或 alert 验证业务状态。
3. 如果验证失败，提出最小回滚路径。

## 7. 推荐回答格式

只读诊断：

```text
已观察事实：
- ...

推断：
- ...

证据：
- ...

下一步：
- 建议调用 <tool_name> 验证 <假设>。
```

修复确认：

```text
我将对 <namespace>/<kind>/<name> 执行 <action>。
当前状态是 <summary>。
预期影响是 <impact>，回滚方式是 <rollback>。
请确认是否执行。
```

修复完成：

```text
已执行 <tool_name>。
验证结果：
- <verification>

仍需关注：
- <residual risk>
```


# 云原生 MCP Agent 知识库

本目录用于给云原生运维类 agent 提供可导入知识库的场景、工具顺序和修复方法。内容基于当前项目的 MCP 工具清单整理；在线执行时以运行时 `tools/list` 为准，离线参考 `docs/TOOLS.md`。

## 建议导入顺序

1. `SKILL.md`
   - 作为 agent 的入口规则，定义“只通过 MCP 工具排查和操作”的基本约束。
2. `01-agent-operating-model.md`
   - 定义问题分类、读优先、安全变更、结果解析和回答格式。
3. `02-tool-routing-and-sequences.md`
   - 定义各类云原生问题该先用哪些工具，以及跨工具排查顺序。
4. `03-kubernetes-scenarios.md`
   - Kubernetes 常见故障场景：Pod、Deployment、Service、Node、RBAC、资源压力等。
5. `04-observability-scenarios.md`
   - Prometheus、Loki、Jaeger、Alertmanager、Grafana、Sentry、Langfuse、OpenTelemetry 场景。
6. `05-release-config-data-scenarios.md`
   - Helm、Argo CD、Nacos、Elasticsearch、Kibana 等发布、配置、数据平台场景。
7. `06-remediation-templates.md`
   - 变更修复模板：patch、restart、scale、delete、rollback、silence、dashboard/datasource 更新。
8. `07-question-bank.md`
   - 用户提问样例与推荐工具链，适合作为检索增强的问答索引。
9. `08-kubernetes-advanced-scenarios.md`
   - Kubernetes 网络、DNS、Ingress、存储、HPA、PDB、Job、Webhook、证书等进阶场景。
10. `09-platform-observability-deep-dive.md`
   - 可观测性和平台组件的深度排查：告警抖动、通知失败、高基数、日志延迟、trace 缺口等。
11. `10-remediation-decision-trees.md`
   - 根据症状选择修复动作的决策树，帮助 agent 避免盲目重启、回滚或删除。

## Agent 使用原则

- 优先使用当前 MCP server 暴露的工具，不使用 `kubectl`、`helm` CLI、直接 HTTP API 或数据库直连来替代 MCP 工具。
- 每次回答先判断用户请求类型：解释、只读诊断、变更请求、验证请求、恢复请求。
- 诊断先读后写，优先 summary、health、paginated 工具，再按需读取 full/detail 工具。
- Kubernetes 资源类工具通常需要 `kind`；用户说“pods/deployments”时要转换成 `Pod`、`Deployment`。
- 变更类动作必须在读当前状态后执行，并在执行前要求明确确认；执行后必须用第二个读工具验证。
- Prometheus、Jaeger、OpenTelemetry 等涉及时间窗口时，尽量使用明确时间范围或 RFC3339 时间。
- 工具返回值可能是已解析对象，也可能是 MCP envelope 中的 `content[0].text`；解析前先检查原始结构，避免重复 `JSON.parse`。

## 回答结构建议

诊断类回答：

1. 已观察事实
2. 最可能原因
3. 证据链
4. 下一步建议或下一个 MCP 工具调用
5. 是否需要用户确认执行修复

修复类回答：

1. 目标对象和作用范围
2. 当前状态摘要
3. 将执行的最小变更
4. 风险和回滚路径
5. 执行结果
6. 验证结果

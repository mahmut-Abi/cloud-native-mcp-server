# 平台与可观测性深度场景

本文件补充告警、指标、日志、链路、Dashboard、Elastic、LLM 观测和 OpenTelemetry 的深度问题。目标是帮助 agent 从“看到现象”推进到“证明原因”。

## 1. 告警抖动或频繁恢复再触发

用户提问示例：

- “这个告警一直 firing/resolved 抖动。”
- “HighLatency 每隔几分钟就报警一次，怎么排？”

推荐工具顺序：

1. `alertmanager_alerts_summary`
2. `alertmanager_query_alerts_advanced`
3. `prometheus_rules_summary`
4. `prometheus_query_range`
5. `grafana_dashboards_summary` 查相关看板。
6. `loki_query_logs_summary` 和 `jaeger_get_traces_summary` 查业务信号。

判断方法：

- 指标在阈值附近波动：阈值或 for duration 不合理。
- scrape 间歇失败：target up/down 导致表达式抖动。
- 业务周期性尖峰：容量或限流策略问题。
- Alertmanager 分组和 repeat interval 设置不合理：通知噪声。

修复建议：

- 调整告警阈值、`for` 时长、聚合窗口。
- 修复 scrape target。
- 对真实周期性尖峰，扩容或限流。
- 告警规则修改应走配置发布流程，不直接在运行时随意改。

验证：

- `prometheus_query_range`
- `alertmanager_alerts_summary`
- `grafana_get_dashboard_panel_queries`

## 2. 告警未通知

用户提问示例：

- “Prometheus 已经 firing，但是没有收到通知。”
- “为什么这个告警没有发到 Slack/钉钉？”

推荐工具顺序：

1. `prometheus_alerts_summary`
2. `alertmanager_alerts_summary`
3. `alertmanager_silences_summary`
4. `alertmanager_receivers_summary`
5. `alertmanager_get_receivers`
6. `alertmanager_test_receiver`

判断方法：

- Prometheus firing 但 Alertmanager 无告警：Prometheus 到 Alertmanager 配置或网络问题。
- Alertmanager 有告警但被 silence：matcher 匹配。
- route 不匹配 receiver：标签或 routing tree 问题。
- receiver 测试失败：通知集成配置或外部平台问题。

修复建议：

- 修复 Alertmanager 路由或 receiver 配置。
- 删除错误 silence 需用户确认。
- 通知渠道 token/webhook 更改需走密钥和配置流程。

验证：

- `alertmanager_alerts_summary`
- `alertmanager_silences_summary`
- `alertmanager_test_receiver`

## 3. Prometheus 查询慢或高基数风险

用户提问示例：

- “PromQL 查询很慢，Grafana panel 超时。”
- “Prometheus 内存飙升，怀疑高基数。”

推荐工具顺序：

1. `prometheus_get_runtime_info`
2. `prometheus_get_tsdb_status`
3. `prometheus_get_tsdb_stats`
4. `prometheus_query`
5. `prometheus_query_range`
6. `prometheus_get_series`
7. `grafana_get_dashboard_panel_queries`

判断方法：

- series 数量异常增长：label 高基数。
- query 使用裸 metric 大范围聚合：查询代价大。
- Grafana panel 时间范围过大或 step 太细：返回数据过多。
- recording rule 缺失：重复执行昂贵表达式。

修复建议：

- 降低 label 基数，避免 user_id、request_id、trace_id 等进入 metric label。
- 给大查询增加聚合、过滤和合理 step。
- 添加 recording rule。
- 调整 dashboard panel 查询。

验证：

- `prometheus_get_tsdb_status`
- `prometheus_query_range`
- `grafana_get_dashboard_panel_queries`

## 4. Prometheus 规则不生效

用户提问示例：

- “规则文件改了但告警没出现。”
- “recording rule 没有生成新指标。”

推荐工具顺序：

1. `prometheus_rules_summary`
2. `prometheus_get_rules`
3. `prometheus_query`
4. `prometheus_get_metrics_metadata`
5. `prometheus_targets_summary`
6. `kubernetes_get_pod_logs` 查 Prometheus 或 rule reload 日志。

判断方法：

- rule group 不存在：配置未加载。
- rule evaluation error：PromQL 错误或数据缺失。
- recording metric 存在但 label 不符合预期：规则表达式问题。
- reload 失败：配置语法错误或挂载未更新。

修复建议：

- 修正规则配置并重新发布。
- 修正 PromQL。
- 修复 Prometheus reload 或配置挂载。

验证：

- `prometheus_rules_summary`
- `prometheus_query`
- `prometheus_get_rules`

## 5. Loki 日志写入延迟或缺口

用户提问示例：

- “日志延迟好几分钟才出现。”
- “某个时间段日志断了。”

推荐工具顺序：

1. `loki_test_connection`
2. `loki_query_logs_summary`
3. `loki_get_series`
4. `kubernetes_get_resource_summary` 查日志采集器。
5. `kubernetes_get_pod_logs` 查 promtail/fluent-bit/otel collector。
6. `prometheus_query_range` 查采集器队列、丢弃、重试指标。
7. `opentelemetry_get_collector_summary` 查 collector。

判断方法：

- 采集器队列积压：后端慢或网络问题。
- 某些 Pod 无日志：采集器 selector 或权限问题。
- 某些 label 缺失：pipeline parse 或 relabel 配置问题。
- 后端限流：日志量、tenant、ingester 压力。

修复建议：

- 修复采集器配置和权限。
- 调整 batch、buffer、重试、限流。
- 降低日志量或拆分租户。
- 修复 Loki 后端容量。

验证：

- `loki_query_logs_summary`
- `prometheus_query_range`
- `kubernetes_get_pod_logs`

## 6. Loki 查询结果过多或费用异常

用户提问示例：

- “LogQL 查询太慢，返回太多日志。”
- “日志成本突然变高。”

推荐工具顺序：

1. `loki_get_label_names`
2. `loki_get_label_values`
3. `loki_get_series`
4. `loki_query_logs_summary`
5. `prometheus_query_range` 查日志写入速率。
6. `grafana_get_dashboard_panel_queries` 查 dashboard 是否有重查询。

判断方法：

- 查询 selector 过宽：扫描大量 streams。
- 日志标签高基数：stream 数爆炸。
- 应用 debug 日志打开：写入量激增。
- Dashboard 自动刷新过频：查询压力高。

修复建议：

- 收窄 LogQL selector。
- 降低日志级别。
- 修正 pipeline label，避免高基数字段作为 label。
- 调整 dashboard 刷新频率。

验证：

- `loki_get_series`
- `loki_query_logs_summary`
- `prometheus_query_range`

## 7. Jaeger 没有 trace 或 trace 不完整

用户提问示例：

- “服务明明有请求，但是 Jaeger 查不到 trace。”
- “trace 里缺少下游服务 span。”

推荐工具顺序：

1. `jaeger_get_services_summary`
2. `jaeger_get_service_ops`
3. `jaeger_get_traces_summary`
4. `opentelemetry_get_collector_summary`
5. `opentelemetry_get_config_summary`
6. `opentelemetry_analyze_pipeline_status`
7. `loki_query_logs_summary` 查 SDK 或 collector 错误。

判断方法：

- service 不存在：应用未接入 SDK、service.name 错误、collector 未收。
- trace 不完整：上下文传播断裂、采样、异步边界未注入。
- 只有入口 span：下游未 instrumentation 或 propagation header 被丢。
- collector exporter 错误：后端不可达。

修复建议：

- 修正 instrumentation 和 `service.name`。
- 修复 propagation header。
- 调整 sampling。
- 修复 OTel collector pipeline 或 Jaeger exporter。

验证：

- `jaeger_get_services_summary`
- `jaeger_get_traces_summary`
- `opentelemetry_query_traces`

## 8. Grafana 变量或时间范围导致误判

用户提问示例：

- “Dashboard 看起来没数据，但 Prometheus 查有数据。”
- “同一个 panel 换变量后就空了。”

推荐工具顺序：

1. `grafana_dashboards_summary`
2. `grafana_dashboard`
3. `grafana_get_dashboard_panel_queries`
4. `grafana_datasources_summary`
5. `grafana_check_datasource_health`
6. `prometheus_query` 或 `loki_query_logs_summary`

判断方法：

- 变量 query 返回空：label、regex、datasource 错误。
- panel 时间范围覆盖 dashboard 时间：查询窗口不一致。
- datasource UID 变更：panel 指向旧 datasource。
- PromQL 使用变量时没有默认值或多选语法错误。

修复建议：

- 修正变量 query、regex、默认值。
- 修正 panel query 和 datasource UID。
- 修改 dashboard 前获取版本，必要时可恢复版本。

验证：

- `grafana_get_dashboard_panel_queries`
- `grafana_check_datasource_health`
- 对应后端 query 工具

## 9. Elasticsearch 写入阻塞或只读

用户提问示例：

- “日志写不进 Elasticsearch。”
- “index 被设置成 read_only_allow_delete 了。”

推荐工具顺序：

1. `elasticsearch_cluster_health_summary`
2. `elasticsearch_nodes_summary`
3. `elasticsearch_indices_summary`
4. `elasticsearch_index_stats`
5. `elasticsearch_get_cluster_detail_advanced`
6. `kibana_query_logs` 或采集器日志

判断方法：

- 磁盘超过 flood stage watermark：index 自动只读。
- primary shard 不可用：写入失败。
- ingest pipeline 报错：文档被拒绝。
- ILM 或 rollover 失败：写入 alias 指向异常。

修复建议：

- 释放磁盘或扩容。
- 修复 shard allocation。
- 修复 ingest pipeline。
- 解除只读设置属于变更动作，必须确认且先解决磁盘根因。

验证：

- `elasticsearch_cluster_health_summary`
- `elasticsearch_indices_summary`
- `kibana_query_logs`

## 10. Kibana Connector 或 Alert 执行失败

用户提问示例：

- “Kibana alert 触发了但 webhook 没收到。”
- “Connector 测试失败，帮我看原因。”

推荐工具顺序：

1. `kibana_health_summary`
2. `kibana_get_alert_rules`
3. `kibana_get_alert_rule`
4. `kibana_get_connectors`
5. `kibana_get_connector`
6. `kibana_test_connector`
7. `kibana_get_alert_rule_history`

判断方法：

- rule disabled/muted：不执行。
- connector auth/url 错误：通知发送失败。
- action frequency 或 throttle：通知少于预期。
- Kibana task manager 异常：rule 调度延迟。

修复建议：

- enable/unmute rule 前要求确认。
- 修正 connector 配置和密钥。
- 调整 action frequency。

验证：

- `kibana_get_alert_rule_history`
- `kibana_test_connector`
- `kibana_get_alerts`

## 11. Sentry issue 没有关联 release 或 commit

用户提问示例：

- “Sentry issue 看不到 release，没法判断哪个版本引入。”
- “source map 好像没生效。”

推荐工具顺序：

1. `sentry_test_connection`
2. `sentry_list_projects`
3. `sentry_list_issues_summary`
4. `sentry_get_issue`
5. `sentry_list_issue_events`
6. `sentry_get_issue_event`
7. 关联 `kubernetes_get_resource` 查镜像 tag、env release。

判断方法：

- event 没有 release 字段：SDK release 配置缺失。
- stack 未符号化或 sourcemap 缺失：构建产物上传问题。
- environment 不一致：查询范围或配置错误。

修复建议：

- 修正 SDK release/environment 配置。
- 修复 sourcemap 或符号文件上传流程。
- 将镜像 tag、git sha 和 Sentry release 对齐。

验证：

- `sentry_list_issue_events`
- `sentry_get_issue_event`
- 新事件是否带 release/environment。

## 12. Langfuse 数据缺失或环境混乱

用户提问示例：

- “Langfuse 没有某个环境的 trace。”
- “同一个会话里的 trace 不完整。”

推荐工具顺序：

1. `langfuse_check_health`
2. `langfuse_list_traces_summary`
3. `langfuse_get_trace`
4. `langfuse_list_sessions`
5. `langfuse_list_observations`
6. `langfuse_get_metrics`
7. `loki_query_logs_summary` 查应用 SDK 错误。

判断方法：

- environment 或 tags 不一致：查询过滤不匹配。
- trace 存在但 observation 缺失：SDK 上报中断或异步失败。
- session 断裂：session id 未传递。
- metrics 异常但 trace 少：采样、flush 或网络问题。

修复建议：

- 修正 SDK environment、session id、trace id。
- 修复 flush、batch、网络和鉴权。
- 对 prompt/score 质量问题继续查 prompt 和 score。

验证：

- `langfuse_list_traces_summary`
- `langfuse_list_sessions`
- `langfuse_get_metrics`


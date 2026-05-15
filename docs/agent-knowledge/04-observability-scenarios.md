# 可观测性场景

本文件覆盖 Prometheus、Loki、Jaeger、Alertmanager、Grafana、Sentry、Langfuse、OpenTelemetry。核心原则是：先确认后端健康，再确认采集链路，最后定位应用或配置问题。

## 1. 告警正在触发

用户提问示例：

- “这个 HighErrorRate 告警为什么触发？”
- “现在 prod 有哪些严重告警？”

推荐工具顺序：

1. `alertmanager_health_summary`
2. `alertmanager_alerts_summary`
3. `alertmanager_query_alerts_advanced`
4. `prometheus_alerts_summary`
5. `prometheus_rules_summary`
6. `prometheus_query` 或 `prometheus_query_range`
7. 关联 `kubernetes_get_unhealthy_resources`、`loki_query_logs_summary`、`jaeger_get_traces_summary`

判断方法：

- Alertmanager 有告警但 Prometheus 没有 firing：可能是状态不同步或旧告警未清理。
- Prometheus alert firing 且指标持续异常：告警是真实症状。
- 有 silence：确认是否预期屏蔽。
- receiver 配置异常：检查 `alertmanager_receivers_summary` 或 `alertmanager_get_receivers`。

修复建议：

- 真实故障：按业务服务继续排查。
- 告警规则误报：建议调整 Prometheus rule，但修改规则需走配置发布流程。
- 临时维护：用户确认后创建 silence。

验证：

- `alertmanager_alerts_summary`
- `prometheus_query_range`
- `alertmanager_silences_summary`

## 2. 指标缺失或 Prometheus 抓取失败

用户提问示例：

- “服务上线了，但是 Prometheus 没有指标。”
- “为什么这个 target down？”

推荐工具顺序：

1. `prometheus_test_connection`
2. `prometheus_targets_summary`
3. `prometheus_get_targets`
4. `prometheus_get_label_names`、`prometheus_get_label_values`
5. `prometheus_get_series`
6. `prometheus_get_target_metadata`
7. `kubernetes_get_resource_summary` 查 exporter 或 service monitor 对应资源
8. `loki_query_logs_summary` 查 exporter/collector 日志

判断方法：

- target down：网络、Service、endpoint、scrape path、TLS、auth。
- target missing：服务发现规则或 label 不匹配。
- target up 但 metric missing：应用未暴露、metric 名变化、采样路径错误。
- 只有 dashboard 无数据：可能是 Grafana datasource、变量或 PromQL 错误。

修复建议：

- 修正 Service labels、scrape annotation、ServiceMonitor、PodMonitor 或 scrape config。
- 修正 exporter path、port、TLS、auth。
- 修正 Grafana panel query 或 datasource。

验证：

- `prometheus_targets_summary`
- `prometheus_query`
- `grafana_check_datasource_health`

## 3. 日志查不到

用户提问示例：

- “Loki 里查不到 `api` 的日志。”
- “Kibana Discover 没有新日志，是采集断了吗？”

推荐工具顺序：

Loki：

1. `loki_test_connection`
2. `loki_get_label_names`
3. `loki_get_label_values`
4. `loki_query_logs_summary`
5. `loki_get_series`
6. `kubernetes_get_resource_summary`、`kubernetes_get_pod_logs`

Kibana：

1. `kibana_health_summary`
2. `kibana_spaces_summary`
3. `kibana_index_patterns_summary` 或 `kibana_get_data_views`
4. `kibana_query_logs`
5. `elasticsearch_indices_summary`

判断方法：

- 应用 Pod 有 stdout 日志但 Loki/Kibana 无日志：采集链路问题。
- Loki label 不符合查询条件：selector 错误。
- Kibana data view 时间字段错误：Discover 无数据。
- Elasticsearch index 缺失或只读：存储侧问题。

修复建议：

- 修正日志采集器 selector、pipeline、index/data stream、tenant。
- 修正 Kibana data view、space 或时间字段。
- 修复 Elasticsearch index 或 ingest pipeline。

验证：

- `loki_query_logs_summary`
- `kibana_query_logs`
- `elasticsearch_indices_summary`

## 4. 延迟升高或 5xx 增加

用户提问示例：

- “最近 30 分钟 checkout 接口变慢，帮我找慢在哪。”
- “API 5xx 飙升，给我一个根因判断。”

推荐工具顺序：

1. `prometheus_query_range` 查错误率、延迟、流量。
2. `jaeger_get_services_summary`
3. `jaeger_get_service_ops`
4. `jaeger_get_traces_summary` 或 `jaeger_search_traces`
5. `jaeger_get_trace`
6. `loki_query_logs_summary`
7. `sentry_list_issues_summary`
8. `kubernetes_get_resource_summary`、`kubernetes_get_rollout_status`

判断方法：

- trace 显示某个 downstream span 慢：定位下游。
- 只有某版本 Pod 错误高：发布回归或配置差异。
- Sentry issue 与时间窗口匹配：应用异常。
- Prometheus 显示 saturation：CPU、内存、连接池或队列瓶颈。

修复建议：

- 下游慢：继续排下游服务、数据库、外部依赖。
- 发布回归：走 release regression 流程。
- 资源瓶颈：调副本、资源或限流配置。

验证：

- `prometheus_query_range`
- `jaeger_get_traces_summary`
- `loki_query_logs_summary`
- `sentry_list_issues_summary`

## 5. Grafana 面板无数据或报错

用户提问示例：

- “Grafana dashboard 空白，是数据源问题还是查询问题？”
- “这个 panel 报错，帮我定位。”

推荐工具顺序：

1. `grafana_test_connection`
2. `grafana_dashboards_summary` 或 `grafana_search_dashboards`
3. `grafana_dashboard`
4. `grafana_get_dashboard_panel_queries`
5. `grafana_datasources_summary`
6. `grafana_check_datasource_health`
7. 根据 datasource 类型关联 Prometheus、Loki、Elasticsearch
8. 必要时 `grafana_render_panel_image`

判断方法：

- datasource health failed：数据源连接、认证、网络问题。
- datasource 正常但 query 无数据：PromQL/LogQL/时间范围/变量错误。
- 只有渲染失败：renderer 或 panel 配置问题。
- dashboard 版本变化：查 `grafana_get_dashboard_versions`。

修复建议：

- 修正 datasource URL、auth、proxy、TLS。
- 修正 panel query、变量或时间范围。
- 恢复 dashboard 版本需用户确认后使用 `grafana_restore_dashboard_version`。

验证：

- `grafana_check_datasource_health`
- `grafana_get_dashboard_panel_queries`
- 后端查询工具验证数据存在

## 6. Alertmanager silence 或通知异常

用户提问示例：

- “为什么这个告警没通知？”
- “帮我给这个维护窗口加 silence。”

推荐工具顺序：

1. `alertmanager_health_summary`
2. `alertmanager_alerts_summary`
3. `alertmanager_silences_summary`
4. `alertmanager_receivers_summary`
5. `alertmanager_get_receivers`
6. `alertmanager_test_receiver`

判断方法：

- 告警被 silence 匹配：说明是预期屏蔽或误屏蔽。
- 告警分组延迟：需要看 group wait/group interval。
- receiver 测试失败：通知渠道配置异常。

修复建议：

- 创建 silence 前确认 matcher、开始/结束时间、创建者、原因。
- 删除 silence 前确认 silence id 和影响范围。
- 修改 receiver 应走配置发布流程。

验证：

- `alertmanager_silences_summary`
- `alertmanager_alerts_summary`
- `alertmanager_test_receiver`

## 7. Sentry 异常突增

用户提问示例：

- “Sentry 上登录接口异常突然变多，帮我查。”
- “这个 issue 影响哪些用户和版本？”

推荐工具顺序：

1. `sentry_test_connection`
2. scope 不明确时：`sentry_list_organizations`、`sentry_list_projects`
3. `sentry_list_issues_summary`
4. `sentry_list_issues`
5. `sentry_get_issue`
6. `sentry_list_issue_events`
7. `sentry_get_issue_event`
8. 关联 `loki_query_logs_summary`、`jaeger_get_trace`、`kubernetes_get_rollout_status`

判断方法：

- issue first seen 与发布窗口重合：发布回归。
- 影响集中在某 environment、release、transaction：范围明确。
- event payload 指向外部依赖或参数：定位请求路径。

修复建议：

- 代码 bug：建议回滚或 hotfix。
- 配置错误：修正配置源。
- 外部依赖异常：关联 trace 和 logs 继续排查。

验证：

- `sentry_list_issues_summary`
- `prometheus_query_range`
- `loki_query_logs_summary`

## 8. Langfuse LLM 质量、延迟或成本异常

用户提问示例：

- “LLM 回答质量变差，帮我看是不是 prompt 版本问题。”
- “最近 token 成本突然升高，查一下是哪条链路。”

推荐工具顺序：

1. `langfuse_check_health`
2. `langfuse_list_traces_summary`
3. `langfuse_get_trace`
4. `langfuse_list_observations`
5. prompt 问题：`langfuse_list_prompts`、`langfuse_get_prompt`
6. 质量问题：`langfuse_list_scores`
7. 成本或延迟：`langfuse_get_metrics`
8. 应用错误关联：`sentry_list_issues_summary`、`loki_query_logs_summary`

判断方法：

- prompt version 或 label 变化后质量下降：prompt 回归。
- token 使用升高：上下文长度、检索结果、模型变化、重试。
- trace latency 高：模型响应慢、工具调用慢、下游服务慢。
- score 下降但错误率没变：质量问题不是可用性问题。

修复建议：

- 回退 prompt label/version 或修正 prompt 内容。
- 限制上下文、调整检索数量、控制重试。
- 调整模型或超时策略。

验证：

- `langfuse_list_scores`
- `langfuse_get_metrics`
- `langfuse_list_traces_summary`

## 9. OpenTelemetry Collector 管道异常

用户提问示例：

- “为什么 trace 没进 Jaeger？”
- “Collector pipeline 有没有配置问题？”

推荐工具顺序：

1. `opentelemetry_get_health`
2. `opentelemetry_get_collector_summary`
3. `opentelemetry_get_config_summary`
4. `opentelemetry_analyze_pipeline_status`
5. `opentelemetry_query_metrics`
6. `opentelemetry_query_logs`
7. `opentelemetry_query_traces`
8. 关联后端：Prometheus、Loki、Jaeger、Langfuse

判断方法：

- receiver 存在但 pipeline 未引用：配置断链。
- exporter 失败：后端不可达、认证或 TLS 问题。
- 缺 batch/memory_limiter：高流量下丢数据或背压。
- sampling 配置变化：trace 数量下降但不是故障。

修复建议：

- 修正 collector 配置。
- 修复 exporter 后端连接。
- 增加 batch、memory_limiter 或合理 sampling。

验证：

- `opentelemetry_analyze_pipeline_status`
- `opentelemetry_query_metrics`
- 对应后端查询工具确认数据恢复


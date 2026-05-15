---
name: cloud-native-mcp-agent-runbook
description: 使用当前 cloud-native-mcp-server 的 MCP 工具排查和修复 Kubernetes、Helm、Argo CD、Prometheus、Loki、Jaeger、Grafana、Alertmanager、Kibana、Elasticsearch、Nacos、Sentry、Langfuse、OpenTelemetry 等云原生问题。
---

# Cloud Native MCP Agent Runbook

## 什么时候使用

当用户询问云原生故障、发布异常、可观测性数据、告警、日志、链路、配置中心、数据平台、LLM 应用观测或需要通过 MCP 工具执行运维动作时，使用本规则。

## 工作方式

1. 先把用户问题转成明确的操作意图：
   - 解释
   - 只读诊断
   - 变更请求
   - 验证请求
   - 恢复请求
2. 优先从运行时 `tools/list` 获取可用工具。若无法获取，参考当前仓库 `docs/TOOLS.md`。
3. 选择最小、只读、summary 优先的工具开始排查。
4. 当第一个信号只能说明现象，必须至少再关联一种信号：
   - Kubernetes 状态与事件
   - Prometheus 指标
   - Loki 或 Kibana 日志
   - Jaeger 或 OpenTelemetry 链路
   - Sentry issue
   - Langfuse trace、score、metrics
   - Alertmanager 告警和 silence
   - Helm/Argo CD 发布状态
5. 写操作前必须先读取当前状态，并要求用户明确确认。
6. 写操作后必须立刻验证，优先使用 rollout、wait、summary、health、metrics、logs、traces。

## MCP 调用约束

- 使用准确的运行时工具名，例如 `kubernetes_get_resource_summary`，不要改成 camelCase 或自造别名。
- 参数优先使用扁平 JSON。
- 对象和数组参数直接传结构化 JSON，不要把对象编码成字符串，除非工具 schema 明确要求。
- Kubernetes 资源工具通常需要 `kind`、`name`、`namespace`。
- Prometheus、Jaeger、OpenTelemetry 的时间字段使用 RFC3339 或工具 schema 支持的时间格式。
- 读取返回值时先看原始结构；如果客户端已经返回对象，不要再次 JSON parse。

## 优先工具

- Kubernetes 总览：`kubernetes_get_unhealthy_resources`、`kubernetes_get_recent_events`
- Kubernetes 单对象：`kubernetes_get_resource_summary`，必要时 `kubernetes_get_resource`
- Kubernetes 日志：`kubernetes_get_pod_logs`
- Kubernetes 发布：`kubernetes_get_rollout_status`、`kubernetes_wait_for_resource`
- Prometheus：`prometheus_targets_summary`、`prometheus_alerts_summary`、`prometheus_query`、`prometheus_query_range`
- Loki：`loki_query_logs_summary`，必要时 `loki_query` 或 `loki_query_range`
- Jaeger：`jaeger_get_services_summary`、`jaeger_get_traces_summary`、`jaeger_get_trace`
- Alertmanager：`alertmanager_health_summary`、`alertmanager_alerts_summary`、`alertmanager_silences_summary`
- Helm：`helm_list_releases_paginated`、`helm_get_release_summary`、`helm_get_release_status`、`helm_get_release_history`
- Argo CD：`argocd_test_connection`、`argocd_list_applications_summary`、`argocd_get_application`
- Grafana：`grafana_dashboards_summary`、`grafana_datasources_summary`、`grafana_check_datasource_health`
- Kibana：`kibana_health_summary`、`kibana_query_logs`、`kibana_dashboards_summary`
- Elasticsearch：`elasticsearch_cluster_health_summary`、`elasticsearch_nodes_summary`、`elasticsearch_indices_summary`
- Nacos：`nacos_test_connection`、`nacos_list_configs_summary`、`nacos_list_services_summary`
- Sentry：`sentry_test_connection`、`sentry_list_issues_summary`、`sentry_get_issue`
- Langfuse：`langfuse_check_health`、`langfuse_list_traces_summary`、`langfuse_get_trace`、`langfuse_list_scores`、`langfuse_list_organization_projects`、`langfuse_list_project_memberships`
- OpenTelemetry：`opentelemetry_get_collector_summary`、`opentelemetry_get_config_summary`、`opentelemetry_analyze_pipeline_status`

## 输出格式

诊断时输出：

- 事实
- 推断
- 证据
- 下一步工具
- 需要确认的修复动作

执行修复后输出：

- 执行了什么
- 影响范围
- 验证结果
- 仍需关注的问题

## 知识库导航

- 通用操作模型：`01-agent-operating-model.md`
- 工具选型和执行顺序：`02-tool-routing-and-sequences.md`
- Kubernetes 基础场景：`03-kubernetes-scenarios.md`
- 可观测性基础场景：`04-observability-scenarios.md`
- 发布、配置和数据平台场景：`05-release-config-data-scenarios.md`
- 修复动作模板：`06-remediation-templates.md`
- 用户问题索引：`07-question-bank.md`
- Kubernetes 进阶场景：`08-kubernetes-advanced-scenarios.md`
- 平台与可观测性深度场景：`09-platform-observability-deep-dive.md`
- 修复决策树：`10-remediation-decision-trees.md`

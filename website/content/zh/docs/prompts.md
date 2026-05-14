---
title: "Prompts"
weight: 9
description: "内置 MCP Prompt，用于故障排查、修复流程和引导式工具调用。"
---

# MCP Prompt

Cloud Native MCP Server 内置了一批 MCP Prompt，客户端可以先拉取 Prompt，再让 Agent 按 Prompt 指引去调用工具。

这些 Prompt 主要覆盖：

- 故障分诊
- 工作负载排查
- 服务连通性排查
- 发布恢复
- 安全修复流程
- 可观测性信号关联
- Argo CD 交付问题诊断
- LLM 应用问题分析

## Prompt 清单

| Prompt | 用途 |
|--------|------|
| `cloud_native_incident_triage` | 跨信号事故分诊 |
| `kubernetes_workload_diagnosis` | 针对单个 Pod / Workload 的排障 |
| `kubernetes_safe_remediation` | 引导 patch / restart / scale / delete 等修复流程 |
| `kubernetes_service_connectivity_diagnosis` | 排查 Service、Pod、EndpointSlice 与请求链路问题 |
| `kubernetes_rollout_recovery` | 排查 rollout 失败并选择恢复路径 |
| `cloud_native_observability_correlation` | 关联 alerts、metrics、logs、traces、Sentry、Langfuse |
| `argocd_delivery_diagnosis` | 排查 Argo CD / GitOps 交付问题 |
| `llm_app_observability_investigation` | 排查 LLM 应用故障与质量回退 |
| `prometheus_metrics_diagnosis` | 排查 Prometheus target、查询、告警与规则问题 |
| `loki_log_investigation` | 排查 Loki 日志与 LogQL 选择器问题 |
| `jaeger_trace_investigation` | 排查 Jaeger trace 与依赖关系问题 |
| `grafana_dashboard_diagnosis` | 排查 Grafana 仪表盘、面板、数据源与渲染问题 |
| `alertmanager_alert_triage` | 排查 Alertmanager 告警、分组、静默与接收器问题 |
| `helm_release_diagnosis` | 排查 Helm release 状态、values、manifest 与回滚选择 |
| `kibana_log_diagnosis` | 排查 Kibana 日志、仪表盘、data view、告警与 saved object |
| `elasticsearch_cluster_diagnosis` | 排查 Elasticsearch 集群健康、节点、索引与搜索问题 |
| `nacos_config_service_diagnosis` | 排查 Nacos 配置、服务发现与节点状态 |
| `sentry_issue_investigation` | 排查 Sentry issue 与 issue event |
| `langfuse_llm_trace_investigation` | 排查 Langfuse trace、prompt、评分、数据集与指标 |
| `opentelemetry_collector_diagnosis` | 排查 OTel Collector 配置、健康和 pipeline 状态 |
| `utilities_helper_usage` | 使用通用辅助 prompt 做时间、等待与轻量抓取 |
| `cloud_native_question_resolution` | 先理解用户问题，再把任务路由到正确 service 和工具 |
| `multi_service_root_cause_analysis` | 面向多服务、多信号的复合根因分析 |
| `release_regression_diagnosis` | 面向发布、同步、升级后的回归问题排查 |
| `telemetry_gap_diagnosis` | 面向 metrics、logs、traces 缺失的复合排查 |
| `end_to_end_request_path_diagnosis` | 面向用户请求链路故障的端到端排查 |

## 设计原则

这些 Prompt 默认遵循：

- 先读后写
- 先 summary，再 full detail
- 使用运行时真实存在的工具名
- 区分事实和推断
- 状态变更前要求显式确认

## 推荐客户端流程

1. 先从 MCP server 列出 prompts。
2. 选择与用户场景匹配的 prompt。
3. 带上 `namespace`、`kind`、`name`、`symptom`、`time_range` 等参数获取 prompt 内容。
4. 让 Agent 按 prompt 中给出的工作流和工具顺序执行。

## 可用性说明

Prompt 会根据启用的服务自动过滤：

- aggregate 端点会暴露完整 prompt 清单，但会按已启用服务做过滤
- service-specific 端点只会暴露和该 service 对应的 prompt
- 若 prompt 依赖的 service 未启用，它会被过滤掉，或在调用时被中间件拒绝

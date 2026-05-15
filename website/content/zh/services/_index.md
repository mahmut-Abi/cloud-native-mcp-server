---
title: "服务概览"
weight: 2
description: "Cloud Native MCP Server 已集成的云原生服务与能力入口。"
---

# 服务概览

Cloud Native MCP Server 集成了 15 个常用云原生服务，覆盖集群管理、监控、日志、告警、GitOps、配置中心、追踪、错误监控与 LLM 可观测性等运维场景。

## 服务列表

- [Kubernetes 服务]({{< relref "kubernetes.md" >}}) - 集群资源与工作负载管理
- [Helm 服务]({{< relref "helm.md" >}}) - Chart 与发布生命周期管理
- [Grafana 服务]({{< relref "grafana.md" >}}) - 仪表盘与告警可视化
- [Prometheus 服务]({{< relref "prometheus.md" >}}) - 指标采集与查询
- [Loki 服务]({{< relref "loki.md" >}}) - LogQL 查询与日志流检查
- [Kibana 服务]({{< relref "kibana.md" >}}) - 日志检索与分析
- [Elasticsearch 服务]({{< relref "elasticsearch.md" >}}) - 索引、搜索与集群操作
- [Argo CD 服务]({{< relref "argocd.md" >}}) - GitOps 应用、项目、集群与 manifest 检查
- [Alertmanager 服务]({{< relref "alertmanager.md" >}}) - 告警路由与静默管理
- [Jaeger 服务]({{< relref "jaeger.md" >}}) - 分布式链路追踪
- [Nacos 服务]({{< relref "nacos.md" >}}) - 命名空间、配置与服务发现检查
- [Langfuse 服务]({{< relref "langfuse.md" >}}) - LLM Trace、Prompt、评分、指标、项目、成员与 API Key 管理
- [Sentry 服务]({{< relref "sentry.md" >}}) - 错误监控、Issue 排障与事件检查
- [OpenTelemetry 服务]({{< relref "opentelemetry.md" >}}) - 遥测链路采集与排查
- [Utilities 服务]({{< relref "utilities.md" >}}) - 通用辅助工具集

## 建议阅读路径

1. 从 [快速开始]({{< relref "/getting-started/_index.md" >}}) 启动服务
2. 在 [配置指南]({{< relref "/docs/configuration.md" >}}) 配置服务地址与认证
3. 在 [工具参考]({{< relref "/docs/tools.md" >}}) 查看完整工具能力

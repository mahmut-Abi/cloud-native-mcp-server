---
title: "OpenTelemetry 服务"
weight: 9
---

# OpenTelemetry 服务

OpenTelemetry 服务提供采集器诊断、配置分析与遥测检查能力，包含 12 个工具来管理可观测性资源。

## 概述

Cloud Native MCP Server 中的 OpenTelemetry 服务使 AI 助手能够高效地从应用程序和基础设施收集和分析遥测数据。它提供用于指标收集、追踪分析和日志聚合的工具。

### 主要功能

{{< columns >}}
### 📊 指标收集
收集和分析应用程序和基础设施指标。
<--->

### 📍 追踪分析
分布式追踪与性能洞察。
{{< /columns >}}

{{< columns >}}
### 📝 日志聚合
统一日志收集和分析。
<--->

### 🛠️ 配置
管理 OpenTelemetry 收集器配置和状态。
{{< /columns >}}

---

## 可用工具 (12)

### 指标管理
- `opentelemetry_get_metrics`: 从采集器获取指标
- `opentelemetry_query_metrics`: 对采集器指标执行 PromQL 风格查询

### 追踪管理
- `opentelemetry_get_traces`: 从采集器获取追踪
- `opentelemetry_query_traces`: 带过滤条件搜索追踪

### 日志管理
- `opentelemetry_get_logs`: 从采集器获取日志
- `opentelemetry_query_logs`: 带过滤条件搜索日志

### 采集器诊断
- `opentelemetry_get_health`: 获取 OpenTelemetry 收集器健康状况
- `opentelemetry_get_status`: 获取详细收集器状态
- `opentelemetry_get_collector_summary`: 获取紧凑的采集器健康与配置概览
- `opentelemetry_get_config`: 获取完整采集器配置
- `opentelemetry_get_config_summary`: 获取配置摘要
- `opentelemetry_analyze_pipeline_status`: 分析 pipeline 中缺失组件与常见错误配置

---

## 快速示例

### 从收集器获取指标

```json
{
  "method": "tools/call",
  "params": {
    "name": "opentelemetry_get_metrics",
    "arguments": {
      "metric_name": "http_requests_total"
    }
  }
}
```

### 分析一条 traces pipeline

```json
{
  "method": "tools/call",
  "params": {
    "name": "opentelemetry_analyze_pipeline_status",
    "arguments": {
      "signal": "traces"
    }
  }
}
```

### 获取收集器配置摘要

```json
{
  "method": "tools/call",
  "params": {
    "name": "opentelemetry_get_config_summary",
    "arguments": {}
  }
}
```

---

## 最佳实践

- 实施适当的资源属性以实现有效过滤
- 为追踪使用适当的采样策略
- 配置适当的指标收集间隔
- 监控收集器健康状况和资源使用情况
- 基于遥测数据模式设置警报

## 下一步

- [Jaeger 服务](/zh/services/jaeger/) 了解详细追踪
- [可观测性指南](/zh/services/opentelemetry/) 了解详细设置
- [指标最佳实践](/zh/services/prometheus/) 了解收集策略

---
title: "OpenTelemetry 服务"
weight: 9
---

# OpenTelemetry 服务

OpenTelemetry 服务提供全面的指标、追踪和日志收集与分析功能，包含 9 个工具来管理可观测性资源。

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

## 可用工具 (9)

### 指标管理
- **otel-get-metrics**: 从 OpenTelemetry 收集器获取指标
- **otel-get-metric-data**: 获取指标数据
- **otel-list-metric-streams**: 列出指标流
- **otel-get-metrics-schema**: 获取指标模式

### 追踪管理
- **otel-get-traces**: 从 OpenTelemetry 收集器获取追踪
- **otel-search-traces**: 搜索追踪
- **otel-get-traces-schema**: 获取追踪模式

### 日志和配置管理
- **otel-get-logs**: 从 OpenTelemetry 收集器获取日志
- **otel-get-logs-schema**: 获取日志模式

### 系统信息
- **otel-get-status**: 获取 OpenTelemetry 收集器状态
- **otel-get-config**: 获取 OpenTelemetry 收集器配置
- **otel-get-health**: 获取 OpenTelemetry 收集器健康状况
- **otel-get-versions**: 获取 OpenTelemetry 组件版本

---

## 快速示例

### 从收集器获取指标

```json
{
  "method": "tools/call",
  "params": {
    "name": "otel-get-metrics",
    "arguments": {
      "metric_name": "http_requests_total",
      "start_time": "1 hour ago",
      "end_time": "now"
    }
  }
}
```

### 获取特定服务的追踪

```json
{
  "method": "tools/call",
  "params": {
    "name": "otel-get-traces",
    "arguments": {
      "service_name": "my-app",
      "limit": 50
    }
  }
}
```

### 获取收集器配置

```json
{
  "method": "tools/call",
  "params": {
    "name": "otel-get-config",
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
- [可观测性指南](/zh/guides/observability/) 了解详细设置
- [指标最佳实践](/zh/guides/metrics/) 了解收集策略
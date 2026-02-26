---
title: "Jaeger 服务"
weight: 8
---

# Jaeger 服务

Jaeger 服务提供全面的分布式追踪和性能分析功能，包含 8 个工具来管理追踪资源。

## 概述

Cloud Native MCP Server 中的 Jaeger 服务使 AI 助手能够高效地分析分布式追踪、服务依赖和性能指标。它提供用于追踪查询、依赖分析和性能监控的工具。

### 主要功能

{{< columns >}}
### 📍 追踪分析
跨微服务查询和分析分布式追踪。
<--->

### 🔗 依赖映射
可视化服务依赖关系和调用图。
{{< /columns >}}

{{< columns >}}
### ⚡ 性能监控
分析性能瓶颈和延迟模式。
<--->

### 📊 指标收集
收集和分析追踪指标和统计信息。
{{< /columns >}}

---

## 可用工具 (8)

### 追踪管理
- **jaeger-get-traces**: 按查询获取追踪
- **jaeger-get-trace**: 获取特定追踪
- **jaeger-search-traces**: 使用过滤器搜索追踪

### 服务和操作分析
- **jaeger-get-services**: 获取所有服务
- **jaeger-get-service-operations**: 获取服务的操作
- **jaeger-get-operations**: 获取所有操作

### 依赖关系和指标
- **jaeger-get-dependencies**: 获取服务依赖关系
- **jaeger-get-metrics**: 获取追踪指标

---

## 快速示例

### 获取特定服务的追踪

```json
{
  "method": "tools/call",
  "params": {
    "name": "jaeger-get-traces",
    "arguments": {
      "service": "my-app",
      "limit": 100
    }
  }
}
```

### 获取服务依赖关系

```json
{
  "method": "tools/call",
  "params": {
    "name": "jaeger-get-dependencies",
    "arguments": {
      "service": "my-app",
      "start": "1 hour ago"
    }
  }
}
```

### 获取所有服务

```json
{
  "method": "tools/call",
  "params": {
    "name": "jaeger-get-services",
    "arguments": {}
  }
}
```

---

## 最佳实践

- 实施适当的追踪头传播
- 使用适当的采样策略以提高性能
- 定期分析追踪以发现性能瓶颈
- 监控服务依赖关系以了解架构变化
- 基于追踪指标和异常设置警报

## 下一步

- [OpenTelemetry 服务](/zh/services/opentelemetry/) 了解指标和日志
- [追踪指南](/zh/services/jaeger/) 了解详细设置
- [性能分析](/zh/guides/performance/) 了解优化
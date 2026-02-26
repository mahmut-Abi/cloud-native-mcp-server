---
title: "Prometheus 服务"
weight: 4
---

# Prometheus 服务

Prometheus 服务提供全面的指标收集、查询和监控功能，包含 20 个工具来管理 Prometheus 资源。

## 概述

Cloud Native MCP Server 中的 Prometheus 服务使 AI 助手能够高效地查询和管理 Prometheus 指标。它提供用于指标查询、警报管理和配置管理的工具。

### 主要功能

{{< columns >}}
### 📈 指标查询
使用 PromQL 对 Prometheus 指标进行强大的查询功能。
<--->

### ⚠️ 警报管理
有效管理 Prometheus 警报和警报规则。
{{< /columns >}}

{{< columns >}}
### 🛠️ 配置
处理 Prometheus 配置和运行时信息。
<--->

### 📊 监控
从 Prometheus 访问详细的监控数据和统计信息。
{{< /columns >}}

---

## 可用工具 (20)

### 查询执行
- **prometheus-query**: 执行即时查询
- **prometheus-query-range**: 执行范围查询
- **prometheus-query-exemplars**: 查询示例数据

### 元数据查询
- **prometheus-label-names**: 获取标签名称
- **prometheus-label-values**: 获取标签值
- **prometheus-series**: 获取时间序列
- **prometheus-metadata**: 获取元数据

### 目标管理
- **prometheus-get-targets**: 获取目标列表
- **prometheus-get-target-metadata**: 获取目标元数据

### 规则和警报管理
- **prometheus-get-rules**: 获取规则列表
- **prometheus-get-alerts**: 获取警报列表
- **prometheus-get-alert-managers**: 获取 Alertmanager 实例

### 配置管理
- **prometheus-get-config**: 获取配置信息
- **prometheus-get-flags**: 获取启动参数

### 状态查询
- **prometheus-get-status**: 获取状态信息
- **prometheus-get-build-info**: 获取构建信息
- **prometheus-get-runtime-info**: 获取运行时信息

### TSDB 操作
- **prometheus-get-tsdb-status**: 获取 TSDB 状态
- **prometheus-get-tsdb-heatmap**: 获取 TSDB 热力图

---

## 快速示例

### 查询过去 5 分钟的指标

```json
{
  "method": "tools/call",
  "params": {
    "name": "prometheus-query-range",
    "arguments": {
      "query": "up",
      "start": "5 minutes ago",
      "end": "now",
      "step": "30s"
    }
  }
}
```

### 获取所有目标

```json
{
  "method": "tools/call",
  "params": {
    "name": "prometheus-get-targets",
    "arguments": {}
  }
}
```

### 查询特定指标

```json
{
  "method": "tools/call",
  "params": {
    "name": "prometheus-query",
    "arguments": {
      "query": "rate(http_requests_total[5m])"
    }
  }
}
```

---

## 最佳实践

- 使用适当的查询范围以避免性能问题
- 定期审查和优化 PromQL 查询
- 监控目标健康状况和可用性
- 为预聚合指标配置适当的记录规则
- 设置具有适当阈值的适当警报规则

## 下一步

- [Grafana 服务](/zh/services/grafana/) 了解可视化
- [监控指南](/zh/services/prometheus/) 了解详细设置
- [性能优化](/zh/guides/performance/) 了解查询优化
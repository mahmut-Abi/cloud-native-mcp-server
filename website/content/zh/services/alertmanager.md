---
title: "Alertmanager 服务"
weight: 7
---

# Alertmanager 服务

Alertmanager 服务提供全面的警报规则管理和通知功能，包含 15 个工具来管理警报资源。

## 概述

Cloud Native MCP Server 中的 Alertmanager 服务使 AI 助手能够高效地管理 Prometheus 警报、静默和通知路由。它提供用于警报管理、静默和配置的工具。

### 主要功能

{{< columns >}}
### ⚠️ 警报管理
对 Prometheus 警报和警报组进行完全控制。
<--->

### 🔕 静默管理
使用精确的时间和匹配规则管理警报静默。
{{< /columns >}}

{{< columns >}}
### 📮 通知路由
配置通知路由和接收器以传递警报。
<--->

### ⚙️ 配置
管理 Alertmanager 配置和状态信息。
{{< /columns >}}

---

## 可用工具 (15)

### 警报管理
- **alertmanager-get-alerts**: 获取所有警报
- **alertmanager-get-alert**: 获取特定警报
- **alertmanager-get-alert-groups**: 获取警报组
- **alertmanager-get-receivers**: 获取所有接收器

### 静默管理
- **alertmanager-get-silences**: 获取所有静默
- **alertmanager-create-silence**: 创建新静默
- **alertmanager-delete-silence**: 删除静默
- **alertmanager-get-alertmanagers**: 获取 Alertmanager 实例

### 配置和状态
- **alertmanager-get-config**: 获取 Alertmanager 配置
- **alertmanager-get-status**: 获取 Alertmanager 状态
- **alertmanager-get-metrics**: 获取 Alertmanager 指标
- **alertmanager-get-templates**: 获取通知模板
- **alertmanager-get-starttime**: 获取 Alertmanager 启动时间
- **alertmanager-get-version**: 获取 Alertmanager 版本
- **alertmanager-get-flags**: 获取 Alertmanager 标志值

---

## 快速示例

### 为高 CPU 警报创建静默

```json
{
  "method": "tools/call",
  "params": {
    "name": "alertmanager-create-silence",
    "arguments": {
      "matcher": {
        "name": "alertname",
        "value": "HighCPUUsage",
        "isRegex": false
      },
      "startsAt": "2023-10-01T10:00:00Z",
      "endsAt": "2023-10-01T12:00:00Z",
      "createdBy": "admin",
      "comment": "计划维护窗口"
    }
  }
}
```

### 获取所有活动警报

```json
{
  "method": "tools/call",
  "params": {
    "name": "alertmanager-get-alerts",
    "arguments": {}
  }
}
```

### 获取 Alertmanager 配置

```json
{
  "method": "tools/call",
  "params": {
    "name": "alertmanager-get-config",
    "arguments": {}
  }
}
```

---

## 最佳实践

- 使用精确匹配器进行有效静默
- 定期审查活动警报和静默
- 配置适当的通知路由和分组
- 监控 Alertmanager 性能和健康状况
- 为关键警报实施适当的升级策略

## 下一步

- [Prometheus 服务](/zh/services/prometheus/) 了解指标
- [警报指南](/zh/guides/alerting/) 了解详细设置
- [通知配置](/zh/guides/notifications/) 了解路由规则
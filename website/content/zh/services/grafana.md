---
title: "Grafana 服务"
weight: 3
---

# Grafana 服务

Grafana 服务提供全面的可视化、监控仪表板和告警功能，包含 36 个工具来创建和管理 Grafana 资源。

## 概述

Cloud Native MCP Server 中的 Grafana 服务使 AI 助手能够高效地管理 Grafana 仪表板、数据源、警报和其他监控资源。它提供用于仪表板创建、可视化管理和警报配置的工具。

### 主要功能

{{< columns >}}
### 📊 仪表板管理
对 Grafana 仪表板进行完全控制，包括创建、更新和共享。
<--->

### 🗂️ 数据源管理
使用工具管理 Grafana 数据源，包括配置和测试。
{{< /columns >}}

{{< columns >}}
### ⚠️ 警报管理
使用配置和监控工具处理 Grafana 警报和警报规则。
<--->

### 📈 可视化
有效创建和管理 Grafana 可视化和面板。
{{< /columns >}}

---

## 可用工具 (36)

### 仪表板管理
- **grafana-get-dashboards**: 获取所有仪表板
- **grafana-get-dashboard**: 获取特定仪表板
- **grafana-create-dashboard**: 创建新仪表板
- **grafana-update-dashboard**: 更新现有仪表板
- **grafana-delete-dashboard**: 删除仪表板
- **grafana-get-folders**: 获取所有文件夹
- **grafana-create-folder**: 创建新文件夹
- **grafana-update-folder**: 更新现有文件夹
- **grafana-delete-folder**: 删除文件夹

### 数据源管理
- **grafana-get-datasources**: 获取所有数据源
- **grafana-create-datasource**: 创建新数据源
- **grafana-update-datasource**: 更新现有数据源
- **grafana-delete-datasource**: 删除数据源
- **grafana-test-datasource**: 测试数据源连接

### 警报管理
- **grafana-get-alerts**: 获取所有警报
- **grafana-get-alert**: 获取特定警报
- **grafana-create-alert**: 创建新警报
- **grafana-update-alert**: 更新现有警报
- **grafana-delete-alert**: 删除警报
- **grafana-get-alert-rules**: 获取警报规则
- **grafana-get-alert-notifications**: 获取警报通知通道

### 用户和组织管理
- **grafana-get-users**: 获取所有用户
- **grafana-create-user**: 创建新用户
- **grafana-update-user**: 更新现有用户
- **grafana-delete-user**: 删除用户
- **grafana-get-orgs**: 获取所有组织
- **grafana-create-org**: 创建新组织
- **grafana-update-org**: 更新现有组织
- **grafana-delete-org**: 删除组织
- **grafana-get-teams**: 获取所有团队
- **grafana-create-team**: 创建新团队
- **grafana-update-team**: 更新现有团队
- **grafana-delete-team**: 删除团队

### 插件和配置管理
- **grafana-get-plugins**: 获取所有插件
- **grafana-install-plugin**: 安装插件
- **grafana-uninstall-plugin**: 卸载插件
- **grafana-get-annotations**: 获取注释
- **grafana-create-annotation**: 创建注释
- **grafana-get-snapshots**: 获取快照

---

## 快速示例

### 创建新仪表板

```json
{
  "method": "tools/call",
  "params": {
    "name": "grafana-create-dashboard",
    "arguments": {
      "dashboard": {
        "title": "我的应用程序仪表板",
        "panels": [
          {
            "id": 1,
            "title": "每秒请求数",
            "type": "graph",
            "targets": [
              {
                "expr": "rate(http_requests_total[5m])"
              }
            ]
          }
        ]
      }
    }
  }
}
```

### 获取所有仪表板

```json
{
  "method": "tools/call",
  "params": {
    "name": "grafana-get-dashboards",
    "arguments": {}
  }
}
```

### 添加数据源

```json
{
  "method": "tools/call",
  "params": {
    "name": "grafana-create-datasource",
    "arguments": {
      "name": "Prometheus",
      "type": "prometheus",
      "url": "http://prometheus:9090",
      "access": "proxy"
    }
  }
}
```

---

## 最佳实践

- 按应用程序或团队在文件夹中组织仪表板
- 为仪表板和面板使用一致的命名约定
- 配置适当的警报阈值和通知通道
- 根据不断变化的需求定期审查和更新仪表板
- 实施适当的用户权限和访问控制

## 下一步

- [Prometheus 服务](/zh/services/prometheus/) 了解指标收集
- [监控指南](/zh/guides/monitoring/) 了解详细设置
- [告警最佳实践](/zh/guides/alerting/) 了解有效告警
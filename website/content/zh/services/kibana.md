---
title: "Kibana 服务"
weight: 5
---

# Kibana 服务

Kibana 服务提供全面的日志分析、可视化和数据探索功能，包含 52 个工具来管理 Kibana 资源。

## 概述

Cloud Native MCP Server 中的 Kibana 服务使 AI 助手能够高效地管理 Kibana 仪表板、可视化和日志数据。它提供用于日志分析、数据可视化和 Elasticsearch 集成的工具。

### 主要功能

{{< columns >}}
### 🔍 日志分析
使用 Elasticsearch 集成的强大日志分析和搜索功能。
<--->

### 📊 数据可视化
创建和管理 Kibana 仪表板和可视化。
{{< /columns >}}

{{< columns >}}
### 🗂️ 对象管理
处理 Kibana 保存的对象，包括搜索、可视化和仪表板。
<--->

### 🌐 Elasticsearch 集成
与 Elasticsearch 无缝集成以进行数据索引和检索。
{{< /columns >}}

---

## 可用工具 (52)

### 索引管理
- **kibana-get-indices**: 获取所有索引
- **kibana-create-index**: 创建新索引
- **kibana-delete-index**: 删除索引
- **kibana-get-index-patterns**: 获取所有索引模式
- **kibana-create-index-pattern**: 创建新索引模式
- **kibana-delete-index-pattern**: 删除索引模式

### 可视化管理
- **kibana-get-visualizations**: 获取所有可视化
- **kibana-create-visualization**: 创建新可视化
- **kibana-update-visualization**: 更新可视化
- **kibana-delete-visualization**: 删除可视化

### 仪表板管理
- **kibana-get-dashboards**: 获取所有仪表板
- **kibana-create-dashboard**: 创建新仪表板
- **kibana-update-dashboard**: 更新仪表板
- **kibana-delete-dashboard**: 删除仪表板

### 搜索管理
- **kibana-get-searches**: 获取所有保存的搜索
- **kibana-create-search**: 创建新保存的搜索
- **kibana-update-search**: 更新保存的搜索
- **kibana-delete-search**: 删除保存的搜索

### 对象管理
- **kibana-get-objects**: 获取所有保存的对象
- **kibana-create-object**: 创建新保存的对象
- **kibana-update-object**: 更新保存的对象
- **kibana-delete-object**: 删除保存的对象

### 用户管理
- **kibana-get-users**: 获取所有用户
- **kibana-create-user**: 创建新用户
- **kibana-update-user**: 更新用户
- **kibana-delete-user**: 删除用户

### 角色管理
- **kibana-get-roles**: 获取所有角色
- **kibana-create-role**: 创建新角色
- **kibana-update-role**: 更新角色
- **kibana-delete-role**: 删除角色

### 空间管理
- **kibana-get-spaces**: 获取所有空间
- **kibana-create-space**: 创建新空间
- **kibana-update-space**: 更新空间
- **kibana-delete-space**: 删除空间

### API 密钥管理
- **kibana-get-api-keys**: 获取所有 API 密钥
- **kibana-create-api-key**: 创建新 API 密钥
- **kibana-delete-api-key**: 删除 API 密钥

### 系统信息
- **kibana-get-status**: 获取 Kibana 状态
- **kibana-get-info**: 获取 Kibana 信息
- **kibana-get-config**: 获取 Kibana 配置
- **kibana-get-plugins**: 获取已安装的插件
- **kibana-get-license**: 获取许可证信息
- **kibana-get-features**: 获取功能标志
- **kibana-get-telemetry**: 获取遥测设置
- **kibana-update-telemetry**: 更新遥测设置

### 导入/导出
- **kibana-get-saved-objects**: 按类型获取保存的对象
- **kibana-export-objects**: 导出保存的对象
- **kibana-import-objects**: 导入保存的对象
- **kibana-get-maps**: 获取地图保存的对象
- **kibana-get-cases**: 获取案例保存的对象

---

## 快速示例

### 创建新索引模式

```json
{
  "method": "tools/call",
  "params": {
    "name": "kibana-create-index-pattern",
    "arguments": {
      "id": "my-app-logs-*",
      "title": "my-app-logs-*",
      "timeFieldName": "@timestamp"
    }
  }
}
```

### 搜索文档

```json
{
  "method": "tools/call",
    "params": {
      "name": "kibana-get-searches",
      "arguments": {}
    }
  }
}
```

### 获取 Kibana 状态

```json
{
  "method": "tools/call",
  "params": {
    "name": "kibana-get-status",
    "arguments": {}
  }
}
```

---

## 最佳实践

- 使用适当的索引模式进行高效日志搜索
- 定期审查和优化可视化性能
- 使用角色和空间实施适当的访问控制
- 监控集群健康状况和资源使用情况
- 使用保存的搜索来访问频繁访问的数据模式

## 下一步

- [Elasticsearch 服务](/zh/services/elasticsearch/) 了解索引
- [日志分析指南](/zh/guides/logging/) 了解详细设置
- [可视化最佳实践](/zh/guides/visualization/) 了解有效仪表板
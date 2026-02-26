---
title: "Elasticsearch 服务"
weight: 6
---

# Elasticsearch 服务

Elasticsearch 服务提供全面的日志存储、搜索和数据索引功能，包含 14 个工具来管理 Elasticsearch 资源。

## 概述

Cloud Native MCP Server 中的 Elasticsearch 服务使 AI 助手能够高效地管理 Elasticsearch 索引、文档和集群。它提供用于索引、搜索和集群管理的工具。

### 主要功能

{{< columns >}}
### 🔍 高级搜索
强大的全文搜索和分析功能。
<--->

### 🗂️ 索引管理
对 Elasticsearch 索引和映射进行完全控制。
{{< /columns >}}

{{< columns >}}
### 📦 文档管理
处理文档索引、检索和批量操作。
<--->

### 🖥️ 集群管理
监控和管理 Elasticsearch 集群健康状况和性能。
{{< /columns >}}

---

## 可用工具 (14)

### 索引管理
- **elasticsearch-get-indices**: 获取所有索引
- **elasticsearch-create-index**: 创建新索引
- **elasticsearch-delete-index**: 删除索引
- **elasticsearch-get-index-settings**: 获取索引设置
- **elasticsearch-update-index-settings**: 更新索引设置
- **elasticsearch-get-mappings**: 获取索引映射
- **elasticsearch-update-mappings**: 更新索引映射

### 文档操作
- **elasticsearch-index-document**: 索引文档
- **elasticsearch-get-document**: 获取文档
- **elasticsearch-search**: 搜索文档
- **elasticsearch-delete-document**: 删除文档

### 集群管理
- **elasticsearch-get-cluster-info**: 获取集群信息
- **elasticsearch-get-cluster-health**: 获取集群健康状况
- **elasticsearch-get-nodes**: 获取集群节点

---

## 快速示例

### 索引文档

```json
{
  "method": "tools/call",
  "params": {
    "name": "elasticsearch-index-document",
    "arguments": {
      "index": "my-app-logs",
      "id": "1",
      "document": {
        "timestamp": "2023-10-01T12:00:00Z",
        "level": "info",
        "message": "应用程序启动成功"
      }
    }
  }
}
```

### 搜索文档

```json
{
  "method": "tools/call",
  "params": {
    "name": "elasticsearch-search",
    "arguments": {
      "index": "my-app-logs",
      "query": {
        "match": {
          "level": "error"
        }
      }
    }
  }
}
```

### 获取集群健康状况

```json
{
  "method": "tools/call",
  "params": {
    "name": "elasticsearch-get-cluster-health",
    "arguments": {}
  }
}
```

---

## 最佳实践

- 使用适当的索引映射以实现高效搜索
- 定期优化索引以提高性能
- 监控集群健康状况和资源使用情况
- 实施适当的索引生命周期管理
- 使用批量操作以实现高效数据摄取

## 下一步

- [Kibana 服务](/zh/services/kibana/) 了解可视化
- [日志管理指南](/zh/services/kibana/) 了解详细设置
- [搜索优化](/zh/services/elasticsearch/) 了解查询性能
---
title: "Elasticsearch 工具"
weight: 60
---

# Elasticsearch 工具

Elasticsearch 服务目前更偏向 **集群与索引排障**，而不是通用文档 CRUD 平台。

## 推荐起点

- `elasticsearch_cluster_health_summary`
- `elasticsearch_list_indices_paginated`
- `elasticsearch_nodes_summary`

## 常用操作

- `elasticsearch_health`
- `elasticsearch_list_indices`
- `elasticsearch_indices_summary`
- `elasticsearch_search_indices`
- `elasticsearch_index_stats`
- `elasticsearch_nodes`
- `elasticsearch_info`
- `elasticsearch_get_index_detail_advanced`
- `elasticsearch_get_cluster_detail_advanced`

## 示例

```json
{
  "name": "elasticsearch_search_indices",
  "arguments": {
    "query": "jaeger-*",
    "healthStatus": "green",
    "limit": 20
  }
}
```

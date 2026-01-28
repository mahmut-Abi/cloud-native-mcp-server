---
title: "Elasticsearch 工具"
weight: 60
---

# Elasticsearch 工具

Cloud Native MCP Server 提供 14 个 Elasticsearch 管理工具，用于索引管理、文档操作和集群管理。

## 索引管理

### list_indices

列出所有索引。

**参数**:
- `pattern` (string, optional) - 索引模式（如：*、logs-*）

**返回**: 索引列表

### get_index

获取单个索引的详细信息。

**参数**:
- `index` (string, required) - 索引名称

**返回**: 索引详细信息，包括映射、设置、别名等

### create_index

创建新的索引。

**参数**:
- `index` (string, required) - 索引名称
- `body` (object, optional) - 索引配置，包括：
  - `settings` - 索引设置
  - `mappings` - 字段映射
  - `aliases` - 别名配置

**返回**: 创建结果

### delete_index

删除指定的索引。

**参数**:
- `index` (string, required) - 索引名称或索引模式

**返回**: 删除结果

### get_index_stats

获取索引统计信息。

**参数**:
- `index` (string, required) - 索引名称

**返回**: 索引统计信息，包括：
  - 文档数量
  - 存储大小
  - 分片信息
  - 搜索性能指标

## 文档操作

### index_document

索引文档（创建或更新）。

**参数**:
- `index` (string, required) - 索引名称
- `body` (object, required) - 文档内容
- `id` (string, optional) - 文档 ID（不指定则自动生成）
- `refresh` (string, optional) - 刷新策略（true, false, wait_for）

**返回**: 索引结果，包括文档 ID 和版本

### get_document

获取单个文档。

**参数**:
- `index` (string, required) - 索引名称
- `id` (string, required) - 文档 ID

**返回**: 文档内容和元数据

### search_documents

搜索文档。

**参数**:
- `index` (string, required) - 索引名称
- `body` (object, required) - 查询 DSL，包括：
  - `query` - 查询条件
  - `aggs` - 聚合配置
  - `sort` - 排序配置
  - `size` - 返回数量
  - `from` - 起始位置

**返回**: 搜索结果，包括匹配的文档和聚合结果

### update_document

更新文档。

**参数**:
- `index` (string, required) - 索引名称
- `id` (string, required) - 文档 ID
- `body` (object, required) - 更新内容，格式为：
  ```json
  {
    "doc": { ... },
    "doc_as_upsert": true
  }
  ```

**返回**: 更新结果

### delete_document

删除文档。

**参数**:
- `index` (string, required) - 索引名称
- `id` (string, required) - 文档 ID

**返回**: 删除结果

## 集群管理

### get_cluster_health

获取集群健康状态。

**参数**:
- `index` (string, optional) - 限制到特定索引

**返回**: 集群健康信息，包括：
  - 集群状态（green, yellow, red）
  - 节点数量
  - 分片状态
  - 活动任务

### get_cluster_stats

获取集群统计信息。

**返回**: 集群统计数据，包括：
  - 节点信息
  - 索引统计
  - 文档数量
  - 存储使用
  - JVM 信息

### get_cluster_info

获取集群基本信息。

**返回**: 集群信息，包括：
  - 集群名称
  - 版本
  - UUID
  - 节点信息

## 别名管理

### get_aliases

获取别名信息。

**参数**:
- `index` (string, optional) - 索引名称
- `name` (string, optional) - 别名名称

**返回**: 别名配置

## 配置

Elasticsearch 工具通过以下配置进行初始化：

```yaml
elasticsearch:
  enabled: false

  # Elasticsearch 服务器地址（支持多节点）
  addresses:
    - "http://localhost:9200"

  # 单个地址（addresses 的替代方案）
  address: ""

  # 认证配置
  username: ""
  password: ""

  # Bearer Token
  bearerToken: ""

  # API Key（最高优先级）
  apiKey: ""

  # 请求超时（秒）
  timeoutSec: 30

  # TLS 配置
  tlsSkipVerify: false
  tlsCertFile: ""
  tlsKeyFile: ""
  tlsCAFile: ""
```

## 查询示例

### 创建索引

```json
{
  "name": "create_index",
  "arguments": {
    "index": "logs-2024",
    "body": {
      "settings": {
        "number_of_shards": 3,
        "number_of_replicas": 1
      },
      "mappings": {
        "properties": {
          "timestamp": {
            "type": "date"
          },
          "level": {
            "type": "keyword"
          },
          "message": {
            "type": "text"
          }
        }
      }
    }
  }
}
```

### 索引文档

```json
{
  "name": "index_document",
  "arguments": {
    "index": "logs-2024",
    "body": {
      "timestamp": "2024-01-01T00:00:00Z",
      "level": "info",
      "message": "Application started"
    }
  }
}
```

### 搜索文档

```json
{
  "name": "search_documents",
  "arguments": {
    "index": "logs-2024",
    "body": {
      "query": {
        "bool": {
          "must": [
            {
              "match": {
                "level": "error"
              }
            }
          ],
          "filter": [
            {
              "range": {
                "timestamp": {
                  "gte": "now-1h"
                }
              }
            }
          ]
        }
      },
      "sort": [
        {
          "timestamp": {
            "order": "desc"
          }
        }
      ],
      "size": 100
    }
  }
}
```

### 聚合查询

```json
{
  "name": "search_documents",
  "arguments": {
    "index": "logs-2024",
    "body": {
      "size": 0,
      "aggs": {
        "by_level": {
          "terms": {
            "field": "level"
          }
        },
        "over_time": {
          "date_histogram": {
            "field": "timestamp",
            "calendar_interval": "1h"
          }
        }
      }
    }
  }
}
```

## 最佳实践

1. **索引设计**:
   - 使用时间序列索引模式（如：logs-YYYY.MM.DD）
   - 合理设置分片和副本数量
   - 使用合适的字段映射和分词器

2. **查询优化**:
   - 使用过滤器查询（filter）提高性能
   - 限制返回的数据量（使用 `size`）
   - 使用聚合进行数据分析

3. **性能考虑**:
   - 定期检查 `get_cluster_health` 监控集群状态
   - 使用 `get_index_stats` 监控索引性能
   - 合理配置刷新策略

4. **数据管理**:
   - 使用索引生命周期管理（ILM）
   - 定期删除过期数据
   - 备份重要索引

5. **高可用性**:
   - 配置多个数据节点
   - 设置适当的副本数量
   - 使用别名进行索引切换

## 相关文档

- [配置指南](/zh/guides/configuration/)
- [部署指南](/zh/guides/deployment/)
- [Kibana 工具](/zh/tools/kibana/)
- [Elasticsearch 文档](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html)
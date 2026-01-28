---
title: "Kibana 工具"
weight: 50
---

# Kibana 工具

Cloud Native MCP Server 提供 52 个 Kibana 管理工具，用于索引管理、文档操作、可视化、仪表板和空间管理。

## 索引管理

### list_indices

列出所有索引。

**参数**:
- `pattern` (string, optional) - 索引模式

**返回**: 索引列表

### get_index

获取单个索引的详细信息。

**参数**:
- `index` (string, required) - 索引名称

**返回**: 索引详细信息

### create_index

创建新的索引。

**参数**:
- `index` (string, required) - 索引名称
- `body` (object, required) - 索引配置

**返回**: 创建结果

### delete_index

删除指定的索引。

**参数**:
- `index` (string, required) - 索引名称

**返回**: 删除结果

### get_index_stats

获取索引统计信息。

**参数**:
- `index` (string, required) - 索引名称

**返回**: 索引统计数据

### get_index_settings

获取索引设置。

**参数**:
- `index` (string, required) - 索引名称

**返回**: 索引设置

### update_index_settings

更新索引设置。

**参数**:
- `index` (string, required) - 索引名称
- `body` (object, required) - 设置配置

**返回**: 更新结果

## 文档操作

### search_documents

搜索文档。

**参数**:
- `index` (string, required) - 索引名称
- `query` (object, required) - 查询配置
- `size` (int, optional) - 返回数量
- `from` (int, optional) - 起始位置

**返回**: 搜索结果

### get_document

获取单个文档。

**参数**:
- `index` (string, required) - 索引名称
- `id` (string, required) - 文档 ID

**返回**: 文档内容

### create_document

创建新文档。

**参数**:
- `index` (string, required) - 索引名称
- `body` (object, required) - 文档内容
- `id` (string, optional) - 文档 ID

**返回**: 创建结果

### update_document

更新文档。

**参数**:
- `index` (string, required) - 索引名称
- `id` (string, required) - 文档 ID
- `body` (object, required) - 更新内容

**返回**: 更新结果

### delete_document

删除文档。

**参数**:
- `index` (string, required) - 索引名称
- `id` (string, required) - 文档 ID

**返回**: 删除结果

### bulk_operations

执行批量操作。

**参数**:
- `operations` (array, required) - 操作数组

**返回**: 批量操作结果

## 查询构建

### build_query

构建查询。

**参数**:
- `queryType` (string, required) - 查询类型
- `field` (string, required) - 字段名
- `value` (string, required) - 值
- `operator` (string, optional) - 操作符

**返回**: 查询配置

### execute_query

执行查询。

**参数**:
- `index` (string, required) - 索引名称
- `query` (object, required) - 查询配置

**返回**: 查询结果

### aggregate_data

聚合数据。

**参数**:
- `index` (string, required) - 索引名称
- `aggs` (object, required) - 聚合配置

**返回**: 聚合结果

### get_query_stats

获取查询统计信息。

**参数**:
- `index` (string, required) - 索引名称

**返回**: 查询统计数据

## 可视化

### list_visualizations

列出所有可视化。

**返回**: 可视化列表

### get_visualization

获取单个可视化的详细信息。

**参数**:
- `id` (string, required) - 可视化 ID

**返回**: 可视化详细信息

### create_visualization

创建新的可视化。

**参数**:
- `visualization` (object, required) - 可视化配置

**返回**: 创建的可视化信息

### update_visualization

更新现有的可视化。

**参数**:
- `id` (string, required) - 可视化 ID
- `visualization` (object, required) - 可视化配置

**返回**: 更新后的可视化信息

### delete_visualization

删除指定的可视化。

**参数**:
- `id` (string, required) - 可视化 ID

**返回**: 删除结果

## 仪表板

### list_dashboards

列出所有仪表板。

**参数**:
- `spaceId` (string, optional) - 空间 ID

**返回**: 仪表板列表

### get_dashboard

获取单个仪表板的详细信息。

**参数**:
- `id` (string, required) - 仪表板 ID

**返回**: 仪表板详细信息

### create_dashboard

创建新的仪表板。

**参数**:
- `dashboard` (object, required) - 仪表板配置
- `options` (object, optional) - 选项

**返回**: 创建的仪表板信息

### update_dashboard

更新现有的仪表板。

**参数**:
- `id` (string, required) - 仪表板 ID
- `dashboard` (object, required) - 仪表板配置

**返回**: 更新后的仪表板信息

### delete_dashboard

删除指定的仪表板。

**参数**:
- `id` (string, required) - 仪表板 ID

**返回**: 删除结果

## 索引模式

### list_index_patterns

列出所有索引模式。

**返回**: 索引模式列表

### get_index_pattern

获取单个索引模式的详细信息。

**参数**:
- `id` (string, required) - 索引模式 ID

**返回**: 索引模式详细信息

### create_index_pattern

创建新的索引模式。

**参数**:
- `pattern` (string, required) - 索引模式
- `timeField` (string, optional) - 时间字段

**返回**: 创建的索引模式信息

### update_index_pattern

更新现有的索引模式。

**参数**:
- `id` (string, required) - 索引模式 ID
- `pattern` (string, required) - 索引模式

**返回**: 更新后的索引模式信息

### delete_index_pattern

删除指定的索引模式。

**参数**:
- `id` (string, required) - 索引模式 ID

**返回**: 删除结果

## 保存查询

### list_saved_queries

列出所有保存的查询。

**返回**: 保存的查询列表

### get_saved_query

获取单个保存的查询。

**参数**:
- `id` (string, required) - 查询 ID

**返回**: 查询详细信息

### create_saved_query

创建新的保存查询。

**参数**:
- `query` (object, required) - 查询配置
- `title` (string, required) - 查询标题

**返回**: 创建的查询信息

### update_saved_query

更新现有的保存查询。

**参数**:
- `id` (string, required) - 查询 ID
- `query` (object, required) - 查询配置

**返回**: 更新后的查询信息

### delete_saved_query

删除指定的保存查询。

**参数**:
- `id` (string, required) - 查询 ID

**返回**: 删除结果

## 空间管理

### list_spaces

列出所有空间。

**返回**: 空间列表

### get_space

获取单个空间的详细信息。

**参数**:
- `id` (string, required) - 空间 ID

**返回**: 空间详细信息

### create_space

创建新的空间。

**参数**:
- `name` (string, required) - 空间名称
- `description` (string, optional) - 描述

**返回**: 创建的空间信息

### update_space

更新现有的空间。

**参数**:
- `id` (string, required) - 空间 ID
- `name` (string, required) - 空间名称
- `description` (string, optional) - 描述

**返回**: 更新后的空间信息

### delete_space

删除指定的空间。

**参数**:
- `id` (string, required) - 空间 ID

**返回**: 删除结果

## 发现

### discover_data

发现数据。

**参数**:
- `indexPattern` (string, required) - 索引模式
- `query` (object, optional) - 查询配置
- `sort` (array, optional) - 排序配置

**返回**: 发现结果

### get_field_capabilities

获取字段能力。

**参数**:
- `indexPattern` (string, required) - 索引模式
- `fields` (array, optional) - 字段列表

**返回**: 字段能力

## 导出导入

### export_objects

导出对象。

**参数**:
- `objects` (array, required) - 对象列表
- `type` (string, required) - 对象类型

**返回**: 导出结果

### import_objects

导入对象。

**参数**:
- `objects` (array, required) - 对象列表
- `createNewCopies` (bool, optional) - 创建新副本

**返回**: 导入结果

## 短链接

### create_short_url

创建短链接。

**参数**:
- `url` (string, required) - 目标 URL
- `name` (string, optional) - 链接名称

**返回**: 短链接信息

## 配置

Kibana 工具通过以下配置进行初始化：

```yaml
kibana:
  enabled: false

  # Kibana 服务器 URL
  url: "https://localhost:5601"

  # API Key
  apiKey: ""

  # Basic Auth
  username: ""
  password: ""

  # 请求超时（秒）
  timeoutSec: 30

  # TLS 配置
  skipVerify: false

  # 空间名称
  space: "default"
```

## 最佳实践

1. **索引管理**:
   - 使用索引模式管理时间序列数据
   - 定期使用 `get_index_stats` 监控索引状态
   - 使用适当的索引生命周期策略

2. **查询优化**:
   - 使用过滤器查询提高性能
   - 限制返回的数据量（使用 `size` 和 `from`）
   - 使用聚合进行数据分析

3. **仪表板组织**:
   - 使用空间组织不同的工作环境
   - 创建可视化并复用到多个仪表板
   - 保存常用查询供团队使用

4. **性能考虑**:
   - 避免在查询中使用通配符
   - 使用时间范围限制数据量
   - 合理设置分片和副本数量

## 相关文档

- [配置指南](/zh/guides/configuration/)
- [部署指南](/zh/guides/deployment/)
- [Elasticsearch 工具](/zh/tools/elasticsearch/)
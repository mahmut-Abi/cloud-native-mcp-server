---
title: "Prometheus 工具"
weight: 40
---

# Prometheus 工具

Cloud Native MCP Server 提供 20 个 Prometheus 管理工具，用于指标查询、规则管理、目标监控和配置管理。

## 查询执行

### query

执行即时查询（Instant Query）。

**参数**:
- `query` (string, required) - PromQL 查询表达式
- `time` (string, optional) - 查询时间戳（RFC3339 或 Unix 时间戳）

**返回**: 查询结果

### query_range

执行范围查询（Range Query）。

**参数**:
- `query` (string, required) - PromQL 查询表达式
- `start` (string, required) - 开始时间（RFC3339 或 Unix 时间戳）
- `end` (string, required) - 结束时间（RFC3339 或 Unix 时间戳）
- `step` (string, required) - 查询步长（如：15s, 1m, 5m）

**返回**: 查询结果（时间序列数据）

### query_exemplars

查询示例数据（Exemplars）。

**参数**:
- `query` (string, required) - PromQL 查询表达式
- `start` (string, required) - 开始时间
- `end` (string, required) - 结束时间

**返回**: 示例数据

## 元数据查询

### label_names

获取所有标签名。

**参数**:
- `match` (array, optional) - 标签选择器
- `start` (string, optional) - 开始时间
- `end` (string, optional) - 结束时间

**返回**: 标签名列表

### label_values

获取指定标签名的所有值。

**参数**:
- `label` (string, required) - 标签名
- `match` (array, optional) - 标签选择器
- `start` (string, optional) - 开始时间
- `end` (string, optional) - 结束时间

**返回**: 标签值列表

### series

获取匹配的时间序列。

**参数**:
- `match` (array, required) - 标签选择器
- `start` (string, required) - 开始时间
- `end` (string, required) - 结束时间

**返回**: 时间序列列表

### metadata

获取指标元数据。

**参数**:
- `metric` (string, optional) - 指标名称
- `limit` (int, optional) - 限制数量

**返回**: 元数据

## 目标管理

### targets

获取所有目标（Targets）列表。

**参数**:
- `state` (string, optional) - 状态过滤（active, dropped, any）

**返回**: 目标列表

### get_target_metadata

获取目标元数据。

**参数**:
- `matchTarget` (string, required) - 目标匹配
- `metric` (string, required) - 指标名称
- `limit` (int, optional) - 限制数量

**返回**: 目标元数据

## 规则管理

### rules

获取所有记录和告警规则。

**返回**: 规则列表

### get_alerts

获取所有活动的告警。

**参数**:
- `silenced` (bool, optional) - 是否包含沉默的告警
- `inhibited` (bool, optional) - 是否包含抑制的告警

**返回**: 告警列表

## 配置管理

### config

获取 Prometheus 配置。

**返回**: 配置信息

### flags

获取 Prometheus 启动参数。

**返回**: 启动参数列表

## 状态查询

### status

获取 Prometheus 状态信息。

**返回**: 状态信息

### query_stats

获取查询统计信息。

**返回**: 查询统计数据

## 快照管理

### snapshot

创建数据快照。

**参数**:
- `head` (bool, optional) - 创建 HEAD 快照

**返回**: 快照信息

## TSDB 操作

### tsdb_stats

获取 TSDB 统计信息。

**返回**: TSDB 统计数据

### tsdb_series

获取 TSDB 序列信息。

**返回**: TSDB 序列数据

## 存储操作

### block_info

获取存储块信息。

**返回**: 存储块信息

## 配置

Prometheus 工具通过以下配置进行初始化：

```yaml
prometheus:
  enabled: false

  # Prometheus 服务器地址
  address: "http://localhost:9090"

  # 请求超时（秒）
  timeoutSec: 30

  # Basic Auth
  username: ""
  password: ""

  # Bearer Token
  bearerToken: ""

  # TLS 配置
  tlsSkipVerify: false
  tlsCertFile: ""
  tlsKeyFile: ""
  tlsCAFile: ""
```

## 查询示例

### 即时查询

```json
{
  "name": "query",
  "arguments": {
    "query": "up{job=\"kubernetes-pods\"}"
  }
}
```

### 范围查询

```json
{
  "name": "query_range",
  "arguments": {
    "query": "rate(http_requests_total[5m])",
    "start": "2024-01-01T00:00:00Z",
    "end": "2024-01-01T01:00:00Z",
    "step": "1m"
  }
}
```

### 查询标签值

```json
{
  "name": "label_values",
  "arguments": {
    "label": "job"
  }
}
```

## 最佳实践

1. **查询优化**:
   - 使用 `query_range` 时设置适当的步长
   - 使用标签过滤减少返回数据量
   - 避免在查询中使用高基数标签

2. **性能考虑**:
   - 限制查询时间范围
   - 使用 `record` 规则预计算复杂查询
   - 监控查询性能指标

3. **告警管理**:
   - 定期检查 `get_alerts` 查看活动告警
   - 使用 `rules` 验证告警规则配置
   - 设置合理的告警阈值

4. **目标监控**:
   - 使用 `targets` 检查采集目标状态
   - 监控目标健康状态和采集延迟

## 相关文档

- [配置指南](/zh/guides/configuration/)
- [部署指南](/zh/guides/deployment/)
- [Prometheus 文档](https://prometheus.io/docs/)
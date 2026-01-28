---
title: "OpenTelemetry 工具"
weight: 90
---

# OpenTelemetry 工具

Cloud Native MCP Server 提供 9 个 OpenTelemetry 管理工具，用于指标、追踪和日志数据的收集和查询。

## 指标管理

### get_metrics

获取指标数据。

**参数**:
- `metricName` (string, optional) - 指标名称
- `labels` (object, optional) - 标签过滤
- `startTime` (string, optional) - 开始时间（RFC3339）
- `endTime` (string, optional) - 结束时间（RFC3339）

**返回**: 指标数据

**示例**:
```json
{
  "name": "get_metrics",
  "arguments": {
    "metricName": "http.requests",
    "labels": {
      "service": "api-gateway",
      "method": "GET"
    },
    "startTime": "2024-01-01T00:00:00Z",
    "endTime": "2024-01-01T01:00:00Z"
  }
}
```

### get_metric_data

获取详细的指标数据点。

**参数**:
- `metricName` (string, required) - 指标名称
- `labels` (object, optional) - 标签过滤
- `startTime` (string, required) - 开始时间（RFC3339）
- `endTime` (string, required) - 结束时间（RFC3339）
- `aggregation` (string, optional) - 聚合方式（sum, avg, min, max）

**返回**: 指标数据点列表

### list_metric_streams

列出所有指标流。

**参数**:
- `filter` (string, optional) - 过滤条件

**返回**: 指标流列表

## 追踪管理

### get_traces

获取追踪数据。

**参数**:
- `traceId` (string, optional) - 追踪 ID
- `serviceName` (string, optional) - 服务名称
- `operationName` (string, optional) - 操作名称
- `startTime` (string, optional) - 开始时间（RFC3339）
- `endTime` (string, optional) - 结束时间（RFC3339）
- `limit` (int, optional) - 返回数量限制

**返回**: 追踪数据

**示例**:
```json
{
  "name": "get_traces",
  "arguments": {
    "serviceName": "payment-service",
    "startTime": "2024-01-01T00:00:00Z",
    "endTime": "2024-01-01T01:00:00Z",
    "limit": 50
  }
}
```

### search_traces

搜索追踪。

**参数**:
- `serviceName` (string, required) - 服务名称
- `operationName` (string, optional) - 操作名称
- `tags` (object, optional) - 标签过滤
- `startTimeMin` (string, required) - 开始时间（RFC3339）
- `startTimeMax` (string, required) - 结束时间（RFC3339）
- `durationMin` (string, optional) - 最小持续时间
- `durationMax` (string, optional) - 最大持续时间
- `limit` (int, optional) - 返回数量限制

**返回**: 追踪列表

## 日志管理

### get_logs

获取日志数据。

**参数**:
- `serviceName` (string, optional) - 服务名称
- `severity` (string, optional) - 日志级别（DEBUG, INFO, WARN, ERROR）
- `startTime` (string, optional) - 开始时间（RFC3339）
- `endTime` (string, optional) - 结束时间（RFC3339）
- `limit` (int, optional) - 返回数量限制

**返回**: 日志列表

**示例**:
```json
{
  "name": "get_logs",
  "arguments": {
    "serviceName": "api-gateway",
    "severity": "ERROR",
    "startTime": "2024-01-01T00:00:00Z",
    "endTime": "2024-01-01T01:00:00Z",
    "limit": 100
  }
}
```

### search_logs

搜索日志。

**参数**:
- `query` (string, required) - 搜索查询
- `serviceName` (string, optional) - 服务名称
- `severity` (string, optional) - 日志级别
- `startTimeMin` (string, required) - 开始时间（RFC3339）
- `startTimeMax` (string, required) - 结束时间（RFC3339）
- `limit` (int, optional) - 返回数量限制

**返回**: 日志列表

## 配置管理

### get_config

获取 OpenTelemetry 配置。

**返回**: 配置信息，包括：
- 采样配置
- 导出器配置
- 资源属性

### get_status

获取 OpenTelemetry 状态。

**返回**: 状态信息，包括：
- 活跃的指标流
- 活跃的追踪
- 活跃的日志流
- 采样率

## 配置

OpenTelemetry 工具通过以下配置进行初始化：

```yaml
opentelemetry:
  enabled: false

  # OpenTelemetry Collector 地址
  address: "http://localhost:4318"

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

## 数据类型

### 指标（Metrics）

OpenTelemetry 支持三种指标类型：

1. **Counter**: 单调递增的计数器
   - 用途: 计数事件（请求数、错误数）
   - 示例: `http.requests.total`

2. **Gauge**: 可增可减的值
   - 用途: 当前状态（内存使用、连接数）
   - 示例: `memory.used.bytes`

3. **Histogram**: 分布统计
   - 用途: 性能分析（请求延迟）
   - 示例: `http.request.duration`

### 追踪（Traces）

追踪表示分布式系统中的请求流程：

```json
{
  "traceId": "abc123def456",
  "spanId": "789xyz",
  "parentSpanId": "123abc",
  "name": "GET /api/users",
  "kind": "SPAN_KIND_SERVER",
  "startTimeUnixNano": 1704067200000000000,
  "endTimeUnixNano": 1704067200500000000,
  "attributes": {
    "http.method": "GET",
    "http.route": "/api/users",
    "http.status_code": 200,
    "service.name": "api-gateway"
  },
  "status": {
    "code": "STATUS_CODE_OK"
  }
}
```

### 日志（Logs）

日志记录应用程序事件：

```json
{
  "timeUnixNano": 1704067200000000000,
  "severityNumber": 9,
  "severityText": "INFO",
  "body": {
    "stringValue": "User login successful"
  },
  "attributes": {
    "user.id": "12345",
    "service.name": "auth-service",
    "service.version": "1.0.0"
  },
  "traceId": "abc123def456",
  "spanId": "789xyz"
  }
}
```

## 使用场景

### 性能监控

```json
{
  "name": "get_metric_data",
  "arguments": {
    "metricName": "http.request.duration",
    "labels": {
      "service": "api-gateway"
    },
    "startTime": "2024-01-01T00:00:00Z",
    "endTime": "2024-01-01T01:00:00Z",
    "aggregation": "avg"
  }
}
```

获取 API 网关的平均请求延迟。

### 错误追踪

```json
{
  "name": "search_logs",
  "arguments": {
    "query": "error",
    "severity": "ERROR",
    "startTimeMin": "2024-01-01T00:00:00Z",
    "startTimeMax": "2024-01-01T01:00:00Z",
    "limit": 50
  }
}
```

搜索最近的错误日志。

### 请求分析

```json
{
  "name": "search_traces",
  "arguments": {
    "serviceName": "payment-service",
    "durationMin": "1s",
    "startTimeMin": "2024-01-01T00:00:00Z",
    "startTimeMax": "2024-01-01T01:00:00Z",
    "limit": 10
  }
}
```

查找支付服务的慢请求。

### 流量监控

```json
{
  "name": "get_metrics",
  "arguments": {
    "metricName": "http.requests.total",
    "startTime": "2024-01-01T00:00:00Z",
    "endTime": "2024-01-01T01:00:00Z"
  }
}
```

获取一小时的请求总数。

## 最佳实践

1. **指标设计**:
   - 使用语义化的指标名称
   - 添加有用的标签（service, version, environment）
   - 选择合适的指标类型（Counter, Gauge, Histogram）

2. **追踪策略**:
   - 为所有外部调用创建 Span
   - 添加有意义的属性
   - 保持合理的采样率

3. **日志规范**:
   - 使用结构化日志
   - 添加相关上下文信息
   - 使用适当的日志级别

4. **数据采样**:
   - 高流量场景使用低采样率
   - 错误追踪始终采样
   - 关键业务流程使用高采样率

5. **标签管理**:
   - 限制标签数量和基数
   - 使用一致的标签命名
   - 避免在标签中存储高基数数据

## 集成建议

### 与 Jaeger 集成

使用 OpenTelemetry 收集追踪数据，发送到 Jaeger 进行可视化。

### 与 Prometheus 集成

使用 OpenTelemetry 收集指标数据，发送到 Prometheus 进行查询和告警。

### 与 Elasticsearch 集成

使用 OpenTelemetry 收集日志数据，发送到 Elasticsearch/Kibana 进行搜索和分析。

## 相关文档

- [配置指南](/zh/guides/configuration/)
- [部署指南](/zh/guides/deployment/)
- [Jaeger 工具](/zh/tools/jaeger/)
- [Prometheus 工具](/zh/tools/prometheus/)
- [OpenTelemetry 文档](https://opentelemetry.io/docs/)
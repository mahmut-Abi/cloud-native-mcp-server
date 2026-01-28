---
title: "Jaeger 工具"
weight: 80
---

# Jaeger 工具

Cloud Native MCP Server 提供 8 个 Jaeger 管理工具，用于分布式追踪、依赖分析和性能监控。

## 追踪查询

### get_trace

获取单个追踪的详细信息。

**参数**:
- `traceId` (string, required) - 追踪 ID

**返回**: 追踪详细信息，包括：
- 追踪 ID
- 开始时间
- 持续时间
- Span 列表
- 服务调用关系

### search_traces

搜索追踪。

**参数**:
- `service` (string, optional) - 服务名称
- `operation` (string, optional) - 操作名称
- `tags` (object, optional) - 标签过滤
- `startTimeMin` (string, optional) - 开始时间（RFC3339）
- `startTimeMax` (string, optional) - 结束时间（RFC3339）
- `durationMin` (string, optional) - 最小持续时间
- `durationMax` (string, optional) - 最大持续时间
- `limit` (int, optional) - 返回数量限制

**返回**: 追踪列表

**示例**:
```json
{
  "name": "search_traces",
  "arguments": {
    "service": "api-gateway",
    "operation": "GET /api/users",
    "startTimeMin": "2024-01-01T00:00:00Z",
    "startTimeMax": "2024-01-01T01:00:00Z",
    "limit": 100
  }
}
```

### get_services

获取所有服务列表。

**返回**: 服务名称列表

### get_operations

获取指定服务的所有操作列表。

**参数**:
- `service` (string, required) - 服务名称

**返回**: 操作名称列表

## 依赖分析

### get_dependencies

获取服务之间的依赖关系。

**参数**:
- `startTimeMin` (string, optional) - 开始时间（RFC3339）
- `endTimeMax` (string, optional) - 结束时间（RFC3339）

**返回**: 依赖关系图，包括：
- 父服务
- 子服务
- 调用次数

## 指标查询

### get_metrics

获取追踪指标。

**参数**:
- `service` (string, optional) - 服务名称
- `operation` (string, optional) - 操作名称
- `spanKind` (string, optional) - Span 类型（client, server, producer, consumer）

**返回**: 指标数据，包括：
- 追踪数量
- 错误率
- 平均延迟
- P50/P95/P99 延迟

## 配置查询

### get_config

获取 Jaeger 配置。

**返回**: 配置信息

### get_status

获取 Jaeger 状态。

**返回**: 状态信息，包括：
- 存储状态
- 采样率
- 启用功能

## 配置

Jaeger 工具通过以下配置进行初始化：

```yaml
jaeger:
  enabled: false

  # Jaeger 服务器地址
  # 通常为 Jaeger Query API
  address: "http://localhost:16686"

  # 请求超时（秒）
  timeoutSec: 30
```

## 追踪数据结构

### Span

Span 是追踪的基本单元，表示单个操作：

```json
{
  "traceID": "abc123def456",
  "spanID": "789xyz",
  "operationName": "GET /api/users",
  "startTime": "2024-01-01T00:00:00Z",
  "duration": 50000000,
  "tags": [
    {
      "key": "http.method",
      "value": "GET"
    },
    {
      "key": "http.status_code",
      "value": "200"
    }
  ],
  "process": {
    "serviceName": "api-gateway",
    "tags": [
      {
        "key": "hostname",
        "value": "server-1"
      }
    ]
  },
  "references": [
    {
      "refType": "CHILD_OF",
      "traceID": "abc123def456",
      "spanID": "123abc"
    }
  ]
}
```

### Trace

Trace 是多个 Span 的集合，表示一个完整的请求流程：

```json
{
  "traceID": "abc123def456",
  "spans": [
    {
      "spanID": "123abc",
      "operationName": "GET /api/users",
      "startTime": "2024-01-01T00:00:00Z",
      "duration": 100000000
    },
    {
      "spanID": "456def",
      "operationName": "database.query",
      "startTime": "2024-01-01T00:00:00.010Z",
      "duration": 50000000,
      "references": [
        {
          "refType": "CHILD_OF",
          "traceID": "abc123def456",
          "spanID": "123abc"
        }
      ]
    }
  ]
}
```

## 使用场景

### 性能分析

```json
{
  "name": "search_traces",
  "arguments": {
    "service": "api-gateway",
    "durationMin": "1s",
    "limit": 10
  }
}
```

查找超过 1 秒的慢请求。

### 错误追踪

```json
{
  "name": "search_traces",
  "arguments": {
    "service": "payment-service",
    "tags": {
      "error": "true"
    },
    "startTimeMin": "2024-01-01T00:00:00Z",
    "limit": 50
  }
}
```

查找支付服务的错误追踪。

### 依赖分析

```json
{
  "name": "get_dependencies",
  "arguments": {
    "startTimeMin": "2024-01-01T00:00:00Z",
    "endTimeMax": "2024-01-01T23:59:59Z"
  }
}
```

分析一天内的服务依赖关系。

### 服务监控

```json
{
  "name": "get_metrics",
  "arguments": {
    "service": "api-gateway"
  }
}
```

获取 API 网关的性能指标。

## 最佳实践

1. **追踪设计**:
   - 为每个重要操作创建 Span
   - 添加有意义的标签（service, operation, error）
   - 记录关键业务指标

2. **性能优化**:
   - 使用适当的采样率
   - 避免在 Span 中存储大量数据
   - 定期清理旧的追踪数据

3. **错误分析**:
   - 在错误发生时记录详细信息
   - 使用标签标识错误类型
   - 关联相关追踪进行根因分析

4. **依赖管理**:
   - 定期使用 `get_dependencies` 检查服务依赖
   - 识别循环依赖和不必要的依赖
   - 优化服务间调用

5. **监控和告警**:
   - 监控追踪成功率
   - 设置 P95/P99 延迟告警
   - 追踪错误率变化

## Span 类型

Jaeger 支持以下 Span 类型：

- **server**: 服务端接收请求
- **client**: 客户端发送请求
- **producer**: 消息生产者
- **consumer**: 消息消费者

## 采样策略

合理的采样策略对于性能和数据量控制：

- **始终采样**: 关键业务流程
- **概率采样**: 一般业务流程（如 1%）
- **动态采样**: 根据负载调整

## 相关文档

- [配置指南](/zh/guides/configuration/)
- [部署指南](/zh/guides/deployment/)
- [OpenTelemetry 工具](/zh/tools/opentelemetry/)
- [Jaeger 文档](https://www.jaegertracing.io/docs/)
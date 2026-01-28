---
title: "数据流"
weight: 30
---

# 数据流

本文档描述 Cloud Native MCP Server 中的请求和响应数据流。

## 请求流

### 完整请求流程

```
客户端
  │
  ├─> HTTP/SSE 连接
  │
  ├─> 认证中间件
  │   ├─> 验证 API Key/Token
  │   └─> 检查权限
  │
  ├─> 速率限制中间件
  │   └─> 检查配额
  │
  ├─> 路由层
  │   ├─> 解析服务和方法
  │   └─> 验证参数
  │
  ├─> 审计中间件
  │   └─> 记录请求开始
  │
  ├─> 服务管理器
  │   └─> 路由到服务
  │
  ├─> 缓存层
  │   ├─> 检查缓存
  │   └─> 返回缓存或继续
  │
  ├─> 服务
  │   ├─> 调用外部 API
  │   ├─> 处理响应
  │   └─> 更新缓存
  │
  ├─> 审计中间件
  │   └─> 记录请求完成
  │
  ├─> 指标中间件
  │   └─> 记录指标
  │
  └─> 响应返回客户端
```

### 详细步骤

#### 1. 客户端连接

客户端通过以下方式之一连接到服务器：

- **SSE**: `GET /api/kubernetes/sse`
- **HTTP**: `POST /api/kubernetes/http`
- **Streamable-HTTP**: `POST /api/kubernetes/streamable-http`
- **stdio**: 标准输入/输出

**请求格式** (JSON-RPC 2.0):

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "list_pods",
    "arguments": {
      "namespace": "default"
    }
  }
}
```

#### 2. 认证中间件

验证客户端身份：

```go
// 检查认证头
apiKey := r.Header.Get("X-API-Key")
if apiKey == "" {
    return errors.New("missing API key")
}

// 验证 API Key
if apiKey != config.Auth.APIKey {
    return errors.New("invalid API key")
}
```

**认证方式**:
- API Key: `X-API-Key: your-key`
- Bearer Token: `Authorization: Bearer your-token`
- Basic Auth: `Authorization: Basic base64(user:pass)`

#### 3. 速率限制中间件

检查客户端请求频率：

```go
// Token Bucket 算法
if rateLimiter.Allow() {
    // 继续处理请求
} else {
    return errors.New("rate limit exceeded")
}
```

**配置**:

```yaml
ratelimit:
  enabled: true
  requests_per_second: 100
  burst: 200
```

#### 4. 路由层

解析请求并路由到正确的服务：

```go
// 解析工具名称
toolName := params["name"].(string)

// 路由到服务
serviceName := getToolService(toolName)
service := manager.GetService(serviceName)

// 调用工具
result, err := service.CallTool(ctx, toolName, params["arguments"].(map[string]interface{}))
```

#### 5. 审计中间件

记录请求信息：

```go
auditLog := &AuditLog{
    Timestamp: time.Now(),
    RequestID: getRequestID(r),
    Tool:      toolName,
    Params:    params["arguments"],
    Status:    "started",
}

auditLogger.Log(auditLog)
```

#### 6. 服务管理器

管理服务生命周期：

```go
// 获取服务
service := registry.GetService(serviceName)

// 检查服务健康
if err := service.HealthCheck(); err != nil {
    return fmt.Errorf("service unavailable: %w", err)
}

// 调用工具
result, err := service.CallTool(ctx, toolName, arguments)
```

#### 7. 缓存层

检查和更新缓存：

```go
// 生成缓存键
cacheKey := fmt.Sprintf("%s:%s:%v", serviceName, toolName, arguments)

// 检查缓存
if cached, found := cache.Get(cacheKey); found {
    return cached, nil
}

// 调用服务
result, err := service.CallTool(ctx, toolName, arguments)

// 更新缓存
if err == nil {
    cache.Set(cacheKey, result, ttl)
}
```

#### 8. 服务执行

调用外部服务：

```go
// Kubernetes 示例
pods, err := k8sClient.CoreV1().Pods(namespace).List(ctx, options)

// Grafana 示例
dashboards, err := grafanaClient.ListDashboards(ctx)

// Prometheus 示例
result, err := prometheusClient.Query(ctx, query)
```

#### 9. 审计记录完成

记录请求完成信息：

```go
auditLog.Status = "completed"
auditLog.Duration = time.Since(startTime)
auditLog.Error = err

auditLogger.Log(auditLog)
```

#### 10. 指标记录

记录性能指标：

```go
metrics.RecordRequest(toolName, "success", duration)
metrics.RecordCacheHit(serviceName, cacheHit)
```

#### 11. 响应返回

返回结果给客户端：

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "pods": [
      {
        "name": "pod-1",
        "status": "Running"
      }
    ]
  }
}
```

## 响应流

### 响应处理流程

```
服务
  │
  ├─> 处理结果
  │
  ├─> 数据转换
  │   ├─> 格式化
  │   └─> 压缩
  │
  ├─> 缓存更新
  │   └─> 存储到缓存
  │
  ├─> 指标更新
  │   └─> 记录性能指标
  │
  └─> 返回响应
```

### 响应优化

#### 数据转换

```go
// 格式化响应
formatted := formatResponse(result)

// 压缩响应
if config.Performance.CompressionEnabled {
    compressed := compressResponse(formatted)
    return compressed
}
```

#### 缓存更新

```go
// 存储到缓存
cache.Set(cacheKey, result, ttl)

// 记录缓存统计
metrics.RecordCacheStore(serviceName)
```

#### 指标更新

```go
// 记录请求指标
metrics.RecordRequest(toolName, status, duration)

// 记录延迟
metrics.RecordLatency(toolName, duration)
```

## SSE 流式响应

### SSE 连接流程

```
客户端
  │
  ├─> GET /api/kubernetes/sse
  │
  ├─> 认证中间件
  │
  ├─> 建立持久连接
  │
  ├─> 发送事件
  │   ├─> Event: data
  │   ├─> Event: update
  │   └─> Event: complete
  │
  └─> 保持连接
```

### SSE 事件格式

```
event: data
data: {"type":"pod","name":"pod-1","status":"Running"}

event: data
data: {"type":"pod","name":"pod-2","status":"Pending"}

event: complete
data: {"message":"All data sent"}
```

### SSE 配置

```yaml
server:
  mode: "sse"
  writeTimeoutSec: 0  # 保持连接
  idleTimeoutSec: 60
```

## Streamable-HTTP 响应

### 流式 HTTP 响应流程

```
客户端
  │
  ├─> POST /api/kubernetes/streamable-http
  │
  ├─> 认证中间件
  │
  ├─> 处理请求
  │
  ├─> 流式返回数据
  │   ├─> Chunk 1
  │   ├─> Chunk 2
  │   └─> Chunk N
  │
  └─> 完成
```

### 流式响应格式

```
HTTP/1.1 200 OK
Content-Type: application/json
Transfer-Encoding: chunked

{"data": "chunk-1"}
{"data": "chunk-2"}
{"data": "chunk-3"}
```

## 错误处理

### 错误响应格式

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32602,
    "message": "Invalid params",
    "data": {
      "details": "namespace is required"
    }
  }
}
```

### 错误类型

- **InvalidParams**: 参数无效
- **NotFound**: 资源不存在
- **PermissionDenied**: 权限不足
- **Timeout**: 请求超时
- **InternalError**: 内部错误
- **ServiceUnavailable**: 服务不可用

### 错误处理流程

```
服务错误
  │
  ├─> 错误转换
  │   └─> 转换为标准错误格式
  │
  ├─> 审计记录
  │   └─> 记录错误信息
  │
  ├─> 指标记录
  │   └─> 记录错误指标
  │
  └─> 返回错误响应
```

## 性能优化

### 缓存策略

```go
// 缓存键生成
cacheKey := generateCacheKey(serviceName, toolName, arguments)

// LRU 缓存
if cached, found := cache.Get(cacheKey); found {
    return cached, nil
}

// 缓存更新
cache.Set(cacheKey, result, ttl)
```

### 连接池

```go
// HTTP 客户端连接池
client := &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}
```

### 批处理

```go
// 批量获取资源
pods, err := k8sClient.CoreV1().Pods(namespace).List(ctx, options)

// 批量处理
for _, pod := range pods.Items {
    processPod(pod)
}
```

## 相关文档

- [架构概览](/zh/concepts/architecture/overview/)
- [核心组件](/zh/concepts/architecture/components/)
- [功能特性](/zh/concepts/features/)
- [性能指南](/zh/guides/performance/)
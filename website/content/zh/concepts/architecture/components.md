---
title: "核心组件"
weight: 20
---

# 核心组件

本文档详细描述 Cloud Native MCP Server 的各个核心组件。

## 1. HTTP 服务器

**职责**: 处理传入的 HTTP/SSE 请求和连接

**特性**:
- 支持多种运行模式 (SSE, HTTP, stdio, Streamable-HTTP)
- 可配置的超时和连接限制
- 优雅关闭
- 健康检查端点

**关键文件**:
- `cmd/server/server.go`
- `internal/middleware/`

**配置**:

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8080"
  readTimeoutSec: 30
  writeTimeoutSec: 0
  idleTimeoutSec: 60
```

## 2. 路由层

**职责**: 将请求路由到正确的服务和工具

**特性**:
- 动态路由注册
- 路径参数解析
- 查询参数验证
- 错误处理

**关键文件**:
- `internal/services/registry.go`

**路由示例**:

```go
// 服务路由
routes := map[string]Service{
    "kubernetes": k8sService,
    "grafana": grafanaService,
    "prometheus": prometheusService,
}

// 工具路由
toolRoute := fmt.Sprintf("/api/%s/http", serviceName)
```

## 3. 中间件层

**职责**: 在请求处理之前和之后执行通用逻辑

### 认证中间件

验证客户端身份：

```yaml
auth:
  enabled: true
  mode: "apikey"
  apiKey: "${MCP_AUTH_API_KEY}"
```

**支持的模式**:
- `apikey`: API Key 认证
- `bearer`: Bearer Token 认证
- `basic`: HTTP Basic Auth

**关键文件**:
- `internal/middleware/auth_middleware.go`
- `internal/middleware/auth.go`

### 审计中间件

记录所有操作：

```yaml
audit:
  enabled: true
  storage: "memory"
  format: "json"
```

**记录内容**:
- 请求时间
- 用户身份
- 工具名称
- 请求参数
- 响应状态
- 执行时间

**关键文件**:
- `internal/middleware/audit_middleware.go`
- `internal/middleware/audit_log.go`

### 速率限制中间件

防止滥用：

```yaml
ratelimit:
  enabled: true
  requests_per_second: 100
  burst: 200
```

**实现**: Token Bucket 算法

**关键文件**:
- `internal/middleware/ratelimit.go`

### 安全中间件

输入清理和安全检查：

```yaml
sanitization:
  enabled: true
  max_length: 1000
```

**清理内容**:
- SQL 注入
- XSS 攻击
- 命令注入
- 危险字符

**关键文件**:
- `internal/middleware/security_middleware.go`
- `internal/util/sanitize/`

### 指标中间件

收集性能指标：

```go
// Prometheus 指标
mcp_requests_total{method="kubernetes_list_pods",status="success"} 1234
mcp_request_duration_seconds{method="kubernetes_list_pods"} 0.123
```

**关键文件**:
- `internal/middleware/metrics_middleware.go`
- `internal/observability/metrics/`

## 4. 服务管理器

**职责**: 管理所有注册的服务和工具

**特性**:
- 服务注册和发现
- 工具调用路由
- 服务生命周期管理
- 健康检查协调

**关键文件**:
- `internal/services/manager/manager.go`
- `internal/services/registry.go`

**服务注册流程**:

```go
// 1. 创建服务实例
service := kubernetes.NewService()

// 2. 初始化服务
service.Initialize(config.Kubernetes)

// 3. 注册到管理器
manager.RegisterService(service)

// 4. 注册工具
for _, tool := range service.GetTools() {
    manager.RegisterTool(tool)
}
```

## 5. 缓存层

**职责**: 提供高性能缓存以减少外部服务调用

### LRU 缓存

最近最少使用缓存：

```go
cache := cache.NewLRUCache(1000, 300*time.Second)
```

**适用场景**:
- 读取密集型操作
- 数据变化不频繁
- 高延迟操作

### 分段缓存

提供更好的并发性能：

```go
cache := cache.NewSegmentedCache(1000, 10, 300*time.Second)
```

**适用场景**:
- 高并发场景
- 需要低延迟
- 多核 CPU

**配置**:

```yaml
cache:
  enabled: true
  type: "lru"
  max_size: 1000
  default_ttl: 300
```

**关键文件**:
- `internal/services/cache/`

## 6. 密钥管理器

**职责**: 安全地存储和管理敏感凭据

**特性**:
- 内存存储
- 密钥轮换
- 密钥生成
- 过期管理

**支持的密钥类型**:
- API 密钥
- Bearer token
- Basic auth 凭据

**关键文件**:
- `internal/secrets/manager.go`

**使用示例**:

```go
// 创建密钥管理器
manager := secrets.NewInMemoryManager()

// 存储密钥
secret := &secrets.Secret{
    Type:  secrets.SecretTypeAPIKey,
    Name:  "my-api-key",
    Value: "Abc123!@#Xyz789!@#",
}
manager.Store(secret)

// 检索密钥
retrieved, err := manager.Retrieve(secret.ID)

// 轮换密钥
rotated, err := manager.Rotate(secret.ID)
```

## 7. 日志系统

**职责**: 结构化日志记录

**特性**:
- 多级别日志 (debug, info, warn, error)
- JSON 和文本格式
- 结构化字段
- 上下文支持

**配置**:

```yaml
logging:
  level: "info"
  json: false
```

**日志格式**:

```json
{
  "level": "info",
  "timestamp": "2024-01-01T00:00:00Z",
  "service": "kubernetes",
  "tool": "list_pods",
  "duration_ms": 123,
  "status": "success"
}
```

**关键文件**:
- `internal/logging/logging.go`

## 8. 指标系统

**职责**: 收集和暴露性能指标

**特性**:
- Prometheus 格式
- 请求计数
- 延迟统计
- 缓存命中率

**关键指标**:

```
# 请求指标
mcp_requests_total{method="kubernetes_list_pods",status="success"} 1234
mcp_request_duration_seconds{method="kubernetes_list_pods"} 0.123

# 缓存指标
mcp_cache_hits_total{service="kubernetes"} 456
mcp_cache_misses_total{service="kubernetes"} 78

# 连接指标
mcp_active_connections 10
mcp_total_connections 100
```

**关键文件**:
- `internal/observability/metrics/`

## 9. OpenTelemetry 集成

**职责**: 分布式追踪和遥测

**特性**:
- 指标收集
- 追踪支持
- 日志关联

**配置**:

```yaml
opentelemetry:
  enabled: false
  address: "http://localhost:4318"
```

**使用示例**:

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

tracer := otel.Tracer("cloud-native-mcp-server")

ctx, span := tracer.Start(ctx, "list_pods")
defer span.End()

// 执行操作
pods, err := k8sClient.ListPods(ctx, namespace)
```

**关键文件**:
- `internal/observability/otel/`

## 组件交互

```
客户端请求
    │
    ▼
HTTP Server
    │
    ▼
认证中间件 ──> 密钥管理器
    │
    ▼
速率限制中间件
    │
    ▼
路由层 ──> 服务管理器
    │
    ▼
审计中间件
    │
    ▼
服务
    │
    ├─> 缓存层
    │
    ├─> 日志系统
    │
    └─> 指标系统
    │
    ▼
响应
```

## 相关文档

- [架构概览](/zh/concepts/architecture/overview/)
- [数据流](/zh/concepts/architecture/dataflow/)
- [功能特性](/zh/concepts/features/)
- [配置指南](/zh/guides/configuration/)
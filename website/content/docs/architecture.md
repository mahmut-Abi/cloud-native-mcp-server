---
title: "架构指南"
---

# 架构指南

本文档描述 Cloud Native MCP Server 的系统架构和设计原则。

## 目录

- [概述](#概述)
- [系统架构](#系统架构)
- [核心组件](#核心组件)
- [服务集成](#服务集成)
- [数据流](#数据流)
- [设计原则](#设计原则)
- [性能优化](#性能优化)
- [扩展性](#扩展性)

---

## 概述

Cloud Native MCP Server 是一个高性能的 Model Context Protocol (MCP) 服务器，用于管理 Kubernetes 和云原生基础设施。它采用模块化设计，支持多种运行模式和协议。

### 架构目标

- **高性能**: 优化的缓存、连接池和资源管理
- **可扩展性**: 模块化设计，易于添加新服务
- **安全性**: 多层认证、输入清理和审计日志
- **可观测性**: 内置指标、日志和追踪
- **可靠性**: 健康检查、重试机制和优雅降级

---

## 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                         客户端                               │
│  (Claude Desktop, Browser, Custom MCP Clients)              │
└────────────────────┬────────────────────────────────────────┘
                     │
                     │ MCP Protocol (SSE/HTTP/stdio)
                     │
┌────────────────────▼────────────────────────────────────────┐
│                    HTTP Server                               │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  路由层 (SSE/HTTP/Streamable-HTTP)                     │ │
│  └────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  中间件层                                               │ │
│  │  - 认证 (API Key/Bearer/Basic)                         │ │
│  │  - 审计日志                                             │ │
│  │  - 速率限制                                             │ │
│  │  - 安全中间件                                           │ │
│  │  - 指标收集                                             │ │
│  └────────────────────────────────────────────────────────┘ │
└────────────────────┬────────────────────────────────────────┘
                     │
                     │
┌────────────────────▼────────────────────────────────────────┐
│                  服务管理层                                 │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │Kubernetes│  │   Helm   │  │ Grafana  │  │Prometheus│  │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘  │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │  Kibana  │  │Elastic   │  │ AlertMgr │  │  Jaeger  │  │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘  │
│  ┌──────────┐  ┌──────────┐                               │
│  │  Otel    │  │Utilities │                               │
│  └──────────┘  └──────────┘                               │
└────────────────────┬────────────────────────────────────────┘
                     │
                     │
┌────────────────────▼────────────────────────────────────────┐
│                  基础设施层                                 │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  缓存层 (LRU/Segmented)                                │ │
│  └────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  密钥管理                                               │ │
│  └────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  日志系统                                               │ │
│  └────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  指标系统                                               │ │
│  └────────────────────────────────────────────────────────┘ │
└────────────────────┬────────────────────────────────────────┘
                     │
                     │
┌────────────────────▼────────────────────────────────────────┐
│                  外部服务                                   │
│  Kubernetes Cluster, Grafana, Prometheus, ES, etc.        │
└─────────────────────────────────────────────────────────────┘
```

---

## 核心组件

### 1. HTTP 服务器

**职责**: 处理传入的 HTTP/SSE 请求和连接

**特性**:
- 支持多种运行模式 (SSE, HTTP, stdio, Streamable-HTTP)
- 可配置的超时和连接限制
- 优雅关闭
- 健康检查端点

**关键文件**:
- `cmd/server/server.go`
- `internal/middleware/`

### 2. 路由层

**职责**: 将请求路由到正确的服务和工具

**特性**:
- 动态路由注册
- 路径参数解析
- 查询参数验证
- 错误处理

**关键文件**:
- `internal/services/registry.go`

### 3. 中间件层

**职责**: 在请求处理之前和之后执行通用逻辑

**中间件**:
- **认证**: API Key, Bearer Token, Basic Auth
- **审计日志**: 记录所有操作
- **速率限制**: 防止滥用
- **安全**: 输入清理和验证
- **指标**: 收集性能指标

**关键文件**:
- `internal/middleware/auth_middleware.go`
- `internal/middleware/audit_middleware.go`
- `internal/middleware/ratelimit.go`
- `internal/middleware/security_middleware.go`
- `internal/middleware/metrics_middleware.go`

### 4. 服务管理器

**职责**: 管理所有注册的服务和工具

**特性**:
- 服务注册和发现
- 工具调用路由
- 服务生命周期管理
- 健康检查协调

**关键文件**:
- `internal/services/manager/manager.go`

### 5. 缓存层

**职责**: 提供高性能缓存以减少外部服务调用

**特性**:
- LRU 缓存
- 分段缓存
- TTL 支持
- 缓存统计

**关键文件**:
- `internal/services/cache/`

### 6. 密钥管理器

**职责**: 安全地存储和管理敏感凭据

**特性**:
- 内存存储
- 密钥轮换
- 密钥生成
- 过期管理

**关键文件**:
- `internal/secrets/manager.go`

### 7. 日志系统

**职责**: 结构化日志记录

**特性**:
- 多级别日志 (debug, info, warn, error)
- JSON 和文本格式
- 结构化字段
- 上下文支持

**关键文件**:
- `internal/logging/logging.go`

### 8. 指标系统

**职责**: 收集和暴露性能指标

**特性**:
- Prometheus 格式
- 请求计数
- 延迟统计
- 缓存命中率

**关键文件**:
- `internal/observability/metrics/`

---

## 服务集成

### 服务接口

所有服务都实现统一的接口：

```go
type Service interface {
    // 服务名称
    Name() string

    // 初始化服务
    Initialize(config interface{}) error

    // 获取工具列表
    GetTools() []mcp.Tool

    // 调用工具
    CallTool(ctx context.Context, name string, arguments map[string]interface{}) (interface{}, error)

    // 健康检查
    HealthCheck() error

    // 关闭服务
    Shutdown() error
}
```

### 服务注册

服务在启动时自动注册：

```go
registry := services.NewRegistry()

// 注册服务
registry.Register(kubernetes.NewService())
registry.Register(grafana.NewService())
registry.Register(prometheus.NewService())
// ... 其他服务
```

### 工具调用流程

1. 客户端发送工具调用请求
2. 路由层解析请求，确定服务和工具
3. 中间件层执行认证、审计等
4. 服务管理器路由到正确的服务
5. 缓存层检查缓存
6. 服务执行工具调用
7. 结果返回给客户端
8. 审计日志记录操作

---

## 数据流

### 请求流

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
  │   └─> 解析服务和方法
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

### 响应流

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

---

## 设计原则

### 1. 模块化

每个服务都是独立的模块，可以单独启用/禁用：

```yaml
enableDisable:
  enabledServices: ["kubernetes", "helm", "prometheus"]
  disabledServices: ["elasticsearch", "kibana"]
```

### 2. 可扩展性

易于添加新服务：

1. 创建服务目录
2. 实现服务接口
3. 注册工具
4. 配置选项

### 3. 配置驱动

所有行为都通过配置控制：

- 服务启用/禁用
- 认证方式
- 缓存策略
- 日志级别

### 4. 故障隔离

服务故障不会影响其他服务：

```go
// 服务健康检查
func (s *Service) HealthCheck() error {
    if err := s.client.Ping(); err != nil {
        return fmt.Errorf("service unavailable: %w", err)
    }
    return nil
}
```

### 5. 优雅降级

服务不可用时返回友好错误：

```json
{
  "error": {
    "code": "SERVICE_UNAVAILABLE",
    "message": "Grafana service is temporarily unavailable",
    "details": {
      "service": "grafana",
      "retry_after": "30s"
    }
  }
}
```

---

## 性能优化

### 1. 缓存策略

#### LRU 缓存

```go
cache := cache.NewLRUCache(1000, 300*time.Second)
```

**适用场景**:
- 读取密集型操作
- 数据变化不频繁
- 高延迟操作

#### 分段缓存

```go
cache := cache.NewSegmentedCache(1000, 10, 300*time.Second)
```

**适用场景**:
- 不同类型的数据
- 需要不同的 TTL
- 并发访问

### 2. 连接池

```yaml
kubernetes:
  qps: 100.0
  burst: 200
  timeoutSec: 30
```

### 3. 响应压缩

```yaml
performance:
  compression_enabled: true
  compression_level: 6
```

### 4. JSON 编码池

```go
pool := json.NewEncoderPool(100, 8192)
```

### 5. 批处理

```go
// 批量获取资源
pods, err := k8sClient.CoreV1().Pods(namespace).List(ctx, options)
```

---

## 扩展性

### 添加新服务

1. **创建服务目录**

```bash
mkdir internal/services/myservice
```

2. **实现服务接口**

```go
package myservice

import (
    "context"
    "github.com/mahmut-Abi/cloud-native-mcp-server/internal/mcp"
)

type Service struct {
    config Config
    client *Client
}

func NewService() *Service {
    return &Service{}
}

func (s *Service) Name() string {
    return "myservice"
}

func (s *Service) Initialize(config interface{}) error {
    s.config = config.(Config)
    s.client = NewClient(s.config)
    return nil
}

func (s *Service) GetTools() []mcp.Tool {
    return []mcp.Tool{
        {
            Name:        "get_data",
            Description: "Get data from MyService",
            InputSchema: map[string]interface{}{
                "type": "object",
                "properties": map[string]interface{}{
                    "id": map[string]interface{}{
                        "type":        "string",
                        "description": "Data ID",
                    },
                },
                "required": []string{"id"},
            },
        },
    }
}

func (s *Service) CallTool(ctx context.Context, name string, args map[string]interface{}) (interface{}, error) {
    switch name {
    case "get_data":
        return s.GetData(ctx, args["id"].(string))
    default:
        return nil, fmt.Errorf("unknown tool: %s", name)
    }
}

func (s *Service) HealthCheck() error {
    return s.client.Ping()
}

func (s *Service) Shutdown() error {
    return s.client.Close()
}
```

3. **注册服务**

```go
// cmd/server/server.go
registry.Register(myservice.NewService())
```

4. **添加配置**

```yaml
# config.example.yaml
myservice:
  enabled: false
  url: "http://myservice:8080"
  apiKey: "${MYSERVICE_API_KEY}"
```

### 自定义工具

```go
// 添加自定义工具
func (s *Service) GetTools() []mcp.Tool {
    return []mcp.Tool{
        {
            Name:        "custom_tool",
            Description: "Custom tool description",
            InputSchema: map[string]interface{}{
                "type": "object",
                "properties": map[string]interface{}{
                    "param1": map[string]interface{}{
                        "type": "string",
                    },
                },
            },
        },
    }
}
```

---

## 可观测性

### 指标

#### 请求指标

```go
mcp_requests_total{method="kubernetes_list_pods",status="success"} 1234
mcp_request_duration_seconds{method="kubernetes_list_pods"} 0.123
```

#### 缓存指标

```go
mcp_cache_hits_total{service="kubernetes"} 456
mcp_cache_misses_total{service="kubernetes"} 78
```

#### 连接指标

```go
mcp_active_connections 10
mcp_total_connections 100
```

### 日志

#### 结构化日志

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

### 追踪

#### OpenTelemetry 集成

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

---

## 部署架构

### 单节点部署

```
┌─────────────────┐
│   MCP Server    │
│  (All Services) │
└────────┬────────┘
         │
         ├─> Kubernetes
         ├─> Grafana
         ├─> Prometheus
         └─> ...
```

### 多节点部署

```
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│  MCP Node 1  │  │  MCP Node 2  │  │  MCP Node 3  │
└──────┬───────┘  └──────┬───────┘  └──────┬───────┘
       │                 │                 │
       └─────────────────┴─────────────────┘
                         │
                         ▼
              ┌──────────────────┐
              │   Load Balancer  │
              └────────┬─────────┘
                       │
                       ▼
              ┌──────────────────┐
              │  External Services│
              └──────────────────┘
```

### 微服务部署

```
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│  MCP Gateway │  │  MCP Service │  │  MCP Service │
│   (Router)   │  │  (Kubernetes) │  │   (Grafana)  │
└──────┬───────┘  └──────────────┘  └──────────────┘
       │
       ▼
┌──────────────────┐
│  Service Mesh    │
│  (mTLS, Routing) │
└──────────────────┘
```

---

## 相关文档

- [完整工具参考](/docs/tools/)
- [配置指南](/docs/configuration/)
- [部署指南](/docs/deployment/)
- [性能指南](/docs/performance/)
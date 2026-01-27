---
title: "Architecture Guide"
---

# Architecture Guide

This document describes the system architecture and design principles of Cloud Native MCP Server.

## Table of Contents

- [Overview](#overview)
- [System Architecture](#system-architecture)
- [Core Components](#core-components)
- [Service Integration](#service-integration)
- [Data Flow](#data-flow)
- [Design Principles](#design-principles)
- [Performance Optimization](#performance-optimization)
- [Scalability](#scalability)

---

## Overview

Cloud Native MCP Server is a high-performance Model Context Protocol (MCP) server for managing Kubernetes and cloud-native infrastructure. It adopts a modular design with support for multiple runtime modes and protocols.

### Architecture Goals

- **High Performance**: Optimized caching, connection pooling, and resource management
- **Scalability**: Modular design, easy to add new services
- **Security**: Multi-layer authentication, input sanitization, and audit logging
- **Observability**: Built-in metrics, logging, and tracing
- **Reliability**: Health checks, retry mechanisms, and graceful degradation

---

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Client                               │
│  (Claude Desktop, Browser, Custom MCP Clients)              │
└────────────────────┬────────────────────────────────────────┘
                     │
                     │ MCP Protocol (SSE/HTTP/stdio)
                     │
┌────────────────────▼────────────────────────────────────────┐
│                    HTTP Server                               │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  Routing Layer (SSE/HTTP/Streamable-HTTP)              │ │
│  └────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  Middleware Layer                                       │ │
│  │  - Authentication (API Key/Bearer/Basic)               │ │
│  │  - Audit Logging                                        │ │
│  │  - Rate Limiting                                        │ │
│  │  - Security Middleware                                  │ │
│  │  - Metrics Collection                                   │ │
│  └────────────────────────────────────────────────────────┘ │
└────────────────────┬────────────────────────────────────────┘
                     │
                     │
┌────────────────────▼────────────────────────────────────────┐
│                  Service Management Layer                   │
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
│                  Infrastructure Layer                       │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  Cache Layer (LRU/Segmented)                           │ │
│  └────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  Secret Management                                      │ │
│  └────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  Logging System                                        │ │
│  └────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  Metrics System                                        │ │
│  └────────────────────────────────────────────────────────┘ │
└────────────────────┬────────────────────────────────────────┘
                     │
                     │
┌────────────────────▼────────────────────────────────────────┐
│                  External Services                          │
│  Kubernetes Cluster, Grafana, Prometheus, ES, etc.        │
└─────────────────────────────────────────────────────────────┘
```

---

## Core Components

### 1. HTTP Server

**Responsibility**: Handle incoming HTTP/SSE requests and connections

**Features**:
- Support for multiple runtime modes (SSE, HTTP, stdio, Streamable-HTTP)
- Configurable timeouts and connection limits
- Graceful shutdown
- Health check endpoints

**Key Files**:
- `cmd/server/server.go`
- `internal/middleware/`

### 2. Routing Layer

**Responsibility**: Route requests to the correct services and tools

**Features**:
- Dynamic routing registration
- Path parameter parsing
- Query parameter validation
- Error handling

**Key Files**:
- `internal/services/registry.go`

### 3. Middleware Layer

**Responsibility**: Execute common logic before and after request processing

**Middlewares**:
- **Authentication**: API Key, Bearer Token, Basic Auth
- **Audit Logging**: Record all operations
- **Rate Limiting**: Prevent abuse
- **Security**: Input sanitization and validation
- **Metrics**: Collect performance metrics

**Key Files**:
- `internal/middleware/auth_middleware.go`
- `internal/middleware/audit_middleware.go`
- `internal/middleware/ratelimit.go`
- `internal/middleware/security_middleware.go`
- `internal/middleware/metrics_middleware.go`

### 4. Service Manager

**Responsibility**: Manage all registered services and tools

**Features**:
- Service registration and discovery
- Tool call routing
- Service lifecycle management
- Health check coordination

**Key Files**:
- `internal/services/manager/manager.go`

### 5. Cache Layer

**Responsibility**: Provide high-performance caching to reduce external service calls

**Features**:
- LRU cache
- Segmented cache
- TTL support
- Cache statistics

**Key Files**:
- `internal/services/cache/`

### 6. Secret Manager

**Responsibility**: Securely store and manage sensitive credentials

**Features**:
- In-memory storage
- Key rotation
- Key generation
- Expiration management

**Key Files**:
- `internal/secrets/manager.go`

### 7. Logging System

**Responsibility**: Structured logging

**Features**:
- Multiple log levels (debug, info, warn, error)
- JSON and text formats
- Structured fields
- Context support

**Key Files**:
- `internal/logging/logging.go`

### 8. Metrics System

**Responsibility**: Collect and expose performance metrics

**Features**:
- Prometheus format
- Request counts
- Latency statistics
- Cache hit rates

**Key Files**:
- `internal/observability/metrics/`

---

## Service Integration

### Service Interface

All services implement a unified interface:

```go
type Service interface {
    // Service name
    Name() string

    // Initialize service
    Initialize(config interface{}) error

    // Get tool list
    GetTools() []mcp.Tool

    // Call tool
    CallTool(ctx context.Context, name string, arguments map[string]interface{}) (interface{}, error)

    // Health check
    HealthCheck() error

    // Shutdown service
    Shutdown() error
}
```

### Service Registration

Services are automatically registered at startup:

```go
registry := services.NewRegistry()

// Register services
registry.Register(kubernetes.NewService())
registry.Register(grafana.NewService())
registry.Register(prometheus.NewService())
// ... other services
```

### Tool Call Flow

1. Client sends tool call request
2. Routing layer parses request, determines service and tool
3. Middleware layer executes authentication, audit, etc.
4. Service manager routes to correct service
5. Cache layer checks cache
6. Service executes tool call
7. Result returned to client
8. Audit log records operation

---

## Data Flow

### Request Flow

```
Client
  │
  ├─> HTTP/SSE Connection
  │
  ├─> Authentication Middleware
  │   ├─> Validate API Key/Token
  │   └─> Check permissions
  │
  ├─> Rate Limiting Middleware
  │   └─> Check quota
  │
  ├─> Routing Layer
  │   └─> Parse service and method
  │
  ├─> Audit Middleware
  │   └─> Record request start
  │
  ├─> Service Manager
  │   └─> Route to service
  │
  ├─> Cache Layer
  │   ├─> Check cache
  │   └─> Return cache or continue
  │
  ├─> Service
  │   ├─> Call external API
  │   ├─> Process response
  │   └─> Update cache
  │
  ├─> Audit Middleware
  │   └─> Record request completion
  │
  ├─> Metrics Middleware
  │   └─> Record metrics
  │
  └─> Response returned to client
```

### Response Flow

```
Service
  │
  ├─> Process result
  │
  ├─> Data Transformation
  │   ├─> Formatting
  │   └─> Compression
  │
  ├─> Cache Update
  │   └─> Store in cache
  │
  ├─> Metrics Update
  │   └─> Record performance metrics
  │
  └─> Return response
```

---

## Design Principles

### 1. Modularity

Each service is an independent module that can be enabled/disabled individually:

```yaml
enableDisable:
  enabledServices: ["kubernetes", "helm", "prometheus"]
  disabledServices: ["elasticsearch", "kibana"]
```

### 2. Scalability

Easy to add new services:

1. Create service directory
2. Implement service interface
3. Register tools
4. Configure options

### 3. Configuration Driven

All behavior is controlled through configuration:

- Service enable/disable
- Authentication method
- Cache strategy
- Log level

### 4. Fault Isolation

Service failures don't affect other services:

```go
// Service health check
func (s *Service) HealthCheck() error {
    if err := s.client.Ping(); err != nil {
        return fmt.Errorf("service unavailable: %w", err)
    }
    return nil
}
```

### 5. Graceful Degradation

Return friendly errors when services are unavailable:

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

## Performance Optimization

### 1. Caching Strategy

#### LRU Cache

```go
cache := cache.NewLRUCache(1000, 300*time.Second)
```

**Use Cases**:
- Read-intensive operations
- Infrequently changing data
- High latency operations

#### Segmented Cache

```go
cache := cache.NewSegmentedCache(1000, 10, 300*time.Second)
```

**Use Cases**:
- Different types of data
- Need different TTLs
- Concurrent access

### 2. Connection Pooling

```yaml
kubernetes:
  qps: 100.0
  burst: 200
  timeoutSec: 30
```

### 3. Response Compression

```yaml
performance:
  compression_enabled: true
  compression_level: 6
```

### 4. JSON Encoding Pool

```go
pool := json.NewEncoderPool(100, 8192)
```

### 5. Batching

```go
// Batch fetch resources
pods, err := k8sClient.CoreV1().Pods(namespace).List(ctx, options)
```

---

## Scalability

### Adding a New Service

1. **Create Service Directory**

```bash
mkdir internal/services/myservice
```

2. **Implement Service Interface**

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

3. **Register Service**

```go
// cmd/server/server.go
registry.Register(myservice.NewService())
```

4. **Add Configuration**

```yaml
# config.example.yaml
myservice:
  enabled: false
  url: "http://myservice:8080"
  apiKey: "${MYSERVICE_API_KEY}"
```

### Custom Tools

```go
// Add custom tools
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

## Observability

### Metrics

#### Request Metrics

```go
mcp_requests_total{method="kubernetes_list_pods",status="success"} 1234
mcp_request_duration_seconds{method="kubernetes_list_pods"} 0.123
```

#### Cache Metrics

```go
mcp_cache_hits_total{service="kubernetes"} 456
mcp_cache_misses_total{service="kubernetes"} 78
```

#### Connection Metrics

```go
mcp_active_connections 10
mcp_total_connections 100
```

### Logging

#### Structured Logging

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

### Tracing

#### OpenTelemetry Integration

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

tracer := otel.Tracer("cloud-native-mcp-server")

ctx, span := tracer.Start(ctx, "list_pods")
defer span.End()

// Execute operation
pods, err := k8sClient.ListPods(ctx, namespace)
```

---

## Deployment Architecture

### Single Node Deployment

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

### Multi-Node Deployment

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

### Microservices Deployment

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

## Related Documentation

- [Complete Tools Reference](/docs/tools/)
- [Configuration Guide](/docs/configuration/)
- [Deployment Guide](/docs/deployment/)
- [Performance Guide](/docs/performance/)
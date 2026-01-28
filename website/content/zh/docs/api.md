---
title: "API 文档"
weight: 10
---

# API 文档

Cloud Native MCP Server 通过 Model Context Protocol (MCP) 提供全面的 API，用于与所有集成服务进行交互。

## 基础 URL

```
http://localhost:8080 (默认)
```

## 认证

服务器支持多种认证方法：

### API 密钥
```bash
curl -X POST http://localhost:8080/v1/mcp/list-tools \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json"
```

### Bearer Token (JWT)
```bash
curl -X POST http://localhost:8080/v1/mcp/list-tools \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json"
```

### Basic Auth
```bash
curl -X POST http://localhost:8080/v1/mcp/list-tools \
  -u "username:password" \
  -H "Content-Type: application/json"
```

## 核心端点

### MCP 端点

#### 列出可用工具
```
POST /v1/mcp/list-tools
```

获取所有服务中所有可用工具的列表。

**请求:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "mcp/list-tools",
  "params": {}
}
```

**响应:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "tools": [
      {
        "name": "kubernetes-get-pods",
        "description": "获取命名空间中的 Pod",
        "input_schema": {
          "type": "object",
          "properties": {
            "namespace": {
              "type": "string",
              "description": "Kubernetes 命名空间"
            }
          }
        }
      }
    ]
  }
}
```

#### 执行工具
```
POST /v1/mcp/call-tool
```

使用给定参数执行特定工具。

**请求:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "kubernetes-get-pods",
    "arguments": {
      "namespace": "default"
    }
  }
}
```

**响应:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "success": true,
    "data": {
      "pods": [
        {
          "name": "my-app-12345",
          "status": "Running",
          "restarts": 0
        }
      ]
    }
  }
}
```

### 健康检查
```
GET /health
```

检查服务器的健康状态。

**响应:**
```json
{
  "status": "healthy",
  "timestamp": "2023-10-01T12:00:00Z",
  "services": {
    "kubernetes": "connected",
    "prometheus": "connected",
    "grafana": "connected"
  }
}
```

### 服务器信息
```
GET /info
```

获取服务器信息和版本详情。

**响应:**
```json
{
  "version": "1.0.0",
  "build_date": "2023-10-01T12:00:00Z",
  "go_version": "go1.21",
  "services_count": 10,
  "tools_count": 220
}
```

## 支持的协议

服务器支持多种通信协议：

### SSE (Server-Sent Events) - 默认
```
POST /api/aggregate/sse
```

### HTTP
```
POST /api/aggregate/http
```

### Stdio
在 stdio 模式下运行时可用。

## 错误处理

API 使用标准 JSON-RPC 2.0 错误格式：

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

### 常见错误代码

- `-32700`: 解析错误
- `-32600`: 无效请求
- `-32601`: 方法未找到
- `-32602`: 无效参数
- `-32603`: 内部错误
- `401`: 未授权
- `403`: 禁止访问
- `404`: 未找到
- `429`: 速率限制
- `500`: 内部服务器错误

## 速率限制

服务器实现速率限制以防止滥用：

- **每分钟请求数**: 1000 (可配置)
- **突发限制**: 100 (可配置)

## 配置选项

### 服务器配置

您可以使用以下环境变量配置服务器：

- `MCP_SERVER_ADDR`: 服务器地址 (默认: `0.0.0.0:8080`)
- `MCP_SERVER_MODE`: 通信模式 (sse, http, stdio) (默认: `sse`)
- `MCP_SERVER_API_KEY`: 认证 API 密钥
- `MCP_SERVER_BEARER_TOKEN`: JWT 认证令牌
- `MCP_SERVER_RATE_LIMIT`: 每分钟请求数 (默认: `1000`)
- `MCP_SERVER_BURST_LIMIT`: 突发限制 (默认: `100`)

### 服务配置

每个集成服务都可以使用特定的环境变量进行配置：

- **Kubernetes**: `KUBECONFIG` 或集群内配置
- **Prometheus**: `MCP_PROMETHEUS_URL`
- **Grafana**: `MCP_GRAFANA_URL`, `MCP_GRAFANA_API_KEY`
- **Elasticsearch**: `MCP_ELASTICSEARCH_URL`
- **Alertmanager**: `MCP_ALERTMANAGER_URL`

## 示例

### 使用 cURL

```bash
curl -X POST http://localhost:8080/v1/mcp/call-tool \
  -H "Authorization: Bearer your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "kubernetes-get-pods",
      "arguments": {
        "namespace": "default"
      }
    }
  }'
```

### 使用 Python

```python
import requests
import json

def call_tool(tool_name, arguments):
    headers = {
        "Authorization": "Bearer YOUR_API_KEY",
        "Content-Type": "application/json"
    }
    
    payload = {
        "jsonrpc": "2.0",
        "id": 1,
        "method": "tools/call",
        "params": {
            "name": tool_name,
            "arguments": arguments
        }
    }
    
    response = requests.post(
        "http://localhost:8080/v1/mcp/call-tool",
        headers=headers,
        data=json.dumps(payload)
    )
    
    return response.json()

# 示例用法
result = call_tool("kubernetes-get-pods", {"namespace": "default"})
print(result)
```

## 下一步

- [工具参考](/docs/tools/) 了解详细工具文档
- [配置指南](/guides/configuration/) 了解设置说明
- [安全最佳实践](/guides/security/) 了解 API 安全
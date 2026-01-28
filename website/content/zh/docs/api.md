---
title: "API Documentation"
weight: 10
---

# API Documentation

The Cloud Native MCP Server provides a comprehensive API for interacting with all integrated services through the Model Context Protocol (MCP).

## Base URL

```
http://localhost:8080 (default)
```

## Authentication

The server supports multiple authentication methods:

### API Key
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

## Core Endpoints

### MCP Endpoints

#### List Available Tools
```
POST /v1/mcp/list-tools
```

Get a list of all available tools across all services.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "mcp/list-tools",
  "params": {}
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "tools": [
      {
        "name": "kubernetes-get-pods",
        "description": "Get pods in a namespace",
        "input_schema": {
          "type": "object",
          "properties": {
            "namespace": {
              "type": "string",
              "description": "Kubernetes namespace"
            }
          }
        }
      }
    ]
  }
}
```

#### Execute Tool
```
POST /v1/mcp/call-tool
```

Execute a specific tool with given parameters.

**Request:**
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

**Response:**
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

### Health Check
```
GET /health
```

Check the health status of the server.

**Response:**
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

### Server Info
```
GET /info
```

Get server information and version details.

**Response:**
```json
{
  "version": "1.0.0",
  "build_date": "2023-10-01T12:00:00Z",
  "go_version": "go1.21",
  "services_count": 10,
  "tools_count": 220
}
```

## Supported Protocols

The server supports multiple communication protocols:

### SSE (Server-Sent Events) - Default
```
POST /api/aggregate/sse
```

### HTTP
```
POST /api/aggregate/http
```

### Stdio
Available when running in stdio mode.

## Error Handling

The API uses standard JSON-RPC 2.0 error format:

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

### Common Error Codes

- `-32700`: Parse error
- `-32600`: Invalid Request
- `-32601`: Method not found
- `-32602`: Invalid params
- `-32603`: Internal error
- `401`: Unauthorized
- `403`: Forbidden
- `404`: Not found
- `429`: Rate limited
- `500`: Internal server error

## Rate Limiting

The server implements rate limiting to prevent abuse:

- **Requests per minute**: 1000 (configurable)
- **Burst limit**: 100 (configurable)

## Configuration Options

### Server Configuration

You can configure the server with the following environment variables:

- `MCP_SERVER_ADDR`: Server address (default: `0.0.0.0:8080`)
- `MCP_SERVER_MODE`: Communication mode (sse, http, stdio) (default: `sse`)
- `MCP_SERVER_API_KEY`: API key for authentication
- `MCP_SERVER_BEARER_TOKEN`: JWT token for authentication
- `MCP_SERVER_RATE_LIMIT`: Requests per minute (default: `1000`)
- `MCP_SERVER_BURST_LIMIT`: Burst limit (default: `100`)

### Service Configuration

Each integrated service can be configured with specific environment variables:

- **Kubernetes**: `KUBECONFIG` or in-cluster configuration
- **Prometheus**: `MCP_PROMETHEUS_URL`
- **Grafana**: `MCP_GRAFANA_URL`, `MCP_GRAFANA_API_KEY`
- **Elasticsearch**: `MCP_ELASTICSEARCH_URL`
- **Alertmanager**: `MCP_ALERTMANAGER_URL`

## Examples

### Using with cURL

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

### Using with Python

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

# Example usage
result = call_tool("kubernetes-get-pods", {"namespace": "default"})
print(result)
```

## Next Steps

- [Tools Reference](/docs/tools/) for detailed tool documentation
- [Configuration Guides](/guides/configuration/) for setup instructions
- [Security Best Practices](/guides/security/) for securing your API
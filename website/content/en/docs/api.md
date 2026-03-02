---
title: "API Documentation"
weight: 80
description: "Transport endpoints, authentication, runtime APIs, and integration patterns for Cloud Native MCP Server."
---

# API Documentation

This page documents the runtime endpoints exposed by Cloud Native MCP Server.

## Base Address

```text
http://localhost:8080
```

## Transport Modes and Endpoints

Cloud Native MCP Server supports four runtime modes:

| Mode | Typical Usage | Aggregate Endpoint |
| --- | --- | --- |
| `sse` | Broad MCP client compatibility | `/api/aggregate/sse` |
| `streamable-http` | Modern MCP transport | `/api/aggregate/streamable-http` |
| `http` | Legacy message endpoint compatibility | `/api/aggregate/sse/message` |
| `stdio` | Local runtime integration | stdin/stdout |

Service-specific endpoint pattern:

- SSE: `/api/<service>/sse`
- Streamable HTTP: `/api/<service>/streamable-http`

Common service names include `kubernetes`, `helm`, `grafana`, `prometheus`, `kibana`, `elasticsearch`, `alertmanager`, `jaeger`, `opentelemetry`, `utilities`, and `aggregate`.

> Note: This server does **not** expose legacy `/v1/mcp/*` REST-style endpoints.

---

## Authentication

Enable auth in runtime configuration:

```bash
export MCP_AUTH_ENABLED=true
export MCP_AUTH_MODE=apikey
export MCP_AUTH_API_KEY='ChangeMe-Strong-Key-123!'
```

Supported auth modes:

- `apikey`
- `bearer`
- `basic`

### API Key

Use `X-Api-Key` header (recommended) or `api_key` query parameter.

```bash
curl -sS -N \
  -H "X-Api-Key: ChangeMe-Strong-Key-123!" \
  http://127.0.0.1:8080/api/aggregate/sse
```

### Bearer Token

```bash
export MCP_AUTH_MODE=bearer
export MCP_AUTH_BEARER_TOKEN='your-jwt-or-token'

curl -sS -N \
  -H "Authorization: Bearer your-jwt-or-token" \
  http://127.0.0.1:8080/api/aggregate/sse
```

### Basic Auth

```bash
export MCP_AUTH_MODE=basic
export MCP_AUTH_USERNAME='admin'
export MCP_AUTH_PASSWORD='strong-password'

curl -sS -N -u "admin:strong-password" \
  http://127.0.0.1:8080/api/aggregate/sse
```

---

## SSE Workflow

In `sse` mode, integration typically follows this flow:

1. Open SSE stream on `/api/aggregate/sse`.
2. Receive endpoint event containing message endpoint.
3. POST JSON-RPC requests (for example `initialize`) to the message endpoint.

Recommended validation command:

```bash
make sse-smoke BASE_URL=http://127.0.0.1:8080
```

If you are outside repository root:

```bash
/path/to/cloud-native-mcp-server/scripts/sse_smoke_test.sh http://127.0.0.1:8080
```

---

## Streamable HTTP Example

In `streamable-http` mode, use the aggregate streamable endpoint directly:

```bash
curl -sS -X POST \
  http://127.0.0.1:8080/api/aggregate/streamable-http \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "initialize",
    "params": {
      "protocolVersion": "2024-11-05",
      "clientInfo": {
        "name": "manual-client",
        "version": "1.0.0"
      }
    }
  }'
```

For production integration, prefer MCP SDK/client implementations over manual curl payload construction.

---

## Runtime Endpoints

| Endpoint | Description |
| --- | --- |
| `GET /health` | Server health check |
| `GET /metrics` | Prometheus metrics endpoint (may require auth if enabled) |
| `GET /api/openapi.json` | OpenAPI schema for HTTP endpoints |
| `GET /api/docs` | Swagger UI for API exploration |
| `GET /api/audit/logs` | Audit query endpoint (when audit is enabled) |
| `GET /api/audit/stats` | Audit stats endpoint (when audit is enabled) |

---

## Key Environment Variables

### Server and Transport

- `MCP_MODE` (`sse`, `streamable-http`, `http`, `stdio`)
- `MCP_ADDR`
- `MCP_READ_TIMEOUT`
- `MCP_WRITE_TIMEOUT`
- `MCP_IDLE_TIMEOUT`

### Authentication

- `MCP_AUTH_ENABLED`
- `MCP_AUTH_MODE`
- `MCP_AUTH_API_KEY`
- `MCP_AUTH_BEARER_TOKEN`
- `MCP_AUTH_USERNAME`
- `MCP_AUTH_PASSWORD`

### Rate Limiting

- `MCP_RATELIMIT_ENABLED`
- `MCP_RATELIMIT_REQUESTS_PER_SECOND`
- `MCP_RATELIMIT_BURST`

### Service Selection

- `MCP_ENABLED_SERVICES`
- `MCP_DISABLED_SERVICES`

### Integration Examples

- `MCP_PROM_ADDRESS`
- `MCP_GRAFANA_URL`
- `MCP_KIBANA_URL`
- `MCP_ELASTICSEARCH_ADDRESS`
- `MCP_ALERTMANAGER_ADDRESS`

---

## Troubleshooting API Access

- `401 unauthorized`: verify `MCP_AUTH_*` settings and credential transmission method.
- stream opens but requests fail: test with `make sse-smoke` and inspect logs.
- missing services/tools: check `MCP_ENABLED_SERVICES` / `MCP_DISABLED_SERVICES`.
- slow calls: narrow query scope and review upstream service connectivity.

## Next Steps

- [Getting Started]({{< relref "/getting-started/_index.md" >}})
- [Getting Started FAQ]({{< relref "/getting-started/faq.md" >}})
- [Troubleshooting]({{< relref "/getting-started/troubleshooting.md" >}})
- [Tools Reference]({{< relref "tools.md" >}})
- [Configuration Guide]({{< relref "configuration.md" >}})
- [Security Guide]({{< relref "security.md" >}})

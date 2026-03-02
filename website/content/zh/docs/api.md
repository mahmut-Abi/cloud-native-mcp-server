---
title: "API 文档"
weight: 80
description: "Cloud Native MCP Server 的传输端点、认证方式、运行时接口与接入示例。"
---

# API 文档

本页面描述 Cloud Native MCP Server 当前版本实际暴露的运行时接口。

## 基础地址

```text
http://localhost:8080
```

## 运行模式与端点

Cloud Native MCP Server 支持四种运行模式：

| 模式 | 典型用途 | 聚合入口 |
| --- | --- | --- |
| `sse` | MCP 客户端兼容性优先 | `/api/aggregate/sse` |
| `streamable-http` | 推荐的现代 MCP 传输方式 | `/api/aggregate/streamable-http` |
| `http` | 历史 message 端点兼容 | `/api/aggregate/sse/message` |
| `stdio` | 本地运行时集成 | stdin/stdout |

服务级端点模式：

- SSE：`/api/<service>/sse`
- Streamable HTTP：`/api/<service>/streamable-http`

常见 service 名称包括 `kubernetes`、`helm`、`grafana`、`prometheus`、`kibana`、`elasticsearch`、`alertmanager`、`jaeger`、`opentelemetry`、`utilities`、`aggregate`。

> 说明：当前服务不提供旧版 `/v1/mcp/*` 风格接口。

---

## 认证方式

运行时启用认证示例：

```bash
export MCP_AUTH_ENABLED=true
export MCP_AUTH_MODE=apikey
export MCP_AUTH_API_KEY='ChangeMe-Strong-Key-123!'
```

支持模式：

- `apikey`
- `bearer`
- `basic`

### API Key

建议优先使用 `X-Api-Key` 请求头，也支持 `api_key` 查询参数。

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

## SSE 接入流程

在 `sse` 模式下，典型流程如下：

1. 连接 `/api/aggregate/sse` 建立事件流。
2. 从返回事件中获取 message endpoint。
3. 向 message endpoint 发送 JSON-RPC 请求（例如 `initialize`）。

推荐使用内置自检命令验证链路：

```bash
make sse-smoke BASE_URL=http://127.0.0.1:8080
```

若不在仓库根目录：

```bash
/path/to/cloud-native-mcp-server/scripts/sse_smoke_test.sh http://127.0.0.1:8080
```

---

## Streamable HTTP 示例

在 `streamable-http` 模式下，可直接使用聚合端点：

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

生产接入建议优先使用 MCP SDK/客户端实现，避免手写协议细节。

---

## 运行时接口

| 端点 | 说明 |
| --- | --- |
| `GET /health` | 服务健康检查 |
| `GET /metrics` | Prometheus 指标接口（启用认证后可能需要鉴权） |
| `GET /api/openapi.json` | HTTP 端点 OpenAPI Schema |
| `GET /api/docs` | Swagger UI 页面 |
| `GET /api/audit/logs` | 审计日志查询接口（启用审计后可用） |
| `GET /api/audit/stats` | 审计统计接口（启用审计后可用） |

---

## 关键环境变量

### 服务与传输

- `MCP_MODE`（`sse`、`streamable-http`、`http`、`stdio`）
- `MCP_ADDR`
- `MCP_READ_TIMEOUT`
- `MCP_WRITE_TIMEOUT`
- `MCP_IDLE_TIMEOUT`

### 认证

- `MCP_AUTH_ENABLED`
- `MCP_AUTH_MODE`
- `MCP_AUTH_API_KEY`
- `MCP_AUTH_BEARER_TOKEN`
- `MCP_AUTH_USERNAME`
- `MCP_AUTH_PASSWORD`

### 限流

- `MCP_RATELIMIT_ENABLED`
- `MCP_RATELIMIT_REQUESTS_PER_SECOND`
- `MCP_RATELIMIT_BURST`

### 服务启用控制

- `MCP_ENABLED_SERVICES`
- `MCP_DISABLED_SERVICES`

### 第三方集成示例

- `MCP_PROM_ADDRESS`
- `MCP_GRAFANA_URL`
- `MCP_KIBANA_URL`
- `MCP_ELASTICSEARCH_ADDRESS`
- `MCP_ALERTMANAGER_ADDRESS`

---

## API 调用排障建议

- 出现 `401 unauthorized`：检查 `MCP_AUTH_*` 配置与凭据传递方式。
- SSE 建连成功但请求失败：执行 `make sse-smoke` 并结合日志定位。
- 服务或工具缺失：检查 `MCP_ENABLED_SERVICES` / `MCP_DISABLED_SERVICES`。
- 调用慢或超时：缩小请求范围并检查上游服务可达性。

## 下一步

- [快速开始]({{< relref "/getting-started/_index.md" >}})
- [快速开始 FAQ]({{< relref "/getting-started/faq.md" >}})
- [故障排除]({{< relref "/getting-started/troubleshooting.md" >}})
- [工具参考]({{< relref "tools.md" >}})
- [配置指南]({{< relref "configuration.md" >}})
- [安全指南]({{< relref "security.md" >}})

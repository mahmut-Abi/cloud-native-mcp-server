---
title: "Getting Started"
weight: 1
description: "Install Cloud Native MCP Server, validate endpoints, and move to production with a safe baseline."
---

# Getting Started

This guide helps you install Cloud Native MCP Server, validate connectivity, and prepare for production rollout.

## What You Will Set Up

- Run the server in one of two modes: `sse`, `streamable-http`
- Enable authentication with the correct environment variables
- Verify runtime health and MCP handshake behavior
- Continue with FAQ and troubleshooting playbooks

---

## Prerequisites

- Kubernetes access (`~/.kube/config` or in-cluster credentials)
- Docker or a Linux host for binary execution
- Go `1.25+` (only needed for source build)
- Network access to observability backends you plan to integrate

---

## Installation Options

{{< tabs >}}
{{< tab "Docker" >}}
```bash
docker run -d \
  --name cloud-native-mcp-server \
  -p 8080:8080 \
  -v ~/.kube:/root/.kube:ro \
  -e MCP_AUTH_ENABLED=true \
  -e MCP_AUTH_MODE=apikey \
  -e MCP_AUTH_API_KEY='ChangeMe-Strong-Key-123!' \
  mahmutabi/cloud-native-mcp-server:latest
```
{{< /tab >}}

{{< tab "Binary" >}}
```bash
curl -LO https://github.com/mahmut-Abi/cloud-native-mcp-server/releases/latest/download/cloud-native-mcp-server-linux-amd64
chmod +x cloud-native-mcp-server-linux-amd64
./cloud-native-mcp-server-linux-amd64 --mode=sse --addr=0.0.0.0:8080
```
{{< /tab >}}

{{< tab "Source" >}}
```bash
git clone https://github.com/mahmut-Abi/cloud-native-mcp-server.git
cd cloud-native-mcp-server
make build
./cloud-native-mcp-server --mode=streamable-http --addr=0.0.0.0:8080
```
{{< /tab >}}
{{< /tabs >}}

---

## Choose a Run Mode

| Mode | Recommended For | Main Endpoint |
| --- | --- | --- |
| `sse` | Broad MCP client compatibility | `/api/aggregate/sse` |
| `streamable-http` | Modern MCP transport in production | `/api/aggregate/streamable-http` |

---

## First Validation

Run these checks after startup:

```bash
# Health check
curl -sS http://127.0.0.1:8080/health

# End-to-end SSE handshake + initialize
make sse-smoke BASE_URL=http://127.0.0.1:8080
```

If you are not in the repository directory, run the script directly:

```bash
/path/to/cloud-native-mcp-server/scripts/sse_smoke_test.sh http://127.0.0.1:8080
```

---

## Authentication Check

When `MCP_AUTH_ENABLED=true` and `MCP_AUTH_MODE=apikey`:

```bash
# SSE stream request with API key
curl -sS -N "http://127.0.0.1:8080/api/aggregate/sse?api_key=ChangeMe-Strong-Key-123!"
```

You can also pass the key via request header:

```bash
curl -sS -N \
  -H "X-Api-Key: ChangeMe-Strong-Key-123!" \
  http://127.0.0.1:8080/api/aggregate/sse
```

---

## Common Runtime Settings

```bash
# Server mode and bind address
export MCP_MODE=sse
export MCP_ADDR=0.0.0.0:8080

# Authentication (apikey mode)
export MCP_AUTH_ENABLED=true
export MCP_AUTH_MODE=apikey
export MCP_AUTH_API_KEY='ChangeMe-Strong-Key-123!'

# Optional: disable non-required services
export MCP_DISABLED_SERVICES="kibana,jaeger"
```

---

## Next Steps

- [Getting Started FAQ]({{< relref "faq.md" >}})
- [Troubleshooting]({{< relref "troubleshooting.md" >}})
- [Security Guide]({{< relref "/docs/security.md" >}})
- [Configuration Guide]({{< relref "/docs/configuration.md" >}})
- [Performance Guide]({{< relref "/docs/performance.md" >}})
- [Tools Reference]({{< relref "/docs/tools.md" >}})

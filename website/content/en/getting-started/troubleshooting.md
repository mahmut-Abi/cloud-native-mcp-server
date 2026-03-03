---
title: "Troubleshooting"
weight: 30
description: "Step-by-step checks for startup, auth, transport, and service integration issues."
---

# Troubleshooting

Use this page as a practical playbook when setup or runtime behavior is not as expected.

## 1. Server fails to start

Check port conflicts and startup logs:

```bash
# Check if port is occupied
ss -lntp | rg 8080

# Start with debug logs
./cloud-native-mcp-server --mode=sse --addr=127.0.0.1:8080 --log-level=debug
```

If using Docker, inspect container logs:

```bash
docker logs --tail=200 cloud-native-mcp-server
```

## 2. `/health` is unreachable or non-200

```bash
curl -sv http://127.0.0.1:8080/health
```

If connection fails:

- Confirm process/container is running.
- Confirm bind address and published port.
- Confirm firewall or security group rules.

## 3. SSE endpoint opens but handshake fails

```bash
curl -svN --connect-timeout 5 --max-time 15 \
  -H "Accept: text/event-stream" \
  "http://127.0.0.1:8080/api/aggregate/sse"
```

Then run the built-in smoke test:

```bash
make sse-smoke BASE_URL=http://127.0.0.1:8080
```

If smoke test fails, inspect the first stream events and message endpoint behavior.

## 4. Authentication keeps returning 401

Verify runtime auth configuration:

```bash
export MCP_AUTH_ENABLED=true
export MCP_AUTH_MODE=apikey
export MCP_AUTH_API_KEY='ChangeMe-Strong-Key-123!'
```

Validate using both supported API key styles:

```bash
curl -sS -N -H "X-Api-Key: ChangeMe-Strong-Key-123!" \
  http://127.0.0.1:8080/api/aggregate/sse

curl -sS -N "http://127.0.0.1:8080/api/aggregate/sse?api_key=ChangeMe-Strong-Key-123!"
```

## 5. Tool calls are slow or time out

Start with high-impact checks:

- Confirm upstream backends (Kubernetes API, Prometheus, Grafana, etc.) are reachable.
- Reduce scope of each request (namespace, time range, object count).
- Enable pagination/summary patterns.
- Tune timeout and rate-limit settings in config.

## 6. Some services are unavailable

List enabled services and tools:

```bash
./cloud-native-mcp-server --list=services --output=table
./cloud-native-mcp-server --list=tools --service=kubernetes --output=table
```

If a service is unexpectedly missing, check:

- `MCP_ENABLED_SERVICES`
- `MCP_DISABLED_SERVICES`
- service-specific credentials and endpoint variables

## 7. Gather useful diagnostics before filing an issue

Include these details in your issue report:

- server mode (`sse` / `streamable-http`)
- startup command and key env vars (mask secrets)
- `curl` output for `/health` and transport endpoint
- relevant logs (last 100-200 lines)

## Related Pages

- [Getting Started]({{< relref "_index.md" >}})
- [FAQ]({{< relref "faq.md" >}})
- [Configuration Guide]({{< relref "/docs/configuration.md" >}})
- [Security Guide]({{< relref "/docs/security.md" >}})

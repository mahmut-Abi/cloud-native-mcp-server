---
title: "Performance Guide"
---

# Performance Guide

This guide focuses on performance tuning that matches the **current** Cloud Native MCP Server implementation.

---

## What You Can Tune Today

### 1. Server Timeouts

```yaml
server:
  readTimeoutSec: 30
  writeTimeoutSec: 0
  idleTimeoutSec: 60
```

Recommendations:
- Keep `writeTimeoutSec: 0` for SSE-heavy workloads.
- Increase `readTimeoutSec` for slow clients or large request bodies.
- Increase `idleTimeoutSec` for long-lived clients behind proxies.

### 2. Kubernetes Client Throughput

```yaml
kubernetes:
  timeoutSec: 30
  qps: 100.0
  burst: 200
```

Recommendations:
- Increase `qps`/`burst` for larger clusters.
- Increase `timeoutSec` for heavy list/watch or expensive queries.

### 3. Request Rate Limiting

```yaml
ratelimit:
  enabled: true
  requests_per_second: 100
  burst: 200
```

Recommendations:
- Enable in multi-tenant or internet-facing deployments.
- Start with conservative values and raise gradually from metrics.

### 4. Reduce Unnecessary Service Overhead

```yaml
enableDisable:
  enabledServices: ["kubernetes", "prometheus", "grafana", "aggregate"]
```

Recommendations:
- Only enable services you actually use.
- Disable unused tools for lower memory and startup overhead.

### 5. Audit Cost Control (When Audit Is Enabled)

```yaml
audit:
  enabled: true
  storage: "database"
  sampling:
    enabled: true
    rate: 0.3
```

Recommendations:
- Use sampling on high-QPS workloads.
- Use `storage: database` for large, query-heavy audit datasets.

---

## Built-In Optimizations (No Public YAML Key)

The server already includes internal optimizations such as:
- response truncation safeguards
- efficient JSON processing paths
- internal caching and pooling in service/tool layers

These are implementation details, not stable public config keys.
Use observable knobs above (`server`, `kubernetes`, `ratelimit`, `enableDisable`, `audit`) first.

---

## Performance Metrics to Track

Key metrics to watch:
- request rate
- p95/p99 latency
- error rate
- active connections
- memory and CPU usage

Example queries:

```promql
rate(mcp_requests_total[5m])
```

```promql
histogram_quantile(0.99, rate(mcp_request_duration_seconds_bucket[5m]))
```

```promql
rate(mcp_errors_total[5m])
```

---

## Benchmarking

### Health Endpoint

```bash
ab -n 10000 -c 100 http://127.0.0.1:8080/health
```

### Tool Call Endpoint (Legacy Message Endpoint Compatibility)

Create payload:

```json
{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}
```

Run benchmark:

```bash
ab -n 1000 -c 10 \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-key" \
  -p payload.json \
  http://127.0.0.1:8080/api/kubernetes/sse/message
```

---

## Production Baseline Example

```yaml
server:
  mode: "streamable-http"
  addr: "0.0.0.0:8080"
  readTimeoutSec: 30
  writeTimeoutSec: 0
  idleTimeoutSec: 60

logging:
  level: "info"
  json: true

kubernetes:
  kubeconfig: ""
  timeoutSec: 30
  qps: 100.0
  burst: 200

ratelimit:
  enabled: true
  requests_per_second: 100
  burst: 200

auth:
  enabled: true
  mode: "apikey"
  apiKey: "${MCP_AUTH_API_KEY}"

audit:
  enabled: true
  storage: "database"
  sampling:
    enabled: true
    rate: 0.3

enableDisable:
  enabledServices: ["kubernetes", "prometheus", "grafana", "aggregate"]
```

---

## Troubleshooting

### High Latency

Actions:
1. Increase `kubernetes.timeoutSec` if backend queries are slow.
2. Increase `kubernetes.qps`/`kubernetes.burst` for API-client bottlenecks.
3. Check p99 latency and correlate with backend service saturation.

### High Error Rate

Actions:
1. Verify auth mode and credentials.
2. Check backend service availability (Prometheus/Grafana/etc.).
3. Inspect audit logs via `/api/audit/logs` if audit is enabled.

### High Memory Usage

Actions:
1. Reduce enabled services/tools.
2. Enable audit sampling.
3. Reduce burst limits if traffic spikes cause memory pressure.

---

## Best Practices

1. Prefer `streamable-http` in production unless a client requires SSE.
2. Tune one dimension at a time (`timeout`, then `qps/burst`, then `ratelimit`).
3. Keep load tests representative of real tool usage.
4. Track p95/p99 latency, not only average latency.

---

## Related Documentation

- [Configuration Guide](/docs/configuration/)
- [Deployment Guide](/docs/deployment/)
- [Security Guide](/docs/security/)
- [Tools Reference](/docs/tools/)

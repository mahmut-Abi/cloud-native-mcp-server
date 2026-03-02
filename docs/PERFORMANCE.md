# Performance

This document describes performance tuning that matches the **current** Cloud Native MCP Server implementation.

For website docs:
- English: `website/content/en/docs/performance.md`
- Chinese: `website/content/zh/docs/performance.md`

---

## What You Can Tune Today

### 1. Server Timeouts

```yaml
server:
  readTimeoutSec: 30
  writeTimeoutSec: 0
  idleTimeoutSec: 60
```

Guidance:
- Keep `writeTimeoutSec: 0` for SSE-heavy traffic.
- Increase `readTimeoutSec` for slow clients or larger request bodies.
- Increase `idleTimeoutSec` for long-lived connections behind reverse proxies.

### 2. Kubernetes Client Throughput

```yaml
kubernetes:
  timeoutSec: 30
  qps: 100.0
  burst: 200
```

Guidance:
- Increase `qps`/`burst` for larger clusters.
- Increase `timeoutSec` for expensive list/query operations.

### 3. Rate Limiting

```yaml
ratelimit:
  enabled: true
  requests_per_second: 100
  burst: 200
```

Guidance:
- Enable in internet-facing or multi-tenant deployments.
- Start conservative, then raise based on real metrics.

### 4. Service Scope Reduction

```yaml
enableDisable:
  enabledServices: ["kubernetes", "prometheus", "grafana", "aggregate"]
```

Guidance:
- Enable only services actually used by clients.
- Disable unused tools/services to reduce startup and runtime overhead.

### 5. Audit Cost Control (When Audit Is Enabled)

```yaml
audit:
  enabled: true
  storage: "database"
  sampling:
    enabled: true
    rate: 0.3
```

Guidance:
- Use sampling for high-QPS deployments.
- Use database storage for large audit datasets and frequent queries.

---

## Built-In Optimizations (No Public YAML Key)

The server includes internal optimizations such as:
- response truncation safeguards
- optimized JSON processing paths
- internal pooling/caching in service and tool layers

These are implementation details, not stable public config keys.
Tune supported keys first: `server`, `kubernetes`, `ratelimit`, `enableDisable`, `audit`.

---

## Metrics to Watch

Track:
- request throughput
- p95/p99 latency
- error rate
- active connections
- memory and CPU usage

Example Prometheus queries:

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

Create payload file (`payload.json`):

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
1. Increase `kubernetes.timeoutSec` for slow backend queries.
2. Increase `kubernetes.qps`/`kubernetes.burst` if client-side throttling appears.
3. Correlate p99 latency with backend service saturation.

### High Error Rate

Actions:
1. Verify auth mode and credentials.
2. Check backend service availability.
3. Inspect `/api/audit/logs` when audit is enabled.

### High Memory Usage

Actions:
1. Reduce enabled services/tools.
2. Enable audit sampling.
3. Lower burst values if traffic spikes create memory pressure.

---

## Best Practices

1. Prefer `streamable-http` in production unless a client requires SSE.
2. Tune one dimension at a time (`timeout` -> `qps/burst` -> `ratelimit`).
3. Keep load tests representative of real tool usage.
4. Track p95/p99, not only average latency.

---

## Related Docs

- `docs/CONFIGURATION.md`
- `docs/DEPLOYMENT.md`
- `website/content/en/docs/performance.md`
- `website/content/zh/docs/performance.md`

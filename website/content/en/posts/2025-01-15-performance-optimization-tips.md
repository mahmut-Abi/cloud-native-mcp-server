---
title: "Performance Optimization Tips for Cloud Native MCP Server"
date: 2025-01-15T10:00:00Z
description: "Practical tuning tips for Cloud Native MCP Server: timeouts, rate limits, metrics, and production performance checks."
tags: ["performance", "optimization", "tutorials"]
---

Learn how to optimize Cloud Native MCP Server for predictable latency and higher throughput in real production workloads.

## Cache and Response Strategy

The server already includes internal cache and response shaping mechanisms. You can improve performance further by reducing response scope per call:

- Prefer namespace-scoped queries over cluster-wide queries.
- Use pagination parameters (for tools that support them) on large datasets.
- Query only the fields you actually need for the current decision.

### Example: Limit payload size

```json
{
  "method": "kubernetes-get-pods",
  "params": {
    "namespace": "default",
    "limit": 50
  }
}
```

## Tune Kubernetes and Service Timeouts

Use runtime variables that are supported by the current server:

```bash
# Kubernetes client tuning
export MCP_K8S_TIMEOUT=30
export MCP_K8S_QPS=100
export MCP_K8S_BURST=200

# Upstream service request timeout (example: Prometheus)
export MCP_PROM_TIMEOUT=30
```

These settings should match your cluster size and backend responsiveness.

## Control Request Pressure

For busy environments, apply built-in rate limiting:

```bash
export MCP_RATELIMIT_ENABLED=true
export MCP_RATELIMIT_REQUESTS_PER_SECOND=25
export MCP_RATELIMIT_BURST=80
```

This helps prevent overload during traffic spikes and protects upstream services.

## Resource Planning

### Memory

- Small environments: 512MB - 1GB
- Medium environments: 1GB - 2GB
- Large/high-concurrency environments: 2GB+

### CPU

Cloud Native MCP Server uses optimized encoding and transport paths. In CPU-constrained environments:

- reduce burst rate
- reduce query fan-out
- disable non-required services

```bash
export MCP_DISABLED_SERVICES="kibana,jaeger"
```

## Monitor Performance with `/metrics`

```bash
curl -sS http://localhost:8080/metrics
```

Useful metrics include:

- `http_request_duration_seconds`
- `http_requests_total`
- `tool_call_duration_seconds`
- `tool_calls_total`
- `cache_hits_total`
- `cache_misses_total`

## Practical Checklist

1. Keep requests narrow and paginated.
2. Tune `MCP_K8S_QPS` / `MCP_K8S_BURST` for your cluster profile.
3. Set realistic upstream timeouts.
4. Enable rate limiting in production.
5. Watch metrics continuously and iterate.

Need deeper guidance? Read the [Performance Guide](/docs/performance/).

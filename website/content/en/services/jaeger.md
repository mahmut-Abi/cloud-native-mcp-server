---
title: "Jaeger Service"
weight: 8
---

# Jaeger Service

The Jaeger service provides comprehensive distributed tracing and performance analysis with 8 tools for managing tracing resources.

## Overview

The Jaeger service in Cloud Native MCP Server enables AI assistants to analyze distributed traces, service dependencies, and performance metrics efficiently. It provides tools for trace querying, dependency analysis, and performance monitoring.

### Key Capabilities

{{< columns >}}
### 📍 Trace Analysis
Query and analyze distributed traces across microservices.
<--->

### 🔗 Dependency Mapping
Visualize service dependencies and call graphs.
{{< /columns >}}

{{< columns >}}
### ⚡ Performance Monitoring
Analyze performance bottlenecks and latency patterns.
<--->

### 📊 Metrics Collection
Collect and analyze tracing metrics and statistics.
{{< /columns >}}

---

## Available Tools (8)

### Trace Management
- **jaeger-get-traces**: Get traces by query
- **jaeger-get-trace**: Get a specific trace
- **jaeger-search-traces**: Search traces with filters

### Service and Operation Analysis
- **jaeger-get-services**: Get all services
- **jaeger-get-service-operations**: Get operations for a service
- **jaeger-get-operations**: Get all operations

### Dependency and Metrics
- **jaeger-get-dependencies**: Get service dependencies
- **jaeger-get-metrics**: Get tracing metrics

---

## Quick Examples

### Get traces for a specific service

```json
{
  "method": "tools/call",
  "params": {
    "name": "jaeger-get-traces",
    "arguments": {
      "service": "my-app",
      "limit": 100
    }
  }
}
```

### Get service dependencies

```json
{
  "method": "tools/call",
  "params": {
    "name": "jaeger-get-dependencies",
    "arguments": {
      "service": "my-app",
      "start": "1 hour ago"
    }
  }
}
```

### Get all services

```json
{
  "method": "tools/call",
  "params": {
    "name": "jaeger-get-services",
    "arguments": {}
  }
}
```

---

## Best Practices

- Implement proper tracing headers propagation
- Use appropriate sampling strategies for performance
- Regularly analyze traces for performance bottlenecks
- Monitor service dependencies for architecture changes
- Set up alerts based on tracing metrics and anomalies

## Next Steps

- [OpenTelemetry Service](/en/services/opentelemetry/) for metrics and logs
- [Tracing Guides](/en/guides/tracing/) for detailed setup
- [Performance Analysis](/en/guides/performance/) for optimization
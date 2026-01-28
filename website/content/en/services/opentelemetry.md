---
title: "OpenTelemetry Service"
weight: 9
---

# OpenTelemetry Service

The OpenTelemetry service provides comprehensive metrics, traces, and logs collection and analysis with 9 tools for managing observability resources.

## Overview

The OpenTelemetry service in Cloud Native MCP Server enables AI assistants to collect and analyze telemetry data from applications and infrastructure efficiently. It provides tools for metrics collection, trace analysis, and log aggregation.

### Key Capabilities

{{< columns >}}
### 📊 Metrics Collection
Collect and analyze application and infrastructure metrics.
<--->

### 📍 Trace Analysis
Distributed tracing with performance insights.
{{< /columns >}}

{{< columns >}}
### 📝 Log Aggregation
Unified log collection and analysis.
<--->

### 🛠️ Configuration
Manage OpenTelemetry collector configuration and status.
{{< /columns >}}

---

## Available Tools (9)

### Metrics Management
- **otel-get-metrics**: Get metrics from OpenTelemetry collector
- **otel-get-metric-data**: Get metric data
- **otel-list-metric-streams**: List metric streams
- **otel-get-metrics-schema**: Get metrics schema

### Trace Management
- **otel-get-traces**: Get traces from OpenTelemetry collector
- **otel-search-traces**: Search traces
- **otel-get-traces-schema**: Get traces schema

### Log and Configuration Management
- **otel-get-logs**: Get logs from OpenTelemetry collector
- **otel-get-logs-schema**: Get logs schema

### System Information
- **otel-get-status**: Get OpenTelemetry collector status
- **otel-get-config**: Get OpenTelemetry collector configuration
- **otel-get-health**: Get OpenTelemetry collector health
- **otel-get-versions**: Get OpenTelemetry component versions

---

## Quick Examples

### Get metrics from the collector

```json
{
  "method": "tools/call",
  "params": {
    "name": "otel-get-metrics",
    "arguments": {
      "metric_name": "http_requests_total",
      "start_time": "1 hour ago",
      "end_time": "now"
    }
  }
}
```

### Get traces for a specific service

```json
{
  "method": "tools/call",
  "params": {
    "name": "otel-get-traces",
    "arguments": {
      "service_name": "my-app",
      "limit": 50
    }
  }
}
```

### Get collector configuration

```json
{
  "method": "tools/call",
  "params": {
    "name": "otel-get-config",
    "arguments": {}
  }
}
```

---

## Best Practices

- Implement proper resource attributes for effective filtering
- Use appropriate sampling strategies for traces
- Configure appropriate metric collection intervals
- Monitor collector health and resource usage
- Set up alerts based on telemetry data patterns

## Next Steps

- [Jaeger Service](/en/services/jaeger/) for detailed tracing
- [Observability Guides](/en/guides/observability/) for detailed setup
- [Metrics Best Practices](/en/guides/metrics/) for collection strategies
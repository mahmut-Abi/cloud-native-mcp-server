---
title: "OpenTelemetry Service"
weight: 9
---

# OpenTelemetry Service

The OpenTelemetry service provides collector diagnostics, configuration analysis, and telemetry inspection with 12 tools for managing observability resources.

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

## Available Tools (12)

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

### Collector Diagnostics
- **opentelemetry_get_health**: Get OpenTelemetry collector health
- **opentelemetry_get_status**: Get OpenTelemetry collector status
- **opentelemetry_get_collector_summary**: Get a compact collector health and config overview
- **opentelemetry_get_config**: Get full OpenTelemetry collector configuration
- **opentelemetry_get_config_summary**: Get a compact config summary
- **opentelemetry_analyze_pipeline_status**: Analyze pipelines for missing components and common misconfigurations

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

- [Jaeger Service](/services/jaeger/) for detailed tracing
- [Observability Guides](/services/opentelemetry/) for detailed setup
- [Metrics Best Practices](/services/prometheus/) for collection strategies

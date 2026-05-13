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
- `opentelemetry_get_metrics`: Get metrics from the collector
- `opentelemetry_query_metrics`: Run a PromQL-style query against collector metrics

### Trace Management
- `opentelemetry_get_traces`: Get traces from the collector
- `opentelemetry_query_traces`: Search traces with filters and time range

### Log Management
- `opentelemetry_get_logs`: Get logs from the collector
- `opentelemetry_query_logs`: Search logs with filters and time range

### Collector Diagnostics
- `opentelemetry_get_health`: Get OpenTelemetry collector health
- `opentelemetry_get_status`: Get detailed collector status
- `opentelemetry_get_collector_summary`: Get a compact collector health and config overview
- `opentelemetry_get_config`: Get full collector configuration
- `opentelemetry_get_config_summary`: Get a compact config summary
- `opentelemetry_analyze_pipeline_status`: Analyze pipelines for missing components and common misconfigurations

---

## Quick Examples

### Get metrics from the collector

```json
{
  "method": "tools/call",
  "params": {
    "name": "opentelemetry_get_metrics",
    "arguments": {
      "metric_name": "http_requests_total"
    }
  }
}
```

### Analyze one collector's pipelines

```json
{
  "method": "tools/call",
  "params": {
    "name": "opentelemetry_analyze_pipeline_status",
    "arguments": {
      "signal": "traces"
    }
  }
}
```

### Get collector configuration summary

```json
{
  "method": "tools/call",
  "params": {
    "name": "opentelemetry_get_config_summary",
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

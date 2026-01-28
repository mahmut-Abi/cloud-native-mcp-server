---
title: "Prometheus Service"
weight: 4
---

# Prometheus Service

The Prometheus service provides comprehensive metrics collection, querying, and monitoring capabilities with 20 tools for managing Prometheus resources.

## Overview

The Prometheus service in Cloud Native MCP Server enables AI assistants to query and manage Prometheus metrics efficiently. It provides tools for metric querying, alert management, and configuration management.

### Key Capabilities

{{< columns >}}
### 📈 Metrics Querying
Powerful querying capabilities for Prometheus metrics using PromQL.
<--->

### ⚠️ Alert Management
Manage Prometheus alerts and alerting rules effectively.
{{< /columns >}}

{{< columns >}}
### 🛠️ Configuration
Handle Prometheus configuration and runtime information.
<--->

### 📊 Monitoring
Access detailed monitoring data and statistics from Prometheus.
{{< /columns >}}

---

## Available Tools (20)

### Query Execution
- **prometheus-query**: Execute instant query
- **prometheus-query-range**: Execute range query
- **prometheus-query-exemplars**: Query exemplar data

### Metadata Queries
- **prometheus-label-names**: Get label names
- **prometheus-label-values**: Get label values
- **prometheus-series**: Get time series
- **prometheus-metadata**: Get metadata

### Target Management
- **prometheus-get-targets**: Get target list
- **prometheus-get-target-metadata**: Get target metadata

### Rules and Alerts Management
- **prometheus-get-rules**: Get rules list
- **prometheus-get-alerts**: Get alerts list
- **prometheus-get-alert-managers**: Get Alertmanager instances

### Configuration Management
- **prometheus-get-config**: Get configuration information
- **prometheus-get-flags**: Get startup flags

### Status Queries
- **prometheus-get-status**: Get status information
- **prometheus-get-build-info**: Get build information
- **prometheus-get-runtime-info**: Get runtime information

### TSDB Operations
- **prometheus-get-tsdb-status**: Get TSDB status
- **prometheus-get-tsdb-heatmap**: Get TSDB heatmap

---

## Quick Examples

### Query metrics for the last 5 minutes

```json
{
  "method": "tools/call",
  "params": {
    "name": "prometheus-query-range",
    "arguments": {
      "query": "up",
      "start": "5 minutes ago",
      "end": "now",
      "step": "30s"
    }
  }
}
```

### Get all targets

```json
{
  "method": "tools/call",
  "params": {
    "name": "prometheus-get-targets",
    "arguments": {}
  }
}
```

### Query specific metric

```json
{
  "method": "tools/call",
  "params": {
    "name": "prometheus-query",
    "arguments": {
      "query": "rate(http_requests_total[5m])"
    }
  }
}
```

---

## Best Practices

- Use appropriate query ranges to avoid performance issues
- Regularly review and optimize PromQL queries
- Monitor target health and availability
- Configure appropriate recording rules for pre-aggregated metrics
- Set up proper alerting rules with appropriate thresholds

## Next Steps

- [Grafana Service](/en/services/grafana/) for visualization
- [Monitoring Guides](/en/guides/monitoring/) for detailed setup
- [Performance Optimization](/en/guides/performance/) for query optimization
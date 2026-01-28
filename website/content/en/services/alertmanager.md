---
title: "Alertmanager Service"
weight: 7
---

# Alertmanager Service

The Alertmanager service provides comprehensive alert rules management and notifications with 15 tools for managing alerting resources.

## Overview

The Alertmanager service in Cloud Native MCP Server enables AI assistants to manage Prometheus alerts, silences, and notification routing efficiently. It provides tools for alert management, silencing, and configuration.

### Key Capabilities

{{< columns >}}
### ⚠️ Alert Management
Complete control over Prometheus alerts and alert groups.
<--->

### 🔕 Silence Management
Manage alert silences with precise timing and matching rules.
{{< /columns >}}

{{< columns >}}
### 📮 Notification Routing
Configure notification routes and receivers for alert delivery.
<--->

### ⚙️ Configuration
Manage Alertmanager configuration and status information.
{{< /columns >}}

---

## Available Tools (15)

### Alert Management
- **alertmanager-get-alerts**: Get all alerts
- **alertmanager-get-alert**: Get a specific alert
- **alertmanager-get-alert-groups**: Get alert groups
- **alertmanager-get-receivers**: Get all receivers

### Silence Management
- **alertmanager-get-silences**: Get all silences
- **alertmanager-create-silence**: Create a new silence
- **alertmanager-delete-silence**: Delete a silence
- **alertmanager-get-alertmanagers**: Get Alertmanager instances

### Configuration and Status
- **alertmanager-get-config**: Get Alertmanager configuration
- **alertmanager-get-status**: Get Alertmanager status
- **alertmanager-get-metrics**: Get Alertmanager metrics
- **alertmanager-get-templates**: Get notification templates
- **alertmanager-get-starttime**: Get Alertmanager start time
- **alertmanager-get-version**: Get Alertmanager version
- **alertmanager-get-flags**: Get Alertmanager flag values

---

## Quick Examples

### Create a silence for high CPU alerts

```json
{
  "method": "tools/call",
  "params": {
    "name": "alertmanager-create-silence",
    "arguments": {
      "matcher": {
        "name": "alertname",
        "value": "HighCPUUsage",
        "isRegex": false
      },
      "startsAt": "2023-10-01T10:00:00Z",
      "endsAt": "2023-10-01T12:00:00Z",
      "createdBy": "admin",
      "comment": "Planned maintenance window"
    }
  }
}
```

### Get all active alerts

```json
{
  "method": "tools/call",
  "params": {
    "name": "alertmanager-get-alerts",
    "arguments": {}
  }
}
```

### Get Alertmanager configuration

```json
{
  "method": "tools/call",
  "params": {
    "name": "alertmanager-get-config",
    "arguments": {}
  }
}
```

---

## Best Practices

- Use precise matchers for effective silencing
- Regularly review active alerts and silences
- Configure appropriate notification routes and grouping
- Monitor Alertmanager performance and health
- Implement proper escalation policies for critical alerts

## Next Steps

- [Prometheus Service](/en/services/prometheus/) for metrics
- [Alerting Guides](/en/guides/alerting/) for detailed setup
- [Notification Configuration](/en/guides/notifications/) for routing rules
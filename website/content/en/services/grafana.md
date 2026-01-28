---
title: "Grafana Service"
weight: 3
---

# Grafana Service

The Grafana service provides comprehensive visualization, monitoring dashboards, and alerting capabilities with 36 tools for creating and managing Grafana resources.

## Overview

The Grafana service in Cloud Native MCP Server enables AI assistants to manage Grafana dashboards, data sources, alerts, and other monitoring resources efficiently. It provides tools for dashboard creation, visualization management, and alert configuration.

### Key Capabilities

{{< columns >}}
### 📊 Dashboard Management
Complete control over Grafana dashboards including creation, updates, and sharing.
<--->

### 🗂️ Data Source Management
Manage Grafana data sources with tools for configuration and testing.
{{< /columns >}}

{{< columns >}}
### ⚠️ Alert Management
Handle Grafana alerts and alert rules with configuration and monitoring tools.
<--->

### 📈 Visualization
Create and manage Grafana visualizations and panels effectively.
{{< /columns >}}

---

## Available Tools (36)

### Dashboard Management
- **grafana-get-dashboards**: Get all dashboards
- **grafana-get-dashboard**: Get a specific dashboard
- **grafana-create-dashboard**: Create a new dashboard
- **grafana-update-dashboard**: Update an existing dashboard
- **grafana-delete-dashboard**: Delete a dashboard
- **grafana-get-folders**: Get all folders
- **grafana-create-folder**: Create a new folder
- **grafana-update-folder**: Update an existing folder
- **grafana-delete-folder**: Delete a folder

### Data Source Management
- **grafana-get-datasources**: Get all data sources
- **grafana-create-datasource**: Create a new data source
- **grafana-update-datasource**: Update an existing data source
- **grafana-delete-datasource**: Delete a data source
- **grafana-test-datasource**: Test a data source connection

### Alert Management
- **grafana-get-alerts**: Get all alerts
- **grafana-get-alert**: Get a specific alert
- **grafana-create-alert**: Create a new alert
- **grafana-update-alert**: Update an existing alert
- **grafana-delete-alert**: Delete an alert
- **grafana-get-alert-rules**: Get alert rules
- **grafana-get-alert-notifications**: Get alert notification channels

### User and Organization Management
- **grafana-get-users**: Get all users
- **grafana-create-user**: Create a new user
- **grafana-update-user**: Update an existing user
- **grafana-delete-user**: Delete a user
- **grafana-get-orgs**: Get all organizations
- **grafana-create-org**: Create a new organization
- **grafana-update-org**: Update an existing organization
- **grafana-delete-org**: Delete an organization
- **grafana-get-teams**: Get all teams
- **grafana-create-team**: Create a new team
- **grafana-update-team**: Update an existing team
- **grafana-delete-team**: Delete a team

### Plugin and Configuration Management
- **grafana-get-plugins**: Get all plugins
- **grafana-install-plugin**: Install a plugin
- **grafana-uninstall-plugin**: Uninstall a plugin
- **grafana-get-annotations**: Get annotations
- **grafana-create-annotation**: Create an annotation
- **grafana-get-snapshots**: Get snapshots

---

## Quick Examples

### Create a new dashboard

```json
{
  "method": "tools/call",
  "params": {
    "name": "grafana-create-dashboard",
    "arguments": {
      "dashboard": {
        "title": "My Application Dashboard",
        "panels": [
          {
            "id": 1,
            "title": "Requests per Second",
            "type": "graph",
            "targets": [
              {
                "expr": "rate(http_requests_total[5m])"
              }
            ]
          }
        ]
      }
    }
  }
}
```

### Get all dashboards

```json
{
  "method": "tools/call",
  "params": {
    "name": "grafana-get-dashboards",
    "arguments": {}
  }
}
```

### Add a data source

```json
{
  "method": "tools/call",
  "params": {
    "name": "grafana-create-datasource",
    "arguments": {
      "name": "Prometheus",
      "type": "prometheus",
      "url": "http://prometheus:9090",
      "access": "proxy"
    }
  }
}
```

---

## Best Practices

- Organize dashboards in folders by application or team
- Use consistent naming conventions for dashboards and panels
- Configure appropriate alert thresholds and notification channels
- Regularly review and update dashboards based on changing requirements
- Implement proper user permissions and access controls

## Next Steps

- [Prometheus Service](/en/services/prometheus/) for metrics collection
- [Monitoring Guides](/en/guides/monitoring/) for detailed setup
- [Alerting Best Practices](/en/guides/alerting/) for effective alerting
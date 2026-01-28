---
title: "Kibana Service"
weight: 5
---

# Kibana Service

The Kibana service provides comprehensive log analysis, visualization, and data exploration capabilities with 52 tools for managing Kibana resources.

## Overview

The Kibana service in Cloud Native MCP Server enables AI assistants to manage Kibana dashboards, visualizations, and log data efficiently. It provides tools for log analysis, data visualization, and Elasticsearch integration.

### Key Capabilities

{{< columns >}}
### 🔍 Log Analysis
Powerful log analysis and search capabilities with Elasticsearch integration.
<--->

### 📊 Data Visualization
Create and manage Kibana dashboards and visualizations.
{{< /columns >}}

{{< columns >}}
### 🗂️ Object Management
Handle Kibana saved objects including searches, visualizations, and dashboards.
<--->

### 🌐 Elasticsearch Integration
Seamless integration with Elasticsearch for data indexing and retrieval.
{{< /columns >}}

---

## Available Tools (52)

### Index Management
- **kibana-get-indices**: Get all indices
- **kibana-create-index**: Create a new index
- **kibana-delete-index**: Delete an index
- **kibana-get-index-patterns**: Get all index patterns
- **kibana-create-index-pattern**: Create a new index pattern
- **kibana-delete-index-pattern**: Delete an index pattern

### Visualization Management
- **kibana-get-visualizations**: Get all visualizations
- **kibana-create-visualization**: Create a new visualization
- **kibana-update-visualization**: Update a visualization
- **kibana-delete-visualization**: Delete a visualization

### Dashboard Management
- **kibana-get-dashboards**: Get all dashboards
- **kibana-create-dashboard**: Create a new dashboard
- **kibana-update-dashboard**: Update a dashboard
- **kibana-delete-dashboard**: Delete a dashboard

### Search Management
- **kibana-get-searches**: Get all saved searches
- **kibana-create-search**: Create a new saved search
- **kibana-update-search**: Update a saved search
- **kibana-delete-search**: Delete a saved search

### Object Management
- **kibana-get-objects**: Get all saved objects
- **kibana-create-object**: Create a new saved object
- **kibana-update-object**: Update a saved object
- **kibana-delete-object**: Delete a saved object

### User Management
- **kibana-get-users**: Get all users
- **kibana-create-user**: Create a new user
- **kibana-update-user**: Update a user
- **kibana-delete-user**: Delete a user

### Role Management
- **kibana-get-roles**: Get all roles
- **kibana-create-role**: Create a new role
- **kibana-update-role**: Update a role
- **kibana-delete-role**: Delete a role

### Space Management
- **kibana-get-spaces**: Get all spaces
- **kibana-create-space**: Create a new space
- **kibana-update-space**: Update a space
- **kibana-delete-space**: Delete a space

### API Key Management
- **kibana-get-api-keys**: Get all API keys
- **kibana-create-api-key**: Create a new API key
- **kibana-delete-api-key**: Delete an API key

### System Information
- **kibana-get-status**: Get Kibana status
- **kibana-get-info**: Get Kibana information
- **kibana-get-config**: Get Kibana configuration
- **kibana-get-plugins**: Get installed plugins
- **kibana-get-license**: Get license information
- **kibana-get-features**: Get feature flags
- **kibana-get-telemetry**: Get telemetry settings
- **kibana-update-telemetry**: Update telemetry settings

### Import/Export
- **kibana-get-saved-objects**: Get saved objects by type
- **kibana-export-objects**: Export saved objects
- **kibana-import-objects**: Import saved objects
- **kibana-get-maps**: Get Maps saved objects
- **kibana-get-cases**: Get Cases saved objects

---

## Quick Examples

### Create a new index pattern

```json
{
  "method": "tools/call",
  "params": {
    "name": "kibana-create-index-pattern",
    "arguments": {
      "id": "my-app-logs-*",
      "title": "my-app-logs-*",
      "timeFieldName": "@timestamp"
    }
  }
}
```

### Search documents

```json
{
  "method": "tools/call",
    "params": {
      "name": "kibana-get-searches",
      "arguments": {}
    }
  }
}
```

### Get Kibana status

```json
{
  "method": "tools/call",
  "params": {
    "name": "kibana-get-status",
    "arguments": {}
  }
}
```

---

## Best Practices

- Use appropriate index patterns for efficient log searching
- Regularly review and optimize visualizations for performance
- Implement proper access controls using roles and spaces
- Monitor cluster health and resource usage
- Use saved searches for frequently accessed data patterns

## Next Steps

- [Elasticsearch Service](/en/services/elasticsearch/) for indexing
- [Log Analysis Guides](/en/guides/logging/) for detailed setup
- [Visualization Best Practices](/en/guides/visualization/) for effective dashboards
---
title: "Elasticsearch Service"
weight: 6
---

# Elasticsearch Service

The Elasticsearch service provides comprehensive log storage, search, and data indexing capabilities with 14 tools for managing Elasticsearch resources.

## Overview

The Elasticsearch service in Cloud Native MCP Server enables AI assistants to manage Elasticsearch indices, documents, and clusters efficiently. It provides tools for indexing, searching, and cluster management.

### Key Capabilities

{{< columns >}}
### 🔍 Advanced Search
Powerful full-text search and analytics capabilities.
<--->

### 🗂️ Index Management
Complete control over Elasticsearch indices and mappings.
{{< /columns >}}

{{< columns >}}
### 📦 Document Management
Handle document indexing, retrieval, and bulk operations.
<--->

### 🖥️ Cluster Management
Monitor and manage Elasticsearch cluster health and performance.
{{< /columns >}}

---

## Available Tools (14)

### Index Management
- **elasticsearch-get-indices**: Get all indices
- **elasticsearch-create-index**: Create a new index
- **elasticsearch-delete-index**: Delete an index
- **elasticsearch-get-index-settings**: Get index settings
- **elasticsearch-update-index-settings**: Update index settings
- **elasticsearch-get-mappings**: Get index mappings
- **elasticsearch-update-mappings**: Update index mappings

### Document Operations
- **elasticsearch-index-document**: Index a document
- **elasticsearch-get-document**: Get a document
- **elasticsearch-search**: Search documents
- **elasticsearch-delete-document**: Delete a document

### Cluster Management
- **elasticsearch-get-cluster-info**: Get cluster information
- **elasticsearch-get-cluster-health**: Get cluster health
- **elasticsearch-get-nodes**: Get cluster nodes

---

## Quick Examples

### Index a document

```json
{
  "method": "tools/call",
  "params": {
    "name": "elasticsearch-index-document",
    "arguments": {
      "index": "my-app-logs",
      "id": "1",
      "document": {
        "timestamp": "2023-10-01T12:00:00Z",
        "level": "info",
        "message": "Application started successfully"
      }
    }
  }
}
```

### Search documents

```json
{
  "method": "tools/call",
  "params": {
    "name": "elasticsearch-search",
    "arguments": {
      "index": "my-app-logs",
      "query": {
        "match": {
          "level": "error"
        }
      }
    }
  }
}
```

### Get cluster health

```json
{
  "method": "tools/call",
  "params": {
    "name": "elasticsearch-get-cluster-health",
    "arguments": {}
  }
}
```

---

## Best Practices

- Use appropriate index mappings for efficient searching
- Regularly optimize indices for performance
- Monitor cluster health and resource usage
- Implement proper index lifecycle management
- Use bulk operations for efficient data ingestion

## Next Steps

- [Kibana Service](/en/services/kibana/) for visualization
- [Log Management Guides](/en/guides/logging/) for detailed setup
- [Search Optimization](/en/guides/search/) for query performance
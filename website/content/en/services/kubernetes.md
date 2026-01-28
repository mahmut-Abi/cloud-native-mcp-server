---
title: "Kubernetes Service"
weight: 1
---

# Kubernetes Service

The Kubernetes service provides comprehensive container orchestration and resource management capabilities with 28 specialized tools for managing your Kubernetes clusters.

## Overview

The Kubernetes service in Cloud Native MCP Server enables AI assistants to manage Kubernetes resources efficiently. It provides tools for deployments, services, configmaps, secrets, and other core Kubernetes resources.

### Key Capabilities

{{< columns >}}
### 🔧 Deployment Management
Complete control over Kubernetes deployments including creation, updates, scaling, and deletion operations.
<--->

### 🗂️ Resource Management
Manage all Kubernetes resources including pods, services, configmaps, secrets, and persistent volumes.
{{< /columns >}}

{{< columns >}}
### 📊 Monitoring
Get detailed information about pods, nodes, and resource usage across your cluster.
<--->

### 🔐 Security
Manage secrets, RBAC configurations, and other security-related Kubernetes resources.
{{< /columns >}}

---

## Available Tools (28)

### Pod Management
- **kubernetes-get-pods**: Get detailed information about pods in a namespace
- **kubernetes-list-pods**: List all pods in a namespace
- **kubernetes-get-pod**: Get specific pod details
- **kubernetes-delete-pod**: Delete a specific pod
- **kubernetes-get-pod-logs**: Get logs from a pod
- **kubernetes-get-pod-events**: Get events related to a pod

### Deployment Management
- **kubernetes-list-deployments**: List all deployments in a namespace
- **kubernetes-get-deployment**: Get specific deployment details
- **kubernetes-create-deployment**: Create a new deployment
- **kubernetes-update-deployment**: Update an existing deployment
- **kubernetes-delete-deployment**: Delete a deployment
- **kubernetes-scale-deployment**: Scale a deployment
- **kubernetes-restart-deployment**: Restart a deployment

### Service Management
- **kubernetes-list-services**: List all services in a namespace
- **kubernetes-get-service**: Get specific service details
- **kubernetes-create-service**: Create a new service
- **kubernetes-update-service**: Update an existing service
- **kubernetes-delete-service**: Delete a service

### Configuration Management
- **kubernetes-list-configmaps**: List all configmaps in a namespace
- **kubernetes-get-configmap**: Get specific configmap details
- **kubernetes-create-configmap**: Create a new configmap
- **kubernetes-update-configmap**: Update an existing configmap
- **kubernetes-delete-configmap**: Delete a configmap
- **kubernetes-list-secrets**: List all secrets in a namespace
- **kubernetes-get-secret**: Get specific secret details
- **kubernetes-create-secret**: Create a new secret
- **kubernetes-update-secret**: Update an existing secret
- **kubernetes-delete-secret**: Delete a secret

---

## Quick Examples

### List all pods in default namespace

```json
{
  "method": "tools/call",
  "params": {
    "name": "kubernetes-list-pods",
    "arguments": {
      "namespace": "default"
    }
  }
}
```

### Get specific deployment details

```json
{
  "method": "tools/call",
  "params": {
    "name": "kubernetes-get-deployment",
    "arguments": {
      "name": "my-app",
      "namespace": "production"
    }
  }
}
```

### Create a new configmap

```json
{
  "method": "tools/call",
  "params": {
    "name": "kubernetes-create-configmap",
    "arguments": {
      "name": "app-config",
      "namespace": "default",
      "data": {
        "config.json": "{\"debug\": true, \"port\": 8080}"
      }
    }
  }
}
```

---

## Best Practices

- Always specify namespaces when working with Kubernetes resources
- Use labels and annotations effectively for resource organization
- Implement proper RBAC policies for security
- Monitor resource usage to optimize cluster performance
- Regularly backup critical configurations

## Next Steps

- [Helm Service](/en/services/helm/) for package management
- [Configuration Guides](/en/guides/configuration/) for detailed setup
- [Security Best Practices](/en/guides/security/) for securing your cluster
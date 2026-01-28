---
title: "Helm Service"
weight: 2
---

# Helm Service

The Helm service provides comprehensive package management and deployment capabilities with 31 tools for managing Helm charts, releases, and repositories.

## Overview

The Helm service in Cloud Native MCP Server enables AI assistants to manage Helm charts and releases efficiently. It provides tools for chart installation, upgrades, rollbacks, and repository management.

### Key Capabilities

{{< columns >}}
### 📦 Chart Management
Complete control over Helm charts including installation, upgrading, and uninstalling releases.
<--->

### 🗄️ Repository Management
Manage Helm chart repositories with tools for adding, updating, and searching charts.
{{< /columns >}}

{{< columns >}}
### 🔄 Release Management
Handle Helm releases with rollback, history, and status checking capabilities.
<--->

### ⚙️ Configuration
Manage chart values, configurations, and dependencies effectively.
{{< /columns >}}

---

## Available Tools (31)

### Chart Management
- **helm-list-releases**: List all releases in all namespaces
- **helm-install-chart**: Install a chart
- **helm-upgrade-release**: Upgrade a release
- **helm-uninstall-release**: Uninstall a release
- **helm-get-release**: Get information about a release
- **helm-rollback-release**: Rollback a release
- **helm-get-history**: Get release history
- **helm-search-repo**: Search for charts in the repository
- **helm-add-repo**: Add a chart repository
- **helm-update-repo**: Update chart repositories
- **helm-repo-list**: List chart repositories
- **helm-get-values**: Get values for a release
- **helm-template**: Template a chart locally
- **helm-package**: Package a chart directory into a chart archive
- **helm-pull**: Download a chart from a repository
- **helm-push**: Push a chart to a registry

### Chart Information
- **helm-get-chart**: Get information about a chart
- **helm-create**: Create a new chart
- **helm-dependency-build**: Build out the charts dependencies
- **helm-dependency-update**: Update the charts dependencies
- **helm-lint**: Examine a chart for possible issues
- **helm-test**: Run tests for a release
- **helm-status**: Show the status of a release
- **helm-history**: Show the history of a release
- **helm-get-manifest**: Show the manifest for a release
- **helm-get-notes**: Show the notes for a release
- **helm-get-hooks**: Show the hooks for a release
- **helm-get-all**: Get all resources for a release
- **helm-verify**: Verify the provenance of a chart
- **helm-show-chart**: Show information about a chart
- **helm-show-readme**: Show the README of a chart

---

## Quick Examples

### Install a chart

```json
{
  "method": "tools/call",
  "params": {
    "name": "helm-install-chart",
    "arguments": {
      "chart": "nginx-ingress",
      "repo": "https://kubernetes.github.io/ingress-nginx",
      "release": "my-nginx",
      "namespace": "ingress-nginx"
    }
  }
}
```

### Upgrade a release

```json
{
  "method": "tools/call",
  "params": {
    "name": "helm-upgrade-release",
    "arguments": {
      "release": "my-nginx",
      "chart": "nginx-ingress",
      "repo": "https://kubernetes.github.io/ingress-nginx",
      "set_values": {
        "controller.replicaCount": 3
      }
    }
  }
}
```

### List all releases

```json
{
  "method": "tools/call",
  "params": {
    "name": "helm-list-releases",
    "arguments": {}
  }
}
```

---

## Best Practices

- Always specify namespaces when installing Helm charts
- Use values files for complex configurations
- Regularly update chart repositories
- Monitor release history for rollback capabilities
- Validate charts before installation using linting tools

## Next Steps

- [Kubernetes Service](/en/services/kubernetes/) for core orchestration
- [Configuration Guides](/en/guides/configuration/) for detailed setup
- [Deployment Best Practices](/en/guides/deployment/) for production deployments
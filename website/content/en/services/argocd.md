---
title: "Argo CD Service"
weight: 10
---

# Argo CD Service

The Argo CD service provides read-only GitOps inspection workflows with 7 tools for applications, manifests, projects, and clusters.

## Overview

Use the Argo CD service when you need to understand what Argo CD currently knows about an application's desired state, sync status, health, project boundaries, and connected clusters without leaving the MCP surface.

## Available Tools (7)

- `argocd_test_connection`
- `argocd_list_applications_summary`
- `argocd_get_application`
- `argocd_get_application_manifests`
- `argocd_list_projects`
- `argocd_get_project`
- `argocd_list_clusters`

## Recommended Workflow

1. Start with `argocd_test_connection`
2. Use `argocd_list_applications_summary` to discover candidate apps
3. Use `argocd_get_application` for sync and health detail
4. Use `argocd_get_application_manifests` when you need rendered resources

## Configuration Example

```yaml
argocd:
  enabled: true
  url: "https://argocd.example.com"
  username: "admin"
  password: ""
  authToken: ""
  timeoutSec: 30
```

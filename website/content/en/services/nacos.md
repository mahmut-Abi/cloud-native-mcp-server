---
title: "Nacos Service"
weight: 11
---

# Nacos Service

The Nacos service provides 9 read-only tools for namespace discovery, config inspection, service discovery, instance listing, and cluster-node visibility.

## Overview

Use the Nacos service when you need to inspect what is stored in Nacos as a configuration center or naming service: namespaces, config entries, services, registered instances, and server-side metrics.

## Available Tools (9)

- `nacos_test_connection`
- `nacos_list_namespaces`
- `nacos_list_configs_summary`
- `nacos_get_config`
- `nacos_list_services_summary`
- `nacos_get_service`
- `nacos_list_instances`
- `nacos_list_cluster_nodes`
- `nacos_get_system_metrics`

## Recommended Workflow

1. Start with `nacos_test_connection`
2. Use `nacos_list_namespaces` to discover namespace IDs
3. Use `nacos_list_configs_summary` or `nacos_list_services_summary`
4. Drill down with `nacos_get_config`, `nacos_get_service`, or `nacos_list_instances`

## Configuration Example

```yaml
nacos:
  enabled: true
  url: "http://localhost:8848/nacos"
  username: ""
  password: ""
  accessToken: ""
  namespaceId: ""
  group: "DEFAULT_GROUP"
  timeoutSec: 30
```

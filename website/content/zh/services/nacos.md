---
title: "Nacos 服务"
weight: 11
---

# Nacos 服务

Nacos 服务提供 9 个只读工具，用于命名空间发现、配置查看、服务发现、实例列表与集群节点检查。

## 概述

当你需要查看 Nacos 里存了哪些配置、服务注册信息、实例状态以及集群节点与指标时，可以直接使用这个服务，而不必切到独立控制台。

## 可用工具 (9)

- `nacos_test_connection`
- `nacos_list_namespaces`
- `nacos_list_configs_summary`
- `nacos_get_config`
- `nacos_list_services_summary`
- `nacos_get_service`
- `nacos_list_instances`
- `nacos_list_cluster_nodes`
- `nacos_get_system_metrics`

## 推荐调用顺序

1. 先用 `nacos_test_connection`
2. 用 `nacos_list_namespaces` 找命名空间 ID
3. 用 `nacos_list_configs_summary` 或 `nacos_list_services_summary` 做轻量发现
4. 再按需下钻到 `nacos_get_config`、`nacos_get_service`、`nacos_list_instances`

## 配置示例

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

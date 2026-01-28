---
title: "Kubernetes 工具"
weight: 10
---

# Kubernetes 工具

Cloud Native MCP Server 提供 28 个 Kubernetes 管理工具，用于 Pod、Deployment、Service 等资源的全面管理。

## Pod 管理

### list_pods

列出指定命名空间中的 Pod。

**参数**:
- `namespace` (string, required) - 命名空间

**返回**: Pod 列表

### get_pod

获取单个 Pod 的详细信息。

**参数**:
- `namespace` (string, required) - 命名空间
- `name` (string, required) - Pod 名称

**返回**: Pod 详细信息

### describe_pod

描述 Pod 的状态和事件。

**参数**:
- `namespace` (string, required) - 命名空间
- `name` (string, required) - Pod 名称

**返回**: Pod 状态描述

### delete_pod

删除指定的 Pod。

**参数**:
- `namespace` (string, required) - 命名空间
- `name` (string, required) - Pod 名称

**返回**: 删除结果

### get_pod_logs

获取 Pod 的日志。

**参数**:
- `namespace` (string, required) - 命名空间
- `name` (string, required) - Pod 名称
- `container` (string, optional) - 容器名称
- `tailLines` (int, optional) - 尾部行数
- `sinceSeconds` (int, optional) - 从多少秒前开始

**返回**: Pod 日志内容

### get_pod_events

获取 Pod 的事件。

**参数**:
- `namespace` (string, required) - 命名空间
- `name` (string, required) - Pod 名称

**返回**: Pod 事件列表

## Deployment 管理

### list_deployments

列出指定命名空间中的 Deployment。

**参数**:
- `namespace` (string, required) - 命名空间

**返回**: Deployment 列表

### get_deployment

获取单个 Deployment 的详细信息。

**参数**:
- `namespace` (string, required) - 命名空间
- `name` (string, required) - Deployment 名称

**返回**: Deployment 详细信息

### create_deployment

创建新的 Deployment。

**参数**:
- `namespace` (string, required) - 命名空间
- `manifest` (object, required) - Deployment 清单

**返回**: 创建的 Deployment 信息

### update_deployment

更新现有的 Deployment。

**参数**:
- `namespace` (string, required) - 命名空间
- `name` (string, required) - Deployment 名称
- `manifest` (object, required) - 更新的 Deployment 清单

**返回**: 更新后的 Deployment 信息

### delete_deployment

删除指定的 Deployment。

**参数**:
- `namespace` (string, required) - 命名空间
- `name` (string, required) - Deployment 名称

**返回**: 删除结果

### scale_deployment

扩缩容 Deployment。

**参数**:
- `namespace` (string, required) - 命名空间
- `name` (string, required) - Deployment 名称
- `replicas` (int, required) - 副本数

**返回**: 扩缩容结果

### restart_deployment

重启 Deployment（滚动重启）。

**参数**:
- `namespace` (string, required) - 命名空间
- `name` (string, required) - Deployment 名称

**返回**: 重启结果

## Service 管理

### list_services

列出指定命名空间中的 Service。

**参数**:
- `namespace` (string, required) - 命名空间

**返回**: Service 列表

### get_service

获取单个 Service 的详细信息。

**参数**:
- `namespace` (string, required) - 命名空间
- `name` (string, required) - Service 名称

**返回**: Service 详细信息

### create_service

创建新的 Service。

**参数**:
- `namespace` (string, required) - 命名空间
- `manifest` (object, required) - Service 清单

**返回**: 创建的 Service 信息

### delete_service

删除指定的 Service。

**参数**:
- `namespace` (string, required) - 命名空间
- `name` (string, required) - Service 名称

**返回**: 删除结果

## ConfigMap 管理

### list_configmaps

列出指定命名空间中的 ConfigMap。

**参数**:
- `namespace` (string, required) - 命名空间

**返回**: ConfigMap 列表

### get_configmap

获取单个 ConfigMap 的详细信息。

**参数**:
- `namespace` (string, required) - 命名空间
- `name` (string, required) - ConfigMap 名称

**返回**: ConfigMap 详细信息

### create_configmap

创建新的 ConfigMap。

**参数**:
- `namespace` (string, required) - 命名空间
- `name` (string, required) - ConfigMap 名称
- `data` (object, required) - ConfigMap 数据

**返回**: 创建的 ConfigMap 信息

## Secret 管理

### list_secrets

列出指定命名空间中的 Secret。

**参数**:
- `namespace` (string, required) - 命名空间

**返回**: Secret 列表

### get_secret

获取单个 Secret 的详细信息。

**参数**:
- `namespace` (string, required) - 命名空间
- `name` (string, required) - Secret 名称

**返回**: Secret 详细信息

### create_secret

创建新的 Secret。

**参数**:
- `namespace` (string, required) - 命名空间
- `name` (string, required) - Secret 名称
- `type` (string, required) - Secret 类型
- `data` (object, required) - Secret 数据（Base64 编码）

**返回**: 创建的 Secret 信息

## 命名空间管理

### list_namespaces

列出所有命名空间。

**返回**: 命名空间列表

### get_namespace

获取单个命名空间的详细信息。

**参数**:
- `name` (string, required) - 命名空间名称

**返回**: 命名空间详细信息

### create_namespace

创建新的命名空间。

**参数**:
- `name` (string, required) - 命名空间名称

**返回**: 创建的命名空间信息

## 节点管理

### list_nodes

列出集群中的所有节点。

**返回**: 节点列表

### get_node

获取单个节点的详细信息。

**参数**:
- `name` (string, required) - 节点名称

**返回**: 节点详细信息

### describe_node

描述节点的状态和资源使用情况。

**参数**:
- `name` (string, required) - 节点名称

**返回**: 节点状态描述

## 资源状态

### get_resource_usage

获取指定命名空间的资源使用情况。

**参数**:
- `namespace` (string, required) - 命名空间

**返回**: 资源使用统计

### get_cluster_info

获取集群的基本信息。

**返回**: 集群信息

## 配置

Kubernetes 工具通过以下配置进行初始化：

```yaml
kubernetes:
  # kubeconfig 文件路径
  kubeconfig: ""

  # 单个 API 调用超时（秒）
  timeoutSec: 30

  # API 客户端每秒查询数 (QPS)
  qps: 100.0

  # API 客户端突发速率
  burst: 200
```

## 最佳实践

1. **使用过滤参数**: 在可能的情况下使用命名空间和标签过滤来减少返回的数据量
2. **监控资源使用**: 定期使用 `get_resource_usage` 检查资源使用情况
3. **日志管理**: 使用 `tailLines` 参数限制日志返回量
4. **故障排查**: 使用 `describe_pod` 和 `get_pod_events` 诊断 Pod 问题
5. **优雅操作**: 使用 `restart_deployment` 而不是删除 Pod 进行重启

## 相关文档

- [配置指南](/zh/guides/configuration/)
- [部署指南](/zh/guides/deployment/)
- [架构指南](/zh/concepts/architecture/)
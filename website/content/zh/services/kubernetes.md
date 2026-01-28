---
title: "Kubernetes 服务"
weight: 1
---

# Kubernetes 服务

Kubernetes 服务提供全面的容器编排和资源管理功能，包含 28 个专门的工具来管理您的 Kubernetes 集群。

## 概述

Cloud Native MCP Server 中的 Kubernetes 服务使 AI 助手能够高效地管理 Kubernetes 资源。它提供用于部署、服务、配置映射、密钥和其他核心 Kubernetes 资源的工具。

### 主要功能

{{< columns >}}
### 🔧 部署管理
对 Kubernetes 部署进行完全控制，包括创建、更新、扩缩容和删除操作。
<--->

### 🗂️ 资源管理
管理所有 Kubernetes 资源，包括 Pod、服务、配置映射、密钥和持久卷。
{{< /columns >}}

{{< columns >}}
### 📊 监控
获取集群中 Pod、节点和资源使用情况的详细信息。
<--->

### 🔐 安全
管理密钥、RBAC 配置和其他安全相关的 Kubernetes 资源。
{{< /columns >}}

---

## 可用工具 (28)

### Pod 管理
- **kubernetes-get-pods**: 获取命名空间中 Pod 的详细信息
- **kubernetes-list-pods**: 列出命名空间中的所有 Pod
- **kubernetes-get-pod**: 获取特定 Pod 详情
- **kubernetes-delete-pod**: 删除特定 Pod
- **kubernetes-get-pod-logs**: 获取 Pod 的日志
- **kubernetes-get-pod-events**: 获取与 Pod 相关的事件

### 部署管理
- **kubernetes-list-deployments**: 列出命名空间中的所有部署
- **kubernetes-get-deployment**: 获取特定部署详情
- **kubernetes-create-deployment**: 创建新部署
- **kubernetes-update-deployment**: 更新现有部署
- **kubernetes-delete-deployment**: 删除部署
- **kubernetes-scale-deployment**: 扩缩容部署
- **kubernetes-restart-deployment**: 重启部署

### 服务管理
- **kubernetes-list-services**: 列出命名空间中的所有服务
- **kubernetes-get-service**: 获取特定服务详情
- **kubernetes-create-service**: 创建新服务
- **kubernetes-update-service**: 更新现有服务
- **kubernetes-delete-service**: 删除服务

### 配置管理
- **kubernetes-list-configmaps**: 列出命名空间中的所有配置映射
- **kubernetes-get-configmap**: 获取特定配置映射详情
- **kubernetes-create-configmap**: 创建新配置映射
- **kubernetes-update-configmap**: 更新现有配置映射
- **kubernetes-delete-configmap**: 删除配置映射
- **kubernetes-list-secrets**: 列出命名空间中的所有密钥
- **kubernetes-get-secret**: 获取特定密钥详情
- **kubernetes-create-secret**: 创建新密钥
- **kubernetes-update-secret**: 更新现有密钥
- **kubernetes-delete-secret**: 删除密钥

---

## 快速示例

### 列出 default 命名空间中的所有 Pod

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

### 获取特定部署详情

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

### 创建新配置映射

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

## 最佳实践

- 在使用 Kubernetes 资源时始终指定命名空间
- 有效使用标签和注解进行资源组织
- 实施适当的 RBAC 策略以确保安全
- 监控资源使用情况以优化集群性能
- 定期备份关键配置

## 下一步

- [Helm 服务](/zh/services/helm/) 用于包管理
- [配置指南](/zh/guides/configuration/) 了解详细设置
- [安全最佳实践](/zh/guides/security/) 保护您的集群
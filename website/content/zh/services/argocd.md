---
title: "Argo CD 服务"
weight: 10
---

# Argo CD 服务

Argo CD 服务提供 7 个只读 GitOps 检查工具，用于应用、manifest、项目与集群状态查看。

## 概述

当你需要确认某个应用在 Argo CD 中的期望状态、同步状态、健康状态、所属项目以及渲染后的资源清单时，可以直接使用这个服务，而不必离开 MCP 工作流。

## 可用工具 (7)

- `argocd_test_connection`
- `argocd_list_applications_summary`
- `argocd_get_application`
- `argocd_get_application_manifests`
- `argocd_list_projects`
- `argocd_get_project`
- `argocd_list_clusters`

## 推荐调用顺序

1. 先用 `argocd_test_connection`
2. 再用 `argocd_list_applications_summary` 做应用发现
3. 用 `argocd_get_application` 看同步与健康详情
4. 需要渲染资源时再用 `argocd_get_application_manifests`

## 配置示例

```yaml
argocd:
  enabled: true
  url: "https://argocd.example.com"
  username: "admin"
  password: ""
  authToken: ""
  timeoutSec: 30
```

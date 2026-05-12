---
title: "Sentry 服务"
weight: 11
---

# Sentry 服务

Sentry 服务提供 9 个只读工具，用于 Issue 排障、项目发现与 Issue 事件检查。

## 概述

Cloud Native MCP Server 中的 Sentry 服务通过官方 Sentry REST API 暴露组织、项目、Issue 与事件查询能力，适合 AI 助手做错误监控排障与告警调查，不包含写操作。

### 主要功能

{{< columns >}}
### 🚩 Issue 排障
先用紧凑视图浏览 issue，再按需下钻到单个 issue 详情。
<--->

### 🧭 项目发现
浏览 organization 和 project，快速定位正确的 Sentry 范围。
{{< /columns >}}

{{< columns >}}
### 🧾 事件检查
查看某个 issue 关联的具体 event，理解报错实例。
<--->

### 🔐 令牌校验
验证当前 Sentry token 与基础地址是否可用。
{{< /columns >}}

---

## 可用工具 (9)

- **sentry_test_connection**: 校验 Sentry 连通性
- **sentry_list_organizations**: 列出当前 token 可见的 organizations
- **sentry_list_projects**: 列出 organization 下的 projects
- **sentry_get_project**: 按 organization / project slug 获取项目详情
- **sentry_list_issues_summary**: 紧凑型 issue 发现视图
- **sentry_list_issues**: 带过滤条件的完整 issue 列表
- **sentry_get_issue**: 获取单个 issue
- **sentry_list_issue_events**: 列出某个 issue 的 events
- **sentry_get_issue_event**: 获取单个 issue event

## 配置示例

```yaml
sentry:
  enabled: true
  url: "https://sentry.io"
  authToken: "sntrys_..."
  organization: "acme"
  project: "frontend"
  timeoutSec: 30
```

## 下一步

- [配置指南](/zh/docs/configuration/) 查看环境变量与路由
- [Jaeger 服务](/zh/services/jaeger/) 了解分布式链路追踪
- [Langfuse 服务](/zh/services/langfuse/) 了解 LLM 可观测性

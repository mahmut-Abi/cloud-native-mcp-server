---
title: "工具参考"
---

# 工具参考

Cloud Native MCP Server 当前暴露 **311 个工具**。这页只保留适合 LLM 的精简入口。
如果你需要精确的运行时清单，优先查看仓库里的 [docs/TOOLS.md](https://github.com/mahmut-Abi/cloud-native-mcp-server/blob/main/docs/TOOLS.md) 或直接运行 `cloud-native-mcp-server --list tools`。

## 推荐起点

### 基础设施

- `kubernetes_list_resources_summary`
- `kubernetes_get_recent_events`
- `helm_list_releases_paginated`
- `helm_get_release_summary`

### 指标、日志、链路

- `prometheus_query`
- `loki_query_logs_summary`
- `jaeger_get_traces_summary`
- `opentelemetry_get_collector_summary`

### 可视化与仪表板

- `grafana_dashboards_summary`
- `grafana_datasources_summary`
- `grafana_update_dashboard`
- `grafana_render_panel_image`

### 错误监控与 LLM 可观测性

- `sentry_list_issues_summary`
- `sentry_get_issue`
- `langfuse_check_health`
- `langfuse_list_traces_summary`

## 调用建议

- 使用运行时暴露的精确 `snake_case` 工具名。
- 优先从摘要类工具开始，再切到详情工具。
- 参数是对象或数组时，优先发送结构化 JSON。
- 把 `docs/TOOLS.md` 视为仓库里的权威参考。

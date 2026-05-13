---
title: "Tools Reference"
---

# Tools Reference

Cloud Native MCP Server exposes **327 tools**. This page is a compact, LLM-friendly shortlist.
For the exact runtime inventory, prefer [docs/TOOLS.md](https://github.com/mahmut-Abi/cloud-native-mcp-server/blob/main/docs/TOOLS.md) or `cloud-native-mcp-server --list tools`.

## Recommended Starting Points

### Infrastructure

- `kubernetes_list_resources_summary`
- `kubernetes_get_recent_events`
- `helm_list_releases_paginated`
- `helm_get_release_summary`
- `argocd_list_applications_summary`
- `nacos_list_services_summary`

### Metrics, Logs, and Traces

- `prometheus_query`
- `loki_query_logs_summary`
- `jaeger_get_traces_summary`
- `opentelemetry_get_collector_summary`

### Visualization and Dashboards

- `grafana_dashboards_summary`
- `grafana_datasources_summary`
- `grafana_update_dashboard`
- `grafana_render_panel_image`

### Error Monitoring and LLM Observability

- `sentry_list_issues_summary`
- `sentry_get_issue`
- `langfuse_check_health`
- `langfuse_list_traces_summary`

## Calling Guidance

- Use exact runtime `snake_case` tool names.
- Prefer summary tools before full-detail reads.
- Send structured JSON for object and array parameters.
- Treat `docs/TOOLS.md` as the canonical repository reference.

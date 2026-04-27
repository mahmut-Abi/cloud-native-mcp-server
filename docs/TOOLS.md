# Cloud Native MCP Server - Tool Reference

This document provides a curated reference aligned with the current cloud-native MCP server inventory.
For the exact runtime inventory, prefer `--list tools`.

## LLM Calling Rules

- Prefer summary and paginated tools before calling full-detail tools.
- When a parameter represents an object or array, send structured JSON if your MCP client supports it.
- Many handlers still accept legacy JSON strings for compatibility, but structured JSON is preferred.
- Prometheus and tracing timestamps should use RFC3339.
- Kibana tools may accept both `camelCase` and `snake_case` forms for some parameters, but the schema field name remains the canonical form.
- Prefer flat tool arguments over nested `params`, even though Kubernetes handlers now accept nested `params` for compatibility.

## Table of Contents

- [Kubernetes (34 tools)](#kubernetes-34-tools)
- [Helm (34 tools)](#helm-34-tools)
- [Grafana (43 tools)](#grafana-43-tools)
- [Prometheus (20 tools)](#prometheus-20-tools)
- [Kibana (73 tools)](#kibana-73-tools)
- [Elasticsearch (12 tools)](#elasticsearch-12-tools)
- [Alertmanager (16 tools)](#alertmanager-16-tools)
- [Jaeger (8 tools)](#jaeger-8-tools)
- [OpenTelemetry (9 tools)](#opentelemetry-9-tools)
- [Utilities (6 tools)](#utilities-6-tools)

---

## Kubernetes (34 tools)

### Common Response Shapes

When calling Kubernetes tools from scripts, do not assume the return value is always a JSON string.
Some MCP client wrappers already parse the tool result into an object for you.
Inspect the raw return value first before calling `JSON.parse(...)`.

Common Kubernetes response shapes:

| Tool | Typical shape |
|------|---------------|
| `kubernetes_get_api_versions` | `["v1","apps/v1",...]` |
| `kubernetes_get_api_resources` | `[{"name":"pods","kind":"Pod",...}, ...]` |
| `kubernetes_list_resources_summary` | `{"items":[...], "count": N, "pagination": {...}}` |
| `kubernetes_list_resources` | `{"data":{"items":[...]}, "count": N, "pagination": {...}}` |
| `kubernetes_list_resources` with `jsonpath` | `{"data":[...], "count": N, "pagination": {...}}` |
| `kubernetes_list_resources` with `jsonpaths` | `{"data":{"expressions":[...], "data":[...]}, "count": N, "pagination": {...}}` |
| `kubernetes_search_resources` | `{"query":"...", "kinds":[...], "matched": N, "resources":[...]}` |
| `kubernetes_wait_for_resource` | `{"kind":"...", "name":"...", "condition":"...", "message":"...", "attempts": N, ...}` |
| `kubernetes_restart_workload` | `{"status":"ok", "message":"workload restart triggered", "resource": {...}, "wait": {...}?}` |

Practical guidance:

- If your client returns an MCP envelope, the JSON payload is usually in `content[0].text`.
- If your client already returns an object or array, do not run `JSON.parse` on it again.
- For `kubernetes_search_resources`, you may provide `kind` or `resourceTypes`, and `query` or `name`.

### Resource Management

| Tool | Description | Priority |
|------|-------------|----------|
| `kubernetes_list_resources_summary` | List resources with summary (90-95% smaller than full). Returns only essential fields (name, namespace, kind, status, age, labels). | ⚠️ PRIORITY |
| `kubernetes_get_resource_summary` | Get single resource summary with essential fields. Optimized for LLM efficiency. | ⚠️ PRIORITY |
| `kubernetes_list_resources` | List resources with filtering, pagination, single `jsonpath`, or multi-column `jsonpaths` extraction. | - |
| `kubernetes_get_resource` | Get resource details with JSONPath support. Accepts full expressions like `{.status.phase}` and bare paths like `status.phase`. | - |
| `kubernetes_describe_resource` | Describe resource in detail (similar to kubectl describe). | - |
| `kubernetes_create_resource` | Create a resource with structured `metadata` and optional `spec` objects. Legacy JSON string payloads are still accepted. | - |
| `kubernetes_patch_resource` | Patch an existing resource with targeted changes. Use object payloads for `merge`/`apply` and RFC 6902 arrays for `json`. | - |
| `kubernetes_delete_resource` | Delete resource. | - |

### Pod Operations

| Tool | Description | Priority |
|------|-------------|----------|
| `kubernetes_get_pod_logs` | Get pod logs with tailLines support. | - |
| `kubernetes_pod_exec` | Execute command in pod container. | - |
| `kubernetes_scale_resource` | Scale deployment/replicaset. | - |
| `kubernetes_get_rollout_status` | Get rollout status for a workload after patch or scale operations. | - |
| `kubernetes_restart_workload` | Trigger a rollout restart for a supported workload. | - |
| `kubernetes_port_forward` | Port forward to pod. | - |

### Events and Troubleshooting

| Tool | Description | Priority |
|------|-------------|----------|
| `kubernetes_get_recent_events` | Get recent critical events (warnings, errors, failed pods) with 80-90% smaller output. | ⚠️ PRIORITY |
| `kubernetes_get_events` | Get cluster events with filtering support. | - |
| `kubernetes_get_unhealthy_resources` | Find unhealthy resources across cluster. | - |
| `kubernetes_analyze_issue` | Analyze issues and provide recommendations. | - |

### Monitoring and Usage

| Tool | Description | Priority |
|------|-------------|----------|
| `kubernetes_get_resource_usage` | Get resource usage (CPU/Memory) for nodes or pods. | - |
| `kubernetes_get_node_conditions` | Get node conditions and status. | - |
| `kubernetes_cordon_node` | Mark a node unschedulable. | - |
| `kubernetes_uncordon_node` | Mark a node schedulable again. | - |
| `kubernetes_drain_node` | Cordon and drain a node for maintenance. | - |
| `kubernetes_wait_for_resource` | Wait until a resource reaches a desired condition. | - |

### API and Permissions

| Tool | Description | Priority |
|------|-------------|----------|
| `kubernetes_get_api_versions` | Get available API versions. | - |
| `kubernetes_get_api_resources` | Get available resources for API version. | - |
| `kubernetes_check_permissions` | Check RBAC permissions. | - |

### Search and Discovery

| Tool | Description | Priority |
|------|-------------|----------|
| `kubernetes_search_resources` | Search resources by name. Accepts `kind` or `resourceTypes`, and `query` or `name`. If no query is provided, it lists matching resources of the selected kinds up to `limit`. | - |

---

## Helm (34 tools)

### Release Management

| Tool | Description | Priority |
|------|-------------|----------|
| `helm_list_releases_paginated` | List Helm releases with pagination and summary (80-90% smaller). | ⚠️ PRIORITY |
| `helm_list_releases_summary` | List all releases with summary information. | - |
| `helm_get_release_summary` | Get brief summary of a release. | - |
| `helm_get_release` | Get release details. | - |
| `helm_get_release_status` | Get release status. | - |
| `helm_get_release_values` | Get release values. | - |
| `helm_get_release_manifest` | Get release manifest. | - |
| `helm_get_release_history` | Get release history. | - |

### Chart Operations

| Tool | Description | Priority |
|------|-------------|----------|
| `helm_search_charts` | Search Helm charts. | - |
| `helm_get_chart_info` | Get chart information. | - |
| `helm_template_chart` | Template a chart. | - |
| `helm_pull_chart` | Pull chart to local. | - |
| `helm_show_chart` | Show chart details. | - |
| `helm_show_values` | Show chart values. | - |
| `helm_show_readme` | Show chart README. | - |

### Release Lifecycle

| Tool | Description | Priority |
|------|-------------|----------|
| `helm_install_release` | Install a release. | - |
| `helm_upgrade_release` | Upgrade a release. | - |
| `helm_uninstall_release` | Uninstall a release. | - |
| `helm_rollback_release` | Rollback a release. | - |

### Discovery and Health

| Tool | Description | Priority |
|------|-------------|----------|
| `helm_cluster_overview` | Get cluster overview. | - |
| `helm_quick_info` | Quick release info. | - |
| `helm_health_check` | Health check releases. | - |
| `helm_find_broken_releases` | Find releases with failed or pending status. | - |
| `helm_find_releases_by_chart` | Find all releases using a specific chart. | - |
| `helm_list_releases_in_namespace` | List releases in specific namespace. | - |

### Validation and Cache

| Tool | Description | Priority |
|------|-------------|----------|
| `helm_validate_release` | Validate release configuration. | - |
| `helm_clear_cache` | Clear Helm cache. | - |
| `helm_cache_stats` | Get cache statistics. | - |

### Repository Operations

| Tool | Description | Priority |
|------|-------------|----------|
| `helm_add_repository` | Add Helm repository. | - |
| `helm_list_repos` | List Helm repositories. | - |
| `helm_remove_repository` | Remove Helm repository. | - |
| `helm_update_repositories` | Update Helm repositories. | - |

---

## Grafana (43 tools)

### Dashboard Management

| Tool | Description | Priority |
|------|-------------|----------|
| `grafana_dashboards_summary` | List dashboards with minimal output (70-85% smaller). | ⚠️ PRIORITY |
| `grafana_dashboards` | List all dashboards with metadata. | - |
| `grafana_dashboard` | Get specific dashboard by UID. | - |
| `grafana_search_dashboards` | Search dashboards by query. | - |
| `grafana_create_dashboard` | Create new dashboard. | - |
| `grafana_update_dashboard` | Update existing dashboard. | - |
| `grafana_delete_dashboard` | Delete dashboard. | - |

### Data Sources

| Tool | Description | Priority |
|------|-------------|----------|
| `grafana_datasources_summary` | List data sources with minimal output (70-85% smaller). | ⚠️ PRIORITY |
| `grafana_datasources` | List all data sources with smart filtering. | - |
| `grafana_datasource` | Get specific data source by UID. | - |
| `grafana_create_datasource` | Create new data source. | - |
| `grafana_update_datasource` | Update data source. | - |
| `grafana_delete_datasource` | Delete data source. | - |
| `grafana_test_datasource` | Test data source connection. | - |

### Folders

| Tool | Description | Priority |
|------|-------------|----------|
| `grafana_folders` | List all folders with metadata. | - |
| `grafana_folder` | Get specific folder by UID. | - |
| `grafana_create_folder` | Create new folder. | - |
| `grafana_update_folder` | Update folder. | - |
| `grafana_delete_folder` | Delete folder. | - |

### Alerts

| Tool | Description | Priority |
|------|-------------|----------|
| `grafana_alerts` | List alert rules with intelligent limits. | - |
| `grafana_alert_rule` | Get specific alert rule. | - |
| `grafana_create_alert_rule` | Create alert rule. | - |
| `grafana_update_alert_rule` | Update alert rule. | - |
| `grafana_delete_alert_rule` | Delete alert rule. | - |

### Annotations

| Tool | Description | Priority |
|------|-------------|----------|
| `grafana_create_annotation` | Create annotation. | - |
| `grafana_get_annotations` | Get annotations. | - |
| `grafana_update_annotation` | Update annotation. | - |
| `grafana_delete_annotation` | Delete annotation. | - |

### Connection and Health

| Tool | Description | Priority |
|------|-------------|----------|
| `grafana_test_connection` | Test connection to Grafana. | - |
| `grafana_health` | Get Grafana health status. | - |

### Snapshots

| Tool | Description | Priority |
|------|-------------|----------|
| `grafana_create_snapshot` | Create dashboard snapshot. | - |
| `grafana_get_snapshot` | Get snapshot. | - |
| `grafana_delete_snapshot` | Delete snapshot. | - |

---

## Prometheus (20 tools)

### Querying

| Tool | Description | Priority |
|------|-------------|----------|
| `prometheus_query` | Execute instant query. | - |
| `prometheus_query_range` | Execute range query (may return large data). | ⚠️ |
| `prometheus_get_series` | Get time series matching label selector. | ⚠️ |

### Targets

| Tool | Description | Priority |
|------|-------------|----------|
| `prometheus_targets_summary` | Get targets summary (70-80% smaller). | ⚠️ PRIORITY |
| `prometheus_get_targets` | Get current state of targets. | - |

### Alerts

| Tool | Description | Priority |
|------|-------------|----------|
| `prometheus_alerts_summary` | Get alerts summary (70-80% smaller). | ⚠️ PRIORITY |
| `prometheus_get_alerts` | Get current active alerts. | - |

### Rules

| Tool | Description | Priority |
|------|-------------|----------|
| `prometheus_rules_summary` | Get rules summary (70-80% smaller). | ⚠️ PRIORITY |
| `prometheus_get_rules` | Get recording and alerting rules. | - |

### Labels

| Tool | Description | Priority |
|------|-------------|----------|
| `prometheus_get_label_names` | Get all available label names. | - |
| `prometheus_get_label_values` | Get values for specific label. | - |

### Metadata

| Tool | Description | Priority |
|------|-------------|----------|
| `prometheus_get_metrics_metadata` | Get metrics metadata. | - |
| `prometheus_get_target_metadata` | Get target metadata. | - |

### Server Information

| Tool | Description | Priority |
|------|-------------|----------|
| `prometheus_test_connection` | Test connection to Prometheus. | - |
| `prometheus_get_server_info` | Get server information. | - |
| `prometheus_get_runtime_info` | Get runtime and build information. | - |

### TSDB Operations

| Tool | Description | Priority |
|------|-------------|----------|
| `prometheus_get_tsdb_stats` | Get TSDB statistics. | - |
| `prometheus_get_tsdb_status` | Get TSDB status. | - |
| `prometheus_create_snapshot` | Create TSDB snapshot. | - |
| `prometheus_get_wal_replay_status` | Get WAL replay status. | - |

---

## Kibana (73 tools)

### Spaces

| Tool | Description | Priority |
|------|-------------|----------|
| `kibana_get_spaces` | Get all Kibana spaces. | - |
| `kibana_get_space` | Get specific space by ID. | - |
| `kibana_create_space` | Create new space. | - |
| `kibana_update_space` | Update space. | - |
| `kibana_delete_space` | Delete space. | - |

### Index Patterns

| Tool | Description | Priority |
|------|-------------|----------|
| `kibana_get_index_patterns` | Get all index patterns. | - |
| `kibana_get_index_pattern` | Get specific index pattern. | - |
| `kibana_create_index_pattern` | Create index pattern. | - |
| `kibana_update_index_pattern` | Update index pattern. | - |
| `kibana_delete_index_pattern` | Delete index pattern. | - |

### Dashboards

| Tool | Description | Priority |
|------|-------------|----------|
| `kibana_get_dashboards` | Get all dashboards. | - |
| `kibana_get_dashboard` | Get specific dashboard. | - |
| `kibana_create_dashboard` | Create dashboard. | - |
| `kibana_update_dashboard` | Update dashboard. | - |
| `kibana_delete_dashboard` | Delete dashboard. | - |

### Visualizations

| Tool | Description | Priority |
|------|-------------|----------|
| `kibana_get_visualizations` | Get all visualizations. | - |
| `kibana_get_visualization` | Get specific visualization. | - |
| `kibana_create_visualization` | Create visualization. | - |
| `kibana_update_visualization` | Update visualization. | - |
| `kibana_delete_visualization` | Delete visualization. | - |

### Saved Objects

| Tool | Description | Priority |
|------|-------------|----------|
| `kibana_search_saved_objects` | Search saved objects with pagination. | - |
| `kibana_get_saved_object` | Get specific saved object. | - |
| `kibana_create_saved_object` | Create saved object. | - |
| `kibana_update_saved_object` | Update saved object. | - |
| `kibana_delete_saved_object` | Delete saved object. | - |

### Discover

| Tool | Description | Priority |
|------|-------------|----------|
| `kibana_search` | Search documents in Kibana. | - |
| `kibana_get_discover_history` | Get discover history. | - |

### Canvas

| Tool | Description | Priority |
|------|-------------|----------|
| `kibana_get_canvas_workpads` | Get canvas workpads. | - |
| `kibana_get_canvas_workpad` | Get specific workpad. | - |

### Maps

| Tool | Description | Priority |
|------|-------------|----------|
| `kibana_get_maps` | Get maps. | - |
| `kibana_get_map` | Get specific map. | - |

### ML (Machine Learning)

| Tool | Description | Priority |
|------|-------------|----------|
| `kibana_get_ml_jobs` | Get ML jobs. | - |
| `kibana_get_ml_job` | Get specific ML job. | - |

### Security

| Tool | Description | Priority |
|------|-------------|----------|
| `kibana_get_saved_queries` | Get saved queries. | - |

### Advanced Operations

| Tool | Description | Priority |
|------|-------------|----------|
| `kibana_export_saved_objects` | Export saved objects. | - |
| `kibana_import_saved_objects` | Import saved objects. | - |
| `kibana_get_status` | Get Kibana server status. | - |
| `kibana_health_summary` | Get a compact Kibana health summary. | ⚠️ PRIORITY |

---

## Elasticsearch (12 tools)

### Cluster Health

| Tool | Description | Priority |
|------|-------------|----------|
| `elasticsearch_cluster_health_summary` | Get cluster health summary (lightweight). | ⚠️ PRIORITY |
| `elasticsearch_health` | Check cluster health status. | - |

### Indices

| Tool | Description | Priority |
|------|-------------|----------|
| `elasticsearch_list_indices_paginated` | List indices with pagination (80-90% smaller). | ⚠️ PRIORITY |
| `elasticsearch_indices_summary` | List indices summary (75-90% smaller). | ⚠️ PRIORITY |
| `elasticsearch_list_indices` | List all indices. | - |
| `elasticsearch_search_indices` | Search indices with filters. | - |
| `elasticsearch_get_index_detail_advanced` | Advanced index detail retrieval. | 🔍 |
| `elasticsearch_index_stats` | Get index statistics. | - |

### Nodes

| Tool | Description | Priority |
|------|-------------|----------|
| `elasticsearch_nodes_summary` | Get nodes summary (75-90% smaller). | ⚠️ PRIORITY |
| `elasticsearch_nodes` | Get cluster nodes information. | - |

### Cluster Information

| Tool | Description | Priority |
|------|-------------|----------|
| `elasticsearch_get_cluster_detail_advanced` | Advanced cluster detail retrieval. | 🔍 |
| `elasticsearch_info` | Get cluster information. | - |

### Search

| Tool | Description | Priority |
|------|-------------|----------|
| `elasticsearch_search` | Search documents. | - |

---

## Alertmanager (16 tools)

### Alerts

| Tool | Description | Priority |
|------|-------------|----------|
| `alertmanager_alerts_summary` | Get alerts summary (85-95% smaller). | ⚠️ PRIORITY |
| `alertmanager_query_alerts` | Query alerts with filters. | - |
| `alertmanager_query_alerts_advanced` | Advanced alert query with pagination. | 🔍 |
| `alertmanager_get_alerts` | Get current alerts. | - |
| `alertmanager_get_alert_groups` | Get alert groups. | - |

### Silences

| Tool | Description | Priority |
|------|-------------|----------|
| `alertmanager_silences_summary` | Get silences summary (85-95% smaller). | ⚠️ PRIORITY |
| `alertmanager_get_silences` | Get current silences. | - |
| `alertmanager_create_silence` | Create new silence. | - |
| `alertmanager_delete_silence` | Delete silence. | - |

### Receivers

| Tool | Description | Priority |
|------|-------------|----------|
| `alertmanager_receivers_summary` | Get receivers summary (85-95% smaller). | ⚠️ PRIORITY |
| `alertmanager_get_receivers` | Get configured receivers. | - |
| `alertmanager_test_receiver` | Test receiver configuration. | - |

### Status and Health

| Tool | Description | Priority |
|------|-------------|----------|
| `alertmanager_health_summary` | Get health and status summary. | ⚠️ PRIORITY |
| `alertmanager_get_status` | Get Alertmanager status. | - |

---

## Jaeger (8 tools)

### Traces

| Tool | Description | Priority |
|------|-------------|----------|
| `jaeger_get_traces_summary` | Get traces summary (70-85% smaller). | ⚠️ PRIORITY |
| `jaeger_get_traces` | Retrieve traces with filtering. | - |
| `jaeger_get_trace` | Get specific trace by ID. | ⚠️ PRIORITY |
| `jaeger_search_traces` | Search traces with advanced filtering. | - |

### Services

| Tool | Description | Priority |
|------|-------------|----------|
| `jaeger_get_services` | Get all registered services. | - |
| `jaeger_get_service_ops` | Get operations for specific service. | - |

### Dependencies

| Tool | Description | Priority |
|------|-------------|----------|
| `jaeger_get_dependencies` | Get service dependency graph. | - |

---

## OpenTelemetry (9 tools)

### Metrics

| Tool | Description | Priority |
|------|-------------|----------|
| `opentelemetry_get_metrics` | Retrieve metrics from OpenTelemetry Collector. Can filter by metric name and time range. | - |
| `opentelemetry_query_metrics` | Execute a PromQL-style query against OpenTelemetry Collector metrics. Useful for aggregations, filtering, and complex metric analysis. | - |

### Traces

| Tool | Description | Priority |
|------|-------------|----------|
| `opentelemetry_get_traces` | Retrieve traces from OpenTelemetry Collector. Can filter by trace ID, service name, and time range. | - |
| `opentelemetry_query_traces` | Search for traces matching custom criteria in OpenTelemetry Collector. Supports filtering by service, tags, and time range. | - |

### Logs

| Tool | Description | Priority |
|------|-------------|----------|
| `opentelemetry_get_logs` | Retrieve logs from OpenTelemetry Collector. Can filter by service name, log level, and time range. | - |
| `opentelemetry_query_logs` | Search for logs matching custom criteria in OpenTelemetry Collector. Supports filtering by service, level, message content, and time range. | - |

### Health and Status

| Tool | Description | Priority |
|------|-------------|----------|
| `opentelemetry_get_health` | Check the health status of OpenTelemetry Collector. Returns overall health and component status. | ⚠️ PRIORITY |
| `opentelemetry_get_status` | Retrieve detailed status information about OpenTelemetry Collector, including components, pipelines, and configuration. | - |

### Configuration

| Tool | Description | Priority |
|------|-------------|----------|
| `opentelemetry_get_config` | Retrieve the current configuration of OpenTelemetry Collector. Shows pipelines, receivers, processors, exporters, and extensions. | - |

---

## Utilities (6 tools)

### Time

| Tool | Description | Priority |
|------|-------------|----------|
| `utilities_get_time` | Get current time in specified format. | - |
| `utilities_get_timestamp` | Get Unix timestamp. | - |
| `utilities_get_date` | Get current date. | - |

### Execution Control

| Tool | Description | Priority |
|------|-------------|----------|
| `utilities_pause` | Pause execution for specified seconds. | - |
| `utilities_sleep` | Sleep for specified duration. | - |

### Web

| Tool | Description | Priority |
|------|-------------|----------|
| `utilities_web_fetch` | Fetch content from URL. | - |

---

## Tool Priority Legend

| Marker | Description |
|--------|-------------|
| ⚠️ PRIORITY | LLM-optimized tool with significantly smaller output (70-95% reduction). Recommended for most use cases. |
| ⚠️ | May return large amounts of data. Use with caution and consider filters/limits. |
| 🔍 | Advanced tool for detailed analysis. Use when comprehensive information needed. |
| - | Standard tool for general use. |

## Best Practices

1. **Start with summary tools**: Always use tools marked with ⚠️ PRIORITY first to get essential information
2. **Use pagination**: For large datasets, use paginated tools to avoid context overflow
3. **Apply filters**: Use filters and selectors to narrow down results
4. **Check limits**: Be mindful of default and maximum limits on tools
5. **Use debug mode**: Enable debug when troubleshooting tool execution issues

## Tool Categories

### Discovery Tools
- `*_summary` - Quick overview with minimal data
- `*_list_*` - List resources with basic info
- `*_get_*` - Get specific resource details

### Analysis Tools
- `*_search_*` - Search with filters
- `*_query_*` - Query with parameters
- `*_analyze_*` - Analyze and provide insights

### Management Tools
- `*_create_*` - Create new resources
- `*_update_*` - Modify existing resources
- `*_delete_*` - Remove resources

### Monitoring Tools
- `*_health` - Health status checks
- `*_status` - Status information
- `*_metrics` - Performance metrics

### Utility Tools
- `*_test_*` - Test connections/configurations
- `*_validate_*` - Validate configurations
- `*_cache_*` - Cache management

---

<!-- BEGIN GENERATED TOOL INVENTORY -->
## Generated Inventory

This section is generated from `internal/services/**/tools/*.go`.
Do not edit this block by hand.

### Kubernetes (34 tools)

- `kubernetes_analyze_issue`
- `kubernetes_check_permissions`
- `kubernetes_cordon_node`
- `kubernetes_create_resource`
- `kubernetes_delete_resource`
- `kubernetes_describe_resource`
- `kubernetes_drain_node`
- `kubernetes_get_api_resources`
- `kubernetes_get_api_versions`
- `kubernetes_get_events`
- `kubernetes_get_events_detail`
- `kubernetes_get_node_conditions`
- `kubernetes_get_pod_logs`
- `kubernetes_get_recent_events`
- `kubernetes_get_resource`
- `kubernetes_get_resource_detail_advanced`
- `kubernetes_get_resource_details`
- `kubernetes_get_resource_summary`
- `kubernetes_get_resource_usage`
- `kubernetes_get_resources_detail`
- `kubernetes_get_rollout_status`
- `kubernetes_get_unhealthy_resources`
- `kubernetes_list_resources`
- `kubernetes_list_resources_full`
- `kubernetes_list_resources_summary`
- `kubernetes_patch_resource`
- `kubernetes_pod_exec`
- `kubernetes_port_forward`
- `kubernetes_restart_workload`
- `kubernetes_scale_resource`
- `kubernetes_search_resources`
- `kubernetes_test_tool`
- `kubernetes_uncordon_node`
- `kubernetes_wait_for_resource`

### Helm (34 tools)

- `helm_add_repository`
- `helm_cache_stats`
- `helm_clear_cache`
- `helm_cluster_overview`
- `helm_compare_revisions`
- `helm_find_broken_releases`
- `helm_find_releases_by_chart`
- `helm_find_releases_by_labels`
- `helm_get_chart_info`
- `helm_get_recent_failures`
- `helm_get_release`
- `helm_get_release_history`
- `helm_get_release_history_paginated`
- `helm_get_release_manifest`
- `helm_get_release_status`
- `helm_get_release_summary`
- `helm_get_release_values`
- `helm_get_resources_of_release`
- `helm_health_check`
- `helm_install_release`
- `helm_list_releases`
- `helm_list_releases_in_namespace`
- `helm_list_releases_paginated`
- `helm_list_releases_summary`
- `helm_list_repos`
- `helm_quick_info`
- `helm_remove_repository`
- `helm_rollback_release`
- `helm_search_charts`
- `helm_template_chart`
- `helm_uninstall_release`
- `helm_update_repositories`
- `helm_upgrade_release`
- `helm_validate_release`

### Grafana (43 tools)

- `grafana_alerts`
- `grafana_check_datasource_health`
- `grafana_create_alert_rule`
- `grafana_create_annotation`
- `grafana_create_datasource`
- `grafana_create_graphite_annotation`
- `grafana_current_user`
- `grafana_dashboard`
- `grafana_dashboards`
- `grafana_dashboards_summary`
- `grafana_datasource_detail`
- `grafana_datasources`
- `grafana_datasources_summary`
- `grafana_delete_alert_rule`
- `grafana_delete_datasource`
- `grafana_folder_detail`
- `grafana_folders`
- `grafana_generate_deeplink`
- `grafana_get_alert_rule_by_uid`
- `grafana_get_annotation_tags`
- `grafana_get_annotations`
- `grafana_get_dashboard_panel_queries`
- `grafana_get_dashboard_property`
- `grafana_get_datasource_by_name`
- `grafana_get_resource_description`
- `grafana_get_resource_permissions`
- `grafana_get_role_assignments`
- `grafana_get_role_details`
- `grafana_list_all_roles`
- `grafana_list_contact_points`
- `grafana_list_team_roles`
- `grafana_list_teams`
- `grafana_list_user_roles`
- `grafana_organization`
- `grafana_patch_annotation`
- `grafana_render_panel_image`
- `grafana_search_dashboards`
- `grafana_test_connection`
- `grafana_update_alert_rule`
- `grafana_update_annotation`
- `grafana_update_dashboard`
- `grafana_update_datasource`
- `grafana_users`

### Prometheus (20 tools)

- `prometheus_alerts_summary`
- `prometheus_create_snapshot`
- `prometheus_get_alerts`
- `prometheus_get_label_names`
- `prometheus_get_label_values`
- `prometheus_get_metrics_metadata`
- `prometheus_get_rules`
- `prometheus_get_runtime_info`
- `prometheus_get_series`
- `prometheus_get_server_info`
- `prometheus_get_target_metadata`
- `prometheus_get_targets`
- `prometheus_get_tsdb_stats`
- `prometheus_get_tsdb_status`
- `prometheus_get_wal_replay_status`
- `prometheus_query`
- `prometheus_query_range`
- `prometheus_rules_summary`
- `prometheus_targets_summary`
- `prometheus_test_connection`

### Kibana (73 tools)

- `kibana_bulk_delete_saved_objects`
- `kibana_clone_dashboard`
- `kibana_clone_visualization`
- `kibana_create_alert_rule`
- `kibana_create_connector`
- `kibana_create_dashboard`
- `kibana_create_data_view`
- `kibana_create_index_pattern`
- `kibana_create_saved_object`
- `kibana_create_space`
- `kibana_create_visualization`
- `kibana_dashboards_paginated`
- `kibana_dashboards_summary`
- `kibana_delete_alert_rule`
- `kibana_delete_connector`
- `kibana_delete_dashboard`
- `kibana_delete_data_view`
- `kibana_delete_index_pattern`
- `kibana_delete_saved_object`
- `kibana_delete_space`
- `kibana_delete_visualization`
- `kibana_disable_alert_rule`
- `kibana_enable_alert_rule`
- `kibana_export_saved_objects`
- `kibana_get_alert_rule`
- `kibana_get_alert_rule_history`
- `kibana_get_alert_rule_types`
- `kibana_get_alert_rules`
- `kibana_get_alerts`
- `kibana_get_canvas_workpads`
- `kibana_get_connector`
- `kibana_get_connector_types`
- `kibana_get_connectors`
- `kibana_get_dashboard`
- `kibana_get_dashboard_detail_advanced`
- `kibana_get_dashboards`
- `kibana_get_data_view`
- `kibana_get_data_views`
- `kibana_get_index_pattern`
- `kibana_get_index_pattern_fields`
- `kibana_get_index_patterns`
- `kibana_get_lens_objects`
- `kibana_get_maps`
- `kibana_get_saved_search`
- `kibana_get_saved_searches`
- `kibana_get_space`
- `kibana_get_spaces`
- `kibana_get_status`
- `kibana_get_visualization`
- `kibana_get_visualizations`
- `kibana_health_summary`
- `kibana_import_saved_objects`
- `kibana_index_patterns_summary`
- `kibana_mute_alert_rule`
- `kibana_query_logs`
- `kibana_refresh_index_pattern_fields`
- `kibana_search_saved_objects`
- `kibana_search_saved_objects_advanced`
- `kibana_set_default_index_pattern`
- `kibana_spaces_summary`
- `kibana_test_connection`
- `kibana_test_connector`
- `kibana_unmute_alert_rule`
- `kibana_update_alert_rule`
- `kibana_update_connector`
- `kibana_update_dashboard`
- `kibana_update_data_view`
- `kibana_update_index_pattern`
- `kibana_update_saved_object`
- `kibana_update_space`
- `kibana_update_visualization`
- `kibana_visualizations_paginated`
- `kibana_visualizations_summary`

### Elasticsearch (12 tools)

- `elasticsearch_cluster_health_summary`
- `elasticsearch_get_cluster_detail_advanced`
- `elasticsearch_get_index_detail_advanced`
- `elasticsearch_health`
- `elasticsearch_index_stats`
- `elasticsearch_indices_summary`
- `elasticsearch_info`
- `elasticsearch_list_indices`
- `elasticsearch_list_indices_paginated`
- `elasticsearch_nodes`
- `elasticsearch_nodes_summary`
- `elasticsearch_search_indices`

### Alertmanager (16 tools)

- `alertmanager_alert_groups_paginated`
- `alertmanager_alerts_summary`
- `alertmanager_create_silence`
- `alertmanager_delete_silence`
- `alertmanager_get_alert_groups`
- `alertmanager_get_alerts`
- `alertmanager_get_receivers`
- `alertmanager_get_silences`
- `alertmanager_get_status`
- `alertmanager_health_summary`
- `alertmanager_query_alerts`
- `alertmanager_query_alerts_advanced`
- `alertmanager_receivers_summary`
- `alertmanager_silences_paginated`
- `alertmanager_silences_summary`
- `alertmanager_test_receiver`

### Jaeger (8 tools)

- `jaeger_get_dependencies`
- `jaeger_get_service_ops`
- `jaeger_get_services`
- `jaeger_get_services_summary`
- `jaeger_get_trace`
- `jaeger_get_traces`
- `jaeger_get_traces_summary`
- `jaeger_search_traces`

### OpenTelemetry (9 tools)

- `opentelemetry_get_config`
- `opentelemetry_get_health`
- `opentelemetry_get_logs`
- `opentelemetry_get_metrics`
- `opentelemetry_get_status`
- `opentelemetry_get_traces`
- `opentelemetry_query_logs`
- `opentelemetry_query_metrics`
- `opentelemetry_query_traces`

### Utilities (6 tools)

- `utilities_get_date`
- `utilities_get_time`
- `utilities_get_timestamp`
- `utilities_pause`
- `utilities_sleep`
- `utilities_web_fetch`

<!-- END GENERATED TOOL INVENTORY -->
---

For more information, see the main [README](../README.md).

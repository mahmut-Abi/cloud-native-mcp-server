# K8s MCP Server - Complete Tools Reference

This document provides a comprehensive reference for all 210+ tools available in the K8s MCP Server.

## Table of Contents

- [Kubernetes (28 tools)](#kubernetes-28-tools)
- [Helm (31 tools)](#helm-31-tools)
- [Grafana (36 tools)](#grafana-36-tools)
- [Prometheus (20 tools)](#prometheus-20-tools)
- [Kibana (52 tools)](#kibana-52-tools)
- [Elasticsearch (14 tools)](#elasticsearch-14-tools)
- [Alertmanager (15 tools)](#alertmanager-15-tools)
- [Jaeger (8 tools)](#jaeger-8-tools)
- [Utilities (6 tools)](#utilities-6-tools)

---

## Kubernetes (28 tools)

### Resource Management

| Tool | Description | Priority |
|------|-------------|----------|
| `kubernetes_list_resources_summary` | List resources with summary (90-95% smaller than full). Returns only essential fields (name, namespace, kind, status, age, labels). | ⚠️ PRIORITY |
| `kubernetes_get_resource_summary` | Get single resource summary with essential fields. Optimized for LLM efficiency. | ⚠️ PRIORITY |
| `kubernetes_list_resources` | List all resources with full details. | - |
| `kubernetes_get_resource` | Get resource details with JSONPath support. | - |
| `kubernetes_describe_resource` | Describe resource in detail (similar to kubectl describe). | - |
| `kubernetes_create_resource` | Create Kubernetes resource from YAML/JSON. | - |
| `kubernetes_update_resource` | Update existing resource. | - |
| `kubernetes_delete_resource` | Delete resource. | - |

### Pod Operations

| Tool | Description | Priority |
|------|-------------|----------|
| `kubernetes_get_pod_logs` | Get pod logs with tailLines support. | - |
| `kubernetes_pod_exec` | Execute command in pod container. | - |
| `kubernetes_scale_resource` | Scale deployment/replicaset. | - |
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

### API and Permissions

| Tool | Description | Priority |
|------|-------------|----------|
| `kubernetes_get_api_versions` | Get available API versions. | - |
| `kubernetes_get_api_resources` | Get available resources for API version. | - |
| `kubernetes_check_permissions` | Check RBAC permissions. | - |

### Search and Discovery

| Tool | Description | Priority |
|------|-------------|----------|
| `kubernetes_search_resources` | Search resources by labels, annotations, or fields. | - |

### Advanced Operations

| Tool | Description | Priority |
|------|-------------|----------|
| `kubernetes_apply_manifest` | Apply Kubernetes manifest. | - |
| `kubernetes_get_config` | Get Kubernetes configuration. | - |
| `kubernetes_cluster_info` | Get cluster information. | - |

---

## Helm (31 tools)

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
| `helm_add_repo` | Add Helm repository. | - |
| `helm_list_repos` | List Helm repositories. | - |
| `helm_remove_repo` | Remove Helm repository. | - |
| `helm_update_repo` | Update Helm repository. | - |

---

## Grafana (36 tools)

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

## Kibana (52 tools)

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
| `kibana_export_objects` | Export saved objects. | - |
| `kibana_import_objects` | Import saved objects. | - |
| `kibana_get_cluster_info` | Get cluster information. | - |
| `kibana_health` | Get Kibana health status. | - |

---

## Elasticsearch (14 tools)

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

## Alertmanager (15 tools)

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

For more information, see the main [README](../README.md).
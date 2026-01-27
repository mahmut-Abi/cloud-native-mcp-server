---
title: "Tools Reference"
---

# Tools Reference

Cloud Native MCP Server provides 220+ powerful tools covering Kubernetes management, application deployment, monitoring, log analysis, and more.

## Kubernetes Tools (28)

### Pod Management
- `list_pods` - List pods
- `get_pod` - Get pod details
- `describe_pod` - Describe pod status
- `delete_pod` - Delete pod
- `get_pod_logs` - Get pod logs
- `get_pod_events` - Get pod events

### Deployment Management
- `list_deployments` - List deployments
- `get_deployment` - Get deployment details
- `create_deployment` - Create deployment
- `update_deployment` - Update deployment
- `delete_deployment` - Delete deployment
- `scale_deployment` - Scale deployment
- `restart_deployment` - Restart deployment

### Service Management
- `list_services` - List services
- `get_service` - Get service details
- `create_service` - Create service
- `delete_service` - Delete service

### ConfigMap & Secret
- `list_configmaps` - List ConfigMaps
- `get_configmap` - Get ConfigMap details
- `create_configmap` - Create ConfigMap
- `list_secrets` - List secrets
- `get_secret` - Get secret details
- `create_secret` - Create secret

### Namespaces
- `list_namespaces` - List namespaces
- `get_namespace` - Get namespace details
- `create_namespace` - Create namespace

### Node Management
- `list_nodes` - List nodes
- `get_node` - Get node details
- `describe_node` - Describe node status

### Resource Status
- `get_resource_usage` - Get resource usage
- `get_cluster_info` - Get cluster information

## Helm Tools (31)

### Chart Management
- `list_repositories` - List Helm repositories
- `add_repository` - Add Helm repository
- `remove_repository` - Remove Helm repository
- `update_repository` - Update Helm repository
- `search_chart` - Search chart
- `show_chart` - Show chart details
- `pull_chart` - Pull chart

### Release Management
- `list_releases` - List releases
- `get_release` - Get release details
- `install_chart` - Install chart
- `upgrade_release` - Upgrade release
- `rollback_release` - Rollback release
- `uninstall_release` - Uninstall release
- `get_release_history` - Get release history
- `get_release_status` - Get release status
- `get_release_values` - Get release configuration values

### Values Management
- `get_values` - Get configuration values
- `set_values` - Set configuration values
- `diff_values` - Compare configuration value differences

### Release Operations
- `test_release` - Test release
- `lint_chart` - Lint chart
- `package_chart` - Package chart
- `verify_chart` - Verify chart
- `template_chart` - Generate template

### Chart Dependencies
- `list_dependencies` - List dependencies
- `update_dependencies` - Update dependencies

### Plugin Management
- `list_plugins` - List plugins
- `install_plugin` - Install plugin

### Version Management
- `list_versions` - List chart versions
- `get_version_info` - Get version information

### Debugging Tools
- `debug_release` - Debug release

## Grafana Tools (36)

### Dashboard Management
- `list_dashboards` - List dashboards
- `get_dashboard` - Get dashboard details
- `create_dashboard` - Create dashboard
- `update_dashboard` - Update dashboard
- `delete_dashboard` - Delete dashboard
- `import_dashboard` - Import dashboard
- `export_dashboard` - Export dashboard
- `search_dashboards` - Search dashboards
- `get_dashboard_by_uid` - Get dashboard by UID
- `get_dashboard_by_tag` - Get dashboard by tag

### Datasource Management
- `list_datasources` - List datasources
- `get_datasource` - Get datasource details
- `create_datasource` - Create datasource
- `update_datasource` - Update datasource
- `delete_datasource` - Delete datasource
- `test_datasource` - Test datasource connection

### Folder Management
- `list_folders` - List folders
- `get_folder` - Get folder details
- `create_folder` - Create folder
- `update_folder` - Update folder
- `delete_folder` - Delete folder

### Query Execution
- `execute_query` - Execute query
- `execute_multiple_queries` - Execute multiple queries
- `query_metrics` - Query metrics

### Alert Management
- `list_alerts` - List alerts
- `get_alert` - Get alert details
- `pause_alert` - Pause alert
- `resume_alert` - Resume alert
- `get_alert_rules` - Get alert rules

### User Management
- `list_users` - List users
- `get_user` - Get user details
- `create_user` - Create user

### Organization Management
- `list_organizations` - List organizations
- `get_organization` - Get organization details

### Health Check
- `get_health` - Get health status
- `get_version` - Get version information

## Prometheus Tools (20)

### Query Execution
- `query` - Execute instant query
- `query_range` - Execute range query
- `query_exemplars` - Query exemplar data

### Metadata Queries
- `label_names` - Get label names
- `label_values` - Get label values
- `series` - Get time series
- `metadata` - Get metadata

### Target Management
- `targets` - Get target list
- `get_target_metadata` - Get target metadata

### Rules Management
- `rules` - Get rules list
- `get_alerts` - Get alerts list

### Configuration Management
- `config` - Get configuration information
- `flags` - Get startup flags

### Status Queries
- `status` - Get status information
- `query_stats` - Get query statistics

### Snapshot Management
- `snapshot` - Create snapshot

### TSDB Operations
- `tsdb_stats` - Get TSDB statistics
- `tsdb_series` - Get TSDB series

### Storage Operations
- `block_info` - Get block information

## Kibana Tools (52)

### Index Management
- `list_indices` - List indices
- `get_index` - Get index details
- `create_index` - Create index
- `delete_index` - Delete index
- `get_index_stats` - Get index statistics
- `get_index_settings` - Get index settings
- `update_index_settings` - Update index settings

### Document Operations
- `search_documents` - Search documents
- `get_document` - Get document
- `create_document` - Create document
- `update_document` - Update document
- `delete_document` - Delete document
- `bulk_operations` - Bulk operations

### Query Building
- `build_query` - Build query
- `execute_query` - Execute query
- `aggregate_data` - Aggregate data
- `get_query_stats` - Get query statistics

### Visualizations
- `list_visualizations` - List visualizations
- `get_visualization` - Get visualization
- `create_visualization` - Create visualization
- `update_visualization` - Update visualization
- `delete_visualization` - Delete visualization

### Dashboards
- `list_dashboards` - List dashboards
- `get_dashboard` - Get dashboard
- `create_dashboard` - Create dashboard
- `update_dashboard` - Update dashboard
- `delete_dashboard` - Delete dashboard

### Index Patterns
- `list_index_patterns` - List index patterns
- `get_index_pattern` - Get index pattern
- `create_index_pattern` - Create index pattern
- `update_index_pattern` - Update index pattern
- `delete_index_pattern` - Delete index pattern

### Saved Queries
- `list_saved_queries` - List saved queries
- `get_saved_query` - Get saved query
- `create_saved_query` - Create saved query
- `update_saved_query` - Update saved query
- `delete_saved_query` - Delete saved query

### Space Management
- `list_spaces` - List spaces
- `get_space` - Get space
- `create_space` - Create space
- `update_space` - Update space
- `delete_space` - Delete space

### Discover
- `discover_data` - Discover data
- `get_field_capabilities` - Get field capabilities

### Export/Import
- `export_objects` - Export objects
- `import_objects` - Import objects

### Short URLs
- `create_short_url` - Create short URL

## Elasticsearch Tools (14)

### Index Management
- `list_indices` - List indices
- `get_index` - Get index
- `create_index` - Create index
- `delete_index` - Delete index
- `get_index_stats` - Get index statistics

### Document Operations
- `index_document` - Index document
- `get_document` - Get document
- `search_documents` - Search documents
- `update_document` - Update document
- `delete_document` - Delete document

### Cluster Management
- `get_cluster_health` - Get cluster health
- `get_cluster_stats` - Get cluster statistics
- `get_cluster_info` - Get cluster information

### Alias Management
- `get_aliases` - Get aliases

## Alertmanager Tools (15)

### Alert Management
- `list_alerts` - List alerts
- `get_alert` - Get alert details
- `get_alert_groups` - Get alert groups
- `get_silences` - Get silences
- `create_silence` - Create silence
- `delete_silence` - Delete silence
- `expire_silence` - Expire silence

### Rules Management
- `get_alert_rules` - Get alert rules
- `list_rule_groups` - List rule groups

### Configuration Management
- `get_config` - Get configuration
- `get_status` - Get status

### Notification Management
- `list_notifications` - List notifications
- `get_receivers` - Get receiver configuration
- `list_routes` - List routes

### Health Check
- `get_health` - Get health status

## Jaeger Tools (8)

### Trace Queries
- `get_trace` - Get trace
- `search_traces` - Search traces
- `get_services` - Get service list
- `get_operations` - Get operation list

### Dependency Analysis
- `get_dependencies` - Get dependencies

### Metrics Queries
- `get_metrics` - Get metrics

### Configuration Queries
- `get_config` - Get configuration
- `get_status` - Get status

## OpenTelemetry Tools (9)

### Metrics Management
- `get_metrics` - Get metrics
- `get_metric_data` - Get metric data
- `list_metric_streams` - List metric streams

### Trace Management
- `get_traces` - Get traces
- `search_traces` - Search traces

### Log Management
- `get_logs` - Get logs
- `search_logs` - Search logs

### Configuration Management
- `get_config` - Get configuration
- `get_status` - Get status

## Utilities Tools (6)

### General Tools
- `base64_encode` - Base64 encode
- `base64_decode` - Base64 decode
- `json_parse` - JSON parse
- `json_stringify` - JSON stringify
- `timestamp` - Get timestamp
- `uuid` - Generate UUID

## Tool Call Examples

### Kubernetes - List Pods

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "list_pods",
    "arguments": {
      "namespace": "default"
    }
  }
}
```

### Helm - Install Chart

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "install_chart",
    "arguments": {
      "chart": "nginx-ingress",
      "repo": "https://kubernetes.github.io/ingress-nginx",
      "release": "my-nginx",
      "namespace": "ingress-nginx"
    }
  }
}
```

### Prometheus - Query Metrics

```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "query",
    "arguments": {
      "query": "up{job=\"kubernetes-pods\"}"
    }
  }
}
```

### Grafana - List Dashboards

```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/call",
  "params": {
    "name": "list_dashboards",
    "arguments": {}
  }
}
```

## Tool Parameter Description

All tools support the following common parameters:

- `timeout` - Request timeout (seconds)
- `dry_run` - Dry run mode, does not actually execute
- `verbose` - Verbose output mode

For tool-specific parameters, please refer to the detailed documentation of each service.

## Error Handling

Tool calls may return the following errors:

- `InvalidParams` - Invalid parameters
- `NotFound` - Resource does not exist
- `PermissionDenied` - Insufficient permissions
- `Timeout` - Request timeout
- `InternalError` - Internal error

Error response format:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32602,
    "message": "Invalid params",
    "data": {
      "details": "namespace is required"
    }
  }
}
```
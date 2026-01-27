---
title: "工具参考"
---

# 工具参考

Cloud Native MCP Server 提供 220+ 个强大的工具，覆盖 Kubernetes 管理和应用部署、监控、日志分析等各个方面。

## Kubernetes 工具（28 个）

### Pod 管理
- `list_pods` - 列出 Pod
- `get_pod` - 获取 Pod 详情
- `describe_pod` - 描述 Pod 状态
- `delete_pod` - 删除 Pod
- `get_pod_logs` - 获取 Pod 日志
- `get_pod_events` - 获取 Pod 事件

### Deployment 管理
- `list_deployments` - 列出 Deployment
- `get_deployment` - 获取 Deployment 详情
- `create_deployment` - 创建 Deployment
- `update_deployment` - 更新 Deployment
- `delete_deployment` - 删除 Deployment
- `scale_deployment` - 扩缩容 Deployment
- `restart_deployment` - 重启 Deployment

### Service 管理
- `list_services` - 列出 Service
- `get_service` - 获取 Service 详情
- `create_service` - 创建 Service
- `delete_service` - 删除 Service

### ConfigMap & Secret
- `list_configmaps` - 列出 ConfigMap
- `get_configmap` - 获取 ConfigMap 详情
- `create_configmap` - 创建 ConfigMap
- `list_secrets` - 列出 Secret
- `get_secret` - 获取 Secret 详情
- `create_secret` - 创建 Secret

### 命名空间
- `list_namespaces` - 列出命名空间
- `get_namespace` - 获取命名空间详情
- `create_namespace` - 创建命名空间

### 节点管理
- `list_nodes` - 列出节点
- `get_node` - 获取节点详情
- `describe_node` - 描述节点状态

### 资源状态
- `get_resource_usage` - 获取资源使用情况
- `get_cluster_info` - 获取集群信息

## Helm 工具（31 个）

### Chart 管理
- `list_repositories` - 列出 Helm 仓库
- `add_repository` - 添加 Helm 仓库
- `remove_repository` - 移除 Helm 仓库
- `update_repository` - 更新 Helm 仓库
- `search_chart` - 搜索 Chart
- `show_chart` - 显示 Chart 详情
- `pull_chart` - 下载 Chart

### Release 管理
- `list_releases` - 列出 Release
- `get_release` - 获取 Release 详情
- `install_chart` - 安装 Chart
- `upgrade_release` - 升级 Release
- `rollback_release` - 回滚 Release
- `uninstall_release` - 卸载 Release
- `get_release_history` - 获取 Release 历史
- `get_release_status` - 获取 Release 状态
- `get_release_values` - 获取 Release 配置值

### Values 管理
- `get_values` - 获取配置值
- `set_values` - 设置配置值
- `diff_values` - 比较配置值差异

### Release 操作
- `test_release` - 测试 Release
- `lint_chart` - 检查 Chart
- `package_chart` - 打包 Chart
- `verify_chart` - 验证 Chart
- `template_chart` - 生成模板

### Chart 依赖
- `list_dependencies` - 列出依赖
- `update_dependencies` - 更新依赖

### 插件管理
- `list_plugins` - 列出插件
- `install_plugin` - 安装插件

### 版本管理
- `list_versions` - 列出 Chart 版本
- `get_version_info` - 获取版本信息

### 调试工具
- `debug_release` - 调试 Release

## Grafana 工具（36 个）

### 仪表板管理
- `list_dashboards` - 列出仪表板
- `get_dashboard` - 获取仪表板详情
- `create_dashboard` - 创建仪表板
- `update_dashboard` - 更新仪表板
- `delete_dashboard` - 删除仪表板
- `import_dashboard` - 导入仪表板
- `export_dashboard` - 导出仪表板
- `search_dashboards` - 搜索仪表板
- `get_dashboard_by_uid` - 通过 UID 获取仪表板
- `get_dashboard_by_tag` - 通过标签获取仪表板

### 数据源管理
- `list_datasources` - 列出数据源
- `get_datasource` - 获取数据源详情
- `create_datasource` - 创建数据源
- `update_datasource` - 更新数据源
- `delete_datasource` - 删除数据源
- `test_datasource` - 测试数据源连接

### 文件夹管理
- `list_folders` - 列出文件夹
- `get_folder` - 获取文件夹详情
- `create_folder` - 创建文件夹
- `update_folder` - 更新文件夹
- `delete_folder` - 删除文件夹

### 查询执行
- `execute_query` - 执行查询
- `execute_multiple_queries` - 执行多个查询
- `query_metrics` - 查询指标

### 告警管理
- `list_alerts` - 列出告警
- `get_alert` - 获取告警详情
- `pause_alert` - 暂停告警
- `resume_alert` - 恢复告警
- `get_alert_rules` - 获取告警规则

### 用户管理
- `list_users` - 列出用户
- `get_user` - 获取用户详情
- `create_user` - 创建用户

### 组织管理
- `list_organizations` - 列出组织
- `get_organization` - 获取组织详情

### 健康检查
- `get_health` - 获取健康状态
- `get_version` - 获取版本信息

## Prometheus 工具（20 个）

### 查询执行
- `query` - 执行即时查询
- `query_range` - 执行范围查询
- `query_exemplars` - 查询示例数据

### 元数据查询
- `label_names` - 获取标签名
- `label_values` - 获取标签值
- `series` - 获取时间序列
- `metadata` - 获取元数据

### 目标管理
- `targets` - 获取目标列表
- `get_target_metadata` - 获取目标元数据

### 规则管理
- `rules` - 获取规则列表
- `get_alerts` - 获取告警列表

### 配置管理
- `config` - 获取配置信息
- `flags` - 获取启动参数

### 状态查询
- `status` - 获取状态信息
- `query_stats` - 获取查询统计

### 快照管理
- `snapshot` - 创建快照

### TSDB 操作
- `tsdb_stats` - 获取 TSDB 统计
- `tsdb_series` - 获取 TSDB 序列

### 存储操作
- `block_info` - 获取块信息

## Kibana 工具（52 个）

### 索引管理
- `list_indices` - 列出索引
- `get_index` - 获取索引详情
- `create_index` - 创建索引
- `delete_index` - 删除索引
- `get_index_stats` - 获取索引统计
- `get_index_settings` - 获取索引设置
- `update_index_settings` - 更新索引设置

### 文档操作
- `search_documents` - 搜索文档
- `get_document` - 获取文档
- `create_document` - 创建文档
- `update_document` - 更新文档
- `delete_document` - 删除文档
- `bulk_operations` - 批量操作

### 查询构建
- `build_query` - 构建查询
- `execute_query` - 执行查询
- `aggregate_data` - 聚合数据
- `get_query_stats` - 获取查询统计

### 可视化
- `list_visualizations` - 列出可视化
- `get_visualization` - 获取可视化
- `create_visualization` - 创建可视化
- `update_visualization` - 更新可视化
- `delete_visualization` - 删除可视化

### 仪表板
- `list_dashboards` - 列出仪表板
- `get_dashboard` - 获取仪表板
- `create_dashboard` - 创建仪表板
- `update_dashboard` - 更新仪表板
- `delete_dashboard` - 删除仪表板

### 索引模式
- `list_index_patterns` - 列出索引模式
- `get_index_pattern` - 获取索引模式
- `create_index_pattern` - 创建索引模式
- `update_index_pattern` - 更新索引模式
- `delete_index_pattern` - 删除索引模式

### 保存查询
- `list_saved_queries` - 列出保存的查询
- `get_saved_query` - 获取保存的查询
- `create_saved_query` - 创建保存的查询
- `update_saved_query` - 更新保存的查询
- `delete_saved_query` - 删除保存的查询

### 空间管理
- `list_spaces` - 列出空间
- `get_space` - 获取空间
- `create_space` - 创建空间
- `update_space` - 更新空间
- `delete_space` - 删除空间

### 发现
- `discover_data` - 发现数据
- `get_field_capabilities` - 获取字段能力

### 导出导入
- `export_objects` - 导出对象
- `import_objects` - 导入对象

### 短链接
- `create_short_url` - 创建短链接

## Elasticsearch 工具（14 个）

### 索引管理
- `list_indices` - 列出索引
- `get_index` - 获取索引
- `create_index` - 创建索引
- `delete_index` - 删除索引
- `get_index_stats` - 获取索引统计

### 文档操作
- `index_document` - 索引文档
- `get_document` - 获取文档
- `search_documents` - 搜索文档
- `update_document` - 更新文档
- `delete_document` - 删除文档

### 集群管理
- `get_cluster_health` - 获取集群健康
- `get_cluster_stats` - 获取集群统计
- `get_cluster_info` - 获取集群信息

### 别名管理
- `get_aliases` - 获取别名

## Alertmanager 工具（15 个）

### 告警管理
- `list_alerts` - 列出告警
- `get_alert` - 获取告警详情
- `get_alert_groups` - 获取告警组
- `get_silences` - 获取沉默规则
- `create_silence` - 创建沉默规则
- `delete_silence` - 删除沉默规则
- `expire_silence` - 过期沉默规则

### 规则管理
- `get_alert_rules` - 获取告警规则
- `list_rule_groups` - 列出规则组

### 配置管理
- `get_config` - 获取配置
- `get_status` - 获取状态

### 通知管理
- `list_notifications` - 列出通知
- `get_receivers` - 获取接收者配置
- `list_routes` - 列出路由

### 健康检查
- `get_health` - 获取健康状态

## Jaeger 工具（8 个）

### 追踪查询
- `get_trace` - 获取追踪
- `search_traces` - 搜索追踪
- `get_services` - 获取服务列表
- `get_operations` - 获取操作列表

### 依赖分析
- `get_dependencies` - 获取依赖关系

### 指标查询
- `get_metrics` - 获取指标

### 配置查询
- `get_config` - 获取配置
- `get_status` - 获取状态

## OpenTelemetry 工具（9 个）

### 指标管理
- `get_metrics` - 获取指标
- `get_metric_data` - 获取指标数据
- `list_metric_streams` - 列出指标流

### 追踪管理
- `get_traces` - 获取追踪
- `search_traces` - 搜索追踪

### 日志管理
- `get_logs` - 获取日志
- `search_logs` - 搜索日志

### 配置管理
- `get_config` - 获取配置
- `get_status` - 获取状态

## Utilities 工具（6 个）

### 通用工具
- `base64_encode` - Base64 编码
- `base64_decode` - Base64 解码
- `json_parse` - JSON 解析
- `json_stringify` - JSON 序列化
- `timestamp` - 获取时间戳
- `uuid` - 生成 UUID

## 工具调用示例

### Kubernetes - 列出 Pod

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

### Helm - 安装 Chart

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

### Prometheus - 查询指标

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

### Grafana - 查询仪表板

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

## 工具参数说明

所有工具都支持以下通用参数：

- `timeout` - 请求超时时间（秒）
- `dry_run` - 试运行模式，不实际执行
- `verbose` - 详细输出模式

工具特定的参数请参考各服务的详细文档。

## 错误处理

工具调用可能返回以下错误：

- `InvalidParams` - 参数无效
- `NotFound` - 资源不存在
- `PermissionDenied` - 权限不足
- `Timeout` - 请求超时
- `InternalError` - 内部错误

错误响应格式：

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
---
title: "Grafana 工具"
weight: 30
---

# Grafana 工具

Cloud Native MCP Server 提供 36 个 Grafana 管理工具，用于仪表板、数据源、文件夹、告警和用户管理。

## 仪表板管理

### list_dashboards

列出所有仪表板。

**参数**:
- `folderId` (int, optional) - 文件夹 ID
- `tag` (string, optional) - 标签过滤
- `query` (string, optional) - 搜索查询

**返回**: 仪表板列表

### get_dashboard

获取单个仪表板的详细信息。

**参数**:
- `uid` (string, required) - 仪表板 UID

**返回**: 仪表板详细信息

### create_dashboard

创建新的仪表板。

**参数**:
- `dashboard` (object, required) - 仪表板配置
- `folderId` (int, optional) - 文件夹 ID
- `overwrite` (bool, optional) - 覆盖已存在的仪表板

**返回**: 创建的仪表板信息

### update_dashboard

更新现有的仪表板。

**参数**:
- `dashboard` (object, required) - 仪表板配置
- `overwrite` (bool, optional) - 覆盖已存在的仪表板

**返回**: 更新后的仪表板信息

### delete_dashboard

删除指定的仪表板。

**参数**:
- `uid` (string, required) - 仪表板 UID

**返回**: 删除结果

### import_dashboard

导入仪表板。

**参数**:
- `dashboard` (object, required) - 仪表板配置
- `folderId` (int, optional) - 文件夹 ID
- `overwrite` (bool, optional) - 覆盖已存在的仪表板

**返回**: 导入结果

### export_dashboard

导出仪表板。

**参数**:
- `uid` (string, required) - 仪表板 UID

**返回**: 仪表板配置

### search_dashboards

搜索仪表板。

**参数**:
- `query` (string, required) - 搜索查询
- `tag` (string, optional) - 标签过滤
- `type` (string, optional) - 类型过滤（dash-db, dash-folder）

**返回**: 搜索结果

### get_dashboard_by_uid

通过 UID 获取仪表板（与 `get_dashboard` 功能相同）。

**参数**:
- `uid` (string, required) - 仪表板 UID

**返回**: 仪表板详细信息

### get_dashboard_by_tag

通过标签获取仪表板。

**参数**:
- `tag` (string, required) - 标签

**返回**: 仪表板列表

## 数据源管理

### list_datasources

列出所有数据源。

**返回**: 数据源列表

### get_datasource

获取单个数据源的详细信息。

**参数**:
- `uid` (string, required) - 数据源 UID

**返回**: 数据源详细信息

### create_datasource

创建新的数据源。

**参数**:
- `datasource` (object, required) - 数据源配置

**返回**: 创建的数据源信息

### update_datasource

更新现有的数据源。

**参数**:
- `datasource` (object, required) - 数据源配置

**返回**: 更新后的数据源信息

### delete_datasource

删除指定的数据源。

**参数**:
- `uid` (string, required) - 数据源 UID

**返回**: 删除结果

### test_datasource

测试数据源连接。

**参数**:
- `uid` (string, required) - 数据源 UID

**返回**: 测试结果

## 文件夹管理

### list_folders

列出所有文件夹。

**返回**: 文件夹列表

### get_folder

获取单个文件夹的详细信息。

**参数**:
- `uid` (string, required) - 文件夹 UID

**返回**: 文件夹详细信息

### create_folder

创建新的文件夹。

**参数**:
- `title` (string, required) - 文件夹标题
- `uid` (string, optional) - 文件夹 UID

**返回**: 创建的文件夹信息

### update_folder

更新现有的文件夹。

**参数**:
- `uid` (string, required) - 文件夹 UID
- `title` (string, required) - 新标题
- `overwrite` (bool, optional) - 覆盖已存在的文件夹

**返回**: 更新后的文件夹信息

### delete_folder

删除指定的文件夹。

**参数**:
- `uid` (string, required) - 文件夹 UID

**返回**: 删除结果

## 查询执行

### execute_query

执行查询。

**参数**:
- `query` (object, required) - 查询配置
- `datasourceUid` (string, required) - 数据源 UID

**返回**: 查询结果

### execute_multiple_queries

执行多个查询。

**参数**:
- `queries` (array, required) - 查询配置数组

**返回**: 查询结果数组

### query_metrics

查询指标数据。

**参数**:
- `query` (string, required) - 查询表达式
- `datasourceUid` (string, required) - 数据源 UID
- `start` (string, optional) - 开始时间
- `end` (string, optional) - 结束时间
- `step` (string, optional) - 步长

**返回**: 指标数据

## 告警管理

### list_alerts

列出所有告警。

**参数**:
- `dashboardId` (int, optional) - 仪表板 ID
- `panelId` (int, optional) - 面板 ID

**返回**: 告警列表

### get_alert

获取单个告警的详细信息。

**参数**:
- `alertId` (int, required) - 告警 ID

**返回**: 告警详细信息

### pause_alert

暂停告警。

**参数**:
- `alertId` (int, required) - 告警 ID

**返回**: 暂停结果

### resume_alert

恢复告警。

**参数**:
- `alertId` (int, required) - 告警 ID

**返回**: 恢复结果

### get_alert_rules

获取告警规则。

**参数**:
- `dashboardId` (int, optional) - 仪表板 ID

**返回**: 告警规则列表

## 用户管理

### list_users

列出所有用户。

**参数**:
- `page` (int, optional) - 页码
- `perPage` (int, optional) - 每页数量

**返回**: 用户列表

### get_user

获取单个用户的详细信息。

**参数**:
- `userId` (int, required) - 用户 ID

**返回**: 用户详细信息

### create_user

创建新用户。

**参数**:
- `name` (string, required) - 用户名
- `email` (string, required) - 邮箱
- `login` (string, required) - 登录名
- `password` (string, required) - 密码

**返回**: 创建的用户信息

## 组织管理

### list_organizations

列出所有组织。

**返回**: 组织列表

### get_organization

获取单个组织的详细信息。

**参数**:
- `orgId` (int, required) - 组织 ID

**返回**: 组织详细信息

## 健康检查

### get_health

获取 Grafana 服务器的健康状态。

**返回**: 健康状态

### get_version

获取 Grafana 版本信息。

**返回**: 版本信息

## 配置

Grafana 工具通过以下配置进行初始化：

```yaml
grafana:
  enabled: false

  # Grafana 服务器 URL
  url: "http://localhost:3000"

  # API Key
  apiKey: ""

  # Basic Auth
  username: ""
  password: ""

  # 请求超时（秒）
  timeoutSec: 30
```

## 最佳实践

1. **使用 API Key**: 在生产环境中使用 API Key 而不是 Basic Auth
2. **文件夹组织**: 使用文件夹组织仪表板
3. **数据源验证**: 使用 `test_datasource` 验证数据源配置
4. **告警管理**: 定期检查告警规则和状态
5. **版本控制**: 使用 `export_dashboard` 备份重要仪表板

## 相关文档

- [配置指南](/zh/guides/configuration/)
- [部署指南](/zh/guides/deployment/)
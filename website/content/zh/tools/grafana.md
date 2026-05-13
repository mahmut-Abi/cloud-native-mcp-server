---
title: "Grafana 工具"
weight: 30
---

# Grafana 工具

Grafana 服务的真实运行时工具名都以 `grafana_` 开头。  
优先参考仓库根目录的 `docs/TOOLS.md`，不要依赖旧的非前缀别名。

## 推荐起点

- `grafana_dashboards_summary`
- `grafana_datasources_summary`
- `grafana_plugins_summary`
- `grafana_test_connection`

## 常用操作

### 仪表板

- `grafana_dashboard`
- `grafana_search_dashboards`
- `grafana_update_dashboard`
- `grafana_get_dashboard_versions`
- `grafana_restore_dashboard_version`
- `grafana_delete_dashboard`

### 文件夹

- `grafana_folders`
- `grafana_folder_detail`
- `grafana_create_folder`
- `grafana_update_folder`
- `grafana_delete_folder`

### 数据源

- `grafana_datasource_detail`
- `grafana_get_datasource_by_name`
- `grafana_check_datasource_health`
- `grafana_create_datasource`
- `grafana_update_datasource`
- `grafana_delete_datasource`

### 告警与注解

- `grafana_alerts`
- `grafana_get_alert_rule_by_uid`
- `grafana_create_alert_rule`
- `grafana_update_alert_rule`
- `grafana_delete_alert_rule`
- `grafana_get_annotations`
- `grafana_create_annotation`
- `grafana_patch_annotation`

### 导航与渲染

- `grafana_generate_deeplink`
- `grafana_generate_logs_drilldown_link`
- `grafana_render_panel_image`

## 示例

```json
{
  "name": "grafana_update_dashboard",
  "arguments": {
    "dashboard": {
      "title": "Example Dashboard",
      "panels": []
    },
    "folderUID": "team-observability",
    "overwrite": true
  }
}
```

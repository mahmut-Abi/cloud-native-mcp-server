---
title: "Kibana 工具"
weight: 50
---

# Kibana 工具

Kibana 服务的真实运行时工具名都以 `kibana_` 开头。  
优先使用分页和摘要工具，不要依赖旧的无前缀别名。

## 推荐起点

- `kibana_health_summary`
- `kibana_dashboards_paginated`
- `kibana_visualizations_paginated`
- `kibana_query_logs`
- `kibana_search_saved_objects`

## 常用操作

### 空间与索引模式

- `kibana_get_spaces`
- `kibana_create_space`
- `kibana_get_index_patterns`
- `kibana_get_index_pattern`
- `kibana_create_index_pattern`
- `kibana_update_index_pattern`

### 仪表板与可视化

- `kibana_get_dashboards`
- `kibana_get_dashboard`
- `kibana_create_dashboard`
- `kibana_update_dashboard`
- `kibana_clone_dashboard`
- `kibana_get_visualizations`
- `kibana_get_visualization`

### Saved Objects 与日志分析

- `kibana_search_saved_objects_advanced`
- `kibana_create_saved_object`
- `kibana_update_saved_object`
- `kibana_query_logs`
- `kibana_get_canvas_workpads`
- `kibana_get_lens_objects`
- `kibana_get_maps`

### Alerting

- `kibana_get_alert_rules`
- `kibana_get_alert_rule`
- `kibana_create_alert_rule`
- `kibana_update_alert_rule`
- `kibana_get_connectors`

## 示例

```json
{
  "name": "kibana_query_logs",
  "arguments": {
    "indexPattern": "logs-*",
    "query": "namespace:dify AND severity:ERROR",
    "size": 20
  }
}
```

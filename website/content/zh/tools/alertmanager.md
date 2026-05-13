---
title: "Alertmanager 工具"
weight: 70
---

# Alertmanager 工具

Alertmanager 服务的高价值入口是摘要类工具，不建议一开始就取全量 alerts/silences。

## 推荐起点

- `alertmanager_alerts_summary`
- `alertmanager_silences_summary`
- `alertmanager_receivers_summary`
- `alertmanager_health_summary`

## 常用操作

- `alertmanager_get_alerts`
- `alertmanager_get_alert_groups`
- `alertmanager_query_alerts`
- `alertmanager_create_silence`
- `alertmanager_delete_silence`
- `alertmanager_get_silences`

## 示例

```json
{
  "name": "alertmanager_alerts_summary",
  "arguments": {
    "active_only": true
  }
}
```

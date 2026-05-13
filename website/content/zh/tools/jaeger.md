---
title: "Jaeger 工具"
weight: 80
---

# Jaeger 工具

Jaeger 服务的真实运行时工具名都以 `jaeger_` 开头。  
追踪搜索必须显式指定 `service`，不能直接查询“全部服务”。

## 推荐起点

- `jaeger_get_services_summary`
- `jaeger_get_traces_summary`
- `jaeger_get_trace`

## 常用操作

- `jaeger_get_services`
- `jaeger_get_service_ops`
- `jaeger_get_traces`
- `jaeger_search_traces`
- `jaeger_get_dependencies`

## 调用要点

- 先用 `jaeger_get_services_summary` 找可用服务名
- 再用 `jaeger_get_traces_summary` 做轻量发现
- 最后用 `jaeger_get_trace` 看具体 trace 详情
- 时间参数支持 RFC3339，也兼容 Unix 秒/毫秒/微秒/纳秒

## 示例

```json
{
  "name": "jaeger_get_traces_summary",
  "arguments": {
    "service": "dify",
    "start_time": "2026-05-13T00:00:00Z",
    "end_time": "2026-05-13T23:59:59Z",
    "limit": 20
  }
}
```

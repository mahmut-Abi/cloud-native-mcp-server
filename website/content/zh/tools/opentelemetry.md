---
title: "OpenTelemetry 工具"
weight: 90
---

# OpenTelemetry 工具

OpenTelemetry 服务的真实运行时工具名都以 `opentelemetry_` 开头。  
这些工具主要面向 **Collector 健康检查、配置摘要、pipeline 排障**，不是完整的通用遥测查询语言。

## 推荐起点

- `opentelemetry_get_collector_summary`
- `opentelemetry_get_config_summary`
- `opentelemetry_analyze_pipeline_status`

## 常用操作

### 指标

- `opentelemetry_get_metrics`
- `opentelemetry_query_metrics`

### 追踪

- `opentelemetry_get_traces`
- `opentelemetry_query_traces`

### 日志

- `opentelemetry_get_logs`
- `opentelemetry_query_logs`

### 采集器诊断

- `opentelemetry_get_health`
- `opentelemetry_get_status`
- `opentelemetry_get_config`
- `opentelemetry_get_config_summary`
- `opentelemetry_get_collector_summary`
- `opentelemetry_analyze_pipeline_status`

## 示例

```json
{
  "name": "opentelemetry_analyze_pipeline_status",
  "arguments": {
    "signal": "traces"
  }
}
```

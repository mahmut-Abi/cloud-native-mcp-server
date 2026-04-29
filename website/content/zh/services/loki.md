---
title: "Loki 服务"
weight: 5
---

# Loki 服务

Loki 服务提供日志聚合与 LogQL 查询能力，包含 7 个工具，适合做标签发现、日志流检查和以日志为中心的排障。

## 概述

Cloud Native MCP Server 中的 Loki 服务帮助 AI 助手和运维人员在不直接倾倒全量日志的前提下，先做紧凑摘要，再按需下钻到目标 LogQL 查询。

### 主要能力

{{< columns >}}
### 🪵 LogQL 查询
执行 Loki 即时查询和范围查询。
<--->

### 🧭 标签发现
在构造更大查询前，先发现标签、标签值和已索引的 series。
{{< /columns >}}

{{< columns >}}
### 🎯 适合 LLM 的摘要
优先返回紧凑的日志流摘要，而不是直接展开原始日志。
<--->

### ✅ 连通性检查
快速验证 Loki 地址与认证配置是否可用。
{{< /columns >}}

---

## 可用工具 (7)

- `loki_query_logs_summary`
- `loki_query`
- `loki_query_range`
- `loki_get_label_names`
- `loki_get_label_values`
- `loki_get_series`
- `loki_test_connection`

---

## 快速示例

### 先获取紧凑日志摘要

```json
{
  "method": "tools/call",
  "params": {
    "name": "loki_query_logs_summary",
    "arguments": {
      "query": "{namespace=\"prod\"} |= \"error\"",
      "limit": 50
    }
  }
}
```

### 查询某个标签的可用值

```json
{
  "method": "tools/call",
  "params": {
    "name": "loki_get_label_values",
    "arguments": {
      "label": "namespace"
    }
  }
}
```

### 在大查询前先检查 series

```json
{
  "method": "tools/call",
  "params": {
    "name": "loki_get_series",
    "arguments": {
      "matchers": ["{app=\"api\"}"]
    }
  }
}
```

---

## 最佳实践

- 先用 `loki_query_logs_summary`，再决定是否调用 `loki_query_range`。
- 尽量缩小时间窗口和 stream selector。
- 在猜标签前，先用标签发现工具。
- 优先使用 `|=`、`|~` 以及解析后的标签过滤，而不是直接拉大范围原始日志。

## 下一步

- [Prometheus 服务](/zh/services/prometheus/) 做指标关联
- [Jaeger 服务](/zh/services/jaeger/) 做链路关联
- [配置指南](/zh/docs/configuration/) 查看运行时配置

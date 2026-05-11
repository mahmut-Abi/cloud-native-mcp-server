---
title: "Langfuse 服务"
weight: 10
---

# Langfuse 服务

Langfuse 服务提供面向 LLM 可观测性与评测分析的 13 个工具，覆盖 Trace、Session、Prompt、评分和指标查询。

## 概述

Cloud Native MCP Server 中的 Langfuse 服务通过 Langfuse Public API 暴露一组只读分析能力，方便 AI 助手在同一个 MCP 接口中排查提示词执行链路、查看会话上下文、分析 observation 与评分结果。

### 主要功能

{{< columns >}}
### 🧭 Trace 发现
先用紧凑视图快速浏览 trace，再按需下钻到完整详情。
<--->

### 🧵 Session 串联
通过 session 将同一用户流程中的多条 trace 关联起来。
{{< /columns >}}

{{< columns >}}
### 📝 Prompt 分析
查看 prompt 版本、label 与解析后的内容。
<--->

### 📏 评测与指标
检查 score，并通过 metrics API 做趋势分析。
{{< /columns >}}

---

## 可用工具 (13)

### 健康检查
- **langfuse_check_health**: 检查 Langfuse API 与数据库健康状态

### Traces
- **langfuse_list_traces_summary**: 紧凑型 trace 发现视图
- **langfuse_list_traces**: 带过滤条件的完整 trace 列表
- **langfuse_get_trace**: 按 ID 获取单条 trace

### Sessions 与 Observations
- **langfuse_list_sessions**: 列出 sessions
- **langfuse_get_session**: 获取单个 session
- **langfuse_list_observations**: 列出 observations
- **langfuse_get_observation**: 获取单个 observation

### Prompts、Scores 与 Metrics
- **langfuse_list_prompts**: 列出 prompt 版本
- **langfuse_get_prompt**: 按名称获取 prompt
- **langfuse_list_scores**: 列出评分与评测结果
- **langfuse_get_score**: 获取单个 score
- **langfuse_get_metrics**: 执行 Langfuse metrics 查询

---

## 配置示例

```yaml
langfuse:
  enabled: true
  url: "https://cloud.langfuse.com"
  publicKey: "pk-lf-..."
  secretKey: "sk-lf-..."
  timeoutSec: 30
```

Langfuse Public API 使用 HTTP Basic Auth：

- 用户名：Langfuse public key
- 密码：Langfuse secret key

## 下一步

- [OpenTelemetry 服务](/zh/services/opentelemetry/) 了解基础设施遥测
- [Jaeger 服务](/zh/services/jaeger/) 了解分布式链路追踪
- [配置指南](/zh/docs/configuration/) 查看环境变量与端点配置

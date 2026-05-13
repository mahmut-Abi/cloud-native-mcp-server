---
title: "Langfuse 服务"
weight: 10
---

# Langfuse 服务

Langfuse 服务提供面向 LLM 可观测性与评测分析的 25 个工具，覆盖 Trace、Session、Prompt、评分、数据集、模型、标注队列与指标查询。

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

## 可用工具 (25)

### 健康检查
- **langfuse_check_health**: 检查 Langfuse API 与数据库健康状态

### Traces
- **langfuse_list_traces_summary**: 紧凑型 trace 发现视图
- **langfuse_list_traces**: 带过滤条件的完整 trace 列表
- **langfuse_get_trace**: 按 ID 获取单条 trace

### 标注队列与数据集
- **langfuse_list_annotation_queues**: 列出标注队列
- **langfuse_get_annotation_queue**: 获取单个标注队列
- **langfuse_list_annotation_queue_items**: 列出标注队列条目
- **langfuse_list_datasets**: 列出数据集
- **langfuse_get_dataset**: 获取单个数据集
- **langfuse_list_dataset_runs**: 列出数据集运行记录
- **langfuse_get_dataset_run**: 获取单个数据集运行记录

### Sessions 与 Observations
- **langfuse_list_sessions**: 列出 sessions
- **langfuse_get_session**: 获取单个 session
- **langfuse_list_observations**: 列出 observations
- **langfuse_get_observation**: 获取单个 observation

### 模型与评分配置
- **langfuse_list_llm_connections**: 列出 LLM 连接
- **langfuse_list_models**: 列出模型
- **langfuse_get_model**: 获取单个模型
- **langfuse_list_score_configs**: 列出评分配置
- **langfuse_get_score_config**: 获取单个评分配置

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

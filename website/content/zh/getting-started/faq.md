---
title: "FAQ"
weight: 20
description: "首次接入与生产落地中最常见问题汇总。"
---

# 快速开始 FAQ

## 应该优先选择哪种运行模式？

建议按以下优先级判断：

- `streamable-http`: 更现代的 MCP 传输方式，生产环境优先
- `sse`: 兼容性最好，适合快速接入和广泛客户端适配
- `stdio`: 本地 Agent 运行时集成
- `http`: 历史兼容模式，仅在必要时使用

如果没有明确约束，先用 `sse` 验证链路，再迁移到 `streamable-http`。

## 为什么会返回 `401 unauthorized`？

按顺序排查：

1. 运行时是否启用了 `MCP_AUTH_ENABLED=true`。
2. `MCP_AUTH_MODE` 是否为 `apikey`、`bearer`、`basic` 之一。
3. `apikey` 模式下 `MCP_AUTH_API_KEY` 是否非空且与请求一致。
4. 请求是否通过 `X-Api-Key` 头或 `api_key` 查询参数携带密钥。
5. 修改配置后是否重启了进程或容器。

## API Key 应该放在 Header 还是 Query？

两种方式都支持，生产环境建议优先使用 Header：

```bash
curl -sS -N \
  -H "X-Api-Key: ChangeMe-Strong-Key-123!" \
  http://127.0.0.1:8080/api/aggregate/sse
```

Query 参数建议只用于本地快速验证。

## 如何减少返回内容，避免模型上下文过大？

- 优先调用摘要类工具，再决定是否拉取明细。
- 大结果集优先使用分页能力。
- 通过 `MCP_DISABLED_SERVICES` 禁用暂不需要的服务。
- 配置限流，避免突发请求堆积。

## 可以只启用部分服务吗？

可以，示例如下：

```bash
export MCP_ENABLED_SERVICES="kubernetes,helm,prometheus"
export MCP_DISABLED_SERVICES="kibana,jaeger"
```

团队内建议统一策略，避免不同环境配置不一致。

## 生产环境最小安全基线是什么？

- 开启认证并定期轮换密钥。
- 妥善保护上游系统凭据（Grafana、Prometheus、Kibana 等）。
- 启用结构化日志；有审计要求时开启审计日志。
- 持续监控 `/health` 和关键服务联通性。
- 用接近真实流量的工具调用组合做压测。

## 接下来阅读什么？

- [故障排除]({{< relref "troubleshooting.md" >}})
- [安全指南]({{< relref "/docs/security.md" >}})
- [配置指南]({{< relref "/docs/configuration.md" >}})
- [性能指南]({{< relref "/docs/performance.md" >}})

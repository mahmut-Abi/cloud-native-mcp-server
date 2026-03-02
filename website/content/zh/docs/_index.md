---
title: "文档中心"
weight: 1
description: "从安装、架构、配置、部署到安全、性能、API 与工具参考的一站式文档入口。"
bookCollapseSection: true
---

# 文档中心

本章节是 Cloud Native MCP Server 的核心文档入口，用于部署、配置与生产运维。

## 按目标阅读

- [完成首次安装并启动]({{< relref "/getting-started/_index.md" >}})
- [查看首次接入 FAQ]({{< relref "/getting-started/faq.md" >}})
- [使用故障排除手册]({{< relref "/getting-started/troubleshooting.md" >}})
- [理解系统架构与请求链路]({{< relref "architecture.md" >}})
- [配置服务、认证与运行参数]({{< relref "configuration.md" >}})
- [执行生产环境部署]({{< relref "deployment.md" >}})
- [落实安全策略与加固措施]({{< relref "security.md" >}})
- [进行性能调优与压测]({{< relref "performance.md" >}})
- [查看完整工具能力清单]({{< relref "tools.md" >}})
- [对接 MCP API 端点]({{< relref "api.md" >}})

## 快速开始

```bash
docker run -d \
  --name cloud-native-mcp-server \
  -p 8080:8080 \
  -v ~/.kube:/root/.kube:ro \
  -e MCP_AUTH_ENABLED=true \
  -e MCP_AUTH_MODE=apikey \
  -e MCP_AUTH_API_KEY='ChangeMe-Strong-Key-123!' \
  mahmutabi/cloud-native-mcp-server:latest
```

启动后可访问：

- `SSE`: `http://localhost:8080/api/aggregate/sse`
- `Streamable-HTTP`: `http://localhost:8080/api/aggregate/streamable-http`

## 文档结构

- `快速开始`: 安装流程与第一条 MCP 调用
- `FAQ`: 首次接入与上线常见问题
- `故障排除`: 启动、认证、链路和服务异常排查
- `架构指南`: 组件职责与数据流模型
- `配置指南`: 全量配置项与服务集成参数
- `部署指南`: Docker、Kubernetes、Helm 生产部署策略
- `安全指南`: 认证、审计、密钥与安全实践
- `性能指南`: 缓存、并发、压测与优化建议
- `工具参考`: 全部服务工具目录与示例
- `API 文档`: 协议端点与请求/响应格式

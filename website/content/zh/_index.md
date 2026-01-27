---
title: Cloud Native MCP Server
weight: 1
---

<div align="center">

# Cloud Native MCP Server

高性能 Kubernetes 和云原生基础设施管理 MCP 服务器

[GitHub](https://github.com/mahmut-Abi/cloud-native-mcp-server) • 
[English](/#)

</div>

---

## 简介

Cloud Native MCP Server 是一个高性能的 Model Context Protocol (MCP) 服务器，用于 Kubernetes 和云原生基础设施管理。它集成了 10 个服务和 220+ 工具，让 AI 助手能够轻松管理您的云原生基础设施。

## 主要特性

- 🚀 **高性能** - LRU 缓存、JSON 编码池、智能响应限制
- 🔒 **安全可靠** - API Key、Bearer Token、Basic Auth 多种认证方式
- 📊 **全面监控** - 集成 Prometheus、Grafana、Jaeger 等监控工具
- 🔧 **灵活配置** - 支持 SSE、HTTP、stdio 多种模式
- 📝 **审计日志** - 完整的操作审计和日志记录
- 🤖 **AI 优化** - 专为 LLM 设计，包含摘要工具和分页功能

## 统计数据

| 项目 | 数量 |
|------|------|
| 集成服务 | 10 |
| MCP 工具 | 220+ |
| 运行模式 | 3 |
| 开源许可 | MIT |

## 快速开始

### Docker 部署

```bash
docker run -d \
  --name cloud-native-mcp-server \
  -p 8080:8080 \
  -v ~/.kube:/root/.kube:ro \
  mahmutabi/cloud-native-mcp-server:latest
```

### 从源码构建

```bash
git clone https://github.com/mahmut-Abi/cloud-native-mcp-server.git
cd cloud-native-mcp-server
make build
./cloud-native-mcp-server --mode=sse --addr=0.0.0.0:8080
```

## 集成服务

- **Kubernetes** - 容器编排和资源管理
- **Helm** - 应用包管理和部署
- **Grafana** - 可视化、监控仪表板和告警
- **Prometheus** - 指标收集、查询和监控
- **Kibana** - 日志分析、可视化和数据探索
- **Elasticsearch** - 日志存储、搜索和数据索引
- **Alertmanager** - 告警规则管理和通知
- **Jaeger** - 分布式追踪和性能分析
- **OpenTelemetry** - 指标、追踪和日志收集分析
- **Utilities** - 通用工具集

## 许可证

MIT License - 详见 [LICENSE](https://github.com/mahmut-Abi/cloud-native-mcp-server/blob/main/LICENSE)

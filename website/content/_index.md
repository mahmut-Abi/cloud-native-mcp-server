---
title: "Cloud Native MCP Server"
---

<div class="hero">
  <h1>Cloud Native MCP Server</h1>
  <p>高性能 Kubernetes 和云原生基础设施管理 MCP 服务器，集成 10 个服务和 220+ 工具，让 AI 助手轻松管理您的云原生基础设施</p>
  <a href="https://github.com/mahmut-Abi/cloud-native-mcp-server" class="cta-button">查看 GitHub 仓库</a>
</div>

<div class="stats-grid">
  <div class="stat-item">
    <div class="stat-number">10</div>
    <div class="stat-label">集成服务</div>
  </div>
  <div class="stat-item">
    <div class="stat-number">220+</div>
    <div class="stat-label">MCP 工具</div>
  </div>
  <div class="stat-item">
    <div class="stat-number">3</div>
    <div class="stat-label">运行模式</div>
  </div>
  <div class="stat-item">
    <div class="stat-number">100%</div>
    <div class="stat-label">开源免费</div>
  </div>
</div>

## 快速开始

### 使用二进制文件

```bash
# 下载最新版本
curl -LO https://github.com/mahmut-Abi/cloud-native-mcp-server/releases/latest/download/cloud-native-mcp-server-linux-amd64
chmod +x cloud-native-mcp-server-linux-amd64

# 以 SSE 模式运行（默认）
./cloud-native-mcp-server-linux-amd64 --mode=sse --addr=0.0.0.0:8080
```

### 使用 Docker

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

## 核心特性

<div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 1.5rem; margin-top: 2rem;">

<div class="feature-card">
  <div class="feature-icon">🚀</div>
  <h3>高性能</h3>
  <p>LRU 缓存、JSON 编码池、智能响应限制，确保最佳性能</p>
</div>

<div class="feature-card">
  <div class="feature-icon">🔒</div>
  <h3>安全可靠</h3>
  <p>API Key、Bearer Token、Basic Auth 多种认证方式，安全的密钥管理</p>
</div>

<div class="feature-card">
  <div class="feature-icon">📊</div>
  <h3>全面监控</h3>
  <p>集成 Prometheus、Grafana、Jaeger 等监控和追踪工具</p>
</div>

<div class="feature-card">
  <div class="feature-icon">🔧</div>
  <h3>灵活配置</h3>
  <p>支持 SSE、HTTP、stdio 多种模式，适配各种使用场景</p>
</div>

<div class="feature-card">
  <div class="feature-icon">📝</div>
  <h3>审计日志</h3>
  <p>完整的操作审计和日志记录，支持多种存储方式</p>
</div>

<div class="feature-card">
  <div class="feature-icon">🤖</div>
  <h3>AI 优化</h3>
  <p>专为 LLM 设计，包含摘要工具和分页功能，防止上下文溢出</p>
</div>

</div>

## 服务概览

| 服务 | 工具数量 | 描述 |
|------|---------|------|
| **Kubernetes** | 28 | 容器编排和资源管理 |
| **Helm** | 31 | 应用包管理和部署 |
| **Grafana** | 36 | 可视化、监控仪表板和告警 |
| **Prometheus** | 20 | 指标收集、查询和监控 |
| **Kibana** | 52 | 日志分析、可视化和数据探索 |
| **Elasticsearch** | 14 | 日志存储、搜索和数据索引 |
| **Alertmanager** | 15 | 告警规则管理和通知 |
| **Jaeger** | 8 | 分布式追踪和性能分析 |
| **OpenTelemetry** | 9 | 指标、追踪和日志收集分析 |
| **Utilities** | 6 | 通用工具集 |

**总计：220+ 工具**

## API 端点

### SSE 模式

| 端点 | 描述 |
|------|------|
| `/api/aggregate/sse` | 所有服务（推荐） |
| `/api/kubernetes/sse` | Kubernetes 服务 |
| `/api/helm/sse` | Helm 服务 |
| `/api/grafana/sse` | Grafana 服务 |
| `/api/prometheus/sse` | Prometheus 服务 |
| `/api/kibana/sse` | Kibana 服务 |
| `/api/elasticsearch/sse` | Elasticsearch 服务 |
| `/api/alertmanager/sse` | Alertmanager 服务 |
| `/api/jaeger/sse` | Jaeger 服务 |
| `/api/opentelemetry/sse` | OpenTelemetry 服务 |
| `/api/utilities/sse` | Utilities 服务 |

### HTTP 模式

将上述端点中的 `/sse` 替换为 `/http` 即可。

## 文档

- [完整工具参考](/docs/tools/) - 所有 220+ 工具的详细文档
- [配置指南](/docs/configuration/) - 配置选项和示例
- [部署指南](/docs/deployment/) - 部署策略和最佳实践
- [安全指南](/docs/security/) - 认证、密钥管理和安全最佳实践
- [架构指南](/docs/architecture/) - 系统架构和设计
- [性能指南](/docs/performance/) - 性能特性和调优

## 构建

```bash
# 构建当前平台
make build

# 运行测试
make test

# 代码检查
make lint

# Docker 构建
make docker-build
```

## 许可证

MIT License - 详见 [LICENSE](https://github.com/mahmut-Abi/cloud-native-mcp-server/blob/main/LICENSE)
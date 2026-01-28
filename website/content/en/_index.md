---
title: Cloud Native MCP Server
weight: 1
---

<div class="hero">
  <h1>Cloud Native MCP Server</h1>
  <p>高性能量子 MCP 服务器，专为 Kubernetes 和云原生基础设施管理而设计</p>
  <div class="hero-buttons">
    <a href="https://github.com/mahmut-Abi/cloud-native-mcp-server" class="cta-button">GitHub 仓库</a>
    <a href="#quick-start" class="cta-button" style="background: transparent; border: 2px solid white; margin-left: 1rem;">快速开始</a>
  </div>
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/mahmut-Abi/cloud-native-mcp-server)](https://goreportcard.com/report/github.com/mahmut-Abi/cloud-native-mcp-server)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org)

## 核心特性

{{< columns >}}
### 🚀 高性能
LRU 缓存、JSON 编码池、智能响应限制，确保最优性能表现
<--->

### 🔒 安全可靠
API Key、Bearer Token、Basic Auth 多重认证，保障系统安全
{{< /columns >}}

{{< columns >}}
### 📊 全面监控
原生集成 Prometheus、Grafana、Jaeger 等云原生监控工具
<--->

### 🤖 AI 优化
专为 LLM 设计，包含摘要工具和分页功能，防止上下文溢出
{{< /columns >}}

---

## 项目统计

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
    <div class="stat-number">MIT</div>
    <div class="stat-label">开源许可</div>
  </div>
</div>

---

## 集成服务

<div class="service-grid">
  <div class="service-card">
    <h3> Kubernetes <span class="tool-count">28 工具</span></h3>
    <p>核心容器编排和资源管理</p>
  </div>
  <div class="service-card">
    <h3> Helm <span class="tool-count">31 工具</span></h3>
    <p>应用包管理与部署</p>
  </div>
  <div class="service-card">
    <h3> Grafana <span class="tool-count">36 工具</span></h3>
    <p>可视化、监控仪表板和告警</p>
  </div>
  <div class="service-card">
    <h3> Prometheus <span class="tool-count">20 工具</span></h3>
    <p>指标收集、查询和监控</p>
  </div>
  <div class="service-card">
    <h3> Kibana <span class="tool-count">52 工具</span></h3>
    <p>日志分析、可视化和数据探索</p>
  </div>
  <div class="service-card">
    <h3> Elasticsearch <span class="tool-count">14 工具</span></h3>
    <p>日志存储、搜索和数据索引</p>
  </div>
  <div class="service-card">
    <h3> Alertmanager <span class="tool-count">15 工具</span></h3>
    <p>告警规则管理和通知</p>
  </div>
  <div class="service-card">
    <h3> Jaeger <span class="tool-count">8 工具</span></h3>
    <p>分布式追踪和性能分析</p>
  </div>
  <div class="service-card">
    <h3> OpenTelemetry <span class="tool-count">9 工具</span></h3>
    <p>指标、追踪和日志收集分析</p>
  </div>
  <div class="service-card">
    <h3> Utilities <span class="tool-count">6 工具</span></h3>
    <p>通用工具集</p>
  </div>
</div>

---

## <span id="quick-start">快速开始</span>

{{< tabs >}}
{{< tab "Docker" >}}
### Docker 部署

```bash
docker run -d \
  --name cloud-native-mcp-server \
  -p 8080:8080 \
  -v ~/.kube:/root/.kube:ro \
  mahmutabi/cloud-native-mcp-server:latest
```
{{< /tab >}}

{{< tab "Binary" >}}
### 二进制部署

```bash
# 下载最新版本
curl -LO https://github.com/mahmut-Abi/cloud-native-mcp-server/releases/latest/download/cloud-native-mcp-server-linux-amd64
chmod +x cloud-native-mcp-server-linux-amd64

# 运行服务
./cloud-native-mcp-server --mode=sse --addr=0.0.0.0:8080
```
{{< /tab >}}

{{< tab "Source" >}}
### 源码构建

```bash
git clone https://github.com/mahmut-Abi/cloud-native-mcp-server.git
cd cloud-native-mcp-server
make build
./cloud-native-mcp-server --mode=sse --addr=0.0.0.0:8080
```
{{< /tab >}}
{{< /tabs >}}

---

## 文档导航

- [快速开始](/zh/getting-started/) - 快速部署和使用
- [核心概念](/zh/concepts/architecture/overview/) - 了解系统架构和设计原理
- [服务概览](/zh/services/) - 探索 10 个集成服务
- [配置指南](/zh/guides/configuration/server/) - 详细配置选项和示例
- [部署指南](/zh/guides/deployment/kubernetes/) - 部署策略和最佳实践
- [安全指南](/zh/guides/security/best-practices/) - 认证、密钥管理和安全最佳实践
- [性能指南](/zh/guides/performance/optimization/) - 性能特性与优化
- [API 文档](/zh/docs/api/) - 完整的 API 参考
- [工具参考](/zh/docs/tools/) - 所有 220+ 工具的详细文档
- [网站地图](/zh/sitemap/) - 完整的网站导航

---

## 更多资源

- [博客](/zh/posts/) - 最新新闻、更新和教程
- [案例展示](/zh/showcase/) - 真实世界用例和用户评价
- [GitHub 仓库](https://github.com/mahmut-Abi/cloud-native-mcp-server) - 源代码和问题跟踪

---

## 开源贡献

Cloud Native MCP Server 是一个开源项目，欢迎提交 Issue 和 Pull Request 来改进项目。

**许可证**: MIT License - 详见 [LICENSE](https://github.com/mahmut-Abi/cloud-native-mcp-server/blob/main/LICENSE)

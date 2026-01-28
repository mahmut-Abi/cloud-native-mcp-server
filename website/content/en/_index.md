---
title: Cloud Native MCP Server
weight: 1
---

<div class="hero">
  <h1>Cloud Native MCP Server</h1>
  <p>High-performance quantum MCP server designed for Kubernetes and cloud-native infrastructure management</p>
  <div class="hero-buttons">
    <a href="https://github.com/mahmut-Abi/cloud-native-mcp-server" class="cta-button">GitHub Repository</a>
    <a href="#quick-start" class="cta-button" style="background: transparent; border: 2px solid white; margin-left: 1rem;">Quick Start</a>
  </div>
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/mahmut-Abi/cloud-native-mcp-server)](https://goreportcard.com/report/github.com/mahmut-Abi/cloud-native-mcp-server)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org)

## Core Features

{{< columns >}}
### 🚀 High Performance
LRU caching, JSON encoding pools, intelligent response limiting ensure optimal performance
--->

### 🔒 Secure & Reliable
API Key, Bearer Token, Basic Auth multiple authentication methods ensure system security
{{< /columns >}}

{{< columns >}}
### 📊 Comprehensive Monitoring
Native integration with Prometheus, Grafana, Jaeger and other cloud-native monitoring tools
--->

### 🤖 AI Optimized
Designed specifically for LLM with summary tools and pagination to prevent context overflow
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

## Documentation Navigation

- [Getting Started](/en/getting-started/) - Quick deployment and usage
- [Core Concepts](/en/concepts/architecture/overview/) - Understand system architecture and design principles
- [Service Overview](/en/services/) - Explore 10 integrated services
- [Configuration Guide](/en/guides/configuration/server/) - Detailed configuration options and examples
- [Deployment Guide](/en/guides/deployment/kubernetes/) - Deployment strategies and best practices
- [Security Guide](/en/guides/security/best-practices/) - Authentication, key management and security best practices
- [Performance Guide](/en/guides/performance/optimization/) - Performance features and optimization
- [API Documentation](/en/docs/api/) - Complete API reference
- [Tools Reference](/en/docs/tools/) - Detailed documentation for all 220+ tools
- [Site Map](/en/sitemap/) - Complete site navigation

---

## Additional Resources

- [Blog](/en/posts/) - Latest news, updates and tutorials
- [Case Studies](/en/showcase/) - Real-world use cases and user testimonials
- [GitHub Repository](https://github.com/mahmut-Abi/cloud-native-mcp-server) - Source code and issue tracking

---

## 开源贡献

Cloud Native MCP Server 是一个开源项目，欢迎提交 Issue 和 Pull Request 来改进项目。

**许可证**: MIT License - 详见 [LICENSE](https://github.com/mahmut-Abi/cloud-native-mcp-server/blob/main/LICENSE)

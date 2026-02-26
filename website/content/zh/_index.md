---
title: Cloud Native MCP Server
weight: 1
description: 面向 Kubernetes 与云原生运维场景的高性能 Model Context Protocol 服务器。
---

<div class="hero">
  <h1>Cloud Native MCP Server</h1>
  <p>面向 Kubernetes 与云原生基础设施管理的生产级 MCP 服务器。</p>
  <div class="hero-buttons">
    <a href="https://github.com/mahmut-Abi/cloud-native-mcp-server" class="cta-button"><span>GitHub 仓库</span></a>
    <a href="#quick-start" class="cta-button transparent"><span>快速开始</span></a>
  </div>
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/mahmut-Abi/cloud-native-mcp-server)](https://goreportcard.com/report/github.com/mahmut-Abi/cloud-native-mcp-server)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org)

## 核心特性

{{< columns >}}
<h3>高性能执行</h3>
<p>LRU 缓存、JSON 编码池与响应裁剪策略，确保稳定吞吐与低延迟。</p>
<--->

<h3>安全与审计</h3>
<p>支持 API Key、Bearer、Basic 认证，并提供审计日志能力。</p>
{{< /columns >}}

{{< columns >}}
<h3>可观测性集成</h3>
<p>原生对接 Prometheus、Grafana、Jaeger、OpenTelemetry。</p>
<--->

<h3>面向 AI 交互优化</h3>
<p>为 LLM 设计的摘要工具与分页策略，降低上下文溢出风险。</p>
{{< /columns >}}

---

## 项目概览

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
    <div class="stat-number">4</div>
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
    <h3>Kubernetes <span class="tool-count">28 工具</span></h3>
    <p>容器编排、资源查询与集群运维任务。</p>
  </div>
  <div class="service-card">
    <h3>Helm <span class="tool-count">31 工具</span></h3>
    <p>应用包生命周期管理与发布流程支持。</p>
  </div>
  <div class="service-card">
    <h3>Grafana <span class="tool-count">36 工具</span></h3>
    <p>仪表盘、可视化与告警配置管理。</p>
  </div>
  <div class="service-card">
    <h3>Prometheus <span class="tool-count">20 工具</span></h3>
    <p>指标查询、规则排查与监控数据运维。</p>
  </div>
  <div class="service-card">
    <h3>Kibana <span class="tool-count">52 工具</span></h3>
    <p>日志分析与 Elastic 生态可观测性探索。</p>
  </div>
  <div class="service-card">
    <h3>Elasticsearch <span class="tool-count">14 工具</span></h3>
    <p>索引管理、搜索调试与集群检查。</p>
  </div>
  <div class="service-card">
    <h3>Alertmanager <span class="tool-count">15 工具</span></h3>
    <p>告警路由、静默策略与通知链路管理。</p>
  </div>
  <div class="service-card">
    <h3>Jaeger <span class="tool-count">8 工具</span></h3>
    <p>分布式链路追踪与延迟路径分析。</p>
  </div>
  <div class="service-card">
    <h3>OpenTelemetry <span class="tool-count">9 工具</span></h3>
    <p>指标、日志、追踪的采集链路检查。</p>
  </div>
  <div class="service-card">
    <h3>Utilities <span class="tool-count">6 工具</span></h3>
    <p>通用辅助工具，覆盖日常运维场景。</p>
  </div>
</div>

---

## <span id="quick-start">快速开始</span>

{{< tabs >}}
{{< tab "Docker" >}}
{{< highlight bash >}}
docker run -d \
  --name cloud-native-mcp-server \
  -p 8080:8080 \
  -v ~/.kube:/root/.kube:ro \
  mahmutabi/cloud-native-mcp-server:latest
{{< /highlight >}}
{{< /tab >}}

{{< tab "二进制" >}}
{{< highlight bash >}}
curl -LO https://github.com/mahmut-Abi/cloud-native-mcp-server/releases/latest/download/cloud-native-mcp-server-linux-amd64
chmod +x cloud-native-mcp-server-linux-amd64
./cloud-native-mcp-server-linux-amd64 --mode=sse --addr=0.0.0.0:8080
{{< /highlight >}}
{{< /tab >}}

{{< tab "源码" >}}
{{< highlight bash >}}
git clone https://github.com/mahmut-Abi/cloud-native-mcp-server.git
cd cloud-native-mcp-server
make build
./cloud-native-mcp-server --mode=sse --addr=0.0.0.0:8080
{{< /highlight >}}
{{< /tab >}}
{{< /tabs >}}

---

## 文档导航

- [快速开始]({{< relref "getting-started/_index.md" >}})
- [架构设计]({{< relref "docs/architecture.md" >}})
- [配置说明]({{< relref "docs/configuration.md" >}})
- [部署指南]({{< relref "docs/deployment.md" >}})
- [安全指南]({{< relref "docs/security.md" >}})
- [性能指南]({{< relref "docs/performance.md" >}})
- [工具参考]({{< relref "docs/tools.md" >}})
- [服务概览]({{< relref "services/_index.md" >}})
- [站点地图]({{< relref "sitemap.md" >}})

---

## 更多资源

- [博客]({{< relref "posts/_index.md" >}})
- [案例展示]({{< relref "showcase.md" >}})
- [GitHub 仓库](https://github.com/mahmut-Abi/cloud-native-mcp-server)

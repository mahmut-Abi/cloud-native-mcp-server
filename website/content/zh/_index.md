---
title: Cloud Native MCP Server
weight: 1
description: 面向 Kubernetes 与云原生运维场景的高性能 Model Context Protocol 服务器。
---

<div class="hero">
  <h1>Cloud Native MCP Server</h1>
  <p>面向 Kubernetes 与云原生基础设施管理的生产级 MCP 服务器，聚合 11 个服务与 250+ 工具，支持 SSE / Streamable-HTTP 两种运行模式。</p>
  <div class="hero-buttons">
    <a href="https://github.com/mahmut-Abi/cloud-native-mcp-server" class="cta-button"><span>GitHub 仓库</span></a>
    <a href="#quick-start" class="cta-button transparent"><span>快速开始</span></a>
  </div>
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/mahmut-Abi/cloud-native-mcp-server)](https://goreportcard.com/report/github.com/mahmut-Abi/cloud-native-mcp-server)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org)

## 核心价值

<div class="value-grid">
  <article class="value-card">
    <h3>统一运维入口</h3>
    <p>将 Kubernetes、Helm、Grafana、Prometheus、Kibana 等能力收敛到单一 MCP 接口，减少跨系统切换成本。</p>
  </article>
  <article class="value-card">
    <h3>面向生产的安全策略</h3>
    <p>支持 apikey / bearer / basic 认证、速率限制与审计日志，便于在企业环境中落地安全合规要求。</p>
  </article>
  <article class="value-card">
    <h3>对 AI Agent 友好</h3>
    <p>内置分页与摘要能力，降低上下文爆炸风险，让大模型在复杂排障中保持稳定输出。</p>
  </article>
</div>

---

## 典型使用场景

<div class="usecase-grid">
  <article class="usecase-card">
    <h3>故障定位与快速止血</h3>
    <p>聚合查看 Pod 状态、事件、日志与监控指标，缩短从告警到定位根因的路径。</p>
  </article>
  <article class="usecase-card">
    <h3>发布与变更管控</h3>
    <p>通过 Helm 与 Kubernetes 工具链执行发布、回滚、扩缩容，并保留审计痕迹。</p>
  </article>
  <article class="usecase-card">
    <h3>可观测性协同分析</h3>
    <p>跨 Prometheus、Grafana、Loki、Jaeger、OpenTelemetry 联动，串联指标、日志与链路数据。</p>
  </article>
</div>

---

## 项目概览

<div class="stats-grid">
  <div class="stat-item">
    <div class="stat-number">11</div>
    <div class="stat-label">集成服务</div>
  </div>
  <div class="stat-item">
    <div class="stat-number">250+</div>
    <div class="stat-label">MCP 工具</div>
  </div>
  <div class="stat-item">
    <div class="stat-number">2</div>
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
    <h3>Kubernetes <span class="tool-count">34 工具</span></h3>
    <p>容器编排、资源查询与集群运维任务。</p>
  </div>
  <div class="service-card">
    <h3>Helm <span class="tool-count">34 工具</span></h3>
    <p>应用包生命周期管理与发布流程支持。</p>
  </div>
  <div class="service-card">
    <h3>Grafana <span class="tool-count">43 工具</span></h3>
    <p>仪表盘、可视化与告警配置管理。</p>
  </div>
  <div class="service-card">
    <h3>Prometheus <span class="tool-count">20 工具</span></h3>
    <p>指标查询、规则排查与监控数据运维。</p>
  </div>
  <div class="service-card">
    <h3>Loki <span class="tool-count">7 工具</span></h3>
    <p>面向日志优先排障的 LogQL 查询、标签发现与日志流检查。</p>
  </div>
  <div class="service-card">
    <h3>Kibana <span class="tool-count">73 工具</span></h3>
    <p>日志分析与 Elastic 生态可观测性探索。</p>
  </div>
  <div class="service-card">
    <h3>Elasticsearch <span class="tool-count">12 工具</span></h3>
    <p>索引管理、搜索调试与集群检查。</p>
  </div>
  <div class="service-card">
    <h3>Alertmanager <span class="tool-count">16 工具</span></h3>
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
  -e MCP_AUTH_ENABLED=true \
  -e MCP_AUTH_MODE=apikey \
  -e MCP_AUTH_API_KEY='ChangeMe-Strong-Key-123!' \
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
./cloud-native-mcp-server --mode=streamable-http --addr=0.0.0.0:8080
{{< /highlight >}}
{{< /tab >}}
{{< /tabs >}}

### 可用性验证

{{< highlight bash >}}
# 1) 健康检查
curl -sS http://127.0.0.1:8080/health

# 2) SSE 握手与 initialize 全链路验证（在仓库根目录执行）
make sse-smoke BASE_URL=http://127.0.0.1:8080
{{< /highlight >}}

### 常用入口

- SSE 聚合端点（`--mode=sse`）: `http://127.0.0.1:8080/api/aggregate/sse`
- Streamable-HTTP 聚合端点（`--mode=streamable-http`）: `http://127.0.0.1:8080/api/aggregate/streamable-http`
- 健康检查: `http://127.0.0.1:8080/health`

---

## 上线前检查清单

<div class="ops-grid">
  <article class="ops-card">
    <h3>认证与权限</h3>
    <ul>
      <li>启用 `MCP_AUTH_ENABLED=true`。</li>
      <li>使用 `apikey` / `bearer` / `basic` 模式之一。</li>
      <li>最小化 Kubernetes 与第三方系统权限范围。</li>
    </ul>
  </article>
  <article class="ops-card">
    <h3>可观测性与审计</h3>
    <ul>
      <li>开启结构化日志与必要的指标采集。</li>
      <li>根据审计要求启用审计日志与存储策略。</li>
      <li>确认 `/health` 与关键服务检查稳定。</li>
    </ul>
  </article>
  <article class="ops-card">
    <h3>性能与稳定性</h3>
    <ul>
      <li>按业务压力调优限流、超时与并发参数。</li>
      <li>优先使用分页与摘要工具降低上下文体积。</li>
      <li>压测时覆盖高峰期真实工具调用组合。</li>
    </ul>
  </article>
</div>

---

## 文档导航

- [快速开始]({{< relref "getting-started/_index.md" >}})
- [快速开始 FAQ]({{< relref "getting-started/faq.md" >}})
- [故障排除]({{< relref "getting-started/troubleshooting.md" >}})
- [架构设计]({{< relref "docs/architecture.md" >}})
- [配置说明]({{< relref "docs/configuration.md" >}})
- [部署指南]({{< relref "docs/deployment.md" >}})
- [安全指南]({{< relref "docs/security.md" >}})
- [性能指南]({{< relref "docs/performance.md" >}})
- [工具参考]({{< relref "docs/tools.md" >}})
- [服务概览]({{< relref "services/_index.md" >}})
- [站点地图]({{< relref "sitemap.md" >}})

---

## 常见问题与排障入口

<div class="resource-grid">
  <article class="resource-card">
    <h3>快速开始 FAQ</h3>
    <p>覆盖认证模式、运行模式选择、连接方式和上线建议，适合首次接入团队快速统一认知。</p>
    <a class="resource-link" href='{{< relref "getting-started/faq.md" >}}'>查看 FAQ</a>
  </article>
  <article class="resource-card">
    <h3>故障排除手册</h3>
    <p>提供启动失败、401、SSE 连接异常、服务不可用等高频问题的排查路径和命令清单。</p>
    <a class="resource-link" href='{{< relref "getting-started/troubleshooting.md" >}}'>进入排障</a>
  </article>
</div>

---

## 更多资源

- [博客]({{< relref "posts/_index.md" >}})
- [案例展示]({{< relref "showcase.md" >}})
- [GitHub 仓库](https://github.com/mahmut-Abi/cloud-native-mcp-server)

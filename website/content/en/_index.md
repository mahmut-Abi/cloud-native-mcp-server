---
title: Cloud Native MCP Server
weight: 1
description: High-performance Model Context Protocol server for Kubernetes and cloud-native operations.
---

<div class="hero">
  <h1>Cloud Native MCP Server</h1>
  <p>A production-focused MCP server for Kubernetes and cloud-native infrastructure management.</p>
  <div class="hero-buttons">
    <a href="https://github.com/mahmut-Abi/cloud-native-mcp-server" class="cta-button"><span>GitHub Repository</span></a>
    <a href="#quick-start" class="cta-button transparent"><span>Quick Start</span></a>
  </div>
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/mahmut-Abi/cloud-native-mcp-server)](https://goreportcard.com/report/github.com/mahmut-Abi/cloud-native-mcp-server)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org)

## Core Features

{{< columns >}}
<h3>High Performance</h3>
<p>LRU caching, JSON encoding pools, and response shaping designed for predictable low-latency behavior.</p>
<--->

<h3>Security First</h3>
<p>API key, bearer token, and basic auth support with audit logging and access control options.</p>
{{< /columns >}}

{{< columns >}}
<h3>Observability Built In</h3>
<p>Native integrations with Prometheus, Grafana, Jaeger, and OpenTelemetry.</p>
<--->

<h3>LLM-Friendly Output</h3>
<p>Summary tools and pagination patterns that help agents avoid context overload.</p>
{{< /columns >}}

---

## Project Snapshot

<div class="stats-grid">
  <div class="stat-item">
    <div class="stat-number">10</div>
    <div class="stat-label">Integrated Services</div>
  </div>
  <div class="stat-item">
    <div class="stat-number">220+</div>
    <div class="stat-label">MCP Tools</div>
  </div>
  <div class="stat-item">
    <div class="stat-number">4</div>
    <div class="stat-label">Run Modes</div>
  </div>
  <div class="stat-item">
    <div class="stat-number">MIT</div>
    <div class="stat-label">License</div>
  </div>
</div>

---

## Integrated Services

<div class="service-grid">
  <div class="service-card">
    <h3>Kubernetes <span class="tool-count">28 tools</span></h3>
    <p>Core orchestration and resource management workflows.</p>
  </div>
  <div class="service-card">
    <h3>Helm <span class="tool-count">31 tools</span></h3>
    <p>Package lifecycle and release operations for Kubernetes apps.</p>
  </div>
  <div class="service-card">
    <h3>Grafana <span class="tool-count">36 tools</span></h3>
    <p>Dashboards, alerting, and visualization management.</p>
  </div>
  <div class="service-card">
    <h3>Prometheus <span class="tool-count">20 tools</span></h3>
    <p>Metrics querying, rules inspection, and monitoring workflows.</p>
  </div>
  <div class="service-card">
    <h3>Kibana <span class="tool-count">52 tools</span></h3>
    <p>Log exploration and analytics for Elastic-based observability.</p>
  </div>
  <div class="service-card">
    <h3>Elasticsearch <span class="tool-count">14 tools</span></h3>
    <p>Index inspection, search, and cluster operation support.</p>
  </div>
  <div class="service-card">
    <h3>Alertmanager <span class="tool-count">15 tools</span></h3>
    <p>Alert routing, silence management, and incident visibility.</p>
  </div>
  <div class="service-card">
    <h3>Jaeger <span class="tool-count">8 tools</span></h3>
    <p>Distributed tracing and request-path diagnostics.</p>
  </div>
  <div class="service-card">
    <h3>OpenTelemetry <span class="tool-count">9 tools</span></h3>
    <p>Telemetry pipeline checks for traces, logs, and metrics.</p>
  </div>
  <div class="service-card">
    <h3>Utilities <span class="tool-count">6 tools</span></h3>
    <p>General-purpose helpers for day-to-day operational tasks.</p>
  </div>
</div>

---

## <span id="quick-start">Quick Start</span>

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

{{< tab "Binary" >}}
{{< highlight bash >}}
curl -LO https://github.com/mahmut-Abi/cloud-native-mcp-server/releases/latest/download/cloud-native-mcp-server-linux-amd64
chmod +x cloud-native-mcp-server-linux-amd64
./cloud-native-mcp-server-linux-amd64 --mode=sse --addr=0.0.0.0:8080
{{< /highlight >}}
{{< /tab >}}

{{< tab "Source" >}}
{{< highlight bash >}}
git clone https://github.com/mahmut-Abi/cloud-native-mcp-server.git
cd cloud-native-mcp-server
make build
./cloud-native-mcp-server --mode=sse --addr=0.0.0.0:8080
{{< /highlight >}}
{{< /tab >}}
{{< /tabs >}}

---

## Documentation Map

- [Getting Started]({{< relref "getting-started/_index.md" >}})
- [Architecture]({{< relref "docs/architecture.md" >}})
- [Configuration]({{< relref "docs/configuration.md" >}})
- [Deployment]({{< relref "docs/deployment.md" >}})
- [Security]({{< relref "docs/security.md" >}})
- [Performance]({{< relref "docs/performance.md" >}})
- [Tools Reference]({{< relref "docs/tools.md" >}})
- [Services Overview]({{< relref "services/_index.md" >}})
- [Sitemap]({{< relref "sitemap.md" >}})

---

## More Resources

- [Blog]({{< relref "posts/_index.md" >}})
- [Showcase]({{< relref "showcase.md" >}})
- [GitHub Repository](https://github.com/mahmut-Abi/cloud-native-mcp-server)

---
title: Cloud Native MCP Server
weight: 1
---

<div class="hero">
  <h1>Cloud Native MCP Server</h1>
  <p>High-performance Model Context Protocol (MCP) server for Kubernetes and cloud-native infrastructure management</p>
  <div class="hero-buttons">
    <a href="https://github.com/mahmut-Abi/cloud-native-mcp-server" class="cta-button">GitHub Repository</a>
    <a href="#quick-start" class="cta-button" style="background: transparent; border: 2px solid white; margin-left: 1rem;">Quick Start</a>
  </div>
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/mahmut-Abi/cloud-native-mcp-server)](https://goreportcard.com/report/github.com/mahmut-Abi/cloud-native-mcp-server)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org)

## Key Features

{{< columns >}}
### 🚀 High Performance
LRU cache, JSON encoding pool, intelligent response limiting for optimal performance
<--->

### 🔒 Secure & Reliable
API Key, Bearer Token, Basic Auth multiple authentication methods for security
{{< /columns >}}

{{< columns >}}
### 📊 Comprehensive Monitoring
Natively integrated with Prometheus, Grafana, Jaeger and other cloud-native tools
<--->

### 🤖 AI Optimized
Designed for LLMs with summary tools and pagination to prevent context overflow
{{< /columns >}}

---

## Project Statistics

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
    <div class="stat-number">3</div>
    <div class="stat-label">Running Modes</div>
  </div>
  <div class="stat-item">
    <div class="stat-number">MIT</div>
    <div class="stat-label">Open Source License</div>
  </div>
</div>

---

## Integrated Services

<div class="service-grid">
  <div class="service-card">
    <h3> Kubernetes <span class="tool-count">28 tools</span></h3>
    <p>Core container orchestration and resource management</p>
  </div>
  <div class="service-card">
    <h3> Helm <span class="tool-count">31 tools</span></h3>
    <p>Application package management and deployment</p>
  </div>
  <div class="service-card">
    <h3> Grafana <span class="tool-count">36 tools</span></h3>
    <p>Visualization, monitoring dashboards, and alerting</p>
  </div>
  <div class="service-card">
    <h3> Prometheus <span class="tool-count">20 tools</span></h3>
    <p>Metrics collection, querying, and monitoring</p>
  </div>
  <div class="service-card">
    <h3> Kibana <span class="tool-count">52 tools</span></h3>
    <p>Log analysis, visualization, and data exploration</p>
  </div>
  <div class="service-card">
    <h3> Elasticsearch <span class="tool-count">14 tools</span></h3>
    <p>Log storage, search, and data indexing</p>
  </div>
  <div class="service-card">
    <h3> Alertmanager <span class="tool-count">15 tools</span></h3>
    <p>Alert rules management and notifications</p>
  </div>
  <div class="service-card">
    <h3> Jaeger <span class="tool-count">8 tools</span></h3>
    <p>Distributed tracing and performance analysis</p>
  </div>
  <div class="service-card">
    <h3> OpenTelemetry <span class="tool-count">9 tools</span></h3>
    <p>Metrics, traces, and logs collection and analysis</p>
  </div>
  <div class="service-card">
    <h3> Utilities <span class="tool-count">6 tools</span></h3>
    <p>General-purpose utility tools</p>
  </div>
</div>

---

## <span id="quick-start">Quick Start</span>

{{< tabs >}}
{{< tab "Docker" >}}
### Docker Deployment

```bash
docker run -d \
  --name cloud-native-mcp-server \
  -p 8080:8080 \
  -v ~/.kube:/root/.kube:ro \
  mahmutabi/cloud-native-mcp-server:latest
```
{{< /tab >}}

{{< tab "Binary" >}}
### Binary Deployment

```bash
# Download the latest release
curl -LO https://github.com/mahmut-Abi/cloud-native-mcp-server/releases/latest/download/cloud-native-mcp-server-linux-amd64
chmod +x cloud-native-mcp-server-linux-amd64

# Run the service
./cloud-native-mcp-server --mode=sse --addr=0.0.0.0:8080
```
{{< /tab >}}

{{< tab "Source" >}}
### Build from Source

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

- [Quick Start](/en/getting-started/) - Quick deployment and usage
- [Core Concepts](/en/concepts/architecture/overview/) - Understand system architecture and design principles
- [Service Overview](/en/services/) - Explore the 10 integrated services
- [Configuration Guide](/en/guides/configuration/server/) - Detailed configuration options and examples
- [Deployment Guide](/en/guides/deployment/kubernetes/) - Deployment strategies and best practices
- [Security Guide](/en/guides/security/best-practices/) - Authentication, secrets management, and security best practices
- [Performance Guide](/en/guides/performance/optimization/) - Performance features and tuning
- [API Documentation](/en/docs/api/) - Complete API reference
- [Tools Reference](/en/docs/tools/) - Detailed documentation for all 220+ tools
- [Sitemap](/en/sitemap/) - Complete website navigation

---

## Additional Resources

- [Blog](/en/posts/) - Latest news, updates, and tutorials
- [Showcase](/en/showcase/) - Real-world use cases and testimonials
- [GitHub Repository](https://github.com/mahmut-Abi/cloud-native-mcp-server) - Source code and issue tracking

---

## Open Source Contribution

Cloud Native MCP Server is an open source project. Contributions are welcome via Issues and Pull Requests.

**License**: MIT License - see [LICENSE](https://github.com/mahmut-Abi/cloud-native-mcp-server/blob/main/LICENSE) for details

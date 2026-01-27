---
title: "Home"
---

<div class="hero">
  <h1>Cloud Native MCP Server</h1>
  <p>A high-performance Model Context Protocol (MCP) server for Kubernetes and cloud-native infrastructure management with 10 integrated services and 220+ tools, enabling AI assistants to effortlessly manage your cloud-native infrastructure</p>
  <a href="https://github.com/mahmut-Abi/cloud-native-mcp-server" class="cta-button">View on GitHub</a>
</div>

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
    <div class="stat-number">100%</div>
    <div class="stat-label">Open Source</div>
  </div>
</div>

## Quick Start

```bash
docker run -d \
  --name cloud-native-mcp-server \
  -p 8080:8080 \
  -v ~/.kube:/root/.kube:ro \
  mahmutabi/cloud-native-mcp-server:latest
```

## Key Features

<div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 1.5rem; margin-top: 2rem;">

<div class="feature-card">
  <div class="feature-icon">🚀</div>
  <h3>High Performance</h3>
  <p>LRU cache, JSON encoding pool, and intelligent response limiting for optimal performance</p>
</div>

<div class="feature-card">
  <div class="feature-icon">🔒</div>
  <h3>Secure & Reliable</h3>
  <p>API Key, Bearer Token, and Basic Auth support with secure credential management</p>
</div>

<div class="feature-card">
  <div class="feature-icon">📊</div>
  <h3>Comprehensive Monitoring</h3>
  <p>Integrated with Prometheus, Grafana, Jaeger, and other monitoring tools</p>
</div>

<div class="feature-card">
  <div class="feature-icon">🔧</div>
  <h3>Flexible Configuration</h3>
  <p>Support for SSE, HTTP, and stdio modes to fit various use cases</p>
</div>

</div>

## Learn More

- [View All Services](/en/services/) - Learn about 10 integrated services
- [Complete Tools Reference](/en/docs/tools/) - Detailed documentation for all 220+ tools
- [Deployment Guide](/en/docs/deployment/) - Deployment strategies and best practices
- [Configuration Guide](/en/docs/configuration/) - Configuration options and examples

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

### Binary

```bash
# Download the latest release
curl -LO https://github.com/mahmut-Abi/cloud-native-mcp-server/releases/latest/download/cloud-native-mcp-server-linux-amd64
chmod +x cloud-native-mcp-server-linux-amd64

# Run in SSE mode (default)
./cloud-native-mcp-server-linux-amd64 --mode=sse --addr=0.0.0.0:8080
```

### Docker

```bash
docker run -d \
  --name cloud-native-mcp-server \
  -p 8080:8080 \
  -v ~/.kube:/root/.kube:ro \
  mahmutabi/cloud-native-mcp-server:latest
```

### From Source

```bash
git clone https://github.com/mahmut-Abi/cloud-native-mcp-server.git
cd cloud-native-mcp-server

make build
./cloud-native-mcp-server --mode=sse --addr=0.0.0.0:8080
```

## Key Features

<div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 1.5rem; margin-top: 2rem;">

<div class="feature-card">
  <div class="feature-icon">🚀</div>
  <h3>High Performance</h3>
  <p>LRU cache with TTL support, JSON encoding pool, intelligent response limits for optimal performance</p>
</div>

<div class="feature-card">
  <div class="feature-icon">🔒</div>
  <h3>Secure & Reliable</h3>
  <p>Multiple authentication methods (API Key, Bearer Token, Basic Auth), secure credentials management</p>
</div>

<div class="feature-card">
  <div class="feature-icon">📊</div>
  <h3>Comprehensive Monitoring</h3>
  <p>Integrated Prometheus, Grafana, Jaeger for monitoring and tracing</p>
</div>

<div class="feature-card">
  <div class="feature-icon">🔧</div>
  <h3>Flexible Configuration</h3>
  <p>Supports SSE, HTTP, stdio multiple modes, adapts to various use cases</p>
</div>

<div class="feature-card">
  <div class="feature-icon">📝</div>
  <h3>Audit Logging</h3>
  <p>Complete operation audit and logging, supports multiple storage methods</p>
</div>

<div class="feature-card">
  <div class="feature-icon">🤖</div>
  <h3>AI Optimized</h3>
  <p>Designed for LLMs with summary tools and pagination to prevent context overflow</p>
</div>

</div>

## Services Overview

| Service | Tools | Description |
|---------|-------|-------------|
| **Kubernetes** | 28 | Core container orchestration and resource management |
| **Helm** | 31 | Application package management and deployment |
| **Grafana** | 36 | Visualization, monitoring dashboards, and alerting |
| **Prometheus** | 20 | Metrics collection, querying, and monitoring |
| **Kibana** | 52 | Log analysis, visualization, and data exploration |
| **Elasticsearch** | 14 | Log storage, search, and data indexing |
| **Alertmanager** | 15 | Alert rules management and notifications |
| **Jaeger** | 8 | Distributed tracing and performance analysis |
| **OpenTelemetry** | 9 | Metrics, traces, and logs collection and analysis |
| **Utilities** | 6 | General-purpose utility tools |

**Total: 220+ tools**

## API Endpoints

### SSE Mode

| Endpoint | Description |
|----------|-------------|
| `/api/aggregate/sse` | All services (recommended) |
| `/api/kubernetes/sse` | Kubernetes service |
| `/api/helm/sse` | Helm service |
| `/api/grafana/sse` | Grafana service |
| `/api/prometheus/sse` | Prometheus service |
| `/api/kibana/sse` | Kibana service |
| `/api/elasticsearch/sse` | Elasticsearch service |
| `/api/alertmanager/sse` | Alertmanager service |
| `/api/jaeger/sse` | Jaeger service |
| `/api/opentelemetry/sse` | OpenTelemetry service |
| `/api/utilities/sse` | Utilities service |

### HTTP Mode

Replace `/sse` with `/http` in the endpoints above.

## Documentation

- [Complete Tools Reference](/en/docs/tools/) - Detailed documentation for all 220+ tools
- [Configuration Guide](/en/docs/configuration/) - Configuration options and examples
- [Deployment Guide](/en/docs/deployment/) - Deployment strategies and best practices
- [Security Guide](/en/docs/security/) - Authentication, secrets management, and security best practices
- [Architecture Guide](/en/docs/architecture/) - System architecture and design
- [Performance Guide](/en/docs/performance/) - Performance features and tuning

## Build

```bash
# Build for current platform
make build

# Run tests
make test

# Code linting
make lint

# Docker build
make docker-build
```

## License

MIT License - see [LICENSE](https://github.com/mahmut-Abi/cloud-native-mcp-server/blob/main/LICENSE) for details
---
title: Cloud Native MCP Server
weight: 1
---

<div align="center">

# Cloud Native MCP Server

High-performance MCP server for Kubernetes and cloud-native infrastructure management

[GitHub](https://github.com/mahmut-Abi/cloud-native-mcp-server) • 
[中文](/#)

</div>

---

## Introduction

Cloud Native MCP Server is a high-performance Model Context Protocol (MCP) server for Kubernetes and cloud-native infrastructure management. It integrates 10 services and 220+ tools, enabling AI assistants to effortlessly manage your cloud-native infrastructure.

## Key Features

- 🚀 **High Performance** - LRU cache, JSON encoding pool, intelligent response limiting
- 🔒 **Secure & Reliable** - API Key, Bearer Token, Basic Auth support
- 📊 **Comprehensive Monitoring** - Integrated with Prometheus, Grafana, Jaeger
- 🔧 **Flexible Configuration** - SSE, HTTP, stdio modes
- 📝 **Audit Logging** - Complete operation audit and logging
- 🤖 **AI Optimized** - Designed for LLMs with summary tools and pagination

## Statistics

| Item | Count |
|------|-------|
| Integrated Services | 10 |
| MCP Tools | 220+ |
| Running Modes | 3 |
| License | MIT |

## Quick Start

### Docker Deployment

```bash
docker run -d \
  --name cloud-native-mcp-server \
  -p 8080:8080 \
  -v ~/.kube:/root/.kube:ro \
  mahmutabi/cloud-native-mcp-server:latest
```

### Build from Source

```bash
git clone https://github.com/mahmut-Abi/cloud-native-mcp-server.git
cd cloud-native-mcp-server
make build
./cloud-native-mcp-server --mode=sse --addr=0.0.0.0:8080
```

## Integrated Services

- **Kubernetes** - Container orchestration and resource management
- **Helm** - Application package management and deployment
- **Grafana** - Visualization, monitoring dashboards, and alerting
- **Prometheus** - Metrics collection, querying, and monitoring
- **Kibana** - Log analysis, visualization, and data exploration
- **Elasticsearch** - Log storage, search, and data indexing
- **Alertmanager** - Alert rules management and notifications
- **Jaeger** - Distributed tracing and performance analysis
- **OpenTelemetry** - Metrics, traces, and logs collection
- **Utilities** - General-purpose utility tools

## License

MIT License - see [LICENSE](https://github.com/mahmut-Abi/cloud-native-mcp-server/blob/main/LICENSE) for details

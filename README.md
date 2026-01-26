# Kubernetes MCP Server

[![Go Report Card](https://goreportcard.com/badge/github.com/mahmut-Abi/k8s-mcp-server)](https://goreportcard.com/report/github.com/mahmut-Abi/k8s-mcp-server)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org)

[🇨🇳 中文文档](README-zh.md) | [🇬🇧 English](README.md)

A high-performance Model Context Protocol (MCP) server for Kubernetes and cloud-native infrastructure management with 9 integrated services and 210+ tools.

---

## Features

- **9 Integrated Services**: Kubernetes, Grafana, Prometheus, Kibana, Elasticsearch, Helm, Alertmanager, Jaeger, Utilities
- **210+ MCP Tools**: Comprehensive toolset for infrastructure operations
- **Multi-Protocol Support**: SSE, HTTP, and stdio modes
- **Smart Caching**: LRU cache with TTL support for optimal performance
- **Performance Optimized**: JSON encoding pool, response size control, intelligent limits
- **Enhanced Authentication**: API Key (with complexity requirements), Bearer Token (JWT validation), Basic Auth
- **Secrets Management**: Secure credential storage and rotation
- **Input Sanitization**: Protection against injection attacks
- **Audit Logging**: Track all tool calls and operations
- **LLM-Optimized**: Summary tools and pagination to prevent context overflow

---

## Services Overview

| Service | Tools | Description |
|---------|-------|-------------|
| **kubernetes** | 28 | Core container orchestration and resource management |
| **helm** | 31 | Application package management and deployment |
| **grafana** | 36 | Visualization, monitoring dashboards, and alerting |
| **prometheus** | 20 | Metrics collection, querying, and monitoring |
| **kibana** | 52 | Log analysis, visualization, and data exploration |
| **elasticsearch** | 14 | Log storage, search, and data indexing |
| **alertmanager** | 15 | Alert rules management and notifications |
| **jaeger** | 8 | Distributed tracing and performance analysis |
| **utilities** | 6 | General-purpose utility tools |

**Total: 210+ tools**

---

## Quick Start

### Binary

```bash
# Download the latest release
curl -LO https://github.com/mahmut-Abi/k8s-mcp-server/releases/latest/download/k8s-mcp-server-linux-amd64
chmod +x k8s-mcp-server-linux-amd64

# Run in SSE mode (default)
./k8s-mcp-server-linux-amd64 --mode=sse --addr=0.0.0.0:8080
```

### Docker

```bash
docker run -d \
  --name k8s-mcp-server \
  -p 8080:8080 \
  -v ~/.kube:/root/.kube:ro \
  mahmutabi/k8s-mcp-server:latest
```

### From Source

```bash
git clone https://github.com/mahmut-Abi/k8s-mcp-server.git
cd k8s-mcp-server

make build
./k8s-mcp-server --mode=sse --addr=0.0.0.0:8080
```

---

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
| `/api/utilities/sse` | Utilities service |

### HTTP Mode

Replace `/sse` with `/http` in the endpoints above.

---

## Documentation

- [Complete Tools Reference](docs/TOOLS.md) - Detailed documentation for all 210+ tools
- [Configuration Guide](docs/CONFIGURATION.md) - Configuration options and examples
- [Deployment Guide](docs/DEPLOYMENT.md) - Deployment strategies and best practices
- [Security Guide](docs/SECURITY.md) - Authentication, secrets management, and security best practices
- [Architecture Guide](docs/ARCHITECTURE.md) - System architecture and design
- [Performance Guide](docs/PERFORMANCE.md) - Performance features and tuning

---

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

---

## Contributing

Contributions are welcome! Please read our contributing guidelines and submit pull requests.

---

## License

MIT License - see [LICENSE](LICENSE) for details.
# Kubernetes MCP Server

[![Go Report Card](https://goreportcard.com/badge/github.com/mahmut-Abi/k8s-mcp-server)](https://goreportcard.com/report/github.com/mahmut-Abi/k8s-mcp-server)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)

A high-performance Model Context Protocol (MCP) server for Kubernetes and cloud-native infrastructure management with 9 integrated services and 210+ tools.

---

## Features

- **9 Integrated Services**: Kubernetes, Grafana, Prometheus, Kibana, Elasticsearch, Helm, Alertmanager, Jaeger, Utilities
- **210+ MCP Tools**: Comprehensive toolset for infrastructure operations
- **Multi-Protocol Support**: SSE, HTTP, and stdio modes
- **Smart Caching**: LRU cache with TTL support for optimal performance
- **Performance Optimized**: JSON encoding pool, response size control, intelligent limits
- **Authentication**: API Key, Bearer Token, Basic Auth support
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

# Or HTTP mode
./k8s-mcp-server-linux-amd64 --mode=http --addr=0.0.0.0:8080
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
| `/api/kubernetes/sse` | Kubernetes service |
| `/api/helm/sse` | Helm service |
| `/api/grafana/sse` | Grafana service |
| `/api/prometheus/sse` | Prometheus service |
| `/api/kibana/sse` | Kibana service |
| `/api/elasticsearch/sse` | Elasticsearch service |
| `/api/alertmanager/sse` | Alertmanager service |
| `/api/jaeger/sse` | Jaeger service |
| `/api/utilities/sse` | Utilities service |
| `/api/aggregate/sse` | All services (recommended) |

### HTTP Mode

Replace `/sse` with `/http` in the endpoints above.

---

## Configuration

### YAML Config File

```yaml
# config.yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8080"

logging:
  level: "info"

kubernetes:
  kubeconfig: ""
  timeoutSec: 30

auth:
  enabled: false
  mode: "apikey"
  apiKey: "your-secret-key"

grafana:
  enabled: false
  url: "http://grafana:3000"
  apiKey: ""

prometheus:
  enabled: false
  address: "http://prometheus:9090"

kibana:
  enabled: false
  url: "http://kibana:5601"

elasticsearch:
  enabled: false
  url: "http://elasticsearch:9200"

alertmanager:
  enabled: false
  url: "http://alertmanager:9093"

jaeger:
  enabled: false
  url: "http://jaeger:16686"

audit:
  enabled: false
  maxLogs: 1000
```

### Environment Variables

```bash
export MCP_MODE=sse
export MCP_ADDR=0.0.0.0:8080
export MCP_LOG_LEVEL=info
export MCP_AUTH_ENABLED=false
export MCP_K8S_KUBECONFIG=~/.kube/config
```

### Command Line Flags

```bash
./k8s-mcp-server \
  --mode=sse \
  --addr=0.0.0.0:8080 \
  --config=config.yaml \
  --log-level=info
```

---

## Available Tools

For a complete list of all 210+ tools with detailed descriptions, see [TOOLS.md](docs/TOOLS.md).

### Quick Reference

#### Kubernetes (28 tools)
- `kubernetes_list_resources_summary` - List resources with optimized output
- `kubernetes_get_resource_summary` - Get single resource summary
- `kubernetes_get_pod_logs` - Get pod logs
- `kubernetes_get_events` - Get cluster events
- `kubernetes_describe_resource` - Describe resource in detail
- And 23 more...

#### Helm (31 tools)
- `helm_list_releases_paginated` - List releases with pagination
- `helm_get_release_summary` - Get release summary
- `helm_search_charts` - Search Helm charts
- `helm_cluster_overview` - Get cluster overview
- And 27 more...

#### Grafana (36 tools)
- `grafana_dashboards_summary` - List dashboards with minimal output
- `grafana_datasources_summary` - List data sources
- `grafana_dashboard` - Get specific dashboard
- `grafana_alerts` - List alert rules
- And 32 more...

#### Prometheus (20 tools)
- `prometheus_query` - Execute instant query
- `prometheus_query_range` - Execute range query
- `prometheus_alerts_summary` - Get alerts summary
- `prometheus_targets_summary` - Get targets summary
- And 16 more...

#### Kibana (52 tools)
- `kibana_search_saved_objects` - Search saved objects
- `kibana_get_index_patterns` - Get index patterns
- `kibana_get_spaces` - Get Kibana spaces
- And 49 more...

#### Elasticsearch (14 tools)
- `elasticsearch_list_indices_paginated` - List indices with pagination
- `elasticsearch_cluster_health_summary` - Get cluster health
- `elasticsearch_search_indices` - Search indices
- And 11 more...

#### Alertmanager (15 tools)
- `alertmanager_alerts_summary` - Get alerts summary
- `alertmanager_silences_summary` - Get silences summary
- `alertmanager_create_silence` - Create silence
- And 12 more...

#### Jaeger (8 tools)
- `jaeger_get_traces_summary` - Get traces summary
- `jaeger_get_trace` - Get specific trace
- `jaeger_get_services` - Get all services
- And 5 more...

#### Utilities (6 tools)
- `utilities_get_time` - Get current time
- `utilities_get_timestamp` - Get Unix timestamp
- `utilities_web_fetch` - Fetch URL content
- And 3 more...

---

## LLM-Optimized Tools

Many tools have LLM-optimized versions marked with ⚠️ PRIORITY that provide:
- 70-95% smaller response sizes
- Essential fields only
- Pagination support
- Context overflow prevention

Examples:
- `kubernetes_list_resources_summary` vs `kubernetes_list_resources`
- `grafana_dashboards_summary` vs `grafana_dashboards`
- `prometheus_alerts_summary` vs `prometheus_get_alerts`
- `elasticsearch_indices_summary` vs `elasticsearch_list_indices`

---

## Project Structure

```
k8s-mcp-server/
├── cmd/
│   └── server/              # Main entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── logging/             # Logging utilities
│   ├── middleware/          # HTTP middleware (auth, audit, metrics)
│   ├── observability/       # Metrics and monitoring
│   ├── services/            # Service implementations
│   │   ├── kubernetes/      # Kubernetes service (28 tools)
│   │   ├── helm/            # Helm service (31 tools)
│   │   ├── grafana/         # Grafana service (36 tools)
│   │   ├── prometheus/      # Prometheus service (20 tools)
│   │   ├── kibana/          # Kibana service (52 tools)
│   │   ├── elasticsearch/   # Elasticsearch service (14 tools)
│   │   ├── alertmanager/    # Alertmanager service (15 tools)
│   │   ├── jaeger/          # Jaeger service (8 tools)
│   │   ├── utilities/       # Utilities service (6 tools)
│   │   ├── cache/           # LRU cache implementation
│   │   ├── framework/       # Service initialization framework
│   │   └── manager/         # Service manager
│   └── util/                # Utilities
│       ├── circuitbreaker/  # Circuit breaker pattern
│       ├── performance/     # Performance optimizations
│       └── pool/            # Object pooling
├── docs/                    # Documentation
│   └── TOOLS.md            # Complete tools reference
└── deploy/                  # Deployment files
    ├── Dockerfile
    ├── helm/
    │   └── k8s-mcp-server/
    └── kubernetes/
```

---

## Build

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run with race detection
make test-race

# Code linting
make lint

# Docker build
make docker-build
```

---

## Performance Features

- **Intelligent Caching**: LRU cache with TTL for frequently accessed data
- **Response Size Control**: Automatic truncation and optimization
- **JSON Encoding Pool**: Reuse JSON encoders for better performance
- **Circuit Breaker**: Prevent cascading failures
- **Pagination**: Support for large datasets
- **Summary Tools**: Optimized tools for LLM consumption

---

## Documentation

- [Complete Tools Reference](docs/TOOLS.md) - Detailed documentation for all 210+ tools
- [Configuration Guide](docs/CONFIGURATION.md) - Configuration options and examples
- [Deployment Guide](docs/DEPLOYMENT.md) - Deployment strategies and best practices

---

## Contributing

Contributions are welcome! Please read our contributing guidelines and submit pull requests.

---

## License

MIT License - see [LICENSE](LICENSE) for details.
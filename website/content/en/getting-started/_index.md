---
title: "Getting Started"
weight: 1
---

# Getting Started with Cloud Native MCP Server

Welcome to the Cloud Native MCP Server! This guide will help you get up and running quickly with the most powerful Model Context Protocol (MCP) server for Kubernetes and cloud-native infrastructure management.

## Overview

Cloud Native MCP Server is a high-performance Model Context Protocol (MCP) server that integrates 10 services and 220+ tools to enable AI assistants to effortlessly manage your cloud-native infrastructure.

### What You'll Learn

- How to install and deploy Cloud Native MCP Server
- Basic configuration options
- How to use the core services
- Best practices for security and performance

---

## Installation Options

Choose the installation method that best fits your environment:

{{< tabs >}}
{{< tab "Docker" >}}
### Docker Installation

The easiest way to get started is using Docker:

```bash
# Pull the latest image
docker pull mahmutabi/cloud-native-mcp-server:latest

# Run the server
docker run -d \
  --name cloud-native-mcp-server \
  -p 8080:8080 \
  -v ~/.kube:/root/.kube:ro \
  -e MCP_SERVER_API_KEY="your-secure-api-key" \
  mahmutabi/cloud-native-mcp-server:latest
```

Once running, you can access the server at `http://localhost:8080`.
{{< /tab >}}

{{< tab "Binary" >}}
### Binary Installation

Download the pre-compiled binary for your OS:

```bash
# Linux (amd64)
curl -LO https://github.com/mahmut-Abi/cloud-native-mcp-server/releases/latest/download/cloud-native-mcp-server-linux-amd64
chmod +x cloud-native-mcp-server-linux-amd64

# Run the server
./cloud-native-mcp-server-linux-amd64 --mode=sse --addr=0.0.0.0:8080
```

The binary includes all 10 integrated services and 220+ tools.
{{< /tab >}}

{{< tab "Source" >}}
### From Source

Build from source for development or customization:

```bash
# Clone the repository
git clone https://github.com/mahmut-Abi/cloud-native-mcp-server.git
cd cloud-native-mcp-server

# Build the server
make build

# Run with default settings
./cloud-native-mcp-server --mode=sse --addr=0.0.0.0:8080
```

Make sure you have Go 1.25+ installed.
{{< /tab >}}
{{< /tabs >}}

---

## Initial Configuration

After installation, you'll need to configure your server with appropriate authentication and service endpoints.

### Authentication Setup

The server supports multiple authentication methods:

```bash
# API Key (recommended for production)
export MCP_SERVER_API_KEY="your-very-secure-api-key-with-32-chars-minimum"

# Or Bearer Token (JWT)
export MCP_SERVER_BEARER_TOKEN="your-jwt-token"

# Or Basic Auth
export MCP_SERVER_BASIC_AUTH_USER="admin"
export MCP_SERVER_BASIC_AUTH_PASS="secure-password"
```

### Service Configuration

The server will automatically detect and configure services if they're accessible:

- Kubernetes: Requires `~/.kube/config` or in-cluster config
- Prometheus: Connects to `http://prometheus:9090` by default
- Grafana: Connects to `http://grafana:3000` by default
- And more...

---

## Your First MCP Call

Once your server is running, you can make your first MCP call:

```bash
curl -X POST http://localhost:8080/v1/mcp/list-tools \
  -H "Authorization: Bearer your-api-key" \
  -H "Content-Type: application/json"
```

This will return a list of all 220+ available tools across the 10 integrated services.

---

## Integrated Services Overview

Cloud Native MCP Server integrates 10 core services:

{{< columns >}}
### 🔧 Kubernetes
Manage your Kubernetes clusters with 28 specialized tools for deployments, services, configmaps, secrets, and more.
<--->

### 📦 Helm
Deploy and manage Helm charts with 31 tools for chart management, releases, and repositories.
{{< /columns >}}

{{< columns >}}
### 📊 Grafana
Create and manage dashboards, alerts, and data sources with 36 monitoring tools.
<--->

### 📈 Prometheus
Query metrics, manage rules, and configure alerting with 20 observability tools.
{{< /columns >}}

{{< columns >}}
### 🔍 Kibana
Analyze logs and visualize data with 52 Elasticsearch integration tools.
<--->

### ⚡ Elasticsearch
Index, search, and analyze data with 14 advanced search tools.
{{< /columns >}}

---

## Next Steps

Now that you've installed and configured Cloud Native MCP Server, you might want to:

- [Configure authentication and security settings](/en/guides/security/)
- [Explore service-specific configurations](/en/guides/configuration/)
- [Learn about performance optimization](/en/guides/performance/)
- [Review the complete tools reference](/en/docs/tools/)

### Quick Links

- [Architecture Overview](/en/concepts/architecture/)
- [Security Best Practices](/en/guides/security/best-practices/)
- [Performance Tuning](/en/guides/performance/optimization/)
- [Troubleshooting](/en/guides/troubleshooting/)

---

## Support and Community

Need help? Check out these resources:

- [GitHub Issues](https://github.com/mahmut-Abi/cloud-native-mcp-server/issues) for bug reports
- [GitHub Discussions](https://github.com/mahmut-Abi/cloud-native-mcp-server/discussions) for questions
- [Documentation](/) for complete reference
- [Contributing Guide](https://github.com/mahmut-Abi/cloud-native-mcp-server/blob/main/CONTRIBUTING.md) to get involved
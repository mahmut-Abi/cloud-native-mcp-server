---
title: "Documentation"
---

# Documentation Center

Welcome to the Cloud Native MCP Server documentation center. This includes complete usage guides, API references, and best practices.

## Quick Start

- [Installation Guide](#installation)
- [Configuration](#configuration)
- [Quick Start](#quick-start)

## Core Documentation

### [Complete Tools Reference](/en/docs/tools/)
Detailed documentation for all 220+ MCP tools, including usage methods and examples.

### [Configuration Guide](/en/docs/configuration/)
Complete configuration options, including server configuration, service configuration, authentication configuration, etc.

### [Deployment Guide](/en/docs/deployment/)
Detailed instructions for various deployment methods, including Docker, Kubernetes, Helm, etc.

### [Security Guide](/en/docs/security/)
Authentication, authorization, secrets management, and security best practices.

### [Architecture Guide](/en/docs/architecture/)
System architecture design, component descriptions, and extension methods.

### [Performance Guide](/en/docs/performance/)
Performance optimization recommendations, tuning parameters, and best practices.

## Installation

### System Requirements

- Go 1.25 or higher version
- Linux, macOS, or Windows
- kubeconfig file for Kubernetes cluster access
- Optional: Docker (for containerized deployment)

### Binary Installation

```bash
# Download latest version
curl -LO https://github.com/mahmut-Abi/cloud-native-mcp-server/releases/latest/download/cloud-native-mcp-server-linux-amd64
chmod +x cloud-native-mcp-server-linux-amd64

# Verify installation
./cloud-native-mcp-server-linux-amd64 --version
```

### Docker Installation

```bash
docker pull mahmutabi/cloud-native-mcp-server:latest

# Verify installation
docker run --rm mahmutabi/cloud-native-mcp-server:latest --version
```

### Build from Source

```bash
git clone https://github.com/mahmut-Abi/cloud-native-mcp-server.git
cd cloud-native-mcp-server

# Build
make build

# Verify
./cloud-native-mcp-server --version
```

## Configuration

### Basic Configuration

Create configuration file `config.yaml`:

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8080"

kubernetes:
  kubeconfig: ""  # Use default kubeconfig

logging:
  level: "info"
  json: false
```

### Enable Services

```yaml
prometheus:
  enabled: true
  address: "http://localhost:9090"

grafana:
  enabled: true
  url: "http://localhost:3000"
  apiKey: "your-api-key"
```

### Configure Authentication

```yaml
auth:
  enabled: true
  mode: "apikey"
  apiKey: "your-secure-api-key"
```

### Environment Variables

All configurations support environment variables:

```bash
export MCP_MODE=sse
export MCP_ADDR=0.0.0.0:8080
export MCP_KUBECONFIG=/path/to/kubeconfig
export MCP_LOG_LEVEL=info
```

## Quick Start

### 1. Start Server

```bash
# Run with default configuration
./cloud-native-mcp-server --mode=sse --addr=0.0.0.0:8080

# Run with configuration file
./cloud-native-mcp-server --config=config.yaml
```

### 2. Connect to Server

**SSE Mode:**

```bash
curl -N http://localhost:8080/api/aggregate/sse
```

**HTTP Mode:**

```bash
curl http://localhost:8080/api/aggregate/http \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}'
```

### 3. Use Tools

Example of calling Kubernetes tool:

```bash
curl -N http://localhost:8080/api/kubernetes/sse \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "list_pods",
      "arguments": {
        "namespace": "default"
      }
    }
  }'
```

### 4. Use in MCP Client

Configure MCP client (e.g., Claude Desktop):

```json
{
  "mcpServers": {
    "cloud-native": {
      "command": "/path/to/cloud-native-mcp-server",
      "args": ["--mode=stdio"]
    }
  }
}
```

## Running Modes

### SSE Mode (Recommended for Production)
Real-time bidirectional communication, suitable for scenarios requiring real-time updates.

```bash
./cloud-native-mcp-server --mode=sse --addr=0.0.0.0:8080
```

### HTTP Mode
Standard REST API, easy to integrate.

```bash
./cloud-native-mcp-server --mode=http --addr=0.0.0.0:8080
```

### stdio Mode (Recommended for Development)
Standard input/output, suitable for MCP clients.

```bash
./cloud-native-mcp-server --mode=stdio
```

### Streamable-HTTP Mode
MCP 2025-11-25 specification, modern communication method.

```bash
./cloud-native-mcp-server --mode=streamable-http --addr=0.0.0.0:8080
```

## Troubleshooting

### Common Issues

**Q: Cannot connect to Kubernetes cluster**
- Check kubeconfig file path
- Verify cluster access permissions
- Check network connection

**Q: Slow service response**
- Enable caching functionality
- Adjust timeout settings
- Check cluster resources

**Q: Authentication failure**
- Verify API key configuration
- Check authentication mode settings
- Confirm token validity

### Debug Logging

Enable debug logs:

```bash
./cloud-native-mcp-server --log-level=debug
```

Or in configuration file:

```yaml
logging:
  level: "debug"
  json: true
```

## Next Steps

- Read complete [Tools Reference Documentation](/en/docs/tools/)
- View [Deployment Guide](/en/docs/deployment/) for production deployment
- Learn [Security Best Practices](/en/docs/security/)
- Explore [Performance Optimization Tips](/en/docs/performance/)

## Get Help

- View [GitHub Issues](https://github.com/mahmut-Abi/cloud-native-mcp-server/issues)
- Read [Project Wiki](https://github.com/mahmut-Abi/cloud-native-mcp-server/wiki)
- Join community discussions
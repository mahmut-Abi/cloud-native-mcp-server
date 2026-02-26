---
title: "Documentation"
weight: 1
description: "Start here for setup, architecture, configuration, deployment, security, performance, API, and tools."
bookCollapseSection: true
---

# Documentation Center

This section is the primary reference for deploying and operating Cloud Native MCP Server.

## Read by Goal

- [Install and run your first instance]({{< relref "/getting-started/_index.md" >}})
- [Understand architecture and request flow]({{< relref "architecture.md" >}})
- [Configure server, auth, and integrations]({{< relref "configuration.md" >}})
- [Deploy in production environments]({{< relref "deployment.md" >}})
- [Apply hardening and security controls]({{< relref "security.md" >}})
- [Tune throughput and latency]({{< relref "performance.md" >}})
- [Explore all available tools]({{< relref "tools.md" >}})
- [Integrate via MCP endpoints]({{< relref "api.md" >}})

## Quick Start

```bash
docker run -d \
  --name cloud-native-mcp-server \
  -p 8080:8080 \
  -v ~/.kube:/root/.kube:ro \
  mahmutabi/cloud-native-mcp-server:latest
```

Then open:

- `SSE`: `http://localhost:8080/api/aggregate/sse`
- `HTTP`: `http://localhost:8080/api/aggregate/http`

## Reference Structure

- `Getting Started`: installation and first-call flow
- `Architecture`: system model and component boundaries
- `Configuration`: all runtime and integration options
- `Deployment`: Docker, Kubernetes, and Helm production patterns
- `Security`: auth and operational hardening practices
- `Performance`: optimization and benchmarking guidance
- `Tools`: complete service tool catalog
- `API`: protocol endpoints and request examples

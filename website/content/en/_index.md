---
title: Cloud Native MCP Server
weight: 1
description: High-performance Model Context Protocol server for Kubernetes and cloud-native operations.
---

<div class="hero">
  <h1>Cloud Native MCP Server</h1>
  <p>A production-grade MCP server for Kubernetes and cloud-native infrastructure management, exposing 10 services and 220+ tools across SSE, Streamable HTTP, HTTP, and stdio modes.</p>
  <div class="hero-buttons">
    <a href="https://github.com/mahmut-Abi/cloud-native-mcp-server" class="cta-button"><span>GitHub Repository</span></a>
    <a href="#quick-start" class="cta-button transparent"><span>Quick Start</span></a>
  </div>
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/mahmut-Abi/cloud-native-mcp-server)](https://goreportcard.com/report/github.com/mahmut-Abi/cloud-native-mcp-server)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org)

## Core Value

<div class="value-grid">
  <article class="value-card">
    <h3>Single Operations Interface</h3>
    <p>Unifies Kubernetes, Helm, Grafana, Prometheus, Kibana, and more behind one MCP surface, reducing context switching for operators and agents.</p>
  </article>
  <article class="value-card">
    <h3>Production Security Controls</h3>
    <p>Supports apikey / bearer / basic authentication, rate limiting, and audit logging for enterprise hardening and compliance workflows.</p>
  </article>
  <article class="value-card">
    <h3>Agent-Friendly Output</h3>
    <p>Pagination and summarization patterns help AI assistants stay efficient during large incident investigations.</p>
  </article>
</div>

---

## Typical Use Cases

<div class="usecase-grid">
  <article class="usecase-card">
    <h3>Incident Triage</h3>
    <p>Correlate pod state, events, logs, and metrics to shorten the path from alert to root cause.</p>
  </article>
  <article class="usecase-card">
    <h3>Release and Change Control</h3>
    <p>Use Helm and Kubernetes tools for rollout, rollback, scaling, and controlled production changes with auditable traces.</p>
  </article>
  <article class="usecase-card">
    <h3>Cross-Observability Analysis</h3>
    <p>Bridge Prometheus, Grafana, Jaeger, and OpenTelemetry signals for end-to-end diagnostics.</p>
  </article>
</div>

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
  -e MCP_AUTH_ENABLED=true \
  -e MCP_AUTH_MODE=apikey \
  -e MCP_AUTH_API_KEY='ChangeMe-Strong-Key-123!' \
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
./cloud-native-mcp-server --mode=streamable-http --addr=0.0.0.0:8080
{{< /highlight >}}
{{< /tab >}}
{{< /tabs >}}

### Availability Check

{{< highlight bash >}}
# 1) Health check
curl -sS http://127.0.0.1:8080/health

# 2) End-to-end SSE handshake + initialize check (run at repo root)
make sse-smoke BASE_URL=http://127.0.0.1:8080
{{< /highlight >}}

### Common Entry Points

- Aggregate SSE endpoint (`--mode=sse`): `http://127.0.0.1:8080/api/aggregate/sse`
- Aggregate Streamable HTTP endpoint (`--mode=streamable-http`): `http://127.0.0.1:8080/api/aggregate/streamable-http`
- Health endpoint: `http://127.0.0.1:8080/health`

---

## Pre-Production Checklist

<div class="ops-grid">
  <article class="ops-card">
    <h3>Authentication and Access</h3>
    <ul>
      <li>Enable `MCP_AUTH_ENABLED=true` in production.</li>
      <li>Choose one mode: `apikey`, `bearer`, or `basic`.</li>
      <li>Apply least-privilege access to Kubernetes and external systems.</li>
    </ul>
  </article>
  <article class="ops-card">
    <h3>Observability and Audit</h3>
    <ul>
      <li>Enable structured logs and core metrics collection.</li>
      <li>Enable audit logs if change tracking is required.</li>
      <li>Continuously validate `/health` and core upstream service checks.</li>
    </ul>
  </article>
  <article class="ops-card">
    <h3>Performance and Resilience</h3>
    <ul>
      <li>Tune rate limits, timeouts, and concurrency for your traffic profile.</li>
      <li>Prefer summary and pagination tools to limit context size.</li>
      <li>Load-test with realistic multi-service tool-call patterns.</li>
    </ul>
  </article>
</div>

---

## Documentation Map

- [Getting Started]({{< relref "getting-started/_index.md" >}})
- [Getting Started FAQ]({{< relref "getting-started/faq.md" >}})
- [Troubleshooting]({{< relref "getting-started/troubleshooting.md" >}})
- [Architecture]({{< relref "docs/architecture.md" >}})
- [Configuration]({{< relref "docs/configuration.md" >}})
- [Deployment]({{< relref "docs/deployment.md" >}})
- [Security]({{< relref "docs/security.md" >}})
- [Performance]({{< relref "docs/performance.md" >}})
- [Tools Reference]({{< relref "docs/tools.md" >}})
- [Services Overview]({{< relref "services/_index.md" >}})
- [Sitemap]({{< relref "sitemap.md" >}})

---

## FAQ and Troubleshooting Entry

<div class="resource-grid">
  <article class="resource-card">
    <h3>Getting Started FAQ</h3>
    <p>Answers common implementation questions around auth mode, transport mode, client integration, and production rollout strategy.</p>
    <a class="resource-link" href='{{< relref "getting-started/faq.md" >}}'>Read FAQ</a>
  </article>
  <article class="resource-card">
    <h3>Troubleshooting Playbook</h3>
    <p>Step-by-step checks for startup failures, 401 responses, SSE handshake issues, and unavailable service integrations.</p>
    <a class="resource-link" href='{{< relref "getting-started/troubleshooting.md" >}}'>Open Troubleshooting</a>
  </article>
</div>

---

## More Resources

- [Blog]({{< relref "posts/_index.md" >}})
- [Showcase]({{< relref "showcase.md" >}})
- [GitHub Repository](https://github.com/mahmut-Abi/cloud-native-mcp-server)

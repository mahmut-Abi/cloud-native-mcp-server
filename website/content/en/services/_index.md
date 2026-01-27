---
title: "Services"
---

# Integrated Services

Cloud Native MCP Server integrates 10 powerful cloud-native services, providing 220+ tools that fully cover Kubernetes management and application deployment, monitoring, log analysis, and more.

<div class="service-grid">

<div class="service-card">
  <h3>☸️ Kubernetes</h3>
  <p><span class="tool-count">28 tools</span></p>
  <p>Core container orchestration and resource management, including complete lifecycle management for Pods, Deployments, Services, ConfigMaps, Secrets, and more.</p>
  <p><strong>Key Features:</strong></p>
  <ul>
    <li>Pod and container management</li>
    <li>Application deployment and scaling</li>
    <li>Service discovery and load balancing</li>
    <li>Configuration and secrets management</li>
    <li>Namespace and node management</li>
  </ul>
</div>

<div class="service-card">
  <h3>⚓ Helm</h3>
  <p><span class="tool-count">31 tools</span></p>
  <p>Kubernetes application package manager, simplifying application deployment, upgrades, and management.</p>
  <p><strong>Key Features:</strong></p>
  <ul>
    <li>Chart repository management</li>
    <li>Release lifecycle management</li>
    <li>Values configuration management</li>
    <li>Dependency management</li>
    <li>Plugin system</li>
  </ul>
</div>

<div class="service-card">
  <h3>📊 Grafana</h3>
  <p><span class="tool-count">36 tools</span></p>
  <p>Open-source analytics and visualization platform for monitoring and metrics visualization.</p>
  <p><strong>Key Features:</strong></p>
  <ul>
    <li>Dashboard management</li>
    <li>Data source configuration</li>
    <li>Visualization creation</li>
    <li>Alert management</li>
    <li>User and organization management</li>
  </ul>
</div>

<div class="service-card">
  <h3>📈 Prometheus</h3>
  <p><span class="tool-count">20 tools</span></p>
  <p>Open-source monitoring and alerting system for collecting and querying time-series data.</p>
  <p><strong>Key Features:</strong></p>
  <ul>
    <li>Instant and range queries</li>
    <li>Label and metadata queries</li>
    <li>Target management</li>
    <li>Rules and alert management</li>
    <li>TSDB and storage management</li>
  </ul>
</div>

<div class="service-card">
  <h3>🔍 Kibana</h3>
  <p><span class="tool-count">52 tools</span></p>
  <p>Elastic Stack data visualization and management interface for log analysis and data exploration.</p>
  <p><strong>Key Features:</strong></p>
  <ul>
    <li>Index and document management</li>
    <li>Data query and aggregation</li>
    <li>Visualization and dashboards</li>
    <li>Index pattern management</li>
    <li>Space and permission management</li>
  </ul>
</div>

<div class="service-card">
  <h3>🔎 Elasticsearch</h3>
  <p><span class="tool-count">14 tools</span></p>
  <p>Distributed search and analytics engine for log storage and full-text search.</p>
  <p><strong>Key Features:</strong></p>
  <ul>
    <li>Index management</li>
    <li>Document operations</li>
    <li>Data search</li>
    <li>Cluster management</li>
    <li>Alias management</li>
  </ul>
</div>

<div class="service-card">
  <h3>🚨 Alertmanager</h3>
  <p><span class="tool-count">15 tools</span></p>
  <p>Prometheus alert handling and routing system for managing alert notifications.</p>
  <p><strong>Key Features:</strong></p>
  <ul>
    <li>Alert management</li>
    <li>Silence rules</li>
    <li>Alert routing</li>
    <li>Notification configuration</li>
    <li>Rule group management</li>
  </ul>
</div>

<div class="service-card">
  <h3>🔗 Jaeger</h3>
  <p><span class="tool-count">8 tools</span></p>
  <p>Distributed tracing platform for monitoring and troubleshooting microservice architectures.</p>
  <p><strong>Key Features:</strong></p>
  <ul>
    <li>Trace queries</li>
    <li>Service discovery</li>
    <li>Dependency analysis</li>
    <li>Performance analysis</li>
    <li>Metric queries</li>
  </ul>
</div>

<div class="service-card">
  <h3>📡 OpenTelemetry</h3>
  <p><span class="tool-count">9 tools</span></p>
  <p>Unified observability framework for collecting metrics, traces, and logs.</p>
  <p><strong>Key Features:</strong></p>
  <ul>
    <li>Metric collection</li>
    <li>Trace management</li>
    <li>Log aggregation</li>
    <li>Unified configuration</li>
    <li>Cross-language support</li>
  </ul>
</div>

<div class="service-card">
  <h3>🛠️ Utilities</h3>
  <p><span class="tool-count">6 tools</span></p>
  <p>General-purpose utility toolkit providing common data processing and transformation functions.</p>
  <p><strong>Key Features:</strong></p>
  <ul>
    <li>Base64 encoding/decoding</li>
    <li>JSON processing</li>
    <li>Timestamp generation</li>
    <li>UUID generation</li>
    <li>Data conversion</li>
  </ul>
</div>

</div>

## Service Configuration

Each service can be independently configured and enabled/disabled. Here's an example configuration:

### Enable Services

```yaml
# Kubernetes (enabled by default)
kubernetes:
  kubeconfig: ""
  timeoutSec: 30
  qps: 100.0
  burst: 200

# Prometheus
prometheus:
  enabled: true
  address: "http://localhost:9090"
  timeoutSec: 30

# Grafana
grafana:
  enabled: true
  url: "http://localhost:3000"
  apiKey: "your-api-key"

# Kibana
kibana:
  enabled: true
  url: "https://localhost:5601"
  apiKey: "your-api-key"

# Elasticsearch
elasticsearch:
  enabled: true
  addresses:
    - "http://localhost:9200"
  timeoutSec: 30

# Helm
helm:
  enabled: true
  namespace: "default"
  timeoutSec: 300

# Alertmanager
alertmanager:
  enabled: true
  address: "http://localhost:9093"
  timeoutSec: 30

# Jaeger
jaeger:
  enabled: true
  address: "http://localhost:16686"
  timeoutSec: 30

# OpenTelemetry
opentelemetry:
  enabled: true
  address: "http://localhost:4318"
  timeoutSec: 30
```

## Service Filtering

Services can be enabled or disabled via configuration file or environment variables:

```yaml
enableDisable:
  # Disabled services (comma-separated)
  disabledServices: []

  # Enabled services (comma-separated, overrides disabled)
  enabledServices: ["kubernetes", "helm", "prometheus", "grafana"]

  # Disabled tools (comma-separated)
  disabledTools: []
```

Or using environment variables:

```bash
export MCP_ENABLED_SERVICES=kubernetes,helm,prometheus,grafana
export MCP_DISABLED_SERVICES=elasticsearch,kibana
```

## Authentication Configuration

Each service supports multiple authentication methods:

### Basic Auth

```yaml
prometheus:
  username: "admin"
  password: "password"
```

### API Key

```yaml
grafana:
  apiKey: "eyJrIjoi..."
```

### Bearer Token

```yaml
kibana:
  bearerToken: "eyJhbGci..."
```

### TLS/mTLS

```yaml
elasticsearch:
  tlsSkipVerify: false
  tlsCertFile: "/path/to/cert.pem"
  tlsKeyFile: "/path/to/key.pem"
  tlsCAFile: "/path/to/ca.pem"
```

## Best Practices

### 1. Enable Services on Demand
Only enable the services you need to reduce resource consumption and attack surface.

### 2. Configure Appropriate Timeouts
Set suitable timeout values based on service response times.

### 3. Use Environment Variables
Use environment variables for sensitive information (like API keys, passwords).

### 4. Enable Caching
Enable caching for frequently accessed services to improve performance.

### 5. Monitor Service Health
Regularly check the health status of each service.

### 6. Use Connection Pools
Configure appropriate QPS and burst parameters.

### 7. Error Handling
Configure appropriate error handling and retry policies for each service.

## Troubleshooting

### Service Connection Failed
**Problem**: Cannot connect to service

**Solution**:
```bash
# Check pod status
kubectl get pods -l app=cloud-native-mcp-server

# Check logs
kubectl logs -l app=cloud-native-mcp-server

# Check service
kubectl get svc cloud-native-mcp-server

# Port forward test
kubectl port-forward svc/cloud-native-mcp-server 8080:8080
curl http://localhost:8080/health
```

### Authentication Failed
**Problem**: 401 Unauthorized

**Solution**:
```bash
# Check auth configuration
kubectl get configmap k8s-mcp-config -o yaml

# Verify secrets
kubectl get secret k8s-mcp-secrets -o yaml

# Test with correct headers
curl -H "X-API-Key: your-key" http://localhost:8080/health
```

### Performance Issues
**Problem**: Slow response times

**Solution**:
```yaml
# Increase timeouts
kubernetes:
  timeoutSec: 60

# Enable caching
config:
  cache:
    enabled: true

# Use summary tools
# Replace kubernetes_list_resources with kubernetes_list_resources_summary
```

## Extension Development

To add a new service, refer to the development documentation:

1. Create service directory structure
2. Implement service interface
3. Register tools
4. Write tests
5. Update documentation

Detailed guide available at [Development Documentation](https://github.com/mahmut-Abi/cloud-native-mcp-server/tree/main/docs/development).
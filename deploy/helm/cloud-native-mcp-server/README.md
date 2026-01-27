# Kubernetes MCP Server Helm Chart

This Helm chart deploys Kubernetes MCP Server in a Kubernetes cluster.

## Introduction

Kubernetes MCP Server is an intelligent cloud-native management platform designed for the AI era. It provides context-aware cloud-native resource management capabilities through the standardized MCP (Model Context Protocol).

## Installation

### Add Helm Repository

```bash
# Add repository
helm repo add cloud-native-mcp-server https://charts.example.com/cloud-native-mcp-server
helm repo update

# View available versions
helm search repo cloud-native-mcp-server
```

### Install Chart

```bash
# Install latest version
helm install my-release cloud-native-mcp-server/cloud-native-mcp-server

# Install specific version
helm install my-release cloud-native-mcp-server/cloud-native-mcp-server --version 0.1.0

# Install to specific namespace
helm install my-release cloud-native-mcp-server/cloud-native-mcp-server --namespace mcp-server --create-namespace
```

### Use Local Chart

```bash
# Clone repository
git clone https://github.com/mahmut-Abi/cloud-native-mcp-server.git
cd cloud-native-mcp-server/helm/cloud-native-mcp-server

# Install
helm install my-release . --namespace mcp-server --create-namespace
```

## Configuration

### View Default Configuration

```bash
helm show values cloud-native-mcp-server/cloud-native-mcp-server
```

### Custom Configuration

Create a `values.yaml` file:

```yaml
# Example configuration
replicaCount: 3

image:
  tag: "v1.0.0"

config:
  server:
    mode: "sse"
    addr: "0.0.0.0:8080"

  logging:
    level: "info"
    json: true

  auth:
    enabled: true
    mode: "apikey"
    apiKey: "your-secret-api-key"

  prometheus:
    enabled: true
    address: "http://prometheus:9090"

  grafana:
    enabled: true
    url: "http://grafana:3000"
    apiKey: "your-grafana-api-key"
```

Install with custom configuration:

```bash
helm install my-release cloud-native-mcp-server/cloud-native-mcp-server -f values.yaml --namespace mcp-server
```

### Override via Command Line

```bash
helm install my-release cloud-native-mcp-server/cloud-native-mcp-server \
  --set replicaCount=5 \
  --set config.auth.apiKey="new-secret-key" \
  --set config.prometheus.enabled=true \
  --namespace mcp-server
```

## Upgrade

```bash
# Upgrade to latest version
helm upgrade my-release cloud-native-mcp-server/cloud-native-mcp-server

# Upgrade with custom configuration
helm upgrade my-release cloud-native-mcp-server/cloud-native-mcp-server -f values.yaml

# Upgrade via command line
helm upgrade my-release cloud-native-mcp-server/cloud-native-mcp-server \
  --set replicaCount=5 \
  --set config.auth.apiKey="updated-secret-key"
```

## Uninstall

```bash
# Uninstall
helm uninstall my-release --namespace mcp-server

# Delete PVC (if persistent storage was used)
kubectl delete pvc --namespace mcp-server -l app.kubernetes.io/instance=my-release
```

## Configuration Parameters

### Global Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `global.imageRegistry` | Global image registry | `""` |
| `global.imagePullSecrets` | Global image pull secrets | `[]` |
| `global.storageClass` | Global storage class | `""` |

### Image Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `image.registry` | Image registry | `"docker.io"` |
| `image.repository` | Image repository | `"mahmutabi/cloud-native-mcp-server"` |
| `image.tag` | Image tag | `"latest"` |
| `image.pullPolicy` | Image pull policy | `"IfNotPresent"` |
| `image.pullSecrets` | Image pull secrets | `[]` |

### Replicas and Resources

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of replicas | `3` |
| `resources.limits.cpu` | CPU limit | `500m` |
| `resources.limits.memory` | Memory limit | `512Mi` |
| `resources.requests.cpu` | CPU request | `100m` |
| `resources.requests.memory` | Memory request | `128Mi` |

### Service Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `service.type` | Service type | `"ClusterIP"` |
| `service.port` | Service port | `8080` |
| `service.targetPort` | Target port | `8080` |
| `service.annotations` | Service annotations | `{}` |

### Ingress Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `ingress.enabled` | Enable Ingress | `false` |
| `ingress.className` | Ingress class name | `""` |
| `ingress.annotations` | Ingress annotations | `{}` |
| `ingress.hosts` | Host configuration | `[{"host": "cloud-native-mcp-server.local", "paths": [{"path": "/", "pathType": "Prefix"}]}]` |
| `ingress.tls` | TLS configuration | `[]` |

### RBAC Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `rbac.create` | Create RBAC rules | `true` |
| `rbac.useClusterAdmin` | Use cluster-admin role | `true` |
| `rbac.rules` | Custom RBAC rules | `[]` |

### Application Configuration

For detailed configuration options, refer to the `config` section in `values.yaml`, including:

- Server configuration (server)
- Logging configuration (logging)
- Kubernetes client configuration (kubernetes)
- Prometheus configuration (prometheus)
- Grafana configuration (grafana)
- Kibana configuration (kibana)
- Helm configuration (helm)
- Elasticsearch configuration (elasticsearch)
- Authentication configuration (auth)
- Audit configuration (audit)

## Advanced Usage

### Enable Autoscaling

```yaml
autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70
```

### Configure Persistent Storage

```yaml
persistence:
  enabled: true
  storageClass: "fast-ssd"
  accessModes: ["ReadWriteOnce"]
  size: 10Gi
```

### Enable Network Policy

```yaml
networkPolicy:
  enabled: true
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              name: monitoring
  egress:
    - to: []
```

### Configure Monitoring

```yaml
monitoring:
  enabled: true
  serviceMonitor:
    enabled: true
    namespace: "monitoring"
    interval: 30s
  prometheusRule:
    enabled: true
    rules:
      - alert: K8sMcpServerDown
        expr: up{job="cloud-native-mcp-server"} == 0
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "K8s MCP Server is down"
```

## Troubleshooting

### Check Deployment Status

```bash
# View Pod status
kubectl get pods -l app.kubernetes.io/name=cloud-native-mcp-server -n mcp-server

# View Service status
kubectl get svc -l app.kubernetes.io/name=cloud-native-mcp-server -n mcp-server

# View events
kubectl get events -n mcp-server --sort-by='.lastTimestamp'
```

### View Logs

```bash
# View all Pod logs
kubectl logs -l app.kubernetes.io/name=cloud-native-mcp-server -n mcp-server -f

# View specific Pod logs
kubectl logs -f deployment/cloud-native-mcp-server -n mcp-server
```

### Debug Deployment

```bash
# Check configuration
helm get values my-release -n mcp-server

# View rendered templates
helm get manifest my-release -n mcp-server

# Test configuration
helm template my-release . -f values.yaml --debug
```

### Common Issues

**Issue**: Pod cannot start
- Check image pull secret configuration
- Verify image tag is correct
- Check if resource limits are reasonable

**Issue**: Cannot access service
- Check service type and port configuration
- Verify Ingress configuration is correct
- Check if network policy blocks access

**Issue**: Authentication failed
- Verify API key configuration is correct
- Check authentication mode settings
- Validate RBAC permission configuration

## Development

### Build and Test Chart

```bash
# Check Chart syntax
helm lint .

# Test template rendering
helm template test-release . --debug

# Package Chart
helm package .

# Install local package
helm install my-release ./cloud-native-mcp-server-0.1.0.tgz
```

### Update Dependencies

```bash
# Add dependency
helm dependency add <chart-name>

# Update dependencies
helm dependency update

# List dependencies
helm dependency list
```

## Contributing

Issues and Pull Requests are welcome to improve this Helm chart.

## License

This project is open source under the MIT License.
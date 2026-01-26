# Deployment Guide

This guide covers various deployment strategies and best practices for the K8s MCP Server.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start Deployments](#quick-start-deployments)
- [Kubernetes Deployment](#kubernetes-deployment)
- [Docker Deployment](#docker-deployment)
- [Helm Deployment](#helm-deployment)
- [Production Considerations](#production-considerations)
- [Monitoring and Observability](#monitoring-and-observability)
- [Security Best Practices](#security-best-practices)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

### System Requirements

- **OS**: Linux, macOS, or Windows
- **CPU**: 1 core minimum, 2+ cores recommended
- **RAM**: 512MB minimum, 1GB+ recommended
- **Disk**: 100MB minimum
- **Network**: Access to Kubernetes cluster and configured services

### Software Requirements

- **Go**: 1.25+ (for building from source)
- **Docker**: 20.10+ (for containerized deployment)
- **kubectl**: Configured with cluster access
- **Helm**: 3.0+ (for Helm deployment)

### Service Dependencies

Optional services to connect to:

- **Grafana** (optional)
- **Prometheus** (optional)
- **Kibana** (optional)
- **Elasticsearch** (optional)
- **Alertmanager** (optional)
- **Jaeger** (optional)

---

## Quick Start Deployments

### Binary Deployment

```bash
# Download latest release
wget https://github.com/mahmut-Abi/k8s-mcp-server/releases/latest/download/k8s-mcp-server-linux-amd64
chmod +x k8s-mcp-server-linux-amd64

# Create config
cat > config.yaml << EOF
server:
  mode: "sse"
  addr: "0.0.0.0:8080"

logging:
  level: "info"

kubernetes:
  enabled: true
  kubeconfig: ""
EOF

# Run
./k8s-mcp-server-linux-amd64
```

### Docker Quick Start

```bash
docker run -d \
  --name k8s-mcp-server \
  -p 8080:8080 \
  -v ~/.kube:/root/.kube:ro \
  mahmutabi/k8s-mcp-server:latest
```

---

## Kubernetes Deployment

### Basic Deployment

Create a deployment manifest:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: k8s-mcp-server
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: k8s-mcp-server
  template:
    metadata:
      labels:
        app: k8s-mcp-server
    spec:
      serviceAccountName: k8s-mcp-server
      containers:
      - name: k8s-mcp-server
        image: mahmutabi/k8s-mcp-server:latest
        ports:
        - containerPort: 8080
        env:
        - name: MCP_MODE
          value: "sse"
        - name: MCP_ADDR
          value: "0.0.0.0:8080"
        - name: MCP_LOG_LEVEL
          value: "info"
        volumeMounts:
        - name: kubeconfig
          mountPath: /root/.kube
          readOnly: true
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
      volumes:
      - name: kubeconfig
        configMap:
          name: kubeconfig
```

### Service Account with RBAC

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: k8s-mcp-server
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: k8s-mcp-server
rules:
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["get", "list", "watch", "describe"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: k8s-mcp-server
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: k8s-mcp-server
subjects:
- kind: ServiceAccount
  name: k8s-mcp-server
  namespace: default
```

### Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: k8s-mcp-server
  namespace: default
spec:
  type: ClusterIP
  ports:
  - port: 8080
    targetPort: 8080
    protocol: TCP
  selector:
    app: k8s-mcp-server
```

### Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: k8s-mcp-server
  namespace: default
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: k8s-mcp.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: k8s-mcp-server
            port:
              number: 8080
```

### Deploy

```bash
# Apply all manifests
kubectl apply -f deploy/kubernetes/

# Verify deployment
kubectl get pods -l app=k8s-mcp-server
kubectl logs -l app=k8s-mcp-server

# Test connection
kubectl port-forward svc/k8s-mcp-server 8080:8080
curl http://localhost:8080/health
```

---

## Docker Deployment

### Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  k8s-mcp-server:
    image: mahmutabi/k8s-mcp-server:latest
    container_name: k8s-mcp-server
    ports:
      - "8080:8080"
    volumes:
      - ~/.kube:/root/.kube:ro
      - ./config.yaml:/app/config.yaml:ro
    environment:
      - MCP_MODE=sse
      - MCP_ADDR=0.0.0.0:8080
      - MCP_LOG_LEVEL=info
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    networks:
      - monitoring

networks:
  monitoring:
    external: true
```

### Run with Docker Compose

```bash
# Start
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down

# Restart
docker-compose restart
```

### Custom Docker Image

Build your own image:

```dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o k8s-mcp-server ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/k8s-mcp-server .

EXPOSE 8080

CMD ["./k8s-mcp-server"]
```

Build and push:

```bash
# Build
docker build -t your-registry/k8s-mcp-server:latest .

# Push
docker push your-registry/k8s-mcp-server:latest
```

---

## Helm Deployment

### Install from Chart Repository

```bash
# Add repository
helm repo add k8s-mcp https://mahmut-Abi.github.io/k8s-mcp-server

# Update repository
helm repo update

# Install
helm install k8s-mcp-server k8s-mcp/k8s-mcp-server

# Upgrade
helm upgrade k8s-mcp-server k8s-mcp/k8s-mcp-server

# Uninstall
helm uninstall k8s-mcp-server
```

### Custom Values

Create `values.yaml`:

```yaml
replicaCount: 2

image:
  repository: mahmutabi/k8s-mcp-server
  tag: latest
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 8080

ingress:
  enabled: true
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
  hosts:
    - host: k8s-mcp.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: k8s-mcp-tls
      hosts:
        - k8s-mcp.example.com

resources:
  requests:
    memory: "256Mi"
    cpu: "250m"
  limits:
    memory: "512Mi"
    cpu: "500m"

autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 80
  targetMemoryUtilizationPercentage: 80

config:
  server:
    mode: "sse"
    addr: "0.0.0.0:8080"
  logging:
    level: "info"
    format: "json"
  kubernetes:
    enabled: true
  grafana:
    enabled: true
    url: "http://grafana:3000"
    apiKey: "${GRAFANA_API_KEY}"
  prometheus:
    enabled: true
    address: "http://prometheus:9090"

rbac:
  create: true
  rules:
  - apiGroups: ["*"]
    resources: ["*"]
    verbs: ["get", "list", "watch", "describe"]

serviceAccount:
  create: true
  name: ""

podAnnotations: {}

podSecurityContext: {}
  # fsGroup: 2000

securityContext: {}
  # capabilities:
  #   drop:
  #   - ALL
  # readOnlyRootFilesystem: true
  # runAsNonRoot: true
  # runAsUser: 1000

nodeSelector: {}

tolerations: []

affinity: {}
```

### Install with Custom Values

```bash
helm install k8s-mcp-server ./deploy/helm/k8s-mcp-server -f values.yaml
```

---

## Production Considerations

### High Availability

Deploy multiple replicas with proper resource limits:

```yaml
replicaCount: 3

resources:
  requests:
    memory: "512Mi"
    cpu: "500m"
  limits:
    memory: "1Gi"
    cpu: "1000m"

autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 10
```

### Resource Optimization

Enable caching and tune parameters:

```yaml
config:
  cache:
    enabled: true
    type: "lru"
    max_size: 2000
    default_ttl: 300

  performance:
    max_response_size: 5242880
    compression_enabled: true
    json_pool_size: 200
```

### Security

1. **Enable Authentication**:
```yaml
config:
  auth:
    enabled: true
    mode: "apikey"
    api_key: "${MCP_API_KEY}"
```

2. **Use Secrets**:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: k8s-mcp-secrets
type: Opaque
stringData:
  mcp-api-key: "your-secret-key"
  grafana-api-key: "your-grafana-key"
```

3. **Network Policies**:
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: k8s-mcp-server
spec:
  podSelector:
    matchLabels:
      app: k8s-mcp-server
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector: {}
    ports:
    - protocol: TCP
      port: 8080
  egress:
  - to:
    - namespaceSelector: {}
    ports:
    - protocol: TCP
      port: 443
```

### Logging and Monitoring

```yaml
config:
  logging:
    level: "info"
    format: "json"
    output: "stdout"

  audit:
    enabled: true
    storage: "file"
    file_path: "/var/log/k8s-mcp-audit.log"
```

Add Prometheus monitoring:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: k8s-mcp-server
  namespace: default
  annotations:
    prometheus.io/scrape: "true"
    prometheus.io/port: "8080"
    prometheus.io/path: "/metrics"
spec:
  type: ClusterIP
  ports:
  - port: 8080
    targetPort: 8080
  selector:
    app: k8s-mcp-server
```

---

## Monitoring and Observability

### Health Checks

The server provides health endpoints:

```bash
# Basic health check
curl http://localhost:8080/health

# Detailed health
curl http://localhost:8080/health/detailed

# Readiness
curl http://localhost:8080/ready
```

### Metrics

Prometheus metrics are available at `/metrics`:

```bash
curl http://localhost:8080/metrics
```

Key metrics:
- `mcp_requests_total` - Total requests
- `mcp_request_duration_seconds` - Request duration
- `mcp_cache_hits_total` - Cache hits
- `mcp_cache_misses_total` - Cache misses
- `mcp_active_connections` - Active connections

### Logging

Structured JSON logs:

```json
{
  "level": "info",
  "timestamp": "2024-01-01T00:00:00Z",
  "message": "Starting K8s MCP Server",
  "version": "1.0.0",
  "mode": "sse"
}
```

### Audit Logs

Audit logs track all operations:

```json
{
  "timestamp": "2024-01-01T00:00:00Z",
  "request_id": "abc123",
  "tool": "kubernetes_list_resources_summary",
  "params": {"kind": "Pod"},
  "duration_ms": 123,
  "status": "success"
}
```

---

## Security Best Practices

### 1. Least Privilege RBAC

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: k8s-mcp-server
rules:
# Allow read-only access to most resources
- apiGroups: [""]
  resources: ["pods", "services", "configmaps", "secrets"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments", "replicasets"]
  verbs: ["get", "list", "watch"]
# Allow describe for troubleshooting
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["describe"]
```

### 2. Secrets Management

Use Kubernetes secrets for sensitive data:

```bash
kubectl create secret generic k8s-mcp-secrets \
  --from-literal=mcp-api-key='your-key' \
  --from-literal=grafana-api-key='your-grafana-key'
```

Mount secrets as environment variables:

```yaml
env:
- name: MCP_AUTH_API_KEY
  valueFrom:
    secretKeyRef:
      name: k8s-mcp-secrets
      key: mcp-api-key
```

### 3. Network Security

- Use TLS for external access
- Implement network policies
- Restrict ingress/egress traffic
- Use service mesh for mTLS

### 4. Pod Security

```yaml
securityContext:
  runAsNonRoot: true
  runAsUser: 1000
  fsGroup: 1000
  capabilities:
    drop:
    - ALL
  readOnlyRootFilesystem: true
```

### 5. Image Security

- Use signed images
- Scan images for vulnerabilities
- Keep images updated
- Use specific version tags

---

## Troubleshooting

### Common Issues

#### 1. Connection Refused

**Problem**: Cannot connect to server

**Solution**:
```bash
# Check pod status
kubectl get pods -l app=k8s-mcp-server

# Check logs
kubectl logs -l app=k8s-mcp-server

# Check service
kubectl get svc k8s-mcp-server

# Port forward test
kubectl port-forward svc/k8s-mcp-server 8080:8080
curl http://localhost:8080/health
```

#### 2. Authentication Failed

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

#### 3. Kubernetes API Access Denied

**Problem**: Cannot access Kubernetes API

**Solution**:
```bash
# Check RBAC
kubectl get clusterrole k8s-mcp-server -o yaml

# Check service account
kubectl get sa k8s-mcp-server

# Verify cluster role binding
kubectl get clusterrolebinding k8s-mcp-server

# Test permissions
kubectl auth can-i list pods --as=system:serviceaccount:default:k8s-mcp-server
```

#### 4. High Memory Usage

**Problem**: Pod OOMKilled

**Solution**:
```yaml
# Increase memory limits
resources:
  limits:
    memory: "1Gi"

# Reduce cache size
config:
  cache:
    max_size: 500

# Enable response compression
config:
  performance:
    compression_enabled: true
```

#### 5. Slow Response Times

**Problem**: Requests timing out

**Solution**:
```yaml
# Increase timeouts
kubernetes:
  timeout_sec: 60

# Enable caching
config:
  cache:
    enabled: true

# Use summary tools
# Replace kubernetes_list_resources with kubernetes_list_resources_summary
```

### Debug Mode

Enable debug logging:

```yaml
logging:
  level: "debug"
```

Or via environment:

```bash
export MCP_LOG_LEVEL=debug
```

### Health Check Script

```bash
#!/bin/bash

echo "Checking K8s MCP Server health..."

# Check endpoint
curl -f http://localhost:8080/health || exit 1

# Check metrics
curl -f http://localhost:8080/metrics > /dev/null || exit 1

# Check readiness
curl -f http://localhost:8080/ready || exit 1

echo "All checks passed!"
```

---

For more information, see:
- [Complete Tools Reference](TOOLS.md)
- [Configuration Guide](CONFIGURATION.md)
- [Main README](../README.md)
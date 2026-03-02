# Deployment Guide

This guide focuses on deployment methods that match the current server behavior and configuration schema.

For website docs:
- English: `website/content/en/docs/deployment.md`
- Chinese: `website/content/zh/docs/deployment.md`

---

## Prerequisites

- Kubernetes access (for Kubernetes/Helm deployment)
- Docker/Podman (for container deployment)
- A valid config file or environment variables for authentication and optional services

---

## Binary Deployment

```bash
# Example
./cloud-native-mcp-server --config=config.yaml
```

Minimal `config.yaml`:

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8080"

logging:
  level: "info"

kubernetes:
  kubeconfig: ""
```

Smoke check:

```bash
curl -sS http://127.0.0.1:8080/health
```

---

## Container Deployment

```bash
docker run -d \
  --name cloud-native-mcp-server \
  -p 8080:8080 \
  -v ~/.kube:/root/.kube:ro \
  -e MCP_MODE=sse \
  -e MCP_ADDR=0.0.0.0:8080 \
  -e MCP_LOG_LEVEL=info \
  -e MCP_AUTH_ENABLED=true \
  -e MCP_AUTH_MODE=apikey \
  -e MCP_AUTH_API_KEY='ChangeMe-Strong-Key-123!' \
  mahmutabi/cloud-native-mcp-server:latest
```

---

## Kubernetes Deployment

### Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cloud-native-mcp-server
  namespace: default
spec:
  replicas: 2
  selector:
    matchLabels:
      app: cloud-native-mcp-server
  template:
    metadata:
      labels:
        app: cloud-native-mcp-server
    spec:
      serviceAccountName: cloud-native-mcp-server
      containers:
        - name: cloud-native-mcp-server
          image: mahmutabi/cloud-native-mcp-server:latest
          ports:
            - containerPort: 8080
          env:
            - name: MCP_MODE
              value: "streamable-http"
            - name: MCP_ADDR
              value: "0.0.0.0:8080"
            - name: MCP_LOG_LEVEL
              value: "info"
            - name: MCP_AUTH_ENABLED
              value: "true"
            - name: MCP_AUTH_MODE
              value: "apikey"
            - name: MCP_AUTH_API_KEY
              valueFrom:
                secretKeyRef:
                  name: cloud-native-mcp-secrets
                  key: mcp-auth-api-key
          resources:
            requests:
              cpu: "250m"
              memory: "256Mi"
            limits:
              cpu: "1000m"
              memory: "1Gi"
          readinessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 10
```

### Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: cloud-native-mcp-server
  namespace: default
spec:
  selector:
    app: cloud-native-mcp-server
  ports:
    - protocol: TCP
      port: 8080
      targetPort: 8080
```

### Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: cloud-native-mcp-secrets
  namespace: default
type: Opaque
stringData:
  mcp-auth-api-key: "ChangeMe-Strong-Key-123!"
```

---

## Helm Deployment

```bash
helm install cloud-native-mcp-server ./deploy/helm/cloud-native-mcp-server \
  --set replicaCount=2
```

Recommended values (aligned with current config keys):

```yaml
replicaCount: 2

resources:
  requests:
    cpu: "250m"
    memory: "256Mi"
  limits:
    cpu: "1000m"
    memory: "1Gi"

config:
  server:
    mode: "streamable-http"
    addr: "0.0.0.0:8080"
    readTimeoutSec: 30
    writeTimeoutSec: 0
    idleTimeoutSec: 60

  logging:
    level: "info"
    json: true

  auth:
    enabled: true
    mode: "apikey"
    apiKey: "${MCP_AUTH_API_KEY}"

  kubernetes:
    timeoutSec: 30
    qps: 100.0
    burst: 200

  ratelimit:
    enabled: true
    requests_per_second: 100
    burst: 200
```

---

## Production Recommendations

### 1. Limit Scope of Enabled Services

```yaml
enableDisable:
  enabledServices: ["kubernetes", "prometheus", "grafana", "aggregate"]
```

### 2. Configure Timeouts and Rate Limit

```yaml
server:
  readTimeoutSec: 30
  writeTimeoutSec: 0
  idleTimeoutSec: 60

ratelimit:
  enabled: true
  requests_per_second: 100
  burst: 200
```

### 3. Secure Authentication Secrets

- Use Kubernetes `Secret` for `MCP_AUTH_API_KEY`
- Avoid committing plain-text credentials into `values.yaml`

### 4. Enable Audit Logs for Traceability

```yaml
audit:
  enabled: true
  storage: "file"
  format: "json"
  file:
    path: "/var/log/cloud-native-mcp-server/audit.log"
    maxSizeMB: 100
    maxBackups: 10
    maxAgeDays: 30
    compress: true
```

---

## Monitoring and Health

### Health Endpoint

```bash
curl -sS http://127.0.0.1:8080/health
```

### OpenAPI/Docs Endpoints

```bash
curl -sS http://127.0.0.1:8080/api/openapi.json
# Swagger UI:
# http://127.0.0.1:8080/api/docs
```

### Audit Endpoints (when audit is enabled)

```bash
curl -sS -H "X-API-Key: ${MCP_AUTH_API_KEY}" \
  "http://127.0.0.1:8080/api/audit/logs?limit=50"

curl -sS -H "X-API-Key: ${MCP_AUTH_API_KEY}" \
  "http://127.0.0.1:8080/api/audit/stats"
```

---

## Troubleshooting

### Config Not Applied

- Confirm startup logs include `Configuration loaded successfully from file`
- Verify the file path passed to `--config`

### Authentication Errors (401)

- Confirm `MCP_AUTH_ENABLED=true`
- Confirm `MCP_AUTH_MODE` and corresponding credential variable are set
- Confirm the client sends matching auth headers/credentials

### Slow Requests

- Increase `kubernetes.timeoutSec` for heavy operations
- Tune `kubernetes.qps` and `kubernetes.burst`
- Enable `ratelimit` to protect service under burst traffic

---

## Related Docs

- `docs/CONFIGURATION.md`
- `README.md`
- `website/content/en/docs/deployment.md`
- `website/content/zh/docs/deployment.md`

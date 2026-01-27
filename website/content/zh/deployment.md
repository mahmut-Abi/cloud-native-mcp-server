---
title: "部署指南"
---

# 部署指南

本指南涵盖 Cloud Native MCP Server 的各种部署策略和最佳实践。

## 目录

- [前提条件](#前提条件)
- [快速部署](#快速部署)
- [Kubernetes 部署](#kubernetes-部署)
- [Docker 部署](#docker-部署)
- [Helm 部署](#helm-部署)
- [生产环境考虑](#生产环境考虑)
- [监控和可观测性](#监控和可观测性)
- [安全最佳实践](#安全最佳实践)
- [故障排查](#故障排查)

---

## 前提条件

### 系统要求

- **操作系统**: Linux, macOS, 或 Windows
- **CPU**: 最低 1 核，推荐 2+ 核
- **内存**: 最低 512MB，推荐 1GB+
- **磁盘**: 最低 100MB
- **网络**: 可访问 Kubernetes 集群和配置的服务

### 软件要求

- **Go**: 1.25+ (从源码构建)
- **Docker**: 20.10+ (容器化部署)
- **kubectl**: 已配置集群访问
- **Helm**: 3.0+ (Helm 部署)

### 服务依赖

可选连接的服务：

- **Grafana** (可选)
- **Prometheus** (可选)
- **Kibana** (可选)
- **Elasticsearch** (可选)
- **Alertmanager** (可选)
- **Jaeger** (可选)
- **OpenTelemetry** (可选)

---

## 快速部署

### 二进制部署

```bash
# 下载最新版本
wget https://github.com/mahmut-Abi/cloud-native-mcp-server/releases/latest/download/cloud-native-mcp-server-linux-amd64
chmod +x cloud-native-mcp-server-linux-amd64

# 创建配置
cat > config.yaml << EOF
server:
  mode: "sse"
  addr: "0.0.0.0:8080"

logging:
  level: "info"

kubernetes:
  kubeconfig: ""
EOF

# 运行
./cloud-native-mcp-server-linux-amd64
```

### Docker 快速启动

```bash
docker run -d \
  --name cloud-native-mcp-server \
  -p 8080:8080 \
  -v ~/.kube:/root/.kube:ro \
  mahmutabi/cloud-native-mcp-server:latest
```

---

## Kubernetes 部署

### 基本部署

创建部署清单：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cloud-native-mcp-server
  namespace: default
spec:
  replicas: 1
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

### Service Account 和 RBAC

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cloud-native-mcp-server
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cloud-native-mcp-server
rules:
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["get", "list", "watch", "describe"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: cloud-native-mcp-server
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cloud-native-mcp-server
subjects:
- kind: ServiceAccount
  name: cloud-native-mcp-server
  namespace: default
```

### Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: cloud-native-mcp-server
  namespace: default
spec:
  type: ClusterIP
  ports:
  - port: 8080
    targetPort: 8080
    protocol: TCP
  selector:
    app: cloud-native-mcp-server
```

### Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: cloud-native-mcp-server
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
            name: cloud-native-mcp-server
            port:
              number: 8080
```

### 部署

```bash
# 应用所有清单
kubectl apply -f deploy/kubernetes/

# 验证部署
kubectl get pods -l app=cloud-native-mcp-server
kubectl logs -l app=cloud-native-mcp-server

# 测试连接
kubectl port-forward svc/cloud-native-mcp-server 8080:8080
curl http://localhost:8080/health
```

---

## Docker 部署

### Docker Compose

创建 `docker-compose.yml`:

```yaml
version: '3.8'

services:
  cloud-native-mcp-server:
    image: mahmutabi/cloud-native-mcp-server:latest
    container_name: cloud-native-mcp-server
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

### 使用 Docker Compose 运行

```bash
# 启动
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止
docker-compose down

# 重启
docker-compose restart
```

### 自定义 Docker 镜像

构建自己的镜像：

```dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o cloud-native-mcp-server ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/cloud-native-mcp-server .

EXPOSE 8080

CMD ["./cloud-native-mcp-server"]
```

构建和推送：

```bash
# 构建
docker build -t your-registry/cloud-native-mcp-server:latest .

# 推送
docker push your-registry/cloud-native-mcp-server:latest
```

---

## Helm 部署

### 从 Chart 仓库安装

```bash
# 添加仓库
helm repo add k8s-mcp https://mahmut-Abi.github.io/cloud-native-mcp-server

# 更新仓库
helm repo update

# 安装
helm install cloud-native-mcp-server k8s-mcp/cloud-native-mcp-server

# 升级
helm upgrade cloud-native-mcp-server k8s-mcp/cloud-native-mcp-server

# 卸载
helm uninstall cloud-native-mcp-server
```

### 自定义 Values

创建 `values.yaml`:

```yaml
replicaCount: 2

image:
  repository: mahmutabi/cloud-native-mcp-server
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
    kubeconfig: ""
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

nodeSelector: {}

tolerations: []

affinity: {}
```

### 使用自定义 Values 安装

```bash
helm install cloud-native-mcp-server ./deploy/helm/cloud-native-mcp-server -f values.yaml
```

---

## 生产环境考虑

### 高可用性

部署多个副本并设置适当的资源限制：

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

### 资源优化

启用缓存并调优参数：

```yaml
config:
  # 启用缓存
  cache:
    enabled: true
    type: "lru"
    max_size: 2000
    default_ttl: 300

  # 性能优化
  performance:
    max_response_size: 5242880
    compression_enabled: true
    json_pool_size: 200
```

### 安全

#### 1. 启用认证

```yaml
config:
  auth:
    enabled: true
    mode: "apikey"
    apiKey: "${MCP_AUTH_API_KEY}"
```

#### 2. 使用 Secrets

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

#### 3. 网络策略

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: cloud-native-mcp-server
spec:
  podSelector:
    matchLabels:
      app: cloud-native-mcp-server
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

### 日志和监控

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

添加 Prometheus 监控：

```yaml
apiVersion: v1
kind: Service
metadata:
  name: cloud-native-mcp-server
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
    app: cloud-native-mcp-server
```

---

## 监控和可观测性

### 健康检查

服务器提供健康检查端点：

```bash
# 基本健康检查
curl http://localhost:8080/health

# 详细健康信息
curl http://localhost:8080/health/detailed

# 就绪检查
curl http://localhost:8080/ready
```

### 指标

Prometheus 指标在 `/metrics` 端点可用：

```bash
curl http://localhost:8080/metrics
```

关键指标：
- `mcp_requests_total` - 总请求数
- `mcp_request_duration_seconds` - 请求持续时间
- `mcp_cache_hits_total` - 缓存命中数
- `mcp_cache_misses_total` - 缓存未命中数
- `mcp_active_connections` - 活动连接数

### 日志

结构化 JSON 日志：

```json
{
  "level": "info",
  "timestamp": "2024-01-01T00:00:00Z",
  "message": "Starting Cloud Native MCP Server",
  "version": "1.0.0",
  "mode": "sse"
}
```

### 审计日志

审计日志跟踪所有操作：

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

## 安全最佳实践

### 1. 最小权限 RBAC

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cloud-native-mcp-server
rules:
# 允许对大多数资源的只读访问
- apiGroups: [""]
  resources: ["pods", "services", "configmaps", "secrets"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments", "replicasets"]
  verbs: ["get", "list", "watch"]
# 允许 describe 用于故障排查
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["describe"]
```

### 2. 密钥管理

使用 Kubernetes secrets 存储敏感数据：

```bash
kubectl create secret generic k8s-mcp-secrets \
  --from-literal=mcp-api-key='your-key' \
  --from-literal=grafana-api-key='your-grafana-key'
```

将 secrets 挂载为环境变量：

```yaml
env:
- name: MCP_AUTH_API_KEY
  valueFrom:
    secretKeyRef:
      name: k8s-mcp-secrets
      key: mcp-api-key
```

### 3. 网络安全

- 对外部访问使用 TLS
- 实施网络策略
- 限制入口/出口流量
- 使用服务网格进行 mTLS

### 4. Pod 安全

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

### 5. 镜像安全

- 使用签名镜像
- 扫描镜像漏洞
- 保持镜像更新
- 使用特定版本标签

---

## 故障排查

### 常见问题

#### 1. 连接被拒绝

**问题**: 无法连接到服务器

**解决方案**:
```bash
# 检查 pod 状态
kubectl get pods -l app=cloud-native-mcp-server

# 检查日志
kubectl logs -l app=cloud-native-mcp-server

# 检查 service
kubectl get svc cloud-native-mcp-server

# 端口转发测试
kubectl port-forward svc/cloud-native-mcp-server 8080:8080
curl http://localhost:8080/health
```

#### 2. 认证失败

**问题**: 401 Unauthorized

**解决方案**:
```bash
# 检查认证配置
kubectl get configmap k8s-mcp-config -o yaml

# 验证 secrets
kubectl get secret k8s-mcp-secrets -o yaml

# 使用正确的头部测试
curl -H "X-API-Key: your-key" http://localhost:8080/health
```

#### 3. Kubernetes API 访问被拒绝

**问题**: 无法访问 Kubernetes API

**解决方案**:
```bash
# 检查 RBAC
kubectl get clusterrole cloud-native-mcp-server -o yaml

# 检查 service account
kubectl get sa cloud-native-mcp-server

# 验证 cluster role binding
kubectl get clusterrolebinding cloud-native-mcp-server

# 测试权限
kubectl auth can-i list pods --as=system:serviceaccount:default:cloud-native-mcp-server
```

#### 4. 高内存使用

**问题**: Pod OOMKilled

**解决方案**:
```yaml
# 增加内存限制
resources:
  limits:
    memory: "1Gi"

# 减少缓存大小
config:
  cache:
    max_size: 500

# 启用响应压缩
config:
  performance:
    compression_enabled: true
```

#### 5. 响应慢

**问题**: 请求超时

**解决方案**:
```yaml
# 增加超时
kubernetes:
  timeoutSec: 60

# 启用缓存
config:
  cache:
    enabled: true

# 使用摘要工具
# 用 kubernetes_list_resources_summary 替换 kubernetes_list_resources
```

### 调试模式

启用调试日志：

```yaml
logging:
  level: "debug"
```

或通过环境变量：

```bash
export MCP_LOG_LEVEL=debug
```

### 健康检查脚本

```bash
#!/bin/bash

echo "检查 Cloud Native MCP Server 健康状况..."

# 检查端点
curl -f http://localhost:8080/health || exit 1

# 检查指标
curl -f http://localhost:8080/metrics > /dev/null || exit 1

# 检查就绪状态
curl -f http://localhost:8080/ready || exit 1

echo "所有检查通过！"
```

---

## 相关文档

- [完整工具参考](/docs/tools/)
- [配置指南](/docs/configuration/)
- [安全指南](/docs/security/)
---
title: "部署"
weight: 20
---

# 部署指南

本指南涵盖 Cloud Native MCP Server 的各种部署策略和最佳实践。

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

## 部署方式

Cloud Native MCP Server 支持多种部署方式：

- **二进制部署** - 直接下载可执行文件运行
- **Docker 部署** - 使用 Docker 容器运行
- **Kubernetes 部署** - 在 Kubernetes 集群中部署
- **Helm 部署** - 使用 Helm Chart 部署

## 内容

- [Kubernetes 部署](/zh/guides/deployment/kubernetes/) - 在 Kubernetes 集群中部署
- [Docker 部署](/zh/guides/deployment/docker/) - 使用 Docker 容器部署
- [Helm 部署](/zh/guides/deployment/helm/) - 使用 Helm Chart 部署
- [生产环境考虑](#生产环境考虑) - 高可用性和安全性
- [监控和可观测性](#监控和可观测性) - 健康检查和指标
- [故障排查](#故障排查) - 常见问题和解决方案

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
```

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

## 相关文档

- [配置指南](/zh/guides/configuration/)
- [安全指南](/zh/guides/security/)
- [性能指南](/zh/guides/performance/)
- [架构指南](/zh/concepts/architecture/)
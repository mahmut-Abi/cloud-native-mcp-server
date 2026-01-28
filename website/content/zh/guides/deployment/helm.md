---
title: "Helm 部署"
weight: 30
---

# Helm 部署

本指南描述如何使用 Helm Chart 部署 Cloud Native MCP Server。

## 前提条件

- Helm 3.0+ 已安装
- Kubernetes 集群已配置
- kubectl 已安装并配置

## 添加仓库

```bash
# 添加仓库
helm repo add k8s-mcp https://mahmut-Abi.github.io/cloud-native-mcp-server

# 更新仓库
helm repo update

# 搜索 Chart
helm search repo k8s-mcp
```

## 基本安装

### 默认安装

```bash
# 安装
helm install cloud-native-mcp-server k8s-mcp/cloud-native-mcp-server

# 升级
helm upgrade cloud-native-mcp-server k8s-mcp/cloud-native-mcp-server

# 卸载
helm uninstall cloud-native-mcp-server
```

### 指定命名空间

```bash
# 安装到指定命名空间
helm install cloud-native-mcp-server k8s-mcp/cloud-native-mcp-server \
  --namespace monitoring \
  --create-namespace
```

## 自定义 Values

### 创建 values.yaml

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
  auth:
    enabled: true
    mode: "apikey"
    apiKey: "${MCP_AUTH_API_KEY}"

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
helm install cloud-native-mcp-server k8s-mcp/cloud-native-mcp-server \
  -f values.yaml \
  --namespace monitoring \
  --create-namespace
```

## 配置选项

### 镜像配置

```yaml
image:
  repository: mahmutabi/cloud-native-mcp-server
  tag: "v1.0.0"
  pullPolicy: IfNotPresent
  pullSecrets: []
```

### 副本配置

```yaml
replicaCount: 3

autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 10
  targetCPUUtilizationPercentage: 80
  targetMemoryUtilizationPercentage: 80
```

### 资源配置

```yaml
resources:
  requests:
    memory: "512Mi"
    cpu: "500m"
  limits:
    memory: "1Gi"
    cpu: "1000m"
```

### 服务配置

```yaml
service:
  type: ClusterIP
  port: 8080
  annotations:
    prometheus.io/scrape: "true"
    prometheus.io/port: "8080"
    prometheus.io/path: "/metrics"
```

### Ingress 配置

```yaml
ingress:
  enabled: true
  className: "nginx"
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
  hosts:
    - host: k8s-mcp.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: k8s-mcp-tls
      hosts:
        - k8s-mcp.example.com
```

### 配置文件

```yaml
config:
  server:
    mode: "sse"
    addr: "0.0.0.0:8080"
  logging:
    level: "info"
    format: "json"
  kubernetes:
    kubeconfig: ""
  # ... 其他服务配置
```

### Secrets

```yaml
secrets:
  create: true
  mcpApiKey: "your-api-key"
  grafanaApiKey: "your-grafana-key"
  # 或者使用现有 Secret
  existingSecret: "mcp-secrets"
```

### RBAC

```yaml
rbac:
  create: true
  rules:
  - apiGroups: [""]
    resources: ["pods", "services", "configmaps", "secrets"]
    verbs: ["get", "list", "watch"]
  - apiGroups: ["apps"]
    resources: ["deployments", "replicasets"]
    verbs: ["get", "list", "watch"]
```

### ServiceAccount

```yaml
serviceAccount:
  create: true
  name: "cloud-native-mcp-server"
  annotations: {}
```

### 节点选择

```yaml
nodeSelector:
  node.kubernetes.io/instance-type: "m5.large"

tolerations:
  - key: "workload"
    operator: "Equal"
    value: "monitoring"
    effect: "NoSchedule"

affinity:
  podAntiAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
    - labelSelector:
        matchExpressions:
        - key: app
          operator: In
          values:
          - cloud-native-mcp-server
      topologyKey: "kubernetes.io/hostname"
```

## 高级用法

### 使用现有 ConfigMap

```yaml
configMap:
  create: false
  name: "mcp-config"
```

### 使用现有 Secret

```yaml
secret:
  create: false
  name: "mcp-secrets"
```

### 环境变量

```yaml
env:
  - name: CUSTOM_VAR
    value: "custom-value"
  - name: SECRET_VAR
    valueFrom:
      secretKeyRef:
        name: mcp-secrets
        key: secret-var
```

### 卷挂载

```yaml
volumes:
  - name: kubeconfig
    configMap:
      name: kubeconfig
      optional: true

volumeMounts:
  - name: kubeconfig
    mountPath: /root/.kube
    readOnly: true
```

## 监控配置

### PodMonitor

```yaml
podMonitor:
  enabled: true
  interval: 30s
  scrapeTimeout: 10s
  namespace: monitoring
```

### ServiceMonitor

```yaml
serviceMonitor:
  enabled: true
  interval: 30s
  scrapeTimeout: 10s
  namespace: monitoring
```

## 故障排查

### 检查安装状态

```bash
helm status cloud-native-mcp-server -n monitoring

# 查看历史版本
helm history cloud-native-mcp-server -n monitoring

# 获取 Values
helm get values cloud-native-mcp-server -n monitoring

# 获取所有信息
helm get all cloud-native-mcp-server -n monitoring
```

### 调试安装

```bash
# Dry run
helm install cloud-native-mcp-server k8s-mcp/cloud-native-mcp-server \
  -f values.yaml \
  --namespace monitoring \
  --dry-run \
  --debug

# 渲染模板
helm template cloud-native-mcp-server k8s-mcp/cloud-native-mcp-server \
  -f values.yaml \
  --namespace monitoring
```

### 回滚

```bash
# 查看历史
helm history cloud-native-mcp-server -n monitoring

# 回滚到上一个版本
helm rollback cloud-native-mcp-server -n monitoring

# 回滚到特定版本
helm rollback cloud-native-mcp-server 2 -n monitoring
```

## 相关文档

- [Kubernetes 部署](/zh/guides/deployment/kubernetes/)
- [Docker 部署](/zh/guides/deployment/docker/)
- [配置指南](/zh/guides/configuration/)
- [Helm 文档](https://helm.sh/docs/)
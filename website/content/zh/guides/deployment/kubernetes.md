---
title: "Kubernetes 部署"
weight: 10
---

# Kubernetes 部署

本指南描述如何在 Kubernetes 集群中部署 Cloud Native MCP Server。

## 前提条件

- 已配置的 Kubernetes 集群
- kubectl 已安装并配置
- 适当的集群权限

## 基本部署

### Deployment

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

### ServiceAccount 和 RBAC

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

## 部署步骤

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

## 高可用配置

### 多副本部署

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cloud-native-mcp-server
spec:
  replicas: 3
  # ... 其他配置
```

### 自动扩缩容

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: cloud-native-mcp-server
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: cloud-native-mcp-server
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 80
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

## 配置管理

### 使用 ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: mcp-config
data:
  config.yaml: |
    server:
      mode: "sse"
      addr: "0.0.0.0:8080"
    logging:
      level: "info"
    kubernetes:
      kubeconfig: ""
    grafana:
      enabled: true
      url: "http://grafana:3000"
      apiKey: "${GRAFANA_API_KEY}"
```

### 使用 Secret

```bash
kubectl create secret generic mcp-secrets \
  --from-literal=mcp-api-key='your-key' \
  --from-literal=grafana-api-key='your-grafana-key'
```

```yaml
env:
- name: MCP_AUTH_API_KEY
  valueFrom:
    secretKeyRef:
      name: mcp-secrets
      key: mcp-api-key
- name: GRAFANA_API_KEY
  valueFrom:
    secretKeyRef:
      name: mcp-secrets
      key: grafana-api-key
```

## 安全配置

### Pod 安全

```yaml
spec:
  securityContext:
    runAsNonRoot: true
    runAsUser: 1000
    fsGroup: 1000
  containers:
  - name: cloud-native-mcp-server
    securityContext:
      allowPrivilegeEscalation: false
      capabilities:
        drop:
        - ALL
      readOnlyRootFilesystem: true
```

### 网络策略

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

## 监控配置

### Service 监控

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
  # ... 其他配置
```

## 故障排查

### 检查 Pod 状态

```bash
kubectl get pods -l app=cloud-native-mcp-server
kubectl describe pod <pod-name>
kubectl logs <pod-name>
```

### 检查服务

```bash
kubectl get svc cloud-native-mcp-server
kubectl describe svc cloud-native-mcp-server
```

### 检查 RBAC

```bash
kubectl get sa cloud-native-mcp-server
kubectl get clusterrole cloud-native-mcp-server
kubectl get clusterrolebinding cloud-native-mcp-server
kubectl auth can-i list pods --as=system:serviceaccount:default:cloud-native-mcp-server
```

## 相关文档

- [Docker 部署](/zh/guides/deployment/docker/)
- [Helm 部署](/zh/guides/deployment/helm/)
- [配置指南](/zh/guides/configuration/)
- [安全指南](/zh/guides/security/)
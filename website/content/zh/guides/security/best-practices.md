---
title: "安全最佳实践"
weight: 30
---

# 安全最佳实践

本文档描述 Cloud Native MCP Server 的安全最佳实践。

## 1. 使用强认证

### API Key 要求

- 最少 16 字符
- 包含至少 3 种字符类型：
  - 大写字母 (A-Z)
  - 小写字母 (a-z)
  - 数字 (0-9)
  - 特殊字符 (!@#$%^&*()_+-=[]{}|;:,.<>?)

### 配置示例

```yaml
auth:
  enabled: true
  mode: "apikey"
  apiKey: "Abc123!@#Xyz789!@#Abc123!@#"
```

## 2. 定期轮换密钥

### 自动轮换配置

```yaml
secrets:
  # 自动轮换间隔（小时）
  rotation_interval: 168  # 7 天

  # 密钥过期时间（天）
  max_age: 30

  # 保留过期密钥（用于审计）
  keep_expired: true
```

### 手动轮换步骤

1. 生成新密钥
2. 更新配置
3. 测试新密钥
4. 删除旧密钥

## 3. 使用 Kubernetes Secrets

### 创建 Secret

```bash
kubectl create secret generic mcp-secrets \
  --from-literal=api-key='Abc123!@#Xyz789!@#' \
  --from-literal=jwt-secret='your-jwt-secret'
```

### 在部署中使用

```yaml
env:
- name: MCP_AUTH_API_KEY
  valueFrom:
    secretKeyRef:
      name: mcp-secrets
      key: api-key
```

## 4. 永远不要硬编码凭据

### 不好的做法

```yaml
auth:
  apiKey: "Abc123!@#Xyz789!@#"
```

### 好的做法

```yaml
auth:
  apiKey: "${MCP_AUTH_API_KEY}"
```

## 5. 启用审计日志

### 配置审计日志

```yaml
audit:
  enabled: true
  storage: "database"
  database:
    type: "sqlite"
    sqlitePath: "/var/lib/cloud-native-mcp-server/audit.db"
    maxRecords: 100000
  masking:
    enabled: true
    maskValue: "***REDACTED***"
```

### 查询审计日志

```bash
# 查询最近的失败认证
curl -H "X-API-Key: your-key" \
  "http://localhost:8080/api/audit/query?status=failed&limit=50"

# 查询特定用户的操作
curl -H "X-API-Key: your-key" \
  "http://localhost:8080/api/audit/query?user=admin&limit=100"
```

## 6. 生产环境使用 HTTPS

### TLS 配置

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8443"
  tls:
    certFile: "/path/to/cert.pem"
    keyFile: "/path/to/key.pem"
    minVersion: "TLS1.2"
    maxVersion: "TLS1.3"
```

### mTLS 配置

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8443"
  tls:
    certFile: "/path/to/server-cert.pem"
    keyFile: "/path/to/server-key.pem"
    clientAuth: "RequireAndVerifyClientCert"
    caFile: "/path/to/ca-cert.pem"
```

## 7. 限制访问

### 防火墙规则

```bash
# 使用 iptables
iptables -A INPUT -p tcp --dport 8080 -s 10.0.0.0/8 -j ACCEPT
iptables -A INPUT -p tcp --dport 8080 -j DROP
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
    - namespaceSelector:
        matchLabels:
          name: monitoring
    ports:
    - protocol: TCP
      port: 8080
```

### 速率限制

```yaml
ratelimit:
  enabled: true
  requests_per_second: 100
  burst: 200
```

## 8. 监控可疑活动

### 告警配置

```yaml
monitoring:
  # 失败认证告警阈值
  auth_failure_threshold: 5
  auth_failure_window: 300  # 5 分钟

  # 异常行为检测
  anomaly_detection:
    enabled: true
    sensitivity: "medium"
```

### 监控指标

- 失败的认证尝试
- 速率限制触发次数
- 异常请求模式
- 资源使用异常

## 9. 保持依赖更新

```bash
# 更新依赖
go get -u ./...
go mod tidy

# 检查漏洞
go list -json -m all | nancy sleuth

# 自动更新
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

## 10. 实施最小权限原则

### RBAC 配置

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cloud-native-mcp-server
rules:
- apiGroups: [""]
  resources: ["pods", "services"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments"]
  verbs: ["get", "list", "watch"]
```

### 容器安全

```yaml
securityContext:
  runAsNonRoot: true
  runAsUser: 1000
  readOnlyRootFilesystem: true
  capabilities:
    drop:
    - ALL
```

## 11. 网络隔离

### 使用网络策略

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
    - namespaceSelector:
        matchLabels:
          name: monitoring
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

### 使用服务网格

- Istio
- Linkerd
- Consul Connect

## 12. 容器安全

### 非 root 用户

```yaml
securityContext:
  runAsNonRoot: true
  runAsUser: 1000
  fsGroup: 1000
```

### 只读文件系统

```yaml
securityContext:
  readOnlyRootFilesystem: true
```

### 删除特权

```yaml
securityContext:
  allowPrivilegeEscalation: false
  capabilities:
    drop:
    - ALL
```

## 13. 镜像安全

### 使用签名镜像

```bash
# 验证镜像签名
cosign verify mahmutabi/cloud-native-mcp-server:latest
```

### 扫描镜像漏洞

```bash
# 使用 Trivy
trivy image mahmutabi/cloud-native-mcp-server:latest

# 使用 Clair
clairctl analyze mahmutabi/cloud-native-mcp-server:latest
```

### 使用特定版本标签

```yaml
image:
  repository: mahmutabi/cloud-native-mcp-server
  tag: "v1.0.0"  # 使用具体版本
```

## 14. 数据加密

### 静态数据加密

```yaml
# 使用加密的 ConfigMap
apiVersion: v1
kind: Secret
metadata:
  name: mcp-secrets
type: Opaque
stringData:
  api-key: "Abc123!@#Xyz789!@#"
```

### 传输中加密

```yaml
server:
  tls:
    certFile: "/path/to/cert.pem"
    keyFile: "/path/to/key.pem"
    minVersion: "TLS1.2"
```

## 15. 备份和恢复

### 定期备份

```bash
# 备份配置
kubectl get configmap mcp-config -o yaml > backup-config.yaml

# 备份 Secrets
kubectl get secret mcp-secrets -o yaml > backup-secrets.yaml
```

### 恢复流程

```bash
# 恢复配置
kubectl apply -f backup-config.yaml

# 恢复 Secrets
kubectl apply -f backup-secrets.yaml
```

## 相关文档

- [认证](/zh/guides/security/authentication/)
- [密钥管理](/zh/guides/security/secrets/)
- [配置指南](/zh/guides/configuration/authentication/)
- [部署指南](/zh/guides/deployment/)
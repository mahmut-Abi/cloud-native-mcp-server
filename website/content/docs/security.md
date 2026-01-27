---
title: "安全指南"
---

# 安全指南

本文档描述了 Cloud Native MCP Server 的安全特性和最佳实践。

## 目录

- [认证](#认证)
- [密钥管理](#密钥管理)
- [输入清理](#输入清理)
- [审计日志](#审计日志)
- [安全最佳实践](#安全最佳实践)
- [安全头部](#安全头部)
- [报告安全问题](#报告安全问题)

---

## 认证

### API Key 认证

API 密钥必须满足以下复杂度要求：

- **最小长度**: 16 个字符
- **字符类型**: 以下 4 种类型中至少包含 3 种：
  - 大写字母 (A-Z)
  - 小写字母 (a-z)
  - 数字 (0-9)
  - 特殊字符 (!@#$%^&*()_+-=[]{}|;:,.<>?)

**有效示例**:
- `Abc123!@#Xyz789!@#` (大写、小写、数字、特殊字符)
- `Abc123Xyz789Abc123` (大写、小写、数字)
- `ABC123!@#XYZ789!@#` (大写、数字、特殊字符)

**无效示例**:
- `Abc123!@#` (少于 16 个字符)
- `abcdefgh12345678` (只有小写和数字，不满足 3 种字符类型)
- `ABCDEFGHIJKLMNOPQRSTUVWXYZ` (只有大写)

### Bearer Token 认证

Bearer token 必须遵循 JWT 结构：

- **格式**: `header.payload.signature`
- **最小长度**: 32 个字符
- **编码**: Base64URL 编码的各部分
- **验证**: 每个部分必须只包含有效的 base64url 字符 (A-Z, a-z, 0-9, -, _, +)

**有效示例**:
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```

**无效示例**:
- `abcdefgh12345678abcdefgh12345678` (没有 JWT 结构)
- `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ` (少于 32 个字符)
- `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c$` (末尾有无效字符)

### Basic 认证

Basic 认证使用用户名和密码：

- **用户名**: 非空字符串
- **密码**: 非空字符串

**示例**:
```bash
curl -u admin:secret http://localhost:8080/api/aggregate/sse
```

### 配置认证

在配置文件中启用认证：

```yaml
auth:
  # 启用/禁用认证
  enabled: true

  # 认证模式: apikey | bearer | basic
  mode: "apikey"

  # API Key（用于 apikey 模式）
  # 最少 16 字符，推荐 32+ 字符
  apiKey: "Abc123!@#Xyz789!@#Abc123!@#"

  # Bearer token（用于 bearer 模式）
  bearerToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

  # Basic Auth 用户名
  username: "admin"

  # Basic Auth 密码
  password: "secret-password"

  # JWT 密钥（用于 JWT 验证）
  jwtSecret: "your-jwt-secret-key"

  # JWT 算法 (HS256, RS256, etc.)
  jwtAlgorithm: "HS256"
```

### 环境变量认证

使用环境变量配置认证：

```bash
export MCP_AUTH_ENABLED=true
export MCP_AUTH_MODE=apikey
export MCP_AUTH_API_KEY="Abc123!@#Xyz789!@#Abc123!@#"

./cloud-native-mcp-server
```

---

## 密钥管理

服务器包含密钥管理模块，用于安全地存储凭据。

### 特性

- **安全存储**: 带过期支持的内存存储
- **密钥轮换**: API 密钥和 bearer token 的自动轮换
- **密钥生成**: 内置生成器，用于复杂的 API 密钥和 JWT 类型的 token
- **环境变量**: 支持从环境变量加载密钥
- **密钥类型**: API 密钥、bearer token、basic auth 凭据

### 使用密钥管理器

```go
import "github.com/mahmut-Abi/cloud-native-mcp-server/internal/secrets"

// 创建新的密钥管理器
manager := secrets.NewInMemoryManager()

// 存储密钥
secret := &secrets.Secret{
    Type:  secrets.SecretTypeAPIKey,
    Name:  "my-api-key",
    Value: "Abc123!@#Xyz789!@#",
}
manager.Store(secret)

// 检索密钥
retrieved, err := manager.Retrieve(secret.ID)

// 轮换密钥
rotated, err := manager.Rotate(secret.ID)

// 生成新的 API 密钥
newKey, err := manager.GenerateAPIKey("my-new-key")

// 生成新的 bearer token
newToken, err := manager.GenerateBearerToken("my-new-token")
```

### 密钥过期

密钥可以有过期时间：

```go
secret := &secrets.Secret{
    Type:       secrets.SecretTypeAPIKey,
    Name:       "temporary-key",
    Value:      "Abc123!@#Xyz789!@#",
    ExpiresAt:  time.Now().Add(24 * time.Hour), // 24 小时后过期
}
```

过期的密钥会自动从列表中排除，无法检索。

### 密钥轮换策略

定期轮换密钥是安全最佳实践：

```yaml
auth:
  enabled: true
  mode: "apikey"
  apiKey: "${MCP_AUTH_API_KEY}"

secrets:
  # 自动轮换间隔（小时）
  rotation_interval: 168  # 7 天

  # 密钥过期时间（天）
  max_age: 30

  # 保留过期密钥（用于审计）
  keep_expired: true
```

---

## 输入清理

所有用户输入都经过清理，以防止注入攻击。

### 清理特性

- **过滤值**: 移除危险字符（SQL 注入、XSS、命令注入）
- **URL 验证**: 只允许 http/https 协议用于 web 获取
- **长度限制**: 最大字符串长度强制（1000 个字符）
- **特殊字符移除**: 移除分号、引号和其他注入向量

### 清理规则

以下字符会从用户输入中移除：

- **SQL 注入**: `;`, `'`, `"`, `--`, `/*`, `*/`
- **命令注入**: `|`, `&`, `$`, `(`, `)`, `<`, `>`, `\``, `\`
- **XSS**: `<script>`, `javascript:`, `onload=`, `onerror=`

### 示例

```go
import "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/sanitize"

// 清理过滤值
cleanValue := sanitize.SanitizeFilterValue("'; DROP TABLE users; --")
// 结果: " DROP TABLE users "

// 清理 JSONPath
cleanPath := sanitize.SanitizeJSONPath("$.data[*].name; rm -rf /")
// 结果: "$.data[*].name rm -rf "

// 验证字符串
isValid := sanitize.ValidateString("normal input")
// 结果: true
```

### 配置输入清理

```yaml
sanitization:
  # 启用输入清理
  enabled: true

  # 最大字符串长度
  max_length: 1000

  # 允许的 URL 协议
  allowed_protocols:
    - http
    - https

  # 禁用的字符模式
  blocked_patterns:
    - "';"
    - "DROP TABLE"
    - "rm -rf"
    - "<script>"
```

---

## 审计日志

审计日志跟踪所有操作，有助于安全监控和合规性。

### 启用审计日志

```yaml
audit:
  enabled: true
  level: "info"
  storage: "database"
  format: "json"

  # 敏感数据掩码
  masking:
    enabled: true
    fields:
      - password
      - token
      - apiKey
      - secret
      - authorization
    maskValue: "***REDACTED***"

  # 存储配置
  database:
    type: "sqlite"
    sqlitePath: "/var/lib/cloud-native-mcp-server/audit.db"
    maxRecords: 100000
```

### 审计事件

以下事件会被记录：

- 认证成功/失败
- 工具调用
- 配置更改
- 错误和异常
- 访问拒绝

### 审计日志格式

```json
{
  "timestamp": "2024-01-01T00:00:00Z",
  "request_id": "abc123",
  "user": "admin",
  "tool": "kubernetes_list_pods",
  "params": {
    "namespace": "default"
  },
  "duration_ms": 123,
  "status": "success",
  "error": ""
}
```

### 查询审计日志

```bash
# 查询最近 100 条审计日志
curl -H "X-API-Key: your-key" \
  "http://localhost:8080/api/audit/query?limit=100"

# 查询特定用户的审计日志
curl -H "X-API-Key: your-key" \
  "http://localhost:8080/api/audit/query?user=admin&limit=50"

# 查询失败的认证尝试
curl -H "X-API-Key: your-key" \
  "http://localhost:8080/api/audit/query?tool=auth_login&status=failed"
```

---

## 安全最佳实践

### 1. 使用强认证

- 始终使用满足复杂度要求的 API 密钥
- 定期轮换 API 密钥
- 使用 bearer token 进行基于 JWT 的认证
- 永远不要将凭据提交到版本控制

### 2. 启用审计日志

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

### 3. 生产环境使用 HTTPS

在生产环境部署时始终使用 HTTPS：

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8443"
  tls:
    certFile: "/path/to/cert.pem"
    keyFile: "/path/to/key.pem"
```

### 4. 限制访问

- 使用防火墙规则限制对服务器的访问
- 在 Kubernetes 中实施网络策略
- 使用 RBAC 控制 Kubernetes 资源的访问
- 实施速率限制防止暴力破解

```yaml
ratelimit:
  enabled: true
  requests_per_second: 100
  burst: 200
```

### 5. 监控可疑活动

- 启用指标和监控
- 为失败的认证尝试设置告警
- 定期审查审计日志
- 实施异常检测

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

### 6. 保持依赖更新

定期更新依赖以修补安全漏洞：

```bash
go get -u ./...
go mod tidy
```

### 7. 使用 Kubernetes Secrets

永远不要在配置文件中硬编码敏感信息：

```yaml
# 不好的做法
auth:
  apiKey: "Abc123!@#Xyz789!@#"

# 好的做法
auth:
  apiKey: "${MCP_AUTH_API_KEY}"
```

创建 Kubernetes Secret：

```bash
kubectl create secret generic mcp-secrets \
  --from-literal=api-key='Abc123!@#Xyz789!@#' \
  --from-literal=jwt-secret='your-jwt-secret'
```

在部署中引用：

```yaml
env:
- name: MCP_AUTH_API_KEY
  valueFrom:
    secretKeyRef:
      name: mcp-secrets
      key: api-key
```

### 8. 实施最小权限原则

- 只授予必要的权限
- 使用 RBAC 限制 Kubernetes 访问
- 定期审查和更新权限
- 使用服务账号隔离

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cloud-native-mcp-server
rules:
- apiGroups: [""]
  resources: ["pods", "services"]
  verbs: ["get", "list", "watch"]
```

### 9. 网络隔离

- 使用网络策略限制 Pod 间通信
- 在不同命名空间中隔离服务
- 使用 Ingress 控制器管理外部访问
- 考虑使用服务网格进行 mTLS

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

### 10. 容器安全

- 使用非 root 用户运行容器
- 使用只读文件系统
- 删除不必要的特权
- 扫描镜像漏洞

```yaml
securityContext:
  runAsNonRoot: true
  runAsUser: 1000
  readOnlyRootFilesystem: true
  capabilities:
    drop:
    - ALL
```

---

## 安全头部

服务器会自动过滤调试日志中的敏感头部：

- `Authorization`
- `Cookie`
- `X-API-Key`
- `X-Api-Key`
- `x-api-key`

这些头部永远不会以明文形式记录。

### 自定义安全头部

```yaml
security:
  # 额外的安全头部
  headers:
    X-Frame-Options: "DENY"
    X-Content-Type-Options: "nosniff"
    X-XSS-Protection: "1; mode=block"
    Strict-Transport-Security: "max-age=31536000; includeSubDomains"
    Content-Security-Policy: "default-src 'self'"

  # 头部过滤
  header_filtering:
    enabled: true
    filtered_headers:
      - authorization
      - cookie
      - x-api-key
      - x-auth-token
```

---

## TLS/SSL 配置

在生产环境中使用 TLS/SSL 加密通信：

### 基本 TLS 配置

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

### Let's Encrypt 集成

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8443"
  tls:
    autoCert:
      enabled: true
      host: "k8s-mcp.example.com"
      email: "admin@example.com"
      cacheDir: "/var/lib/letsencrypt"
```

---

## 速率限制

防止暴力攻击和滥用：

```yaml
ratelimit:
  enabled: true
  requests_per_second: 100
  burst: 200
  cleanup_interval: 60

  # 特定客户端限制
  client_limits:
    default:
      requests_per_second: 100
    authenticated:
      requests_per_second: 200

  # 白名单
  whitelist:
    - "10.0.0.0/8"
    - "192.168.0.0/16"

  # 黑名单
  blacklist:
    - " malicious.example.com"
```

---

## 报告安全问题

如果您发现安全漏洞，请私下报告：

- **Email**: security@example.com
- **GitHub Security Advisories**: https://github.com/mahmut-Abi/cloud-native-mcp-server/security/advisories

请不要为安全漏洞公开创建 issue。

### 安全披露流程

1. 通过私密渠道报告漏洞
2. 我们会在 48 小时内确认收到
3. 评估漏洞的严重程度和影响范围
4. 开发和测试修复
5. 在修复发布前协调披露时间表
6. 发布安全更新

### 致谢

我们会感谢所有负责任地报告安全问题的研究人员。

---

## 合规性

### GDPR 合规

- 数据保护
- 访问控制
- 审计日志
- 数据删除

### SOC 2 合规

- 安全监控
- 访问管理
- 变更管理
- 事件响应

### HIPAA 合规

- PHI 保护
- 访问审计
- 加密传输
- 业务连续性

---

## 相关文档

- [完整工具参考](/docs/tools/)
- [配置指南](/docs/configuration/)
- [部署指南](/docs/deployment/)
---
title: "认证配置"
weight: 30
---

# 认证配置

本文档描述 Cloud Native MCP Server 的认证和授权配置。

## 认证模式

服务器支持三种认证模式：

| 模式 | 描述 | 使用场景 |
|------|------|----------|
| `apikey` | API Key 认证 | 简单场景，快速部署 |
| `bearer` | Bearer Token 认证 | JWT token，更安全 |
| `basic` | HTTP Basic Auth | 传统认证方式 |

## API Key 认证

### 配置

```yaml
auth:
  # 启用/禁用认证
  enabled: true

  # 认证模式
  mode: "apikey"

  # API Key
  apiKey: "Abc123!@#Xyz789!@#Abc123!@#"
```

### 使用方式

```bash
# 使用 curl
curl -H "X-API-Key: Abc123!@#Xyz789!@#" \
  http://localhost:8080/api/kubernetes/http
```

### API Key 复杂度要求

- **最小长度**: 16 个字符
- **字符类型**: 以下 4 种类型中至少包含 3 种：
  - 大写字母 (A-Z)
  - 小写字母 (a-z)
  - 数字 (0-9)
  - 特殊字符 (!@#$%^&*()_+-=[]{}|;:,.<>?)

### 有效示例

- `Abc123!@#Xyz789!@#` (大写、小写、数字、特殊字符)
- `Abc123Xyz789Abc123` (大写、小写、数字)
- `ABC123!@#XYZ789!@#` (大写、数字、特殊字符)

## Bearer Token 认证

### 配置

```yaml
auth:
  enabled: true
  mode: "bearer"
  bearerToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

  # JWT 密钥（用于 JWT 验证）
  jwtSecret: "your-jwt-secret-key"

  # JWT 算法 (HS256, RS256, etc.)
  jwtAlgorithm: "HS256"
```

### 使用方式

```bash
# 使用 curl
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  http://localhost:8080/api/kubernetes/http
```

### Bearer Token 要求

- **格式**: `header.payload.signature`
- **最小长度**: 32 个字符
- **编码**: Base64URL 编码的各部分
- **验证**: 每个部分必须只包含有效的 base64url 字符

### JWT Token 结构

```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "sub": "user123",
    "name": "John Doe",
    "iat": 1516239022
  },
  "signature": "SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}
```

## Basic 认证

### 配置

```yaml
auth:
  enabled: true
  mode: "basic"
  username: "admin"
  password: "secure-password"
```

### 使用方式

```bash
# 使用 curl
curl -u admin:secure-password \
  http://localhost:8080/api/kubernetes/http

# 或者使用 Authorization 头
curl -H "Authorization: Basic YWRtaW46c2VjdXJlLXBhc3N3b3Jk" \
  http://localhost:8080/api/kubernetes/http
```

### Base64 编码

```bash
# 编码用户名和密码
echo -n "admin:secure-password" | base64
# 输出: YWRtaW46c2VjdXJlLXBhc3N3b3Jk
```

## 环境变量认证

使用环境变量配置认证：

```bash
# API Key 模式
export MCP_AUTH_ENABLED=true
export MCP_AUTH_MODE=apikey
export MCP_AUTH_API_KEY="Abc123!@#Xyz789!@#"

# Bearer Token 模式
export MCP_AUTH_ENABLED=true
export MCP_AUTH_MODE=bearer
export MCP_AUTH_API_KEY="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Basic Auth 模式
export MCP_AUTH_ENABLED=true
export MCP_AUTH_MODE=basic
export MCP_AUTH_USERNAME=admin
export MCP_AUTH_PASSWORD=secure-password
```

## 密钥管理

### 使用 Kubernetes Secrets

创建 Secret：

```bash
kubectl create secret generic mcp-auth \
  --from-literal=api-key='Abc123!@#Xyz789!@#' \
  --from-literal=jwt-secret='your-jwt-secret'
```

在部署中引用：

```yaml
env:
- name: MCP_AUTH_API_KEY
  valueFrom:
    secretKeyRef:
      name: mcp-auth
      key: api-key
- name: MCP_AUTH_JWT_SECRET
  valueFrom:
    secretKeyRef:
      name: mcp-auth
      key: jwt-secret
```

### 密钥轮换

定期轮换 API Key 和 Bearer Token：

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

## 速率限制

防止暴力破解和滥用：

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
    - "malicious.example.com"
```

## 安全最佳实践

### 1. 使用强认证凭证

- API Key 最少 16 字符
- Bearer Token 使用 JWT 标准
- 密码包含大小写字母、数字和特殊字符

### 2. 定期轮换密钥

```yaml
secrets:
  rotation_interval: 168  # 7 天
  max_age: 30
```

### 3. 使用环境变量

```yaml
auth:
  apiKey: "${MCP_AUTH_API_KEY}"
```

### 4. 启用速率限制

```yaml
ratelimit:
  enabled: true
  requests_per_second: 100
```

### 5. 监控认证失败

```yaml
audit:
  enabled: true
  level: "info"
```

### 6. 使用 HTTPS

```yaml
server:
  tls:
    certFile: "/path/to/cert.pem"
    keyFile: "/path/to/key.pem"
```

## 错误处理

### 认证失败

```json
{
  "error": {
    "code": "AUTHENTICATION_FAILED",
    "message": "Invalid API key or token"
  }
}
```

### 速率限制

```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many requests, please try again later",
    "data": {
      "retry_after": "60s"
    }
  }
}
```

## 相关文档

- [服务器配置](/zh/guides/configuration/server/)
- [服务配置](/zh/guides/configuration/services/)
- [安全指南](/zh/guides/security/)
- [密钥管理](/zh/guides/security/secrets/)
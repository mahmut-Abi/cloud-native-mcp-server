---
title: "认证"
weight: 10
---

# 认证

本文档描述 Cloud Native MCP Server 的认证配置和使用。

## 认证模式

服务器支持三种认证模式：

| 模式 | 描述 | 使用场景 |
|------|------|----------|
| `apikey` | API Key 认证 | 简单场景，快速部署 |
| `bearer` | Bearer Token 认证 | JWT token，更安全 |
| `basic` | HTTP Basic Auth | 传统认证方式 |

## API Key 认证

### API Key 要求

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

### 配置

```yaml
auth:
  enabled: true
  mode: "apikey"
  apiKey: "Abc123!@#Xyz789!@#Abc123!@#"
```

### 使用方式

```bash
# 使用 curl
curl -H "X-API-Key: Abc123!@#Xyz789!@#" \
  http://localhost:8080/api/kubernetes/http

# 或使用 Authorization 头
curl -H "Authorization: Bearer Abc123!@#Xyz789!@#" \
  http://localhost:8080/api/kubernetes/http
```

## Bearer Token 认证

### Bearer Token 要求

- **格式**: `header.payload.signature`
- **最小长度**: 32 个字符
- **编码**: Base64URL 编码的各部分
- **验证**: 每个部分必须只包含有效的 base64url 字符

### 有效示例

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```

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
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  http://localhost:8080/api/kubernetes/http
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
# 使用 curl -u
curl -u admin:secure-password \
  http://localhost:8080/api/kubernetes/http

# 或使用 Authorization 头
curl -H "Authorization: Basic YWRtaW46c2VjdXJlLXBhc3N3b3Jk" \
  http://localhost:8080/api/kubernetes/http
```

### Base64 编码

```bash
# 编码用户名和密码
echo -n "admin:secure-password" | base64
# 输出: YWRtaW46c2VjdXJlLXBhc3N3b3Jk
```

## 环境变量

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

## 使用 Kubernetes Secrets

### 创建 Secret

```bash
kubectl create secret generic mcp-auth \
  --from-literal=api-key='Abc123!@#Xyz789!@#' \
  --from-literal=username='admin' \
  --from-literal=password='secure-password'
```

### 在部署中引用

```yaml
env:
- name: MCP_AUTH_API_KEY
  valueFrom:
    secretKeyRef:
      name: mcp-auth
      key: api-key
- name: MCP_AUTH_USERNAME
  valueFrom:
    secretKeyRef:
      name: mcp-auth
      key: username
- name: MCP_AUTH_PASSWORD
  valueFrom:
    secretKeyRef:
      name: mcp-auth
      key: password
```

## 认证流程

### API Key 认证流程

```
1. 客户端发送请求
   - Header: X-API-Key: <api-key>

2. 服务器验证 API Key
   - 检查 API Key 是否存在
   - 验证 API Key 是否有效

3. 验证通过
   - 处理请求
   - 返回响应

4. 验证失败
   - 返回 401 Unauthorized
   - 记录审计日志
```

### Bearer Token 认证流程

```
1. 客户端发送请求
   - Header: Authorization: Bearer <token>

2. 服务器验证 Token
   - 解析 JWT token
   - 验证签名
   - 检查过期时间

3. 验证通过
   - 处理请求
   - 返回响应

4. 验证失败
   - 返回 401 Unauthorized
   - 记录审计日志
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

### 缺少认证

```json
{
  "error": {
    "code": "MISSING_AUTHENTICATION",
    "message": "API key or token is required"
  }
}
```

## 速率限制

防止暴力破解和滥用：

```yaml
ratelimit:
  enabled: true
  requests_per_second: 100
  burst: 200
```

## 相关文档

- [密钥管理](/zh/guides/security/secrets/)
- [最佳实践](/zh/guides/security/best-practices/)
- [配置指南](/zh/guides/configuration/authentication/)
- [安全指南](/zh/guides/security/)
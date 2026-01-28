---
title: "密钥管理"
weight: 20
---

# 密钥管理

本文档描述 Cloud Native MCP Server 的密钥管理功能。

## 概述

服务器包含密钥管理模块，用于安全地存储和管理敏感凭据。

## 特性

- **安全存储**: 带过期支持的内存存储
- **密钥轮换**: API 密钥和 bearer token 的自动轮换
- **密钥生成**: 内置生成器，用于复杂的 API 密钥和 JWT 类型的 token
- **环境变量**: 支持从环境变量加载密钥
- **密钥类型**: API 密钥、bearer token、basic auth 凭据

## 密钥类型

### API Key

用于简单认证场景：

```yaml
apiKey: "Abc123!@#Xyz789!@#Abc123!@#"
```

**要求**:
- 最小 16 字符
- 包含至少 3 种字符类型（大写、小写、数字、特殊字符）

### Bearer Token

用于 JWT 认证：

```yaml
bearerToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**要求**:
- 最小 32 字符
- 符合 JWT 标准（header.payload.signature）

### Basic Auth

用于传统认证：

```yaml
username: "admin"
password: "secure-password"
```

**要求**:
- 用户名和密码非空
- 密码足够复杂

## 使用密钥管理器

### 创建密钥管理器

```go
import "github.com/mahmut-Abi/cloud-native-mcp-server/internal/secrets"

// 创建新的密钥管理器
manager := secrets.NewInMemoryManager()
```

### 存储密钥

```go
secret := &secrets.Secret{
    Type:  secrets.SecretTypeAPIKey,
    Name:  "my-api-key",
    Value: "Abc123!@#Xyz789!@#",
}

err := manager.Store(secret)
```

### 检索密钥

```go
retrieved, err := manager.Retrieve(secret.ID)
if err != nil {
    // 处理错误
}

fmt.Printf("Secret value: %s\n", retrieved.Value)
```

### 轮换密钥

```go
rotated, err := manager.Rotate(secret.ID)
if err != nil {
    // 处理错误
}

fmt.Printf("New secret value: %s\n", rotated.Value)
```

### 生成新的 API 密钥

```go
newKey, err := manager.GenerateAPIKey("my-new-key")
if err != nil {
    // 处理错误
}

fmt.Printf("Generated API key: %s\n", newKey.Value)
```

### 生成新的 Bearer Token

```go
newToken, err := manager.GenerateBearerToken("my-new-token")
if err != nil {
    // 处理错误
}

fmt.Printf("Generated bearer token: %s\n", newToken.Value)
```

## 密钥过期

密钥可以有过期时间：

```go
secret := &secrets.Secret{
    Type:       secrets.SecretTypeAPIKey,
    Name:       "temporary-key",
    Value:      "Abc123!@#Xyz789!@#",
    ExpiresAt:  time.Now().Add(24 * time.Hour), // 24 小时后过期
}

err := manager.Store(secret)
```

过期的密钥会自动从列表中排除，无法检索。

## 密钥轮换策略

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

### 轮换流程

```
1. 生成新密钥
   - manager.GenerateAPIKey()

2. 验证新密钥
   - 检查复杂度要求
   - 验证唯一性

3. 更新配置
   - 替换旧密钥
   - 保存到安全存储

4. 测试新密钥
   - 验证新密钥有效
   - 确认服务正常

5. 删除旧密钥
   - 保留过期密钥用于审计
   - 或完全删除
```

## 环境变量支持

### 从环境变量加载

```yaml
auth:
  apiKey: "${MCP_AUTH_API_KEY}"
  bearerToken: "${MCP_AUTH_BEARER_TOKEN}"
  username: "${MCP_AUTH_USERNAME}"
  password: "${MCP_AUTH_PASSWORD}"
```

### 设置环境变量

```bash
# API Key
export MCP_AUTH_API_KEY="Abc123!@#Xyz789!@#"

# Bearer Token
export MCP_AUTH_BEARER_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Basic Auth
export MCP_AUTH_USERNAME="admin"
export MCP_AUTH_PASSWORD="secure-password"
```

## Kubernetes Secrets

### 创建 Secret

```bash
kubectl create secret generic mcp-auth \
  --from-literal=api-key='Abc123!@#Xyz789!@#' \
  --from-literal=bearer-token='eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...' \
  --from-literal=username='admin' \
  --from-literal=password='secure-password'
```

### 在部署中使用

```yaml
env:
- name: MCP_AUTH_API_KEY
  valueFrom:
    secretKeyRef:
      name: mcp-auth
      key: api-key
- name: MCP_AUTH_BEARER_TOKEN
  valueFrom:
    secretKeyRef:
      name: mcp-auth
      key: bearer-token
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

## 安全最佳实践

### 1. 使用强密钥

- API Key 最少 16 字符
- Bearer Token 使用 JWT 标准
- 密码包含大小写字母、数字和特殊字符

### 2. 定期轮换

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

### 4. 永远不要提交到版本控制

```bash
# .gitignore
config.yaml
.env
*.key
*.pem
```

### 5. 限制访问权限

```bash
# 设置文件权限
chmod 600 config.yaml
chmod 600 .env
```

### 6. 使用加密存储

对于生产环境，考虑使用：
- HashiCorp Vault
- AWS Secrets Manager
- Azure Key Vault
- Google Secret Manager

## 相关文档

- [认证](/zh/guides/security/authentication/)
- [最佳实践](/zh/guides/security/best-practices/)
- [配置指南](/zh/guides/configuration/authentication/)
- [安全指南](/zh/guides/security/)
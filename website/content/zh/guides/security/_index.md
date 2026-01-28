---
title: "安全"
weight: 30
---

# 安全指南

本文档描述 Cloud Native MCP Server 的安全特性和最佳实践。

## 概述

Cloud Native MCP Server 提供多层安全保护，包括：

- **认证**: API Key、Bearer Token、Basic Auth
- **密钥管理**: 安全存储、密钥轮换、密钥生成
- **输入清理**: 防止注入攻击
- **审计日志**: 跟踪所有操作
- **安全头部**: 过滤敏感信息

## 内容

- [认证](/zh/guides/security/authentication/) - 认证方式和配置
- [密钥管理](/zh/guides/security/secrets/) - 密钥存储和管理
- [最佳实践](/zh/guides/security/best-practices/) - 安全最佳实践
- [输入清理](#输入清理) - 输入验证和清理
- [审计日志](#审计日志) - 审计日志配置
- [安全头部](#安全头部) - 安全头部配置

## 快速开始

### 启用认证

```yaml
auth:
  enabled: true
  mode: "apikey"
  apiKey: "Abc123!@#Xyz789!@#Abc123!@#"
```

### 启用审计日志

```yaml
audit:
  enabled: true
  storage: "database"
  database:
    type: "sqlite"
    sqlitePath: "/var/lib/cloud-native-mcp-server/audit.db"
  format: "json"
  masking:
    enabled: true
    maskValue: "***REDACTED***"
```

### 启用输入清理

```yaml
sanitization:
  enabled: true
  max_length: 1000
  allowed_protocols:
    - http
    - https
```

## 输入清理

所有用户输入都经过清理，以防止注入攻击。

### 清理特性

- **过滤值**: 移除危险字符（SQL 注入、XSS、命令注入）
- **URL 验证**: 只允许 http/https 协议
- **长度限制**: 最大字符串长度强制（1000 个字符）
- **特殊字符移除**: 移除分号、引号和其他注入向量

### 清理规则

以下字符会从用户输入中移除：

- **SQL 注入**: `;`, `'`, `"`, `--`, `/*`, `*/`
- **命令注入**: `|`, `&`, `$`, `(`, `)`, `<`, `>`, `\``, `\`
- **XSS**: `<script>`, `javascript:`, `onload=`, `onerror=`

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
```

### 审计事件

以下事件会被记录：

- 认证成功/失败
- 工具调用
- 配置更改
- 错误和异常
- 访问拒绝

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
  headers:
    X-Frame-Options: "DENY"
    X-Content-Type-Options: "nosniff"
    X-XSS-Protection: "1; mode=block"
    Strict-Transport-Security: "max-age=31536000; includeSubDomains"
    Content-Security-Policy: "default-src 'self'"
```

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

## 速率限制

防止暴力攻击和滥用：

```yaml
ratelimit:
  enabled: true
  requests_per_second: 100
  burst: 200
  cleanup_interval: 60
```

## 相关文档

- [认证](/zh/guides/security/authentication/)
- [密钥管理](/zh/guides/security/secrets/)
- [最佳实践](/zh/guides/security/best-practices/)
- [配置指南](/zh/guides/configuration/authentication/)
- [部署指南](/zh/guides/deployment/)
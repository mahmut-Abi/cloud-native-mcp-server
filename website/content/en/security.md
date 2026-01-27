---
title: "Security Guide"
---

# Security Guide

This document describes the security features and best practices for Cloud Native MCP Server.

## Table of Contents

- [Authentication](#authentication)
- [Secret Management](#secret-management)
- [Input Sanitization](#input-sanitization)
- [Audit Logging](#audit-logging)
- [Security Best Practices](#security-best-practices)
- [Security Headers](#security-headers)
- [Reporting Security Issues](#reporting-security-issues)

---

## Authentication

### API Key Authentication

API keys must meet the following complexity requirements:

- **Minimum Length**: 16 characters
- **Character Types**: Must include at least 3 of the following 4 types:
  - Uppercase letters (A-Z)
  - Lowercase letters (a-z)
  - Numbers (0-9)
  - Special characters (!@#$%^&*()_+-=[]{}|;:,.<>?)

**Valid Examples**:
- `Abc123!@#Xyz789!@#` (uppercase, lowercase, numbers, special characters)
- `Abc123Xyz789Abc123` (uppercase, lowercase, numbers)
- `ABC123!@#XYZ789!@#` (uppercase, numbers, special characters)

**Invalid Examples**:
- `Abc123!@#` (less than 16 characters)
- `abcdefgh12345678` (only lowercase and numbers, doesn't meet 3 character types)
- `ABCDEFGHIJKLMNOPQRSTUVWXYZ` (only uppercase)

### Bearer Token Authentication

Bearer tokens must follow JWT structure:

- **Format**: `header.payload.signature`
- **Minimum Length**: 32 characters
- **Encoding**: Base64URL encoded parts
- **Validation**: Each part must contain only valid base64url characters (A-Z, a-z, 0-9, -, _, +)

**Valid Example**:
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```

**Invalid Examples**:
- `abcdefgh12345678abcdefgh12345678` (no JWT structure)
- `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ` (less than 32 characters)
- `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c$` (invalid character at end)

### Basic Authentication

Basic authentication uses username and password:

- **Username**: Non-empty string
- **Password**: Non-empty string

**Example**:
```bash
curl -u admin:secret http://localhost:8080/api/aggregate/sse
```

### Configuring Authentication

Enable authentication in the configuration file:

```yaml
auth:
  # Enable/disable authentication
  enabled: true

  # Authentication mode: apikey | bearer | basic
  mode: "apikey"

  # API Key (for apikey mode)
  # Minimum 16 characters, recommended 32+ characters
  apiKey: "Abc123!@#Xyz789!@#Abc123!@#"

  # Bearer token (for bearer mode)
  bearerToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

  # Basic Auth username
  username: "admin"

  # Basic Auth password
  password: "secret-password"

  # JWT secret (for JWT verification)
  jwtSecret: "your-jwt-secret-key"

  # JWT algorithm (HS256, RS256, etc.)
  jwtAlgorithm: "HS256"
```

### Authentication via Environment Variables

Configure authentication using environment variables:

```bash
export MCP_AUTH_ENABLED=true
export MCP_AUTH_MODE=apikey
export MCP_AUTH_API_KEY="Abc123!@#Xyz789!@#Abc123!@#"

./cloud-native-mcp-server
```

---

## Secret Management

The server includes a secret management module for securely storing credentials.

### Features

- **Secure Storage**: In-memory storage with expiration support
- **Key Rotation**: Automatic rotation for API keys and bearer tokens
- **Key Generation**: Built-in generators for complex API keys and JWT-style tokens
- **Environment Variables**: Support for loading secrets from environment variables
- **Secret Types**: API keys, bearer tokens, basic auth credentials

### Using the Secret Manager

```go
import "github.com/mahmut-Abi/cloud-native-mcp-server/internal/secrets"

// Create a new secret manager
manager := secrets.NewInMemoryManager()

// Store a secret
secret := &secrets.Secret{
    Type:  secrets.SecretTypeAPIKey,
    Name:  "my-api-key",
    Value: "Abc123!@#Xyz789!@#",
}
manager.Store(secret)

// Retrieve a secret
retrieved, err := manager.Retrieve(secret.ID)

// Rotate a secret
rotated, err := manager.Rotate(secret.ID)

// Generate a new API key
newKey, err := manager.GenerateAPIKey("my-new-key")

// Generate a new bearer token
newToken, err := manager.GenerateBearerToken("my-new-token")
```

### Secret Expiration

Secrets can have expiration times:

```go
secret := &secrets.Secret{
    Type:       secrets.SecretTypeAPIKey,
    Name:       "temporary-key",
    Value:      "Abc123!@#Xyz789!@#",
    ExpiresAt:  time.Now().Add(24 * time.Hour), // Expires in 24 hours
}
```

Expired secrets are automatically excluded from the list and cannot be retrieved.

### Key Rotation Strategy

Regularly rotating keys is a security best practice:

```yaml
auth:
  enabled: true
  mode: "apikey"
  apiKey: "${MCP_AUTH_API_KEY}"

secrets:
  # Automatic rotation interval (hours)
  rotation_interval: 168  # 7 days

  # Secret expiration time (days)
  max_age: 30

  # Keep expired secrets (for auditing)
  keep_expired: true
```

---

## Input Sanitization

All user input is sanitized to prevent injection attacks.

### Sanitization Features

- **Filter Values**: Removes dangerous characters (SQL injection, XSS, command injection)
- **URL Validation**: Only allows http/https protocols for web fetching
- **Length Limits**: Maximum string length enforced (1000 characters)
- **Special Character Removal**: Removes semicolons, quotes, and other injection vectors

### Sanitization Rules

The following characters are removed from user input:

- **SQL Injection**: `;`, `'`, `"`, `--`, `/*`, `*/`
- **Command Injection**: `|`, `&`, `$`, `(`, `)`, `<`, `>`, `\``, `\`
- **XSS**: `<script>`, `javascript:`, `onload=`, `onerror=`

### Examples

```go
import "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/sanitize"

// Sanitize filter value
cleanValue := sanitize.SanitizeFilterValue("'; DROP TABLE users; --")
// Result: " DROP TABLE users "

// Sanitize JSONPath
cleanPath := sanitize.SanitizeJSONPath("$.data[*].name; rm -rf /")
// Result: "$.data[*].name rm -rf "

// Validate string
isValid := sanitize.ValidateString("normal input")
// Result: true
```

### Configuring Input Sanitization

```yaml
sanitization:
  # Enable input sanitization
  enabled: true

  # Maximum string length
  max_length: 1000

  # Allowed URL protocols
  allowed_protocols:
    - http
    - https

  # Blocked character patterns
  blocked_patterns:
    - "';"
    - "DROP TABLE"
    - "rm -rf"
    - "<script>"
```

---

## Audit Logging

Audit logs track all operations for security monitoring and compliance.

### Enabling Audit Logging

```yaml
audit:
  enabled: true
  level: "info"
  storage: "database"
  format: "json"

  # Sensitive data masking
  masking:
    enabled: true
    fields:
      - password
      - token
      - apiKey
      - secret
      - authorization
    maskValue: "***REDACTED***"

  # Storage configuration
  database:
    type: "sqlite"
    sqlitePath: "/var/lib/cloud-native-mcp-server/audit.db"
    maxRecords: 100000
```

### Audit Events

The following events are logged:

- Authentication success/failure
- Tool calls
- Configuration changes
- Errors and exceptions
- Access denials

### Audit Log Format

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

### Querying Audit Logs

```bash
# Query recent 100 audit logs
curl -H "X-API-Key: your-key" \
  "http://localhost:8080/api/audit/query?limit=100"

# Query audit logs for a specific user
curl -H "X-API-Key: your-key" \
  "http://localhost:8080/api/audit/query?user=admin&limit=50"

# Query failed authentication attempts
curl -H "X-API-Key: your-key" \
  "http://localhost:8080/api/audit/query?tool=auth_login&status=failed"
```

---

## Security Best Practices

### 1. Use Strong Authentication

- Always use API keys that meet complexity requirements
- Rotate API keys regularly
- Use bearer tokens for JWT-based authentication
- Never commit credentials to version control

### 2. Enable Audit Logging

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

### 3. Use HTTPS in Production

Always use HTTPS when deploying in production:

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8443"
  tls:
    certFile: "/path/to/cert.pem"
    keyFile: "/path/to/key.pem"
```

### 4. Restrict Access

- Use firewall rules to limit access to the server
- Implement network policies in Kubernetes
- Use RBAC to control access to Kubernetes resources
- Implement rate limiting to prevent brute force attacks

```yaml
ratelimit:
  enabled: true
  requests_per_second: 100
  burst: 200
```

### 5. Monitor Suspicious Activity

- Enable metrics and monitoring
- Set up alerts for failed authentication attempts
- Regularly review audit logs
- Implement anomaly detection

```yaml
monitoring:
  # Failed authentication alert threshold
  auth_failure_threshold: 5
  auth_failure_window: 300  # 5 minutes

  # Anomaly behavior detection
  anomaly_detection:
    enabled: true
    sensitivity: "medium"
```

### 6. Keep Dependencies Updated

Regularly update dependencies to patch security vulnerabilities:

```bash
go get -u ./...
go mod tidy
```

### 7. Use Kubernetes Secrets

Never hardcode sensitive information in configuration files:

```yaml
# Bad practice
auth:
  apiKey: "Abc123!@#Xyz789!@#"

# Good practice
auth:
  apiKey: "${MCP_AUTH_API_KEY}"
```

Create a Kubernetes Secret:

```bash
kubectl create secret generic mcp-secrets \
  --from-literal=api-key='Abc123!@#Xyz789!@#' \
  --from-literal=jwt-secret='your-jwt-secret'
```

Reference in deployment:

```yaml
env:
- name: MCP_AUTH_API_KEY
  valueFrom:
    secretKeyRef:
      name: mcp-secrets
      key: api-key
```

### 8. Implement Least Privilege Principle

- Grant only necessary permissions
- Use RBAC to limit Kubernetes access
- Regularly review and update permissions
- Use service accounts for isolation

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

### 9. Network Isolation

- Use network policies to restrict pod-to-pod communication
- Isolate services in different namespaces
- Use ingress controllers for external access management
- Consider using service mesh for mTLS

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

### 10. Container Security

- Run containers as non-root user
- Use read-only filesystems
- Remove unnecessary privileges
- Scan images for vulnerabilities

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

## Security Headers

The server automatically filters sensitive headers in debug logs:

- `Authorization`
- `Cookie`
- `X-API-Key`
- `X-Api-Key`
- `x-api-key`

These headers are never logged in plaintext.

### Custom Security Headers

```yaml
security:
  # Additional security headers
  headers:
    X-Frame-Options: "DENY"
    X-Content-Type-Options: "nosniff"
    X-XSS-Protection: "1; mode=block"
    Strict-Transport-Security: "max-age=31536000; includeSubDomains"
    Content-Security-Policy: "default-src 'self'"

  # Header filtering
  header_filtering:
    enabled: true
    filtered_headers:
      - authorization
      - cookie
      - x-api-key
      - x-auth-token
```

---

## TLS/SSL Configuration

Use TLS/SSL for encrypted communication in production environments:

### Basic TLS Configuration

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

### mTLS Configuration

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

### Let's Encrypt Integration

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

## Rate Limiting

Prevent brute force attacks and abuse:

```yaml
ratelimit:
  enabled: true
  requests_per_second: 100
  burst: 200
  cleanup_interval: 60

  # Specific client limits
  client_limits:
    default:
      requests_per_second: 100
    authenticated:
      requests_per_second: 200

  # Whitelist
  whitelist:
    - "10.0.0.0/8"
    - "192.168.0.0/16"

  # Blacklist
  blacklist:
    - "malicious.example.com"
```

---

## Reporting Security Issues

If you discover a security vulnerability, please report it privately:

- **Email**: security@example.com
- **GitHub Security Advisories**: https://github.com/mahmut-Abi/cloud-native-mcp-server/security/advisories

Please do not create public issues for security vulnerabilities.

### Security Disclosure Process

1. Report vulnerability through private channels
2. We will acknowledge receipt within 48 hours
3. Assess severity and impact of the vulnerability
4. Develop and test the fix
5. Coordinate disclosure timeline before release
6. Release security update

### Acknowledgments

We will credit all researchers who responsibly report security issues.

---

## Compliance

### GDPR Compliance

- Data protection
- Access control
- Audit logging
- Data deletion

### SOC 2 Compliance

- Security monitoring
- Access management
- Change management
- Incident response

### HIPAA Compliance

- PHI protection
- Access auditing
- Encrypted transmission
- Business continuity

---

## Related Documentation

- [Complete Tools Reference](/docs/tools/)
- [Configuration Guide](/docs/configuration/)
- [Deployment Guide](/docs/deployment/)
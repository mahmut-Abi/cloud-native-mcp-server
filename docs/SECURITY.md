# Security

This document describes the security features and best practices for the k8s-mcp-server.

## Authentication

### API Key Authentication

API keys must meet the following complexity requirements:

- **Minimum length**: 16 characters
- **Character classes**: At least 3 of the following 4 types:
  - Uppercase letters (A-Z)
  - Lowercase letters (a-z)
  - Digits (0-9)
  - Special characters (!@#$%^&*()_+-=[]{}|;:,.<>?)

**Valid examples**:
- `Abc123!@#Xyz789!@#` (uppercase, lowercase, digits, special)
- `Abc123Xyz789Abc123` (uppercase, lowercase, digits)
- `ABC123!@#XYZ789!@#` (uppercase, digits, special)

**Invalid examples**:
- `Abc123!@#` (less than 16 characters)
- `abcdefgh12345678` (only lowercase and digits, not 3 character classes)
- `ABCDEFGHIJKLMNOPQRSTUVWXYZ` (only uppercase)

### Bearer Token Authentication

Bearer tokens must follow JWT structure:

- **Format**: `header.payload.signature`
- **Minimum length**: 32 characters
- **Encoding**: Base64URL encoded parts
- **Validation**: Each part must contain only valid base64url characters (A-Z, a-z, 0-9, -, _, +)

**Valid example**:
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```

**Invalid examples**:
- `abcdefgh12345678abcdefgh12345678` (no JWT structure)
- `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ` (less than 32 chars)
- `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c$` (invalid character at end)

### Basic Authentication

Basic authentication uses username and password:

- **Username**: Non-empty string
- **Password**: Non-empty string

**Example**:
```bash
curl -u admin:secret http://localhost:8080/api/aggregate/sse
```

## Secrets Management

The server includes a secrets management module for secure credential storage.

### Features

- **Secure Storage**: In-memory storage with expiration support
- **Secret Rotation**: Automatic rotation for API keys and bearer tokens
- **Secret Generation**: Built-in generators for complex API keys and JWT-like tokens
- **Environment Variables**: Support for loading secrets from environment variables
- **Secret Types**: API keys, bearer tokens, basic auth credentials

### Using Secrets Manager

```go
import "github.com/mahmut-Abi/k8s-mcp-server/internal/secrets"

// Create a new secrets manager
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

Expired secrets are automatically excluded from listings and cannot be retrieved.

## Input Sanitization

All user inputs are sanitized to prevent injection attacks.

### Sanitization Features

- **Filter Values**: Remove dangerous characters (SQL injection, XSS, command injection)
- **URL Validation**: Only allow http/https schemes for web fetch
- **Length Limits**: Maximum string length enforcement (1000 characters)
- **Special Character Removal**: Remove semicolons, quotes, and other injection vectors

### Sanitization Rules

The following characters are removed from user inputs:

- SQL injection: `;`, `'`, `"`, `--`, `/*`, `*/`
- Command injection: `|`, `&`, `$`, `(`, `)`, `<`, `>`, `\``, `\`
- XSS: `<script>`, `javascript:`, `onload=`, `onerror=`

### Example

```go
import "github.com/mahmut-Abi/k8s-mcp-server/internal/util/sanitize"

// Sanitize a filter value
cleanValue := sanitize.SanitizeFilterValue("'; DROP TABLE users; --")
// Result: " DROP TABLE users "

// Sanitize a JSONPath
cleanPath := sanitize.SanitizeJSONPath("$.data[*].name; rm -rf /")
// Result: "$.data[*].name rm -rf "

// Validate a string
isValid := sanitize.ValidateString("normal input")
// Result: true
```

## Alert Sorting

Alertmanager alerts can be sorted by multiple fields.

### Supported Sort Fields

| Field | Description | Values |
|-------|-------------|--------|
| `severity` | Alert severity | critical, warning, info |
| `severity_desc` | Alert severity (descending) | critical > warning > info |
| `startsAt` | Alert start time | RFC3339 timestamp |
| `startsAt_desc` | Alert start time (descending) | - |
| `endsAt` | Alert end time | RFC3339 timestamp |
| `endsAt_desc` | Alert end time (descending) | - |
| `fingerprint` | Alert fingerprint | string |
| `fingerprint_desc` | Alert fingerprint (descending) | - |

### Example Usage

```json
{
  "sortBy": "severity_desc"
}
```

This will return alerts sorted by severity in descending order (critical first, then warning, then info).

## Security Best Practices

### 1. Use Strong Authentication

- Always use API keys that meet the complexity requirements
- Rotate API keys regularly
- Use bearer tokens for JWT-based authentication
- Never commit credentials to version control

### 2. Enable Audit Logging

```yaml
audit:
  enabled: true
  maxLogs: 1000
  storage:
    type: "file"
    path: "/var/log/k8s-mcp-server/audit.log"
```

### 3. Use HTTPS in Production

Always use HTTPS when deploying the server in production:

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8443"
  tls:
    certFile: "/path/to/cert.pem"
    keyFile: "/path/to/key.pem"
```

### 4. Limit Access

- Use firewall rules to limit access to the server
- Implement network policies in Kubernetes
- Use RBAC to control access to Kubernetes resources

### 5. Monitor for Suspicious Activity

- Enable metrics and monitoring
- Set up alerts for failed authentication attempts
- Review audit logs regularly

### 6. Keep Dependencies Updated

Regularly update dependencies to patch security vulnerabilities:

```bash
go get -u ./...
go mod tidy
```

## Security Headers

The server automatically filters sensitive headers from debug logs:

- `Authorization`
- `Cookie`
- `X-API-Key`
- `X-Api-Key`
- `x-api-key`

These headers are never logged in plain text.

## Reporting Security Issues

If you discover a security vulnerability, please report it privately:

- Email: security@example.com
- GitHub Security Advisories: https://github.com/mahmut-Abi/k8s-mcp-server/security/advisories

Please do not open a public issue for security vulnerabilities.
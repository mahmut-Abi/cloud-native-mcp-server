# Configuration Guide

This guide covers all configuration options for the K8s MCP Server.

## Table of Contents

- [Configuration Methods](#configuration-methods)
- [Server Configuration](#server-configuration)
- [Service Configuration](#service-configuration)
- [Authentication](#authentication)
- [Logging](#logging)
- [Audit Logging](#audit-logging)
- [Caching](#caching)
- [Performance Tuning](#performance-tuning)
- [Example Configurations](#example-configurations)

---

## Configuration Methods

The K8s MCP Server supports three configuration methods (in order of precedence):

1. **Command line flags** - Highest priority
2. **Environment variables** - Medium priority
3. **YAML config file** - Lowest priority

### Example Configuration Flow

```bash
# Config file sets defaults
# Environment variables override config file
# Command line flags override everything

./k8s-mcp-server \
  --config=config.yaml \
  --log-level=debug
```

---

## Server Configuration

### Basic Settings

```yaml
server:
  mode: "sse"              # Server mode: sse, http, or stdio
  addr: "0.0.0.0:8080"     # Listen address
  read_timeout: 30         # Read timeout in seconds
  write_timeout: 30        # Write timeout in seconds
  max_connections: 1000    # Maximum concurrent connections
```

### Command Line Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--mode` | Server mode (sse, http, stdio) | sse |
| `--addr` | Listen address | 0.0.0.0:8080 |
| `--config` | Config file path | config.yaml |
| `--log-level` | Log level (debug, info, warn, error) | info |

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `MCP_MODE` | Server mode | sse |
| `MCP_ADDR` | Listen address | 0.0.0.0:8080 |
| `MCP_LOG_LEVEL` | Log level | info |

---

## Service Configuration

### Kubernetes

```yaml
kubernetes:
  enabled: true
  kubeconfig: ""           # Path to kubeconfig (empty for in-cluster)
  timeout_sec: 30          # API timeout in seconds
  qps: 50                  # Queries per second limit
  burst: 100               # Burst limit
  cache_ttl: 300           # Cache TTL in seconds
  insecure_skip_tls_verify: false
```

### Grafana

```yaml
grafana:
  enabled: false
  url: "http://grafana:3000"
  api_key: ""              # API key for authentication
  username: ""             # Basic auth username
  password: ""             # Basic auth password
  timeout_sec: 30
  cache_ttl: 180
```

### Prometheus

```yaml
prometheus:
  enabled: false
  address: "http://prometheus:9090"
  timeout_sec: 30
  cache_ttl: 60
  query_timeout: 120       # Query timeout in seconds
```

### Kibana

```yaml
kibana:
  enabled: false
  url: "http://kibana:5601"
  api_key: ""
  username: ""
  password: ""
  timeout_sec: 30
  cache_ttl: 180
```

### Elasticsearch

```yaml
elasticsearch:
  enabled: false
  url: "http://elasticsearch:9200"
  username: ""
  password: ""
  timeout_sec: 30
  cache_ttl: 180
```

### Alertmanager

```yaml
alertmanager:
  enabled: false
  url: "http://alertmanager:9093"
  timeout_sec: 30
  cache_ttl: 60
```

### Jaeger

```yaml
jaeger:
  enabled: false
  url: "http://jaeger:16686"
  timeout_sec: 30
  cache_ttl: 120
```

### Utilities

```yaml
utilities:
  enabled: true
```

---

## Authentication

### API Key Authentication

```yaml
auth:
  enabled: true
  mode: "apikey"
  api_key: "your-secret-api-key"
  header_name: "X-API-Key"  # Custom header name
```

### Bearer Token Authentication

```yaml
auth:
  enabled: true
  mode: "bearer"
  api_key: "your-bearer-token"
```

### Basic Authentication

```yaml
auth:
  enabled: true
  mode: "basic"
  username: "admin"
  password: "secret-password"
```

### Multiple API Keys

```yaml
auth:
  enabled: true
  mode: "apikey"
  api_keys:
    - key: "key-1"
      name: "Service A"
      permissions: ["read", "write"]
    - key: "key-2"
      name: "Service B"
      permissions: ["read"]
```

### Environment Variables for Auth

| Variable | Description |
|----------|-------------|
| `MCP_AUTH_ENABLED` | Enable authentication (true/false) |
| `MCP_AUTH_MODE` | Auth mode (apikey, bearer, basic) |
| `MCP_AUTH_API_KEY` | API key or bearer token |
| `MCP_AUTH_USERNAME` | Username for basic auth |
| `MCP_AUTH_PASSWORD` | Password for basic auth |

---

## Logging

### Log Levels

```yaml
logging:
  level: "info"            # debug, info, warn, error
  format: "json"           # json or text
  output: "stdout"         # stdout, stderr, or file path
```

### Structured Logging

```yaml
logging:
  level: "debug"
  format: "json"
  output: "stdout"
  fields:
    environment: "production"
    cluster: "prod-cluster"
```

### Log Rotation

```yaml
logging:
  output: "/var/log/k8s-mcp-server.log"
  max_size: 100           # MB
  max_age: 30             # days
  max_backups: 10         # number of backups
  compress: true
```

---

## Audit Logging

### Basic Audit Configuration

```yaml
audit:
  enabled: true
  max_logs: 1000          # Maximum logs in memory
  log_level: "info"       # Log level for audit events
```

### File Audit Storage

```yaml
audit:
  enabled: true
  storage: "file"
  file_path: "/var/log/k8s-mcp-audit.log"
  max_size: 100           # MB
  max_age: 30
  max_backups: 10
  compress: true
```

### Database Audit Storage

```yaml
audit:
  enabled: true
  storage: "database"
  database:
    driver: "postgres"
    dsn: "postgres://user:pass@localhost:5432/audit?sslmode=disable"
    table: "audit_logs"
```

### Audit Fields

The following fields are logged for each operation:

- Timestamp
- Request ID
- User/Client info
- Tool name
- Parameters (masked for sensitive data)
- Execution time
- Result status
- Error (if any)

---

## Caching

### Global Cache Configuration

```yaml
cache:
  enabled: true
  type: "lru"             # lru or segmented
  max_size: 1000          # Maximum cache entries
  default_ttl: 300        # Default TTL in seconds
```

### Service-Specific Cache TTL

```yaml
kubernetes:
  cache_ttl: 300          # 5 minutes

grafana:
  cache_ttl: 180          # 3 minutes

prometheus:
  cache_ttl: 60           # 1 minute
```

### Cache Statistics

Enable cache statistics for monitoring:

```yaml
cache:
  enabled: true
  collect_stats: true
  stats_interval: 60      # Stats collection interval (seconds)
```

---

## Performance Tuning

### Connection Pooling

```yaml
server:
  max_connections: 1000
  max_idle_conns: 100
  idle_timeout: 90        # seconds
```

### Request Limits

```yaml
server:
  max_request_size: 10485760  # 10MB
  max_header_size: 1048576    # 1MB
  read_timeout: 30
  write_timeout: 30
```

### Rate Limiting

```yaml
ratelimit:
  enabled: true
  requests_per_second: 100
  burst: 200
  cleanup_interval: 60
```

### Response Size Control

```yaml
performance:
  max_response_size: 5242880  # 5MB
  truncate_large_responses: true
  compression_enabled: true
  compression_level: 6
```

### JSON Encoding Pool

```yaml
performance:
  json_pool_size: 100
  json_buffer_size: 8192
```

---

## Example Configurations

### Minimal Configuration (Kubernetes Only)

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8080"

logging:
  level: "info"

kubernetes:
  enabled: true
  kubeconfig: ""
```

### Full Stack Monitoring

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8080"

logging:
  level: "info"
  format: "json"

kubernetes:
  enabled: true
  kubeconfig: ""

grafana:
  enabled: true
  url: "http://grafana:3000"
  api_key: "${GRAFANA_API_KEY}"

prometheus:
  enabled: true
  address: "http://prometheus:9090"

alertmanager:
  enabled: true
  url: "http://alertmanager:9093"

audit:
  enabled: true
  storage: "file"
  file_path: "/var/log/k8s-mcp-audit.log"
```

### Production with Auth and Caching

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8080"
  max_connections: 1000

logging:
  level: "info"
  format: "json"

kubernetes:
  enabled: true
  kubeconfig: ""
  timeout_sec: 30
  cache_ttl: 300

grafana:
  enabled: true
  url: "http://grafana:3000"
  api_key: "${GRAFANA_API_KEY}"
  cache_ttl: 180

auth:
  enabled: true
  mode: "apikey"
  api_key: "${MCP_API_KEY}"

audit:
  enabled: true
  storage: "database"
  database:
    driver: "postgres"
    dsn: "${AUDIT_DB_DSN}"

cache:
  enabled: true
  type: "lru"
  max_size: 1000
  default_ttl: 300

ratelimit:
  enabled: true
  requests_per_second: 100
  burst: 200
```

---

## Configuration Validation

The server validates configuration on startup. Common validation errors:

### Invalid Server Mode
```
Error: invalid server mode "invalid". Must be one of: sse, http, stdio
```

### Missing Required Fields
```
Error: missing required field "api_key" in auth configuration
```

### Invalid Service URL
```
Error: invalid service URL "grafana:3000". Must include scheme (http/https)
```

### Out of Range Values
```
Error: cache_ttl must be between 0 and 3600 seconds
```

---

## Environment Variable Substitution

You can use environment variables in the YAML config file:

```yaml
grafana:
  url: "${GRAFANA_URL}"
  api_key: "${GRAFANA_API_KEY}"

auth:
  api_key: "${MCP_API_KEY}"
```

Set environment variables before starting the server:

```bash
export GRAFANA_URL="http://grafana:3000"
export GRAFANA_API_KEY="your-api-key"
export MCP_API_KEY="your-mcp-key"

./k8s-mcp-server
```

---

## Testing Configuration

Test your configuration without starting the server:

```bash
./k8s-mcp-server --config=config.yaml --validate-config
```

This will:
- Parse the configuration file
- Validate all fields
- Check service connectivity
- Report any errors

---

## Hot Reload

Hot reload is not supported. Restart the server to apply configuration changes:

```bash
# Send SIGTERM for graceful shutdown
kill -TERM <pid>

# Server will finish in-flight requests and exit
# Then start with new configuration
./k8s-mcp-server --config=new-config.yaml
```

---

For more information, see:
- [Complete Tools Reference](TOOLS.md)
- [Deployment Guide](DEPLOYMENT.md)
- [Main README](../README.md)
---
title: "Configuration Guide"
---

# Configuration Guide

This guide covers all configuration options for Cloud Native MCP Server.

## Table of Contents

- [Configuration Methods](#configuration-methods)
- [Server Configuration](#server-configuration)
- [Service Configuration](#service-configuration)
- [Authentication Configuration](#authentication-configuration)
- [Logging Configuration](#logging-configuration)
- [Audit Logs](#audit-logs)
- [Cache Configuration](#cache-configuration)
- [Performance Tuning](#performance-tuning)
- [Example Configurations](#example-configurations)

---

## Configuration Methods

K8s MCP Server supports three configuration methods (in order of priority):

1. **Command Line Arguments** - Highest priority
2. **Environment Variables** - Medium priority
3. **YAML Configuration File** - Lowest priority

### Configuration Priority Example

```bash
# Configuration file sets default values
# Environment variables override configuration file
# Command line arguments override all settings

./cloud-native-mcp-server \
  --config=config.yaml \
  --log-level=debug
```

---

## Server Configuration

### Basic Settings

```yaml
server:
  # Run mode: sse | streamable-http | http | stdio
  # Recommended: stdio for development, streamable-http for production
  mode: "sse"

  # Server listen address
  addr: "0.0.0.0:8080"

  # HTTP read timeout (seconds)
  # 0 = no timeout (not recommended for production)
  # Recommended: 30-60 seconds
  readTimeoutSec: 30

  # HTTP write timeout (seconds)
  # Should be set to 0 for SSE connections to keep them alive
  writeTimeoutSec: 0

  # HTTP idle timeout (seconds)
  # Default: 60 seconds
  idleTimeoutSec: 60
```

### SSE Path Configuration

```yaml
server:
  ssePaths:
    # Kubernetes SSE endpoint
    kubernetes: "/api/kubernetes/sse"

    # Grafana SSE endpoint
    grafana: "/api/grafana/sse"

    # Prometheus SSE endpoint
    prometheus: "/api/prometheus/sse"

    # Kibana SSE endpoint
    kibana: "/api/kibana/sse"

    # Helm SSE endpoint
    helm: "/api/helm/sse"

    # Alertmanager SSE endpoint
    alertmanager: "/api/alertmanager/sse"

    # Elasticsearch SSE endpoint
    elasticsearch: "/api/elasticsearch/sse"

    # Utilities SSE endpoint
    utilities: "/api/utilities/sse"

    # Aggregated SSE endpoint for all services
    aggregate: "/api/aggregate/sse"
```

### Streamable-HTTP Path Configuration

```yaml
server:
  streamableHttpPaths:
    # Kubernetes Streamable-HTTP endpoint
    kubernetes: "/api/kubernetes/streamable-http"

    # Grafana Streamable-HTTP endpoint
    grafana: "/api/grafana/streamable-http"

    # Prometheus Streamable-HTTP endpoint
    prometheus: "/api/prometheus/streamable-http"

    # Kibana Streamable-HTTP endpoint
    kibana: "/api/kibana/streamable-http"

    # Helm Streamable-HTTP endpoint
    helm: "/api/helm/streamable-http"

    # Alertmanager Streamable-HTTP endpoint
    alertmanager: "/api/alertmanager/streamable-http"

    # Elasticsearch Streamable-HTTP endpoint
    elasticsearch: "/api/elasticsearch/streamable-http"

    # Utilities Streamable-HTTP endpoint
    utilities: "/api/utilities/streamable-http"

    # Aggregated Streamable-HTTP endpoint for all services
    aggregate: "/api/aggregate/streamable-http"
```

### Command Line Arguments

| Parameter | Description | Default |
|-----------|-------------|---------|
| `--mode` | Server mode (sse, streamable-http, http, stdio) | sse |
| `--addr` | Listen address | 0.0.0.0:8080 |
| `--config` | Configuration file path | config.yaml |
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
  # kubeconfig file path
  # If empty, uses default: $KUBECONFIG → ~/.kube/config → service account
  kubeconfig: ""

  # Timeout for a single API call (seconds)
  timeoutSec: 30

  # API client queries per second (QPS)
  qps: 100.0

  # API client burst rate
  burst: 200
```

### Prometheus

```yaml
prometheus:
  # Enable/disable Prometheus service
  enabled: false

  # Prometheus server address
  # Format: http://host:port or https://host:port
  address: "http://localhost:9090"

  # Request timeout (seconds)
  timeoutSec: 30

  # Basic auth username (optional)
  username: ""

  # Basic auth password (optional)
  password: ""

  # Bearer token authentication (optional, higher priority than Basic Auth)
  bearerToken: ""

  # Skip TLS certificate verification
  # Do not use in production!
  tlsSkipVerify: false

  # TLS client certificate file path (for mTLS authentication)
  tlsCertFile: ""

  # TLS client key file path
  tlsKeyFile: ""

  # TLS CA certificate file path
  tlsCAFile: ""
```

### Grafana

```yaml
grafana:
  # Enable/disable Grafana service
  enabled: false

  # Grafana server URL
  # Format: http://host:port or https://host:port
  url: "http://localhost:3000"

  # Grafana API Key (recommended)
  # Create in Grafana: Administration → API Keys
  apiKey: ""

  # Basic auth username (alternative to API Key)
  username: ""

  # Basic auth password
  password: ""

  # Request timeout (seconds)
  timeoutSec: 30
```

### Kibana

```yaml
kibana:
  # Enable/disable Kibana service
  enabled: false

  # Kibana server URL
  # Format: http://host:port or https://host:port
  url: "https://localhost:5601"

  # Kibana API Key (recommended)
  # Create in Kibana: Stack Management → API Keys
  apiKey: ""

  # Basic auth username (alternative to API Key)
  username: ""

  # Basic auth password
  password: ""

  # Request timeout (seconds)
  timeoutSec: 30

  # Skip TLS certificate verification
  # Do not use in production!
  skipVerify: false

  # Kibana space name
  # Default: "default"
  space: "default"
```

### Helm

```yaml
helm:
  # Enable/disable Helm service
  enabled: false

  # Helm operations kubeconfig path
  # If empty, uses the same kubeconfig as Kubernetes client
  kubeconfigPath: ""

  # Default namespace for Helm operations
  namespace: "default"

  # Enable Helm debug mode
  debug: false

  # Repository update timeout (seconds)
  # Default: 300 (5 minutes)
  # Recommended for China: 600-900
  timeoutSec: 300

  # Maximum retry attempts
  # Number of retries for failed repository updates
  # Default: 3
  # Recommended: 3-5
  maxRetries: 3

  # Enable mirrors
  # Used to accelerate Helm repository pulls
  # Default: false
  useMirrors: false

  # Custom mirror mapping
  # Format: original repository URL -> mirror URL
  mirrors: {}
```

### Elasticsearch

```yaml
elasticsearch:
  # Enable/disable Elasticsearch service
  enabled: false

  # Elasticsearch server addresses (supports multi-node HA)
  addresses:
    - "http://localhost:9200"

  # Single Elasticsearch server address (alternative to addresses)
  # Used when addresses is empty
  address: ""

  # Basic auth username
  username: ""

  # Basic auth password
  password: ""

  # Bearer token authentication (optional, higher priority than Basic Auth)
  bearerToken: ""

  # API Key authentication (optional, highest priority)
  # Format: id:api_key
  apiKey: ""

  # Request timeout (seconds)
  timeoutSec: 30

  # Skip TLS certificate verification
  # Do not use in production!
  tlsSkipVerify: false

  # TLS client certificate file path (for mTLS authentication)
  tlsCertFile: ""

  # TLS client key file path
  tlsKeyFile: ""

  # TLS CA certificate file path
  tlsCAFile: ""
```

### Alertmanager

```yaml
alertmanager:
  # Enable/disable Alertmanager service
  enabled: false

  # Alertmanager server address
  # Format: http://host:port or https://host:port
  address: "http://localhost:9093"

  # Request timeout (seconds)
  timeoutSec: 30

  # Basic auth username (optional)
  username: ""

  # Basic auth password (optional)
  password: ""

  # Bearer token authentication (optional, higher priority than Basic Auth)
  bearerToken: ""

  # Skip TLS certificate verification
  # Do not use in production!
  tlsSkipVerify: false

  # TLS client certificate file path (for mTLS authentication)
  tlsCertFile: ""

  # TLS client key file path
  tlsKeyFile: ""

  # TLS CA certificate file path
  tlsCAFile: ""
```

### OpenTelemetry

```yaml
opentelemetry:
  # Enable/disable OpenTelemetry service
  enabled: false

  # OpenTelemetry Collector address
  # Format: http://host:port or https://host:port
  address: "http://localhost:4318"

  # Request timeout (seconds)
  timeoutSec: 30

  # Basic auth username (optional)
  username: ""

  # Basic auth password (optional)
  password: ""

  # Bearer token authentication (optional, higher priority than Basic Auth)
  bearerToken: ""

  # Skip TLS certificate verification
  # Do not use in production!
  tlsSkipVerify: false

  # TLS client certificate file path (for mTLS authentication)
  tlsCertFile: ""

  # TLS client key file path
  tlsKeyFile: ""

  # TLS CA certificate file path
  tlsCAFile: ""
```

### Utilities

```yaml
utilities:
  # Utilities service is always enabled
  enabled: true
```

---

## Authentication Configuration

### API Key Authentication

```yaml
auth:
  # Enable/disable authentication
  enabled: false

  # Authentication mode: apikey | bearer | basic
  # apikey: X-API-Key simple API key authentication
  # bearer: Bearer Token (JWT) authentication
  # basic: HTTP Basic Auth
  mode: "apikey"

  # API Key (for apikey mode)
  # Minimum 8 characters, recommended 16+ characters
  apiKey: ""

  # Bearer token (for bearer mode)
  # Minimum 16 characters recommended (JWT token)
  bearerToken: ""

  # Basic Auth username
  username: ""

  # Basic Auth password
  password: ""

  # JWT secret (for JWT verification)
  jwtSecret: ""

  # JWT algorithm (HS256, RS256, etc.)
  jwtAlgorithm: "HS256"
```

### Authentication Environment Variables

| Variable | Description |
|----------|-------------|
| `MCP_AUTH_ENABLED` | Enable authentication (1, true, yes, on) |
| `MCP_AUTH_MODE` | Authentication mode (apikey, bearer, basic) |
| `MCP_AUTH_API_KEY` | API key or bearer token |
| `MCP_AUTH_USERNAME` | Basic auth username |
| `MCP_AUTH_PASSWORD` | Basic auth password |
| `MCP_AUTH_JWT_SECRET` | JWT secret |
| `MCP_AUTH_JWT_ALGORITHM` | JWT algorithm |

---

## Logging Configuration

```yaml
logging:
  # Log level: debug | info | warn | error
  level: "info"

  # Use JSON format logs
  # Suitable for log aggregation systems (ELK, Splunk, etc.)
  json: false
```

### Log Level Description

- **debug**: Detailed debugging information, including all requests and responses
- **info**: General information, including important operations and status changes
- **warn**: Warning information, does not affect functionality but needs attention
- **error**: Error information, functionality is impaired

---

## Audit Logs

### Basic Configuration

```yaml
audit:
  # Enable/disable audit logging
  enabled: false

  # Audit log level: debug | info | warn | error
  level: "info"

  # Audit log storage: stdout | file | database | all
  storage: "memory"

  # Log format: text | json
  # json: Structured JSON format, suitable for log aggregation
  # text: Human-readable text format
  format: "json"

  # Maximum query results
  maxResults: 1000

  # Query time range (days)
  timeRange: 90
```

### File Storage Configuration

```yaml
audit:
  storage: "file"
  file:
    # Log file path
    path: "/var/log/cloud-native-mcp-server/audit.log"

    # Maximum log file size (MB)
    maxSizeMB: 100

    # Maximum number of backup files
    maxBackups: 10

    # Maximum log file age (days)
    maxAgeDays: 30

    # Compress rotated log files
    compress: true

    # Maximum number of logs in memory storage
    maxLogs: 10000
```

### Database Storage Configuration

```yaml
audit:
  storage: "database"
  database:
    # Database type: sqlite | postgresql | mysql
    type: "sqlite"

    # SQLite database file path
    # Used only when type="sqlite"
    sqlitePath: "/var/lib/cloud-native-mcp-server/audit.db"

    # PostgreSQL connection string
    # Used only when type="postgresql"
    # Format: postgresql://user:password@host:port/dbname
    connectionString: ""

    # Database table name
    tableName: "audit_logs"

    # Maximum number of records
    maxRecords: 100000

    # Cleanup interval (hours)
    cleanupInterval: 24
```

### Query API Configuration

```yaml
audit:
  query:
    # Enable query API
    enabled: true

    # Maximum results per query
    maxResults: 1000

    # Maximum time range (days)
    timeRange: 90
```

### Sensitive Data Masking Configuration

```yaml
audit:
  masking:
    # Enable masking
    enabled: true

    # Fields to mask
    fields:
      - password
      - token
      - apiKey
      - secret
      - passwd
      - pwd
      - authorization

    # Mask replacement value
    maskValue: "***REDACTED***"
```

### Sampling Configuration (High Traffic Scenarios)

```yaml
audit:
  sampling:
    # Enable sampling
    enabled: false

    # Sampling rate (0-1)
    # 1.0 = log all, 0.1 = log 10%
    rate: 1.0
```

---

## Service and Tool Filtering

```yaml
enableDisable:
  # Disabled services (comma-separated)
  disabledServices: []

  # Enabled services (comma-separated, overrides disabled list)
  enabledServices: []

  # Disabled tools (comma-separated)
  disabledTools: []
```

---

## Performance Tuning

### Response Size Control

```yaml
# Implemented in code
performance:
  max_response_size: 5242880  # 5MB
  truncate_large_responses: true
  compression_enabled: true
  compression_level: 6
```

### JSON Encoding Pool

```yaml
# Implemented in code
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
  kubeconfig: ""
```

### Complete Monitoring Stack

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8080"

logging:
  level: "info"
  json: false

kubernetes:
  kubeconfig: ""

grafana:
  enabled: true
  url: "http://localhost:3000"
  apiKey: "${GRAFANA_API_KEY}"

prometheus:
  enabled: true
  address: "http://localhost:9090"

alertmanager:
  enabled: true
  address: "http://localhost:9093"

audit:
  enabled: true
  storage: "memory"
  format: "json"
```

### Production Configuration (Authentication and Caching)

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8080"
  readTimeoutSec: 30
  writeTimeoutSec: 0
  idleTimeoutSec: 60

logging:
  level: "info"
  json: true

kubernetes:
  kubeconfig: ""
  timeoutSec: 30
  qps: 100.0
  burst: 200

grafana:
  enabled: true
  url: "http://grafana:3000"
  apiKey: "${GRAFANA_API_KEY}"
  timeoutSec: 30

prometheus:
  enabled: true
  address: "http://prometheus:9090"
  timeoutSec: 30

auth:
  enabled: true
  mode: "apikey"
  apiKey: "${MCP_AUTH_API_KEY}"

audit:
  enabled: true
  storage: "database"
  database:
    type: "sqlite"
    sqlitePath: "/var/lib/cloud-native-mcp-server/audit.db"
    maxRecords: 100000
    cleanupInterval: 24
  format: "json"
  masking:
    enabled: true
    maskValue: "***REDACTED***"
```

---

## Configuration Validation

The server validates configuration on startup. Common validation errors:

### Invalid Server Mode
```
Error: invalid server mode "invalid". Must be one of: sse, streamable-http, http, stdio
```

### Missing Required Field
```
Error: missing required field "api_key" in auth configuration
```

### Invalid Service URL
```
Error: invalid service URL "grafana:3000". Must include scheme (http/https)
```

---

## Environment Variable Substitution

You can use environment variables in the YAML configuration file:

```yaml
grafana:
  url: "${GRAFANA_URL}"
  apiKey: "${GRAFANA_API_KEY}"

auth:
  apiKey: "${MCP_AUTH_API_KEY}"
```

Set environment variables before starting the server:

```bash
export GRAFANA_URL="http://grafana:3000"
export GRAFANA_API_KEY="your-api-key"
export MCP_AUTH_API_KEY="your-mcp-key"

./cloud-native-mcp-server
```

---

## Testing Configuration

Test configuration without starting the server:

```bash
# Check configuration file syntax
./cloud-native-mcp-server --config=config.yaml --validate-config
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

# Server will complete in-flight requests and exit
# Then start with new configuration
./cloud-native-mcp-server --config=new-config.yaml
```
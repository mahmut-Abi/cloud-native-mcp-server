# Configuration Guide

This document describes the **current** configuration schema used by Cloud Native MCP Server.

For website docs:
- English: `website/content/en/docs/configuration.md`
- Chinese: `website/content/zh/docs/configuration.md`

---

## Configuration Sources and Precedence

Configuration is merged in this order (highest first):

1. CLI flags (`--mode`, `--addr`, etc.)
2. Environment variables (`MCP_*`)
3. YAML config file (`--config`)
4. Built-in defaults

Common startup command:

```bash
./cloud-native-mcp-server --config=config.yaml
```

---

## Server Configuration

```yaml
server:
  # sse | http | streamable-http | stdio
  mode: "sse"

  # listen address
  addr: "0.0.0.0:8080"

  # HTTP timeouts in seconds
  readTimeoutSec: 30
  writeTimeoutSec: 0
  idleTimeoutSec: 60

  ssePaths:
    kubernetes: "/api/kubernetes/sse"
    grafana: "/api/grafana/sse"
    prometheus: "/api/prometheus/sse"
    kibana: "/api/kibana/sse"
    helm: "/api/helm/sse"
    elasticsearch: "/api/elasticsearch/sse"
    alertmanager: "/api/alertmanager/sse"
    jaeger: "/api/jaeger/sse"
    opentelemetry: "/api/opentelemetry/sse"
    utilities: "/api/utilities/sse"
    aggregate: "/api/aggregate/sse"

  streamableHttpPaths:
    kubernetes: "/api/kubernetes/streamable-http"
    grafana: "/api/grafana/streamable-http"
    prometheus: "/api/prometheus/streamable-http"
    kibana: "/api/kibana/streamable-http"
    helm: "/api/helm/streamable-http"
    elasticsearch: "/api/elasticsearch/streamable-http"
    alertmanager: "/api/alertmanager/streamable-http"
    jaeger: "/api/jaeger/streamable-http"
    opentelemetry: "/api/opentelemetry/streamable-http"
    utilities: "/api/utilities/streamable-http"
    aggregate: "/api/aggregate/streamable-http"

  cors:
    allowedOrigins: []
    allowedMethods: ["GET", "POST", "OPTIONS"]
    allowedHeaders: ["Content-Type", "Authorization", "X-API-Key"]
    maxAge: 86400
```

---

## Logging and Rate Limit

```yaml
logging:
  level: "info"
  json: false

ratelimit:
  enabled: false
  requests_per_second: 100
  burst: 200
```

Environment variables:
- `MCP_LOG_LEVEL`
- `MCP_LOG_JSON`
- `MCP_RATELIMIT_ENABLED`
- `MCP_RATELIMIT_REQUESTS_PER_SECOND`
- `MCP_RATELIMIT_BURST`

Note: legacy aliases `MCP_RATE_LIMIT_*` are still accepted.

---

## Authentication

```yaml
auth:
  enabled: true
  # apikey | bearer | basic
  mode: "apikey"

  # for apikey mode
  apiKey: "ChangeMe-Strong-Key-123!"

  # for bearer mode (static token mode)
  bearerToken: ""

  # for basic mode
  username: ""
  password: ""

  # optional JWT verification settings
  jwtSecret: ""
  jwtAlgorithm: "HS256"

  # OIDC discovery (bearer mode): provide issuer URL OR discovery URL
  oidcIssuerUrl: ""
  oidcDiscoveryUrl: ""

  # optional OIDC claim checks
  oidcIssuer: ""
  oidcAudience: ""
  oidcClientId: ""

  # optional OIDC transport/cache tuning
  oidcHttpTimeoutSec: 5
  oidcJwksCacheTtlSec: 600
```

Validation behavior:
- when `mode: apikey`, `auth.apiKey` is required
- when `mode: bearer`, either `auth.bearerToken` OR OIDC discovery config (`auth.oidcIssuerUrl` / `auth.oidcDiscoveryUrl`) is required
- when `mode: basic`, both `auth.username` and `auth.password` are required
- when OIDC discovery is configured, token signature/issuer/audience are validated from discovery + JWKS (OpenID Connect Discovery 1.0)

Environment variables (OIDC):
- `MCP_AUTH_OIDC_ISSUER_URL`
- `MCP_AUTH_OIDC_DISCOVERY_URL`
- `MCP_AUTH_OIDC_ISSUER`
- `MCP_AUTH_OIDC_AUDIENCE`
- `MCP_AUTH_OIDC_CLIENT_ID`
- `MCP_AUTH_OIDC_HTTP_TIMEOUT`
- `MCP_AUTH_OIDC_JWKS_CACHE_TTL`

---

## Service Configuration

```yaml
kubernetes:
  kubeconfig: ""
  timeoutSec: 30
  qps: 100.0
  burst: 200

prometheus:
  enabled: false
  address: "http://prometheus:9090"
  timeoutSec: 30
  username: ""
  password: ""
  bearerToken: ""
  tlsSkipVerify: false
  tlsCertFile: ""
  tlsKeyFile: ""
  tlsCAFile: ""

grafana:
  enabled: false
  url: "http://grafana:3000"
  apiKey: ""
  username: ""
  password: ""
  timeoutSec: 30

kibana:
  enabled: false
  url: "http://kibana:5601"
  apiKey: ""
  username: ""
  password: ""
  timeoutSec: 30
  skipVerify: false
  space: "default"

helm:
  enabled: false
  kubeconfigPath: ""
  namespace: "default"
  debug: false
  timeoutSec: 300
  maxRetries: 3
  httpProxy: ""

elasticsearch:
  enabled: false
  addresses: []
  address: ""
  username: ""
  password: ""
  bearerToken: ""
  apiKey: ""
  timeoutSec: 30
  tlsSkipVerify: false
  tlsCertFile: ""
  tlsKeyFile: ""
  tlsCAFile: ""

alertmanager:
  enabled: false
  address: "http://alertmanager:9093"
  timeoutSec: 30
  username: ""
  password: ""
  bearerToken: ""
  tlsSkipVerify: false
  tlsCertFile: ""
  tlsKeyFile: ""
  tlsCAFile: ""

jaeger:
  enabled: false
  address: "http://jaeger:16686"
  timeoutSec: 30

opentelemetry:
  enabled: false
  address: "http://otel-collector:4318"
  timeoutSec: 30
  username: ""
  password: ""
  bearerToken: ""
  tlsSkipVerify: false
  tlsCertFile: ""
  tlsKeyFile: ""
  tlsCAFile: ""
```

---

## Audit Configuration

```yaml
audit:
  enabled: false
  level: "info"
  storage: "memory" # memory | file | database | all
  format: "json"    # json | text
  maxResults: 1000
  timeRange: 90

  file:
    path: "/var/log/cloud-native-mcp-server/audit.log"
    maxSizeMB: 100
    maxBackups: 10
    maxAgeDays: 30
    compress: true
    maxLogs: 10000

  database:
    type: "sqlite" # sqlite | postgres | mysql
    sqlitePath: "/var/lib/cloud-native-mcp-server/audit.db"
    tableName: "audit_logs"
    maxRecords: 100000
    cleanupInterval: 24

  query:
    enabled: true
    maxResults: 1000
    timeRange: 90

  alerts:
    enabled: false
    failureThreshold: 10
    checkIntervalSec: 60
    method: "none"
    webhookURL: ""

  masking:
    enabled: true
    fields: ["password", "token", "apiKey", "authorization"]
    maskValue: "***REDACTED***"

  sampling:
    enabled: false
    rate: 1.0
```

---

## Service and Tool Filtering

```yaml
enableDisable:
  disabledServices: []
  enabledServices: []
  disabledTools: []
```

---

## Example Configurations

### Minimal

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8080"

logging:
  level: "info"

kubernetes:
  kubeconfig: ""
```

### Production Baseline

```yaml
server:
  mode: "streamable-http"
  addr: "0.0.0.0:8080"
  readTimeoutSec: 30
  writeTimeoutSec: 0
  idleTimeoutSec: 60

logging:
  level: "info"
  json: true

auth:
  enabled: true
  mode: "apikey"
  apiKey: "${MCP_AUTH_API_KEY}"

ratelimit:
  enabled: true
  requests_per_second: 100
  burst: 200

audit:
  enabled: true
  storage: "database"
  format: "json"
  database:
    type: "sqlite"
    sqlitePath: "/var/lib/cloud-native-mcp-server/audit.db"
    maxRecords: 100000
    cleanupInterval: 24

kubernetes:
  kubeconfig: ""
  timeoutSec: 30
  qps: 100.0
  burst: 200
```

---

## Validation and Smoke Checks

The current CLI does **not** provide a `--validate-config` flag.

Recommended quick checks:

```bash
# 1) Load config and print enabled services
./cloud-native-mcp-server --config=config.yaml --list=services --output=table

# 2) Start server and check health endpoint
./cloud-native-mcp-server --config=config.yaml
curl -sS http://127.0.0.1:8080/health
```

If configuration is invalid, startup logs include validation failure details.

---

## Related Docs

- `docs/DEPLOYMENT.md`
- `README.md`
- `website/content/en/docs/configuration.md`
- `website/content/zh/docs/configuration.md`

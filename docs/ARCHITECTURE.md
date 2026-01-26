# Architecture

This document describes the architecture and design of the k8s-mcp-server.

## Project Structure

```
k8s-mcp-server/
├── cmd/
│   └── server/              # Main entry point
│       ├── cli.go           # CLI flag parsing
│       ├── main.go          # Application entry
│       ├── server.go        # Server initialization
│       ├── display.go       # Display utilities
│       ├── shutdown.go      # Graceful shutdown
│       └── utils.go         # Utility functions
├── internal/
│   ├── config/              # Configuration management
│   │   ├── config.go        # Configuration structures
│   │   ├── loader.go        # Config file loading
│   │   ├── validator.go     # Config validation
│   │   ├── env_parser.go    # Environment variable parsing
│   │   └── serverConfig/    # Server-specific config
│   ├── constants/           # Application constants
│   │   └── constants.go     # Magic numbers and defaults
│   ├── errors/              # Error handling
│   │   ├── errors.go        # Error types
│   │   └── converter.go     # Error conversion
│   ├── logging/             # Logging utilities
│   │   ├── logging.go       # Logger initialization
│   │   └── logging_test.go  # Logging tests
│   ├── middleware/          # HTTP middleware
│   │   ├── auth_middleware.go      # Authentication
│   │   ├── auth.go                 # Auth utilities
│   │   ├── metrics_middleware.go   # Metrics collection
│   │   ├── ratelimit.go            # Rate limiting
│   │   ├── audit_middleware.go     # Audit logging
│   │   ├── audit_log.go            # Audit log storage
│   │   ├── security_middleware.go  # Security headers
│   │   └── hook/                   # Middleware hooks
│   ├── observability/       # Metrics and monitoring
│   │   └── metrics/         # Metrics definitions
│   ├── secrets/             # Secrets management
│   │   ├── manager.go       # Secrets manager
│   │   └── manager_test.go  # Secrets tests
│   ├── services/            # Service implementations
│   │   ├── kubernetes/      # Kubernetes service
│   │   ├── helm/            # Helm service
│   │   ├── grafana/         # Grafana service
│   │   ├── prometheus/      # Prometheus service
│   │   ├── kibana/          # Kibana service
│   │   ├── elasticsearch/   # Elasticsearch service
│   │   ├── alertmanager/    # Alertmanager service
│   │   ├── jaeger/          # Jaeger service
│   │   ├── utilities/       # Utilities service
│   │   ├── cache/           # LRU cache
│   │   ├── framework/       # Service framework
│   │   ├── manager/         # Service manager
│   │   ├── common/          # Common helpers
│   │   ├── tools/           # Tool builders
│   │   └── prompts/         # Prompt templates
│   └── util/                # Utilities
│       ├── circuitbreaker/  # Circuit breaker
│       ├── performance/     # Performance optimizations
│       ├── pool/            # Object pooling
│       ├── sanitize/        # Input sanitization
│       ├── retry/           # Retry logic
│       └── log/             # Logging utilities
├── docs/                    # Documentation
│   ├── TOOLS.md            # Tools reference
│   ├── CONFIGURATION.md    # Configuration guide
│   ├── DEPLOYMENT.md      # Deployment guide
│   ├── SECURITY.md        # Security guide
│   ├── ARCHITECTURE.md    # Architecture guide
│   └── PERFORMANCE.md     # Performance guide
└── deploy/                  # Deployment files
    ├── Dockerfile
    ├── helm/
    │   └── k8s-mcp-server/
    └── kubernetes/
```

## Core Components

### 1. Server Initialization

The server is initialized in `cmd/server/server.go`:

1. Parse CLI flags and configuration
2. Initialize logging
3. Initialize metrics system
4. Initialize audit storage
5. Initialize authentication provider
6. Initialize services
7. Create MCP servers for each mode (SSE, HTTP, stdio)
8. Setup HTTP routes
9. Start server

### 2. Service Framework

Each service follows a common pattern:

```
service/
├── client/          # Backend client
├── handlers/        # MCP tool handlers
├── service.go       # Service interface
└── tools/           # Tool definitions
```

### Service Initialization

```go
func InitService(config ServiceConfig) (*Service, error) {
    // 1. Create client
    client := NewClient(config)
    
    // 2. Create handlers
    handlers := NewHandlers(client)
    
    // 3. Register tools
    for name, handler := range handlers {
        RegisterTool(name, handler)
    }
    
    // 4. Return service
    return &Service{
        client:   client,
        handlers: handlers,
    }, nil
}
```

### 3. MCP Protocol

The server implements the Model Context Protocol (MCP):

- **SSE Mode**: Server-Sent Events for real-time communication
- **HTTP Mode**: RESTful HTTP API
- **Stdio Mode**: Standard input/output for CLI tools

### 4. Middleware Pipeline

Requests flow through middleware in this order:

1. **CORS**: Cross-Origin Resource Sharing
2. **Logging**: Request/response logging
3. **Security**: Security headers
4. **Rate Limiting**: Token bucket rate limiter
5. **Authentication**: API key, bearer token, basic auth
6. **Audit**: Audit logging
7. **Metrics**: Metrics collection

## Data Flow

### Request Flow

```
Client Request
    ↓
CORS Middleware
    ↓
Logging Middleware
    ↓
Security Middleware
    ↓
Rate Limiting Middleware
    ↓
Authentication Middleware
    ↓
Audit Middleware
    ↓
Metrics Middleware
    ↓
Service Handler
    ↓
Backend Client
    ↓
Response
    ↓
Audit Logging
    ↓
Metrics Recording
    ↓
Client Response
```

### Tool Execution Flow

```
1. Receive MCP tool request
2. Extract arguments
3. Sanitize inputs
4. Call service handler
5. Handler calls backend client
6. Client makes HTTP request to backend
7. Parse response
8. Apply caching (if enabled)
9. Return formatted response
```

## Caching Strategy

### LRU Cache with TTL

The server uses a segmented LRU cache:

- **Hot Segment**: Recently accessed items (20% of capacity)
- **Cold Segment**: Less frequently accessed items (80% of capacity)
- **TTL**: Time-to-live for cache entries
- **Eviction**: LRU eviction when full

### Cache Keys

Cache keys are constructed from:
- Service name
- Tool name
- Request parameters (sorted and hashed)

Example:
```
kubernetes:list_resources:namespace=default:resource=pods
```

## Authentication Flow

### API Key Authentication

```
1. Extract API key from header (X-API-Key) or query parameter
2. Validate API key format (length, complexity)
3. Look up API key in auth provider
4. Check if API key is valid and not expired
5. Extract user roles and permissions
6. Allow or deny request
```

### Bearer Token Authentication

```
1. Extract bearer token from Authorization header
2. Validate JWT structure (header.payload.signature)
3. Validate base64url encoding
4. Look up token in auth provider
5. Check if token is valid and not expired
6. Extract user roles and permissions
7. Allow or deny request
```

## Error Handling

### Error Types

- `ErrMissingRequiredParam`: Required parameter is missing
- `ErrInvalidParameter`: Parameter validation failed
- `ErrResourceNotFound`: Resource not found
- `ErrUnauthorized`: Authentication failed
- `ErrForbidden`: Permission denied
- `ErrRateLimitExceeded`: Rate limit exceeded

### Error Response Format

```json
{
  "code": 1,
  "data": null,
  "message": "Error description"
}
```

## Metrics

### Collected Metrics

- **HTTP Metrics**: Request count, response time, active connections
- **Service Metrics**: Tool calls, cache hits/misses, errors
- **Circuit Breaker Metrics**: State changes, failures, successes
- **Rate Limiter Metrics**: Requests allowed/denied

### Metrics Endpoint

```
GET /metrics
```

Returns Prometheus-formatted metrics.

## Concurrency Model

### Goroutine Management

- **Goroutine Pool**: Reusable goroutine pool for concurrent operations
- **Context Cancellation**: Proper cleanup on context cancellation
- **Rate Limiting**: Token bucket algorithm for rate limiting
- **Circuit Breaker**: Prevent cascading failures

### Thread Safety

- **Mutex Protection**: Shared state protected by mutexes
- **Atomic Operations**: Use atomic operations for counters
- **Channels**: Use channels for goroutine communication

## Extension Points

### Adding a New Service

1. Create service directory: `internal/services/myservice/`
2. Implement client: `client/client.go`
3. Implement handlers: `handlers/handlers.go`
4. Register service in `internal/services/registry.go`
5. Add configuration to `internal/config/config.go`
6. Update documentation

### Adding a New Tool

1. Define tool in `handlers/handlers.go`
2. Implement handler function
3. Register tool in service initialization
4. Add tool description to `docs/TOOLS.md`

## Performance Optimizations

### JSON Encoding Pool

Reusable JSON encoders to reduce allocations:

```go
var GlobalJSONPool = NewJSONEncoderPool(100)
```

### Object Pooling

Reusable objects to reduce GC pressure:

```go
var StringBuilderPool = sync.Pool{
    New: func() interface{} {
        return &strings.Builder{}
    },
}
```

### Response Compression

Automatic response compression for large responses:

```go
if len(response) > 1024 {
    response = compress(response)
}
```

## Security Considerations

### Input Validation

- All user inputs are sanitized
- URL validation for web fetch
- Length limits on all inputs
- Type checking for parameters

### Secrets Management

- Secrets stored in memory
- Automatic expiration
- Secure generation
- Environment variable support

### Audit Logging

- All tool calls are logged
- Sensitive headers are filtered
- Configurable storage backends

## Deployment Architecture

### Single Instance

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │
       ▼
┌─────────────────┐
│ k8s-mcp-server  │
│                 │
│  ┌───────────┐  │
│  │ Services  │  │
│  └─────┬─────┘  │
└────────┼────────┘
         │
    ┌────┴────┐
    │         │
    ▼         ▼
┌────────┐ ┌──────────┐
│ K8s    │ │ Backends │
│ API    │ │ (Grafana │
│        │ │  Prom   │
└────────┘ └──────────┘
```

### Kubernetes Deployment

```
┌─────────────────────────────────┐
│            Ingress              │
└────────────┬────────────────────┘
             │
             ▼
┌─────────────────────────────────┐
│         Service (Load Balancer)  │
└────────────┬────────────────────┘
             │
    ┌────────┴────────┐
    │                 │
    ▼                 ▼
┌─────────┐      ┌─────────┐
│ Pod 1   │ ...  │ Pod N   │
│         │      │         │
└─────────┘      └─────────┘
```

## Future Enhancements

- [ ] gRPC support
- [ ] WebSocket support
- [ ] Plugin system
- [ ] Custom middleware
- [ ] Distributed tracing
- [ ] Advanced caching strategies
- [ ] Service mesh integration
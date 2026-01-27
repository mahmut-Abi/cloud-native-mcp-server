# Performance

This document describes the performance features and optimizations of the cloud-native-mcp-server.

## Performance Features

- **Intelligent Caching**: LRU cache with TTL for frequently accessed data
- **Response Size Control**: Automatic truncation and optimization
- **JSON Encoding Pool**: Reuse JSON encoders for better performance
- **Circuit Breaker**: Prevent cascading failures
- **Pagination**: Support for large datasets
- **Summary Tools**: Optimized tools for LLM consumption
- **Input Sanitization**: Protection against injection attacks
- **Secrets Management**: Secure credential storage with rotation support
- **Enhanced Validation**: Strict API key and token validation

## Caching

### LRU Cache with TTL

The server uses a segmented LRU cache to optimize performance:

- **Capacity**: Configurable (default: 10,000 entries)
- **TTL**: Configurable (default: 5 minutes)
- **Hot Segment**: 20% of capacity for frequently accessed items
- **Cold Segment**: 80% of capacity for less frequently accessed items

### Cache Configuration

```yaml
cache:
  enabled: true
  maxSize: 10000
  ttl: 300s
  hotSegmentRatio: 0.2
```

### Cache Keys

Cache keys are constructed from:
- Service name
- Tool name
- Request parameters (sorted and hashed)

Example:
```
kubernetes:list_resources:namespace=default:resource=pods
```

### Cache Hit Rate

Typical cache hit rates:
- **Kubernetes resources**: 70-80%
- **Grafana dashboards**: 60-70%
- **Prometheus metrics**: 50-60%
- **Elasticsearch indices**: 40-50%

## Response Size Control

### Automatic Truncation

Large responses are automatically truncated to prevent context overflow:

```yaml
response:
  maxSize: 1000000  # 1MB
  truncate: true
  truncateMessage: "... (truncated)"
```

### Summary Tools

Many tools have LLM-optimized versions that return 70-95% smaller responses:

| Tool | Full Size | Summary Size | Reduction |
|------|-----------|--------------|-----------|
| `kubernetes_list_resources` | 50KB | 5KB | 90% |
| `grafana_dashboards` | 30KB | 3KB | 90% |
| `prometheus_get_alerts` | 40KB | 8KB | 80% |
| `elasticsearch_list_indices` | 60KB | 10KB | 83% |

### Pagination

Support for large datasets with pagination:

```json
{
  "limit": 50,
  "offset": 0
}
```

## JSON Encoding Pool

### Object Pooling

The server uses object pooling to reduce allocations:

```go
var GlobalJSONPool = NewJSONEncoderPool(100)
```

### Benefits

- **Reduced GC Pressure**: Fewer allocations mean less garbage collection
- **Better Performance**: Reused encoders are faster than creating new ones
- **Memory Efficiency**: Pool size is configurable

### StringBuilder Pool

Reusable string builders for string operations:

```go
var StringBuilderPool = sync.Pool{
    New: func() interface{} {
        return &strings.Builder{}
    },
}
```

## Circuit Breaker

### Circuit Breaker Pattern

Prevents cascading failures by stopping requests to failing services:

```yaml
circuitBreaker:
  enabled: true
  failureThreshold: 5
  successThreshold: 2
  timeout: 30s
```

### States

1. **Closed**: Normal operation, requests pass through
2. **Open**: Requests are blocked, circuit is tripped
3. **Half-Open**: Testing if service has recovered

### Metrics

- **Circuit Breaker State Changes**: Track state transitions
- **Failure Count**: Count of failures before tripping
- **Success Count**: Count of successes in half-open state

## Rate Limiting

### Token Bucket Algorithm

Rate limiting using token bucket algorithm:

```yaml
rateLimit:
  enabled: true
  requestsPerSecond: 10
  burstSize: 20
```

### Benefits

- **Protection**: Prevents abuse and DoS attacks
- **Fairness**: Ensures equal access for all clients
- **Configurable**: Can be adjusted per client or globally

## HTTP Connection Pooling

### Connection Pool Configuration

```go
transport := &http.Transport{
    MaxIdleConns:        256,
    MaxIdleConnsPerHost: 100,
    IdleConnTimeout:     120 * time.Second,
}
```

### Benefits

- **Reduced Latency**: Reused connections are faster
- **Lower Resource Usage**: Fewer TCP handshakes
- **Better Scalability**: Handles more concurrent requests

## Goroutine Management

### Goroutine Pool

Reusable goroutines for concurrent operations:

```go
pool := NewWorkerPool(100)
```

### Context Cancellation

Proper cleanup on context cancellation:

```go
select {
case <-ctx.Done():
    return ctx.Err()
case result := <-results:
    return result
}
```

### Benefits

- **Efficient**: Reusable goroutines reduce overhead
- **Scalable**: Handles high concurrency
- **Safe**: Proper cleanup prevents resource leaks

## Performance Metrics

### HTTP Metrics

- **Request Count**: Total number of requests
- **Response Time**: Time to process requests
- **Active Connections**: Number of active connections
- **Error Rate**: Percentage of failed requests

### Service Metrics

- **Tool Calls**: Number of tool invocations
- **Cache Hit Rate**: Percentage of cache hits
- **Backend Latency**: Time spent calling backends
- **Error Rate**: Percentage of tool failures

### Circuit Breaker Metrics

- **State Changes**: Number of state transitions
- **Failure Count**: Count of failures
- **Success Count**: Count of successes

### Rate Limiter Metrics

- **Requests Allowed**: Number of requests allowed
- **Requests Denied**: Number of requests denied
- **Current Rate**: Current request rate

## Performance Tuning

### Cache Tuning

Increase cache size for frequently accessed data:

```yaml
cache:
  maxSize: 20000  # Double the default
  ttl: 600s      # Increase TTL to 10 minutes
```

### Rate Limiting

Adjust rate limits based on your needs:

```yaml
rateLimit:
  requestsPerSecond: 50  # Increase for high-traffic scenarios
  burstSize: 100          # Increase burst size
```

### Connection Pooling

Adjust connection pool for high concurrency:

```yaml
http:
  maxIdleConns: 500
  maxIdleConnsPerHost: 200
```

## Performance Benchmarks

### Tool Execution Time

| Tool | Average Time | P95 Time | P99 Time |
|------|--------------|----------|----------|
| `kubernetes_list_resources` | 50ms | 100ms | 200ms |
| `grafana_dashboards` | 30ms | 60ms | 100ms |
| `prometheus_query` | 40ms | 80ms | 150ms |
| `elasticsearch_search` | 60ms | 120ms | 250ms |

### Throughput

- **Requests per Second**: 1000+ (with caching)
- **Concurrent Connections**: 500+
- **Memory Usage**: 100-200MB (typical)

### Latency

- **P50 Latency**: 30ms
- **P95 Latency**: 100ms
- **P99 Latency**: 200ms

## Performance Monitoring

### Metrics Endpoint

```
GET /metrics
```

Returns Prometheus-formatted metrics:

```
http_requests_total{method="GET",path="/api/kubernetes/sse"} 1234
http_request_duration_seconds{method="GET",path="/api/kubernetes/sse",quantile="0.95"} 0.1
cache_hits_total{service="kubernetes"} 5678
cache_misses_total{service="kubernetes"} 123
```

### Health Check

```
GET /health
```

Returns server health status:

```json
{
  "status": "healthy",
  "uptime": 3600,
  "version": "1.0.0"
}
```

## Performance Best Practices

### 1. Enable Caching

Always enable caching for frequently accessed data:

```yaml
cache:
  enabled: true
```

### 2. Use Summary Tools

Use summary tools to reduce response size:

```json
{
  "tool": "kubernetes_list_resources_summary"
}
```

### 3. Enable Circuit Breaker

Prevent cascading failures:

```yaml
circuitBreaker:
  enabled: true
```

### 4. Configure Rate Limiting

Protect against abuse:

```yaml
rateLimit:
  enabled: true
  requestsPerSecond: 10
```

### 5. Monitor Metrics

Regularly monitor performance metrics:

```bash
curl http://localhost:8080/metrics
```

### 6. Tune for Your Workload

Adjust configuration based on your specific workload:

- High traffic: Increase rate limits and connection pool
- Large datasets: Increase cache size and TTL
- Low latency: Enable caching and use summary tools

## Performance Troubleshooting

### High Latency

**Symptoms**: Slow response times

**Solutions**:
- Enable caching
- Increase cache size
- Use summary tools
- Check backend latency

### High Memory Usage

**Symptoms**: High memory consumption

**Solutions**:
- Reduce cache size
- Reduce TTL
- Enable response truncation
- Check for memory leaks

### High CPU Usage

**Symptoms**: High CPU utilization

**Solutions**:
- Reduce rate limit
- Enable circuit breaker
- Check for inefficient queries
- Profile the application

### Cache Miss Rate

**Symptoms**: Low cache hit rate

**Solutions**:
- Increase TTL
- Increase cache size
- Check cache key generation
- Monitor cache usage

## Future Performance Enhancements

- [ ] Redis caching for distributed deployments
- [ ] Response compression (gzip, brotli)
- [ ] HTTP/2 support
- [ ] gRPC support
- [ ] Query optimization
- [ ] Predictive caching
- [ ] Adaptive rate limiting
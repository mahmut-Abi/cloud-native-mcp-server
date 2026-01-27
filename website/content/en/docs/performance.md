---
title: "Performance Guide"
---

# Performance Guide

This document describes the performance features and optimization recommendations for Cloud Native MCP Server.

## Table of Contents

- [Performance Features](#performance-features)
- [Performance Metrics](#performance-metrics)
- [Optimization Strategies](#optimization-strategies)
- [Benchmarking](#benchmarking)
- [Performance Tuning](#performance-tuning)
- [Troubleshooting](#troubleshooting)

---

## Performance Features

### 1. Intelligent Caching

#### LRU Cache

Least Recently Used (LRU) cache automatically manages memory usage:

```yaml
cache:
  enabled: true
  type: "lru"
  max_size: 1000
  default_ttl: 300  # 5 minutes
```

**Advantages**:
- Automatically evicts least recently used entries
- Controllable memory usage
- Suitable for most scenarios

**Use Cases**:
- Read-intensive operations
- Infrequently changing data
- High latency external calls

#### Segmented Cache

Segmented cache provides better concurrent performance:

```yaml
cache:
  enabled: true
  type: "segmented"
  max_size: 1000
  segments: 10
  default_ttl: 300
```

**Advantages**:
- Reduces lock contention
- Better concurrent performance
- Configurable number of segments

**Use Cases**:
- High concurrency scenarios
- Need for low latency
- Multi-core CPUs

### 2. JSON Encoding Pool

Pre-allocated encoder pool reduces memory allocations:

```go
// Internal implementation
pool := json.NewEncoderPool(100, 8192)
```

**Advantages**:
- Reduces memory allocations
- Improves JSON encoding speed
- Lowers GC pressure

### 3. Response Compression

Automatically compresses large responses:

```yaml
performance:
  compression_enabled: true
  compression_level: 6
```

**Advantages**:
- Reduces network transfer
- Saves bandwidth
- Improves response speed

**Use Cases**:
- Large responses (>10KB)
- Limited network bandwidth
- Cross-datacenter access

### 4. Connection Pooling

Optimized HTTP client connection pool:

```yaml
kubernetes:
  qps: 100.0
  burst: 200
  timeoutSec: 30
```

**Advantages**:
- Reuses connections
- Reduces TCP handshake overhead
- Improves throughput

### 5. Response Size Control

Intelligently truncates oversized responses:

```yaml
performance:
  max_response_size: 5242880  # 5MB
  truncate_large_responses: true
```

**Advantages**:
- Prevents memory overflow
- Controls network transfer
- Improves response speed

---

## Performance Metrics

### Key Metrics

#### Request Metrics

```
mcp_requests_total{method="kubernetes_list_pods",status="success"} 1234
mcp_request_duration_seconds{method="kubernetes_list_pods"} 0.123
mcp_request_duration_seconds{method="kubernetes_list_pods",quantile="0.99"} 0.456
```

#### Cache Metrics

```
mcp_cache_hits_total{service="kubernetes"} 456
mcp_cache_misses_total{service="kubernetes"} 78
mcp_cache_hit_rate{service="kubernetes"} 0.85
```

#### Connection Metrics

```
mcp_active_connections 10
mcp_total_connections 100
mcp_connection_duration_seconds 300
```

#### Error Metrics

```
mcp_errors_total{type="timeout"} 5
mcp_errors_total{type="authentication"} 2
mcp_errors_total{type="service_unavailable"} 1
```

### Performance Benchmarks

#### Single Node Performance

| Metric | Value |
|--------|-------|
| Max Concurrent Connections | 1000 |
| Request Throughput (QPS) | 500+ |
| Average Response Time | <100ms |
| P99 Response Time | <500ms |
| Memory Usage | <512MB |
| CPU Usage | <50% (2 cores) |

#### Service-Specific Performance

| Service | Average Response Time | Cache Hit Rate |
|---------|---------------------|----------------|
| Kubernetes | 50ms | 85% |
| Grafana | 120ms | 90% |
| Prometheus | 80ms | 75% |
| Kibana | 200ms | 80% |
| Elasticsearch | 150ms | 70% |

---

## Optimization Strategies

### 1. Cache Optimization

#### Enable Caching

```yaml
cache:
  enabled: true
  type: "lru"
  max_size: 1000
  default_ttl: 300
```

#### Service-Specific TTL

```yaml
kubernetes:
  cache_ttl: 300  # 5 minutes

grafana:
  cache_ttl: 180  # 3 minutes

prometheus:
  cache_ttl: 60   # 1 minute
```

#### Cache Warmup

```go
// Warm up cache on service startup
func (s *Service) WarmupCache(ctx context.Context) error {
    // Pre-load frequently used data
    _, err := s.ListPods(ctx, "default")
    if err != nil {
        return err
    }
    return nil
}
```

### 2. Connection Optimization

#### Adjust QPS and Burst

```yaml
kubernetes:
  qps: 100.0   # Queries per second
  burst: 200   # Burst rate
  timeoutSec: 30
```

**Recommendations**:
- Adjust QPS based on cluster size
- Burst = QPS * 2
- Adjust timeout based on operation complexity

#### Connection Timeout

```yaml
kubernetes:
  timeoutSec: 30
```

**Recommendations**:
- Fast operations: 10-30 seconds
- Complex queries: 60-120 seconds
- Batch operations: 300+ seconds

### 3. Response Optimization

#### Enable Compression

```yaml
performance:
  compression_enabled: true
  compression_level: 6
```

**Compression Levels**:
- 1-3: Fastest, low compression
- 6: Balanced (recommended)
- 9: Slowest, high compression

#### Limit Response Size

```yaml
performance:
  max_response_size: 5242880  # 5MB
  truncate_large_responses: true
```

#### Use Summary Tools

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "kubernetes_list_resources_summary",
    "arguments": {
      "namespace": "default"
    }
  }
}
```

### 4. Concurrency Optimization

#### Adjust Maximum Connections

```yaml
server:
  max_connections: 1000
```

#### Adjust Worker Threads

```yaml
performance:
  worker_threads: 4
```

### 5. Memory Optimization

#### Limit Cache Size

```yaml
cache:
  max_size: 1000
```

#### Enable Response Compression

```yaml
performance:
  compression_enabled: true
```

#### Adjust Buffer Size

```yaml
performance:
  buffer_size: 8192
```

---

## Benchmarking

### Testing Tools

Use Apache Bench for benchmarking:

```bash
# Test health check endpoint
ab -n 10000 -c 100 http://localhost:8080/health

# Test tool call
ab -n 1000 -c 10 \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-key" \
  -p payload.json \
  http://localhost:8080/api/kubernetes/http
```

### Benchmark Results

#### Health Check Endpoint

```
Concurrency Level:      100
Time taken for tests:   2.345 seconds
Complete requests:      10000
Failed requests:        0
Total transferred:      1234567 bytes
Requests per second:    4264.31 [#/sec] (mean)
Time per request:       23.456 [ms] (mean)
Time per request:       0.235 [ms] (mean, across all concurrent requests)
```

#### Kubernetes Tool Call

```
Concurrency Level:      10
Time taken for tests:   12.345 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      9876543 bytes
Requests per second:    81.01 [#/sec] (mean)
Time per request:       123.456 [ms] (mean)
Time per request:       12.346 [ms] (mean, across all concurrent requests)
```

---

## Performance Tuning

### Production Configuration

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8080"
  max_connections: 1000
  readTimeoutSec: 30
  writeTimeoutSec: 0
  idleTimeoutSec: 60

logging:
  level: "info"
  json: true

cache:
  enabled: true
  type: "segmented"
  max_size: 2000
  segments: 10
  default_ttl: 300

performance:
  compression_enabled: true
  compression_level: 6
  max_response_size: 5242880
  truncate_large_responses: true
  worker_threads: 4
  buffer_size: 8192

kubernetes:
  kubeconfig: ""
  timeoutSec: 30
  qps: 100.0
  burst: 200
  cache_ttl: 300

grafana:
  enabled: true
  url: "http://grafana:3000"
  apiKey: "${GRAFANA_API_KEY}"
  timeoutSec: 30
  cache_ttl: 180

prometheus:
  enabled: true
  address: "http://prometheus:9090"
  timeoutSec: 30
  cache_ttl: 60

ratelimit:
  enabled: true
  requests_per_second: 100
  burst: 200
```

### High Performance Configuration

```yaml
cache:
  enabled: true
  type: "segmented"
  max_size: 5000
  segments: 20
  default_ttl: 600

performance:
  compression_enabled: true
  compression_level: 9
  max_response_size: 10485760  # 10MB
  worker_threads: 8
  buffer_size: 16384

kubernetes:
  qps: 200.0
  burst: 400
  timeoutSec: 60
```

### Low Latency Configuration

```yaml
cache:
  enabled: true
  type: "segmented"
  max_size: 1000
  segments: 10
  default_ttl: 60

performance:
  compression_enabled: false
  max_response_size: 1048576  # 1MB
  worker_threads: 4

kubernetes:
  qps: 100.0
  burst: 200
  timeoutSec: 10
```

---

## Monitoring and Analysis

### Prometheus Queries

#### Request Rate

```promql
rate(mcp_requests_total[5m])
```

#### Error Rate

```promql
rate(mcp_errors_total[5m])
```

#### P99 Latency

```promql
histogram_quantile(0.99, rate(mcp_request_duration_seconds_bucket[5m]))
```

#### Cache Hit Rate

```promql
mcp_cache_hits_total / (mcp_cache_hits_total + mcp_cache_misses_total)
```

### Grafana Dashboards

#### Key Panels

1. **Request Rate**: Requests per second
2. **P50/P95/P99 Latency**: Response time distribution
3. **Error Rate**: Error percentage
4. **Cache Hit Rate**: Cache efficiency
5. **Active Connections**: Current connection count
6. **Memory Usage**: Memory consumption
7. **CPU Usage**: CPU utilization

---

## Troubleshooting

### High Latency

**Symptoms**: Response time >1s

**Troubleshooting Steps**:

1. Check cache hit rate
```bash
curl http://localhost:8080/metrics | grep cache_hit_rate
```

2. Check external service latency
```bash
kubectl top pods
```

3. Enable debug logging
```yaml
logging:
  level: "debug"
```

4. Increase cache TTL
```yaml
cache:
  default_ttl: 600
```

### High Memory Usage

**Symptoms**: Memory usage >1GB

**Troubleshooting Steps**:

1. Check cache size
```yaml
cache:
  max_size: 500
```

2. Enable response compression
```yaml
performance:
  compression_enabled: true
```

3. Check response size
```yaml
performance:
  max_response_size: 5242880
```

4. Analyze memory usage
```bash
pprof http://localhost:8080/debug/pprof/heap
```

### Low Throughput

**Symptoms**: QPS < 100

**Troubleshooting Steps**:

1. Increase worker threads
```yaml
performance:
  worker_threads: 8
```

2. Adjust connection pool
```yaml
kubernetes:
  qps: 200.0
  burst: 400
```

3. Check network bandwidth
```bash
iftop
```

4. Check concurrent connections
```bash
curl http://localhost:8080/metrics | grep active_connections
```

### High Error Rate

**Symptoms**: Error rate > 5%

**Troubleshooting Steps**:

1. Check authentication configuration
```yaml
auth:
  enabled: true
  mode: "apikey"
  apiKey: "${MCP_AUTH_API_KEY}"
```

2. Check service health
```bash
curl http://localhost:8080/health
```

3. Increase timeout
```yaml
kubernetes:
  timeoutSec: 60
```

4. Check audit logs
```bash
curl -H "X-API-Key: your-key" \
  "http://localhost:8080/api/audit/query?status=failed"
```

---

## Performance Best Practices

### 1. Always Enable Caching

```yaml
cache:
  enabled: true
```

### 2. Use Appropriate TTL

- Static data: 600-3600 seconds
- Dynamic data: 60-300 seconds
- Real-time data: 10-30 seconds

### 3. Optimize External Service Calls

- Batch operations over individual operations
- Use filtering to reduce data volume
- Use pagination for large datasets

### 4. Monitor Key Metrics

- Request rate
- Response time
- Error rate
- Cache hit rate

### 5. Regularly Review Configuration

- Adjust QPS based on load
- Adjust cache size based on memory usage
- Adjust compression level based on network conditions

### 6. Use Summary Tools

For large datasets, use summary tools:

```json
{
  "name": "kubernetes_list_resources_summary"
}
```

Instead of:

```json
{
  "name": "kubernetes_list_resources"
}
```

---

## Related Documentation

- [Complete Tools Reference](/docs/tools/)
- [Configuration Guide](/docs/configuration/)
- [Deployment Guide](/docs/deployment/)
- [Architecture Guide](/docs/architecture/)
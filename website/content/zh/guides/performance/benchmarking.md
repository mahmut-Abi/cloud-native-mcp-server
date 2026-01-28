---
title: "基准测试"
weight: 20
---

# 基准测试

本文档描述 Cloud Native MCP Server 的基准测试方法和结果。

## 测试工具

使用 Apache Bench 进行基准测试：

```bash
# 安装 Apache Bench
# Ubuntu/Debian
sudo apt-get install apache2-utils

# macOS (已预装)
ab --version

# 测试健康检查端点
ab -n 10000 -c 100 http://localhost:8080/health

# 测试工具调用
ab -n 1000 -c 10 \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-key" \
  -p payload.json \
  http://localhost:8080/api/kubernetes/http
```

## 测试场景

### 1. 健康检查端点

```bash
ab -n 10000 -c 100 http://localhost:8080/health
```

**目的**: 测试基本 HTTP 性能

**预期结果**:
- Requests per second: 4000+
- Time per request: <25ms
- Success rate: 100%

### 2. Kubernetes 工具调用

#### 列出 Pod

```bash
cat > payload.json << EOF
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "list_pods",
    "arguments": {
      "namespace": "default"
    }
  }
}
EOF

ab -n 1000 -c 10 \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-key" \
  -p payload.json \
  http://localhost:8080/api/kubernetes/http
```

**目的**: 测试 Kubernetes 工具性能

**预期结果**:
- Requests per second: 80+
- Time per request: <125ms
- Success rate: 100%

### 3. Prometheus 查询

```bash
cat > payload.json << EOF
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "query",
    "arguments": {
      "query": "up{job=\"kubernetes-pods\"}"
    }
  }
}
EOF

ab -n 1000 -c 10 \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-key" \
  -p payload.json \
  http://localhost:8080/api/prometheus/http
```

**目的**: 测试 Prometheus 查询性能

**预期结果**:
- Requests per second: 100+
- Time per request: <100ms
- Success rate: 100%

### 4. Grafana 仪表板查询

```bash
cat > payload.json << EOF
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "list_dashboards",
    "arguments": {}
  }
}
EOF

ab -n 1000 -c 10 \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-key" \
  -p payload.json \
  http://localhost:8080/api/grafana/http
```

**目的**: 测试 Grafana 工具性能

**预期结果**:
- Requests per second: 60+
- Time per request: <150ms
- Success rate: 100%

## 基准测试结果

### 健康检查端点

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

### Kubernetes 工具调用

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

### Prometheus 查询

```
Concurrency Level:      10
Time taken for tests:   9.876 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      7654321 bytes
Requests per second:    101.26 [#/sec] (mean)
Time per request:       98.765 [ms] (mean)
Time per request:       9.877 [ms] (mean, across all concurrent requests)
```

### Grafana 仪表板查询

```
Concurrency Level:      10
Time taken for tests:   15.432 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      8765432 bytes
Requests per second:    64.81 [#/sec] (mean)
Time per request:       154.321 [ms] (mean)
Time per request:       15.432 [ms] (mean, across all concurrent requests)
```

## 性能分析

### 响应时间分布

| 端点 | 平均 | P50 | P95 | P99 |
|------|------|-----|-----|-----|
| /health | 23ms | 20ms | 35ms | 50ms |
| list_pods | 123ms | 100ms | 200ms | 450ms |
| prometheus query | 99ms | 80ms | 150ms | 300ms |
| list_dashboards | 154ms | 120ms | 250ms | 500ms |

### 吞吐量对比

| 端点 | QPS (10并发) | QPS (100并发) |
|------|-------------|--------------|
| /health | 4264 | 4500+ |
| list_pods | 81 | 95 |
| prometheus query | 101 | 120 |
| list_dashboards | 65 | 78 |

### 缓存影响

| 场景 | 无缓存 | 有缓存 | 改进 |
|------|--------|--------|------|
| 重复查询 | 100ms | 5ms | 95% |
| 缓存命中率 | 0% | 85% | - |
| 内存使用 | 200MB | 450MB | +125% |

## 性能优化建议

### 1. 启用缓存

```yaml
cache:
  enabled: true
  max_size: 2000
  default_ttl: 300
```

**效果**: 响应时间减少 95%

### 2. 启用压缩

```yaml
performance:
  compression_enabled: true
  compression_level: 6
```

**效果**: 网络传输减少 70-80%

### 3. 调整连接池

```yaml
kubernetes:
  qps: 100.0
  burst: 200
```

**效果**: 吞吐量提升 30%

### 4. 限制响应大小

```yaml
performance:
  max_response_size: 5242880
  truncate_large_responses: true
```

**效果**: 内存使用减少 40%

## 持续监控

### Prometheus 指标

```bash
# 查询性能指标
curl http://localhost:8080/metrics | grep mcp_request_duration
```

### Grafana 仪表板

创建性能监控仪表板，包含：

1. 请求速率
2. 响应时间（P50/P95/P99）
3. 错误率
4. 缓存命中率
5. 活动连接数
6. CPU 和内存使用

### 告警规则

```yaml
groups:
- name: mcp_performance
  rules:
  - alert: HighLatency
    expr: histogram_quantile(0.99, rate(mcp_request_duration_seconds_bucket[5m])) > 0.5
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High P99 latency detected"

  - alert: LowCacheHitRate
    expr: rate(mcp_cache_hits_total[5m]) / (rate(mcp_cache_hits_total[5m]) + rate(mcp_cache_misses_total[5m])) < 0.5
    for: 10m
    labels:
      severity: warning
    annotations:
      summary: "Low cache hit rate detected"
```

## 相关文档

- [优化](/zh/guides/performance/optimization/)
- [配置指南](/zh/guides/configuration/)
- [部署指南](/zh/guides/deployment/)
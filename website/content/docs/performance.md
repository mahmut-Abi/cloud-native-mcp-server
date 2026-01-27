---
title: "性能指南"
---

# 性能指南

本文档描述 Cloud Native MCP Server 的性能特性和优化建议。

## 目录

- [性能特性](#性能特性)
- [性能指标](#性能指标)
- [优化策略](#优化策略)
- [基准测试](#基准测试)
- [性能调优](#性能调优)
- [故障排查](#故障排查)

---

## 性能特性

### 1. 智能缓存

#### LRU 缓存

最近最少使用（LRU）缓存自动管理内存使用：

```yaml
cache:
  enabled: true
  type: "lru"
  max_size: 1000
  default_ttl: 300  # 5 分钟
```

**优势**:
- 自动淘汰最少使用的条目
- 内存使用可控
- 适合大多数场景

**适用场景**:
- 读取密集型操作
- 数据变化不频繁
- 高延迟外部调用

#### 分段缓存

分段缓存提供更好的并发性能：

```yaml
cache:
  enabled: true
  type: "segmented"
  max_size: 1000
  segments: 10
  default_ttl: 300
```

**优势**:
- 减少锁竞争
- 更好的并发性能
- 可配置的段数

**适用场景**:
- 高并发场景
- 需要低延迟
- 多核 CPU

### 2. JSON 编码池

预分配的编码器池减少内存分配：

```go
// 内部实现
pool := json.NewEncoderPool(100, 8192)
```

**优势**:
- 减少内存分配
- 提高 JSON 编码速度
- 降低 GC 压力

### 3. 响应压缩

自动压缩大型响应：

```yaml
performance:
  compression_enabled: true
  compression_level: 6
```

**优势**:
- 减少网络传输
- 节省带宽
- 提高响应速度

**适用场景**:
- 大型响应（>10KB）
- 网络带宽受限
- 跨数据中心访问

### 4. 连接池

优化的 HTTP 客户端连接池：

```yaml
kubernetes:
  qps: 100.0
  burst: 200
  timeoutSec: 30
```

**优势**:
- 重用连接
- 减少 TCP 握手开销
- 提高吞吐量

### 5. 响应大小控制

智能截断过大的响应：

```yaml
performance:
  max_response_size: 5242880  # 5MB
  truncate_large_responses: true
```

**优势**:
- 防止内存溢出
- 控制网络传输
- 提高响应速度

---

## 性能指标

### 关键指标

#### 请求指标

```
mcp_requests_total{method="kubernetes_list_pods",status="success"} 1234
mcp_request_duration_seconds{method="kubernetes_list_pods"} 0.123
mcp_request_duration_seconds{method="kubernetes_list_pods",quantile="0.99"} 0.456
```

#### 缓存指标

```
mcp_cache_hits_total{service="kubernetes"} 456
mcp_cache_misses_total{service="kubernetes"} 78
mcp_cache_hit_rate{service="kubernetes"} 0.85
```

#### 连接指标

```
mcp_active_connections 10
mcp_total_connections 100
mcp_connection_duration_seconds 300
```

#### 错误指标

```
mcp_errors_total{type="timeout"} 5
mcp_errors_total{type="authentication"} 2
mcp_errors_total{type="service_unavailable"} 1
```

### 性能基准

#### 单节点性能

| 指标 | 值 |
|------|-----|
| 最大并发连接 | 1000 |
| 请求吞吐量 (QPS) | 500+ |
| 平均响应时间 | <100ms |
| P99 响应时间 | <500ms |
| 内存使用 | <512MB |
| CPU 使用 | <50% (2核) |

#### 服务特定性能

| 服务 | 平均响应时间 | 缓存命中率 |
|------|------------|-----------|
| Kubernetes | 50ms | 85% |
| Grafana | 120ms | 90% |
| Prometheus | 80ms | 75% |
| Kibana | 200ms | 80% |
| Elasticsearch | 150ms | 70% |

---

## 优化策略

### 1. 缓存优化

#### 启用缓存

```yaml
cache:
  enabled: true
  type: "lru"
  max_size: 1000
  default_ttl: 300
```

#### 服务特定 TTL

```yaml
kubernetes:
  cache_ttl: 300  # 5 分钟

grafana:
  cache_ttl: 180  # 3 分钟

prometheus:
  cache_ttl: 60   # 1 分钟
```

#### 缓存预热

```go
// 在服务启动时预热缓存
func (s *Service) WarmupCache(ctx context.Context) error {
    // 预加载常用数据
    _, err := s.ListPods(ctx, "default")
    if err != nil {
        return err
    }
    return nil
}
```

### 2. 连接优化

#### 调整 QPS 和 Burst

```yaml
kubernetes:
  qps: 100.0   # 每秒查询数
  burst: 200   # 突发速率
  timeoutSec: 30
```

**建议**:
- QPS 根据集群规模调整
- Burst = QPS * 2
- Timeout 根据操作复杂度调整

#### 连接超时

```yaml
kubernetes:
  timeoutSec: 30
```

**建议**:
- 快速操作: 10-30 秒
- 复杂查询: 60-120 秒
- 批量操作: 300 秒+

### 3. 响应优化

#### 启用压缩

```yaml
performance:
  compression_enabled: true
  compression_level: 6
```

**压缩级别**:
- 1-3: 最快，压缩率低
- 6: 平衡（推荐）
- 9: 最慢，压缩率高

#### 限制响应大小

```yaml
performance:
  max_response_size: 5242880  # 5MB
  truncate_large_responses: true
```

#### 使用摘要工具

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

### 4. 并发优化

#### 调整最大连接数

```yaml
server:
  max_connections: 1000
```

#### 调整工作线程

```yaml
performance:
  worker_threads: 4
```

### 5. 内存优化

#### 限制缓存大小

```yaml
cache:
  max_size: 1000
```

#### 启用响应压缩

```yaml
performance:
  compression_enabled: true
```

#### 调整缓冲区大小

```yaml
performance:
  buffer_size: 8192
```

---

## 基准测试

### 测试工具

使用 Apache Bench 进行基准测试：

```bash
# 测试健康检查端点
ab -n 10000 -c 100 http://localhost:8080/health

# 测试工具调用
ab -n 1000 -c 10 \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-key" \
  -p payload.json \
  http://localhost:8080/api/kubernetes/http
```

### 基准测试结果

#### 健康检查端点

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

#### Kubernetes 工具调用

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

## 性能调优

### 生产环境配置

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

### 高性能配置

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

### 低延迟配置

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

## 监控和分析

### Prometheus 查询

#### 请求速率

```promql
rate(mcp_requests_total[5m])
```

#### 错误率

```promql
rate(mcp_errors_total[5m])
```

#### P99 延迟

```promql
histogram_quantile(0.99, rate(mcp_request_duration_seconds_bucket[5m]))
```

#### 缓存命中率

```promql
mcp_cache_hits_total / (mcp_cache_hits_total + mcp_cache_misses_total)
```

### Grafana 仪表板

#### 关键面板

1. **请求速率**: 每秒请求数
2. **P50/P95/P99 延迟**: 响应时间分布
3. **错误率**: 错误百分比
4. **缓存命中率**: 缓存效率
5. **活动连接**: 当前连接数
6. **内存使用**: 内存消耗
7. **CPU 使用**: CPU 利用率

---

## 故障排查

### 高延迟

**症状**: 响应时间 >1s

**排查步骤**:

1. 检查缓存命中率
```bash
curl http://localhost:8080/metrics | grep cache_hit_rate
```

2. 检查外部服务延迟
```bash
kubectl top pods
```

3. 启用调试日志
```yaml
logging:
  level: "debug"
```

4. 增加缓存 TTL
```yaml
cache:
  default_ttl: 600
```

### 高内存使用

**症状**: 内存使用 >1GB

**排查步骤**:

1. 检查缓存大小
```yaml
cache:
  max_size: 500
```

2. 启用响应压缩
```yaml
performance:
  compression_enabled: true
```

3. 检查响应大小
```yaml
performance:
  max_response_size: 5242880
```

4. 分析内存使用
```bash
pprof http://localhost:8080/debug/pprof/heap
```

### 低吞吐量

**症状**: QPS < 100

**排查步骤**:

1. 增加工作线程
```yaml
performance:
  worker_threads: 8
```

2. 调整连接池
```yaml
kubernetes:
  qps: 200.0
  burst: 400
```

3. 检查网络带宽
```bash
iftop
```

4. 检查并发连接
```bash
curl http://localhost:8080/metrics | grep active_connections
```

### 高错误率

**症状**: 错误率 > 5%

**排查步骤**:

1. 检查认证配置
```yaml
auth:
  enabled: true
  mode: "apikey"
  apiKey: "${MCP_AUTH_API_KEY}"
```

2. 检查服务健康
```bash
curl http://localhost:8080/health
```

3. 增加超时
```yaml
kubernetes:
  timeoutSec: 60
```

4. 检查审计日志
```bash
curl -H "X-API-Key: your-key" \
  "http://localhost:8080/api/audit/query?status=failed"
```

---

## 性能最佳实践

### 1. 始终启用缓存

```yaml
cache:
  enabled: true
```

### 2. 使用适当的 TTL

- 静态数据: 600-3600 秒
- 动态数据: 60-300 秒
- 实时数据: 10-30 秒

### 3. 优化外部服务调用

- 批量操作优于单个操作
- 使用过滤减少数据量
- 使用分页处理大量数据

### 4. 监控关键指标

- 请求速率
- 响应时间
- 错误率
- 缓存命中率

### 5. 定期审查配置

- 根据负载调整 QPS
- 根据内存使用调整缓存大小
- 根据网络条件调整压缩级别

### 6. 使用摘要工具

对于大型数据集，使用摘要工具：

```json
{
  "name": "kubernetes_list_resources_summary"
}
```

而不是：

```json
{
  "name": "kubernetes_list_resources"
}
```

---

## 相关文档

- [完整工具参考](/docs/tools/)
- [配置指南](/docs/configuration/)
- [部署指南](/docs/deployment/)
- [架构指南](/docs/architecture/)
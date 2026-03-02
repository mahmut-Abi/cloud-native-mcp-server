---
title: "性能指南"
weight: 60
description: "缓存、并发、压测与优化参数的实践建议。"
---

# 性能指南

本指南仅包含与当前 Cloud Native MCP Server 实现一致的性能调优项。

---

## 当前可调优项

### 1. 服务器超时

```yaml
server:
  readTimeoutSec: 30
  writeTimeoutSec: 0
  idleTimeoutSec: 60
```

建议：
- SSE 场景保持 `writeTimeoutSec: 0`。
- 客户端较慢或请求体较大时，适当提高 `readTimeoutSec`。
- 长连接场景适当提高 `idleTimeoutSec`。

### 2. Kubernetes 客户端吞吐

```yaml
kubernetes:
  timeoutSec: 30
  qps: 100.0
  burst: 200
```

建议：
- 大集群可适当提高 `qps`/`burst`。
- 慢查询或重操作可提高 `timeoutSec`。

### 3. 请求限流

```yaml
ratelimit:
  enabled: true
  requests_per_second: 100
  burst: 200
```

建议：
- 多租户或公网入口建议开启。
- 从保守值起步，根据指标逐步上调。

### 4. 减少不必要服务开销

```yaml
enableDisable:
  enabledServices: ["kubernetes", "prometheus", "grafana", "aggregate"]
```

建议：
- 仅启用实际需要的服务。
- 禁用不用的工具可降低内存和启动开销。

### 5. 审计开销控制（启用审计时）

```yaml
audit:
  enabled: true
  storage: "database"
  sampling:
    enabled: true
    rate: 0.3
```

建议：
- 高 QPS 场景开启采样。
- 审计查询需求高时优先 `storage: database`。

---

## 内置优化（无公开 YAML 键）

服务内部已包含多项优化，例如：
- 响应截断保护
- 更高效的 JSON 处理路径
- 服务/工具层内部缓存与对象池

这些属于实现细节，不是稳定公开配置项。
优先调优上面的可观测配置（`server`、`kubernetes`、`ratelimit`、`enableDisable`、`audit`）。

---

## 建议关注的性能指标

重点指标：
- 请求速率
- p95/p99 延迟
- 错误率
- 活跃连接数
- 内存与 CPU 使用

示例查询：

```promql
rate(mcp_requests_total[5m])
```

```promql
histogram_quantile(0.99, rate(mcp_request_duration_seconds_bucket[5m]))
```

```promql
rate(mcp_errors_total[5m])
```

---

## 基准测试

### 健康检查端点

```bash
ab -n 10000 -c 100 http://127.0.0.1:8080/health
```

### 工具调用端点（message 端点，历史兼容）

先准备请求体：

```json
{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}
```

再执行压测：

```bash
ab -n 1000 -c 10 \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-key" \
  -p payload.json \
  http://127.0.0.1:8080/api/kubernetes/sse/message
```

---

## 生产基线示例

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

kubernetes:
  kubeconfig: ""
  timeoutSec: 30
  qps: 100.0
  burst: 200

ratelimit:
  enabled: true
  requests_per_second: 100
  burst: 200

auth:
  enabled: true
  mode: "apikey"
  apiKey: "${MCP_AUTH_API_KEY}"

audit:
  enabled: true
  storage: "database"
  sampling:
    enabled: true
    rate: 0.3

enableDisable:
  enabledServices: ["kubernetes", "prometheus", "grafana", "aggregate"]
```

---

## 故障排查

### 高延迟

操作建议：
1. 后端查询慢时提高 `kubernetes.timeoutSec`。
2. API 客户端瓶颈时提高 `kubernetes.qps`/`kubernetes.burst`。
3. 用 p99 延迟与后端负载联动分析。

### 高错误率

操作建议：
1. 校验认证模式和凭据。
2. 检查 Prometheus/Grafana 等后端可用性。
3. 启用审计时，使用 `/api/audit/logs` 排查失败请求。

### 内存偏高

操作建议：
1. 减少启用的服务/工具。
2. 开启审计采样。
3. 流量突发时，适当降低 `burst`。

---

## 最佳实践

1. 生产默认优先 `streamable-http`，除非客户端必须使用 SSE。
2. 一次只调整一个维度（`timeout` -> `qps/burst` -> `ratelimit`）。
3. 压测场景尽量贴近真实工具调用。
4. 重点关注 p95/p99，而非仅平均延迟。

---

## 相关文档

- [配置指南](/zh/docs/configuration/)
- [部署指南](/zh/docs/deployment/)
- [安全指南](/zh/docs/security/)
- [工具参考](/zh/docs/tools/)

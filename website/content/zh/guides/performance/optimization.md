---
title: "优化"
weight: 10
---

# 性能优化

本文档给出与当前版本一致的性能优化策略，仅使用实际可配置项。

## 1. 先做基线观测

先建立性能基线，再调参：
- 请求速率（RPS/QPS）
- p95/p99 延迟
- 错误率
- 内存与 CPU

建议至少跑 15-30 分钟真实流量或压测流量后再改配置。

## 2. 服务器超时调优

```yaml
server:
  readTimeoutSec: 30
  writeTimeoutSec: 0
  idleTimeoutSec: 60
```

建议：
- SSE 场景保持 `writeTimeoutSec: 0`。
- 高频断连时提高 `idleTimeoutSec`。
- 慢客户端场景提高 `readTimeoutSec`。

## 3. Kubernetes 调优

```yaml
kubernetes:
  timeoutSec: 30
  qps: 100.0
  burst: 200
```

建议：
- 集群越大，越需要更高 `qps/burst`。
- 超时错误增多时适当提升 `timeoutSec`。

## 4. 限流保护

```yaml
ratelimit:
  enabled: true
  requests_per_second: 100
  burst: 200
```

建议：
- 公网入口、共享集群建议默认开启。
- 先保守，再逐步上调。

## 5. 精简服务范围

```yaml
enableDisable:
  enabledServices: ["kubernetes", "prometheus", "grafana", "aggregate"]
```

建议：
- 只启用业务需要的服务。
- 关闭不用服务可以降低资源占用。

## 6. 审计成本控制

```yaml
audit:
  enabled: true
  storage: "database"
  sampling:
    enabled: true
    rate: 0.3
```

建议：
- 高流量场景建议启用采样。
- 需要长期检索时使用 database 存储。

## 7. 生产配置模板

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

## 8. 常见误区

- 不要使用历史字段，如 `cache_ttl`、`max_connections`、`max_response_size`。
- 不要一次性改太多参数，否则难以定位收益来源。
- 不要只看平均延迟，必须看 p95/p99。

## 相关文档

- [基准测试](/zh/guides/performance/benchmarking/)
- [配置指南](/zh/docs/configuration/)
- [部署指南](/zh/docs/deployment/)

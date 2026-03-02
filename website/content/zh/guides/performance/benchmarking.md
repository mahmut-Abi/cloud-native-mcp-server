---
title: "基准测试"
weight: 20
---

# 基准测试

本文档给出可复现的压测流程，用于验证 Cloud Native MCP Server 在当前配置下的性能表现。

## 1. 测试前准备

- 固定测试环境（CPU/内存配额、节点负载）
- 固定服务配置（避免测试中变更）
- 记录版本与配置摘要

建议记录：
- commit SHA
- 运行模式（`sse` 或 `streamable-http`）
- `kubernetes.timeoutSec/qps/burst`
- `ratelimit` 是否开启

## 2. 测试工具

推荐：Apache Bench (`ab`)。

```bash
# Ubuntu / Debian
sudo apt-get install apache2-utils

# macOS
ab -V
```

## 3. 压测场景

### 场景 A：健康检查

```bash
ab -n 10000 -c 100 http://127.0.0.1:8080/health
```

用途：测基础网络与进程响应能力。

### 场景 B：工具调用（message 端点，历史兼容）

创建 `payload.json`：

```json
{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}
```

压测命令：

```bash
ab -n 1000 -c 10 \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-key" \
  -p payload.json \
  http://127.0.0.1:8080/api/kubernetes/sse/message
```

用途：测真实业务链路下的吞吐与延迟。

## 4. 结果解读

重点字段：
- Requests per second
- Time per request
- Failed requests

建议同时采集：
- p95/p99 延迟
- CPU、内存
- 后端依赖（Kubernetes API、Prometheus/Grafana）延迟

## 5. 调优顺序

建议按以下顺序逐步调参：

1. `kubernetes.timeoutSec`
2. `kubernetes.qps` / `kubernetes.burst`
3. `ratelimit.requests_per_second` / `ratelimit.burst`
4. `server.readTimeoutSec` / `server.idleTimeoutSec`

每次只改一组参数，并重新压测。

## 6. 可用的生产基线配置

```yaml
server:
  mode: "streamable-http"
  addr: "0.0.0.0:8080"
  readTimeoutSec: 30
  writeTimeoutSec: 0
  idleTimeoutSec: 60

kubernetes:
  timeoutSec: 30
  qps: 100.0
  burst: 200

ratelimit:
  enabled: true
  requests_per_second: 100
  burst: 200
```

## 7. 常见问题

### 吞吐低

- 提高 `kubernetes.qps/burst`
- 检查下游服务是否成为瓶颈
- 检查限流是否过严

### 延迟高

- 提高 `kubernetes.timeoutSec`
- 检查 p99 与下游 API 延迟的相关性
- 减少同时启用的服务与工具范围

### 错误率高

- 检查认证配置和凭据
- 查看 `/api/audit/logs`（启用审计时）
- 检查后端服务可用性

## 相关文档

- [性能优化](/zh/guides/performance/optimization/)
- [配置指南](/zh/docs/configuration/)
- [部署指南](/zh/docs/deployment/)

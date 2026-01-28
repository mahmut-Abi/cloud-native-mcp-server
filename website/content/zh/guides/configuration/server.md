---
title: "服务器配置"
weight: 10
---

# 服务器配置

本文档描述 Cloud Native MCP Server 的服务器基本配置。

## 运行模式

服务器支持四种运行模式：

| 模式 | 描述 | 适用场景 |
|------|------|----------|
| `sse` | Server-Sent Events | 实时推送，生产环境推荐 |
| `streamable-http` | 流式 HTTP | 大数据量传输 |
| `http` | 标准 HTTP | 简单请求/响应 |
| `stdio` | 标准输入/输出 | 开发环境，本地测试 |

### SSE 模式

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8080"
  writeTimeoutSec: 0  # 保持连接
```

**优势**:
- 实时推送
- 长连接
- 适合流式数据

### Streamable-HTTP 模式

```yaml
server:
  mode: "streamable-http"
  addr: "0.0.0.0:8080"
  writeTimeoutSec: 0
```

**优势**:
- 大数据量传输
- 分块响应
- 更好的内存管理

### HTTP 模式

```yaml
server:
  mode: "http"
  addr: "0.0.0.0:8080"
  writeTimeoutSec: 30
```

**优势**:
- 简单请求/响应
- 无状态
- 易于缓存

### stdio 模式

```yaml
server:
  mode: "stdio"
```

**优势**:
- 无需网络
- 适合本地开发
- 易于测试

## 基本设置

```yaml
server:
  # 运行模式
  mode: "sse"

  # 服务器监听地址
  addr: "0.0.0.0:8080"

  # HTTP 读取超时（秒）
  # 0 = 无超时（生产环境不推荐）
  # 推荐: 30-60 秒
  readTimeoutSec: 30

  # HTTP 写入超时（秒）
  # SSE 连接应设置为 0 以保持连接
  writeTimeoutSec: 0

  # HTTP 空闲超时（秒）
  # 默认: 60 秒
  idleTimeoutSec: 60
```

## 路径配置

### SSE 路径配置

```yaml
server:
  ssePaths:
    # Kubernetes SSE 端点
    kubernetes: "/api/kubernetes/sse"

    # Grafana SSE 端点
    grafana: "/api/grafana/sse"

    # Prometheus SSE 端点
    prometheus: "/api/prometheus/sse"

    # Kibana SSE 端点
    kibana: "/api/kibana/sse"

    # Helm SSE 端点
    helm: "/api/helm/sse"

    # Alertmanager SSE 端点
    alertmanager: "/api/alertmanager/sse"

    # Elasticsearch SSE 端点
    elasticsearch: "/api/elasticsearch/sse"

    # Utilities SSE 端点
    utilities: "/api/utilities/sse"

    # 聚合所有服务的 SSE 端点
    aggregate: "/api/aggregate/sse"
```

### Streamable-HTTP 路径配置

```yaml
server:
  streamableHttpPaths:
    # Kubernetes Streamable-HTTP 端点
    kubernetes: "/api/kubernetes/streamable-http"

    # Grafana Streamable-HTTP 端点
    grafana: "/api/grafana/streamable-http"

    # Prometheus Streamable-HTTP 端点
    prometheus: "/api/prometheus/streamable-http"

    # Kibana Streamable-HTTP 端点
    kibana: "/api/kibana/streamable-http"

    # Helm Streamable-HTTP 端点
    helm: "/api/helm/streamable-http"

    # Alertmanager Streamable-HTTP 端点
    alertmanager: "/api/alertmanager/streamable-http"

    # Elasticsearch Streamable-HTTP 端点
    elasticsearch: "/api/elasticsearch/streamable-http"

    # Utilities Streamable-HTTP 端点
    utilities: "/api/utilities/streamable-http"

    # 聚合所有服务的 Streamable-HTTP 端点
    aggregate: "/api/aggregate/streamable-http"
```

## 命令行参数

| 参数 | 描述 | 默认值 |
|------|------|--------|
| `--mode` | 服务器模式 (sse, streamable-http, http, stdio) | sse |
| `--addr` | 监听地址 | 0.0.0.0:8080 |
| `--config` | 配置文件路径 | config.yaml |
| `--log-level` | 日志级别 (debug, info, warn, error) | info |

### 使用示例

```bash
# 使用 SSE 模式
./cloud-native-mcp-server --mode=sse --addr=0.0.0.0:8080

# 使用 stdio 模式
./cloud-native-mcp-server --mode=stdio

# 使用自定义配置文件
./cloud-native-mcp-server --config=/path/to/config.yaml

# 设置日志级别
./cloud-native-mcp-server --log-level=debug
```

## 环境变量

| 变量 | 描述 | 默认值 |
|------|------|--------|
| `MCP_MODE` | 服务器模式 | sse |
| `MCP_ADDR` | 监听地址 | 0.0.0.0:8080 |
| `MCP_LOG_LEVEL` | 日志级别 | info |

### 使用示例

```bash
# 设置环境变量
export MCP_MODE=sse
export MCP_ADDR=0.0.0.0:8080
export MCP_LOG_LEVEL=info

# 启动服务器
./cloud-native-mcp-server
```

## 性能配置

```yaml
server:
  # 最大连接数
  maxConnections: 1000

  # 读取缓冲区大小（字节）
  readBufferSize: 4096

  # 写入缓冲区大小（字节）
  writeBufferSize: 4096
```

## TLS/SSL 配置

在生产环境中使用 TLS/SSL 加密通信：

### 基本 TLS 配置

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8443"
  tls:
    certFile: "/path/to/cert.pem"
    keyFile: "/path/to/key.pem"
    minVersion: "TLS1.2"
    maxVersion: "TLS1.3"
```

### mTLS 配置

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8443"
  tls:
    certFile: "/path/to/server-cert.pem"
    keyFile: "/path/to/server-key.pem"
    clientAuth: "RequireAndVerifyClientCert"
    caFile: "/path/to/ca-cert.pem"
```

## 优雅关闭

服务器支持优雅关闭，处理流程：

1. 接收到 SIGTERM 信号
2. 停止接受新连接
3. 等待现有请求完成（最多 30 秒）
4. 关闭所有服务连接
5. 退出

```bash
# 发送 SIGTERM 信号
kill -TERM <pid>

# 或者使用
kill -15 <pid>
```

## 健康检查

服务器提供健康检查端点：

```bash
# 基本健康检查
curl http://localhost:8080/health

# 详细健康信息
curl http://localhost:8080/health/detailed

# 就绪检查
curl http://localhost:8080/ready
```

## 指标端点

Prometheus 指标在 `/metrics` 端点可用：

```bash
curl http://localhost:8080/metrics
```

## 相关文档

- [服务配置](/zh/guides/configuration/services/)
- [认证配置](/zh/guides/configuration/authentication/)
- [部署指南](/zh/guides/deployment/)
- [架构指南](/zh/concepts/architecture/)
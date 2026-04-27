# Cloud Native MCP Server

[![Go Report Card](https://goreportcard.com/badge/github.com/mahmut-Abi/cloud-native-mcp-server)](https://goreportcard.com/report/github.com/mahmut-Abi/cloud-native-mcp-server)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org)

[🇨🇳 中文文档](README-zh.md) | [🇬🇧 English](README.md)

一个高性能的模型上下文协议（MCP）服务器，用于 Kubernetes 和云原生基础设施管理，集成了 10 个服务和 250+ 工具。

## LLM 调用工具建议

- 优先使用摘要类和分页类工具，只有在确实需要更多字段时再切到全量详情工具。
- 当参数语义是对象或数组时，优先传结构化 JSON；很多 handler 仍兼容旧的 JSON 字符串写法。
- Prometheus 和 tracing 相关时间字段使用 RFC3339 格式。
- Kubernetes 的 cluster-scoped 资源不要传 `namespace`，namespaced 资源要传。
- Kibana 的部分 handler 兼容 `camelCase` 和 `snake_case` 两种参数别名，但仍建议优先使用 schema 中展示的字段名。

## 功能特性

- **多服务集成**: Kubernetes、Grafana、Prometheus、Kibana、Elasticsearch、Helm、Alertmanager、Jaeger、OpenTelemetry、Utilities
- **多协议支持**: SSE 和 streamable-http 模式
- **智能缓存**: 支持 TTL 的 LRU 缓存以优化性能
- **性能优化**: JSON 编码池、响应大小控制、智能限制
- **增强的身份验证**: 支持 API Key（复杂度要求）、Bearer Token（JWT 验证）、Basic Auth
- **密钥管理**: 安全的凭证存储和轮换
- **输入清理**: 防止注入攻击
- **审计日志**: 跟踪所有工具调用和操作
- **LLM 优化**: 摘要工具和分页以防止上下文溢出

## 服务概览

| 服务 | 描述 |
|------|------|
| **kubernetes** | 容器编排和资源管理（含节点维护与 rollout 工具） |
| **helm** | 应用包管理和部署 |
| **grafana** | 可视化、监控仪表板和告警 |
| **prometheus** | 指标收集、查询和监控 |
| **kibana** | 日志分析、可视化和数据探索 |
| **elasticsearch** | 日志存储、搜索和数据索引 |
| **alertmanager** | 告警规则管理和通知 |
| **jaeger** | 分布式追踪和性能分析 |
| **opentelemetry** | 指标、追踪和日志收集与分析 |
| **utilities** | 通用工具 |

## 快速开始

### 二进制文件

```bash
# 下载最新版本
curl -LO https://github.com/mahmut-Abi/cloud-native-mcp-server/releases/latest/download/cloud-native-mcp-server-linux-amd64
chmod +x cloud-native-mcp-server-linux-amd64

# 以 SSE 模式运行（默认）
./cloud-native-mcp-server-linux-amd64 --mode=sse --addr=0.0.0.0:8080
```

### Docker

```bash
docker run -d \
--name cloud-native-mcp-server \
-p 8080:8080 \
-v ~/.kube:/root/.kube:ro \
mahmutabi/cloud-native-mcp-server:latest
```

### 从源码构建

```bash
git clone https://github.com/mahmut-Abi/cloud-native-mcp-server.git
cd cloud-native-mcp-server

make build
./cloud-native-mcp-server --mode=sse --addr=0.0.0.0:8080
```

## API 端点

### SSE 模式

| 端点 | 描述 |
|------|------|
| `/api/aggregate/sse` | 所有服务（推荐）|
| `/api/kubernetes/sse` | Kubernetes 服务 |
| `/api/helm/sse` | Helm 服务 |
| `/api/grafana/sse` | Grafana 服务 |
| `/api/prometheus/sse` | Prometheus 服务 |
| `/api/kibana/sse` | Kibana 服务 |
| `/api/elasticsearch/sse` | Elasticsearch 服务 |
| `/api/alertmanager/sse` | Alertmanager 服务 |
| `/api/jaeger/sse` | Jaeger 服务 |
| `/api/opentelemetry/sse` | OpenTelemetry 服务 |
| `/api/utilities/sse` | Utilities 服务 |

### Streamable-HTTP 模式

使用 `--mode=streamable-http` 可暴露 MCP streamable HTTP 端点，例如：
- `/api/aggregate/streamable-http`

### SSE 联调自检

对运行中的服务做端到端校验（SSE 握手 + `initialize`）：

```bash
# 无鉴权
make sse-smoke BASE_URL=http://127.0.0.1:8080

# API Key 鉴权
API_KEY=your-key make sse-smoke BASE_URL=http://127.0.0.1:8080
```

## 文档

- [完整工具参考](docs/TOOLS.md) - 与当前服务清单对齐的工具指南
- [配置指南](docs/CONFIGURATION.md) - 配置选项和示例
- [部署指南](docs/DEPLOYMENT.md) - 部署策略和最佳实践
- [安全指南](docs/SECURITY.md) - 身份验证、密钥管理和安全最佳实践
- [架构指南](docs/ARCHITECTURE.md) - 系统架构和设计
- [性能指南](docs/PERFORMANCE.md) - 性能功能和调优

## 构建

```bash
# 为当前平台构建
make build

# 运行测试
make test

# 代码检查
make lint

# Docker 构建
make docker-build
```

## 贡献

欢迎贡献！请阅读我们的贡献指南并提交拉取请求。

## 许可证

MIT License - see [LICENSE](LICENSE) for details.

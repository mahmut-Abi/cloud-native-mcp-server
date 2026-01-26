# Kubernetes MCP Server

[![Go Report Card](https://goreportcard.com/badge/github.com/mahmut-Abi/k8s-mcp-server)](https://goreportcard.com/report/github.com/mahmut-Abi/k8s-mcp-server)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)

一个高性能的模型上下文协议（MCP）服务器，用于 Kubernetes 和云原生基础设施管理，集成了多个服务和工具。

## 功能特性

- **多服务集成**: Kubernetes、Grafana、Prometheus、Kibana、Elasticsearch、Helm、Alertmanager、Jaeger、Utilities
- **多协议支持**: SSE、HTTP 和 stdio 模式
- **智能缓存**: 支持 TTL 的 LRU 缓存以优化性能
- **性能优化**: JSON 编码池、响应大小控制、智能限制
- **身份验证**: 支持 API Key、Bearer Token、Basic Auth
- **审计日志**: 跟踪所有工具调用和操作
- **LLM 优化**: 摘要工具和分页以防止上下文溢出

## 服务概览

| 服务 | 描述 |
|------|------|
| **kubernetes** | 容器编排和资源管理 |
| **helm** | 应用包管理和部署 |
| **grafana** | 可视化、监控仪表板和告警 |
| **prometheus** | 指标收集、查询和监控 |
| **kibana** | 日志分析、可视化和数据探索 |
| **elasticsearch** | 日志存储、搜索和数据索引 |
| **alertmanager** | 告警规则管理和通知 |
| **jaeger** | 分布式追踪和性能分析 |
| **utilities** | 通用工具 |

## 快速开始

### 二进制文件

```bash
# 下载最新版本
curl -LO https://github.com/mahmut-Abi/k8s-mcp-server/releases/latest/download/k8s-mcp-server-linux-amd64
chmod +x k8s-mcp-server-linux-amd64

# 以 SSE 模式运行（默认）
./k8s-mcp-server-linux-amd64 --mode=sse --addr=0.0.0.0:8080

# 或 HTTP 模式
./k8s-mcp-server-linux-amd64 --mode=http --addr=0.0.0.0:8080
```

### Docker

```bash
docker run -d \
--name k8s-mcp-server \
-p 8080:8080 \
-v ~/.kube:/root/.kube:ro \
mahmutabi/k8s-mcp-server:latest
```

### 从源码构建

```bash
git clone https://github.com/mahmut-Abi/k8s-mcp-server.git
cd k8s-mcp-server

make build
./k8s-mcp-server --mode=sse --addr=0.0.0.0:8080
```

## API 端点

### SSE 模式

| 端点 | 描述 |
|------|------|
| `/api/kubernetes/sse` | Kubernetes 服务 |
| `/api/helm/sse` | Helm 服务 |
| `/api/grafana/sse` | Grafana 服务 |
| `/api/prometheus/sse` | Prometheus 服务 |
| `/api/kibana/sse` | Kibana 服务 |
| `/api/elasticsearch/sse` | Elasticsearch 服务 |
| `/api/alertmanager/sse` | Alertmanager 服务 |
| `/api/jaeger/sse` | Jaeger 服务 |
| `/api/utilities/sse` | Utilities 服务 |
| `/api/aggregate/sse` | 所有服务（推荐）|

### HTTP 模式

将上述端点中的 `/sse` 替换为 `/http`。

## 配置

### YAML 配置文件

```yaml
# config.yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8080"

logging:
  level: "info"

kubernetes:
  kubeconfig: ""
  timeoutSec: 30

auth:
  enabled: false
  mode: "apikey"
  apiKey: "your-secret-key"

grafana:
  enabled: false
  url: "http://grafana:3000"
  apiKey: ""

prometheus:
  enabled: false
  address: "http://prometheus:9090"

kibana:
  enabled: false
  url: "http://kibana:5601"

elasticsearch:
  enabled: false
  url: "http://elasticsearch:9200"

alertmanager:
  enabled: false
  url: "http://alertmanager:9093"

jaeger:
  enabled: false
  url: "http://jaeger:16686"

audit:
  enabled: false
  maxLogs: 1000
```

### 环境变量

```bash
export MCP_MODE=sse
export MCP_ADDR=0.0.0.0:8080
export MCP_LOG_LEVEL=info
export MCP_AUTH_ENABLED=false
export MCP_K8S_KUBECONFIG=~/.kube/config
```

### 命令行参数

```bash
./k8s-mcp-server \
  --mode=sse \
  --addr=0.0.0.0:8080 \
  --config=config.yaml \
  --log-level=info
```

## 可用工具

完整的工具列表和详细说明，请参阅 [TOOLS.md](docs/TOOLS.md)。

### 快速参考

#### Kubernetes 工具
- `kubernetes_list_resources_summary` - 列出资源（优化输出）
- `kubernetes_get_resource_summary` - 获取单个资源摘要
- `kubernetes_get_pod_logs` - 获取 Pod 日志
- `kubernetes_get_events` - 获取集群事件
- `kubernetes_describe_resource` - 详细描述资源

#### Helm 工具
- `helm_list_releases_paginated` - 列出发布（分页）
- `helm_get_release_summary` - 获取发布摘要
- `helm_search_charts` - 搜索 Helm charts
- `helm_cluster_overview` - 获取集群概览

#### Grafana 工具
- `grafana_dashboards_summary` - 列出仪表板（最小输出）
- `grafana_datasources_summary` - 列出数据源
- `grafana_dashboard` - 获取特定仪表板
- `grafana_alerts` - 列出告警规则

#### Prometheus 工具
- `prometheus_query` - 执行即时查询
- `prometheus_query_range` - 执行范围查询
- `prometheus_alerts_summary` - 获取告警摘要
- `prometheus_targets_summary` - 获取目标摘要

#### Kibana 工具
- `kibana_search_saved_objects` - 搜索保存的对象
- `kibana_get_index_patterns` - 获取索引模式
- `kibana_get_spaces` - 获取 Kibana 空间

#### Elasticsearch 工具
- `elasticsearch_list_indices_paginated` - 列出索引（分页）
- `elasticsearch_cluster_health_summary` - 获取集群健康状态
- `elasticsearch_search_indices` - 搜索索引

#### Alertmanager 工具
- `alertmanager_alerts_summary` - 获取告警摘要
- `alertmanager_silences_summary` - 获取静默摘要
- `alertmanager_create_silence` - 创建静默

#### Jaeger 工具
- `jaeger_get_traces_summary` - 获取追踪摘要
- `jaeger_get_trace` - 获取特定追踪
- `jaeger_get_services` - 获取所有服务

#### Utilities 工具
- `utilities_get_time` - 获取当前时间
- `utilities_get_timestamp` - 获取 Unix 时间戳
- `utilities_web_fetch` - 获取 URL 内容

## LLM 优化工具

许多工具都有 LLM 优化版本，标记为 ⚠️ PRIORITY，提供：
- 70-95% 更小的响应大小
- 仅包含必要字段
- 分页支持
- 防止上下文溢出

示例：
- `kubernetes_list_resources_summary` vs `kubernetes_list_resources`
- `grafana_dashboards_summary` vs `grafana_dashboards`
- `prometheus_alerts_summary` vs `prometheus_get_alerts`

## 项目结构

```
k8s-mcp-server/
├── cmd/
│   └── server/              # 主入口
├── internal/
│   ├── config/              # 配置管理
│   ├── logging/             # 日志工具
│   ├── middleware/          # HTTP 中间件（auth、audit、metrics）
│   ├── observability/       # 指标和监控
│   ├── services/            # 服务实现
│   │   ├── kubernetes/      # Kubernetes 服务
│   │   ├── helm/            # Helm 服务
│   │   ├── grafana/         # Grafana 服务
│   │   ├── prometheus/      # Prometheus 服务
│   │   ├── kibana/          # Kibana 服务
│   │   ├── elasticsearch/   # Elasticsearch 服务
│   │   ├── alertmanager/    # Alertmanager 服务
│   │   ├── jaeger/          # Jaeger 服务
│   │   ├── utilities/       # Utilities 服务
│   │   ├── cache/           # LRU 缓存实现
│   │   ├── framework/       # 服务初始化框架
│   │   └── manager/         # 服务管理器
│   └── util/                # 工具
│       ├── circuitbreaker/  # 熔断器模式
│       ├── performance/     # 性能优化
│       └── pool/            # 对象池
├── docs/                    # 文档
│   └── TOOLS.md            # 完整工具参考
└── deploy/                  # 部署文件
    ├── Dockerfile
    ├── helm/
    │   └── k8s-mcp-server/
    └── kubernetes/
```

## 构建

```bash
# 为当前平台构建
make build

# 为所有平台构建
make build-all

# 运行测试
make test

# 运行竞态检测
make test-race

# 代码检查
make lint

# Docker 构建
make docker-build
```

## 性能特性

- **智能缓存**: 支持 TTL 的 LRU 缓存用于频繁访问的数据
- **响应大小控制**: 自动截断和优化
- **JSON 编码池**: 重用 JSON 编码器以提升性能
- **熔断器**: 防止级联故障
- **分页**: 支持大数据集
- **摘要工具**: 为 LLM 消费优化的工具

## 文档

- [完整工具参考](docs/TOOLS.md) - 所有工具的详细文档
- [配置指南](docs/CONFIGURATION.md) - 配置选项和示例
- [部署指南](docs/DEPLOYMENT.md) - 部署策略和最佳实践

## 贡献

欢迎贡献！请阅读我们的贡献指南并提交拉取请求。

## 许可证

MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。
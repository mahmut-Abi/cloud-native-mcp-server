---
title: "配置指南"
---

# 配置指南

本指南涵盖 Cloud Native MCP Server 的所有配置选项。

## 目录

- [配置方法](#配置方法)
- [服务器配置](#服务器配置)
- [服务配置](#服务配置)
- [认证配置](#认证配置)
- [日志配置](#日志配置)
- [审计日志](#审计日志)
- [缓存配置](#缓存配置)
- [性能调优](#性能调优)
- [示例配置](#示例配置)

---

## 配置方法

K8s MCP Server 支持三种配置方法（按优先级排序）：

1. **命令行参数** - 最高优先级
2. **环境变量** - 中等优先级
3. **YAML 配置文件** - 最低优先级

### 配置优先级示例

```bash
# 配置文件设置默认值
# 环境变量覆盖配置文件
# 命令行参数覆盖所有设置

./cloud-native-mcp-server \
  --config=config.yaml \
  --log-level=debug
```

---

## 服务器配置

### 基本设置

```yaml
server:
  # 运行模式: sse | streamable-http | http | stdio
  # 推荐开发环境使用 stdio，生产环境使用 streamable-http
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

### 命令行参数

| 参数 | 描述 | 默认值 |
|------|------|--------|
| `--mode` | 服务器模式 (sse, streamable-http, http, stdio) | sse |
| `--addr` | 监听地址 | 0.0.0.0:8080 |
| `--config` | 配置文件路径 | config.yaml |
| `--log-level` | 日志级别 (debug, info, warn, error) | info |

### 环境变量

| 变量 | 描述 | 默认值 |
|------|------|--------|
| `MCP_MODE` | 服务器模式 | sse |
| `MCP_ADDR` | 监听地址 | 0.0.0.0:8080 |
| `MCP_LOG_LEVEL` | 日志级别 | info |

---

## 服务配置

### Kubernetes

```yaml
kubernetes:
  # kubeconfig 文件路径
  # 如果为空，使用默认: $KUBECONFIG → ~/.kube/config → service account
  kubeconfig: ""

  # 单个 API 调用超时（秒）
  timeoutSec: 30

  # API 客户端每秒查询数 (QPS)
  qps: 100.0

  # API 客户端突发速率
  burst: 200
```

### Prometheus

```yaml
prometheus:
  # 启用/禁用 Prometheus 服务
  enabled: false

  # Prometheus 服务器地址
  # 格式: http://host:port 或 https://host:port
  address: "http://localhost:9090"

  # 请求超时（秒）
  timeoutSec: 30

  # Basic auth 用户名（可选）
  username: ""

  # Basic auth 密码（可选）
  password: ""

  # Bearer token 认证（可选，优先级高于 Basic Auth）
  bearerToken: ""

  # 跳过 TLS 证书验证
  # 生产环境不要使用！
  tlsSkipVerify: false

  # TLS 客户端证书文件路径（用于 mTLS 认证）
  tlsCertFile: ""

  # TLS 客户端密钥文件路径
  tlsKeyFile: ""

  # TLS CA 证书文件路径
  tlsCAFile: ""
```

### Grafana

```yaml
grafana:
  # 启用/禁用 Grafana 服务
  enabled: false

  # Grafana 服务器 URL
  # 格式: http://host:port 或 https://host:port
  url: "http://localhost:3000"

  # Grafana API Key（推荐）
  # 在 Grafana 中创建: Administration → API Keys
  apiKey: ""

  # Basic auth 用户名（API Key 的替代方案）
  username: ""

  # Basic auth 密码
  password: ""

  # 请求超时（秒）
  timeoutSec: 30
```

### Kibana

```yaml
kibana:
  # 启用/禁用 Kibana 服务
  enabled: false

  # Kibana 服务器 URL
  # 格式: http://host:port 或 https://host:port
  url: "https://localhost:5601"

  # Kibana API Key（推荐）
  # 在 Kibana 中创建: Stack Management → API Keys
  apiKey: ""

  # Basic auth 用户名（API Key 的替代方案）
  username: ""

  # Basic auth 密码
  password: ""

  # 请求超时（秒）
  timeoutSec: 30

  # 跳过 TLS 证书验证
  # 生产环境不要使用！
  skipVerify: false

  # Kibana 空间名称
  # 默认: "default"
  space: "default"
```

### Helm

```yaml
helm:
  # 启用/禁用 Helm 服务
  enabled: false

  # Helm 操作 kubeconfig 路径
  # 如果为空，使用与 Kubernetes 客户端相同的 kubeconfig
  kubeconfigPath: ""

  # Helm 操作默认命名空间
  namespace: "default"

  # 启用 Helm 调试模式
  debug: false

  # 仓库更新超时（秒）
  # 默认: 300 (5 分钟)
  # 国内环境推荐: 600-900
  timeoutSec: 300

  # 最大重试次数
  # 失败的仓库更新的重试次数
  # 默认: 3
  # 推荐: 3-5
  maxRetries: 3

  # 启用镜像
  # 用于加速 Helm 仓库拉取
  # 默认: false
  useMirrors: false

  # 自定义镜像映射
  # 格式: 原始仓库 URL -> 镜像 URL
  mirrors: {}
```

### Elasticsearch

```yaml
elasticsearch:
  # 启用/禁用 Elasticsearch 服务
  enabled: false

  # Elasticsearch 服务器地址（支持多节点高可用）
  addresses:
    - "http://localhost:9200"

  # 单个 Elasticsearch 服务器地址（addresses 的替代方案）
  # 当 addresses 为空时使用
  address: ""

  # Basic auth 用户名
  username: ""

  # Basic auth 密码
  password: ""

  # Bearer token 认证（可选，优先级高于 Basic Auth）
  bearerToken: ""

  # API Key 认证（可选，最高优先级）
  # 格式: id:api_key
  apiKey: ""

  # 请求超时（秒）
  timeoutSec: 30

  # 跳过 TLS 证书验证
  # 生产环境不要使用！
  tlsSkipVerify: false

  # TLS 客户端证书文件路径（用于 mTLS 认证）
  tlsCertFile: ""

  # TLS 客户端密钥文件路径
  tlsKeyFile: ""

  # TLS CA 证书文件路径
  tlsCAFile: ""
```

### Alertmanager

```yaml
alertmanager:
  # 启用/禁用 Alertmanager 服务
  enabled: false

  # Alertmanager 服务器地址
  # 格式: http://host:port 或 https://host:port
  address: "http://localhost:9093"

  # 请求超时（秒）
  timeoutSec: 30

  # Basic auth 用户名（可选）
  username: ""

  # Basic auth 密码（可选）
  password: ""

  # Bearer token 认证（可选，优先级高于 Basic Auth）
  bearerToken: ""

  # 跳过 TLS 证书验证
  # 生产环境不要使用！
  tlsSkipVerify: false

  # TLS 客户端证书文件路径（用于 mTLS 认证）
  tlsCertFile: ""

  # TLS 客户端密钥文件路径
  tlsKeyFile: ""

  # TLS CA 证书文件路径
  tlsCAFile: ""
```

### OpenTelemetry

```yaml
opentelemetry:
  # 启用/禁用 OpenTelemetry 服务
  enabled: false

  # OpenTelemetry Collector 地址
  # 格式: http://host:port 或 https://host:port
  address: "http://localhost:4318"

  # 请求超时（秒）
  timeoutSec: 30

  # Basic auth 用户名（可选）
  username: ""

  # Basic auth 密码（可选）
  password: ""

  # Bearer token 认证（可选，优先级高于 Basic Auth）
  bearerToken: ""

  # 跳过 TLS 证书验证
  # 生产环境不要使用！
  tlsSkipVerify: false

  # TLS 客户端证书文件路径（用于 mTLS 认证）
  tlsCertFile: ""

  # TLS 客户端密钥文件路径
  tlsKeyFile: ""

  # TLS CA 证书文件路径
  tlsCAFile: ""
```

### Utilities

```yaml
utilities:
  # Utilities 服务始终启用
  enabled: true
```

---

## 认证配置

### API Key 认证

```yaml
auth:
  # 启用/禁用认证
  enabled: false

  # 认证模式: apikey | bearer | basic
  # apikey: X-API-Key 简单 API 密钥认证
  # bearer: Bearer Token (JWT) 认证
  # basic: HTTP Basic Auth
  mode: "apikey"

  # API Key（用于 apikey 模式）
  # 最少 8 字符，推荐 16+ 字符
  apiKey: ""

  # Bearer token（用于 bearer 模式）
  # 最少 16 字符推荐（JWT token）
  bearerToken: ""

  # Basic Auth 用户名
  username: ""

  # Basic Auth 密码
  password: ""

  # JWT 密钥（用于 JWT 验证）
  jwtSecret: ""

  # JWT 算法 (HS256, RS256, etc.)
  jwtAlgorithm: "HS256"
```

### 认证环境变量

| 变量 | 描述 |
|------|------|
| `MCP_AUTH_ENABLED` | 启用认证 (1, true, yes, on) |
| `MCP_AUTH_MODE` | 认证模式 (apikey, bearer, basic) |
| `MCP_AUTH_API_KEY` | API key 或 bearer token |
| `MCP_AUTH_USERNAME` | Basic auth 用户名 |
| `MCP_AUTH_PASSWORD` | Basic auth 密码 |
| `MCP_AUTH_JWT_SECRET` | JWT 密钥 |
| `MCP_AUTH_JWT_ALGORITHM` | JWT 算法 |

---

## 日志配置

```yaml
logging:
  # 日志级别: debug | info | warn | error
  level: "info"

  # 使用 JSON 格式日志
  # 适用于日志聚合系统 (ELK, Splunk, etc.)
  json: false
```

### 日志级别说明

- **debug**: 详细的调试信息，包括所有请求和响应
- **info**: 一般信息，包括重要操作和状态变化
- **warn**: 警告信息，不影响功能但需要注意
- **error**: 错误信息，功能受损

---

## 审计日志

### 基本配置

```yaml
audit:
  # 启用/禁用审计日志
  enabled: false

  # 审计日志级别: debug | info | warn | error
  level: "info"

  # 审计日志存储: stdout | file | database | all
  storage: "memory"

  # 日志格式: text | json
  # json: 结构化 JSON 格式，适合日志聚合
  # text: 人类可读的文本格式
  format: "json"

  # 最大查询结果数
  maxResults: 1000

  # 查询时间范围（天）
  timeRange: 90
```

### 文件存储配置

```yaml
audit:
  storage: "file"
  file:
    # 日志文件路径
    path: "/var/log/cloud-native-mcp-server/audit.log"

    # 最大日志文件大小 (MB)
    maxSizeMB: 100

    # 最大备份文件数
    maxBackups: 10

    # 最大日志文件年龄（天）
    maxAgeDays: 30

    # 压缩轮转的日志文件
    compress: true

    # 内存存储的最大日志数
    maxLogs: 10000
```

### 数据库存储配置

```yaml
audit:
  storage: "database"
  database:
    # 数据库类型: sqlite | postgresql | mysql
    type: "sqlite"

    # SQLite 数据库文件路径
    # 仅当 type="sqlite" 时使用
    sqlitePath: "/var/lib/cloud-native-mcp-server/audit.db"

    # PostgreSQL 连接字符串
    # 仅当 type="postgresql" 时使用
    # 格式: postgresql://user:password@host:port/dbname
    connectionString: ""

    # 数据库表名
    tableName: "audit_logs"

    # 最大记录数
    maxRecords: 100000

    # 清理间隔（小时）
    cleanupInterval: 24
```

### 查询 API 配置

```yaml
audit:
  query:
    # 启用查询 API
    enabled: true

    # 每个查询的最大结果数
    maxResults: 1000

    # 最大时间范围（天）
    timeRange: 90
```

### 敏感数据掩码配置

```yaml
audit:
  masking:
    # 启用掩码
    enabled: true

    # 要掩码的字段
    fields:
      - password
      - token
      - apiKey
      - secret
      - passwd
      - pwd
      - authorization

    # 掩码替换值
    maskValue: "***REDACTED***"
```

### 采样配置（高流量场景）

```yaml
audit:
  sampling:
    # 启用采样
    enabled: false

    # 采样率 (0-1)
    # 1.0 = 记录所有, 0.1 = 记录 10%
    rate: 1.0
```

---

## 服务和工具过滤

```yaml
enableDisable:
  # 禁用的服务（逗号分隔）
  disabledServices: []

  # 启用的服务（逗号分隔，覆盖禁用列表）
  enabledServices: []

  # 禁用的工具（逗号分隔）
  disabledTools: []
```

---

## 性能调优

### 响应大小控制

```yaml
# 在代码中实现
performance:
  max_response_size: 5242880  # 5MB
  truncate_large_responses: true
  compression_enabled: true
  compression_level: 6
```

### JSON 编码池

```yaml
# 在代码中实现
performance:
  json_pool_size: 100
  json_buffer_size: 8192
```

---

## 示例配置

### 最小配置（仅 Kubernetes）

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8080"

logging:
  level: "info"

kubernetes:
  kubeconfig: ""
```

### 完整监控栈

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8080"

logging:
  level: "info"
  json: false

kubernetes:
  kubeconfig: ""

grafana:
  enabled: true
  url: "http://localhost:3000"
  apiKey: "${GRAFANA_API_KEY}"

prometheus:
  enabled: true
  address: "http://localhost:9090"

alertmanager:
  enabled: true
  address: "http://localhost:9093"

audit:
  enabled: true
  storage: "memory"
  format: "json"
```

### 生产环境配置（认证和缓存）

```yaml
server:
  mode: "sse"
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

grafana:
  enabled: true
  url: "http://grafana:3000"
  apiKey: "${GRAFANA_API_KEY}"
  timeoutSec: 30

prometheus:
  enabled: true
  address: "http://prometheus:9090"
  timeoutSec: 30

auth:
  enabled: true
  mode: "apikey"
  apiKey: "${MCP_AUTH_API_KEY}"

audit:
  enabled: true
  storage: "database"
  database:
    type: "sqlite"
    sqlitePath: "/var/lib/cloud-native-mcp-server/audit.db"
    maxRecords: 100000
    cleanupInterval: 24
  format: "json"
  masking:
    enabled: true
    maskValue: "***REDACTED***"
```

---

## 配置验证

服务器在启动时验证配置。常见验证错误：

### 无效的服务器模式
```
Error: invalid server mode "invalid". Must be one of: sse, streamable-http, http, stdio
```

### 缺少必需字段
```
Error: missing required field "api_key" in auth configuration
```

### 无效的服务 URL
```
Error: invalid service URL "grafana:3000". Must include scheme (http/https)
```

---

## 环境变量替换

可以在 YAML 配置文件中使用环境变量：

```yaml
grafana:
  url: "${GRAFANA_URL}"
  apiKey: "${GRAFANA_API_KEY}"

auth:
  apiKey: "${MCP_AUTH_API_KEY}"
```

在启动服务器前设置环境变量：

```bash
export GRAFANA_URL="http://grafana:3000"
export GRAFANA_API_KEY="your-api-key"
export MCP_AUTH_API_KEY="your-mcp-key"

./cloud-native-mcp-server
```

---

## 测试配置

在不启动服务器的情况下测试配置：

```bash
# 检查配置文件语法
./cloud-native-mcp-server --config=config.yaml --validate-config
```

这将会：
- 解析配置文件
- 验证所有字段
- 检查服务连通性
- 报告任何错误

---

## 热重载

不支持热重载。重启服务器以应用配置更改：

```bash
# 发送 SIGTERM 以优雅关闭
kill -TERM <pid>

# 服务器将完成进行中的请求并退出
# 然后使用新配置启动
./cloud-native-mcp-server --config=new-config.yaml
```
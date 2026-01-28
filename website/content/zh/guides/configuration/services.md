---
title: "服务配置"
weight: 20
---

# 服务配置

本文档描述 Cloud Native MCP Server 集成的各个服务的配置选项。

## Kubernetes

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

### 配置建议

- **开发环境**: `qps: 50.0, burst: 100`
- **生产环境**: `qps: 100.0, burst: 200`
- **大规模集群**: `qps: 200.0, burst: 400`

## Prometheus

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

### 查询示例

```json
{
  "name": "query",
  "arguments": {
    "query": "up{job=\"kubernetes-pods\"}"
  }
}
```

## Grafana

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

### 创建 API Key

1. 登录 Grafana
2. 导航到 Configuration → API Keys
3. 点击 "Add API Key"
4. 输入名称和过期时间
5. 复制生成的 API Key

## Kibana

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

## Helm

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

### 镜像配置示例

```yaml
helm:
  useMirrors: true
  mirrors:
    "https://kubernetes-charts.storage.googleapis.com": "https://mirror.example.com/kubernetes-charts"
    "https://charts.bitnami.com/bitnami": "https://mirror.example.com/bitnami"
```

## Elasticsearch

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

### 多节点配置

```yaml
elasticsearch:
  enabled: true
  addresses:
    - "http://es-node-1:9200"
    - "http://es-node-2:9200"
    - "http://es-node-3:9200"
  username: "elastic"
  password: "${ES_PASSWORD}"
```

## Alertmanager

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

## OpenTelemetry

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

## Utilities

```yaml
utilities:
  # Utilities 服务始终启用
  enabled: true
```

## 服务启用/禁用

```yaml
enableDisable:
  # 禁用的服务（逗号分隔）
  disabledServices: []

  # 启用的服务（逗号分隔，覆盖禁用列表）
  enabledServices: []

  # 禁用的工具（逗号分隔）
  disabledTools: []
```

### 示例

```yaml
# 仅启用 Kubernetes 和 Prometheus
enableDisable:
  enabledServices: ["kubernetes", "prometheus"]
  disabledServices: ["elasticsearch", "kibana", "grafana"]

# 禁用特定工具
enableDisable:
  disabledTools: ["delete_pod", "delete_deployment"]
```

## 配置验证

服务器在启动时验证每个服务的配置：

### 检查项

1. 服务 URL 格式
2. 认证凭证
3. 网络连接性
4. 权限验证

### 常见错误

```
Error: failed to connect to grafana service: connection refused
Error: invalid elasticsearch address: missing scheme
Error: authentication failed for prometheus: invalid API key
```

## 相关文档

- [服务器配置](/zh/guides/configuration/server/)
- [认证配置](/zh/guides/configuration/authentication/)
- [工具参考](/zh/tools/)
- [部署指南](/zh/guides/deployment/)
---
title: "服务"
---

# 集成服务

Cloud Native MCP Server 集成了 10 个强大的云原生服务，提供 220+ 个工具，全面覆盖 Kubernetes 管理和应用部署、监控、日志分析等各个方面。

<div class="service-grid">

<div class="service-card">
  <h3>☸️ Kubernetes</h3>
  <p><span class="tool-count">28 个工具</span></p>
  <p>核心容器编排和资源管理，包括 Pod、Deployment、Service、ConfigMap、Secret 等资源的完整生命周期管理。</p>
  <p><strong>主要功能：</strong></p>
  <ul>
    <li>Pod 和容器管理</li>
    <li>应用部署和扩缩容</li>
    <li>服务发现和负载均衡</li>
    <li>配置和密钥管理</li>
    <li>命名空间和节点管理</li>
  </ul>
</div>

<div class="service-card">
  <h3>⚓ Helm</h3>
  <p><span class="tool-count">31 个工具</span></p>
  <p>Kubernetes 应用包管理器，简化应用的部署、升级和管理流程。</p>
  <p><strong>主要功能：</strong></p>
  <ul>
    <li>Chart 仓库管理</li>
    <li>Release 生命周期管理</li>
    <li>Values 配置管理</li>
    <li>依赖管理</li>
    <li>插件系统</li>
  </ul>
</div>

<div class="service-card">
  <h3>📊 Grafana</h3>
  <p><span class="tool-count">36 个工具</span></p>
  <p>开源的分析和可视化平台，用于监控和指标可视化。</p>
  <p><strong>主要功能：</strong></p>
  <ul>
    <li>仪表板管理</li>
    <li>数据源配置</li>
    <li>可视化创建</li>
    <li>告警管理</li>
    <li>用户和组织管理</li>
  </ul>
</div>

<div class="service-card">
  <h3>📈 Prometheus</h3>
  <p><span class="tool-count">20 个工具</span></p>
  <p>开源的监控和告警系统，用于收集和查询时间序列数据。</p>
  <p><strong>主要功能：</strong></p>
  <ul>
    <li>即时和范围查询</li>
    <li>标签和元数据查询</li>
    <li>目标管理</li>
    <li>规则和告警管理</li>
    <li>TSDB 和存储管理</li>
  </ul>
</div>

<div class="service-card">
  <h3>🔍 Kibana</h3>
  <p><span class="tool-count">52 个工具</span></p>
  <p>Elastic Stack 的数据可视化和管理界面，用于日志分析和数据探索。</p>
  <p><strong>主要功能：</strong></p>
  <ul>
    <li>索引和文档管理</li>
    <li>数据查询和聚合</li>
    <li>可视化和仪表板</li>
    <li>索引模式管理</li>
    <li>空间和权限管理</li>
  </ul>
</div>

<div class="service-card">
  <h3>🔎 Elasticsearch</h3>
  <p><span class="tool-count">14 个工具</span></p>
  <p>分布式搜索和分析引擎，用于日志存储和全文搜索。</p>
  <p><strong>主要功能：</strong></p>
  <ul>
    <li>索引管理</li>
    <li>文档操作</li>
    <li>数据搜索</li>
    <li>集群管理</li>
    <li>别名管理</li>
  </ul>
</div>

<div class="service-card">
  <h3>🚨 Alertmanager</h3>
  <p><span class="tool-count">15 个工具</span></p>
  <p>Prometheus 告警处理和路由系统，用于管理告警通知。</p>
  <p><strong>主要功能：</strong></p>
  <ul>
    <li>告警管理</li>
    <li>沉默规则</li>
    <li>告警路由</li>
    <li>通知配置</li>
    <li>规则组管理</li>
  </ul>
</div>

<div class="service-card">
  <h3>🔗 Jaeger</h3>
  <p><span class="tool-count">8 个工具</span></p>
  <p>分布式追踪平台，用于监控和排查微服务架构中的问题。</p>
  <p><strong>主要功能：</strong></p>
  <ul>
    <li>追踪查询</li>
    <li>服务发现</li>
    <li>依赖分析</li>
    <li>性能分析</li>
    <li>指标查询</li>
  </ul>
</div>

<div class="service-card">
  <h3>📡 OpenTelemetry</h3>
  <p><span class="tool-count">9 个工具</span></p>
  <p>统一的可观测性框架，用于收集指标、追踪和日志。</p>
  <p><strong>主要功能：</strong></p>
  <ul>
    <li>指标收集</li>
    <li>追踪管理</li>
    <li>日志聚合</li>
    <li>统一配置</li>
    <li>跨语言支持</li>
  </ul>
</div>

<div class="service-card">
  <h3>🛠️ Utilities</h3>
  <p><span class="tool-count">6 个工具</span></p>
  <p>通用工具集，提供常用的数据处理和转换功能。</p>
  <p><strong>主要功能：</strong></p>
  <ul>
    <li>Base64 编解码</li>
    <li>JSON 处理</li>
    <li>时间戳生成</li>
    <li>UUID 生成</li>
    <li>数据转换</li>
  </ul>
</div>

</div>

## 服务配置

每个服务都可以独立配置和启用/禁用。以下是配置示例：

### 启用服务

```yaml
# Kubernetes（默认启用）
kubernetes:
  kubeconfig: ""
  timeoutSec: 30
  qps: 100.0
  burst: 200

# Prometheus
prometheus:
  enabled: true
  address: "http://localhost:9090"
  timeoutSec: 30

# Grafana
grafana:
  enabled: true
  url: "http://localhost:3000"
  apiKey: "your-api-key"

# Kibana
kibana:
  enabled: true
  url: "https://localhost:5601"
  apiKey: "your-api-key"

# Elasticsearch
elasticsearch:
  enabled: true
  addresses:
    - "http://localhost:9200"
  timeoutSec: 30

# Helm
helm:
  enabled: true
  namespace: "default"
  timeoutSec: 300

# Alertmanager
alertmanager:
  enabled: true
  address: "http://localhost:9093"
  timeoutSec: 30

# Jaeger
jaeger:
  enabled: true
  address: "http://localhost:16686"
  timeoutSec: 30

# OpenTelemetry
opentelemetry:
  enabled: true
  address: "http://localhost:4318"
  timeoutSec: 30
```

## 服务过滤

可以通过配置文件或环境变量启用或禁用特定服务：

```yaml
enableDisable:
  # 禁用的服务
  disabledServices: []

  # 启用的服务（覆盖禁用列表）
  enabledServices: ["kubernetes", "helm", "prometheus", "grafana"]

  # 禁用的工具
  disabledTools: []
```

或使用环境变量：

```bash
export MCP_ENABLED_SERVICES=kubernetes,helm,prometheus,grafana
export MCP_DISABLED_SERVICES=elasticsearch,kibana
```

## 认证配置

每个服务都支持多种认证方式：

### Basic Auth

```yaml
prometheus:
  username: "admin"
  password: "password"
```

### API Key

```yaml
grafana:
  apiKey: "eyJrIjoi..."
```

### Bearer Token

```yaml
kibana:
  bearerToken: "eyJhbGci..."
```

### TLS/mTLS

```yaml
elasticsearch:
  tlsSkipVerify: false
  tlsCertFile: "/path/to/cert.pem"
  tlsKeyFile: "/path/to/key.pem"
  tlsCAFile: "/path/to/ca.pem"
```

## 最佳实践

### 1. 按需启用服务
只启用需要的服务，减少资源消耗和攻击面。

### 2. 合理配置超时
根据服务响应时间设置合适的超时值。

### 3. 使用环境变量
敏感信息（如 API Key、密码）使用环境变量配置。

### 4. 启用缓存
对频繁查询的服务启用缓存以提高性能。

### 5. 监控服务健康
定期检查各服务的健康状态。

### 6. 使用连接池
配置合适的 QPS 和 burst 参数。

### 7. 错误处理
为每个服务配置适当的错误处理和重试策略。

## 故障排查

### 服务连接失败
- 检查服务地址和端口是否正确
- 验证网络连通性
- 确认认证信息是否正确

### 认证失败
- 验证 API Key、用户名密码等认证信息
- 检查认证模式配置
- 确认服务端是否启用了认证

### 性能问题
- 启用缓存功能
- 调整超时设置
- 优化查询语句
- 增加 QPS 和 burst 配置

### 数据不一致
- 检查服务端数据源
- 验证查询参数
- 清除缓存后重试

## 扩展开发

如需添加新的服务，请参考开发文档：

1. 创建服务目录结构
2. 实现服务接口
3. 注册工具
4. 编写测试
5. 更新文档

详细指南请查看 [开发文档](https://github.com/mahmut-Abi/cloud-native-mcp-server/tree/main/docs/development)。
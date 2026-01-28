---
title: "工具参考"
weight: 1
---

# 工具参考

Cloud Native MCP Server 提供 220+ 个强大的工具，覆盖 Kubernetes 管理和应用部署、监控、日志分析等各个方面。

## 工具概览

服务器集成以下云原生服务和工具：

- **Kubernetes** - 28 个工具用于 Pod、Deployment、Service 等资源管理
- **Helm** - 31 个工具用于 Chart 和 Release 管理
- **Grafana** - 36 个工具用于仪表板、数据源和用户管理
- **Prometheus** - 20 个工具用于指标查询和规则管理
- **Kibana** - 52 个工具用于索引、文档和可视化管理
- **Elasticsearch** - 14 个工具用于索引、文档和集群管理
- **Alertmanager** - 15 个工具用于告警和通知管理
- **Jaeger** - 8 个工具用于分布式追踪和依赖分析
- **OpenTelemetry** - 9 个工具用于指标、追踪和日志管理
- **Utilities** - 6 个通用工具（Base64、JSON、UUID 等）

## 工具调用示例

所有工具都通过标准的 JSON-RPC 2.0 协议调用：

### Kubernetes - 列出 Pod

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "list_pods",
    "arguments": {
      "namespace": "default"
    }
  }
}
```

### Helm - 安装 Chart

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "install_chart",
    "arguments": {
      "chart": "nginx-ingress",
      "repo": "https://kubernetes.github.io/ingress-nginx",
      "release": "my-nginx",
      "namespace": "ingress-nginx"
    }
  }
}
```

### Prometheus - 查询指标

```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "query",
    "arguments": {
      "query": "up{job=\"kubernetes-pods\"}"
    }
  }
}
```

### Grafana - 查询仪表板

```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/call",
  "params": {
    "name": "list_dashboards",
    "arguments": {}
  }
}
```

## 通用参数

所有工具都支持以下通用参数：

- `timeout` - 请求超时时间（秒）
- `dry_run` - 试运行模式，不实际执行
- `verbose` - 详细输出模式

工具特定的参数请参考各服务的详细文档。

## 错误处理

工具调用可能返回以下错误：

- `InvalidParams` - 参数无效
- `NotFound` - 资源不存在
- `PermissionDenied` - 权限不足
- `Timeout` - 请求超时
- `InternalError` - 内部错误

### 错误响应格式

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32602,
    "message": "Invalid params",
    "data": {
      "details": "namespace is required"
    }
  }
}
```

## 下一步

- [Kubernetes 工具](/zh/tools/kubernetes/) - 28 个 Kubernetes 管理工具
- [Helm 工具](/zh/tools/helm/) - 31 个 Helm Chart 管理工具
- [Grafana 工具](/zh/tools/grafana/) - 36 个 Grafana 可视化工具
- [Prometheus 工具](/zh/tools/prometheus/) - 20 个 Prometheus 指标工具
- [Kibana 工具](/zh/tools/kibana/) - 52 个 Kibana 日志分析工具
- [Elasticsearch 工具](/zh/tools/elasticsearch/) - 14 个 Elasticsearch 搜索工具
- [Alertmanager 工具](/zh/tools/alertmanager/) - 15 个 Alertmanager 告警工具
- [Jaeger 工具](/zh/tools/jaeger/) - 8 个 Jaeger 追踪工具
- [OpenTelemetry 工具](/zh/tools/opentelemetry/) - 9 个 OpenTelemetry 遥测工具
- [Utilities 工具](/zh/tools/utilities/) - 6 个通用实用工具
---
title: "配置"
weight: 10
---

# 配置指南

本指南涵盖 Cloud Native MCP Server 的所有配置选项。

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

## 内容

- [服务器配置](/zh/guides/configuration/server/) - 服务器基本设置和运行模式
- [服务配置](/zh/guides/configuration/services/) - 各服务的配置选项
- [认证配置](/zh/guides/configuration/authentication/) - 认证和授权设置
- [日志配置](#日志配置) - 日志级别和格式
- [审计日志](#审计日志) - 审计日志配置
- [性能调优](#性能调优) - 性能优化配置

## 快速开始

### 最小配置

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

audit:
  enabled: true
  storage: "memory"
  format: "json"
```

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

## 热重载

不支持热重载。重启服务器以应用配置更改：

```bash
# 发送 SIGTERM 以优雅关闭
kill -TERM <pid>

# 服务器将完成进行中的请求并退出
# 然后使用新配置启动
./cloud-native-mcp-server --config=new-config.yaml
```

## 相关文档

- [服务器配置](/zh/guides/configuration/server/)
- [服务配置](/zh/guides/configuration/services/)
- [认证配置](/zh/guides/configuration/authentication/)
- [部署指南](/zh/guides/deployment/)
- [架构指南](/zh/concepts/architecture/)
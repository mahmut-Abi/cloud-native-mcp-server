---
title: "文档"
---

# 文档中心

欢迎来到 Cloud Native MCP Server 文档中心。这里包含了完整的使用指南、API 参考和最佳实践。

## 快速开始

- [安装指南](#安装)
- [配置指南](#配置)
- [快速入门](#快速入门)

## 核心文档

### [完整工具参考](/docs/tools/)
详细文档所有 220+ MCP 工具，包括使用方法和示例。

### [配置指南](/docs/configuration/)
完整的配置选项说明，包括服务器配置、服务配置、认证配置等。

### [部署指南](/docs/deployment/)
各种部署方式的详细说明，包括 Docker、Kubernetes、Helm 等。

### [安全指南](/docs/security/)
认证、授权、密钥管理和安全最佳实践。

### [架构指南](/docs/architecture/)
系统架构设计、组件说明和扩展方式。

### [性能指南](/docs/performance/)
性能优化建议、调优参数和最佳实践。

## 安装

### 系统要求

- Go 1.25 或更高版本
- Linux、macOS 或 Windows
- 用于 Kubernetes 集群的 kubeconfig 文件
- 可选：Docker（用于容器化部署）

### 二进制安装

```bash
# 下载最新版本
curl -LO https://github.com/mahmut-Abi/cloud-native-mcp-server/releases/latest/download/cloud-native-mcp-server-linux-amd64
chmod +x cloud-native-mcp-server-linux-amd64

# 验证安装
./cloud-native-mcp-server-linux-amd64 --version
```

### Docker 安装

```bash
docker pull mahmutabi/cloud-native-mcp-server:latest

# 验证安装
docker run --rm mahmutabi/cloud-native-mcp-server:latest --version
```

### 从源码构建

```bash
git clone https://github.com/mahmut-Abi/cloud-native-mcp-server.git
cd cloud-native-mcp-server

# 构建
make build

# 验证
./cloud-native-mcp-server --version
```

## 配置

### 基本配置

创建配置文件 `config.yaml`：

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8080"

kubernetes:
  kubeconfig: ""  # 使用默认 kubeconfig

logging:
  level: "info"
  json: false
```

### 启用服务

```yaml
prometheus:
  enabled: true
  address: "http://localhost:9090"

grafana:
  enabled: true
  url: "http://localhost:3000"
  apiKey: "your-api-key"
```

### 配置认证

```yaml
auth:
  enabled: true
  mode: "apikey"
  apiKey: "your-secure-api-key"
```

### 环境变量

所有配置都支持环境变量：

```bash
export MCP_MODE=sse
export MCP_ADDR=0.0.0.0:8080
export MCP_KUBECONFIG=/path/to/kubeconfig
export MCP_LOG_LEVEL=info
```

## 快速入门

### 1. 启动服务器

```bash
# 使用默认配置
./cloud-native-mcp-server --mode=sse --addr=0.0.0.0:8080

# 使用配置文件
./cloud-native-mcp-server --config=config.yaml
```

### 2. 连接到服务器

**SSE 模式：**

```bash
curl -N http://localhost:8080/api/aggregate/sse
```

**HTTP 模式：**

```bash
curl http://localhost:8080/api/aggregate/http \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}'
```

### 3. 使用工具

调用 Kubernetes 工具示例：

```bash
curl -N http://localhost:8080/api/kubernetes/sse \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "list_pods",
      "arguments": {
        "namespace": "default"
      }
    }
  }'
```

### 4. 在 MCP 客户端中使用

配置 MCP 客户端（如 Claude Desktop）：

```json
{
  "mcpServers": {
    "cloud-native": {
      "command": "/path/to/cloud-native-mcp-server",
      "args": ["--mode=stdio"]
    }
  }
}
```

## 运行模式

### SSE 模式（推荐生产环境）

实时双向通信，适合需要实时更新的场景。

```bash
./cloud-native-mcp-server --mode=sse --addr=0.0.0.0:8080
```

### HTTP 模式

标准 REST API，易于集成。

```bash
./cloud-native-mcp-server --mode=http --addr=0.0.0.0:8080
```

### stdio 模式（推荐开发环境）

标准输入输出，适合 MCP 客户端。

```bash
./cloud-native-mcp-server --mode=stdio
```

### Streamable-HTTP 模式

MCP 2025-11-25 规范，现代化通信方式。

```bash
./cloud-native-mcp-server --mode=streamable-http --addr=0.0.0.0:8080
```

## 故障排查

### 常见问题

**Q: 无法连接到 Kubernetes 集群**
- 检查 kubeconfig 文件路径
- 验证集群访问权限
- 检查网络连接

**Q: 服务响应慢**
- 启用缓存功能
- 调整超时设置
- 检查集群资源

**Q: 认证失败**
- 验证 API Key 配置
- 检查认证模式设置
- 确认令牌有效性

### 日志调试

启用调试日志：

```bash
./cloud-native-mcp-server --log-level=debug
```

或在配置文件中：

```yaml
logging:
  level: "debug"
  json: true
```

## 下一步

- 阅读完整 [工具参考文档](/docs/tools/)
- 查看 [部署指南](/docs/deployment/) 了解生产部署
- 学习 [安全最佳实践](/docs/security/)
- 探索 [性能优化技巧](/docs/performance/)

## 获取帮助

- 查看 [GitHub Issues](https://github.com/mahmut-Abi/cloud-native-mcp-server/issues)
- 阅读 [项目 Wiki](https://github.com/mahmut-Abi/cloud-native-mcp-server/wiki)
- 加入社区讨论
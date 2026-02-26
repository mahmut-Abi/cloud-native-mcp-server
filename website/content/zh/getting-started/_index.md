---
title: "快速开始"
weight: 1
---

# Cloud Native MCP Server 快速开始

欢迎使用 Cloud Native MCP Server！本指南将帮助您快速上手这个最强大的用于 Kubernetes 和云原生基础设施管理的 Model Context Protocol (MCP) 服务器。

## 概述

Cloud Native MCP Server 是一个高性能的 Model Context Protocol (MCP) 服务器，集成了 10 个服务和 220+ 工具，让 AI 助手能够轻松管理您的云原生基础设施。

### 您将学到

- 如何安装和部署 Cloud Native MCP Server
- 基本配置选项
- 如何使用核心服务
- 安全和性能的最佳实践

---

## 安装选项

选择最适合您环境的安装方法：

{{< tabs >}}
{{< tab "Docker" >}}
### Docker 安装

最简单的开始方式是使用 Docker：

```bash
# 拉取最新镜像
docker pull mahmutabi/cloud-native-mcp-server:latest

# 运行服务器
docker run -d \
  --name cloud-native-mcp-server \
  -p 8080:8080 \
  -v ~/.kube:/root/.kube:ro \
  -e MCP_SERVER_API_KEY="your-secure-api-key" \
  mahmutabi/cloud-native-mcp-server:latest
```

运行后，您可以通过 `http://localhost:8080` 访问服务器。
{{< /tab >}}

{{< tab "二进制" >}}
### 二进制安装

下载预编译的二进制文件：

```bash
# Linux (amd64)
curl -LO https://github.com/mahmut-Abi/cloud-native-mcp-server/releases/latest/download/cloud-native-mcp-server-linux-amd64
chmod +x cloud-native-mcp-server-linux-amd64

# 运行服务器
./cloud-native-mcp-server-linux-amd64 --mode=sse --addr=0.0.0.0:8080
```

二进制文件包含所有 10 个集成服务和 220+ 工具。
{{< /tab >}}

{{< tab "源码" >}}
### 源码编译

从源码构建用于开发或自定义：

```bash
# 克隆仓库
git clone https://github.com/mahmut-Abi/cloud-native-mcp-server.git
cd cloud-native-mcp-server

# 构建服务器
make build

# 使用默认设置运行
./cloud-native-mcp-server --mode=sse --addr=0.0.0.0:8080
```

确保您已安装 Go 1.25+。
{{< /tab >}}
{{< /tabs >}}

---

## 初始配置

安装后，您需要使用适当的认证和服务端点配置服务器。

### 认证设置

服务器支持多种认证方法：

```bash
# API 密钥（推荐用于生产环境）
export MCP_SERVER_API_KEY="your-very-secure-api-key-with-32-chars-minimum"

# 或 Bearer Token (JWT)
export MCP_SERVER_BEARER_TOKEN="your-jwt-token"

# 或 Basic Auth
export MCP_SERVER_BASIC_AUTH_USER="admin"
export MCP_SERVER_BASIC_AUTH_PASS="secure-password"
```

### 服务配置

如果服务可访问，服务器将自动检测和配置服务：

- Kubernetes: 需要 `~/.kube/config` 或集群内配置
- Prometheus: 默认连接到 `http://prometheus:9090`
- Grafana: 默认连接到 `http://grafana:3000`
- 等等...

---

## 您的第一个 MCP 调用

服务器运行后，您可以进行第一个 MCP 调用：

```bash
curl -X POST http://localhost:8080/v1/mcp/list-tools \
  -H "Authorization: Bearer your-api-key" \
  -H "Content-Type: application/json"
```

这将返回 10 个集成服务中所有 220+ 可用工具的列表。

---

## 集成服务概述

Cloud Native MCP Server 集成了 10 个核心服务：

{{< columns >}}
### 🔧 Kubernetes
使用 28 个专门的工具管理您的 Kubernetes 集群，包括部署、服务、配置映射、密钥等。
<--->

### 📦 Helm
使用 31 个工具部署和管理 Helm 图表，用于图表管理、发布和仓库。
{{< /columns >}}

{{< columns >}}
### 📊 Grafana
使用 36 个监控工具创建和管理仪表板、警报和数据源。
<--->

### 📈 Prometheus
使用 20 个可观测性工具查询指标、管理规则和配置警报。
{{< /columns >}}

{{< columns >}}
### 🔍 Kibana
使用 52 个 Elasticsearch 集成工具分析日志和可视化数据。
<--->

### ⚡ Elasticsearch
使用 14 个高级搜索工具进行索引、搜索和分析数据。
{{< /columns >}}

---

## 下一步

现在您已经安装和配置了 Cloud Native MCP Server，您可能想要：

- [配置认证和安全设置](/zh/guides/security/)
- [探索服务特定配置](/zh/guides/configuration/)
- [了解性能优化](/zh/guides/performance/)
- [查看完整的工具参考](/zh/docs/tools/)

### 快速链接

- [架构概述](/zh/concepts/architecture/)
- [安全最佳实践](/zh/guides/security/best-practices/)
- [性能调优](/zh/guides/performance/optimization/)
- [故障排除](/zh/docs/getting-started/)

---

## 支持和社区

需要帮助？查看这些资源：

- [GitHub Issues](https://github.com/mahmut-Abi/cloud-native-mcp-server/issues) 用于错误报告
- [GitHub Discussions](https://github.com/mahmut-Abi/cloud-native-mcp-server/discussions) 用于提问
- [文档](/) 用于完整参考
- [贡献指南](https://github.com/mahmut-Abi/cloud-native-mcp-server/blob/main/CONTRIBUTING.md) 了解如何参与
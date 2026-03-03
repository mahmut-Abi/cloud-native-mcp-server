---
title: "快速开始"
weight: 1
description: "完成安装、端点连通性验证与生产基线配置。"
---

# 快速开始

本指南用于快速完成 Cloud Native MCP Server 的安装、联通验证和生产基线配置。

## 本文将完成

- 在 `sse`、`streamable-http` 两种模式中选择运行方式
- 使用正确的环境变量启用认证
- 验证健康状态与 MCP 握手链路
- 进入 FAQ 与排障手册继续深化

---

## 前置条件

- 可用的 Kubernetes 访问凭据（`~/.kube/config` 或集群内凭据）
- Docker 环境或可执行 Linux 二进制的主机
- Go `1.25+`（仅源码构建时需要）
- 可访问将要集成的可观测性后端服务

---

## 安装方式

{{< tabs >}}
{{< tab "Docker" >}}
```bash
docker run -d \
  --name cloud-native-mcp-server \
  -p 8080:8080 \
  -v ~/.kube:/root/.kube:ro \
  -e MCP_AUTH_ENABLED=true \
  -e MCP_AUTH_MODE=apikey \
  -e MCP_AUTH_API_KEY='ChangeMe-Strong-Key-123!' \
  mahmutabi/cloud-native-mcp-server:latest
```
{{< /tab >}}

{{< tab "二进制" >}}
```bash
curl -LO https://github.com/mahmut-Abi/cloud-native-mcp-server/releases/latest/download/cloud-native-mcp-server-linux-amd64
chmod +x cloud-native-mcp-server-linux-amd64
./cloud-native-mcp-server-linux-amd64 --mode=sse --addr=0.0.0.0:8080
```
{{< /tab >}}

{{< tab "源码" >}}
```bash
git clone https://github.com/mahmut-Abi/cloud-native-mcp-server.git
cd cloud-native-mcp-server
make build
./cloud-native-mcp-server --mode=streamable-http --addr=0.0.0.0:8080
```
{{< /tab >}}
{{< /tabs >}}

---

## 运行模式选择

| 模式 | 适用场景 | 主要入口 |
| --- | --- | --- |
| `sse` | 兼容性优先的 MCP 客户端接入 | `/api/aggregate/sse` |
| `streamable-http` | 推荐的现代 MCP 传输（生产优先） | `/api/aggregate/streamable-http` |

---

## 首次验证

启动后执行以下检查：

```bash
# 健康检查
curl -sS http://127.0.0.1:8080/health

# SSE 握手与 initialize 全链路检查
make sse-smoke BASE_URL=http://127.0.0.1:8080
```

如果不在仓库根目录，可直接执行脚本：

```bash
/path/to/cloud-native-mcp-server/scripts/sse_smoke_test.sh http://127.0.0.1:8080
```

---

## 认证验证

当启用 `MCP_AUTH_ENABLED=true` 且 `MCP_AUTH_MODE=apikey` 时：

```bash
# 通过 query 参数传递 API Key
curl -sS -N "http://127.0.0.1:8080/api/aggregate/sse?api_key=ChangeMe-Strong-Key-123!"
```

也可以通过请求头传递：

```bash
curl -sS -N \
  -H "X-Api-Key: ChangeMe-Strong-Key-123!" \
  http://127.0.0.1:8080/api/aggregate/sse
```

---

## 常用启动参数

```bash
# 服务模式与监听地址
export MCP_MODE=sse
export MCP_ADDR=0.0.0.0:8080

# 认证（apikey 模式）
export MCP_AUTH_ENABLED=true
export MCP_AUTH_MODE=apikey
export MCP_AUTH_API_KEY='ChangeMe-Strong-Key-123!'

# 可选：禁用暂不需要的服务
export MCP_DISABLED_SERVICES="kibana,jaeger"
```

---

## 下一步

- [快速开始 FAQ]({{< relref "faq.md" >}})
- [故障排除]({{< relref "troubleshooting.md" >}})
- [安全指南]({{< relref "/docs/security.md" >}})
- [配置指南]({{< relref "/docs/configuration.md" >}})
- [性能指南]({{< relref "/docs/performance.md" >}})
- [工具参考]({{< relref "/docs/tools.md" >}})

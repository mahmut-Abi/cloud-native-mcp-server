---
title: "故障排除"
weight: 30
description: "启动、认证、传输链路与服务集成异常的排查清单。"
---

# 故障排除

当安装或运行出现异常时，可按本页步骤快速定位问题。

## 1. 服务无法启动

先排查端口占用和启动日志：

```bash
# 查看 8080 端口占用
ss -lntp | rg 8080

# 以 debug 日志启动
./cloud-native-mcp-server --mode=sse --addr=127.0.0.1:8080 --log-level=debug
```

如果使用 Docker：

```bash
docker logs --tail=200 cloud-native-mcp-server
```

## 2. `/health` 无法访问或非 200

```bash
curl -sv http://127.0.0.1:8080/health
```

若连接失败，请确认：

- 进程/容器是否正常运行
- 监听地址与端口映射是否正确
- 防火墙或安全组策略是否放行

## 3. SSE 能连接但握手异常

```bash
curl -svN --connect-timeout 5 --max-time 15 \
  -H "Accept: text/event-stream" \
  "http://127.0.0.1:8080/api/aggregate/sse"
```

再执行内置自检：

```bash
make sse-smoke BASE_URL=http://127.0.0.1:8080
```

如果失败，重点查看首个事件和 message 端点返回。

## 4. 认证持续返回 401

先确认运行时认证配置：

```bash
export MCP_AUTH_ENABLED=true
export MCP_AUTH_MODE=apikey
export MCP_AUTH_API_KEY='ChangeMe-Strong-Key-123!'
```

然后分别验证两种 API Key 传递方式：

```bash
curl -sS -N -H "X-Api-Key: ChangeMe-Strong-Key-123!" \
  http://127.0.0.1:8080/api/aggregate/sse

curl -sS -N "http://127.0.0.1:8080/api/aggregate/sse?api_key=ChangeMe-Strong-Key-123!"
```

## 5. 工具调用慢或超时

优先检查以下项：

- 上游服务（Kubernetes API、Prometheus、Grafana 等）是否可达
- 单次请求范围是否过大（namespace、时间区间、对象数量）
- 是否优先使用分页/摘要能力
- 超时与限流参数是否符合实际流量

## 6. 部分服务不可用

查看当前可用服务与工具：

```bash
./cloud-native-mcp-server --list=services --output=table
./cloud-native-mcp-server --list=tools --service=kubernetes --output=table
```

如果服务缺失，重点检查：

- `MCP_ENABLED_SERVICES`
- `MCP_DISABLED_SERVICES`
- 对应服务地址、认证、TLS 等配置

## 7. 提交 Issue 前建议收集的信息

- 运行模式（`sse` / `streamable-http` / `http` / `stdio`）
- 启动命令与关键环境变量（注意脱敏）
- `/health` 与传输端点的 `curl` 输出
- 关键日志（最近 100-200 行）

## 相关页面

- [快速开始]({{< relref "_index.md" >}})
- [快速开始 FAQ]({{< relref "faq.md" >}})
- [配置指南]({{< relref "/docs/configuration.md" >}})
- [安全指南]({{< relref "/docs/security.md" >}})

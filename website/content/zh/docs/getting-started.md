---
title: "快速开始"
weight: 10
description: "快速进入安装、FAQ 和故障排除入口。"
---

# 快速开始

该页面位于 `docs/` 路径下，作为快速开始入口页使用。

如果你希望按完整流程阅读，请优先进入：

- [快速开始主指南]({{< relref "/getting-started/_index.md" >}})
- [快速开始 FAQ]({{< relref "/getting-started/faq.md" >}})
- [故障排除手册]({{< relref "/getting-started/troubleshooting.md" >}})

---

## 30 秒可用性检查

```bash
# 健康检查
curl -sS http://127.0.0.1:8080/health

# SSE 握手 + initialize 全链路检查（仓库根目录执行）
make sse-smoke BASE_URL=http://127.0.0.1:8080
```

---

## 常用入口

- SSE 聚合入口（`--mode=sse`）: `http://127.0.0.1:8080/api/aggregate/sse`
- Streamable HTTP 聚合入口（`--mode=streamable-http`）: `http://127.0.0.1:8080/api/aggregate/streamable-http`
- 健康检查: `http://127.0.0.1:8080/health`

---

## 推荐阅读顺序

1. [快速开始主指南]({{< relref "/getting-started/_index.md" >}})
2. [安全指南]({{< relref "security.md" >}})
3. [配置指南]({{< relref "configuration.md" >}})
4. [部署指南]({{< relref "deployment.md" >}})
5. [性能指南]({{< relref "performance.md" >}})

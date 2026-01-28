---
title: "Helm 服务"
weight: 2
---

# Helm 服务

Helm 服务提供全面的包管理和部署功能，包含 31 个工具来管理 Helm 图表、发布和仓库。

## 概述

Cloud Native MCP Server 中的 Helm 服务使 AI 助手能够高效地管理 Helm 图表和发布。它提供用于图表安装、升级、回滚和仓库管理的工具。

### 主要功能

{{< columns >}}
### 📦 图表管理
对 Helm 图表进行完全控制，包括安装、升级和卸载发布。
<--->

### 🗄️ 仓库管理
使用工具管理 Helm 图表仓库，包括添加、更新和搜索图表。
{{< /columns >}}

{{< columns >}}
### 🔄 发布管理
使用回滚、历史记录和状态检查功能处理 Helm 发布。
<--->

### ⚙️ 配置
有效管理图表值、配置和依赖关系。
{{< /columns >}}

---

## 可用工具 (31)

### 图表管理
- **helm-list-releases**: 列出所有命名空间中的所有发布
- **helm-install-chart**: 安装图表
- **helm-upgrade-release**: 升级发布
- **helm-uninstall-release**: 卸载发布
- **helm-get-release**: 获取发布信息
- **helm-rollback-release**: 回滚发布
- **helm-get-history**: 获取发布历史
- **helm-search-repo**: 在仓库中搜索图表
- **helm-add-repo**: 添加图表仓库
- **helm-update-repo**: 更新图表仓库
- **helm-repo-list**: 列出图表仓库
- **helm-get-values**: 获取发布的值
- **helm-template**: 在本地生成图表模板
- **helm-package**: 将图表目录打包成图表归档
- **helm-pull**: 从仓库下载图表
- **helm-push**: 将图表推送到注册表

### 图表信息
- **helm-get-chart**: 获取图表信息
- **helm-create**: 创建新图表
- **helm-dependency-build**: 构建图表依赖关系
- **helm-dependency-update**: 更新图表依赖关系
- **helm-lint**: 检查图表可能存在的问题
- **helm-test**: 为发布运行测试
- **helm-status**: 显示发布的状态
- **helm-history**: 显示发布的历史记录
- **helm-get-manifest**: 显示发布的清单
- **helm-get-notes**: 显示发布的注释
- **helm-get-hooks**: 显示发布的钩子
- **helm-get-all**: 获取发布的所有资源
- **helm-verify**: 验证图表的来源
- **helm-show-chart**: 显示图表信息
- **helm-show-readme**: 显示图表的 README

---

## 快速示例

### 安装图表

```json
{
  "method": "tools/call",
  "params": {
    "name": "helm-install-chart",
    "arguments": {
      "chart": "nginx-ingress",
      "repo": "https://kubernetes.github.io/ingress-nginx",
      "release": "my-nginx",
      "namespace": "ingress-nginx"
    }
  }
}
```

### 升级发布

```json
{
  "method": "tools/call",
  "params": {
    "name": "helm-upgrade-release",
    "arguments": {
      "release": "my-nginx",
      "chart": "nginx-ingress",
      "repo": "https://kubernetes.github.io/ingress-nginx",
      "set_values": {
        "controller.replicaCount": 3
      }
    }
  }
}
```

### 列出所有发布

```json
{
  "method": "tools/call",
  "params": {
    "name": "helm-list-releases",
    "arguments": {}
  }
}
```

---

## 最佳实践

- 安装 Helm 图表时始终指定命名空间
- 为复杂配置使用值文件
- 定期更新图表仓库
- 监控发布历史以备回滚功能
- 在安装前使用检查工具验证图表

## 下一步

- [Kubernetes 服务](/zh/services/kubernetes/) 了解核心编排
- [配置指南](/zh/guides/configuration/) 了解详细设置
- [部署最佳实践](/zh/guides/deployment/) 了解生产部署
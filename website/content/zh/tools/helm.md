---
title: "Helm 工具"
weight: 20
---

# Helm 工具

Cloud Native MCP Server 提供 31 个 Helm 管理工具，用于 Chart 仓库、Release 管理和版本控制。

## Chart 管理

### list_repositories

列出所有已添加的 Helm 仓库。

**返回**: 仓库列表

### add_repository

添加新的 Helm 仓库。

**参数**:
- `name` (string, required) - 仓库名称
- `url` (string, required) - 仓库 URL
- `username` (string, optional) - 用户名（私有仓库）
- `password` (string, optional) - 密码（私有仓库）

**返回**: 添加结果

### remove_repository

移除指定的 Helm 仓库。

**参数**:
- `name` (string, required) - 仓库名称

**返回**: 移除结果

### update_repository

更新指定的 Helm 仓库。

**参数**:
- `name` (string, required) - 仓库名称

**返回**: 更新结果

### search_chart

在仓库中搜索 Chart。

**参数**:
- `keyword` (string, required) - 搜索关键词
- `repo` (string, optional) - 仓库名称（可选，不指定则搜索所有仓库）
- `version` (string, optional) - 版本约束

**返回**: 搜索结果

### show_chart

显示 Chart 的详细信息。

**参数**:
- `chart` (string, required) - Chart 名称（格式：repo/chart）
- `version` (string, optional) - Chart 版本

**返回**: Chart 详细信息

### pull_chart

下载 Chart 到本地。

**参数**:
- `chart` (string, required) - Chart 名称（格式：repo/chart）
- `version` (string, optional) - Chart 版本
- `destination` (string, optional) - 下载目录

**返回**: 下载结果

## Release 管理

### list_releases

列出所有 Release。

**参数**:
- `namespace` (string, optional) - 命名空间
- `all` (bool, optional) - 列出所有命名空间的 Release
- `filter` (string, optional) - 过滤条件

**返回**: Release 列表

### get_release

获取单个 Release 的详细信息。

**参数**:
- `name` (string, required) - Release 名称
- `namespace` (string, required) - 命名空间

**返回**: Release 详细信息

### install_chart

安装新的 Release。

**参数**:
- `chart` (string, required) - Chart 名称或路径
- `release` (string, required) - Release 名称
- `namespace` (string, required) - 命名空间
- `repo` (string, optional) - 仓库 URL
- `values` (object, optional) - 配置值
- `version` (string, optional) - Chart 版本
- `createNamespace` (bool, optional) - 创建命名空间

**返回**: 安装结果

### upgrade_release

升级现有的 Release。

**参数**:
- `release` (string, required) - Release 名称
- `namespace` (string, required) - 命名空间
- `chart` (string, required) - 新 Chart 名称或路径
- `values` (object, optional) - 配置值
- `version` (string, optional) - Chart 版本

**返回**: 升级结果

### rollback_release

回滚 Release 到之前的版本。

**参数**:
- `release` (string, required) - Release 名称
- `namespace` (string, required) - 命名空间
- `revision` (int, optional) - 回滚到的版本号

**返回**: 回滚结果

### uninstall_release

卸载指定的 Release。

**参数**:
- `release` (string, required) - Release 名称
- `namespace` (string, required) - 命名空间
- `keepHistory` (bool, optional) - 保留历史记录

**返回**: 卸载结果

### get_release_history

获取 Release 的历史记录。

**参数**:
- `release` (string, required) - Release 名称
- `namespace` (string, required) - 命名空间

**返回**: 历史记录列表

### get_release_status

获取 Release 的当前状态。

**参数**:
- `release` (string, required) - Release 名称
- `namespace` (string, required) - 命名空间

**返回**: Release 状态

### get_release_values

获取 Release 的配置值。

**参数**:
- `release` (string, required) - Release 名称
- `namespace` (string, required) - 命名空间
- `all` (bool, optional) - 显示所有值（包括默认值）

**返回**: 配置值

## Values 管理

### get_values

获取 Release 的配置值（与 `get_release_values` 功能相同）。

**参数**:
- `release` (string, required) - Release 名称
- `namespace` (string, required) - 命名空间

**返回**: 配置值

### set_values

设置 Release 的配置值。

**参数**:
- `release` (string, required) - Release 名称
- `namespace` (string, required) - 命名空间
- `values` (object, required) - 配置值
- `reset` (bool, optional) - 重置为默认值

**返回**: 设置结果

### diff_values

比较两个配置值的差异。

**参数**:
- `values1` (object, required) - 第一个配置值
- `values2` (object, required) - 第二个配置值

**返回**: 差异列表

## Release 操作

### test_release

测试 Release 的配置。

**参数**:
- `release` (string, required) - Release 名称
- `namespace` (string, required) - 命名空间

**返回**: 测试结果

### lint_chart

检查 Chart 的配置。

**参数**:
- `chart` (string, required) - Chart 名称或路径
- `values` (object, optional) - 配置值

**返回**: 检查结果

### package_chart

打包 Chart。

**参数**:
- `chart` (string, required) - Chart 目录路径
- `destination` (string, optional) - 输出目录
- `version` (string, optional) - 版本
- `appVersion` (string, optional) - 应用版本

**返回**: 打包结果

### verify_chart

验证 Chart 的签名。

**参数**:
- `chart` (string, required) - Chart 文件路径
- `key` (string, optional) - 公钥文件路径

**返回**: 验证结果

### template_chart

生成 Chart 的模板。

**参数**:
- `chart` (string, required) - Chart 名称或路径
- `release` (string, required) - Release 名称
- `namespace` (string, required) - 命名空间
- `values` (object, optional) - 配置值

**返回**: 生成的模板

## Chart 依赖

### list_dependencies

列出 Chart 的依赖。

**参数**:
- `chart` (string, required) - Chart 名称或路径

**返回**: 依赖列表

### update_dependencies

更新 Chart 的依赖。

**参数**:
- `chart` (string, required) - Chart 目录路径
- `skipDeps` (bool, optional) - 跳过依赖更新

**返回**: 更新结果

## 插件管理

### list_plugins

列出所有已安装的 Helm 插件。

**返回**: 插件列表

### install_plugin

安装 Helm 插件。

**参数**:
- `url` (string, required) - 插件 URL
- `version` (string, optional) - 插件版本

**返回**: 安装结果

## 版本管理

### list_versions

列出 Chart 的所有可用版本。

**参数**:
- `chart` (string, required) - Chart 名称（格式：repo/chart）

**返回**: 版本列表

### get_version_info

获取指定版本的详细信息。

**参数**:
- `chart` (string, required) - Chart 名称（格式：repo/chart）
- `version` (string, required) - 版本号

**返回**: 版本详细信息

## 调试工具

### debug_release

调试 Release 的配置。

**参数**:
- `release` (string, required) - Release 名称
- `namespace` (string, required) - 命名空间

**返回**: 调试信息

## 配置

Helm 工具通过以下配置进行初始化：

```yaml
helm:
  enabled: false

  # kubeconfig 路径
  kubeconfigPath: ""

  # 默认命名空间
  namespace: "default"

  # 调试模式
  debug: false

  # 仓库更新超时（秒）
  timeoutSec: 300

  # 最大重试次数
  maxRetries: 3

  # 启用镜像
  useMirrors: false

  # 自定义镜像映射
  mirrors: {}
```

## 最佳实践

1. **使用版本固定**: 在生产环境中始终指定 Chart 版本
2. **配置管理**: 使用 `values.yaml` 文件管理配置
3. **依赖检查**: 使用 `list_dependencies` 检查依赖关系
4. **测试验证**: 使用 `test_release` 和 `lint_chart` 验证配置
5. **版本控制**: 使用 `get_release_history` 跟踪变更历史

## 相关文档

- [配置指南](/zh/guides/configuration/)
- [部署指南](/zh/guides/deployment/)
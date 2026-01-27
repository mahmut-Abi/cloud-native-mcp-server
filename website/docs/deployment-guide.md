# GitHub Pages 部署指南

本指南说明如何将 Cloud Native MCP Server 的文档网站部署到 GitHub Pages。

## 前提条件

1. 已有 GitHub 账户
2. 已创建 GitHub 仓库（mahmut-Abi/cloud-native-mcp-server）
3. 已完成网站内容的创建和本地测试

## 自动部署

### 1. 启用 GitHub Pages

在 GitHub 仓库页面：
1. 进入 **Settings** → **Pages**
2. 在 **Build and deployment** 部分，选择 **Source** 为 **GitHub Actions**
3. 保存设置

### 2. 配置工作流

GitHub Actions 工作流文件已创建在 `.github/workflows/gh-pages.yml`。

该工作流会：
- 每次推送到 main 分支时自动触发
- 使用最新的 Hugo 版本构建网站
- 将构建结果部署到 GitHub Pages

### 3. 部署网站

将代码推送到 main 分支：

```bash
git add .
git commit -m "feat: add Hugo website for documentation"
git push origin main
```

GitHub Actions 会自动开始构建和部署过程。

### 4. 查看部署状态

1. 进入 GitHub 仓库的 **Actions** 标签
2. 查看工作流运行状态
3. 部署完成后，网站将在以下地址可访问：
   - `https://mahmut-abi.github.io/cloud-native-mcp-server/`

## 手动部署

如果需要手动触发部署：

1. 进入 GitHub 仓库的 **Actions** 标签
2. 选择 **Deploy Hugo to GitHub Pages** 工作流
3. 点击 **Run workflow** 按钮
4. 选择分支（通常是 main）
5. 点击 **Run workflow** 按钮确认

## 本地测试

在部署前，建议在本地测试网站：

```bash
cd website
hugo server -D
```

访问 http://localhost:1313 查看网站。

## 配置说明

### baseURL 配置

在 `website/hugo.toml` 中，baseURL 已设置为：

```toml
baseURL = 'https://mahmut-abi.github.io/cloud-native-mcp-server/'
```

如果你的仓库名不同，需要相应修改。

### 自定义域名（可选）

如果需要使用自定义域名：

1. 在 `website/hugo.toml` 中修改 baseURL：
   ```toml
   baseURL = 'https://your-custom-domain.com/'
   ```

2. 在 GitHub Pages 设置中添加自定义域名

3. 配置 DNS 记录

## 更新网站

更新网站内容：

1. 修改 `website/content/` 目录中的 Markdown 文件
2. 本地测试：`cd website && hugo server -D`
3. 提交并推送：
   ```bash
   git add website/
   git commit -m "docs: update website content"
   git push origin main
   ```

GitHub Actions 会自动重新部署。

## 故障排查

### 构建失败

查看 GitHub Actions 日志：
1. 进入 **Actions** 标签
2. 点击失败的工作流运行
3. 查看错误日志

常见问题：
- 主题子模块未正确拉取
- Hugo 配置文件有语法错误
- 内容文件有格式错误

### 网站样式不正确

确保：
1. `static/css/custom.css` 文件存在
2. 主题子模块已正确克隆
3. 清除浏览器缓存

### 404 错误

检查：
1. baseURL 配置是否正确
2. 文件路径是否正确
3. 是否使用了正确的相对路径

## 维护建议

1. **定期更新主题**：保持 Ananke 主题最新
   ```bash
   cd themes/ananke
   git pull origin main
   ```

2. **测试本地构建**：部署前在本地测试
   ```bash
   hugo --minify
   ```

3. **监控部署**：定期检查 GitHub Actions 运行状态

4. **备份内容**：定期备份网站内容

## 相关资源

- [Hugo 官方文档](https://gohugo.io/documentation/)
- [GitHub Pages 文档](https://docs.github.com/en/pages)
- [Ananke 主题文档](https://github.com/theNewDynamic/gohugo-theme-ananke)
- [GitHub Actions 文档](https://docs.github.com/en/actions)
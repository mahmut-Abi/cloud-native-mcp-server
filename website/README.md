# Cloud Native MCP Server Website

这是 Cloud Native MCP Server 项目的官方文档网站，使用 Hugo 框架构建，部署在 GitHub Pages 上。

## 本地开发

### 安装 Hugo

```bash
# macOS
brew install hugo

# Linux (Ubuntu/Debian)
sudo apt-get install hugo

# Windows (Chocolatey)
choco install hugo-extended
```

### 运行开发服务器

```bash
cd website
hugo server -D
```

访问 http://localhost:1313 查看网站。

### 构建生产版本

```bash
cd website
hugo --minify
```

构建后的文件将输出到 `public/` 目录。

## 项目结构

```
website/
├── content/           # 网站内容
│   ├── _index.md     # 首页
│   ├── docs/         # 文档
│   ├── features/     # 特性介绍
│   └── services/     # 服务介绍
├── layouts/          # 自定义布局
│   └── partials/     # 局部模板
├── static/           # 静态资源
│   └── css/          # 自定义样式
├── themes/           # Hugo 主题
│   └── ananke/       # Ananke 主题（子模块）
├── hugo.toml         # Hugo 配置文件
└── README.md         # 本文件
```

## 添加新内容

### 创建新页面

```bash
hugo new content/docs/new-page.md
```

### 编辑内容

编辑 Markdown 文件，使用 Hugo 的 shortcode 和 front matter。

## 部署

网站通过 GitHub Actions 自动部署到 GitHub Pages。

### 手动部署

1. 推送到 main 分支会自动触发部署
2. 或者在 GitHub Actions 页面手动触发工作流

### 部署配置

部署配置在 `.github/workflows/gh-pages.yml` 文件中。

## 自定义

### 修改样式

编辑 `static/css/custom.css` 文件。

### 修改配置

编辑 `hugo.toml` 文件。

### 自定义布局

在 `layouts/` 目录中创建或修改模板文件。

## 主题

本项目使用 [Ananke](https://github.com/theNewDynamic/gohugo-theme-ananke) 主题。

主题作为 Git 子模块管理：

```bash
# 更新主题
cd themes/ananke
git pull origin main
```

## 贡献

欢迎提交 Issue 和 Pull Request 来改进网站内容和样式。

## 许可证

与主项目保持一致，MIT License。
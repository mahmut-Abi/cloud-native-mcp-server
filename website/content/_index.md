---
title: "Cloud Native MCP Server"
---

<style>
.hero {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 4rem 2rem;
  text-align: center;
  border-radius: 12px;
  margin-bottom: 3rem;
  box-shadow: 0 10px 40px rgba(102, 126, 234, 0.3);
}
.hero h1 {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 1rem;
  text-shadow: 0 2px 4px rgba(0,0,0,0.2);
}
.hero p {
  font-size: 1.1rem;
  opacity: 0.95;
  max-width: 800px;
  margin: 0 auto 2rem;
  line-height: 1.6;
}
.cta-button {
  display: inline-block;
  background: white;
  color: #667eea;
  padding: 0.875rem 2rem;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 600;
  transition: all 0.3s ease;
  box-shadow: 0 4px 6px rgba(0,0,0,0.1);
}
.cta-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(0,0,0,0.2);
}
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 2rem;
  margin: 3rem 0;
}
.stat-item {
  text-align: center;
  padding: 1.5rem;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0,0,0,0.1);
  transition: transform 0.3s ease;
}
.stat-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0,0,0,0.15);
}
.stat-number {
  font-size: 2.5rem;
  font-weight: 700;
  color: #667eea;
  margin-bottom: 0.5rem;
}
.stat-label {
  font-size: 1rem;
  color: #2d3748;
  font-weight: 500;
}
.feature-card {
  background: white;
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 4px 6px rgba(0,0,0,0.1);
  transition: all 0.3s ease;
  height: 100%;
}
.feature-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 24px rgba(0,0,0,0.15);
}
.feature-icon {
  font-size: 2.5rem;
  margin-bottom: 1rem;
}
pre {
  background: #2d3748;
  color: #e2e8f0;
  padding: 1.5rem;
  border-radius: 8px;
  overflow-x: auto;
  margin: 1.5rem 0;
  border: 1px solid rgba(255,255,255,0.1);
}
pre code {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 0.9rem;
  line-height: 1.6;
}
h1, h2, h3 {
  color: #2d3748;
  margin-top: 2rem;
  margin-bottom: 1rem;
}
h1 {
  font-size: 2.5rem;
  font-weight: 700;
}
h2 {
  font-size: 2rem;
  font-weight: 600;
  border-bottom: 2px solid #667eea;
  padding-bottom: 0.5rem;
}
h3 {
  font-size: 1.5rem;
  font-weight: 600;
}
p {
  line-height: 1.8;
  margin-bottom: 1rem;
}
a {
  color: #667eea;
  text-decoration: none;
  transition: color 0.3s ease;
}
a:hover {
  color: #764ba2;
  text-decoration: underline;
}
</style>

<div class="hero">
  <h1>Cloud Native MCP Server</h1>
  <p>高性能 Kubernetes 和云原生基础设施管理 MCP 服务器，集成 10 个服务和 220+ 工具，让 AI 助手轻松管理您的云原生基础设施</p>
  <a href="https://github.com/mahmut-Abi/cloud-native-mcp-server" class="cta-button">查看 GitHub 仓库</a>
</div>

<div class="stats-grid">
  <div class="stat-item">
    <div class="stat-number">10</div>
    <div class="stat-label">集成服务</div>
  </div>
  <div class="stat-item">
    <div class="stat-number">220+</div>
    <div class="stat-label">MCP 工具</div>
  </div>
  <div class="stat-item">
    <div class="stat-number">3</div>
    <div class="stat-label">运行模式</div>
  </div>
  <div class="stat-item">
    <div class="stat-number">100%</div>
    <div class="stat-label">开源免费</div>
  </div>
</div>

## 快速开始

```bash
docker run -d \
  --name cloud-native-mcp-server \
  -p 8080:8080 \
  -v ~/.kube:/root/.kube:ro \
  mahmutabi/cloud-native-mcp-server:latest
```

## 核心特性

<div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 1.5rem; margin-top: 2rem;">

<div class="feature-card">
  <div class="feature-icon">🚀</div>
  <h3>高性能</h3>
  <p>LRU 缓存、JSON 编码池、智能响应限制，确保最佳性能</p>
</div>

<div class="feature-card">
  <div class="feature-icon">🔒</div>
  <h3>安全可靠</h3>
  <p>API Key、Bearer Token、Basic Auth 多种认证方式，安全的密钥管理</p>
</div>

<div class="feature-card">
  <div class="feature-icon">📊</div>
  <h3>全面监控</h3>
  <p>集成 Prometheus、Grafana、Jaeger 等监控和追踪工具</p>
</div>

<div class="feature-card">
  <div class="feature-icon">🔧</div>
  <h3>灵活配置</h3>
  <p>支持 SSE、HTTP、stdio 多种模式，适配各种使用场景</p>
</div>

</div>

## 了解更多

- [查看所有服务](/services/) - 了解 10 个集成服务的详细信息
- [完整工具参考](/docs/tools/) - 所有 220+ 工具的详细文档
- [部署指南](/docs/deployment/) - 部署策略和最佳实践
- [配置指南](/docs/configuration/) - 配置选项和示例
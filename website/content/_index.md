---
title: "Cloud Native MCP Server"
---

<div class="hero">
  <h1>Cloud Native MCP Server</h1>
  <p>高性能 Kubernetes 和云原生基础设施管理 MCP 服务器，集成 10 个服务和 220+ 工具，让 AI 助手轻松管理您的云原生基础设施</p>
  <a href="https://github.com/mahmut-Abi/cloud-native-mcp-server" class="cta-button">查看 GitHub 仓库</a>
</div>

<div class="container">
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

  <h2 class="section-title">核心特性</h2>
  <div class="features-grid">
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

  <div class="code-section">
    <h2>快速开始</h2>
    <pre>docker run -d \
  --name cloud-native-mcp-server \
  -p 8080:8080 \
  -v ~/.kube:/root/.kube:ro \
  mahmutabi/cloud-native-mcp-server:latest</pre>
  </div>

  <h2 class="section-title">了解更多</h2>
  <div class="links-section">
    <ul>
      <li><a href="/services/">查看所有服务</a> - 了解 10 个集成服务的详细信息</li>
      <li><a href="/docs/tools/">完整工具参考</a> - 所有 220+ 工具的详细文档</li>
      <li><a href="/docs/deployment/">部署指南</a> - 部署策略和最佳实践</li>
      <li><a href="/docs/configuration/">配置指南</a> - 配置选项和示例</li>
    </ul>
  </div>
</div>

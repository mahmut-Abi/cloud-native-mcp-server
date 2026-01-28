---
title: "Alertmanager 工具"
weight: 70
---

# Alertmanager 工具

Cloud Native MCP Server 提供 15 个 Alertmanager 管理工具，用于告警管理、沉默规则、通知路由和配置管理。

## 告警管理

### list_alerts

列出所有活动告警。

**参数**:
- `silenced` (bool, optional) - 是否包含沉默的告警
- `inhibited` (bool, optional) - 是否包含抑制的告警
- `filter` (string, optional) - 告警过滤器

**返回**: 告警列表

### get_alert

获取单个告警的详细信息。

**参数**:
- `fingerprint` (string, required) - 告警指纹

**返回**: 告警详细信息

### get_alert_groups

获取告警组。

**参数**:
- `filter` (string, optional) - 告警组过滤器

**返回**: 告警组列表

## 沉默规则管理

### get_silences

获取所有沉默规则。

**参数**:
- `filter` (string, optional) - 沉默规则过滤器

**返回**: 沉默规则列表

### create_silence

创建新的沉默规则。

**参数**:
- `matchers` (array, required) - 匹配器数组
- `startsAt` (string, required) - 开始时间（RFC3339）
- `endsAt` (string, required) - 结束时间（RFC3339）
- `createdBy` (string, required) - 创建者
- `comment` (string, required) - 注释

**返回**: 创建的沉默规则信息

**示例**:
```json
{
  "name": "create_silence",
  "arguments": {
    "matchers": [
      {
        "name": "alertname",
        "value": "HighCPUUsage",
        "isRegex": false
      }
    ],
    "startsAt": "2024-01-01T00:00:00Z",
    "endsAt": "2024-01-01T06:00:00Z",
    "createdBy": "admin",
    "comment": "Maintenance window"
  }
}
```

### delete_silence

删除沉默规则。

**参数**:
- `silenceId` (string, required) - 沉默规则 ID

**返回**: 删除结果

### expire_silence

使沉默规则过期。

**参数**:
- `silenceId` (string, required) - 沉默规则 ID

**返回**: 过期结果

## 规则管理

### get_alert_rules

获取告警规则。

**参数**:
- `filter` (string, optional) - 规则过滤器

**返回**: 告警规则列表

### list_rule_groups

列出规则组。

**返回**: 规则组列表

## 配置管理

### get_config

获取 Alertmanager 配置。

**返回**: 配置信息，包括：
- 全局配置
- 路由配置
- 接收者配置
- 抑制规则

### get_status

获取 Alertmanager 状态。

**返回**: 状态信息，包括：
- 集群状态
- 配置版本
- 启动时间

## 通知管理

### list_notifications

列出通知历史。

**参数**:
- `receiver` (string, optional) - 接收者名称
- `limit` (int, optional) - 限制数量

**返回**: 通知列表

### get_receivers

获取所有接收者配置。

**返回**: 接收者列表

### list_routes

列出所有路由。

**返回**: 路由配置列表

## 健康检查

### get_health

获取 Alertmanager 健康状态。

**返回**: 健康状态

## 配置

Alertmanager 工具通过以下配置进行初始化：

```yaml
alertmanager:
  enabled: false

  # Alertmanager 服务器地址
  address: "http://localhost:9093"

  # 请求超时（秒）
  timeoutSec: 30

  # Basic Auth
  username: ""
  password: ""

  # Bearer Token
  bearerToken: ""

  # TLS 配置
  tlsSkipVerify: false
  tlsCertFile: ""
  tlsKeyFile: ""
  tlsCAFile: ""
```

## 告警生命周期

告警在 Alertmanager 中的生命周期：

1. **接收**: Prometheus 发送告警到 Alertmanager
2. **去重**: 根据标签去重
3. **分组**: 按照路由规则分组
4. **抑制**: 检查抑制规则
5. **沉默**: 检查沉默规则
6. **路由**: 路由到对应的接收者
7. **通知**: 发送通知
8. **静默**: 告警恢复后静默一段时间

## 沉默规则使用场景

### 维护窗口

```json
{
  "matchers": [
    {
      "name": "severity",
      "value": "critical",
      "isRegex": false
    }
  ],
  "startsAt": "2024-01-01T02:00:00Z",
  "endsAt": "2024-01-01T04:00:00Z",
  "createdBy": "admin",
  "comment": "Scheduled maintenance"
}
```

### 测试环境

```json
{
  "matchers": [
    {
      "name": "env",
      "value": "test",
      "isRegex": false
    }
  ],
  "startsAt": "2024-01-01T00:00:00Z",
  "endsAt": "2024-12-31T23:59:59Z",
  "createdBy": "admin",
  "comment": "Test environment"
}
```

### 特定告警

```json
{
  "matchers": [
    {
      "name": "alertname",
      "value": "DiskSpaceLow",
      "isRegex": false
    },
    {
      "name": "instance",
      "value": "server1",
      "isRegex": false
    }
  ],
  "startsAt": "2024-01-01T00:00:00Z",
  "endsAt": "2024-01-02T00:00:00Z",
  "createdBy": "admin",
  "comment": "Known issue, being addressed"
}
```

## 最佳实践

1. **告警设计**:
   - 使用清晰的告警名称和描述
   - 设置合理的严重级别
   - 添加相关的标签（severity, team, service）

2. **沉默规则管理**:
   - 为维护窗口创建临时沉默规则
   - 为测试环境创建永久沉默规则
   - 定期清理过期的沉默规则

3. **路由配置**:
   - 按服务和团队组织路由
   - 设置合理的告警分组
   - 配置告警抑制避免告警风暴

4. **通知管理**:
   - 使用多种通知渠道（邮件、Slack、Webhook）
   - 设置合理的告警间隔
   - 配置告警升级策略

5. **监控和维护**:
   - 定期检查 `get_health` 确保服务正常
   - 使用 `get_status` 监控配置变更
   - 定期审查和优化告警规则

## 相关文档

- [配置指南](/zh/guides/configuration/)
- [部署指南](/zh/guides/deployment/)
- [Prometheus 工具](/zh/tools/prometheus/)
- [Alertmanager 文档](https://prometheus.io/docs/alerting/latest/alertmanager/)
---
title: "Utilities 工具"
weight: 100
---

# Utilities 工具

Cloud Native MCP Server 提供 6 个通用实用工具，用于数据转换、编码和标识生成。

## Base64 编码

### base64_encode

对字符串进行 Base64 编码。

**参数**:
- `text` (string, required) - 要编码的文本

**返回**: Base64 编码的字符串

**示例**:
```json
{
  "name": "base64_encode",
  "arguments": {
    "text": "Hello, World!"
  }
}
```

**返回**: `"SGVsbG8sIFdvcmxkIQ=="`

### base64_decode

对 Base64 编码的字符串进行解码。

**参数**:
- `encoded` (string, required) - Base64 编码的字符串

**返回**: 解码后的文本

**示例**:
```json
{
  "name": "base64_decode",
  "arguments": {
    "encoded": "SGVsbG8sIFdvcmxkIQ=="
  }
}
```

**返回**: `"Hello, World!"`

## JSON 处理

### json_parse

解析 JSON 字符串。

**参数**:
- `json` (string, required) - JSON 字符串

**返回**: 解析后的对象

**示例**:
```json
{
  "name": "json_parse",
  "arguments": {
    "json": "{\"name\":\"John\",\"age\":30}"
  }
}
```

**返回**:
```json
{
  "name": "John",
  "age": 30
}
```

### json_stringify

将对象序列化为 JSON 字符串。

**参数**:
- `object` (object, required) - 要序列化的对象
- `pretty` (boolean, optional) - 是否美化输出（默认：false）

**返回**: JSON 字符串

**示例**:
```json
{
  "name": "json_stringify",
  "arguments": {
    "object": {
      "name": "John",
      "age": 30
    },
    "pretty": true
  }
}
```

**返回**:
```json
{
  "name": "John",
  "age": 30
}
```

## 时间戳

### timestamp

获取当前时间戳。

**参数**:
- `format` (string, optional) - 返回格式（unix, iso8601, rfc3339，默认：unix）

**返回**: 时间戳

**示例**:
```json
{
  "name": "timestamp",
  "arguments": {
    "format": "iso8601"
  }
}
```

**返回**: `"2024-01-01T12:00:00Z"`

**支持的格式**:
- `unix` - Unix 时间戳（秒）
- `iso8601` - ISO 8601 格式
- `rfc3339` - RFC 3339 格式

## UUID 生成

### uuid

生成 UUID（通用唯一标识符）。

**参数**:
- `version` (string, optional) - UUID 版本（v4，默认：v4）

**返回**: UUID 字符串

**示例**:
```json
{
  "name": "uuid",
  "arguments": {}
}
```

**返回**: `"f47ac10b-58cc-4372-a567-0e02b2c3d479"`

**支持的版本**:
- `v4` - 随机 UUID（推荐）

## 使用场景

### 数据编码

在 Kubernetes Secret 中存储配置：

```json
{
  "name": "base64_encode",
  "arguments": {
    "text": "my-secret-password"
  }
}
```

### 数据转换

在服务和工具之间传递数据：

```json
{
  "name": "json_stringify",
  "arguments": {
    "object": {
      "key1": "value1",
      "key2": "value2"
    }
  }
}
```

### 数据解析

解析 API 响应：

```json
{
  "name": "json_parse",
  "arguments": {
    "json": "{\"status\":\"success\",\"data\":{...}}"
  }
}
```

### 时间戳生成

生成唯一标识符：

```json
{
  "name": "timestamp",
  "arguments": {
    "format": "rfc3339"
  }
}
```

### 唯一标识

生成请求 ID 或资源 ID：

```json
{
  "name": "uuid",
  "arguments": {}
}
```

## 实用示例

### 创建 Kubernetes Secret

使用 Base64 编码创建 Secret：

```json
{
  "name": "base64_encode",
  "arguments": {
    "text": "YWRtaW4="
  }
}
```

### 解析配置文件

读取和解析 JSON 配置：

```json
{
  "name": "json_parse",
  "arguments": {
    "json": "{\"service\":\"api\",\"port\":8080}"
  }
}
```

### 生成唯一 ID

为日志条目生成唯一 ID：

```json
{
  "name": "uuid",
  "arguments": {}
}
```

组合使用工具：

```json
{
  "name": "json_stringify",
  "arguments": {
    "object": {
      "timestamp": "2024-01-01T12:00:00Z",
      "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
      "data": {
        "message": "Request processed"
      }
    },
    "pretty": true
  }
}
```

## 最佳实践

1. **Base64 编码**:
   - 仅用于编码二进制数据或特殊字符
   - 不用于加密（仅编码）
   - 解码前验证数据格式

2. **JSON 处理**:
   - 验证 JSON 格式后再解析
   - 处理解析错误
   - 使用 `pretty` 参数提高可读性

3. **时间戳**:
   - 使用 ISO 8601 或 RFC 3339 格式用于跨系统兼容
   - 时区使用 UTC
   - 考虑时区转换

4. **UUID 生成**:
   - 使用 v4 UUID 用于大多数场景
   - 避免依赖 UUID 的顺序
   - 适合作为数据库主键或请求 ID

## 相关文档

- [配置指南](/zh/guides/configuration/)
- [部署指南](/zh/guides/deployment/)
- [Kubernetes 工具](/zh/tools/kubernetes/)
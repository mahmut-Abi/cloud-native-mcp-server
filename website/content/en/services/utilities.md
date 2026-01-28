---
title: "Utilities Service"
weight: 10
---

# Utilities Service

The Utilities service provides general-purpose utility tools with 6 tools for common operations and system utilities.

## Overview

The Utilities service in Cloud Native MCP Server provides general-purpose tools for common operations, data transformation, and system utilities. It provides tools for encoding, decoding, and basic data manipulation.

### Key Capabilities

{{< columns >}}
### 🔧 System Utilities
Common system operations and utility functions.
<--->

### 🔄 Data Transformation
Encoding, decoding, and data format conversion.
{{< /columns >}}

{{< columns >}}
### ⚙️ Configuration Tools
System configuration and metadata utilities.
<--->

### 📊 Information Tools
System information and statistics utilities.
{{< /columns >}}

---

## Available Tools (6)

### Encoding and Transformation
- **utils-base64-encode**: Base64 encode
- **utils-base64-decode**: Base64 decode
- **utils-json-parse**: JSON parse
- **utils-json-stringify**: JSON stringify

### System Utilities
- **utils-get-time**: Get current server time
- **utils-get-uuid**: Generate UUID
- **utils-get-stats**: Get server statistics
- **utils-get-versions**: Get component versions
- **utils-get-uptime**: Get server uptime
- **utils-get-config**: Get server configuration
- **utils-get-env**: Get environment variables

---

## Quick Examples

### Base64 encode a string

```json
{
  "method": "tools/call",
  "params": {
    "name": "utils-base64-encode",
    "arguments": {
      "input": "Hello, World!"
    }
  }
}
```

### Get server time

```json
{
  "method": "tools/call",
  "params": {
    "name": "utils-get-time",
    "arguments": {}
  }
}
```

### Get server statistics

```json
{
  "method": "tools/call",
  "params": {
    "name": "utils-get-stats",
    "arguments": {}
  }
}
```

---

## Best Practices

- Use utility tools for common data transformations
- Monitor system statistics for performance insights
- Regularly check server configuration and versions
- Use UUID generation for unique identifiers
- Implement proper error handling for utility operations

## Next Steps

- [Getting Started](/en/getting-started/) for quick setup
- [Configuration Guides](/en/guides/configuration/) for detailed setup
- [Tools Reference](/en/docs/tools/) for complete documentation
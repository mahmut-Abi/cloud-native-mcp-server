---
title: "Loki Service"
weight: 5
---

# Loki Service

The Loki service provides log aggregation and LogQL querying with 7 tools for label discovery, stream inspection, and log-first troubleshooting.

## Overview

The Loki service in Cloud Native MCP Server helps AI assistants and operators inspect logs without starting from raw full-text dumps. It emphasizes compact summaries first, then targeted LogQL queries when more detail is needed.

### Key Capabilities

{{< columns >}}
### 🪵 LogQL Queries
Run instant and range LogQL queries against Loki streams.
<--->

### 🧭 Label Discovery
Discover labels, values, and indexed series before building larger queries.
{{< /columns >}}

{{< columns >}}
### 🎯 LLM-Friendly Summaries
Start with compact stream summaries instead of full log payloads.
<--->

### ✅ Connectivity Checks
Verify the configured Loki endpoint and auth settings quickly.
{{< /columns >}}

---

## Available Tools (7)

- `loki_query_logs_summary`
- `loki_query`
- `loki_query_range`
- `loki_get_label_names`
- `loki_get_label_values`
- `loki_get_series`
- `loki_test_connection`

---

## Quick Examples

### Start with a compact log summary

```json
{
  "method": "tools/call",
  "params": {
    "name": "loki_query_logs_summary",
    "arguments": {
      "query": "{namespace=\"prod\"} |= \"error\"",
      "limit": 50
    }
  }
}
```

### List available values for a label

```json
{
  "method": "tools/call",
  "params": {
    "name": "loki_get_label_values",
    "arguments": {
      "label": "namespace"
    }
  }
}
```

### Inspect indexed series before a wider query

```json
{
  "method": "tools/call",
  "params": {
    "name": "loki_get_series",
    "arguments": {
      "matchers": ["{app=\"api\"}"]
    }
  }
}
```

---

## Best Practices

- Start with `loki_query_logs_summary` before `loki_query_range`.
- Narrow time windows and stream selectors aggressively.
- Use label discovery tools before guessing selector keys and values.
- Prefer targeted pipelines like `|=`, `|~`, and parsed labels over dumping broad streams.

## Next Steps

- [Prometheus Service](/services/prometheus/) for metrics correlation
- [Jaeger Service](/services/jaeger/) for trace correlation
- [Configuration Guide](/docs/configuration/) for runtime setup

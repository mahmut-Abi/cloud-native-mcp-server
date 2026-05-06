# Kubernetes Calling Notes

This file is for tasks that use `kubernetes_*` tools exposed by a compatible cloud-native MCP server.

## Canonical Argument Style

- Prefer flat arguments:

```json
{
  "kind": "Pod",
  "namespace": "open-telemetry"
}
```

- Nested `params` is tolerated by many handlers for compatibility, but it is not the preferred shape:

```json
{
  "params": {
    "kind": "Pod",
    "namespace": "open-telemetry"
  }
}
```

## `kind` Rules

- `kind` is required by most resource-oriented tools.
- Good examples:
  - `Pod`
  - `Deployment`
  - `StatefulSet`
  - `DaemonSet`
  - `Service`
  - `ConfigMap`
  - `Secret`
  - `Namespace`
  - `Node`
- If the user speaks in plural resource names, convert them to a concrete kind before calling the tool.

## Common Response Shapes

- `kubernetes_get_api_versions`: array of strings
- `kubernetes_get_api_resources`: array of resource metadata objects
- `kubernetes_list_resources_summary`: object with `items`, `count`, and `pagination`
- `kubernetes_list_resources`: object with `data`, `count`, and `pagination`
- `kubernetes_search_resources`: object with `query`, `kinds`, `matched`, and `resources`
- `kubernetes_wait_for_resource`: object describing the condition reached and attempts made
- `kubernetes_restart_workload`: object with restart status and optional wait result

## Parsing Rules

- Inspect the raw return value before deciding whether to parse.
- If the client already returns an object or array, use it directly.
- If the client returns an MCP envelope, the JSON payload is usually in `content[0].text`.
- Double parsing is the usual cause of errors such as `Unexpected token 'o'` or text that starts with `[object Object]`.

## Tool Selection

- Need an overview: `kubernetes_get_unhealthy_resources`, `kubernetes_get_recent_events`
- Need a compact list: `kubernetes_list_resources_summary`
- Need one compact object: `kubernetes_get_resource_summary`
- Need full object or field extraction: `kubernetes_get_resource`, `kubernetes_list_resources`
- Need fuzzy discovery: `kubernetes_search_resources`
- Need restart or rollout: `kubernetes_restart_workload`, `kubernetes_get_rollout_status`
- Need to wait for a condition: `kubernetes_wait_for_resource`

## Failure Patterns

- `missing required parameter: kind`
  - Add `kind`
  - Or use `kubernetes_search_resources` if the task is name-based discovery

- `Tool not found`
  - Use the exact runtime name returned by `tools/list`, such as `kubernetes_restart_workload`
  - Confirm the server has been restarted after recent tool additions

- JSON parse failure
  - Inspect the raw value first
  - Do not call `JSON.parse(...)` on an object that the wrapper already parsed

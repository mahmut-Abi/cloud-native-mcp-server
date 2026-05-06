---
name: cloud-native-mcp-tools
description: Use when an agent should troubleshoot or operate a cloud-native system through MCP tools instead of direct shell or raw API calls. This skill is agent-neutral and guides issue diagnosis, tool selection, argument formatting, cross-signal investigation, common Kubernetes pitfalls, and response parsing across Kubernetes, Helm, Grafana, Prometheus, Loki, Kibana, Elasticsearch, Alertmanager, Jaeger, OpenTelemetry, and utility tools.
---

# Cloud Native MCP Tools

## Overview

Use this skill when the user wants troubleshooting, incident analysis, health checks, or controlled operations performed through a cloud-native MCP tool layer.
It is intentionally agent-neutral. The only assumptions are that the host agent can list tools, call a tool by name with JSON arguments, and inspect the returned tool result.

The primary goal is to help the agent diagnose problems with evidence before changing state.
The secondary goal is to pick the smallest correct tool first, send arguments in the shape the server expects, and avoid common parse or naming mistakes.
If the host supports multiple skills, pair this skill with `mcp-tool-operator` for generic MCP operating rules and use this skill as the server-specific profile.

## Quick Start

1. Identify the symptom, scope, and suspected layer.
2. Start with read-only summary tools to build a quick health snapshot.
3. Narrow to the affected workload, namespace, service, metric, log stream, or trace.
4. Correlate at least two signal types when possible: resource state, events, metrics, logs, traces, or alerts.
5. Use the exact runtime tool name returned by the server, usually `snake_case`.
6. Send flat structured arguments unless a client wrapper forces a different shape.
7. Inspect the raw tool return before applying `JSON.parse(...)`.

## Troubleshooting Workflow

1. Establish scope:
   cluster-wide, namespace, workload, pod, service, alert, metric, log stream, or trace.
2. Build a health snapshot:
   prefer summary tools such as unhealthy resources, recent events, firing alerts, target summaries, or trace service summaries.
3. Identify the narrowest failing object:
   one workload, one pod set, one alert rule, one metric query, one log query, or one trace family.
4. Correlate signals:
   combine Kubernetes state with logs, metrics, or traces before concluding root cause.
5. Form a concrete diagnosis:
   describe what is failing, where it is failing, and the strongest evidence.
6. Propose or perform the smallest corrective action:
   restart, patch, scale, drain, silence, or configuration change only when justified.
7. Verify after action:
   re-run the relevant read tools to confirm recovery.

Read [references/troubleshooting-workflow.md](references/troubleshooting-workflow.md) for common issue patterns.

## Operations Workflow

1. Confirm the requested change and target scope:
   cluster, namespace, workload, release, dashboard, alert, silence, or tracing component.
2. Read current state first:
   inspect summaries or the current object before changing anything.
3. Choose the smallest mutation:
   create, patch, scale, restart, install, upgrade, update, or delete.
4. Prefer targeted changes over broad replacements:
   patch one field or one object before replacing an entire resource.
5. Apply the change with explicit arguments and the exact runtime tool name.
6. Verify post-change state:
   use read tools, rollout tools, metrics, logs, or traces.
7. If the change failed or worsened the state, propose or execute the smallest rollback path.

Read [references/operations-workflow.md](references/operations-workflow.md) for create, update, delete, rollout, and rollback patterns.

## Portability

- Treat this skill as a tool-usage guide, not as a client-specific implementation.
- Adapt invocation syntax to the host agent or MCP runtime.
- Trust the live `tools/list` output over guessed names or stale local notes.
- If the target workspace includes a `docs/TOOLS.md` file, use it as a secondary reference for inventory and response shapes.
- `agents/openai.yaml` is optional metadata for OpenAI-compatible runtimes. Other agents can ignore it.

## Service Selection

- Use `kubernetes_*` for cluster resources, logs, events, exec, rollout, restart, wait, RBAC, and node maintenance.
- Use `helm_*` for release lifecycle, chart lookup, manifests, values, and repository operations.
- Use `prometheus_*` for metrics, labels, series, range queries, target status, rules, alerts, and server metadata.
- Use `loki_*` for log queries, log label discovery, and log series discovery.
- Use `jaeger_*` for trace search, service discovery, operations, dependencies, and trace detail.
- Use `alertmanager_*` for active alerts, silences, routing, and cluster health.
- Use `grafana_*` for dashboards, folders, data sources, annotations, snapshots, and Grafana-managed alert rules.
- Use `kibana_*` and `elasticsearch_*` for the Elastic stack.
- Use `opentelemetry_*` for collector and instrumentation workflows exposed by this server.
- Use `utilities_*` only for generic helper tasks that do not belong to a domain service.

Read [references/service-map.md](references/service-map.md) when you need first-choice tools for a user intent.

## Signal Correlation

- Kubernetes gives state, events, rollout, scheduling, pod health, and node health.
- Prometheus gives current and historical metric evidence.
- Loki gives workload or component logs.
- Jaeger gives request path and latency or error traces.
- Alertmanager gives the current alert surface and silencing state.
- Grafana helps discover dashboards and alert rules that explain what operators are looking at.
- Helm helps confirm release history or recent chart changes.

Use multiple signals when the first tool only shows symptoms, not cause.

## Operation Categories

- Read and discovery:
  list, get, search, summary, describe, status, health, labels, series, traces, dashboards.
- Create:
  create Kubernetes resources, install Helm releases, create Grafana objects, create silences.
- Update:
  patch resources, scale workloads, restart workloads, upgrade Helm releases, update dashboards, rules, or data sources.
- Delete:
  delete resources, uninstall releases, delete dashboards, annotations, or silences.
- Verify:
  rollout status, wait conditions, alerts, metrics, logs, traces, and connection tests.

## Working Rules

- Start read-only unless the user clearly asked for a state-changing action.
- Prefer summary tools before full object or full log retrieval.
- Prefer diagnosis before remediation. Do not restart or patch first unless the user explicitly asked for that action.
- For work tasks, read current state before create, patch, scale, restart, or delete unless the task is trivial and the target is fully specified.
- When a field represents an object or array, send structured JSON instead of a JSON-encoded string when the client supports it.
- For Prometheus, Jaeger, and OpenTelemetry time parameters, use RFC3339 timestamps.
- If a tool returns a large object and you only need a few fields, switch to summary tools or JSONPath extraction instead of post-filtering locally.
- If a tool name is missing, confirm it from the running server's inventory instead of guessing aliases.
- When answering, separate observed facts from inference.
- After a state-changing action, verify with a second read or health tool.

## Kubernetes Rules

- Most Kubernetes resource tools require `kind`. Do not omit it for `list`, `get`, `summary`, `describe`, `patch`, `delete`, or `wait`.
- If the user says "pods", "deployments", or another resource class informally, convert that into a concrete `kind` such as `Pod` or `Deployment`.
- `kubernetes_search_resources` is the main exception: it accepts `kind` or `resourceTypes`, and `query` or `name`.
- Prefer these tools first:
  - `kubernetes_list_resources_summary`
  - `kubernetes_get_resource_summary`
  - `kubernetes_get_recent_events`
  - `kubernetes_get_unhealthy_resources`
- Use `kubernetes_get_resource` or `kubernetes_list_resources` only when you need the full object or precise JSONPath output.
- Use `kubernetes_patch_resource` for targeted changes. For `merge` or `apply`, send an object. For `json`, send an RFC 6902 array.
- For restart and rollout workflows, the usual order is `kubernetes_restart_workload` then `kubernetes_get_rollout_status` or `kubernetes_wait_for_resource`.
- Some cloud-native MCP servers tolerate nested `params` for compatibility, but flat arguments remain the canonical form.

Read [references/kubernetes.md](references/kubernetes.md) when the task touches Kubernetes.

## Investigation Patterns

- Pod or workload unhealthy:
  start with `kubernetes_get_unhealthy_resources`, `kubernetes_get_recent_events`, then inspect the specific workload or pod summary, then logs.
- Deployment changed but traffic still failing:
  check rollout status, pod readiness, recent events, service endpoints, then logs and traces.
- Alert firing:
  inspect the alert first, then the backing metrics, then the affected workload state and logs.
- Metrics missing:
  check Prometheus target summaries, then Kubernetes service or pod state, then logs from the exporter or collector.
- High latency or request failure:
  start from traces or alerting metrics, then narrow to workload logs and rollout history.
- Suspected bad rollout:
  inspect Helm release status or history, workload rollout status, pod readiness, and recent events before any restart or rollback.

## Work Task Patterns

- Create a Kubernetes object:
  inspect namespace and nearby resources first, then use `kubernetes_create_resource`, then verify with summary or get tools.
- Update a Kubernetes object:
  read current state, prefer `kubernetes_patch_resource`, then verify with get, summary, rollout, or wait tools.
- Delete a Kubernetes object:
  confirm exact target identity and namespace, then use `kubernetes_delete_resource`, then verify absence or replacement behavior.
- Restart or scale a workload:
  read workload summary first, apply `kubernetes_restart_workload` or `kubernetes_scale_resource`, then verify rollout and readiness.
- Install or upgrade a Helm release:
  inspect release status or values first, use install or upgrade, then verify rollout, events, and release status.
- Change dashboards, alerts, or data sources:
  fetch the existing object first, apply the update, then re-read and run connection or health checks where available.
- Silence or unsilence alerts:
  inspect current alerts and silences first, create or remove the silence, then confirm effective alert state.

## Response Handling

- Do not assume the tool result is always a JSON string.
- Some clients return an already-parsed object or array.
- Some clients return an MCP envelope, where the actual JSON payload sits in a text field such as `content[0].text`.
- Parse only after inspecting the raw shape.
- When a wrapper returns `[object Object]` parse errors, the usual cause is double parsing.

## Troubleshooting

- `Tool not found`: use the exact `snake_case` runtime name and confirm the server has been restarted after recent tool additions.
- `missing required parameter: kind`: add `kind` for Kubernetes resource tools, or use `kubernetes_search_resources` if the task is discovery-oriented.
- `Unexpected token` while parsing: inspect the raw return value first; it may already be an object.
- Excessive output: back up to a summary tool, add filters, or use JSONPath.
- No obvious root cause from one signal: correlate Kubernetes plus one of metrics, logs, or traces before concluding.

## References

- Use [references/service-map.md](references/service-map.md) for intent-to-tool selection.
- Use [references/kubernetes.md](references/kubernetes.md) for Kubernetes-specific argument and parsing rules.
- Use [references/operations-workflow.md](references/operations-workflow.md) for change-management and CRUD sequences.
- Use [references/troubleshooting-workflow.md](references/troubleshooting-workflow.md) for symptom-driven investigation sequences.
- If the target repository includes `docs/TOOLS.md`, use it for exact inventory and response-shape notes.

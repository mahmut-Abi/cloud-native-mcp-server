# Troubleshooting Workflow

Use this file when the agent should diagnose an issue, not just call a tool.

## Default Sequence

1. Define the symptom precisely.
   Examples: pod crash loop, failed rollout, high latency, no logs, missing metrics, firing alert, bad trace latency.
2. Define scope.
   Cluster, namespace, workload, service, pod, alert rule, dashboard, metric, log stream, or trace.
3. Start with summary tools.
   Avoid full object dumps or broad log pulls until the scope is narrow.
4. Find the narrow failing unit.
   One workload, one pod group, one metric series, one alert, or one request path.
5. Correlate signals.
   Do not stop at one symptom source if the cause is still unclear.
6. State the diagnosis and evidence.
7. Only then propose or perform remediation.

## Symptom-Driven Paths

### Kubernetes workload unhealthy

- `kubernetes_get_unhealthy_resources`
- `kubernetes_get_recent_events`
- `kubernetes_list_resources_summary` or `kubernetes_get_resource_summary`
- `kubernetes_get_pod_logs`
- `kubernetes_get_rollout_status` or `kubernetes_wait_for_resource`

### Alert firing

- `alertmanager_get_alerts`
- `prometheus_alerts_summary` or `prometheus_rules_summary`
- `prometheus_query` or `prometheus_query_range`
- Correlate with Kubernetes state and logs

### Missing metrics or scrape failures

- `prometheus_targets_summary`
- `prometheus_get_target_metadata` or `prometheus_get_metrics_metadata`
- Kubernetes summaries for the exporter or collector
- Relevant logs from workload or collector

### Log-based failures

- `loki_query_logs_summary`
- `loki_query` or `loki_query_range`
- Kubernetes workload or pod summaries
- Prometheus metrics if the logs imply resource or traffic issues

### High latency or request failures

- `jaeger_get_traces_summary` or `jaeger_search_traces`
- `jaeger_get_trace`
- `prometheus_query_range`
- Workload logs and rollout status

### Suspected bad deployment

- `helm_get_release_status` or `helm_get_release_history`
- `kubernetes_get_rollout_status`
- `kubernetes_get_recent_events`
- `kubernetes_get_pod_logs`

## Diagnostic Discipline

- Facts first, inference second.
- Prefer the narrowest tool that can confirm or reject a hypothesis.
- Avoid broad restarts as a first response.
- After remediation, verify with the same signal that showed the problem.

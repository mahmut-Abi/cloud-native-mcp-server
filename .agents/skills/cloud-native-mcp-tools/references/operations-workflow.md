# Operations Workflow

Use this file when the agent should perform cloud-native work tasks, not just diagnose issues.

## Default Change Sequence

1. Confirm intent.
   Understand whether the task is create, update, delete, restart, scale, install, upgrade, or rollback.
2. Read current state.
   Inspect the current object, release, alert, dashboard, or service before mutation.
3. Minimize blast radius.
   Prefer one target, one namespace, one workload, or one release at a time.
4. Choose the smallest mutation.
   Prefer patch over replace, scale over redeploy, restart over broader changes when it satisfies the request.
5. Apply the change.
6. Verify immediately.
   Use a read, rollout, wait, health, metric, log, or trace tool.
7. Roll back or escalate if verification fails.

## CRUD Patterns

### Create

- Kubernetes resources:
  use `kubernetes_create_resource`, then verify with `kubernetes_get_resource_summary` or `kubernetes_get_resource`
- Helm releases:
  use `helm_install_release`, then verify with release status and workload rollout
- Grafana objects:
  use the relevant create tool, then re-read the object and test the dependent connection where possible
- Alertmanager silences:
  create the silence, then verify silence presence and alert effect

### Read

- Prefer summary, search, health, and paginated tools first
- Use full object reads only when details are needed for a mutation or diagnosis

### Update

- Kubernetes resources:
  prefer `kubernetes_patch_resource`
- Kubernetes scaling or restart:
  use `kubernetes_scale_resource` or `kubernetes_restart_workload`, then verify rollout
- Helm releases:
  use `helm_upgrade_release`, then inspect status, events, and workload readiness
- Grafana, Kibana, or Elastic objects:
  fetch the current object first, then update with the smallest necessary change

### Delete

- Confirm exact identity and scope before deleting
- For Kubernetes, use `kubernetes_delete_resource`
- For Helm, use `helm_uninstall_release`
- For dashboards, annotations, silences, or similar objects, use the dedicated delete tool
- Verify the object is gone and that the surrounding system remains healthy

## Verification Patterns

- Resource change:
  summary or get tool
- Workload change:
  rollout status, wait, events, logs
- Release change:
  release status, workload health, recent events
- Alerting change:
  alerts, silences, rules
- Observability change:
  target status, connection tests, recent logs, traces

## Rollback Guidance

- If the change created a broken rollout, verify whether restart, scale adjustment, or Helm rollback is the smallest corrective action.
- If the change was a patch, read the current object before proposing a compensating patch.
- If the change touched dashboards, alert rules, or silences, restore the previous known-good object if available.
- Always report what changed, what verification failed, and what rollback path is available.

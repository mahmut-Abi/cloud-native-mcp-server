# Kubernetes Handlers Refactoring Plan

## Current Status
- `handlers/handlers.go` - 1,768 lines with 30+ handler functions

## Proposed Split

### 1. `common.go`
Utility functions and helpers

### 2. `pods.go` (~300 lines)
Pod handlers:
- `HandleGetPods()`
- `HandleGetPod()`
- `HandleCreatePod()`
- `HandleUpdatePod()`
- `HandleDeletePod()`
- `HandleGetPodLogs()`
- `HandleGetPodMetrics()`

### 3. `deployments.go` (~300 lines)
Deployment handlers:
- `HandleGetDeployments()`
- `HandleGetDeployment()`
- `HandleCreateDeployment()`
- `HandleUpdateDeployment()`
- `HandleDeleteDeployment()`
- `HandleScaleDeployment()`
- `HandleRestartDeployment()`

### 4. `services.go` (~250 lines)
Service handlers:
- `HandleGetServices()`
- `HandleGetService()`
- `HandleCreateService()`
- `HandleUpdateService()`
- `HandleDeleteService()`

### 5. `configmaps.go` (~200 lines)
ConfigMap handlers:
- `HandleGetConfigMaps()`
- `HandleGetConfigMap()`
- `HandleCreateConfigMap()`
- `HandleUpdateConfigMap()`
- `HandleDeleteConfigMap()`

### 6. `secrets.go` (~200 lines)
Secret handlers:
- `HandleGetSecrets()`
- `HandleGetSecret()`
- `HandleCreateSecret()`
- `HandleUpdateSecret()`
- `HandleDeleteSecret()`

### 7. `namespaces.go` (~150 lines)
Namespace handlers:
- `HandleGetNamespaces()`
- `HandleGetNamespace()`
- `HandleCreateNamespace()`
- `HandleDeleteNamespace()`

### 8. `nodes.go` (~150 lines)
Node handlers:
- `HandleGetNodes()`
- `HandleGetNode()`
- `HandleGetNodeMetrics()`

### 9. `events.go` (~150 lines)
Event handlers:
- `HandleGetEvents()`

## Benefits
- Better organization by Kubernetes resource type
- Easier to maintain and extend
- Clear separation of concerns
- Faster compilation

# Prometheus Handlers Refactoring Plan

## Current Status
- `handlers.go` - 1,010 lines with 25+ handler functions

## Proposed Split

### 1. `common.go`
Utility functions and helpers

### 2. `query.go` (~300 lines)
Query handlers:
- `HandleQuery()`
- `HandleQueryRange()`

### 3. `targets.go` (~150 lines)
Target handlers:
- `HandleGetTargets()`
- `HandleGetTargetsSummary()`

### 4. `alerts.go` (~150 lines)
Alert handlers:
- `HandleGetAlerts()`
- `HandleGetAlertsSummary()`

### 5. `rules.go` (~150 lines)
Rule handlers:
- `HandleGetRules()`
- `HandleGetRulesSummary()`

### 6. `labels.go` (~150 lines)
Label handlers:
- `HandleGetLabelNames()`
- `HandleGetLabelValues()`

### 7. `series.go` (~100 lines)
Series handlers:
- `HandleGetSeries()`

### 8. `tsdb.go` (~150 lines)
TSDB handlers:
- `HandleGetTSDBStats()`
- `HandleGetTSDBStatus()`
- `HandleCreateSnapshot()`
- `HandleGetWALReplayStatus()`

### 9. `metadata.go` (~100 lines)
Metadata handlers:
- `HandleGetMetricsMetadata()`
- `HandleGetTargetMetadata()`

### 10. `info.go` (~50 lines)
Info handlers:
- `HandleTestConnection()`
- `HandleGetServerInfo()`
- `HandleGetRuntimeInfo()`

## Benefits
- Better organization by functionality
- Easier to maintain and test
- Clear separation of concerns
- Faster compilation

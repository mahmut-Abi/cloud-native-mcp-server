# Kibana Handlers Refactoring Plan

## Current Status
- `handlers.go` is 3,204 lines with 70+ handler functions
- File is too large and difficult to maintain

## Proposed Split
The handlers should be split into the following files:

### 1. `common.go`
Utility functions used across all handlers:
- `marshalIndentJSON`
- `parseLimitWithWarnings`
- `getOptionalIntParam`
- `getOptionalBoolParam`
- `getOptionalStringParam`
- `marshalOptimizedResponse`
- `requireStringParam`
- `HandleTestConnection`

### 2. `spaces.go` (~200 lines)
Space-related handlers:
- `HandleGetSpaces`
- `HandleGetSpace`
- `HandleCreateSpace`
- `HandleUpdateSpace`
- `HandleDeleteSpace`
- `HandleSpacesSummary`

### 3. `index_patterns.go` (~300 lines)
Index pattern handlers:
- `HandleGetIndexPatterns`
- `HandleGetIndexPattern`
- `HandleCreateIndexPattern`
- `HandleUpdateIndexPattern`
- `HandleDeleteIndexPattern`
- `HandleSetDefaultIndexPattern`
- `HandleRefreshIndexPatternFields`
- `HandleGetIndexPatternFields`

### 4. `dashboards.go` (~400 lines)
Dashboard handlers:
- `HandleGetDashboards`
- `HandleGetDashboard`
- `HandleCreateDashboard`
- `HandleUpdateDashboard`
- `HandleDeleteDashboard`
- `HandleCloneDashboard`
- `HandleDashboardsPaginated`
- `HandleGetDashboardDetailAdvanced`

### 5. `visualizations.go` (~350 lines)
Visualization handlers:
- `HandleGetVisualizations`
- `HandleGetVisualization`
- `HandleCreateVisualization`
- `HandleUpdateVisualization`
- `HandleDeleteVisualization`
- `HandleCloneVisualization`
- `HandleVisualizationsPaginated`

### 6. `saved_objects.go` (~500 lines)
Saved object handlers:
- `HandleSearchSavedObjects`
- `HandleCreateSavedObject`
- `HandleUpdateSavedObject`
- `HandleDeleteSavedObject`
- `HandleBulkDeleteSavedObjects`
- `HandleExportSavedObjects`
- `HandleImportSavedObjects`
- `HandleSearchSavedObjectsAdvanced`
- `HandleGetSavedSearches`
- `HandleGetSavedSearch`

### 7. `alerts.go` (~450 lines)
Alert rule handlers:
- `HandleGetKibanaAlerts`
- `HandleGetAlertRules`
- `HandleGetAlertRule`
- `HandleCreateAlertRule`
- `HandleUpdateAlertRule`
- `HandleDeleteAlertRule`
- `HandleEnableAlertRule`
- `HandleDisableAlertRule`
- `HandleMuteAlertRule`
- `HandleUnmuteAlertRule`
- `HandleGetAlertRuleTypes`
- `HandleGetAlertRuleHistory`

### 8. `connectors.go` (~350 lines)
Connector handlers:
- `HandleGetConnectors`
- `HandleGetConnector`
- `HandleCreateConnector`
- `HandleUpdateConnector`
- `HandleDeleteConnector`
- `HandleTestConnector`
- `HandleGetConnectorTypes`

### 9. `data_views.go` (~200 lines)
Data view handlers:
- `HandleGetDataViews`
- `HandleGetDataView`
- `HandleCreateDataView`
- `HandleUpdateDataView`
- `HandleDeleteDataView`

### 10. `other.go` (~150 lines)
Other handlers:
- `HandleGetKibanaStatus`
- `HandleQueryLogs`
- `HandleGetCanvasWorkpads`
- `HandleGetLensObjects`
- `HandleGetMaps`
- `HandleGetHealthSummary`

## Implementation Strategy
1. Start with `common.go` - extract utility functions first
2. Create `spaces.go` - simplest and most independent
3. Continue with other files one at a time
4. Test compilation after each split
5. Ensure all imports are correct
6. Update any external references if needed

## Benefits
- Easier to navigate and maintain
- Better code organization
- Easier to test individual components
- Smaller files faster to compile
- Better code review process

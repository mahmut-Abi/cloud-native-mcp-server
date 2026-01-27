# Kibana Tools Refactoring Plan

## Current Status
- `tools.go` - 1,754 lines with 50+ tool definitions

## Proposed Split

### 1. `common.go`
Common utility functions

### 2. `spaces.go` (~200 lines)
Space tools:
- `GetSpacesTool()`
- `GetSpaceTool()`
- `CreateSpaceTool()`
- `UpdateSpaceTool()`
- `DeleteSpaceTool()`
- `SpacesSummaryTool()`

### 3. `index_patterns.go` (~300 lines)
Index pattern tools:
- `GetIndexPatternsTool()`
- `GetIndexPatternTool()`
- `CreateIndexPatternTool()`
- `UpdateIndexPatternTool()`
- `DeleteIndexPatternTool()`
- `SetDefaultIndexPatternTool()`
- `RefreshIndexPatternFieldsTool()`
- `GetIndexPatternFieldsTool()`

### 4. `dashboards.go` (~350 lines)
Dashboard tools:
- `GetDashboardsTool()`
- `GetDashboardTool()`
- `CreateDashboardTool()`
- `UpdateDashboardTool()`
- `DeleteDashboardTool()`
- `CloneDashboardTool()`
- `DashboardsPaginatedTool()`
- `GetDashboardDetailAdvancedTool()`

### 5. `visualizations.go` (~300 lines)
Visualization tools:
- `GetVisualizationsTool()`
- `GetVisualizationTool()`
- `CreateVisualizationTool()`
- `UpdateVisualizationTool()`
- `DeleteVisualizationTool()`
- `CloneVisualizationTool()`
- `VisualizationsPaginatedTool()`

### 6. `saved_objects.go` (~400 lines)
Saved object tools:
- `SearchSavedObjectsTool()`
- `SearchSavedObjectsAdvancedTool()`
- `CreateSavedObjectTool()`
- `UpdateSavedObjectTool()`
- `DeleteSavedObjectTool()`
- `BulkDeleteSavedObjectsTool()`
- `ExportSavedObjectsTool()`
- `ImportSavedObjectsTool()`
- `GetSavedSearchesTool()`
- `GetSavedSearchTool()`

### 7. `alerts.go` (~350 lines)
Alert tools:
- `GetKibanaAlertsTool()`
- `GetAlertRulesTool()`
- `GetAlertRuleTool()`
- `CreateAlertRuleTool()`
- `UpdateAlertRuleTool()`
- `DeleteAlertRuleTool()`
- `EnableAlertRuleTool()`
- `DisableAlertRuleTool()`
- `MuteAlertRuleTool()`
- `UnmuteAlertRuleTool()`
- `GetAlertRuleTypesTool()`
- `GetAlertRuleHistoryTool()`

### 8. `connectors.go` (~300 lines)
Connector tools:
- `GetConnectorsTool()`
- `GetConnectorTool()`
- `CreateConnectorTool()`
- `UpdateConnectorTool()`
- `DeleteConnectorTool()`
- `TestConnectorTool()`
- `GetConnectorTypesTool()`

### 9. `data_views.go` (~200 lines)
Data view tools:
- `GetDataViewsTool()`
- `GetDataViewTool()`
- `CreateDataViewTool()`
- `UpdateDataViewTool()`
- `DeleteDataViewTool()`

### 10. `other.go` (~150 lines)
Other tools:
- `QueryLogsTool()`
- `GetCanvasWorkpadsTool()`
- `GetLensObjectsTool()`
- `GetMapsTool()`
- `GetHealthSummaryTool()`

## Benefits
- Better organization by functionality
- Easier to maintain and extend
- Clear separation of concerns
- Faster compilation
- Better tool management

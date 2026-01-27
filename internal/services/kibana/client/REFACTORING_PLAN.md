# Kibana Client Refactoring Plan

## Current Status
- `client.go` is 2,905 lines with 50+ client methods
- File is too large and difficult to maintain

## Proposed Split

### 1. `client.go` (main file)
Core client structure and HTTP methods:
- `Client` struct
- `NewClient()`
- `makeRequest()`
- `handleResponse()`
- `TestConnection()`
- `GetKibanaStatus()`

### 2. `spaces.go` (~150 lines)
Space operations:
- `GetSpaces()`
- `GetSpace()`
- `CreateSpace()`
- `UpdateSpace()`
- `DeleteSpace()`
- `SpacesSummary()`

### 3. `index_patterns.go` (~200 lines)
Index pattern operations:
- `GetIndexPatterns()`
- `GetIndexPattern()`
- `CreateIndexPattern()`
- `UpdateIndexPattern()`
- `DeleteIndexPattern()`
- `SetDefaultIndexPattern()`
- `RefreshIndexPatternFields()`
- `GetIndexPatternFields()`

### 4. `dashboards.go` (~300 lines)
Dashboard operations:
- `GetDashboards()`
- `GetDashboard()`
- `CreateDashboard()`
- `UpdateDashboard()`
- `DeleteDashboard()`
- `CloneDashboard()`
- `DashboardsPaginated()`
- `GetDashboardDetailAdvanced()`

### 5. `visualizations.go` (~250 lines)
Visualization operations:
- `GetVisualizations()`
- `GetVisualization()`
- `CreateVisualization()`
- `UpdateVisualization()`
- `DeleteVisualization()`
- `CloneVisualization()`
- `VisualizationsPaginated()`

### 6. `saved_objects.go` (~400 lines)
Saved object operations:
- `SearchSavedObjects()`
- `SearchSavedObjectsAdvanced()`
- `CreateSavedObject()`
- `UpdateSavedObject()`
- `DeleteSavedObject()`
- `BulkDeleteSavedObjects()`
- `ExportSavedObjects()`
- `ImportSavedObjects()`
- `GetSavedSearches()`
- `GetSavedSearch()`

### 7. `alerts.go` (~350 lines)
Alert operations:
- `GetAlerts()`
- `GetAlertRules()`
- `GetAlertRule()`
- `CreateAlertRule()`
- `UpdateAlertRule()`
- `DeleteAlertRule()`
- `EnableAlertRule()`
- `DisableAlertRule()`
- `MuteAlertRule()`
- `UnmuteAlertRule()`
- `GetAlertRuleTypes()`
- `GetAlertRuleHistory()`

### 8. `connectors.go` (~300 lines)
Connector operations:
- `GetConnectors()`
- `GetConnector()`
- `CreateConnector()`
- `UpdateConnector()`
- `DeleteConnector()`
- `TestConnector()`
- `GetConnectorTypes()`

### 9. `data_views.go` (~200 lines)
Data view operations:
- `GetDataViews()`
- `GetDataView()`
- `CreateDataView()`
- `UpdateDataView()`
- `DeleteDataView()`

### 10. `other.go` (~150 lines)
Other operations:
- `QueryLogs()`
- `GetCanvasWorkpads()`
- `GetLensObjects()`
- `GetMaps()`
- `GetHealthSummary()`

## Benefits
- Better code organization by functionality
- Easier to maintain and test
- Smaller files compile faster
- Clear separation of concerns

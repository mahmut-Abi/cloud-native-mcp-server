# Grafana Services Refactoring Plan

## Current Status
Multiple files exceed 1,000 lines:
- `client/client.go` - 1,623 lines
- `handlers/handlers.go` - 1,711 lines

## client.go Refactoring (1,623 lines)

### Proposed Split

#### 1. `client.go` (main file)
Core client structure:
- `Client` struct
- `NewClient()`
- `makeRequest()`
- `TestConnection()`

#### 2. `dashboards.go` (~400 lines)
Dashboard operations:
- `GetDashboards()`
- `GetDashboard()`
- `CreateDashboard()`
- `UpdateDashboard()`
- `DeleteDashboard()`

#### 3. `datasources.go` (~350 lines)
Datasource operations:
- `GetDataSources()`
- `GetDataSource()`
- `CreateDataSource()`
- `UpdateDataSource()`
- `DeleteDataSource()`

#### 4. `folders.go` (~200 lines)
Folder operations:
- `GetFolders()`
- `GetFolder()`
- `CreateFolder()`
- `UpdateFolder()`
- `DeleteFolder()`

#### 5. `alerts.go` (~300 lines)
Alert operations:
- `GetAlerts()`
- `GetAlert()`
- `CreateAlert()`
- `UpdateAlert()`
- `DeleteAlert()`

#### 6. `panels.go` (~250 lines)
Panel operations:
- `GetPanels()`
- `GetPanel()`

## handlers.go Refactoring (1,711 lines)

### Proposed Split

#### 1. `common.go`
Utility functions

#### 2. `dashboards.go` (~400 lines)
Dashboard handlers

#### 3. `datasources.go` (~350 lines)
Datasource handlers

#### 4. `folders.go` (~200 lines)
Folder handlers

#### 5. `alerts.go` (~300 lines)
Alert handlers

#### 6. `panels.go` (~250 lines)
Panel handlers

## Benefits
- Improved maintainability
- Better code organization
- Easier testing
- Faster compilation

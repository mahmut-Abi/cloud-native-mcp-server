# Helm Services Refactoring Plan

## Current Status
Multiple files exceed 1,000 lines:
- `client/client.go` - 2,118 lines
- `handlers/handlers.go` - 1,391 lines

## client.go Refactoring (2,118 lines)

### Proposed Split

#### 1. `client.go` (main file)
Core client structure:
- `Client` struct
- `NewClient()`
- `TestConnection()`

#### 2. `repositories.go` (~400 lines)
Repository operations:
- `ListRepositories()`
- `AddRepository()`
- `UpdateRepository()`
- `RemoveRepository()`
- `GetRepository()`

#### 3. `charts.go` (~500 lines)
Chart operations:
- `ListCharts()`
- `SearchCharts()`
- `GetChart()`
- `GetChartVersion()`
- `GetChartValues()`
- `GetChartReadme()`

#### 4. `releases.go` (~600 lines)
Release operations:
- `ListReleases()`
- `GetRelease()`
- `InstallRelease()`
- `UpgradeRelease()`
- `RollbackRelease()`
- `UninstallRelease()`
- `GetReleaseHistory()`
- `GetReleaseStatus()`
- `GetReleaseManifest()`

#### 5. `values.go` (~300 lines)
Value operations:
- `GetValues()`
- `UpdateValues()`
- `DownloadValues()`

## handlers.go Refactoring (1,391 lines)

### Proposed Split

#### 1. `common.go`
Utility functions

#### 2. `repositories.go` (~350 lines)
Repository handlers

#### 3. `charts.go` (~450 lines)
Chart handlers

#### 4. `releases.go` (~550 lines)
Release handlers

#### 5. `values.go` (~300 lines)
Value handlers

## Benefits
- Better organization by functionality
- Easier to maintain and extend
- Clear separation of concerns
- Faster compilation
- Better test coverage potential

// Package handlers provides HTTP handlers for Prometheus MCP operations.
// It implements request handling for Prometheus queries, metrics, targets, and alerts.
//
// This package has been refactored into smaller, focused files:
// - common.go: Utility functions and helpers
// - query.go: Query handlers (HandleQuery, HandleQueryRange)
// - targets.go: Target handlers (HandleGetTargets, HandleGetTargetsSummary)
// - alerts.go: Alert handlers (HandleGetAlerts, HandleGetAlertsSummary)
// - rules.go: Rule handlers (HandleGetRules, HandleGetRulesSummary)
// - labels.go: Label handlers (HandleGetLabelNames, HandleGetLabelValues)
// - series.go: Series handlers (HandleGetSeries)
// - tsdb.go: TSDB handlers (HandleGetTSDBStats, HandleGetTSDBStatus, HandleCreateSnapshot, HandleGetWALReplayStatus)
// - metadata.go: Metadata handlers (HandleGetMetricsMetadata, HandleGetTargetMetadata)
// - info.go: Info handlers (HandleTestConnection, HandleGetServerInfo, HandleGetRuntimeInfo)
package handlers

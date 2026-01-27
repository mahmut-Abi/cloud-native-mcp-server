// Package handlers provides HTTP handlers for Kibana MCP operations.
// It implements request handling for Kibana spaces, dashboards, visualizations, and saved objects.
//
// This file has been refactored from a single 3,224-line file into multiple smaller files:
//   - common.go: Utility functions
//   - spaces.go: Space operations
//   - index_patterns.go: Index pattern operations
//   - dashboards.go: Dashboard operations
//   - visualizations.go: Visualization operations
//   - saved_objects.go: Saved object operations
//   - alerts.go: Alert rule operations
//   - connectors.go: Connector operations
//   - data_views.go: Data view operations
//   - other.go: Other operations (Status, Logs, Canvas, Lens, Maps, Health)
//
// See REFACTORING_PLAN.md for more details about the refactoring.
package handlers

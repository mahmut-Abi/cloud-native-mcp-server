// Package client provides summary extraction utilities for Kibana resources.
package client

type DashboardSummary struct {
	ID    string
	Title string
	Type  string
}

type IndexPatternSummary struct {
	ID      string
	Title   string
	Pattern string
}

type VisualizationSummary struct {
	ID    string
	Title string
	Type  string
}

type SpaceSummary struct {
	ID   string
	Name string
}

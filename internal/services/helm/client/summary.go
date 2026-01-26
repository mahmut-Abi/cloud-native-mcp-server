// Package client provides Helm client operations for the MCP server.
package client

import (
	"time"

	"helm.sh/helm/v3/pkg/release"
)

// ReleaseSummary contains essential information about a Helm release
type ReleaseSummary struct {
	Name        string    `json:"name" yaml:"name"`
	Namespace   string    `json:"namespace" yaml:"namespace"`
	Chart       string    `json:"chart" yaml:"chart"`     // Format: name:version
	Version     int       `json:"version" yaml:"version"` // Revision number
	Status      string    `json:"status" yaml:"status"`   // deployed, failed, pending
	Updated     time.Time `json:"updated" yaml:"updated"` // Last update time
	AppVersion  string    `json:"app_version" yaml:"app_version"`
	Description string    `json:"description" yaml:"description"` // Max 200 chars
}

// ChartSummary contains essential information about a Helm chart
type ChartSummary struct {
	Name        string `json:"name" yaml:"name"`
	Version     string `json:"version" yaml:"version"`
	AppVersion  string `json:"app_version" yaml:"app_version"`
	Description string `json:"description" yaml:"description"` // Max 200 chars
	Repository  string `json:"repository" yaml:"repository"`
}

// ListOptions defines options for list operations
type ListOptions struct {
	Limit  int    `json:"limit" yaml:"limit"`     // Default 50
	Offset int    `json:"offset" yaml:"offset"`   // Default 0
	SortBy string `json:"sort_by" yaml:"sort_by"` // name, updated, status
	Order  string `json:"order" yaml:"order"`     // asc, desc
	Filter string `json:"filter" yaml:"filter"`   // Label filter expression
}

// ExtractReleaseSummary extracts a summary from a release
func ExtractReleaseSummary(rel *release.Release) *ReleaseSummary {
	if rel == nil {
		return nil
	}

	description := rel.Info.Description
	if len(description) > 200 {
		description = description[:200]
	}

	chart := rel.Chart.ChartFullPath()
	if chart == "" {
		chart = rel.Chart.Metadata.Name + ":" + rel.Chart.Metadata.Version
	}

	return &ReleaseSummary{
		Name:        rel.Name,
		Namespace:   rel.Namespace,
		Chart:       chart,
		Version:     rel.Version,
		Status:      rel.Info.Status.String(),
		Updated:     rel.Info.LastDeployed.Time,
		AppVersion:  rel.Chart.Metadata.AppVersion,
		Description: description,
	}
}

// ExtractReleaseSummaries extracts summaries from multiple releases
func ExtractReleaseSummaries(releases []*release.Release) []*ReleaseSummary {
	var summaries []*ReleaseSummary
	for _, rel := range releases {
		if summary := ExtractReleaseSummary(rel); summary != nil {
			summaries = append(summaries, summary)
		}
	}
	return summaries
}

// Paginate applies pagination to a slice of summaries
func Paginate(summaries []*ReleaseSummary, opts ListOptions) []*ReleaseSummary {
	if opts.Limit <= 0 {
		opts.Limit = 50
	}

	start := opts.Offset
	if start >= len(summaries) {
		return []*ReleaseSummary{}
	}

	end := start + opts.Limit
	if end > len(summaries) {
		end = len(summaries)
	}

	return summaries[start:end]
}

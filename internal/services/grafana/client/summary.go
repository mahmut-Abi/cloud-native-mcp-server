// Package client provides summary extraction utilities for Grafana resources.
// This file contains helper functions to extract essential information from
// full Grafana objects for LLM-friendly output.
package client

// DashboardSummary represents a lightweight summary of a Grafana dashboard
type DashboardSummary struct {
	ID        int      `json:"id"`
	UID       string   `json:"uid"`
	Title     string   `json:"title"`
	FolderID  int      `json:"folder_id,omitempty"`
	FolderUID string   `json:"folder_uid,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	IsStarred bool     `json:"is_starred"`
	URL       string   `json:"url,omitempty"`
}

// DataSourceSummary represents a lightweight summary of a Grafana data source
type DataSourceSummary struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	IsDefault bool   `json:"is_default"`
	URL       string `json:"url,omitempty"`
	IsProxy   bool   `json:"is_proxy,omitempty"`
}

// FolderSummary represents a lightweight summary of a Grafana folder
type FolderSummary struct {
	ID    int    `json:"id"`
	UID   string `json:"uid"`
	Title string `json:"title"`
	URL   string `json:"url,omitempty"`
}

// ExtractDashboardSummary extracts essential summary information from a dashboard
func ExtractDashboardSummary(dashboard *Dashboard) *DashboardSummary {
	if dashboard == nil {
		return nil
	}

	return &DashboardSummary{
		ID:        dashboard.ID,
		UID:       dashboard.UID,
		Title:     dashboard.Title,
		FolderID:  dashboard.FolderID,
		FolderUID: dashboard.FolderUID,
		Tags:      dashboard.Tags,
		IsStarred: dashboard.IsStarred,
		URL:       dashboard.URL,
	}
}

// ExtractDashboardSummaries extracts summaries from multiple dashboards
func ExtractDashboardSummaries(dashboards []*Dashboard) []*DashboardSummary {
	var summaries []*DashboardSummary

	for _, dashboard := range dashboards {
		if summary := ExtractDashboardSummary(dashboard); summary != nil {
			summaries = append(summaries, summary)
		}
	}

	return summaries
}

// ExtractDataSourceSummary extracts essential summary information from a data source
func ExtractDataSourceSummary(ds *DataSource) *DataSourceSummary {
	if ds == nil {
		return nil
	}

	return &DataSourceSummary{
		ID:        ds.ID,
		Name:      ds.Name,
		Type:      ds.Type,
		IsDefault: ds.IsDefault,
		URL:       ds.URL,
	}
}

// ExtractDataSourceSummaries extracts summaries from multiple data sources
func ExtractDataSourceSummaries(dataSources []*DataSource) []*DataSourceSummary {
	var summaries []*DataSourceSummary

	for _, ds := range dataSources {
		if summary := ExtractDataSourceSummary(ds); summary != nil {
			summaries = append(summaries, summary)
		}
	}

	return summaries
}

// ExtractFolderSummary extracts essential summary information from a folder
func ExtractFolderSummary(folder *Folder) *FolderSummary {
	if folder == nil {
		return nil
	}

	return &FolderSummary{
		ID:    folder.ID,
		UID:   folder.UID,
		Title: folder.Title,
		URL:   folder.URL,
	}
}

// ExtractFolderSummaries extracts summaries from multiple folders
func ExtractFolderSummaries(folders []*Folder) []*FolderSummary {
	var summaries []*FolderSummary

	for _, folder := range folders {
		if summary := ExtractFolderSummary(folder); summary != nil {
			summaries = append(summaries, summary)
		}
	}

	return summaries
}

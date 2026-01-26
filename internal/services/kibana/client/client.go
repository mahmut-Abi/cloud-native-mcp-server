// Package client provides Kibana HTTP API client functionality.
// It offers operations for interacting with Kibana spaces, index patterns,
// visualizations, dashboards, and other Kibana resources through REST API calls.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	optimize "github.com/mahmut-Abi/k8s-mcp-server/internal/util/performance"
)

// ClientOptions holds configuration parameters for creating a Kibana client.
type ClientOptions struct {
	URL        string        // Kibana server URL
	APIKey     string        // Kibana API key for authentication
	Username   string        // Username for basic authentication
	Password   string        // Password for basic authentication
	Timeout    time.Duration // HTTP request timeout
	SkipVerify bool          // Skip TLS certificate verification
	Space      string        // Kibana space (default: default)
}

// Client provides operations for interacting with Kibana API.
type Client struct {
	baseURL    string            // Base URL for Kibana API
	httpClient *http.Client      // HTTP client for API requests
	apiKey     string            // API key for authentication
	username   string            // Username for basic auth
	password   string            // Password for basic auth
	space      string            // Kibana space
	headers    map[string]string // Additional headers
}

// Space represents a Kibana space.
type Space struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Description      string   `json:"description,omitempty"`
	Color            string   `json:"color,omitempty"`
	Initials         string   `json:"initials,omitempty"`
	DisabledFeatures []string `json:"disabledFeatures,omitempty"`
	ImageURL         string   `json:"imageUrl,omitempty"`
}

// IndexPattern represents a Kibana index pattern.
type IndexPattern struct {
	ID            string                 `json:"id,omitempty"`
	Title         string                 `json:"title"`
	TimeField     string                 `json:"timeFieldName,omitempty"`
	Fields        []IndexPatternField    `json:"fields,omitempty"`
	SourceFilters []SourceFilter         `json:"sourceFilters,omitempty"`
	FieldFormats  map[string]interface{} `json:"fieldFormats,omitempty"`
	Type          string                 `json:"type,omitempty"`
	Attributes    map[string]interface{} `json:"attributes,omitempty"`
}

// IndexPatternField represents a field in an index pattern.
type IndexPatternField struct {
	Name              string `json:"name"`
	Type              string `json:"type"`
	Searchable        bool   `json:"searchable"`
	Aggregatable      bool   `json:"aggregatable"`
	ReadFromDocValues bool   `json:"readFromDocValues,omitempty"`
	Scripted          bool   `json:"scripted,omitempty"`
	Script            string `json:"script,omitempty"`
	Lang              string `json:"lang,omitempty"`
}

// SourceFilter represents a source filter in an index pattern.
type SourceFilter struct {
	Value string `json:"value"`
}

// Dashboard represents a Kibana dashboard.
type Dashboard struct {
	ID                    string                 `json:"id,omitempty"`
	Title                 string                 `json:"title"`
	Description           string                 `json:"description,omitempty"`
	PanelsJSON            string                 `json:"panelsJSON,omitempty"`
	OptionsJSON           string                 `json:"optionsJSON,omitempty"`
	UIStateJSON           string                 `json:"uiStateJSON,omitempty"`
	Version               int                    `json:"version,omitempty"`
	TimeRestore           bool                   `json:"timeRestore,omitempty"`
	TimeTo                string                 `json:"timeTo,omitempty"`
	TimeFrom              string                 `json:"timeFrom,omitempty"`
	RefreshInterval       map[string]interface{} `json:"refreshInterval,omitempty"`
	KibanaSavedObjectMeta map[string]interface{} `json:"kibanaSavedObjectMeta,omitempty"`
	Type                  string                 `json:"type,omitempty"`
	Attributes            map[string]interface{} `json:"attributes,omitempty"`
}

// Visualization represents a Kibana visualization.
type Visualization struct {
	ID                    string                 `json:"id,omitempty"`
	Title                 string                 `json:"title"`
	VisState              string                 `json:"visState,omitempty"`
	UIStateJSON           string                 `json:"uiStateJSON,omitempty"`
	Description           string                 `json:"description,omitempty"`
	Version               int                    `json:"version,omitempty"`
	SavedSearchRefName    string                 `json:"savedSearchRefName,omitempty"`
	KibanaSavedObjectMeta map[string]interface{} `json:"kibanaSavedObjectMeta,omitempty"`
	Type                  string                 `json:"type,omitempty"`
	Attributes            map[string]interface{} `json:"attributes,omitempty"`
}

// SavedObject represents a generic Kibana saved object.
type SavedObject struct {
	ID         string                 `json:"id,omitempty"`
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
	References []Reference            `json:"references,omitempty"`
	Namespaces []string               `json:"namespaces,omitempty"`
	Version    string                 `json:"version,omitempty"`
	Updated    string                 `json:"updated_at,omitempty"`
	Created    string                 `json:"created_at,omitempty"`
}

// Reference represents a reference between saved objects.
type Reference struct {
	Name string `json:"name"`
	Type string `json:"type"`
	ID   string `json:"id"`
}

// SearchResult represents search results from Kibana.
type SearchResult struct {
	Page         int           `json:"page"`
	PerPage      int           `json:"per_page"`
	Total        int           `json:"total"`
	SavedObjects []SavedObject `json:"saved_objects"`
}

// NewClient creates a new Kibana client with the specified options.
func NewClient(opts *ClientOptions) (*Client, error) {
	if opts.URL == "" {
		return nil, fmt.Errorf("kibana URL is required")
	}

	// Parse and validate URL
	baseURL, err := url.Parse(opts.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid kibana URL: %w", err)
	}

	// Ensure URL has proper path
	if !strings.HasSuffix(baseURL.Path, "/") {
		baseURL.Path += "/"
	}
	baseURL.Path += "api/"

	// Create optimized HTTP client with timeout
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	httpClient := optimize.NewOptimizedHTTPClientWithTimeout(timeout)

	// Set default space if not provided
	space := opts.Space
	if space == "" {
		space = "default"
	}

	// Prepare headers
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["kbn-xsrf"] = "true" // Required by Kibana API

	client := &Client{
		baseURL:    baseURL.String(),
		httpClient: httpClient,
		apiKey:     opts.APIKey,
		username:   opts.Username,
		password:   opts.Password,
		space:      space,
		headers:    headers,
	}

	return client, nil
}

// makeRequest performs an HTTP request to the Kibana API.
func (c *Client) makeRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	// Build URL with space prefix if not default
	var reqURL string
	if c.space != "default" {
		reqURL = c.baseURL + "spaces/" + c.space + "/" + endpoint
	} else {
		reqURL = c.baseURL + endpoint
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// Set authentication
	if c.apiKey != "" {
		req.Header.Set("Authorization", "ApiKey "+c.apiKey)
	} else if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	logrus.WithFields(logrus.Fields{
		"method":   method,
		"url":      reqURL,
		"has_auth": c.apiKey != "" || (c.username != "" && c.password != ""),
		"space":    c.space,
	}).Debug("Making Kibana API request")

	return c.httpClient.Do(req)
}

// handleResponse processes the HTTP response and returns the body.
func (c *Client) handleResponse(resp *http.Response) ([]byte, error) {
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("kibana API error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// GetSpaces retrieves all Kibana spaces.
func (c *Client) GetSpaces(ctx context.Context) ([]Space, error) {
	logrus.Debug("Getting Kibana spaces")

	resp, err := c.makeRequest(ctx, "GET", "spaces/space", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var spaces []Space
	if err := json.Unmarshal(body, &spaces); err != nil {
		return nil, fmt.Errorf("failed to unmarshal spaces: %w", err)
	}

	logrus.WithField("count", len(spaces)).Debug("Retrieved Kibana spaces")
	return spaces, nil
}

// GetSpace retrieves a specific Kibana space.
func (c *Client) GetSpace(ctx context.Context, spaceID string) (*Space, error) {
	logrus.WithField("space_id", spaceID).Debug("Getting Kibana space")

	resp, err := c.makeRequest(ctx, "GET", "spaces/space/"+spaceID, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var space Space
	if err := json.Unmarshal(body, &space); err != nil {
		return nil, fmt.Errorf("failed to unmarshal space: %w", err)
	}

	logrus.Debug("Retrieved Kibana space")
	return &space, nil
}

// GetIndexPatterns retrieves all index patterns.
func (c *Client) GetIndexPatterns(ctx context.Context) ([]IndexPattern, error) {
	logrus.Debug("Getting Kibana index patterns")

	resp, err := c.makeRequest(ctx, "GET", "saved_objects/_find?type=index-pattern", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal index patterns: %w", err)
	}

	var indexPatterns []IndexPattern
	for _, obj := range result.SavedObjects {
		var indexPattern IndexPattern
		indexPattern.ID = obj.ID
		indexPattern.Type = obj.Type
		indexPattern.Attributes = obj.Attributes

		// Extract specific fields from attributes
		if title, ok := obj.Attributes["title"].(string); ok {
			indexPattern.Title = title
		}
		if timeField, ok := obj.Attributes["timeFieldName"].(string); ok {
			indexPattern.TimeField = timeField
		}

		indexPatterns = append(indexPatterns, indexPattern)
	}

	logrus.WithField("count", len(indexPatterns)).Debug("Retrieved Kibana index patterns")
	return indexPatterns, nil
}

// GetDashboards retrieves all dashboards.
func (c *Client) GetDashboards(ctx context.Context) ([]Dashboard, error) {
	logrus.Debug("Getting Kibana dashboards")

	resp, err := c.makeRequest(ctx, "GET", "saved_objects/_find?type=dashboard", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dashboards: %w", err)
	}

	var dashboards []Dashboard
	for _, obj := range result.SavedObjects {
		var dashboard Dashboard
		dashboard.ID = obj.ID
		dashboard.Type = obj.Type
		dashboard.Attributes = obj.Attributes

		// Extract specific fields from attributes
		if title, ok := obj.Attributes["title"].(string); ok {
			dashboard.Title = title
		}
		if description, ok := obj.Attributes["description"].(string); ok {
			dashboard.Description = description
		}

		dashboards = append(dashboards, dashboard)
	}

	logrus.WithField("count", len(dashboards)).Debug("Retrieved Kibana dashboards")
	return dashboards, nil
}

// GetDashboard retrieves a specific dashboard by ID.
func (c *Client) GetDashboard(ctx context.Context, dashboardID string) (*Dashboard, error) {
	logrus.WithField("dashboard_id", dashboardID).Debug("Getting Kibana dashboard")

	resp, err := c.makeRequest(ctx, "GET", "saved_objects/dashboard/"+dashboardID, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var obj SavedObject
	if err := json.Unmarshal(body, &obj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dashboard: %w", err)
	}

	var dashboard Dashboard
	dashboard.ID = obj.ID
	dashboard.Type = obj.Type
	dashboard.Attributes = obj.Attributes

	// Extract specific fields from attributes
	if title, ok := obj.Attributes["title"].(string); ok {
		dashboard.Title = title
	}
	if description, ok := obj.Attributes["description"].(string); ok {
		dashboard.Description = description
	}

	logrus.Debug("Retrieved Kibana dashboard")
	return &dashboard, nil
}

// GetVisualizations retrieves all visualizations.
func (c *Client) GetVisualizations(ctx context.Context) ([]Visualization, error) {
	logrus.Debug("Getting Kibana visualizations")

	resp, err := c.makeRequest(ctx, "GET", "saved_objects/_find?type=visualization", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal visualizations: %w", err)
	}

	var visualizations []Visualization
	for _, obj := range result.SavedObjects {
		var visualization Visualization
		visualization.ID = obj.ID
		visualization.Type = obj.Type
		visualization.Attributes = obj.Attributes

		// Extract specific fields from attributes
		if title, ok := obj.Attributes["title"].(string); ok {
			visualization.Title = title
		}
		if description, ok := obj.Attributes["description"].(string); ok {
			visualization.Description = description
		}

		visualizations = append(visualizations, visualization)
	}

	logrus.WithField("count", len(visualizations)).Debug("Retrieved Kibana visualizations")
	return visualizations, nil
}

// SearchSavedObjects performs a search across saved objects.
func (c *Client) SearchSavedObjects(ctx context.Context, objectType, search string, page, perPage int) (*SearchResult, error) {
	logrus.WithFields(logrus.Fields{
		"type":     objectType,
		"search":   search,
		"page":     page,
		"per_page": perPage,
	}).Debug("Searching Kibana saved objects")

	endpoint := "saved_objects/_find"
	params := url.Values{}

	if objectType != "" {
		params.Set("type", objectType)
	}
	if search != "" {
		params.Set("search", search)
	}
	if page > 0 {
		params.Set("page", fmt.Sprintf("%d", page))
	}
	if perPage > 0 {
		params.Set("per_page", fmt.Sprintf("%d", perPage))
	}

	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search results: %w", err)
	}

	logrus.WithField("count", len(result.SavedObjects)).Debug("Retrieved search results")
	return &result, nil
}

// TestConnection tests the connection to Kibana API.
func (c *Client) TestConnection(ctx context.Context) error {
	logrus.Debug("Testing Kibana connection")

	resp, err := c.makeRequest(ctx, "GET", "status", nil)
	if err != nil {
		return fmt.Errorf("failed to connect to kibana: %w", err)
	}

	_, err = c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("kibana health check failed: %w", err)
	}

	logrus.Debug("Kibana connection test successful")
	return nil
}

// KibanaStatus represents Kibana health status.
type KibanaStatus struct {
	State   string                 `json:"state"`
	Version map[string]interface{} `json:"version,omitempty"`
	Metrics map[string]interface{} `json:"metrics,omitempty"`
}

// SavedSearch represents a Kibana saved search.
type SavedSearch struct {
	ID           string                 `json:"id,omitempty"`
	Title        string                 `json:"title"`
	Description  string                 `json:"description,omitempty"`
	SearchSource map[string]interface{} `json:"searchSource,omitempty"`
	Columns      []string               `json:"columns,omitempty"`
	Sort         map[string]interface{} `json:"sort,omitempty"`
	Attributes   map[string]interface{} `json:"attributes,omitempty"`
}

// GetVisualization retrieves a specific visualization by ID.
func (c *Client) GetVisualization(ctx context.Context, visID string) (*Visualization, error) {
	logrus.WithField("id", visID).Debug("Getting Kibana visualization")

	resp, err := c.makeRequest(ctx, "GET", "saved_objects/visualization/"+visID, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var obj SavedObject
	if err := json.Unmarshal(body, &obj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal visualization: %w", err)
	}

	var visualization Visualization
	visualization.ID = obj.ID
	visualization.Type = obj.Type
	visualization.Attributes = obj.Attributes

	if title, ok := obj.Attributes["title"].(string); ok {
		visualization.Title = title
	}
	if description, ok := obj.Attributes["description"].(string); ok {
		visualization.Description = description
	}

	logrus.Debug("Retrieved visualization")
	return &visualization, nil
}

// GetIndexPattern retrieves a specific index pattern by ID.
func (c *Client) GetIndexPattern(ctx context.Context, indexPatternID string) (*IndexPattern, error) {
	logrus.WithField("id", indexPatternID).Debug("Getting Kibana index pattern")

	resp, err := c.makeRequest(ctx, "GET", "saved_objects/index-pattern/"+indexPatternID, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var obj SavedObject
	if err := json.Unmarshal(body, &obj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal index pattern: %w", err)
	}

	var indexPattern IndexPattern
	indexPattern.ID = obj.ID
	indexPattern.Type = obj.Type
	indexPattern.Attributes = obj.Attributes

	if title, ok := obj.Attributes["title"].(string); ok {
		indexPattern.Title = title
	}

	logrus.Debug("Retrieved index pattern")
	return &indexPattern, nil
}

// GetSavedSearches retrieves all saved searches.
func (c *Client) GetSavedSearches(ctx context.Context) ([]SavedSearch, error) {
	logrus.Debug("Getting Kibana saved searches")

	resp, err := c.makeRequest(ctx, "GET", "saved_objects/_find?type=search", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal saved searches: %w", err)
	}

	var searches []SavedSearch
	for _, obj := range result.SavedObjects {
		var search SavedSearch
		search.ID = obj.ID
		search.Attributes = obj.Attributes

		if title, ok := obj.Attributes["title"].(string); ok {
			search.Title = title
		}
		if description, ok := obj.Attributes["description"].(string); ok {
			search.Description = description
		}

		searches = append(searches, search)
	}

	logrus.WithField("count", len(searches)).Debug("Retrieved saved searches")
	return searches, nil
}

// GetSavedSearch retrieves a specific saved search by ID.
func (c *Client) GetSavedSearch(ctx context.Context, searchID string) (*SavedSearch, error) {
	logrus.WithField("id", searchID).Debug("Getting Kibana saved search")

	resp, err := c.makeRequest(ctx, "GET", "saved_objects/search/"+searchID, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var obj SavedObject
	if err := json.Unmarshal(body, &obj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal saved search: %w", err)
	}

	var search SavedSearch
	search.ID = obj.ID
	search.Attributes = obj.Attributes

	if title, ok := obj.Attributes["title"].(string); ok {
		search.Title = title
	}
	if description, ok := obj.Attributes["description"].(string); ok {
		search.Description = description
	}

	logrus.Debug("Retrieved saved search")
	return &search, nil
}

// GetKibanaStatus retrieves Kibana status and health information.
func (c *Client) GetKibanaStatus(ctx context.Context) (*KibanaStatus, error) {
	logrus.Debug("Getting Kibana status")

	resp, err := c.makeRequest(ctx, "GET", "status", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var status KibanaStatus
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal status: %w", err)
	}

	logrus.Debug("Retrieved Kibana status")
	return &status, nil
}

// Response size limits for Kibana - similar to other optimizations
const (
	defaultLimit      = 20  // Conservative default limit
	maxLimit          = 100 // Maximum allowed limit
	warningLimit      = 50  // Warning threshold for large requests
	defaultLimitNodes = 50  // Default limit for nodes
)

// PaginationInfo represents pagination metadata for Kibana responses
type PaginationInfo struct {
	CurrentPage     int   `json:"currentPage"`
	PerPage         int   `json:"perPage"`
	TotalCount      int64 `json:"totalCount"`
	TotalPages      int   `json:"totalPages"`
	HasNextPage     bool  `json:"hasNextPage"`
	HasPreviousPage bool  `json:"hasPreviousPage"`
}

// SpacesSummary returns optimized spaces information
func (c *Client) SpacesSummary(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	spaces, err := c.GetSpaces(ctx)
	if err != nil {
		return nil, err
	}

	// Apply limit
	if limit <= 0 {
		limit = defaultLimitNodes
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	var summaries []map[string]interface{}
	count := 0
	for _, space := range spaces {
		if count >= limit {
			break
		}
		summary := map[string]interface{}{
			"id":               space.ID,
			"name":             space.Name,
			"description":      space.Description,
			"disabledFeatures": space.DisabledFeatures,
		}
		summaries = append(summaries, summary)
		count++
	}

	return summaries, nil
}

// DashboardsPaginated returns dashboards with pagination support and optimization
func (c *Client) DashboardsPaginated(ctx context.Context, page, perPage int, search string, includeDescription bool) ([]map[string]interface{}, *PaginationInfo, error) {
	// Validate and normalize pagination parameters
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = defaultLimit
	}
	if perPage > maxLimit {
		perPage = maxLimit
	}

	// Build query parameters
	params := url.Values{}
	params.Set("type", "dashboard")
	if search != "" {
		params.Set("search", search)
	}
	params.Set("page", fmt.Sprintf("%d", page))
	params.Set("per_page", fmt.Sprintf("%d", perPage))

	endpoint := "saved_objects/_find?" + params.Encode()

	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, nil, err
	}

	var result SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal dashboards: %w", err)
	}

	// Create optimized summaries
	var summaries []map[string]interface{}
	for _, obj := range result.SavedObjects {
		summary := map[string]interface{}{
			"id":   obj.ID,
			"type": obj.Type,
		}

		// Extract title
		if title, ok := obj.Attributes["title"].(string); ok {
			summary["title"] = title
		}

		// Optionally include description
		if includeDescription {
			if description, ok := obj.Attributes["description"].(string); ok {
				summary["description"] = description
			}
		}

		// Include updated timestamp
		if updated, ok := obj.Attributes["updated_at"].(string); ok {
			summary["updated_at"] = updated
		}

		summaries = append(summaries, summary)
	}

	// Calculate pagination info
	totalPages := (result.Total + perPage - 1) / perPage
	pagination := &PaginationInfo{
		CurrentPage:     page,
		PerPage:         perPage,
		TotalCount:      int64(result.Total),
		TotalPages:      totalPages,
		HasNextPage:     page < totalPages,
		HasPreviousPage: page > 1,
	}

	return summaries, pagination, nil
}

// VisualizationsPaginated returns visualizations with pagination support and optimization
func (c *Client) VisualizationsPaginated(ctx context.Context, page, perPage int, search, visType string) ([]map[string]interface{}, *PaginationInfo, error) {
	// Validate and normalize pagination parameters
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = defaultLimit
	}
	if perPage > maxLimit {
		perPage = maxLimit
	}

	// Build query parameters
	params := url.Values{}
	params.Set("type", "visualization")
	if search != "" {
		params.Set("search", search)
	}
	params.Set("page", fmt.Sprintf("%d", page))
	params.Set("per_page", fmt.Sprintf("%d", perPage))

	endpoint := "saved_objects/_find?" + params.Encode()

	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, nil, err
	}

	var result SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal visualizations: %w", err)
	}

	// Create optimized summaries
	var summaries []map[string]interface{}
	for _, obj := range result.SavedObjects {
		summary := map[string]interface{}{
			"id":   obj.ID,
			"type": obj.Type,
		}

		// Extract title
		if title, ok := obj.Attributes["title"].(string); ok {
			summary["title"] = title
		}

		// Extract visualization type from visState
		if visState, ok := obj.Attributes["visState"].(string); ok {
			var visData map[string]interface{}
			if err := json.Unmarshal([]byte(visState), &visData); err == nil {
				if typ, ok := visData["type"].(string); ok {
					summary["vis_type"] = typ
					// Filter by type if specified
					if visType != "" && typ != visType {
						continue
					}
				}
			}
		}

		// Include updated timestamp
		if updated, ok := obj.Attributes["updated_at"].(string); ok {
			summary["updated_at"] = updated
		}

		summaries = append(summaries, summary)
	}

	// Calculate pagination info
	totalPages := (result.Total + perPage - 1) / perPage
	pagination := &PaginationInfo{
		CurrentPage:     page,
		PerPage:         perPage,
		TotalCount:      int64(result.Total),
		TotalPages:      totalPages,
		HasNextPage:     page < totalPages,
		HasPreviousPage: page > 1,
	}

	return summaries, pagination, nil
}

// SearchSavedObjectsAdvanced provides advanced saved objects search with enhanced filters
func (c *Client) SearchSavedObjectsAdvanced(ctx context.Context, objectType, search string, page, perPage int, sortField, sortOrder, hasReference string, fields []string) (*SearchResult, error) {
	// Validate parameters
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 30
	}
	if perPage > 200 {
		perPage = 200
	}
	if sortField == "" {
		sortField = "title"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	// Build query parameters
	params := url.Values{}
	if objectType != "" {
		params.Set("type", objectType)
	}
	if search != "" {
		params.Set("search", search)
	}
	if hasReference != "" {
		params.Set("has_reference", hasReference)
	}
	if len(fields) > 0 {
		params.Set("fields", strings.Join(fields, ","))
	}
	params.Set("page", fmt.Sprintf("%d", page))
	params.Set("per_page", fmt.Sprintf("%d", perPage))
	params.Set("sort", sortField+":"+sortOrder)

	endpoint := "saved_objects/_find?" + params.Encode()

	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search results: %w", err)
	}

	return &result, nil
}

// GetDashboardDetailAdvanced returns comprehensive dashboard information with optimization
func (c *Client) GetDashboardDetailAdvanced(ctx context.Context, dashboardID string, includePanels, includeUIState, includeTimeOptions bool, outputFormat string) (map[string]interface{}, error) {
	dashboard, err := c.GetDashboard(ctx, dashboardID)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"id":      dashboard.ID,
		"type":    dashboard.Type,
		"queryAt": time.Now().Format(time.RFC3339),
		"metadata": map[string]interface{}{
			"includePanels":      includePanels,
			"includeUIState":     includeUIState,
			"includeTimeOptions": includeTimeOptions,
			"outputFormat":       outputFormat,
		},
	}

	// Basic information
	result["title"] = dashboard.Title
	if dashboard.Description != "" {
		result["description"] = dashboard.Description
	}
	result["version"] = dashboard.Version

	// Include panels information if requested
	if includePanels && dashboard.PanelsJSON != "" {
		var panels []interface{}
		if err := json.Unmarshal([]byte(dashboard.PanelsJSON), &panels); err == nil {
			result["panels"] = panels
		}
	}

	// Include UI state if requested
	if includeUIState && dashboard.UIStateJSON != "" {
		var uiState map[string]interface{}
		if err := json.Unmarshal([]byte(dashboard.UIStateJSON), &uiState); err == nil {
			result["uiState"] = uiState
		}
	}

	// Include time options if requested
	if includeTimeOptions {
		result["timeRestore"] = dashboard.TimeRestore
		result["timeTo"] = dashboard.TimeTo
		result["timeFrom"] = dashboard.TimeFrom
		if dashboard.RefreshInterval != nil {
			result["refreshInterval"] = dashboard.RefreshInterval
		}
	}

	// Apply output format optimization
	switch outputFormat {
	case "compact":
		// Remove verbose fields for compact output
		delete(result, "panels")
		delete(result, "uiState")
		delete(result, "timeTo")
		delete(result, "timeFrom")
		delete(result, "refreshInterval")
	case "verbose":
		// Add raw data for complete analysis
		result["rawAttributes"] = dashboard.Attributes
		result["panelsJSON"] = dashboard.PanelsJSON
		result["uiStateJSON"] = dashboard.UIStateJSON
		result["optionsJSON"] = dashboard.OptionsJSON
	}

	return result, nil
}

// GetHealthSummary returns lightweight health information
func (c *Client) GetHealthSummary(ctx context.Context, level string, includeSavedObjects bool) (map[string]interface{}, error) {
	status, err := c.GetKibanaStatus(ctx)
	if err != nil {
		return nil, err
	}

	summary := map[string]interface{}{
		"state":     status.State,
		"timestamp": time.Now().Format(time.RFC3339),
		"metadata": map[string]interface{}{
			"level":               level,
			"includeSavedObjects": includeSavedObjects,
			"optimizedFor":        "monitoring and LLM efficiency",
		},
	}

	// Include version information
	if status.Version != nil {
		summary["version"] = status.Version
	}

	// Include detailed status for higher levels
	if level == "detailed" || level == "metrics" {
		if status.Metrics != nil {
			summary["metrics"] = status.Metrics
		}
	}

	// Include saved objects statistics if requested
	if includeSavedObjects && level == "metrics" {
		// Get count of saved objects by type
		types := []string{"dashboard", "visualization", "index-pattern", "search"}
		savedObjectStats := make(map[string]interface{})

		for _, objectType := range types {
			result, err := c.SearchSavedObjects(ctx, objectType, "", 1, 1)
			if err == nil {
				savedObjectStats[objectType] = map[string]interface{}{
					"count": result.Total,
				}
			}
		}

		if len(savedObjectStats) > 0 {
			summary["savedObjectsStats"] = savedObjectStats
		}
	}

	return summary, nil
}

// helper function to get string field from attributes
func getStringField(attributes map[string]interface{}, key string) string {
	if v, ok := attributes[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// ============ Analysis & Discover Tools ============

// LogQueryResult represents a log search result.
type LogQueryResult struct {
	Hits     []map[string]interface{} `json:"hits"`
	Total    int                      `json:"total"`
	Took     int64                    `json:"took"`
	TimedOut bool                     `json:"timedOut"`
}

// QueryLogs performs a direct log search against Elasticsearch through Kibana API.
func (c *Client) QueryLogs(ctx context.Context, indexPattern string, query string, size int, sortBy string, sortOrder string) (*LogQueryResult, error) {
	logrus.WithFields(logrus.Fields{
		"indexPattern": indexPattern,
		"query":        query,
		"size":         size,
	}).Debug("Querying logs through Kibana API")

	// Build the search body
	searchBody := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"query_string": map[string]interface{}{
							"query": query,
						},
					},
				},
			},
		},
		"size": size,
		"sort": []map[string]interface{}{
			{
				sortBy: map[string]interface{}{
					"order": sortOrder,
				},
			},
		},
	}

	body, err := json.Marshal(searchBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search body: %w", err)
	}

	// Use Elasticsearch search API through Kibana proxy
	path := "_search"
	resp, err := c.makeRequest(ctx, "POST", path, body)
	if err != nil {
		return nil, err
	}

	respBody, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result struct {
		Took int64 `json:"took"`
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []map[string]interface{} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search result: %w", err)
	}

	return &LogQueryResult{
		Hits:     result.Hits.Hits,
		Total:    result.Hits.Total.Value,
		Took:     result.Took,
		TimedOut: false,
	}, nil
}

// CanvasWorkpad represents a Canvas workpad.
type CanvasWorkpad struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Attributes map[string]interface{} `json:"attributes"`
	TotalPages int                    `json:"totalPages,omitempty"`
	Width      int                    `json:"width,omitempty"`
	Height     int                    `json:"height,omitempty"`
	Type       string                 `json:"type"`
}

// GetCanvasWorkpads retrieves all Canvas workpads.
func (c *Client) GetCanvasWorkpads(ctx context.Context) ([]CanvasWorkpad, error) {
	logrus.Debug("Getting Canvas workpads")

	result, err := c.SearchSavedObjects(ctx, "canvas-workpad", "", 1, 100)
	if err != nil {
		return nil, err
	}

	workpads := make([]CanvasWorkpad, 0, len(result.SavedObjects))
	for _, obj := range result.SavedObjects {
		workpad := CanvasWorkpad{
			ID:         obj.ID,
			Name:       getStringField(obj.Attributes, "title"),
			Attributes: obj.Attributes,
			Type:       obj.Type,
		}
		workpads = append(workpads, workpad)
	}

	logrus.WithField("count", len(workpads)).Debug("Retrieved Canvas workpads")
	return workpads, nil
}

// LensObject represents a Lens visualization.
type LensObject struct {
	ID         string                 `json:"id"`
	Title      string                 `json:"title"`
	Attributes map[string]interface{} `json:"attributes"`
	Type       string                 `json:"type"`
}

// GetLensObjects retrieves all Lens visualizations.
func (c *Client) GetLensObjects(ctx context.Context) ([]LensObject, error) {
	logrus.Debug("Getting Lens objects")

	result, err := c.SearchSavedObjects(ctx, "lens", "", 1, 100)
	if err != nil {
		return nil, err
	}

	lenses := make([]LensObject, 0, len(result.SavedObjects))
	for _, obj := range result.SavedObjects {
		lens := LensObject{
			ID:         obj.ID,
			Title:      getStringField(obj.Attributes, "title"),
			Attributes: obj.Attributes,
			Type:       obj.Type,
		}
		lenses = append(lenses, lens)
	}

	logrus.WithField("count", len(lenses)).Debug("Retrieved Lens objects")
	return lenses, nil
}

// MapObject represents a Kibana Map.
type MapObject struct {
	ID         string                 `json:"id"`
	Title      string                 `json:"title"`
	Attributes map[string]interface{} `json:"attributes"`
	Type       string                 `json:"type"`
}

// GetMaps retrieves all Maps.
func (c *Client) GetMaps(ctx context.Context) ([]MapObject, error) {
	logrus.Debug("Getting Maps")

	result, err := c.SearchSavedObjects(ctx, "map", "", 1, 100)
	if err != nil {
		return nil, err
	}

	maps := make([]MapObject, 0, len(result.SavedObjects))
	for _, obj := range result.SavedObjects {
		m := MapObject{
			ID:         obj.ID,
			Title:      getStringField(obj.Attributes, "title"),
			Attributes: obj.Attributes,
			Type:       obj.Type,
		}
		maps = append(maps, m)
	}

	logrus.WithField("count", len(maps)).Debug("Retrieved Maps")
	return maps, nil
}

// KibanaAlert represents a Kibana alerting rule.
type KibanaAlert struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	AlertTypeID string                   `json:"alertTypeId"`
	Schedule    map[string]interface{}   `json:"schedule"`
	Consumer    string                   `json:"consumer"`
	Tags        []string                 `json:"tags,omitempty"`
	Enabled     bool                     `json:"enabled"`
	Actions     []map[string]interface{} `json:"actions,omitempty"`
	Params      map[string]interface{}   `json:"params"`
	CreatedAt   string                   `json:"createdAt"`
	UpdatedAt   string                   `json:"updatedAt"`
}

// GetAlerts retrieves all Kibana alerting rules.
func (c *Client) GetAlerts(ctx context.Context) ([]KibanaAlert, error) {
	logrus.Debug("Getting Kibana alerts")

	resp, err := c.makeRequest(ctx, "GET", "api/alerts", nil)
	if err != nil {
		return nil, err
	}

	respBody, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var alerts []KibanaAlert
	if err := json.Unmarshal(respBody, &alerts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal alerts: %w", err)
	}

	logrus.WithField("count", len(alerts)).Debug("Retrieved Kibana alerts")
	return alerts, nil
}

// GetIndexPatternFields retrieves fields for an index pattern.
func (c *Client) GetIndexPatternFields(ctx context.Context, patternID string) ([]IndexPatternField, error) {
	logrus.WithField("patternID", patternID).Debug("Getting index pattern fields")

	resp, err := c.makeRequest(ctx, "GET", fmt.Sprintf("api/index_patterns/%s/fields", patternID), nil)
	if err != nil {
		return nil, err
	}

	respBody, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result struct {
		Fields []IndexPatternField `json:"fields"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal fields: %w", err)
	}

	logrus.WithField("count", len(result.Fields)).Debug("Retrieved index pattern fields")
	return result.Fields, nil
}

// ============ Write Operations: Spaces ============

// CreateSpace creates a new Kibana space
func (c *Client) CreateSpace(ctx context.Context, space Space) (*Space, error) {
	logrus.WithField("space_id", space.ID).Debug("Creating Kibana space")

	resp, err := c.makeRequest(ctx, "POST", "spaces/space", space)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var createdSpace Space
	if err := json.Unmarshal(body, &createdSpace); err != nil {
		return nil, fmt.Errorf("failed to unmarshal created space: %w", err)
	}

	logrus.WithField("space_id", createdSpace.ID).Debug("Created Kibana space")
	return &createdSpace, nil
}

// UpdateSpace updates an existing Kibana space
func (c *Client) UpdateSpace(ctx context.Context, spaceID string, space Space) (*Space, error) {
	logrus.WithField("space_id", spaceID).Debug("Updating Kibana space")

	resp, err := c.makeRequest(ctx, "PUT", "spaces/space/"+spaceID, space)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var updatedSpace Space
	if err := json.Unmarshal(body, &updatedSpace); err != nil {
		return nil, fmt.Errorf("failed to unmarshal updated space: %w", err)
	}

	logrus.WithField("space_id", spaceID).Debug("Updated Kibana space")
	return &updatedSpace, nil
}

// DeleteSpace deletes a Kibana space
func (c *Client) DeleteSpace(ctx context.Context, spaceID string, force bool) error {
	logrus.WithFields(logrus.Fields{
		"space_id": spaceID,
		"force":    force,
	}).Debug("Deleting Kibana space")

	// Build URL with query parameters
	endpoint := "spaces/space/" + spaceID
	if force {
		endpoint += "?force=true"
	}

	resp, err := c.makeRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	_, err = c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to delete space: %w", err)
	}

	logrus.WithField("space_id", spaceID).Debug("Deleted Kibana space")
	return nil
}

// ============ Write Operations: Index Patterns ============

// CreateIndexPattern creates a new index pattern
func (c *Client) CreateIndexPattern(ctx context.Context, title string, timeField string) (*IndexPattern, error) {
	logrus.WithField("title", title).Debug("Creating index pattern")

	attributes := map[string]interface{}{
		"title": title,
	}
	if timeField != "" {
		attributes["timeFieldName"] = timeField
	}

	obj := SavedObject{
		Type:       "index-pattern",
		Attributes: attributes,
	}

	resp, err := c.makeRequest(ctx, "POST", "saved_objects/index-pattern", obj)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result SavedObject
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal created index pattern: %w", err)
	}

	indexPattern := &IndexPattern{
		ID:         result.ID,
		Title:      title,
		TimeField:  timeField,
		Attributes: result.Attributes,
	}

	logrus.WithField("id", result.ID).Debug("Created index pattern")
	return indexPattern, nil
}

// UpdateIndexPattern updates an index pattern
func (c *Client) UpdateIndexPattern(ctx context.Context, patternID string, title string, timeField string) (*IndexPattern, error) {
	logrus.WithField("pattern_id", patternID).Debug("Updating index pattern")

	attributes := map[string]interface{}{}
	if title != "" {
		attributes["title"] = title
	}
	if timeField != "" {
		attributes["timeFieldName"] = timeField
	}

	obj := map[string]interface{}{
		"attributes": attributes,
	}

	resp, err := c.makeRequest(ctx, "PUT", "saved_objects/index-pattern/"+patternID, obj)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result SavedObject
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal updated index pattern: %w", err)
	}

	indexPattern := &IndexPattern{
		ID:         result.ID,
		Title:      getStringField(result.Attributes, "title"),
		TimeField:  getStringField(result.Attributes, "timeFieldName"),
		Attributes: result.Attributes,
	}

	logrus.WithField("pattern_id", patternID).Debug("Updated index pattern")
	return indexPattern, nil
}

// DeleteIndexPattern deletes an index pattern
func (c *Client) DeleteIndexPattern(ctx context.Context, patternID string) error {
	logrus.WithField("pattern_id", patternID).Debug("Deleting index pattern")

	resp, err := c.makeRequest(ctx, "DELETE", "saved_objects/index-pattern/"+patternID, nil)
	if err != nil {
		return err
	}

	_, err = c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to delete index pattern: %w", err)
	}

	logrus.WithField("pattern_id", patternID).Debug("Deleted index pattern")
	return nil
}

// SetDefaultIndexPattern sets an index pattern as the default
func (c *Client) SetDefaultIndexPattern(ctx context.Context, patternID string) error {
	logrus.WithField("pattern_id", patternID).Debug("Setting default index pattern")

	// Kibana doesn't have a direct API to set default pattern
	// This is typically stored in a config object, so we'll create/update a config
	config := map[string]interface{}{
		"attributes": map[string]interface{}{
			"defaultIndex": patternID,
		},
	}

	resp, err := c.makeRequest(ctx, "POST", "saved_objects/config", config)
	if err != nil {
		return err
	}

	_, err = c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to set default index pattern: %w", err)
	}

	logrus.WithField("pattern_id", patternID).Debug("Set default index pattern")
	return nil
}

// RefreshIndexPatternFields refreshes fields for an index pattern
func (c *Client) RefreshIndexPatternFields(ctx context.Context, patternID string) error {
	logrus.WithField("pattern_id", patternID).Debug("Refreshing index pattern fields")

	resp, err := c.makeRequest(ctx, "POST", fmt.Sprintf("api/index_patterns/%s/fields", patternID), nil)
	if err != nil {
		return err
	}

	_, err = c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to refresh index pattern fields: %w", err)
	}

	logrus.WithField("pattern_id", patternID).Debug("Refreshed index pattern fields")
	return nil
}

// ============ Write Operations: Dashboards ============

// CreateDashboard creates a new dashboard
func (c *Client) CreateDashboard(ctx context.Context, title, description string, timeRestore bool, timeFrom, timeTo string, refreshInterval map[string]interface{}) (*Dashboard, error) {
	logrus.WithField("title", title).Debug("Creating dashboard")

	attributes := map[string]interface{}{
		"title": title,
	}
	if description != "" {
		attributes["description"] = description
	}
	if timeRestore {
		attributes["timeRestore"] = true
	}
	if timeFrom != "" {
		attributes["timeFrom"] = timeFrom
	}
	if timeTo != "" {
		attributes["timeTo"] = timeTo
	}
	if refreshInterval != nil {
		attributes["refreshInterval"] = refreshInterval
	}

	obj := SavedObject{
		Type:       "dashboard",
		Attributes: attributes,
	}

	resp, err := c.makeRequest(ctx, "POST", "saved_objects/dashboard", obj)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result SavedObject
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal created dashboard: %w", err)
	}

	dashboard := &Dashboard{
		ID:          result.ID,
		Title:       title,
		Description: description,
		Type:        "dashboard",
		Attributes:  result.Attributes,
	}

	logrus.WithField("id", result.ID).Debug("Created dashboard")
	return dashboard, nil
}

// UpdateDashboard updates a dashboard
func (c *Client) UpdateDashboard(ctx context.Context, dashboardID string, title, description, panelsJSON string, timeFrom, timeTo string, version int) (*Dashboard, error) {
	logrus.WithField("dashboard_id", dashboardID).Debug("Updating dashboard")

	attributes := map[string]interface{}{}
	if title != "" {
		attributes["title"] = title
	}
	if description != "" {
		attributes["description"] = description
	}
	if panelsJSON != "" {
		attributes["panelsJSON"] = panelsJSON
	}
	if timeFrom != "" {
		attributes["timeFrom"] = timeFrom
	}
	if timeTo != "" {
		attributes["timeTo"] = timeTo
	}

	obj := map[string]interface{}{
		"attributes": attributes,
	}
	if version > 0 {
		obj["version"] = version
	}

	resp, err := c.makeRequest(ctx, "PUT", "saved_objects/dashboard/"+dashboardID, obj)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result SavedObject
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal updated dashboard: %w", err)
	}

	dashboard := &Dashboard{
		ID:         result.ID,
		Title:      getStringField(result.Attributes, "title"),
		Attributes: result.Attributes,
	}

	logrus.WithField("dashboard_id", dashboardID).Debug("Updated dashboard")
	return dashboard, nil
}

// DeleteDashboard deletes a dashboard
func (c *Client) DeleteDashboard(ctx context.Context, dashboardID string) error {
	logrus.WithField("dashboard_id", dashboardID).Debug("Deleting dashboard")

	resp, err := c.makeRequest(ctx, "DELETE", "saved_objects/dashboard/"+dashboardID, nil)
	if err != nil {
		return err
	}

	_, err = c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to delete dashboard: %w", err)
	}

	logrus.WithField("dashboard_id", dashboardID).Debug("Deleted dashboard")
	return nil
}

// CloneDashboard clones a dashboard
func (c *Client) CloneDashboard(ctx context.Context, dashboardID string, newTitle string) (*Dashboard, error) {
	logrus.WithFields(logrus.Fields{
		"dashboard_id": dashboardID,
		"new_title":    newTitle,
	}).Debug("Cloning dashboard")

	// First get the original dashboard
	original, err := c.GetDashboard(ctx, dashboardID)
	if err != nil {
		return nil, err
	}

	// Create new dashboard with same attributes but new title
	attributes := map[string]interface{}{
		"title": newTitle,
	}
	if original.Description != "" {
		attributes["description"] = original.Description
	}
	if original.PanelsJSON != "" {
		attributes["panelsJSON"] = original.PanelsJSON
	}
	if original.OptionsJSON != "" {
		attributes["optionsJSON"] = original.OptionsJSON
	}
	if original.TimeRestore {
		attributes["timeRestore"] = true
	}
	if original.TimeFrom != "" {
		attributes["timeFrom"] = original.TimeFrom
	}
	if original.TimeTo != "" {
		attributes["timeTo"] = original.TimeTo
	}

	obj := SavedObject{
		Type:       "dashboard",
		Attributes: attributes,
	}

	resp, err := c.makeRequest(ctx, "POST", "saved_objects/dashboard", obj)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result SavedObject
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cloned dashboard: %w", err)
	}

	dashboard := &Dashboard{
		ID:         result.ID,
		Title:      newTitle,
		Attributes: result.Attributes,
	}

	logrus.WithField("id", result.ID).Debug("Cloned dashboard")
	return dashboard, nil
}

// ============ Write Operations: Visualizations ============

// CreateVisualization creates a new visualization
func (c *Client) CreateVisualization(ctx context.Context, title string, visState map[string]interface{}, description string, savedSearchRefName string) (*Visualization, error) {
	logrus.WithField("title", title).Debug("Creating visualization")

	attributes := map[string]interface{}{
		"title": title,
	}
	if visState != nil {
		visStateJSON, err := json.Marshal(visState)
		if err == nil {
			attributes["visState"] = string(visStateJSON)
		}
	}
	if description != "" {
		attributes["description"] = description
	}
	if savedSearchRefName != "" {
		attributes["savedSearchRefName"] = savedSearchRefName
	}

	obj := SavedObject{
		Type:       "visualization",
		Attributes: attributes,
	}

	resp, err := c.makeRequest(ctx, "POST", "saved_objects/visualization", obj)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result SavedObject
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal created visualization: %w", err)
	}

	visualization := &Visualization{
		ID:         result.ID,
		Title:      title,
		VisState:   getStringField(result.Attributes, "visState"),
		Attributes: result.Attributes,
	}

	logrus.WithField("id", result.ID).Debug("Created visualization")
	return visualization, nil
}

// UpdateVisualization updates a visualization
func (c *Client) UpdateVisualization(ctx context.Context, visualizationID string, title string, visState map[string]interface{}, description string, version int) (*Visualization, error) {
	logrus.WithField("visualization_id", visualizationID).Debug("Updating visualization")

	attributes := map[string]interface{}{}
	if title != "" {
		attributes["title"] = title
	}
	if visState != nil {
		visStateJSON, err := json.Marshal(visState)
		if err == nil {
			attributes["visState"] = string(visStateJSON)
		}
	}
	if description != "" {
		attributes["description"] = description
	}

	obj := map[string]interface{}{
		"attributes": attributes,
	}
	if version > 0 {
		obj["version"] = version
	}

	resp, err := c.makeRequest(ctx, "PUT", "saved_objects/visualization/"+visualizationID, obj)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result SavedObject
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal updated visualization: %w", err)
	}

	visualization := &Visualization{
		ID:         result.ID,
		Title:      getStringField(result.Attributes, "title"),
		VisState:   getStringField(result.Attributes, "visState"),
		Attributes: result.Attributes,
	}

	logrus.WithField("visualization_id", visualizationID).Debug("Updated visualization")
	return visualization, nil
}

// DeleteVisualization deletes a visualization
func (c *Client) DeleteVisualization(ctx context.Context, visualizationID string) error {
	logrus.WithField("visualization_id", visualizationID).Debug("Deleting visualization")

	resp, err := c.makeRequest(ctx, "DELETE", "saved_objects/visualization/"+visualizationID, nil)
	if err != nil {
		return err
	}

	_, err = c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to delete visualization: %w", err)
	}

	logrus.WithField("visualization_id", visualizationID).Debug("Deleted visualization")
	return nil
}

// CloneVisualization clones a visualization
func (c *Client) CloneVisualization(ctx context.Context, visualizationID string, newTitle string) (*Visualization, error) {
	logrus.WithFields(logrus.Fields{
		"visualization_id": visualizationID,
		"new_title":        newTitle,
	}).Debug("Cloning visualization")

	// First get the original visualization
	original, err := c.GetVisualization(ctx, visualizationID)
	if err != nil {
		return nil, err
	}

	// Create new visualization with same attributes but new title
	attributes := map[string]interface{}{
		"title": newTitle,
	}
	if original.VisState != "" {
		attributes["visState"] = original.VisState
	}
	if original.Description != "" {
		attributes["description"] = original.Description
	}

	obj := SavedObject{
		Type:       "visualization",
		Attributes: attributes,
	}

	resp, err := c.makeRequest(ctx, "POST", "saved_objects/visualization", obj)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result SavedObject
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cloned visualization: %w", err)
	}

	visualization := &Visualization{
		ID:         result.ID,
		Title:      newTitle,
		Attributes: result.Attributes,
	}

	logrus.WithField("id", result.ID).Debug("Cloned visualization")
	return visualization, nil
}

// ============ Write Operations: Saved Objects (Generic) ============

// CreateSavedObject creates a generic saved object
func (c *Client) CreateSavedObject(ctx context.Context, objectType string, attributes map[string]interface{}, references []Reference, initialObjectType string) (*SavedObject, error) {
	logrus.WithFields(logrus.Fields{
		"type": objectType,
	}).Debug("Creating saved object")

	obj := SavedObject{
		Type:       objectType,
		Attributes: attributes,
	}
	if len(references) > 0 {
		obj.References = references
	}
	if initialObjectType != "" {
		obj.Attributes["type"] = initialObjectType
	}

	resp, err := c.makeRequest(ctx, "POST", "saved_objects/"+objectType, obj)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result SavedObject
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal created saved object: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"type": objectType,
		"id":   result.ID,
	}).Debug("Created saved object")
	return &result, nil
}

// UpdateSavedObject updates a saved object
func (c *Client) UpdateSavedObject(ctx context.Context, objectType string, objectID string, attributes map[string]interface{}, references []Reference, version string) (*SavedObject, error) {
	logrus.WithFields(logrus.Fields{
		"type": objectType,
		"id":   objectID,
	}).Debug("Updating saved object")

	obj := map[string]interface{}{
		"attributes": attributes,
	}
	if len(references) > 0 {
		refs := make([]map[string]interface{}, 0, len(references))
		for _, ref := range references {
			refs = append(refs, map[string]interface{}{
				"name": ref.Name,
				"type": ref.Type,
				"id":   ref.ID,
			})
		}
		obj["references"] = refs
	}
	if version != "" {
		obj["version"] = version
	}

	resp, err := c.makeRequest(ctx, "PUT", "saved_objects/"+objectType+"/"+objectID, obj)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result SavedObject
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal updated saved object: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"type": objectType,
		"id":   objectID,
	}).Debug("Updated saved object")
	return &result, nil
}

// DeleteSavedObject deletes a saved object
func (c *Client) DeleteSavedObject(ctx context.Context, objectType string, objectID string, force bool) error {
	logrus.WithFields(logrus.Fields{
		"type":  objectType,
		"id":    objectID,
		"force": force,
	}).Debug("Deleting saved object")

	endpoint := "saved_objects/" + objectType + "/" + objectID
	if force {
		endpoint += "?force=true"
	}

	resp, err := c.makeRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	_, err = c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to delete saved object: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"type": objectType,
		"id":   objectID,
	}).Debug("Deleted saved object")
	return nil
}

// BulkDeleteSavedObjects deletes multiple saved objects
func (c *Client) BulkDeleteSavedObjects(ctx context.Context, objects []SavedObject) error {
	logrus.WithField("count", len(objects)).Debug("Bulk deleting saved objects")

	objectsToDelete := make([]map[string]interface{}, 0, len(objects))
	for _, obj := range objects {
		objectsToDelete = append(objectsToDelete, map[string]interface{}{
			"type": obj.Type,
			"id":   obj.ID,
		})
	}

	resp, err := c.makeRequest(ctx, "POST", "saved_objects/_bulk_delete", map[string]interface{}{
		"objects": objectsToDelete,
	})
	if err != nil {
		return err
	}

	_, err = c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to bulk delete saved objects: %w", err)
	}

	logrus.WithField("count", len(objects)).Debug("Bulk deleted saved objects")
	return nil
}

// ExportSavedObjects exports saved objects
func (c *Client) ExportSavedObjects(ctx context.Context, objects []SavedObject, includeReferences bool) ([]byte, error) {
	logrus.WithField("count", len(objects)).Debug("Exporting saved objects")

	objectsToExport := make([]map[string]interface{}, 0, len(objects))
	for _, obj := range objects {
		objectsToExport = append(objectsToExport, map[string]interface{}{
			"type": obj.Type,
			"id":   obj.ID,
		})
	}

	params := map[string]interface{}{
		"objects": objectsToExport,
	}
	if includeReferences {
		params["includeReferencesDeep"] = true
	}

	resp, err := c.makeRequest(ctx, "POST", "saved_objects/_export", params)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to export saved objects: %w", err)
	}

	logrus.WithField("count", len(objects)).Debug("Exported saved objects")
	return body, nil
}

// ImportSavedObjects imports saved objects from JSON
func (c *Client) ImportSavedObjects(ctx context.Context, fileContent string, createNewCopies bool) error {
	logrus.Debug("Importing saved objects")

	// fileContent can be either base64-encoded or raw JSON
	fileData := []byte(fileContent)

	// Validate it's valid JSON
	var exportData map[string]interface{}
	if err := json.Unmarshal(fileData, &exportData); err != nil {
		return fmt.Errorf("invalid import file format: %w", err)
	}

	params := map[string]interface{}{
		"file": fileContent,
	}
	if createNewCopies {
		params["createNewCopies"] = true
	}

	resp, err := c.makeRequest(ctx, "POST", "saved_objects/_import", params)
	if err != nil {
		return err
	}

	respBody, err := c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to import saved objects: %w", err)
	}

	// Check for errors in import result
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err == nil {
		if errors, ok := result["errors"].([]interface{}); ok && len(errors) > 0 {
			return fmt.Errorf("import completed with errors: %v", errors)
		}
	}

	logrus.Debug("Imported saved objects successfully")
	return nil
}

// ============ Alert Rules ============

// KibanaAlertRule represents a Kibana alerting rule.
type KibanaAlertRule struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	AlertTypeID string                   `json:"alertTypeId"`
	Consumer    string                   `json:"consumer"`
	Schedule    map[string]interface{}   `json:"schedule"`
	Enabled     bool                     `json:"enabled"`
	Actions     []map[string]interface{} `json:"actions,omitempty"`
	Tags        []string                 `json:"tags,omitempty"`
	Params      map[string]interface{}   `json:"params"`
	CreatedAt   string                   `json:"createdAt"`
	UpdatedAt   string                   `json:"updatedAt"`
	NotifyWhen  string                   `json:"notifyWhen,omitempty"`
	Throttle    string                   `json:"throttle,omitempty"`
	MuteAll     bool                     `json:"muteAll"`
}

// KibanaAlertRuleType represents an available alert rule type.
type KibanaAlertRuleType struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Producer    string                 `json:"producer"`
	Description string                 `json:"description"`
	Params      map[string]interface{} `json:"params"`
	Actions     map[string]interface{} `json:"actions"`
}

// KibanaConnector represents a Kibana action connector.
type KibanaConnector struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	ConnectorTypeID string                 `json:"connectorTypeId"`
	Config          map[string]interface{} `json:"config"`
	Secrets         map[string]interface{} `json:"secrets,omitempty"`
	Preconfigured   bool                   `json:"preconfigured,omitempty"`
}

// KibanaConnectorType represents an available connector type.
type KibanaConnectorType struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Inputs      []map[string]interface{} `json:"inputs"`
	Description string                   `json:"description"`
}

// KibanaDataView represents a Kibana data view (formerly index pattern).
type KibDataView struct {
	ID            string                 `json:"id,omitempty"`
	Name          string                 `json:"name"`
	Title         string                 `json:"title"`
	TimeField     string                 `json:"timeFieldName,omitempty"`
	SourceFilters []SourceFilter         `json:"sourceFilters,omitempty"`
	Fields        []IndexPatternField    `json:"fields,omitempty"`
	FieldFormats  map[string]interface{} `json:"fieldFormats,omitempty"`
	Type          string                 `json:"type,omitempty"`
	Origin        string                 `json:"origin,omitempty"`
	AllowNoIndex  bool                   `json:"allowNoIndex,omitempty"`
}

// GetAlertRules retrieves all alert rules with pagination.
func (c *Client) GetAlertRules(ctx context.Context, page, perPage int, filter string, enabled *bool) ([]KibanaAlertRule, error) {
	logrus.WithFields(logrus.Fields{
		"page":    page,
		"perPage": perPage,
		"filter":  filter,
		"enabled": enabled,
	}).Debug("Getting alert rules")

	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	params := url.Values{}
	params.Set("page", fmt.Sprintf("%d", page))
	params.Set("per_page", fmt.Sprintf("%d", perPage))
	if filter != "" {
		params.Set("search", filter)
	}
	if enabled != nil {
		params.Set("enabled", fmt.Sprintf("%t", *enabled))
	}

	endpoint := "api/alerting/rules?" + params.Encode()
	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	respBody, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var rules []KibanaAlertRule
	if err := json.Unmarshal(respBody, &rules); err != nil {
		return nil, fmt.Errorf("failed to unmarshal alert rules: %w", err)
	}

	logrus.WithField("count", len(rules)).Debug("Retrieved alert rules")
	return rules, nil
}

// GetAlertRule retrieves a specific alert rule by ID.
func (c *Client) GetAlertRule(ctx context.Context, ruleID string) (*KibanaAlertRule, error) {
	logrus.WithField("rule_id", ruleID).Debug("Getting alert rule")

	resp, err := c.makeRequest(ctx, "GET", "api/alerting/rules/"+ruleID, nil)
	if err != nil {
		return nil, err
	}

	respBody, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var rule KibanaAlertRule
	if err := json.Unmarshal(respBody, &rule); err != nil {
		return nil, fmt.Errorf("failed to unmarshal alert rule: %w", err)
	}

	logrus.WithField("rule_id", ruleID).Debug("Retrieved alert rule")
	return &rule, nil
}

// CreateAlertRule creates a new alert rule.
func (c *Client) CreateAlertRule(ctx context.Context, name, alertTypeID string, schedule, params map[string]interface{}, actions []map[string]interface{}, tags []string) (*KibanaAlertRule, error) {
	logrus.WithField("name", name).Debug("Creating alert rule")

	rule := map[string]interface{}{
		"name":        name,
		"alertTypeId": alertTypeID,
		"schedule":    schedule,
		"params":      params,
		"actions":     actions,
		"tags":        tags,
		"enabled":     true,
	}

	resp, err := c.makeRequest(ctx, "POST", "api/alerting/rules", rule)
	if err != nil {
		return nil, err
	}

	respBody, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var createdRule KibanaAlertRule
	if err := json.Unmarshal(respBody, &createdRule); err != nil {
		return nil, fmt.Errorf("failed to unmarshal created alert rule: %w", err)
	}

	logrus.WithField("rule_id", createdRule.ID).Debug("Created alert rule")
	return &createdRule, nil
}

// UpdateAlertRule updates an existing alert rule.
func (c *Client) UpdateAlertRule(ctx context.Context, ruleID string, name, schedule string, params, actions map[string]interface{}, tags []string) (*KibanaAlertRule, error) {
	logrus.WithField("rule_id", ruleID).Debug("Updating alert rule")

	rule := map[string]interface{}{}
	if name != "" {
		rule["name"] = name
	}
	if schedule != "" {
		// Parse schedule if it's a string
		var scheduleMap map[string]interface{}
		if err := json.Unmarshal([]byte(schedule), &scheduleMap); err == nil {
			rule["schedule"] = scheduleMap
		} else {
			rule["schedule"] = map[string]interface{}{"interval": schedule}
		}
	}
	if params != nil {
		rule["params"] = params
	}
	if actions != nil {
		rule["actions"] = actions
	}
	if tags != nil {
		rule["tags"] = tags
	}

	resp, err := c.makeRequest(ctx, "PUT", "api/alerting/rules/"+ruleID, rule)
	if err != nil {
		return nil, err
	}

	respBody, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var updatedRule KibanaAlertRule
	if err := json.Unmarshal(respBody, &updatedRule); err != nil {
		return nil, fmt.Errorf("failed to unmarshal updated alert rule: %w", err)
	}

	logrus.WithField("rule_id", ruleID).Debug("Updated alert rule")
	return &updatedRule, nil
}

// DeleteAlertRule deletes an alert rule.
func (c *Client) DeleteAlertRule(ctx context.Context, ruleID string) error {
	logrus.WithField("rule_id", ruleID).Debug("Deleting alert rule")

	resp, err := c.makeRequest(ctx, "DELETE", "api/alerting/rules/"+ruleID, nil)
	if err != nil {
		return err
	}

	_, err = c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to delete alert rule: %w", err)
	}

	logrus.WithField("rule_id", ruleID).Debug("Deleted alert rule")
	return nil
}

// EnableAlertRule enables a disabled alert rule.
func (c *Client) EnableAlertRule(ctx context.Context, ruleID string) error {
	logrus.WithField("rule_id", ruleID).Debug("Enabling alert rule")

	resp, err := c.makeRequest(ctx, "POST", "api/alerting/rules/"+ruleID+"/_enable", nil)
	if err != nil {
		return err
	}

	_, err = c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to enable alert rule: %w", err)
	}

	logrus.WithField("rule_id", ruleID).Debug("Enabled alert rule")
	return nil
}

// DisableAlertRule disables an alert rule.
func (c *Client) DisableAlertRule(ctx context.Context, ruleID string) error {
	logrus.WithField("rule_id", ruleID).Debug("Disabling alert rule")

	resp, err := c.makeRequest(ctx, "POST", "api/alerting/rules/"+ruleID+"/_disable", nil)
	if err != nil {
		return err
	}

	_, err = c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to disable alert rule: %w", err)
	}

	logrus.WithField("rule_id", ruleID).Debug("Disabled alert rule")
	return nil
}

// MuteAlertRule mutes an alert rule for a specified duration.
func (c *Client) MuteAlertRule(ctx context.Context, ruleID, duration string) error {
	logrus.WithFields(logrus.Fields{
		"rule_id":  ruleID,
		"duration": duration,
	}).Debug("Muting alert rule")

	resp, err := c.makeRequest(ctx, "POST", "api/alerting/rules/"+ruleID+"/_mute", map[string]interface{}{
		"duration": duration,
	})
	if err != nil {
		return err
	}

	_, err = c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to mute alert rule: %w", err)
	}

	logrus.WithField("rule_id", ruleID).Debug("Muted alert rule")
	return nil
}

// UnmuteAlertRule unmutes a previously muted alert rule.
func (c *Client) UnmuteAlertRule(ctx context.Context, ruleID string) error {
	logrus.WithField("rule_id", ruleID).Debug("Unmuting alert rule")

	resp, err := c.makeRequest(ctx, "POST", "api/alerting/rules/"+ruleID+"/_unmute", nil)
	if err != nil {
		return err
	}

	_, err = c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to unmute alert rule: %w", err)
	}

	logrus.WithField("rule_id", ruleID).Debug("Unmuted alert rule")
	return nil
}

// GetAlertRuleTypes retrieves all available alert rule types.
func (c *Client) GetAlertRuleTypes(ctx context.Context) ([]KibanaAlertRuleType, error) {
	logrus.Debug("Getting alert rule types")

	resp, err := c.makeRequest(ctx, "GET", "api/alerting/rule_types", nil)
	if err != nil {
		return nil, err
	}

	respBody, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var ruleTypes []KibanaAlertRuleType
	if err := json.Unmarshal(respBody, &ruleTypes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal alert rule types: %w", err)
	}

	logrus.WithField("count", len(ruleTypes)).Debug("Retrieved alert rule types")
	return ruleTypes, nil
}

// GetAlertRuleHistory retrieves execution history for an alert rule.
func (c *Client) GetAlertRuleHistory(ctx context.Context, ruleID string, page, perPage int) ([]map[string]interface{}, error) {
	logrus.WithFields(logrus.Fields{
		"rule_id": ruleID,
		"page":    page,
		"perPage": perPage,
	}).Debug("Getting alert rule history")

	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	params := url.Values{}
	params.Set("page", fmt.Sprintf("%d", page))
	params.Set("per_page", fmt.Sprintf("%d", perPage))

	endpoint := "api/alerting/rules/" + ruleID + "/execution?" + params.Encode()
	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	respBody, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var history []map[string]interface{}
	if err := json.Unmarshal(respBody, &history); err != nil {
		return nil, fmt.Errorf("failed to unmarshal alert rule history: %w", err)
	}

	logrus.WithField("count", len(history)).Debug("Retrieved alert rule history")
	return history, nil
}

// ============ Connectors ============

// GetConnectors retrieves all action connectors.
func (c *Client) GetConnectors(ctx context.Context, page, perPage int) ([]KibanaConnector, error) {
	logrus.WithFields(logrus.Fields{
		"page":    page,
		"perPage": perPage,
	}).Debug("Getting connectors")

	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	params := url.Values{}
	params.Set("page", fmt.Sprintf("%d", page))
	params.Set("per_page", fmt.Sprintf("%d", perPage))

	endpoint := "api/actions/connector?" + params.Encode()
	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	respBody, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var connectors []KibanaConnector
	if err := json.Unmarshal(respBody, &connectors); err != nil {
		return nil, fmt.Errorf("failed to unmarshal connectors: %w", err)
	}

	logrus.WithField("count", len(connectors)).Debug("Retrieved connectors")
	return connectors, nil
}

// GetConnector retrieves a specific connector by ID.
func (c *Client) GetConnector(ctx context.Context, connectorID string) (*KibanaConnector, error) {
	logrus.WithField("connector_id", connectorID).Debug("Getting connector")

	resp, err := c.makeRequest(ctx, "GET", "api/actions/connector/"+connectorID, nil)
	if err != nil {
		return nil, err
	}

	respBody, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var connector KibanaConnector
	if err := json.Unmarshal(respBody, &connector); err != nil {
		return nil, fmt.Errorf("failed to unmarshal connector: %w", err)
	}

	logrus.WithField("connector_id", connectorID).Debug("Retrieved connector")
	return &connector, nil
}

// CreateConnector creates a new action connector.
func (c *Client) CreateConnector(ctx context.Context, name, connectorTypeID string, config, secrets map[string]interface{}) (*KibanaConnector, error) {
	logrus.WithField("name", name).Debug("Creating connector")

	connector := map[string]interface{}{
		"name":            name,
		"connectorTypeId": connectorTypeID,
		"config":          config,
	}
	if secrets != nil {
		connector["secrets"] = secrets
	}

	resp, err := c.makeRequest(ctx, "POST", "api/actions/connector", connector)
	if err != nil {
		return nil, err
	}

	respBody, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var createdConnector KibanaConnector
	if err := json.Unmarshal(respBody, &createdConnector); err != nil {
		return nil, fmt.Errorf("failed to unmarshal created connector: %w", err)
	}

	logrus.WithField("connector_id", createdConnector.ID).Debug("Created connector")
	return &createdConnector, nil
}

// UpdateConnector updates an existing connector.
func (c *Client) UpdateConnector(ctx context.Context, connectorID string, name string, config, secrets map[string]interface{}) (*KibanaConnector, error) {
	logrus.WithField("connector_id", connectorID).Debug("Updating connector")

	connector := map[string]interface{}{}
	if name != "" {
		connector["name"] = name
	}
	if config != nil {
		connector["config"] = config
	}
	if secrets != nil {
		connector["secrets"] = secrets
	}

	resp, err := c.makeRequest(ctx, "PUT", "api/actions/connector/"+connectorID, connector)
	if err != nil {
		return nil, err
	}

	respBody, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var updatedConnector KibanaConnector
	if err := json.Unmarshal(respBody, &updatedConnector); err != nil {
		return nil, fmt.Errorf("failed to unmarshal updated connector: %w", err)
	}

	logrus.WithField("connector_id", connectorID).Debug("Updated connector")
	return &updatedConnector, nil
}

// DeleteConnector deletes a connector.
func (c *Client) DeleteConnector(ctx context.Context, connectorID string) error {
	logrus.WithField("connector_id", connectorID).Debug("Deleting connector")

	resp, err := c.makeRequest(ctx, "DELETE", "api/actions/connector/"+connectorID, nil)
	if err != nil {
		return err
	}

	_, err = c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to delete connector: %w", err)
	}

	logrus.WithField("connector_id", connectorID).Debug("Deleted connector")
	return nil
}

// TestConnector tests a connector by sending a test notification.
func (c *Client) TestConnector(ctx context.Context, connectorID string, body map[string]interface{}) error {
	logrus.WithField("connector_id", connectorID).Debug("Testing connector")

	resp, err := c.makeRequest(ctx, "POST", "api/actions/connector/"+connectorID+"/_test", body)
	if err != nil {
		return err
	}

	_, err = c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to test connector: %w", err)
	}

	logrus.WithField("connector_id", connectorID).Debug("Tested connector")
	return nil
}

// GetConnectorTypes retrieves all available connector types.
func (c *Client) GetConnectorTypes(ctx context.Context) ([]KibanaConnectorType, error) {
	logrus.Debug("Getting connector types")

	resp, err := c.makeRequest(ctx, "GET", "api/actions/connector_types", nil)
	if err != nil {
		return nil, err
	}

	respBody, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var connectorTypes []KibanaConnectorType
	if err := json.Unmarshal(respBody, &connectorTypes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal connector types: %w", err)
	}

	logrus.WithField("count", len(connectorTypes)).Debug("Retrieved connector types")
	return connectorTypes, nil
}

// ============ Data Views ============

// GetDataViews retrieves all data views.
func (c *Client) GetDataViews(ctx context.Context, page, perPage int) ([]KibDataView, error) {
	logrus.WithFields(logrus.Fields{
		"page":    page,
		"perPage": perPage,
	}).Debug("Getting data views")

	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	params := url.Values{}
	params.Set("page", fmt.Sprintf("%d", page))
	params.Set("per_page", fmt.Sprintf("%d", perPage))

	endpoint := "api/data_views?" + params.Encode()
	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	respBody, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result struct {
		DataViews []KibDataView `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data views: %w", err)
	}

	logrus.WithField("count", len(result.DataViews)).Debug("Retrieved data views")
	return result.DataViews, nil
}

// GetDataView retrieves a specific data view by ID.
func (c *Client) GetDataView(ctx context.Context, dataViewID string) (*KibDataView, error) {
	logrus.WithField("data_view_id", dataViewID).Debug("Getting data view")

	resp, err := c.makeRequest(ctx, "GET", "api/data_views/"+dataViewID, nil)
	if err != nil {
		return nil, err
	}

	respBody, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var dataView KibDataView
	if err := json.Unmarshal(respBody, &dataView); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data view: %w", err)
	}

	logrus.WithField("data_view_id", dataViewID).Debug("Retrieved data view")
	return &dataView, nil
}

// CreateDataView creates a new data view.
func (c *Client) CreateDataView(ctx context.Context, title, name, timeField string, sourceFilters []SourceFilter, fieldFormats map[string]interface{}, allowNoIndex bool) (*KibDataView, error) {
	logrus.WithField("title", title).Debug("Creating data view")

	dataView := map[string]interface{}{
		"title": title,
	}
	if name != "" {
		dataView["name"] = name
	}
	if timeField != "" {
		dataView["timeFieldName"] = timeField
	}
	if sourceFilters != nil {
		sourceFiltersJSON, err := json.Marshal(sourceFilters)
		if err == nil {
			dataView["sourceFilters"] = string(sourceFiltersJSON)
		}
	}
	if fieldFormats != nil {
		dataView["fieldFormats"] = fieldFormats
	}
	if allowNoIndex {
		dataView["allowNoIndex"] = true
	}

	resp, err := c.makeRequest(ctx, "POST", "api/data_views", dataView)
	if err != nil {
		return nil, err
	}

	respBody, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var createdDataView KibDataView
	if err := json.Unmarshal(respBody, &createdDataView); err != nil {
		return nil, fmt.Errorf("failed to unmarshal created data view: %w", err)
	}

	logrus.WithField("data_view_id", createdDataView.ID).Debug("Created data view")
	return &createdDataView, nil
}

// UpdateDataView updates an existing data view.
func (c *Client) UpdateDataView(ctx context.Context, dataViewID string, title, name, timeField string) (*KibDataView, error) {
	logrus.WithField("data_view_id", dataViewID).Debug("Updating data view")

	dataView := map[string]interface{}{}
	if title != "" {
		dataView["title"] = title
	}
	if name != "" {
		dataView["name"] = name
	}
	if timeField != "" {
		dataView["timeFieldName"] = timeField
	}

	resp, err := c.makeRequest(ctx, "PUT", "api/data_views/"+dataViewID, dataView)
	if err != nil {
		return nil, err
	}

	respBody, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var updatedDataView KibDataView
	if err := json.Unmarshal(respBody, &updatedDataView); err != nil {
		return nil, fmt.Errorf("failed to unmarshal updated data view: %w", err)
	}

	logrus.WithField("data_view_id", dataViewID).Debug("Updated data view")
	return &updatedDataView, nil
}

// DeleteDataView deletes a data view.
func (c *Client) DeleteDataView(ctx context.Context, dataViewID string) error {
	logrus.WithField("data_view_id", dataViewID).Debug("Deleting data view")

	resp, err := c.makeRequest(ctx, "DELETE", "api/data_views/"+dataViewID, nil)
	if err != nil {
		return err
	}

	_, err = c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to delete data view: %w", err)
	}

	logrus.WithField("data_view_id", dataViewID).Debug("Deleted data view")
	return nil
}

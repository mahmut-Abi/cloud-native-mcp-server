// Package client provides HTTP client for Alertmanager API operations.
// It handles authentication, request/response processing, and connection management.
package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/mahmut-Abi/k8s-mcp-server/internal/constants"
	optimize "github.com/mahmut-Abi/k8s-mcp-server/internal/util/performance"
	"github.com/sirupsen/logrus"
)

var logger = logrus.WithField("component", "alertmanager-client")

// ClientOptions configures the Alertmanager HTTP client
type ClientOptions struct {
	Address       string        // Alertmanager server address
	Timeout       time.Duration // Request timeout
	Username      string        // Basic auth username
	Password      string        // Basic auth password
	BearerToken   string        // Bearer token for authentication
	TLSSkipVerify bool          // Skip TLS certificate verification
	TLSCertFile   string        // TLS certificate file
	TLSKeyFile    string        // TLS key file
	TLSCAFile     string        // TLS CA file
}

// DefaultClientOptions returns default client options
func DefaultClientOptions() *ClientOptions {
	return &ClientOptions{
		Address: "http://localhost:9093",
		Timeout: 30 * time.Second,
	}
}

// Client provides HTTP client for Alertmanager API operations
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	username   string
	password   string
	token      string
}

// NewClient creates a new Alertmanager client with default options
func NewClient() (*Client, error) {
	return NewClientWithOptions(DefaultClientOptions())
}

// NewClientWithOptions creates a new Alertmanager client with custom options
func NewClientWithOptions(opts *ClientOptions) (*Client, error) {
	if opts == nil {
		opts = DefaultClientOptions()
	}

	// Parse base URL
	baseURL, err := url.Parse(strings.TrimSuffix(opts.Address, "/"))
	if err != nil {
		return nil, fmt.Errorf("invalid Alertmanager address: %w", err)
	}

	// Configure TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: opts.TLSSkipVerify,
	}

	// Load TLS certificates if specified
	if opts.TLSCertFile != "" && opts.TLSKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(opts.TLSCertFile, opts.TLSKeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS certificates: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	// Create optimized HTTP client with TLS configuration
	httpClient := optimize.NewOptimizedHTTPClientWithTimeout(opts.Timeout)
	if transport, ok := httpClient.Transport.(*http.Transport); ok {
		transport.TLSClientConfig = tlsConfig
	}

	client := &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
		username:   opts.Username,
		password:   opts.Password,
		token:      opts.BearerToken,
	}

	return client, nil
}

// buildURL constructs a full URL for an API endpoint
func (c *Client) buildURL(endpoint string) string {
	u := *c.baseURL
	u.Path = path.Join(u.Path, "api/v2", endpoint)
	return u.String()
}

// doRequest executes an HTTP request with authentication
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.buildURL(endpoint), reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add authentication
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	} else if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	logger.WithFields(logrus.Fields{
		"method":   method,
		"endpoint": endpoint,
		"url":      req.URL.String(),
	}).Debug("Making HTTP request")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// parseResponse parses HTTP response into target struct
func (c *Client) parseResponse(resp *http.Response, target interface{}) error {
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	if target == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

// GetStatus retrieves Alertmanager status information
func (c *Client) GetStatus(ctx context.Context) (map[string]interface{}, error) {
	resp, err := c.doRequest(ctx, "GET", "status", nil)
	if err != nil {
		return nil, err
	}

	var status map[string]interface{}
	if err := c.parseResponse(resp, &status); err != nil {
		return nil, err
	}

	return status, nil
}

// GetAlerts retrieves current alerts with optional filters
func (c *Client) GetAlerts(ctx context.Context, filters map[string]string) ([]map[string]interface{}, error) {
	endpoint := "alerts"

	// Add query parameters if filters are provided
	if len(filters) > 0 {
		params := url.Values{}
		for key, value := range filters {
			params.Add(key, value)
		}
		endpoint += "?" + params.Encode()
	}

	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var alerts []map[string]interface{}
	if err := c.parseResponse(resp, &alerts); err != nil {
		return nil, err
	}

	return alerts, nil
}

// GetAlertGroups retrieves alert groups
func (c *Client) GetAlertGroups(ctx context.Context) ([]map[string]interface{}, error) {
	resp, err := c.doRequest(ctx, "GET", "alerts/groups", nil)
	if err != nil {
		return nil, err
	}

	var groups []map[string]interface{}
	if err := c.parseResponse(resp, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

// GetSilences retrieves current silences
func (c *Client) GetSilences(ctx context.Context) ([]map[string]interface{}, error) {
	resp, err := c.doRequest(ctx, "GET", "silences", nil)
	if err != nil {
		return nil, err
	}

	var silences []map[string]interface{}
	if err := c.parseResponse(resp, &silences); err != nil {
		return nil, err
	}

	return silences, nil
}

// CreateSilence creates a new silence
func (c *Client) CreateSilence(ctx context.Context, silence map[string]interface{}) (map[string]interface{}, error) {
	resp, err := c.doRequest(ctx, "POST", "silences", silence)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteSilence deletes a silence by ID
func (c *Client) DeleteSilence(ctx context.Context, silenceID string) error {
	endpoint := fmt.Sprintf("silence/%s", silenceID)
	resp, err := c.doRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	return c.parseResponse(resp, nil)
}

// GetReceivers retrieves configured receivers
func (c *Client) GetReceivers(ctx context.Context) ([]map[string]interface{}, error) {
	resp, err := c.doRequest(ctx, "GET", "receivers", nil)
	if err != nil {
		return nil, err
	}

	var receivers []map[string]interface{}
	if err := c.parseResponse(resp, &receivers); err != nil {
		return nil, err
	}

	return receivers, nil
}

// TestReceiver tests a receiver configuration
func (c *Client) TestReceiver(ctx context.Context, receiver map[string]interface{}) (map[string]interface{}, error) {
	resp, err := c.doRequest(ctx, "POST", "receivers/test", receiver)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Response size limits for Alertmanager - similar to other optimizations
// Use constants from internal/constants package for consistency
const (
	defaultLimit      = constants.DefaultPageSize
	maxLimit          = constants.MaxPageSize
	warningLimit      = constants.WarningPageSize
	defaultLimitNodes = constants.DefaultPageSizeNodes
)

// PaginationInfo represents pagination metadata for Alertmanager responses
type PaginationInfo struct {
	CurrentPage     int   `json:"currentPage"`
	PerPage         int   `json:"perPage"`
	TotalCount      int64 `json:"totalCount"`
	TotalPages      int   `json:"totalPages"`
	HasNextPage     bool  `json:"hasNextPage"`
	HasPreviousPage bool  `json:"hasPreviousPage"`
}

// AlertsSummary returns optimized alerts information
func (c *Client) AlertsSummary(ctx context.Context, filter string, receiver string, silenced *bool, activeOnly *bool, limit int) ([]map[string]interface{}, error) {
	// Build filters
	filters := make(map[string]string)

	if filter != "" {
		filters["filter"] = filter
	}
	if receiver != "" {
		filters["receiver"] = receiver
	}
	if silenced != nil {
		filters["silenced"] = fmt.Sprintf("%t", *silenced)
	}
	if activeOnly != nil && *activeOnly {
		filters["active"] = "true"
	}

	// Get alerts
	alerts, err := c.GetAlerts(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Apply limit
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	var summaries []map[string]interface{}
	count := 0
	for _, alert := range alerts {
		if count >= limit {
			break
		}

		summary := map[string]interface{}{
			"labels":      alert["labels"],
			"annotations": alert["annotations"],
			"startsAt":    alert["startsAt"],
			"endsAt":      alert["endsAt"],
			"status":      alert["status"],
			"fingerprint": alert["fingerprint"],
		}
		summaries = append(summaries, summary)
		count++
	}

	return summaries, nil
}

// SilencesSummary returns optimized silences information
func (c *Client) SilencesSummary(ctx context.Context, status string, limit int) ([]map[string]interface{}, error) {
	silences, err := c.GetSilences(ctx)
	if err != nil {
		return nil, err
	}

	// Apply limit
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	var summaries []map[string]interface{}
	count := 0
	for _, silence := range silences {
		if count >= limit {
			break
		}

		// Filter by status if specified
		if status != "" {
			if silenceStatus, ok := silence["status"].(map[string]interface{}); ok {
				if state, ok := silenceStatus["state"].(string); ok && state != status {
					continue
				}
			}
		}

		summary := map[string]interface{}{
			"id":        silence["id"],
			"status":    silence["status"],
			"matchers":  silence["matchers"],
			"startsAt":  silence["startsAt"],
			"endsAt":    silence["endsAt"],
			"createdBy": silence["createdBy"],
			"comment":   silence["comment"],
		}
		summaries = append(summaries, summary)
		count++
	}

	return summaries, nil
}

// AlertGroupsPaginated returns alert groups with pagination support and optimization
func (c *Client) AlertGroupsPaginated(ctx context.Context, page, perPage int, receiver string, activeOnly *bool, sortBy string) ([]map[string]interface{}, *PaginationInfo, error) {
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

	groups, err := c.GetAlertGroups(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Filter groups
	var filteredGroups []map[string]interface{}
	for _, group := range groups {
		// Filter by receiver if specified
		if receiver != "" {
			if groupReceiver, ok := group["receiver"].(string); ok && groupReceiver != receiver {
				continue
			}
		}

		// Filter by active status if specified
		if activeOnly != nil && *activeOnly {
			if blocks, ok := group["blocks"].([]interface{}); ok {
				hasActive := false
				for _, block := range blocks {
					if blockMap, ok := block.(map[string]interface{}); ok {
						if blockStatus, ok := blockMap["status"].(map[string]interface{}); ok {
							if state, ok := blockStatus["state"].(string); ok && state == "active" {
								hasActive = true
								break
							}
						}
					}
				}
				if !hasActive {
					continue
				}
			}
		}

		filteredGroups = append(filteredGroups, group)
	}

	// Note: Sorting implementation can be added here in the future if needed
	// Currently returning filtered groups as-is
	_ = sortBy // Acknowledge that sortBy is intentionally not used yet

	// Apply pagination
	total := len(filteredGroups)
	start := (page - 1) * perPage
	end := start + perPage

	if start >= total {
		return []map[string]interface{}{}, &PaginationInfo{
			CurrentPage:     page,
			PerPage:         perPage,
			TotalCount:      int64(total),
			TotalPages:      (total + perPage - 1) / perPage,
			HasNextPage:     false,
			HasPreviousPage: page > 1,
		}, nil
	}

	if end > total {
		end = total
	}

	paginatedGroups := filteredGroups[start:end]

	// Create optimized summaries
	var summaries []map[string]interface{}
	for _, group := range paginatedGroups {
		summary := map[string]interface{}{
			"id":       group["id"],
			"receiver": group["receiver"],
			"labels":   group["labels"],
		}

		// Include block count instead of full blocks
		if blocks, ok := group["blocks"].([]interface{}); ok {
			summary["blockCount"] = len(blocks)
			// Include active block count
			activeCount := 0
			for _, block := range blocks {
				if blockMap, ok := block.(map[string]interface{}); ok {
					if blockStatus, ok := blockMap["status"].(map[string]interface{}); ok {
						if state, ok := blockStatus["state"].(string); ok && state == "active" {
							activeCount++
						}
					}
				}
			}
			summary["activeBlockCount"] = activeCount
		}

		summaries = append(summaries, summary)
	}

	// Calculate pagination info
	totalPages := (total + perPage - 1) / perPage
	pagination := &PaginationInfo{
		CurrentPage:     page,
		PerPage:         perPage,
		TotalCount:      int64(total),
		TotalPages:      totalPages,
		HasNextPage:     page < totalPages,
		HasPreviousPage: page > 1,
	}

	return summaries, pagination, nil
}

// SilencesPaginated returns silences with pagination support and optimization
func (c *Client) SilencesPaginated(ctx context.Context, page, perPage int, status, createdBy, commentFilter string) ([]map[string]interface{}, *PaginationInfo, error) {
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

	silences, err := c.GetSilences(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Filter silences
	var filteredSilences []map[string]interface{}
	for _, silence := range silences {
		// Filter by status if specified
		if status != "" {
			if silenceStatus, ok := silence["status"].(map[string]interface{}); ok {
				if state, ok := silenceStatus["state"].(string); ok && state != status {
					continue
				}
			}
		}

		// Filter by creator if specified
		if createdBy != "" {
			if silenceCreator, ok := silence["createdBy"].(string); ok && silenceCreator != createdBy {
				continue
			}
		}

		// Filter by comment content if specified
		if commentFilter != "" {
			if comment, ok := silence["comment"].(string); ok && !strings.Contains(strings.ToLower(comment), strings.ToLower(commentFilter)) {
				continue
			}
		}

		filteredSilences = append(filteredSilences, silence)
	}

	// Apply pagination
	total := len(filteredSilences)
	start := (page - 1) * perPage
	end := start + perPage

	if start >= total {
		return []map[string]interface{}{}, &PaginationInfo{
			CurrentPage:     page,
			PerPage:         perPage,
			TotalCount:      int64(total),
			TotalPages:      (total + perPage - 1) / perPage,
			HasNextPage:     false,
			HasPreviousPage: page > 1,
		}, nil
	}

	if end > total {
		end = total
	}

	paginatedSilences := filteredSilences[start:end]

	// Create optimized summaries
	var summaries []map[string]interface{}
	for _, silence := range paginatedSilences {
		summary := map[string]interface{}{
			"id":        silence["id"],
			"status":    silence["status"],
			"startsAt":  silence["startsAt"],
			"endsAt":    silence["endsAt"],
			"createdBy": silence["createdBy"],
			"comment":   silence["comment"],
		}

		// Include matchers count instead of full matchers
		if matchers, ok := silence["matchers"].([]interface{}); ok {
			summary["matchersCount"] = len(matchers)
		}

		summaries = append(summaries, summary)
	}

	// Calculate pagination info
	totalPages := (total + perPage - 1) / perPage
	pagination := &PaginationInfo{
		CurrentPage:     page,
		PerPage:         perPage,
		TotalCount:      int64(total),
		TotalPages:      totalPages,
		HasNextPage:     page < totalPages,
		HasPreviousPage: page > 1,
	}

	return summaries, pagination, nil
}

// ReceiversSummary returns optimized receivers information
func (c *Client) ReceiversSummary(ctx context.Context, testInfo *bool) ([]map[string]interface{}, error) {
	receivers, err := c.GetReceivers(ctx)
	if err != nil {
		return nil, err
	}

	var summaries []map[string]interface{}
	for _, receiver := range receivers {
		summary := map[string]interface{}{
			"name":   receiver["name"],
			"status": "configured", // Alertmanager receivers are always configured
		}

		// Include test information if requested
		includeTestInfo := false
		if testInfo != nil {
			includeTestInfo = *testInfo
		}

		if includeTestInfo {
			if settings, ok := receiver["settings"]; ok {
				summary["settingsType"] = fmt.Sprintf("%T", settings)
			}
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}

// QueryAlertsAdvanced provides advanced alert querying with enhanced filters
func (c *Client) QueryAlertsAdvanced(ctx context.Context, filter, receiver string, silenced, active, inhibited *bool, timeRange string, page, perPage int, sortBy, sortOrder string, includeLabels *bool) ([]map[string]interface{}, *PaginationInfo, error) {
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
	// Note: Sorting functionality not yet implemented, but parameters validated for future use
	_ = sortBy    // Will be used when sorting is implemented
	_ = sortOrder // Will be used when sorting is implemented

	// Build filters
	filters := make(map[string]string)

	if filter != "" {
		filters["filter"] = filter
	}
	if receiver != "" {
		filters["receiver"] = receiver
	}
	if silenced != nil {
		filters["silenced"] = fmt.Sprintf("%t", *silenced)
	}
	if active != nil {
		filters["active"] = fmt.Sprintf("%t", *active)
	}
	if inhibited != nil {
		filters["inhibited"] = fmt.Sprintf("%t", *inhibited)
	}

	alerts, err := c.GetAlerts(ctx, filters)
	if err != nil {
		return nil, nil, err
	}

	// Apply pagination
	total := len(alerts)
	start := (page - 1) * perPage
	end := start + perPage

	if start >= total {
		return []map[string]interface{}{}, &PaginationInfo{
			CurrentPage:     page,
			PerPage:         perPage,
			TotalCount:      int64(total),
			TotalPages:      (total + perPage - 1) / perPage,
			HasNextPage:     false,
			HasPreviousPage: page > 1,
		}, nil
	}

	if end > total {
		end = total
	}

	paginatedAlerts := alerts[start:end]

	// Create optimized summaries
	var summaries []map[string]interface{}
	for _, alert := range paginatedAlerts {
		summary := map[string]interface{}{
			"startsAt":    alert["startsAt"],
			"endsAt":      alert["endsAt"],
			"status":      alert["status"],
			"fingerprint": alert["fingerprint"],
		}

		// Include labels if requested
		includeLabelsVal := false
		if includeLabels != nil {
			includeLabelsVal = *includeLabels
		}

		if includeLabelsVal {
			summary["labels"] = alert["labels"]
		} else {
			// Include only essential labels
			if labels, ok := alert["labels"].(map[string]interface{}); ok {
				essentialLabels := make(map[string]interface{})
				for key, value := range labels {
					if key == "alertname" || key == "severity" || key == "instance" {
						essentialLabels[key] = value
					}
				}
				summary["labels"] = essentialLabels
			}
		}

		// Include annotations summary
		if annotations, ok := alert["annotations"].(map[string]interface{}); ok {
			summary["annotationsCount"] = len(annotations)
			if summaryVal, ok := annotations["summary"].(string); ok {
				summary["summary"] = summaryVal
			}
		}

		summaries = append(summaries, summary)
	}

	// Calculate pagination info
	totalPages := (total + perPage - 1) / perPage
	pagination := &PaginationInfo{
		CurrentPage:     page,
		PerPage:         perPage,
		TotalCount:      int64(total),
		TotalPages:      totalPages,
		HasNextPage:     page < totalPages,
		HasPreviousPage: page > 1,
	}

	return summaries, pagination, nil
}

// GetHealthSummary returns lightweight health information
func (c *Client) GetHealthSummary(ctx context.Context, level string, includeCluster *bool) (map[string]interface{}, error) {
	status, err := c.GetStatus(ctx)
	if err != nil {
		return nil, err
	}

	summary := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"metadata": map[string]interface{}{
			"level":          level,
			"includeCluster": includeCluster,
			"optimizedFor":   "monitoring and LLM efficiency",
		},
	}

	// Include basic status information
	if statusInfo, ok := status["status"].(map[string]interface{}); ok {
		if state, ok := statusInfo["state"].(string); ok {
			summary["state"] = state
		}
	}

	// Include version information
	if version, ok := status["version"]; ok {
		summary["version"] = version
	}

	// Include cluster information if requested
	includeClusterVal := false
	if includeCluster != nil {
		includeClusterVal = *includeCluster
	}

	if includeClusterVal && level == "detailed" {
		if clusterInfo, ok := status["cluster"]; ok {
			summary["cluster"] = clusterInfo
		}
	}

	// Include detailed status for higher levels
	if level == "detailed" || level == "metrics" {
		summary["fullStatus"] = status
	}

	return summary, nil
}

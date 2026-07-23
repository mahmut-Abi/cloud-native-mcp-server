// Package client provides VictoriaMetrics HTTP API client functionality.
// It offers operations for querying VictoriaMetrics metrics, exporting/importing data,
// and managing the VictoriaMetrics instance through its HTTP API.
package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
	"github.com/sirupsen/logrus"
)

// ClientOptions holds configuration parameters for creating a VictoriaMetrics client.
type ClientOptions struct {
	Address        string        // VictoriaMetrics server address
	Username       string        // Username for basic authentication
	Password       string        // Password for basic authentication
	BearerToken    string        // Bearer token for authentication
	Timeout        time.Duration // HTTP request timeout
	TLSSkipVerify  bool          // Skip TLS certificate verification
	TLSCertFile    string        // TLS certificate file
	TLSKeyFile     string        // TLS key file
	TLSCAFile      string        // TLS CA file
	MaxRetries     int           // Retries for idempotent requests
	RetryBaseDelay time.Duration // Base delay for exponential backoff
	RetryMaxDelay  time.Duration // Maximum delay between retries
}

// Client provides operations for interacting with VictoriaMetrics API.
type Client struct {
	baseURL        string            // Base URL for VictoriaMetrics API
	vmURL          string            // Root URL (without /api/v1/) for VM-specific endpoints
	httpClient     *http.Client      // HTTP client for API requests
	username       string            // Username for basic auth
	password       string            // Password for basic auth
	bearerToken    string            // Bearer token for auth
	headers        map[string]string // Additional headers
	maxRetries     int
	retryBaseDelay time.Duration
	retryMaxDelay  time.Duration
}

// QueryResult represents a VictoriaMetrics query result.
type QueryResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string      `json:"resultType"`
		Result     interface{} `json:"result"`
	} `json:"data"`
	ErrorType string `json:"errorType,omitempty"`
	Error     string `json:"error,omitempty"`
}

// ExportResult represents a VictoriaMetrics export data entry.
type ExportResult struct {
	Metric     map[string]string `json:"metric"`
	Values     []float64         `json:"values"`
	Timestamps []int64           `json:"timestamps"`
}

// TSDBStats represents VictoriaMetrics TSDB statistics.
type TSDBStats struct {
	Status string `json:"status"`
	Data   struct {
		SeriesCount         uint64 `json:"totalSeries"`
		SeriesCountPerHour  uint64 `json:"seriesCountByMetricName,omitempty"`
		ChunksCount         uint64 `json:"chunksCount,omitempty"`
		IndexSizeBytes      uint64 `json:"indexSizeBytes,omitempty"`
		DiskUsageBytes      uint64 `json:"diskUsageBytes,omitempty"`
		SlowQueriesCount    uint64 `json:"slowQueriesCount,omitempty"`
		ActiveQueriesCount  uint64 `json:"activeQueriesCount,omitempty"`
	} `json:"data"`
}

// ActiveQuery represents a currently running query in VictoriaMetrics.
type ActiveQuery struct {
	Status string `json:"status"`
	Data   []struct {
		Query    string  `json:"query"`
		Duration float64 `json:"duration"`
		ID       string  `json:"id"`
		Remote   string  `json:"remote_addr"`
		Start    string  `json:"start"`
		End      string  `json:"end"`
		Step     string  `json:"step"`
	} `json:"data"`
}

// LabelValuesResult represents a label values response.
type LabelValuesResult struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

// MetadataResult represents metric metadata.
type MetadataResult struct {
	Status string              `json:"status"`
	Data   map[string][]Metric `json:"data"`
}

// Metric represents a single metric metadata entry.
type Metric struct {
	Type string `json:"type"`
	Help string `json:"help"`
	Unit string `json:"unit"`
}

// HealthResult represents VictoriaMetrics health status.
type HealthResult struct {
	Status string `json:"status"`
}

// ImportRequest represents a bulk import request.
type ImportRequest struct {
	Metric     map[string]string `json:"metric"`
	Values     []float64         `json:"values,omitempty"`
	Timestamps []int64           `json:"timestamps,omitempty"`
}

// NewClient creates a new VictoriaMetrics client with the specified options.
func NewClient(opts *ClientOptions) (*Client, error) {
	if opts.Address == "" {
		return nil, fmt.Errorf("victoriametrics address is required")
	}

	// Parse and validate URL
	baseURL, err := url.Parse(opts.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid victoriametrics address: %w", err)
	}

	// Store root URL for VM-specific endpoints (e.g., /health, /api/v1/admin/...)
	rootURL := opts.Address
	if !strings.HasSuffix(rootURL, "/") {
		rootURL += "/"
	}

	// Ensure URL has proper path for prometheus-compatible API
	if !strings.HasSuffix(baseURL.Path, "/") {
		baseURL.Path += "/"
	}
	baseURL.Path += "api/v1/"

	// Create HTTP client with timeout and TLS configuration
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: opts.TLSSkipVerify,
		},
	}

	// Configure TLS if certificates are provided
	if opts.TLSCertFile != "" && opts.TLSKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(opts.TLSCertFile, opts.TLSKeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS certificate: %w", err)
		}
		transport.TLSClientConfig.Certificates = []tls.Certificate{cert}
	}

	// Create optimized HTTP client with TLS configuration
	httpClient := optimize.NewOptimizedHTTPClientWithTimeout(timeout)
	if clientTransport, ok := httpClient.Transport.(*http.Transport); ok {
		clientTransport.TLSClientConfig = transport.TLSClientConfig
	}

	// Prepare headers
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["Accept"] = "application/json"

	client := &Client{
		baseURL:        baseURL.String(),
		vmURL:          rootURL,
		httpClient:     httpClient,
		username:       opts.Username,
		password:       opts.Password,
		bearerToken:    opts.BearerToken,
		headers:        headers,
		maxRetries:     opts.MaxRetries,
		retryBaseDelay: opts.RetryBaseDelay,
		retryMaxDelay:  opts.RetryMaxDelay,
	}
	client.maxRetries, client.retryBaseDelay, client.retryMaxDelay = optimize.NormalizeRetryConfig(client.maxRetries, client.retryBaseDelay, client.retryMaxDelay)

	return client, nil
}

// makeRequest performs an HTTP request to the VictoriaMetrics API.
func (c *Client) makeRequest(ctx context.Context, method, endpoint string, params url.Values) (*http.Response, error) {
	reqURL := c.baseURL + endpoint

	var encodedParams string
	if method == "GET" && params != nil {
		reqURL += "?" + params.Encode()
	}
	if method == "POST" && params != nil {
		encodedParams = params.Encode()
	}

	return optimize.DoWithHTTPRetry(
		ctx,
		method,
		c.maxRetries,
		c.retryBaseDelay,
		c.retryMaxDelay,
		func(attempt int) (*http.Response, error) {
			var reqBody io.Reader
			if method == "POST" && encodedParams != "" {
				reqBody = strings.NewReader(encodedParams)
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
			if c.bearerToken != "" {
				req.Header.Set("Authorization", "Bearer "+c.bearerToken)
			} else if c.username != "" && c.password != "" {
				req.SetBasicAuth(c.username, c.password)
			}

			logrus.WithFields(logrus.Fields{
				"method":   method,
				"url":      reqURL,
				"has_auth": c.bearerToken != "" || (c.username != "" && c.password != ""),
				"attempt":  attempt,
			}).Debug("Making VictoriaMetrics API request")

			return c.httpClient.Do(req)
		},
		func(event optimize.HTTPRetryEvent) {
			fields := logrus.Fields{
				"method":   method,
				"url":      reqURL,
				"attempt":  event.Attempt,
				"retry_in": event.Delay,
			}
			if event.Err != nil {
				logrus.WithFields(fields).WithError(event.Err).Warn("Retrying VictoriaMetrics API request after transient transport error")
				return
			}
			fields["status_code"] = event.StatusCode
			logrus.WithFields(fields).Warn("Retrying VictoriaMetrics API request after retryable status")
		},
	)
}

// makeRawRequest performs an HTTP request to the VictoriaMetrics root URL (not /api/v1/).
// Used for VM-specific endpoints like /health.
func (c *Client) makeRawRequest(ctx context.Context, method, path string, body io.Reader, contentType string) (*http.Response, error) {
	reqURL := c.vmURL + strings.TrimPrefix(path, "/")

	req, err := http.NewRequestWithContext(ctx, method, reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// Set authentication
	if c.bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	} else if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

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
		return nil, fmt.Errorf("victoriametrics API error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Query executes a VictoriaMetrics instant query using MetricsQL at a single point in time.
func (c *Client) Query(ctx context.Context, query string, timestamp *time.Time) (*QueryResult, error) {
	logrus.WithField("query", query).Debug("Executing VictoriaMetrics query")

	params := url.Values{}
	params.Set("query", query)
	if timestamp != nil {
		params.Set("time", strconv.FormatInt(timestamp.Unix(), 10))
	}

	resp, err := c.makeRequest(ctx, "GET", "query", params)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result QueryResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse query result: %w", err)
	}

	if result.Status == "error" {
		return nil, fmt.Errorf("victoriametrics query error: %s (type: %s)", result.Error, result.ErrorType)
	}

	return &result, nil
}

// QueryRange executes a VictoriaMetrics range query using MetricsQL over a time range.
func (c *Client) QueryRange(ctx context.Context, query, start, end, step string) (*QueryResult, error) {
	logrus.WithFields(logrus.Fields{
		"query": query,
		"start": start,
		"end":   end,
		"step":  step,
	}).Debug("Executing VictoriaMetrics range query")

	params := url.Values{}
	params.Set("query", query)
	params.Set("start", start)
	params.Set("end", end)
	params.Set("step", step)

	resp, err := c.makeRequest(ctx, "GET", "query_range", params)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result QueryResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse query result: %w", err)
	}

	if result.Status == "error" {
		return nil, fmt.Errorf("victoriametrics query error: %s", result.Error)
	}

	return &result, nil
}

// Export exports raw metrics data in bulk from VictoriaMetrics using the native export API.
func (c *Client) Export(ctx context.Context, match []string, start, end time.Time) ([]ExportResult, error) {
	logrus.WithFields(logrus.Fields{
		"match_count": len(match),
		"start":       start,
		"end":         end,
	}).Debug("Exporting metrics from VictoriaMetrics")

	params := url.Values{}
	for _, m := range match {
		params.Add("match[]", m)
	}
	params.Set("start", strconv.FormatInt(start.Unix(), 10))
	params.Set("end", strconv.FormatInt(end.Unix(), 10))

	resp, err := c.makeRequest(ctx, "GET", "export", params)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	// VictoriaMetrics export returns JSON lines format
	var results []ExportResult
	lines := strings.Split(strings.TrimSpace(string(body)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		var entry ExportResult
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			logrus.WithError(err).WithField("line", line).Warn("Failed to parse export line, skipping")
			continue
		}
		results = append(results, entry)
	}

	return results, nil
}

// Import imports metrics data in bulk into VictoriaMetrics using the native import API.
func (c *Client) Import(ctx context.Context, data []ImportRequest) error {
	logrus.WithField("count", len(data)).Debug("Importing metrics into VictoriaMetrics")

	// Build JSON lines payload
	var lines []string
	for _, entry := range data {
		jsonBytes, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("failed to marshal import entry: %w", err)
		}
		lines = append(lines, string(jsonBytes))
	}
	payload := strings.Join(lines, "\n")

	resp, err := c.makeRawRequest(ctx, "POST", "api/v1/import", strings.NewReader(payload), "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}

	_, err = c.handleResponse(resp)
	return err
}

// GetSeriesCount returns the total number of time series stored in VictoriaMetrics.
func (c *Client) GetSeriesCount(ctx context.Context) (uint64, error) {
	logrus.Debug("Getting VictoriaMetrics series count")

	resp, err := c.makeRequest(ctx, "GET", "series/count", nil)
	if err != nil {
		return 0, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return 0, err
	}

	var result struct {
		Status string `json:"status"`
		Data   []uint64 `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("failed to parse series count: %w", err)
	}

	if len(result.Data) == 0 {
		return 0, nil
	}
	return result.Data[0], nil
}

// GetLabelNames returns all label names from VictoriaMetrics.
func (c *Client) GetLabelNames(ctx context.Context, start, end *time.Time) ([]string, error) {
	logrus.Debug("Getting VictoriaMetrics label names")

	params := url.Values{}
	if start != nil {
		params.Set("start", strconv.FormatInt(start.Unix(), 10))
	}
	if end != nil {
		params.Set("end", strconv.FormatInt(end.Unix(), 10))
	}

	resp, err := c.makeRequest(ctx, "GET", "labels", params)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result LabelValuesResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse label names: %w", err)
	}

	return result.Data, nil
}

// GetLabelValues returns all values for a specific label name in VictoriaMetrics.
func (c *Client) GetLabelValues(ctx context.Context, labelName string, start, end *time.Time) ([]string, error) {
	logrus.WithField("label", labelName).Debug("Getting VictoriaMetrics label values")

	params := url.Values{}
	if start != nil {
		params.Set("start", strconv.FormatInt(start.Unix(), 10))
	}
	if end != nil {
		params.Set("end", strconv.FormatInt(end.Unix(), 10))
	}

	resp, err := c.makeRequest(ctx, "GET", "label/"+labelName+"/values", params)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result LabelValuesResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse label values: %w", err)
	}

	return result.Data, nil
}

// GetTSDBStats returns TSDB storage statistics.
func (c *Client) GetTSDBStats(ctx context.Context) (*TSDBStats, error) {
	logrus.Debug("Getting VictoriaMetrics TSDB stats")

	resp, err := c.makeRequest(ctx, "GET", "status/tsdb", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result TSDBStats
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse TSDB stats: %w", err)
	}

	return &result, nil
}

// GetActiveQueries returns currently running queries.
func (c *Client) GetActiveQueries(ctx context.Context) (*ActiveQuery, error) {
	logrus.Debug("Getting VictoriaMetrics active queries")

	resp, err := c.makeRequest(ctx, "GET", "status/active_queries", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result ActiveQuery
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse active queries: %w", err)
	}

	return &result, nil
}

// DeleteSeries deletes time series matching the given label matchers.
func (c *Client) DeleteSeries(ctx context.Context, match []string) (int, error) {
	logrus.WithField("match", match).Debug("Deleting VictoriaMetrics series")

	params := url.Values{}
	for _, m := range match {
		params.Add("match[]", m)
	}

	resp, err := c.makeRequest(ctx, "POST", "admin/tsdb/delete_series", params)
	if err != nil {
		return 0, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return 0, err
	}

	var result struct {
		Status   string `json:"status"`
		Data     struct {
			DeletedSeries int `json:"deletedSeries"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("failed to parse delete series result: %w", err)
	}

	return result.Data.DeletedSeries, nil
}

// HealthCheck checks VictoriaMetrics health status.
func (c *Client) HealthCheck(ctx context.Context) (*HealthResult, error) {
	logrus.Debug("Checking VictoriaMetrics health")

	resp, err := c.makeRawRequest(ctx, "GET", "health", nil, "")
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result HealthResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse health check: %w", err)
	}

	return &result, nil
}

// GetMetadata returns metadata about metrics in VictoriaMetrics.
func (c *Client) GetMetadata(ctx context.Context, metric string, limit int) (*MetadataResult, error) {
	logrus.WithField("metric", metric).Debug("Getting VictoriaMetrics metric metadata")

	params := url.Values{}
	if metric != "" {
		params.Set("metric", metric)
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	resp, err := c.makeRequest(ctx, "GET", "metadata", params)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result MetadataResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return &result, nil
}

// Package client provides Prometheus HTTP API client functionality.
// It offers operations for querying Prometheus metrics, targets, and alerts
// through the Prometheus HTTP API.
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

	optimize "github.com/mahmut-Abi/k8s-mcp-server/internal/util/performance"
	"github.com/sirupsen/logrus"
)

// ClientOptions holds configuration parameters for creating a Prometheus client.
type ClientOptions struct {
	Address       string        // Prometheus server address
	Username      string        // Username for basic authentication
	Password      string        // Password for basic authentication
	BearerToken   string        // Bearer token for authentication
	Timeout       time.Duration // HTTP request timeout
	TLSSkipVerify bool          // Skip TLS certificate verification
	TLSCertFile   string        // TLS certificate file
	TLSKeyFile    string        // TLS key file
	TLSCAFile     string        // TLS CA file
}

// Client provides operations for interacting with Prometheus API.
type Client struct {
	baseURL     string            // Base URL for Prometheus API
	httpClient  *http.Client      // HTTP client for API requests
	username    string            // Username for basic auth
	password    string            // Password for basic auth
	bearerToken string            // Bearer token for auth
	headers     map[string]string // Additional headers
}

// QueryResult represents a Prometheus query result.
type QueryResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string      `json:"resultType"`
		Result     interface{} `json:"result"`
	} `json:"data"`
	ErrorType string `json:"errorType,omitempty"`
	Error     string `json:"error,omitempty"`
}

// MetricSample represents a single metric sample.
type MetricSample struct {
	Metric    map[string]string `json:"metric"`
	Value     []interface{}     `json:"value,omitempty"`
	Values    [][]interface{}   `json:"values,omitempty"`
	Timestamp float64           `json:"timestamp,omitempty"`
}

// Target represents a Prometheus target.
type Target struct {
	DiscoveredLabels   map[string]string `json:"discoveredLabels"`
	Labels             map[string]string `json:"labels"`
	ScrapePool         string            `json:"scrapePool"`
	ScrapeURL          string            `json:"scrapeUrl"`
	GlobalURL          string            `json:"globalUrl"`
	LastError          string            `json:"lastError"`
	LastScrape         time.Time         `json:"lastScrape"`
	LastScrapeDuration float64           `json:"lastScrapeDuration"`
	Health             string            `json:"health"`
	ScrapeInterval     string            `json:"scrapeInterval"`
	ScrapeTimeout      string            `json:"scrapeTimeout"`
}

// Alert represents a Prometheus alert.
type Alert struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	State       string            `json:"state"`
	ActiveAt    *time.Time        `json:"activeAt,omitempty"`
	Value       string            `json:"value"`
}

// Rule represents a Prometheus recording or alerting rule.
type Rule struct {
	Name           string            `json:"name"`
	Query          string            `json:"query"`
	Duration       float64           `json:"duration,omitempty"`
	Labels         map[string]string `json:"labels,omitempty"`
	Annotations    map[string]string `json:"annotations,omitempty"`
	Alerts         []Alert           `json:"alerts,omitempty"`
	Health         string            `json:"health"`
	LastError      string            `json:"lastError,omitempty"`
	Type           string            `json:"type"`
	EvaluationTime float64           `json:"evaluationTime"`
	LastEvaluation time.Time         `json:"lastEvaluation"`
}

// RuleGroup represents a group of Prometheus rules.
type RuleGroup struct {
	Name                    string    `json:"name"`
	File                    string    `json:"file"`
	Rules                   []Rule    `json:"rules"`
	Interval                float64   `json:"interval"`
	Limit                   int       `json:"limit"`
	EvaluationTime          float64   `json:"evaluationTime"`
	LastEvaluation          time.Time `json:"lastEvaluation"`
	PartialResponseStrategy string    `json:"partialResponseStrategy,omitempty"`
}

// NewClient creates a new Prometheus client with the specified options.
func NewClient(opts *ClientOptions) (*Client, error) {
	if opts.Address == "" {
		return nil, fmt.Errorf("prometheus address is required")
	}

	// Parse and validate URL
	baseURL, err := url.Parse(opts.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid prometheus address: %w", err)
	}

	// Ensure URL has proper path
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
		baseURL:     baseURL.String(),
		httpClient:  httpClient,
		username:    opts.Username,
		password:    opts.Password,
		bearerToken: opts.BearerToken,
		headers:     headers,
	}

	return client, nil
}

// makeRequest performs an HTTP request to the Prometheus API.
func (c *Client) makeRequest(ctx context.Context, method, endpoint string, params url.Values) (*http.Response, error) {
	var reqURL string
	if method == "GET" && params != nil {
		reqURL = c.baseURL + endpoint + "?" + params.Encode()
	} else {
		reqURL = c.baseURL + endpoint
	}

	var reqBody io.Reader
	if method == "POST" && params != nil {
		reqBody = strings.NewReader(params.Encode())
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
	}).Debug("Making Prometheus API request")

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
		return nil, fmt.Errorf("prometheus API error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Query executes a Prometheus query at a single point in time.
func (c *Client) Query(ctx context.Context, query string, timestamp *time.Time) (*QueryResult, error) {
	logrus.WithField("query", query).Debug("Executing Prometheus query")

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
		return nil, fmt.Errorf("failed to unmarshal query result: %w", err)
	}

	logrus.Debug("Prometheus query executed successfully")
	return &result, nil
}

// QueryRange executes a Prometheus query over a range of time.
func (c *Client) QueryRange(ctx context.Context, query string, start, end time.Time, step string) (*QueryResult, error) {
	logrus.WithFields(logrus.Fields{
		"query": query,
		"start": start,
		"end":   end,
		"step":  step,
	}).Debug("Executing Prometheus range query")

	params := url.Values{}
	params.Set("query", query)
	params.Set("start", strconv.FormatInt(start.Unix(), 10))
	params.Set("end", strconv.FormatInt(end.Unix(), 10))
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
		return nil, fmt.Errorf("failed to unmarshal range query result: %w", err)
	}

	logrus.Debug("Prometheus range query executed successfully")
	return &result, nil
}

// GetTargets retrieves the current targets from Prometheus.
func (c *Client) GetTargets(ctx context.Context, state string) ([]Target, error) {
	logrus.WithField("state", state).Debug("Getting Prometheus targets")

	params := url.Values{}
	if state != "" {
		params.Set("state", state)
	}

	resp, err := c.makeRequest(ctx, "GET", "targets", params)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status string `json:"status"`
		Data   struct {
			ActiveTargets  []Target `json:"activeTargets"`
			DroppedTargets []Target `json:"droppedTargets"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal targets: %w", err)
	}

	// Combine active and dropped targets
	var targets []Target
	targets = append(targets, result.Data.ActiveTargets...)
	targets = append(targets, result.Data.DroppedTargets...)

	logrus.WithField("count", len(targets)).Debug("Retrieved Prometheus targets")
	return targets, nil
}

// GetRules retrieves the current recording and alerting rules from Prometheus.
func (c *Client) GetRules(ctx context.Context, ruleType string) ([]RuleGroup, error) {
	logrus.WithField("type", ruleType).Debug("Getting Prometheus rules")

	params := url.Values{}
	if ruleType != "" {
		params.Set("type", ruleType)
	}

	resp, err := c.makeRequest(ctx, "GET", "rules", params)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status string `json:"status"`
		Data   struct {
			Groups []RuleGroup `json:"groups"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rules: %w", err)
	}

	logrus.WithField("count", len(result.Data.Groups)).Debug("Retrieved Prometheus rules")
	return result.Data.Groups, nil
}

// GetAlerts retrieves the current alerts from Prometheus.
func (c *Client) GetAlerts(ctx context.Context) ([]Alert, error) {
	logrus.Debug("Getting Prometheus alerts")

	resp, err := c.makeRequest(ctx, "GET", "alerts", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status string `json:"status"`
		Data   struct {
			Alerts []Alert `json:"alerts"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal alerts: %w", err)
	}

	logrus.WithField("count", len(result.Data.Alerts)).Debug("Retrieved Prometheus alerts")
	return result.Data.Alerts, nil
}

// GetLabelNames retrieves the list of label names.
func (c *Client) GetLabelNames(ctx context.Context, start, end *time.Time) ([]string, error) {
	logrus.Debug("Getting Prometheus label names")

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

	var result struct {
		Status string   `json:"status"`
		Data   []string `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal label names: %w", err)
	}

	logrus.WithField("count", len(result.Data)).Debug("Retrieved Prometheus label names")
	return result.Data, nil
}

// GetLabelValues retrieves the list of label values for a given label name.
func (c *Client) GetLabelValues(ctx context.Context, labelName string, start, end *time.Time) ([]string, error) {
	logrus.WithField("label", labelName).Debug("Getting Prometheus label values")

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

	var result struct {
		Status string   `json:"status"`
		Data   []string `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal label values: %w", err)
	}

	logrus.WithField("count", len(result.Data)).Debug("Retrieved Prometheus label values")
	return result.Data, nil
}

// GetSeries retrieves series matching label selectors.
func (c *Client) GetSeries(ctx context.Context, matches []string, start, end *time.Time) ([]map[string]string, error) {
	logrus.WithField("matches", matches).Debug("Getting Prometheus series")

	params := url.Values{}
	for _, match := range matches {
		params.Add("match[]", match)
	}
	if start != nil {
		params.Set("start", strconv.FormatInt(start.Unix(), 10))
	}
	if end != nil {
		params.Set("end", strconv.FormatInt(end.Unix(), 10))
	}

	resp, err := c.makeRequest(ctx, "GET", "series", params)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status string              `json:"status"`
		Data   []map[string]string `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal series: %w", err)
	}

	logrus.WithField("count", len(result.Data)).Debug("Retrieved Prometheus series")
	return result.Data, nil
}

// TestConnection tests the connection to Prometheus API.
func (c *Client) TestConnection(ctx context.Context) error {
	logrus.Debug("Testing Prometheus connection")

	resp, err := c.makeRequest(ctx, "GET", "query", url.Values{"query": {"up"}})
	if err != nil {
		return fmt.Errorf("failed to connect to prometheus: %w", err)
	}

	_, err = c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("prometheus health check failed: %w", err)
	}

	logrus.Debug("Prometheus connection test successful")
	return nil
}

// ServerInfo represents Prometheus server information.
type ServerInfo struct {
	Status    string                 `json:"status"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Version   string                 `json:"version,omitempty"`
	BuildInfo map[string]interface{} `json:"buildInfo,omitempty"`
}

// MetricsMetadata represents metadata for a metric.
type MetricsMetadata struct {
	Type string `json:"type"`
	Help string `json:"help"`
	Unit string `json:"unit,omitempty"`
}

// GetServerInfo retrieves Prometheus server information.
func (c *Client) GetServerInfo(ctx context.Context) (*ServerInfo, error) {
	logrus.Debug("Getting Prometheus server info")

	resp, err := c.makeRequest(ctx, "GET", "status/config", url.Values{})
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status string            `json:"status"`
		Data   map[string]string `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal server info: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err == nil {
		if dataMap, ok := data["data"].(map[string]interface{}); ok {
			data = dataMap
		}
	}

	serverInfo := &ServerInfo{
		Status: result.Status,
		Data:   data,
	}

	logrus.Debug("Retrieved Prometheus server info")
	return serverInfo, nil
}

// GetMetricsMetadata retrieves metadata for all metrics or a specific metric.
func (c *Client) GetMetricsMetadata(ctx context.Context, metric string) (map[string][]MetricsMetadata, error) {
	logrus.WithField("metric", metric).Debug("Getting Prometheus metrics metadata")

	params := url.Values{}
	if metric != "" {
		params.Set("metric", metric)
	}

	resp, err := c.makeRequest(ctx, "GET", "metadata", params)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status string                       `json:"status"`
		Data   map[string][]MetricsMetadata `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metrics metadata: %w", err)
	}

	logrus.WithField("count", len(result.Data)).Debug("Retrieved Prometheus metrics metadata")
	return result.Data, nil
}

// GetTargetMetadata retrieves metadata for targets.
func (c *Client) GetTargetMetadata(ctx context.Context, metric string) ([]map[string]interface{}, error) {
	logrus.WithField("metric", metric).Debug("Getting Prometheus target metadata")

	params := url.Values{}
	if metric != "" {
		params.Set("metric", metric)
	}

	resp, err := c.makeRequest(ctx, "GET", "targets/metadata", params)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status string                   `json:"status"`
		Data   []map[string]interface{} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal target metadata: %w", err)
	}

	logrus.WithField("count", len(result.Data)).Debug("Retrieved Prometheus target metadata")
	return result.Data, nil
}

// ============ TSDB Operations ============

// TSDBStats represents Prometheus TSDB statistics.
type TSDBStats struct {
	HeadStats struct {
		NumSeries            int64 `json:"numSeries"`
		NumLabelValuePairs   int64 `json:"numLabelValuePairs"`
		SeriesCreatedSeries  int64 `json:"seriesCreatedSeries"`
		SeriesRemovedSeries  int64 `json:"seriesRemovedSeries"`
		SeriesNotFoundSeries int64 `json:"seriesNotFoundSeries"`
		SampleIngestedSeries int64 `json:"sampleIngestedSeries"`
		OutOfOrderSamples    int64 `json:"outOfOrderSamples"`
		TooOldSamples        int64 `json:"tooOldSamples"`
	} `json:"headStats"`
	ChunkCount       int64 `json:"chunkCount"`
	ChunksSize       int64 `json:"chunksSize"`
	MinTime          int64 `json:"minTime"`
	MaxTime          int64 `json:"maxTime"`
	NumSeriesCreated int64 `json:"numSeriesCreated"`
	NumSeriesRemoved int64 `json:"numSeriesRemoved"`
	NumSeriesActive  int64 `json:"numSeriesActive"`
}

// GetTSDBStats retrieves TSDB statistics.
func (c *Client) GetTSDBStats(ctx context.Context) (*TSDBStats, error) {
	logrus.Debug("Getting Prometheus TSDB stats")

	resp, err := c.makeRequest(ctx, "GET", "tsdb/stat", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var stats TSDBStats
	if err := json.Unmarshal(body, &stats); err != nil {
		return nil, fmt.Errorf("failed to unmarshal TSDB stats: %w", err)
	}

	logrus.Debug("Retrieved Prometheus TSDB stats")
	return &stats, nil
}

// TSDBStatus represents TSDB status information.
type TSDBStatus struct {
	Head struct {
		NumSeries          int64 `json:"numSeries"`
		MaxTime            int64 `json:"maxTime"`
		MinTime            int64 `json:"minTime"`
		NumLabelValuePairs int64 `json:"numLabelValuePairs"`
	} `json:"head"`
	Blocks []TSDBBlock `json:"blocks,omitempty"`
}

// TSDBBlock represents a TSDB block.
type TSDBBlock struct {
	ULID       string `json:"ulid"`
	MinTime    int64  `json:"minTime"`
	MaxTime    int64  `json:"maxTime"`
	Duration   int64  `json:"duration"`
	NumSamples int64  `json:"numSamples"`
	NumSeries  int64  `json:"numSeries"`
	NumChunks  int64  `json:"numChunks"`
	Size       int64  `json:"size"`
}

// GetTSDBStatus retrieves TSDB status information.
func (c *Client) GetTSDBStatus(ctx context.Context) (*TSDBStatus, error) {
	logrus.Debug("Getting Prometheus TSDB status")

	resp, err := c.makeRequest(ctx, "GET", "status/tsdb", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var status TSDBStatus
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal TSDB status: %w", err)
	}

	logrus.Debug("Retrieved Prometheus TSDB status")
	return &status, nil
}

// RuntimeInfo represents Prometheus runtime/build information.
type RuntimeInfo struct {
	Version   string `json:"version"`
	Revision  string `json:"revision"`
	Branch    string `json:"branch"`
	BuildUser string `json:"buildUser"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`
}

// GetRuntimeInfo retrieves Prometheus runtime and build information.
func (c *Client) GetRuntimeInfo(ctx context.Context) (*RuntimeInfo, error) {
	logrus.Debug("Getting Prometheus runtime info")

	resp, err := c.makeRequest(ctx, "GET", "status/buildinfo", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status string `json:"status"`
		Data   struct {
			VersionInfo RuntimeInfo `json:"versionInfo"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal runtime info: %w", err)
	}

	logrus.Debug("Retrieved Prometheus runtime info")
	return &result.Data.VersionInfo, nil
}

// SnapshotResult represents the result of creating a TSDB snapshot.
type SnapshotResult struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// CreateSnapshot creates a TSDB snapshot.
func (c *Client) CreateSnapshot(ctx context.Context, skipHead bool) (*SnapshotResult, error) {
	logrus.Debug("Creating Prometheus TSDB snapshot")

	params := url.Values{}
	if skipHead {
		params.Set("skipHead", "true")
	}

	resp, err := c.makeRequest(ctx, "POST", "tsdb/snapshot", params)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result SnapshotResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal snapshot result: %w", err)
	}

	logrus.WithField("name", result.Name).Debug("Created Prometheus TSDB snapshot")
	return &result, nil
}

// WALReplayStatus represents the WAL replay status.
type WALReplayStatus struct {
	MinSegment int   `json:"minSegment"`
	MaxSegment int   `json:"maxSegment"`
	Duration   int64 `json:"duration"`
}

// GetWALReplayStatus retrieves WAL replay status.
func (c *Client) GetWALReplayStatus(ctx context.Context) (*WALReplayStatus, error) {
	logrus.Debug("Getting Prometheus WAL replay status")

	resp, err := c.makeRequest(ctx, "GET", "tsdb/wal/replay", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var status WALReplayStatus
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal WAL replay status: %w", err)
	}

	logrus.Debug("Retrieved Prometheus WAL replay status")
	return &status, nil
}

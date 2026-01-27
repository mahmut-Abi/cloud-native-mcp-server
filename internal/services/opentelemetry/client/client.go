// Package client provides HTTP client for OpenTelemetry Collector API operations.
package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/errors"
	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
)

// ClientOptions defines configuration options for the OpenTelemetry HTTP client.
type ClientOptions struct {
	Address       string        // OpenTelemetry Collector address (e.g., http://localhost:4318)
	Username      string        // Basic auth username
	Password      string        // Basic auth password
	BearerToken   string        // Bearer token for authentication
	Timeout       time.Duration // Request timeout
	TLSSkipVerify bool          // Skip TLS certificate verification
	TLSCertFile   string        // Path to TLS certificate file
	TLSKeyFile    string        // Path to TLS key file
	TLSCAFile     string        // Path to TLS CA file
}

// Client represents an HTTP client for OpenTelemetry Collector API operations.
type Client struct {
	httpClient  *http.Client
	baseURL     string
	username    string
	password    string
	bearerToken string
}

// NewClient creates a new OpenTelemetry HTTP client with the provided options.
func NewClient(opts *ClientOptions) (*Client, error) {
	if opts == nil {
		return nil, errors.InvalidParamError("options", "cannot be nil")
	}

	if opts.Address == "" {
		return nil, errors.InvalidParamError("address", "is required")
	}

	// Create HTTP client with optimized configuration
	httpClient := optimize.NewOptimizedHTTPClientWithTimeout(opts.Timeout)

	// Configure TLS if needed
	if opts.TLSSkipVerify || opts.TLSCertFile != "" || opts.TLSCAFile != "" {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: opts.TLSSkipVerify,
		}

		// Load CA certificate if provided
		if opts.TLSCAFile != "" {
			caCert, err := os.ReadFile(opts.TLSCAFile)
			if err != nil {
				return nil, errors.InternalError(fmt.Errorf("failed to read CA certificate: %w", err))
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			tlsConfig.RootCAs = caCertPool
		}

		// Load client certificate if provided
		if opts.TLSCertFile != "" && opts.TLSKeyFile != "" {
			cert, err := tls.LoadX509KeyPair(opts.TLSCertFile, opts.TLSKeyFile)
			if err != nil {
				return nil, errors.InternalError(fmt.Errorf("failed to load client certificate: %w", err))
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}

		httpClient.Transport.(*http.Transport).TLSClientConfig = tlsConfig
	}

	return &Client{
		httpClient:  httpClient,
		baseURL:     opts.Address,
		username:    opts.Username,
		password:    opts.Password,
		bearerToken: opts.BearerToken,
	}, nil
}

// makeRequest performs an HTTP request to the OpenTelemetry Collector API.
func (c *Client) makeRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authentication headers
	if c.bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	} else if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	return c.handleResponse(resp)
}

// handleResponse processes the HTTP response from OpenTelemetry Collector API.
func (c *Client) handleResponse(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("OpenTelemetry API error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// GetHealth retrieves the health status of the OpenTelemetry Collector.
func (c *Client) GetHealth(ctx context.Context) (map[string]interface{}, error) {
	body, err := c.makeRequest(ctx, http.MethodGet, "/healthz", nil)
	if err != nil {
		return nil, err
	}

	var health map[string]interface{}
	if err := json.Unmarshal(body, &health); err != nil {
		return nil, fmt.Errorf("failed to unmarshal health response: %w", err)
	}

	return health, nil
}

// GetStatus retrieves the status information of the OpenTelemetry Collector.
func (c *Client) GetStatus(ctx context.Context) (map[string]interface{}, error) {
	body, err := c.makeRequest(ctx, http.MethodGet, "/status", nil)
	if err != nil {
		return nil, err
	}

	var status map[string]interface{}
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal status response: %w", err)
	}

	return status, nil
}

// GetConfig retrieves the configuration of the OpenTelemetry Collector.
func (c *Client) GetConfig(ctx context.Context) (map[string]interface{}, error) {
	body, err := c.makeRequest(ctx, http.MethodGet, "/config", nil)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	if err := json.Unmarshal(body, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config response: %w", err)
	}

	return config, nil
}

// GetMetrics retrieves metrics from the OpenTelemetry Collector.
func (c *Client) GetMetrics(ctx context.Context, metricName *string, startTime *time.Time, endTime *time.Time) (map[string]interface{}, error) {
	path := "/metrics"

	queryParams := make(map[string]string)
	if metricName != nil {
		queryParams["name"] = *metricName
	}
	if startTime != nil {
		queryParams["start_time"] = startTime.Format(time.RFC3339)
	}
	if endTime != nil {
		queryParams["end_time"] = endTime.Format(time.RFC3339)
	}

	// Build query string
	query := ""
	if len(queryParams) > 0 {
		query = "?"
		first := true
		for k, v := range queryParams {
			if !first {
				query += "&"
			}
			query += fmt.Sprintf("%s=%s", k, v)
			first = false
		}
	}

	body, err := c.makeRequest(ctx, http.MethodGet, path+query, nil)
	if err != nil {
		return nil, err
	}

	var metrics map[string]interface{}
	if err := json.Unmarshal(body, &metrics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metrics response: %w", err)
	}

	return metrics, nil
}

// QueryMetrics executes a PromQL-style query against the OpenTelemetry Collector.
func (c *Client) QueryMetrics(ctx context.Context, query string, queryTime *time.Time) (map[string]interface{}, error) {
	if query == "" {
		return nil, errors.InvalidParamError("query", "cannot be empty")
	}

	path := "/api/v1/query"

	requestBody := map[string]string{
		"query": query,
	}
	if queryTime != nil {
		requestBody["time"] = queryTime.Format(time.RFC3339)
	}

	body, err := c.makeRequest(ctx, http.MethodPost, path, requestBody)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal query response: %w", err)
	}

	return result, nil
}

// GetTraces retrieves traces from the OpenTelemetry Collector.
func (c *Client) GetTraces(ctx context.Context, traceID *string, service *string, startTime *time.Time, endTime *time.Time, limit *int) (map[string]interface{}, error) {
	path := "/traces"

	queryParams := make(map[string]string)
	if traceID != nil {
		queryParams["trace_id"] = *traceID
	}
	if service != nil {
		queryParams["service"] = *service
	}
	if startTime != nil {
		queryParams["start_time"] = startTime.Format(time.RFC3339)
	}
	if endTime != nil {
		queryParams["end_time"] = endTime.Format(time.RFC3339)
	}
	if limit != nil {
		queryParams["limit"] = fmt.Sprintf("%d", *limit)
	}

	// Build query string
	query := ""
	if len(queryParams) > 0 {
		query = "?"
		first := true
		for k, v := range queryParams {
			if !first {
				query += "&"
			}
			query += fmt.Sprintf("%s=%s", k, v)
			first = false
		}
	}

	body, err := c.makeRequest(ctx, http.MethodGet, path+query, nil)
	if err != nil {
		return nil, err
	}

	var traces map[string]interface{}
	if err := json.Unmarshal(body, &traces); err != nil {
		return nil, fmt.Errorf("failed to unmarshal traces response: %w", err)
	}

	return traces, nil
}

// QueryTraces searches for traces matching criteria in the OpenTelemetry Collector.
func (c *Client) QueryTraces(ctx context.Context, query string, service *string, startTime *time.Time, endTime *time.Time, limit *int) (map[string]interface{}, error) {
	if query == "" {
		return nil, errors.InvalidParamError("query", "cannot be empty")
	}

	path := "/api/v1/traces"

	requestBody := map[string]interface{}{
		"query": query,
	}
	if service != nil {
		requestBody["service"] = *service
	}
	if startTime != nil {
		requestBody["start_time"] = startTime.Format(time.RFC3339)
	}
	if endTime != nil {
		requestBody["end_time"] = endTime.Format(time.RFC3339)
	}
	if limit != nil {
		requestBody["limit"] = *limit
	}

	body, err := c.makeRequest(ctx, http.MethodPost, path, requestBody)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal traces query response: %w", err)
	}

	return result, nil
}

// GetLogs retrieves logs from the OpenTelemetry Collector.
func (c *Client) GetLogs(ctx context.Context, service *string, level *string, startTime *time.Time, endTime *time.Time, limit *int) (map[string]interface{}, error) {
	path := "/logs"

	queryParams := make(map[string]string)
	if service != nil {
		queryParams["service"] = *service
	}
	if level != nil {
		queryParams["level"] = *level
	}
	if startTime != nil {
		queryParams["start_time"] = startTime.Format(time.RFC3339)
	}
	if endTime != nil {
		queryParams["end_time"] = endTime.Format(time.RFC3339)
	}
	if limit != nil {
		queryParams["limit"] = fmt.Sprintf("%d", *limit)
	}

	// Build query string
	query := ""
	if len(queryParams) > 0 {
		query = "?"
		first := true
		for k, v := range queryParams {
			if !first {
				query += "&"
			}
			query += fmt.Sprintf("%s=%s", k, v)
			first = false
		}
	}

	body, err := c.makeRequest(ctx, http.MethodGet, path+query, nil)
	if err != nil {
		return nil, err
	}

	var logs map[string]interface{}
	if err := json.Unmarshal(body, &logs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal logs response: %w", err)
	}

	return logs, nil
}

// QueryLogs searches for logs matching criteria in the OpenTelemetry Collector.
func (c *Client) QueryLogs(ctx context.Context, query string, service *string, level *string, startTime *time.Time, endTime *time.Time, limit *int) (map[string]interface{}, error) {
	if query == "" {
		return nil, errors.InvalidParamError("query", "cannot be empty")
	}

	path := "/api/v1/logs"

	requestBody := map[string]interface{}{
		"query": query,
	}
	if service != nil {
		requestBody["service"] = *service
	}
	if level != nil {
		requestBody["level"] = *level
	}
	if startTime != nil {
		requestBody["start_time"] = startTime.Format(time.RFC3339)
	}
	if endTime != nil {
		requestBody["end_time"] = endTime.Format(time.RFC3339)
	}
	if limit != nil {
		requestBody["limit"] = *limit
	}

	body, err := c.makeRequest(ctx, http.MethodPost, path, requestBody)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal logs query response: %w", err)
	}

	return result, nil
}

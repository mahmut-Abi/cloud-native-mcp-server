package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mahmut-Abi/k8s-mcp-server/internal/errors"
	"github.com/sirupsen/logrus"
)

// ClientOptions holds configuration parameters for creating a Jaeger client.
type ClientOptions struct {
	BaseURL string        // Jaeger server base URL
	Timeout time.Duration // HTTP request timeout
}

// Client provides operations for interacting with Jaeger API.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// Trace represents a Jaeger trace.
type Trace struct {
	TraceID   string    `json:"traceID"`
	Spans     []Span    `json:"spans"`
	Processes []Process `json:"processes"`
	Warnings  []string  `json:"warnings"`
}

// Span represents a span in a trace.
type Span struct {
	TraceID       string      `json:"traceID"`
	SpanID        string      `json:"spanID"`
	OperationName string      `json:"operationName"`
	References    []Reference `json:"references"`
	StartTime     int64       `json:"startTime"`
	Duration      int64       `json:"duration"`
	Tags          []KeyValue  `json:"tags"`
	Logs          []Log       `json:"logs"`
	ProcessID     string      `json:"processID"`
	Warnings      []string    `json:"warnings"`
}

// Process represents a process in a trace.
type Process struct {
	ServiceName string     `json:"serviceName"`
	Tags        []KeyValue `json:"tags"`
}

// Reference represents a span reference.
type Reference struct {
	RefType string `json:"refType"`
	TraceID string `json:"traceID"`
	SpanID  string `json:"spanID"`
}

// KeyValue represents a tag or log field.
type KeyValue struct {
	Key   string      `json:"key"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// Log represents a log entry.
type Log struct {
	Timestamp int64      `json:"timestamp"`
	Fields    []KeyValue `json:"fields"`
}

// TraceQueryParameters represents parameters for trace search.
type TraceQueryParameters struct {
	Service     string            `json:"service"`
	Operation   string            `json:"operation"`
	Tags        map[string]string `json:"tags"`
	StartTime   string            `json:"startTime"`
	EndTime     string            `json:"endTime"`
	Limit       int               `json:"limit"`
	MinDuration string            `json:"minDuration"`
	MaxDuration string            `json:"maxDuration"`
}

// Service represents a Jaeger service.
type Service struct {
	Name       string   `json:"name"`
	Operations []string `json:"operations"`
}

// Dependency represents a service dependency.
type Dependency struct {
	Parent    string `json:"parent"`
	Child     string `json:"child"`
	CallCount int64  `json:"callCount"`
}

// NewClient creates a new Jaeger client with the specified options.
func NewClient(opts *ClientOptions) (*Client, error) {
	if opts.BaseURL == "" {
		return nil, fmt.Errorf("jaeger base URL is required")
	}

	// Parse and validate URL
	baseURL, err := url.Parse(opts.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid jaeger base URL: %w", err)
	}

	// Ensure URL has proper path
	if !strings.HasSuffix(baseURL.Path, "/") {
		baseURL.Path += "/"
	}

	// Create HTTP client with timeout
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	httpClient := &http.Client{
		Timeout: timeout,
	}

	client := &Client{
		baseURL:    baseURL.String(),
		httpClient: httpClient,
	}

	return client, nil
}

// makeRequest performs an HTTP request to the Jaeger API.
func (c *Client) makeRequest(ctx context.Context, method, endpoint string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + strings.TrimPrefix(endpoint, "/")

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	logrus.WithFields(logrus.Fields{
		"method": method,
		"url":    url,
	}).Debug("Making Jaeger API request")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "JAEGER_CONNECTION_FAILED", "failed to connect to Jaeger").
			WithHTTPStatus(503)
	}

	return resp, nil
}

// handleResponse processes the HTTP response and returns the body.
func (c *Client) handleResponse(resp *http.Response) ([]byte, error) {
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "JAEGER_INVALID_RESPONSE", "invalid Jaeger API response").
			WithHTTPStatus(502)
	}

	if resp.StatusCode >= 400 {
		switch resp.StatusCode {
		case 401:
			return nil, errors.New("JAEGER_UNAUTHORIZED", "unauthorized access to Jaeger").
				WithHTTPStatus(401)
		case 403:
			return nil, errors.New("JAEGER_FORBIDDEN", "forbidden access to Jaeger resource").
				WithHTTPStatus(403)
		case 404:
			return nil, errors.NotFoundError("resource")
		case 429:
			return nil, errors.New("JAEGER_RATE_LIMITED", "Jaeger API rate limit exceeded").
				WithHTTPStatus(429)
		default:
			return nil, errors.New("JAEGER_API_ERROR", fmt.Sprintf("Jaeger API error (status %d): %s", resp.StatusCode, string(body))).
				WithHTTPStatus(resp.StatusCode).
				WithContext("status_code", resp.StatusCode)
		}
	}

	return body, nil
}

// GetTrace retrieves a specific trace by ID.
func (c *Client) GetTrace(ctx context.Context, traceID string) (*Trace, error) {
	logrus.WithField("trace_id", traceID).Debug("Getting Jaeger trace")

	resp, err := c.makeRequest(ctx, "GET", fmt.Sprintf("api/traces/%s", traceID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data []Trace `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal trace: %w", err)
	}

	if len(result.Data) == 0 {
		return nil, errors.NotFoundError(fmt.Sprintf("trace: %s", traceID)).
			WithContext("trace_id", traceID)
	}

	return &result.Data[0], nil
}

// SearchTraces searches for traces based on query parameters.
func (c *Client) SearchTraces(ctx context.Context, params TraceQueryParameters) ([]Trace, error) {
	logrus.Debug("Searching Jaeger traces")

	query := url.Values{}
	if params.Service != "" {
		query.Set("service", params.Service)
	}
	if params.Operation != "" {
		query.Set("operation", params.Operation)
	}
	if params.StartTime != "" {
		query.Set("start", params.StartTime)
	}
	if params.EndTime != "" {
		query.Set("end", params.EndTime)
	}
	if params.Limit > 0 {
		query.Set("limit", fmt.Sprintf("%d", params.Limit))
	}
	if params.MinDuration != "" {
		query.Set("minDuration", params.MinDuration)
	}
	if params.MaxDuration != "" {
		query.Set("maxDuration", params.MaxDuration)
	}

	// Add tags
	for key, value := range params.Tags {
		query.Set(fmt.Sprintf("tag:%s", key), value)
	}

	resp, err := c.makeRequest(ctx, "GET", fmt.Sprintf("api/traces?%s", query.Encode()), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data []Trace `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal traces: %w", err)
	}

	return result.Data, nil
}

// GetServices retrieves all services from Jaeger.
func (c *Client) GetServices(ctx context.Context) ([]Service, error) {
	logrus.Debug("Getting Jaeger services")

	resp, err := c.makeRequest(ctx, "GET", "api/services", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data []string `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal services: %w", err)
	}

	services := make([]Service, len(result.Data))
	for i, name := range result.Data {
		services[i] = Service{Name: name}
	}

	logrus.WithField("count", len(services)).Debug("Retrieved services")
	return services, nil
}

// GetOperations retrieves operations for a specific service.
func (c *Client) GetOperations(ctx context.Context, service string) ([]string, error) {
	logrus.WithField("service", service).Debug("Getting Jaeger operations")

	resp, err := c.makeRequest(ctx, "GET", fmt.Sprintf("api/operations?service=%s", service), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data []string `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal operations: %w", err)
	}

	return result.Data, nil
}

// GetDependencies retrieves service dependencies.
func (c *Client) GetDependencies(ctx context.Context, startTime, endTime string) ([]Dependency, error) {
	logrus.Debug("Getting Jaeger dependencies")

	query := url.Values{}
	if startTime != "" {
		query.Set("startTime", startTime)
	}
	if endTime != "" {
		query.Set("endTime", endTime)
	}

	resp, err := c.makeRequest(ctx, "GET", fmt.Sprintf("api/dependencies?%s", query.Encode()), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result []Dependency
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dependencies: %w", err)
	}

	return result, nil
}

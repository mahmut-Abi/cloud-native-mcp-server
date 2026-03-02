package client

import (
	"bytes"
	"context"
	"encoding/json"
	stderrs "errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/errors"
	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
	"github.com/sirupsen/logrus"
)

const (
	defaultRequestTimeout = 30 * time.Second
	defaultMaxRetries     = 2
	defaultRetryDelay     = 200 * time.Millisecond
	defaultRetryMaxDelay  = 2 * time.Second
)

// ClientOptions holds configuration parameters for creating a Jaeger client.
type ClientOptions struct {
	BaseURL        string        // Jaeger server base URL
	Timeout        time.Duration // HTTP request timeout
	MaxRetries     int           // Retries after the first request attempt
	RetryBaseDelay time.Duration // Base delay for exponential backoff
	RetryMaxDelay  time.Duration // Maximum delay between retries
}

// Client provides operations for interacting with Jaeger API.
type Client struct {
	baseURL        string
	httpClient     *http.Client
	maxRetries     int
	retryBaseDelay time.Duration
	retryMaxDelay  time.Duration
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
	if opts == nil {
		return nil, fmt.Errorf("jaeger client options are required")
	}

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
		timeout = defaultRequestTimeout
	}

	maxRetries := opts.MaxRetries
	if maxRetries < 0 {
		maxRetries = 0
	}
	if maxRetries == 0 {
		maxRetries = defaultMaxRetries
	}

	retryBaseDelay := opts.RetryBaseDelay
	if retryBaseDelay <= 0 {
		retryBaseDelay = defaultRetryDelay
	}

	retryMaxDelay := opts.RetryMaxDelay
	if retryMaxDelay <= 0 {
		retryMaxDelay = defaultRetryMaxDelay
	}
	if retryMaxDelay < retryBaseDelay {
		retryMaxDelay = retryBaseDelay
	}

	httpClient := optimize.NewOptimizedHTTPClientWithTimeout(timeout)

	client := &Client{
		baseURL:        baseURL.String(),
		httpClient:     httpClient,
		maxRetries:     maxRetries,
		retryBaseDelay: retryBaseDelay,
		retryMaxDelay:  retryMaxDelay,
	}

	return client, nil
}

// makeRequest performs an HTTP request to the Jaeger API.
func (c *Client) makeRequest(ctx context.Context, method, endpoint string, body io.Reader) (*http.Response, error) {
	requestURL := c.baseURL + strings.TrimPrefix(endpoint, "/")
	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = io.ReadAll(body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
	}

	totalAttempts := c.maxRetries + 1
	for attempt := 1; attempt <= totalAttempts; attempt++ {
		var reqBody io.Reader
		if len(bodyBytes) > 0 {
			reqBody = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, method, requestURL, reqBody)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Accept", "application/json")
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		logrus.WithFields(logrus.Fields{
			"method":  method,
			"url":     requestURL,
			"attempt": attempt,
		}).Debug("Making Jaeger API request")

		resp, reqErr := c.httpClient.Do(req)
		if reqErr == nil {
			if c.shouldRetryStatusCode(resp.StatusCode) && attempt < totalAttempts {
				_ = resp.Body.Close()
				delay := c.nextRetryDelay(attempt)
				logrus.WithFields(logrus.Fields{
					"method":      method,
					"url":         requestURL,
					"status_code": resp.StatusCode,
					"attempt":     attempt,
					"retry_in":    delay,
				}).Warn("Retrying Jaeger API request after retryable status")
				if waitErr := waitForRetry(ctx, delay); waitErr != nil {
					return nil, waitErr
				}
				continue
			}
			return resp, nil
		}

		if attempt < totalAttempts && c.shouldRetryError(reqErr) {
			delay := c.nextRetryDelay(attempt)
			logrus.WithFields(logrus.Fields{
				"method":   method,
				"url":      requestURL,
				"attempt":  attempt,
				"retry_in": delay,
			}).WithError(reqErr).Warn("Retrying Jaeger API request after transient transport error")
			if waitErr := waitForRetry(ctx, delay); waitErr != nil {
				return nil, waitErr
			}
			continue
		}

		return nil, errors.Wrap(reqErr, "JAEGER_CONNECTION_FAILED", "failed to connect to Jaeger").
			WithHTTPStatus(503)
	}

	return nil, errors.New("JAEGER_RETRY_EXHAUSTED", "retry attempts exhausted for Jaeger API request").
		WithHTTPStatus(503).
		WithContext("url", requestURL).
		WithContext("max_retries", c.maxRetries)
}

func (c *Client) shouldRetryError(err error) bool {
	if err == nil {
		return false
	}
	if stderrs.Is(err, context.Canceled) {
		return false
	}
	var netErr net.Error
	if stderrs.As(err, &netErr) {
		return netErr.Timeout() || netErr.Temporary()
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "timeout")
}

func (c *Client) shouldRetryStatusCode(statusCode int) bool {
	switch statusCode {
	case http.StatusRequestTimeout, http.StatusTooManyRequests, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

func (c *Client) nextRetryDelay(attempt int) time.Duration {
	delay := c.retryBaseDelay
	for i := 1; i < attempt; i++ {
		delay *= 2
		if delay >= c.retryMaxDelay {
			delay = c.retryMaxDelay
			break
		}
	}
	if delay > c.retryMaxDelay {
		delay = c.retryMaxDelay
	}

	// Add bounded jitter so repeated calls don't synchronize retries.
	jitter := time.Duration(rand.Int63n(int64(delay/4 + 1)))
	return delay + jitter
}

func waitForRetry(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
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

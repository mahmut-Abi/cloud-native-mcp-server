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

// ClientOptions holds configuration parameters for creating a Loki client.
type ClientOptions struct {
	Address        string
	Username       string
	Password       string
	BearerToken    string
	Timeout        time.Duration
	TLSSkipVerify  bool
	TLSCertFile    string
	TLSKeyFile     string
	TLSCAFile      string
	MaxRetries     int
	RetryBaseDelay time.Duration
	RetryMaxDelay  time.Duration
}

// Client provides operations for interacting with Loki.
type Client struct {
	baseURL        string
	httpClient     *http.Client
	username       string
	password       string
	bearerToken    string
	headers        map[string]string
	maxRetries     int
	retryBaseDelay time.Duration
	retryMaxDelay  time.Duration
}

// NewClient creates a new Loki client.
func NewClient(opts *ClientOptions) (*Client, error) {
	if opts == nil {
		return nil, fmt.Errorf("loki client options are required")
	}
	if opts.Address == "" {
		return nil, fmt.Errorf("loki address is required")
	}

	baseURL, err := url.Parse(opts.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid loki address: %w", err)
	}
	if !strings.HasSuffix(baseURL.Path, "/") {
		baseURL.Path += "/"
	}
	baseURL.Path += "loki/api/v1/"

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	httpClient := optimize.NewOptimizedHTTPClientWithTimeout(timeout)
	if transport, ok := httpClient.Transport.(*http.Transport); ok {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: opts.TLSSkipVerify,
		}
		if opts.TLSCertFile != "" && opts.TLSKeyFile != "" {
			cert, err := tls.LoadX509KeyPair(opts.TLSCertFile, opts.TLSKeyFile)
			if err != nil {
				return nil, fmt.Errorf("failed to load TLS certificate: %w", err)
			}
			transport.TLSClientConfig.Certificates = []tls.Certificate{cert}
		}
	}

	client := &Client{
		baseURL:     baseURL.String(),
		httpClient:  httpClient,
		username:    opts.Username,
		password:    opts.Password,
		bearerToken: opts.BearerToken,
		headers: map[string]string{
			"Accept": "application/json",
		},
		maxRetries:     opts.MaxRetries,
		retryBaseDelay: opts.RetryBaseDelay,
		retryMaxDelay:  opts.RetryMaxDelay,
	}
	client.maxRetries, client.retryBaseDelay, client.retryMaxDelay = optimize.NormalizeRetryConfig(
		client.maxRetries,
		client.retryBaseDelay,
		client.retryMaxDelay,
	)

	return client, nil
}

func (c *Client) makeRequest(ctx context.Context, endpoint string, params url.Values) (*http.Response, error) {
	requestURL := c.baseURL + endpoint
	if len(params) > 0 {
		requestURL += "?" + params.Encode()
	}

	return optimize.DoWithHTTPRetry(
		ctx,
		http.MethodGet,
		c.maxRetries,
		c.retryBaseDelay,
		c.retryMaxDelay,
		func(attempt int) (*http.Response, error) {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create request: %w", err)
			}

			for key, value := range c.headers {
				req.Header.Set(key, value)
			}
			if c.bearerToken != "" {
				req.Header.Set("Authorization", "Bearer "+c.bearerToken)
			} else if c.username != "" && c.password != "" {
				req.SetBasicAuth(c.username, c.password)
			}

			logrus.WithFields(logrus.Fields{
				"attempt": attempt,
				"url":     requestURL,
			}).Debug("Making Loki API request")

			return c.httpClient.Do(req)
		},
		func(event optimize.HTTPRetryEvent) {
			fields := logrus.Fields{
				"url":      requestURL,
				"attempt":  event.Attempt,
				"retry_in": event.Delay,
			}
			if event.Err != nil {
				logrus.WithFields(fields).WithError(event.Err).Warn("Retrying Loki API request after transient transport error")
				return
			}
			fields["status_code"] = event.StatusCode
			logrus.WithFields(fields).Warn("Retrying Loki API request after retryable status")
		},
	)
}

func (c *Client) readResponse(resp *http.Response, out interface{}) error {
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("loki API error (status %d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("failed to unmarshal Loki response: %w", err)
	}
	return nil
}

// Query executes a Loki instant query.
func (c *Client) Query(ctx context.Context, query string, queryTime *time.Time, limit int, direction string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("query", query)
	if queryTime != nil {
		params.Set("time", strconv.FormatInt(queryTime.UnixNano(), 10))
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if direction != "" {
		params.Set("direction", direction)
	}

	resp, err := c.makeRequest(ctx, "query", params)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := c.readResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// QueryRange executes a Loki range query.
func (c *Client) QueryRange(ctx context.Context, query string, start, end time.Time, limit int, direction, step string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("start", strconv.FormatInt(start.UnixNano(), 10))
	params.Set("end", strconv.FormatInt(end.UnixNano(), 10))
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if direction != "" {
		params.Set("direction", direction)
	}
	if step != "" {
		params.Set("step", step)
	}

	resp, err := c.makeRequest(ctx, "query_range", params)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := c.readResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetLabelNames retrieves Loki label names.
func (c *Client) GetLabelNames(ctx context.Context, query string, start, end *time.Time) ([]string, error) {
	params := url.Values{}
	if query != "" {
		params.Set("query", query)
	}
	if start != nil {
		params.Set("start", strconv.FormatInt(start.UnixNano(), 10))
	}
	if end != nil {
		params.Set("end", strconv.FormatInt(end.UnixNano(), 10))
	}

	resp, err := c.makeRequest(ctx, "labels", params)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status string   `json:"status"`
		Data   []string `json:"data"`
	}
	if err := c.readResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetLabelValues retrieves values for a Loki label.
func (c *Client) GetLabelValues(ctx context.Context, labelName, query string, start, end *time.Time) ([]string, error) {
	params := url.Values{}
	if query != "" {
		params.Set("query", query)
	}
	if start != nil {
		params.Set("start", strconv.FormatInt(start.UnixNano(), 10))
	}
	if end != nil {
		params.Set("end", strconv.FormatInt(end.UnixNano(), 10))
	}

	resp, err := c.makeRequest(ctx, "label/"+url.PathEscape(labelName)+"/values", params)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status string   `json:"status"`
		Data   []string `json:"data"`
	}
	if err := c.readResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// GetSeries retrieves Loki series for one or more selectors.
func (c *Client) GetSeries(ctx context.Context, matchers []string, start, end *time.Time) ([]map[string]string, error) {
	params := url.Values{}
	for _, matcher := range matchers {
		matcher = strings.TrimSpace(matcher)
		if matcher != "" {
			params.Add("match[]", matcher)
		}
	}
	if start != nil {
		params.Set("start", strconv.FormatInt(start.UnixNano(), 10))
	}
	if end != nil {
		params.Set("end", strconv.FormatInt(end.UnixNano(), 10))
	}

	resp, err := c.makeRequest(ctx, "series", params)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status string              `json:"status"`
		Data   []map[string]string `json:"data"`
	}
	if err := c.readResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// TestConnection checks whether Loki is reachable.
func (c *Client) TestConnection(ctx context.Context) error {
	resp, err := c.makeRequest(ctx, "labels", nil)
	if err != nil {
		return err
	}

	var result struct {
		Status string `json:"status"`
	}
	return c.readResponse(resp, &result)
}

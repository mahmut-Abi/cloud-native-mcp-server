package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
	"github.com/sirupsen/logrus"
)

const defaultRequestTimeout = 30 * time.Second

// ClientOptions holds configuration for the Nacos client.
type ClientOptions struct {
	URL            string
	Username       string
	Password       string
	AccessToken    string
	Timeout        time.Duration
	MaxRetries     int
	RetryBaseDelay time.Duration
	RetryMaxDelay  time.Duration
}

// Client provides read-only access to the Nacos OpenAPI.
type Client struct {
	baseURL        *url.URL
	httpClient     *http.Client
	username       string
	password       string
	accessToken    string
	tokenExpiry    time.Time
	tokenMutex     sync.Mutex
	maxRetries     int
	retryBaseDelay time.Duration
	retryMaxDelay  time.Duration
}

// NewClient creates a Nacos client.
func NewClient(opts *ClientOptions) (*Client, error) {
	if opts == nil {
		return nil, fmt.Errorf("nacos client options are required")
	}
	if strings.TrimSpace(opts.URL) == "" {
		return nil, fmt.Errorf("nacos URL is required")
	}

	parsedURL, err := url.Parse(strings.TrimSpace(opts.URL))
	if err != nil {
		return nil, fmt.Errorf("invalid nacos URL: %w", err)
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, fmt.Errorf("invalid nacos URL: missing scheme or host")
	}

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = defaultRequestTimeout
	}

	maxRetries, retryBaseDelay, retryMaxDelay := optimize.NormalizeRetryConfig(
		opts.MaxRetries,
		opts.RetryBaseDelay,
		opts.RetryMaxDelay,
	)

	return &Client{
		baseURL:        parsedURL,
		httpClient:     optimize.NewOptimizedHTTPClientWithTimeout(timeout),
		username:       strings.TrimSpace(opts.Username),
		password:       strings.TrimSpace(opts.Password),
		accessToken:    strings.TrimSpace(opts.AccessToken),
		maxRetries:     maxRetries,
		retryBaseDelay: retryBaseDelay,
		retryMaxDelay:  retryMaxDelay,
	}, nil
}

// ListNamespaces returns all namespaces visible to the current client.
func (c *Client) ListNamespaces(ctx context.Context) ([]map[string]interface{}, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, "/v1/console/namespaces", nil, nil, true)
	if err != nil {
		return nil, err
	}

	payload, err := c.readJSONMap(resp)
	if err != nil {
		return nil, err
	}

	for _, key := range []string{"data", "pageItems", "namespaces"} {
		if items, ok := mapSlice(payload[key]); ok {
			return items, nil
		}
	}
	return []map[string]interface{}{}, nil
}

// ListConfigs returns configuration summaries from Nacos.
func (c *Client) ListConfigs(ctx context.Context, params url.Values) (map[string]interface{}, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, "/v1/cs/configs", params, nil, true)
	if err != nil {
		return nil, err
	}
	return c.readJSONMap(resp)
}

// GetConfig returns one configuration entry from Nacos.
func (c *Client) GetConfig(ctx context.Context, params url.Values) (map[string]interface{}, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, "/v1/cs/configs", params, nil, true)
	if err != nil {
		return nil, err
	}

	body, contentType, err := c.readResponse(resp)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"dataId":   params.Get("dataId"),
		"group":    params.Get("group"),
		"tenant":   params.Get("tenant"),
		"content":  strings.TrimSpace(string(body)),
		"format":   inferConfigFormat(params.Get("dataId"), contentType, body),
		"metadata": map[string]interface{}{},
	}
	return result, nil
}

// ListServices returns service summaries from Nacos.
func (c *Client) ListServices(ctx context.Context, params url.Values) (map[string]interface{}, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, "/v1/ns/service/list", params, nil, true)
	if err != nil {
		return nil, err
	}
	return c.readJSONMap(resp)
}

// GetService returns one service from Nacos.
func (c *Client) GetService(ctx context.Context, params url.Values) (map[string]interface{}, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, "/v1/ns/service", params, nil, true)
	if err != nil {
		return nil, err
	}
	return c.readJSONMap(resp)
}

// ListInstances returns instances for one service.
func (c *Client) ListInstances(ctx context.Context, params url.Values) (map[string]interface{}, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, "/v1/ns/instance/list", params, nil, true)
	if err != nil {
		return nil, err
	}
	return c.readJSONMap(resp)
}

// ListClusterNodes returns Nacos cluster node information.
func (c *Client) ListClusterNodes(ctx context.Context) (map[string]interface{}, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, "/v1/ns/operator/servers", nil, nil, true)
	if err != nil {
		return nil, err
	}
	return c.readJSONMap(resp)
}

// GetSystemMetrics returns Nacos server metrics.
func (c *Client) GetSystemMetrics(ctx context.Context) (map[string]interface{}, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, "/v1/ns/operator/metrics", nil, nil, true)
	if err != nil {
		return nil, err
	}
	return c.readJSONMap(resp)
}

// TestConnection verifies that Nacos is reachable and authentication works.
func (c *Client) TestConnection(ctx context.Context) (map[string]interface{}, error) {
	namespaces, err := c.ListNamespaces(ctx)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"status":           "ok",
		"namespaceCount":   len(namespaces),
		"authenticated":    c.username != "" || c.accessToken != "",
		"baseURL":          c.baseURL.String(),
		"sampleNamespaces": sampleMaps(namespaces, 5),
	}, nil
}

func (c *Client) makeRequest(ctx context.Context, method, endpoint string, params url.Values, body io.Reader, withAuthRetry bool) (*http.Response, error) {
	makeOne := func() (*http.Response, error) {
		u := *c.baseURL
		u.Path = joinURLPath(c.baseURL.Path, endpoint)
		if len(params) > 0 {
			q := u.Query()
			for key, values := range params {
				for _, value := range values {
					q.Add(key, value)
				}
			}
			u.RawQuery = q.Encode()
		}

		return optimize.DoWithHTTPRetry(
			ctx,
			method,
			c.maxRetries,
			c.retryBaseDelay,
			c.retryMaxDelay,
			func(attempt int) (*http.Response, error) {
				req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
				if err != nil {
					return nil, fmt.Errorf("failed to create request: %w", err)
				}
				req.Header.Set("Accept", "application/json")

				if method == http.MethodPost && body != nil {
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				}

				token, err := c.ensureAccessToken(ctx)
				if err != nil {
					return nil, err
				}
				if token != "" {
					q := req.URL.Query()
					q.Set("accessToken", token)
					req.URL.RawQuery = q.Encode()
				}

				logrus.WithFields(logrus.Fields{
					"attempt": attempt,
					"method":  method,
					"url":     req.URL.String(),
				}).Debug("Making Nacos API request")

				return c.httpClient.Do(req)
			},
			func(event optimize.HTTPRetryEvent) {
				fields := logrus.Fields{
					"attempt":  event.Attempt,
					"method":   method,
					"endpoint": endpoint,
					"retry_in": event.Delay,
				}
				if event.Err != nil {
					logrus.WithFields(fields).WithError(event.Err).Warn("Retrying Nacos API request after transient transport error")
					return
				}
				fields["status_code"] = event.StatusCode
				logrus.WithFields(fields).Warn("Retrying Nacos API request after retryable status")
			},
		)
	}

	resp, err := makeOne()
	if err != nil {
		return nil, err
	}
	if !withAuthRetry || resp.StatusCode != http.StatusForbidden {
		return resp, nil
	}

	bodyBytes, _, readErr := c.readResponse(resp)
	if readErr != nil {
		return nil, readErr
	}
	if !isTokenExpiredResponse(bodyBytes) {
		return responseFromBytes(resp.StatusCode, resp.Header, bodyBytes), nil
	}

	c.clearToken()
	return makeOne()
}

func (c *Client) ensureAccessToken(ctx context.Context) (string, error) {
	if c.accessToken != "" && (c.tokenExpiry.IsZero() || time.Now().Before(c.tokenExpiry)) {
		return c.accessToken, nil
	}
	if c.username == "" || c.password == "" {
		return c.accessToken, nil
	}

	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	if c.accessToken != "" && (c.tokenExpiry.IsZero() || time.Now().Before(c.tokenExpiry)) {
		return c.accessToken, nil
	}

	token, expiry, err := c.login(ctx)
	if err != nil {
		return "", err
	}
	c.accessToken = token
	c.tokenExpiry = expiry
	return c.accessToken, nil
}

func (c *Client) login(ctx context.Context) (string, time.Time, error) {
	form := url.Values{}
	form.Set("username", c.username)
	form.Set("password", c.password)

	endpoints := []string{"/v1/auth/users/login", "/v1/auth/login"}
	var lastErr error
	for _, endpoint := range endpoints {
		u := *c.baseURL
		u.Path = joinURLPath(c.baseURL.Path, endpoint)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBufferString(form.Encode()))
		if err != nil {
			return "", time.Time{}, fmt.Errorf("failed to create login request: %w", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		body, _, err := c.readResponse(resp)
		if err != nil {
			lastErr = err
			continue
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			lastErr = fmt.Errorf("failed to decode nacos login response: %w", err)
			continue
		}

		token := firstString(payload, "accessToken", "token")
		if token == "" {
			lastErr = fmt.Errorf("nacos login response did not include access token")
			continue
		}

		expiry := time.Time{}
		if ttlStr := firstString(payload, "tokenTtl"); ttlStr != "" {
			if ttl, err := strconv.Atoi(ttlStr); err == nil && ttl > 0 {
				expiry = time.Now().Add(time.Duration(ttl-30) * time.Second)
			}
		}
		if ttlFloat, ok := payload["tokenTtl"].(float64); ok && ttlFloat > 0 {
			expiry = time.Now().Add(time.Duration(int(ttlFloat)-30) * time.Second)
		}

		return token, expiry, nil
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("nacos login failed")
	}
	return "", time.Time{}, lastErr
}

func (c *Client) readJSONMap(resp *http.Response) (map[string]interface{}, error) {
	body, _, err := c.readResponse(resp)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode nacos response: %w", err)
	}
	return result, nil
}

func (c *Client) readResponse(resp *http.Response) ([]byte, string, error) {
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read nacos response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, "", fmt.Errorf("nacos API error (status %d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return body, resp.Header.Get("Content-Type"), nil
}

func (c *Client) clearToken() {
	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()
	c.accessToken = ""
	c.tokenExpiry = time.Time{}
}

func joinURLPath(basePath, endpoint string) string {
	if strings.TrimSpace(basePath) == "" {
		return endpoint
	}
	return path.Join(strings.TrimSuffix(basePath, "/"), endpoint)
}

func firstString(values map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if value, ok := values[key].(string); ok && strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func inferConfigFormat(dataID, contentType string, body []byte) string {
	lowerDataID := strings.ToLower(dataID)
	switch {
	case strings.HasSuffix(lowerDataID, ".yaml"), strings.HasSuffix(lowerDataID, ".yml"):
		return "yaml"
	case strings.HasSuffix(lowerDataID, ".json"):
		return "json"
	case strings.HasSuffix(lowerDataID, ".properties"):
		return "properties"
	case strings.Contains(strings.ToLower(contentType), "json"):
		return "json"
	case strings.Contains(strings.ToLower(contentType), "yaml"):
		return "yaml"
	case bytes.Contains(body, []byte("=")):
		return "properties"
	default:
		return "text"
	}
}

func mapSlice(value interface{}) ([]map[string]interface{}, bool) {
	items, ok := value.([]interface{})
	if !ok {
		return nil, false
	}
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		if typed, ok := item.(map[string]interface{}); ok {
			result = append(result, typed)
		}
	}
	return result, true
}

func sampleMaps(items []map[string]interface{}, limit int) []map[string]interface{} {
	if len(items) <= limit {
		return items
	}
	return items[:limit]
}

func isTokenExpiredResponse(body []byte) bool {
	lower := strings.ToLower(string(body))
	return strings.Contains(lower, "token invalid") ||
		strings.Contains(lower, "token expired") ||
		strings.Contains(lower, "invalid token")
}

func responseFromBytes(status int, headers http.Header, body []byte) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     headers.Clone(),
		Body:       io.NopCloser(bytes.NewReader(body)),
	}
}

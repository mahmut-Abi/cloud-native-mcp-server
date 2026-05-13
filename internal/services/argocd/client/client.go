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
	"strings"
	"sync"
	"time"

	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
	"github.com/sirupsen/logrus"
)

const defaultRequestTimeout = 30 * time.Second

// ClientOptions holds configuration for the Argo CD client.
type ClientOptions struct {
	URL            string
	Username       string
	Password       string
	AuthToken      string
	Timeout        time.Duration
	MaxRetries     int
	RetryBaseDelay time.Duration
	RetryMaxDelay  time.Duration
}

// Client provides read-only access to the Argo CD API.
type Client struct {
	baseURL        *url.URL
	httpClient     *http.Client
	username       string
	password       string
	authToken      string
	tokenMutex     sync.Mutex
	maxRetries     int
	retryBaseDelay time.Duration
	retryMaxDelay  time.Duration
}

// NewClient creates a new Argo CD client.
func NewClient(opts *ClientOptions) (*Client, error) {
	if opts == nil {
		return nil, fmt.Errorf("argocd client options are required")
	}
	if strings.TrimSpace(opts.URL) == "" {
		return nil, fmt.Errorf("argocd URL is required")
	}

	parsedURL, err := url.Parse(strings.TrimSpace(opts.URL))
	if err != nil {
		return nil, fmt.Errorf("invalid argocd URL: %w", err)
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, fmt.Errorf("invalid argocd URL: missing scheme or host")
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
		authToken:      strings.TrimSpace(opts.AuthToken),
		maxRetries:     maxRetries,
		retryBaseDelay: retryBaseDelay,
		retryMaxDelay:  retryMaxDelay,
	}, nil
}

// TestConnection verifies that Argo CD is reachable and auth works.
func (c *Client) TestConnection(ctx context.Context) (map[string]interface{}, error) {
	apps, err := c.ListApplications(ctx, url.Values{})
	if err != nil {
		return nil, err
	}
	items, _ := ItemsSlice(apps)
	return map[string]interface{}{
		"status":           "ok",
		"applicationCount": len(items),
		"authenticated":    c.username != "" || c.authToken != "",
	}, nil
}

// ListApplications returns the Argo CD application list payload.
func (c *Client) ListApplications(ctx context.Context, params url.Values) (map[string]interface{}, error) {
	return c.getJSON(ctx, "/api/v1/applications", params)
}

// GetApplication returns one application.
func (c *Client) GetApplication(ctx context.Context, name string, params url.Values) (map[string]interface{}, error) {
	return c.getJSON(ctx, fmt.Sprintf("/api/v1/applications/%s", url.PathEscape(strings.TrimSpace(name))), params)
}

// GetApplicationManifests returns manifests for one application.
func (c *Client) GetApplicationManifests(ctx context.Context, name string, params url.Values) (map[string]interface{}, error) {
	return c.getJSON(ctx, fmt.Sprintf("/api/v1/applications/%s/manifests", url.PathEscape(strings.TrimSpace(name))), params)
}

// ListProjects returns projects list.
func (c *Client) ListProjects(ctx context.Context) (map[string]interface{}, error) {
	return c.getJSON(ctx, "/api/v1/projects", nil)
}

// GetProject returns one project.
func (c *Client) GetProject(ctx context.Context, name string) (map[string]interface{}, error) {
	return c.getJSON(ctx, fmt.Sprintf("/api/v1/projects/%s", url.PathEscape(strings.TrimSpace(name))), nil)
}

// ListClusters returns clusters list.
func (c *Client) ListClusters(ctx context.Context) (map[string]interface{}, error) {
	return c.getJSON(ctx, "/api/v1/clusters", nil)
}

func (c *Client) getJSON(ctx context.Context, endpoint string, params url.Values) (map[string]interface{}, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, endpoint, params, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read argocd response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("argocd API error (status %d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var result map[string]interface{}
	if len(body) == 0 {
		return map[string]interface{}{}, nil
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode argocd response: %w", err)
	}
	return result, nil
}

func (c *Client) makeRequest(ctx context.Context, method, endpoint string, params url.Values, body []byte) (*http.Response, error) {
	return optimize.DoWithHTTPRetry(
		ctx,
		method,
		c.maxRetries,
		c.retryBaseDelay,
		c.retryMaxDelay,
		func(attempt int) (*http.Response, error) {
			reqURL := *c.baseURL
			reqURL.Path = joinURLPath(c.baseURL.Path, endpoint)
			if len(params) > 0 {
				q := reqURL.Query()
				for key, values := range params {
					for _, value := range values {
						q.Add(key, value)
					}
				}
				reqURL.RawQuery = q.Encode()
			}

			var reader io.Reader
			if len(body) > 0 {
				reader = bytes.NewReader(body)
			}
			req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), reader)
			if err != nil {
				return nil, fmt.Errorf("failed to create request: %w", err)
			}

			req.Header.Set("Accept", "application/json")
			if len(body) > 0 {
				req.Header.Set("Content-Type", "application/json")
			}

			token, err := c.ensureToken(ctx)
			if err != nil {
				return nil, err
			}
			if token != "" {
				req.Header.Set("Authorization", "Bearer "+token)
			}

			logrus.WithFields(logrus.Fields{
				"attempt": attempt,
				"method":  method,
				"url":     req.URL.String(),
			}).Debug("Making Argo CD API request")

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
				logrus.WithFields(fields).WithError(event.Err).Warn("Retrying Argo CD API request after transient transport error")
				return
			}
			fields["status_code"] = event.StatusCode
			logrus.WithFields(fields).Warn("Retrying Argo CD API request after retryable status")
		},
	)
}

func (c *Client) ensureToken(ctx context.Context) (string, error) {
	if c.authToken != "" {
		return c.authToken, nil
	}
	if c.username == "" || c.password == "" {
		return "", nil
	}

	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	if c.authToken != "" {
		return c.authToken, nil
	}

	payload, err := json.Marshal(map[string]string{
		"username": c.username,
		"password": c.password,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal argocd session payload: %w", err)
	}

	reqURL := *c.baseURL
	reqURL.Path = joinURLPath(c.baseURL.Path, "/api/v1/session")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL.String(), bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("failed to create argocd session request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read argocd session response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("argocd session error (status %d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to decode argocd session response: %w", err)
	}
	token, ok := result["token"].(string)
	if !ok || strings.TrimSpace(token) == "" {
		return "", fmt.Errorf("argocd session response did not include token")
	}
	c.authToken = strings.TrimSpace(token)
	return c.authToken, nil
}

// ItemsSlice extracts Kubernetes-style list items from an Argo CD list response.
func ItemsSlice(payload map[string]interface{}) ([]map[string]interface{}, bool) {
	items, ok := payload["items"].([]interface{})
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

func joinURLPath(basePath, endpoint string) string {
	if strings.TrimSpace(basePath) == "" {
		return endpoint
	}
	return path.Join(strings.TrimSuffix(basePath, "/"), endpoint)
}

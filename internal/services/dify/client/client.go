package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
	"github.com/sirupsen/logrus"
)

const defaultRequestTimeout = 30 * time.Second

// PaginationCursor describes a cursor from a Link header.
type PaginationCursor struct {
	Cursor  string `json:"cursor"`
	Rel     string `json:"rel"`
	Results bool   `json:"results"`
	URL     string `json:"url,omitempty"`
}

// Pagination captures cursor pagination metadata from Link headers.
type Pagination struct {
	Next *PaginationCursor `json:"next,omitempty"`
	Prev *PaginationCursor `json:"prev,omitempty"`
}

// ClientOptions holds configuration for creating a Dify client.
type ClientOptions struct {
	ConsoleURL      string
	ConsoleEmail    string
	ConsolePassword string
	ServiceURL      string
	APIKey          string
	Timeout         time.Duration
	MaxRetries      int
	RetryBaseDelay  time.Duration
	RetryMaxDelay   time.Duration
}

// Client provides access to the Dify Console API (session-based auth)
// and the Dify Service API (Bearer token auth).
type Client struct {
	httpClient     *http.Client
	consoleURL     string
	consoleEmail   string
	consolePassword string
	serviceURL     string
	apiKey         string
	maxRetries     int
	retryBaseDelay time.Duration
	retryMaxDelay  time.Duration

	mu        sync.Mutex
	cookies   map[string]string
	csrfToken string
}

// NewClient creates a new Dify client.
// At least one of ConsoleURL or ServiceURL must be provided.
// If ConsoleURL is provided, ConsoleEmail and ConsolePassword are required.
// If ServiceURL is provided, APIKey is required.
func NewClient(opts *ClientOptions) (*Client, error) {
	if opts == nil {
		return nil, fmt.Errorf("dify client options are required")
	}
	if strings.TrimSpace(opts.ConsoleURL) == "" && strings.TrimSpace(opts.ServiceURL) == "" {
		return nil, fmt.Errorf("dify console URL or service URL is required")
	}
	if opts.ConsoleURL != "" {
		if strings.TrimSpace(opts.ConsoleEmail) == "" {
			return nil, fmt.Errorf("dify console email is required when console URL is provided")
		}
		if strings.TrimSpace(opts.ConsolePassword) == "" {
			return nil, fmt.Errorf("dify console password is required when console URL is provided")
		}
	}
	if opts.ServiceURL != "" && strings.TrimSpace(opts.APIKey) == "" {
		return nil, fmt.Errorf("dify API key is required when service URL is provided")
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

	cli := &Client{
		httpClient:      optimize.NewOptimizedHTTPClientWithTimeout(timeout),
		consoleURL:      strings.TrimRight(strings.TrimSpace(opts.ConsoleURL), "/"),
		consoleEmail:    strings.TrimSpace(opts.ConsoleEmail),
		consolePassword: opts.ConsolePassword,
		serviceURL:      strings.TrimRight(strings.TrimSpace(opts.ServiceURL), "/"),
		apiKey:          strings.TrimSpace(opts.APIKey),
		maxRetries:      maxRetries,
		retryBaseDelay:  retryBaseDelay,
		retryMaxDelay:   retryMaxDelay,
	}

	return cli, nil
}

// ensureConsoleSession handles the Dify Console login flow.
// It acquires a mutex to prevent concurrent login attempts and
// stores the session cookies and CSRF token on the client struct.
func (c *Client) ensureConsoleSession(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cookies != nil && c.csrfToken != "" {
		return nil
	}

	logrus.Debug("Starting Dify Console login flow")

	loginURL := c.consoleURL + "/login"

	loginBody := map[string]interface{}{
		"email":       c.consoleEmail,
		"password":    base64.StdEncoding.EncodeToString([]byte(c.consolePassword)),
		"remember_me": true,
	}
	bodyBytes, err := json.Marshal(loginBody)
	if err != nil {
		return fmt.Errorf("failed to marshal login body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, loginURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("console login request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read login response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("console login failed (status %d): %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	c.cookies = make(map[string]string)
	for _, cookie := range resp.Cookies() {
		c.cookies[cookie.Name] = cookie.Value
	}

	if token, ok := c.cookies["csrf_token"]; ok {
		c.csrfToken = token
	} else if token, ok := c.cookies["__Host-csrf_token"]; ok {
		c.csrfToken = token
	}

	if c.csrfToken == "" {
		return fmt.Errorf("console login succeeded but no csrf_token cookie was returned")
	}

	logrus.Debug("Dify Console login completed successfully")
	return nil
}

// ConsoleRequest sends a request to the Dify Console API using cookie-based
// session authentication. It automatically performs the login flow if no
// active session exists.
func (c *Client) ConsoleRequest(ctx context.Context, method, path string, query url.Values, body interface{}) (json.RawMessage, error) {
	if c.consoleURL == "" {
		return nil, fmt.Errorf("console URL is not configured")
	}

	if err := c.ensureConsoleSession(ctx); err != nil {
		return nil, fmt.Errorf("failed to establish console session: %w", err)
	}

	c.mu.Lock()
	cookieStr := buildCookieHeader(c.cookies)
	csrf := c.csrfToken
	c.mu.Unlock()

	requestURL := c.consoleURL + "/" + strings.TrimPrefix(path, "/")
	if len(query) > 0 {
		requestURL += "?" + query.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	extraHeaders := map[string]string{
		"X-CSRF-Token": csrf,
		"Cookie":       cookieStr,
	}

	return c.doRequest(ctx, method, requestURL, bodyReader, extraHeaders)
}

// ServiceRequest sends a request to the Dify Service API using Bearer token
// authentication.
func (c *Client) ServiceRequest(ctx context.Context, method, path string, query url.Values, body interface{}) (json.RawMessage, error) {
	if c.serviceURL == "" {
		return nil, fmt.Errorf("service URL is not configured")
	}
	if c.apiKey == "" {
		return nil, fmt.Errorf("API key is not configured")
	}

	requestURL := c.serviceURL + "/" + strings.TrimPrefix(strings.TrimPrefix(path, "/v1"), "/")
	if len(query) > 0 {
		requestURL += "?" + query.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	extraHeaders := map[string]string{
		"Authorization": "Bearer " + c.apiKey,
	}

	return c.doRequest(ctx, method, requestURL, bodyReader, extraHeaders)
}

// doRequest executes an HTTP request with retry support for idempotent
// methods (GET, HEAD). Non-idempotent methods are executed directly without
// retry. Returns the parsed JSON response body.
func (c *Client) doRequest(ctx context.Context, method, requestURL string, body io.Reader, extraHeaders map[string]string) (json.RawMessage, error) {
	resp, err := optimize.DoWithHTTPRetry(
		ctx,
		method,
		c.maxRetries,
		c.retryBaseDelay,
		c.retryMaxDelay,
		func(attempt int) (*http.Response, error) {
			var bodyReader io.Reader
			if body != nil {
				// Reset the body reader for each attempt by reading it
				// into memory and creating a new reader.
				// For the first attempt, use the provided reader directly.
				if attempt > 1 {
					return nil, fmt.Errorf("cannot retry non-idempotent request with body")
				}
				bodyReader = body
			}

			req, err := http.NewRequestWithContext(ctx, method, requestURL, bodyReader)
			if err != nil {
				return nil, fmt.Errorf("failed to create request: %w", err)
			}

			req.Header.Set("Accept", "application/json")
			if bodyReader != nil {
				req.Header.Set("Content-Type", "application/json")
			}
			for k, v := range extraHeaders {
				req.Header.Set(k, v)
			}

			logrus.WithFields(logrus.Fields{
				"attempt": attempt,
				"method":  method,
				"url":     requestURL,
			}).Debug("Making Dify API request")

			return c.httpClient.Do(req)
		},
		func(event optimize.HTTPRetryEvent) {
			fields := logrus.Fields{
				"url":      requestURL,
				"method":   method,
				"attempt":  event.Attempt,
				"retry_in": event.Delay,
			}
			if event.Err != nil {
				logrus.WithFields(fields).WithError(event.Err).Warn("Retrying Dify API request after transient transport error")
				return
			}
			fields["status_code"] = event.StatusCode
			logrus.WithFields(fields).Warn("Retrying Dify API request after retryable status")
		},
	)
	if err != nil {
		return nil, fmt.Errorf("dify API request failed: %w", err)
	}

	return readResponse(resp)
}

// readResponse reads the response body and checks for HTTP errors.
func readResponse(resp *http.Response) (json.RawMessage, error) {
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read dify response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("dify API error (status %d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if len(body) == 0 {
		return json.RawMessage("{}"), nil
	}

	return json.RawMessage(body), nil
}

// buildCookieHeader builds a Cookie header value from a cookie map.
func buildCookieHeader(cookies map[string]string) string {
	if len(cookies) == 0 {
		return ""
	}
	parts := make([]string, 0, len(cookies))
	for name, value := range cookies {
		parts = append(parts, name+"="+value)
	}
	return strings.Join(parts, "; ")
}

// ParsePagination extracts pagination cursors from an HTTP Link header.
// Follows the same Link header parsing pattern as Sentry.
func ParsePagination(linkHeader string) *Pagination {
	linkHeader = strings.TrimSpace(linkHeader)
	if linkHeader == "" {
		return nil
	}

	pagination := &Pagination{}
	for _, part := range strings.Split(linkHeader, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		segments := strings.Split(part, ";")
		if len(segments) < 2 {
			continue
		}

		target := strings.TrimSpace(segments[0])
		target = strings.TrimPrefix(target, "<")
		target = strings.TrimSuffix(target, ">")

		cursor := &PaginationCursor{URL: target}
		for _, segment := range segments[1:] {
			segment = strings.TrimSpace(segment)
			key, value, ok := strings.Cut(segment, "=")
			if !ok {
				continue
			}
			key = strings.TrimSpace(key)
			value = strings.Trim(strings.TrimSpace(value), "\"")

			switch key {
			case "rel":
				cursor.Rel = value
			case "results":
				cursor.Results = strings.EqualFold(value, "true")
			case "cursor":
				cursor.Cursor = value
			}
		}

		switch cursor.Rel {
		case "next":
			pagination.Next = cursor
		case "previous", "prev":
			pagination.Prev = cursor
		}
	}

	if pagination.Next == nil && pagination.Prev == nil {
		return nil
	}
	return pagination
}

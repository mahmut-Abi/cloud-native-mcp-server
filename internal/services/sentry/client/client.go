package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
	"github.com/sirupsen/logrus"
)

const defaultRequestTimeout = 30 * time.Second

// PaginationCursor describes a Sentry cursor from the Link header.
type PaginationCursor struct {
	Cursor  string `json:"cursor"`
	Rel     string `json:"rel"`
	Results bool   `json:"results"`
	URL     string `json:"url,omitempty"`
}

// Pagination captures Sentry cursor pagination metadata.
type Pagination struct {
	Next *PaginationCursor `json:"next,omitempty"`
	Prev *PaginationCursor `json:"prev,omitempty"`
}

// ClientOptions holds configuration for creating a Sentry client.
type ClientOptions struct {
	URL            string
	AuthToken      string
	Timeout        time.Duration
	MaxRetries     int
	RetryBaseDelay time.Duration
	RetryMaxDelay  time.Duration
}

// Client provides read-only access to the Sentry REST API.
type Client struct {
	baseURL        string
	httpClient     *http.Client
	authToken      string
	maxRetries     int
	retryBaseDelay time.Duration
	retryMaxDelay  time.Duration
}

// NewClient creates a new Sentry client.
func NewClient(opts *ClientOptions) (*Client, error) {
	if opts == nil {
		return nil, fmt.Errorf("sentry client options are required")
	}
	if strings.TrimSpace(opts.URL) == "" {
		return nil, fmt.Errorf("sentry URL is required")
	}
	if strings.TrimSpace(opts.AuthToken) == "" {
		return nil, fmt.Errorf("sentry auth token is required")
	}

	parsedURL, err := url.Parse(strings.TrimSpace(opts.URL))
	if err != nil {
		return nil, fmt.Errorf("invalid sentry URL: %w", err)
	}

	cleanPath := strings.TrimSuffix(parsedURL.Path, "/")
	if !strings.HasSuffix(cleanPath, "/api/0") {
		parsedURL.Path = path.Join(cleanPath, "/api/0")
	}
	if !strings.HasSuffix(parsedURL.Path, "/") {
		parsedURL.Path += "/"
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
		baseURL:        parsedURL.String(),
		httpClient:     optimize.NewOptimizedHTTPClientWithTimeout(timeout),
		authToken:      strings.TrimSpace(opts.AuthToken),
		maxRetries:     maxRetries,
		retryBaseDelay: retryBaseDelay,
		retryMaxDelay:  retryMaxDelay,
	}, nil
}

// ListOrganizations returns organizations visible to the token.
func (c *Client) ListOrganizations(ctx context.Context, params url.Values) ([]map[string]interface{}, *Pagination, error) {
	return c.getJSONArray(ctx, "organizations/", params)
}

// ListProjects returns projects for an organization.
func (c *Client) ListProjects(ctx context.Context, organization string, params url.Values) ([]map[string]interface{}, *Pagination, error) {
	endpoint := fmt.Sprintf("organizations/%s/projects/", url.PathEscape(strings.TrimSpace(organization)))
	return c.getJSONArray(ctx, endpoint, params)
}

// GetProject returns details for a specific project.
func (c *Client) GetProject(ctx context.Context, organization, project string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("projects/%s/%s/", url.PathEscape(strings.TrimSpace(organization)), url.PathEscape(strings.TrimSpace(project)))
	return c.getJSONObject(ctx, endpoint, nil)
}

// ListIssues returns issues for an organization.
func (c *Client) ListIssues(ctx context.Context, organization string, params url.Values) ([]map[string]interface{}, *Pagination, error) {
	endpoint := fmt.Sprintf("organizations/%s/issues/", url.PathEscape(strings.TrimSpace(organization)))
	return c.getJSONArray(ctx, endpoint, params)
}

// GetIssue returns issue details by issue ID.
func (c *Client) GetIssue(ctx context.Context, issueID string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("issues/%s/", url.PathEscape(strings.TrimSpace(issueID)))
	return c.getJSONObject(ctx, endpoint, nil)
}

// ListIssueEvents returns events for an issue.
func (c *Client) ListIssueEvents(ctx context.Context, issueID string, params url.Values) ([]map[string]interface{}, *Pagination, error) {
	endpoint := fmt.Sprintf("issues/%s/events/", url.PathEscape(strings.TrimSpace(issueID)))
	return c.getJSONArray(ctx, endpoint, params)
}

// GetIssueEvent returns a specific event for an issue.
func (c *Client) GetIssueEvent(ctx context.Context, issueID, eventID string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("issues/%s/events/%s/", url.PathEscape(strings.TrimSpace(issueID)), url.PathEscape(strings.TrimSpace(eventID)))
	return c.getJSONObject(ctx, endpoint, nil)
}

func (c *Client) getJSONObject(ctx context.Context, endpoint string, params url.Values) (map[string]interface{}, error) {
	resp, err := c.makeRequest(ctx, endpoint, params)
	if err != nil {
		return nil, err
	}

	body, _, err := c.readResponse(resp)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return map[string]interface{}{}, nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sentry response: %w", err)
	}

	return result, nil
}

func (c *Client) getJSONArray(ctx context.Context, endpoint string, params url.Values) ([]map[string]interface{}, *Pagination, error) {
	resp, err := c.makeRequest(ctx, endpoint, params)
	if err != nil {
		return nil, nil, err
	}

	body, headers, err := c.readResponse(resp)
	if err != nil {
		return nil, nil, err
	}

	var result []map[string]interface{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal sentry response array: %w", err)
		}
	}

	return result, parsePagination(headers.Get("Link")), nil
}

func (c *Client) makeRequest(ctx context.Context, endpoint string, params url.Values) (*http.Response, error) {
	requestURL := c.baseURL + strings.TrimPrefix(endpoint, "/")
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

			req.Header.Set("Accept", "application/json")
			req.Header.Set("Authorization", "Bearer "+c.authToken)

			logrus.WithFields(logrus.Fields{
				"attempt": attempt,
				"url":     requestURL,
			}).Debug("Making Sentry API request")

			return c.httpClient.Do(req)
		},
		func(event optimize.HTTPRetryEvent) {
			fields := logrus.Fields{
				"url":      requestURL,
				"attempt":  event.Attempt,
				"retry_in": event.Delay,
			}
			if event.Err != nil {
				logrus.WithFields(fields).WithError(event.Err).Warn("Retrying Sentry API request after transient transport error")
				return
			}
			fields["status_code"] = event.StatusCode
			logrus.WithFields(fields).Warn("Retrying Sentry API request after retryable status")
		},
	)
}

func (c *Client) readResponse(resp *http.Response) ([]byte, http.Header, error) {
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read sentry response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, nil, fmt.Errorf("sentry API error (status %d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return body, resp.Header.Clone(), nil
}

func parsePagination(linkHeader string) *Pagination {
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

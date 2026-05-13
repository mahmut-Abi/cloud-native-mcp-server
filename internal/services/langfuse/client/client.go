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

// ClientOptions holds configuration for creating a Langfuse client.
type ClientOptions struct {
	URL            string
	PublicKey      string
	SecretKey      string
	Timeout        time.Duration
	MaxRetries     int
	RetryBaseDelay time.Duration
	RetryMaxDelay  time.Duration
}

// Client provides read-only access to the Langfuse Public API.
type Client struct {
	baseURL        string
	httpClient     *http.Client
	publicKey      string
	secretKey      string
	maxRetries     int
	retryBaseDelay time.Duration
	retryMaxDelay  time.Duration
}

// NewClient creates a new Langfuse client.
func NewClient(opts *ClientOptions) (*Client, error) {
	if opts == nil {
		return nil, fmt.Errorf("langfuse client options are required")
	}
	if strings.TrimSpace(opts.URL) == "" {
		return nil, fmt.Errorf("langfuse URL is required")
	}
	if strings.TrimSpace(opts.PublicKey) == "" {
		return nil, fmt.Errorf("langfuse public key is required")
	}
	if strings.TrimSpace(opts.SecretKey) == "" {
		return nil, fmt.Errorf("langfuse secret key is required")
	}

	parsedURL, err := url.Parse(strings.TrimSpace(opts.URL))
	if err != nil {
		return nil, fmt.Errorf("invalid langfuse URL: %w", err)
	}

	cleanPath := strings.TrimSuffix(parsedURL.Path, "/")
	if !strings.HasSuffix(cleanPath, "/api/public") {
		parsedURL.Path = path.Join(cleanPath, "/api/public")
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
		publicKey:      strings.TrimSpace(opts.PublicKey),
		secretKey:      strings.TrimSpace(opts.SecretKey),
		maxRetries:     maxRetries,
		retryBaseDelay: retryBaseDelay,
		retryMaxDelay:  retryMaxDelay,
	}, nil
}

// CheckHealth returns Langfuse API and database health information.
func (c *Client) CheckHealth(ctx context.Context) (map[string]interface{}, error) {
	return c.getJSON(ctx, "health", nil)
}

// ListTraces returns paginated traces.
func (c *Client) ListTraces(ctx context.Context, params url.Values) (map[string]interface{}, error) {
	return c.getJSON(ctx, "traces", params)
}

// GetTrace returns a specific trace by ID.
func (c *Client) GetTrace(ctx context.Context, traceID string, fields string) (map[string]interface{}, error) {
	params := url.Values{}
	if trimmed := strings.TrimSpace(fields); trimmed != "" {
		params.Set("fields", trimmed)
	}
	return c.getJSON(ctx, "traces/"+url.PathEscape(strings.TrimSpace(traceID)), params)
}

// ListSessions returns paginated sessions.
func (c *Client) ListSessions(ctx context.Context, params url.Values) (map[string]interface{}, error) {
	return c.getJSON(ctx, "sessions", params)
}

// GetSession returns a single session by ID.
func (c *Client) GetSession(ctx context.Context, sessionID string) (map[string]interface{}, error) {
	return c.getJSON(ctx, "sessions/"+url.PathEscape(strings.TrimSpace(sessionID)), nil)
}

// ListObservations returns paginated observations.
func (c *Client) ListObservations(ctx context.Context, params url.Values) (map[string]interface{}, error) {
	return c.getJSON(ctx, "observations", params)
}

// GetObservation returns a single observation by ID.
func (c *Client) GetObservation(ctx context.Context, observationID string) (map[string]interface{}, error) {
	return c.getJSON(ctx, "observations/"+url.PathEscape(strings.TrimSpace(observationID)), nil)
}

// ListPrompts returns paginated prompt metadata via the v2 prompts endpoint.
func (c *Client) ListPrompts(ctx context.Context, params url.Values) (map[string]interface{}, error) {
	return c.getJSON(ctx, "v2/prompts", params)
}

// ListAnnotationQueues returns paginated annotation queues.
func (c *Client) ListAnnotationQueues(ctx context.Context, params url.Values) (map[string]interface{}, error) {
	return c.getJSON(ctx, "annotation-queues", params)
}

// GetAnnotationQueue returns a single annotation queue by ID.
func (c *Client) GetAnnotationQueue(ctx context.Context, queueID string) (map[string]interface{}, error) {
	return c.getJSON(ctx, "annotation-queues/"+url.PathEscape(strings.TrimSpace(queueID)), nil)
}

// ListAnnotationQueueItems returns items for a specific annotation queue.
func (c *Client) ListAnnotationQueueItems(ctx context.Context, queueID string, params url.Values) (map[string]interface{}, error) {
	return c.getJSON(ctx, "annotation-queues/"+url.PathEscape(strings.TrimSpace(queueID))+"/items", params)
}

// ListDatasets returns paginated datasets.
func (c *Client) ListDatasets(ctx context.Context, params url.Values) (map[string]interface{}, error) {
	return c.getJSON(ctx, "v2/datasets", params)
}

// GetDataset returns a specific dataset by name.
func (c *Client) GetDataset(ctx context.Context, datasetName string) (map[string]interface{}, error) {
	return c.getJSON(ctx, "v2/datasets/"+url.PathEscape(strings.TrimSpace(datasetName)), nil)
}

// ListDatasetRuns returns paginated runs for a dataset.
func (c *Client) ListDatasetRuns(ctx context.Context, datasetName string, params url.Values) (map[string]interface{}, error) {
	return c.getJSON(ctx, "datasets/"+url.PathEscape(strings.TrimSpace(datasetName))+"/runs", params)
}

// GetDatasetRun returns a single dataset run by name.
func (c *Client) GetDatasetRun(ctx context.Context, datasetName, runName string) (map[string]interface{}, error) {
	return c.getJSON(ctx, "datasets/"+url.PathEscape(strings.TrimSpace(datasetName))+"/runs/"+url.PathEscape(strings.TrimSpace(runName)), nil)
}

// ListLLMConnections returns paginated LLM connections.
func (c *Client) ListLLMConnections(ctx context.Context, params url.Values) (map[string]interface{}, error) {
	return c.getJSON(ctx, "llm-connections", params)
}

// ListModels returns paginated model definitions.
func (c *Client) ListModels(ctx context.Context, params url.Values) (map[string]interface{}, error) {
	return c.getJSON(ctx, "models", params)
}

// GetModel returns a single model definition by ID.
func (c *Client) GetModel(ctx context.Context, modelID string) (map[string]interface{}, error) {
	return c.getJSON(ctx, "models/"+url.PathEscape(strings.TrimSpace(modelID)), nil)
}

// ListScoreConfigs returns paginated score configurations.
func (c *Client) ListScoreConfigs(ctx context.Context, params url.Values) (map[string]interface{}, error) {
	return c.getJSON(ctx, "score-configs", params)
}

// GetScoreConfig returns a single score configuration by ID.
func (c *Client) GetScoreConfig(ctx context.Context, configID string) (map[string]interface{}, error) {
	return c.getJSON(ctx, "score-configs/"+url.PathEscape(strings.TrimSpace(configID)), nil)
}

// GetPrompt returns a prompt by name, label, or version.
func (c *Client) GetPrompt(ctx context.Context, promptName string, params url.Values) (map[string]interface{}, error) {
	return c.getJSON(ctx, "v2/prompts/"+url.PathEscape(strings.TrimSpace(promptName)), params)
}

// ListScores returns paginated scores via the v2 scores endpoint.
func (c *Client) ListScores(ctx context.Context, params url.Values) (map[string]interface{}, error) {
	return c.getJSON(ctx, "v2/scores", params)
}

// GetScore returns a single score by ID.
func (c *Client) GetScore(ctx context.Context, scoreID string) (map[string]interface{}, error) {
	return c.getJSON(ctx, "v2/scores/"+url.PathEscape(strings.TrimSpace(scoreID)), nil)
}

// GetMetrics executes a metrics query against Langfuse.
func (c *Client) GetMetrics(ctx context.Context, queryJSON string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("query", queryJSON)
	return c.getJSON(ctx, "metrics", params)
}

func (c *Client) getJSON(ctx context.Context, endpoint string, params url.Values) (map[string]interface{}, error) {
	resp, err := c.makeRequest(ctx, endpoint, params)
	if err != nil {
		return nil, err
	}
	return c.readJSON(resp)
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
			req.SetBasicAuth(c.publicKey, c.secretKey)

			logrus.WithFields(logrus.Fields{
				"attempt": attempt,
				"url":     requestURL,
			}).Debug("Making Langfuse API request")

			return c.httpClient.Do(req)
		},
		func(event optimize.HTTPRetryEvent) {
			fields := logrus.Fields{
				"url":      requestURL,
				"attempt":  event.Attempt,
				"retry_in": event.Delay,
			}
			if event.Err != nil {
				logrus.WithFields(fields).WithError(event.Err).Warn("Retrying Langfuse API request after transient transport error")
				return
			}
			fields["status_code"] = event.StatusCode
			logrus.WithFields(fields).Warn("Retrying Langfuse API request after retryable status")
		},
	)
}

func (c *Client) readJSON(resp *http.Response) (map[string]interface{}, error) {
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Langfuse response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("langfuse API error (status %d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if len(body) == 0 {
		return map[string]interface{}{}, nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Langfuse response: %w", err)
	}

	return result, nil
}

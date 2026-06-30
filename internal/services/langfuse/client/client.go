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

// ClientOptions holds configuration for creating a Langfuse client.
type ClientOptions struct {
	URL            string
	Username       string
	Password       string
	ProjectID      string        // admin API key target project (self-hosted only)
	Timeout        time.Duration
	MaxRetries     int
	RetryBaseDelay time.Duration
	RetryMaxDelay  time.Duration
}

// Client provides access to the Langfuse Public API.
type Client struct {
	baseURL        string
	httpClient     *http.Client
	username       string
	password       string
	projectID      string
	maxRetries     int
	retryBaseDelay time.Duration
	retryMaxDelay  time.Duration

	adminKey    string             // original admin API key for project switching
	projectKeys map[string]struct {
		publicKey string
		secretKey string
	}
	projectKeysMu sync.Mutex
}

// NewClient creates a new Langfuse client.
func NewClient(opts *ClientOptions) (*Client, error) {
	if opts == nil {
		return nil, fmt.Errorf("langfuse client options are required")
	}
	if strings.TrimSpace(opts.URL) == "" {
		return nil, fmt.Errorf("langfuse URL is required")
	}

	username := strings.TrimSpace(opts.Username)
	password := strings.TrimSpace(opts.Password)
	if username == "" {
		return nil, fmt.Errorf("langfuse username is required")
	}
	if password == "" {
		return nil, fmt.Errorf("langfuse password is required")
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
		username:       username,
		password:       password,
		projectID:      strings.TrimSpace(opts.ProjectID),
		maxRetries:     maxRetries,
		retryBaseDelay: retryBaseDelay,
		retryMaxDelay:  retryMaxDelay,
		adminKey:       password,
		projectKeys:    make(map[string]struct{ publicKey, secretKey string }),
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

// GetProject returns the project associated with the configured project-scoped credentials.
func (c *Client) GetProject(ctx context.Context) (map[string]interface{}, error) {
	return c.getJSON(ctx, "projects", nil)
}

// ListOrganizationProjects returns all projects visible to the configured organization-scoped credentials.
func (c *Client) ListOrganizationProjects(ctx context.Context) (map[string]interface{}, error) {
	return c.getJSON(ctx, "organizations/projects", nil)
}

// CreateProject creates a Langfuse project.
func (c *Client) CreateProject(ctx context.Context, name string, metadata map[string]interface{}, retentionDays int) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"name":      strings.TrimSpace(name),
		"retention": retentionDays,
	}
	if metadata != nil {
		payload["metadata"] = metadata
	}
	return c.postJSON(ctx, "projects", nil, payload)
}

// UpdateProject updates a Langfuse project.
func (c *Client) UpdateProject(ctx context.Context, projectID, name string, metadata map[string]interface{}, retentionDays *int) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"name": strings.TrimSpace(name),
	}
	if metadata != nil {
		payload["metadata"] = metadata
	}
	if retentionDays != nil {
		payload["retention"] = *retentionDays
	}
	return c.putJSON(ctx, "projects/"+url.PathEscape(strings.TrimSpace(projectID)), nil, payload)
}

// DeleteProject deletes a Langfuse project asynchronously.
func (c *Client) DeleteProject(ctx context.Context, projectID string) (map[string]interface{}, error) {
	return c.deleteJSON(ctx, "projects/"+url.PathEscape(strings.TrimSpace(projectID)), nil)
}

// ListProjectMemberships returns memberships for a project.
func (c *Client) ListProjectMemberships(ctx context.Context, projectID string) (map[string]interface{}, error) {
	return c.getJSON(ctx, "projects/"+url.PathEscape(strings.TrimSpace(projectID))+"/memberships", nil)
}

// UpsertProjectMembership creates or updates a user's project membership.
func (c *Client) UpsertProjectMembership(ctx context.Context, projectID, userID, role string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"userId": strings.TrimSpace(userID),
		"role":   strings.TrimSpace(role),
	}
	return c.putJSON(ctx, "projects/"+url.PathEscape(strings.TrimSpace(projectID))+"/memberships", nil, payload)
}

// DeleteProjectMembership deletes a user's membership from a project.
func (c *Client) DeleteProjectMembership(ctx context.Context, projectID, userID string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"userId": strings.TrimSpace(userID),
	}
	return c.requestJSON(ctx, http.MethodDelete, "projects/"+url.PathEscape(strings.TrimSpace(projectID))+"/memberships", nil, payload)
}

// ListOrganizationAPIKeys returns organization API keys visible to the configured organization-scoped credentials.
func (c *Client) ListOrganizationAPIKeys(ctx context.Context) (map[string]interface{}, error) {
	return c.getJSON(ctx, "organizations/apiKeys", nil)
}

// ListProjectAPIKeys returns API keys for a project.
func (c *Client) ListProjectAPIKeys(ctx context.Context, projectID string) (map[string]interface{}, error) {
	return c.getJSON(ctx, "projects/"+url.PathEscape(strings.TrimSpace(projectID))+"/apiKeys", nil)
}

// CreateProjectAPIKey creates a new API key for a project.
func (c *Client) CreateProjectAPIKey(ctx context.Context, projectID, note, publicKey, secretKey string) (map[string]interface{}, error) {
	publicKey = strings.TrimSpace(publicKey)
	secretKey = strings.TrimSpace(secretKey)
	if (publicKey == "") != (secretKey == "") {
		return nil, fmt.Errorf("publicKey and secretKey must be provided together when predefining credentials")
	}

	payload := map[string]interface{}{}
	if note = strings.TrimSpace(note); note != "" {
		payload["note"] = note
	}
	if publicKey != "" {
		payload["publicKey"] = publicKey
		payload["secretKey"] = secretKey
	}

	return c.postJSON(ctx, "projects/"+url.PathEscape(strings.TrimSpace(projectID))+"/apiKeys", nil, payload)
}

// DeleteProjectAPIKey deletes an API key from a project.
func (c *Client) DeleteProjectAPIKey(ctx context.Context, projectID, apiKeyID string) (map[string]interface{}, error) {
	endpoint := "projects/" + url.PathEscape(strings.TrimSpace(projectID)) + "/apiKeys/" + url.PathEscape(strings.TrimSpace(apiKeyID))
	return c.deleteJSON(ctx, endpoint, nil)
}

func (c *Client) getJSON(ctx context.Context, endpoint string, params url.Values) (map[string]interface{}, error) {
	return c.requestJSON(ctx, http.MethodGet, endpoint, params, nil)
}

func (c *Client) postJSON(ctx context.Context, endpoint string, params url.Values, payload interface{}) (map[string]interface{}, error) {
	return c.requestJSON(ctx, http.MethodPost, endpoint, params, payload)
}

func (c *Client) putJSON(ctx context.Context, endpoint string, params url.Values, payload interface{}) (map[string]interface{}, error) {
	return c.requestJSON(ctx, http.MethodPut, endpoint, params, payload)
}

func (c *Client) deleteJSON(ctx context.Context, endpoint string, params url.Values) (map[string]interface{}, error) {
	return c.requestJSON(ctx, http.MethodDelete, endpoint, params, nil)
}

func (c *Client) requestJSON(ctx context.Context, method, endpoint string, params url.Values, payload interface{}) (map[string]interface{}, error) {
	resp, err := c.makeRequest(ctx, method, endpoint, params, payload)
	if err != nil {
		return nil, err
	}
	return c.readJSON(resp)
}

func (c *Client) makeRequest(ctx context.Context, method, endpoint string, params url.Values, payload interface{}) (*http.Response, error) {
	requestURL := c.baseURL + strings.TrimPrefix(endpoint, "/")
	if len(params) > 0 {
		requestURL += "?" + params.Encode()
	}

	var bodyBytes []byte
	if payload != nil {
		var err error
		bodyBytes, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal Langfuse request body: %w", err)
		}
	}

	return optimize.DoWithHTTPRetry(
		ctx,
		method,
		c.maxRetries,
		c.retryBaseDelay,
		c.retryMaxDelay,
		func(attempt int) (*http.Response, error) {
			var body io.Reader
			if bodyBytes != nil {
				body = bytes.NewReader(bodyBytes)
			}
			req, err := http.NewRequestWithContext(ctx, method, requestURL, body)
			if err != nil {
				return nil, fmt.Errorf("failed to create request: %w", err)
			}

			req.Header.Set("Accept", "application/json")
			if bodyBytes != nil {
				req.Header.Set("Content-Type", "application/json")
			}

			if c.projectID != "" {
				req.Header.Set("Authorization", "Bearer "+c.password)
				req.Header.Set("x-langfuse-admin-api-key", c.password)
				req.Header.Set("x-langfuse-project-id", c.projectID)
			} else {
				req.SetBasicAuth(c.username, c.password)
			}

			logrus.WithFields(logrus.Fields{
				"attempt": attempt,
				"method":  method,
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

// SwitchProject switches the client to use a project-level API key for the
// given project. Uses the stored admin key to fetch or create a project key.
func (c *Client) SwitchProject(ctx context.Context, targetProjectID string) error {
	c.projectKeysMu.Lock()
	defer c.projectKeysMu.Unlock()

	if pk, ok := c.projectKeys[targetProjectID]; ok {
		c.username = pk.publicKey
		c.password = pk.secretKey
		c.projectID = ""
		return nil
	}

	pk, err := fetchOrCreateProjectKey(ctx, c.baseURL, c.adminKey, targetProjectID)
	if err != nil {
		return err
	}

	c.projectKeys[targetProjectID] = struct {
		publicKey string
		secretKey string
	}{pk.PublicKey, pk.SecretKey}
	c.username = pk.PublicKey
	c.password = pk.SecretKey
	c.projectID = ""
	return nil
}

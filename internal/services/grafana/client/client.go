// Package client provides Grafana HTTP API client functionality.
// It offers operations for interacting with Grafana dashboards, data sources,
// and other Grafana resources through REST API calls.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/errors"
	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
	"github.com/sirupsen/logrus"
)

// ClientOptions holds configuration parameters for creating a Grafana client.
type ClientOptions struct {
	URL                  string               // Grafana server URL
	APIKey               string               // Grafana API key for authentication
	Username             string               // Username for basic authentication
	Password             string               // Password for basic authentication
	Timeout              time.Duration        // HTTP request timeout
	SkipVerify           bool                 // Skip TLS certificate verification
	CircuitBreakerConfig CircuitBreakerConfig // Circuit breaker configuration
}

// Client provides operations for interacting with Grafana API.
type Client struct {
	baseURL               string                 // Base URL for Grafana API
	httpClient            *http.Client           // HTTP client for API requests
	apiKey                string                 // API key for authentication
	username              string                 // Username for basic auth
	password              string                 // Password for basic auth
	headers               map[string]string      // Additional headers
	circuitBreakerManager *CircuitBreakerManager // Circuit breaker manager
}

// Dashboard represents a Grafana dashboard.
type Dashboard struct {
	ID        int                    `json:"id,omitempty"`
	UID       string                 `json:"uid,omitempty"`
	Title     string                 `json:"title"`
	Tags      []string               `json:"tags,omitempty"`
	FolderID  int                    `json:"folderId,omitempty"`
	FolderUID string                 `json:"folderUid,omitempty"`
	IsStarred bool                   `json:"isStarred,omitempty"`
	Dashboard map[string]interface{} `json:"dashboard,omitempty"`
	URL       string                 `json:"url,omitempty"`
	Version   int                    `json:"version,omitempty"`
	Meta      map[string]interface{} `json:"meta,omitempty"`
}

// DataSource represents a Grafana data source.
type DataSource struct {
	ID        int                    `json:"id,omitempty"`
	UID       string                 `json:"uid,omitempty"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	URL       string                 `json:"url"`
	Access    string                 `json:"access,omitempty"`
	Database  string                 `json:"database,omitempty"`
	User      string                 `json:"user,omitempty"`
	Password  string                 `json:"password,omitempty"`
	IsDefault bool                   `json:"isDefault,omitempty"`
	JSONData  map[string]interface{} `json:"jsonData,omitempty"`
}

// Folder represents a Grafana folder.
type Folder struct {
	ID       int    `json:"id,omitempty"`
	UID      string `json:"uid,omitempty"`
	Title    string `json:"title"`
	URL      string `json:"url,omitempty"`
	HasACL   bool   `json:"hasAcl,omitempty"`
	CanSave  bool   `json:"canSave,omitempty"`
	CanEdit  bool   `json:"canEdit,omitempty"`
	CanAdmin bool   `json:"canAdmin,omitempty"`
	Version  int    `json:"version,omitempty"`
}

// AlertRule represents a Grafana alert rule.
type AlertRule struct {
	ID              int                      `json:"id,omitempty"`
	UID             string                   `json:"uid,omitempty"`
	Title           string                   `json:"title"`
	Condition       string                   `json:"condition"`
	Data            []map[string]interface{} `json:"data"`
	IntervalSeconds int                      `json:"intervalSeconds"`
	NoDataState     string                   `json:"noDataState"`
	ExecErrState    string                   `json:"execErrState"`
	For             string                   `json:"for"`
	Annotations     map[string]string        `json:"annotations,omitempty"`
	Labels          map[string]string        `json:"labels,omitempty"`
	FolderUID       string                   `json:"folderUID,omitempty"`
	RuleGroup       string                   `json:"ruleGroup,omitempty"`
}

// NewClient creates a new Grafana client with the specified options.
func NewClient(opts *ClientOptions) (*Client, error) {
	if opts.URL == "" {
		return nil, fmt.Errorf("grafana URL is required")
	}

	// Parse and validate URL
	baseURL, err := url.Parse(opts.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid grafana URL: %w", err)
	}

	// Ensure URL has proper path
	if !strings.HasSuffix(baseURL.Path, "/") {
		baseURL.Path += "/"
	}
	baseURL.Path += "api/"

	// Create HTTP client with timeout
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	httpClient := optimize.NewOptimizedHTTPClientWithTimeout(timeout)

	// Prepare headers
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"

	// Initialize circuit breaker manager
	cbConfig := opts.CircuitBreakerConfig
	if !cbConfig.Enabled {
		cbConfig = CircuitBreakerConfig{Enabled: false}
	}

	client := &Client{
		baseURL:               baseURL.String(),
		httpClient:            httpClient,
		apiKey:                opts.APIKey,
		username:              opts.Username,
		password:              opts.Password,
		headers:               headers,
		circuitBreakerManager: NewCircuitBreakerManager(cbConfig),
	}

	return client, nil
}

// makeRequest performs an HTTP request to the Grafana API.
func (c *Client) makeRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	url := c.baseURL + strings.TrimPrefix(endpoint, "/")
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// Set authentication
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	} else if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	logrus.WithFields(logrus.Fields{
		"method":   method,
		"url":      url,
		"has_body": body != nil,
	}).Debug("Making Grafana API request")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "GRAFANA_CONNECTION_FAILED", "failed to connect to Grafana").
			WithHTTPStatus(503)
	}

	return resp, nil
}

// handleResponse processes the HTTP response and returns the body.
func (c *Client) handleResponse(resp *http.Response) ([]byte, error) {
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "GRAFANA_INVALID_RESPONSE", "invalid Grafana API response").
			WithHTTPStatus(502)
	}

	if resp.StatusCode >= 400 {
		// Handle specific status codes
		switch resp.StatusCode {
		case 401:
			return nil, errors.New("GRAFANA_UNAUTHORIZED", "unauthorized access to Grafana").
				WithHTTPStatus(401)
		case 403:
			return nil, errors.New("GRAFANA_FORBIDDEN", "forbidden access to Grafana resource").
				WithHTTPStatus(403)
		case 404:
			return nil, errors.NotFoundError("resource")
		case 429:
			return nil, errors.New("GRAFANA_RATE_LIMITED", "Grafana API rate limit exceeded").
				WithHTTPStatus(429)
		case 500:
			return nil, errors.Wrap(fmt.Errorf("%s", string(body)), "GRAFANA_SERVER_ERROR", "Grafana server error").
				WithHTTPStatus(500)
		default:
			return nil, errors.New("GRAFANA_API_ERROR", fmt.Sprintf("Grafana API error (status %d): %s", resp.StatusCode, string(body))).
				WithHTTPStatus(resp.StatusCode).
				WithContext("status_code", resp.StatusCode)
		}
	}

	return body, nil
}

// GetDashboards retrieves all dashboards from Grafana.
func (c *Client) GetDashboards(ctx context.Context) ([]Dashboard, error) {
	logrus.Debug("Getting Grafana dashboards")

	var dashboards []Dashboard
	err := c.wrapWithCircuitBreaker(ctx, OpGetDashboards, func(ctx context.Context) error {
		resp, err := c.makeRequest(ctx, "GET", "search?type=dash-db", nil)
		if err != nil {
			return err
		}

		body, err := c.handleResponse(resp)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(body, &dashboards); err != nil {
			return fmt.Errorf("failed to unmarshal dashboards: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	logrus.WithField("count", len(dashboards)).Debug("Retrieved dashboards")
	return dashboards, nil
}

// GetDashboard retrieves a specific dashboard by UID.
func (c *Client) GetDashboard(ctx context.Context, uid string) (*Dashboard, error) {
	logrus.WithField("uid", uid).Debug("Getting Grafana dashboard")

	var dashboard *Dashboard
	err := c.wrapWithCircuitBreaker(ctx, OpGetDashboard, func(ctx context.Context) error {
		resp, err := c.makeRequest(ctx, "GET", "dashboards/uid/"+uid, nil)
		if err != nil {
			return err
		}

		body, err := c.handleResponse(resp)
		if err != nil {
			return err
		}

		var result struct {
			Dashboard Dashboard              `json:"dashboard"`
			Meta      map[string]interface{} `json:"meta"`
		}
		if err := json.Unmarshal(body, &result); err != nil {
			return fmt.Errorf("failed to unmarshal dashboard: %w", err)
		}

		result.Dashboard.Meta = result.Meta
		dashboard = &result.Dashboard
		return nil
	})

	if err != nil {
		return nil, err
	}

	logrus.Debug("Retrieved dashboard")
	return dashboard, nil
}

// GetDataSources retrieves all data sources from Grafana.
func (c *Client) GetDataSources(ctx context.Context) ([]DataSource, error) {
	logrus.Debug("Getting Grafana data sources")

	resp, err := c.makeRequest(ctx, "GET", "datasources", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var dataSources []DataSource
	if err := json.Unmarshal(body, &dataSources); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data sources: %w", err)
	}

	logrus.WithField("count", len(dataSources)).Debug("Retrieved data sources")
	return dataSources, nil
}

// GetFolders retrieves all folders from Grafana.
func (c *Client) GetFolders(ctx context.Context) ([]Folder, error) {
	logrus.Debug("Getting Grafana folders")

	resp, err := c.makeRequest(ctx, "GET", "folders", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var folders []Folder
	if err := json.Unmarshal(body, &folders); err != nil {
		return nil, fmt.Errorf("failed to unmarshal folders: %w", err)
	}

	logrus.WithField("count", len(folders)).Debug("Retrieved folders")
	return folders, nil
}

// GetAlertRules retrieves alert rules from Grafana.
func (c *Client) GetAlertRules(ctx context.Context) ([]AlertRule, error) {
	logrus.Debug("Getting Grafana alert rules")

	resp, err := c.makeRequest(ctx, "GET", "ruler/grafana/api/v1/rules", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	// Alert rules API returns a complex nested structure
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal alert rules: %w", err)
	}

	// Extract alert rules from the nested structure
	var rules []AlertRule
	// This is a simplified extraction - the actual structure is more complex
	logrus.WithField("raw_result", result).Debug("Alert rules raw response")

	logrus.WithField("count", len(rules)).Debug("Retrieved alert rules")
	return rules, nil
}

// TestConnection tests the connection to Grafana API.
func (c *Client) TestConnection(ctx context.Context) error {
	logrus.Debug("Testing Grafana connection")

	resp, err := c.makeRequest(ctx, "GET", "health", nil)
	if err != nil {
		return fmt.Errorf("failed to connect to grafana: %w", err)
	}

	_, err = c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("grafana health check failed: %w", err)
	}

	logrus.Debug("Grafana connection test successful")
	return nil
}

// User represents a Grafana user.
type User struct {
	ID            int    `json:"id,omitempty"`
	Username      string `json:"login"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	IsAdmin       bool   `json:"isAdmin"`
	IsDisabled    bool   `json:"isDisabled"`
	LastSeenAtUTC int64  `json:"lastSeenAtUtc,omitempty"`
	LastSeenAt    string `json:"lastSeenAt,omitempty"`
}

// Organization represents a Grafana organization.
type Organization struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// DatasourceHealth represents datasource health status.
type DatasourceHealth struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	Message string `json:"message"`
	UID     string `json:"uid"`
}

// GetFolder retrieves a specific folder by UID.
func (c *Client) GetFolder(ctx context.Context, uid string) (*Folder, error) {
	logrus.WithField("uid", uid).Debug("Getting Grafana folder")

	resp, err := c.makeRequest(ctx, "GET", "folders/"+uid, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var folder Folder
	if err := json.Unmarshal(body, &folder); err != nil {
		return nil, fmt.Errorf("failed to unmarshal folder: %w", err)
	}

	logrus.Debug("Retrieved folder")
	return &folder, nil
}

// GetDataSource retrieves a specific data source by UID.
func (c *Client) GetDataSource(ctx context.Context, uid string) (*DataSource, error) {
	logrus.WithField("uid", uid).Debug("Getting Grafana data source")

	resp, err := c.makeRequest(ctx, "GET", "datasources/uid/"+uid, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var dataSource DataSource
	if err := json.Unmarshal(body, &dataSource); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data source: %w", err)
	}

	logrus.Debug("Retrieved data source")
	return &dataSource, nil
}

// CheckDatasourceHealth tests the health of a specific datasource.
func (c *Client) CheckDatasourceHealth(ctx context.Context, uid string) (*DatasourceHealth, error) {
	logrus.WithField("uid", uid).Debug("Checking datasource health")

	resp, err := c.makeRequest(ctx, "GET", "datasources/uid/"+uid+"/health", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var health DatasourceHealth
	if err := json.Unmarshal(body, &health); err != nil {
		return nil, fmt.Errorf("failed to unmarshal datasource health: %w", err)
	}

	logrus.Debug("Retrieved datasource health")
	return &health, nil
}

// GetOrganization retrieves organization information.
func (c *Client) GetOrganization(ctx context.Context) (*Organization, error) {
	logrus.Debug("Getting Grafana organization")

	resp, err := c.makeRequest(ctx, "GET", "org", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var org Organization
	if err := json.Unmarshal(body, &org); err != nil {
		return nil, fmt.Errorf("failed to unmarshal organization: %w", err)
	}

	logrus.Debug("Retrieved organization")
	return &org, nil
}

// GetCurrentUser retrieves the current authenticated user.
func (c *Client) GetCurrentUser(ctx context.Context) (*User, error) {
	logrus.Debug("Getting current Grafana user")

	resp, err := c.makeRequest(ctx, "GET", "user", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	logrus.Debug("Retrieved current user")
	return &user, nil
}

// GetUsers retrieves all users in the organization.
func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	logrus.Debug("Getting Grafana users")

	resp, err := c.makeRequest(ctx, "GET", "users", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var users []User
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal users: %w", err)
	}

	logrus.WithField("count", len(users)).Debug("Retrieved users")
	return users, nil
}

// ============ Admin Tools ============

// Team represents a Grafana team.
type Team struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	MemberCount int    `json:"memberCount"`
}

// Role represents a Grafana role.
type Role struct {
	UID         string       `json:"uid"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions"`
	BuiltIn     bool         `json:"builtIn"`
}

// Permission represents a Grafana permission.
type Permission struct {
	Action string `json:"action"`
	Scope  string `json:"scope"`
}

// RoleAssignment represents a role assignment in Grafana.
type RoleAssignment struct {
	RoleUID string `json:"roleUID"`
	UserID  int    `json:"userId,omitempty"`
	TeamID  int    `json:"teamId,omitempty"`
}

// GetTeams retrieves all teams in the organization.
func (c *Client) GetTeams(ctx context.Context) ([]Team, error) {
	logrus.Debug("Getting Grafana teams")

	resp, err := c.makeRequest(ctx, "GET", "teams", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result struct {
		Teams []Team `json:"teams"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal teams: %w", err)
	}

	logrus.WithField("count", len(result.Teams)).Debug("Retrieved teams")
	return result.Teams, nil
}

// GetRoles retrieves all available roles.
func (c *Client) GetRoles(ctx context.Context) ([]Role, error) {
	logrus.Debug("Getting Grafana roles")

	resp, err := c.makeRequest(ctx, "GET", "access-control/roles", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var roles []Role
	if err := json.Unmarshal(body, &roles); err != nil {
		return nil, fmt.Errorf("failed to unmarshal roles: %w", err)
	}

	logrus.WithField("count", len(roles)).Debug("Retrieved roles")
	return roles, nil
}

// GetRoleDetails retrieves details of a specific role.
func (c *Client) GetRoleDetails(ctx context.Context, roleUID string) (*Role, error) {
	logrus.WithField("roleUID", roleUID).Debug("Getting role details")

	resp, err := c.makeRequest(ctx, "GET", "access-control/roles/"+roleUID, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var role Role
	if err := json.Unmarshal(body, &role); err != nil {
		return nil, fmt.Errorf("failed to unmarshal role: %w", err)
	}

	logrus.Debug("Retrieved role details")
	return &role, nil
}

// GetRoleAssignments retrieves assignments for a specific role.
func (c *Client) GetRoleAssignments(ctx context.Context, roleUID string) ([]RoleAssignment, error) {
	logrus.WithField("roleUID", roleUID).Debug("Getting role assignments")

	resp, err := c.makeRequest(ctx, "GET", "access-control/roles/"+roleUID+"/assignments", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var assignments []RoleAssignment
	if err := json.Unmarshal(body, &assignments); err != nil {
		return nil, fmt.Errorf("failed to unmarshal role assignments: %w", err)
	}

	logrus.WithField("count", len(assignments)).Debug("Retrieved role assignments")
	return assignments, nil
}

// GetUserRoles retrieves roles for a specific user.
func (c *Client) GetUserRoles(ctx context.Context, userID int) ([]Role, error) {
	logrus.WithField("userID", userID).Debug("Getting user roles")

	resp, err := c.makeRequest(ctx, "GET", fmt.Sprintf("users/%d/roles", userID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var roles []Role
	if err := json.Unmarshal(body, &roles); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user roles: %w", err)
	}

	logrus.WithField("count", len(roles)).Debug("Retrieved user roles")
	return roles, nil
}

// GetTeamRoles retrieves roles for a specific team.
func (c *Client) GetTeamRoles(ctx context.Context, teamID int) ([]Role, error) {
	logrus.WithField("teamID", teamID).Debug("Getting team roles")

	resp, err := c.makeRequest(ctx, "GET", fmt.Sprintf("teams/%d/roles", teamID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var roles []Role
	if err := json.Unmarshal(body, &roles); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team roles: %w", err)
	}

	logrus.WithField("count", len(roles)).Debug("Retrieved team roles")
	return roles, nil
}

// GetResourcePermissions retrieves permissions for a specific resource.
func (c *Client) GetResourcePermissions(ctx context.Context, resourceType string, resourceUID string) ([]Permission, error) {
	logrus.WithFields(logrus.Fields{
		"resourceType": resourceType,
		"resourceUID":  resourceUID,
	}).Debug("Getting resource permissions")

	resp, err := c.makeRequest(ctx, "GET", fmt.Sprintf("access-control/resource/%s/%s/permissions", resourceType, resourceUID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var permissions []Permission
	if err := json.Unmarshal(body, &permissions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal permissions: %w", err)
	}

	logrus.WithField("count", len(permissions)).Debug("Retrieved resource permissions")
	return permissions, nil
}

// ResourceDescription describes a Grafana resource type.
type ResourceDescription struct {
	Resource     string `json:"resource"`
	Descriptions string `json:"descriptions"`
	CanCreate    bool   `json:"canCreate"`
	CanRead      bool   `json:"canRead"`
	CanUpdate    bool   `json:"canUpdate"`
	CanDelete    bool   `json:"canDelete"`
}

// GetResourceDescription retrieves description for a resource type.
func (c *Client) GetResourceDescription(ctx context.Context, resourceType string) (*ResourceDescription, error) {
	logrus.WithField("resourceType", resourceType).Debug("Getting resource description")

	resp, err := c.makeRequest(ctx, "GET", "access-control/resource/"+resourceType+"/description", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var desc ResourceDescription
	if err := json.Unmarshal(body, &desc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal resource description: %w", err)
	}

	logrus.Debug("Retrieved resource description")
	return &desc, nil
}

// ============ Dashboard Update Tools ============

// DashboardUpdateRequest represents a request to create or update a dashboard.
type DashboardUpdateRequest struct {
	Dashboard map[string]interface{} `json:"dashboard"`
	FolderUID string                 `json:"folderUid"`
	Overwrite bool                   `json:"overwrite"`
	Message   string                 `json:"message,omitempty"`
}

// UpdateDashboard creates or updates a dashboard.
func (c *Client) UpdateDashboard(ctx context.Context, req DashboardUpdateRequest) (*Dashboard, error) {
	logrus.Debug("Updating Grafana dashboard")

	resp, err := c.makeRequest(ctx, "POST", "dashboards", req)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result struct {
		Dashboard Dashboard `json:"dashboard"`
		Slug      string    `json:"slug"`
		Version   int       `json:"version"`
		ID        int       `json:"id"`
		UID       string    `json:"uid"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dashboard update response: %w", err)
	}

	logrus.WithField("uid", result.UID).Debug("Dashboard updated")
	return &result.Dashboard, nil
}

// PanelQuery represents a panel query with datasource info.
type PanelQuery struct {
	PanelID        int    `json:"panelId"`
	PanelTitle     string `json:"panelTitle"`
	Query          string `json:"query"`
	DatasourceUID  string `json:"datasourceUid"`
	DatasourceType string `json:"datasourceType"`
	DatasourceName string `json:"datasourceName"`
}

// DashboardPanelInfo represents panel information with queries.
type DashboardPanelInfo struct {
	PanelID        int          `json:"id"`
	PanelTitle     string       `json:"title"`
	Type           string       `json:"type"`
	Queries        []PanelQuery `json:"queries"`
	DatasourceUID  string       `json:"datasourceUid,omitempty"`
	DatasourceType string       `json:"datasourceType,omitempty"`
}

// GetDashboardPanelQueries retrieves panel queries and datasource info from a dashboard.
func (c *Client) GetDashboardPanelQueries(ctx context.Context, dashboardUID string) ([]DashboardPanelInfo, error) {
	logrus.WithField("dashboardUID", dashboardUID).Debug("Getting dashboard panel queries")

	// First get the dashboard
	dashboard, err := c.GetDashboard(ctx, dashboardUID)
	if err != nil {
		return nil, err
	}

	var panels []DashboardPanelInfo

	// Extract panels from dashboard JSON
	if dashboard.Dashboard != nil {
		if panelList, ok := dashboard.Dashboard["panels"].([]interface{}); ok {
			for _, panel := range panelList {
				if panelMap, ok := panel.(map[string]interface{}); ok {
					panelInfo := DashboardPanelInfo{
						PanelID:    int(panelMap["id"].(float64)),
						PanelTitle: panelMap["title"].(string),
						Type:       panelMap["type"].(string),
						Queries:    []PanelQuery{},
					}

					// Extract targets (queries)
					if targets, ok := panelMap["targets"].([]interface{}); ok {
						for _, target := range targets {
							if targetMap, ok := target.(map[string]interface{}); ok {
								query := PanelQuery{
									PanelID:    panelInfo.PanelID,
									PanelTitle: panelInfo.PanelTitle,
								}

								// Extract datasource info
								if ds, ok := panelMap["datasource"].(map[string]interface{}); ok {
									if uid, ok := ds["uid"].(string); ok {
										panelInfo.DatasourceUID = uid
									}
									if typeStr, ok := ds["type"].(string); ok {
										panelInfo.DatasourceType = typeStr
									}
								}

								// Extract raw query
								if expr, ok := targetMap["expr"].(string); ok {
									query.Query = expr
								} else if queryStr, ok := targetMap["query"].(string); ok {
									query.Query = queryStr
								}

								panelInfo.Queries = append(panelInfo.Queries, query)
							}
						}
					}

					panels = append(panels, panelInfo)
				}
			}
		}
	}

	logrus.WithField("count", len(panels)).Debug("Retrieved dashboard panel queries")
	return panels, nil
}

// GetDashboardProperty extracts specific parts of a dashboard using JSONPath-like expressions.
func (c *Client) GetDashboardProperty(ctx context.Context, dashboardUID string, propertyPath string) (interface{}, error) {
	logrus.WithFields(logrus.Fields{
		"dashboardUID": dashboardUID,
		"propertyPath": propertyPath,
	}).Debug("Getting dashboard property")

	dashboard, err := c.GetDashboard(ctx, dashboardUID)
	if err != nil {
		return nil, err
	}

	// Simple property extraction (simplified JSONPath)
	// Supports simple paths like "title", "panels[0].title", "tags"
	var value interface{} = dashboard.Dashboard

	// Handle dashboard-level fields
	if propertyPath == "title" {
		return dashboard.Title, nil
	}
	if propertyPath == "uid" {
		return dashboard.UID, nil
	}
	if propertyPath == "tags" {
		return dashboard.Tags, nil
	}
	if propertyPath == "version" {
		return dashboard.Version, nil
	}

	// Navigate nested dashboard data
	parts := strings.Split(propertyPath, ".")
	for _, part := range parts {
		if value == nil {
			return nil, fmt.Errorf("property path not found: %s", propertyPath)
		}

		if mapVal, ok := value.(map[string]interface{}); ok {
			value = mapVal[part]
		} else {
			return nil, fmt.Errorf("cannot navigate to %s in property path", part)
		}
	}

	logrus.Debug("Retrieved dashboard property")
	return value, nil
}

// ============ Alerting Tools ============

// AlertRuleDetail represents detailed alert rule information.
type AlertRuleDetail struct {
	ID              int                      `json:"id"`
	UID             string                   `json:"uid"`
	Title           string                   `json:"title"`
	Condition       string                   `json:"condition"`
	Data            []map[string]interface{} `json:"data"`
	IntervalSeconds int                      `json:"intervalSeconds"`
	NoDataState     string                   `json:"noDataState"`
	ExecErrState    string                   `json:"execErrState"`
	For             string                   `json:"for"`
	Annotations     map[string]string        `json:"annotations"`
	Labels          map[string]string        `json:"labels"`
	FolderUID       string                   `json:"folderUID"`
	RuleGroup       string                   `json:"ruleGroup"`
	OrgID           int                      `json:"orgId"`
}

// GetAlertRuleByUID retrieves a specific alert rule by UID.
func (c *Client) GetAlertRuleByUID(ctx context.Context, ruleUID string) (*AlertRuleDetail, error) {
	logrus.WithField("ruleUID", ruleUID).Debug("Getting alert rule by UID")

	resp, err := c.makeRequest(ctx, "GET", "alerting/rules/"+ruleUID, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var rule AlertRuleDetail
	if err := json.Unmarshal(body, &rule); err != nil {
		return nil, fmt.Errorf("failed to unmarshal alert rule: %w", err)
	}

	logrus.Debug("Retrieved alert rule")
	return &rule, nil
}

// CreateAlertRuleRequest represents a request to create an alert rule.
type CreateAlertRuleRequest struct {
	Title           string                   `json:"title"`
	Condition       string                   `json:"condition"`
	Data            []map[string]interface{} `json:"data"`
	IntervalSeconds int                      `json:"intervalSeconds"`
	NoDataState     string                   `json:"noDataState,omitempty"`
	ExecErrState    string                   `json:"execErrState,omitempty"`
	For             string                   `json:"for,omitempty"`
	Annotations     map[string]string        `json:"annotations,omitempty"`
	Labels          map[string]string        `json:"labels,omitempty"`
	FolderUID       string                   `json:"folderUID"`
	RuleGroup       string                   `json:"ruleGroup"`
}

// CreateAlertRule creates a new alert rule.
func (c *Client) CreateAlertRule(ctx context.Context, req CreateAlertRuleRequest) (*AlertRuleDetail, error) {
	logrus.WithField("title", req.Title).Debug("Creating alert rule")

	resp, err := c.makeRequest(ctx, "POST", "alerting/rules", req)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var rule AlertRuleDetail
	if err := json.Unmarshal(body, &rule); err != nil {
		return nil, fmt.Errorf("failed to unmarshal alert rule: %w", err)
	}

	logrus.WithField("uid", rule.UID).Debug("Alert rule created")
	return &rule, nil
}

// UpdateAlertRuleRequest represents a request to update an alert rule.
type UpdateAlertRuleRequest struct {
	Title           string                   `json:"title"`
	Condition       string                   `json:"condition"`
	Data            []map[string]interface{} `json:"data"`
	IntervalSeconds int                      `json:"intervalSeconds"`
	NoDataState     string                   `json:"noDataState,omitempty"`
	ExecErrState    string                   `json:"execErrState,omitempty"`
	For             string                   `json:"for,omitempty"`
	Annotations     map[string]string        `json:"annotations,omitempty"`
	Labels          map[string]string        `json:"labels,omitempty"`
	FolderUID       string                   `json:"folderUID"`
	RuleGroup       string                   `json:"ruleGroup"`
}

// UpdateAlertRule updates an existing alert rule.
func (c *Client) UpdateAlertRule(ctx context.Context, ruleUID string, req UpdateAlertRuleRequest) (*AlertRuleDetail, error) {
	logrus.WithField("ruleUID", ruleUID).Debug("Updating alert rule")

	resp, err := c.makeRequest(ctx, "PUT", "alerting/rules/"+ruleUID, req)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var rule AlertRuleDetail
	if err := json.Unmarshal(body, &rule); err != nil {
		return nil, fmt.Errorf("failed to unmarshal alert rule: %w", err)
	}

	logrus.Debug("Alert rule updated")
	return &rule, nil
}

// DeleteAlertRule deletes an alert rule by UID.
func (c *Client) DeleteAlertRule(ctx context.Context, ruleUID string) error {
	logrus.WithField("ruleUID", ruleUID).Debug("Deleting alert rule")

	resp, err := c.makeRequest(ctx, "DELETE", "alerting/rules/"+ruleUID, nil)
	if err != nil {
		return err
	}

	_, err = c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to delete alert rule: %w", err)
	}

	logrus.Debug("Alert rule deleted")
	return nil
}

// ContactPoint represents a Grafana contact point.
type ContactPoint struct {
	UID        string                 `json:"uid"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Settings   map[string]interface{} `json:"settings"`
	Provenance string                 `json:"provenance,omitempty"`
}

// GetContactPoints retrieves all contact points.
func (c *Client) GetContactPoints(ctx context.Context) ([]ContactPoint, error) {
	logrus.Debug("Getting contact points")

	resp, err := c.makeRequest(ctx, "GET", "alerting/notifications/channels", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var contactPoints []ContactPoint
	if err := json.Unmarshal(body, &contactPoints); err != nil {
		return nil, fmt.Errorf("failed to unmarshal contact points: %w", err)
	}

	logrus.WithField("count", len(contactPoints)).Debug("Retrieved contact points")
	return contactPoints, nil
}

// ============ Annotation Tools ============

// Annotation represents a Grafana annotation.
type Annotation struct {
	ID           int64    `json:"id"`
	Annotation   int      `json:"annotation"`
	DashboardUID string   `json:"dashboardUid,omitempty"`
	PanelID      int      `json:"panelId,omitempty"`
	UserID       int      `json:"userId"`
	UserLogin    string   `json:"userLogin"`
	NewState     string   `json:"newState"`
	PrevState    string   `json:"prevState"`
	Time         int64    `json:"time"`
	TimeEnd      int64    `json:"timeEnd,omitempty"`
	Text         string   `json:"text"`
	Tags         []string `json:"tags"`
	Type         string   `json:"type"`
	RegionID     string   `json:"regionId,omitempty"`
	Created      int64    `json:"created"`
	Updated      int64    `json:"updated"`
	Region       string   `json:"region,omitempty"`
}

// GetAnnotations retrieves annotations with optional filters.
func (c *Client) GetAnnotations(ctx context.Context, params map[string]string) ([]Annotation, error) {
	logrus.Debug("Getting annotations")

	// Build query params
	queryParams := ""
	if len(params) > 0 {
		parts := []string{}
		for k, v := range params {
			parts = append(parts, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
		}
		queryParams = "?" + strings.Join(parts, "&")
	}

	resp, err := c.makeRequest(ctx, "GET", "annotations"+queryParams, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var annotations []Annotation
	if err := json.Unmarshal(body, &annotations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal annotations: %w", err)
	}

	logrus.WithField("count", len(annotations)).Debug("Retrieved annotations")
	return annotations, nil
}

// CreateAnnotationRequest represents a request to create an annotation.
type CreateAnnotationRequest struct {
	DashboardUID string   `json:"dashboardUid,omitempty"`
	PanelID      int      `json:"panelId,omitempty"`
	Time         int64    `json:"time"`
	TimeEnd      int64    `json:"timeEnd,omitempty"`
	Text         string   `json:"text"`
	Tags         []string `json:"tags"`
}

// CreateAnnotation creates a new annotation.
func (c *Client) CreateAnnotation(ctx context.Context, req CreateAnnotationRequest) (*Annotation, error) {
	logrus.Debug("Creating annotation")

	resp, err := c.makeRequest(ctx, "POST", "annotations", req)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var annotation Annotation
	if err := json.Unmarshal(body, &annotation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal annotation: %w", err)
	}

	logrus.WithField("id", annotation.ID).Debug("Annotation created")
	return &annotation, nil
}

// UpdateAnnotationRequest represents a request to update an annotation.
type UpdateAnnotationRequest struct {
	Time    int64    `json:"time,omitempty"`
	TimeEnd int64    `json:"timeEnd,omitempty"`
	Text    string   `json:"text,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}

// UpdateAnnotation updates an existing annotation.
func (c *Client) UpdateAnnotation(ctx context.Context, annotationID int64, req UpdateAnnotationRequest) (*Annotation, error) {
	logrus.WithField("annotationID", annotationID).Debug("Updating annotation")

	resp, err := c.makeRequest(ctx, "PUT", fmt.Sprintf("annotations/%d", annotationID), req)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var annotation Annotation
	if err := json.Unmarshal(body, &annotation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal annotation: %w", err)
	}

	logrus.Debug("Annotation updated")
	return &annotation, nil
}

// PatchAnnotationRequest represents a request to patch an annotation.
type PatchAnnotationRequest struct {
	Time    int64    `json:"time,omitempty"`
	TimeEnd int64    `json:"timeEnd,omitempty"`
	Text    string   `json:"text,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}

// PatchAnnotation partially updates an annotation.
func (c *Client) PatchAnnotation(ctx context.Context, annotationID int64, req PatchAnnotationRequest) (*Annotation, error) {
	logrus.WithField("annotationID", annotationID).Debug("Patching annotation")

	resp, err := c.makeRequest(ctx, "PATCH", fmt.Sprintf("annotations/%d", annotationID), req)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var annotation Annotation
	if err := json.Unmarshal(body, &annotation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal annotation: %w", err)
	}

	logrus.Debug("Annotation patched")
	return &annotation, nil
}

// AnnotationTag represents an annotation tag.
type AnnotationTag struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

// GetAnnotationTags retrieves available annotation tags.
func (c *Client) GetAnnotationTags(ctx context.Context, tag string) ([]AnnotationTag, error) {
	logrus.Debug("Getting annotation tags")

	queryParams := ""
	if tag != "" {
		queryParams = "?tag=" + url.QueryEscape(tag)
	}

	resp, err := c.makeRequest(ctx, "GET", "annotations/tags"+queryParams, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var result struct {
		Tags []AnnotationTag `json:"tags"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal annotation tags: %w", err)
	}

	logrus.WithField("count", len(result.Tags)).Debug("Retrieved annotation tags")
	return result.Tags, nil
}

// DeleteAnnotation deletes an annotation by ID.
func (c *Client) DeleteAnnotation(ctx context.Context, annotationID int64) error {
	logrus.WithField("annotationID", annotationID).Debug("Deleting annotation")

	resp, err := c.makeRequest(ctx, "DELETE", fmt.Sprintf("annotations/%d", annotationID), nil)
	if err != nil {
		return err
	}

	_, err = c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to delete annotation: %w", err)
	}

	logrus.Debug("Annotation deleted")
	return nil
}

// ============ Navigation/Deeplink Tools ============

// Deeplink represents a generated deeplink URL.
type Deeplink struct {
	URL   string `json:"url"`
	Label string `json:"label"`
	Path  string `json:"path"`
}

// GenerateDeeplink generates accurate deeplink URLs for Grafana resources.
func (c *Client) GenerateDeeplink(ctx context.Context, resourceType string, resourceUID string, params map[string]string) (*Deeplink, error) {
	logrus.WithFields(logrus.Fields{
		"resourceType": resourceType,
		"resourceUID":  resourceUID,
	}).Debug("Generating deeplink")

	baseURL := strings.TrimSuffix(strings.TrimSuffix(c.baseURL, "api/"), "/")

	var path string
	switch resourceType {
	case "dashboard":
		path = fmt.Sprintf("/d/%s", resourceUID)
	case "panel":
		if panelID, ok := params["panelId"]; ok {
			path = fmt.Sprintf("/d/%s?viewPanel=%s", resourceUID, panelID)
		} else {
			path = fmt.Sprintf("/d/%s", resourceUID)
		}
	case "explore":
		if dsUID, ok := params["datasource"]; ok {
			path = fmt.Sprintf("/explore?left={\"datasource\":\"%s\"}", dsUID)
		} else {
			path = "/explore"
		}
	default:
		return nil, fmt.Errorf("unknown resource type: %s", resourceType)
	}

	// Add time range if specified
	if from, ok := params["from"]; ok {
		if to, ok := params["to"]; ok {
			if strings.Contains(path, "?") {
				path += fmt.Sprintf("&from=%s&to=%s", from, to)
			} else {
				path += fmt.Sprintf("?from=%s&to=%s", from, to)
			}
		}
	}

	deeplink := &Deeplink{
		URL:   baseURL + path,
		Label: fmt.Sprintf("Grafana %s: %s", resourceType, resourceUID),
		Path:  path,
	}

	logrus.WithField("url", deeplink.URL).Debug("Generated deeplink")
	return deeplink, nil
}

// ============ Datasource by Name ============

// GetDataSourceByName retrieves a data source by name.
func (c *Client) GetDataSourceByName(ctx context.Context, name string) (*DataSource, error) {
	logrus.WithField("name", name).Debug("Getting data source by name")

	resp, err := c.makeRequest(ctx, "GET", "datasources/name/"+url.QueryEscape(name), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var dataSource DataSource
	if err := json.Unmarshal(body, &dataSource); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data source: %w", err)
	}

	logrus.Debug("Retrieved data source by name")
	return &dataSource, nil
}

// ============ Panel Image Rendering ============

// PanelImageResponse represents the response from rendering a panel.
type PanelImageResponse struct {
	ImageData   []byte `json:"-"`
	ContentType string `json:"contentType"`
	Size        int    `json:"size"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
}

// RenderDashboardPanel renders a dashboard panel to an image.
func (c *Client) RenderDashboardPanel(ctx context.Context, dashboardUID string, panelID int, params map[string]string) (*PanelImageResponse, error) {
	logrus.WithFields(logrus.Fields{
		"dashboardUID": dashboardUID,
		"panelID":      panelID,
	}).Debug("Rendering dashboard panel")

	queryParams := url.Values{}
	queryParams.Add("panelId", fmt.Sprintf("%d", panelID))
	queryParams.Add("width", "800")
	queryParams.Add("height", "400")

	for key, value := range params {
		queryParams.Add(key, value)
	}

	queryParams.Add("key", "render")

	path := fmt.Sprintf("render/d-solo/%s?%s", dashboardUID, queryParams.Encode())

	resp, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to render panel: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		return nil, fmt.Errorf("render failed with status %d: %s", resp.StatusCode, string(body))
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read panel image: %w", err)
	}

	logrus.WithField("size", len(imageData)).Debug("Rendered panel successfully")

	return &PanelImageResponse{
		ImageData:   imageData,
		ContentType: resp.Header.Get("Content-Type"),
		Size:        len(imageData),
	}, nil
}

// ============ Graphite Annotation ============

// GraphiteAnnotation represents a Graphite annotation.
type GraphiteAnnotation struct {
	ID        int64  `json:"id"`
	What      string `json:"what"`
	Data      string `json:"data"`
	Timestamp int64  `json:"timestamp"`
	Tags      string `json:"tags,omitempty"`
}

// CreateGraphiteAnnotation creates a Graphite annotation.
func (c *Client) CreateGraphiteAnnotation(ctx context.Context, what string, data string, timestamp int64, tags string) (*GraphiteAnnotation, error) {
	logrus.WithField("what", what).Debug("Creating Graphite annotation")

	payload := map[string]interface{}{
		"what": what,
		"data": data,
	}

	if timestamp > 0 {
		payload["when"] = timestamp
	}

	if tags != "" {
		payload["tags"] = tags
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal annotation: %w", err)
	}

	resp, err := c.makeRequest(ctx, "POST", "api/annotations/graphite", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	respBody, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var annotation GraphiteAnnotation
	if err := json.Unmarshal(respBody, &annotation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal annotation: %w", err)
	}

	logrus.WithField("id", annotation.ID).Debug("Created Graphite annotation")
	return &annotation, nil
}

// ============ Datasource Management Tools ============

// CreateDatasourceRequest represents a request to create a datasource.
type CreateDatasourceRequest struct {
	Name           string                 `json:"name"`
	Type           string                 `json:"type"`
	URL            string                 `json:"url"`
	Access         string                 `json:"access,omitempty"`
	Database       string                 `json:"database,omitempty"`
	User           string                 `json:"user,omitempty"`
	Password       string                 `json:"password,omitempty"`
	JSONData       map[string]interface{} `json:"jsonData,omitempty"`
	SecureJSONData map[string]interface{} `json:"secureJsonData,omitempty"`
	IsDefault      bool                   `json:"isDefault,omitempty"`
}

// CreateDatasource creates a new datasource in Grafana.
func (c *Client) CreateDatasource(ctx context.Context, req CreateDatasourceRequest) (*DataSource, error) {
	logrus.WithField("name", req.Name).Debug("Creating datasource")

	resp, err := c.makeRequest(ctx, "POST", "datasources", req)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var dataSource DataSource
	if err := json.Unmarshal(body, &dataSource); err != nil {
		return nil, fmt.Errorf("failed to unmarshal datasource: %w", err)
	}

	logrus.WithField("uid", dataSource.UID).Debug("Datasource created successfully")
	return &dataSource, nil
}

// UpdateDatasourceRequest represents a request to update a datasource.
type UpdateDatasourceRequest struct {
	UID            string                 `json:"uid"`
	Name           string                 `json:"name"`
	Type           string                 `json:"type"`
	URL            string                 `json:"url"`
	Access         string                 `json:"access,omitempty"`
	Database       string                 `json:"database,omitempty"`
	User           string                 `json:"user,omitempty"`
	Password       string                 `json:"password,omitempty"`
	JSONData       map[string]interface{} `json:"jsonData,omitempty"`
	SecureJSONData map[string]interface{} `json:"secureJsonData,omitempty"`
	IsDefault      bool                   `json:"isDefault,omitempty"`
}

// UpdateDatasource updates an existing datasource in Grafana.
func (c *Client) UpdateDatasource(ctx context.Context, req UpdateDatasourceRequest) (*DataSource, error) {
	logrus.WithField("uid", req.UID).Debug("Updating datasource")

	resp, err := c.makeRequest(ctx, "PUT", "datasources/"+req.UID, req)
	if err != nil {
		return nil, err
	}

	body, err := c.handleResponse(resp)
	if err != nil {
		return nil, err
	}

	var dataSource DataSource
	if err := json.Unmarshal(body, &dataSource); err != nil {
		return nil, fmt.Errorf("failed to unmarshal datasource: %w", err)
	}

	logrus.WithField("uid", dataSource.UID).Debug("Datasource updated successfully")
	return &dataSource, nil
}

// DeleteDatasource deletes a datasource by UID.
func (c *Client) DeleteDatasource(ctx context.Context, uid string) error {
	logrus.WithField("uid", uid).Debug("Deleting datasource")

	resp, err := c.makeRequest(ctx, "DELETE", "datasources/uid/"+uid, nil)
	if err != nil {
		return err
	}

	_, err = c.handleResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to delete datasource: %w", err)
	}

	logrus.WithField("uid", uid).Debug("Datasource deleted successfully")
	return nil
}

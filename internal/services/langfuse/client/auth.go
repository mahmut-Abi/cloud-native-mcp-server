package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const csrfEndpoint = "/api/auth/csrf"
const signinEndpoint = "/api/auth/callback/credentials"
const sessionEndpoint = "/api/auth/session"

// Session represents the Langfuse session response.
type Session struct {
	User struct {
		ID                  string         `json:"id"`
		Email               string         `json:"email"`
		Name                string         `json:"name"`
		Organizations       []Organization `json:"organizations"`
		CanCreateOrganizations bool         `json:"canCreateOrganizations"`
	} `json:"user"`
}

// Organization represents a Langfuse organization.
type Organization struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Role     string    `json:"role"`
	Projects []Project `json:"projects"`
}

// Project represents a Langfuse project.
type Project struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
}

// APIKey represents a created API key.
type APIKey struct {
	ID               string `json:"id"`
	CreatedAt        string `json:"createdAt"`
	Note             string `json:"note"`
	PublicKey        string `json:"publicKey"`
	SecretKey        string `json:"secretKey"`
	DisplaySecretKey string `json:"displaySecretKey"`
}

// AuthError is returned when console auth fails.
type AuthError struct {
	Step    string
	Message string
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("langfuse console auth: %s: %s", e.Step, e.Message)
}

// ConsoleAuthenticator handles NextAuth-based console credential authentication.
type ConsoleAuthenticator struct {
	baseURL    string
	httpClient *http.Client
	jar        *cookiejar.Jar
	Projects   []Project
}

// NewConsoleAuthenticator creates a new authenticator for console credentials.
func NewConsoleAuthenticator(baseURL string) (*ConsoleAuthenticator, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	trimmed := strings.TrimRight(baseURL, "/")
	return &ConsoleAuthenticator{
		baseURL: trimmed,
		httpClient: &http.Client{
			Jar: jar,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 1 && req.URL.Host != via[0].URL.Host {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
		jar: jar,
	}, nil
}

func (a *ConsoleAuthenticator) base(path string) string {
	return a.baseURL + path
}

// Login authenticates with email/password via NextAuth credentials flow.
func (a *ConsoleAuthenticator) Login(email, password string) (*Session, error) {
	csrfToken, err := a.getCSRF()
	if err != nil {
		return nil, &AuthError{Step: "csrf", Message: err.Error()}
	}

	if err := a.postCredentials(csrfToken, email, password); err != nil {
		return nil, &AuthError{Step: "signin", Message: err.Error()}
	}

	return a.getSession()
}

func (a *ConsoleAuthenticator) getCSRF() (string, error) {
	resp, err := a.httpClient.Get(a.base(csrfEndpoint))
	if err != nil {
		return "", fmt.Errorf("csrf request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading csrf response: %w", err)
	}

	var csrfResp struct {
		CsrfToken string `json:"csrfToken"`
	}
	if err := json.Unmarshal(body, &csrfResp); err != nil {
		return "", fmt.Errorf("parsing csrf response: %w", err)
	}
	if csrfResp.CsrfToken == "" {
		return "", fmt.Errorf("empty csrf token")
	}

	return csrfResp.CsrfToken, nil
}

func (a *ConsoleAuthenticator) postCredentials(csrfToken, email, password string) error {
	form := url.Values{
		"csrfToken": {csrfToken},
		"email":     {email},
		"password":  {password},
	}

	resp, err := a.httpClient.PostForm(a.base(signinEndpoint), form)
	if err != nil {
		return fmt.Errorf("signin request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("signin returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (a *ConsoleAuthenticator) getSession() (*Session, error) {
	resp, err := a.httpClient.Get(a.base(sessionEndpoint))
	if err != nil {
		return nil, fmt.Errorf("session request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading session: %w", err)
	}

	var session Session
	if err := json.Unmarshal(body, &session); err != nil {
		return nil, fmt.Errorf("parsing session: %w", err)
	}

	if session.User.Email == "" {
		return nil, fmt.Errorf("not authenticated (empty user)")
	}

	return &session, nil
}

// GetProjects extracts all projects from the session.
func (a *ConsoleAuthenticator) GetProjects() []Project {
	return a.Projects
}

// CreateAPIKey creates a new project-level API key via TRPC mutation.
func (a *ConsoleAuthenticator) CreateAPIKey(projectID string) (*APIKey, error) {
	return a.trpcCreate("projectApiKeys.create", map[string]interface{}{
		"projectId": projectID,
		"note":      "mcp-server-auto",
	})
}

// CreateOrgAPIKey creates a new organization-level API key via TRPC mutation.
func (a *ConsoleAuthenticator) CreateOrgAPIKey(orgID string) (*APIKey, error) {
	return a.trpcCreate("organizationApiKeys.create", map[string]interface{}{
		"orgId": orgID,
		"note":  "mcp-server-auto-org",
	})
}

func (a *ConsoleAuthenticator) trpcCreate(procedure string, params map[string]interface{}) (*APIKey, error) {
	batchInput := map[string]interface{}{
		"0": map[string]interface{}{
			procedure: params,
		},
	}
	inputJSON, _ := json.Marshal(batchInput)

	reqURL := a.base("/api/trpc?batch=1")

	req, err := http.NewRequest("POST", reqURL, strings.NewReader(string(inputJSON)))
	if err != nil {
		return nil, fmt.Errorf("creating TRPC request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("TRPC request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading TRPC response: %w", err)
	}

	if resp.StatusCode >= 400 {
		contentType := resp.Header.Get("Content-Type")
		if strings.Contains(contentType, "text/html") || strings.HasPrefix(strings.TrimSpace(string(body)), "<!") {
			return nil, fmt.Errorf("TRPC endpoint /api/trpc returned %d (HTML) — TRPC not available on this Langfuse instance. Use pk-lf-* API keys directly or enable TRPC", resp.StatusCode)
		}
		return nil, fmt.Errorf("TRPC returned status %d: %s", resp.StatusCode, string(body))
	}

	logger.Printf("TRPC %s status=%d body=%s", procedure, resp.StatusCode, string(body))

	// Non-batch mutation response: {"result":{"data":{"json":{...}}}}
	var singleResult struct {
		Result struct {
			Data struct {
				JSON APIKey `json:"json"`
			} `json:"data"`
		} `json:"result"`
		Error struct {
			JSON struct {
				Message string `json:"message"`
				Code    int    `json:"code"`
			} `json:"json"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &singleResult); err != nil {
		// Try batch response format as fallback
		var batchResult []struct {
			Result struct {
				Data struct {
					JSON APIKey `json:"json"`
				} `json:"data"`
			} `json:"result"`
			Error struct {
				JSON struct {
					Message string `json:"message"`
					Code    int    `json:"code"`
				} `json:"json"`
			} `json:"error"`
		}
		if err2 := json.Unmarshal(body, &batchResult); err2 != nil {
			return nil, fmt.Errorf("parsing TRPC response: %w (body: %s)", err, string(body))
		}
		if len(batchResult) == 0 {
			return nil, fmt.Errorf("empty TRPC batch response")
		}
		if batchResult[0].Error.JSON.Message != "" {
			return nil, &AuthError{
				Step:    "trpc_create",
				Message: batchResult[0].Error.JSON.Message,
			}
		}
		key := &batchResult[0].Result.Data.JSON
		if key.PublicKey == "" || key.SecretKey == "" {
			return nil, fmt.Errorf("empty key data in TRPC response: %s", string(body))
		}
		return key, nil
	}

	if singleResult.Error.JSON.Message != "" {
		return nil, &AuthError{
			Step:    "trpc_create",
			Message: singleResult.Error.JSON.Message,
		}
	}

	key := &singleResult.Result.Data.JSON
	if key.PublicKey == "" || key.SecretKey == "" {
		return nil, fmt.Errorf("empty key data in TRPC response: %s", string(body))
	}

	return key, nil
}

// TryConsoleAuth attempts to authenticate with console credentials and create an API key.
// Tries organization-level key first, falls back to project-level key.
// Returns the API key (pk:sk) or an error.
func TryConsoleAuth(baseURL, username, password string) (*APIKey, string, error) {
	auth, err := NewConsoleAuthenticator(baseURL)
	if err != nil {
		return nil, "", fmt.Errorf("console auth init: %w", err)
	}

	session, err := auth.Login(username, password)
	if err != nil {
		return nil, "", err
	}

	// Collect projects from session
	for _, org := range session.User.Organizations {
		auth.Projects = append(auth.Projects, org.Projects...)
	}

	if len(session.User.Organizations) == 0 {
		return nil, "", fmt.Errorf("no organizations found for user %s", username)
	}
	if len(auth.Projects) == 0 {
		return nil, "", fmt.Errorf("no projects found for user %s", username)
	}

	// Try org-level API key first (best scope for global access)
	orgID := session.User.Organizations[0].ID
	orgName := session.User.Organizations[0].Name
	key, err := auth.CreateOrgAPIKey(orgID)
	if err == nil {
		return key, orgName + " (org-level)", nil
	}
	orgErr := err.Error()

	// Fall back to project-level key
	key, err = auth.CreateAPIKey(auth.Projects[0].ID)
	if err == nil {
		return key, auth.Projects[0].Name + " (project-level)", nil
	}

	return nil, "", fmt.Errorf("creating API key failed: org(%s), project: %s", orgErr, err.Error())
}

// TryConsoleAuthViaREST logs in via NextAuth with email/password, extracts
// projects from the session, and returns the first project's info.
// Does NOT create API keys (community edition restriction).
func TryConsoleAuthViaREST(baseURL, email, password string) (projectID, projectName string, _ error) {
	auth, err := NewConsoleAuthenticator(baseURL)
	if err != nil {
		return "", "", fmt.Errorf("console auth init: %w", err)
	}

	session, err := auth.Login(email, password)
	if err != nil {
		return "", "", fmt.Errorf("console login: %w", err)
	}

	// Collect all projects from all organizations
	var projects []struct{ id, name string }
	for _, org := range session.User.Organizations {
		for _, p := range org.Projects {
			projects = append(projects, struct{ id, name string }{p.ID, p.Name})
		}
	}

	if len(projects) == 0 {
		return "", "", fmt.Errorf("no projects found for user %s", email)
	}

	logger.Printf("Langfuse console auth: logged in as %s, found %d projects across %d orgs",
		email, len(projects), len(session.User.Organizations))

	return projects[0].id, projects[0].name, nil
}

// IsConsoleCredential returns true if the username is NOT a pk-lf-* API key.
func IsConsoleCredential(username string) bool {
	return !strings.HasPrefix(strings.TrimSpace(username), "pk-lf-")
}

// TryAdminKeyAuth uses an admin/org-level key to auto-discover and create a
// project-level API key via the Langfuse REST API (no TRPC required).
// Returns the project API key or an error.
func TryAdminKeyAuth(baseURL, username, password string) (*APIKey, string, error) {
	apiBase := strings.TrimRight(baseURL, "/") + "/api/public"
	client := &http.Client{Timeout: 30 * time.Second}

	adminHeaders := func(req *http.Request, projectID string) {
		req.Header.Set("Accept", "application/json")
		req.Header.Set("x-langfuse-admin-api-key", password)
		if projectID != "" {
			req.Header.Set("x-langfuse-project-id", projectID)
		}
	}

	// Step 1: list projects
	req, _ := http.NewRequest("GET", apiBase+"/projects", nil)
	adminHeaders(req, "")
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("admin key: listing projects: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("admin key: listing projects: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var projects struct {
		Data []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, "", fmt.Errorf("admin key: parsing projects: %w", err)
	}
	if len(projects.Data) == 0 {
		return nil, "", fmt.Errorf("admin key: no projects found")
	}

	projectID := projects.Data[0].ID
	projectName := projects.Data[0].Name

	// Step 2: get or create project API key
	pk, err := fetchOrCreateProjectKey(context.Background(), apiBase+"/", password, projectID)
	if err != nil {
		return nil, "", fmt.Errorf("admin key: %w", err)
	}
	logger.Printf("Langfuse admin key: project '%s' ready (pk=%s)", projectName, pk.PublicKey[:20]+"...")
	return pk, projectName, nil
}

func basicAuth(username, password string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
}

// fetchOrCreateProjectKey gets or creates a project-level API key using the admin key.
func fetchOrCreateProjectKey(ctx context.Context, apiBase, adminKey, projectID string) (*APIKey, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	hdr := func(req *http.Request) {
		req.Header.Set("Accept", "application/json")
		req.Header.Set("x-langfuse-admin-api-key", adminKey)
		req.Header.Set("x-langfuse-project-id", projectID)
	}

	// Try to get existing project keys
	req, _ := http.NewRequest("GET", apiBase+"projects/"+projectID+"/apiKeys", nil)
	hdr(req)
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode < 400 {
		defer resp.Body.Close()
		var result struct{ Data []APIKey `json:"data"` }
		if json.NewDecoder(resp.Body).Decode(&result) == nil && len(result.Data) > 0 {
			k := result.Data[0]
			if k.PublicKey != "" && k.SecretKey != "" {
				return &k, nil
			}
		}
	}
	if resp != nil {
		resp.Body.Close()
	}

	// Create a new project API key
	createReq, _ := http.NewRequest("POST", apiBase+"projects/"+projectID+"/apiKeys", nil)
	hdr(createReq)
	resp2, err := client.Do(createReq)
	if err != nil {
		return nil, fmt.Errorf("creating project API key: %w", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode >= 400 {
		body, _ := io.ReadAll(resp2.Body)
		return nil, fmt.Errorf("status %d: %s", resp2.StatusCode, strings.TrimSpace(string(body)))
	}

	var created APIKey
	if err := json.NewDecoder(resp2.Body).Decode(&created); err != nil {
		return nil, fmt.Errorf("decoding create response: %w", err)
	}
	if created.PublicKey == "" || created.SecretKey == "" {
		return nil, fmt.Errorf("empty key data in create response")
	}
	return &created, nil
}

func doJSON[T any](client *http.Client, method, url, authHeader string) (T, error) {
	var result T
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return result, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", authHeader)

	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return result, fmt.Errorf("API error (status %d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("decoding response: %w", err)
	}
	return result, nil
}

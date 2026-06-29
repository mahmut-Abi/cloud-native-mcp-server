package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
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
		baseURL:    trimmed,
		httpClient: &http.Client{Jar: jar},
		jar:        jar,
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
	defer resp.Body.Close()

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
	defer resp.Body.Close()

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
	defer resp.Body.Close()

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

// CreateAPIKey creates a new project-level API key via TRPC batch mutation.
func (a *ConsoleAuthenticator) CreateAPIKey(projectID string) (*APIKey, error) {
	input := map[string]interface{}{
		"0": map[string]interface{}{
			"projectId": projectID,
			"note":      "mcp-server-auto",
		},
	}

	inputJSON, _ := json.Marshal(input)
	encoded := url.QueryEscape(string(inputJSON))

	reqURL := a.base("/api/trpc/projectApiKeys.create") + "?batch=1&input=" + encoded

	req, err := http.NewRequest("POST", reqURL, strings.NewReader("{}"))
	if err != nil {
		return nil, fmt.Errorf("creating TRPC request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("TRPC request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading TRPC response: %w", err)
	}

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

	if err := json.Unmarshal(body, &batchResult); err != nil {
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

// TryConsoleAuth attempts to authenticate with console credentials and create an API key.
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

	if len(auth.Projects) == 0 {
		return nil, "", fmt.Errorf("no projects found for user %s", username)
	}

	// Try to create API key for the first project
	key, err := auth.CreateAPIKey(auth.Projects[0].ID)
	if err != nil {
		return nil, "", fmt.Errorf("creating API key: %w", err)
	}

	return key, auth.Projects[0].Name, nil
}

// IsConsoleCredential returns true if the username is NOT a pk-lf-* API key.
func IsConsoleCredential(username string) bool {
	return !strings.HasPrefix(strings.TrimSpace(username), "pk-lf-")
}

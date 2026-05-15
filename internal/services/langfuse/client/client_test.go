package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ClientOptions
		wantErr bool
	}{
		{
			name: "valid root URL",
			opts: &ClientOptions{
				URL:      "https://langfuse.example.com",
				Username: "pk-test",
				Password: "sk-test",
			},
		},
		{
			name: "valid api public URL",
			opts: &ClientOptions{
				URL:      "https://langfuse.example.com/api/public",
				Username: "pk-test",
				Password: "sk-test",
			},
		},
		{
			name: "legacy public and secret key aliases",
			opts: &ClientOptions{
				URL:       "https://langfuse.example.com",
				PublicKey: "pk-test",
				SecretKey: "sk-test",
			},
		},
		{
			name: "missing URL",
			opts: &ClientOptions{
				Username: "pk-test",
				Password: "sk-test",
			},
			wantErr: true,
		},
		{
			name: "missing username",
			opts: &ClientOptions{
				URL:      "https://langfuse.example.com",
				Password: "sk-test",
			},
			wantErr: true,
		},
		{
			name: "missing password",
			opts: &ClientOptions{
				URL:      "https://langfuse.example.com",
				Username: "pk-test",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && client == nil {
				t.Fatal("NewClient() returned nil client")
			}
		})
	}
}

func TestCheckHealthUsesBasicAuthAndPublicAPIPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			t.Fatal("expected basic auth")
		}
		if user != "pk-test" || pass != "sk-test" {
			t.Fatalf("unexpected basic auth credentials: %q / %q", user, pass)
		}
		if r.URL.Path != "/api/public/health" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		URL:      server.URL,
		Username: "pk-test",
		Password: "sk-test",
		Timeout:  5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	result, err := client.CheckHealth(context.Background())
	if err != nil {
		t.Fatalf("CheckHealth() error = %v", err)
	}
	if result["status"] != "ok" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestGetPromptEncodesNestedPromptNames(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/public/v2/prompts/folder/subfolder/prompt-a" {
			t.Fatalf("unexpected decoded path: %s", r.URL.Path)
		}
		if r.URL.RawPath != "/api/public/v2/prompts/folder%2Fsubfolder%2Fprompt-a" {
			t.Fatalf("unexpected raw path: %s", r.URL.RawPath)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"folder/subfolder/prompt-a"}`))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		URL:      server.URL,
		Username: "pk-test",
		Password: "sk-test",
		Timeout:  5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	result, err := client.GetPrompt(context.Background(), "folder/subfolder/prompt-a", nil)
	if err != nil {
		t.Fatalf("GetPrompt() error = %v", err)
	}
	if result["name"] != "folder/subfolder/prompt-a" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestGetDatasetEncodesDatasetName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/public/v2/datasets/folder/dataset-a" {
			t.Fatalf("unexpected decoded path: %s", r.URL.Path)
		}
		if r.URL.RawPath != "/api/public/v2/datasets/folder%2Fdataset-a" {
			t.Fatalf("unexpected raw path: %s", r.URL.RawPath)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"folder/dataset-a"}`))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		URL:      server.URL,
		Username: "pk-test",
		Password: "sk-test",
		Timeout:  5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	result, err := client.GetDataset(context.Background(), "folder/dataset-a")
	if err != nil {
		t.Fatalf("GetDataset() error = %v", err)
	}
	if result["name"] != "folder/dataset-a" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestListAnnotationQueueItemsIncludesStatusAndPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/public/annotation-queues/queue-1/items" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("status"); got != "PENDING" {
			t.Fatalf("unexpected status query: %s", got)
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Fatalf("unexpected page query: %s", got)
		}
		if got := r.URL.Query().Get("limit"); got != "10" {
			t.Fatalf("unexpected limit query: %s", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[],"meta":{"page":2,"limit":10}}`))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		URL:      server.URL,
		Username: "pk-test",
		Password: "sk-test",
		Timeout:  5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	params := url.Values{}
	params.Set("status", "PENDING")
	params.Set("page", "2")
	params.Set("limit", "10")
	result, err := client.ListAnnotationQueueItems(context.Background(), "queue-1", params)
	if err != nil {
		t.Fatalf("ListAnnotationQueueItems() error = %v", err)
	}
	if _, ok := result["meta"]; !ok {
		t.Fatalf("expected meta in result, got %#v", result)
	}
}

func TestProjectAPIKeyManagement(t *testing.T) {
	var seenCreate bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			t.Fatal("expected basic auth")
		}
		if user != "org-pk" || pass != "org-sk" {
			t.Fatalf("unexpected basic auth credentials: %q / %q", user, pass)
		}

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/public/organizations/apiKeys":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"apiKeys":[]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/public/projects/project-1/apiKeys":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"apiKeys":[{"id":"key-1","publicKey":"pk-lf-1","displaySecretKey":"sk-lf-..."}]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/public/projects/project-1/apiKeys":
			seenCreate = true
			defer func() { _ = r.Body.Close() }()
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode request body: %v", err)
			}
			if body["note"] != "created by test" {
				t.Fatalf("unexpected note: %#v", body)
			}
			if body["publicKey"] != "pk-lf-custom" || body["secretKey"] != "sk-lf-custom" {
				t.Fatalf("unexpected predefined keys: %#v", body)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"key-2","publicKey":"pk-lf-custom","secretKey":"sk-lf-custom","displaySecretKey":"sk-lf-..."}`))
		case r.Method == http.MethodDelete && r.URL.Path == "/api/public/projects/project-1/apiKeys/key-2":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"success":true}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		URL:      server.URL,
		Username: "org-pk",
		Password: "org-sk",
		Timeout:  5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if _, err := client.ListOrganizationAPIKeys(context.Background()); err != nil {
		t.Fatalf("ListOrganizationAPIKeys() error = %v", err)
	}
	if _, err := client.ListProjectAPIKeys(context.Background(), "project-1"); err != nil {
		t.Fatalf("ListProjectAPIKeys() error = %v", err)
	}
	created, err := client.CreateProjectAPIKey(context.Background(), "project-1", "created by test", "pk-lf-custom", "sk-lf-custom")
	if err != nil {
		t.Fatalf("CreateProjectAPIKey() error = %v", err)
	}
	if created["secretKey"] != "sk-lf-custom" {
		t.Fatalf("expected secret key in creation response, got %#v", created)
	}
	if _, err := client.DeleteProjectAPIKey(context.Background(), "project-1", "key-2"); err != nil {
		t.Fatalf("DeleteProjectAPIKey() error = %v", err)
	}
	if !seenCreate {
		t.Fatal("expected create request")
	}
}

func TestCreateProjectAPIKeyRequiresPredefinedPair(t *testing.T) {
	client, err := NewClient(&ClientOptions{
		URL:      "https://langfuse.example.com",
		Username: "org-pk",
		Password: "org-sk",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if _, err := client.CreateProjectAPIKey(context.Background(), "project-1", "", "pk-lf-custom", ""); err == nil {
		t.Fatal("expected error for partial predefined credentials")
	}
}

func TestProjectManagement(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			t.Fatal("expected basic auth")
		}
		if user != "org-pk" || pass != "org-sk" {
			t.Fatalf("unexpected basic auth credentials: %q / %q", user, pass)
		}

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/public/projects":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":[{"id":"project-1","name":"Current"}]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/public/organizations/projects":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":[{"id":"project-1","name":"Project One"}]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/public/projects":
			defer func() { _ = r.Body.Close() }()
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode create body: %v", err)
			}
			if body["name"] != "Project One" || body["retention"].(float64) != 7 {
				t.Fatalf("unexpected create body: %#v", body)
			}
			metadata, ok := body["metadata"].(map[string]interface{})
			if !ok || metadata["team"] != "platform" {
				t.Fatalf("unexpected metadata: %#v", body)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"project-1","name":"Project One","metadata":{"team":"platform"},"retentionDays":7}`))
		case r.Method == http.MethodPut && r.URL.Path == "/api/public/projects/project-1":
			defer func() { _ = r.Body.Close() }()
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode update body: %v", err)
			}
			if body["name"] != "Project One Updated" || body["retention"].(float64) != 0 {
				t.Fatalf("unexpected update body: %#v", body)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"project-1","name":"Project One Updated","retentionDays":0}`))
		case r.Method == http.MethodDelete && r.URL.Path == "/api/public/projects/project-1":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte(`{"success":true,"message":"Project deletion scheduled"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/public/projects/project-1/memberships":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"memberships":[{"userId":"user-1","role":"MEMBER","email":"user@example.com","name":"User One"}]}`))
		case r.Method == http.MethodPut && r.URL.Path == "/api/public/projects/project-1/memberships":
			defer func() { _ = r.Body.Close() }()
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode membership upsert body: %v", err)
			}
			if body["userId"] != "user-1" || body["role"] != "ADMIN" {
				t.Fatalf("unexpected membership upsert body: %#v", body)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"userId":"user-1","role":"ADMIN","email":"user@example.com","name":"User One"}`))
		case r.Method == http.MethodDelete && r.URL.Path == "/api/public/projects/project-1/memberships":
			defer func() { _ = r.Body.Close() }()
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode membership delete body: %v", err)
			}
			if body["userId"] != "user-1" {
				t.Fatalf("unexpected membership delete body: %#v", body)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"message":"Membership deleted","userId":"user-1"}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		URL:      server.URL,
		Username: "org-pk",
		Password: "org-sk",
		Timeout:  5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if _, err := client.GetProject(context.Background()); err != nil {
		t.Fatalf("GetProject() error = %v", err)
	}
	if _, err := client.ListOrganizationProjects(context.Background()); err != nil {
		t.Fatalf("ListOrganizationProjects() error = %v", err)
	}
	if _, err := client.CreateProject(context.Background(), "Project One", map[string]interface{}{"team": "platform"}, 7); err != nil {
		t.Fatalf("CreateProject() error = %v", err)
	}
	retention := 0
	if _, err := client.UpdateProject(context.Background(), "project-1", "Project One Updated", nil, &retention); err != nil {
		t.Fatalf("UpdateProject() error = %v", err)
	}
	if _, err := client.DeleteProject(context.Background(), "project-1"); err != nil {
		t.Fatalf("DeleteProject() error = %v", err)
	}
	if _, err := client.ListProjectMemberships(context.Background(), "project-1"); err != nil {
		t.Fatalf("ListProjectMemberships() error = %v", err)
	}
	if _, err := client.UpsertProjectMembership(context.Background(), "project-1", "user-1", "ADMIN"); err != nil {
		t.Fatalf("UpsertProjectMembership() error = %v", err)
	}
	if _, err := client.DeleteProjectMembership(context.Background(), "project-1", "user-1"); err != nil {
		t.Fatalf("DeleteProjectMembership() error = %v", err)
	}
}

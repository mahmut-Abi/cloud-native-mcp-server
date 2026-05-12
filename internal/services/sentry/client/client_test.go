package client

import (
	"context"
	"net/http"
	"net/http/httptest"
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
				URL:       "https://sentry.example.com",
				AuthToken: "token",
			},
		},
		{
			name: "valid api URL",
			opts: &ClientOptions{
				URL:       "https://sentry.example.com/api/0",
				AuthToken: "token",
			},
		},
		{
			name: "missing URL",
			opts: &ClientOptions{
				AuthToken: "token",
			},
			wantErr: true,
		},
		{
			name: "missing token",
			opts: &ClientOptions{
				URL: "https://sentry.example.com",
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

func TestListProjectsUsesBearerAuthAndPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/0/organizations/acme/projects/" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer token" {
			t.Fatalf("unexpected auth header: %q", got)
		}
		w.Header().Set("Link", `<https://example.invalid/api/0/organizations/acme/projects/?cursor=abc>; rel="next"; results="true"; cursor="abc"`)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"slug":"frontend"}]`))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		URL:       server.URL,
		AuthToken: "token",
		Timeout:   5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	result, pagination, err := client.ListProjects(context.Background(), "acme", nil)
	if err != nil {
		t.Fatalf("ListProjects() error = %v", err)
	}
	if len(result) != 1 || result[0]["slug"] != "frontend" {
		t.Fatalf("unexpected result: %#v", result)
	}
	if pagination == nil || pagination.Next == nil || pagination.Next.Cursor != "abc" {
		t.Fatalf("unexpected pagination: %#v", pagination)
	}
}

func TestGetIssueEventUsesExpectedPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/0/issues/123/events/abc123/" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"abc123","title":"event"}`))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		URL:       server.URL,
		AuthToken: "token",
		Timeout:   5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	result, err := client.GetIssueEvent(context.Background(), "123", "abc123")
	if err != nil {
		t.Fatalf("GetIssueEvent() error = %v", err)
	}
	if result["id"] != "abc123" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

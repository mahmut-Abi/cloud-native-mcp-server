package client

import (
	"context"
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
				URL:       "https://langfuse.example.com",
				PublicKey: "pk-test",
				SecretKey: "sk-test",
			},
		},
		{
			name: "valid api public URL",
			opts: &ClientOptions{
				URL:       "https://langfuse.example.com/api/public",
				PublicKey: "pk-test",
				SecretKey: "sk-test",
			},
		},
		{
			name: "missing URL",
			opts: &ClientOptions{
				PublicKey: "pk-test",
				SecretKey: "sk-test",
			},
			wantErr: true,
		},
		{
			name: "missing public key",
			opts: &ClientOptions{
				URL:       "https://langfuse.example.com",
				SecretKey: "sk-test",
			},
			wantErr: true,
		},
		{
			name: "missing secret key",
			opts: &ClientOptions{
				URL:       "https://langfuse.example.com",
				PublicKey: "pk-test",
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
		URL:       server.URL,
		PublicKey: "pk-test",
		SecretKey: "sk-test",
		Timeout:   5 * time.Second,
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
		URL:       server.URL,
		PublicKey: "pk-test",
		SecretKey: "sk-test",
		Timeout:   5 * time.Second,
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
		URL:       server.URL,
		PublicKey: "pk-test",
		SecretKey: "sk-test",
		Timeout:   5 * time.Second,
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
		URL:       server.URL,
		PublicKey: "pk-test",
		SecretKey: "sk-test",
		Timeout:   5 * time.Second,
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

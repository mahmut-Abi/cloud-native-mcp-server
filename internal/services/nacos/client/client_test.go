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
			name: "valid url",
			opts: &ClientOptions{
				URL: "http://localhost:8848/nacos",
			},
		},
		{
			name:    "missing url",
			opts:    &ClientOptions{},
			wantErr: true,
		},
		{
			name:    "invalid url",
			opts:    &ClientOptions{URL: "://bad"},
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
				t.Fatal("expected non-nil client")
			}
		})
	}
}

func TestListNamespacesUsesLoginToken(t *testing.T) {
	var sawLogin bool
	var sawNamespaces bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/nacos/v1/auth/users/login", "/nacos/v1/auth/login":
			sawLogin = true
			if err := r.ParseForm(); err != nil {
				t.Fatalf("ParseForm() error = %v", err)
			}
			if r.Form.Get("username") != "nacos" || r.Form.Get("password") != "secret" {
				t.Fatalf("unexpected login form: %v", r.Form)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"accessToken":"token-1","tokenTtl":1800}`))
		case "/nacos/v1/console/namespaces":
			sawNamespaces = true
			if got := r.URL.Query().Get("accessToken"); got != "token-1" {
				t.Fatalf("expected accessToken token-1, got %q", got)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":[{"namespace":"public","namespaceShowName":"Public"}]}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		URL:      server.URL + "/nacos",
		Username: "nacos",
		Password: "secret",
		Timeout:  5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	result, err := client.ListNamespaces(context.Background())
	if err != nil {
		t.Fatalf("ListNamespaces() error = %v", err)
	}
	if !sawLogin || !sawNamespaces {
		t.Fatalf("expected login and namespaces request, got login=%v namespaces=%v", sawLogin, sawNamespaces)
	}
	if len(result) != 1 || result[0]["namespace"] != "public" {
		t.Fatalf("unexpected namespaces result: %#v", result)
	}
}

func TestListServices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/nacos/v1/ns/service/list" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		query := r.URL.Query()
		if query.Get("namespaceId") != "ns-a" || query.Get("groupName") != "DEFAULT_GROUP" {
			t.Fatalf("unexpected query: %v", query)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"count":2,"doms":["svc-a","svc-b"]}`))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		URL:         server.URL + "/nacos",
		AccessToken: "token",
		Timeout:     5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	params := url.Values{}
	params.Set("namespaceId", "ns-a")
	params.Set("groupName", "DEFAULT_GROUP")
	result, err := client.ListServices(context.Background(), params)
	if err != nil {
		t.Fatalf("ListServices() error = %v", err)
	}
	if result["count"].(float64) != 2 {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestGetConfigReturnsTextPayload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/nacos/v1/cs/configs" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("spring.application.name=dify"))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		URL:         server.URL + "/nacos",
		AccessToken: "token",
		Timeout:     5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	params := url.Values{}
	params.Set("dataId", "application.properties")
	params.Set("group", "DEFAULT_GROUP")
	result, err := client.GetConfig(context.Background(), params)
	if err != nil {
		t.Fatalf("GetConfig() error = %v", err)
	}
	if result["format"] != "properties" {
		t.Fatalf("expected properties format, got %#v", result["format"])
	}
	if result["content"] != "spring.application.name=dify" {
		t.Fatalf("unexpected content: %#v", result["content"])
	}
}

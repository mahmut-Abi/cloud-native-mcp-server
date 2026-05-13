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
		{name: "valid url", opts: &ClientOptions{URL: "https://argocd.example.com"}},
		{name: "missing url", opts: &ClientOptions{}, wantErr: true},
		{name: "invalid url", opts: &ClientOptions{URL: "://bad"}, wantErr: true},
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

func TestListApplicationsUsesSessionToken(t *testing.T) {
	var sawSession bool
	var sawList bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/argocd/api/v1/session":
			sawSession = true
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"token":"argocd-token"}`))
		case "/argocd/api/v1/applications":
			sawList = true
			if got := r.Header.Get("Authorization"); got != "Bearer argocd-token" {
				t.Fatalf("expected bearer token, got %q", got)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"items":[{"metadata":{"name":"demo"}}]}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		URL:      server.URL + "/argocd",
		Username: "admin",
		Password: "secret",
		Timeout:  5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	result, err := client.ListApplications(context.Background(), url.Values{})
	if err != nil {
		t.Fatalf("ListApplications() error = %v", err)
	}
	if !sawSession || !sawList {
		t.Fatalf("expected session and list requests, got session=%v list=%v", sawSession, sawList)
	}
	if _, ok := result["items"]; !ok {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestGetApplicationManifests(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/applications/demo/manifests" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("revision") != "main" {
			t.Fatalf("unexpected query: %v", r.URL.Query())
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"manifests":["apiVersion: v1\nkind: ConfigMap"]}`))
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

	params := url.Values{}
	params.Set("revision", "main")
	result, err := client.GetApplicationManifests(context.Background(), "demo", params)
	if err != nil {
		t.Fatalf("GetApplicationManifests() error = %v", err)
	}
	if _, ok := result["manifests"]; !ok {
		t.Fatalf("unexpected result: %#v", result)
	}
}

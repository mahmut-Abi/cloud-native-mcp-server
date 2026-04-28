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
			name: "valid client",
			opts: &ClientOptions{
				Address: "http://localhost:3100",
				Timeout: 5 * time.Second,
			},
		},
		{
			name:    "missing address",
			opts:    &ClientOptions{},
			wantErr: true,
		},
		{
			name:    "invalid address",
			opts:    &ClientOptions{Address: "://bad"},
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

func TestQueryRange(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/loki/api/v1/query_range" {
			t.Fatalf("expected /loki/api/v1/query_range, got %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("query"); got != `{app="api"}` {
			t.Fatalf("unexpected query param: %q", got)
		}
		if got := r.URL.Query().Get("direction"); got != "backward" {
			t.Fatalf("unexpected direction: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"success","data":{"resultType":"streams","result":[]}}`))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{Address: server.URL, Timeout: 5 * time.Second})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	start := time.Unix(10, 0).UTC()
	end := time.Unix(20, 0).UTC()
	result, err := client.QueryRange(context.Background(), `{app="api"}`, start, end, 50, "backward", "")
	if err != nil {
		t.Fatalf("QueryRange() error = %v", err)
	}
	if result["status"] != "success" {
		t.Fatalf("unexpected status: %#v", result["status"])
	}
}

func TestGetLabelNames(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/loki/api/v1/labels" {
			t.Fatalf("expected /loki/api/v1/labels, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"success","data":["app","namespace"]}`))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{Address: server.URL, Timeout: 5 * time.Second})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	labels, err := client.GetLabelNames(context.Background(), "", nil, nil)
	if err != nil {
		t.Fatalf("GetLabelNames() error = %v", err)
	}
	if len(labels) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(labels))
	}
}

func TestGetSeries(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/loki/api/v1/series" {
			t.Fatalf("expected /loki/api/v1/series, got %s", r.URL.Path)
		}
		if got := r.URL.Query()["match[]"]; len(got) != 2 {
			t.Fatalf("expected 2 matchers, got %v", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"success","data":[{"app":"api","namespace":"prod"}]}`))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{Address: server.URL, Timeout: 5 * time.Second})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	series, err := client.GetSeries(context.Background(), []string{`{app="api"}`, `{namespace="prod"}`}, nil, nil)
	if err != nil {
		t.Fatalf("GetSeries() error = %v", err)
	}
	if len(series) != 1 || series[0]["app"] != "api" {
		t.Fatalf("unexpected series: %#v", series)
	}
}

func TestTestConnection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"success","data":[]}`))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{Address: server.URL, Timeout: 5 * time.Second})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if err := client.TestConnection(context.Background()); err != nil {
		t.Fatalf("TestConnection() error = %v", err)
	}
}

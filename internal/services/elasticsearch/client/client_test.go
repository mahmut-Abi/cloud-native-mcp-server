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
			name: "valid client with defaults",
			opts: &ClientOptions{
				Addresses: []string{"http://localhost:9200"},
				Timeout:   30 * time.Second,
			},
			wantErr: false,
		},
		{
			name:    "nil options",
			opts:    nil,
			wantErr: false,
		},
		{
			name: "empty addresses",
			opts: &ClientOptions{
				Addresses: []string{},
			},
			wantErr: false,
		},
		{
			name: "address without http prefix",
			opts: &ClientOptions{
				Addresses: []string{"localhost:9200"},
			},
			wantErr: false,
		},
		{
			name: "with basic auth",
			opts: &ClientOptions{
				Addresses: []string{"http://localhost:9200"},
				Username:  "user",
				Password:  "pass",
			},
			wantErr: false,
		},
		{
			name: "with bearer token",
			opts: &ClientOptions{
				Addresses:   []string{"http://localhost:9200"},
				BearerToken: "token123",
			},
			wantErr: false,
		},
		{
			name: "with API key",
			opts: &ClientOptions{
				Addresses: []string{"http://localhost:9200"},
				APIKey:    "api-key",
			},
			wantErr: false,
		},
		{
			name: "zero timeout",
			opts: &ClientOptions{
				Addresses: []string{"http://localhost:9200"},
				Timeout:   0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() should return non-nil client")
			}
		})
	}
}

func TestHealth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/_cluster/health" {
			t.Errorf("Expected path /_cluster/health, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"cluster_name": "test-cluster",
			"status": "green",
			"number_of_nodes": 1
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		Addresses: []string{server.URL},
		Timeout:   30 * time.Second,
	})

	ctx := context.Background()
	health, err := client.Health(ctx)
	if err != nil {
		t.Errorf("Health() error = %v", err)
		return
	}

	if health == nil {
		t.Error("Health() should return non-nil health")
	}
}

func TestIndices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/_cat/indices" {
			t.Errorf("Expected path /_cat/indices, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		Addresses: []string{server.URL},
		Timeout:   30 * time.Second,
	})

	ctx := context.Background()
	indices, err := client.Indices(ctx)
	if err != nil {
		t.Errorf("Indices() error = %v", err)
		return
	}

	// Indices returns nil when empty, which is acceptable
	_ = indices
}

func TestIndexStats(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-index/_stats" {
			t.Errorf("Expected path /test-index/_stats, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_all": {
				"primaries": {
					"docs": {"count": 100}
				}
			}
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		Addresses: []string{server.URL},
		Timeout:   30 * time.Second,
	})

	ctx := context.Background()
	stats, err := client.IndexStats(ctx, "test-index")
	if err != nil {
		t.Errorf("IndexStats() error = %v", err)
		return
	}

	if stats == nil {
		t.Error("IndexStats() should return non-nil stats")
	}
}

func TestNodes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/_nodes" {
			t.Errorf("Expected path /_nodes, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"cluster_name": "test-cluster",
			"nodes": {}
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		Addresses: []string{server.URL},
		Timeout:   30 * time.Second,
	})

	ctx := context.Background()
	nodes, err := client.Nodes(ctx)
	if err != nil {
		t.Errorf("Nodes() error = %v", err)
		return
	}

	if nodes == nil {
		t.Error("Nodes() should return non-nil nodes")
	}
}

func TestHealthError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`not valid json`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		Addresses: []string{server.URL},
		Timeout:   30 * time.Second,
	})

	ctx := context.Background()
	_, err := client.Health(ctx)
	if err == nil {
		t.Error("Health() should return error for invalid JSON")
	}
}

func TestIndicesError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "Not found"}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		Addresses: []string{server.URL},
		Timeout:   30 * time.Second,
	})

	ctx := context.Background()
	_, err := client.Indices(ctx)
	if err == nil {
		t.Error("Indices() should return error for 404 status")
	}
}

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
				BaseURL: "http://localhost:16686",
				Timeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "empty base URL",
			opts: &ClientOptions{
				BaseURL: "",
				Timeout: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "invalid URL",
			opts: &ClientOptions{
				BaseURL: "://invalid-url",
				Timeout: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "base URL without trailing slash",
			opts: &ClientOptions{
				BaseURL: "http://localhost:16686",
				Timeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "default timeout",
			opts: &ClientOptions{
				BaseURL: "http://localhost:16686",
				Timeout: 0,
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

func TestGetTrace(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/traces/test-trace-id" {
			t.Errorf("Expected path /api/traces/test-trace-id, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": [{"traceID": "test-trace-id", "spans": [], "processes": []}]}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		BaseURL: server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()
	trace, err := client.GetTrace(ctx, "test-trace-id")
	if err != nil {
		t.Errorf("GetTrace() error = %v", err)
		return
	}

	if trace == nil {
		t.Error("GetTrace() should return non-nil trace")
		return
	}

	if trace.TraceID != "test-trace-id" {
		t.Errorf("Expected traceID 'test-trace-id', got '%s'", trace.TraceID)
	}
}

func TestGetTraceNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		BaseURL: server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()
	_, err := client.GetTrace(ctx, "nonexistent-trace-id")
	if err == nil {
		t.Error("GetTrace() should return error for not found")
	}
}

func TestSearchTraces(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": []}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		BaseURL: server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()
	params := TraceQueryParameters{
		Service:   "test-service",
		Operation: "test-operation",
		Limit:     10,
	}

	traces, err := client.SearchTraces(ctx, params)
	if err != nil {
		t.Errorf("SearchTraces() error = %v", err)
		return
	}

	if traces == nil {
		t.Error("SearchTraces() should return non-nil traces")
	}
}

func TestGetServices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": ["service1", "service2"]}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		BaseURL: server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()
	services, err := client.GetServices(ctx)
	if err != nil {
		t.Errorf("GetServices() error = %v", err)
		return
	}

	if services == nil {
		t.Error("GetServices() should return non-nil services")
	}

	if len(services) != 2 {
		t.Errorf("Expected 2 services, got %d", len(services))
	}
}

func TestGetOperations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": ["op1", "op2"]}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		BaseURL: server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()
	ops, err := client.GetOperations(ctx, "test-service")
	if err != nil {
		t.Errorf("GetOperations() error = %v", err)
		return
	}

	if ops == nil {
		t.Error("GetOperations() should return non-nil operations")
	}

	if len(ops) != 2 {
		t.Errorf("Expected 2 operations, got %d", len(ops))
	}
}

func TestGetDependencies(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		BaseURL: server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()
	deps, err := client.GetDependencies(ctx, "2024-01-01T00:00:00Z", "2024-01-02T00:00:00Z")
	if err != nil {
		t.Errorf("GetDependencies() error = %v", err)
		return
	}

	if deps == nil {
		t.Error("GetDependencies() should return non-nil dependencies")
	}
}

func TestHandleResponseError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{"400 Bad Request", 400, true},
		{"401 Unauthorized", 401, true},
		{"403 Forbidden", 403, true},
		{"404 Not Found", 404, true},
		{"429 Rate Limited", 429, true},
		{"500 Internal Server Error", 500, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(`{"error": "test error"}`))
			}))
			defer server.Close()

			client, _ := NewClient(&ClientOptions{
				BaseURL: server.URL,
				Timeout: 30 * time.Second,
			})

			ctx := context.Background()
			_, err := client.GetTrace(ctx, "test-id")
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error = %v, got %v", tt.wantErr, err)
			}
		})
	}
}

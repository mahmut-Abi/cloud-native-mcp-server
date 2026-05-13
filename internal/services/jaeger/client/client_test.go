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
		_, _ = w.Write([]byte(`{"data": [{"traceID": "test-trace-id", "spans": [{"processID":"p1"}], "processes": {"p1":{"serviceName":"test-service","tags":[]}}}]}`))
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
	if trace.Processes["p1"].ServiceName != "test-service" {
		t.Errorf("Expected process service 'test-service', got '%s'", trace.Processes["p1"].ServiceName)
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
		_, _ = w.Write([]byte(`{"data": [{"traceID":"trace-1","spans":[{"processID":"p1"}],"processes":{"p1":{"serviceName":"svc-a","tags":[]}}}]}`))
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
	if len(traces) != 1 {
		t.Fatalf("expected 1 trace, got %d", len(traces))
	}
	if traces[0].Processes["p1"].ServiceName != "svc-a" {
		t.Fatalf("expected process service svc-a, got %q", traces[0].Processes["p1"].ServiceName)
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

func TestGetOperations_ObjectEntries(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[{"name":"GET /health","spanKind":"server"},{"name":"GET /health","spanKind":"consumer"},{"name":"POST /v1/chat","spanKind":"server"}]}`))
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

	if len(ops) != 2 {
		t.Fatalf("Expected 2 deduplicated operations, got %d", len(ops))
	}
	if ops[0] != "GET /health" || ops[1] != "POST /v1/chat" {
		t.Fatalf("Unexpected operations: %#v", ops)
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

func TestGetDependencies_WrappedData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[{"parent":"frontend","child":"api","callCount":42}]}`))
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

	if len(deps) != 1 {
		t.Fatalf("Expected 1 dependency, got %d", len(deps))
	}
	if deps[0].Parent != "frontend" || deps[0].Child != "api" || deps[0].CallCount != 42 {
		t.Fatalf("Unexpected dependency: %#v", deps[0])
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

func TestRetryOnTransientStatusThenSuccess(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts <= 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"error":"temporary unavailable"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": ["service1"]}`))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		BaseURL:        server.URL,
		Timeout:        2 * time.Second,
		MaxRetries:     2,
		RetryBaseDelay: 1 * time.Millisecond,
		RetryMaxDelay:  5 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	services, err := client.GetServices(context.Background())
	if err != nil {
		t.Fatalf("GetServices() unexpected error = %v", err)
	}
	if len(services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(services))
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestRetryStopsAfterMaxRetries(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"error":"temporary unavailable"}`))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		BaseURL:        server.URL,
		Timeout:        2 * time.Second,
		MaxRetries:     1,
		RetryBaseDelay: 1 * time.Millisecond,
		RetryMaxDelay:  5 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.GetServices(context.Background())
	if err == nil {
		t.Fatalf("expected error after retries are exhausted")
	}
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
}

func TestNoRetryOnClientErrorStatus(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"bad request"}`))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		BaseURL:        server.URL,
		Timeout:        2 * time.Second,
		MaxRetries:     3,
		RetryBaseDelay: 1 * time.Millisecond,
		RetryMaxDelay:  5 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.GetServices(context.Background())
	if err == nil {
		t.Fatalf("expected client error")
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt for non-retryable status, got %d", attempts)
	}
}

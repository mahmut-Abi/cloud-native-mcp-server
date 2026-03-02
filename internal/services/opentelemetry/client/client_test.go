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
				Address: "http://localhost:4318",
				Timeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name:    "nil options",
			opts:    nil,
			wantErr: true,
		},
		{
			name: "empty address",
			opts: &ClientOptions{
				Address: "",
				Timeout: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "with basic auth",
			opts: &ClientOptions{
				Address:  "http://localhost:4318",
				Username: "user",
				Password: "pass",
				Timeout:  30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "with bearer token",
			opts: &ClientOptions{
				Address:     "http://localhost:4318",
				BearerToken: "token123",
				Timeout:     30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "with TLS skip verify",
			opts: &ClientOptions{
				Address:       "https://localhost:4318",
				TLSSkipVerify: true,
				Timeout:       30 * time.Second,
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

func TestGetHealth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/healthz" {
			t.Errorf("Expected path /healthz, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "healthy"}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		Address: server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()
	health, err := client.GetHealth(ctx)
	if err != nil {
		t.Errorf("GetHealth() error = %v", err)
		return
	}

	if health == nil {
		t.Error("GetHealth() should return non-nil health")
	}
}

func TestGetHealthError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"error": "service unavailable"}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		Address: server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()
	_, err := client.GetHealth(ctx)
	if err == nil {
		t.Error("GetHealth() should return error for 503 status")
	}
}

func TestGetTraces(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/traces" {
			t.Errorf("Expected path /traces, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"resourceSpans": []}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		Address: server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()

	traces, err := client.GetTraces(ctx, nil, nil, nil, nil, nil)
	if err != nil {
		t.Errorf("GetTraces() error = %v", err)
		return
	}

	if traces == nil {
		t.Error("GetTraces() should return non-nil traces")
	}
}

func TestGetMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/metrics" {
			t.Errorf("Expected path /metrics, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"resourceMetrics": []}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		Address: server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()

	metrics, err := client.GetMetrics(ctx, nil, nil, nil)
	if err != nil {
		t.Errorf("GetMetrics() error = %v", err)
		return
	}

	if metrics == nil {
		t.Error("GetMetrics() should return non-nil metrics")
	}
}

func TestGetLogs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/logs" {
			t.Errorf("Expected path /logs, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"resourceLogs": []}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		Address: server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()

	logs, err := client.GetLogs(ctx, nil, nil, nil, nil, nil)
	if err != nil {
		t.Errorf("GetLogs() error = %v", err)
		return
	}

	if logs == nil {
		t.Error("GetLogs() should return non-nil logs")
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"healthy"}`))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		Address:        server.URL,
		Timeout:        2 * time.Second,
		MaxRetries:     2,
		RetryBaseDelay: 1 * time.Millisecond,
		RetryMaxDelay:  5 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.GetHealth(context.Background())
	if err != nil {
		t.Fatalf("GetHealth() unexpected error = %v", err)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestNoRetryForNonIdempotentMethod(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"error":"temporary unavailable"}`))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		Address:        server.URL,
		Timeout:        2 * time.Second,
		MaxRetries:     3,
		RetryBaseDelay: 1 * time.Millisecond,
		RetryMaxDelay:  5 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.QueryMetrics(context.Background(), "up", nil)
	if err == nil {
		t.Fatalf("expected error for non-idempotent request")
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt for non-idempotent request, got %d", attempts)
	}
}

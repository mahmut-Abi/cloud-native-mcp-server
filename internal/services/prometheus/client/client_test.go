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
				Address: "http://localhost:9090",
				Timeout: 30 * time.Second,
			},
			wantErr: false,
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
			name: "invalid URL",
			opts: &ClientOptions{
				Address: "://invalid-url",
				Timeout: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "with basic auth",
			opts: &ClientOptions{
				Address:  "http://localhost:9090",
				Username: "user",
				Password: "pass",
				Timeout:  30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "with bearer token",
			opts: &ClientOptions{
				Address:     "http://localhost:9090",
				BearerToken: "token123",
				Timeout:     30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "default timeout",
			opts: &ClientOptions{
				Address: "http://localhost:9090",
				Timeout: 0,
			},
			wantErr: false,
		},
		{
			name: "with TLS skip verify",
			opts: &ClientOptions{
				Address:       "https://localhost:9090",
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

func TestQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/query" {
			t.Errorf("Expected path /api/v1/query, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status": "success",
			"data": {
				"resultType": "vector",
				"result": []
			}
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		Address: server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()
	now := time.Now()
	result, err := client.Query(ctx, "up", &now)
	if err != nil {
		t.Errorf("Query() error = %v", err)
		return
	}

	if result == nil {
		t.Error("Query() should return non-nil result")
	}
}

func TestQueryRange(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/query_range" {
			t.Errorf("Expected path /api/v1/query_range, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status": "success",
			"data": {
				"resultType": "matrix",
				"result": []
			}
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		Address: server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()
	end := time.Now()
	start := end.Add(-1 * time.Hour)

	result, err := client.QueryRange(ctx, "up", start, end, "1m")
	if err != nil {
		t.Errorf("QueryRange() error = %v", err)
		return
	}

	if result == nil {
		t.Error("QueryRange() should return non-nil result")
	}
}

func TestGetTargets(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status": "success",
			"data": {
				"activeTargets": [
					{
						"discoveredLabels": {},
						"labels": {"job": "test"},
						"scrapePool": "test",
						"scrapeUrl": "http://localhost:8080/metrics",
						"globalUrl": "http://localhost:8080/metrics",
						"lastError": "",
						"lastScrape": "2024-01-01T00:00:00Z",
						"lastScrapeDuration": 0.001,
						"health": "up",
						"scrapeInterval": "15s",
						"scrapeTimeout": "10s"
					}
				],
				"droppedTargets": []
			}
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		Address: server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()
	targets, err := client.GetTargets(ctx, "active")
	if err != nil {
		t.Errorf("GetTargets() error = %v", err)
		return
	}

	if targets == nil {
		t.Error("GetTargets() should return non-nil targets")
	}
}

func TestGetAlerts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/alerts" {
			t.Errorf("Expected path /api/v1/alerts, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status": "success",
			"data": {
				"alerts": []
			}
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		Address: server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()
	alerts, err := client.GetAlerts(ctx)
	if err != nil {
		t.Errorf("GetAlerts() error = %v", err)
		return
	}

	if alerts == nil {
		t.Error("GetAlerts() should return non-nil alerts")
	}
}

func TestGetRules(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/rules" {
			t.Errorf("Expected path /api/v1/rules, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status": "success",
			"data": {
				"groups": []
			}
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		Address: server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()
	rules, err := client.GetRules(ctx, "all")
	if err != nil {
		t.Errorf("GetRules() error = %v", err)
		return
	}

	if rules == nil {
		t.Error("GetRules() should return non-nil rules")
	}
}

func TestQueryError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{
			"status": "error",
			"errorType": "BadData",
			"error": "invalid query"
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(&ClientOptions{
		Address: server.URL,
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()
	now := time.Now()
	_, err := client.Query(ctx, "invalid_query", &now)
	if err == nil {
		t.Error("Query() should return error for invalid query")
	}
}

func TestRetryOnTransientStatusThenSuccess(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts <= 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"status":"error","error":"temporary unavailable"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"success","data":{"alerts":[]}}`))
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

	_, err = client.GetAlerts(context.Background())
	if err != nil {
		t.Fatalf("GetAlerts() unexpected error = %v", err)
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
		_, _ = w.Write([]byte(`{"status":"error","error":"temporary unavailable"}`))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		Address:        server.URL,
		Timeout:        2 * time.Second,
		MaxRetries:     1,
		RetryBaseDelay: 1 * time.Millisecond,
		RetryMaxDelay:  5 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.GetAlerts(context.Background())
	if err == nil {
		t.Fatalf("expected retry exhaustion error")
	}
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
}

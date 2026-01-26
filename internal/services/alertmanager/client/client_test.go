package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	if client == nil {
		t.Fatal("Client is nil")
	}

	if client.baseURL.String() != "http://localhost:9093" {
		t.Errorf("Expected base URL to be http://localhost:9093, got %s", client.baseURL.String())
	}
}

func TestNewClientWithOptions(t *testing.T) {
	opts := &ClientOptions{
		Address:     "https://alertmanager.example.com:9093",
		Timeout:     60 * time.Second,
		Username:    "user",
		Password:    "pass",
		BearerToken: "token123",
	}

	client, err := NewClientWithOptions(opts)
	if err != nil {
		t.Fatalf("Failed to create client with options: %v", err)
	}

	if client.baseURL.String() != "https://alertmanager.example.com:9093" {
		t.Errorf("Expected base URL to be https://alertmanager.example.com:9093, got %s", client.baseURL.String())
	}

	if client.username != "user" {
		t.Errorf("Expected username to be 'user', got %s", client.username)
	}

	if client.password != "pass" {
		t.Errorf("Expected password to be 'pass', got %s", client.password)
	}

	if client.token != "token123" {
		t.Errorf("Expected token to be 'token123', got %s", client.token)
	}
}

func TestGetStatus(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/status" {
			t.Errorf("Expected path /api/v2/status, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"cluster": {
				"name": "test-cluster",
				"status": "ready"
			},
			"uptime": "24h0m0s",
			"versionInfo": {
				"version": "0.24.0"
			}
		}`))
	}))
	defer server.Close()

	// Create client with test server URL
	opts := &ClientOptions{
		Address: server.URL,
		Timeout: 5 * time.Second,
	}

	client, err := NewClientWithOptions(opts)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test GetStatus
	ctx := context.Background()
	status, err := client.GetStatus(ctx)
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}

	if status == nil {
		t.Fatal("Status is nil")
	}

	// Verify status content
	if cluster, ok := status["cluster"].(map[string]interface{}); ok {
		if name, ok := cluster["name"].(string); !ok || name != "test-cluster" {
			t.Errorf("Expected cluster name to be 'test-cluster', got %v", name)
		}
	} else {
		t.Error("Expected cluster information in status")
	}
}

func TestGetAlerts(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/alerts" {
			t.Errorf("Expected path /api/v2/alerts, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[
			{
				"labels": {
					"alertname": "HighCPU",
					"instance": "server1"
				},
				"annotations": {
					"description": "CPU usage is high"
				},
				"startsAt": "2023-01-01T00:00:00Z",
				"status": {
					"state": "active"
				}
			}
		]`))
	}))
	defer server.Close()

	// Create client with test server URL
	opts := &ClientOptions{
		Address: server.URL,
		Timeout: 5 * time.Second,
	}

	client, err := NewClientWithOptions(opts)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test GetAlerts
	ctx := context.Background()
	alerts, err := client.GetAlerts(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to get alerts: %v", err)
	}

	if len(alerts) != 1 {
		t.Fatalf("Expected 1 alert, got %d", len(alerts))
	}

	// Verify alert content
	alert := alerts[0]
	if labels, ok := alert["labels"].(map[string]interface{}); ok {
		if alertname, ok := labels["alertname"].(string); !ok || alertname != "HighCPU" {
			t.Errorf("Expected alertname to be 'HighCPU', got %v", alertname)
		}
	} else {
		t.Error("Expected labels in alert")
	}
}

func TestInvalidURL(t *testing.T) {
	opts := &ClientOptions{
		Address: "ht tp://invalid url with spaces",
		Timeout: 5 * time.Second,
	}

	_, err := NewClientWithOptions(opts)
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}

package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/jaeger/client"
	"github.com/mark3labs/mcp-go/mcp"
)

type mockJaegerService struct {
	client *client.Client
}

func (m *mockJaegerService) GetClient() *client.Client {
	return m.client
}

func newMockService(t *testing.T, handler http.HandlerFunc) *mockJaegerService {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	c, err := client.NewClient(&client.ClientOptions{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("failed to create jaeger client: %v", err)
	}

	return &mockJaegerService{client: c}
}

func TestGetTracesHandler_UsesDefaultsWithoutPanicking(t *testing.T) {
	var capturedQuery url.Values

	service := newMockService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/traces" {
			t.Fatalf("expected path /api/traces, got %s", r.URL.Path)
		}
		capturedQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[]}`))
	})

	handler := GetTracesHandler(service)
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if got := capturedQuery.Get("limit"); got != "20" {
		t.Fatalf("expected default limit=20, got %q", got)
	}
	if capturedQuery.Get("start") == "" {
		t.Fatal("expected default start query parameter")
	}
	if capturedQuery.Get("end") == "" {
		t.Fatal("expected default end query parameter")
	}
}

func TestGetTraceHandler_MissingTraceIDReturnsError(t *testing.T) {
	handler := GetTraceHandler(&mockJaegerService{})
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	_, err := handler(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for missing trace_id")
	}
	if !strings.Contains(err.Error(), "missing required parameter: trace_id") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetServiceOperationsHandler_MissingServiceReturnsError(t *testing.T) {
	handler := GetServiceOperationsHandler(&mockJaegerService{})
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	_, err := handler(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for missing service")
	}
	if !strings.Contains(err.Error(), "missing required parameter: service") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSearchTracesHandler_ParsesLimitAndTags(t *testing.T) {
	var capturedQuery url.Values

	service := newMockService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/traces" {
			t.Fatalf("expected path /api/traces, got %s", r.URL.Path)
		}
		capturedQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[]}`))
	})

	handler := SearchTracesHandler(service)
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"limit": "150",
				"tags": map[string]interface{}{
					"http.method":      " GET ",
					"http.status_code": 500,
				},
			},
		},
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if got := capturedQuery.Get("limit"); got != "100" {
		t.Fatalf("expected limit capped to 100, got %q", got)
	}
	if got := capturedQuery.Get("tag:http.method"); got != "GET" {
		t.Fatalf("expected tag:http.method=GET, got %q", got)
	}
	if got := capturedQuery.Get("tag:http.status_code"); got != "500" {
		t.Fatalf("expected tag:http.status_code=500, got %q", got)
	}
}

func TestGetServicesSummaryHandler_WithNilClientReturnsError(t *testing.T) {
	handler := GetServicesSummaryHandler(&mockJaegerService{})
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	_, err := handler(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for nil jaeger client")
	}
	if !strings.Contains(err.Error(), "jaeger client is not initialized") {
		t.Fatalf("unexpected error: %v", err)
	}
}

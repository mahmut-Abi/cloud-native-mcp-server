package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
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

	svc := newMockService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/traces" {
			t.Fatalf("expected path /api/traces, got %s", r.URL.Path)
		}
		capturedQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[]}`))
	})

	handler := GetTracesHandler()
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"service": "checkout",
			},
		},
	}

	ctx := client.NewContext(context.Background(), svc.GetClient())
	result, err := handler(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if got := capturedQuery.Get("limit"); got != "20" {
		t.Fatalf("expected default limit=20, got %q", got)
	}
	if got := capturedQuery.Get("service"); got != "checkout" {
		t.Fatalf("expected service=checkout, got %q", got)
	}
	if capturedQuery.Get("start") == "" {
		t.Fatal("expected default start query parameter")
	}
	if capturedQuery.Get("end") == "" {
		t.Fatal("expected default end query parameter")
	}
	if _, err := strconv.ParseInt(capturedQuery.Get("start"), 10, 64); err != nil {
		t.Fatalf("expected numeric microsecond start, got %q: %v", capturedQuery.Get("start"), err)
	}
	if _, err := strconv.ParseInt(capturedQuery.Get("end"), 10, 64); err != nil {
		t.Fatalf("expected numeric microsecond end, got %q: %v", capturedQuery.Get("end"), err)
	}
}

func TestGetTraceHandler_MissingTraceIDReturnsError(t *testing.T) {
	handler := GetTraceHandler()
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

func TestGetTracesHandler_MissingServiceReturnsError(t *testing.T) {
	handler := GetTracesHandler()
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

func TestGetServiceOperationsHandler_MissingServiceReturnsError(t *testing.T) {
	handler := GetServiceOperationsHandler()
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

	svc := newMockService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/traces" {
			t.Fatalf("expected path /api/traces, got %s", r.URL.Path)
		}
		capturedQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[]}`))
	})

	handler := SearchTracesHandler()
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"service": "payments",
				"limit":   "150",
				"tags": map[string]interface{}{
					"http.method":      " GET ",
					"http.status_code": 500,
				},
			},
		},
	}

	ctx := client.NewContext(context.Background(), svc.GetClient())
	result, err := handler(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if got := capturedQuery.Get("limit"); got != "100" {
		t.Fatalf("expected limit capped to 100, got %q", got)
	}
	if got := capturedQuery.Get("service"); got != "payments" {
		t.Fatalf("expected service=payments, got %q", got)
	}
	if got := capturedQuery.Get("tag:http.method"); got != "GET" {
		t.Fatalf("expected tag:http.method=GET, got %q", got)
	}
	if got := capturedQuery.Get("tag:http.status_code"); got != "500" {
		t.Fatalf("expected tag:http.status_code=500, got %q", got)
	}
}

func TestGetTracesSummaryHandler_UsesProcessMapForService(t *testing.T) {
	svc := newMockService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"traceID":"trace-1","spans":[{"operationName":"GET /checkout","duration":1234,"processID":"p1"}],"processes":{"p1":{"serviceName":"checkout-service","tags":[]}}}]}`))
	})

	handler := GetTracesSummaryHandler()
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"service": "checkout-service",
			},
		},
	}

	ctx := client.NewContext(context.Background(), svc.GetClient())
	result, err := handler(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil || len(result.Content) == 0 {
		t.Fatal("expected non-empty result")
	}
	text := result.Content[0].(mcp.TextContent).Text
	if !strings.Contains(text, "checkout-service") {
		t.Fatalf("expected summary to include service name from process map, got %s", text)
	}
}

func TestGetServicesSummaryHandler_WithNilClientReturnsError(t *testing.T) {
	handler := GetServicesSummaryHandler()
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	_, err := handler(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for nil jaeger client")
	}
	if !strings.Contains(err.Error(), "jaeger client not found in context") {
		t.Fatalf("unexpected error: %v", err)
	}
}

package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/langfuse/client"
	"github.com/mark3labs/mcp-go/mcp"
)

type mockLangfuseService struct {
	client *client.Client
}

func (m *mockLangfuseService) GetClient() *client.Client {
	return m.client
}

func newMockLangfuseService(t *testing.T, handler http.HandlerFunc) *mockLangfuseService {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	c, err := client.NewClient(&client.ClientOptions{
		URL:      server.URL,
		Username: "pk-test",
		Password: "sk-test",
		Timeout:  5 * time.Second,
	})
	if err != nil {
		t.Fatalf("failed to create langfuse client: %v", err)
	}

	return &mockLangfuseService{client: c}
}

func TestHandleGetMetricsNormalizesTimestampFilters(t *testing.T) {
	service := newMockLangfuseService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/public/metrics" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		payload := r.URL.Query().Get("query")
		if !strings.Contains(payload, `"fromTimestamp":"2026-05-13T00:00:00Z"`) {
			t.Fatalf("expected normalized fromTimestamp, got %s", payload)
		}
		if !strings.Contains(payload, `"toTimestamp":"2026-05-13T13:29:29Z"`) {
			t.Fatalf("expected normalized toTimestamp, got %s", payload)
		}
		if strings.Contains(payload, `"field":"Timestamp"`) {
			t.Fatalf("expected timestamp filters to be lifted out, got %s", payload)
		}
		if !strings.Contains(payload, `"metrics":[{"aggregation":"sum","measure":"usage"},{"aggregation":"count","measure":"count"}]`) {
			t.Fatalf("expected metrics payload preserved, got %s", payload)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[]}`))
	})

	handler := HandleGetMetrics(service)
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"query": map[string]interface{}{
					"view": "traces",
					"dimensions": []interface{}{
						map[string]interface{}{"field": "name"},
					},
					"metrics": []interface{}{
						map[string]interface{}{"measure": "usage", "aggregation": "sum"},
						map[string]interface{}{"measure": "count", "aggregation": "count"},
					},
					"filters": []interface{}{
						map[string]interface{}{"field": "Timestamp", "operator": ">=", "value": "2026-05-13T00:00:00Z"},
						map[string]interface{}{"field": "Timestamp", "operator": "<=", "value": "2026-05-13T13:29:29Z"},
					},
				},
			},
		},
	}

	if _, err := handler(context.Background(), req); err != nil {
		t.Fatalf("HandleGetMetrics() error = %v", err)
	}
}

func TestNormalizeMetricsQueryAddsTypeAndColumn(t *testing.T) {
	query, err := normalizeMetricsQuery(map[string]interface{}{
		"view":          "traces",
		"fromTimestamp": "2026-05-13T00:00:00Z",
		"toTimestamp":   "2026-05-13T01:00:00Z",
		"filters": []interface{}{
			map[string]interface{}{"field": "name", "operator": "=", "value": "chat"},
		},
	})
	if err != nil {
		t.Fatalf("normalizeMetricsQuery() error = %v", err)
	}

	filters, ok := query["filters"].([]map[string]interface{})
	if !ok || len(filters) != 1 {
		t.Fatalf("expected one normalized filter, got %#v", query["filters"])
	}
	if filters[0]["column"] != "name" {
		t.Fatalf("expected field to normalize to column, got %#v", filters[0])
	}
	if filters[0]["type"] != "string" {
		t.Fatalf("expected inferred type string, got %#v", filters[0]["type"])
	}
}

func TestNormalizeMembershipRole(t *testing.T) {
	role, err := normalizeMembershipRole("admin")
	if err != nil {
		t.Fatalf("normalizeMembershipRole() error = %v", err)
	}
	if role != "ADMIN" {
		t.Fatalf("expected ADMIN, got %q", role)
	}

	if _, err := normalizeMembershipRole("superuser"); err == nil {
		t.Fatal("expected invalid role error")
	}
}

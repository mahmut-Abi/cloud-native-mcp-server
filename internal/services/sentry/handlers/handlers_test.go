package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/sentry/client"
	"github.com/mark3labs/mcp-go/mcp"
)

type mockSentryService struct {
	client              *client.Client
	defaultOrganization string
	defaultProject      string
}

func (m *mockSentryService) GetClient() *client.Client {
	return m.client
}

func (m *mockSentryService) GetDefaultOrganization() string {
	return m.defaultOrganization
}

func (m *mockSentryService) GetDefaultProject() string {
	return m.defaultProject
}

func newMockService(t *testing.T, defaultOrg, defaultProject string, handler http.HandlerFunc) *mockSentryService {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	c, err := client.NewClient(&client.ClientOptions{
		URL:       server.URL,
		AuthToken: "token",
		Timeout:   5 * time.Second,
	})
	if err != nil {
		t.Fatalf("failed to create sentry client: %v", err)
	}

	return &mockSentryService{
		client:              c,
		defaultOrganization: defaultOrg,
		defaultProject:      defaultProject,
	}
}

func decodeResult(t *testing.T, result *mcp.CallToolResult) map[string]interface{} {
	t.Helper()

	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result.Content) == 0 {
		t.Fatal("expected at least one content item")
	}

	textContent, ok := mcp.AsTextContent(result.Content[0])
	if !ok {
		t.Fatalf("expected text content, got %T", result.Content[0])
	}

	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(textContent.Text), &payload); err != nil {
		t.Fatalf("failed to decode tool result: %v", err)
	}
	return payload
}

func TestHandleTestConnectionListsOrganizationsWithoutDefaultScope(t *testing.T) {
	service := newMockService(t, "", "", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/0/organizations/" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"slug":"acme"}]`))
	})

	handler := HandleTestConnection(service)
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	ctx := client.NewContext(context.Background(), service.GetClient())
	result, err := handler(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	payload := decodeResult(t, result)
	if payload["status"] != "ok" {
		t.Fatalf("unexpected payload: %#v", payload)
	}
	if payload["mode"] != "organizations" {
		t.Fatalf("unexpected mode: %#v", payload["mode"])
	}
}

func TestHandleGetProjectUsesDefaultOrganizationAndProject(t *testing.T) {
	service := newMockService(t, "acme", "frontend", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/0/projects/acme/frontend/" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"slug":"frontend","id":"42"}`))
	})

	handler := HandleGetProject(service)
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	ctx := client.NewContext(context.Background(), service.GetClient())
	result, err := handler(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	payload := decodeResult(t, result)
	if payload["slug"] != "frontend" {
		t.Fatalf("unexpected payload: %#v", payload)
	}
}

func TestHandleListIssuesSummaryCompactsPayloadAndUsesDefaults(t *testing.T) {
	var capturedQuery url.Values

	service := newMockService(t, "acme", "", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/0/organizations/acme/issues/" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		capturedQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{
			"id":"123",
			"shortId":"ACME-1",
			"title":"panic in checkout",
			"level":"error",
			"status":"unresolved",
			"count":"17",
			"userCount":3,
			"project":{"slug":"frontend"},
			"metadata":{"type":"should-be-dropped"}
		}]`))
	})

	handler := HandleListIssuesSummary(service)
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"query":       "is:unresolved",
				"environment": []interface{}{"prod"},
				"project_ids": []interface{}{"11"},
			},
		},
	}

	ctx := client.NewContext(context.Background(), service.GetClient())
	result, err := handler(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if capturedQuery.Get("query") != "is:unresolved" {
		t.Fatalf("unexpected query params: %#v", capturedQuery)
	}
	if got := capturedQuery.Get("environment"); got != "prod" {
		t.Fatalf("expected environment=prod, got %q", got)
	}
	if got := capturedQuery.Get("project"); got != "11" {
		t.Fatalf("expected project=11, got %q", got)
	}

	payload := decodeResult(t, result)
	data, ok := payload["data"].([]interface{})
	if !ok || len(data) != 1 {
		t.Fatalf("unexpected payload data: %#v", payload)
	}
	item, ok := data[0].(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected item type: %T", data[0])
	}
	if _, exists := item["metadata"]; exists {
		t.Fatalf("summary should not include metadata: %#v", item)
	}
	if item["shortId"] != "ACME-1" || item["title"] != "panic in checkout" {
		t.Fatalf("unexpected summary item: %#v", item)
	}
}

func TestHandleListIssueEventsEncodesFilters(t *testing.T) {
	var capturedQuery url.Values

	service := newMockService(t, "acme", "", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/0/issues/777/events/" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		capturedQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"eventID":"abc123"}]`))
	})

	handler := HandleListIssueEvents(service)
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"issue_id":    "777",
				"full":        true,
				"sample":      false,
				"environment": []interface{}{"prod"},
				"limit":       25.0,
			},
		},
	}

	ctx := client.NewContext(context.Background(), service.GetClient())
	result, err := handler(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if capturedQuery.Get("full") != "true" {
		t.Fatalf("expected full=true, got %#v", capturedQuery)
	}
	if capturedQuery.Get("sample") != "false" {
		t.Fatalf("expected sample=false, got %#v", capturedQuery)
	}
	if capturedQuery.Get("limit") != "25" {
		t.Fatalf("expected limit=25, got %#v", capturedQuery)
	}
	if capturedQuery.Get("environment") != "prod" {
		t.Fatalf("expected environment=prod, got %#v", capturedQuery)
	}
}

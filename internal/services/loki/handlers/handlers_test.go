package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/loki/client"
	"github.com/mark3labs/mcp-go/mcp"
)

type mockLokiService struct {
	client *client.Client
}

func (m *mockLokiService) GetClient() *client.Client {
	return m.client
}

func newMockService(t *testing.T, handler http.HandlerFunc) *mockLokiService {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	c, err := client.NewClient(&client.ClientOptions{
		Address: server.URL,
		Timeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("failed to create loki client: %v", err)
	}

	return &mockLokiService{client: c}
}

func TestQueryRangeHandlerAcceptsNestedParams(t *testing.T) {
	var captured url.Values

	service := newMockService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/loki/api/v1/query_range" {
			t.Fatalf("expected path /loki/api/v1/query_range, got %s", r.URL.Path)
		}
		captured = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"success","data":{"resultType":"streams","result":[]}}`))
	})

	handler := QueryRangeHandler(service)
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"params": map[string]interface{}{
					"query": `{app="api"}`,
					"limit": "25",
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
	if got := captured.Get("limit"); got != "25" {
		t.Fatalf("expected limit=25, got %q", got)
	}
	if captured.Get("start") == "" || captured.Get("end") == "" {
		t.Fatal("expected default start/end values to be set")
	}
}

func TestGetSeriesHandlerMissingMatchers(t *testing.T) {
	handler := GetSeriesHandler(&mockLokiService{})
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	_, err := handler(context.Background(), req)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "missing required parameter: matchers") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestQueryLogsSummaryHandlerBuildsSummary(t *testing.T) {
	service := newMockService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"status":"success",
			"data":{
				"resultType":"streams",
				"result":[
					{
						"stream":{"app":"api","namespace":"prod"},
						"values":[
							["1714300000000000000","first line"],
							["1714300001000000000","second line"]
						]
					}
				]
			}
		}`))
	})

	handler := QueryLogsSummaryHandler(service)
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"query": `{app="api"}`,
			},
		},
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil || len(result.Content) == 0 {
		t.Fatal("expected non-empty result")
	}
	text := result.Content[0].(mcp.TextContent).Text
	if !strings.Contains(text, `"stream_count":1`) {
		t.Fatalf("unexpected summary payload: %s", text)
	}
	if !strings.Contains(text, `"sample_lines":["first line","second line"]`) {
		t.Fatalf("unexpected sample lines: %s", text)
	}
}

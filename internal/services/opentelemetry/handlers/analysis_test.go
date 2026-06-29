package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/opentelemetry/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func newOTelClientForTest(t *testing.T, handler http.HandlerFunc) *client.Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	c, err := client.NewClient(&client.ClientOptions{
		Address: server.URL,
		Timeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("failed to create otel client: %v", err)
	}
	return c
}

func TestHandleGetConfigSummary_YAMLConfig(t *testing.T) {
	c := newOTelClientForTest(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/config" {
			t.Fatalf("expected /config, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/x-yaml")
		_, _ = w.Write([]byte(`
receivers:
  otlp:
    protocols:
      grpc: {}
processors:
  batch: {}
  memory_limiter: {}
exporters:
  otlphttp:
    endpoint: https://tempo.example.com
service:
  telemetry:
    metrics:
      address: 0.0.0.0:8888
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch, memory_limiter]
      exporters: [otlphttp]
`))
	})

	handler := HandleGetConfigSummary()
	result, err := handler(client.NewContext(context.Background(), c), mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	text := result.Content[0].(mcp.TextContent).Text

	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(text), &payload); err != nil {
		t.Fatalf("expected JSON response, got %v", err)
	}

	components := payload["components"].(map[string]interface{})
	receivers := components["receivers"].(map[string]interface{})
	if receivers["count"].(float64) != 1 {
		t.Fatalf("expected receiver count 1, got %#v", receivers["count"])
	}
	signals := payload["signals"].(map[string]interface{})
	if signals["traces"].(float64) != 1 {
		t.Fatalf("expected one traces pipeline, got %#v", signals["traces"])
	}
}

func TestAnalyzePipelines_FindsBrokenRefsAndSamplingGap(t *testing.T) {
	model := buildCollectorModel(map[string]interface{}{
		"receivers": map[string]interface{}{
			"otlp": map[string]interface{}{},
		},
		"processors": map[string]interface{}{
			"batch": map[string]interface{}{},
		},
		"exporters": map[string]interface{}{
			"otlphttp": map[string]interface{}{},
		},
		"service": map[string]interface{}{
			"pipelines": map[string]interface{}{
				"traces": map[string]interface{}{
					"receivers":  []interface{}{"otlp"},
					"processors": []interface{}{"batch", "missing_proc"},
					"exporters":  []interface{}{"missing_exporter"},
				},
			},
		},
	})

	result := analyzePipelines(model, map[string]interface{}{"status": "ok"}, nil, "", "")
	pipelines := result["pipelines"].([]map[string]interface{})
	if len(pipelines) != 1 {
		t.Fatalf("expected 1 pipeline analysis, got %d", len(pipelines))
	}

	findings := pipelines[0]["findings"].([]finding)
	joined := make([]string, 0, len(findings))
	for _, item := range findings {
		joined = append(joined, item.Message)
	}
	all := strings.Join(joined, "\n")
	if !strings.Contains(all, "missing_proc") {
		t.Fatalf("expected missing processor finding, got %s", all)
	}
	if !strings.Contains(all, "missing_exporter") {
		t.Fatalf("expected missing exporter finding, got %s", all)
	}
	if !strings.Contains(all, "sampling") {
		t.Fatalf("expected sampling guidance, got %s", all)
	}
}

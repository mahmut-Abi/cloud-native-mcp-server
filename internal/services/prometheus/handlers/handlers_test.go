package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/prometheus/client"
)

type testResponse struct {
	path       string
	response   string
	statusCode int
}

func setupTestServer(t *testing.T, responses []testResponse) (*httptest.Server, *client.Client) {
	responseMap := make(map[string]testResponse)
	for _, resp := range responses {
		responseMap[resp.path] = resp
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if resp, ok := responseMap[r.URL.Path]; ok {
			w.WriteHeader(resp.statusCode)
			_, _ = w.Write([]byte(resp.response))
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{}"))
		}
	}))

	opts := &client.ClientOptions{
		Address: server.URL,
	}
	c, err := client.NewClient(opts)
	require.NoError(t, err)

	return server, c
}

func TestHandleQuery(t *testing.T) {
	response := `{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"up"},"value":[1234567890,"1"]}]}}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/api/v1/query", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleQuery(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"query": "up",
			},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
}

func TestHandleQueryMissingQuery(t *testing.T) {
	server, c := setupTestServer(t, []testResponse{})
	defer server.Close()

	handler := HandleQuery(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsError)
}

func TestHandleQueryRange(t *testing.T) {
	response := `{"status":"success","data":{"resultType":"matrix","result":[{"metric":{"__name__":"up"},"values":[]}]}}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/api/v1/query_range", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleQueryRange(c)
	ctx := context.Background()
	start := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	end := time.Now().Format(time.RFC3339)
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"query": "up",
				"start": start,
				"end":   end,
				"step":  "15s",
			},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	if result.IsError {
		t.Logf("Error result: %v", result.Content)
	}
	assert.False(t, result.IsError)
}

func TestHandleQueryRangeMissingParameters(t *testing.T) {
	server, c := setupTestServer(t, []testResponse{})
	defer server.Close()

	handler := HandleQueryRange(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsError)
}

func TestHandleGetTargets(t *testing.T) {
	response := `{"status":"success","data":{"activeTargets":[{"labels":{"job":"prometheus"},"health":"up","scrapeUrl":"http://localhost:9090/metrics"}]}}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/api/v1/targets", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleGetTargets(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"state": "active",
			},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
}

func TestHandleGetAlerts(t *testing.T) {
	response := `{"status":"success","data":{"alerts":[{"labels":{"alertname":"TestAlert"},"state":"firing"}]}}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/api/v1/alerts", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleGetAlerts(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
}

func TestHandleGetRules(t *testing.T) {
	response := `{"status":"success","data":{"groups":[{"name":"default","rules":[{"name":"up","type":"recording","health":"ok"}]}]}}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/api/v1/rules", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleGetRules(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"type": "alert",
			},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
}

func TestHandleGetLabelNames(t *testing.T) {
	response := `{"status":"success","data":["__name__","instance","job"]}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/api/v1/labels", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleGetLabelNames(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
}

func TestHandleGetLabelValues(t *testing.T) {
	response := `{"status":"success","data":["localhost:9090","server1:9090"]}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/api/v1/label/instance/values", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleGetLabelValues(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"label": "instance",
			},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
}

func TestHandleGetLabelValuesMissingLabel(t *testing.T) {
	server, c := setupTestServer(t, []testResponse{})
	defer server.Close()

	handler := HandleGetLabelValues(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsError)
}

func TestHandleGetSeries(t *testing.T) {
	response := `{"status":"success","data":[{"__name__":"up","instance":"localhost:9090"},{"__name__":"up","instance":"server1:9090"}]}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/api/v1/series", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleGetSeries(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"match": []interface{}{"up"},
			},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	if result.IsError {
		t.Logf("Error result: %v", result.Content)
	}
	assert.False(t, result.IsError)
}

func TestHandleGetSeriesMissingMatch(t *testing.T) {
	server, c := setupTestServer(t, []testResponse{})
	defer server.Close()

	handler := HandleGetSeries(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsError)
}

func TestHandleTestConnection(t *testing.T) {
	server, c := setupTestServer(t, []testResponse{})
	defer server.Close()

	handler := HandleTestConnection(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
}

func TestHandleGetServerInfo(t *testing.T) {
	response := `{"status":"success","data":{"version":"2.45.0","revision":"main"}}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/api/v1/status/buildinfo", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleGetServerInfo(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
}

func TestHandleGetMetricsMetadata(t *testing.T) {
	response := `{"status":"success","data":{"up":[{"type":"gauge","help":"Whether the target is up","unit":""}]}}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/api/v1/metadata", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleGetMetricsMetadata(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"metric": "up",
			},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
}

func TestHandleGetTargetMetadata(t *testing.T) {
	response := `{"status":"success","data":[{"target":{"labels":{"__name__":"up"}}}]}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/api/v1/targets/metadata", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleGetTargetMetadata(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"metric": "up",
			},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
}

func TestHandleGetTSDBStats(t *testing.T) {
	response := `{"status":"success","data":{"headStats":{"numSeries":1000,"numChunks":5000}}}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/api/v1/tsdb/stats", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleGetTSDBStats(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
}

func TestHandleGetTSDBStatus(t *testing.T) {
	response := `{"status":"success","data":{"headStats":{"minTime":1234567890,"maxTime":1234567900}}}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/api/v1/tsdb/status", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleGetTSDBStatus(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
}

func TestHandleGetRuntimeInfo(t *testing.T) {
	response := `{"status":"success","data":{"version":"2.45.0","goVersion":"go1.21.0"}}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/api/v1/status/runtimeinfo", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleGetRuntimeInfo(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
}

func TestHandleCreateSnapshot(t *testing.T) {
	response := `{"status":"success","data":{"name":"snapshot-12345"}}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/api/v1/admin/tsdb/snapshot", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleCreateSnapshot(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"skipHead": true,
			},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
}

func TestHandleGetWALReplayStatus(t *testing.T) {
	response := `{"status":"success","data":{"minTime":1234567890,"maxTime":1234567900}}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/api/v1/status/walreplay", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleGetWALReplayStatus(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
}

func TestHandleGetTargetsSummary(t *testing.T) {
	response := `{"status":"success","data":{"activeTargets":[{"labels":{"job":"prometheus","instance":"localhost:9090"},"health":"up","scrapeUrl":"http://localhost:9090/metrics"}]}}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/api/v1/targets", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleGetTargetsSummary(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"state": "active",
			},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
}

func TestHandleGetAlertsSummary(t *testing.T) {
	response := `{"status":"success","data":{"alerts":[{"labels":{"alertname":"TestAlert","instance":"localhost:9090"},"state":"firing"}]}}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/api/v1/alerts", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleGetAlertsSummary(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
}

func TestHandleGetRulesSummary(t *testing.T) {
	response := `{"status":"success","data":{"groups":[{"name":"default","rules":[{"name":"up","type":"recording","health":"ok"}]}]}}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/api/v1/rules", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleGetRulesSummary(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"type": "alert",
			},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
}

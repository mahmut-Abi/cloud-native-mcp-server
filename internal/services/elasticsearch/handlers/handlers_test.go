package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/elasticsearch/client"
)

func TestParseLimitWithWarnings(t *testing.T) {
	tests := []struct {
		name     string
		request  mcp.CallToolRequest
		toolName string
		want     int64
	}{
		{
			name: "default limit when not provided",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{},
				},
			},
			toolName: "test_tool",
			want:     defaultLimit,
		},
		{
			name: "valid limit",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{
						"limit": float64(50),
					},
				},
			},
			toolName: "test_tool",
			want:     50,
		},
		{
			name: "limit too high, capped at max",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{
						"limit": float64(2000),
					},
				},
			},
			toolName: "test_tool",
			want:     maxLimit,
		},
		{
			name: "limit zero, use default",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{
						"limit": float64(0),
					},
				},
			},
			toolName: "test_tool",
			want:     defaultLimit,
		},
		{
			name: "limit negative, use default",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{
						"limit": float64(-10),
					},
				},
			},
			toolName: "test_tool",
			want:     defaultLimit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseLimitWithWarnings(tt.request, tt.toolName)
			if got != tt.want {
				t.Errorf("parseLimitWithWarnings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequireStringParam(t *testing.T) {
	tests := []struct {
		name    string
		request mcp.CallToolRequest
		param   string
		want    string
		wantErr bool
	}{
		{
			name: "valid string param",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{
						"index": "test-index",
					},
				},
			},
			param:   "index",
			want:    "test-index",
			wantErr: false,
		},
		{
			name: "missing param",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{},
				},
			},
			param:   "index",
			want:    "",
			wantErr: true,
		},
		{
			name: "empty string param",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{
						"index": "",
					},
				},
			},
			param:   "index",
			want:    "",
			wantErr: true,
		},
		{
			name: "wrong type param",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]interface{}{
						"index": 123,
					},
				},
			},
			param:   "index",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := requireStringParam(tt.request, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("requireStringParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("requireStringParam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetOptionalStringParam(t *testing.T) {
	tests := []struct {
		name string
		args map[string]interface{}
		param string
		want string
	}{
		{
			name: "valid string param",
			args: map[string]interface{}{
				"index": "test-index",
			},
			param: "index",
			want:  "test-index",
		},
		{
			name:  "missing param",
			args:  map[string]interface{}{},
			param: "index",
			want:  "",
		},
		{
			name: "empty string param",
			args: map[string]interface{}{
				"index": "",
			},
			param: "index",
			want:  "",
		},
		{
			name: "wrong type param",
			args: map[string]interface{}{
				"index": 123,
			},
			param: "index",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.args,
				},
			}
			got := getOptionalStringParam(request, tt.param)
			if got != tt.want {
				t.Errorf("getOptionalStringParam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetOptionalBoolParam(t *testing.T) {
	tests := []struct {
		name string
		args map[string]interface{}
		param string
		want *bool
	}{
		{
			name: "true value",
			args: map[string]interface{}{
				"enabled": true,
			},
			param: "enabled",
			want: func() *bool { b := true; return &b }(),
		},
		{
			name: "false value",
			args: map[string]interface{}{
				"enabled": false,
			},
			param: "enabled",
			want: func() *bool { b := false; return &b }(),
		},
		{
			name:  "missing param",
			args:  map[string]interface{}{},
			param: "enabled",
			want:  nil,
		},
		{
			name: "wrong type param",
			args: map[string]interface{}{
				"enabled": "true",
			},
			param: "enabled",
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.args,
				},
			}
			got := getOptionalBoolParam(request, tt.param)
			if tt.want == nil {
				if got != nil {
					t.Errorf("getOptionalBoolParam() = %v, want nil", got)
				}
			} else {
				if got == nil {
					t.Errorf("getOptionalBoolParam() = nil, want %v", *tt.want)
				} else if *got != *tt.want {
					t.Errorf("getOptionalBoolParam() = %v, want %v", *got, *tt.want)
				}
			}
		})
	}
}

func TestMarshalOptimizedResponse(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		toolName string
		wantErr  bool
	}{
		{
			name:     "valid data",
			data:     map[string]interface{}{"key": "value"},
			toolName: "test_tool",
			wantErr:  false,
		},
		{
			name:     "nil data",
			data:     nil,
			toolName: "test_tool",
			wantErr:  false,
		},
		{
			name:     "complex data",
			data:     map[string]interface{}{"nested": map[string]interface{}{"key": "value"}},
			toolName: "test_tool",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := marshalOptimizedResponse(tt.data, tt.toolName)
			if (err != nil) != tt.wantErr {
				t.Errorf("marshalOptimizedResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result == nil {
				t.Error("marshalOptimizedResponse() should return non-nil result")
			}
		})
	}
}

type testResponse struct {
	path        string
	response    string
	statusCode  int
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
		Addresses: []string{server.URL},
	}
	c, err := client.NewClient(opts)
	require.NoError(t, err)

	return server, c
}

func TestHandleHealthCheck(t *testing.T) {
	response := `{"cluster_name":"test","status":"green","number_of_nodes":1}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/_cluster/health", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleHealthCheck(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestHandleHealthCheckError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	opts := &client.ClientOptions{
		Addresses: []string{server.URL},
	}
	c, err := client.NewClient(opts)
	require.NoError(t, err)

	handler := HandleHealthCheck(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestHandleListIndices(t *testing.T) {
	response := `[{"index":"test-index","health":"green"}]`
	server, c := setupTestServer(t, []testResponse{
		{path: "/_cat/indices", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleListIndices(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestHandleIndexStats(t *testing.T) {
	response := `{"indices":{"test-index":{"primaries":{"docs":{"count":100}}}}}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/test-index/_stats", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleIndexStats(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"index": "test-index",
			},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestHandleIndexStatsMissingIndex(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()

	opts := &client.ClientOptions{
		Addresses: []string{server.URL},
	}
	c, err := client.NewClient(opts)
	require.NoError(t, err)

	handler := HandleIndexStats(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestHandleNodes(t *testing.T) {
	response := `{"nodes":{"node1":{"name":"node-1"}}}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/_nodes", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleNodes(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestHandleInfo(t *testing.T) {
	response := `{"name":"test-cluster","version":{"number":"8.0.0"}}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleInfo(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestHandleListIndicesPaginated(t *testing.T) {
	response := `[{"index":"test-index","health":"green","status":"open","uuid":"1234","pri":"1","rep":"1","docs.count":"100","store.size":"1kb"}]`
	server, c := setupTestServer(t, []testResponse{
		{path: "/_cat/indices", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleListIndicesPaginated(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"limit":          float64(10),
				"includeHealth":  true,
				"indexPattern":   "test*",
				"continueToken":  "",
			},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestHandleGetNodesSummary(t *testing.T) {
	response := `{"nodes":[{"name":"node-1","roles":["data"]}],"count":1}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/_nodes", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleGetNodesSummary(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"role":           "data",
				"includeMetrics": false,
				"limit":          float64(10),
			},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestHandleGetClusterHealthSummary(t *testing.T) {
	response := `{"status":"green","number_of_nodes":1,"active_shards":100}`
	server, c := setupTestServer(t, []testResponse{
		{path: "/_cluster/health", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleGetClusterHealthSummary(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"level":          "basic",
				"includeIndices": false,
			},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestHandleGetIndexDetailAdvanced(t *testing.T) {
	responses := []testResponse{
		{path: "/test-index/_mapping", response: `{"test-index":{"mappings":{}}}`, statusCode: http.StatusOK},
		{path: "/test-index/_settings", response: `{"test-index":{"settings":{}}}`, statusCode: http.StatusOK},
		{path: "/test-index/_stats", response: `{"indices":{"test-index":{"primaries":{}}}}`, statusCode: http.StatusOK},
	}
	server, c := setupTestServer(t, responses)
	defer server.Close()

	handler := HandleGetIndexDetailAdvanced(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"index":           "test-index",
				"includeMappings": true,
				"includeSettings": true,
				"includeStats":    true,
				"includeSegments": false,
				"outputFormat":    "structured",
			},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestHandleGetIndexDetailAdvancedMissingIndex(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()

	opts := &client.ClientOptions{
		Addresses: []string{server.URL},
	}
	c, err := client.NewClient(opts)
	require.NoError(t, err)

	handler := HandleGetIndexDetailAdvanced(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := handler(ctx, request)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestHandleGetClusterDetailAdvanced(t *testing.T) {
	responses := []testResponse{
		{path: "/", response: `{"cluster_name":"test-cluster","version":{"number":"8.0.0"}}`, statusCode: http.StatusOK},
		{path: "/_cluster/health", response: `{"status":"green","number_of_nodes":1}`, statusCode: http.StatusOK},
		{path: "/_nodes", response: `{"nodes":{"node1":{"name":"node-1"}}}`, statusCode: http.StatusOK},
		{path: "/_cluster/stats", response: `{"indices":{"count":1}}`, statusCode: http.StatusOK},
		{path: "/_cluster/settings", response: `{"persistent":{},"transient":{}}`, statusCode: http.StatusOK},
	}
	server, c := setupTestServer(t, responses)
	defer server.Close()

	handler := HandleGetClusterDetailAdvanced(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"includeNodes":    true,
				"includeIndices":  false,
				"includeSettings": true,
				"includeStats":    true,
				"includeShards":   false,
				"outputFormat":    "structured",
			},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestHandleSearchIndices(t *testing.T) {
	response := `[{"index":"test-index","health":"green","status":"open","docs.count":"100","store.size":"1kb"}]`
	server, c := setupTestServer(t, []testResponse{
		{path: "/_cat/indices", response: response, statusCode: http.StatusOK},
	})
	defer server.Close()

	handler := HandleSearchIndices(c)
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"query":         "test*",
				"healthStatus":  "green",
				"indexStatus":   "open",
				"minDocCount":   "10",
				"maxDocCount":   "1000",
				"sortBy":        "name",
				"sortOrder":     "asc",
				"limit":         float64(10),
				"continueToken": "",
			},
		},
	}

	result, err := handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}
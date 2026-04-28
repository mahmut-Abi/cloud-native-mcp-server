package serverConfig

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	transport "github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	healthCheckHandler(w, r)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("unexpected content-type: %s", ct)
	}
}

func TestCORSAndLoggingMiddleware_PassThrough(t *testing.T) {
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	sc := &ServerConfig{}
	h := sc.corsMiddleware(loggingMiddleware(base))

	r := httptest.NewRequest(http.MethodGet, "/api/dev/message", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusTeapot {
		t.Fatalf("expected 418 from base handler, got %d", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no CORS allow-origin header by default, got %q", got)
	}
}

func TestCORSOptionsShortCircuit(t *testing.T) {
	sc := &ServerConfig{}
	h := sc.corsMiddleware(http.NotFoundHandler())
	r := httptest.NewRequest(http.MethodOptions, "/any", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for OPTIONS, got %d", w.Code)
	}
}

// Test InitSSEServers - only test that servers are created correctly
func TestInitSSEServers(t *testing.T) {
	sc := &ServerConfig{}
	mcpServer := server.NewMCPServer("test", "1.0.0")
	appConfig := &config.AppConfig{}

	sseServers := sc.InitSSEServers(mcpServer, "127.0.0.1:8080", appConfig)

	assert.NotNil(t, sseServers)
	assert.Len(t, sseServers, 12) // kubernetes, grafana, prometheus, loki, kibana, helm, elasticsearch, alertmanager, jaeger, opentelemetry, aggregate, utilities
	assert.Contains(t, sseServers, "kubernetes")
	assert.Contains(t, sseServers, "grafana")
	assert.Contains(t, sseServers, "prometheus")
	assert.Contains(t, sseServers, "loki")
	assert.Contains(t, sseServers, "kibana")
	assert.Contains(t, sseServers, "helm")
	assert.Contains(t, sseServers, "elasticsearch")
	assert.Contains(t, sseServers, "alertmanager")
	assert.Contains(t, sseServers, "jaeger")
	assert.Contains(t, sseServers, "opentelemetry")
	assert.Contains(t, sseServers, "aggregate")
	assert.Contains(t, sseServers, "utilities")
}

// Test InitSSEServers with custom paths
func TestInitSSEServersWithCustomPaths(t *testing.T) {
	sc := &ServerConfig{}
	mcpServer := server.NewMCPServer("test", "1.0.0")
	appConfig := &config.AppConfig{
		Server: struct {
			Mode            string `yaml:"mode"`
			Addr            string `yaml:"addr"`
			ReadTimeoutSec  int    `yaml:"readTimeoutSec"`
			WriteTimeoutSec int    `yaml:"writeTimeoutSec"`
			IdleTimeoutSec  int    `yaml:"idleTimeoutSec"`
			SSEPaths        struct {
				Kubernetes    string `yaml:"kubernetes"`
				Grafana       string `yaml:"grafana"`
				Prometheus    string `yaml:"prometheus"`
				Loki          string `yaml:"loki"`
				Kibana        string `yaml:"kibana"`
				Helm          string `yaml:"helm"`
				Elasticsearch string `yaml:"elasticsearch"`
				Alertmanager  string `yaml:"alertmanager"`
				Jaeger        string `yaml:"jaeger"`
				OpenTelemetry string `yaml:"opentelemetry"`
				Aggregate     string `yaml:"aggregate"`
				Utilities     string `yaml:"utilities"`
			} `yaml:"ssePaths"`
			StreamableHTTPPaths struct {
				Kubernetes    string `yaml:"kubernetes"`
				Grafana       string `yaml:"grafana"`
				Prometheus    string `yaml:"prometheus"`
				Loki          string `yaml:"loki"`
				Kibana        string `yaml:"kibana"`
				Helm          string `yaml:"helm"`
				Elasticsearch string `yaml:"elasticsearch"`
				Alertmanager  string `yaml:"alertmanager"`
				Jaeger        string `yaml:"jaeger"`
				OpenTelemetry string `yaml:"opentelemetry"`
				Aggregate     string `yaml:"aggregate"`
				Utilities     string `yaml:"utilities"`
			} `yaml:"streamableHttpPaths"`
			CORS struct {
				AllowedOrigins []string `yaml:"allowedOrigins"`
				AllowedMethods []string `yaml:"allowedMethods"`
				AllowedHeaders []string `yaml:"allowedHeaders"`
				MaxAge         int      `yaml:"maxAge"`
			} `yaml:"cors"`
		}{
			SSEPaths: struct {
				Kubernetes    string `yaml:"kubernetes"`
				Grafana       string `yaml:"grafana"`
				Prometheus    string `yaml:"prometheus"`
				Loki          string `yaml:"loki"`
				Kibana        string `yaml:"kibana"`
				Helm          string `yaml:"helm"`
				Elasticsearch string `yaml:"elasticsearch"`
				Alertmanager  string `yaml:"alertmanager"`
				Jaeger        string `yaml:"jaeger"`
				OpenTelemetry string `yaml:"opentelemetry"`
				Aggregate     string `yaml:"aggregate"`
				Utilities     string `yaml:"utilities"`
			}{
				Kubernetes: "/custom/kubernetes/sse",
				Grafana:    "/custom/grafana/sse",
				Aggregate:  "/custom/aggregate/sse",
			},
		},
	}

	sseServers := sc.InitSSEServers(mcpServer, "127.0.0.1:8080", appConfig)

	assert.NotNil(t, sseServers)
	assert.Len(t, sseServers, 12)
}

// Test InitStreamableHTTPServers
func TestInitStreamableHTTPServers(t *testing.T) {
	sc := &ServerConfig{}
	mcpServer := server.NewMCPServer("test", "1.0.0")
	appConfig := &config.AppConfig{}

	httpServers := sc.InitStreamableHTTPServers(mcpServer, "127.0.0.1:8080", appConfig)

	assert.NotNil(t, httpServers)
	assert.Len(t, httpServers, 12) // Same services as SSE
	assert.Contains(t, httpServers, "kubernetes")
	assert.Contains(t, httpServers, "grafana")
	assert.Contains(t, httpServers, "prometheus")
	assert.Contains(t, httpServers, "loki")
	assert.Contains(t, httpServers, "kibana")
	assert.Contains(t, httpServers, "helm")
	assert.Contains(t, httpServers, "elasticsearch")
	assert.Contains(t, httpServers, "alertmanager")
	assert.Contains(t, httpServers, "jaeger")
	assert.Contains(t, httpServers, "opentelemetry")
	assert.Contains(t, httpServers, "aggregate")
	assert.Contains(t, httpServers, "utilities")
}

// Test InitStreamableHTTPServers with custom paths
func TestInitStreamableHTTPServersWithCustomPaths(t *testing.T) {
	sc := &ServerConfig{}
	mcpServer := server.NewMCPServer("test", "1.0.0")
	appConfig := &config.AppConfig{
		Server: struct {
			Mode            string `yaml:"mode"`
			Addr            string `yaml:"addr"`
			ReadTimeoutSec  int    `yaml:"readTimeoutSec"`
			WriteTimeoutSec int    `yaml:"writeTimeoutSec"`
			IdleTimeoutSec  int    `yaml:"idleTimeoutSec"`
			SSEPaths        struct {
				Kubernetes    string `yaml:"kubernetes"`
				Grafana       string `yaml:"grafana"`
				Prometheus    string `yaml:"prometheus"`
				Loki          string `yaml:"loki"`
				Kibana        string `yaml:"kibana"`
				Helm          string `yaml:"helm"`
				Elasticsearch string `yaml:"elasticsearch"`
				Alertmanager  string `yaml:"alertmanager"`
				Jaeger        string `yaml:"jaeger"`
				OpenTelemetry string `yaml:"opentelemetry"`
				Aggregate     string `yaml:"aggregate"`
				Utilities     string `yaml:"utilities"`
			} `yaml:"ssePaths"`
			StreamableHTTPPaths struct {
				Kubernetes    string `yaml:"kubernetes"`
				Grafana       string `yaml:"grafana"`
				Prometheus    string `yaml:"prometheus"`
				Loki          string `yaml:"loki"`
				Kibana        string `yaml:"kibana"`
				Helm          string `yaml:"helm"`
				Elasticsearch string `yaml:"elasticsearch"`
				Alertmanager  string `yaml:"alertmanager"`
				Jaeger        string `yaml:"jaeger"`
				OpenTelemetry string `yaml:"opentelemetry"`
				Aggregate     string `yaml:"aggregate"`
				Utilities     string `yaml:"utilities"`
			} `yaml:"streamableHttpPaths"`
			CORS struct {
				AllowedOrigins []string `yaml:"allowedOrigins"`
				AllowedMethods []string `yaml:"allowedMethods"`
				AllowedHeaders []string `yaml:"allowedHeaders"`
				MaxAge         int      `yaml:"maxAge"`
			} `yaml:"cors"`
		}{
			StreamableHTTPPaths: struct {
				Kubernetes    string `yaml:"kubernetes"`
				Grafana       string `yaml:"grafana"`
				Prometheus    string `yaml:"prometheus"`
				Loki          string `yaml:"loki"`
				Kibana        string `yaml:"kibana"`
				Helm          string `yaml:"helm"`
				Elasticsearch string `yaml:"elasticsearch"`
				Alertmanager  string `yaml:"alertmanager"`
				Jaeger        string `yaml:"jaeger"`
				OpenTelemetry string `yaml:"opentelemetry"`
				Aggregate     string `yaml:"aggregate"`
				Utilities     string `yaml:"utilities"`
			}{
				Kubernetes: "/custom/kubernetes/http",
				Grafana:    "/custom/grafana/http",
				Aggregate:  "/custom/aggregate/http",
			},
		},
	}

	httpServers := sc.InitStreamableHTTPServers(mcpServer, "127.0.0.1:8080", appConfig)

	assert.NotNil(t, httpServers)
	assert.Len(t, httpServers, 12)
}

// Test SetupMultipleRoutes with SSE mode - only test mux creation, not actual HTTP handling
func TestSetupMultipleRoutes_SSEMode(t *testing.T) {
	sc := &ServerConfig{}
	mcpServer := server.NewMCPServer("test", "1.0.0")
	appConfig := &config.AppConfig{}

	mux := http.NewServeMux()
	sc.SetupMultipleRoutes(mux, nil, nil, "sse", appConfig, mcpServer)

	// Verify mux is created and health check is registered
	assert.NotNil(t, mux)

	// Test that health check route is registered
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// Test SetupMultipleRoutes with streamable-http mode
func TestSetupMultipleRoutes_StreamableHTTPMode(t *testing.T) {
	sc := &ServerConfig{}
	mcpServer := server.NewMCPServer("test", "1.0.0")
	appConfig := &config.AppConfig{}

	mux := http.NewServeMux()
	sc.SetupMultipleRoutes(mux, nil, nil, "streamable-http", appConfig, mcpServer)

	assert.NotNil(t, mux)

	// Verify health check works
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSetupMultipleRoutes_StreamableHTTPSetsStreamingHeaders(t *testing.T) {
	sc := &ServerConfig{}
	mcpServer := server.NewMCPServer("test", "1.0.0")
	appConfig := &config.AppConfig{}

	streamableHTTPServers := sc.InitStreamableHTTPServers(mcpServer, "127.0.0.1:8080", appConfig)

	mux := http.NewServeMux()
	sc.SetupMultipleRoutes(mux, nil, streamableHTTPServers, "streamable-http", appConfig, mcpServer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := httptest.NewRequest(http.MethodGet, "/api/aggregate/streamable-http", nil).WithContext(ctx)
	rec := newStreamingResponseRecorder()
	done := make(chan struct{})
	go func() {
		mux.ServeHTTP(rec, req)
		close(done)
	}()

	time.Sleep(25 * time.Millisecond)
	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("streamable-http handler did not stop after context cancellation")
	}

	_, headers, _ := rec.snapshot()
	assert.Equal(t, "no", headers.Get("X-Accel-Buffering"))
	assert.Equal(t, "no-cache, no-transform", headers.Get("Cache-Control"))
}

// Test SetupMultipleRoutes with SSE mode for health route registration
func TestSetupMultipleRoutes_HealthInSSEMode(t *testing.T) {
	sc := &ServerConfig{}
	mcpServer := server.NewMCPServer("test", "1.0.0")
	appConfig := &config.AppConfig{}

	mux := http.NewServeMux()
	sc.SetupMultipleRoutes(mux, nil, nil, "sse", appConfig, mcpServer)

	assert.NotNil(t, mux)

	// Verify health check works
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// Test health check route is always registered
func TestHealthCheckRoute_AlwaysRegistered(t *testing.T) {
	sc := &ServerConfig{}
	mcpServer := server.NewMCPServer("test", "1.0.0")
	appConfig := &config.AppConfig{}

	mux := http.NewServeMux()
	sc.SetupMultipleRoutes(mux, nil, nil, "sse", appConfig, mcpServer)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Test createServiceMCPServer
func TestCreateServiceMCPServer(t *testing.T) {
	sc := &ServerConfig{}

	serviceServer := sc.createServiceMCPServer("kubernetes")

	assert.NotNil(t, serviceServer)
}

func TestSetupMultipleRoutes_RateLimitAppliedToSSEMessageEndpoint(t *testing.T) {
	sc := &ServerConfig{}
	mcpServer := server.NewMCPServer("test", "1.0.0")

	sseServers := map[string]*server.SSEServer{
		"kubernetes": server.NewSSEServer(mcpServer,
			server.WithStaticBasePath(""),
			server.WithSSEEndpoint("/api/kubernetes/sse"),
			server.WithMessageEndpoint("/api/kubernetes/sse/message"),
		),
	}

	appConfig := &config.AppConfig{}
	appConfig.RateLimit.Enabled = true
	appConfig.RateLimit.RequestsPerSecond = 0.0001
	appConfig.RateLimit.Burst = 1

	mux := http.NewServeMux()
	sc.SetupMultipleRoutes(mux, sseServers, nil, "sse", appConfig, mcpServer)

	req1 := httptest.NewRequest(http.MethodPost, "/api/kubernetes/sse/message", nil)
	req1.RemoteAddr = "127.0.0.1:12345"
	w1 := httptest.NewRecorder()
	mux.ServeHTTP(w1, req1)

	req2 := httptest.NewRequest(http.MethodPost, "/api/kubernetes/sse/message", nil)
	req2.RemoteAddr = "127.0.0.1:12345"
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)

	assert.NotEqual(t, http.StatusTooManyRequests, w1.Code, "first request should pass limiter")
	assert.Equal(t, http.StatusTooManyRequests, w2.Code, "second immediate request should be rate limited")

	assert.NotNil(t, sc.rateLimiter)
	assert.NoError(t, sc.Shutdown())
	assert.Nil(t, sc.rateLimiter)
}

// Test createAggregateMCPServer
func TestCreateAggregateMCPServer(t *testing.T) {
	sc := &ServerConfig{}

	aggregateServer := sc.createAggregateMCPServer()

	assert.NotNil(t, aggregateServer)
}

// Test that services are correctly named
func TestAllServiceNames(t *testing.T) {
	expectedServices := []string{
		"kubernetes",
		"grafana",
		"prometheus",
		"loki",
		"kibana",
		"helm",
		"aggregate",
		"utilities",
	}

	sc := &ServerConfig{}
	mcpServer := server.NewMCPServer("test", "1.0.0")
	appConfig := &config.AppConfig{}

	sseServers := sc.InitSSEServers(mcpServer, "127.0.0.1:8080", appConfig)

	for _, serviceName := range expectedServices {
		assert.Contains(t, sseServers, serviceName, "Service %s should be in SSE servers", serviceName)
	}
}

func TestSetupMultipleRoutes_SSEEmitsEndpointAndPreservesQuery(t *testing.T) {
	sc := &ServerConfig{}
	mcpServer := server.NewMCPServer("test", "1.0.0")
	appConfig := &config.AppConfig{}

	sseServers := sc.InitSSEServers(mcpServer, "127.0.0.1:8080", appConfig)

	mux := http.NewServeMux()
	sc.SetupMultipleRoutes(mux, sseServers, nil, "sse", appConfig, mcpServer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := httptest.NewRequest(http.MethodGet, "/api/aggregate/sse?api_key=test-key", nil).WithContext(ctx)
	req.Header.Set("Accept", "text/event-stream")

	rec := newStreamingResponseRecorder()
	done := make(chan struct{})
	go func() {
		mux.ServeHTTP(rec, req)
		close(done)
	}()

	status, headers, body := waitForEndpointFrame(t, rec, 2*time.Second)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "no", headers.Get("X-Accel-Buffering"))
	assert.Empty(t, headers.Get("Cross-Origin-Embedder-Policy"), "security middleware should not wrap SSE route")
	assert.Contains(t, body, "event: endpoint")
	assert.Contains(t, body, "sessionId=")
	assert.Contains(t, body, "api_key=test-key")

	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("SSE handler did not stop after context cancellation")
	}
}

func TestSetupMultipleRoutes_SSECustomPathForElasticsearch(t *testing.T) {
	sc := &ServerConfig{}
	mcpServer := server.NewMCPServer("test", "1.0.0")
	appConfig := &config.AppConfig{}
	appConfig.Server.SSEPaths.Elasticsearch = "/custom/elasticsearch/sse"

	sseServers := sc.InitSSEServers(mcpServer, "127.0.0.1:8080", appConfig)

	mux := http.NewServeMux()
	sc.SetupMultipleRoutes(mux, sseServers, nil, "sse", appConfig, mcpServer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	customReq := httptest.NewRequest(http.MethodGet, "/custom/elasticsearch/sse", nil).WithContext(ctx)
	customReq.Header.Set("Accept", "text/event-stream")
	rec := newStreamingResponseRecorder()
	done := make(chan struct{})
	go func() {
		mux.ServeHTTP(rec, customReq)
		close(done)
	}()

	_, _, body := waitForEndpointFrame(t, rec, 2*time.Second)
	assert.Contains(t, body, "event: endpoint")
	assert.Contains(t, body, "/custom/elasticsearch/sse/message")

	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("custom elasticsearch SSE handler did not stop after context cancellation")
	}

	defaultReq := httptest.NewRequest(http.MethodGet, "/api/elasticsearch/sse", nil)
	defaultW := httptest.NewRecorder()
	mux.ServeHTTP(defaultW, defaultReq)
	assert.Equal(t, http.StatusNotFound, defaultW.Code)
}

func TestSetupMultipleRoutes_SSEAuthQueryForwardingAllowsInitialize(t *testing.T) {
	sc := &ServerConfig{}
	mcpServer := server.NewMCPServer("test", "1.0.0")
	appConfig := &config.AppConfig{}
	appConfig.Auth.Enabled = true
	appConfig.Auth.Mode = "apikey"
	appConfig.Auth.APIKey = "test-key"

	sseServers := sc.InitSSEServers(mcpServer, "127.0.0.1:8080", appConfig)
	mux := http.NewServeMux()
	sc.SetupMultipleRoutes(mux, sseServers, nil, "sse", appConfig, mcpServer)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	client, err := transport.NewSSE(ts.URL + "/api/aggregate/sse?api_key=test-key")
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Start(ctx)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, client.Close())
	}()

	req := transport.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      mcp.NewRequestId(int64(1)),
		Method:  "initialize",
		Params: map[string]any{
			"protocolVersion": "2024-11-05",
			"clientInfo": map[string]any{
				"name":    "test-client",
				"version": "1.0.0",
			},
		},
	}

	resp, err := client.SendRequest(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	if resp != nil {
		assert.Nil(t, resp.Error)
	}
}

type streamingResponseRecorder struct {
	header     http.Header
	statusCode int
	body       bytes.Buffer
	flushed    chan struct{}
	mu         sync.Mutex
}

func newStreamingResponseRecorder() *streamingResponseRecorder {
	return &streamingResponseRecorder{
		header:  make(http.Header),
		flushed: make(chan struct{}, 1),
	}
}

func (r *streamingResponseRecorder) Header() http.Header {
	return r.header
}

func (r *streamingResponseRecorder) WriteHeader(statusCode int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.statusCode == 0 {
		r.statusCode = statusCode
	}
}

func (r *streamingResponseRecorder) Write(p []byte) (int, error) {
	r.mu.Lock()
	if r.statusCode == 0 {
		r.statusCode = http.StatusOK
	}
	n, err := r.body.Write(p)
	r.mu.Unlock()

	if n > 0 {
		r.notifyFlush()
	}

	return n, err
}

func (r *streamingResponseRecorder) Flush() {
	r.notifyFlush()
}

func (r *streamingResponseRecorder) snapshot() (int, http.Header, string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	status := r.statusCode
	if status == 0 {
		status = http.StatusOK
	}

	return status, r.header.Clone(), r.body.String()
}

func (r *streamingResponseRecorder) notifyFlush() {
	select {
	case r.flushed <- struct{}{}:
	default:
	}
}

func waitForEndpointFrame(t *testing.T, rec *streamingResponseRecorder, timeout time.Duration) (int, http.Header, string) {
	t.Helper()

	deadline := time.After(timeout)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-deadline:
			_, _, body := rec.snapshot()
			t.Fatalf("timed out waiting for endpoint SSE frame, body=%q", body)
		case <-rec.flushed:
		case <-ticker.C:
		}

		status, headers, body := rec.snapshot()
		if strings.Contains(body, "event: endpoint") {
			return status, headers, body
		}
	}
}

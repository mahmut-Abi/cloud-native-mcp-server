package serverConfig

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mahmut-Abi/k8s-mcp-server/internal/config"
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
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("expected CORS headers")
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
	assert.Len(t, sseServers, 10) // kubernetes, grafana, prometheus, kibana, helm, elasticsearch, alertmanager, jaeger, aggregate, utilities
	assert.Contains(t, sseServers, "kubernetes")
	assert.Contains(t, sseServers, "grafana")
	assert.Contains(t, sseServers, "prometheus")
	assert.Contains(t, sseServers, "kibana")
	assert.Contains(t, sseServers, "helm")
	assert.Contains(t, sseServers, "elasticsearch")
	assert.Contains(t, sseServers, "alertmanager")
	assert.Contains(t, sseServers, "jaeger")
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
				Kibana        string `yaml:"kibana"`
				Helm          string `yaml:"helm"`
				Elasticsearch string `yaml:"elasticsearch"`
				Alertmanager  string `yaml:"alertmanager"`
				Jaeger        string `yaml:"jaeger"`
				Aggregate     string `yaml:"aggregate"`
				Utilities     string `yaml:"utilities"`
			} `yaml:"ssePaths"`
			StreamableHTTPPaths struct {
				Kubernetes    string `yaml:"kubernetes"`
				Grafana       string `yaml:"grafana"`
				Prometheus    string `yaml:"prometheus"`
				Kibana        string `yaml:"kibana"`
				Helm          string `yaml:"helm"`
				Elasticsearch string `yaml:"elasticsearch"`
				Alertmanager  string `yaml:"alertmanager"`
				Jaeger        string `yaml:"jaeger"`
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
				Kibana        string `yaml:"kibana"`
				Helm          string `yaml:"helm"`
				Elasticsearch string `yaml:"elasticsearch"`
				Alertmanager  string `yaml:"alertmanager"`
				Jaeger        string `yaml:"jaeger"`
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
	assert.Len(t, sseServers, 10)
}

// Test InitStreamableHTTPServers
func TestInitStreamableHTTPServers(t *testing.T) {
	sc := &ServerConfig{}
	mcpServer := server.NewMCPServer("test", "1.0.0")
	appConfig := &config.AppConfig{}

	httpServers := sc.InitStreamableHTTPServers(mcpServer, "127.0.0.1:8080", appConfig)

	assert.NotNil(t, httpServers)
	assert.Len(t, httpServers, 10) // Same services as SSE
	assert.Contains(t, httpServers, "kubernetes")
	assert.Contains(t, httpServers, "grafana")
	assert.Contains(t, httpServers, "prometheus")
	assert.Contains(t, httpServers, "kibana")
	assert.Contains(t, httpServers, "helm")
	assert.Contains(t, httpServers, "elasticsearch")
	assert.Contains(t, httpServers, "alertmanager")
	assert.Contains(t, httpServers, "jaeger")
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
				Kibana        string `yaml:"kibana"`
				Helm          string `yaml:"helm"`
				Elasticsearch string `yaml:"elasticsearch"`
				Alertmanager  string `yaml:"alertmanager"`
				Jaeger        string `yaml:"jaeger"`
				Aggregate     string `yaml:"aggregate"`
				Utilities     string `yaml:"utilities"`
			} `yaml:"ssePaths"`
			StreamableHTTPPaths struct {
				Kubernetes    string `yaml:"kubernetes"`
				Grafana       string `yaml:"grafana"`
				Prometheus    string `yaml:"prometheus"`
				Kibana        string `yaml:"kibana"`
				Helm          string `yaml:"helm"`
				Elasticsearch string `yaml:"elasticsearch"`
				Alertmanager  string `yaml:"alertmanager"`
				Jaeger        string `yaml:"jaeger"`
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
				Kibana        string `yaml:"kibana"`
				Helm          string `yaml:"helm"`
				Elasticsearch string `yaml:"elasticsearch"`
				Alertmanager  string `yaml:"alertmanager"`
				Jaeger        string `yaml:"jaeger"`
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
	assert.Len(t, httpServers, 10)
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

// Test SetupMultipleRoutes with stdio mode
func TestSetupMultipleRoutes_StdioMode(t *testing.T) {
	sc := &ServerConfig{}
	mcpServer := server.NewMCPServer("test", "1.0.0")
	appConfig := &config.AppConfig{}

	mux := http.NewServeMux()
	sc.SetupMultipleRoutes(mux, nil, nil, "stdio", appConfig, mcpServer)

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
	sc.SetupMultipleRoutes(mux, nil, nil, "stdio", appConfig, mcpServer)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Test createServiceMCPServer
func TestCreateServiceMCPServer(t *testing.T) {
	sc := &ServerConfig{}
	mcpServer := server.NewMCPServer("test", "1.0.0")

	serviceServer := sc.createServiceMCPServer("kubernetes", mcpServer)

	assert.NotNil(t, serviceServer)
}

// Test createAggregateMCPServer
func TestCreateAggregateMCPServer(t *testing.T) {
	sc := &ServerConfig{}
	mcpServer := server.NewMCPServer("test", "1.0.0")

	k8sServer := sc.createServiceMCPServer("kubernetes", mcpServer)
	grafanaServer := sc.createServiceMCPServer("grafana", mcpServer)
	promServer := sc.createServiceMCPServer("prometheus", mcpServer)
	kibanaServer := sc.createServiceMCPServer("kibana", mcpServer)
	helmServer := sc.createServiceMCPServer("helm", mcpServer)
	elasticsearchServer := sc.createServiceMCPServer("elasticsearch", mcpServer)
	alertmanagerServer := sc.createServiceMCPServer("alertmanager", mcpServer)
	jaegerServer := sc.createServiceMCPServer("jaeger", mcpServer)

	aggregateServer := sc.createAggregateMCPServer(mcpServer, k8sServer, grafanaServer, promServer, kibanaServer, helmServer, elasticsearchServer, alertmanagerServer, jaegerServer)

	assert.NotNil(t, aggregateServer)
}

// Test that services are correctly named
func TestAllServiceNames(t *testing.T) {
	expectedServices := []string{
		"kubernetes",
		"grafana",
		"prometheus",
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

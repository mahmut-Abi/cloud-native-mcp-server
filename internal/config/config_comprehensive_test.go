package config

import (
	"os"
	"testing"
)

func TestLoadConfigFromYAML(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"nonexistent file", "/nonexistent/config.yaml", true},
		{"empty path", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Load(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnvironmentVariableOverrides(t *testing.T) {
	originalMode := os.Getenv("MCP_MODE")
	originalAddr := os.Getenv("MCP_ADDR")
	originalLogLevel := os.Getenv("MCP_LOG_LEVEL")

	defer func() {
		if originalMode != "" {
			_ = os.Setenv("MCP_MODE", originalMode)
		} else {
			_ = os.Unsetenv("MCP_MODE")
		}
		if originalAddr != "" {
			_ = os.Setenv("MCP_ADDR", originalAddr)
		} else {
			_ = os.Unsetenv("MCP_ADDR")
		}
		if originalLogLevel != "" {
			_ = os.Setenv("MCP_LOG_LEVEL", originalLogLevel)
		} else {
			_ = os.Unsetenv("MCP_LOG_LEVEL")
		}
	}()

	// Test environment variable override
	_ = os.Setenv("MCP_MODE", "http")
	_ = os.Setenv("MCP_ADDR", "127.0.0.1:9090")
	_ = os.Setenv("MCP_LOG_LEVEL", "debug")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Server.Mode != "http" {
		t.Errorf("Expected mode 'http', got '%s'", cfg.Server.Mode)
	}

	if cfg.Server.Addr != "127.0.0.1:9090" {
		t.Errorf("Expected addr '127.0.0.1:9090', got '%s'", cfg.Server.Addr)
	}

	if cfg.Logging.Level != "debug" {
		t.Errorf("Expected log level 'debug', got '%s'", cfg.Logging.Level)
	}
}

func TestKubernetesConfig(t *testing.T) {
	originalKubeconfig := os.Getenv("MCP_KUBECONFIG")
	originalQPS := os.Getenv("MCP_K8S_QPS")
	originalBurst := os.Getenv("MCP_K8S_BURST")

	defer func() {
		if originalKubeconfig != "" {
			_ = os.Setenv("MCP_KUBECONFIG", originalKubeconfig)
		} else {
			_ = os.Unsetenv("MCP_KUBECONFIG")
		}
		if originalQPS != "" {
			_ = os.Setenv("MCP_K8S_QPS", originalQPS)
		} else {
			_ = os.Unsetenv("MCP_K8S_QPS")
		}
		if originalBurst != "" {
			_ = os.Setenv("MCP_K8S_BURST", originalBurst)
		} else {
			_ = os.Unsetenv("MCP_K8S_BURST")
		}
	}()

	_ = os.Setenv("MCP_KUBECONFIG", "/custom/kubeconfig")
	_ = os.Setenv("MCP_K8S_QPS", "200")
	_ = os.Setenv("MCP_K8S_BURST", "400")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Kubernetes.Kubeconfig != "/custom/kubeconfig" {
		t.Errorf("Expected kubeconfig '/custom/kubeconfig', got '%s'", cfg.Kubernetes.Kubeconfig)
	}

	if cfg.Kubernetes.QPS != 200 {
		t.Errorf("Expected QPS 200, got %f", cfg.Kubernetes.QPS)
	}

	if cfg.Kubernetes.Burst != 400 {
		t.Errorf("Expected Burst 400, got %d", cfg.Kubernetes.Burst)
	}
}

func TestPrometheusConfig(t *testing.T) {
	originalAddr := os.Getenv("MCP_PROM_ADDRESS")
	originalEnabled := os.Getenv("MCP_PROM_ENABLED")

	defer func() {
		if originalAddr != "" {
			_ = os.Setenv("MCP_PROM_ADDRESS", originalAddr)
		} else {
			_ = os.Unsetenv("MCP_PROM_ADDRESS")
		}
		if originalEnabled != "" {
			_ = os.Setenv("MCP_PROM_ENABLED", originalEnabled)
		} else {
			_ = os.Unsetenv("MCP_PROM_ENABLED")
		}
	}()

	_ = os.Setenv("MCP_PROM_ENABLED", "true")
	_ = os.Setenv("MCP_PROM_ADDRESS", "http://prometheus:9090")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Prometheus.Address != "http://prometheus:9090" {
		t.Errorf("Expected address 'http://prometheus:9090', got '%s'", cfg.Prometheus.Address)
	}
}

func TestGrafanaConfig(t *testing.T) {
	originalURL := os.Getenv("MCP_GRAFANA_URL")
	originalEnabled := os.Getenv("MCP_GRAFANA_ENABLED")

	defer func() {
		if originalURL != "" {
			_ = os.Setenv("MCP_GRAFANA_URL", originalURL)
		} else {
			_ = os.Unsetenv("MCP_GRAFANA_URL")
		}
		if originalEnabled != "" {
			_ = os.Setenv("MCP_GRAFANA_ENABLED", originalEnabled)
		} else {
			_ = os.Unsetenv("MCP_GRAFANA_ENABLED")
		}
	}()

	_ = os.Setenv("MCP_GRAFANA_ENABLED", "true")
	_ = os.Setenv("MCP_GRAFANA_URL", "http://grafana:3000")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Grafana.URL != "http://grafana:3000" {
		t.Errorf("Expected URL 'http://grafana:3000', got '%s'", cfg.Grafana.URL)
	}
}

func TestKibanaConfig(t *testing.T) {
	originalURL := os.Getenv("MCP_KIBANA_URL")
	originalEnabled := os.Getenv("MCP_KIBANA_ENABLED")

	defer func() {
		if originalURL != "" {
			_ = os.Setenv("MCP_KIBANA_URL", originalURL)
		} else {
			_ = os.Unsetenv("MCP_KIBANA_URL")
		}
		if originalEnabled != "" {
			_ = os.Setenv("MCP_KIBANA_ENABLED", originalEnabled)
		} else {
			_ = os.Unsetenv("MCP_KIBANA_ENABLED")
		}
	}()

	_ = os.Setenv("MCP_KIBANA_ENABLED", "true")
	_ = os.Setenv("MCP_KIBANA_URL", "http://kibana:5601")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Kibana.URL != "http://kibana:5601" {
		t.Errorf("Expected URL 'http://kibana:5601', got '%s'", cfg.Kibana.URL)
	}
}

func TestHelmConfig(t *testing.T) {
	originalNamespace := os.Getenv("MCP_HELM_NAMESPACE")
	originalDebug := os.Getenv("MCP_HELM_DEBUG")

	defer func() {
		if originalNamespace != "" {
			_ = os.Setenv("MCP_HELM_NAMESPACE", originalNamespace)
		} else {
			_ = os.Unsetenv("MCP_HELM_NAMESPACE")
		}
		if originalDebug != "" {
			_ = os.Setenv("MCP_HELM_DEBUG", originalDebug)
		} else {
			_ = os.Unsetenv("MCP_HELM_DEBUG")
		}
	}()

	_ = os.Setenv("MCP_HELM_NAMESPACE", "helm-system")
	_ = os.Setenv("MCP_HELM_DEBUG", "true")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Helm.Namespace != "helm-system" {
		t.Errorf("Expected namespace 'helm-system', got '%s'", cfg.Helm.Namespace)
	}

}

func TestOpenTelemetryServiceConfigFromEnv(t *testing.T) {
	t.Setenv("MCP_OPENTELEMETRY_ENABLED", "true")
	t.Setenv("MCP_OPENTELEMETRY_ADDRESS", "http://otel-collector:4318")
	t.Setenv("MCP_OPENTELEMETRY_TIMEOUT", "45")
	t.Setenv("MCP_OPENTELEMETRY_TLS_SKIP_VERIFY", "true")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if !cfg.OpenTelemetry.Enabled {
		t.Fatal("Expected OpenTelemetry service to be enabled")
	}
	if cfg.OpenTelemetry.Address != "http://otel-collector:4318" {
		t.Errorf("Expected OpenTelemetry address override, got %q", cfg.OpenTelemetry.Address)
	}
	if cfg.OpenTelemetry.TimeoutSec != 45 {
		t.Errorf("Expected OpenTelemetry timeout 45, got %d", cfg.OpenTelemetry.TimeoutSec)
	}
	if !cfg.OpenTelemetry.TLSSkipVerify {
		t.Fatal("Expected OpenTelemetry tlsSkipVerify to be true")
	}
}

func TestServerPathOverridesFromEnv(t *testing.T) {
	t.Setenv("MCP_SSE_PATH_ELASTICSEARCH", "/custom/elasticsearch/sse")
	t.Setenv("MCP_SSE_PATH_JAEGER", "/custom/jaeger/sse")
	t.Setenv("MCP_SSE_PATH_OPENTELEMETRY", "/custom/opentelemetry/sse")
	t.Setenv("MCP_STREAMABLE_HTTP_PATH_ELASTICSEARCH", "/custom/elasticsearch/stream")
	t.Setenv("MCP_STREAMABLE_HTTP_PATH_JAEGER", "/custom/jaeger/stream")
	t.Setenv("MCP_STREAMABLE_HTTP_PATH_OPENTELEMETRY", "/custom/opentelemetry/stream")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Server.SSEPaths.Elasticsearch != "/custom/elasticsearch/sse" {
		t.Errorf("Expected elasticsearch SSE path override, got %q", cfg.Server.SSEPaths.Elasticsearch)
	}
	if cfg.Server.SSEPaths.Jaeger != "/custom/jaeger/sse" {
		t.Errorf("Expected jaeger SSE path override, got %q", cfg.Server.SSEPaths.Jaeger)
	}
	if cfg.Server.SSEPaths.OpenTelemetry != "/custom/opentelemetry/sse" {
		t.Errorf("Expected opentelemetry SSE path override, got %q", cfg.Server.SSEPaths.OpenTelemetry)
	}
	if cfg.Server.StreamableHTTPPaths.Elasticsearch != "/custom/elasticsearch/stream" {
		t.Errorf("Expected elasticsearch streamable-http path override, got %q", cfg.Server.StreamableHTTPPaths.Elasticsearch)
	}
	if cfg.Server.StreamableHTTPPaths.Jaeger != "/custom/jaeger/stream" {
		t.Errorf("Expected jaeger streamable-http path override, got %q", cfg.Server.StreamableHTTPPaths.Jaeger)
	}
	if cfg.Server.StreamableHTTPPaths.OpenTelemetry != "/custom/opentelemetry/stream" {
		t.Errorf("Expected opentelemetry streamable-http path override, got %q", cfg.Server.StreamableHTTPPaths.OpenTelemetry)
	}
}

func TestServerOTELConfigFromEnv(t *testing.T) {
	t.Setenv("MCP_OTEL_ENABLED", "true")
	t.Setenv("MCP_OTEL_SERVICE_NAME", "mcp-server")
	t.Setenv("MCP_OTEL_TRACING_ENABLED", "true")
	t.Setenv("MCP_OTEL_TRACING_SAMPLE_RATE", "0.75")
	t.Setenv("MCP_OTEL_METRICS_TEMPORALITY", "delta")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if !cfg.OTEL.Enabled {
		t.Fatal("Expected OTEL to be enabled")
	}
	if cfg.OTEL.ServiceName != "mcp-server" {
		t.Errorf("Expected OTEL service name override, got %q", cfg.OTEL.ServiceName)
	}
	if !cfg.OTEL.Tracing.Enabled {
		t.Fatal("Expected OTEL tracing to be enabled")
	}
	if cfg.OTEL.Tracing.SampleRate != 0.75 {
		t.Errorf("Expected OTEL tracing sample rate 0.75, got %f", cfg.OTEL.Tracing.SampleRate)
	}
	if cfg.OTEL.Metrics.Temporality != "delta" {
		t.Errorf("Expected OTEL metrics temporality override, got %q", cfg.OTEL.Metrics.Temporality)
	}
}

func TestRateLimitConfigFromEnv(t *testing.T) {
	t.Setenv("MCP_RATELIMIT_ENABLED", "true")
	t.Setenv("MCP_RATELIMIT_REQUESTS_PER_SECOND", "25.5")
	t.Setenv("MCP_RATELIMIT_BURST", "80")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if !cfg.RateLimit.Enabled {
		t.Fatal("Expected rate limit to be enabled")
	}
	if cfg.RateLimit.RequestsPerSecond != 25.5 {
		t.Errorf("Expected rate limit requests_per_second 25.5, got %f", cfg.RateLimit.RequestsPerSecond)
	}
	if cfg.RateLimit.Burst != 80 {
		t.Errorf("Expected rate limit burst 80, got %d", cfg.RateLimit.Burst)
	}
}

func TestRateLimitAliasConfigFromEnv(t *testing.T) {
	t.Setenv("MCP_RATE_LIMIT_ENABLED", "true")
	t.Setenv("MCP_RATE_LIMIT_REQUESTS_PER_SECOND", "12")
	t.Setenv("MCP_RATE_LIMIT_BURST", "24")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if !cfg.RateLimit.Enabled {
		t.Fatal("Expected rate limit to be enabled via alias env vars")
	}
	if cfg.RateLimit.RequestsPerSecond != 12 {
		t.Errorf("Expected rate limit requests_per_second 12, got %f", cfg.RateLimit.RequestsPerSecond)
	}
	if cfg.RateLimit.Burst != 24 {
		t.Errorf("Expected rate limit burst 24, got %d", cfg.RateLimit.Burst)
	}
}

func TestRateLimitInvalidConfigFromEnv(t *testing.T) {
	t.Setenv("MCP_RATELIMIT_ENABLED", "true")
	t.Setenv("MCP_RATELIMIT_REQUESTS_PER_SECOND", "-1")
	t.Setenv("MCP_RATELIMIT_BURST", "10")

	_, err := Load("")
	if err == nil {
		t.Fatal("Expected error for invalid negative ratelimit requests_per_second")
	}
}

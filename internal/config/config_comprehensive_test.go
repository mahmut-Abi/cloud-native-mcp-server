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

package main

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/mahmut-Abi/k8s-mcp-server/internal/config"
)

func TestParseFlags(t *testing.T) {
	// Save original args and flag state
	oldArgs := os.Args
	oldCommandLine := flag.CommandLine
	defer func() {
		os.Args = oldArgs
		flag.CommandLine = oldCommandLine
	}()

	// Reset flag state for this test
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Test default values
	os.Args = []string{"cmd"}
	config := parseFlags()

	if config.Addr != "0.0.0.0:8080" {
		t.Errorf("Expected default addr '0.0.0.0:8080', got '%s'", config.Addr)
	}

	if config.LogLevel != "info" {
		t.Errorf("Expected default log level 'info', got '%s'", config.LogLevel)
	}

	if config.ReadTimeout != 0 {
		t.Errorf("Expected default read timeout 0, got %v", config.ReadTimeout)
	}

	if config.WriteTimeout != 0 {
		t.Errorf("Expected default write timeout 0, got %v", config.WriteTimeout)
	}

	if config.IdleTimeout != 60*time.Second {
		t.Errorf("Expected default idle timeout 60s, got %v", config.IdleTimeout)
	}
}

func TestParseFlagsWithCustomValues(t *testing.T) {
	// Save original args and flag state
	oldArgs := os.Args
	oldCommandLine := flag.CommandLine
	defer func() {
		os.Args = oldArgs
		flag.CommandLine = oldCommandLine
	}()

	// Reset flag state for this test
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Test custom values
	os.Args = []string{
		"cmd",
		"-addr", "127.0.0.1:9090",
		"-log-level", "debug",
		"-read-timeout", "30",
		"-write-timeout", "60",
		"-idle-timeout", "120",
		"-mode", "http",
	}

	config := parseFlags()

	if config.Addr != "127.0.0.1:9090" {
		t.Errorf("Expected addr '127.0.0.1:9090', got '%s'", config.Addr)
	}

	if config.LogLevel != "debug" {
		t.Errorf("Expected log level 'debug', got '%s'", config.LogLevel)
	}

	if config.ReadTimeout != 30*time.Second {
		t.Errorf("Expected read timeout 30s, got %v", config.ReadTimeout)
	}

	if config.WriteTimeout != 60*time.Second {
		t.Errorf("Expected write timeout 60s, got %v", config.WriteTimeout)
	}

	if config.IdleTimeout != 120*time.Second {
		t.Errorf("Expected idle timeout 120s, got %v", config.IdleTimeout)
	}

	if config.Mode != "http" {
		t.Errorf("Expected mode 'http', got '%s'", config.Mode)
	}
}

func TestApplyAppConfigWithNilConfig(t *testing.T) {
	cliConfig := &CLIConfig{
		Addr:     "0.0.0.0:8080",
		LogLevel: "info",
		Mode:     "",
	}

	applyAppConfig(cliConfig, nil)

	// Values should remain unchanged
	if cliConfig.Addr != "0.0.0.0:8080" {
		t.Errorf("Addr should remain unchanged with nil config, got '%s'", cliConfig.Addr)
	}

	if cliConfig.LogLevel != "info" {
		t.Errorf("LogLevel should remain unchanged with nil config, got '%s'", cliConfig.LogLevel)
	}
}

func TestApplyAppConfigWithValidConfig(t *testing.T) {
	cliConfig := &CLIConfig{
		Addr:         "0.0.0.0:8080",   // Default value
		LogLevel:     "info",           // Default value
		Mode:         "",               // Empty value
		ReadTimeout:  0,                // Default value
		WriteTimeout: 0,                // Default value
		IdleTimeout:  60 * time.Second, // Default value
	}

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
			Mode:            "sse",
			Addr:            "127.0.0.1:9090",
			ReadTimeoutSec:  30,
			WriteTimeoutSec: 60,
			IdleTimeoutSec:  120,
		},
		Logging: struct {
			Level string `yaml:"level"`
			JSON  bool   `yaml:"json"`
		}{
			Level: "debug",
			JSON:  true,
		},
	}

	applyAppConfig(cliConfig, appConfig)

	// CLI config should be updated with values from app config
	if cliConfig.Addr != "127.0.0.1:9090" {
		t.Errorf("Expected addr to be updated to '127.0.0.1:9090', got '%s'", cliConfig.Addr)
	}

	if cliConfig.LogLevel != "debug" {
		t.Errorf("Expected log level to be updated to 'debug', got '%s'", cliConfig.LogLevel)
	}

	if cliConfig.Mode != "sse" {
		t.Errorf("Expected mode to be updated to 'sse', got '%s'", cliConfig.Mode)
	}

	if cliConfig.ReadTimeout != 30*time.Second {
		t.Errorf("Expected read timeout to be updated to 30s, got %v", cliConfig.ReadTimeout)
	}

	if cliConfig.WriteTimeout != 60*time.Second {
		t.Errorf("Expected write timeout to be updated to 60s, got %v", cliConfig.WriteTimeout)
	}

	if cliConfig.IdleTimeout != 120*time.Second {
		t.Errorf("Expected idle timeout to be updated to 120s, got %v", cliConfig.IdleTimeout)
	}
}

func TestApplyAppConfigCLITakesPrecedence(t *testing.T) {
	// CLI config with non-default values
	cliConfig := &CLIConfig{
		Addr:        "custom-addr:8080",
		LogLevel:    "warn",
		Mode:        "stdio",
		addrSet:     true,
		logLevelSet: true,
		modeSet:     true,
	}

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
			Mode: "sse",
			Addr: "app-config-addr:9090",
		},
		Logging: struct {
			Level string `yaml:"level"`
			JSON  bool   `yaml:"json"`
		}{
			Level: "debug",
		},
	}

	applyAppConfig(cliConfig, appConfig)

	// CLI values should take precedence over app config
	if cliConfig.Addr != "custom-addr:8080" {
		t.Errorf("CLI addr should take precedence, expected 'custom-addr:8080', got '%s'", cliConfig.Addr)
	}

	if cliConfig.LogLevel != "warn" {
		t.Errorf("CLI log level should take precedence, expected 'warn', got '%s'", cliConfig.LogLevel)
	}

	if cliConfig.Mode != "stdio" {
		t.Errorf("CLI mode should take precedence, expected 'stdio', got '%s'", cliConfig.Mode)
	}
}

func TestGetDefaultKubeconfig(t *testing.T) {
	// Save original env
	originalKubeconfig := os.Getenv("KUBECONFIG")
	originalHome := os.Getenv("HOME")
	defer func() {
		_ = os.Setenv("KUBECONFIG", originalKubeconfig)
		_ = os.Setenv("HOME", originalHome)
	}()

	// Test with KUBECONFIG env var set
	_ = os.Setenv("KUBECONFIG", "/custom/kubeconfig")
	_ = os.Setenv("HOME", "/home/user")

	kubeconfig := getDefaultKubeconfig()
	if kubeconfig != "/custom/kubeconfig" {
		t.Errorf("Expected '/custom/kubeconfig', got '%s'", kubeconfig)
	}

	// Test without KUBECONFIG env var
	_ = os.Unsetenv("KUBECONFIG")
	kubeconfig = getDefaultKubeconfig()
	expected := "/home/user/.kube/config"
	if kubeconfig != expected {
		t.Errorf("Expected '%s', got '%s'", expected, kubeconfig)
	}
}

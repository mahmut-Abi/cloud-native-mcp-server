package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_EmptyPath(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if cfg == nil {
		t.Fatal("Expected config to be non-nil")
	}
}

func TestLoad_ValidYAML(t *testing.T) {
	// Create temp file with valid YAML
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")

	yamlContent := `
server:
  mode: "http"
  addr: "localhost:8080"
  readTimeoutSec: 30
  writeTimeoutSec: 30
  idleTimeoutSec: 120

logging:
  level: "debug"
  json: true

kubernetes:
  kubeconfig: "/path/to/kubeconfig"
  timeoutSec: 60
  qps: 20.5
  burst: 30
`

	err := os.WriteFile(configFile, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp config file: %v", err)
	}

	cfg, err := Load(configFile)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify config values
	if cfg.Server.Mode != "http" {
		t.Errorf("Expected mode 'http', got '%s'", cfg.Server.Mode)
	}
	if cfg.Server.Addr != "localhost:8080" {
		t.Errorf("Expected addr 'localhost:8080', got '%s'", cfg.Server.Addr)
	}
	if cfg.Logging.Level != "debug" {
		t.Errorf("Expected log level 'debug', got '%s'", cfg.Logging.Level)
	}
	if !cfg.Logging.JSON {
		t.Error("Expected JSON logging to be true")
	}
	if cfg.Kubernetes.QPS != 20.5 {
		t.Errorf("Expected QPS 20.5, got %f", cfg.Kubernetes.QPS)
	}
}

func TestLoad_InvalidPath(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("Expected error for nonexistent file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "invalid.yaml")

	// Write invalid YAML
	err := os.WriteFile(configFile, []byte("invalid: yaml: content: ["), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp config file: %v", err)
	}

	_, err = Load(configFile)
	if err == nil {
		t.Fatal("Expected error for invalid YAML")
	}
}

func TestEnvironmentOverrides(t *testing.T) {
	// Set environment variables
	_ = os.Setenv("MCP_MODE", "sse")
	_ = os.Setenv("MCP_ADDR", "0.0.0.0:9000")
	_ = os.Setenv("MCP_LOG_LEVEL", "error")
	_ = os.Setenv("MCP_LOG_JSON", "true")
	_ = os.Setenv("MCP_K8S_QPS", "15.5")
	defer func() {
		_ = os.Unsetenv("MCP_MODE")
		_ = os.Unsetenv("MCP_ADDR")
		_ = os.Unsetenv("MCP_LOG_LEVEL")
		_ = os.Unsetenv("MCP_LOG_JSON")
		_ = os.Unsetenv("MCP_K8S_QPS")
	}()

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.Server.Mode != "sse" {
		t.Errorf("Expected mode 'sse', got '%s'", cfg.Server.Mode)
	}
	if cfg.Server.Addr != "0.0.0.0:9000" {
		t.Errorf("Expected addr '0.0.0.0:9000', got '%s'", cfg.Server.Addr)
	}
	if cfg.Logging.Level != "error" {
		t.Errorf("Expected log level 'error', got '%s'", cfg.Logging.Level)
	}
	if !cfg.Logging.JSON {
		t.Error("Expected JSON logging to be true")
	}
	if cfg.Kubernetes.QPS != 15.5 {
		t.Errorf("Expected QPS 15.5, got %f", cfg.Kubernetes.QPS)
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test atoiDefault
	if atoiDefault("", 10) != 10 {
		t.Error("atoiDefault should return default for empty string")
	}
	if atoiDefault("42", 10) != 42 {
		t.Error("atoiDefault should parse valid integer")
	}
	if atoiDefault("invalid", 10) != 10 {
		t.Error("atoiDefault should return default for invalid integer")
	}

	// Test atofDefault
	if atofDefault("", 5.5) != 5.5 {
		t.Error("atofDefault should return default for empty string")
	}
	if atofDefault("3.14", 5.5) != 3.14 {
		t.Error("atofDefault should parse valid float")
	}
	if atofDefault("invalid", 5.5) != 5.5 {
		t.Error("atofDefault should return default for invalid float")
	}

	// Test isTrue
	trueCases := []string{"1", "true", "TRUE", "yes", "YES", "on", "ON", " true "}
	for _, tc := range trueCases {
		if !isTrue(tc) {
			t.Errorf("isTrue(%q) should return true", tc)
		}
	}

	falseCases := []string{"0", "false", "no", "off", "", "invalid"}
	for _, fc := range falseCases {
		if isTrue(fc) {
			t.Errorf("isTrue(%q) should return false", fc)
		}
	}
}

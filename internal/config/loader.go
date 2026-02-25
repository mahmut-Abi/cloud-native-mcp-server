package config

import (
	"fmt"
	"os"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/constants"
	"gopkg.in/yaml.v3"
)

// ConfigLoader handles loading configuration from files
type ConfigLoader struct {
	validator *ConfigValidator
	parser    *EnvParser
}

// NewConfigLoader creates a new configuration loader
func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{
		validator: NewConfigValidator(),
		parser:    NewEnvParser(),
	}
}

// Load loads configuration from YAML file (if provided) and merges environment overrides.
// It also validates the configuration before returning it.
func (l *ConfigLoader) Load(path string) (*AppConfig, error) {
	cfg := &AppConfig{}

	// Load from file if path is provided
	if path != "" {
		if err := l.loadFromFile(path, cfg); err != nil {
			return nil, err
		}
	}

	// Apply environment variable overrides
	l.parser.Parse(cfg)

	// Set default values
	l.setDefaults(cfg)

	// Validate configuration
	if err := l.validator.Validate(cfg); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

// setDefaults sets reasonable default values for configuration
func (l *ConfigLoader) setDefaults(cfg *AppConfig) {
	// Server defaults
	if cfg.Server.Mode == "" {
		cfg.Server.Mode = "sse"
	}
	if cfg.Server.Addr == "" {
		cfg.Server.Addr = ":8080"
	}
	if cfg.Server.ReadTimeoutSec == 0 {
		cfg.Server.ReadTimeoutSec = 30
	}
	if cfg.Server.WriteTimeoutSec == 0 {
		cfg.Server.WriteTimeoutSec = 30
	}
	if cfg.Server.IdleTimeoutSec == 0 {
		cfg.Server.IdleTimeoutSec = 60
	}

	// Logging defaults
	if cfg.Logging.Level == "" {
		cfg.Logging.Level = "info"
	}

	// Rate limit defaults
	if cfg.RateLimit.RequestsPerSecond == 0 {
		cfg.RateLimit.RequestsPerSecond = constants.DefaultRateLimitRPS
	}
	if cfg.RateLimit.Burst == 0 {
		cfg.RateLimit.Burst = constants.DefaultRateLimitBurst
	}

	// Kubernetes defaults
	if cfg.Kubernetes.TimeoutSec == 0 {
		cfg.Kubernetes.TimeoutSec = 30
	}
	if cfg.Kubernetes.QPS == 0 {
		cfg.Kubernetes.QPS = 50
	}
	if cfg.Kubernetes.Burst == 0 {
		cfg.Kubernetes.Burst = 100
	}

	// Prometheus defaults
	if cfg.Prometheus.TimeoutSec == 0 {
		cfg.Prometheus.TimeoutSec = 30
	}

	// Grafana defaults
	if cfg.Grafana.TimeoutSec == 0 {
		cfg.Grafana.TimeoutSec = 30
	}

	// Kibana defaults
	if cfg.Kibana.TimeoutSec == 0 {
		cfg.Kibana.TimeoutSec = 30
	}
	if cfg.Kibana.Space == "" {
		cfg.Kibana.Space = "default"
	}

	// Helm defaults
	if cfg.Helm.TimeoutSec == 0 {
		cfg.Helm.TimeoutSec = 300
	}
	if cfg.Helm.MaxRetries == 0 {
		cfg.Helm.MaxRetries = 3
	}

	// Alertmanager defaults
	if cfg.Alertmanager.TimeoutSec == 0 {
		cfg.Alertmanager.TimeoutSec = 30
	}

	// Jaeger defaults
	if cfg.Jaeger.TimeoutSec == 0 {
		cfg.Jaeger.TimeoutSec = 30
	}

	// Elasticsearch defaults
	if cfg.Elasticsearch.TimeoutSec == 0 {
		cfg.Elasticsearch.TimeoutSec = 30
	}

	// Audit defaults
	if cfg.Audit.Level == "" {
		cfg.Audit.Level = "info"
	}
	if cfg.Audit.MaxLogs == 0 {
		cfg.Audit.MaxLogs = 10000
	}
	if cfg.Audit.Storage == "" {
		cfg.Audit.Storage = "memory"
	}
	if cfg.Audit.Format == "" {
		cfg.Audit.Format = "json"
	}
	if cfg.Audit.Query.MaxResults == 0 {
		cfg.Audit.Query.MaxResults = 100
	}
	if cfg.Audit.Query.TimeRange == 0 {
		cfg.Audit.Query.TimeRange = 7
	}
	if cfg.Audit.Alerts.FailureThreshold == 0 {
		cfg.Audit.Alerts.FailureThreshold = 5
	}
	if cfg.Audit.Alerts.CheckIntervalSec == 0 {
		cfg.Audit.Alerts.CheckIntervalSec = 300
	}
	if cfg.Audit.Alerts.Method == "" {
		cfg.Audit.Alerts.Method = "none"
	}
	if cfg.Audit.Masking.MaskValue == "" {
		cfg.Audit.Masking.MaskValue = "***"
	}
	if cfg.Audit.Sampling.Rate == 0 {
		cfg.Audit.Sampling.Rate = 1.0
	}
}

// loadFromFile loads configuration from a YAML file
func (l *ConfigLoader) loadFromFile(path string, cfg *AppConfig) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config file failed: %w", err)
	}

	if err := yaml.Unmarshal(b, cfg); err != nil {
		return fmt.Errorf("yaml unmarshal failed: %w", err)
	}

	return nil
}

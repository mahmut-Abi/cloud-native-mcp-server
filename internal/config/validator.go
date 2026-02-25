package config

import (
	"fmt"
	"net"
	"regexp"
)

// ConfigValidator validates application configuration
type ConfigValidator struct{}

// NewConfigValidator creates a new configuration validator
func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{}
}

// Validate performs comprehensive validation of the configuration
func (v *ConfigValidator) Validate(cfg *AppConfig) error {
	if err := v.validateServerConfig(cfg); err != nil {
		return fmt.Errorf("server config validation failed: %w", err)
	}

	if err := v.validateLoggingConfig(cfg); err != nil {
		return fmt.Errorf("logging config validation failed: %w", err)
	}

	if err := v.validateKubernetesConfig(cfg); err != nil {
		return fmt.Errorf("kubernetes config validation failed: %w", err)
	}

	if err := v.validatePrometheusConfig(cfg); err != nil {
		return fmt.Errorf("prometheus config validation failed: %w", err)
	}

	if err := v.validateGrafanaConfig(cfg); err != nil {
		return fmt.Errorf("grafana config validation failed: %w", err)
	}

	if err := v.validateKibanaConfig(cfg); err != nil {
		return fmt.Errorf("kibana config validation failed: %w", err)
	}

	if err := v.validateHelmConfig(cfg); err != nil {
		return fmt.Errorf("helm config validation failed: %w", err)
	}

	if err := v.validateElasticsearchConfig(cfg); err != nil {
		return fmt.Errorf("elasticsearch config validation failed: %w", err)
	}

	if err := v.validateAlertmanagerConfig(cfg); err != nil {
		return fmt.Errorf("alertmanager config validation failed: %w", err)
	}

	if err := v.validateJaegerConfig(cfg); err != nil {
		return fmt.Errorf("jaeger config validation failed: %w", err)
	}

	if err := v.validateAuditConfig(cfg); err != nil {
		return fmt.Errorf("audit config validation failed: %w", err)
	}

	if err := v.validateAuthConfig(cfg); err != nil {
		return fmt.Errorf("auth config validation failed: %w", err)
	}

	if err := v.validateRateLimitConfig(cfg); err != nil {
		return fmt.Errorf("ratelimit config validation failed: %w", err)
	}

	return nil
}

func (v *ConfigValidator) validateServerConfig(cfg *AppConfig) error {
	validModes := map[string]bool{
		"sse":             true,
		"http":            true,
		"streamable-http": true,
		"stdio":           true,
	}

	// Set default mode if not specified
	if cfg.Server.Mode == "" {
		cfg.Server.Mode = "sse"
	}

	if !validModes[cfg.Server.Mode] {
		return fmt.Errorf("invalid server mode: %s, must be one of: sse, http, streamable-http, stdio", cfg.Server.Mode)
	}

	// Set default address if not specified and mode is not stdio
	if cfg.Server.Mode != "stdio" && cfg.Server.Addr == "" {
		cfg.Server.Addr = "0.0.0.0:8080"
	}

	// Validate timeout values are non-negative and within reasonable limits
	if cfg.Server.ReadTimeoutSec < 0 {
		return fmt.Errorf("read timeout must be non-negative")
	}
	if cfg.Server.ReadTimeoutSec > 3600 {
		return fmt.Errorf("read timeout too large: %d (max: 3600 seconds)", cfg.Server.ReadTimeoutSec)
	}
	if cfg.Server.WriteTimeoutSec < 0 {
		return fmt.Errorf("write timeout must be non-negative")
	}
	if cfg.Server.WriteTimeoutSec > 3600 {
		return fmt.Errorf("write timeout too large: %d (max: 3600 seconds)", cfg.Server.WriteTimeoutSec)
	}
	if cfg.Server.IdleTimeoutSec < 0 {
		return fmt.Errorf("idle timeout must be non-negative")
	}
	if cfg.Server.IdleTimeoutSec > 3600 {
		return fmt.Errorf("idle timeout too large: %d (max: 3600 seconds)", cfg.Server.IdleTimeoutSec)
	}

	// Validate address format if specified
	if cfg.Server.Addr != "" {
		_, _, err := net.SplitHostPort(cfg.Server.Addr)
		if err != nil {
			return fmt.Errorf("invalid address format '%s': %w", cfg.Server.Addr, err)
		}
	}

	return nil
}

func (v *ConfigValidator) validateLoggingConfig(cfg *AppConfig) error {
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	// Set default log level if not specified
	if cfg.Logging.Level == "" {
		cfg.Logging.Level = "info"
	}

	if !validLevels[cfg.Logging.Level] {
		return fmt.Errorf("invalid log level: %s, must be one of: debug, info, warn, error", cfg.Logging.Level)
	}

	return nil
}

func (v *ConfigValidator) validateKubernetesConfig(cfg *AppConfig) error {
	if cfg.Kubernetes.TimeoutSec < 0 {
		return fmt.Errorf("kubernetes timeout must be non-negative")
	}

	if cfg.Kubernetes.QPS < 0 {
		return fmt.Errorf("kubernetes QPS must be non-negative")
	}

	if cfg.Kubernetes.Burst < 0 {
		return fmt.Errorf("kubernetes burst must be non-negative")
	}

	return nil
}

func (v *ConfigValidator) validatePrometheusConfig(cfg *AppConfig) error {
	if !cfg.Prometheus.Enabled {
		return nil
	}

	if cfg.Prometheus.Address == "" {
		return fmt.Errorf("prometheus address is required when enabled")
	}

	// Set default timeout if not specified
	if cfg.Prometheus.TimeoutSec <= 0 {
		cfg.Prometheus.TimeoutSec = 30
	}

	// Validate TLS configuration
	if cfg.Prometheus.TLSSkipVerify && (cfg.Prometheus.TLSCertFile != "" || cfg.Prometheus.TLSKeyFile != "" || cfg.Prometheus.TLSCAFile != "") {
		return fmt.Errorf("prometheus TLS skip verify cannot be enabled with certificate files")
	}

	return nil
}

func (v *ConfigValidator) validateGrafanaConfig(cfg *AppConfig) error {
	if !cfg.Grafana.Enabled {
		return nil
	}

	if cfg.Grafana.URL == "" {
		return fmt.Errorf("grafana URL is required when enabled")
	}

	// Validate URL format
	if !isValidURL(cfg.Grafana.URL) {
		return fmt.Errorf("invalid grafana URL format: %s", cfg.Grafana.URL)
	}

	// Set default timeout if not specified
	if cfg.Grafana.TimeoutSec <= 0 {
		cfg.Grafana.TimeoutSec = 30
	}

	// Note: Authentication is optional for Grafana - it can be configured later
	// or accessed without authentication if the Grafana instance allows it

	return nil
}

func (v *ConfigValidator) validateKibanaConfig(cfg *AppConfig) error {
	if !cfg.Kibana.Enabled {
		return nil
	}

	if cfg.Kibana.URL == "" {
		return fmt.Errorf("kibana URL is required when enabled")
	}

	// Validate URL format
	if !isValidURL(cfg.Kibana.URL) {
		return fmt.Errorf("invalid kibana URL format: %s", cfg.Kibana.URL)
	}

	// Set default timeout if not specified
	if cfg.Kibana.TimeoutSec <= 0 {
		cfg.Kibana.TimeoutSec = 30
	}

	// Note: Authentication is optional for Kibana - it can be configured later
	// or accessed without authentication if the Kibana instance allows it

	return nil
}

func (v *ConfigValidator) validateHelmConfig(cfg *AppConfig) error {
	if !cfg.Helm.Enabled {
		return nil
	}

	if cfg.Helm.TimeoutSec <= 0 {
		return fmt.Errorf("helm timeout must be positive")
	}

	if cfg.Helm.MaxRetries < 0 {
		return fmt.Errorf("helm max retries must be non-negative")
	}

	return nil
}

func (v *ConfigValidator) validateElasticsearchConfig(cfg *AppConfig) error {
	if !cfg.Elasticsearch.Enabled {
		return nil
	}

	if len(cfg.Elasticsearch.Addresses) == 0 && cfg.Elasticsearch.Address == "" {
		return fmt.Errorf("elasticsearch addresses are required when enabled")
	}

	if cfg.Elasticsearch.TimeoutSec <= 0 {
		return fmt.Errorf("elasticsearch timeout must be positive")
	}

	// Validate TLS configuration
	if cfg.Elasticsearch.TLSSkipVerify && (cfg.Elasticsearch.TLSCertFile != "" || cfg.Elasticsearch.TLSKeyFile != "" || cfg.Elasticsearch.TLSCAFile != "") {
		return fmt.Errorf("elasticsearch TLS skip verify cannot be enabled with certificate files")
	}

	return nil
}

func (v *ConfigValidator) validateAlertmanagerConfig(cfg *AppConfig) error {
	if !cfg.Alertmanager.Enabled {
		return nil
	}

	if cfg.Alertmanager.Address == "" {
		return fmt.Errorf("alertmanager address is required when enabled")
	}

	if cfg.Alertmanager.TimeoutSec <= 0 {
		return fmt.Errorf("alertmanager timeout must be positive")
	}

	// Validate TLS configuration
	if cfg.Alertmanager.TLSSkipVerify && (cfg.Alertmanager.TLSCertFile != "" || cfg.Alertmanager.TLSKeyFile != "" || cfg.Alertmanager.TLSCAFile != "") {
		return fmt.Errorf("alertmanager TLS skip verify cannot be enabled with certificate files")
	}

	return nil
}

func (v *ConfigValidator) validateJaegerConfig(cfg *AppConfig) error {
	if !cfg.Jaeger.Enabled {
		return nil
	}

	if cfg.Jaeger.Address == "" {
		return fmt.Errorf("jaeger address is required when enabled")
	}

	if cfg.Jaeger.TimeoutSec <= 0 {
		return fmt.Errorf("jaeger timeout must be positive")
	}

	return nil
}

func (v *ConfigValidator) validateAuditConfig(cfg *AppConfig) error {
	if !cfg.Audit.Enabled {
		return nil
	}

	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if cfg.Audit.Level != "" && !validLevels[cfg.Audit.Level] {
		return fmt.Errorf("invalid audit log level: %s, must be one of: debug, info, warn, error", cfg.Audit.Level)
	}

	validStorage := map[string]bool{
		"memory":   true,
		"file":     true,
		"database": true,
		"all":      true,
	}

	if cfg.Audit.Storage != "" && !validStorage[cfg.Audit.Storage] {
		return fmt.Errorf("invalid audit storage type: %s, must be one of: memory, file, database, all", cfg.Audit.Storage)
	}

	validFormat := map[string]bool{
		"json": true,
		"text": true,
	}

	if cfg.Audit.Format != "" && !validFormat[cfg.Audit.Format] {
		return fmt.Errorf("invalid audit format: %s, must be one of: json, text", cfg.Audit.Format)
	}

	// Validate sampling rate
	if cfg.Audit.Sampling.Enabled {
		if cfg.Audit.Sampling.Rate < 0 || cfg.Audit.Sampling.Rate > 1 {
			return fmt.Errorf("audit sampling rate must be between 0 and 1")
		}
	}

	return nil
}

func (v *ConfigValidator) validateAuthConfig(cfg *AppConfig) error {
	if !cfg.Auth.Enabled {
		return nil
	}

	validModes := map[string]bool{
		"apikey": true,
		"bearer": true,
		"basic":  true,
	}

	if cfg.Auth.Mode == "" {
		return fmt.Errorf("auth mode is required when authentication is enabled")
	}

	if !validModes[cfg.Auth.Mode] {
		return fmt.Errorf("invalid auth mode: %s, must be one of: apikey, bearer, basic", cfg.Auth.Mode)
	}

	// Validate authentication credentials based on mode
	switch cfg.Auth.Mode {
	case "apikey":
		if cfg.Auth.APIKey == "" {
			return fmt.Errorf("API key is required for apikey auth mode")
		}
	case "bearer":
		if cfg.Auth.BearerToken == "" {
			return fmt.Errorf("bearer token is required for bearer auth mode")
		}
	case "basic":
		if cfg.Auth.Username == "" || cfg.Auth.Password == "" {
			return fmt.Errorf("username and password are required for basic auth mode")
		}
	}

	// Validate JWT configuration
	if cfg.Auth.JWTSecret != "" && cfg.Auth.JWTAlgorithm == "" {
		return fmt.Errorf("JWT algorithm is required when JWT secret is provided")
	}

	validJWTAlgorithms := map[string]bool{
		"HS256": true,
		"HS384": true,
		"HS512": true,
		"RS256": true,
		"RS384": true,
		"RS512": true,
		"ES256": true,
		"ES384": true,
		"ES512": true,
	}

	if cfg.Auth.JWTAlgorithm != "" && !validJWTAlgorithms[cfg.Auth.JWTAlgorithm] {
		return fmt.Errorf("invalid JWT algorithm: %s", cfg.Auth.JWTAlgorithm)
	}

	return nil
}

func (v *ConfigValidator) validateRateLimitConfig(cfg *AppConfig) error {
	if cfg.RateLimit.RequestsPerSecond < 0 {
		return fmt.Errorf("ratelimit requests_per_second must be non-negative")
	}

	if cfg.RateLimit.Burst < 0 {
		return fmt.Errorf("ratelimit burst must be non-negative")
	}

	if !cfg.RateLimit.Enabled {
		return nil
	}

	if cfg.RateLimit.RequestsPerSecond <= 0 {
		return fmt.Errorf("ratelimit requests_per_second must be greater than 0 when enabled")
	}

	if cfg.RateLimit.Burst <= 0 {
		return fmt.Errorf("ratelimit burst must be greater than 0 when enabled")
	}

	return nil
}

// isValidURL checks if a string is a valid URL format
func isValidURL(url string) bool {
	// Simple URL validation regex
	regex := regexp.MustCompile(`^https?://[a-zA-Z0-9\-._~:/?#[\]@!$&'()*+,;=]+$`)
	return regex.MatchString(url)
}

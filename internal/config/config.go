package config

import "fmt"

// AppConfig represents application configuration loaded from YAML and environment variables.
type AppConfig struct {
	Server struct {
		Mode            string `yaml:"mode"` // sse | http | streamable-http | stdio
		Addr            string `yaml:"addr"`
		ReadTimeoutSec  int    `yaml:"readTimeoutSec"`  // 0 disables
		WriteTimeoutSec int    `yaml:"writeTimeoutSec"` // 0 disables
		IdleTimeoutSec  int    `yaml:"idleTimeoutSec"`  // default 60
		SSEPaths        struct {
			Kubernetes    string `yaml:"kubernetes"`    // SSE path for Kubernetes service
			Grafana       string `yaml:"grafana"`       // SSE path for Grafana service
			Prometheus    string `yaml:"prometheus"`    // SSE path for Prometheus service
			Kibana        string `yaml:"kibana"`        // SSE path for Kibana service
			Helm          string `yaml:"helm"`          // SSE path for Helm service
			Elasticsearch string `yaml:"elasticsearch"` // SSE path for Elasticsearch service
			Alertmanager  string `yaml:"alertmanager"`  // SSE path for Alertmanager service
			Jaeger        string `yaml:"jaeger"`        // SSE path for Jaeger service
			OpenTelemetry string `yaml:"opentelemetry"` // SSE path for OpenTelemetry service
			Aggregate     string `yaml:"aggregate"`     // SSE path for aggregated service
			Utilities     string `yaml:"utilities"`     // SSE path for Utilities service
		} `yaml:"ssePaths"`
		StreamableHTTPPaths struct {
			Kubernetes    string `yaml:"kubernetes"`    // Streamable-HTTP path for Kubernetes service
			Grafana       string `yaml:"grafana"`       // Streamable-HTTP path for Grafana service
			Prometheus    string `yaml:"prometheus"`    // Streamable-HTTP path for Prometheus service
			Kibana        string `yaml:"kibana"`        // Streamable-HTTP path for Kibana service
			Helm          string `yaml:"helm"`          // Streamable-HTTP path for Helm service
			Elasticsearch string `yaml:"elasticsearch"` // Streamable-HTTP path for Elasticsearch service
			Alertmanager  string `yaml:"alertmanager"`  // Streamable-HTTP path for Alertmanager service
			Jaeger        string `yaml:"jaeger"`        // Streamable-HTTP path for Jaeger service
			OpenTelemetry string `yaml:"opentelemetry"` // Streamable-HTTP path for OpenTelemetry service
			Aggregate     string `yaml:"aggregate"`     // Streamable-HTTP path for aggregated service
			Utilities     string `yaml:"utilities"`     // Streamable-HTTP path for Utilities service
		} `yaml:"streamableHttpPaths"`
		CORS struct {
			AllowedOrigins []string `yaml:"allowedOrigins"` // List of allowed CORS origins
			AllowedMethods []string `yaml:"allowedMethods"` // List of allowed CORS methods
			AllowedHeaders []string `yaml:"allowedHeaders"` // List of allowed CORS headers
			MaxAge         int      `yaml:"maxAge"`         // CORS preflight cache max age in seconds
		} `yaml:"cors"`
	} `yaml:"server"`

	Logging struct {
		Level string `yaml:"level"` // debug|info|warn|error
		JSON  bool   `yaml:"json"`
	} `yaml:"logging"`

	Kubernetes struct {
		Kubeconfig string  `yaml:"kubeconfig"`
		TimeoutSec int     `yaml:"timeoutSec"`
		QPS        float32 `yaml:"qps"`
		Burst      int     `yaml:"burst"`
	} `yaml:"kubernetes"`

	Prometheus struct {
		Enabled       bool   `yaml:"enabled"`       // Enable Prometheus service
		Address       string `yaml:"address"`       // Prometheus server address
		TimeoutSec    int    `yaml:"timeoutSec"`    // Query timeout in seconds
		Username      string `yaml:"username"`      // Basic auth username
		Password      string `yaml:"password"`      // Basic auth password
		BearerToken   string `yaml:"bearerToken"`   // Bearer token for auth
		TLSSkipVerify bool   `yaml:"tlsSkipVerify"` // Skip TLS verification
		TLSCertFile   string `yaml:"tlsCertFile"`   // TLS certificate file
		TLSKeyFile    string `yaml:"tlsKeyFile"`    // TLS key file
		TLSCAFile     string `yaml:"tlsCAFile"`     // TLS CA file
	} `yaml:"prometheus"`

	Grafana struct {
		Enabled    bool   `yaml:"enabled"`    // Enable Grafana service
		URL        string `yaml:"url"`        // Grafana URL
		APIKey     string `yaml:"apiKey"`     // Grafana API key
		Username   string `yaml:"username"`   // Grafana username for basic auth
		Password   string `yaml:"password"`   // Grafana password for basic auth
		TimeoutSec int    `yaml:"timeoutSec"` // Request timeout in seconds
	} `yaml:"grafana"`

	Kibana struct {
		Enabled    bool   `yaml:"enabled"`    // Enable Kibana service
		URL        string `yaml:"url"`        // Kibana URL
		APIKey     string `yaml:"apiKey"`     // Kibana API key
		Username   string `yaml:"username"`   // Kibana username for basic auth
		Password   string `yaml:"password"`   // Kibana password for basic auth
		TimeoutSec int    `yaml:"timeoutSec"` // Request timeout in seconds
		SkipVerify bool   `yaml:"skipVerify"` // Skip TLS certificate verification
		Space      string `yaml:"space"`      // Kibana space (default: default)
	} `yaml:"kibana"`

	Helm struct {
		Enabled        bool              `yaml:"enabled"`        // Enable Helm service
		KubeconfigPath string            `yaml:"kubeconfigPath"` // Path to kubeconfig file for Helm operations
		Namespace      string            `yaml:"namespace"`      // Default namespace for Helm operations
		Debug          bool              `yaml:"debug"`          // Enable Helm debug mode
		TimeoutSec     int               `yaml:"timeoutSec"`     // Repository update timeout in seconds (default: 300)
		MaxRetries     int               `yaml:"maxRetries"`     // Max retries for failed repository updates (default: 3)
		UseMirrors     bool              `yaml:"useMirrors"`     // Use Chinese mirrors for overseas repositories (default: true)
		Mirrors        map[string]string `yaml:"mirrors"`        // Mirror URL mappings (original URL -> mirror URL)
	} `yaml:"helm"`

	Alertmanager struct {
		Enabled       bool   `yaml:"enabled"`       // Enable Alertmanager service
		Address       string `yaml:"address"`       // Alertmanager server address
		TimeoutSec    int    `yaml:"timeoutSec"`    // Request timeout in seconds
		Username      string `yaml:"username"`      // Basic auth username
		Password      string `yaml:"password"`      // Basic auth password
		BearerToken   string `yaml:"bearerToken"`   // Bearer token for auth
		TLSSkipVerify bool   `yaml:"tlsSkipVerify"` // Skip TLS verification
		TLSCertFile   string `yaml:"tlsCertFile"`   // TLS certificate file
		TLSKeyFile    string `yaml:"tlsKeyFile"`    // TLS key file
		TLSCAFile     string `yaml:"tlsCAFile"`     // TLS CA file
	} `yaml:"alertmanager"`

	Jaeger struct {
		Enabled    bool   `yaml:"enabled"`    // Enable Jaeger service
		Address    string `yaml:"address"`    // Jaeger server address
		TimeoutSec int    `yaml:"timeoutSec"` // Request timeout in seconds
	} `yaml:"jaeger"`

	OpenTelemetry struct {
		Enabled       bool   `yaml:"enabled"`       // Enable OpenTelemetry service
		Address       string `yaml:"address"`       // OpenTelemetry Collector address
		TimeoutSec    int    `yaml:"timeoutSec"`    // Request timeout in seconds
		Username      string `yaml:"username"`      // Basic auth username
		Password      string `yaml:"password"`      // Basic auth password
		BearerToken   string `yaml:"bearerToken"`   // Bearer token for auth
		TLSSkipVerify bool   `yaml:"tlsSkipVerify"` // Skip TLS verification
		TLSCertFile   string `yaml:"tlsCertFile"`   // TLS certificate file
		TLSKeyFile    string `yaml:"tlsKeyFile"`    // TLS key file
		TLSCAFile     string `yaml:"tlsCAFile"`     // TLS CA file
	} `yaml:"opentelemetry"`

	// Service and tool filtering

	// Authentication configuration
	Auth struct {
		Enabled      bool   `yaml:"enabled"`      // Enable authentication
		Mode         string `yaml:"mode"`         // auth mode: apikey | bearer | basic
		APIKey       string `yaml:"apiKey"`       // API key for apikey mode
		BearerToken  string `yaml:"bearerToken"`  // Bearer token for bearer mode
		Username     string `yaml:"username"`     // Username for basic auth
		Password     string `yaml:"password"`     // Password for basic auth
		JWTSecret    string `yaml:"jwtSecret"`    // Secret key for JWT validation
		JWTAlgorithm string `yaml:"jwtAlgorithm"` // JWT algorithm (HS256, RS256, etc.)
	} `yaml:"auth"`

	// Audit configuration
	Audit struct {
		Enabled    bool   `yaml:"enabled"`    // Enable audit logging
		Level      string `yaml:"level"`      // Log level: debug|info|warn|error
		MaxLogs    int    `yaml:"maxLogs"`    // Maximum number of audit logs to retain (default: 10000)
		Storage    string `yaml:"storage"`    // Storage type: memory | file | database | all
		Format     string `yaml:"format"`     // Log format: json | text
		MaxResults int    `yaml:"maxResults"` // Max query results
		TimeRange  int    `yaml:"timeRange"`  // Query time range in days
		File       struct {
			Path       string `yaml:"path"`       // Log file path
			MaxSizeMB  int    `yaml:"maxSizeMB"`  // Max file size in MB
			MaxBackups int    `yaml:"maxBackups"` // Max backup files
			MaxAgeDays int    `yaml:"maxAgeDays"` // Max age in days
			Compress   bool   `yaml:"compress"`   // Compress rotated files
			MaxLogs    int    `yaml:"maxLogs"`    // Max logs for memory storage
		} `yaml:"file"`
		Database struct {
			Type             string `yaml:"type"`             // Database type: sqlite | postgres | mysql
			ConnectionString string `yaml:"connectionString"` // Database connection string
			SQLitePath       string `yaml:"sqlitePath"`       // SQLite database file path
			TableName        string `yaml:"tableName"`        // Audit table name
			MaxRecords       int    `yaml:"maxRecords"`       // Max records to keep
			CleanupInterval  int    `yaml:"cleanupInterval"`  // Cleanup interval in hours
		} `yaml:"database"`
		Query struct {
			Enabled    bool `yaml:"enabled"`    // Enable query API
			MaxResults int  `yaml:"maxResults"` // Max results per query
			TimeRange  int  `yaml:"timeRange"`  // Max time range in days
		} `yaml:"query"`
		Alerts struct {
			Enabled          bool   `yaml:"enabled"`          // Enable alerts
			FailureThreshold int    `yaml:"failureThreshold"` // Failure threshold
			CheckIntervalSec int    `yaml:"checkIntervalSec"` // Check interval
			Method           string `yaml:"method"`           // Alert method: email|webhook|slack|none
			WebhookURL       string `yaml:"webhookURL"`       // Webhook URL
		} `yaml:"alerts"`
		Masking struct {
			Enabled   bool     `yaml:"enabled"`   // Enable masking
			Fields    []string `yaml:"fields"`    // Fields to mask
			MaskValue string   `yaml:"maskValue"` // Mask replacement value
		} `yaml:"masking"`
		Sampling struct {
			Enabled bool    `yaml:"enabled"` // Enable sampling
			Rate    float32 `yaml:"rate"`    // Sampling rate (0-1)
		} `yaml:"sampling"`
	} `yaml:"audit"`

	EnableDisable struct {
		DisabledServices []string `yaml:"disabledServices"` // Disabled services
		EnabledServices  []string `yaml:"enabledServices"`  // Enabled services
		DisabledTools    []string `yaml:"disabledTools"`    // Disabled tools
	} `yaml:"enableDisable"`

	Elasticsearch struct {
		Enabled       bool     `yaml:"enabled"`       // Enable Elasticsearch service
		Addresses     []string `yaml:"addresses"`     // Elasticsearch addresses
		Address       string   `yaml:"address"`       // Single address
		Username      string   `yaml:"username"`      // Basic auth username
		Password      string   `yaml:"password"`      // Basic auth password
		BearerToken   string   `yaml:"bearerToken"`   // Bearer token
		APIKey        string   `yaml:"apiKey"`        // API key
		TimeoutSec    int      `yaml:"timeoutSec"`    // Request timeout
		TLSSkipVerify bool     `yaml:"tlsSkipVerify"` // Skip TLS verify
		TLSCertFile   string   `yaml:"tlsCertFile"`   // TLS cert file
		TLSKeyFile    string   `yaml:"tlsKeyFile"`    // TLS key file
		TLSCAFile     string   `yaml:"tlsCAFile"`     // TLS CA file
	} `yaml:"elasticsearch"`
}

// Load loads configuration from YAML file (if provided) and merges environment overrides.
// It also validates the configuration before returning it.
//
// Environment variables:
//
//	MCP_MODE, MCP_ADDR, MCP_READ_TIMEOUT, MCP_WRITE_TIMEOUT, MCP_IDLE_TIMEOUT,
//	MCP_SSE_PATH_KUBERNETES, MCP_SSE_PATH_GRAFANA, MCP_SSE_PATH_PROMETHEUS, MCP_SSE_PATH_KIBANA,
//	MCP_SSE_PATH_HELM, MCP_SSE_PATH_ALERTMANAGER, MCP_SSE_PATH_AGGREGATE, MCP_SSE_PATH_UTILITIES,
//	MCP_LOG_LEVEL, MCP_LOG_JSON,
//	MCP_KUBECONFIG, MCP_K8S_TIMEOUT, MCP_K8S_QPS, MCP_K8S_BURST,
//	MCP_PROM_ENABLED, MCP_PROM_ADDRESS, MCP_PROM_TIMEOUT, MCP_PROM_USERNAME, MCP_PROM_PASSWORD,
//	MCP_PROM_BEARER_TOKEN, MCP_PROM_TLS_SKIP_VERIFY, MCP_PROM_TLS_CERT_FILE,
//	MCP_PROM_TLS_KEY_FILE, MCP_PROM_TLS_CA_FILE,
//	MCP_GRAFANA_ENABLED, MCP_GRAFANA_URL, MCP_GRAFANA_API_KEY,
//	MCP_GRAFANA_USERNAME, MCP_GRAFANA_PASSWORD, MCP_GRAFANA_TIMEOUT,
//	MCP_KIBANA_ENABLED, MCP_KIBANA_URL, MCP_KIBANA_API_KEY,
//	MCP_KIBANA_USERNAME, MCP_KIBANA_PASSWORD, MCP_KIBANA_TIMEOUT,
//	MCP_KIBANA_SKIP_VERIFY, MCP_KIBANA_SPACE,
//	MCP_HELM_ENABLED, MCP_HELM_KUBECONFIG, MCP_HELM_NAMESPACE, MCP_HELM_DEBUG,
//	MCP_HELM_TIMEOUT, MCP_HELM_MAX_RETRIES, MCP_HELM_USE_MIRRORS,
//	MCP_ELASTICSEARCH_ENABLED, MCP_ELASTICSEARCH_ADDRESSES, MCP_ELASTICSEARCH_ADDRESS,
//	MCP_ELASTICSEARCH_USERNAME, MCP_ELASTICSEARCH_PASSWORD, MCP_ELASTICSEARCH_BEARER_TOKEN,
//	MCP_ELASTICSEARCH_API_KEY, MCP_ELASTICSEARCH_TIMEOUT, MCP_ELASTICSEARCH_TLS_SKIP_VERIFY,
//	MCP_ELASTICSEARCH_TLS_CERT_FILE, MCP_ELASTICSEARCH_TLS_KEY_FILE, MCP_ELASTICSEARCH_TLS_CA_FILE,
//	MCP_ALERTMANAGER_ENABLED, MCP_ALERTMANAGER_ADDRESS, MCP_ALERTMANAGER_TIMEOUT,
//	MCP_ALERTMANAGER_USERNAME, MCP_ALERTMANAGER_PASSWORD, MCP_ALERTMANAGER_BEARER_TOKEN,
//	MCP_ALERTMANAGER_TLS_SKIP_VERIFY, MCP_ALERTMANAGER_TLS_CERT_FILE,
//	MCP_ALERTMANAGER_TLS_KEY_FILE, MCP_ALERTMANAGER_TLS_CA_FILE,
//	MCP_AUDIT_ENABLED, MCP_AUDIT_LEVEL, MCP_AUDIT_MAX_RESULTS, MCP_AUDIT_TIME_RANGE,
//	MCP_AUDIT_QUERY_ENABLED, MCP_AUDIT_QUERY_MAX_RESULTS, MCP_AUDIT_QUERY_TIME_RANGE,
//	MCP_AUDIT_ALERTS_ENABLED, MCP_AUDIT_ALERTS_FAILURE_THRESHOLD, MCP_AUDIT_ALERTS_CHECK_INTERVAL,
//	MCP_AUDIT_ALERTS_METHOD, MCP_AUDIT_ALERTS_WEBHOOK_URL,
//	MCP_AUDIT_MASKING_ENABLED, MCP_AUDIT_MASKING_FIELDS, MCP_AUDIT_MASKING_VALUE,
//	MCP_AUDIT_SAMPLING_ENABLED, MCP_AUDIT_SAMPLING_RATE,
//	MCP_AUTH_ENABLED, MCP_AUTH_MODE, MCP_AUTH_API_KEY, MCP_AUTH_BEARER_TOKEN,
//	MCP_AUTH_USERNAME, MCP_AUTH_PASSWORD, MCP_AUTH_JWT_SECRET, MCP_AUTH_JWT_ALGORITHM,
//	MCP_DISABLED_SERVICES, MCP_ENABLED_SERVICES, MCP_DISABLED_TOOLS
func Load(path string) (*AppConfig, error) {
	loader := NewConfigLoader()
	config, err := loader.Load(path)
	if err != nil {
		return nil, err
	}
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}
	return config, nil
}

// Validate validates the configuration and returns an error if invalid
func (c *AppConfig) Validate() error {
	// Validate authentication configuration
	if c.Auth.Enabled {
		if c.Auth.Mode == "" {
			return fmt.Errorf("auth mode is required when auth is enabled")
		}
		switch c.Auth.Mode {
		case "apikey":
			if c.Auth.APIKey == "" {
				return fmt.Errorf("auth API key is required for apikey mode")
			}
		case "bearer":
			if c.Auth.BearerToken == "" {
				return fmt.Errorf("auth bearer token is required for bearer mode")
			}
		case "basic":
			if c.Auth.Username == "" || c.Auth.Password == "" {
				return fmt.Errorf("auth username and password are required for basic mode")
			}
		default:
			return fmt.Errorf("invalid auth mode: %s (must be apikey, bearer, or basic)", c.Auth.Mode)
		}
	}

	// Validate service configurations
	if c.Prometheus.Enabled && c.Prometheus.Address == "" {
		return fmt.Errorf("prometheus address is required when service is enabled")
	}
	if c.Grafana.Enabled && c.Grafana.URL == "" {
		return fmt.Errorf("grafana URL is required when service is enabled")
	}
	if c.Kibana.Enabled && c.Kibana.URL == "" {
		return fmt.Errorf("kibana URL is required when service is enabled")
	}
	if c.Alertmanager.Enabled && c.Alertmanager.Address == "" {
		return fmt.Errorf("alertmanager address is required when service is enabled")
	}
	if c.Jaeger.Enabled && c.Jaeger.Address == "" {
		return fmt.Errorf("jaeger address is required when service is enabled")
	}
	if c.OpenTelemetry.Enabled && c.OpenTelemetry.Address == "" {
		return fmt.Errorf("opentelemetry address is required when service is enabled")
	}
	if c.Elasticsearch.Enabled && len(c.Elasticsearch.Addresses) == 0 && c.Elasticsearch.Address == "" {
		return fmt.Errorf("elasticsearch address or addresses is required when service is enabled")
	}

	// Validate audit configuration
	if c.Audit.Enabled {
		if c.Audit.Storage == "" {
			return fmt.Errorf("audit storage type is required when audit is enabled")
		}
		if c.Audit.Storage == "database" && c.Audit.Database.Type == "" {
			return fmt.Errorf("audit database type is required when database storage is enabled")
		}
		if c.Audit.Storage == "file" && c.Audit.File.Path == "" {
			return fmt.Errorf("audit file path is required when file storage is enabled")
		}
	}

	// Validate timeout values are reasonable
	if c.Server.ReadTimeoutSec < 0 || c.Server.ReadTimeoutSec > 3600 {
		return fmt.Errorf("server read timeout must be between 0 and 3600 seconds")
	}
	if c.Server.WriteTimeoutSec < 0 || c.Server.WriteTimeoutSec > 3600 {
		return fmt.Errorf("server write timeout must be between 0 and 3600 seconds")
	}
	if c.Server.IdleTimeoutSec < 0 || c.Server.IdleTimeoutSec > 3600 {
		return fmt.Errorf("server idle timeout must be between 0 and 3600 seconds")
	}

	return nil
}

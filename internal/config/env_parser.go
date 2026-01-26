package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// EnvParser handles environment variable parsing and overrides
type EnvParser struct{}

// NewEnvParser creates a new environment variable parser
func NewEnvParser() *EnvParser {
	return &EnvParser{}
}

// Parse applies environment variable overrides to the configuration
func (p *EnvParser) Parse(cfg *AppConfig) *AppConfig {
	over := func(key string) (string, bool) {
		v, ok := os.LookupEnv(key)
		return strings.TrimSpace(v), ok
	}

	p.parseServerConfig(cfg, over)
	p.parseLoggingConfig(cfg, over)
	p.parseKubernetesConfig(cfg, over)
	p.parsePrometheusConfig(cfg, over)
	p.parseGrafanaConfig(cfg, over)
	p.parseKibanaConfig(cfg, over)
	p.parseHelmConfig(cfg, over)
	p.parseElasticsearchConfig(cfg, over)
	p.parseAlertmanagerConfig(cfg, over)
	p.parseJaegerConfig(cfg, over)
	p.parseAuditConfig(cfg, over)
	p.parseAuthConfig(cfg, over)
	p.parseEnableDisableConfig(cfg, over)

	return cfg
}

func (p *EnvParser) parseServerConfig(cfg *AppConfig, over func(string) (string, bool)) {
	if v, ok := over("MCP_MODE"); ok {
		cfg.Server.Mode = v
	}
	if v, ok := over("MCP_ADDR"); ok {
		cfg.Server.Addr = v
	}
	if v, ok := over("MCP_READ_TIMEOUT"); ok {
		cfg.Server.ReadTimeoutSec = atoiDefault(v, cfg.Server.ReadTimeoutSec)
	}
	if v, ok := over("MCP_WRITE_TIMEOUT"); ok {
		cfg.Server.WriteTimeoutSec = atoiDefault(v, cfg.Server.WriteTimeoutSec)
	}
	if v, ok := over("MCP_IDLE_TIMEOUT"); ok {
		cfg.Server.IdleTimeoutSec = atoiDefault(v, cfg.Server.IdleTimeoutSec)
	}

	// SSE paths configuration
	if v, ok := over("MCP_SSE_PATH_KUBERNETES"); ok {
		cfg.Server.SSEPaths.Kubernetes = v
	}
	if v, ok := over("MCP_SSE_PATH_GRAFANA"); ok {
		cfg.Server.SSEPaths.Grafana = v
	}
	if v, ok := over("MCP_SSE_PATH_PROMETHEUS"); ok {
		cfg.Server.SSEPaths.Prometheus = v
	}
	if v, ok := over("MCP_SSE_PATH_KIBANA"); ok {
		cfg.Server.SSEPaths.Kibana = v
	}
	if v, ok := over("MCP_SSE_PATH_HELM"); ok {
		cfg.Server.SSEPaths.Helm = v
	}
	if v, ok := over("MCP_SSE_PATH_ALERTMANAGER"); ok {
		cfg.Server.SSEPaths.Alertmanager = v
	}
	if v, ok := over("MCP_SSE_PATH_AGGREGATE"); ok {
		cfg.Server.SSEPaths.Aggregate = v
	}
	if v, ok := over("MCP_SSE_PATH_UTILITIES"); ok {
		cfg.Server.SSEPaths.Utilities = v
	}

	// Streamable-HTTP paths configuration
	if v, ok := over("MCP_STREAMABLE_HTTP_PATH_KUBERNETES"); ok {
		cfg.Server.StreamableHTTPPaths.Kubernetes = v
	}
	if v, ok := over("MCP_STREAMABLE_HTTP_PATH_GRAFANA"); ok {
		cfg.Server.StreamableHTTPPaths.Grafana = v
	}
	if v, ok := over("MCP_STREAMABLE_HTTP_PATH_PROMETHEUS"); ok {
		cfg.Server.StreamableHTTPPaths.Prometheus = v
	}
	if v, ok := over("MCP_STREAMABLE_HTTP_PATH_KIBANA"); ok {
		cfg.Server.StreamableHTTPPaths.Kibana = v
	}
	if v, ok := over("MCP_STREAMABLE_HTTP_PATH_HELM"); ok {
		cfg.Server.StreamableHTTPPaths.Helm = v
	}
	if v, ok := over("MCP_STREAMABLE_HTTP_PATH_ALERTMANAGER"); ok {
		cfg.Server.StreamableHTTPPaths.Alertmanager = v
	}
	if v, ok := over("MCP_STREAMABLE_HTTP_PATH_AGGREGATE"); ok {
		cfg.Server.StreamableHTTPPaths.Aggregate = v
	}
	if v, ok := over("MCP_STREAMABLE_HTTP_PATH_UTILITIES"); ok {
		cfg.Server.StreamableHTTPPaths.Utilities = v
	}
}

func (p *EnvParser) parseLoggingConfig(cfg *AppConfig, over func(string) (string, bool)) {
	if v, ok := over("MCP_LOG_LEVEL"); ok {
		cfg.Logging.Level = v
	}
	if v, ok := over("MCP_LOG_JSON"); ok {
		cfg.Logging.JSON = isTrue(v)
	}
}

func (p *EnvParser) parseKubernetesConfig(cfg *AppConfig, over func(string) (string, bool)) {
	if v, ok := over("MCP_KUBECONFIG"); ok {
		cfg.Kubernetes.Kubeconfig = v
	}
	if v, ok := over("MCP_K8S_TIMEOUT"); ok {
		cfg.Kubernetes.TimeoutSec = atoiDefault(v, cfg.Kubernetes.TimeoutSec)
	}
	if v, ok := over("MCP_K8S_QPS"); ok {
		cfg.Kubernetes.QPS = atofDefault(v, cfg.Kubernetes.QPS)
	}
	if v, ok := over("MCP_K8S_BURST"); ok {
		cfg.Kubernetes.Burst = atoiDefault(v, cfg.Kubernetes.Burst)
	}
}

func (p *EnvParser) parsePrometheusConfig(cfg *AppConfig, over func(string) (string, bool)) {
	if v, ok := over("MCP_PROM_ENABLED"); ok {
		cfg.Prometheus.Enabled = isTrue(v)
	}
	if v, ok := over("MCP_PROM_ADDRESS"); ok {
		cfg.Prometheus.Address = v
	}
	if v, ok := over("MCP_PROM_TIMEOUT"); ok {
		cfg.Prometheus.TimeoutSec = atoiDefault(v, cfg.Prometheus.TimeoutSec)
	}
	if v, ok := over("MCP_PROM_USERNAME"); ok {
		cfg.Prometheus.Username = v
	}
	if v, ok := over("MCP_PROM_PASSWORD"); ok {
		cfg.Prometheus.Password = v
	}
	if v, ok := over("MCP_PROM_BEARER_TOKEN"); ok {
		cfg.Prometheus.BearerToken = v
	}
	if v, ok := over("MCP_PROM_TLS_SKIP_VERIFY"); ok {
		cfg.Prometheus.TLSSkipVerify = isTrue(v)
	}
	if v, ok := over("MCP_PROM_TLS_CERT_FILE"); ok {
		cfg.Prometheus.TLSCertFile = v
	}
	if v, ok := over("MCP_PROM_TLS_KEY_FILE"); ok {
		cfg.Prometheus.TLSKeyFile = v
	}
	if v, ok := over("MCP_PROM_TLS_CA_FILE"); ok {
		cfg.Prometheus.TLSCAFile = v
	}
}

func (p *EnvParser) parseGrafanaConfig(cfg *AppConfig, over func(string) (string, bool)) {
	if v, ok := over("MCP_GRAFANA_ENABLED"); ok {
		cfg.Grafana.Enabled = isTrue(v)
	}
	if v, ok := over("MCP_GRAFANA_URL"); ok {
		cfg.Grafana.URL = v
	}
	if v, ok := over("MCP_GRAFANA_API_KEY"); ok {
		cfg.Grafana.APIKey = v
	}
	if v, ok := over("MCP_GRAFANA_USERNAME"); ok {
		cfg.Grafana.Username = v
	}
	if v, ok := over("MCP_GRAFANA_PASSWORD"); ok {
		cfg.Grafana.Password = v
	}
	if v, ok := over("MCP_GRAFANA_TIMEOUT"); ok {
		cfg.Grafana.TimeoutSec = atoiDefault(v, cfg.Grafana.TimeoutSec)
	}
}

func (p *EnvParser) parseKibanaConfig(cfg *AppConfig, over func(string) (string, bool)) {
	if v, ok := over("MCP_KIBANA_ENABLED"); ok {
		cfg.Kibana.Enabled = isTrue(v)
	}
	if v, ok := over("MCP_KIBANA_URL"); ok {
		cfg.Kibana.URL = v
	}
	if v, ok := over("MCP_KIBANA_API_KEY"); ok {
		cfg.Kibana.APIKey = v
	}
	if v, ok := over("MCP_KIBANA_USERNAME"); ok {
		cfg.Kibana.Username = v
	}
	if v, ok := over("MCP_KIBANA_PASSWORD"); ok {
		cfg.Kibana.Password = v
	}
	if v, ok := over("MCP_KIBANA_TIMEOUT"); ok {
		cfg.Kibana.TimeoutSec = atoiDefault(v, cfg.Kibana.TimeoutSec)
	}
	if v, ok := over("MCP_KIBANA_SKIP_VERIFY"); ok {
		cfg.Kibana.SkipVerify = isTrue(v)
	}
	if v, ok := over("MCP_KIBANA_SPACE"); ok {
		cfg.Kibana.Space = v
	}
}

func (p *EnvParser) parseHelmConfig(cfg *AppConfig, over func(string) (string, bool)) {
	if v, ok := over("MCP_HELM_ENABLED"); ok {
		cfg.Helm.Enabled = isTrue(v)
	}
	if v, ok := over("MCP_HELM_KUBECONFIG"); ok {
		cfg.Helm.KubeconfigPath = v
	}
	if v, ok := over("MCP_HELM_NAMESPACE"); ok {
		cfg.Helm.Namespace = v
	}
	if v, ok := over("MCP_HELM_DEBUG"); ok {
		cfg.Helm.Debug = isTrue(v)
	}
	if v, ok := over("MCP_HELM_TIMEOUT"); ok {
		cfg.Helm.TimeoutSec = atoiDefault(v, cfg.Helm.TimeoutSec)
	}
	if v, ok := over("MCP_HELM_MAX_RETRIES"); ok {
		cfg.Helm.MaxRetries = atoiDefault(v, cfg.Helm.MaxRetries)
	}
	if v, ok := over("MCP_HELM_USE_MIRRORS"); ok {
		cfg.Helm.UseMirrors = isTrue(v)
	}
	if cfg.Helm.Mirrors == nil {
		cfg.Helm.Mirrors = make(map[string]string)
	}
}

func (p *EnvParser) parseElasticsearchConfig(cfg *AppConfig, over func(string) (string, bool)) {
	if v, ok := over("MCP_ELASTICSEARCH_ENABLED"); ok {
		cfg.Elasticsearch.Enabled = isTrue(v)
	}
	if v, ok := over("MCP_ELASTICSEARCH_ADDRESSES"); ok {
		cfg.Elasticsearch.Addresses = strings.Split(v, ",")
	}
	if v, ok := over("MCP_ELASTICSEARCH_ADDRESS"); ok {
		cfg.Elasticsearch.Address = v
	}
	if v, ok := over("MCP_ELASTICSEARCH_USERNAME"); ok {
		cfg.Elasticsearch.Username = v
	}
	if v, ok := over("MCP_ELASTICSEARCH_PASSWORD"); ok {
		cfg.Elasticsearch.Password = v
	}
	if v, ok := over("MCP_ELASTICSEARCH_BEARER_TOKEN"); ok {
		cfg.Elasticsearch.BearerToken = v
	}
	if v, ok := over("MCP_ELASTICSEARCH_API_KEY"); ok {
		cfg.Elasticsearch.APIKey = v
	}
	if v, ok := over("MCP_ELASTICSEARCH_TIMEOUT"); ok {
		cfg.Elasticsearch.TimeoutSec = atoiDefault(v, cfg.Elasticsearch.TimeoutSec)
	}
	if v, ok := over("MCP_ELASTICSEARCH_TLS_SKIP_VERIFY"); ok {
		cfg.Elasticsearch.TLSSkipVerify = isTrue(v)
	}
	if v, ok := over("MCP_ELASTICSEARCH_TLS_CERT_FILE"); ok {
		cfg.Elasticsearch.TLSCertFile = v
	}
	if v, ok := over("MCP_ELASTICSEARCH_TLS_KEY_FILE"); ok {
		cfg.Elasticsearch.TLSKeyFile = v
	}
	if v, ok := over("MCP_ELASTICSEARCH_TLS_CA_FILE"); ok {
		cfg.Elasticsearch.TLSCAFile = v
	}
}

func (p *EnvParser) parseAlertmanagerConfig(cfg *AppConfig, over func(string) (string, bool)) {
	if v, ok := over("MCP_ALERTMANAGER_ENABLED"); ok {
		cfg.Alertmanager.Enabled = isTrue(v)
	}
	if v, ok := over("MCP_ALERTMANAGER_ADDRESS"); ok {
		cfg.Alertmanager.Address = v
	}
	if v, ok := over("MCP_ALERTMANAGER_TIMEOUT"); ok {
		cfg.Alertmanager.TimeoutSec = atoiDefault(v, cfg.Alertmanager.TimeoutSec)
	}
	if v, ok := over("MCP_ALERTMANAGER_USERNAME"); ok {
		cfg.Alertmanager.Username = v
	}
	if v, ok := over("MCP_ALERTMANAGER_PASSWORD"); ok {
		cfg.Alertmanager.Password = v
	}
	if v, ok := over("MCP_ALERTMANAGER_BEARER_TOKEN"); ok {
		cfg.Alertmanager.BearerToken = v
	}
	if v, ok := over("MCP_ALERTMANAGER_TLS_SKIP_VERIFY"); ok {
		cfg.Alertmanager.TLSSkipVerify = isTrue(v)
	}
	if v, ok := over("MCP_ALERTMANAGER_TLS_CERT_FILE"); ok {
		cfg.Alertmanager.TLSCertFile = v
	}
	if v, ok := over("MCP_ALERTMANAGER_TLS_KEY_FILE"); ok {
		cfg.Alertmanager.TLSKeyFile = v
	}
	if v, ok := over("MCP_ALERTMANAGER_TLS_CA_FILE"); ok {
		cfg.Alertmanager.TLSCAFile = v
	}
}

func (p *EnvParser) parseJaegerConfig(cfg *AppConfig, over func(string) (string, bool)) {
	if v, ok := over("MCP_JAEGER_ENABLED"); ok {
		cfg.Jaeger.Enabled = isTrue(v)
	}
	if v, ok := over("MCP_JAEGER_ADDRESS"); ok {
		cfg.Jaeger.Address = v
	}
	if v, ok := over("MCP_JAEGER_TIMEOUT"); ok {
		cfg.Jaeger.TimeoutSec = atoiDefault(v, cfg.Jaeger.TimeoutSec)
	}
}

func (p *EnvParser) parseAuditConfig(cfg *AppConfig, over func(string) (string, bool)) {
	if v, ok := over("MCP_AUDIT_ENABLED"); ok {
		cfg.Audit.Enabled = isTrue(v)
	}
	if v, ok := over("MCP_AUDIT_LEVEL"); ok {
		cfg.Audit.Level = v
	}
	if v, ok := over("MCP_AUDIT_MAX_RESULTS"); ok {
		cfg.Audit.MaxResults = atoiDefault(v, cfg.Audit.MaxResults)
	}
	if v, ok := over("MCP_AUDIT_TIME_RANGE"); ok {
		cfg.Audit.TimeRange = atoiDefault(v, cfg.Audit.TimeRange)
	}
	if v, ok := over("MCP_AUDIT_QUERY_ENABLED"); ok {
		cfg.Audit.Query.Enabled = isTrue(v)
	}
	if v, ok := over("MCP_AUDIT_QUERY_MAX_RESULTS"); ok {
		cfg.Audit.Query.MaxResults = atoiDefault(v, cfg.Audit.Query.MaxResults)
	}
	if v, ok := over("MCP_AUDIT_QUERY_TIME_RANGE"); ok {
		cfg.Audit.Query.TimeRange = atoiDefault(v, cfg.Audit.Query.TimeRange)
	}
	if v, ok := over("MCP_AUDIT_ALERTS_ENABLED"); ok {
		cfg.Audit.Alerts.Enabled = isTrue(v)
	}
	if v, ok := over("MCP_AUDIT_ALERTS_FAILURE_THRESHOLD"); ok {
		cfg.Audit.Alerts.FailureThreshold = atoiDefault(v, cfg.Audit.Alerts.FailureThreshold)
	}
	if v, ok := over("MCP_AUDIT_ALERTS_CHECK_INTERVAL"); ok {
		cfg.Audit.Alerts.CheckIntervalSec = atoiDefault(v, cfg.Audit.Alerts.CheckIntervalSec)
	}
	if v, ok := over("MCP_AUDIT_ALERTS_METHOD"); ok {
		cfg.Audit.Alerts.Method = v
	}
	if v, ok := over("MCP_AUDIT_ALERTS_WEBHOOK_URL"); ok {
		cfg.Audit.Alerts.WebhookURL = v
	}
	if v, ok := over("MCP_AUDIT_MASKING_ENABLED"); ok {
		cfg.Audit.Masking.Enabled = isTrue(v)
	}
	if v, ok := over("MCP_AUDIT_MASKING_FIELDS"); ok {
		cfg.Audit.Masking.Fields = strings.Split(v, ",")
	}
	if v, ok := over("MCP_AUDIT_MASKING_VALUE"); ok {
		cfg.Audit.Masking.MaskValue = v
	}
	if v, ok := over("MCP_AUDIT_SAMPLING_ENABLED"); ok {
		cfg.Audit.Sampling.Enabled = isTrue(v)
	}
	if v, ok := over("MCP_AUDIT_SAMPLING_RATE"); ok {
		cfg.Audit.Sampling.Rate = atofDefault(v, cfg.Audit.Sampling.Rate)
	}
}

func (p *EnvParser) parseAuthConfig(cfg *AppConfig, over func(string) (string, bool)) {
	if v, ok := over("MCP_AUTH_ENABLED"); ok {
		cfg.Auth.Enabled = isTrue(v)
	}
	if v, ok := over("MCP_AUTH_MODE"); ok {
		cfg.Auth.Mode = v
	}
	if v, ok := over("MCP_AUTH_API_KEY"); ok {
		cfg.Auth.APIKey = v
	}
	if v, ok := over("MCP_AUTH_BEARER_TOKEN"); ok {
		cfg.Auth.BearerToken = v
	}
	if v, ok := over("MCP_AUTH_USERNAME"); ok {
		cfg.Auth.Username = v
	}
	if v, ok := over("MCP_AUTH_PASSWORD"); ok {
		cfg.Auth.Password = v
	}
	if v, ok := over("MCP_AUTH_JWT_SECRET"); ok {
		cfg.Auth.JWTSecret = v
	}
	if v, ok := over("MCP_AUTH_JWT_ALGORITHM"); ok {
		cfg.Auth.JWTAlgorithm = v
	}
}

func (p *EnvParser) parseEnableDisableConfig(cfg *AppConfig, over func(string) (string, bool)) {
	if v, ok := over("MCP_DISABLED_SERVICES"); ok {
		cfg.EnableDisable.DisabledServices = strings.Split(v, ",")
	}
	if v, ok := over("MCP_ENABLED_SERVICES"); ok {
		cfg.EnableDisable.EnabledServices = strings.Split(v, ",")
	}
	if v, ok := over("MCP_DISABLED_TOOLS"); ok {
		cfg.EnableDisable.DisabledTools = strings.Split(v, ",")
	}
}

// Helper functions
func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		logrus.WithField("value", s).WithField("default", def).Warnf("Invalid integer value, using default: %v", err)
		return def
	}
	return i
}

func atofDefault(s string, def float32) float32 {
	if s == "" {
		return def
	}
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		logrus.WithField("value", s).WithField("default", def).Warnf("Invalid float value, using default: %v", err)
		return def
	}
	return float32(f)
}

func isTrue(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "1" || s == "true" || s == "yes" || s == "on"
}

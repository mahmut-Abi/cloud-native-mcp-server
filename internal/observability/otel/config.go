package otel

import "time"

// Config defines configuration for OpenTelemetry
type Config struct {
	// Enable enables or disables OpenTelemetry
	Enabled bool `yaml:"enabled" json:"enabled"`

	// ServiceName is the name of this service
	ServiceName string `yaml:"serviceName" json:"serviceName"`

	// ServiceVersion is the version of this service
	ServiceVersion string `yaml:"serviceVersion" json:"serviceVersion"`

	// Environment is the deployment environment (e.g., production, staging, dev)
	Environment string `yaml:"environment" json:"environment"`

	// Endpoint is the OTLP collector endpoint (e.g., http://localhost:4317)
	Endpoint string `yaml:"endpoint" json:"endpoint"`

	// Insecure disables TLS verification for OTLP connection
	Insecure bool `yaml:"insecure" json:"insecure"`

	// Headers are additional headers to send with OTLP requests
	Headers map[string]string `yaml:"headers" json:"headers"`

	// TracingConfig configures distributed tracing
	TracingConfig TracingConfig `yaml:"tracing" json:"tracing"`

	// MetricsConfig configures metrics export
	MetricsConfig MetricsConfig `yaml:"metrics" json:"metrics"`
}

// TracingConfig configures distributed tracing
type TracingConfig struct {
	// Enabled enables or disables tracing
	Enabled bool `yaml:"enabled" json:"enabled"`

	// SampleRate determines the percentage of traces to sample (0.0 to 1.0)
	SampleRate float64 `yaml:"sampleRate" json:"sampleRate"`

	// ExportTimeout is the timeout for exporting traces
	ExportTimeout time.Duration `yaml:"exportTimeout" json:"exportTimeout"`

	// BatchTimeout is the timeout for batching trace exports
	BatchTimeout time.Duration `yaml:"batchTimeout" json:"batchTimeout"`

	// MaxExportBatchSize is the maximum batch size for trace exports
	MaxExportBatchSize int `yaml:"maxExportBatchSize" json:"maxExportBatchSize"`
}

// MetricsConfig configures metrics export
type MetricsConfig struct {
	// Enabled enables or disables metrics export
	Enabled bool `yaml:"enabled" json:"enabled"`

	// ExportInterval is the interval for exporting metrics
	ExportInterval time.Duration `yaml:"exportInterval" json:"exportInterval"`

	// ExportTimeout is the timeout for exporting metrics
	ExportTimeout time.Duration `yaml:"exportTimeout" json:"exportTimeout"`

	// Temporality specifies the temporality of metrics (cumulative or delta)
	Temporality string `yaml:"temporality" json:"temporality"`
}

// DefaultConfig returns default OpenTelemetry configuration
func DefaultConfig() *Config {
	return &Config{
		Enabled:        false,
		ServiceName:    "cloud-native-mcp-server",
		ServiceVersion: "1.0.0",
		Environment:    "production",
		Endpoint:       "http://localhost:4317",
		Insecure:       true,
		Headers:        make(map[string]string),
		TracingConfig: TracingConfig{
			Enabled:            true,
			SampleRate:         0.1, // Sample 10% of traces
			ExportTimeout:      30 * time.Second,
			BatchTimeout:       5 * time.Second,
			MaxExportBatchSize: 512,
		},
		MetricsConfig: MetricsConfig{
			Enabled:        true,
			ExportInterval: 15 * time.Second,
			ExportTimeout:  30 * time.Second,
			Temporality:    "cumulative",
		},
	}
}

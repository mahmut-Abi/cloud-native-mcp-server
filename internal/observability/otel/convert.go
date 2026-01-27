package otel

import (
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
)

// FromAppConfig converts AppConfig to OTEL Config
func FromAppConfig(appConfig *config.AppConfig, version string) *Config {
	if appConfig == nil || !appConfig.OTEL.Enabled {
		return nil
	}

	// Set default values
	serviceName := appConfig.OTEL.ServiceName
	if serviceName == "" {
		serviceName = "cloud-native-mcp-server"
	}

	serviceVersion := appConfig.OTEL.ServiceVersion
	if serviceVersion == "" {
		serviceVersion = version
	}

	environment := appConfig.OTEL.Environment
	if environment == "" {
		environment = "production"
	}

	endpoint := appConfig.OTEL.Endpoint
	if endpoint == "" {
		endpoint = "http://localhost:4317"
	}

	return &Config{
		Enabled:        appConfig.OTEL.Enabled,
		ServiceName:    serviceName,
		ServiceVersion: serviceVersion,
		Environment:    environment,
		Endpoint:       endpoint,
		Insecure:       appConfig.OTEL.Insecure,
		Headers:        make(map[string]string),
		TracingConfig: TracingConfig{
			Enabled:            appConfig.OTEL.Tracing.Enabled,
			SampleRate:         appConfig.OTEL.Tracing.SampleRate,
			ExportTimeout:      time.Duration(appConfig.OTEL.Tracing.ExportTimeoutSec) * time.Second,
			BatchTimeout:       time.Duration(appConfig.OTEL.Tracing.BatchTimeoutSec) * time.Second,
			MaxExportBatchSize: appConfig.OTEL.Tracing.MaxExportBatchSize,
		},
		MetricsConfig: MetricsConfig{
			Enabled:        appConfig.OTEL.Metrics.Enabled,
			ExportInterval: time.Duration(appConfig.OTEL.Metrics.ExportIntervalSec) * time.Second,
			ExportTimeout:  time.Duration(appConfig.OTEL.Metrics.ExportTimeoutSec) * time.Second,
			Temporality:    appConfig.OTEL.Metrics.Temporality,
		},
	}
}

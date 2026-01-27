package otel

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"github.com/sirupsen/logrus"
)

var (
	globalTracerProvider *sdktrace.TracerProvider
	globalMeterProvider  *sdkmetric.MeterProvider
	tracer              trace.Tracer
	meter               metric.Meter
)

// Init initializes OpenTelemetry with the given configuration
func Init(cfg *Config) error {
	if !cfg.Enabled {
		logrus.Info("OpenTelemetry is disabled")
		return nil
	}

	logrus.Info("Initializing OpenTelemetry...")

	// Create resource with service information
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
			attribute.String("host.name", getHostname()),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	// Initialize tracing if enabled
	if cfg.TracingConfig.Enabled {
		if err := initTracing(cfg, res); err != nil {
			logrus.WithError(err).Warn("Failed to initialize tracing")
		}
	}

	// Initialize metrics if enabled
	if cfg.MetricsConfig.Enabled {
		if err := initMetrics(cfg, res); err != nil {
			logrus.WithError(err).Warn("Failed to initialize metrics")
		}
	}

	// Set global providers
	if globalTracerProvider != nil {
		otel.SetTracerProvider(globalTracerProvider)
		tracer = otel.Tracer(cfg.ServiceName)
		logrus.Info("Tracing initialized")
	}

	if globalMeterProvider != nil {
		otel.SetMeterProvider(globalMeterProvider)
		meter = otel.Meter(cfg.ServiceName)
		logrus.Info("Metrics initialized")
	}

	logrus.Info("OpenTelemetry initialized successfully")
	return nil
}

// initTracing initializes distributed tracing
func initTracing(cfg *Config, res *resource.Resource) error {
	ctx := context.Background()

	// Create OTLP trace exporter
	var opts []otlptracegrpc.Option
	opts = append(opts, otlptracegrpc.WithEndpoint(cfg.Endpoint))
	if cfg.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	// Create trace provider with sampler
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter,
			sdktrace.WithBatchTimeout(cfg.TracingConfig.BatchTimeout),
			sdktrace.WithMaxExportBatchSize(cfg.TracingConfig.MaxExportBatchSize),
			sdktrace.WithExportTimeout(cfg.TracingConfig.ExportTimeout),
		),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.TracingConfig.SampleRate))),
	)

	globalTracerProvider = tp
	return nil
}

// initMetrics initializes metrics export
func initMetrics(cfg *Config, res *resource.Resource) error {
	ctx := context.Background()

	// Create OTLP metrics exporter
	var opts []otlpmetricgrpc.Option
	opts = append(opts, otlpmetricgrpc.WithEndpoint(cfg.Endpoint))
	if cfg.Insecure {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}

	exporter, err := otlpmetricgrpc.New(ctx, opts...)
	if err != nil {
		return fmt.Errorf("failed to create OTLP metrics exporter: %w", err)
	}

	// Create meter provider
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter,
			sdkmetric.WithInterval(cfg.MetricsConfig.ExportInterval),
			sdkmetric.WithTimeout(cfg.MetricsConfig.ExportTimeout),
		)),
	)

	globalMeterProvider = mp
	return nil
}

// Shutdown gracefully shuts down OpenTelemetry
func Shutdown(ctx context.Context) error {
	logrus.Info("Shutting down OpenTelemetry...")

	var errs []error

	if globalTracerProvider != nil {
		if err := globalTracerProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to shutdown tracer provider: %w", err))
		}
	}

	if globalMeterProvider != nil {
		if err := globalMeterProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to shutdown meter provider: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}

	logrus.Info("OpenTelemetry shutdown complete")
	return nil
}

// GetTracer returns the global tracer
func GetTracer() trace.Tracer {
	return tracer
}

// GetMeter returns the global meter
func GetMeter() metric.Meter {
	return meter
}

// getHostname returns the hostname
func getHostname() string {
	hostname := "unknown"
	// Try to get hostname (simplified implementation)
	// In production, use os.Hostname()
	return hostname
}
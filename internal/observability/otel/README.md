# OpenTelemetry Integration

This directory contains the OpenTelemetry integration for the MCP server.

## Overview

The MCP server now supports exporting telemetry data (traces and metrics) to an OpenTelemetry Collector using the OpenTelemetry Protocol (OTLP).

## Features

- **Distributed Tracing**: Automatic HTTP request tracing and manual span creation
- **Metrics Export**: Export custom metrics alongside existing Prometheus metrics
- **Service Information**: Automatic service metadata (name, version, environment)
- **Flexible Configuration**: Support for sampling, batching, and custom endpoints

## Architecture

```
MCP Server → OpenTelemetry SDK → OTLP Collector → Backend
                                          ├─ Prometheus (metrics)
                                          ├─ Jaeger/Tempo (traces)
                                          └─ Loki (logs)
```

## Configuration

Add OTEL configuration to your `config.yaml`:

```yaml
otel:
  enabled: true
  serviceName: "cloud-native-mcp-server"
  serviceVersion: "1.0.0"
  environment: "production"
  endpoint: "localhost:4317"
  insecure: true

  tracing:
    enabled: true
    sampleRate: 0.1  # Sample 10% of traces

  metrics:
    enabled: true
    exportIntervalSec: 15
```

## Usage

### Automatic HTTP Tracing

The middleware automatically traces all HTTP requests:

```go
import "github.com/mahmut-Abi/cloud-native-mcp-server/internal/observability/otel"

// Wrap your HTTP handler
handler := otel.Middleware("mcp-server")(httpHandler)
```

### Manual Span Creation

Create spans for custom operations:

```go
import (
    "github.com/mahmut-Abi/cloud-native-mcp-server/internal/observability/otel"
    "go.opentelemetry.io/otel/attribute"
)

helper := otel.NewSpanHelper()

// Using WithSpan helper
err := otel.WithSpan(ctx, "custom-operation", func(ctx context.Context, span trace.Span) error {
    // Add attributes
    otel.SetAttributes(span,
        attribute.String("key", "value"),
    )

    // Do work...
    return nil
})

// Manual span creation
ctx, span := helper.StartSpan(ctx, "manual-operation")
defer span.End()

// Add event
helper.AddEvent(span, "important-event",
    attribute.String("event.type", "user-action"),
)

// Record error
if err != nil {
    helper.RecordError(span, err)
}
```

### Using the Tracer

```go
tracer := otel.GetTracer()
if tracer != nil {
    ctx, span := tracer.Start(ctx, "operation")
    defer span.End()
    // ...
}
```

## Components

- `config.go`: Configuration structures
- `otel.go`: Initialization and lifecycle management
- `middleware.go`: HTTP middleware for automatic tracing
- `span.go`: Helper utilities for span operations
- `convert.go`: Convert from AppConfig to OTEL Config

## Environment Variables

All configuration can be set via environment variables:

```bash
export MCP_OTEL_ENABLED=true
export MCP_OTEL_ENDPOINT=localhost:4317
export MCP_OTEL_TRACING_ENABLED=true
export MCP_OTEL_TRACING_SAMPLE_RATE=0.1
export MCP_OTEL_METRICS_ENABLED=true
```

## Local Development

### Using OpenTelemetry Collector

Create `otel-collector-config.yaml`:

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: localhost:4317

processors:
  batch:

exporters:
  logging:
    loglevel: debug

  prometheusremotewrite:
    endpoint: http://localhost:9090/api/v1/write

  jaeger:
    endpoint: jaeger:4317
    tls:
      insecure: true

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [jaeger, logging]

    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [prometheusremotewrite, logging]
```

Run the collector:

```bash
docker run -d --name otel-collector \
  -p 4317:4317 \
  -v $(pwd)/otel-collector-config.yaml:/etc/otelcol-contrib/config.yaml \
  otel/opentelemetry-collector-contrib:latest
```

## Production Considerations

1. **Sampling**: Use appropriate sampling rates (0.1 for production, 1.0 for dev)
2. **TLS**: Enable TLS in production (`insecure: false`)
3. **Batching**: Tune batch sizes based on traffic volume
4. **Endpoint**: Use a dedicated OTLP collector instance
5. **Backends**: Route to Jaeger/Tempo for traces, Prometheus for metrics

## Testing

Run tests:

```bash
go test ./internal/observability/otel/...
```

## Troubleshooting

### No traces appearing
- Check if `otel.enabled: true` in config
- Verify OTLP collector endpoint is reachable
- Check sampling rate (might be too low)
- Verify tracer provider is initialized

### Metrics not exporting
- Check `otel.metrics.enabled: true`
- Verify export interval is reasonable
- Check OTLP collector metrics pipeline configuration

### Performance impact
- Reduce sampling rate for high traffic
- Increase batch timeout for better batching
- Consider using delta temporality for lower cardinality

## Related Services

- **OpenTelemetry Service**: Query external OpenTelemetry Collector data
- **Jaeger Service**: Query distributed traces
- **Prometheus Service**: Query metrics data
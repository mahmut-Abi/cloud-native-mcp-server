package otel

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Middleware returns an HTTP middleware that adds OpenTelemetry tracing
func Middleware(serviceName string, opts ...otelhttp.Option) func(http.Handler) http.Handler {
	// Set default options
	defaultOpts := []otelhttp.Option{
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return r.Method + " " + r.URL.Path
		}),
	}

	opts = append(defaultOpts, opts...)

	return func(next http.Handler) http.Handler {
		// Wrap the next handler with otelhttp
		return otelhttp.NewHandler(next, serviceName, opts...)
	}
}

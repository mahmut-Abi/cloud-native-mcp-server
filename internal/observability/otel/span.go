package otel

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SpanHelper provides helper functions for working with spans
type SpanHelper struct {
	tracer trace.Tracer
}

// NewSpanHelper creates a new SpanHelper
func NewSpanHelper() *SpanHelper {
	return &SpanHelper{
		tracer: GetTracer(),
	}
}

// StartSpan starts a new span
func (h *SpanHelper) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	if h.tracer == nil {
		return ctx, trace.SpanFromContext(ctx)
	}
	return h.tracer.Start(ctx, name, opts...)
}

// RecordError records an error in the span
func (h *SpanHelper) RecordError(span trace.Span, err error, attrs ...attribute.KeyValue) {
	if span == nil {
		return
	}

	span.SetStatus(codes.Error, err.Error())
	span.SetAttributes(attrs...)
	span.RecordError(err)
}

// SetAttributes sets attributes on the span
func (h *SpanHelper) SetAttributes(span trace.Span, attrs ...attribute.KeyValue) {
	if span == nil {
		return
	}
	span.SetAttributes(attrs...)
}

// AddEvent adds an event to the span
func (h *SpanHelper) AddEvent(span trace.Span, name string, attrs ...attribute.KeyValue) {
	if span == nil {
		return
	}
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// WithSpan runs a function within a span
func WithSpan(ctx context.Context, name string, fn func(ctx context.Context, span trace.Span) error, opts ...trace.SpanStartOption) error {
	if GetTracer() == nil {
		return fn(ctx, nil)
	}

	ctx, span := GetTracer().Start(ctx, name, opts...)
	defer span.End()

	err := fn(ctx, span)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return err
}

// WithSpanAsync runs a function within a span asynchronously
func WithSpanAsync(ctx context.Context, name string, fn func(ctx context.Context, span trace.Span), opts ...trace.SpanStartOption) {
	if GetTracer() == nil {
		fn(ctx, nil)
		return
	}

	go func() {
		ctx, span := GetTracer().Start(ctx, name, opts...)
		defer span.End()

		fn(ctx, span)
	}()
}
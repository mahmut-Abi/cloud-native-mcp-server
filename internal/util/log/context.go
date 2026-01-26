package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

type contextKey string

const (
	traceIDKey contextKey = "trace_id"
	spanIDKey  contextKey = "span_id"
	userIDKey  contextKey = "user_id"
)

func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, traceIDKey, id)
}

func WithSpanID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, spanIDKey, id)
}

func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

func Fields(ctx context.Context) logrus.Fields {
	fields := logrus.Fields{}
	if v := ctx.Value(traceIDKey); v != nil {
		fields["trace_id"] = v
	}
	if v := ctx.Value(spanIDKey); v != nil {
		fields["span_id"] = v
	}
	if v := ctx.Value(userIDKey); v != nil {
		fields["user_id"] = v
	}
	return fields
}

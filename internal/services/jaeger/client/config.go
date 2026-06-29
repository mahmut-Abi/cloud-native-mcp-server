// Package client provides Jaeger HTTP API client functionality.
package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/middleware"
)

const (
	hdrAddress    = "X-Mcp-Backend-Jaeger-Address"
	hdrTimeoutSec = "X-Mcp-Backend-Jaeger-Timeout-Sec"
)

type jaegerContextKey struct{}

func init() {
	middleware.RegisterBackendAuthHandler("jaeger", parseHeadersAndInjectClient)
}

func parseHeadersAndInjectClient(r *http.Request) (*http.Request, error) {
	opts := parseRequestHeaders(r.Header)
	if opts.BaseURL == "" {
		return r, fmt.Errorf("no jaeger address in headers")
	}
	cli, err := NewClient(opts)
	if err != nil {
		return r, err
	}
	ctx := context.WithValue(r.Context(), jaegerContextKey{}, cli)
	return r.WithContext(ctx), nil
}

func parseRequestHeaders(h http.Header) *ClientOptions {
	opts := &ClientOptions{Timeout: 30 * time.Second}
	if v := h.Get(hdrAddress); v != "" {
		opts.BaseURL = v
	}
	if v := h.Get(hdrTimeoutSec); v != "" {
		if sec, err := strconv.Atoi(v); err == nil && sec > 0 {
			opts.Timeout = time.Duration(sec) * time.Second
		}
	}
	return opts
}

// NewContext returns a new context with the given Jaeger client injected.
func NewContext(ctx context.Context, cli *Client) context.Context {
	return context.WithValue(ctx, jaegerContextKey{}, cli)
}

// FromContext extracts the Jaeger client from the request context.
// Returns an error if no client was injected by the backend auth middleware.
func FromContext(ctx context.Context) (*Client, error) {
	cli, ok := ctx.Value(jaegerContextKey{}).(*Client)
	if !ok || cli == nil {
		return nil, fmt.Errorf("jaeger client not found in context")
	}
	return cli, nil
}

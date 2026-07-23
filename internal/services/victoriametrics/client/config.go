// Package client provides VictoriaMetrics HTTP API client functionality.
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
	hdrURL        = "X-Mcp-Backend-VictoriaMetrics-Url"
	hdrToken      = "X-Mcp-Backend-VictoriaMetrics-Token"
	hdrUsername   = "X-Mcp-Backend-VictoriaMetrics-Username"
	hdrPassword   = "X-Mcp-Backend-VictoriaMetrics-Password"
	hdrTLSSkip    = "X-Mcp-Backend-VictoriaMetrics-Tls-Skip-Verify"
	hdrTimeoutSec = "X-Mcp-Backend-VictoriaMetrics-Timeout-Sec"
)

type victoriaMetricsContextKey struct{}

func init() {
	middleware.RegisterBackendAuthHandler("victoriametrics", parseHeadersAndInjectClient)
}

func parseHeadersAndInjectClient(r *http.Request) (*http.Request, error) {
	opts := parseRequestHeaders(r.Header)
	if opts.Address == "" {
		return r, fmt.Errorf("no VictoriaMetrics URL in headers")
	}
	cli, err := NewClient(opts)
	if err != nil {
		return r, err
	}
	ctx := context.WithValue(r.Context(), victoriaMetricsContextKey{}, cli)
	return r.WithContext(ctx), nil
}

func parseRequestHeaders(h http.Header) *ClientOptions {
	opts := &ClientOptions{Timeout: 30 * time.Second}
	if v := h.Get(hdrURL); v != "" {
		opts.Address = v
	}
	if v := h.Get(hdrToken); v != "" {
		opts.BearerToken = v
	}
	if v := h.Get(hdrUsername); v != "" {
		opts.Username = v
	}
	if v := h.Get(hdrPassword); v != "" {
		opts.Password = v
	}
	if v := h.Get(hdrTLSSkip); v != "" {
		opts.TLSSkipVerify, _ = strconv.ParseBool(v)
	}
	if v := h.Get(hdrTimeoutSec); v != "" {
		if sec, err := strconv.Atoi(v); err == nil && sec > 0 {
			opts.Timeout = time.Duration(sec) * time.Second
		}
	}
	return opts
}

// FromContext extracts the VictoriaMetrics client from the request context.
// Returns an error if no client was injected by the backend auth middleware.
func FromContext(ctx context.Context) (*Client, error) {
	cli, ok := ctx.Value(victoriaMetricsContextKey{}).(*Client)
	if !ok || cli == nil {
		return nil, fmt.Errorf("victoriametrics client not found in context")
	}
	return cli, nil
}

// NewContext returns a new context with the given VictoriaMetrics client injected.
func NewContext(ctx context.Context, cli *Client) context.Context {
	return context.WithValue(ctx, victoriaMetricsContextKey{}, cli)
}

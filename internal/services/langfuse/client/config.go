// Package client provides Langfuse HTTP API client functionality.
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
	hdrURL        = "X-Mcp-Backend-Langfuse-Url"
	hdrUsername   = "X-Mcp-Backend-Langfuse-Username"
	hdrPassword   = "X-Mcp-Backend-Langfuse-Password"
	hdrTimeoutSec = "X-Mcp-Backend-Langfuse-Timeout-Sec"
	hdrProjectID  = "X-Mcp-Backend-Langfuse-Project-Id"
)

type langfuseContextKey struct{}

func init() {
	middleware.RegisterBackendAuthHandler("langfuse", parseHeadersAndInjectClient)
}

func parseHeadersAndInjectClient(r *http.Request) (*http.Request, error) {
	opts := parseRequestHeaders(r.Header)
	if opts.URL == "" {
		return r, fmt.Errorf("no langfuse URL in headers")
	}
	cli, err := NewClient(opts)
	if err != nil {
		return r, err
	}
	ctx := context.WithValue(r.Context(), langfuseContextKey{}, cli)
	return r.WithContext(ctx), nil
}

func parseRequestHeaders(h http.Header) *ClientOptions {
	opts := &ClientOptions{Timeout: 30 * time.Second}
	if v := h.Get(hdrURL); v != "" {
		opts.URL = v
	}
	if v := h.Get(hdrUsername); v != "" {
		opts.Username = v
	}
	if v := h.Get(hdrPassword); v != "" {
		opts.Password = v
	}
	if v := h.Get(hdrTimeoutSec); v != "" {
		if sec, err := strconv.Atoi(v); err == nil && sec > 0 {
			opts.Timeout = time.Duration(sec) * time.Second
		}
	}
	if v := h.Get(hdrProjectID); v != "" {
		opts.ProjectID = v
	}
	return opts
}

// NewContext returns a new context with the given Langfuse client injected.
func NewContext(ctx context.Context, cli *Client) context.Context {
	return context.WithValue(ctx, langfuseContextKey{}, cli)
}

// FromContext extracts the Langfuse client from the request context.
// Returns an error if no client was injected by the backend auth middleware.
func FromContext(ctx context.Context) (*Client, error) {
	cli, ok := ctx.Value(langfuseContextKey{}).(*Client)
	if !ok || cli == nil {
		return nil, fmt.Errorf("langfuse client not found in context")
	}
	return cli, nil
}

// Package client provides Sentry HTTP API client functionality.
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
	hdrURL        = "X-Mcp-Backend-Sentry-Url"
	hdrAuthToken  = "X-Mcp-Backend-Sentry-Auth-Token"
	hdrOrg        = "X-Mcp-Backend-Sentry-Organization"
	hdrProject    = "X-Mcp-Backend-Sentry-Project"
	hdrTimeoutSec = "X-Mcp-Backend-Sentry-Timeout-Sec"
)

type sentryContextKey struct{}

func init() {
	middleware.RegisterBackendAuthHandler("sentry", parseHeadersAndInjectClient)
}

func parseHeadersAndInjectClient(r *http.Request) (*http.Request, error) {
	opts := parseRequestHeaders(r.Header)
	if opts.URL == "" {
		return r, fmt.Errorf("no sentry URL in headers")
	}
	cli, err := NewClient(opts)
	if err != nil {
		return r, err
	}
	ctx := context.WithValue(r.Context(), sentryContextKey{}, cli)
	return r.WithContext(ctx), nil
}

func parseRequestHeaders(h http.Header) *ClientOptions {
	opts := &ClientOptions{Timeout: 30 * time.Second}
	if v := h.Get(hdrURL); v != "" {
		opts.URL = v
	}
	if v := h.Get(hdrAuthToken); v != "" {
		opts.AuthToken = v
	}
	if v := h.Get(hdrTimeoutSec); v != "" {
		if sec, err := strconv.Atoi(v); err == nil && sec > 0 {
			opts.Timeout = time.Duration(sec) * time.Second
		}
	}
	return opts
}

// FromContext extracts the Sentry client from the request context.
// Returns an error if no client was injected by the backend auth middleware.
func FromContext(ctx context.Context) (*Client, error) {
	cli, ok := ctx.Value(sentryContextKey{}).(*Client)
	if !ok || cli == nil {
		return nil, fmt.Errorf("sentry client not found in context")
	}
	return cli, nil
}

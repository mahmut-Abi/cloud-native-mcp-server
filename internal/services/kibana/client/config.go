// Package client provides Kibana HTTP API client functionality.
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
	hdrKibanaURL        = "X-Mcp-Backend-Kibana-Url"
	hdrKibanaAPIKey     = "X-Mcp-Backend-Kibana-Api-Key"
	hdrKibanaUsername   = "X-Mcp-Backend-Kibana-Username"
	hdrKibanaPassword   = "X-Mcp-Backend-Kibana-Password"
	hdrKibanaSpace      = "X-Mcp-Backend-Kibana-Space"
	hdrKibanaSkipVerify = "X-Mcp-Backend-Kibana-Skip-Verify"
	hdrKibanaTimeoutSec = "X-Mcp-Backend-Kibana-Timeout-Sec"
)

type kibanaContextKey struct{}

func init() {
	middleware.RegisterBackendAuthHandler("kibana", parseHeadersAndInjectClient)
}

func parseHeadersAndInjectClient(r *http.Request) (*http.Request, error) {
	opts := parseRequestHeaders(r.Header)
	if opts.URL == "" {
		return r, fmt.Errorf("no kibana URL in headers")
	}
	cli, err := NewClient(opts)
	if err != nil {
		return r, err
	}
	ctx := context.WithValue(r.Context(), kibanaContextKey{}, cli)
	return r.WithContext(ctx), nil
}

func parseRequestHeaders(h http.Header) *ClientOptions {
	opts := &ClientOptions{Timeout: 30 * time.Second}
	if v := h.Get(hdrKibanaURL); v != "" {
		opts.URL = v
	}
	if v := h.Get(hdrKibanaAPIKey); v != "" {
		opts.APIKey = v
	}
	if v := h.Get(hdrKibanaUsername); v != "" {
		opts.Username = v
	}
	if v := h.Get(hdrKibanaPassword); v != "" {
		opts.Password = v
	}
	if v := h.Get(hdrKibanaSpace); v != "" {
		opts.Space = v
	}
	if v := h.Get(hdrKibanaSkipVerify); v != "" {
		opts.SkipVerify, _ = strconv.ParseBool(v)
	}
	if v := h.Get(hdrKibanaTimeoutSec); v != "" {
		if sec, err := strconv.Atoi(v); err == nil && sec > 0 {
			opts.Timeout = time.Duration(sec) * time.Second
		}
	}
	return opts
}

// FromContext extracts the Kibana client from the request context.
// Returns an error if no client was injected by the backend auth middleware.
func FromContext(ctx context.Context) (*Client, error) {
	cli, ok := ctx.Value(kibanaContextKey{}).(*Client)
	if !ok || cli == nil {
		return nil, fmt.Errorf("kibana client not found in context")
	}
	return cli, nil
}

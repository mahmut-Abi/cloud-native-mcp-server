// Package client provides Grafana HTTP API client functionality.
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
	hdrGrafanaURL        = "X-Mcp-Backend-Grafana-Url"
	hdrGrafanaAPIKey     = "X-Mcp-Backend-Grafana-Api-Key"
	hdrGrafanaUsername   = "X-Mcp-Backend-Grafana-Username"
	hdrGrafanaPassword   = "X-Mcp-Backend-Grafana-Password"
	hdrGrafanaTimeoutSec = "X-Mcp-Backend-Grafana-Timeout-Sec"
)

type grafanaContextKey struct{}

func init() {
	middleware.RegisterBackendAuthHandler("grafana", parseHeadersAndInjectClient)
}

func parseHeadersAndInjectClient(r *http.Request) (*http.Request, error) {
	opts := parseRequestHeaders(r.Header)
	if opts.URL == "" {
		return r, fmt.Errorf("no grafana URL in headers")
	}
	cli, err := NewClient(opts)
	if err != nil {
		return r, err
	}
	ctx := context.WithValue(r.Context(), grafanaContextKey{}, cli)
	return r.WithContext(ctx), nil
}

func parseRequestHeaders(h http.Header) *ClientOptions {
	opts := &ClientOptions{Timeout: 30 * time.Second}
	if v := h.Get(hdrGrafanaURL); v != "" {
		opts.URL = v
	}
	if v := h.Get(hdrGrafanaAPIKey); v != "" {
		opts.APIKey = v
	}
	if v := h.Get(hdrGrafanaUsername); v != "" {
		opts.Username = v
	}
	if v := h.Get(hdrGrafanaPassword); v != "" {
		opts.Password = v
	}
	if v := h.Get(hdrGrafanaTimeoutSec); v != "" {
		if sec, err := strconv.Atoi(v); err == nil && sec > 0 {
			opts.Timeout = time.Duration(sec) * time.Second
		}
	}
	return opts
}

// FromContext extracts the Grafana client from the request context.
// Returns an error if no client was injected by the backend auth middleware.
func FromContext(ctx context.Context) (*Client, error) {
	cli, ok := ctx.Value(grafanaContextKey{}).(*Client)
	if !ok || cli == nil {
		return nil, fmt.Errorf("grafana client not found in context")
	}
	return cli, nil
}

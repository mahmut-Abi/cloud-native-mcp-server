// Package client provides Alertmanager HTTP API client functionality.
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
	hdrURL        = "X-Mcp-Backend-Alertmanager-Url"
	hdrToken      = "X-Mcp-Backend-Alertmanager-Token"
	hdrUsername   = "X-Mcp-Backend-Alertmanager-Username"
	hdrPassword   = "X-Mcp-Backend-Alertmanager-Password"
	hdrTLSSkip    = "X-Mcp-Backend-Alertmanager-Tls-Skip-Verify"
	hdrTimeoutSec = "X-Mcp-Backend-Alertmanager-Timeout-Sec"
)

type alertmanagerContextKey struct{}

func init() {
	middleware.RegisterBackendAuthHandler("alertmanager", parseHeadersAndInjectClient)
}

func parseHeadersAndInjectClient(r *http.Request) (*http.Request, error) {
	opts := parseRequestHeaders(r.Header)
	if opts.Address == "" {
		return r, fmt.Errorf("no alertmanager URL in headers")
	}
	cli, err := NewClientWithOptions(opts)
	if err != nil {
		return r, err
	}
	ctx := context.WithValue(r.Context(), alertmanagerContextKey{}, cli)
	return r.WithContext(ctx), nil
}

func parseRequestHeaders(h http.Header) *ClientOptions {
	opts := DefaultClientOptions()
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

// FromContext extracts the Alertmanager client from the request context.
// Returns an error if no client was injected by the backend auth middleware.
func FromContext(ctx context.Context) (*Client, error) {
	cli, ok := ctx.Value(alertmanagerContextKey{}).(*Client)
	if !ok || cli == nil {
		return nil, fmt.Errorf("alertmanager client not found in context")
	}
	return cli, nil
}

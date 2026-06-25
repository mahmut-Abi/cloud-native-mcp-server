// Package client provides ArgoCD HTTP API client functionality.
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
	hdrURL        = "X-Mcp-Backend-Argocd-Url"
	hdrUsername   = "X-Mcp-Backend-Argocd-Username"
	hdrPassword   = "X-Mcp-Backend-Argocd-Password"
	hdrAuthToken  = "X-Mcp-Backend-Argocd-Auth-Token"
	hdrTimeoutSec = "X-Mcp-Backend-Argocd-Timeout-Sec"
)

type argocdContextKey struct{}

func init() {
	middleware.RegisterBackendAuthHandler("argocd", parseHeadersAndInjectClient)
}

func parseHeadersAndInjectClient(r *http.Request) (*http.Request, error) {
	opts := parseRequestHeaders(r.Header)
	if opts.URL == "" {
		return r, fmt.Errorf("no argocd URL in headers")
	}
	cli, err := NewClient(opts)
	if err != nil {
		return r, err
	}
	ctx := context.WithValue(r.Context(), argocdContextKey{}, cli)
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

// FromContext extracts the ArgoCD client from the request context.
// Returns an error if no client was injected by the backend auth middleware.
func FromContext(ctx context.Context) (*Client, error) {
	cli, ok := ctx.Value(argocdContextKey{}).(*Client)
	if !ok || cli == nil {
		return nil, fmt.Errorf("argocd client not found in context")
	}
	return cli, nil
}

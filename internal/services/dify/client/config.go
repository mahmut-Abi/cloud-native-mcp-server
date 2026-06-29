// Package client provides Dify HTTP API client functionality.
// It supports both the Dify Console API (session-based auth with email/password)
// and the Dify Service API (Bearer token auth).
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
	hdrConsoleURL     = "X-Mcp-Backend-Dify-Console-Url"
	hdrConsoleEmail   = "X-Mcp-Backend-Dify-Console-Email"
	hdrConsolePass    = "X-Mcp-Backend-Dify-Console-Password"
	hdrServiceURL     = "X-Mcp-Backend-Dify-Service-Url"
	hdrAPIKey         = "X-Mcp-Backend-Dify-Api-Key"
	hdrTimeoutSec     = "X-Mcp-Backend-Dify-Timeout-Sec"
)

type difyContextKey struct{}

func init() {
	middleware.RegisterBackendAuthHandler("dify", parseHeadersAndInjectClient)
}

func parseHeadersAndInjectClient(r *http.Request) (*http.Request, error) {
	opts := parseRequestHeaders(r.Header)
	if opts.ConsoleURL == "" && opts.ServiceURL == "" {
		return r, fmt.Errorf("no dify console URL or service URL in headers")
	}
	if opts.ConsoleURL != "" && (opts.ConsoleEmail == "" || opts.ConsolePassword == "") {
		return r, fmt.Errorf("dify console URL provided but missing email or password")
	}
	if opts.ServiceURL != "" && opts.APIKey == "" {
		return r, fmt.Errorf("dify service URL provided but missing API key")
	}
	cli, err := NewClient(opts)
	if err != nil {
		return r, err
	}
	ctx := context.WithValue(r.Context(), difyContextKey{}, cli)
	return r.WithContext(ctx), nil
}

func parseRequestHeaders(h http.Header) *ClientOptions {
	opts := &ClientOptions{Timeout: 30 * time.Second}
	if v := h.Get(hdrConsoleURL); v != "" {
		opts.ConsoleURL = v
	}
	if v := h.Get(hdrConsoleEmail); v != "" {
		opts.ConsoleEmail = v
	}
	if v := h.Get(hdrConsolePass); v != "" {
		opts.ConsolePassword = v
	}
	if v := h.Get(hdrServiceURL); v != "" {
		opts.ServiceURL = v
	}
	if v := h.Get(hdrAPIKey); v != "" {
		opts.APIKey = v
	}
	if v := h.Get(hdrTimeoutSec); v != "" {
		if sec, err := strconv.Atoi(v); err == nil && sec > 0 {
			opts.Timeout = time.Duration(sec) * time.Second
		}
	}
	return opts
}

// NewContext stores a Dify client in the given context so it can be
// retrieved later via FromContext.
func NewContext(ctx context.Context, cli *Client) context.Context {
	return context.WithValue(ctx, difyContextKey{}, cli)
}

// FromContext extracts the Dify client from the request context.
// Returns an error if no client was injected by the backend auth middleware.
func FromContext(ctx context.Context) (*Client, error) {
	cli, ok := ctx.Value(difyContextKey{}).(*Client)
	if !ok || cli == nil {
		return nil, fmt.Errorf("dify client not found in context")
	}
	return cli, nil
}

// Package client provides Elasticsearch HTTP client operations for the MCP server.
package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/middleware"
)

const (
	hdrAddresses    = "X-Mcp-Backend-Elasticsearch-Addresses"
	hdrUsername     = "X-Mcp-Backend-Elasticsearch-Username"
	hdrPassword     = "X-Mcp-Backend-Elasticsearch-Password"
	hdrBearerToken  = "X-Mcp-Backend-Elasticsearch-Bearer-Token"
	hdrAPIKey       = "X-Mcp-Backend-Elasticsearch-Api-Key"
	hdrTLSSkip      = "X-Mcp-Backend-Elasticsearch-Tls-Skip-Verify"
	hdrTimeoutSec   = "X-Mcp-Backend-Elasticsearch-Timeout-Sec"
)

type elasticsearchContextKey struct{}

func init() {
	middleware.RegisterBackendAuthHandler("elasticsearch", parseHeadersAndInjectClient)
}

func parseHeadersAndInjectClient(r *http.Request) (*http.Request, error) {
	opts := parseRequestHeaders(r.Header)
	if len(opts.Addresses) == 0 {
		return r, fmt.Errorf("no Elasticsearch addresses in headers")
	}
	cli, err := NewClient(opts)
	if err != nil {
		return r, err
	}
	ctx := context.WithValue(r.Context(), elasticsearchContextKey{}, cli)
	return r.WithContext(ctx), nil
}

func parseRequestHeaders(h http.Header) *ClientOptions {
	opts := &ClientOptions{Timeout: 30 * time.Second}
	if v := h.Get(hdrAddresses); v != "" {
		parts := strings.Split(v, ",")
		addrs := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				addrs = append(addrs, p)
			}
		}
		opts.Addresses = addrs
	}
	if v := h.Get(hdrUsername); v != "" {
		opts.Username = v
	}
	if v := h.Get(hdrPassword); v != "" {
		opts.Password = v
	}
	if v := h.Get(hdrBearerToken); v != "" {
		opts.BearerToken = v
	}
	if v := h.Get(hdrAPIKey); v != "" {
		opts.APIKey = v
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

// FromContext extracts the Elasticsearch client from the request context.
// Returns an error if no client was injected by the backend auth middleware.
func FromContext(ctx context.Context) (*Client, error) {
	cli, ok := ctx.Value(elasticsearchContextKey{}).(*Client)
	if !ok || cli == nil {
		return nil, fmt.Errorf("elasticsearch client not found in context")
	}
	return cli, nil
}

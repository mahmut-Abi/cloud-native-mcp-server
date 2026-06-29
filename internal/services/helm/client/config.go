// Package client provides Helm client operations for the MCP server.
package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/middleware"
)

const (
	hdrKubeconfigPath = "X-Mcp-Backend-Helm-Kubeconfig-Path"
	hdrNamespace      = "X-Mcp-Backend-Helm-Namespace"
	hdrHTTPProxy      = "X-Mcp-Backend-Helm-Http-Proxy"
	hdrDebug          = "X-Mcp-Backend-Helm-Debug"
	hdrTimeoutSec     = "X-Mcp-Backend-Helm-Timeout-Sec"
)

type helmContextKey struct{}

func init() {
	middleware.RegisterBackendAuthHandler("helm", parseHeadersAndInjectClient)
}

func parseHeadersAndInjectClient(r *http.Request) (*http.Request, error) {
	opts := parseRequestHeaders(r.Header)
	cli, err := NewClient(opts)
	if err != nil {
		return r, err
	}
	ctx := context.WithValue(r.Context(), helmContextKey{}, cli)
	return r.WithContext(ctx), nil
}

func parseRequestHeaders(h http.Header) *ClientOptions {
	opts := &ClientOptions{}
	timeoutSec := 300 // default

	if v := h.Get(hdrKubeconfigPath); v != "" {
		if content, isBase64 := tryDecodeKubeconfig(v); isBase64 {
			opts.KubeconfigPath = writeTempKubeconfig(content)
		} else if _, err := os.Stat(v); err == nil {
			opts.KubeconfigPath = v
		} else {
			opts.KubeconfigPath = writeTempKubeconfig(v)
		}
	}
	if v := h.Get(hdrNamespace); v != "" {
		opts.Namespace = v
	}
	if v := h.Get(hdrHTTPProxy); v != "" {
		opts.HTTPProxy = v
	}
	if v := h.Get(hdrDebug); v != "" {
		opts.Debug = v == "true"
	}
	if v := h.Get(hdrTimeoutSec); v != "" {
		if sec, err := strconv.Atoi(v); err == nil && sec > 0 {
			timeoutSec = sec
		}
	}

	opts.Optimizer = NewRepositoryOptimizer(timeoutSec, 3, opts.HTTPProxy)
	return opts
}

func tryDecodeKubeconfig(v string) (string, bool) {
	decoded, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return "", false
	}
	result := string(decoded)
	if strings.Contains(result, "apiVersion") || strings.Contains(result, "kind:") {
		return result, true
	}
	return "", false
}

func writeTempKubeconfig(content string) string {
	tmpFile, err := os.CreateTemp("", "mcp-helm-kubeconfig-*.yaml")
	if err != nil {
		return ""
	}
	if _, err := tmpFile.WriteString(content); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
		return ""
	}
	_ = tmpFile.Close()
	return tmpFile.Name()
}

// FromContext extracts the Helm client from the request context.
// Returns an error if no client was injected by the backend auth middleware.
func FromContext(ctx context.Context) (*Client, error) {
	cli, ok := ctx.Value(helmContextKey{}).(*Client)
	if !ok || cli == nil {
		return nil, fmt.Errorf("helm client not found in context")
	}
	return cli, nil
}

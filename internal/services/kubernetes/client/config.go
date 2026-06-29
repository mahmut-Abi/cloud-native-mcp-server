// Package client provides Kubernetes API client functionality.
package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/middleware"
)

const (
	hdrKubeconfig = "X-Mcp-Backend-Kubernetes-Kubeconfig"
	hdrQPS        = "X-Mcp-Backend-Kubernetes-Qps"
	hdrBurst      = "X-Mcp-Backend-Kubernetes-Burst"
	hdrTimeoutSec = "X-Mcp-Backend-Kubernetes-Timeout-Sec"
)

type kubernetesContextKey struct{}

func init() {
	middleware.RegisterBackendAuthHandler("kubernetes", parseHeadersAndInjectClient)
}

func parseHeadersAndInjectClient(r *http.Request) (*http.Request, error) {
	opts := parseRequestHeaders(r.Header)
	if opts.KubeconfigPath == "" {
		// Try in-cluster config or default kubeconfig - still create a client
		opts.KubeconfigPath = os.Getenv("KUBECONFIG")
	}
	cli, err := NewClientWithOptions(opts)
	if err != nil {
		return r, err
	}
	ctx := context.WithValue(r.Context(), kubernetesContextKey{}, cli)
	return r.WithContext(ctx), nil
}

func parseRequestHeaders(h http.Header) *ClientOptions {
	opts := DefaultClientOptions()
	if v := h.Get(hdrKubeconfig); v != "" {
		if content, isBase64 := tryDecode(v); isBase64 {
			opts.KubeconfigPath = writeTempKubeconfig(content)
		} else if _, err := os.Stat(v); err == nil {
			opts.KubeconfigPath = v
		} else {
			opts.KubeconfigPath = writeTempKubeconfig(v)
		}
	}
	if v := h.Get(hdrQPS); v != "" {
		if qps, err := strconv.ParseFloat(v, 32); err == nil && qps > 0 {
			opts.QPS = float32(qps)
		}
	}
	if v := h.Get(hdrBurst); v != "" {
		if burst, err := strconv.Atoi(v); err == nil && burst > 0 {
			opts.Burst = burst
		}
	}
	if v := h.Get(hdrTimeoutSec); v != "" {
		if sec, err := strconv.Atoi(v); err == nil && sec > 0 {
			opts.Timeout = time.Duration(sec) * time.Second
		}
	}
	return opts
}

func tryDecode(v string) (string, bool) {
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
	tmpFile, err := os.CreateTemp("", "mcp-kubeconfig-*.yaml")
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

// FromContext extracts the Kubernetes client from the request context.
// Returns an error if no client was injected by the backend auth middleware.
func FromContext(ctx context.Context) (*Client, error) {
	cli, ok := ctx.Value(kubernetesContextKey{}).(*Client)
	if !ok || cli == nil {
		return nil, fmt.Errorf("kubernetes client not found in context")
	}
	return cli, nil
}

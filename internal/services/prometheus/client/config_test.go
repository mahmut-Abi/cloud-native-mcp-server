package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/middleware"
)

func TestParseRequestHeaders(t *testing.T) {
	h := http.Header{}
	h.Set(hdrURL, "https://prometheus.example.com:9090")
	h.Set(hdrToken, "bearer-token-123")
	h.Set(hdrUsername, "admin")
	h.Set(hdrPassword, "secret")
	h.Set(hdrTLSSkip, "true")
	h.Set(hdrTimeoutSec, "60")

	opts := parseRequestHeaders(h)

	if opts.Address != "https://prometheus.example.com:9090" {
		t.Errorf("expected address, got %q", opts.Address)
	}
	if opts.BearerToken != "bearer-token-123" {
		t.Errorf("expected bearer token, got %q", opts.BearerToken)
	}
	if opts.Username != "admin" {
		t.Errorf("expected username, got %q", opts.Username)
	}
	if opts.Password != "secret" {
		t.Errorf("expected password, got %q", opts.Password)
	}
	if !opts.TLSSkipVerify {
		t.Error("expected TLS skip verify to be true")
	}
	if opts.Timeout.String() != "1m0s" {
		t.Errorf("expected 60s timeout, got %v", opts.Timeout)
	}
}

func TestParseRequestHeadersMinimal(t *testing.T) {
	h := http.Header{}
	h.Set(hdrURL, "http://localhost:9090")

	opts := parseRequestHeaders(h)

	if opts.Address != "http://localhost:9090" {
		t.Errorf("expected address, got %q", opts.Address)
	}
	if opts.BearerToken != "" {
		t.Error("expected no bearer token")
	}
	if opts.Username != "" {
		t.Error("expected no username")
	}
	if opts.TLSSkipVerify {
		t.Error("expected TLS skip verify to be false")
	}
	if opts.Timeout.Seconds() != 30 {
		t.Errorf("expected 30s default timeout, got %v", opts.Timeout)
	}
}

func TestParseRequestHeadersEmpty(t *testing.T) {
	opts := parseRequestHeaders(http.Header{})

	if opts.Address != "" {
		t.Errorf("expected empty address, got %q", opts.Address)
	}
	if opts.Timeout.Seconds() != 30 {
		t.Errorf("expected 30s default timeout, got %v", opts.Timeout)
	}
}

func TestParseRequestHeadersInvalidTimeout(t *testing.T) {
	h := http.Header{}
	h.Set(hdrURL, "http://localhost:9090")
	h.Set(hdrTimeoutSec, "not-a-number")

	opts := parseRequestHeaders(h)

	if opts.Timeout.Seconds() != 30 {
		t.Errorf("expected 30s default timeout for invalid input, got %v", opts.Timeout)
	}
}

func TestNewContextAndFromContext(t *testing.T) {
	cli, err := NewClient(&ClientOptions{Address: "http://localhost:9090"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()
	ctx = NewContext(ctx, cli)

	retrieved, err := FromContext(ctx)
	if err != nil {
		t.Fatalf("FromContext failed: %v", err)
	}
	if retrieved != cli {
		t.Error("FromContext returned different client")
	}
}

func TestFromContextMissing(t *testing.T) {
	_, err := FromContext(context.Background())
	if err == nil {
		t.Error("expected error for missing client in context")
	}
}

func TestFromContextNil(t *testing.T) {
	ctx := context.WithValue(context.Background(), prometheusContextKey{}, (*Client)(nil))
	_, err := FromContext(ctx)
	if err == nil {
		t.Error("expected error for nil client in context")
	}
}

func TestParseHeadersAndInjectClient(t *testing.T) {
	h := http.Header{}
	h.Set(hdrURL, "http://localhost:9090")
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header = h

	newReq, err := parseHeadersAndInjectClient(req)
	if err != nil {
		t.Fatalf("parseHeadersAndInjectClient failed: %v", err)
	}

	cli, err := FromContext(newReq.Context())
	if err != nil {
		t.Fatalf("client not in context after injection: %v", err)
	}
	if cli == nil {
		t.Fatal("client is nil")
	}
}

func TestParseHeadersAndInjectClientMissingURL(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	_, err := parseHeadersAndInjectClient(req)
	if err == nil {
		t.Error("expected error for missing URL")
	}
}

func TestInitHandlerRegistered(t *testing.T) {
	var called bool
	handler := middleware.BackendAuthMiddleware("prometheus")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		_, err := FromContext(r.Context())
		if err != nil {
			t.Logf("client not found in context (expected without valid headers): %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))

	h := http.Header{}
	h.Set(hdrURL, "http://localhost:9090")
	req := httptest.NewRequest("GET", "/api/prometheus/sse/message", nil)
	req.Header = h
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !called {
		t.Error("handler was not called, init() registration may have failed")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

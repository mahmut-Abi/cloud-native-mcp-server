package config

import (
	"strings"
	"testing"
)

func TestAuthOIDCConfigFromEnv(t *testing.T) {
	t.Setenv("MCP_AUTH_ENABLED", "true")
	t.Setenv("MCP_AUTH_MODE", "bearer")
	t.Setenv("MCP_AUTH_BEARER_TOKEN", "")
	t.Setenv("MCP_AUTH_OIDC_ISSUER_URL", "https://issuer.example.com/realms/dev")
	t.Setenv("MCP_AUTH_OIDC_AUDIENCE", "mcp-client")
	t.Setenv("MCP_AUTH_OIDC_CLIENT_ID", "mcp-client-id")
	t.Setenv("MCP_AUTH_OIDC_HTTP_TIMEOUT", "7")
	t.Setenv("MCP_AUTH_OIDC_JWKS_CACHE_TTL", "120")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if !cfg.Auth.Enabled {
		t.Fatal("expected auth to be enabled")
	}
	if cfg.Auth.Mode != "bearer" {
		t.Fatalf("expected auth mode bearer, got %q", cfg.Auth.Mode)
	}
	if cfg.Auth.OIDCIssuerURL != "https://issuer.example.com/realms/dev" {
		t.Fatalf("unexpected OIDC issuer URL: %q", cfg.Auth.OIDCIssuerURL)
	}
	if cfg.Auth.OIDCAudience != "mcp-client" {
		t.Fatalf("unexpected OIDC audience: %q", cfg.Auth.OIDCAudience)
	}
	if cfg.Auth.OIDCClientID != "mcp-client-id" {
		t.Fatalf("unexpected OIDC client ID: %q", cfg.Auth.OIDCClientID)
	}
	if cfg.Auth.OIDCHTTPTimeoutSec != 7 {
		t.Fatalf("unexpected OIDC HTTP timeout: %d", cfg.Auth.OIDCHTTPTimeoutSec)
	}
	if cfg.Auth.OIDCJWKSCacheTTLSec != 120 {
		t.Fatalf("unexpected OIDC cache TTL: %d", cfg.Auth.OIDCJWKSCacheTTLSec)
	}
}

func TestAuthBearerRequiresTokenOrOIDC(t *testing.T) {
	t.Setenv("MCP_AUTH_ENABLED", "true")
	t.Setenv("MCP_AUTH_MODE", "bearer")
	t.Setenv("MCP_AUTH_BEARER_TOKEN", "")
	t.Setenv("MCP_AUTH_OIDC_ISSUER_URL", "")
	t.Setenv("MCP_AUTH_OIDC_DISCOVERY_URL", "")

	_, err := Load("")
	if err == nil {
		t.Fatal("expected error when bearer mode has neither static token nor OIDC discovery config")
	}
	if !strings.Contains(err.Error(), "bearer token is required for bearer auth mode when OIDC discovery is not configured") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAuthOIDCInvalidIssuerURL(t *testing.T) {
	t.Setenv("MCP_AUTH_ENABLED", "true")
	t.Setenv("MCP_AUTH_MODE", "bearer")
	t.Setenv("MCP_AUTH_BEARER_TOKEN", "")
	t.Setenv("MCP_AUTH_OIDC_ISSUER_URL", "://invalid-url")

	_, err := Load("")
	if err == nil {
		t.Fatal("expected validation error for invalid OIDC issuer URL")
	}
	if !strings.Contains(err.Error(), "invalid auth OIDC issuer URL") {
		t.Fatalf("unexpected error: %v", err)
	}
}

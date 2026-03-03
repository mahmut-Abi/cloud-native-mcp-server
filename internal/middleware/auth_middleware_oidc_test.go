package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestAuthMiddlewareBearerOIDC_Success(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	var discoveryRequests atomic.Int32
	var jwksRequests atomic.Int32

	oidcServer, issuerURL := startOIDCTestServer(t, &privateKey.PublicKey, &discoveryRequests, &jwksRequests)
	defer oidcServer.Close()

	handler := AuthMiddleware(AuthConfig{
		Enabled:       true,
		Mode:          "bearer",
		OIDCIssuerURL: issuerURL,
		OIDCAudience:  "mcp-client",
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	token := signOIDCTestToken(t, privateKey, issuerURL, "mcp-client")
	req := httptest.NewRequest(http.MethodGet, "/api/aggregate/sse", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%q", rec.Code, rec.Body.String())
	}
}

func TestAuthMiddlewareBearerOIDC_AudienceMismatch(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	var discoveryRequests atomic.Int32
	var jwksRequests atomic.Int32

	oidcServer, issuerURL := startOIDCTestServer(t, &privateKey.PublicKey, &discoveryRequests, &jwksRequests)
	defer oidcServer.Close()

	handler := AuthMiddleware(AuthConfig{
		Enabled:       true,
		Mode:          "bearer",
		OIDCIssuerURL: issuerURL,
		OIDCAudience:  "expected-audience",
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	token := signOIDCTestToken(t, privateKey, issuerURL, "different-audience")
	req := httptest.NewRequest(http.MethodGet, "/api/aggregate/sse", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d body=%q", rec.Code, rec.Body.String())
	}
}

func TestAuthMiddlewareBearerOIDC_UsesDiscoveryAndJWKSCache(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	var discoveryRequests atomic.Int32
	var jwksRequests atomic.Int32

	oidcServer, issuerURL := startOIDCTestServer(t, &privateKey.PublicKey, &discoveryRequests, &jwksRequests)
	defer oidcServer.Close()

	handler := AuthMiddleware(AuthConfig{
		Enabled:             true,
		Mode:                "bearer",
		OIDCIssuerURL:       issuerURL,
		OIDCAudience:        "mcp-client",
		OIDCJWKSCacheTTLSec: 300,
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	token := signOIDCTestToken(t, privateKey, issuerURL, "mcp-client")

	req1 := httptest.NewRequest(http.MethodGet, "/api/aggregate/sse", nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Fatalf("expected status 200 for first request, got %d", rec1.Code)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/api/aggregate/sse", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("expected status 200 for second request, got %d", rec2.Code)
	}

	if got := discoveryRequests.Load(); got != 1 {
		t.Fatalf("expected discovery endpoint called once, got %d", got)
	}
	if got := jwksRequests.Load(); got != 1 {
		t.Fatalf("expected jwks endpoint called once, got %d", got)
	}
}

func startOIDCTestServer(t *testing.T, publicKey *rsa.PublicKey, discoveryRequests, jwksRequests *atomic.Int32) (*httptest.Server, string) {
	t.Helper()

	const (
		issuerPath = "/issuer/test"
		keyID      = "test-key-id"
	)

	n := base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes())

	var testServer *httptest.Server
	testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case issuerPath + "/.well-known/openid-configuration":
			discoveryRequests.Add(1)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"issuer":   testServer.URL + issuerPath,
				"jwks_uri": testServer.URL + "/jwks",
			})
		case "/jwks":
			jwksRequests.Add(1)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"keys": []map[string]any{
					{
						"kty": "RSA",
						"kid": keyID,
						"use": "sig",
						"alg": "RS256",
						"n":   n,
						"e":   e,
					},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))

	return testServer, testServer.URL + issuerPath
}

func signOIDCTestToken(t *testing.T, privateKey *rsa.PrivateKey, issuer, audience string) string {
	t.Helper()

	claims := jwt.RegisteredClaims{
		Issuer:    issuer,
		Audience:  []string{audience},
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		Subject:   "test-user",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "test-key-id"

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	return signedToken
}

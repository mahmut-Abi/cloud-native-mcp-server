package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddlewareDisabled(t *testing.T) {
	config := AuthConfig{Enabled: false}

	handler := AuthMiddleware(config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthMiddlewareAPIKey(t *testing.T) {
	config := AuthConfig{
		Enabled: true,
		Mode:    "apikey",
		APIKey:  "test-api-key-123",
	}

	handler := AuthMiddleware(config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Valid API key
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "test-api-key-123")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for valid key, got %d", http.StatusOK, w.Code)
	}

	// Invalid API key
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.Header.Set("X-API-Key", "wrong-key")
	w2 := httptest.NewRecorder()

	handler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d for invalid key, got %d", http.StatusUnauthorized, w2.Code)
	}
}

func TestAuthMiddlewareBearer(t *testing.T) {
	config := AuthConfig{
		Enabled:     true,
		Mode:        "bearer",
		BearerToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
	}

	handler := AuthMiddleware(config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Valid bearer token
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for valid token, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthMiddlewareBasic(t *testing.T) {
	config := AuthConfig{
		Enabled:  true,
		Mode:     "basic",
		Username: "admin",
		Password: "password123",
	}

	handler := AuthMiddleware(config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Valid basic auth
	req := httptest.NewRequest("GET", "/test", nil)
	req.SetBasicAuth("admin", "password123")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for valid auth, got %d", http.StatusOK, w.Code)
	}
}

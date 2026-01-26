package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiterAllow(t *testing.T) {
	limiter := NewRateLimiter(10.0, 5)

	clientID := "client-1"

	// First request should be allowed
	if !limiter.Allow(clientID) {
		t.Error("First request should be allowed")
	}

	// Additional requests within burst
	for i := 1; i < 5; i++ {
		if !limiter.Allow(clientID) {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}
}

func TestRateLimiterTokenRefill(t *testing.T) {
	limiter := NewRateLimiter(1.0, 1) // 1 request per second, burst 1

	clientID := "client-2"

	// Use initial token
	if !limiter.Allow(clientID) {
		t.Error("First request should be allowed")
	}

	// Wait for token to refill
	time.Sleep(1100 * time.Millisecond)

	// Should have at least one token now
	if !limiter.Allow(clientID) {
		t.Error("Request after refill should be allowed")
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	wrapped := RateLimitMiddleware(10.0, 2)(handler)

	// Make requests within burst
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "192.0.2.1:12345"
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Request %d should succeed, got %d", i+1, rec.Code)
		}
	}
}

func TestRateLimiterMultipleClients(t *testing.T) {
	limiter := NewRateLimiter(10.0, 2)

	// Different clients should have independent limits
	if !limiter.Allow("client-1") {
		t.Error("client-1 first request should be allowed")
	}

	if !limiter.Allow("client-2") {
		t.Error("client-2 first request should be allowed")
	}

	// Both can make second request due to burst
	if !limiter.Allow("client-1") {
		t.Error("client-1 second request should be allowed (within burst)")
	}

	if !limiter.Allow("client-2") {
		t.Error("client-2 second request should be allowed (within burst)")
	}
}

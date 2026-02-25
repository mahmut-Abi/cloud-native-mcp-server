package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiterAllow(t *testing.T) {
	limiter := NewRateLimiter(10.0, 5)
	defer limiter.Close()

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
	defer limiter.Close()

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

	limiter := NewRateLimiter(10.0, 2)
	defer limiter.Close()
	wrapped := RateLimitMiddlewareWithLimiter(limiter)(handler)

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
	defer limiter.Close()

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

func TestRateLimitMiddleware_SameIPDifferentPorts(t *testing.T) {
	limiter := NewRateLimiter(0.0001, 1)
	defer limiter.Close()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := RateLimitMiddlewareWithLimiter(limiter)(handler)

	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req1.RemoteAddr = "192.0.2.10:11111"
	rec1 := httptest.NewRecorder()
	wrapped.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Fatalf("first request should succeed, got %d", rec1.Code)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.RemoteAddr = "192.0.2.10:22222"
	rec2 := httptest.NewRecorder()
	wrapped.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusTooManyRequests {
		t.Fatalf("second request from same IP with different port should be rate limited, got %d", rec2.Code)
	}
}

func TestRateLimitMiddleware_UsesFirstXForwardedForHop(t *testing.T) {
	limiter := NewRateLimiter(0.0001, 1)
	defer limiter.Close()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := RateLimitMiddlewareWithLimiter(limiter)(handler)

	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req1.Header.Set("X-Forwarded-For", "203.0.113.5, 10.0.0.1")
	req1.RemoteAddr = "10.0.0.2:12345"
	rec1 := httptest.NewRecorder()
	wrapped.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Fatalf("first request should succeed, got %d", rec1.Code)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.Header.Set("X-Forwarded-For", "203.0.113.5, 10.0.0.99")
	req2.RemoteAddr = "10.0.0.3:22222"
	rec2 := httptest.NewRecorder()
	wrapped.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusTooManyRequests {
		t.Fatalf("second request with same first XFF hop should be rate limited, got %d", rec2.Code)
	}
}

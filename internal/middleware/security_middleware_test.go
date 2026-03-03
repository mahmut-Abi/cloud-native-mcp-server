package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type flushCountingRecorder struct {
	*httptest.ResponseRecorder
	flushCount int
}

func (f *flushCountingRecorder) Flush() {
	f.flushCount++
	f.ResponseRecorder.Flush()
}

func TestSecurityMiddleware_PreservesFlusherForStreaming(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("expected writer to implement http.Flusher")
		}
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: ping\n\n"))
		flusher.Flush()
	})

	wrapped := SecurityMiddleware(DefaultSecurityConfig())(handler)
	req := httptest.NewRequest(http.MethodGet, "/api/aggregate/sse", nil)
	rec := &flushCountingRecorder{ResponseRecorder: httptest.NewRecorder()}

	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if rec.flushCount == 0 {
		t.Fatal("expected at least one flush call to pass through security middleware")
	}
	if got := rec.Header().Get("Content-Security-Policy"); got == "" {
		t.Fatal("expected security headers to be set")
	}
}

func TestSecurityMiddleware_PreservesExistingCacheControl(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-transform")
		w.WriteHeader(http.StatusOK)
	})

	wrapped := SecurityMiddleware(DefaultSecurityConfig())(handler)
	req := httptest.NewRequest(http.MethodGet, "/api/aggregate/streamable-http", nil)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if got := rec.Header().Get("Cache-Control"); got != "no-cache, no-transform" {
		t.Fatalf("expected Cache-Control to be preserved, got %q", got)
	}
}

package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testContextKey string

func TestRegisterBackendAuthHandler(t *testing.T) {
	key := testContextKey("test-key")
	RegisterBackendAuthHandler("test-service", func(r *http.Request) (*http.Request, error) {
		return r.WithContext(context.WithValue(r.Context(), key, "test-value")), nil
	})
	RegisterBackendAuthHandler("test-service", func(r *http.Request) (*http.Request, error) {
		return r.WithContext(context.WithValue(r.Context(), key, "overwritten")), nil
	})

	backendHandlersMu.RLock()
	h, ok := backendHandlers["test-service"]
	backendHandlersMu.RUnlock()

	if !ok {
		t.Fatal("expected handler to be registered")
	}

	r, err := h(httptest.NewRequest("GET", "/", nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := r.Context().Value(key); v != "overwritten" {
		t.Errorf("expected context value 'overwritten', got %v", v)
	}
}

func TestBackendAuthMiddlewareNoHandler(t *testing.T) {
	// Clear registered handlers
	backendHandlersMu.Lock()
	backendHandlers = make(map[string]BackendAuthHandler)
	backendHandlersMu.Unlock()

	mw := BackendAuthMiddleware("nonexistent")
	called := false
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/test/sse/message", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !called {
		t.Error("expected handler to be called even without registered backend handler")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestBackendAuthMiddlewareWithHandler(t *testing.T) {
	key := testContextKey("prom-client")
	RegisterBackendAuthHandler("prometheus", func(r *http.Request) (*http.Request, error) {
		ctx := context.WithValue(r.Context(), key, "mock-client")
		return r.WithContext(ctx), nil
	})

	mw := BackendAuthMiddleware("prometheus")
	var capturedCtx context.Context
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedCtx = r.Context()
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/prometheus/sse/message", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if v := capturedCtx.Value(key); v != "mock-client" {
		t.Errorf("expected context value 'mock-client', got %v", v)
	}
}

func TestBackendAuthMiddlewareHandlerErrorPassthrough(t *testing.T) {
	RegisterBackendAuthHandler("failing-service", func(r *http.Request) (*http.Request, error) {
		return r, http.ErrServerClosed
	})

	mw := BackendAuthMiddleware("failing-service")
	called := false
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/test/sse/message", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !called {
		t.Error("expected handler to be called even when backend handler errors")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestBackendAuthMiddlewareConcurrent(t *testing.T) {
	RegisterBackendAuthHandler("concurrent-service", func(r *http.Request) (*http.Request, error) {
		return r.WithContext(context.WithValue(r.Context(), "concurrent", true)), nil
	})

	mw := BackendAuthMiddleware("concurrent-service")

	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Context().Value("concurrent") != true {
					t.Error("expected concurrent context value")
				}
				w.WriteHeader(http.StatusOK)
			}))
			req := httptest.NewRequest("GET", "/api/test/sse/message", nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				t.Errorf("expected 200, got %d", rec.Code)
			}
			done <- true
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}

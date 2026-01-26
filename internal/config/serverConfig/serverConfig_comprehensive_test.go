package serverConfig

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestInitMCPServer tests MCP server initialization
func TestInitMCPServer(t *testing.T) {
	config := &ServerConfig{}
	hooks := config.InitHooks()
	if hooks == nil {
		t.Error("InitHooks should return non-nil hooks")
	}

	mcpServer := config.InitMCPServer(hooks)
	if mcpServer == nil {
		t.Error("InitMCPServer should return non-nil MCP server")
	}
}

// TestResponseWriter tests custom response writer
func TestResponseWriter_WriteHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: rec}

	rw.WriteHeader(http.StatusOK)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

// TestLoggingMiddlewareWithError tests logging middleware with error status
func TestLoggingMiddlewareWithError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error"))
	})

	wrapped := loggingMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", rec.Code)
	}
}

// TestCORSMiddlewareAllowsRequests tests CORS headers are set
func TestCORSMiddlewareAllowsRequests(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	s := &ServerConfig{}
	wrapped := s.corsMiddleware(handler)

	req := httptest.NewRequest(http.MethodPost, "/api/test", nil)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("CORS Allow-Origin header not set")
	}
	if rec.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("CORS Allow-Methods header not set")
	}
}

// TestContextManagement tests context operations
func TestContextManagement(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if ctx.Err() != nil {
		t.Error("Context should not have error after creation")
	}

	cancel()

	if ctx.Err() == nil {
		t.Error("Context should have error after cancel")
	}
}

// TestHealthCheckResponseFormat tests health check response format
func TestHealthCheckResponseFormat(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	healthCheckHandler(rec, req)

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", ct)
	}

	body := rec.Body.String()
	if body != "{\"status\":\"healthy\"}" {
		t.Errorf("Unexpected response body: %s", body)
	}
}

// TestCORSOptionsRequest tests OPTIONS request handling
func TestCORSOptionsRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called for OPTIONS request")
	})

	s := &ServerConfig{}
	wrapped := s.corsMiddleware(handler)

	req := httptest.NewRequest(http.MethodOptions, "/api/test", nil)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 for OPTIONS, got %d", rec.Code)
	}
}

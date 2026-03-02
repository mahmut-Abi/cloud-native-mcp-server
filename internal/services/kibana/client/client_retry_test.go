package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRetryOnTransientStatusThenSuccess(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts <= 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"error":"temporary unavailable"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		URL:            server.URL,
		Timeout:        2 * time.Second,
		MaxRetries:     2,
		RetryBaseDelay: 1 * time.Millisecond,
		RetryMaxDelay:  5 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	spaces, err := client.GetSpaces(context.Background())
	if err != nil {
		t.Fatalf("GetSpaces() unexpected error = %v", err)
	}
	if len(spaces) != 0 {
		t.Fatalf("expected empty spaces, got %d", len(spaces))
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestNoRetryForNonIdempotentMethod(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"error":"temporary unavailable"}`))
	}))
	defer server.Close()

	client, err := NewClient(&ClientOptions{
		URL:            server.URL,
		Timeout:        2 * time.Second,
		MaxRetries:     3,
		RetryBaseDelay: 1 * time.Millisecond,
		RetryMaxDelay:  5 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.CreateSpace(context.Background(), Space{
		ID:   "ops",
		Name: "Operations",
	})
	if err == nil {
		t.Fatalf("expected error for non-idempotent request")
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt for non-idempotent request, got %d", attempts)
	}
}

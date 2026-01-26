package optimize

import (
	"net/http"
	"testing"
	"time"
)

func TestNewOptimizedHTTPClient(t *testing.T) {
	client := NewOptimizedHTTPClient()
	if client == nil {
		t.Error("Expected client, got nil")
		return
	}
	if _, ok := client.Transport.(*http.Transport); !ok {
		t.Error("Expected http.Transport")
	}
}

func TestNewOptimizedHTTPClientWithTimeout(t *testing.T) {
	timeout := 5 * time.Second
	client := NewOptimizedHTTPClientWithTimeout(timeout)
	if client == nil {
		t.Error("Expected client, got nil")
		return
	}
	if client.Timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, client.Timeout)
	}
}

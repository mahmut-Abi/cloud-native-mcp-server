package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestInMemoryAuditStorage tests the in-memory audit storage
func TestInMemoryAuditStorage(t *testing.T) {
	storage := NewInMemoryAuditStorage(100)

	// Test logging entries
	entry1 := &AuditLogEntry{
		Timestamp:   time.Now(),
		CallerIP:    "192.168.1.1",
		UserID:      "user1",
		ToolName:    "k8s_get_pods",
		ServiceName: "kubernetes",
		Action:      "GET /api/v1/pods",
		InputParams: map[string]interface{}{"namespace": "default"},
		Output:      "pod list",
		Duration:    150,
		Status:      "success",
	}

	entry2 := &AuditLogEntry{
		Timestamp:   time.Now(),
		CallerIP:    "192.168.1.2",
		UserID:      "user2",
		ToolName:    "k8s_create_deployment",
		ServiceName: "kubernetes",
		Action:      "POST /api/v1/deployments",
		InputParams: map[string]interface{}{"name": "test-app"},
		Output:      "deployment created",
		Duration:    500,
		Status:      "failure",
		ErrorMsg:    "insufficient permissions",
	}

	// Log entries
	if err := storage.Log(entry1); err != nil {
		t.Errorf("Failed to log entry1: %v", err)
	}

	if err := storage.Log(entry2); err != nil {
		t.Errorf("Failed to log entry2: %v", err)
	}

	// Test querying all entries
	results, err := storage.Query(map[string]interface{}{})
	if err != nil {
		t.Errorf("Failed to query entries: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(results))
	}

	// Test querying by user ID
	results, err = storage.Query(map[string]interface{}{"user_id": "user1"})
	if err != nil {
		t.Errorf("Failed to query by user_id: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 entry for user1, got %d", len(results))
	}

	if results[0].UserID != "user1" {
		t.Errorf("Expected user1, got %s", results[0].UserID)
	}

	// Test querying by service name
	results, err = storage.Query(map[string]interface{}{"service_name": "kubernetes"})
	if err != nil {
		t.Errorf("Failed to query by service_name: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 entries for kubernetes service, got %d", len(results))
	}

	// Test querying by status
	results, err = storage.Query(map[string]interface{}{"status": "failure"})
	if err != nil {
		t.Errorf("Failed to query by status: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 failure entry, got %d", len(results))
	}

	// Test statistics
	startTime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now().Add(1 * time.Hour)
	stats, err := storage.GetStats(startTime, endTime)
	if err != nil {
		t.Errorf("Failed to get stats: %v", err)
	}

	if stats["total_logs"] != 2 {
		t.Errorf("Expected 2 total logs, got %v", stats["total_logs"])
	}

	if stats["success_count"] != 1 {
		t.Errorf("Expected 1 success, got %v", stats["success_count"])
	}

	if stats["failure_count"] != 1 {
		t.Errorf("Expected 1 failure, got %v", stats["failure_count"])
	}
}

// TestAuditMiddleware tests the audit middleware functionality
func TestAuditMiddleware(t *testing.T) {
	storage := NewInMemoryAuditStorage(100)

	config := AuditMiddlewareConfig{
		Enabled: true,
		Storage: storage,
	}

	// Create test handler that simulates a successful operation
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	})

	// Wrap with audit middleware
	auditedHandler := AuditMiddleware(config)(handler)

	// Create test request
	req := httptest.NewRequest("GET", "/api/kubernetes/pods?namespace=default", nil)
	req.Header.Set("X-User-ID", "test-user")
	req.Header.Set("X-Tool-Name", "k8s_get_pods")
	req.Header.Set("X-Service-Name", "kubernetes")
	req.Header.Set("X-Forwarded-For", "192.168.1.100")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute request
	auditedHandler.ServeHTTP(rr, req)

	// Verify response
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Verify audit log was created
	results, err := storage.Query(map[string]interface{}{})
	if err != nil {
		t.Errorf("Failed to query audit logs: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 audit log entry, got %d", len(results))
	}

	entry := results[0]
	if entry.UserID != "test-user" {
		t.Errorf("Expected user_id 'test-user', got '%s'", entry.UserID)
	}

	if entry.ToolName != "k8s_get_pods" {
		t.Errorf("Expected tool_name 'k8s_get_pods', got '%s'", entry.ToolName)
	}

	if entry.ServiceName != "kubernetes" {
		t.Errorf("Expected service_name 'kubernetes', got '%s'", entry.ServiceName)
	}

	if entry.CallerIP != "192.168.1.100" {
		t.Errorf("Expected caller_ip '192.168.1.100', got '%s'", entry.CallerIP)
	}

	if entry.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", entry.Status)
	}

	expectedAction := "GET /api/kubernetes/pods?namespace=default"
	if entry.Action != expectedAction {
		t.Errorf("Expected action '%s', got '%s'", expectedAction, entry.Action)
	}
}

// TestAuditMiddleware_Failure tests audit logging for failed requests
func TestAuditMiddleware_Failure(t *testing.T) {
	storage := NewInMemoryAuditStorage(100)

	config := AuditMiddlewareConfig{
		Enabled: true,
		Storage: storage,
	}

	// Create test handler that simulates a failed operation
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal error"))
	})

	// Wrap with audit middleware
	auditedHandler := AuditMiddleware(config)(handler)

	// Create test request
	req := httptest.NewRequest("POST", "/api/kubernetes/deployments", nil)
	req.Header.Set("X-User-ID", "test-user-2")
	req.Header.Set("X-Tool-Name", "k8s_create_deployment")
	req.Header.Set("X-Service-Name", "kubernetes")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute request
	auditedHandler.ServeHTTP(rr, req)

	// Verify response
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", rr.Code)
	}

	// Verify audit log was created with failure status
	results, err := storage.Query(map[string]interface{}{})
	if err != nil {
		t.Errorf("Failed to query audit logs: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 audit log entry, got %d", len(results))
	}

	entry := results[0]
	if entry.Status != "failure" {
		t.Errorf("Expected status 'failure', got '%s'", entry.Status)
	}

	if entry.UserID != "test-user-2" {
		t.Errorf("Expected user_id 'test-user-2', got '%s'", entry.UserID)
	}
}

// TestAuditMiddleware_Disabled tests that disabled audit middleware doesn't log
func TestAuditMiddleware_Disabled(t *testing.T) {
	storage := NewInMemoryAuditStorage(100)

	config := AuditMiddlewareConfig{
		Enabled: false, // Disabled
		Storage: storage,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	})

	auditedHandler := AuditMiddleware(config)(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	auditedHandler.ServeHTTP(rr, req)

	// Verify response is successful
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Verify no audit log was created
	results, err := storage.Query(map[string]interface{}{})
	if err != nil {
		t.Errorf("Failed to query audit logs: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 audit log entries, got %d", len(results))
	}
}

// TestAuditMiddleware_NilStorage tests behavior with nil storage
func TestAuditMiddleware_NilStorage(t *testing.T) {
	config := AuditMiddlewareConfig{
		Enabled: true,
		Storage: nil, // Nil storage
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	})

	auditedHandler := AuditMiddleware(config)(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Should not panic and should pass through the request
	auditedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

// TestAuditLogEntry_String tests the String method of AuditLogEntry
func TestAuditLogEntry_String(t *testing.T) {
	entry := &AuditLogEntry{
		Timestamp:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		CallerIP:    "192.168.1.1",
		UserID:      "test-user",
		ToolName:    "k8s_get_pods",
		ServiceName: "kubernetes",
		Action:      "GET /api/v1/pods",
		InputParams: map[string]interface{}{"namespace": "default"},
		Output:      "pod list",
		Duration:    150,
		Status:      "success",
	}

	jsonStr := entry.String()
	if jsonStr == "" {
		t.Error("Expected non-empty JSON string")
	}

	// Should be valid JSON
	if !isValidJSON(jsonStr) {
		t.Errorf("Expected valid JSON, got: %s", jsonStr)
	}
}

// TestInMemoryAuditStorage_MaxSize tests the max size limit
func TestInMemoryAuditStorage_MaxSize(t *testing.T) {
	maxSize := 3
	storage := NewInMemoryAuditStorage(maxSize)

	// Add more entries than the max size
	for i := 0; i < 5; i++ {
		entry := &AuditLogEntry{
			Timestamp:   time.Now(),
			CallerIP:    "192.168.1.1",
			UserID:      fmt.Sprintf("user%d", i),
			ToolName:    "test_tool",
			ServiceName: "test_service",
			Action:      "test action",
			Duration:    100,
			Status:      "success",
		}
		if err := storage.Log(entry); err != nil {
			t.Errorf("Failed to log entry %d: %v", i, err)
		}
	}

	// Should only keep the last maxSize entries
	results, err := storage.Query(map[string]interface{}{})
	if err != nil {
		t.Errorf("Failed to query entries: %v", err)
	}

	if len(results) != maxSize {
		t.Errorf("Expected %d entries, got %d", maxSize, len(results))
	}

	// Should have kept the last 3 entries (user2, user3, user4)
	expectedUsers := []string{"user2", "user3", "user4"}
	actualUsers := make([]string, len(results))
	for i, result := range results {
		actualUsers[i] = result.UserID
	}

	for i, expected := range expectedUsers {
		if actualUsers[i] != expected {
			t.Errorf("Expected user %s at position %d, got %s", expected, i, actualUsers[i])
		}
	}
}

// Helper function to check if a string is valid JSON
func isValidJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

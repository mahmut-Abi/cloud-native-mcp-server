package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const auditComponent = "audit"

// AuditStorage is an alias for AuditLogger
type AuditStorage = AuditLogger

// InMemoryAuditStorage stores audit logs in memory
type InMemoryAuditStorage struct {
	mu      sync.RWMutex
	logs    []*AuditLogEntry
	maxSize int
}

// Close is a no-op for memory storage
func (s *InMemoryAuditStorage) Close() error {
	return nil
}

// NewInMemoryAuditStorage creates a new in-memory audit storage
func NewInMemoryAuditStorage(maxSize int) *InMemoryAuditStorage {
	if maxSize == 0 {
		maxSize = 10000
	}
	return &InMemoryAuditStorage{
		logs:    make([]*AuditLogEntry, 0, maxSize),
		maxSize: maxSize,
	}
}

// Log stores an audit log entry
func (s *InMemoryAuditStorage) Log(entry *AuditLogEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.logs) >= s.maxSize {
		// Remove oldest entries if we exceed max size
		s.logs = s.logs[1:]
	}

	s.logs = append(s.logs, entry)
	return nil
}

// Query retrieves audit logs based on criteria with pagination support
// Supports pagination via criteria["page"] and criteria["pageSize"]
func (s *InMemoryAuditStorage) Query(criteria map[string]interface{}) ([]AuditLogEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Extract pagination parameters
	page := 1
	pageSize := 100 // default page size
	if p, ok := criteria["page"]; ok {
		if val, ok := p.(int); ok && val > 0 {
			page = val
		}
	}
	if ps, ok := criteria["pageSize"]; ok {
		if val, ok := ps.(int); ok && val > 0 && val <= 1000 {
			pageSize = val
		}
	}

	var results []AuditLogEntry
	for _, log := range s.logs {
		if matchesCriteria(log, criteria) {
			results = append(results, *log)
		}
	}

	// Apply pagination
	total := len(results)
	start := (page - 1) * pageSize
	if start >= total {
		return []AuditLogEntry{}, nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}

	return results[start:end], nil
}

// GetStats returns statistics about logged operations
func (s *InMemoryAuditStorage) GetStats(startTime, endTime time.Time) (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := map[string]interface{}{
		"total_logs":    len(s.logs),
		"success_count": 0,
		"failure_count": 0,
		"tools_used":    make(map[string]int),
		"services_used": make(map[string]int),
	}

	for _, log := range s.logs {
		if log.Timestamp.Before(startTime) || log.Timestamp.After(endTime) {
			continue
		}

		if log.Status == "success" {
			stats["success_count"] = stats["success_count"].(int) + 1
		} else {
			stats["failure_count"] = stats["failure_count"].(int) + 1
		}

		toolsMap := stats["tools_used"].(map[string]int)
		toolsMap[log.ToolName]++

		servicesMap := stats["services_used"].(map[string]int)
		servicesMap[log.ServiceName]++
	}

	return stats, nil
}

// matchesCriteria checks if a log entry matches the given criteria
func matchesCriteria(entry *AuditLogEntry, criteria map[string]interface{}) bool {
	for key, value := range criteria {
		switch key {
		case "user_id":
			if entry.UserID != value.(string) {
				return false
			}
		case "service_name":
			if entry.ServiceName != value.(string) {
				return false
			}
		case "tool_name":
			if entry.ToolName != value.(string) {
				return false
			}
		case "status":
			if entry.Status != value.(string) {
				return false
			}
		}
	}
	return true
}

// extractToolNameFromMCPRequest tries to extract the tool name from an MCP request body
func extractToolNameFromMCPRequest(requestBody interface{}) string {
	// Try to convert to map and look for MCP tool call structure
	if reqMap, ok := requestBody.(map[string]interface{}); ok {
		// Handle MCP request format: {"method": "tools/call", "params": {"name": "tool_name", ...}}
		if method, exists := reqMap["method"]; exists {
			if methodStr, ok := method.(string); ok && methodStr == "tools/call" {
				if params, exists := reqMap["params"]; exists {
					if paramsMap, ok := params.(map[string]interface{}); ok {
						if name, exists := paramsMap["name"]; exists {
							if nameStr, ok := name.(string); ok {
								return nameStr
							}
						}
					}
				}
			}
		}
	}
	return ""
}

// isSSERequest checks if the HTTP request is for an SSE connection.
func isSSERequest(r *http.Request) bool {
	// Check the Accept header first - most reliable indicator
	if r.Header.Get("Accept") == "text/event-stream" {
		return true
	}

	// Check URL path as backup method
	path := r.URL.Path
	return strings.HasSuffix(path, "/sse") || strings.HasSuffix(path, "/sse/message")
}

// AuditMiddlewareConfig is the configuration for audit middleware
type AuditMiddlewareConfig struct {
	Enabled           bool
	Storage           AuditStorage
	Masker            *SensitiveDataMasker
	EnableDataMasking bool
}

// maxBodySize is the maximum size of request body to capture in bytes (10KB)
const maxBodySize = 10 * 1024

// requestBodyReader is an io.ReadCloser that limits the body size and keeps a copy
type requestBodyReader struct {
	reader   io.Reader // Use io.Reader to wrap with LimitReader
	bodyCopy *bytes.Buffer
	maxSize  int
}

// Read implements io.ReadCloser
func (r *requestBodyReader) Read(p []byte) (n int, err error) {
	// Limit the read to maxSize to prevent excessive memory allocation
	remaining := r.maxSize - r.bodyCopy.Len()
	if remaining <= 0 {
		// Already captured maxSize bytes, just pass through without capturing
		return r.reader.Read(p)
	}

	// Create a limited reader for this read operation
	limitedReader := io.LimitReader(r.reader, int64(remaining))
	n, err = limitedReader.Read(p)

	if n > 0 {
		r.bodyCopy.Write(p[:n])
	}

	return n, err
}

// Close implements io.ReadCloser
func (r *requestBodyReader) Close() error {
	if closer, ok := r.reader.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// newRequestBodyReader creates a new request body reader
func newRequestBodyReader(original io.ReadCloser, maxSize int) *requestBodyReader {
	return &requestBodyReader{
		reader:   original,
		bodyCopy: &bytes.Buffer{},
		maxSize:  maxSize,
	}
}

// AuditMiddleware creates an HTTP middleware for audit logging
func AuditMiddleware(config AuditMiddlewareConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !config.Enabled || config.Storage == nil {
				next.ServeHTTP(w, r)
				return
			}

			startTime := time.Now()

			// Get client IP
			clientIP := r.RemoteAddr
			if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
				clientIP = xff
			}

			// Get user ID from header if available
			userID := r.Header.Get("X-User-ID")
			if userID == "" {
				userID = "anonymous"
			}

			// Capture request body if it's not a GET or HEAD request (without consuming it)
			var requestBody interface{}
			var bodyReader io.ReadCloser

			if r.Method != "GET" && r.Method != "HEAD" && r.Body != nil {
				// Use a body reader that captures a copy
				bodyReader = newRequestBodyReader(r.Body, maxBodySize)
				r.Body = bodyReader
			}

			// Wrap response writer to capture response status and possibly body
			writeWrapper := newResponseWriterWrapper(w)

			// Call next handler
			next.ServeHTTP(writeWrapper, r)

			// After request is processed, try to parse and include the body if available
			if bodyReader != nil {
				if br, ok := bodyReader.(*requestBodyReader); ok && br.bodyCopy.Len() > 0 {
					// Try to parse as JSON first
					var jsonData interface{}
					if err := json.Unmarshal(br.bodyCopy.Bytes(), &jsonData); err == nil {
						requestBody = jsonData
					} else {
						// If not valid JSON, use as string
						requestBody = br.bodyCopy.String()
					}
				}
			}

			// Create audit log entry
			duration := time.Since(startTime).Milliseconds()
			status := "success"
			// Special handling for SSE connections
			isSSE := isSSERequest(r)

			// Debug the status code capture
			logrus.WithFields(logrus.Fields{
				"status_code":   writeWrapper.statusCode,
				"path":          r.URL.Path,
				"is_sse":        isSSE,
				"accept_header": r.Header.Get("Accept"),
			}).Debug("Audit middleware status check")

			// For SSE connections, always mark as success regardless of status code
			// This is because SSE connections can sometimes return non-200 status codes
			// but still be successful from an application perspective
			if writeWrapper.statusCode >= 400 && !isSSE {
				logrus.WithFields(logrus.Fields{
					"status_code": writeWrapper.statusCode,
					"path":        r.URL.Path,
				}).Debug("Audit marking request as failure")
				status = "failure"
			}

			// Try to extract service name from path if not provided in header
			serviceName := r.Header.Get("X-Service-Name")
			if serviceName == "" {
				// Parse service name from URL path
				path := r.URL.Path
				// Look for paths like /api/{service}/...
				if strings.HasPrefix(path, "/api/") {
					parts := strings.Split(path, "/")
					if len(parts) >= 3 {
						// parts[0] is empty because path starts with "/", parts[1] is "api"
						candidateService := parts[2]
						if candidateService != "" && candidateService != "openapi.json" && candidateService != "docs" {
							// Remove any "sse" suffix for SSE endpoints
							serviceName = strings.TrimSuffix(candidateService, "/sse")
						}
					}
				}
			}

			// Try to extract tool name from path or query parameters if not provided in header
			toolName := r.Header.Get("X-Tool-Name")
			if toolName == "" && serviceName != "" {
				// Try to infer tool name from the path
				path := r.URL.Path
				// Skip "/api/{service}/" part
				prefix := fmt.Sprintf("/api/%s/", serviceName)
				if strings.HasPrefix(path, prefix) {
					pathRemainder := strings.TrimPrefix(path, prefix)
					parts := strings.Split(pathRemainder, "/")
					if len(parts) > 0 && parts[0] != "" && parts[0] != "sse" && parts[0] != "message" {
						// Use the first segment after service as potential tool category
						toolName = fmt.Sprintf("%s_%s", serviceName, parts[0])
					}
				}

				// For MCP requests, try to extract tool name from request body
				if toolName == "" && requestBody != nil {
					toolName = extractToolNameFromMCPRequest(requestBody)
				}
			}

			// Extract session ID if available from query parameter
			sessionID := r.URL.Query().Get("sessionId")

			// Prepare input parameters
			inputParams := map[string]interface{}{
				"method": r.Method,
				"path":   r.RequestURI,
			}

			// Add session ID to input parameters if available
			if sessionID != "" {
				inputParams["session_id"] = sessionID
			}

			// Add body to input parameters if available
			if requestBody != nil {
				inputParams["body"] = requestBody
			}

			// Prepare output if any was captured
			var output interface{}
			if writeWrapper.bodyCapture && writeWrapper.buffer.Len() > 0 {
				// Try to parse as JSON first
				var jsonOutput interface{}
				if err := json.Unmarshal(writeWrapper.buffer.Bytes(), &jsonOutput); err == nil {
					output = jsonOutput
				} else {
					// Use as string if not valid JSON
					output = writeWrapper.buffer.String()
				}
			}

			entry := &AuditLogEntry{
				Timestamp:   time.Now(),
				CallerIP:    clientIP,
				UserID:      userID,
				ToolName:    toolName,
				ServiceName: serviceName,
				Action:      fmt.Sprintf("%s %s", r.Method, r.RequestURI),
				InputParams: inputParams,
				Output:      output,
				Duration:    duration,
				Status:      status,
			}

			// Apply sensitive data masking if enabled
			if config.EnableDataMasking && config.Masker != nil {
				config.Masker.MaskAuditEntry(entry)
			}

			// Log the entry
			if config.Storage != nil {
				if err := config.Storage.Log(entry); err != nil {
					logrus.WithFields(logrus.Fields{
						"component": auditComponent,
						"operation": "log_entry",
						"error":     err,
					}).Error("Failed to log audit entry")
				} else {
					logrus.WithFields(logrus.Fields{
						"component":    auditComponent,
						"operation":    "log_entry",
						"user_id":      userID,
						"service_name": entry.ServiceName,
						"tool_name":    entry.ToolName,
						"status":       status,
						"duration_ms":  duration,
					}).Debug("Audit entry logged")
				}
			}
		})
	}
}

// responseWriterWrapper wraps http.ResponseWriter to capture status code and response body
type responseWriterWrapper struct {
	http.ResponseWriter
	headerWritten bool
	statusCode    int
	buffer        bytes.Buffer
	bodyCapture   bool // Only capture body for specific content types
	size          int  // Response size in bytes
}

// WriteHeader captures the status code
func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	if !w.headerWritten {
		w.statusCode = statusCode
		w.headerWritten = true

		// Check content type to determine if we should capture body
		contentType := w.ResponseWriter.Header().Get("Content-Type")
		w.bodyCapture = strings.Contains(contentType, "json") ||
			strings.Contains(contentType, "text") ||
			strings.Contains(contentType, "xml")
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write ensures status code is captured even if WriteHeader isn't called explicitly
func (w *responseWriterWrapper) Write(data []byte) (int, error) {
	if !w.headerWritten {
		// If WriteHeader hasn't been called, assume 200 OK
		w.WriteHeader(http.StatusOK)
	}

	// If content type is one we want to capture, and buffer size is reasonable
	if w.bodyCapture && w.buffer.Len() < maxBodySize {
		remaining := maxBodySize - w.buffer.Len()
		if len(data) <= remaining {
			w.buffer.Write(data)
		} else {
			w.buffer.Write(data[:remaining])
		}
	}

	n, err := w.ResponseWriter.Write(data)
	w.size += n
	return n, err
}

// Flush implements the http.Flusher interface to ensure SSE streaming works correctly
func (w *responseWriterWrapper) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// newResponseWriterWrapper creates a new response writer wrapper
func newResponseWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

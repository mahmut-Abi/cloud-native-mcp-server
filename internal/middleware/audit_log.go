package middleware

import (
	"encoding/json"
	"fmt"
	"time"
)

// AuditLogEntry represents a single audit log entry
type AuditLogEntry struct {
	Timestamp   time.Time              `json:"timestamp"   db:"timestamp"`      // Time of operation
	CallerIP    string                 `json:"caller_ip"   db:"caller_ip"`      // IP address of caller
	UserID      string                 `json:"user_id"     db:"user_id"`        // User identifier
	ToolName    string                 `json:"tool_name"   db:"tool_name"`      // Tool/command name
	ServiceName string                 `json:"service_name" db:"service_name"`  // Service name
	Action      string                 `json:"action"      db:"action"`         // Action performed
	InputParams map[string]interface{} `json:"input_params"  db:"input_params"` // Input parameters
	Output      interface{}            `json:"output"      db:"output"`         // Output/result
	Duration    int64                  `json:"duration_ms"  db:"duration_ms"`   // Execution duration in ms
	Status      string                 `json:"status"      db:"status"`         // success/failure
	ErrorMsg    string                 `json:"error_msg"    db:"error_msg"`     // Error message if failed
}

// AuditLogger interface for different audit logging backends
type AuditLogger interface {
	// Log records an audit log entry
	Log(entry *AuditLogEntry) error

	// Query retrieves audit logs based on criteria
	Query(criteria map[string]interface{}) ([]AuditLogEntry, error)

	// GetStats returns statistics about logged operations
	GetStats(startTime, endTime time.Time) (map[string]interface{}, error)

	// Close closes the logger and releases resources
	Close() error
}

// String returns JSON representation of the audit log entry
func (entry *AuditLogEntry) String() string {
	byte, err := json.Marshal(entry)
	if err != nil {
		return fmt.Sprintf("{error: %v}", err)
	}
	return string(byte)
}

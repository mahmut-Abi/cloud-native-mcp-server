package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"           // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/sirupsen/logrus"
)

// DatabaseAuditStorage implements audit storage using database
type DatabaseAuditStorage struct {
	mu              sync.RWMutex
	db              *sql.DB
	tableName       string
	maxRecords      int
	cleanupInterval time.Duration
	stopCleanup     chan struct{}
}

// NewDatabaseAuditStorage creates a new database-based audit storage
func NewDatabaseAuditStorage(dbType, connectionString, sqlitePath, tableName string, maxRecords, cleanupIntervalHours int) (*DatabaseAuditStorage, error) {
	var db *sql.DB
	var err error

	if tableName == "" {
		tableName = "audit_logs"
	}

	switch dbType {
	case "sqlite":
		if sqlitePath != "" {
			// Ensure directory exists
			dir := filepath.Dir(sqlitePath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create SQLite directory: %w", err)
			}
			connectionString = sqlitePath
		}
		db, err = sql.Open("sqlite3", connectionString)
	case "postgresql", "postgres":
		db, err = sql.Open("postgres", connectionString)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	storage := &DatabaseAuditStorage{
		db:              db,
		tableName:       tableName,
		maxRecords:      maxRecords,
		cleanupInterval: time.Duration(cleanupIntervalHours) * time.Hour,
		stopCleanup:     make(chan struct{}),
	}

	// Create table if not exists
	if err := storage.createTable(); err != nil {
		return nil, fmt.Errorf("failed to create audit table: %w", err)
	}

	// Start cleanup routine if configured
	if cleanupIntervalHours > 0 {
		go storage.cleanupRoutine()
	}

	logrus.WithFields(logrus.Fields{
		"component":  "audit",
		"db_type":    dbType,
		"table_name": tableName,
	}).Info("Database audit storage initialized")

	return storage, nil
}

// isValidTableName validates that the table name contains only safe characters
// to prevent SQL injection attacks
func isValidTableName(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '_' {
			return false
		}
	}
	return true
}

// createTable creates the audit logs table
func (s *DatabaseAuditStorage) createTable() error {
	// Validate table name to prevent SQL injection
	if !isValidTableName(s.tableName) {
		return fmt.Errorf("invalid table name: %s (only alphanumeric characters and underscore allowed)", s.tableName)
	}

	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME NOT NULL,
			caller_ip VARCHAR(45) NOT NULL,
			user_id VARCHAR(255) NOT NULL,
			tool_name VARCHAR(255),
			service_name VARCHAR(255),
			action TEXT NOT NULL,
			input_params TEXT,
			output TEXT,
			duration_ms INTEGER,
			status VARCHAR(50) NOT NULL,
			error_msg TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`, s.tableName)

	_, err := s.db.Exec(query)
	return err
}

// Log stores an audit log entry in database
func (s *DatabaseAuditStorage) Log(entry *AuditLogEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Serialize input_params and output
	inputParamsJSON, _ := json.Marshal(entry.InputParams)
	outputJSON, _ := json.Marshal(entry.Output)

	query := fmt.Sprintf(`
		INSERT INTO %s (
			timestamp, caller_ip, user_id, tool_name, service_name, 
			action, input_params, output, duration_ms, status, error_msg
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, s.tableName)

	_, err := s.db.Exec(query,
		entry.Timestamp,
		entry.CallerIP,
		entry.UserID,
		entry.ToolName,
		entry.ServiceName,
		entry.Action,
		string(inputParamsJSON),
		string(outputJSON),
		entry.Duration,
		entry.Status,
		entry.ErrorMsg,
	)

	if err != nil {
		return fmt.Errorf("failed to insert audit log: %w", err)
	}

	return nil
}

// Query retrieves audit logs from database
func (s *DatabaseAuditStorage) Query(criteria map[string]interface{}) ([]AuditLogEntry, error) {
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

	// Backward compatibility: support max_results
	if max, ok := criteria["max_results"]; ok {
		if maxInt, ok := max.(int); ok && maxInt > 0 && maxInt < pageSize {
			pageSize = maxInt
		}
	}

	// Calculate offset for pagination
	offset := (page - 1) * pageSize

	// Build query
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(fmt.Sprintf("SELECT timestamp, caller_ip, user_id, tool_name, service_name, action, input_params, output, duration_ms, status, error_msg FROM %s WHERE 1=1", s.tableName))

	var args []interface{}
	argIndex := 1

	// Add criteria
	if userID, ok := criteria["user_id"]; ok {
		queryBuilder.WriteString(fmt.Sprintf(" AND user_id = $%d", argIndex))
		args = append(args, userID)
		argIndex++
	}

	if serviceName, ok := criteria["service_name"]; ok {
		queryBuilder.WriteString(fmt.Sprintf(" AND service_name = $%d", argIndex))
		args = append(args, serviceName)
		argIndex++
	}

	if toolName, ok := criteria["tool_name"]; ok {
		queryBuilder.WriteString(fmt.Sprintf(" AND tool_name = $%d", argIndex))
		args = append(args, toolName)
		argIndex++
	}

	if status, ok := criteria["status"]; ok {
		queryBuilder.WriteString(fmt.Sprintf(" AND status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}

	if startTime, ok := criteria["start_time"]; ok {
		queryBuilder.WriteString(fmt.Sprintf(" AND timestamp >= $%d", argIndex))
		args = append(args, startTime)
		argIndex++
	}

	if endTime, ok := criteria["end_time"]; ok {
		queryBuilder.WriteString(fmt.Sprintf(" AND timestamp <= $%d", argIndex))
		args = append(args, endTime)
	}

	queryBuilder.WriteString(" ORDER BY timestamp DESC")
	queryBuilder.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset))

	rows, err := s.db.Query(queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var results []AuditLogEntry
	for rows.Next() {
		var entry AuditLogEntry
		var inputParamsJSON, outputJSON string

		err := rows.Scan(
			&entry.Timestamp,
			&entry.CallerIP,
			&entry.UserID,
			&entry.ToolName,
			&entry.ServiceName,
			&entry.Action,
			&inputParamsJSON,
			&outputJSON,
			&entry.Duration,
			&entry.Status,
			&entry.ErrorMsg,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log row: %w", err)
		}

		// Deserialize JSON fields
		if inputParamsJSON != "" {
			_ = json.Unmarshal([]byte(inputParamsJSON), &entry.InputParams)
		}
		if outputJSON != "" {
			_ = json.Unmarshal([]byte(outputJSON), &entry.Output)
		}

		results = append(results, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating audit log rows: %w", err)
	}

	return results, nil
}

// GetStats returns statistics from database
func (s *DatabaseAuditStorage) GetStats(startTime, endTime time.Time) (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := map[string]interface{}{
		"storage_type": "database",
		"table_name":   s.tableName,
		"start_time":   startTime,
		"end_time":     endTime,
	}

	// Get total count
	var totalCount int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE timestamp BETWEEN ? AND ?", s.tableName)
	err := s.db.QueryRow(query, startTime, endTime).Scan(&totalCount)
	if err != nil {
		logrus.WithError(err).Error("Failed to get total audit log count")
	} else {
		stats["total_logs"] = totalCount
	}

	// Get success/failure counts
	var successCount, failureCount int
	successQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE timestamp BETWEEN ? AND ? AND status = 'success'", s.tableName)
	failureQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE timestamp BETWEEN ? AND ? AND status = 'failure'", s.tableName)

	_ = s.db.QueryRow(successQuery, startTime, endTime).Scan(&successCount)
	_ = s.db.QueryRow(failureQuery, startTime, endTime).Scan(&failureCount)

	stats["success_count"] = successCount
	stats["failure_count"] = failureCount

	return stats, nil
}

// cleanupRoutine periodically cleans up old records
func (s *DatabaseAuditStorage) cleanupRoutine() {
	ticker := time.NewTicker(s.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanup()
		case <-s.stopCleanup:
			return
		}
	}
}

// cleanup removes old records if maxRecords is exceeded
func (s *DatabaseAuditStorage) cleanup() {
	if s.maxRecords <= 0 {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Count current records
	var count int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", s.tableName)
	err := s.db.QueryRow(countQuery).Scan(&count)
	if err != nil {
		logrus.WithError(err).Error("Failed to count audit records for cleanup")
		return
	}

	if count <= s.maxRecords {
		return // No cleanup needed
	}

	// Delete oldest records
	toDelete := count - s.maxRecords
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id IN (SELECT id FROM %s ORDER BY timestamp ASC LIMIT %d)", s.tableName, s.tableName, toDelete)

	result, err := s.db.Exec(deleteQuery)
	if err != nil {
		logrus.WithError(err).Error("Failed to cleanup old audit records")
		return
	}

	deleted, _ := result.RowsAffected()
	logrus.WithFields(logrus.Fields{
		"component": "audit",
		"operation": "cleanup",
		"deleted":   deleted,
	}).Info("Cleaned up old audit records")
}

// Close closes the database connection
func (s *DatabaseAuditStorage) Close() error {
	close(s.stopCleanup)
	return s.db.Close()
}

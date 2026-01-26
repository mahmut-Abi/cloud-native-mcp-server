package middleware

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// FileAuditStorage implements audit storage to files with rotation
type FileAuditStorage struct {
	mu          sync.RWMutex
	filePath    string
	maxSizeMB   int
	maxBackups  int
	maxAgeDays  int
	compress    bool
	format      string // "json" or "text"
	currentFile *os.File
}

// NewFileAuditStorage creates a new file-based audit storage
func NewFileAuditStorage(filePath string, maxSizeMB, maxBackups, maxAgeDays int, compress bool, format string) (*FileAuditStorage, error) {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create audit log directory: %w", err)
	}

	storage := &FileAuditStorage{
		filePath:   filePath,
		maxSizeMB:  maxSizeMB,
		maxBackups: maxBackups,
		maxAgeDays: maxAgeDays,
		compress:   compress,
		format:     format,
	}

	// Open the current log file
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open audit log file: %w", err)
	}
	storage.currentFile = file

	logrus.WithFields(logrus.Fields{
		"component": "audit",
		"path":      filePath,
		"format":    format,
	}).Info("File audit storage initialized")

	return storage, nil
}

// Log writes an audit log entry to file
func (s *FileAuditStorage) Log(entry *AuditLogEntry) error {
	if s == nil {
		return fmt.Errorf("file audit storage is nil")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	var logLine string

	if s.format == "json" {
		bytes, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("failed to marshal audit entry: %w", err)
		}
		logLine = string(bytes) + "\n"
	} else {
		// Text format
		logLine = fmt.Sprintf("%s [%s] %s@%s %s:%s %s (duration: %dms, status: %s)\n",
			entry.Timestamp.Format(time.RFC3339),
			entry.ServiceName,
			entry.UserID,
			entry.CallerIP,
			entry.ToolName,
			entry.Action,
			entry.ErrorMsg,
			entry.Duration,
			entry.Status,
		)
	}

	// Check if rotation is needed
	if err := s.rotateIfNeeded(); err != nil {
		logrus.WithError(err).Error("Failed to rotate audit log file")
	}

	// Write to current file
	if _, err := s.currentFile.WriteString(logLine); err != nil {
		return fmt.Errorf("failed to write audit log: %w", err)
	}

	// Sync to disk
	if err := s.currentFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync audit log: %w", err)
	}

	return nil
}

// Query retrieves audit logs from file (basic implementation)
func (s *FileAuditStorage) Query(criteria map[string]interface{}) ([]AuditLogEntry, error) {
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
	maxResults := 1000
	if max, ok := criteria["max_results"]; ok {
		if maxInt, ok := max.(int); ok {
			maxResults = maxInt
		}
	}

	// Open file for reading
	file, err := os.Open(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []AuditLogEntry{}, nil // No logs yet
		}
		return nil, fmt.Errorf("failed to open audit log file for reading: %w", err)
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	var allResults []AuditLogEntry

	for scanner.Scan() && len(allResults) < maxResults {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var entry AuditLogEntry
		if s.format == "json" {
			if err := json.Unmarshal([]byte(line), &entry); err != nil {
				logrus.WithError(err).Debug("Failed to parse audit log line")
				continue
			}
		} else {
			// Simple text parsing (basic implementation)
			continue // Skip for now
		}

		if matchesCriteria(&entry, criteria) {
			allResults = append(allResults, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading audit log file: %w", err)
	}

	// Apply pagination
	total := len(allResults)
	start := (page - 1) * pageSize
	if start >= total {
		return []AuditLogEntry{}, nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}

	return allResults[start:end], nil
}

// GetStats returns basic statistics (placeholder implementation)
func (s *FileAuditStorage) GetStats(startTime, endTime time.Time) (map[string]interface{}, error) {
	stats := map[string]interface{}{
		"storage_type": "file",
		"file_path":    s.filePath,
		"start_time":   startTime,
		"end_time":     endTime,
	}

	// Get file info
	if info, err := os.Stat(s.filePath); err == nil {
		stats["file_size"] = info.Size()
		stats["last_modified"] = info.ModTime()
	}

	return stats, nil
}

// rotateIfNeeded checks if file rotation is needed and performs it
func (s *FileAuditStorage) rotateIfNeeded() error {
	if s.maxSizeMB <= 0 {
		return nil // No rotation configured
	}

	info, err := s.currentFile.Stat()
	if err != nil {
		return err
	}

	// Check if file size exceeds limit
	if info.Size() < int64(s.maxSizeMB*1024*1024) {
		return nil // No rotation needed
	}

	// Close current file
	if err := s.currentFile.Close(); err != nil {
		logrus.WithError(err).Error("Failed to close audit log file for rotation")
	}

	// Rotate files
	if err := s.rotateFiles(); err != nil {
		return err
	}

	// Open new file
	file, err := os.OpenFile(s.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open new audit log file: %w", err)
	}
	s.currentFile = file

	return nil
}

// rotateFiles performs the actual file rotation
func (s *FileAuditStorage) rotateFiles() error {
	// Remove oldest backup if we have too many
	for i := s.maxBackups; i >= 1; i-- {
		oldFile := fmt.Sprintf("%s.%d", s.filePath, i)
		newFile := fmt.Sprintf("%s.%d", s.filePath, i+1)

		if i == s.maxBackups {
			// Remove the oldest file
			_ = os.Remove(oldFile)
		} else {
			// Rename file
			if _, err := os.Stat(oldFile); err == nil {
				_ = os.Rename(oldFile, newFile)
			}
		}
	}

	// Move current file to .1
	firstBackup := fmt.Sprintf("%s.1", s.filePath)
	if err := os.Rename(s.filePath, firstBackup); err != nil {
		return fmt.Errorf("failed to rotate audit log file: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"component": "audit",
		"operation": "rotate",
		"file":      s.filePath,
	}).Info("Audit log file rotated")

	return nil
}

// Close closes the audit storage
func (s *FileAuditStorage) Close() error {
	// Prevent nil pointer dereference
	if s == nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.currentFile != nil {
		return s.currentFile.Close()
	}
	return nil
}

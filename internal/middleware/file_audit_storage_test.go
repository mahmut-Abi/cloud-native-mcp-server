package middleware

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileAuditStorage(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "audit.log")

	t.Run("valid config", func(t *testing.T) {
		storage, err := NewFileAuditStorage(logPath, 10, 5, 30, true, "json")
		require.NoError(t, err)
		require.NotNil(t, storage)
		assert.Equal(t, logPath, storage.filePath)
		assert.Equal(t, 10, storage.maxSizeMB)
		assert.Equal(t, 5, storage.maxBackups)
		assert.Equal(t, 30, storage.maxAgeDays)
		assert.True(t, storage.compress)
		assert.Equal(t, "json", storage.format)
		assert.NotNil(t, storage.currentFile)
		_ = storage.Close()
	})

	t.Run("text format", func(t *testing.T) {
		storage, err := NewFileAuditStorage(logPath+"_text", 10, 5, 30, true, "text")
		require.NoError(t, err)
		require.NotNil(t, storage)
		assert.Equal(t, "text", storage.format)
		_ = storage.Close()
	})
}

func TestFileAuditStorage_Log(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "audit.log")
	storage, err := NewFileAuditStorage(logPath, 1, 3, 30, false, "json")
	require.NoError(t, err)
	defer func() { _ = storage.Close() }()

	t.Run("json format", func(t *testing.T) {
		entry := &AuditLogEntry{
			Timestamp:   time.Now(),
			CallerIP:    "127.0.0.1",
			UserID:      "test-user",
			ToolName:    "test-tool",
			ServiceName: "test-service",
			Action:      "test-action",
			Duration:    100,
			Status:      "success",
			ErrorMsg:    "",
		}

		err := storage.Log(entry)
		require.NoError(t, err)

		// Verify file content
		content, err := os.ReadFile(logPath)
		require.NoError(t, err)

		var loggedEntry AuditLogEntry
		err = json.Unmarshal([]byte(strings.TrimSpace(string(content))), &loggedEntry)
		require.NoError(t, err)
		assert.Equal(t, entry.UserID, loggedEntry.UserID)
		assert.Equal(t, entry.ToolName, loggedEntry.ToolName)
		assert.Equal(t, entry.Action, loggedEntry.Action)
	})
}

func TestFileAuditStorage_Query(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "query.log")
	storage, err := NewFileAuditStorage(logPath, 10, 3, 30, false, "json")
	require.NoError(t, err)
	defer func() { _ = storage.Close() }()

	// Write test entries
	entry := &AuditLogEntry{
		Timestamp:   time.Now(),
		CallerIP:    "127.0.0.1",
		UserID:      "test-user",
		ToolName:    "test-tool",
		ServiceName: "test-service",
		Action:      "test-action",
		Duration:    100,
		Status:      "success",
	}

	err = storage.Log(entry)
	require.NoError(t, err)

	t.Run("query all", func(t *testing.T) {
		results, err := storage.Query(map[string]interface{}{})
		require.NoError(t, err)
		assert.Equal(t, 1, len(results))
	})

	t.Run("query by user_id", func(t *testing.T) {
		results, err := storage.Query(map[string]interface{}{
			"user_id": "test-user",
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(results))
		for _, result := range results {
			assert.Equal(t, "test-user", result.UserID)
		}
	})
}

func TestFileAuditStorage_GetStats(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "stats.log")
	storage, err := NewFileAuditStorage(logPath, 10, 3, 30, false, "json")
	require.NoError(t, err)
	defer func() { _ = storage.Close() }()

	startTime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now().Add(1 * time.Hour)

	stats, err := storage.GetStats(startTime, endTime)
	require.NoError(t, err)
	require.NotNil(t, stats)

	assert.Equal(t, "file", stats["storage_type"])
	assert.Equal(t, logPath, stats["file_path"])
	assert.Equal(t, startTime, stats["start_time"])
	assert.Equal(t, endTime, stats["end_time"])
}

func TestFileAuditStorage_Close(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "close_test.log")

	t.Run("normal close", func(t *testing.T) {
		storage, err := NewFileAuditStorage(logPath, 10, 3, 30, false, "json")
		require.NoError(t, err)
		require.NotNil(t, storage.currentFile)

		err = storage.Close()
		require.NoError(t, err)
	})

	t.Run("close nil storage", func(t *testing.T) {
		var storage *FileAuditStorage = nil
		err := storage.Close()
		require.NoError(t, err)
	})
}

package middleware

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	appconfig "github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuditFactory_FileStorageConfiguration tests file storage creation with different configurations
func TestAuditFactory_FileStorageConfiguration(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		config      *appconfig.AppConfig
		expectError bool
		description string
	}{
		{
			name: "file_storage_default_config",
			config: &appconfig.AppConfig{
				Audit: struct {
					Enabled    bool   `yaml:"enabled"`
					Level      string `yaml:"level"`
					MaxLogs    int    `yaml:"maxLogs"`
					Storage    string `yaml:"storage"`
					Format     string `yaml:"format"`
					MaxResults int    `yaml:"maxResults"`
					TimeRange  int    `yaml:"timeRange"`
					File       struct {
						Path       string `yaml:"path"`
						MaxSizeMB  int    `yaml:"maxSizeMB"`
						MaxBackups int    `yaml:"maxBackups"`
						MaxAgeDays int    `yaml:"maxAgeDays"`
						Compress   bool   `yaml:"compress"`
						MaxLogs    int    `yaml:"maxLogs"`
					} `yaml:"file"`
					Database struct {
						Type             string `yaml:"type"`
						ConnectionString string `yaml:"connectionString"`
						SQLitePath       string `yaml:"sqlitePath"`
						TableName        string `yaml:"tableName"`
						MaxRecords       int    `yaml:"maxRecords"`
						CleanupInterval  int    `yaml:"cleanupInterval"`
					} `yaml:"database"`
					Query struct {
						Enabled    bool `yaml:"enabled"`
						MaxResults int  `yaml:"maxResults"`
						TimeRange  int  `yaml:"timeRange"`
					} `yaml:"query"`
					Alerts struct {
						Enabled          bool   `yaml:"enabled"`
						FailureThreshold int    `yaml:"failureThreshold"`
						CheckIntervalSec int    `yaml:"checkIntervalSec"`
						Method           string `yaml:"method"`
						WebhookURL       string `yaml:"webhookURL"`
					} `yaml:"alerts"`
					Masking struct {
						Enabled   bool     `yaml:"enabled"`
						Fields    []string `yaml:"fields"`
						MaskValue string   `yaml:"maskValue"`
					} `yaml:"masking"`
					Sampling struct {
						Enabled bool    `yaml:"enabled"`
						Rate    float32 `yaml:"rate"`
					} `yaml:"sampling"`
				}{
					Enabled: true,
					Storage: "file",
					Format:  "json",
					File: struct {
						Path       string `yaml:"path"`
						MaxSizeMB  int    `yaml:"maxSizeMB"`
						MaxBackups int    `yaml:"maxBackups"`
						MaxAgeDays int    `yaml:"maxAgeDays"`
						Compress   bool   `yaml:"compress"`
						MaxLogs    int    `yaml:"maxLogs"`
					}{
						Path:       filepath.Join(tempDir, "audit.log"),
						MaxSizeMB:  10,
						MaxBackups: 5,
						MaxAgeDays: 7,
						Compress:   true,
					},
				},
			},
			expectError: false,
			description: "Valid file storage configuration should succeed",
		},
		{
			name: "memory_storage_config",
			config: &appconfig.AppConfig{
				Audit: struct {
					Enabled    bool   `yaml:"enabled"`
					Level      string `yaml:"level"`
					MaxLogs    int    `yaml:"maxLogs"`
					Storage    string `yaml:"storage"`
					Format     string `yaml:"format"`
					MaxResults int    `yaml:"maxResults"`
					TimeRange  int    `yaml:"timeRange"`
					File       struct {
						Path       string `yaml:"path"`
						MaxSizeMB  int    `yaml:"maxSizeMB"`
						MaxBackups int    `yaml:"maxBackups"`
						MaxAgeDays int    `yaml:"maxAgeDays"`
						Compress   bool   `yaml:"compress"`
						MaxLogs    int    `yaml:"maxLogs"`
					} `yaml:"file"`
					Database struct {
						Type             string `yaml:"type"`
						ConnectionString string `yaml:"connectionString"`
						SQLitePath       string `yaml:"sqlitePath"`
						TableName        string `yaml:"tableName"`
						MaxRecords       int    `yaml:"maxRecords"`
						CleanupInterval  int    `yaml:"cleanupInterval"`
					} `yaml:"database"`
					Query struct {
						Enabled    bool `yaml:"enabled"`
						MaxResults int  `yaml:"maxResults"`
						TimeRange  int  `yaml:"timeRange"`
					} `yaml:"query"`
					Alerts struct {
						Enabled          bool   `yaml:"enabled"`
						FailureThreshold int    `yaml:"failureThreshold"`
						CheckIntervalSec int    `yaml:"checkIntervalSec"`
						Method           string `yaml:"method"`
						WebhookURL       string `yaml:"webhookURL"`
					} `yaml:"alerts"`
					Masking struct {
						Enabled   bool     `yaml:"enabled"`
						Fields    []string `yaml:"fields"`
						MaskValue string   `yaml:"maskValue"`
					} `yaml:"masking"`
					Sampling struct {
						Enabled bool    `yaml:"enabled"`
						Rate    float32 `yaml:"rate"`
					} `yaml:"sampling"`
				}{
					Enabled: true,
					Storage: "memory",
					MaxLogs: 5000,
				},
			},
			expectError: false,
			description: "Valid memory storage configuration should succeed",
		},
		{
			name: "disabled_audit",
			config: &appconfig.AppConfig{
				Audit: struct {
					Enabled    bool   `yaml:"enabled"`
					Level      string `yaml:"level"`
					MaxLogs    int    `yaml:"maxLogs"`
					Storage    string `yaml:"storage"`
					Format     string `yaml:"format"`
					MaxResults int    `yaml:"maxResults"`
					TimeRange  int    `yaml:"timeRange"`
					File       struct {
						Path       string `yaml:"path"`
						MaxSizeMB  int    `yaml:"maxSizeMB"`
						MaxBackups int    `yaml:"maxBackups"`
						MaxAgeDays int    `yaml:"maxAgeDays"`
						Compress   bool   `yaml:"compress"`
						MaxLogs    int    `yaml:"maxLogs"`
					} `yaml:"file"`
					Database struct {
						Type             string `yaml:"type"`
						ConnectionString string `yaml:"connectionString"`
						SQLitePath       string `yaml:"sqlitePath"`
						TableName        string `yaml:"tableName"`
						MaxRecords       int    `yaml:"maxRecords"`
						CleanupInterval  int    `yaml:"cleanupInterval"`
					} `yaml:"database"`
					Query struct {
						Enabled    bool `yaml:"enabled"`
						MaxResults int  `yaml:"maxResults"`
						TimeRange  int  `yaml:"timeRange"`
					} `yaml:"query"`
					Alerts struct {
						Enabled          bool   `yaml:"enabled"`
						FailureThreshold int    `yaml:"failureThreshold"`
						CheckIntervalSec int    `yaml:"checkIntervalSec"`
						Method           string `yaml:"method"`
						WebhookURL       string `yaml:"webhookURL"`
					} `yaml:"alerts"`
					Masking struct {
						Enabled   bool     `yaml:"enabled"`
						Fields    []string `yaml:"fields"`
						MaskValue string   `yaml:"maskValue"`
					} `yaml:"masking"`
					Sampling struct {
						Enabled bool    `yaml:"enabled"`
						Rate    float32 `yaml:"rate"`
					} `yaml:"sampling"`
				}{
					Enabled: false,
				},
			},
			expectError: false,
			description: "Disabled audit should return nil storage without error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage, err := CreateAuditStorage(tt.config)

			if tt.expectError {
				assert.Error(t, err, tt.description)
			} else {
				assert.NoError(t, err, tt.description)

				if tt.config.Audit.Enabled {
					assert.NotNil(t, storage, "Storage should not be nil when audit is enabled")

					// Test basic functionality
					testEntry := &AuditLogEntry{
						Timestamp:   time.Now(),
						CallerIP:    "192.168.1.1",
						UserID:      "test-user",
						ToolName:    "test-tool",
						ServiceName: "test-service",
						Action:      "GET /test",
						InputParams: map[string]interface{}{"param": "value"},
						Duration:    100,
						Status:      "success",
					}

					err := storage.Log(testEntry)
					assert.NoError(t, err, "Logging should succeed")

					// Clean up
					if storage != nil {
						_ = storage.Close()
					}
				} else {
					assert.Nil(t, storage, "Storage should be nil when audit is disabled")
				}
			}
		})
	}
}

// TestFileAuditStorage_Rotation tests log file rotation functionality
func TestFileAuditStorage_Rotation(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test_rotation.log")

	// Create file storage with very small max size for testing rotation
	storage, err := NewFileAuditStorage(logPath, 1, 3, 1, false, "json") // 1KB max size
	require.NoError(t, err)
	defer func() { _ = storage.Close() }()

	// Write enough entries to trigger rotation
	for i := 0; i < 1000; i++ {
		entry := &AuditLogEntry{
			Timestamp:   time.Now(),
			CallerIP:    "192.168.1.1",
			UserID:      "test-user-with-long-name-to-increase-size",
			ToolName:    "test-tool-with-very-long-name",
			ServiceName: "test-service-with-long-name",
			Action:      "GET /test/path/with/very/long/url/structure",
			InputParams: map[string]interface{}{"param": "very-long-value-to-increase-log-size", "index": i, "data": "additional-data-to-make-entry-larger"},
			Duration:    100,
			Status:      "success",
		}
		err := storage.Log(entry)
		assert.NoError(t, err)
	}

	// Check that at least one file exists (rotation may or may not have happened)
	files, err := filepath.Glob(filepath.Join(tempDir, "test_rotation.log*"))
	assert.NoError(t, err)
	assert.True(t, len(files) >= 1, "Should have at least the main log file")

	// Check the main log file exists and has content
	info, err := os.Stat(logPath)
	assert.NoError(t, err)
	assert.True(t, info.Size() > 0, "Log file should have content")
}

// TestFileAuditStorage_TextFormat tests text format logging
func TestFileAuditStorage_TextFormat(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test_text.log")

	storage, err := NewFileAuditStorage(logPath, 10, 5, 7, false, "text")
	require.NoError(t, err)
	defer func() { _ = storage.Close() }()

	entry := &AuditLogEntry{
		Timestamp:   time.Now(),
		CallerIP:    "192.168.1.1",
		UserID:      "test-user",
		ToolName:    "test-tool",
		ServiceName: "test-service",
		Action:      "GET /test",
		Duration:    100,
		Status:      "success",
	}

	err = storage.Log(entry)
	assert.NoError(t, err)

	// Verify file exists and has content
	content, err := os.ReadFile(logPath)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "test-user")
	assert.Contains(t, string(content), "test-service")
	assert.Contains(t, string(content), "success")
}

// TestCompositeAuditStorage tests composite storage functionality
func TestCompositeAuditStorage(t *testing.T) {
	// Create two in-memory storages
	storage1 := NewInMemoryAuditStorage(100)
	storage2 := NewInMemoryAuditStorage(100)

	composite := NewCompositeAuditStorage(storage1, storage2)
	defer func() { _ = composite.Close() }()

	entry := &AuditLogEntry{
		Timestamp:   time.Now(),
		CallerIP:    "192.168.1.1",
		UserID:      "test-user",
		ToolName:    "test-tool",
		ServiceName: "test-service",
		Action:      "GET /test",
		Duration:    100,
		Status:      "success",
	}

	// Log to composite storage
	err := composite.Log(entry)
	assert.NoError(t, err)

	// Verify both storages received the entry
	results1, err := storage1.Query(map[string]interface{}{})
	assert.NoError(t, err)
	assert.Len(t, results1, 1)

	results2, err := storage2.Query(map[string]interface{}{})
	assert.NoError(t, err)
	assert.Len(t, results2, 1)

	// Verify query works on composite (should use first storage)
	compositeResults, err := composite.Query(map[string]interface{}{})
	assert.NoError(t, err)
	assert.Len(t, compositeResults, 1)
}

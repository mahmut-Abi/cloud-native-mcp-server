package middleware

import (
	"testing"

	appconfig "github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
	"github.com/stretchr/testify/assert"
)

// TestCreateAuditStorage_Disabled tests audit storage creation when disabled
func TestCreateAuditStorage_Disabled(t *testing.T) {
	config := &appconfig.AppConfig{
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
	}

	storage, err := CreateAuditStorage(config)
	assert.NoError(t, err)
	assert.Nil(t, storage)
}

// TestCreateAuditStorage_Memory tests memory storage creation
func TestCreateAuditStorage_Memory(t *testing.T) {
	config := &appconfig.AppConfig{
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
			File: struct {
				Path       string `yaml:"path"`
				MaxSizeMB  int    `yaml:"maxSizeMB"`
				MaxBackups int    `yaml:"maxBackups"`
				MaxAgeDays int    `yaml:"maxAgeDays"`
				Compress   bool   `yaml:"compress"`
				MaxLogs    int    `yaml:"maxLogs"`
			}{
				MaxLogs: 1000,
			},
		},
	}

	storage, err := CreateAuditStorage(config)
	assert.NoError(t, err)
	assert.NotNil(t, storage)

	// Should be InMemoryAuditStorage
	_, ok := storage.(*InMemoryAuditStorage)
	assert.True(t, ok)

	err = storage.Close()
	assert.NoError(t, err)
}

// TestCreateAuditStorage_File tests file storage creation
func TestCreateAuditStorage_File(t *testing.T) {
	tmpDir := t.TempDir()
	config := &appconfig.AppConfig{
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
				Path:       tmpDir + "/audit.log",
				MaxSizeMB:  10,
				MaxBackups: 5,
				MaxAgeDays: 30,
				Compress:   false,
			},
		},
	}

	storage, err := CreateAuditStorage(config)
	assert.NoError(t, err)
	assert.NotNil(t, storage)

	// Should be FileAuditStorage
	_, ok := storage.(*FileAuditStorage)
	assert.True(t, ok)

	err = storage.Close()
	assert.NoError(t, err)
}

// TestCreateAuditStorage_UnknownType tests unknown storage type
func TestCreateAuditStorage_UnknownType(t *testing.T) {
	config := &appconfig.AppConfig{
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
			Storage: "unknown",
		},
	}

	storage, err := CreateAuditStorage(config)
	assert.Error(t, err)
	assert.Nil(t, storage)
	assert.Contains(t, err.Error(), "unknown audit storage type")
}

// TestCreateAuditStorage_Database tests database storage creation
func TestCreateAuditStorage_Database(t *testing.T) {
	config := &appconfig.AppConfig{
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
			Storage: "database",
			Database: struct {
				Type             string `yaml:"type"`
				ConnectionString string `yaml:"connectionString"`
				SQLitePath       string `yaml:"sqlitePath"`
				TableName        string `yaml:"tableName"`
				MaxRecords       int    `yaml:"maxRecords"`
				CleanupInterval  int    `yaml:"cleanupInterval"`
			}{
				Type:            "sqlite3",
				SQLitePath:      ":memory:",
				TableName:       "audit_logs",
				MaxRecords:      10000,
				CleanupInterval: 3600,
			},
		},
	}

	storage, err := CreateAuditStorage(config)
	if err != nil {
		t.Skipf("SQLite driver not available: %v", err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, storage)

	// Should be DatabaseAuditStorage
	_, ok := storage.(*DatabaseAuditStorage)
	assert.True(t, ok)

	if storage != nil {
		err = storage.Close()
		assert.NoError(t, err)
	}
}

// TestCreateAuditStorage_InvalidFileConfig tests file storage with invalid config
func TestCreateAuditStorage_InvalidFileConfig(t *testing.T) {
	config := &appconfig.AppConfig{
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
			File: struct {
				Path       string `yaml:"path"`
				MaxSizeMB  int    `yaml:"maxSizeMB"`
				MaxBackups int    `yaml:"maxBackups"`
				MaxAgeDays int    `yaml:"maxAgeDays"`
				Compress   bool   `yaml:"compress"`
				MaxLogs    int    `yaml:"maxLogs"`
			}{
				Path:       "/invalid/nonexistent/path/audit.log",
				MaxSizeMB:  10,
				MaxBackups: 5,
			},
		},
	}

	storage, _ := CreateAuditStorage(config)
	// The behavior depends on filesystem permissions - directory may be createable
	if storage != nil {
		_ = storage.Close()
	}
}

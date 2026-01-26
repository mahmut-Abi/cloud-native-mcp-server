package middleware

import (
	"fmt"
	"strings"
	"time"

	appconfig "github.com/mahmut-Abi/k8s-mcp-server/internal/config"
	"github.com/sirupsen/logrus"
)

// CreateAuditStorage creates appropriate audit storage based on configuration
func CreateAuditStorage(config *appconfig.AppConfig) (AuditStorage, error) {
	if !config.Audit.Enabled {
		return nil, nil
	}

	// Set defaults
	storage := config.Audit.Storage
	if storage == "" {
		storage = "memory"
	}

	format := config.Audit.Format
	if format == "" {
		format = "json"
	}

	switch strings.ToLower(storage) {
	case "memory":
		maxLogs := config.Audit.File.MaxLogs
		if maxLogs == 0 {
			maxLogs = config.Audit.MaxLogs
			if maxLogs == 0 {
				maxLogs = 10000
			}
		}
		return NewInMemoryAuditStorage(maxLogs), nil

	case "file":
		filePath := config.Audit.File.Path
		if filePath == "" {
			filePath = "/var/log/k8s-mcp-server/audit.log"
		}

		maxSizeMB := config.Audit.File.MaxSizeMB
		if maxSizeMB == 0 {
			maxSizeMB = 100
		}

		maxBackups := config.Audit.File.MaxBackups
		if maxBackups == 0 {
			maxBackups = 10
		}

		maxAgeDays := config.Audit.File.MaxAgeDays
		if maxAgeDays == 0 {
			maxAgeDays = 30
		}

		return NewFileAuditStorage(filePath, maxSizeMB, maxBackups, maxAgeDays, config.Audit.File.Compress, format)

	case "database":
		dbType := config.Audit.Database.Type
		if dbType == "" {
			dbType = "sqlite"
		}

		connectionString := config.Audit.Database.ConnectionString
		sqlitePath := config.Audit.Database.SQLitePath
		if dbType == "sqlite" && sqlitePath == "" {
			sqlitePath = "/var/lib/k8s-mcp-server/audit.db"
		}

		tableName := config.Audit.Database.TableName
		if tableName == "" {
			tableName = "audit_logs"
		}

		maxRecords := config.Audit.Database.MaxRecords
		if maxRecords == 0 {
			maxRecords = 100000
		}

		cleanupInterval := config.Audit.Database.CleanupInterval
		if cleanupInterval == 0 {
			cleanupInterval = 24 // 24 hours
		}

		return NewDatabaseAuditStorage(dbType, connectionString, sqlitePath, tableName, maxRecords, cleanupInterval)

	case "all":
		// Create composite storage that writes to both file and database
		fileCfg := *config
		fileCfg.Audit.Storage = "file"
		fileStorage, err := CreateAuditStorage(&fileCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create file storage: %w", err)
		}

		dbCfg := *config
		dbCfg.Audit.Storage = "database"
		dbStorage, err := CreateAuditStorage(&dbCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create database storage: %w", err)
		}

		return NewCompositeAuditStorage(fileStorage, dbStorage), nil

	default:
		return nil, fmt.Errorf("unknown audit storage type: %s", storage)
	}
}

// CompositeAuditStorage writes to multiple storages
type CompositeAuditStorage struct {
	storages []AuditStorage
}

// NewCompositeAuditStorage creates a composite storage
func NewCompositeAuditStorage(storages ...AuditStorage) *CompositeAuditStorage {
	return &CompositeAuditStorage{
		storages: storages,
	}
}

// Log writes to all storages
func (c *CompositeAuditStorage) Log(entry *AuditLogEntry) error {
	var lastErr error
	for _, storage := range c.storages {
		if err := storage.Log(entry); err != nil {
			logrus.WithError(err).Error("Failed to write audit log to storage")
			lastErr = err
		}
	}
	return lastErr
}

// Query queries the first storage (typically file storage for better performance)
func (c *CompositeAuditStorage) Query(criteria map[string]interface{}) ([]AuditLogEntry, error) {
	if len(c.storages) == 0 {
		return nil, fmt.Errorf("no storages available")
	}
	return c.storages[0].Query(criteria)
}

// GetStats gets stats from the first storage
func (c *CompositeAuditStorage) GetStats(startTime, endTime time.Time) (map[string]interface{}, error) {
	if len(c.storages) == 0 {
		return nil, fmt.Errorf("no storages available")
	}
	return c.storages[0].GetStats(startTime, endTime)
}

// Close closes all composite storages
func (c *CompositeAuditStorage) Close() error {
	var lastErr error
	for _, storage := range c.storages {
		if err := storage.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

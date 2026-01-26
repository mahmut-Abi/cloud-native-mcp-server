package handlers

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// GrafanaCacheConfig cache configuration
type GrafanaCacheConfig struct {
	DefaultTTL    time.Duration // Default TTL
	MaxEntries    int           // Maximum cache entries
	MaxMemorySize int64         // Maximum memory usage (bytes)
	CleanupDelay  time.Duration // Cleanup delay
}

// GrafanaCacheEntry cache entry
type GrafanaCacheEntry struct {
	Data      interface{}
	CreatedAt time.Time
	ExpiresAt time.Time
	Size      int64
	HitCount  int64
	Metadata  map[string]interface{}
}

// GrafanaSmartCache smart cache system - optimized specifically for Grafana
type GrafanaSmartCache struct {
	entries     map[string]*GrafanaCacheEntry
	mutex       sync.RWMutex
	config      GrafanaCacheConfig
	memoryUsed  int64
	currentSize int
}

// NewGrafanaSmartCache creates a new Grafana smart cache
func NewGrafanaSmartCache(config GrafanaCacheConfig) *GrafanaSmartCache {
	cache := &GrafanaSmartCache{
		entries: make(map[string]*GrafanaCacheEntry),
		config:  config,
	}

	// Start periodic cleanup
	go cache.startCleanupRoutine()

	return cache
}

// DefaultGrafanaCacheConfig default Grafana cache configuration
var DefaultGrafanaCacheConfig = GrafanaCacheConfig{
	DefaultTTL:    3 * time.Minute,  // 3 minute default TTL - Grafana data is relatively stable
	MaxEntries:    500,              // Grafana has less data, 500 entries is enough
	MaxMemorySize: 30 * 1024 * 1024, // 30MB max memory - Grafana has many configurations
	CleanupDelay:  2 * time.Minute,  // 2 minute cleanup interval
}

// Get retrieves cached data
func (cache *GrafanaSmartCache) Get(key string) (interface{}, bool) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	entry, exists := cache.entries[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		delete(cache.entries, key)
		cache.currentSize--
		cache.memoryUsed -= entry.Size
		return nil, false
	}

	// Update hit statistics
	entry.HitCount++

	logrus.WithFields(logrus.Fields{
		"cacheKey": key[:16] + "...",
		"hitCount": entry.HitCount,
		"age":      time.Since(entry.CreatedAt),
	}).Debug("Grafana cache hit")

	return entry.Data, true
}

// Set sets cached data
func (cache *GrafanaSmartCache) Set(key string, data interface{}, ttl time.Duration) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// Check if space needs to be cleaned
	if cache.needsEviction() {
		cache.evictLRU()
	}

	// Serialize data to calculate size
	dataBytes, err := json.Marshal(data)
	if err != nil {
		logrus.WithError(err).Error("Failed to serialize Grafana cache data")
		return
	}

	// If data is too large, don't cache - Grafana has stricter limits
	if int64(len(dataBytes)) > cache.config.MaxMemorySize/20 { // Single entry not exceeding 1/20 of total memory
		logrus.WithField("size", len(dataBytes)).Warn("Grafana data too large for cache, skipping")
		return
	}

	// Delete existing entry (if exists)
	if existingEntry, exists := cache.entries[key]; exists {
		cache.memoryUsed -= existingEntry.Size
	} else {
		cache.currentSize++
	}

	// Create new entry
	now := time.Now()
	if ttl == 0 {
		ttl = cache.config.DefaultTTL
	}

	entry := &GrafanaCacheEntry{
		Data:      data,
		CreatedAt: now,
		ExpiresAt: now.Add(ttl),
		Size:      int64(len(dataBytes)),
		HitCount:  0,
		Metadata: map[string]interface{}{
			"setAt": now,
		},
	}

	cache.entries[key] = entry
	cache.memoryUsed += entry.Size

	logrus.WithFields(logrus.Fields{
		"cacheKey":     key[:16] + "...",
		"size":         entry.Size,
		"ttl":          ttl,
		"totalEntries": cache.currentSize,
		"memoryUsed":   cache.memoryUsed,
	}).Debug("Grafana data cached")
}

// needsEviction checks if cache needs to be cleaned
func (cache *GrafanaSmartCache) needsEviction() bool {
	return cache.currentSize >= cache.config.MaxEntries ||
		cache.memoryUsed >= cache.config.MaxMemorySize
}

// evictLRU cleans cache using LRU strategy
func (cache *GrafanaSmartCache) evictLRU() {
	var oldestKey string
	var oldestTime time.Time

	// Find the oldest entry
	for key, entry := range cache.entries {
		if oldestKey == "" || entry.CreatedAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.CreatedAt
		}
	}

	if oldestKey != "" {
		if entry, exists := cache.entries[oldestKey]; exists {
			cache.memoryUsed -= entry.Size
			cache.currentSize--
			delete(cache.entries, oldestKey)

			logrus.WithFields(logrus.Fields{
				"cacheKey": oldestKey[:16] + "...",
				"age":      time.Since(entry.CreatedAt),
				"reason":   "LRU eviction",
			}).Debug("Grafana cache entry evicted")
		}
	}
}

// startCleanupRoutine starts periodic cleanup routine
func (cache *GrafanaSmartCache) startCleanupRoutine() {
	ticker := time.NewTicker(cache.config.CleanupDelay)
	defer ticker.Stop()

	for range ticker.C {
		cache.cleanup()
	}
}

// cleanup cleans up expired entries
func (cache *GrafanaSmartCache) cleanup() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	now := time.Now()
	expiredKeys := make([]string, 0)

	for key, entry := range cache.entries {
		if now.After(entry.ExpiresAt) {
			expiredKeys = append(expiredKeys, key)
		}
	}

	for _, key := range expiredKeys {
		if entry, exists := cache.entries[key]; exists {
			cache.memoryUsed -= entry.Size
			cache.currentSize--
			delete(cache.entries, key)

			logrus.WithFields(logrus.Fields{
				"cacheKey": key[:16] + "...",
				"age":      time.Since(entry.CreatedAt),
				"reason":   "expired",
			}).Debug("Grafana cache entry cleaned up")
		}
	}

	if len(expiredKeys) > 0 {
		logrus.WithFields(logrus.Fields{
			"cleanedEntries":   len(expiredKeys),
			"remainingEntries": cache.currentSize,
			"memoryUsed":       cache.memoryUsed,
		}).Info("Grafana cache cleanup completed")
	}
}

// GetStats gets cache statistics
func (cache *GrafanaSmartCache) GetStats() map[string]interface{} {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	totalHits := int64(0)
	var totalAge time.Duration
	validEntries := 0

	for _, entry := range cache.entries {
		totalHits += entry.HitCount
		totalAge += time.Since(entry.CreatedAt)
		validEntries++
	}

	var avgAge time.Duration
	if validEntries > 0 {
		avgAge = totalAge / time.Duration(validEntries)
	}

	return map[string]interface{}{
		"totalEntries":      cache.currentSize,
		"maxEntries":        cache.config.MaxEntries,
		"memoryUsed":        cache.memoryUsed,
		"maxMemorySize":     cache.config.MaxMemorySize,
		"totalHits":         totalHits,
		"averageAge":        avgAge.String(),
		"memoryUtilization": fmt.Sprintf("%.2f%%", float64(cache.memoryUsed)/float64(cache.config.MaxMemorySize)*100),
		"entryUtilization":  fmt.Sprintf("%.2f%%", float64(cache.currentSize)/float64(cache.config.MaxEntries)*100),
	}
}

// Clear clears the cache
func (cache *GrafanaSmartCache) Clear() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.memoryUsed = 0
	cache.currentSize = 0
	cache.entries = make(map[string]*GrafanaCacheEntry)

	logrus.Info("Grafana cache cleared")
}

// Global Grafana cache instance
var DefaultGrafanaSmartCache = NewGrafanaSmartCache(DefaultGrafanaCacheConfig)

// Grafana tool-specific TTL configuration
var GrafanaCacheTTLByTool = map[string]time.Duration{
	"grafana_dashboards_summary":  2 * time.Minute,  // Dashboard summary cache is shorter
	"grafana_get_dashboards":      2 * time.Minute,  // Dashboard list cache is shorter
	"grafana_datasources_summary": 5 * time.Minute,  // Data sources are usually stable
	"grafana_get_datasources":     5 * time.Minute,  // Data source configuration is stable
	"grafana_folders":             10 * time.Minute, // Folders are very stable
	"grafana_alerts":              1 * time.Minute,  // Alerts may change frequently
	"grafana_dashboard":           1 * time.Minute,  // Single dashboard may change
	"grafana_search_dashboards":   2 * time.Minute,  // Search results are relatively stable
	"grafana_test_connection":     30 * time.Second, // Connection test results are cached briefly
}

// GetTTLForGrafanaTool gets the cache TTL for Grafana tool
func GetTTLForGrafanaTool(toolName string) time.Duration {
	if ttl, exists := GrafanaCacheTTLByTool[toolName]; exists {
		return ttl
	}
	return DefaultGrafanaCacheConfig.DefaultTTL
}

// IsGrafanaToolCacheable checks if Grafana tool should be cached
func IsGrafanaToolCacheable(toolName string) bool {
	// These tools have stable results and are suitable for caching
	cacheableTools := map[string]bool{
		"grafana_dashboards_summary":  true,
		"grafana_get_dashboards":      true,
		"grafana_datasources_summary": true,
		"grafana_get_datasources":     true,
		"grafana_folders":             true,
		"grafana_search_dashboards":   true, // Short-term cache
		"grafana_test_connection":     true, // Very short cache
	}

	return cacheableTools[toolName]
}

// GrafanaCacheParamsFilter filters parameters to create more effective cache keys
func GrafanaCacheParamsFilter(toolName string, params map[string]interface{}) map[string]interface{} {
	filtered := make(map[string]interface{})

	// Select important cache parameters based on tool type
	switch toolName {
	case "grafana_dashboards", "grafana_dashboards_summary", "grafana_get_dashboards":
		filtered["limit"] = params["limit"]
		filtered["offset"] = params["offset"]
		// debug parameter doesn't affect result, so not included

	case "grafana_datasources", "grafana_datasources_summary", "grafana_get_datasources":
		filtered["limit"] = params["limit"]
		// debug parameter doesn't affect result

	case "grafana_folders", "grafana_get_folders":
		filtered["limit"] = params["limit"]

	case "grafana_alerts", "grafana_get_alerts":
		filtered["limit"] = params["limit"]

	case "grafana_search_dashboards":
		filtered["query"] = params["query"]
		filtered["tag"] = params["tag"]
		filtered["folderUID"] = params["folderUID"]
		filtered["starred"] = params["starred"]
		filtered["limit"] = params["limit"]

	case "grafana_dashboard", "grafana_get_dashboard":
		filtered["uid"] = params["uid"]
		// debug parameter doesn't affect result

	case "grafana_test_connection":
		// Connection test has no parameters that affect result

	default:
		return params
	}

	return filtered
}

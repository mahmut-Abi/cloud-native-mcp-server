package handlers

import (
	"crypto/md5"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
)

// CacheEntry cache entry
type CacheEntry struct {
	Data      interface{}
	CreatedAt time.Time
	ExpiresAt time.Time
	Size      int64
	HitCount  int64
	Metadata  map[string]interface{}
}

// CacheConfig cache configuration
type CacheConfig struct {
	DefaultTTL    time.Duration // Default TTL
	MaxEntries    int           // Maximum cache entries
	MaxMemorySize int64         // Maximum memory usage (bytes)
	CleanupDelay  time.Duration // Cleanup delay
}

// SmartCache smart cache system
type SmartCache struct {
	entries     map[string]*CacheEntry
	mutex       sync.RWMutex
	config      CacheConfig
	memoryUsed  int64
	currentSize int
}

// NewSmartCache creates new smart cache
func NewSmartCache(config CacheConfig) *SmartCache {
	cache := &SmartCache{
		entries: make(map[string]*CacheEntry),
		config:  config,
	}

	// Start periodic cleanup
	go cache.startCleanupRoutine()

	return cache
}

// DefaultCacheConfig default cache configuration
var DefaultCacheConfig = CacheConfig{
	DefaultTTL:    5 * time.Minute,  // 5 minutes default TTL
	MaxEntries:    1000,             // Maximum 1000 entries
	MaxMemorySize: 50 * 1024 * 1024, // 50MB maximum memory
	CleanupDelay:  1 * time.Minute,  // 1 minute cleanup interval
}

// generateCacheKey generates cache key - optimized version
func (cache *SmartCache) generateCacheKey(toolName string, params map[string]interface{}) string {
	// For simple parameters, use string concatenation to avoid JSON serialization overhead
	if len(params) == 0 {
		return toolName
	}

	// For small number of simple type parameters, use efficient string concatenation
	if len(params) <= 3 && cache.isSimpleParams(params) {
		return fmt.Sprintf("%s:%s", toolName, cache.buildSimpleParams(params))
	}

	// For complex parameters, keep MD5 hash but use more efficient serialization
	keyData := map[string]interface{}{
		"tool":   toolName,
		"params": params,
	}

	keyBytes, _ := optimize.GlobalJSONPool.MarshalToBytes(keyData)
	return fmt.Sprintf("%x", md5.Sum(keyBytes))
}

// isSimpleParams checks if parameters are all simple types
func (cache *SmartCache) isSimpleParams(params map[string]interface{}) bool {
	for _, v := range params {
		switch v.(type) {
		case string, int, int32, int64, bool, float32, float64:
			continue
		default:
			return false
		}
	}
	return true
}

// buildSimpleParams builds string representation of simple parameters
func (cache *SmartCache) buildSimpleParams(params map[string]interface{}) string {
	var parts []string
	for k, v := range params {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	return strings.Join(parts, "&")
}

// Get retrieves cached data
func (cache *SmartCache) Get(key string) (interface{}, bool) {
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
	}).Debug("Cache hit")

	return entry.Data, true
}

// Set sets cached data
func (cache *SmartCache) Set(key string, data interface{}, ttl time.Duration) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// Check if eviction is needed
	if cache.needsEviction() {
		cache.evictLRU()
	}

	// Serialize data to calculate size
	dataBytes, err := optimize.GlobalJSONPool.MarshalToBytes(data)
	if err != nil {
		logrus.WithError(err).Error("Failed to serialize cache data")
		return
	}

	// If data is too large, don't cache
	if int64(len(dataBytes)) > cache.config.MaxMemorySize/10 { // Single entry not exceeding 1/10 of total memory
		logrus.WithField("size", len(dataBytes)).Warn("Data too large for cache, skipping")
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

	entry := &CacheEntry{
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
	}).Debug("Data cached")
}

// needsEviction checks if cache eviction is needed
func (cache *SmartCache) needsEviction() bool {
	return cache.currentSize >= cache.config.MaxEntries ||
		cache.memoryUsed >= cache.config.MaxMemorySize
}

// evictLRU uses LRU strategy to evict cache
func (cache *SmartCache) evictLRU() {
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
			}).Debug("Cache entry evicted")
		}
	}
}

// startCleanupRoutine starts periodic cleanup routine
func (cache *SmartCache) startCleanupRoutine() {
	ticker := time.NewTicker(cache.config.CleanupDelay)
	defer ticker.Stop()

	for range ticker.C {
		cache.cleanup()
	}
}

// cleanup cleans up expired entries
func (cache *SmartCache) cleanup() {
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
			}).Debug("Cache entry cleaned up")
		}
	}

	if len(expiredKeys) > 0 {
		logrus.WithFields(logrus.Fields{
			"cleanedEntries":   len(expiredKeys),
			"remainingEntries": cache.currentSize,
			"memoryUsed":       cache.memoryUsed,
		}).Info("Cache cleanup completed")
	}
}

// GetStats gets cache statistics
func (cache *SmartCache) GetStats() map[string]interface{} {
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

// Clear clears cache
func (cache *SmartCache) Clear() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.memoryUsed = 0
	cache.currentSize = 0
	cache.entries = make(map[string]*CacheEntry)

	logrus.Info("Cache cleared")
}

// CacheableResponse cacheable response format
type CacheableResponse struct {
	Data      interface{}            `json:"data"`
	CacheInfo map[string]interface{} `json:"cache,omitempty"`
}

// CreateCacheableResponse creates cacheable response
func CreateCacheableResponse(data interface{}, toolName string, cacheHit bool, cacheInfo map[string]interface{}) interface{} {
	if cacheHit && cacheInfo != nil {
		return CacheableResponse{
			Data: data,
			CacheInfo: map[string]interface{}{
				"source":   "cache",
				"tool":     toolName,
				"cachedAt": cacheInfo["cachedAt"],
				"ttl":      cacheInfo["ttl"],
				"hitCount": cacheInfo["hitCount"],
			},
		}
	}

	return data
}

// CacheToolResponse tool response cache decorator
func CacheToolResponse(cache *SmartCache, toolName string, params map[string]interface{},
	execFunc func() (interface{}, error), ttl time.Duration) (interface{}, bool, error) {

	// Generate cache key
	cacheKey := cache.generateCacheKey(toolName, params)

	// Try to get from cache
	if cachedData, found := cache.Get(cacheKey); found {
		logrus.WithFields(logrus.Fields{
			"tool":     toolName,
			"cacheKey": cacheKey[:16] + "...",
		}).Info("Cache hit for tool response")

		// Cache info can be used in future versions

		return cachedData, true, nil
	}

	// Execute original function
	logrus.WithFields(logrus.Fields{
		"tool":     toolName,
		"cacheKey": cacheKey[:16] + "...",
	}).Debug("Cache miss, executing function")

	data, err := execFunc()
	if err != nil {
		return nil, false, err
	}

	// Cache result
	cache.Set(cacheKey, data, ttl)

	return data, false, nil
}

// Global cache instance
var DefaultSmartCache = NewSmartCache(DefaultCacheConfig)

// Cache-friendly tool TTL configuration
var CacheTTLByTool = map[string]time.Duration{
	"kubernetes_list_resources_summary": 2 * time.Minute,  // Resource list cache shorter
	"kubernetes_get_resource_summary":   3 * time.Minute,  // Single resource slightly longer
	"kubernetes_get_recent_events":      1 * time.Minute,  // Event cache shortest
	"kubernetes_get_api_versions":       30 * time.Minute, // API version information stable
	"kubernetes_get_api_resources":      15 * time.Minute, // API resources stable
	"kubernetes_check_permissions":      10 * time.Minute, // Permission check moderate
}

// GetTTLForTool gets tool cache TTL
func GetTTLForTool(toolName string) time.Duration {
	if ttl, exists := CacheTTLByTool[toolName]; exists {
		return ttl
	}
	return DefaultCacheConfig.DefaultTTL
}

// IsToolCacheable checks if tool should be cached
func IsToolCacheable(toolName string) bool {
	// These tools have stable results, suitable for caching
	cacheableTools := map[string]bool{
		"kubernetes_list_resources_summary": true,
		"kubernetes_get_resource_summary":   true,
		"kubernetes_get_api_versions":       true,
		"kubernetes_get_api_resources":      true,
		"kubernetes_check_permissions":      true,
		"kubernetes_get_recent_events":      true, // Short-term cache
	}

	return cacheableTools[toolName]
}

// CacheParamsFilter filters parameters to create more effective cache keys
func CacheParamsFilter(toolName string, params map[string]interface{}) map[string]interface{} {
	filtered := make(map[string]interface{})

	// Select important cache parameters based on tool type
	switch toolName {
	case "kubernetes_list_resources_summary", "kubernetes_list_resources", "kubernetes_list_resources_full":
		filtered["kind"] = params["kind"]
		filtered["namespace"] = params["namespace"]
		filtered["labelSelector"] = params["labelSelector"]
		filtered["fieldSelector"] = params["fieldSelector"]
		// limit parameter not included in cache key because different limits but same data

	case "kubernetes_get_resource_summary", "kubernetes_get_resource":
		filtered["kind"] = params["kind"]
		filtered["name"] = params["name"]
		filtered["namespace"] = params["namespace"]

	case "kubernetes_get_recent_events", "kubernetes_get_events", "kubernetes_get_events_detail":
		filtered["namespace"] = params["namespace"]
		filtered["fieldSelector"] = params["fieldSelector"]
		// Parameter includeNormalEvents affects result, so included

	case "kubernetes_get_pod_logs":
		// Logs not suitable for caching because they change frequently
		return filtered

	default:
		return params
	}

	return filtered
}

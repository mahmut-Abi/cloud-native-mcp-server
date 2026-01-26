// Package cache provides caching functionality for Grafana client operations.
package cache

import (
	"fmt"
	"sync"
	"time"
)

// Entry represents a cached item with expiration time.
type Entry struct {
	Value     interface{}
	ExpiresAt time.Time
}

// IsExpired checks if the cache entry has expired.
func (e *Entry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// Cache provides thread-safe caching with TTL support.
type Cache struct {
	mu         sync.RWMutex
	data       map[string]*Entry
	maxSize    int
	defaultTTL time.Duration
}

// NewCache creates a new cache with specified max size and default TTL.
func NewCache(maxSize int, defaultTTL time.Duration) *Cache {
	if maxSize <= 0 {
		maxSize = 100
	}
	if defaultTTL <= 0 {
		defaultTTL = 5 * time.Minute
	}

	return &Cache{
		data:       make(map[string]*Entry),
		maxSize:    maxSize,
		defaultTTL: defaultTTL,
	}
}

// Set stores a value in cache with default TTL.
func (c *Cache) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

// SetWithTTL stores a value in cache with custom TTL.
func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if we need to evict entries
	if len(c.data) >= c.maxSize {
		c.evictOne()
	}

	c.data[key] = &Entry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Get retrieves a value from cache.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.data[key]
	if !exists {
		return nil, false
	}

	if entry.IsExpired() {
		return nil, false
	}

	return entry.Value, true
}

// Delete removes a key from cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}

// Clear removes all entries from cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]*Entry)
}

// evictOne removes one expired entry, or the oldest if no expired entries.
func (c *Cache) evictOne() {
	// First try to remove expired entries
	for key, entry := range c.data {
		if entry.IsExpired() {
			delete(c.data, key)
			return
		}
	}

	// If no expired entries, remove any entry (simple eviction)
	for key := range c.data {
		delete(c.data, key)
		return
	}
}

// Size returns the number of entries in cache.
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.data)
}

// CacheKey generates a cache key for Grafana resources.
func CacheKey(resourceType, id string) string {
	return fmt.Sprintf("%s:%s", resourceType, id)
}

package cache

import (
	"context"
	"hash/fnv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/constants"
)

// SegmentedCache implements Cache interface with segmented locking for better concurrency
type SegmentedCache struct {
	segments   []*cacheSegment
	segmentNum int
	hits       atomic.Int64
	misses     atomic.Int64
	evictions  atomic.Int64
	maxSize    int
}

type cacheSegment struct {
	mu    sync.RWMutex
	store map[string]CacheEntry
}

const defaultSegmentCount = 16 // Powers of 2 are optimal for hash distribution

// NewSegmentedCache creates a new segmented cache for improved concurrency
func NewSegmentedCache() *SegmentedCache {
	cache := &SegmentedCache{
		segments:   make([]*cacheSegment, defaultSegmentCount),
		segmentNum: defaultSegmentCount,
		maxSize:    constants.DefaultCacheSize,
	}

	for i := 0; i < defaultSegmentCount; i++ {
		cache.segments[i] = &cacheSegment{
			store: make(map[string]CacheEntry),
		}
	}

	return cache
}

// getSegment returns the segment for a given key
func (c *SegmentedCache) getSegment(key string) *cacheSegment {
	hash := fnv.New64a()
	hash.Write([]byte(key))
	index := int(hash.Sum64() % uint64(c.segmentNum))
	return c.segments[index]
}

// Get retrieves a value from cache
func (c *SegmentedCache) Get(ctx context.Context, key string) (interface{}, bool) {
	segment := c.getSegment(key)
	segment.mu.RLock()
	entry, exists := segment.store[key]
	segment.mu.RUnlock()

	if !exists {
		c.misses.Add(1)
		return nil, false
	}

	// Check expiration
	if time.Now().After(entry.ExpiresAt) {
		c.misses.Add(1)
		// Clean up expired entry
		segment.mu.Lock()
		delete(segment.store, key)
		c.evictions.Add(1)
		segment.mu.Unlock()
		return nil, false
	}

	c.hits.Add(1)
	return entry.Value, true
}

// Set stores a value in cache with TTL
func (c *SegmentedCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	segment := c.getSegment(key)
	segment.mu.Lock()

	segment.store[key] = CacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}

	// Check if we need to evict old entries from this segment
	if len(segment.store) > (c.maxSize / c.segmentNum) {
		c.evictOldestInSegment(segment)
	}

	segment.mu.Unlock()
}

// evictOldestInSegment removes the oldest entry from a segment
func (c *SegmentedCache) evictOldestInSegment(segment *cacheSegment) {
	var oldestKey string
	oldestTime := time.Now()

	for key, entry := range segment.store {
		if entry.ExpiresAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.ExpiresAt
		}
	}

	if oldestKey != "" {
		delete(segment.store, oldestKey)
		c.evictions.Add(1)
	}
}

// Delete removes a value from cache
func (c *SegmentedCache) Delete(ctx context.Context, key string) {
	segment := c.getSegment(key)
	segment.mu.Lock()
	delete(segment.store, key)
	segment.mu.Unlock()
}

// Clear removes all entries from cache
func (c *SegmentedCache) Clear(ctx context.Context) {
	for _, segment := range c.segments {
		segment.mu.Lock()
		segment.store = make(map[string]CacheEntry)
		segment.mu.Unlock()
	}
}

// GetStats returns cache statistics
func (c *SegmentedCache) GetStats() CacheStats {
	return CacheStats{
		Hits:      c.hits.Load(),
		Misses:    c.misses.Load(),
		Evictions: c.evictions.Load(),
	}
}

// Cleanup removes expired entries from all segments
func (c *SegmentedCache) Cleanup(ctx context.Context) {
	now := time.Now()

	for _, segment := range c.segments {
		segment.mu.Lock()

		var keysToDelete []string
		for key, entry := range segment.store {
			if now.After(entry.ExpiresAt) {
				keysToDelete = append(keysToDelete, key)
			}
		}

		for _, key := range keysToDelete {
			delete(segment.store, key)
			c.evictions.Add(1)
		}

		segment.mu.Unlock()
	}
}

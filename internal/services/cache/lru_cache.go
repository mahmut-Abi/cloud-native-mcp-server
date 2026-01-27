package cache

import (
	"container/list"
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/constants"
)

// LRUCacheEntry represents an entry in the LRU cache
type LRUCacheEntry struct {
	Key       string
	Value     interface{}
	ExpiresAt time.Time
}

// LRUMemoryCache implements Cache interface with LRU eviction policy
type LRUMemoryCache struct {
	mu        sync.Mutex
	list      *list.List
	map_      map[string]*list.Element
	maxSize   int
	hits      atomic.Int64
	misses    atomic.Int64
	evictions atomic.Int64
}

// NewLRUMemoryCache creates a new LRU in-memory cache
func NewLRUMemoryCache() *LRUMemoryCache {
	return &LRUMemoryCache{
		list:    list.New(),
		map_:    make(map[string]*list.Element, 1024),
		maxSize: constants.DefaultCacheSize,
	}
}

// Get retrieves a value from cache with O(1) operation
func (c *LRUMemoryCache) Get(ctx context.Context, key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, exists := c.map_[key]
	if !exists {
		c.misses.Add(1)
		return nil, false
	}

	entry := elem.Value.(*LRUCacheEntry)
	if time.Now().After(entry.ExpiresAt) {
		c.list.Remove(elem)
		delete(c.map_, key)
		c.misses.Add(1)
		return nil, false
	}

	// Move to front (most recently used)
	c.list.MoveToFront(elem)
	c.hits.Add(1)
	return entry.Value, true
}

// Set stores a value in cache with O(1) operation
func (c *LRUMemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If key already exists, update it
	if elem, exists := c.map_[key]; exists {
		elem.Value = &LRUCacheEntry{
			Key:       key,
			Value:     value,
			ExpiresAt: time.Now().Add(ttl),
		}
		c.list.MoveToFront(elem)
		return
	}

	// Add new entry
	entry := &LRUCacheEntry{
		Key:       key,
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
	elem := c.list.PushFront(entry)
	c.map_[key] = elem

	// Evict least recently used if cache is full
	if len(c.map_) > c.maxSize {
		c.evictLRU()
	}
}

// evictLRU removes the least recently used entry
func (c *LRUMemoryCache) evictLRU() {
	if lruElem := c.list.Back(); lruElem != nil {
		c.list.Remove(lruElem)
		lruEntry := lruElem.Value.(*LRUCacheEntry)
		delete(c.map_, lruEntry.Key)
		c.evictions.Add(1)
	}
}

// Delete removes a value from cache
func (c *LRUMemoryCache) Delete(ctx context.Context, key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, exists := c.map_[key]; exists {
		c.list.Remove(elem)
		delete(c.map_, key)
	}
}

// Clear removes all entries from cache
func (c *LRUMemoryCache) Clear(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.list = list.New()
	c.map_ = make(map[string]*list.Element)
	c.hits.Store(0)
	c.misses.Store(0)
	c.evictions.Store(0)
}

// Cleanup removes expired entries
func (c *LRUMemoryCache) Cleanup(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var keysToDelete []string
	now := time.Now()
	for key, elem := range c.map_ {
		entry := elem.Value.(*LRUCacheEntry)
		if now.After(entry.ExpiresAt) {
			keysToDelete = append(keysToDelete, key)
		}
	}

	for _, key := range keysToDelete {
		if elem, exists := c.map_[key]; exists {
			c.list.Remove(elem)
			delete(c.map_, key)
			c.evictions.Add(1)
		}
	}
}

// GetStats returns cache statistics
func (c *LRUMemoryCache) GetStats() CacheStats {
	return CacheStats{
		Hits:      c.hits.Load(),
		Misses:    c.misses.Load(),
		Evictions: c.evictions.Load(),
	}
}

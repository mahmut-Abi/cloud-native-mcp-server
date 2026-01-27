package cache

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mahmut-Abi/k8s-mcp-server/internal/constants"
)

// lruNode represents a node in the LRU doubly-linked list
type lruNode struct {
	key        string
	entry      CacheEntry
	prev, next *lruNode
}

// MemoryCache implements Cache interface using in-memory storage with true LRU eviction
type MemoryCache struct {
	mu         sync.RWMutex
	store      map[string]*lruNode
	head, tail *lruNode
	maxSize    int
	hits       atomic.Int64
	misses     atomic.Int64
	evictions  atomic.Int64
	stopChan   chan struct{} // Channel to stop the cleanup goroutine
}

// NewMemoryCache creates a new in-memory cache with pre-allocated capacity
func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		store:    make(map[string]*lruNode, 1024),
		maxSize:  constants.DefaultCacheSize,
		stopChan: make(chan struct{}),
	}
	// Start background cleanup task
	go cache.startCleanupTask()
	return cache
}

// Stop stops the background cleanup task
func (c *MemoryCache) Stop() {
	close(c.stopChan)
}

// startCleanupTask runs periodic cleanup of expired entries
func (c *MemoryCache) startCleanupTask() {
	ticker := time.NewTicker(constants.CacheCleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.Cleanup(context.Background())
		case <-c.stopChan:
			return
		}
	}
}

// Get retrieves a value from cache and moves it to the front (most recently used)
func (c *MemoryCache) Get(ctx context.Context, key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, exists := c.store[key]
	if !exists {
		c.misses.Add(1)
		return nil, false
	}

	// Check expiration
	if time.Now().After(node.entry.ExpiresAt) {
		c.removeNode(node)
		delete(c.store, key)
		c.misses.Add(1)
		return nil, false
	}

	// Move to front (mark as recently used)
	c.moveToFront(node)
	c.hits.Add(1)
	return node.entry.Value, true
}

// Set stores a value in cache with TTL and moves it to the front
func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	expiresAt := time.Now().Add(ttl)

	c.mu.Lock()
	defer c.mu.Unlock()

	// If key already exists, update and move to front
	if node, exists := c.store[key]; exists {
		node.entry = CacheEntry{
			Value:     value,
			ExpiresAt: expiresAt,
		}
		c.moveToFront(node)
		return
	}

	// Create new node
	node := &lruNode{
		key: key,
		entry: CacheEntry{
			Value:     value,
			ExpiresAt: expiresAt,
		},
	}

	// Add to front
	c.store[key] = node
	c.addToFront(node)

	// Check if we need to evict
	if len(c.store) > c.maxSize {
		c.evictLRU()
	}
}

// Delete removes a value from cache
func (c *MemoryCache) Delete(ctx context.Context, key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.store[key]; exists {
		c.removeNode(node)
		delete(c.store, key)
	}
}

// Clear removes all entries from cache
func (c *MemoryCache) Clear(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store = make(map[string]*lruNode, 1024)
	c.head = nil
	c.tail = nil
}

// GetStats returns cache statistics
func (c *MemoryCache) GetStats() CacheStats {
	return CacheStats{
		Hits:      c.hits.Load(),
		Misses:    c.misses.Load(),
		Evictions: c.evictions.Load(),
	}
}

// Cleanup removes expired entries
func (c *MemoryCache) Cleanup(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, node := range c.store {
		if now.After(node.entry.ExpiresAt) {
			c.removeNode(node)
			delete(c.store, key)
			c.evictions.Add(1)
		}
	}
}

// evictLRU removes the least recently used entry (tail of the list)
func (c *MemoryCache) evictLRU() {
	if c.tail == nil {
		return
	}

	delete(c.store, c.tail.key)
	c.removeNode(c.tail)
	c.evictions.Add(1)
}

// addToFront adds a node to the front of the list
func (c *MemoryCache) addToFront(node *lruNode) {
	node.prev = nil
	node.next = c.head

	if c.head != nil {
		c.head.prev = node
	}
	c.head = node

	if c.tail == nil {
		c.tail = node
	}
}

// moveToFront moves an existing node to the front
func (c *MemoryCache) moveToFront(node *lruNode) {
	if node == c.head {
		return
	}

	c.removeNode(node)
	c.addToFront(node)
}

// removeNode removes a node from the list
func (c *MemoryCache) removeNode(node *lruNode) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		c.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		c.tail = node.prev
	}
}

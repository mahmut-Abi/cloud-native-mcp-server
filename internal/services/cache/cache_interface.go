package cache

import (
	"context"
	"time"
)

// CacheEntry represents a cached value with metadata
type CacheEntry struct {
	Value     interface{}
	ExpiresAt time.Time
}

// Cache defines a standardized caching interface for all services
type Cache interface {
	// Get retrieves a value from cache
	Get(ctx context.Context, key string) (interface{}, bool)

	// Set stores a value in cache with TTL
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration)

	// Delete removes a value from cache
	Delete(ctx context.Context, key string)

	// Clear removes all entries from cache
	Clear(ctx context.Context)
}

// CacheStats tracks cache performance metrics
type CacheStats struct {
	Hits      int64
	Misses    int64
	Evictions int64
}

// GetHitRate returns cache hit rate percentage
func (s *CacheStats) GetHitRate() float64 {
	total := s.Hits + s.Misses
	if total == 0 {
		return 0
	}
	return float64(s.Hits) / float64(total) * 100
}

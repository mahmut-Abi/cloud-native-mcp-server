package cache

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestSegmentedCacheGet(t *testing.T) {
	cache := NewSegmentedCache()
	ctx := context.Background()

	// Test Get on non-existent key
	_, exists := cache.Get(ctx, "key1")
	if exists {
		t.Error("Expected key to not exist")
	}

	// Test Set and Get
	cache.Set(ctx, "key1", "value1", 1*time.Hour)
	val, exists := cache.Get(ctx, "key1")
	if !exists || val != "value1" {
		t.Errorf("Expected value1, got %v", val)
	}
}

func TestSegmentedCacheConcurrency(t *testing.T) {
	cache := NewSegmentedCache()
	ctx := context.Background()
	var wg sync.WaitGroup

	// Concurrent writes
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			key := "key" + string(rune(index))
			cache.Set(ctx, key, index, 1*time.Hour)
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			key := "key" + string(rune(index))
			_, _ = cache.Get(ctx, key)
		}(i)
	}

	wg.Wait()

	stats := cache.GetStats()
	if stats.Hits+stats.Misses == 0 {
		t.Error("Expected some cache operations")
	}
}

func TestSegmentedCacheExpiration(t *testing.T) {
	cache := NewSegmentedCache()
	ctx := context.Background()

	// Set with short TTL
	cache.Set(ctx, "expiring", "value", 10*time.Millisecond)

	// Should exist immediately
	_, exists := cache.Get(ctx, "expiring")
	if !exists {
		t.Error("Expected key to exist immediately")
	}

	// Wait for expiration
	time.Sleep(20 * time.Millisecond)

	// Should not exist after expiration
	_, exists = cache.Get(ctx, "expiring")
	if exists {
		t.Error("Expected key to be expired")
	}
}

func TestSegmentedCacheDelete(t *testing.T) {
	cache := NewSegmentedCache()
	ctx := context.Background()

	cache.Set(ctx, "key", "value", 1*time.Hour)
	cache.Delete(ctx, "key")

	_, exists := cache.Get(ctx, "key")
	if exists {
		t.Error("Expected key to be deleted")
	}
}

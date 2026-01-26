package cache

import (
	"context"
	"testing"
	"time"
)

func TestMemoryCacheGet(t *testing.T) {
	cache := NewMemoryCache()
	ctx := context.Background()

	// Test get non-existent key
	_, exists := cache.Get(ctx, "key1")
	if exists {
		t.Error("Expected key not to exist")
	}

	// Test set and get
	cache.Set(ctx, "key1", "value1", 1*time.Hour)
	val, exists := cache.Get(ctx, "key1")
	if !exists {
		t.Error("Expected key to exist")
	}
	if val != "value1" {
		t.Errorf("Expected 'value1', got %v", val)
	}
}

func TestMemoryCacheExpiration(t *testing.T) {
	cache := NewMemoryCache()
	ctx := context.Background()

	// Set with short TTL
	cache.Set(ctx, "key1", "value1", 100*time.Millisecond)

	// Should exist immediately
	_, exists := cache.Get(ctx, "key1")
	if !exists {
		t.Error("Expected key to exist")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should be expired
	_, exists = cache.Get(ctx, "key1")
	if exists {
		t.Error("Expected key to be expired")
	}
}

func TestMemoryCacheDelete(t *testing.T) {
	cache := NewMemoryCache()
	ctx := context.Background()

	cache.Set(ctx, "key1", "value1", 1*time.Hour)
	cache.Delete(ctx, "key1")

	_, exists := cache.Get(ctx, "key1")
	if exists {
		t.Error("Expected key to be deleted")
	}
}

func TestMemoryCacheClear(t *testing.T) {
	cache := NewMemoryCache()
	ctx := context.Background()

	cache.Set(ctx, "key1", "value1", 1*time.Hour)
	cache.Set(ctx, "key2", "value2", 1*time.Hour)
	cache.Clear(ctx)

	_, exists := cache.Get(ctx, "key1")
	if exists {
		t.Error("Expected cache to be cleared")
	}
}

func TestMemoryCacheStats(t *testing.T) {
	cache := NewMemoryCache()
	ctx := context.Background()

	cache.Set(ctx, "key1", "value1", 1*time.Hour)

	// Generate hits
	for i := 0; i < 5; i++ {
		cache.Get(ctx, "key1")
	}

	// Generate misses
	for i := 0; i < 3; i++ {
		cache.Get(ctx, "missing-key")
	}

	stats := cache.GetStats()
	if stats.Hits != 5 {
		t.Errorf("Expected 5 hits, got %d", stats.Hits)
	}
	if stats.Misses != 3 {
		t.Errorf("Expected 3 misses, got %d", stats.Misses)
	}

	hitRate := stats.GetHitRate()
	expectedRate := 62.5
	if hitRate != expectedRate {
		t.Errorf("Expected hit rate %.1f, got %.1f", expectedRate, hitRate)
	}
}

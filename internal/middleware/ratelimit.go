package middleware

import (
	"container/heap"
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/mahmut-Abi/k8s-mcp-server/internal/constants"
)

// rateLimitEntry represents a rate limit entry for a client
type rateLimitEntry struct {
	clientID string
	tokens   float64
	lastSeen time.Time
	index    int // Index in the heap
}

// rateLimitHeap implements heap.Interface for rate limit entries
type rateLimitHeap []*rateLimitEntry

func (h rateLimitHeap) Len() int           { return len(h) }
func (h rateLimitHeap) Less(i, j int) bool { return h[i].lastSeen.Before(h[j].lastSeen) }
func (h rateLimitHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *rateLimitHeap) Push(x interface{}) {
	n := len(*h)
	entry := x.(*rateLimitEntry)
	entry.index = n
	*h = append(*h, entry)
}

func (h *rateLimitHeap) Pop() interface{} {
	old := *h
	n := len(old)
	entry := old[n-1]
	old[n-1] = nil
	entry.index = -1
	*h = old[0 : n-1]
	return entry
}

// RateLimiter implements token bucket algorithm for rate limiting with optimized cleanup
type RateLimiter struct {
	mu               sync.RWMutex
	entries          map[string]*rateLimitEntry
	heap             *rateLimitHeap
	rps              float64       // requests per second
	burst            int           // burst size
	cleanupTicker    *time.Ticker  // Periodic cleanup of old entries
	cleanupThreshold int           // Threshold for triggering cleanup
	cleanupBatch     int           // Number of entries to clean in one batch
	staleDuration    time.Duration // Duration after which an entry is considered stale
	ctx              context.Context
	cancel           context.CancelFunc
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rps float64, burst int) *RateLimiter {
	h := &rateLimitHeap{}
	heap.Init(h)

	ctx, cancel := context.WithCancel(context.Background())

	rl := &RateLimiter{
		entries:          make(map[string]*rateLimitEntry),
		heap:             h,
		rps:              rps,
		burst:            burst,
		cleanupThreshold: constants.RateLimitCleanupThreshold,
		cleanupBatch:     constants.RateLimitCleanupBatch,
		staleDuration:    constants.RateLimitStaleDuration,
		ctx:              ctx,
		cancel:           cancel,
	}
	// Start background cleanup
	rl.cleanupTicker = time.NewTicker(constants.RateLimitCleanupInterval)
	go rl.cleanupOldEntries()
	return rl
}

// cleanupOldEntries removes stale entries using the heap for efficient cleanup
func (rl *RateLimiter) cleanupOldEntries() {
	for {
		select {
		case <-rl.cleanupTicker.C:
			rl.mu.Lock()
			if len(rl.entries) > rl.cleanupThreshold {
				now := time.Now()
				deleted := 0

				// Remove stale entries from the heap (oldest first)
				for rl.heap.Len() > 0 && deleted < rl.cleanupBatch {
					entry := (*rl.heap)[0]
					if now.Sub(entry.lastSeen) > rl.staleDuration {
						heap.Pop(rl.heap)
						delete(rl.entries, entry.clientID)
						deleted++
					} else {
						// Since heap is ordered by lastSeen, if the oldest entry is not stale,
						// no other entries will be stale either
						break
					}
				}
			}
			rl.mu.Unlock()
		case <-rl.ctx.Done():
			// Context cancelled, stop cleanup goroutine
			return
		}
	}
}

// Allow checks if a request from the given client is allowed
func (rl *RateLimiter) Allow(clientID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	entry, exists := rl.entries[clientID]

	if !exists {
		// Create new entry
		entry = &rateLimitEntry{
			clientID: clientID,
			tokens:   float64(rl.burst) - 1,
			lastSeen: now,
		}
		rl.entries[clientID] = entry
		heap.Push(rl.heap, entry)
		return true
	}

	// Add tokens based on time elapsed
	elapsed := now.Sub(entry.lastSeen).Seconds()
	entry.tokens += elapsed * rl.rps
	if entry.tokens > float64(rl.burst) {
		entry.tokens = float64(rl.burst)
	}

	// Update lastSeen and fix heap position
	entry.lastSeen = now
	heap.Fix(rl.heap, entry.index)

	if entry.tokens >= 1 {
		entry.tokens--
		return true
	}

	return false
}

// Close stops the cleanup goroutine and releases resources
func (rl *RateLimiter) Close() {
	if rl.cancel != nil {
		rl.cancel()
	}
	if rl.cleanupTicker != nil {
		rl.cleanupTicker.Stop()
	}
}

// RateLimitMiddleware creates a middleware that enforces rate limiting
func RateLimitMiddleware(rps float64, burst int) func(http.Handler) http.Handler {
	limiter := NewRateLimiter(rps, burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientID := getClientIP(r)

			if !limiter.Allow(clientID) {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP extracts client IP from request
func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	return r.RemoteAddr
}

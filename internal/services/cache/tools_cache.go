package cache

import (
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
)

// ToolsCache provides a simple caching mechanism for MCP tools
type ToolsCache struct {
	mu     sync.RWMutex
	tools  []mcp.Tool
	cached bool
}

// NewToolsCache creates a new tools cache
func NewToolsCache() *ToolsCache {
	return &ToolsCache{}
}

// Get returns cached tools if available, otherwise calls the provider function
func (c *ToolsCache) Get(provider func() []mcp.Tool) []mcp.Tool {
	c.mu.RLock()
	if c.cached {
		defer c.mu.RUnlock()
		return c.tools
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if c.cached {
		return c.tools
	}

	c.tools = provider()
	c.cached = true
	return c.tools
}

// Invalidate clears the cache
func (c *ToolsCache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cached = false
	c.tools = nil
}

// IsCached returns whether the cache is populated
func (c *ToolsCache) IsCached() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cached
}

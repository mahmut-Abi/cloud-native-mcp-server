// Package tools provides shared utilities for building and managing MCP tool definitions.
package tools

import (
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
)

// ToolParameter defines a structured parameter with all metadata.
type ToolParameter struct {
	Name        string
	Type        string
	Description string
	Required    bool
	Default     interface{}
	Enum        []string
}

// ToolCache implements lazy loading and caching of tool definitions.
// This reduces initialization overhead and provides thread-safe access.
type ToolCache struct {
	mu         sync.RWMutex
	tools      map[string]mcp.Tool
	accessList []string // Track access order for LRU eviction
	maxSize    int      // Maximum number of tools to cache
}

const defaultMaxCacheSize = 200 // Default maximum cache size

var globalToolCache = &ToolCache{
	tools:      make(map[string]mcp.Tool),
	accessList: make([]string, 0, defaultMaxCacheSize),
	maxSize:    defaultMaxCacheSize,
}

// GetOrCreate retrieves a cached tool or creates it using the builder function.
// This ensures each tool is built only once, reducing memory and improving startup performance.
func (tc *ToolCache) GetOrCreate(name string, builder func() mcp.Tool) mcp.Tool {
	tc.mu.RLock()
	tool, exists := tc.tools[name]
	tc.mu.RUnlock()

	if exists {
		// Update access order for LRU
		tc.updateAccessOrder(name)
		return tool
	}

	// Tool not cached, build it
	tool = builder()

	// Store in cache (double-check to avoid duplicate creation)
	tc.mu.Lock()
	defer tc.mu.Unlock()

	// Check again in case another goroutine created it while we were building
	if existingTool, ok := tc.tools[name]; ok {
		tc.updateAccessOrder(name)
		return existingTool
	}

	// Evict oldest entry if cache is full
	if len(tc.tools) >= tc.maxSize {
		tc.evictOldest()
	}

	tc.tools[name] = tool
	tc.accessList = append(tc.accessList, name)
	return tool
}

// updateAccessOrder updates the access order for LRU eviction
// Must be called with write lock held
func (tc *ToolCache) updateAccessOrder(name string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	// Remove from current position
	for i, n := range tc.accessList {
		if n == name {
			tc.accessList = append(tc.accessList[:i], tc.accessList[i+1:]...)
			break
		}
	}
	// Add to end (most recently used)
	tc.accessList = append(tc.accessList, name)
}

// evictOldest removes the least recently used tool from cache
// Must be called with write lock held
func (tc *ToolCache) evictOldest() {
	if len(tc.accessList) == 0 {
		return
	}
	oldest := tc.accessList[0]
	delete(tc.tools, oldest)
	tc.accessList = tc.accessList[1:]
}

// GlobalToolCache returns the global tool cache instance.
func GlobalToolCache() *ToolCache {
	return globalToolCache
}

// NewStringProperty creates a standard string property schema.
func NewStringProperty(description string) map[string]interface{} {
	return map[string]interface{}{
		"type":        "string",
		"description": description,
	}
}

// NewNumberProperty creates a standard number property schema.
func NewNumberProperty(description string, enum []float64) map[string]interface{} {
	prop := map[string]interface{}{
		"type":        "number",
		"description": description,
	}
	if len(enum) > 0 {
		prop["enum"] = enum
	}
	return prop
}

// NewBooleanProperty creates a standard boolean property schema.
func NewBooleanProperty(description string) map[string]interface{} {
	return map[string]interface{}{
		"type":        "boolean",
		"description": description,
	}
}

// NewEnumProperty creates a string property with enum values.
func NewEnumProperty(description string, enum []string) map[string]interface{} {
	return map[string]interface{}{
		"type":        "string",
		"description": description,
		"enum":        enum,
	}
}

// CreateObjectSchema builds a complete tool input schema from parameters.
func CreateObjectSchema(params map[string]map[string]interface{}, required []string) mcp.ToolInputSchema {
	properties := make(map[string]any)
	for k, v := range params {
		properties[k] = v
	}
	return mcp.ToolInputSchema{
		Type:       "object",
		Properties: properties,
		Required:   required,
	}
}

// GetCommonDescription retrieves a standardized description for a parameter.
// Returns empty string if key not found.
func GetCommonDescription(key string) string {
	if desc, exists := CommonDescriptions[key]; exists {
		return desc
	}
	return ""
}

// GetParameterDefault retrieves a standardized default value for a parameter.
func GetParameterDefault(key string) (interface{}, bool) {
	val, exists := ParameterDefaults[key]
	return val, exists
}

// GetEnumValues retrieves standard enum values for a parameter.
func GetEnumValues(key string) ([]string, bool) {
	vals, exists := ParameterEnums[key]
	return vals, exists
}

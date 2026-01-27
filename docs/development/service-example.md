# Minimal Service Example

This document provides a complete, minimal example of creating a new service from scratch.

## Example: Simple Counter Service

We'll create a simple "Counter" service that maintains a counter value through MCP tools.

### Step 1: Service Structure

```
internal/services/counter/
├── client/
│   └── client.go          # In-memory counter client
├── handlers/
│   └── handlers.go        # Tool handlers
├── tools/
│   └── tools.go          # Tool definitions
└── service.go            # Service implementation
```

### Step 2: Service Implementation

`internal/services/counter/service.go`:

```go
package counter

import (
    "context"

    "github.com/mark3labs/mcp-go/mcp"
    server "github.com/mark3labs/mcp-go/server"

    "github.com/mahmut-Abi/k8s-mcp-server/internal/config"
    "github.com/mahmut-Abi/k8s-mcp-server/internal/services/cache"
    "github.com/mahmut-Abi/k8s-mcp-server/internal/services/counter/client"
    "github.com/mahmut-Abi/k8s-mcp-server/internal/services/counter/handlers"
    "github.com/mahmut-Abi/k8s-mcp-server/internal/services/counter/tools"
)

// Service implements the Counter MCP service
type Service struct {
    client     *client.Client
    enabled    bool
    toolsCache *cache.ToolsCache
}

// NewService creates a new Counter service instance
func NewService() *Service {
    return &Service{
        enabled:    true, // Counter is always enabled
        toolsCache: cache.NewToolsCache(),
    }
}

// Name returns the service identifier
func (s *Service) Name() string {
    return "counter"
}

// Initialize configures the service
func (s *Service) Initialize(cfg interface{}) error {
    // Counter doesn't need configuration, just create client
    s.client = client.NewClient()
    s.enabled = true
    return nil
}

// GetTools returns all available Counter tools
func (s *Service) GetTools() []mcp.Tool {
    if !s.enabled || s.client == nil {
        return nil
    }

    return s.toolsCache.Get(func() []mcp.Tool {
        return []mcp.Tool{
            tools.GetCounterTool(),
            tools.IncrementCounterTool(),
            tools.DecrementCounterTool(),
            tools.ResetCounterTool(),
        }
    })
}

// GetHandlers returns all tool handlers
func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
    if !s.enabled || s.client == nil {
        return nil
    }

    return map[string]server.ToolHandlerFunc{
        "counter_get":       handlers.HandleGetCounter(s.client),
        "counter_increment": handlers.HandleIncrementCounter(s.client),
        "counter_decrement": handlers.HandleDecrementCounter(s.client),
        "counter_reset":     handlers.HandleResetCounter(s.client),
    }
}

// IsEnabled returns whether the service is enabled
func (s *Service) IsEnabled() bool {
    return s.enabled && s.client != nil
}

// GetClient returns the underlying client
func (s *Service) GetClient() *client.Client {
    return s.client
}
```

### Step 3: Client Implementation

`internal/services/counter/client/client.go`:

```go
package client

import (
    "sync"
)

// Client provides in-memory counter operations
type Client struct {
    mu     sync.RWMutex
    value  int64
}

// NewClient creates a new counter client
func NewClient() *Client {
    return &Client{
        value: 0,
    }
}

// Get returns the current counter value
func (c *Client) Get() int64 {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.value
}

// Increment increments the counter by 1
func (c *Client) Increment() int64 {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
    return c.value
}

// Decrement decrements the counter by 1
func (c *Client) Decrement() int64 {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value--
    return c.value
}

// Reset resets the counter to 0
func (c *Client) Reset() int64 {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value = 0
    return c.value
}
```

### Step 4: Tool Definitions

`internal/services/counter/tools/tools.go`:

```go
package tools

import (
    "github.com/mark3labs/mcp-go/mcp"
)

// GetCounterTool returns a tool for getting the counter value
func GetCounterTool() mcp.Tool {
    return mcp.NewTool("counter_get",
        mcp.WithDescription("Get the current counter value"),
    )
}

// IncrementCounterTool returns a tool for incrementing the counter
func IncrementCounterTool() mcp.Tool {
    return mcp.NewTool("counter_increment",
        mcp.WithDescription("Increment the counter by 1"),
    )
}

// DecrementCounterTool returns a tool for decrementing the counter
func DecrementCounterTool() mcp.Tool {
    return mcp.NewTool("counter_decrement",
        mcp.WithDescription("Decrement the counter by 1"),
    )
}

// ResetCounterTool returns a tool for resetting the counter
func ResetCounterTool() mcp.Tool {
    return mcp.NewTool("counter_reset",
        mcp.WithDescription("Reset the counter to 0"),
    )
}
```

### Step 5: Tool Handlers

`internal/services/counter/handlers/handlers.go`:

```go
package handlers

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/mark3labs/mcp-go/mcp"

    "github.com/mahmut-Abi/k8s-mcp-server/internal/services/counter/client"
)

// HandleGetCounter handles getting the counter value
func HandleGetCounter(svcClient *client.Client) mcp.ToolHandlerFunc {
    return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        value := svcClient.Get()
        result := map[string]interface{}{
            "value": value,
        }
        resultJSON, _ := json.MarshalIndent(result, "", "  ")
        return mcp.NewToolResultText(string(resultJSON)), nil
    }
}

// HandleIncrementCounter handles incrementing the counter
func HandleIncrementCounter(svcClient *client.Client) mcp.ToolHandlerFunc {
    return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        value := svcClient.Increment()
        result := map[string]interface{}{
            "value": value,
        }
        resultJSON, _ := json.MarshalIndent(result, "", "  ")
        return mcp.NewToolResultText(string(resultJSON)), nil
    }
}

// HandleDecrementCounter handles decrementing the counter
func HandleDecrementCounter(svcClient *client.Client) mcp.ToolHandlerFunc {
    return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        value := svcClient.Decrement()
        result := map[string]interface{}{
            "value": value,
        }
        resultJSON, _ := json.MarshalIndent(result, "", "  ")
        return mcp.NewToolResultText(string(resultJSON)), nil
    }
}

// HandleResetCounter handles resetting the counter
func HandleResetCounter(svcClient *client.Client) mcp.ToolHandlerFunc {
    return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        value := svcClient.Reset()
        result := map[string]interface{}{
            "value": value,
        }
        resultJSON, _ := json.MarshalIndent(result, "", "  ")
        return mcp.NewToolResultText(string(resultJSON)), nil
    }
}
```

### Step 6: Register the Service

Update `internal/services/manager/manager.go`:

```go
// Add field to Manager struct
type Manager struct {
    // ... existing fields ...
    counterService *counter.Service
    // ... other fields ...
}

// In NewManager(), initialize the service
func NewManager() *Manager {
    return &Manager{
        // ... existing initialization ...
        counterService: counter.NewService(),
    }
}

// In Initialize(), register the service
func (m *Manager) Initialize(appConfig *config.AppConfig) error {
    // ... existing code ...
    
    // Create and register counter service
    m.counterService = counter.NewService()
    if err := m.counterService.Initialize(appConfig); err != nil {
        logger.WithError(err).Warn("Counter service initialization failed")
    }
    m.registry.Register(m.counterService)
    
    // ... rest of initialization ...
}

// In Shutdown(), clean up the service
func (m *Manager) Shutdown() error {
    // ... existing shutdown code ...
    
    // Counter doesn't need cleanup, but we can log it
    if m.counterService != nil {
        logger.Debug("Counter service shutdown")
    }
    
    // ... rest of shutdown ...
}

// Optional: Add getter for counter service
func (m *Manager) GetCounterService() *counter.Service {
    return m.counterService
}
```

### Step 7: Test the Service

Create `internal/services/counter/service_test.go`:

```go
package counter

import (
    "testing"

    "github.com/mark3labs/mcp-go/mcp"
    server "github.com/mark3labs/mcp-go/server"
)

func TestNewService(t *testing.T) {
    service := NewService()
    
    if service == nil {
        t.Fatal("NewService() returned nil")
    }
    
    if service.Name() != "counter" {
        t.Errorf("Expected name 'counter', got '%s'", service.Name())
    }
}

func TestServiceInitialization(t *testing.T) {
    service := NewService()
    
    err := service.Initialize(nil)
    if err != nil {
        t.Fatalf("Failed to initialize service: %v", err)
    }
    
    if !service.IsEnabled() {
        t.Error("Service should be enabled after initialization")
    }
}

func TestGetTools(t *testing.T) {
    service := NewService()
    service.Initialize(nil)
    
    tools := service.GetTools()
    
    if len(tools) != 4 {
        t.Errorf("Expected 4 tools, got %d", len(tools))
    }
    
    toolNames := make(map[string]bool)
    for _, tool := range tools {
        toolNames[tool.Name] = true
    }
    
    expectedTools := []string{
        "counter_get",
        "counter_increment",
        "counter_decrement",
        "counter_reset",
    }
    
    for _, name := range expectedTools {
        if !toolNames[name] {
            t.Errorf("Missing tool: %s", name)
        }
    }
}

func TestGetHandlers(t *testing.T) {
    service := NewService()
    service.Initialize(nil)
    
    handlers := service.GetHandlers()
    
    if len(handlers) != 4 {
        t.Errorf("Expected 4 handlers, got %d", len(handlers))
    }
    
    handlerNames := make(map[string]bool)
    for name := range handlers {
        handlerNames[name] = true
    }
    
    expectedHandlers := []string{
        "counter_get",
        "counter_increment",
        "counter_decrement",
        "counter_reset",
    }
    
    for _, name := range expectedHandlers {
        if !handlerNames[name] {
            t.Errorf("Missing handler: %s", name)
        }
    }
}

func TestClientOperations(t *testing.T) {
    client := client.NewClient()
    
    // Test initial value
    if value := client.Get(); value != 0 {
        t.Errorf("Expected initial value 0, got %d", value)
    }
    
    // Test increment
    if value := client.Increment(); value != 1 {
        t.Errorf("Expected value 1 after increment, got %d", value)
    }
    
    // Test increment again
    if value := client.Increment(); value != 2 {
        t.Errorf("Expected value 2 after second increment, got %d", value)
    }
    
    // Test decrement
    if value := client.Decrement(); value != 1 {
        t.Errorf("Expected value 1 after decrement, got %d", value)
    }
    
    // Test reset
    if value := client.Reset(); value != 0 {
        t.Errorf("Expected value 0 after reset, got %d", value)
    }
}

func TestHandlerExecution(t *testing.T) {
    service := NewService()
    service.Initialize(nil)
    
    handlers := service.GetHandlers()
    
    // Test get handler
    result, err := handlers["counter_get"].(server.ToolHandlerFunc)(nil, mcp.CallToolRequest{})
    if err != nil {
        t.Fatalf("Handler returned error: %v", err)
    }
    
    if result.IsError {
        t.Error("Handler returned error result")
    }
    
    // Test increment handler
    result, err = handlers["counter_increment"].(server.ToolHandlerFunc)(nil, mcp.CallToolRequest{})
    if err != nil {
        t.Fatalf("Handler returned error: %v", err)
    }
    
    if result.IsError {
        t.Error("Handler returned error result")
    }
}
```

### Step 8: Usage Example

Start the server and use the counter tools:

```bash
# Start the server
go run cmd/server/main.go --config config.yaml
```

Then interact with the MCP tools:

```
# Get counter value
Tool: counter_get
Result: {"value": 0}

# Increment counter
Tool: counter_increment
Result: {"value": 1}

# Increment again
Tool: counter_increment
Result: {"value": 2}

# Decrement
Tool: counter_decrement
Result: {"value": 1}

# Reset
Tool: counter_reset
Result: {"value": 0}
```

## Key Points

1. **Minimal Implementation**: This example shows the minimum code needed for a service
2. **In-Memory Client**: No external dependencies, easy to test
3. **Thread-Safe**: Uses mutex for concurrent access
4. **Complete Example**: Includes service, client, tools, handlers, and tests
5. **Easy to Extend**: Add more operations by following the same pattern

## Extending the Example

### Add Configuration

If you want to make the counter configurable:

```go
// In service.go
type Service struct {
    client     *client.Client
    enabled    bool
    toolsCache *cache.ToolsCache
    initialValue int64
}

func NewService(initialValue int64) *Service {
    return &Service{
        enabled:       true,
        toolsCache:    cache.NewToolsCache(),
        initialValue:  initialValue,
    }
}

func (s *Service) Initialize(cfg interface{}) error {
    s.client = client.NewClient(s.initialValue)
    s.enabled = true
    return nil
}
```

### Add Logging

```go
// In handlers.go
import "github.com/sirupsen/logrus"

func HandleIncrementCounter(svcClient *client.Client) mcp.ToolHandlerFunc {
    return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        value := svcClient.Increment()
        logrus.WithField("value", value).Info("Counter incremented")
        
        result := map[string]interface{}{"value": value}
        resultJSON, _ := json.MarshalIndent(result, "", "  ")
        return mcp.NewToolResultText(string(resultJSON)), nil
    }
}
```

### Add Metrics

The hook system will automatically record metrics for all tool calls:
- Tool name: `counter_increment`
- Service name: `counter`
- Status: success/error
- Duration: execution time

No additional code needed!

## Summary

This minimal example demonstrates:
- Service structure and initialization
- Client implementation with thread-safety
- Tool definitions with parameters
- Handler implementation
- Service registration
- Comprehensive testing
- Usage examples

Use this as a template for creating more complex services that interact with external APIs, databases, or other systems.
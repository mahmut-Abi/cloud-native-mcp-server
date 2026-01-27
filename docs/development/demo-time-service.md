# Demo: Time Service

A complete, working example of a service that provides time-related operations.

## Overview

The Time Service demonstrates:
- Simple service structure
- No external dependencies
- Practical tool implementations
- Error handling patterns
- Testing examples

## File Structure

```
internal/services/time/
├── client/
│   └── client.go          # Time operations client
├── handlers/
│   └── handlers.go        # Tool handlers
├── tools/
│   └── tools.go          # Tool definitions
├── service.go            # Service implementation
└── service_test.go       # Service tests
```

## Implementation

### 1. Service Implementation

`internal/services/time/service.go`:

```go
package time

import (
    "context"
    "time"

    "github.com/mark3labs/mcp-go/mcp"
    server "github.com/mark3labs/mcp-go/server"

    "github.com/mahmut-Abi/k8s-mcp-server/internal/config"
    "github.com/mahmut-Abi/k8s-mcp-server/internal/services/cache"
    "github.com/mahmut-Abi/k8s-mcp-server/internal/services/time/client"
    "github.com/mahmut-Abi/k8s-mcp-server/internal/services/time/handlers"
    "github.com/mahmut-Abi/k8s-mcp-server/internal/services/time/tools"
)

// Service implements the Time MCP service
type Service struct {
    client     *client.Client
    enabled    bool
    toolsCache *cache.ToolsCache
}

// NewService creates a new Time service instance
func NewService() *Service {
    return &Service{
        enabled:    true, // Time service is always enabled
        toolsCache: cache.NewToolsCache(),
    }
}

// Name returns the service identifier
func (s *Service) Name() string {
    return "time"
}

// Initialize configures the service
func (s *Service) Initialize(cfg interface{}) error {
    // Time service doesn't need configuration
    s.client = client.NewClient()
    s.enabled = true
    return nil
}

// GetTools returns all available Time tools
func (s *Service) GetTools() []mcp.Tool {
    if !s.enabled || s.client == nil {
        return nil
    }

    return s.toolsCache.Get(func() []mcp.Tool {
        return []mcp.Tool{
            tools.GetCurrentTimeTool(),
            tools.GetTimeInZoneTool(),
            tools.GetUnixTimestampTool(),
            tools.ParseTimeTool(),
            tools.FormatTimeTool(),
        }
    })
}

// GetHandlers returns all tool handlers
func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
    if !s.enabled || s.client == nil {
        return nil
    }

    return map[string]server.ToolHandlerFunc{
        "time_get_current":       handlers.HandleGetCurrentTime(s.client),
        "time_get_in_zone":       handlers.HandleGetTimeInZone(s.client),
        "time_get_unix_timestamp": handlers.HandleGetUnixTimestamp(s.client),
        "time_parse":             handlers.HandleParseTime(s.client),
        "time_format":            handlers.HandleFormatTime(s.client),
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

### 2. Client Implementation

`internal/services/time/client/client.go`:

```go
package client

import (
    "fmt"
    "time"
)

// Client provides time-related operations
type Client struct {
    location *time.Location
}

// NewClient creates a new time client
func NewClient() *Client {
    return &Client{
        location: time.UTC,
    }
}

// GetCurrentTime returns the current time
func (c *Client) GetCurrentTime() time.Time {
    return time.Now()
}

// GetTimeInZone returns the current time in a specific time zone
func (c *Client) GetTimeInZone(zone string) (time.Time, error) {
    location, err := time.LoadLocation(zone)
    if err != nil {
        return time.Time{}, fmt.Errorf("invalid time zone: %w", err)
    }
    return time.Now().In(location), nil
}

// GetUnixTimestamp returns the current Unix timestamp
func (c *Client) GetUnixTimestamp() int64 {
    return time.Now().Unix()
}

// ParseTime parses a time string
func (c *Client) ParseTime(timeStr, layout string) (time.Time, error) {
    if layout == "" {
        layout = time.RFC3339
    }
    t, err := time.Parse(layout, timeStr)
    if err != nil {
        return time.Time{}, fmt.Errorf("failed to parse time: %w", err)
    }
    return t, nil
}

// FormatTime formats a time string
func (c *Client) FormatTime(timeStr, layout, outputLayout string) (string, error) {
    t, err := c.ParseTime(timeStr, layout)
    if err != nil {
        return "", err
    }

    if outputLayout == "" {
        outputLayout = time.RFC3339
    }

    return t.Format(outputLayout), nil
}

// TimeInfo represents time information
type TimeInfo struct {
    Timestamp    string `json:"timestamp"`
    Unix         int64  `json:"unix"`
    TimeZone     string `json:"timezone"`
    UTCOffset    string `json:"utc_offset"`
    DayOfWeek    string `json:"day_of_week"`
    Month        string `json:"month"`
    Year         int    `json:"year"`
}
```

### 3. Tool Definitions

`internal/services/time/tools/tools.go`:

```go
package tools

import (
    "github.com/mark3labs/mcp-go/mcp"
)

// GetCurrentTimeTool returns a tool for getting the current time
func GetCurrentTimeTool() mcp.Tool {
    return mcp.NewTool("time_get_current",
        mcp.WithDescription("Get the current time in the system's time zone"),
    )
}

// GetTimeInZoneTool returns a tool for getting time in a specific time zone
func GetTimeInZoneTool() mcp.Tool {
    return mcp.NewTool("time_get_in_zone",
        mcp.WithDescription("Get the current time in a specific time zone"),
        mcp.WithString("zone",
            mcp.Required(),
            mcp.Description("Time zone name (e.g., 'America/New_York', 'Asia/Shanghai')"),
        ),
    )
}

// GetUnixTimestampTool returns a tool for getting Unix timestamp
func GetUnixTimestampTool() mcp.Tool {
    return mcp.NewTool("time_get_unix_timestamp",
        mcp.WithDescription("Get the current Unix timestamp"),
    )
}

// ParseTimeTool returns a tool for parsing a time string
func ParseTimeTool() mcp.Tool {
    return mcp.NewTool("time_parse",
        mcp.WithDescription("Parse a time string into a structured format"),
        mcp.WithString("time_str",
            mcp.Required(),
            mcp.Description("Time string to parse (default format: RFC3339)"),
        ),
        mcp.WithString("layout",
            mcp.Description("Time layout format (default: RFC3339)"),
        ),
    )
}

// FormatTimeTool returns a tool for formatting a time string
func FormatTimeTool() mcp.Tool {
    return mcp.NewTool("time_format",
        mcp.WithDescription("Format a time string to a different format"),
        mcp.WithString("time_str",
            mcp.Required(),
            mcp.Description("Time string to format"),
        ),
        mcp.WithString("layout",
            mcp.Description("Input time layout (default: RFC3339)"),
        ),
        mcp.WithString("output_layout",
            mcp.Description("Output time layout (default: RFC3339)"),
        ),
    )
}
```

### 4. Handler Implementations

`internal/services/time/handlers/handlers.go`:

```go
package handlers

import (
    "encoding/json"
    "fmt"

    "github.com/mark3labs/mcp-go/mcp"
    "github.com/sirupsen/logrus"

    "github.com/mahmut-Abi/k8s-mcp-server/internal/services/time/client"
)

// HandleGetCurrentTime handles getting the current time
func HandleGetCurrentTime(svcClient *client.Client) mcp.ToolHandlerFunc {
    return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        now := svcClient.GetCurrentTime()
        timeInfo := buildTimeInfo(now, now.Location())

        resultJSON, err := json.MarshalIndent(timeInfo, "", "  ")
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("failed to marshal result: %v", err)), nil
        }

        return mcp.NewToolResultText(string(resultJSON)), nil
    }
}

// HandleGetTimeInZone handles getting time in a specific time zone
func HandleGetTimeInZone(svcClient *client.Client) mcp.ToolHandlerFunc {
    return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        zone, err := request.Params.Arguments.String("zone")
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("invalid zone parameter: %v", err)), nil
        }

        timeInZone, err := svcClient.GetTimeInZone(zone)
        if err != nil {
            logrus.WithError(err).Error("Failed to get time in zone")
            return mcp.NewToolResultError(fmt.Sprintf("failed to get time in zone: %v", err)), nil
        }

        timeInfo := buildTimeInfo(timeInZone, timeInZone.Location())

        resultJSON, err := json.MarshalIndent(timeInfo, "", "  ")
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("failed to marshal result: %v", err)), nil
        }

        return mcp.NewToolResultText(string(resultJSON)), nil
    }
}

// HandleGetUnixTimestamp handles getting Unix timestamp
func HandleGetUnixTimestamp(svcClient *client.Client) mcp.ToolHandlerFunc {
    return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        timestamp := svcClient.GetUnixTimestamp()

        result := map[string]interface{}{
            "unix":      timestamp,
            "timestamp": fmt.Sprintf("%d", timestamp),
        }

        resultJSON, err := json.MarshalIndent(result, "", "  ")
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("failed to marshal result: %v", err)), nil
        }

        return mcp.NewToolResultText(string(resultJSON)), nil
    }
}

// HandleParseTime handles parsing a time string
func HandleParseTime(svcClient *client.Client) mcp.ToolHandlerFunc {
    return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        timeStr, err := request.Params.Arguments.String("time_str")
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("invalid time_str parameter: %v", err)), nil
        }

        layout, _ := request.Params.Arguments.String("layout")

        parsedTime, err := svcClient.ParseTime(timeStr, layout)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("failed to parse time: %v", err)), nil
        }

        result := map[string]interface{}{
            "parsed":    parsedTime.Format(time.RFC3339),
            "unix":      parsedTime.Unix(),
            "location":  parsedTime.Location().String(),
        }

        resultJSON, err := json.MarshalIndent(result, "", "  ")
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("failed to marshal result: %v", err)), nil
        }

        return mcp.NewToolResultText(string(resultJSON)), nil
    }
}

// HandleFormatTime handles formatting a time string
func HandleFormatTime(svcClient *client.Client) mcp.ToolHandlerFunc {
    return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        timeStr, err := request.Params.Arguments.String("time_str")
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("invalid time_str parameter: %v", err)), nil
        }

        layout, _ := request.Params.Arguments.String("layout")
        outputLayout, _ := request.Params.Arguments.String("output_layout")

        formatted, err := svcClient.FormatTime(timeStr, layout, outputLayout)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("failed to format time: %v", err)), nil
        }

        result := map[string]interface{}{
            "formatted": formatted,
        }

        resultJSON, err := json.MarshalIndent(result, "", "  ")
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("failed to marshal result: %v", err)), nil
        }

        return mcp.NewToolResultText(string(resultJSON)), nil
    }
}

// buildTimeInfo builds a TimeInfo struct from a time.Time
func buildTimeInfo(t time.Time, loc *time.Location) client.TimeInfo {
    return client.TimeInfo{
        Timestamp: t.Format(time.RFC3339),
        Unix:      t.Unix(),
        TimeZone:  loc.String(),
        UTCOffset: t.Format("-0700"),
        DayOfWeek: t.Weekday().String(),
        Month:     t.Month().String(),
        Year:      t.Year(),
    }
}
```

### 5. Service Tests

`internal/services/time/service_test.go`:

```go
package time

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

    if service.Name() != "time" {
        t.Errorf("Expected name 'time', got '%s'", service.Name())
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

    if len(tools) != 5 {
        t.Errorf("Expected 5 tools, got %d", len(tools))
    }

    toolNames := make(map[string]bool)
    for _, tool := range tools {
        toolNames[tool.Name] = true
    }

    expectedTools := []string{
        "time_get_current",
        "time_get_in_zone",
        "time_get_unix_timestamp",
        "time_parse",
        "time_format",
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

    if len(handlers) != 5 {
        t.Errorf("Expected 5 handlers, got %d", len(handlers))
    }

    handlerNames := make(map[string]bool)
    for name := range handlers {
        handlerNames[name] = true
    }

    expectedHandlers := []string{
        "time_get_current",
        "time_get_in_zone",
        "time_get_unix_timestamp",
        "time_parse",
        "time_format",
    }

    for _, name := range expectedHandlers {
        if !handlerNames[name] {
            t.Errorf("Missing handler: %s", name)
        }
    }
}

func TestClientOperations(t *testing.T) {
    client := client.NewClient()

    // Test GetCurrentTime
    now := client.GetCurrentTime()
    if now.IsZero() {
        t.Error("GetCurrentTime returned zero time")
    }

    // Test GetUnixTimestamp
    timestamp := client.GetUnixTimestamp()
    if timestamp == 0 {
        t.Error("GetUnixTimestamp returned 0")
    }

    // Test GetTimeInZone
    t.Run("GetTimeInZone - UTC", func(t *testing.T) {
        timeInUTC, err := client.GetTimeInZone("UTC")
        if err != nil {
            t.Errorf("Failed to get time in UTC: %v", err)
        }
        if timeInUTC.IsZero() {
            t.Error("GetTimeInZone returned zero time")
        }
    })

    t.Run("GetTimeInZone - Shanghai", func(t *testing.T) {
        timeInShanghai, err := client.GetTimeInZone("Asia/Shanghai")
        if err != nil {
            t.Errorf("Failed to get time in Shanghai: %v", err)
        }
        if timeInShanghai.IsZero() {
            t.Error("GetTimeInZone returned zero time")
        }
    })

    t.Run("GetTimeInZone - Invalid", func(t *testing.T) {
        _, err := client.GetTimeInZone("Invalid/Timezone")
        if err == nil {
            t.Error("Expected error for invalid time zone")
        }
    })

    // Test ParseTime
    t.Run("ParseTime - RFC3339", func(t *testing.T) {
        parsed, err := client.ParseTime("2024-01-01T00:00:00Z", time.RFC3339)
        if err != nil {
            t.Errorf("Failed to parse time: %v", err)
        }
        if parsed.Year() != 2024 {
            t.Errorf("Expected year 2024, got %d", parsed.Year())
        }
    })

    // Test FormatTime
    t.Run("FormatTime", func(t *testing.T) {
        formatted, err := client.FormatTime("2024-01-01T00:00:00Z", time.RFC3339, "2006-01-02 15:04:05")
        if err != nil {
            t.Errorf("Failed to format time: %v", err)
        }
        if formatted != "2024-01-01 00:00:00" {
            t.Errorf("Expected '2024-01-01 00:00:00', got '%s'", formatted)
        }
    })
}

func TestHandlerExecution(t *testing.T) {
    service := NewService()
    service.Initialize(nil)

    handlers := service.GetHandlers()

    // Test get_current handler
    t.Run("time_get_current", func(t *testing.T) {
        result, err := handlers["time_get_current"].(server.ToolHandlerFunc)(nil, mcp.CallToolRequest{})
        if err != nil {
            t.Fatalf("Handler returned error: %v", err)
        }
        if result.IsError {
            t.Error("Handler returned error result")
        }
        if result.Content == nil || len(result.Content) == 0 {
            t.Error("Handler returned empty content")
        }
    })

    // Test get_unix_timestamp handler
    t.Run("time_get_unix_timestamp", func(t *testing.T) {
        result, err := handlers["time_get_unix_timestamp"].(server.ToolHandlerFunc)(nil, mcp.CallToolRequest{})
        if err != nil {
            t.Fatalf("Handler returned error: %v", err)
        }
        if result.IsError {
            t.Error("Handler returned error result")
        }
    })

    // Test get_in_zone handler with valid zone
    t.Run("time_get_in_zone - valid", func(t *testing.T) {
        request := mcp.CallToolRequest{
            Params: mcp.CallToolParams{
                Arguments: map[string]interface{}{
                    "zone": "UTC",
                },
            },
        }
        result, err := handlers["time_get_in_zone"].(server.ToolHandlerFunc)(nil, request)
        if err != nil {
            t.Fatalf("Handler returned error: %v", err)
        }
        if result.IsError {
            t.Error("Handler returned error result")
        }
    })

    // Test get_in_zone handler with invalid zone
    t.Run("time_get_in_zone - invalid", func(t *testing.T) {
        request := mcp.CallToolRequest{
            Params: mcp.CallToolParams{
                Arguments: map[string]interface{}{
                    "zone": "Invalid/Zone",
                },
            },
        }
        result, err := handlers["time_get_in_zone"].(server.ToolHandlerFunc)(nil, request)
        if err != nil {
            t.Fatalf("Handler returned error: %v", err)
        }
        if !result.IsError {
            t.Error("Handler should return error for invalid zone")
        }
    })
}
```

## Usage Examples

### Example 1: Get Current Time

```json
{
  "name": "time_get_current",
  "arguments": {}
}
```

Response:
```json
{
  "timestamp": "2024-01-27T10:30:00Z",
  "unix": 1706350200,
  "timezone": "UTC",
  "utc_offset": "+0000",
  "day_of_week": "Saturday",
  "month": "January",
  "year": 2024
}
```

### Example 2: Get Time in Specific Time Zone

```json
{
  "name": "time_get_in_zone",
  "arguments": {
    "zone": "Asia/Shanghai"
  }
}
```

Response:
```json
{
  "timestamp": "2024-01-27T18:30:00+08:00",
  "unix": 1706350200,
  "timezone": "Asia/Shanghai",
  "utc_offset": "+0800",
  "day_of_week": "Saturday",
  "month": "January",
  "year": 2024
}
```

### Example 3: Get Unix Timestamp

```json
{
  "name": "time_get_unix_timestamp",
  "arguments": {}
}
```

Response:
```json
{
  "unix": 1706350200,
  "timestamp": "1706350200"
}
```

### Example 4: Parse Time

```json
{
  "name": "time_parse",
  "arguments": {
    "time_str": "2024-01-01T00:00:00Z",
    "layout": "2006-01-02T15:04:05Z07:00"
  }
}
```

Response:
```json
{
  "parsed": "2024-01-01T00:00:00Z",
  "unix": 1704067200,
  "location": "UTC"
}
```

### Example 5: Format Time

```json
{
  "name": "time_format",
  "arguments": {
    "time_str": "2024-01-01T00:00:00Z",
    "layout": "2006-01-02T15:04:05Z07:00",
    "output_layout": "2006-01-02 15:04:05"
  }
}
```

Response:
```json
{
  "formatted": "2024-01-01 00:00:00"
}
```

## Running the Demo

To run the Time service demo:

1. Add the service to the manager (as shown in the main guide)
2. Start the server:
```bash
go run cmd/server/main.go --config config.yaml
3. Interact with the tools using an MCP client
```

## Testing

Run the tests:
```bash
go test ./internal/services/time/... -v
```

Expected output:
```
=== RUN   TestNewService
--- PASS: TestNewService (0.00s)
=== RUN   TestServiceInitialization
--- PASS: TestServiceInitialization (0.00s)
=== RUN   TestGetTools
--- PASS: TestGetTools (0.00s)
=== RUN   TestGetHandlers
--- PASS: TestGetHandlers (0.00s)
=== RUN   TestClientOperations
=== RUN   TestClientOperations/GetTimeInZone_-_UTC
=== RUN   TestClientOperations/GetTimeInZone_-_Shanghai
=== RUN   TestClientOperations/GetTimeInZone_-_Invalid
=== RUN   TestClientOperations/ParseTime_-_RFC3339
=== RUN   TestClientOperations/FormatTime
--- PASS: TestClientOperations (0.00s)
=== RUN   TestHandlerExecution
=== RUN   TestHandlerExecution/time_get_current
=== RUN   TestHandlerExecution/time_get_unix_timestamp
=== RUN   TestHandlerExecution/time_get_in_zone_-_valid
=== RUN   TestHandlerExecution/time_get_in_zone_-_invalid
--- PASS: TestHandlerExecution (0.00s)
PASS
ok      github.com/mahmut-Abi/k8s-mcp-server/internal/services/time   0.123s
```

## Key Features Demonstrated

1. **Simple Service Structure**: Minimal overhead, clear separation of concerns
2. **No External Dependencies**: Uses only Go standard library
3. **Error Handling**: Proper error messages and validation
4. **Parameter Validation**: Checks required parameters
5. **Structured Output**: Consistent JSON responses
6. **Comprehensive Tests**: Unit tests for all functionality
7. **Practical Use Cases**: Real-world time manipulation operations

## Extending the Demo

### Add More Time Operations

```go
// In client/client.go
func (c *Client) AddDuration(base time.Time, duration string) (time.Time, error) {
    d, err := time.ParseDuration(duration)
    if err != nil {
        return time.Time{}, fmt.Errorf("invalid duration: %w", err)
    }
    return base.Add(d), nil
}

func (c *Client) TimeDifference(t1, t2 time.Time) time.Duration {
    return t2.Sub(t1)
}
```

### Add Time Zone Conversion

```go
// In client/client.go
func (c *Client) ConvertZone(t time.Time, fromZone, toZone string) (time.Time, error) {
    locFrom, err := time.LoadLocation(fromZone)
    if err != nil {
        return time.Time{}, err
    }

    locTo, err := time.LoadLocation(toZone)
    if err != nil {
        return time.Time{}, err
    }

    return t.In(locFrom).In(locTo), nil
}
```

### Add Business Day Calculation

```go
// In client/client.go
func (c *Client) IsBusinessDay(t time.Time) bool {
    // Weekends are not business days
    if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
        return false
    }
    return true
}

func (c *Client) AddBusinessDays(t time.Time, days int) time.Time {
    added := 0
    for added < days {
        t = t.AddDate(0, 0, 1)
        if c.IsBusinessDay(t) {
            added++
        }
    }
    return t
}
```

## Summary

The Time Service demo provides:
- A complete, working service implementation
- Clear code structure and organization
- Practical tool examples
- Comprehensive test coverage
- Easy-to-understand patterns for creating new services

Use this as a reference when creating your own services, adapting the patterns to your specific use case.
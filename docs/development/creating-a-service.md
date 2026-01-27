# Creating a New Service

This guide explains how to create a new service for the cloud-native-mcp-server.

## Overview

A service in cloud-native-mcp-server is a module that exposes MCP tools for interacting with an external system (like Grafana, Prometheus, etc.). Each service typically consists of:

- A `service.go` file that implements the `Service` interface
- A `client/` directory with HTTP client code
- A `handlers/` directory with tool handlers
- A `tools/` directory with tool definitions

## Service Structure

```
internal/services/yourservice/
├── client/
│   └── client.go          # HTTP client implementation
├── handlers/
│   └── handlers.go        # Tool handlers
├── tools/
│   └── tools.go          # Tool definitions
├── service.go            # Service implementation
└── service_test.go       # Service tests
```

## Step-by-Step Guide

### Step 1: Create the Service Implementation

Create `internal/services/yourservice/service.go`:

```go
package yourservice

import (
    "context"

    "github.com/mark3labs/mcp-go/mcp"
    server "github.com/mark3labs/mcp-go/server"

    "github.com/mahmut-Abi/cloud-native-mcp-server/internal/config"
    "github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/cache"
    "github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/framework"
    "github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/yourservice/client"
    "github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/yourservice/handlers"
    "github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/yourservice/tools"
)

// Service implements the YourService MCP service
type Service struct {
    client        *client.Client
    enabled       bool
    toolsCache    *cache.ToolsCache
    initFramework *framework.CommonServiceInit
}

// NewService creates a new YourService instance
func NewService() *Service {
    // Create service enable checker
    checker := framework.NewServiceEnabled(
        func(cfg *config.AppConfig) bool { return cfg.YourService.Enabled },
        func(cfg *config.AppConfig) string { return cfg.YourService.URL },
    )

    // Create init configuration
    initConfig := &framework.InitConfig{
        Required:     false,
        URLValidator: framework.SimpleURLValidator,
        ClientBuilder: func(cfg *config.AppConfig) (interface{}, error) {
            opts := client.DefaultClientOptions()
            if cfg.YourService.URL != "" {
                opts.URL = cfg.YourService.URL
            }
            if cfg.YourService.APIKey != "" {
                opts.APIKey = cfg.YourService.APIKey
            }
            if cfg.YourService.TimeoutSec > 0 {
                opts.Timeout = time.Duration(cfg.YourService.TimeoutSec) * time.Second
            }
            return client.NewClient(opts)
        },
    }

    return &Service{
        enabled:       false,
        toolsCache:    cache.NewToolsCache(),
        initFramework: framework.NewCommonServiceInit("YourService", initConfig, checker),
    }
}

// Name returns the service identifier
func (s *Service) Name() string {
    return "yourservice"
}

// Initialize configures the service with the provided configuration
func (s *Service) Initialize(cfg interface{}) error {
    return s.initFramework.Initialize(cfg,
        func(enabled bool) { s.enabled = enabled },
        func(clientIface interface{}) {
            if svcClient, ok := clientIface.(*client.Client); ok {
                s.client = svcClient
            }
        },
    )
}

// GetTools returns all available YourService MCP tools
func (s *Service) GetTools() []mcp.Tool {
    if !s.enabled || s.client == nil {
        return nil
    }

    return s.toolsCache.Get(func() []mcp.Tool {
        return []mcp.Tool{
            tools.GetResourceTool(),
            tools.ListResourcesTool(),
            tools.CreateResourceTool(),
            tools.UpdateResourceTool(),
            tools.DeleteResourceTool(),
        }
    })
}

// GetHandlers returns all tool handlers mapped to their respective tool names
func (s *Service) GetHandlers() map[string]server.ToolHandlerFunc {
    if !s.enabled || s.client == nil {
        return nil
    }

    return map[string]server.ToolHandlerFunc{
        "yourservice_get_resource":    handlers.HandleGetResource(s.client),
        "yourservice_list_resources":  handlers.HandleListResources(s.client),
        "yourservice_create_resource": handlers.HandleCreateResource(s.client),
        "yourservice_update_resource": handlers.HandleUpdateResource(s.client),
        "yourservice_delete_resource": handlers.HandleDeleteResource(s.client),
    }
}

// IsEnabled returns whether the service is enabled and ready for use
func (s *Service) IsEnabled() bool {
    return s.enabled && s.client != nil
}

// GetClient returns the underlying YourService client
func (s *Service) GetClient() *client.Client {
    return s.client
}
```

### Step 2: Create the HTTP Client

Create `internal/services/yourservice/client/client.go`:

```go
package client

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "time"

    "github.com/sirupsen/logrus"

    optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
)

var logger = logrus.WithField("component", "yourservice-client")

// ClientOptions holds configuration parameters for creating a YourService client
type ClientOptions struct {
    URL     string        // YourService server URL
    APIKey  string        // API key for authentication
    Timeout time.Duration // HTTP request timeout
}

// DefaultClientOptions returns default client options
func DefaultClientOptions() *ClientOptions {
    return &ClientOptions{
        Timeout: 30 * time.Second,
    }
}

// Client provides operations for interacting with YourService API
type Client struct {
    baseURL    string
    httpClient *http.Client
    apiKey     string
}

// NewClient creates a new YourService client
func NewClient(opts *ClientOptions) (*Client, error) {
    if opts.URL == "" {
        return nil, fmt.Errorf("yourservice URL is required")
    }

    // Parse and validate URL
    baseURL, err := url.Parse(opts.URL)
    if err != nil {
        return nil, fmt.Errorf("invalid yourservice URL: %w", err)
    }

    // Ensure URL has proper path
    if !baseURL.Path.EndsWith("/") {
        baseURL.Path += "/"
    }
    baseURL.Path += "api/"

    // Create HTTP client
    timeout := opts.Timeout
    if timeout == 0 {
        timeout = 30 * time.Second
    }

    httpClient := optimize.NewOptimizedHTTPClientWithTimeout(timeout)

    return &Client{
        baseURL:    baseURL.String(),
        httpClient: httpClient,
        apiKey:     opts.APIKey,
    }, nil
}

// GetResource retrieves a single resource
func (c *Client) GetResource(ctx context.Context, id string) (*Resource, error) {
    path := fmt.Sprintf("resources/%s", id)
    resp, err := c.makeRequest(ctx, "GET", path, nil)
    if err != nil {
        return nil, err
    }

    body, err := c.handleResponse(resp)
    if err != nil {
        return nil, err
    }

    var resource Resource
    if err := json.Unmarshal(body, &resource); err != nil {
        return nil, fmt.Errorf("failed to unmarshal response: %w", err)
    }

    return &resource, nil
}

// makeRequest performs an HTTP request
func (c *Client) makeRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
    var reqBody io.Reader
    if body != nil {
        jsonBody, err := json.Marshal(body)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal request body: %w", err)
        }
        reqBody = bytes.NewReader(jsonBody)
    }

    req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    // Set headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    if c.apiKey != "" {
        req.Header.Set("Authorization", "Bearer "+c.apiKey)
    }

    return c.httpClient.Do(req)
}

// handleResponse processes the HTTP response
func (c *Client) handleResponse(resp *http.Response) ([]byte, error) {
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %w", err)
    }

    if resp.StatusCode >= 400 {
        return nil, fmt.Errorf("yourservice API error (status %d): %s", resp.StatusCode, string(body))
    }

    return body, nil
}

// Resource represents a YourService resource
type Resource struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    // Add other fields as needed
}
```

### Step 3: Create Tool Definitions

Create `internal/services/yourservice/tools/tools.go`:

```go
package tools

import (
    "github.com/mark3labs/mcp-go/mcp"
)

// GetResourceTool returns a tool for getting a single resource
func GetResourceTool() mcp.Tool {
    return mcp.NewTool("yourservice_get_resource",
        mcp.WithDescription("Get a single resource from YourService"),
        mcp.WithString("id",
            mcp.Required(),
            mcp.Description("The ID of the resource to retrieve"),
        ),
    )
}

// ListResourcesTool returns a tool for listing resources
func ListResourcesTool() mcp.Tool {
    return mcp.NewTool("yourservice_list_resources",
        mcp.WithDescription("List all resources from YourService"),
        mcp.WithNumber("limit",
            mcp.Description("Maximum number of resources to return"),
            mcp.DefaultNumber(20),
        ),
    )
}

// CreateResourceTool returns a tool for creating a resource
func CreateResourceTool() mcp.Tool {
    return mcp.NewTool("yourservice_create_resource",
        mcp.WithDescription("Create a new resource in YourService"),
        mcp.WithString("name",
            mcp.Required(),
            mcp.Description("The name of the resource"),
        ),
        mcp.WithString("description",
            mcp.Description("Description of the resource"),
        ),
    )
}

// UpdateResourceTool returns a tool for updating a resource
func UpdateResourceTool() mcp.Tool {
    return mcp.NewTool("yourservice_update_resource",
        mcp.WithDescription("Update an existing resource in YourService"),
        mcp.WithString("id",
            mcp.Required(),
            mcp.Description("The ID of the resource to update"),
        ),
        mcp.WithString("name",
            mcp.Description("The new name of the resource"),
        ),
    )
}

// DeleteResourceTool returns a tool for deleting a resource
func DeleteResourceTool() mcp.Tool {
    return mcp.NewTool("yourservice_delete_resource",
        mcp.WithDescription("Delete a resource from YourService"),
        mcp.WithString("id",
            mcp.Required(),
            mcp.Description("The ID of the resource to delete"),
        ),
    )
}
```

### Step 4: Create Tool Handlers

Create `internal/services/yourservice/handlers/handlers.go`:

```go
package handlers

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/mark3labs/mcp-go/mcp"
    "github.com/sirupsen/logrus"

    "github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/yourservice/client"
)

// HandleGetResource handles the get_resource tool
func HandleGetResource(svcClient *client.Client) mcp.ToolHandlerFunc {
    return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        // Get parameters
        id, err := request.Params.Arguments.String("id")
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("invalid id: %v", err)), nil
        }

        // Call client
        resource, err := svcClient.GetResource(ctx, id)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("failed to get resource: %v", err)), nil
        }

        // Return result
        resultJSON, _ := json.MarshalIndent(resource, "", "  ")
        return mcp.NewToolResultText(string(resultJSON)), nil
    }
}

// HandleListResources handles the list_resources tool
func HandleListResources(svcClient *client.Client) mcp.ToolHandlerFunc {
    return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        // Get parameters
        limit := 20
        if request.Params.Arguments != nil {
            if l, ok := request.Params.Arguments["limit"].(float64); ok {
                limit = int(l)
            }
        }

        // Call client (implement this method)
        resources, err := svcClient.ListResources(ctx, limit)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("failed to list resources: %v", err)), nil
        }

        // Return result
        resultJSON, _ := json.MarshalIndent(resources, "", "  ")
        return mcp.NewToolResultText(string(resultJSON)), nil
    }
}

// HandleCreateResource handles the create_resource tool
func HandleCreateResource(svcClient *client.Client) mcp.ToolHandlerFunc {
    return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        // Get parameters
        name, err := request.Params.Arguments.String("name")
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("invalid name: %v", err)), nil
        }

        description, _ := request.Params.Arguments.String("description")

        // Call client (implement this method)
        resource, err := svcClient.CreateResource(ctx, name, description)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("failed to create resource: %v", err)), nil
        }

        // Return result
        resultJSON, _ := json.MarshalIndent(resource, "", "  ")
        return mcp.NewToolResultText(string(resultJSON)), nil
    }
}

// HandleUpdateResource handles the update_resource tool
func HandleUpdateResource(svcClient *client.Client) mcp.ToolHandlerFunc {
    return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        // Get parameters
        id, err := request.Params.Arguments.String("id")
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("invalid id: %v", err)), nil
        }

        name, _ := request.Params.Arguments.String("name")

        // Call client (implement this method)
        resource, err := svcClient.UpdateResource(ctx, id, name)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("failed to update resource: %v", err)), nil
        }

        // Return result
        resultJSON, _ := json.MarshalIndent(resource, "", "  ")
        return mcp.NewToolResultText(string(resultJSON)), nil
    }
}

// HandleDeleteResource handles the delete_resource tool
func HandleDeleteResource(svcClient *client.Client) mcp.ToolHandlerFunc {
    return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        // Get parameters
        id, err := request.Params.Arguments.String("id")
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("invalid id: %v", err)), nil
        }

        // Call client (implement this method)
        err = svcClient.DeleteResource(ctx, id)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("failed to delete resource: %v", err)), nil
        }

        return mcp.NewToolResultText(fmt.Sprintf("Resource %s deleted successfully", id)), nil
    }
}
```

### Step 5: Register the Service

Update `internal/services/manager/manager.go` to register your service:

```go
// In NewManager(), add:
m.yourServiceService = yourservice.NewService()

// In Initialize(), add:
m.yourServiceService = yourservice.NewService()
m.registry.Register(m.yourServiceService)

// In Shutdown(), add:
if m.yourServiceService != nil {
    if closer, ok := m.yourServiceService.(Closer); ok {
        if err := closer.Close(); err != nil {
            errs = append(errs, fmt.Errorf("%s service close error: %w", "yourservice", err))
        }
    }
}
```

### Step 6: Add Configuration

Update `internal/config/config.go` to add configuration for your service:

```go
type AppConfig struct {
    // ... existing fields ...

    YourService struct {
        Enabled    bool   `yaml:"enabled"`
        URL        string `yaml:"url"`
        APIKey     string `yaml:"apiKey"`
        TimeoutSec int    `yaml:"timeoutSec"`
    } `yaml:"yourservice"`
}
```

### Step 7: Add Configuration Loader

Update `internal/config/env_parser.go` to parse environment variables:

```go
// Add to Parse() function:
if v, ok := over("MCP_YOURSERVICE_ENABLED"); ok {
    cfg.YourService.Enabled = v == "true"
}
if v, ok := over("MCP_YOURSERVICE_URL"); ok {
    cfg.YourService.URL = v
}
if v, ok := over("MCP_YOURSERVICE_API_KEY"); ok {
    cfg.YourService.APIKey = v
}
if v, ok := over("MCP_YOURSERVICE_TIMEOUT"); ok {
    if val, err := strconv.Atoi(v); err == nil {
        cfg.YourService.TimeoutSec = val
    }
}
```

### Step 8: Add Configuration Validation

Update `internal/config/config.go` in the `Validate()` method:

```go
// Validate YourService configuration
if c.YourService.Enabled && c.YourService.URL == "" {
    return fmt.Errorf("yourservice URL is required when service is enabled")
}
```

## Testing

Create tests for your service:

```go
package yourservice

import (
    "testing"
)

func TestNewService(t *testing.T) {
    service := NewService()
    if service == nil {
        t.Fatal("NewService() returned nil")
    }
    if service.Name() != "yourservice" {
        t.Errorf("Expected name 'yourservice', got '%s'", service.Name())
    }
}

func TestServiceInitialization(t *testing.T) {
    service := NewService()
    config := &config.AppConfig{
        YourService: config.YourService{
            Enabled: true,
            URL:     "http://localhost:8080",
        },
    }
    
    err := service.Initialize(config)
    if err != nil {
        t.Fatalf("Failed to initialize service: %v", err)
    }
    
    if !service.IsEnabled() {
        t.Error("Service should be enabled after initialization")
    }
}
```

## Best Practices

1. **Use the Common Framework**: Always use `framework.CommonServiceInit` for consistent initialization
2. **Cache Tools**: Use `cache.ToolsCache` to avoid recreating tools on each call
3. **Error Handling**: Use the project's error handling utilities (`internal/errors`)
4. **Logging**: Use structured logging with `logrus`
5. **Metrics**: Use the hook system to automatically record metrics
6. **Testing**: Write comprehensive tests for your service
7. **Documentation**: Document all public functions and types

## Common Patterns

### Authentication

If your service uses API key authentication:

```go
func (c *Client) makeRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
    req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
    if err != nil {
        return nil, err
    }
    
    // Add API key header
    if c.apiKey != "" {
        req.Header.Set("Authorization", "Bearer "+c.apiKey)
    }
    
    return c.httpClient.Do(req)
}
```

### Pagination

Implement pagination using the common types:

```go
import (
    "github.com/mahmut-Abi/cloud-native-mcp-server/internal/common"
)

func (c *Client) ListResources(ctx context.Context, opts *common.PaginationOptions) (*PaginatedResult, error) {
    // Validate options
    if err := common.ValidatePaginationOptions(opts); err != nil {
        return nil, err
    }
    
    // Build query parameters
    params := url.Values{}
    params.Set("limit", strconv.Itoa(opts.Limit))
    if opts.Continue != "" {
        params.Set("continue", opts.Continue)
    }
    
    // Make request
    // ...
}
```

### Response Handling

Use the common response handler:

```go
import (
    "github.com/mahmut-Abi/cloud-native-mcp-server/internal/common"
)

func (c *Client) makeRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
    // ... make request ...
    
    // Handle response
    handler := common.NewResponseHandler("yourservice")
    body, err := handler.HandleResponse(resp)
    if err != nil {
        return nil, err
    }
    
    return body, nil
}
```

## Next Steps

After creating your service:

1. Add integration tests
2. Update the documentation
3. Add examples to the README
4. Test with actual MCP clients
5. Monitor metrics and logs in production

## Examples

For complete examples, refer to existing services:
- `internal/services/grafana/` - Basic HTTP service with API key auth
- `internal/services/prometheus/` - Service with bearer token auth
- `internal/services/kubernetes/` - Complex service with caching

## Support

If you encounter issues:
1. Check existing service implementations for patterns
2. Review the test files for usage examples
3. Consult the MCP specification for protocol details
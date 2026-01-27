# Development Guide

This directory contains guides and examples for developing services in cloud-native-mcp-server.

## Documentation

### [Creating a New Service](./creating-a-service.md)
A comprehensive guide that explains:
- Service structure and architecture
- Step-by-step service creation process
- Best practices and patterns
- Configuration management
- Testing strategies

### [Minimal Service Example](./service-example.md)
A complete, minimal example that demonstrates:
- Service implementation from scratch
- Client code patterns
- Tool and handler definitions
- Service registration
- Testing examples

## Quick Start

1. **Read the guide**: Start with [Creating a New Service](./creating-a-service.md) to understand the architecture
2. **Study the example**: Review [Minimal Service Example](./service-example.md) for a working reference
3. **Create your service**: Follow the steps to implement your own service
4. **Test thoroughly**: Use the testing patterns shown in the examples

## Service Architecture

A typical service consists of:

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

## Common Patterns

### Service Initialization
- Use `framework.CommonServiceInit` for consistent initialization
- Implement the `Service` interface
- Use `cache.ToolsCache` for tool caching

### Client Implementation
- Use the common HTTP client builder when possible
- Implement proper error handling
- Add connection pooling and timeouts

### Tool Definitions
- Use `mcp.NewTool()` to create tools
- Add clear descriptions
- Define required and optional parameters

### Handler Implementation
- Parse parameters from `request.Params.Arguments`
- Call client methods
- Return structured results using `mcp.NewToolResultText()`

## Testing

Write tests for:
- Service initialization
- Tool registration
- Handler execution
- Client operations
- Error scenarios

## Best Practices

1. **Use existing services as reference** - Check Grafana, Prometheus, Kubernetes services
2. **Follow naming conventions** - Use `{service}_{action}` pattern for tool names
3. **Handle errors gracefully** - Use the project's error handling utilities
4. **Log appropriately** - Use structured logging with context
5. **Test thoroughly** - Write unit and integration tests
6. **Document your code** - Add comments for complex logic

## Support

For questions or issues:
- Review existing service implementations
- Check test files for usage examples
- Consult the MCP specification
- Look at the project architecture documentation

## Related Documentation

- [Architecture Overview](../ARCHITECTURE.md) - Overall system architecture
- [Configuration Guide](../CONFIGURATION.md) - How to configure services
- [Deployment Guide](../DEPLOYMENT.md) - How to deploy the server
- [Performance Guide](../PERFORMANCE.md) - Performance optimization tips
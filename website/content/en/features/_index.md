---
title: "Features"
---

# Core Features

## 🚀 High-Performance Architecture

### Smart Caching
- **LRU Cache**: Least Recently Used algorithm for automatic cache management
- **TTL Support**: Configurable cache expiration times
- **Intelligent Pre-warming**: Automatic caching of frequently accessed data

### Optimized Design
- **JSON Encoding Pool**: Reduces memory allocation, improves encoding efficiency
- **Response Size Control**: Intelligent limit on response sizes to prevent context overflow
- **Concurrency Optimization**: Supports high concurrent request processing

## 🔒 Secure & Reliable

### Multiple Authentication Methods
- **API Key Authentication**: Simple and efficient API key auth with complexity requirements
- **Bearer Token Authentication**: JWT-based token authentication
- **Basic Authentication**: HTTP Basic Auth

### Secrets Management
- **Secure Storage**: In-memory storage with expiration support
- **Secret Rotation**: Automatic rotation for API keys and bearer tokens
- **Secret Generation**: Built-in generators for complex API keys and JWT-like tokens
- **Environment Variables**: Support for loading secrets from environment variables

### Input Sanitization
- **Filter Values**: Remove dangerous characters (SQL injection, XSS, command injection)
- **URL Validation**: Only allow http/https schemes for web fetch
- **Length Limits**: Maximum string length enforcement (1000 characters)
- **Special Character Removal**: Remove semicolons, quotes, and other injection vectors

## 📊 Comprehensive Cloud-Native Services

### Kubernetes Management
- 28 tools covering core container orchestration
- Pod, Deployment, Service management
- ConfigMap, Secret operations
- Namespace, Node management

### Application Deployment
- 31 Helm tools
- Chart management, installation, upgrades
- Release operations
- Repository management

### Monitoring and Observability
- **Prometheus**: 20 tools, metrics collection and querying
- **Grafana**: 36 tools, dashboards and visualization
- **Jaeger**: 8 tools, distributed tracing
- **OpenTelemetry**: 9 tools, unified observability

### Log Management
- **Elasticsearch**: 14 tools, log storage and search
- **Kibana**: 52 tools, log analysis and visualization
- Index management, data exploration

### Alert Management
- **Alertmanager**: 15 tools
- Alert rules management
- Notification routing
- Silence management

## 🔧 Flexible Running Modes

### SSE Mode (Server-Sent Events)
- Real-time event push
- Low latency communication
- Suitable for real-time monitoring scenarios

### HTTP Mode
- Standard REST API
- Easy integration
- Suitable for traditional applications

### stdio Mode
- Standard input/output communication
- Suitable for development environments
- MCP standard protocol

### Streamable-HTTP Mode
- MCP 2025-11-25 specification
- Single endpoint for both request/response and streaming
- Modern communication method

## 🤖 AI-Optimized Design

### LLM-Friendly
- **Summary Tools**: Automatically generate tool result summaries
- **Pagination Support**: Prevent context overflow
- **Smart Filtering**: Sort results based on relevance

### Context Management
- **Response Size Limits**: Automatically truncate overly long responses
- **Key Information Extraction**: Highlight important data
- **Multi-turn Conversation Support**: Maintain context coherence

## 📝 Observability

### Metrics Collection
- Built-in Prometheus metrics
- Custom metrics support
- Performance monitoring

### Logging
- Structured logging
- Multiple log levels (debug, info, warn, error)
- JSON format support

### Distributed Tracing
- OpenTelemetry integration
- Jaeger tracing support
- Performance analysis

## 🛠️ Developer-Friendly

### Flexible Configuration
- YAML configuration file
- Environment variable overrides
- Runtime configuration updates

### Easy Testing
- Complete unit tests
- Integration test support
- Mock tools

### Comprehensive Documentation
- Detailed tool documentation
- Configuration examples
- Best practices guide

## 🌟 Additional Features

- **Service Filtering**: Enable/disable services on demand
- **Tool Filtering**: Fine-grained control over available tools
- **Rate Limiting**: Prevent abuse
- **Timeout Control**: Protect system stability
- **Retry Mechanism**: Improve reliability
- **Error Handling**: Friendly error messages
- **Health Checks**: Service status monitoring
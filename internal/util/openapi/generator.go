package openapi

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/mahmut-Abi/k8s-mcp-server/internal/services"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Generator generates OpenAPI specifications for MCP tools
type Generator struct {
	registry *services.Registry
	spec     *OpenAPISpec
}

// NewGenerator creates a new OpenAPI generator
func NewGenerator(registry *services.Registry) *Generator {
	return &Generator{
		registry: registry,
	}
}

// Generate generates the complete OpenAPI specification
func (g *Generator) Generate() (*OpenAPISpec, error) {
	logrus.Debug("Generating OpenAPI specification")

	// Initialize the spec with enhanced information
	g.spec = &OpenAPISpec{
		OpenAPI: "3.0.0",
		Info: Info{
			Title:       "Kubernetes MCP Server API",
			Description: "Model Context Protocol Server for Kubernetes cluster management with integrated observability tools. Supports SSE, HTTP, and stdio communication modes.",
			Version:     "1.0.0",
			Contact: &map[string]interface{}{
				"name": "Kubernetes MCP Team",
				"url":  "https://github.com/mahmut-Abi/k8s-mcp-server",
			},
			License: &map[string]interface{}{
				"name": "MIT",
				"url":  "https://opensource.org/licenses/MIT",
			},
		},
		Servers: []Server{
			{
				URL:         "http://localhost:8080",
				Description: "Local development server",
				Variables: map[string]map[string]interface{}{
					"port": {
						"default":     "8080",
						"description": "Server port",
					},
				},
			},
			{
				URL:         "https://api.example.com",
				Description: "Production server",
			},
		},
		Paths: make(map[string]PathItem),
		Components: Components{
			Schemas:         make(map[string]interface{}),
			SecuritySchemes: make(map[string]SecurityScheme),
		},
		Tags: []Tag{
			{Name: "Tools", Description: "MCP Tool management and execution with detailed schema information"},
			{Name: "System", Description: "Server status, health checks, and SSE streams"},
			{Name: "Observability", Description: "Integration with Prometheus, Grafana, Kibana, and Elasticsearch"},
		},
		Security: []map[string][]string{
			{"bearerAuth": {}},
			{"apiKey": {}},
		},
	}

	// Add health check endpoint
	g.addHealthEndpoint()

	// Add tools endpoints
	g.addToolsEndpoints()

	// Add SSE endpoints
	g.addSSEEndpoints()

	// Add security schemes
	g.addSecuritySchemes()

	// Add mode-specific information
	g.addModeInformation()

	logrus.Debug("OpenAPI specification generated successfully")
	return g.spec, nil
}

// addHealthEndpoint adds the health check endpoint
func (g *Generator) addHealthEndpoint() {
	g.spec.Paths["/health"] = PathItem{
		Get: &Operation{
			Summary:     "Health check",
			Description: "Returns server health status and information about enabled services. This endpoint can be used for health monitoring and readiness checks.",
			OperationID: "healthCheck",
			Tags:        []string{"System"},
			Responses: map[string]Response{
				"200": {
					Description: "Server is healthy and operational",
					Content: map[string]MediaType{
						"application/json": {
							Schema: map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"status": map[string]interface{}{
										"type":        "string",
										"description": "Overall health status of the server",
										"enum":        []string{"healthy", "degraded", "unhealthy"},
									},
									"timestamp": map[string]interface{}{
										"type":        "string",
										"description": "ISO 8601 timestamp of when the health check was performed",
										"format":      "date-time",
									},
									"services": map[string]interface{}{
										"type":        "object",
										"description": "Health status of individual services",
										"additionalProperties": map[string]interface{}{
											"type": "object",
											"properties": map[string]interface{}{
												"status": map[string]interface{}{
													"type": "string",
													"enum": []string{"healthy", "degraded", "unhealthy"},
												},
												"version": map[string]interface{}{
													"type": "string",
												},
											},
										},
									},
									"version": map[string]interface{}{
										"type":        "string",
										"description": "Server version",
									},
								},
								"required": []string{"status", "timestamp"},
							},
							Example: map[string]interface{}{
								"status":    "healthy",
								"timestamp": "2023-01-01T12:00:00Z",
								"services": map[string]interface{}{
									"kubernetes": map[string]interface{}{
										"status":  "healthy",
										"version": "1.25.0",
									},
									"prometheus": map[string]interface{}{
										"status":  "healthy",
										"version": "2.37.0",
									},
								},
								"version": "1.0.0",
							},
						},
					},
				},
				"503": {
					Description: "Server is unhealthy or in maintenance mode",
					Content: map[string]MediaType{
						"application/json": {
							Schema: map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"status": map[string]interface{}{
										"type": "string",
										"enum": []string{"unhealthy", "maintenance"},
									},
									"timestamp": map[string]interface{}{
										"type":   "string",
										"format": "date-time",
									},
									"services": map[string]interface{}{
										"type": "object",
										"additionalProperties": map[string]interface{}{
											"type": "object",
											"properties": map[string]interface{}{
												"status": map[string]interface{}{
													"type": "string",
													"enum": []string{"healthy", "degraded", "unhealthy"},
												},
												"error": map[string]interface{}{
													"type": "string",
												},
											},
										},
									},
									"error": map[string]interface{}{
										"type": "string",
									},
								},
								"required": []string{"status", "timestamp"},
							},
							Example: map[string]interface{}{
								"status":    "unhealthy",
								"timestamp": "2023-01-01T12:00:00Z",
								"services": map[string]interface{}{
									"kubernetes": map[string]interface{}{
										"status": "unhealthy",
										"error":  "Connection to Kubernetes API failed",
									},
								},
								"error": "One or more services are unhealthy",
							},
						},
					},
				},
			},
		},
	}
}

// getServiceFromToolName extracts service name from tool name
func getServiceFromToolName(toolName string) string {
	parts := strings.Split(toolName, "_")
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}

// addToolsEndpoints adds the tools-related endpoints
func (g *Generator) addToolsEndpoints() {
	// Get all tools from enabled services
	allTools := g.registry.GetAllTools()

	// Group tools by service
	toolsByService := make(map[string][]mcp.Tool)
	serviceNames := make(map[string]bool)

	for _, tool := range allTools {
		// Extract service name from tool name (prefix before first underscore)
		parts := strings.Split(tool.Name, "_")
		if len(parts) > 0 {
			service := parts[0]
			toolsByService[service] = append(toolsByService[service], tool)
			serviceNames[service] = true
		}
	}

	// Add tools list endpoint with enhanced documentation
	g.spec.Paths["/tools"] = PathItem{
		Get: &Operation{
			Summary:     "List all available MCP tools",
			Description: "Returns comprehensive list of all MCP tools exposed by the server, including detailed schema for each tool. Tools are organized by service categories.",
			OperationID: "listTools",
			Tags:        []string{"Tools"},
			Parameters: []Parameter{
				{
					Name:        "service",
					In:          "query",
					Description: "Filter tools by service name (kubernetes, prometheus, grafana, kibana, helm, elasticsearch)",
					Required:    false,
					Schema: map[string]interface{}{
						"type": "string",
					},
				},
				{
					Name:        "category",
					In:          "query",
					Description: "Filter tools by category",
					Required:    false,
					Schema: map[string]interface{}{
						"type": "string",
						"enum": []string{"kubernetes", "prometheus", "grafana", "kibana", "helm", "elasticsearch"},
					},
				},
			},
			Responses: map[string]Response{
				"200": {
					Description: "Successfully retrieved comprehensive tools list with detailed information",
					Content: map[string]MediaType{
						"application/json": {
							Schema: map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"tools": map[string]interface{}{
										"type": "array",
										"items": map[string]interface{}{
											"$ref": "#/components/schemas/ToolDetail",
										},
									},
									"total": map[string]interface{}{
										"type":        "integer",
										"description": "Total number of tools available",
									},
									"categories": map[string]interface{}{
										"type":        "object",
										"description": "Count of tools per category",
										"additionalProperties": map[string]interface{}{
											"type": "integer",
										},
									},
									"services": map[string]interface{}{
										"type":        "object",
										"description": "Count of tools per service",
										"additionalProperties": map[string]interface{}{
											"type": "integer",
										},
									},
								},
							},
							Example: map[string]interface{}{
								"tools": []interface{}{
									map[string]interface{}{
										"name":        "kubernetes_get_pod",
										"description": "Get details of a specific pod in a namespace",
										"service":     "kubernetes",
										"inputSchema": map[string]interface{}{
											"type": "object",
											"properties": map[string]interface{}{
												"namespace": map[string]interface{}{
													"type": "string",
												},
												"name": map[string]interface{}{
													"type": "string",
												},
											},
											"required": []string{"namespace", "name"},
										},
										"examples": []interface{}{
											map[string]interface{}{
												"description": "Get a pod named 'my-pod' in namespace 'default'",
												"arguments": map[string]interface{}{
													"namespace": "default",
													"name":      "my-pod",
												},
											},
										},
									},
								},
								"total": 1,
								"categories": map[string]interface{}{
									"kubernetes": 1,
								},
								"services": map[string]interface{}{
									"kubernetes": 1,
								},
							},
						},
					},
				},
				"400": {
					Description: "Invalid request parameters",
					Content: map[string]MediaType{
						"application/json": {
							Schema: map[string]interface{}{
								"$ref": "#/components/schemas/Error",
							},
						},
					},
				},
				"500": {
					Description: "Server error occurred while retrieving tools",
					Content: map[string]MediaType{
						"application/json": {
							Schema: map[string]interface{}{
								"$ref": "#/components/schemas/Error",
							},
						},
					},
				},
			},
		},
	}

	// Add tool execution endpoint with enhanced documentation
	g.spec.Paths["/tools/call"] = PathItem{
		Post: &Operation{
			Summary:     "Execute a MCP tool",
			Description: "Invokes a specific MCP tool with provided arguments and returns the result. Supports both synchronous and asynchronous execution modes.",
			OperationID: "callTool",
			Tags:        []string{"Tools"},
			RequestBody: &RequestBody{
				Required: true,
				Content: map[string]MediaType{
					"application/json": {
						Schema: map[string]interface{}{
							"$ref": "#/components/schemas/ToolCall",
						},
						Example: map[string]interface{}{
							"name": "kubernetes_get_pod",
							"arguments": map[string]interface{}{
								"namespace": "default",
								"name":      "my-pod",
							},
							"timeout": 30000,
						},
					},
				},
			},
			Responses: map[string]Response{
				"200": {
					Description: "Tool executed successfully",
					Content: map[string]MediaType{
						"application/json": {
							Schema: map[string]interface{}{
								"$ref": "#/components/schemas/ToolResult",
							},
							Example: map[string]interface{}{
								"content": []interface{}{
									map[string]interface{}{
										"type": "text",
										"text": "{\n  \"metadata\": {\n    \"name\": \"my-pod\",\n    \"namespace\": \"default\"\n  },\n  \"spec\": {\n    \"containers\": [\n      {\n        \"name\": \"my-container\",\n        \"image\": \"nginx:latest\"\n      }\n    ]\n  }\n}",
									},
								},
								"isError":       false,
								"executionTime": 150,
								"cached":        false,
							},
						},
					},
				},
				"400": {
					Description: "Invalid tool call request",
					Content: map[string]MediaType{
						"application/json": {
							Schema: map[string]interface{}{
								"$ref": "#/components/schemas/Error",
							},
							Example: map[string]interface{}{
								"code":    400,
								"message": "Invalid tool name or arguments",
								"details": map[string]interface{}{
									"error": "Tool 'invalid_tool' not found",
								},
							},
						},
					},
				},
				"404": {
					Description: "Tool not found",
					Content: map[string]MediaType{
						"application/json": {
							Schema: map[string]interface{}{
								"$ref": "#/components/schemas/Error",
							},
							Example: map[string]interface{}{
								"code":    404,
								"message": "Tool not found",
								"details": map[string]interface{}{
									"toolName": "nonexistent_tool",
								},
							},
						},
					},
				},
				"500": {
					Description: "Tool execution failed",
					Content: map[string]MediaType{
						"application/json": {
							Schema: map[string]interface{}{
								"$ref": "#/components/schemas/Error",
							},
							Example: map[string]interface{}{
								"code":    500,
								"message": "Tool execution failed",
								"details": map[string]interface{}{
									"toolName": "kubernetes_get_pod",
									"error":    "pods \"my-pod\" not found",
								},
							},
						},
					},
				},
			},
		},
	}

	// Add individual tool endpoints with detailed documentation
	for _, tool := range allTools {
		path := fmt.Sprintf("/tools/%s", tool.Name)
		g.spec.Paths[path] = PathItem{
			Get: &Operation{
				Summary:     fmt.Sprintf("Get detailed information for tool: %s", tool.Name),
				Description: fmt.Sprintf("Returns complete schema, examples, and documentation for the %s tool.", tool.Name),
				OperationID: fmt.Sprintf("getTool_%s", tool.Name),
				Tags:        []string{"Tools"},
				Responses: map[string]Response{
					"200": {
						Description: "Tool details retrieved successfully",
						Content: map[string]MediaType{
							"application/json": {
								Schema: map[string]interface{}{
									"$ref": "#/components/schemas/ToolDetail",
								},
								Example: map[string]interface{}{
									"name":        tool.Name,
									"description": tool.Description,
									"service":     getServiceFromToolName(tool.Name),
									"inputSchema": tool.InputSchema,
									"examples": []interface{}{
										map[string]interface{}{
											"description": "Example usage",
											"arguments":   map[string]interface{}{},
										},
									},
									"performanceHints": map[string]interface{}{
										"averageExecutionTime": "150ms",
										"cacheability":         true,
										"recommendations": []string{
											"Use with appropriate timeouts for long-running operations",
										},
									},
								},
							},
						},
					},
					"404": {
						Description: "Tool not found",
						Content: map[string]MediaType{
							"application/json": {
								Schema: map[string]interface{}{
									"$ref": "#/components/schemas/Error",
								},
								Example: map[string]interface{}{
									"code":    404,
									"message": "Tool not found",
									"details": map[string]interface{}{
										"toolName": tool.Name,
									},
								},
							},
						},
					},
				},
			},
		}
	}

	// Add schemas
	g.addSchemas()
}

// addSSEEndpoints adds the SSE endpoints
func (g *Generator) addSSEEndpoints() {
	services := []string{"kubernetes", "grafana", "prometheus", "kibana", "helm", "aggregate"}

	for _, service := range services {
		path := fmt.Sprintf("/api/%s/sse", service)
		messagePath := fmt.Sprintf("/api/%s/sse/message", service)

		// SSE endpoint with enhanced documentation
		g.spec.Paths[path] = PathItem{
			Get: &Operation{
				Summary:     fmt.Sprintf("%s SSE stream", cases.Title(language.English).String(service)),
				Description: fmt.Sprintf("Server-Sent Events stream for %s service. Provides real-time updates and streaming responses for long-running operations.", service),
				OperationID: fmt.Sprintf("sse_%s", service),
				Tags:        []string{"System"},
				Parameters: []Parameter{
					{
						Name:        "Authorization",
						In:          "header",
						Description: "Bearer token for authentication",
						Required:    false,
						Schema: map[string]interface{}{
							"type": "string",
						},
					},
				},
				Responses: map[string]Response{
					"200": {
						Description: "SSE stream established successfully. Events will be sent as they occur.",
						Content: map[string]MediaType{
							"text/event-stream": {
								Schema: map[string]interface{}{
									"type": "string",
								},
								Example: "data: {\"type\": \"text\", \"text\": \"Operation in progress...\"}\n\n" +
									"data: {\"type\": \"text\", \"text\": \"Operation completed\"}\n\n",
							},
						},
					},
					"401": {
						Description: "Authentication required",
						Content: map[string]MediaType{
							"application/json": {
								Schema: map[string]interface{}{
									"$ref": "#/components/schemas/Error",
								},
								Example: map[string]interface{}{
									"code":    401,
									"message": "Unauthorized",
								},
							},
						},
					},
					"500": {
						Description: "Server error occurred while establishing SSE stream",
						Content: map[string]MediaType{
							"application/json": {
								Schema: map[string]interface{}{
									"$ref": "#/components/schemas/Error",
								},
								Example: map[string]interface{}{
									"code":    500,
									"message": "Failed to establish SSE connection",
								},
							},
						},
					},
				},
			},
		}

		// Message endpoint with enhanced documentation
		g.spec.Paths[messagePath] = PathItem{
			Post: &Operation{
				Summary:     fmt.Sprintf("Send message to %s SSE stream", service),
				Description: fmt.Sprintf("Send a message to the %s SSE stream. Used for testing and debugging purposes.", service),
				OperationID: fmt.Sprintf("sendMessage_%s", service),
				Tags:        []string{"System"},
				RequestBody: &RequestBody{
					Required: true,
					Content: map[string]MediaType{
						"application/json": {
							Schema: map[string]interface{}{
								"type":     "object",
								"required": []string{"message"},
								"properties": map[string]interface{}{
									"message": map[string]interface{}{
										"type":        "string",
										"description": "The message to send to the SSE stream",
									},
								},
							},
							Example: map[string]interface{}{
								"message": "Test message for SSE stream",
							},
						},
					},
				},
				Responses: map[string]Response{
					"200": {
						Description: "Message sent successfully to SSE stream",
						Content: map[string]MediaType{
							"application/json": {
								Schema: map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"status": map[string]interface{}{
											"type": "string",
										},
										"message": map[string]interface{}{
											"type": "string",
										},
									},
								},
								Example: map[string]interface{}{
									"status":  "success",
									"message": "Message sent to SSE stream",
								},
							},
						},
					},
					"400": {
						Description: "Invalid message format",
						Content: map[string]MediaType{
							"application/json": {
								Schema: map[string]interface{}{
									"$ref": "#/components/schemas/Error",
								},
								Example: map[string]interface{}{
									"code":    400,
									"message": "Invalid request body",
									"details": map[string]interface{}{
										"error": "Missing required field 'message'",
									},
								},
							},
						},
					},
					"401": {
						Description: "Authentication required",
						Content: map[string]MediaType{
							"application/json": {
								Schema: map[string]interface{}{
									"$ref": "#/components/schemas/Error",
								},
								Example: map[string]interface{}{
									"code":    401,
									"message": "Unauthorized",
								},
							},
						},
					},
				},
			},
		}
	}
}

// addSchemas adds the component schemas
func (g *Generator) addSchemas() {
	// Enhanced Tool schema with more details
	g.spec.Components.Schemas["Tool"] = map[string]interface{}{
		"type":     "object",
		"required": []string{"name", "description"},
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Unique identifier for the tool",
			},
			"description": map[string]interface{}{
				"type":        "string",
				"description": "Human-readable description of what the tool does",
			},
			"service": map[string]interface{}{
				"type":        "string",
				"description": "Service that provides this tool (kubernetes, prometheus, grafana, etc.)",
			},
			"category": map[string]interface{}{
				"type":        "string",
				"description": "Category of the tool",
				"enum":        []string{"kubernetes", "prometheus", "grafana", "kibana", "helm", "elasticsearch"},
			},
			"inputSchema": map[string]interface{}{
				"type":        "object",
				"description": "JSON Schema defining the expected input parameters",
			},
			"version": map[string]interface{}{
				"type":        "string",
				"description": "Tool version",
			},
			"deprecated": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether this tool is deprecated",
				"default":     false,
			},
		},
		"example": map[string]interface{}{
			"name":        "kubernetes_get_pod",
			"description": "Get details of a specific pod in a namespace",
			"service":     "kubernetes",
			"category":    "kubernetes",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"namespace": map[string]interface{}{
						"type": "string",
					},
					"name": map[string]interface{}{
						"type": "string",
					},
				},
				"required": []string{"namespace", "name"},
			},
			"version":    "1.0.0",
			"deprecated": false,
		},
	}

	// Enhanced ToolDetail schema with comprehensive information
	g.spec.Components.Schemas["ToolDetail"] = map[string]interface{}{
		"type": "object",
		"allOf": []interface{}{
			map[string]interface{}{
				"$ref": "#/components/schemas/Tool",
			},
		},
		"properties": map[string]interface{}{
			"examples": map[string]interface{}{
				"type":        "array",
				"description": "Usage examples for this tool",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"description": map[string]interface{}{
							"type": "string",
						},
						"arguments": map[string]interface{}{
							"type": "object",
						},
					},
				},
			},
			"performanceHints": map[string]interface{}{
				"type":        "object",
				"description": "Performance-related information about the tool",
				"properties": map[string]interface{}{
					"averageExecutionTime": map[string]interface{}{
						"type": "string",
					},
					"cacheability": map[string]interface{}{
						"type": "boolean",
					},
					"recommendations": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
				},
			},
			"documentation": map[string]interface{}{
				"type":        "string",
				"description": "Detailed markdown documentation for the tool",
			},
			"relatedTools": map[string]interface{}{
				"type":        "array",
				"description": "Related tool names",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
		},
		"example": map[string]interface{}{
			"name":        "kubernetes_get_pod",
			"description": "Get details of a specific pod in a namespace",
			"service":     "kubernetes",
			"category":    "kubernetes",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"namespace": map[string]interface{}{
						"type": "string",
					},
					"name": map[string]interface{}{
						"type": "string",
					},
				},
				"required": []string{"namespace", "name"},
			},
			"examples": []interface{}{
				map[string]interface{}{
					"description": "Get a pod named 'my-pod' in namespace 'default'",
					"arguments": map[string]interface{}{
						"namespace": "default",
						"name":      "my-pod",
					},
				},
			},
			"performanceHints": map[string]interface{}{
				"averageExecutionTime": "150ms",
				"cacheability":         true,
				"recommendations": []string{
					"Use with appropriate timeouts for long-running operations",
				},
			},
			"version":    "1.0.0",
			"deprecated": false,
		},
	}

	// Enhanced ToolCall schema with timeout support
	g.spec.Components.Schemas["ToolCall"] = map[string]interface{}{
		"type":     "object",
		"required": []string{"name", "arguments"},
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Name of the tool to execute",
			},
			"arguments": map[string]interface{}{
				"type":        "object",
				"description": "Arguments to pass to the tool, matching the tool's inputSchema",
			},
			"timeout": map[string]interface{}{
				"type":        "integer",
				"description": "Execution timeout in milliseconds",
				"minimum":     1,
				"maximum":     300000,
			},
		},
		"example": map[string]interface{}{
			"name": "kubernetes_get_pod",
			"arguments": map[string]interface{}{
				"namespace": "default",
				"name":      "my-pod",
			},
			"timeout": 30000,
		},
	}

	// Enhanced ToolResult schema with detailed information
	g.spec.Components.Schemas["ToolResult"] = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"content": map[string]interface{}{
				"type":        "array",
				"description": "Array of content items returned by the tool",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"type": map[string]interface{}{
							"type": "string",
							"enum": []string{"text", "error", "image"},
						},
						"text": map[string]interface{}{
							"type": "string",
						},
						"data": map[string]interface{}{
							"type": "string",
						},
					},
				},
			},
			"isError": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether the tool execution resulted in an error",
			},
			"executionTime": map[string]interface{}{
				"type":        "integer",
				"description": "Execution time in milliseconds",
			},
			"cached": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether the result was served from cache",
			},
		},
		"example": map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "{\n  \"metadata\": {\n    \"name\": \"my-pod\",\n    \"namespace\": \"default\"\n  },\n  \"spec\": {\n    \"containers\": [\n      {\n        \"name\": \"my-container\",\n        \"image\": \"nginx:latest\"\n      }\n    ]\n  }\n}",
				},
			},
			"isError":       false,
			"executionTime": 150,
			"cached":        false,
		},
	}

	// Error schema for consistent error responses
	g.spec.Components.Schemas["Error"] = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"code": map[string]interface{}{
				"type":        "integer",
				"description": "HTTP status code",
			},
			"message": map[string]interface{}{
				"type":        "string",
				"description": "Human-readable error message",
			},
			"details": map[string]interface{}{
				"type":        "object",
				"description": "Additional error details",
			},
		},
		"required": []string{"code", "message"},
		"example": map[string]interface{}{
			"code":    500,
			"message": "Internal server error",
			"details": map[string]interface{}{
				"error": "Connection timeout",
			},
		},
	}
}

// addSecuritySchemes adds security schemes
func (g *Generator) addSecuritySchemes() {
	g.spec.Components.SecuritySchemes = map[string]SecurityScheme{
		"bearerAuth": {
			Type:         "http",
			Description:  "Bearer token authentication using JWT",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
		"apiKey": {
			Type:        "apiKey",
			Description: "API key authentication",
			Name:        "X-API-Key",
			In:          "header",
		},
	}
}

// addModeInformation adds information about different communication modes
func (g *Generator) addModeInformation() {
	// Add a tag for mode-specific information
	g.spec.Tags = append(g.spec.Tags, Tag{
		Name:        "Communication Modes",
		Description: "Different communication modes supported by the server: SSE, HTTP, and stdio",
	})

	// Add mode information to the info description
	modeInfo := map[string]interface{}{
		"sse": map[string]interface{}{
			"description": "Server-Sent Events mode for real-time streaming communication",
			"endpoints":   []string{"/api/*/sse", "/api/*/sse/message"},
		},
		"http": map[string]interface{}{
			"description": "Traditional HTTP REST API mode",
			"endpoints":   []string{"/tools", "/tools/call", "/tools/{toolName}"},
		},
		"stdio": map[string]interface{}{
			"description": "Standard input/output mode for direct command-line interaction",
			"usage":       "Use --mode=stdio flag when starting the server",
		},
	}

	// Add mode information as an extension to the spec
	if g.spec.Components.Schemas == nil {
		g.spec.Components.Schemas = make(map[string]interface{})
	}
	g.spec.Components.Schemas["CommunicationModes"] = modeInfo
}

// SaveToFile saves the OpenAPI specification to a file
func (g *Generator) SaveToFile(filename string) error {
	spec, err := g.Generate()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

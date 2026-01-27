// Package handlers provides HTTP handlers for Kibana MCP operations.
// This file contains common utility functions used across all handlers.
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/kibana/client"
	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
)

// marshalIndentJSON performs indented JSON encoding using object pool
func marshalIndentJSON(data interface{}) ([]byte, error) {
	// First encode to compact format using object pool
	compactBytes, err := optimize.GlobalJSONPool.MarshalToBytes(data)
	if err != nil {
		return nil, err
	}

	// For scenarios requiring indented display, still use standard library but reduce allocations
	// This is a trade-off between performance and readability
	var result bytes.Buffer
	err = json.Indent(&result, compactBytes, "", "  ")
	return result.Bytes(), err
}

// parseLimitWithWarnings validates and parses limit parameter with warnings
func parseLimitWithWarnings(request mcp.CallToolRequest, toolName string) int {
	limit := 20
	if v, ok := request.GetArguments()["limit"]; ok {
		if f, ok := v.(float64); ok {
			limit = int(f)
			if limit <= 0 {
				limit = 20
			} else if limit > 100 {
				logrus.WithField("requested", limit).WithField("max", 100).Warn("Limit too high, resetting to safe maximum")
				limit = 100
			}
		}
	}

	if limit > 50 {
		logrus.WithFields(logrus.Fields{
			"tool":  toolName,
			"limit": limit,
		}).Warn("Large limit may cause context overflow, consider using pagination")
	}

	return limit
}

// getOptionalIntParam gets optional numeric parameter
func getOptionalIntParam(request mcp.CallToolRequest, param string, defaultValue int) int {
	if v, ok := request.GetArguments()[param]; ok {
		if f, ok := v.(float64); ok {
			val := int(f)
			if val > 0 {
				return val
			}
		} else if s, ok := v.(string); ok {
			if val, err := strconv.Atoi(s); err == nil && val > 0 {
				return val
			}
		}
	}
	return defaultValue
}

// getOptionalBoolParam gets optional boolean parameter
func getOptionalBoolParam(request mcp.CallToolRequest, param string) *bool {
	if value, ok := request.GetArguments()[param].(bool); ok {
		return &value
	}
	return nil
}

// getOptionalStringParam gets optional string parameter
func getOptionalStringParam(request mcp.CallToolRequest, param string) string {
	value, _ := request.GetArguments()[param].(string)
	return value
}

// marshalOptimizedResponse marshals optimized response with size warning
func marshalOptimizedResponse(data any, toolName string) (*mcp.CallToolResult, error) {
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize response: %w", err)
	}

	// Add size warning for large responses
	if len(jsonResponse) > 100*1024 { // 100KB
		logrus.WithFields(logrus.Fields{
			"tool":      toolName,
			"sizeBytes": len(jsonResponse),
			"sizeKB":    len(jsonResponse) / 1024,
		}).Warn("Large response generated")
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(jsonResponse)),
		},
	}, nil
}

// requireStringParam validates required string parameter
func requireStringParam(request mcp.CallToolRequest, param string) (string, error) {
	value, ok := request.GetArguments()[param].(string)
	if !ok || value == "" {
		return "", fmt.Errorf("missing required parameter: %s", param)
	}
	return value, nil
}

// getStringFieldFromMap gets string field from map
func getStringFieldFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

// HandleTestConnection handles Kibana connection test requests.
func HandleTestConnection(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Testing Kibana connection")

		// Test connection
		err := c.TestConnection(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Connection test failed: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Connection test successful"),
			},
		}, nil
	}
}

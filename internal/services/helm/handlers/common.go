// Package handlers provides MCP tool handlers for the Helm service.
// It implements handlers for managing Helm releases, charts, repositories, and their integration with other services.
package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
	"github.com/mark3labs/mcp-go/mcp"
)

var (
	ErrMissingRequiredParam = errors.New("missing required parameter")
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

// Helper function to validate required string parameter
func requireStringParam(request mcp.CallToolRequest, param string) (string, error) {
	value, ok := request.GetArguments()[param].(string)
	if !ok || value == "" {
		return "", fmt.Errorf("%w: %s", ErrMissingRequiredParam, param)
	}
	return value, nil
}

// Helper function to get optional string parameter
func getOptionalStringParam(request mcp.CallToolRequest, param string) string {
	value, _ := request.GetArguments()[param].(string)
	return value
}

// Helper function to get optional bool parameter
func getOptionalBoolParam(request mcp.CallToolRequest, param string) bool {
	value, _ := request.GetArguments()[param].(bool)
	return value
}

// Helper function to get optional int parameter
func getOptionalIntParam(request mcp.CallToolRequest, param string) int {
	value, ok := request.GetArguments()[param].(float64)
	if !ok {
		return 0
	}
	return int(value)
}

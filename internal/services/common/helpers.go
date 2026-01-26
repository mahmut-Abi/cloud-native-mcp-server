package common

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	optimize "github.com/mahmut-Abi/k8s-mcp-server/internal/util/performance"
)

// MarshalJSON outputs JSON response using pooled encoder
func MarshalJSON(data interface{}) (*mcp.CallToolResult, error) {
	if data == nil {
		return mcp.NewToolResultText("{}"), nil
	}

	// If already a string, return directly
	if str, ok := data.(string); ok {
		return mcp.NewToolResultText(str), nil
	}

	// Serialize to JSON using pool
	jsonData, err := optimize.GlobalJSONPool.MarshalToBytes(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return mcp.NewToolResultText(string(jsonData)), nil
}

// TextResponse returns plain text response
func TextResponse(text string) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText(text), nil
}

// ErrorResponse returns error response
func ErrorResponse(err error) (*mcp.CallToolResult, error) {
	return nil, err
}

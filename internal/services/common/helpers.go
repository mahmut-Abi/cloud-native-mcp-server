package common

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/sanitize"
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

// GetSanitizedStringArg safely extracts and sanitizes a string argument from request
func GetSanitizedStringArg(args map[string]interface{}, key string) (string, bool) {
	if value, ok := args[key].(string); ok && value != "" {
		return sanitize.SanitizeFilterValue(value), true
	}
	return "", false
}

// GetSanitizedStringArgWithDefault safely extracts and sanitizes a string argument with default value
func GetSanitizedStringArgWithDefault(args map[string]interface{}, key, defaultValue string) string {
	if value, ok := args[key].(string); ok && value != "" {
		return sanitize.SanitizeFilterValue(value)
	}
	return defaultValue
}

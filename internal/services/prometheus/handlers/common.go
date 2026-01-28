// Package handlers provides HTTP handlers for Prometheus MCP operations.
// It implements request handling for Prometheus queries, metrics, targets, and alerts.
package handlers

import (
	"bytes"
	"encoding/json"

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

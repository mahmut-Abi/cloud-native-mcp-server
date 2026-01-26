package handlers

import (
	"fmt"
	"strings"

	optimize "github.com/mahmut-Abi/k8s-mcp-server/internal/util/performance"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
)

const (
	// Response size limits (bytes)
	MaxResponseSize      = 1024 * 1024 // 1MB
	WarningResponseSize  = 512 * 1024  // 512KB
	TruncateResponseSize = 256 * 1024  // 256KB
)

// ResponseSizeMonitor monitors and limits response size
type ResponseSizeMonitor struct {
	OriginalSize     int64    `json:"originalSize"`
	FinalSize        int64    `json:"finalSize"`
	WasTruncated     bool     `json:"wasTruncated"`
	TruncationReason string   `json:"truncationReason,omitempty"`
	ContentEncoding  string   `json:"contentEncoding"`
	EstimatedTokens  int      `json:"estimatedTokens"`
	Warnings         []string `json:"warnings,omitempty"`
}

// LimitResponseSize limits response size and provides truncation information
func LimitResponseSize(data interface{}, toolName string) (interface{}, *ResponseSizeMonitor) {
	monitor := &ResponseSizeMonitor{
		ContentEncoding: "json",
		Warnings:        []string{},
	}

	// Serialize data to calculate size
	jsonData, err := optimize.GlobalJSONPool.MarshalToBytes(data)
	if err != nil {
		monitor.OriginalSize = 0
		monitor.Warnings = append(monitor.Warnings, fmt.Sprintf("Failed to serialize response: %v", err))
		return data, monitor
	}

	monitor.OriginalSize = int64(len(jsonData))
	monitor.EstimatedTokens = estimateTokens(string(jsonData))

	// If response size is within safe range, return directly
	if monitor.OriginalSize <= WarningResponseSize {
		monitor.FinalSize = monitor.OriginalSize
		return data, monitor
	}

	// If response is too large, truncate
	if monitor.OriginalSize > MaxResponseSize {
		monitor.WasTruncated = true
		monitor.TruncationReason = "Response exceeded maximum size limit"
		truncated := truncateResponse(data, TruncateResponseSize, toolName)

		// Calculate size after truncation
		truncatedJson, _ := optimize.GlobalJSONPool.MarshalToBytes(truncated)
		monitor.FinalSize = int64(len(truncatedJson))

		return truncated, monitor
	}

	// Response is large but not exceeding maximum limit, add warning
	monitor.Warnings = append(monitor.Warnings,
		fmt.Sprintf("Large response (%d bytes, ~%d tokens) may impact performance",
			monitor.OriginalSize, monitor.EstimatedTokens))
	monitor.FinalSize = monitor.OriginalSize

	return data, monitor
}

// truncateResponse truncates response to specified size
func truncateResponse(data interface{}, maxSize int, toolName string) interface{} {
	// For array types, truncate array elements
	if slice, ok := data.([]interface{}); ok {
		if len(slice) == 0 {
			return data
		}

		truncated := []interface{}{}
		currentSize := 0

		for i, item := range slice {
			itemBytes, _ := optimize.GlobalJSONPool.MarshalToBytes(item)
			itemSize := len(itemBytes)

			// Estimate metadata overhead
			metadataOverhead := 200 // JSON structure overhead

			if currentSize+itemSize+metadataOverhead > maxSize {
				logrus.WithFields(logrus.Fields{
					"tool":      toolName,
					"original":  len(slice),
					"truncated": i,
					"maxSize":   maxSize,
				}).Info("Truncated array response")

				// Add truncation information
				truncated = append(truncated, map[string]interface{}{
					"_truncated":      true,
					"_info":           fmt.Sprintf("Response truncated from %d to %d items to fit size limit", len(slice), i),
					"_original_count": len(slice),
					"_truncated_at":   i,
				})
				break
			}

			truncated = append(truncated, item)
			currentSize += itemSize + metadataOverhead
		}

		return truncated
	}

	// For map types, try to remove non-essential fields
	if m, ok := data.(map[string]interface{}); ok {
		// For large items arrays, first try to truncate items
		if items, exists := m["items"]; exists {
			if itemsSlice, ok := items.([]interface{}); ok {
				truncatedItems := truncateResponse(itemsSlice, maxSize-1024, toolName) // Leave space for other fields
				m["items"] = truncatedItems

				// Check total size
				totalBytes, _ := optimize.GlobalJSONPool.MarshalToBytes(m)
				if len(totalBytes) <= maxSize {
					return m
				}
			}
		}

		// If still too large, remove some optional fields
		optionalFields := []string{"metadata", "status", "events", "podTemplates", "spec"}
		for _, field := range optionalFields {
			if _, exists := m[field]; exists {
				delete(m, field)

				totalBytes, _ := optimize.GlobalJSONPool.MarshalToBytes(m)
				if len(totalBytes) <= maxSize {
					logrus.WithFields(logrus.Fields{
						"tool":      toolName,
						"removed":   field,
						"finalSize": len(totalBytes),
					}).Info("Removed field to reduce response size")

					// Add removed field information
					m["_truncated"] = true
					m["_info"] = fmt.Sprintf("Removed field '%s' to fit size limit", field)
					return m
				}
			}
		}

		// If still too large, return basic error information
		return map[string]interface{}{
			"_error":        true,
			"_truncated":    true,
			"_info":         fmt.Sprintf("Response too large (%d bytes), content removed", len(m)),
			"_originalSize": len(m),
			"tool":          toolName,
		}
	}

	// For other types, return error information
	return map[string]interface{}{
		"_error":     true,
		"_truncated": true,
		"_info":      "Response too large, content not available",
		"tool":       toolName,
	}
}

// estimateTokens estimates the token count of text
func estimateTokens(text string) int {
	// Simple token estimation: average 4 characters per token (including spaces)
	// For English technical documentation, this is a reasonable estimate
	return len(text) / 4
}

// AddSizeInfoToResponse adds size information to response
func AddSizeInfoToResponse(response map[string]interface{}, monitor *ResponseSizeMonitor) {
	if response == nil {
		response = make(map[string]interface{})
	}

	// Add size metadata
	response["_sizeInfo"] = map[string]interface{}{
		"originalBytes":   monitor.OriginalSize,
		"finalBytes":      monitor.FinalSize,
		"estimatedTokens": monitor.EstimatedTokens,
		"wasTruncated":    monitor.WasTruncated,
	}

	// Add truncation information
	if monitor.WasTruncated {
		response["_sizeInfo"].(map[string]interface{})["truncationReason"] = monitor.TruncationReason
	}

	// Add warnings
	if len(monitor.Warnings) > 0 {
		response["_warnings"] = monitor.Warnings
	}
}

// SanitizeForLLM sanitizes response content to optimize LLM processing
func SanitizeForLLM(data interface{}, toolName string) interface{} {
	// Remove null values and empty strings
	cleanData := removeNullsAndEmpties(data)

	// Limit nesting depth
	cleanData = limitNestingDepth(cleanData, 5)

	// Check size
	monitor := &ResponseSizeMonitor{}
	jsonData, _ := optimize.GlobalJSONPool.MarshalToBytes(cleanData)
	monitor.OriginalSize = int64(len(jsonData))
	monitor.EstimatedTokens = estimateTokens(string(jsonData))

	if monitor.OriginalSize > WarningResponseSize {
		logrus.WithFields(logrus.Fields{
			"tool":   toolName,
			"bytes":  monitor.OriginalSize,
			"tokens": monitor.EstimatedTokens,
		}).Warn("Response size is large, consider using summary tools or pagination")
	}

	return cleanData
}

// removeNullsAndEmpties recursively removes null values and empty strings
func removeNullsAndEmpties(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		cleaned := make(map[string]interface{})
		for key, value := range v {
			// Skip metadata starting with underscore
			if strings.HasPrefix(key, "_") {
				continue
			}

			cleanedValue := removeNullsAndEmpties(value)
			if cleanedValue != nil && cleanedValue != "" {
				cleaned[key] = cleanedValue
			}
		}
		return cleaned

	case []interface{}:
		cleaned := make([]interface{}, 0, len(v))
		for _, item := range v {
			cleanedItem := removeNullsAndEmpties(item)
			if cleanedItem != nil {
				cleaned = append(cleaned, cleanedItem)
			}
		}
		return cleaned

	default:
		return v
	}
}

// limitNestingDepth limits nesting depth
func limitNestingDepth(data interface{}, maxDepth int) interface{} {
	if maxDepth <= 0 {
		return "[max depth reached]"
	}

	switch v := data.(type) {
	case map[string]interface{}:
		cleaned := make(map[string]interface{})
		for key, value := range v {
			cleaned[key] = limitNestingDepth(value, maxDepth-1)
		}
		return cleaned

	case []interface{}:
		cleaned := make([]interface{}, len(v))
		for i, item := range v {
			cleaned[i] = limitNestingDepth(item, maxDepth-1)
		}
		return cleaned

	default:
		return v
	}
}

// FormatResponseForLLM formats response to optimize LLM processing
func FormatResponseForLLM(data interface{}, toolName string) *mcp.CallToolResult {
	// First sanitize the data
	cleanData := SanitizeForLLM(data, toolName)

	// Check and limit size
	finalData, monitor := LimitResponseSize(cleanData, toolName)

	// If map, add size information
	if resultMap, ok := finalData.(map[string]interface{}); ok {
		AddSizeInfoToResponse(resultMap, monitor)
	}

	// Encode as JSON
	jsonResponse, err := optimize.GlobalJSONPool.MarshalToBytes(finalData)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf(`{"code": 1, "data": null, "message": "Failed to serialize response: %s"}`, err.Error()),
				},
			},
			IsError: false,
		}
	}

	// Log warning information
	if len(monitor.Warnings) > 0 {
		logrus.WithFields(logrus.Fields{
			"tool":     toolName,
			"warnings": monitor.Warnings,
		}).Warn("Response processing warnings")
	}

	if monitor.WasTruncated {
		logrus.WithFields(logrus.Fields{
			"tool":             toolName,
			"originalSize":     monitor.OriginalSize,
			"finalSize":        monitor.FinalSize,
			"truncationReason": monitor.TruncationReason,
		}).Info("Response was truncated")
	}

	return mcp.NewToolResultText(string(jsonResponse))
}

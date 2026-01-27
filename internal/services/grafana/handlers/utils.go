package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
)

// Response size limits for Grafana - same as Kubernetes but tuned for Grafana content
const (
	MaxResponseSize      = 1024 * 1024 // 1MB maximum
	WarningResponseSize  = 512 * 1024  // 512KB warning threshold
	TruncateResponseSize = 256 * 1024  // 256KB truncation threshold
	DashboardWarningSize = 512 * 1024  // 512KB for dashboard content specifically
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

// ResponseSizeMonitor tracks response size and operations
type ResponseSizeMonitor struct {
	OriginalSize     int64
	FinalSize        int64
	WasTruncated     bool
	TruncationReason string
	EstimatedTokens  int
	Warnings         []string
}

// FormatResponseForLLM optimizes response for LLM context management
func FormatResponseForLLM(data any, toolName string) (*mcp.CallToolResult, error) {
	monitor := &ResponseSizeMonitor{
		Warnings: make([]string, 0),
	}

	// Serialize original data to check size
	originalBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize original response: %w", err)
	}

	monitor.OriginalSize = int64(len(originalBytes))
	monitor.EstimatedTokens = estimateTokens(originalBytes)

	// Check if we need to truncate
	processedData := data
	if monitor.OriginalSize > TruncateResponseSize {
		logrus.WithFields(logrus.Fields{
			"tool":    toolName,
			"size":    monitor.OriginalSize,
			"sizeKB":  monitor.OriginalSize / 1024,
			"warning": "Response size exceeds safe limit",
		}).Warn("Large response detected, applying smart truncation")

		processedData, monitor = applySmartTruncation(data, toolName, monitor)
		processedBytes, _ := json.Marshal(processedData)
		monitor.FinalSize = int64(len(processedBytes))
		monitor.WasTruncated = true
		monitor.TruncationReason = "Context overflow prevention"
	} else {
		monitor.FinalSize = monitor.OriginalSize
	}

	// Add warning if still large
	if monitor.FinalSize > WarningResponseSize {
		warning := fmt.Sprintf("Large response (%.1fKB) may cause context overflow", float64(monitor.FinalSize)/1024)
		monitor.Warnings = append(monitor.Warnings, warning)
	}

	// Construct final response with metadata
	finalResponse := map[string]any{
		"data": processedData,
	}

	// Add size and truncation metadata
	metadata := map[string]any{
		"originalSizeBytes": monitor.OriginalSize,
		"finalSizeBytes":    monitor.FinalSize,
		"estimatedTokens":   monitor.EstimatedTokens,
		"sizeReductionPct":  calculateSizeReduction(monitor.OriginalSize, monitor.FinalSize),
	}

	if monitor.WasTruncated {
		metadata["truncated"] = true
		metadata["truncationReason"] = monitor.TruncationReason
	}

	if len(monitor.Warnings) > 0 {
		metadata["warnings"] = monitor.Warnings
	}

	// Add metadata only if we have useful information
	if monitor.WasTruncated || len(monitor.Warnings) > 0 || toolName == "grafana_dashboard" {
		finalResponse["metadata"] = metadata
	}

	// Convert to final response
	jsonResponse, err := json.Marshal(finalResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize final response: %w", err)
	}

	// Log the optimization
	logrus.WithFields(logrus.Fields{
		"tool":             toolName,
		"originalSize":     monitor.OriginalSize,
		"finalSize":        monitor.FinalSize,
		"sizeReductionPct": metadata["sizeReductionPct"],
		"wasTruncated":     monitor.WasTruncated,
		"warnings":         len(monitor.Warnings),
	}).Debug("Response optimization completed")

	return mcp.NewToolResultText(string(jsonResponse)), nil
}

// applySmartTruncation applies context-aware truncation based on tool and data type
func applySmartTruncation(data any, toolName string, monitor *ResponseSizeMonitor) (any, *ResponseSizeMonitor) {
	// Convert to map for processing
	dataMap, ok := data.(map[string]any)
	if !ok {
		// If not a map, create a simple truncation
		originalStr := fmt.Sprintf("%v", data)
		if len(originalStr) > int(TruncateResponseSize) {
			truncated := originalStr[:int(TruncateResponseSize-100)] + "...[truncated due to size]"
			monitor.TruncationReason = "String data too large"
			return truncated, monitor
		}
		return data, monitor
	}

	// Apply tool-specific truncation
	switch toolName {
	case "grafana_dashboards":
		return truncateDashboardList(dataMap, monitor)
	case "grafana_dashboard":
		return truncateSingleDashboard(dataMap, monitor)
	case "grafana_datasources":
		return truncateDatasourceList(dataMap, monitor)
	case "grafana_datasource_detail":
		return truncateSingleDatasource(dataMap, monitor)
	case "grafana_folders":
		return truncateFolderList(dataMap, monitor)
	case "grafana_alerts":
		return truncateAlertList(dataMap, monitor)
	default:
		return truncateGeneric(dataMap, monitor)
	}
}

// truncateDashboardList truncates dashboard list responses
func truncateDashboardList(dataMap map[string]any, monitor *ResponseSizeMonitor) (any, *ResponseSizeMonitor) {
	if dashboards, exists := dataMap["dashboards"]; exists {
		if dashboardsList, ok := dashboards.([]any); ok {
			// Keep fewer dashboards if list is too long
			maxDashboards := 10 // Reduce from potentially 20-50
			if len(dashboardsList) > maxDashboards {
				truncatedDashboards := dashboardsList[:maxDashboards]
				dataMap["dashboards"] = truncatedDashboards
				dataMap["count"] = len(truncatedDashboards)

				// Add truncation metadata
				if metadata, exists := dataMap["metadata"]; exists {
					if metaMap, ok := metadata.(map[string]any); ok {
						metaMap["dashboardListTruncated"] = true
						metaMap["originalCount"] = len(dashboardsList)
						metaMap["truncationReason"] = "Dashboard list too long for context"
					}
				}
				monitor.TruncationReason = "Dashboard list truncated to prevent overflow"
			}
		}
	}
	return dataMap, monitor
}

// truncateSingleDashboard handles single dashboard truncation
func truncateSingleDashboard(dataMap map[string]any, monitor *ResponseSizeMonitor) (any, *ResponseSizeMonitor) {
	if dashboard, exists := dataMap["dashboard"]; exists {
		if dashboardMap, ok := dashboard.(map[string]any); ok {
			// Check dashboard size specifically
			if dashboardBytes, err := json.Marshal(dashboardMap); err == nil {
				if int64(len(dashboardBytes)) > DashboardWarningSize {
					// Very large dashboard - remove heavy panels
					if panelsData, exists := dashboardMap["panels"]; exists {
						if panelsList, ok := panelsData.([]any); ok {
							// Keep only first 5 panels or panels that indicate issues
							keepPanels := make([]any, 0)
							maxPanels := 5
							for i, panel := range panelsList {
								if i >= maxPanels {
									break
								}
								if panelMap, ok := panel.(map[string]any); ok {
									// Keep essential panel info only
									simplifiedPanel := map[string]any{
										"id":      panelMap["id"],
										"title":   panelMap["title"],
										"type":    panelMap["type"],
										"gridPos": panelMap["gridPos"],
									}
									keepPanels = append(keepPanels, simplifiedPanel)
								}
							}
							dashboardMap["panels"] = keepPanels

							monitor.TruncationReason = fmt.Sprintf("Dashboard too large (%.1fKB), panels truncated",
								float64(len(dashboardBytes))/1024)
						}
					}

					// Remove other heavy fields
					delete(dashboardMap, "templating")  // Variables can be large
					delete(dashboardMap, "annotations") // Annotations can be large
					delete(dashboardMap, "revision")    // Revision history
				}
			}
		}
	}
	return dataMap, monitor
}

// truncateDatasourceList truncates datasource list responses
func truncateDatasourceList(dataMap map[string]any, monitor *ResponseSizeMonitor) (any, *ResponseSizeMonitor) {
	if datasources, exists := dataMap["datasources"]; exists {
		if datasourcesList, ok := datasources.([]any); ok {
			// Keep fewer datasources if list is too long
			maxDatasources := 10
			if len(datasourcesList) > maxDatasources {
				truncatedDatasources := datasourcesList[:maxDatasources]
				dataMap["datasources"] = truncatedDatasources
				dataMap["count"] = len(truncatedDatasources)

				monitor.TruncationReason = "Datasource list truncated to prevent overflow"
			}
		}
	}
	return dataMap, monitor
}

// truncateSingleDatasource truncates single datasource responses
func truncateSingleDatasource(dataMap map[string]any, monitor *ResponseSizeMonitor) (any, *ResponseSizeMonitor) {
	// For datasource details, remove sensitive and large configuration
	if datasource, exists := dataMap["datasource"]; exists {
		if datasourceMap, ok := datasource.(map[string]any); ok {
			// Remove sensitive and large fields
			delete(datasourceMap, "secureJsonData")
			delete(datasourceMap, "jsonData") // Can be very large
			delete(datasourceMap, "user")
			delete(datasourceMap, "password")
			delete(datasourceMap, "database") // May contain sensitive info

			monitor.TruncationReason = "Removed sensitive and large configuration data"
		}
	}
	return dataMap, monitor
}

// truncateFolderList truncates folder list responses
func truncateFolderList(dataMap map[string]any, monitor *ResponseSizeMonitor) (any, *ResponseSizeMonitor) {
	if folders, exists := dataMap["folders"]; exists {
		if foldersList, ok := folders.([]any); ok {
			// Folders are usually smaller, but still apply limit
			maxFolders := 20
			if len(foldersList) > maxFolders {
				truncatedFolders := foldersList[:maxFolders]
				dataMap["folders"] = truncatedFolders
				dataMap["count"] = len(truncatedFolders)

				monitor.TruncationReason = "Folder list truncated"
			}
		}
	}
	return dataMap, monitor
}

// truncateAlertList truncates alert list responses
func truncateAlertList(dataMap map[string]any, monitor *ResponseSizeMonitor) (any, *ResponseSizeMonitor) {
	if alerts, exists := dataMap["alertRules"]; exists {
		if alertsList, ok := alerts.([]any); ok {
			// Alerts can be complex, keep fewer
			maxAlerts := 10
			if len(alertsList) > maxAlerts {
				truncatedAlerts := alertsList[:maxAlerts]
				dataMap["alertRules"] = truncatedAlerts
				dataMap["count"] = len(truncatedAlerts)

				monitor.TruncationReason = "Alert rules list truncated"
			}
		}
	}
	return dataMap, monitor
}

// truncateGeneric applies generic truncation to any response
func truncateGeneric(dataMap map[string]any, monitor *ResponseSizeMonitor) (any, *ResponseSizeMonitor) {
	// Remove null values and trim strings
	cleaned := make(map[string]any)
	for key, value := range dataMap {
		if value == nil {
			continue
		}

		if strValue, ok := value.(string); ok && len(strValue) > 1000 {
			// Trim very long strings
			cleaned[key] = strValue[:1000] + "...[truncated]"
		} else if arrValue, ok := value.([]any); ok && len(arrValue) > 10 {
			// Truncate long arrays
			cleaned[key] = arrValue[:10]
		} else {
			cleaned[key] = value
		}
	}

	monitor.TruncationReason = "Generic truncation applied"
	return cleaned, monitor
}

// estimateTokens estimates token count from character count
func estimateTokens(data []byte) int {
	// Rough estimation: 1 token ≈ 4 characters for English/mixed content
	// Add some buffer for technical content
	return int(float64(len(data)) / 3.5)
}

// calculateSizeReduction calculates percentage size reduction
func calculateSizeReduction(original, final int64) float64 {
	if original == 0 {
		return 0
	}
	reduction := float64(original-final) / float64(original) * 100
	if reduction < 0 {
		return 0
	}
	return reduction
}

// RemoveNullFields removes null and empty fields from response recursively
func RemoveNullFields(data any) any {
	switch v := data.(type) {
	case map[string]any:
		cleaned := make(map[string]any)
		for key, value := range v {
			if value == nil || value == "" {
				continue
			}
			cleaned[key] = RemoveNullFields(value)
		}
		return cleaned
	case []any:
		var cleaned []any
		for _, item := range v {
			if item != nil {
				cleaned = append(cleaned, RemoveNullFields(item))
			}
		}
		return cleaned
	default:
		return v
	}
}

// LimitJSONDepth limits nested JSON depth to prevent overly complex structures
func LimitJSONDepth(data any, maxDepth int) any {
	return limitDepth(data, maxDepth, 0)
}

func limitDepth(data any, maxDepth, currentDepth int) any {
	if currentDepth >= maxDepth {
		return "[depth_limited]"
	}

	switch v := data.(type) {
	case map[string]any:
		cleaned := make(map[string]any)
		for key, value := range v {
			cleaned[key] = limitDepth(value, maxDepth, currentDepth+1)
		}
		return cleaned
	case []any:
		var cleaned []any
		for _, item := range v {
			cleaned = append(cleaned, limitDepth(item, maxDepth, currentDepth+1))
		}
		return cleaned
	default:
		return v
	}
}

// PrettyPrintJSON creates a readable JSON string with size limits
func PrettyPrintJSON(data any, maxLines int) string {
	jsonBytes, err := marshalIndentJSON(data)
	if err != nil {
		return fmt.Sprintf("Error marshaling JSON: %v", err)
	}

	lines := strings.Split(string(jsonBytes), "\n")
	if len(lines) <= maxLines {
		return string(jsonBytes)
	}

	// Keep first part and truncate
	keepLines := lines[:maxLines-1]
	truncatedLines := append(keepLines, "... [truncated due to size]")
	return strings.Join(truncatedLines, "\n")
}

// SanitizeText removes excessive whitespace and long lines
func SanitizeText(text string) string {
	if text == "" {
		return text
	}

	// Remove excessive whitespace
	spaceRegex := regexp.MustCompile(`\s+`)
	text = spaceRegex.ReplaceAllString(text, " ")

	// Split into lines and process
	lines := strings.Split(text, "\n")
	var processedLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Truncate very long lines
		if len(line) > 500 {
			line = line[:497] + "..."
		}

		processedLines = append(processedLines, line)

		// Limit total lines
		if len(processedLines) >= 100 {
			break
		}
	}

	return strings.Join(processedLines, "\n")
}

// CreateSizeWarning creates a standardized size warning message
func CreateSizeWarning(originalSize, finalSize int64, toolName string) string {
	reduction := calculateSizeReduction(originalSize, finalSize)
	return fmt.Sprintf("📊 %s: %.1fKB → %.1fKB (%.1f%% reduction) to prevent context overflow",
		toolName,
		float64(originalSize)/1024,
		float64(finalSize)/1024,
		reduction)
}

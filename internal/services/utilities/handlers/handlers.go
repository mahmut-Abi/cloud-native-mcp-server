// Package handlers provides MCP tool handlers for the utilities service.
package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/sanitize"
)

// HandleGetTime returns the current time in a specified format.
func HandleGetTime(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	format := "2006-01-02 15:04:05 MST"

	// Get optional format parameter
	if f, ok := request.GetArguments()["format"].(string); ok && f != "" {
		format = f
	}

	logrus.WithFields(logrus.Fields{
		"tool":   "utilities_get_time",
		"format": format,
	}).Debug("Handler invoked")

	now := time.Now()
	currentTime := now.Format(format)

	response := map[string]interface{}{
		"time":   currentTime,
		"format": format,
		"unix":   now.Unix(),
	}

	data, err := optimize.GlobalJSONPool.MarshalToBytes(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return mcp.NewToolResultText(string(data)), nil
}

// HandleGetTimestamp returns the current Unix timestamp.
func HandleGetTimestamp(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	unit := "seconds"
	includeMillis := false

	// Get optional parameters
	if u, ok := request.GetArguments()["unit"].(string); ok && u != "" {
		unit = u
	}
	if m, ok := request.GetArguments()["include_millis"].(bool); ok {
		includeMillis = m
	}

	logrus.WithFields(logrus.Fields{
		"tool":           "utilities_get_timestamp",
		"unit":           unit,
		"include_millis": includeMillis,
	}).Debug("Handler invoked")

	now := time.Now()
	var timestamp interface{}

	switch unit {
	case "milliseconds", "millis", "ms":
		timestamp = now.UnixMilli()
		if includeMillis {
			timestamp = float64(now.UnixNano()) / 1e6
		}
	case "nanoseconds", "nanos", "ns":
		timestamp = now.UnixNano()
	case "minutes":
		timestamp = now.Unix() / 60
	case "hours":
		timestamp = now.Unix() / 3600
	case "days":
		timestamp = now.Unix() / 86400
	default:
		timestamp = now.Unix()
	}

	response := map[string]interface{}{
		"timestamp":      timestamp,
		"unit":           unit,
		"include_millis": includeMillis,
		"iso8601":        now.Format(time.RFC3339),
	}

	data, err := optimize.GlobalJSONPool.MarshalToBytes(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return mcp.NewToolResultText(string(data)), nil
}

// HandleGetDate returns the current date in various formats.
func HandleGetDate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	format := "2006-01-02"
	includeTime := false

	// Get optional parameters
	if f, ok := request.GetArguments()["format"].(string); ok && f != "" {
		format = f
	}
	if t, ok := request.GetArguments()["include_time"].(bool); ok {
		includeTime = t
	}

	logrus.WithFields(logrus.Fields{
		"tool":         "utilities_get_date",
		"format":       format,
		"include_time": includeTime,
	}).Debug("Handler invoked")

	now := time.Now()

	if includeTime {
		format = format + " 15:04:05"
	}

	currentDate := now.Format(format)

	response := map[string]interface{}{
		"date":        currentDate,
		"format":      format,
		"unix":        now.Unix(),
		"year":        now.Year(),
		"month":       int(now.Month()),
		"day":         now.Day(),
		"weekday":     now.Weekday().String(),
		"day_of_year": now.YearDay(),
		"iso8601":     now.Format(time.RFC3339),
	}

	data, err := optimize.GlobalJSONPool.MarshalToBytes(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return mcp.NewToolResultText(string(data)), nil
}

// HandlePause pauses execution for a specified duration.
func HandlePause(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	duration := 5 // seconds

	if d, ok := request.GetArguments()["seconds"].(float64); ok {
		duration = int(d)
	}

	if duration <= 0 {
		duration = 1
	}
	if duration > 300 {
		duration = 300 // Max 5 minutes
	}

	logrus.WithFields(logrus.Fields{
		"tool":    "utilities_pause",
		"seconds": duration,
	}).Debug("Handler invoked")

	// Pause for specified duration
	time.Sleep(time.Duration(duration) * time.Second)

	response := map[string]interface{}{
		"paused":   true,
		"duration": duration,
		"unit":     "seconds",
		"message":  fmt.Sprintf("Paused for %d seconds", duration),
	}

	data, err := optimize.GlobalJSONPool.MarshalToBytes(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return mcp.NewToolResultText(string(data)), nil
}

// HandleSleep pauses execution for a specified duration (alias for pause).
func HandleSleep(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	duration := 5 // seconds
	unit := "seconds"

	if d, ok := request.GetArguments()["duration"].(float64); ok {
		duration = int(d)
	}
	if u, ok := request.GetArguments()["unit"].(string); ok && u != "" {
		unit = u
	}

	logrus.WithFields(logrus.Fields{
		"tool":     "utilities_sleep",
		"duration": duration,
		"unit":     unit,
	}).Debug("Handler invoked")

	// Calculate sleep duration
	var sleepDuration time.Duration
	switch unit {
	case "milliseconds", "millis", "ms":
		sleepDuration = time.Duration(duration) * time.Millisecond
	case "minutes", "minute", "min":
		sleepDuration = time.Duration(duration) * time.Minute
	case "hours", "hour":
		sleepDuration = time.Duration(duration) * time.Hour
	default:
		sleepDuration = time.Duration(duration) * time.Second
	}

	if sleepDuration <= 0 {
		sleepDuration = time.Second
	}
	if sleepDuration > 5*time.Minute {
		sleepDuration = 5 * time.Minute
	}

	time.Sleep(sleepDuration)

	response := map[string]interface{}{
		"slept":    true,
		"duration": duration,
		"unit":     unit,
		"message":  fmt.Sprintf("Slept for %d %s", duration, unit),
	}

	data, err := optimize.GlobalJSONPool.MarshalToBytes(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return mcp.NewToolResultText(string(data)), nil
}

// HandleWebFetch fetches content from a URL.
func HandleWebFetch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	targetURL := ""
	timeout := 30 // seconds
	maxChars := 10000

	// Get required URL parameter
	if u, ok := request.GetArguments()["url"].(string); ok && u != "" {
		// Validate URL format
		parsedURL, err := url.Parse(u)
		if err != nil {
			return nil, fmt.Errorf("invalid URL format: %w", err)
		}
		// Only allow http and https schemes
		if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
			return nil, fmt.Errorf("only http and https URLs are allowed")
		}
		// Sanitize URL to prevent injection
		targetURL = sanitize.SanitizeFilterValue(u)
	} else {
		return nil, fmt.Errorf("url parameter is required")
	}

	// Get optional parameters
	if t, ok := request.GetArguments()["timeout"].(float64); ok {
		timeout = int(t)
	}
	if m, ok := request.GetArguments()["max_chars"].(float64); ok {
		maxChars = int(m)
	}

	logrus.WithFields(logrus.Fields{
		"tool":      "utilities_web_fetch",
		"url":       targetURL,
		"timeout":   timeout,
		"max_chars": maxChars,
	}).Debug("Handler invoked")

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// Make request
	resp, err := client.Get(targetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Truncate if needed
	content := string(body)
	if len(content) > maxChars {
		content = content[:maxChars]
	}

	response := map[string]interface{}{
		"url":            targetURL,
		"status_code":    resp.StatusCode,
		"content_length": len(body),
		"content_type":   resp.Header.Get("Content-Type"),
		"truncated":      len(body) > maxChars,
		"content":        content,
	}

	data, err := optimize.GlobalJSONPool.MarshalToBytes(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return mcp.NewToolResultText(string(data)), nil
}

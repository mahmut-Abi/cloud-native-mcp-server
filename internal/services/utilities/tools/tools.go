// Package tools provides MCP tool definitions for the utilities service.
package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
)

// GetTimeTool returns the current time in a specified format.
func GetTimeTool() mcp.Tool {
	logrus.Debug("Creating GetTimeTool")
	return mcp.NewTool("utilities_get_time",
		mcp.WithDescription("Get the current time in a specified format. Useful for timestamps, logging, and time-based operations. Default format: '2006-01-02 15:04:05 MST'. Common formats: '2006-01-02 15:04:05' (YYYY-MM-DD HH:MM:SS), '15:04:05' (HH:MM:SS), '2006-01-02' (YYYY-MM-DD), 'Mon Jan 2 15:04:05 2006' (full date), '2006-01-02T15:04:05Z07:00' (ISO 8601)."),
		mcp.WithString("format",
			mcp.Description("Go time format string. Default: '2006-01-02 15:04:05 MST'. See https://pkg.go.dev/time#pkg-constants for format options.")),
	)
}

// GetTimestampTool returns the current Unix timestamp.
func GetTimestampTool() mcp.Tool {
	logrus.Debug("Creating GetTimestampTool")
	return mcp.NewTool("utilities_get_timestamp",
		mcp.WithDescription("Get the current Unix timestamp. Useful for API authentication, cache keys, unique identifiers, and time-sensitive operations. Returns timestamp in specified unit with ISO 8601 format for reference. Default unit: seconds."),
		mcp.WithString("unit",
			mcp.Description("Timestamp unit: 'seconds' (default), 'milliseconds', 'nanoseconds', 'minutes', 'hours', 'days'.")),
		mcp.WithBoolean("include_millis",
			mcp.Description("Include milliseconds in timestamp (only applies to seconds unit).")),
	)
}

// GetDateTool returns the current date.
func GetDateTool() mcp.Tool {
	logrus.Debug("Creating GetDateTool")
	return mcp.NewTool("utilities_get_date",
		mcp.WithDescription("Get the current date with optional time. Returns year, month, day, weekday, and day of year. Useful for date calculations, scheduling, and date-based operations. Default format: '2006-01-02'."),
		mcp.WithString("format",
			mcp.Description("Date format string. Common formats: '2006-01-02' (YYYY-MM-DD, default), '01/02/2006' (MM/DD/YYYY), '02/01/2006' (DD/MM/YYYY), 'Jan 2, 2006' (Month DD, YYYY).")),
		mcp.WithBoolean("include_time",
			mcp.Description("Include time in the response (default: false).")),
	)
}

// PauseTool pauses execution for a specified number of seconds.
func PauseTool() mcp.Tool {
	logrus.Debug("Creating PauseTool")
	return mcp.NewTool("utilities_pause",
		mcp.WithDescription("Pause execution for a specified number of seconds. Useful for waiting for external processes, rate limiting, or giving time for background tasks to complete. Maximum pause: 300 seconds (5 minutes). Default: 5 seconds."),
		mcp.WithNumber("seconds",
			mcp.Description("Number of seconds to pause. Range: 1-300. Default: 5.")),
	)
}

// SleepTool pauses execution for a specified duration.
func SleepTool() mcp.Tool {
	logrus.Debug("Creating SleepTool")
	return mcp.NewTool("utilities_sleep",
		mcp.WithDescription("Pause execution for a specified duration. Similar to pause but supports different time units. Useful for flexible waiting periods. Maximum sleep: 5 minutes. Default: 5 seconds."),
		mcp.WithNumber("duration",
			mcp.Description("Duration value. Default: 5.")),
		mcp.WithString("unit",
			mcp.Description("Time unit: 'seconds' (default), 'milliseconds', 'minutes', 'hours'.")),
	)
}

// WebFetchTool fetches content from a URL.
func WebFetchTool() mcp.Tool {
	logrus.Debug("Creating WebFetchTool")
	return mcp.NewTool("utilities_web_fetch",
		mcp.WithDescription("Fetch content from a URL. Useful for retrieving web pages, APIs, or any HTTP-accessible content. Returns status code, content type, and truncated content. Supports GET requests only. Maximum content: 10000 characters (truncated)."),
		mcp.WithString("url", mcp.Required(),
			mcp.Description("URL to fetch (required). Must be a valid HTTP/HTTPS URL.")),
		mcp.WithNumber("timeout",
			mcp.Description("Request timeout in seconds. Range: 1-120. Default: 30.")),
		mcp.WithNumber("max_chars",
			mcp.Description("Maximum characters to return. Range: 100-50000. Default: 10000.")),
	)
}

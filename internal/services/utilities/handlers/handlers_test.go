package handlers

import (
	"context"
	"testing"

	mcp "github.com/mark3labs/mcp-go/mcp"
)

func TestHandleGetTime(t *testing.T) {
	ctx := context.Background()
	
	// Test with default format
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}
	
	result, err := HandleGetTime(ctx, request)
	if err != nil {
		t.Errorf("HandleGetTime() error = %v", err)
		return
	}
	
	if result == nil {
		t.Error("HandleGetTime() should return non-nil result")
	}
}

func TestHandleGetTimeWithCustomFormat(t *testing.T) {
	ctx := context.Background()
	
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"format": "2006-01-02",
			},
		},
	}
	
	result, err := HandleGetTime(ctx, request)
	if err != nil {
		t.Errorf("HandleGetTime() error = %v", err)
		return
	}
	
	if result == nil {
		t.Error("HandleGetTime() should return non-nil result")
	}
}

func TestHandleGetTimestamp(t *testing.T) {
	ctx := context.Background()
	
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"unit": "seconds",
			},
		},
	}
	
	result, err := HandleGetTimestamp(ctx, request)
	if err != nil {
		t.Errorf("HandleGetTimestamp() error = %v", err)
		return
	}
	
	if result == nil {
		t.Error("HandleGetTimestamp() should return non-nil result")
	}
}

func TestHandleGetTimestampWithMillis(t *testing.T) {
	ctx := context.Background()
	
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"unit": "milliseconds",
			},
		},
	}
	
	result, err := HandleGetTimestamp(ctx, request)
	if err != nil {
		t.Errorf("HandleGetTimestamp() error = %v", err)
		return
	}
	
	if result == nil {
		t.Error("HandleGetTimestamp() should return non-nil result")
	}
}

func TestHandleGetDate(t *testing.T) {
	ctx := context.Background()
	
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}
	
	result, err := HandleGetDate(ctx, request)
	if err != nil {
		t.Errorf("HandleGetDate() error = %v", err)
		return
	}
	
	if result == nil {
		t.Error("HandleGetDate() should return non-nil result")
	}
}

func TestHandleGetDateWithTime(t *testing.T) {
	ctx := context.Background()
	
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"include_time": true,
			},
		},
	}
	
	result, err := HandleGetDate(ctx, request)
	if err != nil {
		t.Errorf("HandleGetDate() error = %v", err)
		return
	}
	
	if result == nil {
		t.Error("HandleGetDate() should return non-nil result")
	}
}

func TestHandlePause(t *testing.T) {
	ctx := context.Background()
	
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"seconds": float64(1),
			},
		},
	}
	
	result, err := HandlePause(ctx, request)
	if err != nil {
		t.Errorf("HandlePause() error = %v", err)
		return
	}
	
	if result == nil {
		t.Error("HandlePause() should return non-nil result")
	}
}

func TestHandleSleep(t *testing.T) {
	ctx := context.Background()
	
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"duration": float64(1),
				"unit":     "seconds",
			},
		},
	}
	
	result, err := HandleSleep(ctx, request)
	if err != nil {
		t.Errorf("HandleSleep() error = %v", err)
		return
	}
	
	if result == nil {
		t.Error("HandleSleep() should return non-nil result")
	}
}

func TestHandleWebFetch(t *testing.T) {
	ctx := context.Background()
	
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"url": "https://example.com",
			},
		},
	}
	
	result, err := HandleWebFetch(ctx, request)
	if err != nil {
		t.Errorf("HandleWebFetch() error = %v", err)
		return
	}
	
	if result == nil {
		t.Error("HandleWebFetch() should return non-nil result")
	}
}

func TestHandleWebFetchWithoutURL(t *testing.T) {
	ctx := context.Background()
	
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}
	
	_, err := HandleWebFetch(ctx, request)
	if err == nil {
		t.Error("HandleWebFetch() should return error when URL is missing")
	}
}

func TestHandleWebFetchWithInvalidURL(t *testing.T) {
	ctx := context.Background()
	
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"url": "ftp://example.com",
			},
		},
	}
	
	_, err := HandleWebFetch(ctx, request)
	if err == nil {
		t.Error("HandleWebFetch() should return error for non-HTTP URL")
	}
}
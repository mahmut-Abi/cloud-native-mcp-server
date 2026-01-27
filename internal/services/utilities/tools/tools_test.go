package tools

import (
	"testing"
)

func TestGetTimeTool(t *testing.T) {
	tool := GetTimeTool()
	if tool.Name != "utilities_get_time" {
		t.Errorf("GetTimeTool() name = %s, want 'utilities_get_time'", tool.Name)
	}

	if tool.Description == "" {
		t.Error("GetTimeTool() should have a description")
	}
}

func TestGetTimestampTool(t *testing.T) {
	tool := GetTimestampTool()
	if tool.Name != "utilities_get_timestamp" {
		t.Errorf("GetTimestampTool() name = %s, want 'utilities_get_timestamp'", tool.Name)
	}

	if tool.Description == "" {
		t.Error("GetTimestampTool() should have a description")
	}
}

func TestGetDateTool(t *testing.T) {
	tool := GetDateTool()
	if tool.Name != "utilities_get_date" {
		t.Errorf("GetDateTool() name = %s, want 'utilities_get_date'", tool.Name)
	}

	if tool.Description == "" {
		t.Error("GetDateTool() should have a description")
	}
}

func TestPauseTool(t *testing.T) {
	tool := PauseTool()
	if tool.Name != "utilities_pause" {
		t.Errorf("PauseTool() name = %s, want 'utilities_pause'", tool.Name)
	}

	if tool.Description == "" {
		t.Error("PauseTool() should have a description")
	}
}

func TestSleepTool(t *testing.T) {
	tool := SleepTool()
	if tool.Name != "utilities_sleep" {
		t.Errorf("SleepTool() name = %s, want 'utilities_sleep'", tool.Name)
	}

	if tool.Description == "" {
		t.Error("SleepTool() should have a description")
	}
}

func TestWebFetchTool(t *testing.T) {
	tool := WebFetchTool()
	if tool.Name != "utilities_web_fetch" {
		t.Errorf("WebFetchTool() name = %s, want 'utilities_web_fetch'", tool.Name)
	}

	if tool.Description == "" {
		t.Error("WebFetchTool() should have a description")
	}
}

func TestGetTimeToolInputSchema(t *testing.T) {
	tool := GetTimeTool()

	_ = tool.InputSchema
}

func TestGetTimestampToolInputSchema(t *testing.T) {
	tool := GetTimestampTool()

	_ = tool.InputSchema
}

func TestWebFetchToolInputSchema(t *testing.T) {
	tool := WebFetchTool()

	_ = tool.InputSchema
}

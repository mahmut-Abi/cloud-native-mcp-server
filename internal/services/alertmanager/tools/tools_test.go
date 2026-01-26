package tools

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestGetStatusTool(t *testing.T) {
	tool := GetStatusTool()

	if tool.Name != "alertmanager_get_status" {
		t.Errorf("Expected tool name to be 'alertmanager_get_status', got %s", tool.Name)
	}

	if tool.Description == "" {
		t.Error("Tool description should not be empty")
	}

	if tool.InputSchema.Type != "object" {
		t.Errorf("Expected input schema type to be 'object', got %s", tool.InputSchema.Type)
	}
}

func TestGetAlertsTool(t *testing.T) {
	tool := GetAlertsTool()

	if tool.Name != "alertmanager_get_alerts" {
		t.Errorf("Expected tool name to be 'alertmanager_get_alerts', got %s", tool.Name)
	}

	if tool.Description == "" {
		t.Error("Tool description should not be empty")
	}

	// Check that filters property exists
	if props, ok := tool.InputSchema.Properties["filters"]; ok {
		if propMap, ok := props.(map[string]interface{}); ok {
			if propMap["type"] != "object" {
				t.Errorf("Expected filters type to be 'object', got %v", propMap["type"])
			}
		} else {
			t.Error("filters property should be a map")
		}
	} else {
		t.Error("Expected filters property in input schema")
	}
}

func TestCreateSilenceTool(t *testing.T) {
	tool := CreateSilenceTool()

	if tool.Name != "alertmanager_create_silence" {
		t.Errorf("Expected tool name to be 'alertmanager_create_silence', got %s", tool.Name)
	}

	// Check required fields
	expectedRequired := []string{"matchers", "endsAt", "comment", "createdBy"}
	for _, field := range expectedRequired {
		found := false
		for _, req := range tool.InputSchema.Required {
			if req == field {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected required field '%s' not found", field)
		}
	}

	// Check matchers property structure
	if props, ok := tool.InputSchema.Properties["matchers"]; ok {
		if propMap, ok := props.(map[string]interface{}); ok {
			if propMap["type"] != "array" {
				t.Errorf("Expected matchers type to be 'array', got %v", propMap["type"])
			}
		} else {
			t.Error("matchers property should be a map")
		}
	} else {
		t.Error("Expected matchers property in input schema")
	}
}

func TestDeleteSilenceTool(t *testing.T) {
	tool := DeleteSilenceTool()

	if tool.Name != "alertmanager_delete_silence" {
		t.Errorf("Expected tool name to be 'alertmanager_delete_silence', got %s", tool.Name)
	}

	// Check required fields
	if len(tool.InputSchema.Required) != 1 || tool.InputSchema.Required[0] != "silenceId" {
		t.Errorf("Expected required field 'silenceId', got %v", tool.InputSchema.Required)
	}

	// Check silenceId property
	if props, ok := tool.InputSchema.Properties["silenceId"]; ok {
		if propMap, ok := props.(map[string]interface{}); ok {
			if propMap["type"] != "string" {
				t.Errorf("Expected silenceId type to be 'string', got %v", propMap["type"])
			}
		} else {
			t.Error("silenceId property should be a map")
		}
	} else {
		t.Error("Expected silenceId property in input schema")
	}
}

func TestQueryAlertsTool(t *testing.T) {
	tool := QueryAlertsTool()

	if tool.Name != "alertmanager_query_alerts" {
		t.Errorf("Expected tool name to be 'alertmanager_query_alerts', got %s", tool.Name)
	}

	// Check optional filter properties
	filterProps := []string{"receiver", "silenced", "active", "unprocessed", "inhibited", "filter", "sortBy", "sortOrder"}
	for _, prop := range filterProps {
		if _, ok := tool.InputSchema.Properties[prop]; !ok {
			t.Errorf("Expected property '%s' in input schema", prop)
		}
	}

	// Check sortBy description (since we're using string instead of enum now)
	if props, ok := tool.InputSchema.Properties["sortBy"]; ok {
		if propMap, ok := props.(map[string]interface{}); ok {
			if propMap["type"] != "string" {
				t.Errorf("Expected sortBy type to be 'string', got %v", propMap["type"])
			}
		}
	}
}

func TestAllToolsHaveValidStructure(t *testing.T) {
	tools := []mcp.Tool{
		GetStatusTool(),
		GetAlertsTool(),
		GetAlertGroupsTool(),
		GetSilencesTool(),
		CreateSilenceTool(),
		DeleteSilenceTool(),
		GetReceiversTool(),
		TestReceiverTool(),
		QueryAlertsTool(),
	}

	for _, tool := range tools {
		// Check that all tools have names starting with "alertmanager_"
		if len(tool.Name) < 12 || tool.Name[:12] != "alertmanager" {
			t.Errorf("Tool name '%s' should start with 'alertmanager'", tool.Name)
		}

		// Check that all tools have descriptions
		if tool.Description == "" {
			t.Errorf("Tool '%s' should have a description", tool.Name)
		}

		// Check that all tools have valid input schema
		if tool.InputSchema.Type != "object" {
			t.Errorf("Tool '%s' should have object input schema type", tool.Name)
		}

		// Properties should be a map
		if tool.InputSchema.Properties == nil {
			t.Errorf("Tool '%s' should have properties in input schema", tool.Name)
		}
	}
}

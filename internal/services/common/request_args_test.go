package common

import (
	"testing"
)

func TestLookupArgSupportsSnakeAndCamelCase(t *testing.T) {
	args := map[string]interface{}{
		"alert_type_id": "cpu_threshold",
		"spaceId":       "prod-space",
	}

	if got, ok := GetStringArg(args, "alertTypeId"); !ok || got != "cpu_threshold" {
		t.Fatalf("GetStringArg(alertTypeId) = %q, %v", got, ok)
	}

	if got, ok := GetStringArg(args, "space_id"); !ok || got != "prod-space" {
		t.Fatalf("GetStringArg(space_id) = %q, %v", got, ok)
	}
}

func TestGetBoolArgSupportsBooleanStrings(t *testing.T) {
	args := map[string]interface{}{
		"enabled":          "true",
		"include_panels":   "0",
		"include_ui_state": false,
	}

	value, err := GetBoolArg(args, "enabled")
	if err != nil || value == nil || !*value {
		t.Fatalf("GetBoolArg(enabled) = %v, %v", value, err)
	}

	value, err = GetBoolArg(args, "include_panels")
	if err != nil || value == nil || *value {
		t.Fatalf("GetBoolArg(include_panels) = %v, %v", value, err)
	}

	value, err = GetBoolArg(args, "include_ui_state")
	if err != nil || value == nil || *value {
		t.Fatalf("GetBoolArg(include_ui_state) = %v, %v", value, err)
	}
}

func TestGetObjectArgSupportsJSONString(t *testing.T) {
	args := map[string]interface{}{
		"schedule": `{"interval":"5m"}`,
	}

	value, ok, err := GetObjectArg(args, "schedule")
	if err != nil || !ok {
		t.Fatalf("GetObjectArg(schedule) error = %v, ok = %v", err, ok)
	}
	if value["interval"] != "5m" {
		t.Fatalf("GetObjectArg(schedule).interval = %v", value["interval"])
	}
}

func TestGetObjectSliceArgSupportsJSONString(t *testing.T) {
	args := map[string]interface{}{
		"actions": `[{"group":"default","id":"connector-1"}]`,
	}

	value, ok, err := GetObjectSliceArg(args, "actions")
	if err != nil || !ok {
		t.Fatalf("GetObjectSliceArg(actions) error = %v, ok = %v", err, ok)
	}
	if len(value) != 1 || value[0]["id"] != "connector-1" {
		t.Fatalf("GetObjectSliceArg(actions) = %#v", value)
	}
}

func TestGetStringSliceArgSupportsCSVAndJSON(t *testing.T) {
	csvArgs := map[string]interface{}{
		"tags": "prod,critical",
	}

	value, ok, err := GetStringSliceArg(csvArgs, "tags")
	if err != nil || !ok || len(value) != 2 {
		t.Fatalf("GetStringSliceArg(csv) = %#v, %v, %v", value, ok, err)
	}

	jsonArgs := map[string]interface{}{
		"fields": `["title","updated_at"]`,
	}

	value, ok, err = GetStringSliceArg(jsonArgs, "fields")
	if err != nil || !ok || len(value) != 2 {
		t.Fatalf("GetStringSliceArg(json) = %#v, %v, %v", value, ok, err)
	}
}

func TestGetJSONStringArgSupportsStructuredValues(t *testing.T) {
	args := map[string]interface{}{
		"panelsJSON": []interface{}{
			map[string]interface{}{"id": float64(1)},
		},
	}

	value, ok, err := GetJSONStringArg(args, "panels_json")
	if err != nil || !ok {
		t.Fatalf("GetJSONStringArg(panels_json) error = %v, ok = %v", err, ok)
	}
	if value != `[{"id":1}]` {
		t.Fatalf("GetJSONStringArg(panels_json) = %s", value)
	}
}

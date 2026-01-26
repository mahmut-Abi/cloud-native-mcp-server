package optimize

import (
	"encoding/json"
	"testing"
)

func TestNewJSONEncoderPool(t *testing.T) {
	pool := NewJSONEncoderPool()
	if pool == nil {
		t.Error("Expected pool, got nil")
	}
}

func TestJSONEncoderPoolMarshalToBytes(t *testing.T) {
	pool := NewJSONEncoderPool()
	data := map[string]interface{}{"key": "value", "count": 42}

	bytes, err := pool.MarshalToBytes(data)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(bytes) == 0 {
		t.Error("Expected non-empty bytes")
	}
}

func TestJSONEncoderPoolMarshalToString(t *testing.T) {
	pool := NewJSONEncoderPool()
	data := map[string]interface{}{"key": "value", "count": 42}

	str, err := pool.MarshalToString(data)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(str) == 0 {
		t.Error("Expected non-empty string")
	}

	// Verify it can be unmarshalled
	var result map[string]interface{}
	err = json.Unmarshal([]byte(str), &result)
	if err != nil {
		t.Errorf("Expected valid JSON, got error: %v", err)
	}
}

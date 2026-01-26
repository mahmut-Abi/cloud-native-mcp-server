package optimize

import (
	"testing"
)

func TestNewStringBuilderPool(t *testing.T) {
	pool := NewStringBuilderPool()
	if pool == nil {
		t.Error("Expected pool, got nil")
	}
}

func TestStringBuilderPoolGetAndPut(t *testing.T) {
	pool := NewStringBuilderPool()

	builder := pool.Get()
	if builder == nil {
		t.Error("Expected builder, got nil")
	}

	builder.WriteString("test")
	pool.Put(builder)
}

package optimize

import (
	"strings"
	"sync"
)

// StringBuilderPool provides a pool of reusable string builders to reduce allocations
type StringBuilderPool struct {
	pool sync.Pool
}

// NewStringBuilderPool creates a new string builder pool
func NewStringBuilderPool() *StringBuilderPool {
	return &StringBuilderPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &strings.Builder{}
			},
		},
	}
}

// Get retrieves a builder from the pool
func (p *StringBuilderPool) Get() *strings.Builder {
	return p.pool.Get().(*strings.Builder)
}

// Put returns a builder to the pool after resetting it
func (p *StringBuilderPool) Put(sb *strings.Builder) {
	sb.Reset()
	p.pool.Put(sb)
}

// GlobalStringBuilderPool is a global string builder pool
var GlobalStringBuilderPool = NewStringBuilderPool()

// BuildString builds a string using a pooled builder
func BuildString(fn func(*strings.Builder)) string {
	sb := GlobalStringBuilderPool.Get()
	defer GlobalStringBuilderPool.Put(sb)
	fn(sb)
	return sb.String()
}

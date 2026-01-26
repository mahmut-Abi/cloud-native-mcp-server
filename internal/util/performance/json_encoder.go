package optimize

import (
	"bytes"
	"encoding/json"
	"io"
	"sync"
)

// JSONEncoderPool provides a pool of reusable JSON encoders to reduce allocations
// This is especially useful for high-frequency JSON encoding operations
type JSONEncoderPool struct {
	pool sync.Pool
}

// NewJSONEncoderPool creates a new JSON encoder pool
func NewJSONEncoderPool() *JSONEncoderPool {
	return &JSONEncoderPool{
		pool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

// MarshalToBytes encodes data to JSON bytes using a pooled buffer
// This reduces heap allocations by reusing buffers
func (p *JSONEncoderPool) MarshalToBytes(data interface{}) ([]byte, error) {
	buf := p.pool.Get().(*bytes.Buffer)
	defer func() {
		// Always reset and return buffer to pool
		buf.Reset()
		p.pool.Put(buf)
	}()

	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}

	// Remove trailing newline added by Encoder
	result := buf.Bytes()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		// Avoid allocating if possible by using slicing
		result = bytes.TrimRight(result, "\n")
	}

	// Avoid cloning when possible - only clone if needed
	// Return copy to ensure buffer can be reused
	returnData := make([]byte, len(result))
	copy(returnData, result)
	return returnData, nil
}

// MarshalToString encodes data to JSON string using a pooled buffer
func (p *JSONEncoderPool) MarshalToString(data interface{}) (string, error) {
	bytes, err := p.MarshalToBytes(data)
	return string(bytes), err
}

// GlobalJSONPool is a global JSON encoder pool for convenience
// StreamEncode encodes data directly to a writer for large payloads
func (p *JSONEncoderPool) StreamEncode(w io.Writer, data interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(data)
}

var GlobalJSONPool = NewJSONEncoderPool()

// JSONEncoderFactory provides efficient JSON encoding with minimal allocations
type JSONEncoderFactory struct {
	pool *JSONEncoderPool
}

// NewJSONEncoderFactory creates a new factory
func NewJSONEncoderFactory() *JSONEncoderFactory {
	return &JSONEncoderFactory{
		pool: NewJSONEncoderPool(),
	}
}

// EncodeToBytes encodes with factory's pool
func (f *JSONEncoderFactory) EncodeToBytes(v interface{}) ([]byte, error) {
	return f.pool.MarshalToBytes(v)
}

// EncodeToString encodes with factory's pool
func (f *JSONEncoderFactory) EncodeToString(v interface{}) (string, error) {
	return f.pool.MarshalToString(v)
}

// MarshalJSON is a convenience function using the global pool
func MarshalJSON(data interface{}) ([]byte, error) {
	return GlobalJSONPool.MarshalToBytes(data)
}

// MarshalJSONString is a convenience function for string encoding
func MarshalJSONString(data interface{}) (string, error) {
	return GlobalJSONPool.MarshalToString(data)
}

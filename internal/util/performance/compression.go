package optimize

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"
	"sync"
)

// CompressionLevel defines the compression level
type CompressionLevel int

const (
	// CompressionDefault uses the default compression level
	CompressionDefault CompressionLevel = 0
	// CompressionFast optimizes for speed
	CompressionFast CompressionLevel = 1
	// CompressionBest optimizes for compression ratio
	CompressionBest CompressionLevel = 9
)

// GzipCompressor provides gzip compression utilities
type GzipCompressor struct {
	bufferPool  sync.Pool
	writerPool  sync.Pool
	level       CompressionLevel
	minSizeComp int64 // Minimum size in bytes for compression to be applied
}

// NewGzipCompressor creates a new gzip compressor with pooled resources
func NewGzipCompressor() *GzipCompressor {
	return NewGzipCompressorWithLevel(CompressionDefault)
}

// NewGzipCompressorWithLevel creates a new gzip compressor with specified compression level
func NewGzipCompressorWithLevel(level CompressionLevel) *GzipCompressor {
	gc := &GzipCompressor{
		level:       level,
		minSizeComp: 1024, // Default to 1KB minimum for compression
	}

	// Initialize buffer pool
	gc.bufferPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	// Initialize writer pool with the specified level
	gc.writerPool = sync.Pool{
		New: func() interface{} {
			w, err := gzip.NewWriterLevel(nil, int(level))
			if err != nil {
				// Fall back to default level if invalid
				w = gzip.NewWriter(nil)
			}
			return w
		},
	}

	return gc
}

// Compress compresses data using gzip
func (gc *GzipCompressor) Compress(data []byte) ([]byte, error) {
	// Skip compression for small data
	if int64(len(data)) < gc.minSizeComp {
		return data, nil
	}

	// Get buffer from pool
	buf := gc.bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer func() {
		buf.Reset()
		gc.bufferPool.Put(buf)
	}()

	// Get writer from pool
	w := gc.writerPool.Get().(*gzip.Writer)
	w.Reset(buf)
	defer gc.writerPool.Put(w)

	// Write data
	if _, err := w.Write(data); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	// Make a copy of the compressed data
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
}

// Decompress decompresses gzip-compressed data
func (gc *GzipCompressor) Decompress(data []byte) ([]byte, error) {
	// Check if data is too small to be compressed data
	if int64(len(data)) < gc.minSizeComp {
		return data, nil
	}

	// Try to detect gzip header
	if len(data) < 2 || data[0] != 0x1f || data[1] != 0x8b {
		// Not gzip compressed
		return data, nil
	}

	// Get buffer from pool
	buf := gc.bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer func() {
		buf.Reset()
		gc.bufferPool.Put(buf)
	}()

	// Create reader
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		// If it's not valid gzip data, return as-is
		return data, nil
	}
	defer func() { _ = r.Close() }()

	// Decompress
	if _, err := io.Copy(buf, r); err != nil {
		return nil, err
	}

	// Make a copy of the decompressed data
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
}

// CompressionRatio calculates the compression ratio
func CompressionRatio(originalSize, compressedSize int64) float64 {
	if originalSize == 0 {
		return 0
	}
	return float64(compressedSize) / float64(originalSize) * 100
}

// SetMinimumSize sets the minimum size for compression to be applied
func (gc *GzipCompressor) SetMinimumSize(sizeBytes int64) {
	gc.minSizeComp = sizeBytes
}

// ShouldCompress determines if a data type is likely to benefit from compression
func ShouldCompress(contentType string) bool {
	// Text formats almost always compress well
	textTypes := []string{
		"text/",
		"application/json",
		"application/xml",
		"application/javascript",
		"application/x-javascript",
		"application/xhtml+xml",
		"application/soap+xml",
		"application/x-yaml",
		"application/yaml",
	}

	for _, textType := range textTypes {
		if strings.HasPrefix(contentType, textType) {
			return true
		}
	}

	// Already compressed formats
	compressedTypes := []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/webp",
		"audio/mp3",
		"audio/mp4",
		"video/mp4",
		"application/zip",
		"application/x-zip-compressed",
		"application/gzip",
		"application/x-gzip",
		"application/x-bzip2",
		"application/x-7z-compressed",
	}

	for _, compressedType := range compressedTypes {
		if strings.HasPrefix(contentType, compressedType) {
			return false
		}
	}

	// Default to compressing if unknown type
	return true
}

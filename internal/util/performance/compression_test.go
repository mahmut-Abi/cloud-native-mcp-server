package optimize

import (
	"testing"
)

func TestGzipCompression(t *testing.T) {
	compressor := NewGzipCompressor()

	// Test that the compressor works correctly - focus on correctness rather than compression ratio
	originalData := []byte("This is test data for compression. This data should be compressed and then decompressed correctly.")

	// For small data, it should return the original unchanged due to minSizeComp threshold
	compressed, err := compressor.Compress(originalData)
	if err != nil {
		t.Fatalf("Compress() error = %v", err)
	}

	// For small data, compressed should equal original due to minimum size threshold
	if string(compressed) != string(originalData) {
		t.Error("Small data should not be compressed due to minimum size threshold")
	}

	// Decompress should return the original data
	decompressed, err := compressor.Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompress() error = %v", err)
	}

	if string(decompressed) != string(originalData) {
		t.Error("Decompressed data doesn't match original")
	}

	// Test with larger data that should actually be compressed
	largeData := make([]byte, 2048)
	for i := range largeData {
		// Create somewhat repetitive but compressible data
		largeData[i] = byte('A' + (i % 26))
	}

	compressedLarge, err := compressor.Compress(largeData)
	if err != nil {
		t.Fatalf("Compress() error = %v", err)
	}

	// For large data, it should actually be processed (though may not compress well depending on content)
	decompressedLarge, err := compressor.Decompress(compressedLarge)
	if err != nil {
		t.Fatalf("Decompress() error = %v", err)
	}

	if string(decompressedLarge) != string(largeData) {
		t.Error("Decompressed large data doesn't match original")
	}
}

func TestCompressionRatio(t *testing.T) {
	tests := []struct {
		original   int64
		compressed int64
		expected   float64
	}{
		{100, 50, 50.0},
		{100, 75, 75.0},
		{0, 0, 0.0},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			ratio := CompressionRatio(tt.original, tt.compressed)
			if ratio != tt.expected {
				t.Errorf("CompressionRatio() = %f, want %f", ratio, tt.expected)
			}
		})
	}
}

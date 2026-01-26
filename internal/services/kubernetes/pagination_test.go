package kubernetes

import (
	"testing"
)

func TestValidatePaginationOptions(t *testing.T) {
	tests := []struct {
		name    string
		opts    *PaginationOptions
		wantErr bool
	}{
		{"nil options", nil, true},
		{"valid options", &PaginationOptions{Limit: 100}, false},
		{"limit too small", &PaginationOptions{Limit: 0}, true},
		{"limit too large", &PaginationOptions{Limit: 600}, true},
		{"limit at max boundary", &PaginationOptions{Limit: 500}, false},
		{"limit at min boundary", &PaginationOptions{Limit: 1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePaginationOptions(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePaginationOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultPaginationOptions(t *testing.T) {
	opts := DefaultPaginationOptions()

	if opts == nil {
		t.Fatal("DefaultPaginationOptions() returned nil")
	}

	if opts.Limit != 100 {
		t.Errorf("Expected default limit 100, got %d", opts.Limit)
	}
}

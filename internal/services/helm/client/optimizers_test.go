package client

import (
	"testing"
	"time"
)

func TestNewRepositoryOptimizer(t *testing.T) {
	tests := []struct {
		name         string
		timeoutSec   int
		maxRetry     int
		wantTimeout  time.Duration
		wantMaxRetry int
	}{
		{
			name:         "default values",
			timeoutSec:   0,
			maxRetry:     0,
			wantTimeout:  300 * time.Second,
			wantMaxRetry: 3,
		},
		{
			name:         "custom timeout and retry",
			timeoutSec:   600,
			maxRetry:     5,
			wantTimeout:  600 * time.Second,
			wantMaxRetry: 5,
		},
		{
			name:         "negative values use defaults",
			timeoutSec:   -1,
			maxRetry:     -1,
			wantTimeout:  300 * time.Second,
			wantMaxRetry: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := NewRepositoryOptimizer(nil, tt.timeoutSec, tt.maxRetry, true)

			if opt.GetTimeout() != tt.wantTimeout {
				t.Errorf("GetTimeout() = %v, want %v", opt.GetTimeout(), tt.wantTimeout)
			}
			if opt.GetMaxRetry() != tt.wantMaxRetry {
				t.Errorf("GetMaxRetry() = %d, want %d", opt.GetMaxRetry(), tt.wantMaxRetry)
			}
		})
	}
}

func TestResolveRepositoryURL(t *testing.T) {
	tests := []struct {
		name      string
		useMirror bool
		inputURL  string
		wantURL   string
	}{
		{
			name:      "mirror enabled, non-matching URL",
			useMirror: true,
			inputURL:  "https://unknown.repo.io",
			wantURL:   "https://unknown.repo.io",
		},
		{
			name:      "mirror disabled",
			useMirror: false,
			inputURL:  "https://kubernetes.github.io",
			wantURL:   "https://kubernetes.github.io",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := NewRepositoryOptimizer(nil, 300, 3, tt.useMirror)
			got := opt.ResolveRepositoryURL(tt.inputURL)

			if got != tt.wantURL {
				t.Errorf("ResolveRepositoryURL(%q) = %q, want %q", tt.inputURL, got, tt.wantURL)
			}
		})
	}
}

func TestHasMirror(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantHas bool
	}{
		{
			name:    "no mirror for kubernetes",
			url:     "https://kubernetes.github.io",
			wantHas: false,
		},
		{
			name:    "no mirror",
			url:     "https://unknown.repo.io",
			wantHas: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := NewRepositoryOptimizer(nil, 300, 3, true)
			got := opt.HasMirror(tt.url)

			if got != tt.wantHas {
				t.Errorf("HasMirror(%q) = %v, want %v", tt.url, got, tt.wantHas)
			}
		})
	}
}

func TestListMirrors(t *testing.T) {
	customMirrors := map[string]string{
		"https://example.com": "https://mirror.example.com",
	}

	opt := NewRepositoryOptimizer(customMirrors, 300, 3, true)
	mirrors := opt.ListMirrors()

	if len(mirrors) == 0 {
		t.Error("ListMirrors() returned empty map")
	}

	if val, ok := mirrors["https://example.com"]; !ok || val != "https://mirror.example.com" {
		t.Error("Custom mirror not found in ListMirrors()")
	}
}

func TestCreateOptimizedHTTPClient(t *testing.T) {
	opt := NewRepositoryOptimizer(nil, 600, 3, true)
	client := opt.CreateOptimizedHTTPClient()

	if client == nil {
		t.Fatal("CreateOptimizedHTTPClient() returned nil")
	}

	if client.Timeout != 600*time.Second {
		t.Errorf("Client timeout = %v, want %v", client.Timeout, 600*time.Second)
	}

	if client.Transport == nil {
		t.Error("Client transport is nil")
	}
}

func TestDefaultMirrorsEmpty(t *testing.T) {
	opt := NewRepositoryOptimizer(nil, 300, 3, true)
	mirrors := opt.ListMirrors()

	if len(mirrors) != 0 {
		t.Errorf("Expected empty default mirrors, got %d entries", len(mirrors))
	}
}

func TestOptimizeWithCustomMirrors(t *testing.T) {
	customMirrors := map[string]string{
		"https://original.com": "https://mirror.com",
	}

	opt := NewRepositoryOptimizer(customMirrors, 300, 3, true)

	resolved := opt.ResolveRepositoryURL("https://original.com")
	if resolved != "https://mirror.com" {
		t.Errorf("Custom mirror not applied: got %s, want https://mirror.com", resolved)
	}
}

func TestGetOptimizer(t *testing.T) {
	opt := NewRepositoryOptimizer(nil, 300, 3, true)

	if opt.GetTimeout() <= 0 {
		t.Error("GetTimeout() returned zero or negative value")
	}

	if opt.GetMaxRetry() <= 0 {
		t.Error("GetMaxRetry() returned zero or negative value")
	}
}

func TestIsMirrorEnabled(t *testing.T) {
	tests := []struct {
		name      string
		useMirror bool
		want      bool
	}{
		{
			name:      "mirrors enabled",
			useMirror: true,
			want:      true,
		},
		{
			name:      "mirrors disabled",
			useMirror: false,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := NewRepositoryOptimizer(nil, 300, 3, tt.useMirror)
			got := opt.IsMirrorEnabled()

			if got != tt.want {
				t.Errorf("IsMirrorEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMirrorConfiguration(t *testing.T) {
	customMirrors := map[string]string{
		"https://custom.repo": "https://custom.mirror",
	}

	opt := NewRepositoryOptimizer(customMirrors, 600, 5, true)

	// Test that custom mirrors are merged with default mirrors
	mirrors := opt.ListMirrors()

	// Check custom mirror
	if val, ok := mirrors["https://custom.repo"]; !ok || val != "https://custom.mirror" {
		t.Error("Custom mirror configuration not preserved")
	}

	// Test mirror enabled flag
	if !opt.IsMirrorEnabled() {
		t.Error("IsMirrorEnabled() returned false when useMirror=true")
	}
}

func TestDisabledMirrorsNoResolution(t *testing.T) {
	customMirrors := map[string]string{
		"https://example.com": "https://mirror.example.com",
	}

	opt := NewRepositoryOptimizer(customMirrors, 300, 3, false)

	// Even with custom mirrors defined, if useMirror is false, no resolution should happen
	resolved := opt.ResolveRepositoryURL("https://example.com")
	if resolved != "https://example.com" {
		t.Errorf("URL resolved when mirrors disabled: got %s, want https://example.com", resolved)
	}

	if opt.IsMirrorEnabled() {
		t.Error("IsMirrorEnabled() returned true when useMirror=false")
	}
}

func TestResolveMultipleMirrors(t *testing.T) {
	customMirrors := map[string]string{
		"https://repo1.com": "https://mirror1.com",
		"https://repo2.com": "https://mirror2.com",
		"https://repo3.com": "https://mirror3.com",
	}

	opt := NewRepositoryOptimizer(customMirrors, 300, 3, true)

	tests := []struct {
		input    string
		expected string
	}{
		{"https://repo1.com", "https://mirror1.com"},
		{"https://repo2.com", "https://mirror2.com"},
		{"https://repo3.com", "https://mirror3.com"},
		{"https://unknown.com", "https://unknown.com"},
	}

	for _, tt := range tests {
		resolved := opt.ResolveRepositoryURL(tt.input)
		if resolved != tt.expected {
			t.Errorf("ResolveRepositoryURL(%q) = %q, want %q", tt.input, resolved, tt.expected)
		}
	}
}

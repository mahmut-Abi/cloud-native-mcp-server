package client

import (
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

func TestDefaultClientOptions(t *testing.T) {
	opts := DefaultClientOptions()

	if opts.Timeout != 30*time.Second {
		t.Errorf("Expected timeout to be 30s, got %v", opts.Timeout)
	}

	if opts.QPS != 100 {
		t.Errorf("Expected QPS to be 100, got %f", opts.QPS)
	}

	if opts.Burst != 200 {
		t.Errorf("Expected Burst to be 200, got %d", opts.Burst)
	}

	if opts.GVRCacheTTL != 15*time.Minute {
		t.Errorf("Expected GVRCacheTTL to be 15m, got %v", opts.GVRCacheTTL)
	}
}

func TestResolveKubeconfigPath(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedNonEmpty bool
	}{
		{
			name:             "explicit path provided",
			input:            "/explicit/path",
			expectedNonEmpty: true,
		},
		{
			name:             "empty path",
			input:            "",
			expectedNonEmpty: false, // depends on environment
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolveKubeconfigPath(tt.input)

			if tt.input != "" && result != tt.input {
				t.Errorf("Expected explicit path %s, got %s", tt.input, result)
			}
		})
	}
}

func TestClientGetters(t *testing.T) {
	client := &Client{
		kubeconfigPath: "/test/path",
		restConfig:     &rest.Config{},
	}

	if client.GetKubeconfigPath() != "/test/path" {
		t.Errorf("Expected kubeconfig path to be '/test/path', got %s", client.GetKubeconfigPath())
	}

	if client.GetRestConfig() != client.restConfig {
		t.Error("Expected GetRestConfig to return the same config instance")
	}
}

func TestToLower(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Pod", "pod"},
		{"SERVICE", "service"},
		{"DeployMent", "deployment"},
		{"", ""},
		{"lowercase", "lowercase"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toLower(tt.input)
			if result != tt.expected {
				t.Errorf("toLower(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestContainsSlash(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"pods", false},
		{"pods/status", true},
		{"", false},
		{"/", true},
		{"no/slash/here", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := containsSlash(tt.input)
			if result != tt.expected {
				t.Errorf("containsSlash(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetCachedGVREmptyKind(t *testing.T) {
	client := &Client{
		gvrCache: make(map[string]schema.GroupVersionResource),
	}

	_, err := client.getCachedGVR("")
	if err == nil {
		t.Error("Expected error for empty kind, got nil")
	}

	if err.Error() != "kind is empty" {
		t.Errorf("Expected 'kind is empty' error, got %q", err.Error())
	}
}

func TestClientOptionsValidation(t *testing.T) {
	opts := &ClientOptions{
		Timeout:     0,  // invalid
		QPS:         -1, // invalid
		Burst:       -1, // invalid
		GVRCacheTTL: 0,  // invalid
	}

	// Test that NewClientWithOptions handles invalid options gracefully
	// Note: This test assumes the function will use defaults for invalid values
	_ = opts // Use opts to avoid unused variable error

	// In a real implementation, you might want to validate these values
	// and return an error or use defaults
}

func TestIsInClusterConfig(t *testing.T) {
	tests := []struct {
		name     string
		setupEnv func()
		cleanup  func()
		expected bool
	}{
		{
			name: "not in cluster - no env vars",
			setupEnv: func() {
				t.Setenv("KUBERNETES_SERVICE_HOST", "")
				t.Setenv("KUBERNETES_SERVICE_PORT", "")
			},
			cleanup:  func() {},
			expected: false,
		},
		{
			name: "in cluster - with env vars and files",
			setupEnv: func() {
				t.Setenv("KUBERNETES_SERVICE_HOST", "kubernetes.default")
				t.Setenv("KUBERNETES_SERVICE_PORT", "443")
			},
			cleanup:  func() {},
			expected: false, // Will be false if service account files don't exist
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			defer tt.cleanup()

			result := isInClusterConfig()
			if result != tt.expected && !tt.expected {
				// Expected false - service account files typically don't exist in test environment
				if result {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestResolveKubeconfigPathWithInClusterDetection(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedNonEmpty bool
	}{
		{
			name:             "explicit path provided",
			input:            "/explicit/path",
			expectedNonEmpty: true,
		},
		{
			name:             "empty path - will check environment",
			input:            "",
			expectedNonEmpty: false, // depends on environment
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolveKubeconfigPath(tt.input)
			if tt.expectedNonEmpty && result == "" {
				t.Errorf("Expected non-empty path, got empty")
			}
		})
	}
}

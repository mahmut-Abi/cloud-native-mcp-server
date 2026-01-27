package sanitize

import (
	"testing"
)

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple text",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "dangerous characters",
			input:    "test;drop table",
			expected: "testdrop table",
		},
		{
			name:     "shell metachars",
			input:    "test|command",
			expected: "testcommand",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "long string",
			input:    "a" + string(make([]byte, 1200)),
			expected: "a", // Null bytes are removed, leaving just "a"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeString(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeString() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSanitizeQuery(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple query",
			input:    "SELECT * FROM users",
			expected: "* FROM users",
		},
		{
			name:     "SQL injection",
			input:    "'; DROP TABLE users; --",
			expected: "TABLE users --",
		},
		{
			name:     "empty query",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeQuery(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeQuery() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSanitizeFilterValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "string input",
			input:    "test;value",
			expected: "testvalue",
		},
		{
			name:     "nil input",
			input:    nil,
			expected: "",
		},
		{
			name:     "number input",
			input:    42,
			expected: "42",
		},
		{
			name:     "special chars",
			input:    "test|command",
			expected: "testcommand",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeFilterValue(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeFilterValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSanitizeIdentifier(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid identifier",
			input:    "test_user-123",
			expected: "test_user-123",
		},
		{
			name:     "invalid characters",
			input:    "test@user",
			expected: "testuser",
		},
		{
			name:     "spaces",
			input:    "test user",
			expected: "testuser",
		},
		{
			name:     "empty",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeIdentifier(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeIdentifier() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSanitizeJSONPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid JSONPath",
			input:    "$.users[0].name",
			expected: ".users[0].name",
		},
		{
			name:     "invalid characters",
			input:    "$.users; DROP TABLE",
			expected: ".users DROP TABLE",
		},
		{
			name:     "empty",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeJSONPath(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeJSONPath() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestValidateString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxLength int
		wantError  bool
	}{
		{
			name:      "valid short string",
			input:     "test",
			maxLength: 100,
			wantError:  false,
		},
		{
			name:      "too long string",
			input:     string(make([]byte, 200)),
			maxLength: 100,
			wantError:  true,
		},
		{
			name:      "string with control chars",
			input:     "test\x00data",
			maxLength: 100,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateString(tt.input, tt.maxLength)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateString() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateIdentifier(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantError bool
	}{
		{
			name:     "valid identifier",
			input:    "test_user-123",
			wantError: false,
		},
		{
			name:     "empty identifier",
			input:    "",
			wantError: true,
		},
		{
			name:     "invalid characters",
			input:    "test@user",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIdentifier(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateIdentifier() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
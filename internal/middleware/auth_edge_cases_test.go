package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAuthMiddleware_EdgeCases tests various edge cases and configuration scenarios
func TestAuthMiddleware_EdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		config         AuthConfig
		setupRequest   func(*http.Request)
		expectedStatus int
		description    string
	}{
		{
			name: "empty_config_disabled",
			config: AuthConfig{
				Enabled: false,
			},
			setupRequest: func(r *http.Request) {
				// No auth headers
			},
			expectedStatus: http.StatusOK,
			description:    "Disabled auth should allow all requests",
		},
		{
			name: "apikey_mode_empty_config_key",
			config: AuthConfig{
				Enabled: true,
				Mode:    "apikey",
				APIKey:  "", // Empty expected key
			},
			setupRequest: func(r *http.Request) {
				r.Header.Set("X-API-Key", "any-key")
			},
			expectedStatus: http.StatusUnauthorized,
			description:    "Empty config API key should reject even if header provided",
		},
		{
			name: "apikey_query_parameter",
			config: AuthConfig{
				Enabled: true,
				Mode:    "apikey",
				APIKey:  "query-key-123",
			},
			setupRequest: func(r *http.Request) {
				// Set API key via query parameter instead of header
				q := r.URL.Query()
				q.Set("api_key", "query-key-123")
				r.URL.RawQuery = q.Encode()
			},
			expectedStatus: http.StatusOK,
			description:    "API key via query parameter should work",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test handler
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("success"))
			})

			// Wrap with auth middleware
			authedHandler := AuthMiddleware(tt.config)(handler)

			// Create test request
			req := httptest.NewRequest("GET", "/test", nil)
			tt.setupRequest(req)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute request
			authedHandler.ServeHTTP(rr, req)

			// Check result
			assert.Equal(t, tt.expectedStatus, rr.Code, tt.description)
		})
	}
}

// TestAuthMiddleware_QuoteHandling tests API key with quotes
func TestAuthMiddleware_QuoteHandling(t *testing.T) {
	tests := []struct {
		name           string
		headerValue    string
		configAPIKey   string
		expectedStatus int
		description    string
	}{
		{
			name:           "API key with single quotes",
			headerValue:    "'12345678'",
			configAPIKey:   "12345678",
			expectedStatus: http.StatusOK,
			description:    "Single quotes should be stripped",
		},
		{
			name:           "API key with double quotes",
			headerValue:    "\"12345678\"",
			configAPIKey:   "12345678",
			expectedStatus: http.StatusOK,
			description:    "Double quotes should be stripped",
		},
		{
			name:           "API key with whitespace",
			headerValue:    "  12345678  ",
			configAPIKey:   "12345678",
			expectedStatus: http.StatusOK,
			description:    "Whitespace should be trimmed",
		},
		{
			name:           "API key with quotes and whitespace",
			headerValue:    "  '12345678'  ",
			configAPIKey:   "12345678",
			expectedStatus: http.StatusOK,
			description:    "Both quotes and whitespace should be cleaned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create auth config
			config := AuthConfig{
				Enabled: true,
				Mode:    "apikey",
				APIKey:  tt.configAPIKey,
			}

			// Create test handler
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("success"))
			})

			// Wrap with auth middleware
			authedHandler := AuthMiddleware(config)(handler)

			// Create test request
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("X-Api-Key", tt.headerValue)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute request
			authedHandler.ServeHTTP(rr, req)

			// Check result
			assert.Equal(t, tt.expectedStatus, rr.Code, tt.description)
		})
	}
}

// TestValidateAPIKey tests API key validation
func TestValidateAPIKey(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{"valid_key", "12345678", true},
		{"long_key", "this-is-a-very-long-api-key-123456789", true},
		{"short_key", "1234567", false}, // Less than 8 chars
		{"empty_key", "", false},
		{"exactly_8_chars", "abcdefgh", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateAPIKey(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestValidateBearerToken tests bearer token validation
func TestValidateBearerToken(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected bool
	}{
		{"valid_token", "1234567890123456", true}, // 16 chars
		{"long_token", "this-is-a-very-long-bearer-token-123456789", true},
		{"short_token", "123456789012345", false}, // Less than 16 chars
		{"empty_token", "", false},
		{"exactly_16_chars", "abcdefgh12345678", true},
		{"jwt_like_token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateBearerToken(tt.token)
			assert.Equal(t, tt.expected, result)
		})
	}
}

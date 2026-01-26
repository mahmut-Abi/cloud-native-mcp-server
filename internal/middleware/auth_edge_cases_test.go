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

// TestValidateAPIKey tests API key validation with complexity requirements
func TestValidateAPIKey(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{"valid_key_complex", "Abc123!@#Xyz789!@#", true},             // Has uppercase, lowercase, digits, special chars
		{"valid_key_upper_lower_digit", "Abc123Xyz789Abc123", true},   // Has uppercase, lowercase, digits
		{"valid_key_upper_lower_special", "Abc!@#Xyz!@#Abc!@#", true}, // Has uppercase, lowercase, special chars
		{"valid_key_upper_digit_special", "ABC123!@#XYZ789!@#", true}, // Has uppercase, digits, special chars
		{"valid_key_lower_digit_special", "abc123!@#xyz789!@#", true}, // Has lowercase, digits, special chars
		{"long_key", "ThisIsAVeryLongAPIKey123!@#WithSpecialChars", true},
		{"short_key", "Ab1!", false}, // Less than 16 chars
		{"empty_key", "", false},
		{"exactly_16_chars_complex", "Abc123!@#Xyz789!@", true}, // Exactly 16 chars with complexity
		{"exactly_16_chars_simple", "abcdefgh12345678", false},  // Exactly 16 chars but only lowercase and digits
		{"only_uppercase", "ABCDEFGHIJKLMNOPQRSTUVWXYZ", false}, // Only uppercase
		{"only_lowercase", "abcdefghijklmnopqrstuvwxyz", false}, // Only lowercase
		{"only_digits", "12345678901234567890", false},          // Only digits
		{"only_special", "!@#$%^&*()_+-=[]{}|;:,.<>?", false},   // Only special chars
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateAPIKey(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestValidateBearerToken tests bearer token validation with JWT structure checks
func TestValidateBearerToken(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected bool
	}{
		{"valid_jwt_token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c", true},
		{"long_jwt_token", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.NHVaYe26MbtOYhSKkoKYdFVomg4i8ZJd8_-RU8VNbftc4TSMb4bXP3l3YlNWACwyXPGffz5aXHc6lty1Y2t4SWRqGteragsVdZufDn5BlnJl9pdR_kdVFUsra2rWKEofkZeIC4yWytE58sMIihvo9H1ScmmVwBcQP6XETqYd0aSHp1gOa9RdUPDvoXQ5oqygTqVtxaDr6wUFKrKItgBMzWIdNZ6y7O9E0DhEPTbE9rfBo6KTFsHAZnMg4k68CDp2woYIaXbmYTWcvbzIuHO7_37GT79XdIwkm95QJ7hYC9RiwrV7mesbY4PAahERJawntho0my942XheVLmGwLMBkQ", true},
		{"short_token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ", false}, // Less than 32 chars total
		{"empty_token", "", false},
		{"simple_token_no_dots", "abcdefgh12345678abcdefgh12345678", false}, // No JWT structure
		{"token_with_invalid_chars", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c$", false}, // Invalid char at end
		{"token_with_space", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c ", false},         // Trailing space
		{"token_with_plus", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV+adQssw5c", true},            // Plus char is valid in base64url
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateBearerToken(tt.token)
			assert.Equal(t, tt.expected, result)
		})
	}
}

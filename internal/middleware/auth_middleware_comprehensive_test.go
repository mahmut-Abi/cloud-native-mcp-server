package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestAuthMiddleware_APIKey tests API key authentication with different header cases
func TestAuthMiddleware_APIKey(t *testing.T) {
	tests := []struct {
		name           string
		headerName     string
		headerValue    string
		configAPIKey   string
		expectedStatus int
	}{
		{
			name:           "Valid API key with X-API-Key",
			headerName:     "X-API-Key",
			headerValue:    "test-api-key-123",
			configAPIKey:   "test-api-key-123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid API key with X-Api-Key (mixed case)",
			headerName:     "X-Api-Key",
			headerValue:    "test-api-key-123",
			configAPIKey:   "test-api-key-123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid API key with x-api-key (lowercase)",
			headerName:     "x-api-key",
			headerValue:    "test-api-key-123",
			configAPIKey:   "test-api-key-123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid API key",
			headerName:     "X-API-Key",
			headerValue:    "wrong-key",
			configAPIKey:   "test-api-key-123",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Missing API key",
			headerName:     "",
			headerValue:    "",
			configAPIKey:   "test-api-key-123",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Empty API key in header",
			headerName:     "X-API-Key",
			headerValue:    "",
			configAPIKey:   "test-api-key-123",
			expectedStatus: http.StatusUnauthorized,
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
			if tt.headerName != "" {
				req.Header.Set(tt.headerName, tt.headerValue)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute request
			authedHandler.ServeHTTP(rr, req)

			// Check result
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

// TestAuthMiddleware_Bearer tests Bearer token authentication
func TestAuthMiddleware_Bearer(t *testing.T) {
	tests := []struct {
		name           string
		headerValue    string
		configToken    string
		expectedStatus int
	}{
		{
			name:           "Valid Bearer token",
			headerValue:    "Bearer valid-token-123",
			configToken:    "valid-token-123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Bearer token",
			headerValue:    "Bearer wrong-token",
			configToken:    "valid-token-123",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Missing Bearer prefix",
			headerValue:    "valid-token-123",
			configToken:    "valid-token-123",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Empty Authorization header",
			headerValue:    "",
			configToken:    "valid-token-123",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create auth config
			config := AuthConfig{
				Enabled:     true,
				Mode:        "bearer",
				BearerToken: tt.configToken,
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
			if tt.headerValue != "" {
				req.Header.Set("Authorization", tt.headerValue)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute request
			authedHandler.ServeHTTP(rr, req)

			// Check result
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

// TestAuthMiddleware_Basic tests Basic authentication
func TestAuthMiddleware_Basic(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		password       string
		configUser     string
		configPass     string
		expectedStatus int
	}{
		{
			name:           "Valid credentials",
			username:       "admin",
			password:       "secret123",
			configUser:     "admin",
			configPass:     "secret123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Wrong username",
			username:       "wronguser",
			password:       "secret123",
			configUser:     "admin",
			configPass:     "secret123",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Wrong password",
			username:       "admin",
			password:       "wrongpass",
			configUser:     "admin",
			configPass:     "secret123",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Empty credentials",
			username:       "",
			password:       "",
			configUser:     "admin",
			configPass:     "secret123",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create auth config
			config := AuthConfig{
				Enabled:  true,
				Mode:     "basic",
				Username: tt.configUser,
				Password: tt.configPass,
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
			if tt.username != "" || tt.password != "" {
				req.SetBasicAuth(tt.username, tt.password)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute request
			authedHandler.ServeHTTP(rr, req)

			// Check result
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

// TestAuthMiddleware_Disabled tests disabled authentication
func TestAuthMiddleware_Disabled(t *testing.T) {
	config := AuthConfig{
		Enabled: false,
		Mode:    "apikey",
		APIKey:  "test-key",
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	})

	authedHandler := AuthMiddleware(config)(handler)

	// Request without any auth headers should succeed
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	authedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

// TestGetAPIKeyFromRequest tests the new case-insensitive header extraction
func TestGetAPIKeyFromRequest(t *testing.T) {
	tests := []struct {
		name        string
		headerName  string
		headerValue string
		queryParam  string
		expectedKey string
	}{
		{
			name:        "X-API-Key header",
			headerName:  "X-API-Key",
			headerValue: "header-key-123",
			expectedKey: "header-key-123",
		},
		{
			name:        "X-Api-Key header (mixed case)",
			headerName:  "X-Api-Key",
			headerValue: "mixed-case-key",
			expectedKey: "mixed-case-key",
		},
		{
			name:        "x-api-key header (lowercase)",
			headerName:  "x-api-key",
			headerValue: "lower-case-key",
			expectedKey: "lower-case-key",
		},
		{
			name:        "Query parameter fallback",
			queryParam:  "query-key-123",
			expectedKey: "query-key-123",
		},
		{
			name:        "Header takes precedence over query",
			headerName:  "X-API-Key",
			headerValue: "header-key",
			queryParam:  "query-key",
			expectedKey: "header-key",
		},
		{
			name:        "No key provided",
			expectedKey: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)

			if tt.headerName != "" {
				req.Header.Set(tt.headerName, tt.headerValue)
			}

			if tt.queryParam != "" {
				q := req.URL.Query()
				q.Set("api_key", tt.queryParam)
				req.URL.RawQuery = q.Encode()
			}

			actualKey := getAPIKeyFromRequest(req)

			if actualKey != tt.expectedKey {
				t.Errorf("Expected key '%s', got '%s'", tt.expectedKey, actualKey)
			}
		})
	}
}

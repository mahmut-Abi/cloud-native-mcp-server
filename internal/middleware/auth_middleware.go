package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

const authComponent = "auth"

// Pre-compiled loggers with common fields
var (
	authLogger = logrus.WithField("component", authComponent)
)

type AuthConfig struct {
	Enabled     bool
	Mode        string // apikey, bearer, basic
	APIKey      string
	BearerToken string
	Username    string
	Password    string
}

// AuthMiddleware creates an HTTP middleware for authentication
func AuthMiddleware(config AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Log auth middleware entry
			authLogger.WithFields(logrus.Fields{
				"enabled": config.Enabled,
				"mode":    config.Mode,
				"path":    r.URL.Path,
			}).Debug("Auth middleware processing request")

			if !config.Enabled {
				authLogger.Debug("Authentication disabled, proceeding")
				next.ServeHTTP(w, r)
				return
			}

			if !authenticate(r, config) {
				authLogger.Error("Authentication failed - returning 401")
				authLogger.Warnf("Authentication failed for request from %s to %s", r.RemoteAddr, r.RequestURI)
				w.WriteHeader(http.StatusUnauthorized)
				authLogger.WithFields(logrus.Fields{
					"operation":   "authenticate",
					"http_method": r.Method,
					"http_path":   r.URL.Path,
					"remote_addr": r.RemoteAddr,
					"status":      "failed",
				}).Warn("Authentication failed")
				_, _ = fmt.Fprint(w, "{\"error\":\"unauthorized\"}")
				return
			}

			next.ServeHTTP(w, r)
			authLogger.Debug("Authentication successful, proceeding")
		})
	}
}

// authenticate checks if the request has valid authentication
func authenticate(r *http.Request, config AuthConfig) bool {
	authLogger.WithFields(logrus.Fields{
		"mode":        config.Mode,
		"has_api_key": getAPIKeyFromRequest(r) != "",
	}).Debug("Starting authentication process")

	switch config.Mode {
	case "apikey":
		providedKey := getAPIKeyFromRequest(r)
		authLogger.WithFields(logrus.Fields{
			"mode":                config.Mode,
			"provided_key_length": len(providedKey),
			"expected_key_length": len(config.APIKey),
			"match":               providedKey == config.APIKey,
		}).Debug("API Key auth attempt")
		return authenticateAPIKey(r, config.APIKey)
	case "bearer":
		authLogger.Debug("Processing Bearer token authentication")
		return authenticateBearer(r, config.BearerToken)
	case "basic":
		authLogger.Debug("Processing Basic authentication")
		return authenticateBasic(r, config.Username, config.Password)
	default:
		authLogger.WithField("mode", config.Mode).Error("Unknown authentication mode")
		return false
	}
}

// getAPIKeyFromRequest extracts API key from request headers or query params
// Uses canonical header name for efficiency and cleans the key by removing quotes
func getAPIKeyFromRequest(r *http.Request) string {
	// Check canonical header name first
	if key := r.Header.Get("X-Api-Key"); key != "" {
		// Clean the key by trimming whitespace and quotes
		cleanKey := strings.TrimSpace(key)
		cleanKey = strings.Trim(cleanKey, "'\"") // Remove single and double quotes
		return cleanKey
	}

	// Fallback to query parameter
	if queryKey := r.URL.Query().Get("api_key"); queryKey != "" {
		// Clean the key by trimming whitespace and quotes
		cleanKey := strings.TrimSpace(queryKey)
		cleanKey = strings.Trim(cleanKey, "'\"") // Remove single and double quotes
		return cleanKey
	}

	return ""
}

// authenticateAPIKey checks API key authentication
func authenticateAPIKey(r *http.Request, expectedKey string) bool {
	key := getAPIKeyFromRequest(r)
	authLogger.WithFields(logrus.Fields{
		"provided_key_length": len(key),
		"expected_key_length": len(expectedKey),
		"keys_match":          key == expectedKey,
		"expected_empty":      expectedKey == "",
	}).Debug("API key authentication check")
	return key == expectedKey && expectedKey != ""
}

// authenticateBearer checks Bearer token authentication
func authenticateBearer(r *http.Request, expectedToken string) bool {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return false
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	return token == expectedToken && expectedToken != ""
}

// authenticateBasic checks Basic authentication
func authenticateBasic(r *http.Request, username, password string) bool {
	user, pass, ok := r.BasicAuth()
	if !ok {
		return false
	}
	return user == username && pass == password && username != "" && password != ""
}

// ValidateAPIKey validates an API key format with complexity requirements
func ValidateAPIKey(key string) bool {
	if key == "" || len(key) < 16 {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	// Check for required character classes
	for _, char := range key {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case char == '!' || char == '@' || char == '#' || char == '$' ||
			char == '%' || char == '^' || char == '&' || char == '*' ||
			char == '(' || char == ')' || char == '-' || char == '_' ||
			char == '+' || char == '=' || char == '[' || char == ']' ||
			char == '{' || char == '}' || char == '|' || char == '\\' ||
			char == ';' || char == ':' || char == '\'' || char == '"' ||
			char == '<' || char == '>' || char == ',' || char == '.' ||
			char == '?' || char == '/' || char == '~' || char == '`':
			hasSpecial = true
		}
	}

	// Require at least 3 of 4 character classes
	classCount := 0
	if hasUpper {
		classCount++
	}
	if hasLower {
		classCount++
	}
	if hasDigit {
		classCount++
	}
	if hasSpecial {
		classCount++
	}

	return classCount >= 3
}

// ValidateBearerToken validates a bearer token format with JWT structure checks
func ValidateBearerToken(token string) bool {
	if token == "" || len(token) < 32 {
		return false
	}

	// Check if token follows JWT structure (header.payload.signature)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}

	// Validate each part is base64url encoded
	for _, part := range parts {
		if part == "" {
			return false
		}
		// Check for valid base64url characters
		for _, char := range part {
			valid := (char >= 'A' && char <= 'Z') ||
				(char >= 'a' && char <= 'z') ||
				(char >= '0' && char <= '9') ||
				char == '-' || char == '_' || char == '+'
			if !valid {
				return false
			}
		}
	}

	return true
}

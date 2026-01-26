package middleware

import (
	"testing"
	"time"
)

func TestValidateJWT(t *testing.T) {
	// Test JWT validation
	secret := "test-secret"

	// Test JWT token format validation
	token := "invalid.token"
	_, err := ValidateJWT(token, secret)
	if err == nil {
		t.Error("Expected error for invalid token format")
	}

	token2 := "a.b.c.d"
	_, err2 := ValidateJWT(token2, secret)
	if err2 == nil {
		t.Error("Expected error for token with too many parts")
	}
}

func TestJWTClaims(t *testing.T) {
	claims := JWTClaims{
		Sub: "user123",
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(1 * time.Hour).Unix(),
		Iss: "mcp-server",
	}

	if claims.Sub != "user123" {
		t.Errorf("Expected subject user123, got %s", claims.Sub)
	}

	if claims.Iss != "mcp-server" {
		t.Errorf("Expected issuer mcp-server, got %s", claims.Iss)
	}
}

package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type JWTClaims struct {
	Sub string `json:"sub"` // Subject
	Iat int64  `json:"iat"` // Issued at
	Exp int64  `json:"exp"` // Expiration
	Iss string `json:"iss"` // Issuer
}

// ValidateJWT validates a JWT token
func ValidateJWT(token string, secret string) (*JWTClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	// Verify signature
	message := parts[0] + "." + parts[1]
	expectedSig := createHMAC(message, secret)
	actualSig := parts[2]

	if expectedSig != actualSig {
		return nil, fmt.Errorf("invalid signature")
	}

	// Decode payload
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}

	var claims JWTClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	// Check expiration
	if claims.Exp > 0 && time.Now().Unix() > claims.Exp {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}

// createHMAC creates HMAC signature
func createHMAC(message, secret string) string {
	if secret == "" {
		return ""
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))

	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

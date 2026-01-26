// Package secrets provides secure secrets management for API keys, tokens, and other sensitive data.
package secrets

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// SecretType defines the type of secret
type SecretType string

const (
	// SecretTypeAPIKey represents an API key secret
	SecretTypeAPIKey SecretType = "apikey"
	// SecretTypeBearerToken represents a bearer token secret
	SecretTypeBearerToken SecretType = "bearertoken"
	// SecretTypeBasicAuth represents basic auth credentials
	SecretTypeBasicAuth SecretType = "basicauth"
	// SecretTypeGeneric represents a generic secret
	SecretTypeGeneric SecretType = "generic"
)

// Secret represents a stored secret with metadata
type Secret struct {
	ID        string            `json:"id"`
	Type      SecretType        `json:"type"`
	Name      string            `json:"name"`
	Value     string            `json:"value,omitempty"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
	ExpiresAt *time.Time        `json:"expiresAt,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// Manager defines the interface for secrets management
type Manager interface {
	// Store stores a secret securely
	Store(secret *Secret) error

	// Retrieve retrieves a secret by ID
	Retrieve(id string) (*Secret, error)

	// List lists all secrets of a given type
	List(secretType SecretType) ([]*Secret, error)

	// Delete deletes a secret by ID
	Delete(id string) error

	// Rotate rotates a secret's value
	Rotate(id string) (*Secret, error)

	// GenerateAPIKey generates a new API key with complexity requirements
	GenerateAPIKey(name string) (*Secret, error)

	// GenerateBearerToken generates a new bearer token with JWT structure
	GenerateBearerToken(name string) (*Secret, error)
}

// InMemoryManager provides an in-memory implementation of secrets manager
// In production, this should be replaced with a secure backend like HashiCorp Vault
type InMemoryManager struct {
	secrets map[string]*Secret
	mu      sync.RWMutex
	logger  *logrus.Logger
}

// NewInMemoryManager creates a new in-memory secrets manager
func NewInMemoryManager() *InMemoryManager {
	return &InMemoryManager{
		secrets: make(map[string]*Secret),
		logger:  logrus.WithField("component", "secrets-manager").Logger,
	}
}

// Store stores a secret securely
func (m *InMemoryManager) Store(secret *Secret) error {
	if secret == nil {
		return fmt.Errorf("secret cannot be nil")
	}

	if secret.ID == "" {
		secret.ID = generateID()
	}

	if secret.CreatedAt.IsZero() {
		secret.CreatedAt = time.Now()
	}

	secret.UpdatedAt = time.Now()

	m.mu.Lock()
	defer m.mu.Unlock()

	m.secrets[secret.ID] = secret
	m.logger.WithFields(logrus.Fields{
		"id":   secret.ID,
		"type": secret.Type,
		"name": secret.Name,
	}).Info("Secret stored")

	return nil
}

// Retrieve retrieves a secret by ID
func (m *InMemoryManager) Retrieve(id string) (*Secret, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	secret, exists := m.secrets[id]
	if !exists {
		return nil, fmt.Errorf("secret not found: %s", id)
	}

	// Check if secret has expired
	if secret.ExpiresAt != nil && time.Now().After(*secret.ExpiresAt) {
		return nil, fmt.Errorf("secret has expired: %s", id)
	}

	// Return a copy to prevent modification
	copy := *secret
	return &copy, nil
}

// List lists all secrets of a given type
func (m *InMemoryManager) List(secretType SecretType) ([]*Secret, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*Secret
	for _, secret := range m.secrets {
		if secret.Type == secretType {
			// Check if secret has expired
			if secret.ExpiresAt == nil || time.Now().Before(*secret.ExpiresAt) {
				// Return a copy to prevent modification
				copy := *secret
				results = append(results, &copy)
			}
		}
	}

	return results, nil
}

// Delete deletes a secret by ID
func (m *InMemoryManager) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.secrets[id]; !exists {
		return fmt.Errorf("secret not found: %s", id)
	}

	delete(m.secrets, id)
	m.logger.WithField("id", id).Info("Secret deleted")

	return nil
}

// Rotate rotates a secret's value
func (m *InMemoryManager) Rotate(id string) (*Secret, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	secret, exists := m.secrets[id]
	if !exists {
		return nil, fmt.Errorf("secret not found: %s", id)
	}

	// Generate new value based on type
	switch secret.Type {
	case SecretTypeAPIKey:
		newValue, err := generateComplexAPIKey()
		if err != nil {
			return nil, fmt.Errorf("failed to generate new API key: %w", err)
		}
		secret.Value = newValue
	case SecretTypeBearerToken:
		newValue, err := generateJWTLikeToken()
		if err != nil {
			return nil, fmt.Errorf("failed to generate new bearer token: %w", err)
		}
		secret.Value = newValue
	default:
		return nil, fmt.Errorf("rotation not supported for secret type: %s", secret.Type)
	}

	secret.UpdatedAt = time.Now()

	m.logger.WithField("id", id).Info("Secret rotated")

	// Return a copy to prevent modification
	copy := *secret
	return &copy, nil
}

// GenerateAPIKey generates a new API key with complexity requirements
func (m *InMemoryManager) GenerateAPIKey(name string) (*Secret, error) {
	value, err := generateComplexAPIKey()
	if err != nil {
		return nil, err
	}

	secret := &Secret{
		ID:        generateID(),
		Type:      SecretTypeAPIKey,
		Name:      name,
		Value:     value,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata: map[string]string{
			"generated": "true",
		},
	}

	if err := m.Store(secret); err != nil {
		return nil, err
	}

	return secret, nil
}

// GenerateBearerToken generates a new bearer token with JWT structure
func (m *InMemoryManager) GenerateBearerToken(name string) (*Secret, error) {
	value, err := generateJWTLikeToken()
	if err != nil {
		return nil, err
	}

	secret := &Secret{
		ID:        generateID(),
		Type:      SecretTypeBearerToken,
		Name:      name,
		Value:     value,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata: map[string]string{
			"generated": "true",
		},
	}

	if err := m.Store(secret); err != nil {
		return nil, err
	}

	return secret, nil
}

// generateComplexAPIKey generates a complex API key meeting security requirements
// Minimum 16 characters with at least 3 of 4 character classes:
// - Uppercase letters
// - Lowercase letters
// - Digits
// - Special characters
func generateComplexAPIKey() (string, error) {
	const (
		upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lowerChars   = "abcdefghijklmnopqrstuvwxyz"
		digitChars   = "0123456789"
		specialChars = "!@#$%^&*()_+-=[]{}|;:,.<>?"
		minLength    = 24
	)

	// Ensure we have at least one of each character class
	var result strings.Builder

	// Add one character from each class
	result.WriteByte(upperChars[randomInt(len(upperChars))])
	result.WriteByte(lowerChars[randomInt(len(lowerChars))])
	result.WriteByte(digitChars[randomInt(len(digitChars))])
	result.WriteByte(specialChars[randomInt(len(specialChars))])

	// Add remaining characters from all classes
	allChars := upperChars + lowerChars + digitChars + specialChars
	for i := 4; i < minLength; i++ {
		result.WriteByte(allChars[randomInt(len(allChars))])
	}

	// Shuffle the result
	resultStr := result.String()
	shuffled := []byte(resultStr)
	for i := range shuffled {
		j := randomInt(len(shuffled))
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return string(shuffled), nil
}

// generateJWTLikeToken generates a JWT-like token with proper structure
func generateJWTLikeToken() (string, error) {
	// Generate random header
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))

	// Generate random payload
	payload := make([]byte, 32)
	if _, err := rand.Read(payload); err != nil {
		return "", err
	}
	payloadB64 := base64.RawURLEncoding.EncodeToString(payload)

	// Generate random signature
	signature := make([]byte, 32)
	if _, err := rand.Read(signature); err != nil {
		return "", err
	}
	signatureB64 := base64.RawURLEncoding.EncodeToString(signature)

	return fmt.Sprintf("%s.%s.%s", header, payloadB64, signatureB64), nil
}

// generateID generates a unique ID for secrets
func generateID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("secret-%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("sec-%x", b)
}

// randomInt generates a random integer in [0, n)
func randomInt(n int) int {
	if n <= 0 {
		return 0
	}
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return 0
	}
	return int(b[0]) % n
}

// GetSecretFromEnv retrieves a secret from environment variables
func GetSecretFromEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("environment variable not set: %s", key)
	}
	return value, nil
}

// GetSecretFromEnvWithDefault retrieves a secret from environment variables with a default value
func GetSecretFromEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

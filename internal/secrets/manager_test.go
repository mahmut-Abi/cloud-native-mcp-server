package secrets

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewInMemoryManager(t *testing.T) {
	manager := NewInMemoryManager()
	assert.NotNil(t, manager)
}

func TestStoreSecret(t *testing.T) {
	manager := NewInMemoryManager()

	secret := &Secret{
		Type:  SecretTypeAPIKey,
		Name:  "test-key",
		Value: "test-value",
	}

	err := manager.Store(secret)
	assert.NoError(t, err)
	assert.NotEmpty(t, secret.ID)
	assert.False(t, secret.CreatedAt.IsZero())
}

func TestStoreNilSecret(t *testing.T) {
	manager := NewInMemoryManager()

	err := manager.Store(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

func TestRetrieveSecret(t *testing.T) {
	manager := NewInMemoryManager()

	secret := &Secret{
		Type:  SecretTypeAPIKey,
		Name:  "test-key",
		Value: "test-value",
	}

	err := manager.Store(secret)
	assert.NoError(t, err)

	retrieved, err := manager.Retrieve(secret.ID)
	assert.NoError(t, err)
	assert.Equal(t, secret.ID, retrieved.ID)
	assert.Equal(t, secret.Type, retrieved.Type)
	assert.Equal(t, secret.Name, retrieved.Name)
	assert.Equal(t, secret.Value, retrieved.Value)
}

func TestRetrieveNonExistentSecret(t *testing.T) {
	manager := NewInMemoryManager()

	_, err := manager.Retrieve("non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestRetrieveExpiredSecret(t *testing.T) {
	manager := NewInMemoryManager()

	expiresAt := time.Now().Add(-1 * time.Hour)
	secret := &Secret{
		Type:      SecretTypeAPIKey,
		Name:      "test-key",
		Value:     "test-value",
		ExpiresAt: &expiresAt,
	}

	err := manager.Store(secret)
	assert.NoError(t, err)

	_, err = manager.Retrieve(secret.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expired")
}

func TestListSecrets(t *testing.T) {
	manager := NewInMemoryManager()

	// Add API key secrets
	err := manager.Store(&Secret{
		Type:  SecretTypeAPIKey,
		Name:  "api-key-1",
		Value: "value-1",
	})
	assert.NoError(t, err)

	err = manager.Store(&Secret{
		Type:  SecretTypeAPIKey,
		Name:  "api-key-2",
		Value: "value-2",
	})
	assert.NoError(t, err)

	// Add bearer token secrets
	err = manager.Store(&Secret{
		Type:  SecretTypeBearerToken,
		Name:  "token-1",
		Value: "token-value-1",
	})
	assert.NoError(t, err)

	// List API keys
	apiKeys, err := manager.List(SecretTypeAPIKey)
	assert.NoError(t, err)
	assert.Len(t, apiKeys, 2)

	// List bearer tokens
	tokens, err := manager.List(SecretTypeBearerToken)
	assert.NoError(t, err)
	assert.Len(t, tokens, 1)
}
func TestListSecretsExcludesExpired(t *testing.T) {
	manager := NewInMemoryManager()

	// Add non-expired secret
	err := manager.Store(&Secret{
		Type:  SecretTypeAPIKey,
		Name:  "valid-key",
		Value: "value",
	})
	assert.NoError(t, err)

	// Add expired secret
	expiresAt := time.Now().Add(-1 * time.Hour)
	err = manager.Store(&Secret{
		Type:      SecretTypeAPIKey,
		Name:      "expired-key",
		Value:     "value",
		ExpiresAt: &expiresAt,
	})
	assert.NoError(t, err)

	secrets, err := manager.List(SecretTypeAPIKey)
	assert.NoError(t, err)
	assert.Len(t, secrets, 1)
	assert.Equal(t, "valid-key", secrets[0].Name)
}
func TestDeleteSecret(t *testing.T) {
	manager := NewInMemoryManager()

	secret := &Secret{
		Type:  SecretTypeAPIKey,
		Name:  "test-key",
		Value: "test-value",
	}

	err := manager.Store(secret)
	assert.NoError(t, err)

	err = manager.Delete(secret.ID)
	assert.NoError(t, err)

	_, err = manager.Retrieve(secret.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestDeleteNonExistentSecret(t *testing.T) {
	manager := NewInMemoryManager()

	err := manager.Delete("non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestRotateAPIKey(t *testing.T) {
	manager := NewInMemoryManager()

	secret := &Secret{
		Type:  SecretTypeAPIKey,
		Name:  "test-key",
		Value: "original-value",
	}

	err := manager.Store(secret)
	assert.NoError(t, err)

	// Add a small delay to ensure timestamp difference
	time.Sleep(1 * time.Millisecond)

	rotated, err := manager.Rotate(secret.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, "original-value", rotated.Value)
	assert.True(t, rotated.UpdatedAt.After(secret.UpdatedAt) || rotated.UpdatedAt.Equal(secret.UpdatedAt))
}

func TestRotateBearerToken(t *testing.T) {
	manager := NewInMemoryManager()

	secret := &Secret{
		Type:  SecretTypeBearerToken,
		Name:  "test-token",
		Value: "original-token",
	}

	err := manager.Store(secret)
	assert.NoError(t, err)

	// Add a small delay to ensure timestamp difference
	time.Sleep(1 * time.Millisecond)

	rotated, err := manager.Rotate(secret.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, "original-token", rotated.Value)
	assert.True(t, rotated.UpdatedAt.After(secret.UpdatedAt) || rotated.UpdatedAt.Equal(secret.UpdatedAt))
}

func TestRotateUnsupportedType(t *testing.T) {
	manager := NewInMemoryManager()

	secret := &Secret{
		Type:  SecretTypeBasicAuth,
		Name:  "test-basic",
		Value: "value",
	}

	err := manager.Store(secret)
	assert.NoError(t, err)

	_, err = manager.Rotate(secret.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rotation not supported")
}

func TestGenerateAPIKey(t *testing.T) {
	manager := NewInMemoryManager()

	secret, err := manager.GenerateAPIKey("generated-key")
	assert.NoError(t, err)
	assert.NotEmpty(t, secret.ID)
	assert.Equal(t, SecretTypeAPIKey, secret.Type)
	assert.Equal(t, "generated-key", secret.Name)
	assert.NotEmpty(t, secret.Value)
	assert.True(t, len(secret.Value) >= 16)

	// Verify complexity
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range secret.Value {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

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

	assert.True(t, classCount >= 3, "API key should have at least 3 character classes")
}

func TestGenerateBearerToken(t *testing.T) {
	manager := NewInMemoryManager()

	secret, err := manager.GenerateBearerToken("generated-token")
	assert.NoError(t, err)
	assert.NotEmpty(t, secret.ID)
	assert.Equal(t, SecretTypeBearerToken, secret.Type)
	assert.Equal(t, "generated-token", secret.Name)
	assert.NotEmpty(t, secret.Value)

	// Verify JWT structure (header.payload.signature)
	parts := strings.Split(secret.Value, ".")
	assert.Len(t, parts, 3, "Bearer token should have JWT structure with 3 parts")
}

func TestGetSecretFromEnv(t *testing.T) {
	// Set environment variable
	t.Setenv("TEST_SECRET", "secret-value")

	value, err := GetSecretFromEnv("TEST_SECRET")
	assert.NoError(t, err)
	assert.Equal(t, "secret-value", value)
}

func TestGetSecretFromEnvNotSet(t *testing.T) {
	_, err := GetSecretFromEnv("NON_EXISTENT_SECRET")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not set")
}

func TestGetSecretFromEnvWithDefault(t *testing.T) {
	// Test with environment variable set
	t.Setenv("TEST_SECRET", "secret-value")

	value := GetSecretFromEnvWithDefault("TEST_SECRET", "default-value")
	assert.Equal(t, "secret-value", value)

	// Test with environment variable not set
	value = GetSecretFromEnvWithDefault("NON_EXISTENT_SECRET", "default-value")
	assert.Equal(t, "default-value", value)
}

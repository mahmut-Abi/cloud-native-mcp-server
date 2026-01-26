package middleware

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMultipleAPIKeys tests registration and authentication with multiple API keys
func TestMultipleAPIKeys(t *testing.T) {
	provider := NewSimpleAuthProvider()

	// Register multiple users with different API keys
	user1 := &AuthContext{
		UserID:   "user1",
		Username: "alice",
		Roles:    []string{"viewer"},
		Token:    "Bearer key-alice-123",
	}
	user2 := &AuthContext{
		UserID:   "user2",
		Username: "bob",
		Roles:    []string{"editor"},
		Token:    "Bearer key-bob-456",
	}
	user3 := &AuthContext{
		UserID:   "user3",
		Username: "charlie",
		Roles:    []string{"admin"},
		Token:    "Bearer key-charlie-789",
	}

	// Register all keys
	provider.RegisterAPIKey("key-alice-123", user1)
	provider.RegisterAPIKey("key-bob-456", user2)
	provider.RegisterAPIKey("key-charlie-789", user3)

	// Test authentication with different keys
	tests := []struct {
		name         string
		token        string
		expectedUser *AuthContext
		expectError  bool
	}{
		{
			name:         "alice_key",
			token:        "Bearer key-alice-123",
			expectedUser: user1,
			expectError:  false,
		},
		{
			name:         "bob_key",
			token:        "Bearer key-bob-456",
			expectedUser: user2,
			expectError:  false,
		},
		{
			name:         "charlie_key",
			token:        "Bearer key-charlie-789",
			expectedUser: user3,
			expectError:  false,
		},
		{
			name:         "unknown_key",
			token:        "Bearer unknown-key",
			expectedUser: nil,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := provider.Authenticate(tt.token)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.UserID, user.UserID)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Roles, user.Roles)
			}
		})
	}
}

// TestAPIKeyOverwrite tests overwriting an existing API key
func TestAPIKeyOverwrite(t *testing.T) {
	provider := NewSimpleAuthProvider()

	// Register initial user
	user1 := &AuthContext{
		UserID:   "user1",
		Username: "alice",
		Roles:    []string{"viewer"},
	}
	provider.RegisterAPIKey("shared-key", user1)

	// Verify initial registration
	user, err := provider.Authenticate("Bearer shared-key")
	assert.NoError(t, err)
	assert.Equal(t, "alice", user.Username)

	// Overwrite with new user
	user2 := &AuthContext{
		UserID:   "user2",
		Username: "bob",
		Roles:    []string{"admin"},
	}
	provider.RegisterAPIKey("shared-key", user2)

	// Verify overwrite worked
	user, err = provider.Authenticate("Bearer shared-key")
	assert.NoError(t, err)
	assert.Equal(t, "bob", user.Username)
	assert.Equal(t, []string{"admin"}, user.Roles)
}

// TestConcurrentAPIKeyRegistration tests concurrent registration of API keys
func TestConcurrentAPIKeyRegistration(t *testing.T) {
	provider := NewSimpleAuthProvider()

	// Number of concurrent goroutines
	numGoroutines := 100
	done := make(chan bool, numGoroutines)

	// Launch concurrent registrations
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			user := &AuthContext{
				UserID:   fmt.Sprintf("user%d", id),
				Username: fmt.Sprintf("user%d", id),
				Roles:    []string{"viewer"},
			}
			key := fmt.Sprintf("key-%d", id)
			provider.RegisterAPIKey(key, user)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify all keys were registered
	for i := 0; i < numGoroutines; i++ {
		key := fmt.Sprintf("key-%d", i)
		token := fmt.Sprintf("Bearer key-%d", i)
		user, err := provider.Authenticate(token)
		assert.NoError(t, err, "Failed to authenticate key %s", key)
		assert.Equal(t, fmt.Sprintf("user%d", i), user.UserID)
	}
}

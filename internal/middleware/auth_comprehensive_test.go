package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSimpleAuthProvider_RegisterAPIKey tests API key registration
func TestSimpleAuthProvider_RegisterAPIKey(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		user        *AuthContext
		expectCount int
	}{
		{
			name: "valid_api_key",
			key:  "test-key-123",
			user: &AuthContext{
				UserID:   "user123",
				Username: "testuser",
				Roles:    []string{"viewer"},
				Token:    "Bearer test-key-123",
			},
			expectCount: 1,
		},
		{
			name:        "empty_key",
			key:         "",
			user:        &AuthContext{UserID: "user123"},
			expectCount: 0,
		},
		{
			name:        "nil_user",
			key:         "test-key",
			user:        nil,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewSimpleAuthProvider()
			provider.RegisterAPIKey(tt.key, tt.user)

			// Check internal map size
			provider.mu.RLock()
			actualCount := len(provider.apiKeys)
			provider.mu.RUnlock()

			assert.Equal(t, tt.expectCount, actualCount)
		})
	}
}

// TestSimpleAuthProvider_Authenticate tests token authentication
func TestSimpleAuthProvider_Authenticate(t *testing.T) {
	provider := NewSimpleAuthProvider()
	testUser := &AuthContext{
		UserID:   "user123",
		Username: "testuser",
		Roles:    []string{"admin"},
		Token:    "Bearer valid-key",
	}
	provider.RegisterAPIKey("valid-key", testUser)

	tests := []struct {
		name        string
		token       string
		expectUser  *AuthContext
		expectError bool
	}{
		{
			name:        "valid_bearer_token",
			token:       "Bearer valid-key",
			expectUser:  testUser,
			expectError: false,
		},
		{
			name:        "invalid_api_key",
			token:       "Bearer invalid-key",
			expectUser:  nil,
			expectError: true,
		},
		{
			name:        "empty_token",
			token:       "",
			expectUser:  nil,
			expectError: true,
		},
		{
			name:        "malformed_token",
			token:       "invalid-format",
			expectUser:  nil,
			expectError: true,
		},
		{
			name:        "wrong_prefix",
			token:       "Basic valid-key",
			expectUser:  nil,
			expectError: true,
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
				assert.Equal(t, tt.expectUser, user)
			}
		})
	}
}

// TestSimpleAuthProvider_ValidatePermission tests permission validation
func TestSimpleAuthProvider_ValidatePermission(t *testing.T) {
	provider := NewSimpleAuthProvider()

	tests := []struct {
		name     string
		ctx      *AuthContext
		resource string
		action   string
		expected bool
	}{
		{
			name:     "nil_context",
			ctx:      nil,
			resource: "pods",
			action:   "get",
			expected: false,
		},
		{
			name: "admin_role_all_permissions",
			ctx: &AuthContext{
				UserID:   "admin1",
				Username: "admin",
				Roles:    []string{"admin"},
			},
			resource: "pods",
			action:   "delete",
			expected: true,
		},
		{
			name: "viewer_role_get_permission",
			ctx: &AuthContext{
				UserID:   "viewer1",
				Username: "viewer",
				Roles:    []string{"viewer"},
			},
			resource: "pods",
			action:   "get",
			expected: true,
		},
		{
			name: "viewer_role_no_update_permission",
			ctx: &AuthContext{
				UserID:   "viewer1",
				Username: "viewer",
				Roles:    []string{"viewer"},
			},
			resource: "pods",
			action:   "update",
			expected: false,
		},
		{
			name: "editor_role_get_and_update_permission",
			ctx: &AuthContext{
				UserID:   "editor1",
				Username: "editor",
				Roles:    []string{"editor"},
			},
			resource: "deployments",
			action:   "update",
			expected: true,
		},
		{
			name: "editor_role_no_delete_permission",
			ctx: &AuthContext{
				UserID:   "editor1",
				Username: "editor",
				Roles:    []string{"editor"},
			},
			resource: "deployments",
			action:   "delete",
			expected: false,
		},
		{
			name: "multiple_roles_with_admin",
			ctx: &AuthContext{
				UserID:   "multiuser",
				Username: "multiuser",
				Roles:    []string{"viewer", "admin", "editor"},
			},
			resource: "secrets",
			action:   "delete",
			expected: true,
		},
		{
			name: "unknown_role",
			ctx: &AuthContext{
				UserID:   "unknown1",
				Username: "unknown",
				Roles:    []string{"unknown"},
			},
			resource: "pods",
			action:   "get",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.ValidatePermission(tt.ctx, tt.resource, tt.action)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestHasPermission tests the hasPermission helper function
func TestHasPermission(t *testing.T) {
	tests := []struct {
		name     string
		roles    []string
		resource string
		action   string
		expected bool
	}{
		{
			name:     "viewer_can_get",
			roles:    []string{"viewer"},
			resource: "pods",
			action:   "get",
			expected: true,
		},
		{
			name:     "editor_can_get",
			roles:    []string{"editor"},
			resource: "deployments",
			action:   "get",
			expected: true,
		},
		{
			name:     "editor_can_update",
			roles:    []string{"editor"},
			resource: "services",
			action:   "update",
			expected: true,
		},
		{
			name:     "viewer_cannot_update",
			roles:    []string{"viewer"},
			resource: "pods",
			action:   "update",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasPermission(tt.roles, tt.resource, tt.action)
			assert.Equal(t, tt.expected, result)
		})
	}
}

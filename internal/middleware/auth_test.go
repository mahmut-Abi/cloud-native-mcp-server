package middleware

import (
	"testing"
)

func TestSimpleAuthProvider(t *testing.T) {
	provider := NewSimpleAuthProvider()

	// Register test user
	user := &AuthContext{
		UserID:   "user1",
		Username: "testuser",
		Roles:    []string{"editor"},
		Token:    "test-token-123",
	}
	provider.RegisterAPIKey("test-key-123", user)

	// Test successful authentication
	authCtx, err := provider.Authenticate("Bearer test-key-123")
	if err != nil {
		t.Fatalf("Authenticate() error = %v", err)
	}

	if authCtx.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %s", authCtx.Username)
	}
}

func TestAuthenticationFailure(t *testing.T) {
	provider := NewSimpleAuthProvider()

	tests := []struct {
		name  string
		token string
	}{
		{"empty token", ""},
		{"invalid format", "InvalidToken"},
		{"invalid key", "Bearer invalid-key"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := provider.Authenticate(tt.token)
			if err == nil {
				t.Error("Expected authentication to fail")
			}
		})
	}
}

func TestPermissionValidation(t *testing.T) {
	provider := NewSimpleAuthProvider()

	// Test admin role
	adminCtx := &AuthContext{
		Roles: []string{"admin"},
	}
	if !provider.ValidatePermission(adminCtx, "any", "any") {
		t.Error("Admin should have all permissions")
	}

	// Test viewer role
	viewerCtx := &AuthContext{
		Roles: []string{"viewer"},
	}
	if !provider.ValidatePermission(viewerCtx, "resource", "get") {
		t.Error("Viewer should be able to get resources")
	}

	// Test insufficient permissions
	if provider.ValidatePermission(viewerCtx, "resource", "delete") {
		t.Error("Viewer should not be able to delete")
	}
}

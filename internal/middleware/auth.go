package middleware

import (
	"fmt"
	"strings"
	"sync"
)

// AuthContext represents authentication context
type AuthContext struct {
	UserID   string
	Username string
	Roles    []string
	Token    string
}

// AuthenticationProvider defines authentication interface
type AuthenticationProvider interface {
	Authenticate(token string) (*AuthContext, error)
	ValidatePermission(ctx *AuthContext, resource, action string) bool
}

// SimpleAuthProvider implements basic API key authentication
type SimpleAuthProvider struct {
	mu      sync.RWMutex
	apiKeys map[string]*AuthContext
}

// NewSimpleAuthProvider creates a new simple auth provider
func NewSimpleAuthProvider() *SimpleAuthProvider {
	return &SimpleAuthProvider{
		apiKeys: make(map[string]*AuthContext),
	}
}

// RegisterAPIKey registers an API key
func (p *SimpleAuthProvider) RegisterAPIKey(key string, user *AuthContext) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if key == "" || user == nil {
		return // Prevent empty keys or nil users
	}
	p.apiKeys[key] = user
}

// Authenticate authenticates using API key
func (p *SimpleAuthProvider) Authenticate(token string) (*AuthContext, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if token == "" {
		return nil, fmt.Errorf("token cannot be empty")
	}

	// Extract Bearer token
	parts := strings.Split(token, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("invalid token format")
	}

	authCtx, exists := p.apiKeys[parts[1]]
	if !exists {
		return nil, fmt.Errorf("invalid API key")
	}

	return authCtx, nil
}

// ValidatePermission validates user permission
func (p *SimpleAuthProvider) ValidatePermission(ctx *AuthContext, resource, action string) bool {
	if ctx == nil {
		return false
	}

	// Admin role has all permissions
	for _, role := range ctx.Roles {
		if role == "admin" {
			return true
		}
	}

	// Check specific role-based permissions
	return hasPermission(ctx.Roles, resource, action)
}

// hasPermission checks if role has permission for resource action
func hasPermission(roles []string, resource, action string) bool {
	for _, role := range roles {
		if role == "viewer" && action == "get" {
			return true
		}
		if role == "editor" && (action == "get" || action == "update") {
			return true
		}
	}
	return false
}

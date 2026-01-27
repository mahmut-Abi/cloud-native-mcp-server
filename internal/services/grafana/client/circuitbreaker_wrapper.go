package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/circuitbreaker"
	"github.com/sirupsen/logrus"
)

// CircuitBreakerConfig holds configuration for circuit breaker
type CircuitBreakerConfig struct {
	Enabled         bool          // Enable circuit breaker
	MaxFailures     int           // Maximum failures before opening
	Timeout         time.Duration // Time before attempting to close circuit
	SuccessesNeeded int           // Successes needed to close circuit in half-open state
}

// DefaultCircuitBreakerConfig returns default circuit breaker configuration
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		Enabled:         true,
		MaxFailures:     5,
		Timeout:         30 * time.Second,
		SuccessesNeeded: 2,
	}
}

// CircuitBreakerManager manages circuit breakers for different operations
type CircuitBreakerManager struct {
	mu            sync.RWMutex
	breakers      map[string]*circuitbreaker.CircuitBreaker
	config        CircuitBreakerConfig
	onStateChange func(operation string, oldState, newState circuitbreaker.State)
}

// NewCircuitBreakerManager creates a new circuit breaker manager
func NewCircuitBreakerManager(config CircuitBreakerConfig) *CircuitBreakerManager {
	if !config.Enabled {
		return &CircuitBreakerManager{
			breakers: make(map[string]*circuitbreaker.CircuitBreaker),
			config:   config,
		}
	}

	manager := &CircuitBreakerManager{
		breakers: make(map[string]*circuitbreaker.CircuitBreaker),
		config:   config,
	}

	manager.onStateChange = func(operation string, oldState, newState circuitbreaker.State) {
		logrus.WithFields(logrus.Fields{
			"operation": operation,
			"old_state": oldState,
			"new_state": newState,
			"service":   "grafana",
		}).Warn("Circuit breaker state changed")
	}

	return manager
}

// getOrCreateBreaker gets or creates a circuit breaker for an operation
func (m *CircuitBreakerManager) getOrCreateBreaker(operation string) *circuitbreaker.CircuitBreaker {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.config.Enabled {
		return nil
	}

	if breaker, exists := m.breakers[operation]; exists {
		return breaker
	}

	breaker := circuitbreaker.NewCircuitBreaker(m.config.MaxFailures, m.config.Timeout)
	if m.onStateChange != nil {
		breaker.SetOnStateChange(func(oldState, newState circuitbreaker.State) {
			m.onStateChange(operation, oldState, newState)
		})
	}

	m.breakers[operation] = breaker
	return breaker
}

// Execute executes a function with circuit breaker protection
func (m *CircuitBreakerManager) Execute(ctx context.Context, operation string, f func(context.Context) error) error {
	if !m.config.Enabled {
		return f(ctx)
	}

	breaker := m.getOrCreateBreaker(operation)
	if breaker == nil {
		return f(ctx)
	}

	return breaker.Do(ctx, f)
}

// GetState returns the current state of a circuit breaker
func (m *CircuitBreakerManager) GetState(operation string) circuitbreaker.State {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.config.Enabled {
		return circuitbreaker.StateClosed
	}

	breaker, exists := m.breakers[operation]
	if !exists {
		return circuitbreaker.StateClosed
	}

	return breaker.GetState()
}

// Reset resets a specific circuit breaker
func (m *CircuitBreakerManager) Reset(operation string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.config.Enabled {
		return
	}

	breaker, exists := m.breakers[operation]
	if exists {
		breaker.Reset()
		logrus.WithFields(logrus.Fields{
			"operation": operation,
			"service":   "grafana",
		}).Info("Circuit breaker reset")
	}
}

// ResetAll resets all circuit breakers
func (m *CircuitBreakerManager) ResetAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.config.Enabled {
		return
	}

	for operation := range m.breakers {
		m.breakers[operation].Reset()
	}

	logrus.WithField("service", "grafana").Info("All circuit breakers reset")
}

// GetStats returns statistics for all circuit breakers
func (m *CircuitBreakerManager) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make(map[string]interface{})
	stats["enabled"] = m.config.Enabled

	if !m.config.Enabled {
		return stats
	}

	breakerStats := make(map[string]string)
	for operation, breaker := range m.breakers {
		breakerStats[operation] = string(breaker.GetState())
	}
	stats["breakers"] = breakerStats

	return stats
}

// Operation names for circuit breaker
const (
	OpGetDashboards    = "get_dashboards"
	OpGetDashboard     = "get_dashboard"
	OpGetDataSources   = "get_datasources"
	OpGetFolders       = "get_folders"
	OpGetAlertRules    = "get_alert_rules"
	OpCreateDashboard  = "create_dashboard"
	OpUpdateDashboard  = "update_dashboard"
	OpDeleteDashboard  = "delete_dashboard"
	OpCreateDatasource = "create_datasource"
	OpUpdateDatasource = "update_datasource"
	OpDeleteDatasource = "delete_datasource"
	OpTestDatasource   = "test_datasource"
	OpCreateAlertRule  = "create_alert_rule"
	OpUpdateAlertRule  = "update_alert_rule"
	OpDeleteAlertRule  = "delete_alert_rule"
	OpGetAnnotations   = "get_annotations"
	OpCreateAnnotation = "create_annotation"
	OpUpdateAnnotation = "update_annotation"
	OpDeleteAnnotation = "delete_annotation"
	OpRenderPanel      = "render_panel"
	OpSearchDashboards = "search_dashboards"
	OpGetUsers         = "get_users"
	OpGetTeams         = "get_teams"
	OpGetRoles         = "get_roles"
)

// wrapWithCircuitBreaker wraps a function with circuit breaker protection
func (c *Client) wrapWithCircuitBreaker(ctx context.Context, operation string, f func(context.Context) error) error {
	if c.circuitBreakerManager == nil {
		return f(ctx)
	}

	return c.circuitBreakerManager.Execute(ctx, operation, f)
}

// CircuitBreakerError creates a circuit breaker error
func CircuitBreakerError(operation string) error {
	return fmt.Errorf("circuit breaker is open for operation: %s", operation)
}

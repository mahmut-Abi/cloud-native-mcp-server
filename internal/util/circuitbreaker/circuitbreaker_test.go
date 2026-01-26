package circuitbreaker

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCircuitBreakerClosed(t *testing.T) {
	cb := NewCircuitBreaker(3, 100*time.Millisecond)

	if cb.GetState() != StateClosed {
		t.Errorf("Initial state should be closed, got %s", cb.GetState())
	}
}

func TestCircuitBreakerTransitionToClosed(t *testing.T) {
	cb := NewCircuitBreaker(2, 100*time.Millisecond)

	err := cb.Do(context.Background(), func(ctx context.Context) error {
		return errors.New("fail1")
	})
	if err == nil {
		t.Error("Expected error on first failure")
	}

	err = cb.Do(context.Background(), func(ctx context.Context) error {
		return errors.New("fail2")
	})
	if err == nil {
		t.Error("Expected error on second failure")
	}

	if cb.GetState() != StateOpen {
		t.Errorf("State should be open after max failures, got %s", cb.GetState())
	}
}

func TestCircuitBreakerOpen(t *testing.T) {
	cb := NewCircuitBreaker(1, 100*time.Millisecond)

	err := cb.Do(context.Background(), func(ctx context.Context) error {
		return errors.New("fail")
	})
	if err == nil {
		t.Error("Expected error on failure")
	}

	if cb.GetState() != StateOpen {
		t.Errorf("State should be open, got %s", cb.GetState())
	}

	err = cb.Do(context.Background(), func(ctx context.Context) error {
		return nil
	})
	if err == nil || err.Error() != "circuit breaker is open, cannot execute operation" {
		t.Errorf("Expected circuit breaker open error, got %v", err)
	}
}

func TestCircuitBreakerReset(t *testing.T) {
	cb := NewCircuitBreaker(1, 100*time.Millisecond)

	err := cb.Do(context.Background(), func(ctx context.Context) error {
		return errors.New("fail")
	})
	if err == nil {
		t.Error("Expected error on failure")
	}

	if cb.GetState() != StateOpen {
		t.Errorf("State should be open, got %s", cb.GetState())
	}

	cb.Reset()

	if cb.GetState() != StateClosed {
		t.Errorf("State should be closed after reset, got %s", cb.GetState())
	}
}

func TestCircuitBreakerHalfOpen(t *testing.T) {
	cb := NewCircuitBreaker(1, 50*time.Millisecond)

	err := cb.Do(context.Background(), func(ctx context.Context) error {
		return errors.New("fail")
	})
	if err == nil {
		t.Error("Expected error on failure")
	}

	if cb.GetState() != StateOpen {
		t.Errorf("State should be open, got %s", cb.GetState())
	}

	time.Sleep(100 * time.Millisecond)

	err = cb.Do(context.Background(), func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected success in half-open state, got error: %v", err)
	}

	if cb.GetState() != StateClosed {
		t.Errorf("State should be closed after successful half-open attempt, got %s", cb.GetState())
	}
}

func TestCircuitBreakerStateChangeCallback(t *testing.T) {
	cb := NewCircuitBreaker(1, 100*time.Millisecond)
	var stateChanges []State

	cb.SetOnStateChange(func(oldState, newState State) {
		stateChanges = append(stateChanges, newState)
	})

	_ = cb.Do(context.Background(), func(ctx context.Context) error {
		return errors.New("fail")
	})

	if len(stateChanges) != 1 || stateChanges[0] != StateOpen {
		t.Errorf("Expected state change to Open, got %v", stateChanges)
	}
}

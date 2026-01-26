package circuitbreaker

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type State string

const (
	StateClosed   State = "closed"
	StateOpen     State = "open"
	StateHalfOpen State = "half-open"
)

type CircuitBreaker struct {
	mu              sync.RWMutex
	state           State
	failures        int
	maxFailures     int
	successesNeeded int
	successesInHalf int
	lastFailureTime time.Time
	timeout         time.Duration
	onStateChange   func(oldState, newState State)
}

func NewCircuitBreaker(maxFailures int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:           StateClosed,
		maxFailures:     maxFailures,
		timeout:         timeout,
		successesNeeded: maxFailures,
	}
}

func (cb *CircuitBreaker) SetOnStateChange(callback func(oldState, newState State)) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.onStateChange = callback
}

func (cb *CircuitBreaker) Do(ctx context.Context, f func(context.Context) error) error {
	cb.mu.Lock()

	if cb.state == StateOpen {
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.transitionTo(StateHalfOpen)
		} else {
			cb.mu.Unlock()
			return fmt.Errorf("circuit breaker is open, cannot execute operation")
		}
	}

	if cb.state == StateOpen {
		cb.mu.Unlock()
		return fmt.Errorf("circuit breaker is open")
	}

	cb.mu.Unlock()

	err := f(ctx)

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.recordFailure()
		return err
	}

	cb.recordSuccess()
	return nil
}

func (cb *CircuitBreaker) recordFailure() {
	cb.failures++
	cb.lastFailureTime = time.Now()

	if cb.state == StateHalfOpen {
		cb.transitionTo(StateOpen)
	} else if cb.state == StateClosed && cb.failures >= cb.maxFailures {
		cb.transitionTo(StateOpen)
	}
}

func (cb *CircuitBreaker) recordSuccess() {
	switch cb.state {
	case StateHalfOpen:
		cb.successesInHalf++
		if cb.successesInHalf >= cb.successesNeeded {
			cb.transitionTo(StateClosed)
		}
	case StateClosed:
		cb.failures = 0
	}
}

func (cb *CircuitBreaker) transitionTo(newState State) {
	oldState := cb.state
	cb.state = newState

	switch newState {
	case StateClosed:
		cb.failures = 0
		cb.successesInHalf = 0
	case StateOpen:
		cb.lastFailureTime = time.Now()
	case StateHalfOpen:
		cb.failures = 0
		cb.successesInHalf = 0
	}

	if cb.onStateChange != nil {
		cb.onStateChange(oldState, newState)
	}
}

func (cb *CircuitBreaker) GetState() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.transitionTo(StateClosed)
}

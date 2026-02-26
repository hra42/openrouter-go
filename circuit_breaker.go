package openrouter

import (
	"sync"
	"time"
)

// CircuitState represents the current state of the circuit breaker.
type CircuitState int

const (
	// CircuitClosed allows requests to pass through normally.
	CircuitClosed CircuitState = iota
	// CircuitOpen blocks requests due to repeated failures.
	CircuitOpen
	// CircuitHalfOpen allows a limited number of test requests.
	CircuitHalfOpen
)

// String returns the string representation of the circuit state.
func (s CircuitState) String() string {
	switch s {
	case CircuitClosed:
		return "closed"
	case CircuitOpen:
		return "open"
	case CircuitHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// CircuitBreakerConfig configures the circuit breaker behavior.
type CircuitBreakerConfig struct {
	// FailureThreshold is the number of consecutive failures before opening the circuit. Default: 5.
	FailureThreshold int
	// ResetTimeout is how long to wait before transitioning from open to half-open. Default: 30s.
	ResetTimeout time.Duration
	// HalfOpenMaxAttempts is the number of test requests allowed in half-open state. Default: 1.
	HalfOpenMaxAttempts int
}

// DefaultCircuitBreakerConfig returns the default circuit breaker configuration.
func DefaultCircuitBreakerConfig() *CircuitBreakerConfig {
	return &CircuitBreakerConfig{
		FailureThreshold:    5,
		ResetTimeout:        30 * time.Second,
		HalfOpenMaxAttempts: 1,
	}
}

// CircuitBreaker implements the circuit breaker pattern for stream reconnection.
type CircuitBreaker struct {
	mu               sync.Mutex
	state            CircuitState
	failures         int
	lastFailure      time.Time
	halfOpenAttempts int
	config           *CircuitBreakerConfig
}

// NewCircuitBreaker creates a new CircuitBreaker with the given configuration.
// If config is nil, DefaultCircuitBreakerConfig is used.
func NewCircuitBreaker(config *CircuitBreakerConfig) *CircuitBreaker {
	if config == nil {
		config = DefaultCircuitBreakerConfig()
	}
	return &CircuitBreaker{
		state:  CircuitClosed,
		config: config,
	}
}

// AllowRequest returns true if the circuit breaker allows a request to proceed.
func (cb *CircuitBreaker) AllowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case CircuitClosed:
		return true
	case CircuitOpen:
		if time.Since(cb.lastFailure) >= cb.config.ResetTimeout {
			cb.state = CircuitHalfOpen
			cb.halfOpenAttempts = 1
			return true
		}
		return false
	case CircuitHalfOpen:
		if cb.halfOpenAttempts < cb.config.HalfOpenMaxAttempts {
			cb.halfOpenAttempts++
			return true
		}
		return false
	}
	return false
}

// RecordSuccess records a successful request and resets the circuit breaker.
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.state = CircuitClosed
	cb.failures = 0
	cb.halfOpenAttempts = 0
}

// RecordFailure records a failed request and may open the circuit.
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailure = time.Now()

	if cb.state == CircuitHalfOpen {
		cb.state = CircuitOpen
		return
	}

	if cb.failures >= cb.config.FailureThreshold {
		cb.state = CircuitOpen
	}
}

// State returns the current state of the circuit breaker.
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// Reset resets the circuit breaker to its initial closed state.
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.state = CircuitClosed
	cb.failures = 0
	cb.halfOpenAttempts = 0
	cb.lastFailure = time.Time{}
}
